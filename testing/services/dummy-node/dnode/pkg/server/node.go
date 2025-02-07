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

	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/ukama"
	pb "github.com/ukama/ukama/testing/services/dummy-node/dnode/pb/gen"
	"github.com/ukama/ukama/testing/services/dummy-node/dnode/pkg"
	"github.com/ukama/ukama/testing/services/dummy-node/dnode/pkg/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type NodeServer struct {
	pb.UnimplementedNodeServiceServer
	orgName        string
	amqpConfig     pkg.AmqpConfig
	msgbus         mb.MsgBusServiceClient
	baseRoutingKey msgbus.RoutingKeyBuilder
}

func NewNodeServer(orgName string, msgBus mb.MsgBusServiceClient, amqpConfig pkg.AmqpConfig) *NodeServer {
	return &NodeServer{
		msgbus:         msgBus,
		orgName:        orgName,
		amqpConfig:     amqpConfig,
		baseRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem("messaging").SetOrgName(orgName).SetService("mesh"),
	}
}

func (s *NodeServer) ResetNode(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	nodeID, err := ukama.ValidateNodeId(req.GetNodeId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", err.Error())
	}

	utils.PushNodeResetViaREST(s.amqpConfig, s.orgName, nodeID.String(), s.msgbus)

	return &pb.Response{
		NodeId: nodeID.String(),
	}, nil
}

func (s *NodeServer) NodeRFOn(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	nodeID, err := ukama.ValidateNodeId(req.GetNodeId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", err.Error())
	}

	utils.PushNodeRFOnViaREST(s.amqpConfig, s.orgName, nodeID.String(), s.msgbus)

	return &pb.Response{
		NodeId: nodeID.String(),
	}, nil
}

func (s *NodeServer) NodeRFOff(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	nodeID, err := ukama.ValidateNodeId(req.GetNodeId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", err.Error())
	}

	utils.PushNodeRFOffViaREST(s.amqpConfig, s.orgName, nodeID.String(), s.msgbus)

	return &pb.Response{
		NodeId: nodeID.String(),
	}, nil
}

func (s *NodeServer) TurnNodeOff(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	nodeID, err := ukama.ValidateNodeId(req.GetNodeId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", err.Error())
	}

	utils.PushNodeOffViaREST(s.amqpConfig, s.orgName, nodeID.String(), s.msgbus)

	return &pb.Response{
		NodeId: nodeID.String(),
	}, nil
}

func (s *NodeServer) TurnNodeOnline(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	nodeID, err := ukama.ValidateNodeId(req.GetNodeId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", err.Error())
	}

	utils.PushNodeOnlineViaREST(s.amqpConfig, s.orgName, nodeID.String(), s.msgbus)

	return &pb.Response{
		NodeId: nodeID.String(),
	}, nil
}
