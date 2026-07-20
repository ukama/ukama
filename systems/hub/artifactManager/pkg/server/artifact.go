/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package server

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"path"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/minio/minio-go/v7"
	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/errors"
	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/hub/artifactmanager/pkg"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	uuid "github.com/ukama/ukama/systems/common/uuid"
	pb "github.com/ukama/ukama/systems/hub/artifactmanager/pb/gen"
	dpb "github.com/ukama/ukama/systems/hub/distributor/pb/gen"
)

const UrlPath = "/v1/hub"
const ChunksPath = "/v1/distributor/"

type ArtifcatServer struct {
	pb.ArtifactServiceServer
	//distributorClient      cnucl.OrgClient
	msgbus                mb.MsgBusServiceClient
	baseRoutingKey        msgbus.RoutingKeyBuilder
	pushGateway           string
	OrgId                 uuid.UUID
	OrgName               string
	storage               pkg.Storage
	storageRequestTimeout time.Duration
	chunker               chunkServer
}

type chunkServer interface {
	CreateChunk(in *dpb.CreateChunkRequest) (*dpb.CreateChunkResponse, error)
}

func NewArtifactServer(orgId uuid.UUID, orgName string, storage pkg.Storage, chunk chunkServer, storageTimeout time.Duration,
	msgBus mb.MsgBusServiceClient, pushGateway string) *ArtifcatServer {

	rotuingKey := msgbus.NewRoutingKeyBuilder().SetCloudSource().SetGlobalScope().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName)

	return &ArtifcatServer{
		OrgId:                 orgId,
		OrgName:               orgName,
		msgbus:                msgBus,
		baseRoutingKey:        rotuingKey,
		pushGateway:           pushGateway,
		chunker:               chunk,
		storage:               storage,
		storageRequestTimeout: storageTimeout,
	}
}

func (s *ArtifcatServer) parseVersion(version string) (*semver.Version, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid version format. Refer to https://semver.org/ for more information")
	}
	return v, err
}

func (s *ArtifcatServer) parseArtifactName(name string) (ver *semver.Version, ext string, err error) {
	if strings.HasSuffix(name, pkg.TarGzExtension) {
		name = strings.TrimSuffix(name, pkg.TarGzExtension)
		ext = pkg.TarGzExtension
	} else if strings.HasSuffix(name, pkg.ChunkIndexExtension) {
		name = strings.TrimSuffix(name, pkg.ChunkIndexExtension)
		ext = pkg.ChunkIndexExtension
	} else {
		return nil, "", fmt.Errorf("unsupported extension")
	}

	ver, err = semver.NewVersion(name)
	if err != nil {
		return nil, "", errors.Wrap(err, "failed to parse version")
	}

	return ver, ext, nil
}

func (s *ArtifcatServer) StoreArtifact(ctx context.Context, in *pb.StoreArtifactRequest) (*pb.StoreArtifactResponse, error) {
	log.Infof("Storing artifact: %s %s of type %d", in.Name, in.Version, in.Type)

	v, err := s.parseVersion(in.Version)
	if err != nil {
		return nil, err
	}

	aType := strings.ToLower(in.Type.String())
	newDigest := sha256Hex(in.Data)

	// Immutability: a version, once stored, cannot be overwritten with different content.
	exists, existingDigest, err := s.storage.StatFile(ctx, in.Name, aType, v, pkg.TarGzExtension)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to stat existing artifact: %v", err)
	}
	if exists {
		switch existingDigest {
		case "":
			// Legacy object with no recorded digest — cannot prove equality; refuse to overwrite.
			return nil, status.Errorf(codes.AlreadyExists,
				"artifact %s version %s already exists (digest unverifiable); versions are immutable", in.Name, v.String())
		case newDigest:
			log.Infof("Artifact %s version %s already present with identical content; idempotent no-op", in.Name, v.String())
			return &pb.StoreArtifactResponse{Name: in.Name, Type: in.Type}, nil
		default:
			return nil, status.Errorf(codes.AlreadyExists,
				"artifact %s version %s already exists with different content; versions are immutable", in.Name, v.String())
		}
	}

	log.Infof("Got file %s with size %d", in.Name, len(in.Data))
	if _, err := s.storage.PutFile(ctx, in.Name, aType, v, pkg.TarGzExtension,
		bytes.NewReader(in.Data), map[string]string{pkg.ContentDigestMetaKey: newDigest}); err != nil {
		log.Errorf("Error storing artifact: %s %s", in.Name, in.Version)
		return nil, err
	}

	go func() {
		if err := s.chunkAndIndex(context.Background(), in.Name, aType, v); err != nil {
			log.Errorf("chunk/index for %s %s failed: %v", in.Name, v.String(), err)
		}
	}()

	return &pb.StoreArtifactResponse{
		Name: in.Name,
		Type: in.Type,
	}, nil
}

func sha256Hex(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}

