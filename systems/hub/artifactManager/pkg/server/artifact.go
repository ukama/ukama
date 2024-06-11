/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package server

import (
	"context"
	"time"

	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/hub/artifactmanager/pkg"

	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	uuid "github.com/ukama/ukama/systems/common/uuid"
	pb "github.com/ukama/ukama/systems/hub/artifactmanager/pb/gen"
)

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
	chunker               pkg.Chunker
	orgName               string
	IsGlobal              bool
}

func NewArtifactServer(orgId uuid.UUID, orgName string, storage pkg.Storage, chunker pkg.Chunker, storageTimeout time.Duration,
	msgBus mb.MsgBusServiceClient, pushGateway string, isGlobal bool) *ArtifcatServer {

	return &ArtifcatServer{
		OrgId:          orgId,
		OrgName:        orgName,
		IsGlobal:       isGlobal,
		msgbus:         msgBus,
		baseRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
		pushGateway:    pushGateway,
	}
}

func (s *ArtifcatServer) StoreArtifact(ctx context.Context, in *pb.StoreArtifactRequest) (*pb.StoreArtifactResponse, error) {
	return nil, nil
}
func (s *ArtifcatServer) GetArtifactLocation(ctx context.Context, in *pb.GetArtifactLocationRequest) (*pb.GetArtifactLocationResponse, error) {
	return nil, nil
}

func (s *ArtifcatServer) GetArtifact(ctx context.Context, in *pb.GetArtifactRequest) (*pb.GetArtifactResponse, error) {
	return nil, nil
}

func (s *ArtifcatServer) GetArtifcatVersionList(ctx context.Context, in *pb.GetArtifactVersionListRequest) (*pb.GetArtifactVersionListResponse, error) {
	return nil, nil
}

func (s *ArtifcatServer) ListArtifacts(ctx context.Context, in *pb.ListArtifactRequest) (*pb.ListArtifactResponse, error) {
	return nil, nil
}
