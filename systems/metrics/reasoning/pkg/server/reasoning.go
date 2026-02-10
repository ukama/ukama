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
	"github.com/ukama/ukama/systems/metrics/reasoning/scheduler"

	log "github.com/sirupsen/logrus"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	pb "github.com/ukama/ukama/systems/metrics/reasoning/pb/gen"
)

type ReasoningServer struct {
	pb.UnimplementedReasoningServiceServer
	msgbus             mb.MsgBusServiceClient
	baseRoutingKey     msgbus.RoutingKeyBuilder
	reasoningScheduler scheduler.ReasoningScheduler
	config             *pkg.Config
}

func NewReasoningServer(msgBus mb.MsgBusServiceClient, reasoningScheduler scheduler.ReasoningScheduler, config *pkg.Config) *ReasoningServer {
	c := &ReasoningServer{
		msgbus:             msgBus,
		config:             config,
		reasoningScheduler: reasoningScheduler,
		baseRoutingKey:  msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(config.OrgName).SetService(pkg.ServiceName),
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

