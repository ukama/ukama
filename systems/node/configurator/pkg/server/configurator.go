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

	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/node/configurator/pkg"
	cfgclient "github.com/ukama/ukama/systems/node/configurator/pkg/client"
	"github.com/ukama/ukama/systems/node/configurator/pkg/db"

	log "github.com/sirupsen/logrus"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	pb "github.com/ukama/ukama/systems/node/configurator/pb/gen"
	configstore "github.com/ukama/ukama/systems/node/configurator/pkg/configStore"
	opmgrpb "github.com/ukama/ukama/systems/operation/manager/pb/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ConfiguratorServer struct {
	pb.UnimplementedConfiguratorServiceServer
	msgbus                 mb.MsgBusServiceClient
	configuratorRoutingKey msgbus.RoutingKeyBuilder
	debug                  bool
	orgName                string
	configStore            configstore.ConfigStoreProvider
	commitRepo             db.CommitRepo
	configRepo             db.ConfigRepo
	opManager              cfgclient.OperationManager
	opLeaseSecs            uint32
}

func NewConfiguratorServer(msgBus mb.MsgBusServiceClient, cfgDb db.ConfigRepo, cmtDb db.CommitRepo, configStore configstore.ConfigStoreProvider, orgName string, debug bool, opMgr cfgclient.OperationManager, leaseSecs uint32) *ConfiguratorServer {

	log.Infof("Config store created: %+v", configStore)
	return &ConfiguratorServer{
		configuratorRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
		msgbus:                 msgBus,
		debug:                  pkg.IsDebugMode,
		orgName:                orgName,
		configStore:            configStore,
		commitRepo:             cmtDb,
		configRepo:             cfgDb,
		opManager:              opMgr,
		opLeaseSecs:            leaseSecs,
	}
}

func (c *ConfiguratorServer) ConfigEvent(ctx context.Context, req *pb.ConfigStoreEvent) (*pb.ConfigStoreEventResponse, error) {
	log.Infof("Received a event from config store %v", req)
	err := c.configStore.HandleConfigStoreEvent(ctx)
	if err != nil {
		log.Errorf("Error while handling config store event.Error: %s", err.Error())
	}
	return &pb.ConfigStoreEventResponse{}, err
}

func (c *ConfiguratorServer) ApplyConfig(ctx context.Context, req *pb.ApplyConfigRequest) (*pb.ApplyConfigResponse, error) {
	log.Infof("Received a request to apply config  %v", req)

	startResp, err := c.opManager.Start(&opmgrpb.StartOperationRequest{
		Type:         "ApplyConfig",
		System:       "node",
		ResourceKey:  "config:apply",
		RequestedBy:  pkg.ServiceName,
		LeaseSeconds: c.opLeaseSecs,
	})
	if err != nil {
		log.Warnf("ApplyConfig lock acquire rejected: %v", err)
		return nil, err
	}
	op := startResp.Operation
	log.Infof("ApplyConfig acquired lock op=%s token=%d", op.Id, op.FencingToken)

	if err := c.configStore.HandleConfigCommitReq(ctx, req.Hash); err != nil {
		log.Errorf("Error while handling apply config req commit %s.Error: %s", req.Hash, err.Error())
		return &pb.ApplyConfigResponse{}, status.Errorf(codes.Internal, "%v", err)
	}
	if _, err := c.opManager.MarkRunning(op.Id, op.FencingToken); err != nil {
		log.Warnf("ApplyConfig mark running failed for op %s: %v", op.Id, err)
	}
	return &pb.ApplyConfigResponse{}, nil
}

func (c *ConfiguratorServer) GetConfigVersion(ctx context.Context, req *pb.ConfigVersionRequest) (*pb.ConfigVersionResponse, error) {
	log.Infof("Received a request to get config for node  %v", req)
	cfg, err := c.configRepo.Get(req.NodeId)
	if err != nil {
		log.Errorf("Error while reading config for node %s. Error: %s", req.NodeId, err.Error())
	}

	return &pb.ConfigVersionResponse{
		NodeId:     req.NodeId,
		Status:     cfg.State.String(),
		Commit:     cfg.Commit.Hash,
		LastStatus: cfg.LastCommitState.String(),
		LastCommit: cfg.LastCommit.Hash,
	}, err
}
