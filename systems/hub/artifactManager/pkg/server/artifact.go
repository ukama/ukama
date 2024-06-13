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
	uuid "github.com/ukama/ukama/systems/common/uuid"
	pb "github.com/ukama/ukama/systems/hub/artifactmanager/pb/gen"
	dpb "github.com/ukama/ukama/systems/hub/distributor/pb/gen"
)

const CappsPath = "/v1/capps"
const ChunksPath = "/v1/chunks"

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
	IsGlobal              bool
}

type chunkServer interface {
	CreateChunk(in *dpb.CreateChunkRequest) (*dpb.CreateChunkResponse, error)
}

func NewArtifactServer(orgId uuid.UUID, orgName string, storage pkg.Storage, chunk chunkServer, storageTimeout time.Duration,
	msgBus mb.MsgBusServiceClient, pushGateway string, isGlobal bool) *ArtifcatServer {

	return &ArtifcatServer{
		OrgId:          orgId,
		OrgName:        orgName,
		IsGlobal:       isGlobal,
		msgbus:         msgBus,
		baseRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
		pushGateway:    pushGateway,
		chunker:        chunk,
		storage:        storage,
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

	log.Infof("Got file %s with size %d", in.Name, len(in.Data))
	loc, err := s.storage.PutFile(ctx, in.Name, strings.ToLower(in.Type.String()), v, pkg.TarGzExtension, bytes.NewReader(in.Data))
	if err != nil {
		log.Errorf("Error storing artifact: %s %s", in.Name, in.Version)
		return nil, err
	}

	go func() {
		aType := strings.ToLower(in.Type.String())
		fPath := fmt.Sprintf("%s/%s/%s", in.Name, v.String(), pkg.TarGzExtension)
		cReq := &dpb.CreateChunkRequest{
			Name:    in.Name,
			Type:    aType,
			Version: v.String(),
			Store:   "s3+" + strings.TrimSuffix(loc, fPath),
		}
		log.Infof("Sending chunking request %+v", cReq)
		resp, err := s.chunker.CreateChunk(cReq)
		if err != nil {
			log.Errorf("Error chunking artifact: %s %s. Error: %+v", in.Name, in.Version, err)
			return
		}

		nctx, cancel := context.WithTimeout(context.Background(),
			s.storageRequestTimeout)
		defer cancel()

		_, err = s.storage.PutFile(nctx, in.Name, aType, v, pkg.ChunkIndexExtension, bytes.NewReader(resp.Index))
		if err != nil {
			log.Errorf("Failed to store artifact index file %+v", in)
		}

	}()

	return nil, nil
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
	defer rd.Close()

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
	log.Infof("Getting version list: %s of type %s", in.Name, in.Type)

	ls, err := s.storage.ListVersions(ctx, in.Name, strings.ToLower(in.Type.String()))
	if err != nil {
		return nil, err
	}

	if len(*ls) == 0 {
		return nil, status.Error(codes.NotFound, "Artifact name is not valid")
	}

	vers := []*pb.VersionInfo{}
	for _, v := range *ls {
		formats := []*pb.FormatInfo{
			{
				Url:       path.Join(CappsPath, in.Name, v.Version+pkg.TarGzExtension),
				CreatedAt: timestamppb.New(v.CreatedAt),
				Size:      v.SizeBytes,
				Type:      "tar.gz",
			},
		}

		if v.Chunked {
			formats = append(formats, &pb.FormatInfo{
				Url:  path.Join(CappsPath, in.Name, v.Version+pkg.ChunkIndexExtension),
				Type: "chunk",
				ExtraInfo: []*pb.ExtraInfoMap{
					{
						Key:   "chunks",
						Value: fmt.Sprintf("%s/", ChunksPath),
					},
				},
			})
		}

		vers = append(vers, &pb.VersionInfo{
			Version: v.Version,
			Formats: formats,
		})
	}

	return &pb.GetArtifactVersionListResponse{
		Name:     in.Name,
		Type:     in.Type,
		Versions: vers,
	}, nil

}

func (s *ArtifcatServer) ListArtifacts(ctx context.Context, in *pb.ListArtifactRequest) (*pb.ListArtifactResponse, error) {
	log.Infof("Getting list of %s artifacts", in.Type)

	ls, err := s.storage.ListApps(ctx, strings.ToLower(in.Type.String()))
	if err != nil {
		return nil, err
	}

	return &pb.ListArtifactResponse{
		Artifact: ls,
	}, nil
}
