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

	"github.com/Masterminds/semver/v3"
	casync "github.com/folbricht/desync"
	log "github.com/sirupsen/logrus"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	uuid "github.com/ukama/ukama/systems/common/uuid"
	pb "github.com/ukama/ukama/systems/hub/distributor/pb/gen"
	"github.com/ukama/ukama/systems/hub/distributor/pkg"
	"github.com/ukama/ukama/systems/hub/distributor/pkg/chunk"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const CappsPath = "/v1/capps"
const ChunksPath = "/v1/chunks"

type DistributorServer struct {
	pb.DistributorServiceServer
	msgbus         mb.MsgBusServiceClient
	baseRoutingKey msgbus.RoutingKeyBuilder
	pushGateway    string
	OrgId          uuid.UUID
	OrgName        string
	Store          pkg.StoreConfig
	ChunkConfig    pkg.ChunkConfig
	IsGlobal       bool
	castore        casync.Store
	converters     casync.Converters
}

func NewDistributionServer(orgId uuid.UUID, orgName string, config *pkg.Config,
	msgBus mb.MsgBusServiceClient, pushGateway string, isGlobal bool) *DistributorServer {

	// s, c, err := distribution.InitDistribution(&config.Distribution)
	// if err != nil {
	// 	log.Fatalf("Distribution initialization failed.")
	// }
	//compressed := c.hasCompression()

	return &DistributorServer{
		OrgId:          orgId,
		OrgName:        orgName,
		IsGlobal:       isGlobal,
		msgbus:         msgBus,
		baseRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
		pushGateway:    pushGateway,
		Store:          config.Distribution.StoreCfg,
		ChunkConfig:    config.Distribution.Chunk,
		// castore:        s,
		// converters:     c,
	}
}

func (s *DistributorServer) parseVersion(version string) (*semver.Version, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid version format. Refer to https://semver.org/ for more information")
	}
	return v, err
}

func (s *DistributorServer) CreateChunk(ctx context.Context, in *pb.CreateChunkRequest) (*pb.CreateChunkResponse, error) {

	var bSize int64
	fname := in.Name

	ver, err := s.parseVersion(in.Version)
	if err != nil {
		return nil, err
	}

	log.Debugf("Handling chunking request %+v.", in)

	buf := new(bytes.Buffer)

	index, err := chunk.CreateChunks(ctx, &s.Store, &s.ChunkConfig, fname, in.Type, ver, in.Store)
	if err != nil {
		log.Errorf("Error while chunking the file %s: %s", in.Name, err.Error())
		return nil, status.Error(codes.Internal, "Error while creating chunks:"+err.Error())
	} else {

		if index != nil {
			bSize, err = index.WriteTo(buf)
			if err != nil {
				return nil, status.Error(codes.Internal, "Error while copying index file:"+err.Error())
			}
		}

		if err != nil {
			log.Errorf("Error while creating index file.")
			return nil, status.Error(codes.Internal, "Error while creating index file:"+err.Error())
		}

		capp := &epb.CappCreatedEvent{
			Name:    fname,
			Version: ver.String(),
		}

		route := s.baseRoutingKey.SetAction("create").SetObject("capp").MustBuild()

		err = s.msgbus.PublishRequest(route, capp)
		if err != nil {
			log.Errorf("Failed to publish message %+v with key %+v. Errors %s", capp, route, err.Error())
		}
	}

	return &pb.CreateChunkResponse{
		Index: buf.Bytes(),
		Size:  bSize,
	}, nil
}

func (s *DistributorServer) GetChunk(ctx context.Context, in *pb.GetChunkRequest) (*pb.GetChunkResponse, error) {
	log.Debugf("Handling get chunk request %+v.", in)
	var b []byte
	// chunk, err := s.castore.GetChunk(in.Id)
	// if err == nil {
	// 	// Optimization for when the chunk modifiers match those
	// 	// of the chunk server. In that case it's not necessary
	// 	// to convert back and forth. Just use the raw data as loaded
	// 	// from the store.
	// 	if len(chunk.storage) > 0 && s.converters.equal(chunk.converters) {
	// 		b = chunk.storage
	// 	} else {
	// 		b, err = chunk.Data()
	// 		if err == nil {
	// 			b, err = s.converters.toStorage(b)
	// 		}
	// 	}
	// }

	return &pb.GetChunkResponse{
		Id:   in.Id,
		Data: b,
	}, nil
}

// func (s *DistributorServer) idFromPath(p string) (casync.ChunkID, error) {
// 	ext := casync.CompressedChunkExt
// 	if !s.compressed {
// 		if strings.HasSuffix(p, casync.CompressedChunkExt) {
// 			return casync.ChunkID{}, errors.New("compressed chunk requested from http chunk store serving uncompressed chunks")
// 		}
// 		ext = casync.UncompressedChunkExt
// 	}
// 	sID := strings.TrimSuffix(path.Base(p), ext)
// 	if len(sID) < 4 {
// 		return casync.ChunkID{}, fmt.Errorf("expected format '/<prefix>/<chunkid>%s", ext)
// 	}

// 	// Make sure the prefix does match the first characters of the ID.
// 	if p != path.Join("/", sID[0:4], sID+ext) {
// 		return casync.ChunkID{}, fmt.Errorf("expected format '/<prefix>/<chunkid>%s", ext)
// 	}
// 	return casync.ChunkIDFromString(sID)
// }