// chunkAndIndex asks the distributor to chunk an already-stored tar.gz, persists the
// returned index, and announces availability. Shared by StoreArtifact and the sweep.
func (s *ArtifcatServer) chunkAndIndex(ctx context.Context, name, aType string, v *semver.Version) error {
	cReq := &dpb.CreateChunkRequest{
		Name:    name,
		Type:    aType,
		Version: v.String(),
		Store:   "s3+" + s.storage.StoreBaseURL(aType) + "?lookup=path",
	}
	log.Infof("Sending chunking request %+v", cReq)

	resp, err := s.chunker.CreateChunk(cReq)
	if err != nil {
		return fmt.Errorf("create chunk for %s %s: %w", name, v.String(), err)
	}

	nctx, cancel := context.WithTimeout(ctx, s.storageRequestTimeout)
	defer cancel()
	if _, err := s.storage.PutFile(nctx, name, aType, v, pkg.ChunkIndexExtension, bytes.NewReader(resp.Index), nil); err != nil {
		return fmt.Errorf("store index for %s %s: %w", name, v.String(), err)
	}

	capp := &epb.EventArtifactUploaded{Name: name, Version: v.String()}
	route := s.baseRoutingKey.SetAction("uploaded").SetObject(aType).MustBuild()
	if err := s.msgbus.PublishRequest(route, capp); err != nil {
		log.Errorf("Failed to publish uploaded event %+v key %+v: %v", capp, route, err)
	}
	return nil
}

func (s *ArtifcatServer) GetArtifactLocation(ctx context.Context, in *pb.GetArtifactLocationRequest) (*pb.GetArtifactLocationResponse, error) {
	log.Infof("Getting apps storage endpoint")
	return &pb.GetArtifactLocationResponse{
		Url: s.storage.GetEndpoint(),
	}, nil
}

func (s *ArtifcatServer) GetArtifact(ctx context.Context, in *pb.GetArtifactRequest) (*pb.GetArtifactResponse, error) {
	log.Infof("Getting artifact: %s of type %s with filename %s", in.Name, in.Type, in.FileName)

	v, ext, err := s.parseArtifactName(in.FileName)
	if err != nil {
		return nil, status.Error(codes.NotFound, "Artifact file name is not valid")
	}

	rd, err := s.storage.GetFile(ctx, in.Name, strings.ToLower(in.Type.String()), v, ext)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := rd.Close(); cerr != nil {
			log.Errorf("Failed to close reader: %v", cerr)
		}
	}()

	data, err := io.ReadAll(rd)
	if err != nil {
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			return nil, status.Error(codes.NotFound, "Artifact not found")
		}
	}

	return &pb.GetArtifactResponse{
		FileName: fmt.Sprintf("%s-%s%s", in.Name, v.String(), ext),
		Name:     in.Name,
		Type:     in.Type,
		Version:  v.String(),
		Data:     data,
	}, nil
}

func (s *ArtifcatServer) GetArtifactVersionList(ctx context.Context, in *pb.GetArtifactVersionListRequest) (*pb.GetArtifactVersionListResponse, error) {
	log.Infof("Getting version list: %s of type %s", in.Name, in.Type.String())

	ls, err := s.storage.ListVersions(ctx, in.Name, strings.ToLower(in.Type.String()))
	if err != nil {
		return nil, err
	}

	if len(*ls) == 0 {
		return nil, status.Error(codes.NotFound, "Artifact name is not valid")
	}

	aType := strings.ToLower(in.Type.String())
	vers := []*pb.VersionInfo{}
	for _, v := range *ls {
		vers = append(vers, &pb.VersionInfo{
			Version: v.Version,
			Formats: buildFormats(aType, in.Name, v),
		})
	}

	return &pb.GetArtifactVersionListResponse{
		Name:     in.Name,
		Type:     in.Type,
		Versions: vers,
	}, nil

}

// buildFormats builds the tar.gz (+ chunk, when present) format entries for a version.
func buildFormats(aType, name string, info pkg.AritfactInfo) []*pb.FormatInfo {
	formats := []*pb.FormatInfo{
		{
			Url:       path.Join(UrlPath, aType, name, info.Version+pkg.TarGzExtension),
			CreatedAt: timestamppb.New(info.CreatedAt),
			Size:      info.SizeBytes,
			Type:      "tar.gz",
		},
	}
	if info.Chunked {
		formats = append(formats, &pb.FormatInfo{
			Url:       path.Join(UrlPath, aType, name, info.Version+pkg.ChunkIndexExtension),
			Type:      "chunk",
			CreatedAt: timestamppb.New(info.CreatedAt),
			Size:      info.SizeBytes,
			ExtraInfo: []*pb.ExtraInfoMap{
				{Key: "chunks", Value: fmt.Sprintf("%s/", ChunksPath)},
			},
		})
	}
	return formats
}

func (s *ArtifcatServer) ListArtifacts(ctx context.Context, in *pb.ListArtifactRequest) (*pb.ListArtifactResponse, error) {
	aType := strings.ToLower(in.Type.String())
	log.Infof("Getting list of %s artifacts (latest=%v)", aType, in.Latest)

	if in.Latest {
		ls, err := s.storage.ListLatestPerApp(ctx, aType)
		if err != nil {
			return nil, err
		}
		out := make([]*pb.LatestArtifact, 0, len(ls))
		for i := range ls {
			out = append(out, &pb.LatestArtifact{
				Name: ls[i].Name,
				Latest: &pb.VersionInfo{
					Version: ls[i].Latest.Version,
					Formats: buildFormats(aType, ls[i].Name, ls[i].Latest),
				},
			})
		}
		return &pb.ListArtifactResponse{LatestArtifacts: out}, nil
	}

	ls, err := s.storage.ListApps(ctx, aType)
	if err != nil {
		return nil, err
	}
	return &pb.ListArtifactResponse{Artifact: ls}, nil
}
