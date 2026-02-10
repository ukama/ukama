/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package server

import (
	"context"

	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/metrics/reasoning/pkg"
	"github.com/ukama/ukama/systems/metrics/reasoning/pkg/store"
	"github.com/ukama/ukama/systems/metrics/reasoning/pkg/utils"
	"github.com/ukama/ukama/systems/metrics/reasoning/scheduler"

	log "github.com/sirupsen/logrus"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	creg "github.com/ukama/ukama/systems/common/rest/client/registry"
	ukama "github.com/ukama/ukama/systems/common/ukama"
	pb "github.com/ukama/ukama/systems/metrics/reasoning/pb/gen"
)

const (
	jobTag               = "reasoning-job"
)

type ReasoningServer struct {
	pb.UnimplementedReasoningServiceServer
	msgbus             mb.MsgBusServiceClient
	baseRoutingKey     msgbus.RoutingKeyBuilder
	reasoningScheduler scheduler.ReasoningScheduler
	config             *pkg.Config
	nodeClient         creg.NodeClient
	store    		   *store.Store
}

func NewReasoningServer(msgBus mb.MsgBusServiceClient, nodeClient creg.NodeClient, config *pkg.Config, store *store.Store) *ReasoningServer {
	scheduler := scheduler.NewReasoningScheduler(config.SchedulerInterval)
	c := &ReasoningServer{
		store:              store,
		msgbus:             msgBus,
		config:             config,
		reasoningScheduler: scheduler,
		nodeClient:         nodeClient,
		baseRoutingKey:  msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(config.OrgName).SetService(pkg.ServiceName),
	}

	// Start the scheduler
	if err := c.reasoningScheduler.Start(jobTag, func() { c.ReasoningJob(context.Background()) }); err != nil {
		log.Errorf("Failed to start the initial scheduler: %v", err)
	}

	return c
}

func (c *ReasoningServer) GetStats(ctx context.Context, req *pb.GetStatsRequest) (*pb.GetStatsResponse, error) {
	return &pb.GetStatsResponse{}, nil
}

func (c *ReasoningServer) GetDomains(ctx context.Context, req *pb.GetDomainsRequest) (*pb.GetDomainsResponse, error) {
	return &pb.GetDomainsResponse{}, nil
}

func (c *ReasoningServer) StartScheduler(ctx context.Context, req *pb.StartSchedulerRequest) (*pb.StartSchedulerResponse, error) {
	log.Info("Starting scheduler")

	return &pb.StartSchedulerResponse{}, nil
}

func (c *ReasoningServer) StopScheduler(ctx context.Context, req *pb.StopSchedulerRequest) (*pb.StopSchedulerResponse, error) {
	log.Info("Stopping scheduler")

	return &pb.StopSchedulerResponse{}, nil
}


func (c *ReasoningServer) ReasoningJob(ctx context.Context) {
	log.Infof("Reasoning job started")
	nodes, err := c.nodeClient.List(creg.ListNodesRequest{
		Connectivity: ukama.NodeConnectivityOnline.String(),
		State: ukama.NodeStateConfigured.String(),
		Type: ukama.NodeType(ukama.NODE_ID_TYPE_TOWERNODE).String(),
	})
	if err != nil {
		log.Errorf("Failed to get nodes: %v", err)
		return
	}

	log.Infof("Node registry nodes: %v", nodes.Nodes)

	for _, node := range nodes.Nodes {
		n, err := utils.SortNodeIds(node.Id)
		if err != nil {
			log.Errorf("Failed to sort nodes: %v", err)
			continue
		}

		log.Infof("Sorted nodes: %v", n)

		toN, fromN, err := utils.GetToNFromStore(c.store, n.TNode)
		if err != nil {
			log.Errorf("Failed to get To and From value: %v", err)
			continue
		}
		log.Infof("To and From value: %s, %s", toN, fromN)
	}
}

