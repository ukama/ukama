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
	"errors"
	"fmt"
	"strings"

	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/node/configurator/pkg"
	"github.com/ukama/ukama/systems/node/configurator/pkg/db"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	log "github.com/sirupsen/logrus"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	copr "github.com/ukama/ukama/systems/common/rest/client/operation"
	pb "github.com/ukama/ukama/systems/node/configurator/pb/gen"
	configstore "github.com/ukama/ukama/systems/node/configurator/pkg/configStore"
)

const (
	applyConfigType = "ApplyConfig"
	nodeKeyPrefix   = "node:"
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
	opManager              copr.ManagerClient
	opLeaseSecs            uint32
}

func NewConfiguratorServer(msgBus mb.MsgBusServiceClient, cfgDb db.ConfigRepo,
	cmtDb db.CommitRepo, configStore configstore.ConfigStoreProvider, orgName string,
	debug bool, opMgr copr.ManagerClient, leaseSecs uint32) *ConfiguratorServer {

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
	plan, err := c.configStore.PlanConfigStoreEvent(ctx)
	if err != nil {
		log.Errorf("Error while handling config store event.Error: %s", err.Error())
		return &pb.ConfigStoreEventResponse{}, err
	}

	if err := c.applyPlan(plan); err != nil {
		return nil, err
	}
	return &pb.ConfigStoreEventResponse{}, nil
}

func (c *ConfiguratorServer) ApplyConfig(ctx context.Context, req *pb.ApplyConfigRequest) (*pb.ApplyConfigResponse, error) {
	log.Infof("Received a request to apply config  %v", req)

	plan, err := c.configStore.PlanConfigCommit(ctx, req.Hash)
	if err != nil {
		log.Errorf("Error while handling apply config req commit %s.Error: %s",
			req.Hash, err.Error())
		return &pb.ApplyConfigResponse{}, status.Errorf(codes.Internal, "%v", err)
	}

	if err := c.applyPlan(plan); err != nil {
		return nil, err
	}
	return &pb.ApplyConfigResponse{}, nil
}

func (c *ConfiguratorServer) applyPlan(plan *configstore.CommitPlan) error {
	nodes := plan.Nodes()
	if len(nodes) == 0 {
		log.Infof("No node config to apply")
		return nil
	}

	ops, err := c.lockNodes(nodes)
	if err != nil {
		return err
	}

	for _, op := range ops {
		if _, err := c.opManager.MarkRunning(op.Id, op.FencingToken,
			pkg.ServiceName); err != nil {
			log.Warnf("%s mark running failed for op %s: %v", applyConfigType, op.Id, err)
			c.failAll(ops, fmt.Sprintf("mark running failed: %v", err))
			return status.Errorf(markRunningCode(err), "mark running: %v", err)
		}
	}

	results := c.configStore.Apply(plan)

	var failed []string
	for _, op := range ops {
		node := strings.TrimPrefix(op.ResourceKey, nodeKeyPrefix)
		if err := results[node]; err != nil {
			failed = append(failed, node)
			c.failOperation(op, err.Error())
			continue
		}
		c.completeOperation(op)
	}

	if len(failed) > 0 {
		return status.Errorf(codes.Internal, "config not applied to nodes %v", failed)
	}
	return nil
}

func (c *ConfiguratorServer) lockNodes(nodes []string) ([]*copr.OperationInfo, error) {
	keys := make([]string, 0, len(nodes))
	for _, n := range nodes {
		keys = append(keys, nodeKeyPrefix+n)
	}

	resp, err := c.opManager.Start(copr.StartRequest{
		Type:         applyConfigType,
		System:       pkg.SystemName,
		ResourceKeys: keys,
		RequestedBy:  pkg.ServiceName,
		LeaseSeconds: c.opLeaseSecs,
	})
	if err == nil && (resp == nil || len(resp.Operations) == 0) {
		err = fmt.Errorf("start returned no operations")
	}
	if errors.Is(err, copr.ErrBusy) {
		log.Warnf("%s rejected for %v: %v", applyConfigType, nodes, err)
		return nil, status.Errorf(codes.AlreadyExists, "%s rejected: %v",
			applyConfigType, err)
	}
	if err != nil {
		log.Errorf("%s lock for %v failed: %v", applyConfigType, nodes, err)
		return nil, status.Errorf(codes.Unavailable, "%s: lock nodes: %v",
			applyConfigType, err)
	}

	return resp.Operations, nil
}

func markRunningCode(err error) codes.Code {
	switch {
	case errors.Is(err, copr.ErrNotRunnable):
		return codes.Aborted
	case errors.Is(err, copr.ErrUnavailable):
		return codes.Unavailable
	default:
		return codes.Internal
	}
}

func (c *ConfiguratorServer) failAll(ops []*copr.OperationInfo, reason string) {
	for _, op := range ops {
		c.failOperation(op, reason)
	}
}

func (c *ConfiguratorServer) failOperation(op *copr.OperationInfo, reason string) {
	if _, err := c.opManager.Fail(op.Id, op.FencingToken, pkg.ServiceName,
		reason); err != nil {
		log.Errorf("%s failed to fail operation %s: %v", applyConfigType, op.Id, err)
	}
}

func (c *ConfiguratorServer) completeOperation(op *copr.OperationInfo) {
	if _, err := c.opManager.Complete(op.Id, op.FencingToken, pkg.ServiceName,
		"dispatched"); err != nil {
		log.Errorf("%s failed to complete operation %s: %v", applyConfigType, op.Id, err)
	}
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
