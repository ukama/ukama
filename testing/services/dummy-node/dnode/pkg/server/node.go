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
	"github.com/ukama/ukama/testing/services/dummy-node/dnode/pkg/db"
	"github.com/ukama/ukama/testing/services/dummy-node/dnode/pkg/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type NodeServer struct {
	pb.UnimplementedNodeServiceServer
	orgName        string
	nodeId         string
	nodeRepo       db.NodeRepo
	msgbus         mb.MsgBusServiceClient
	baseRoutingKey msgbus.RoutingKeyBuilder
}

func NewNodeServer(orgName string, nodeId string, nodeRepo db.NodeRepo, msgBus mb.MsgBusServiceClient) *NodeServer {
	return &NodeServer{
		msgbus:         msgBus,
		nodeId:         nodeId,
		orgName:        orgName,
		nodeRepo:       nodeRepo,
		baseRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem("messaging").SetOrgName(orgName).SetService("mesh"),
	}
}

func (s *NodeServer) ResetNode(ctx context.Context, req *pb.ResetRequest) (*pb.ResetResponse, error) {
	nodeID, err := ukama.ValidateNodeId(req.GetNodeId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", err.Error())
	}

	utils.PushNodeResetViaREST(s.orgName, s.nodeId, s.msgbus)

	return &pb.ResetResponse{
		NodeId: nodeID.String(),
	}, nil
}

func (s *NodeServer) NodeRFOn(ctx context.Context, req *pb.NodeRFOnRequest) (*pb.NodeRFOnResponse, error) {
	nodeID, err := ukama.ValidateNodeId(req.GetNodeId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", err.Error())
	}

	utils.PushNodeRFOnViaREST(s.orgName, s.nodeId, s.msgbus)

	return &pb.NodeRFOnResponse{
		NodeId: nodeID.String(),
	}, nil
}

func (s *NodeServer) NodeRFOff(ctx context.Context, req *pb.NodeRFOffRequest) (*pb.NodeRFOffResponse, error) {
	nodeID, err := ukama.ValidateNodeId(req.GetNodeId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", err.Error())
	}

	utils.PushNodeRFOffViaREST(s.orgName, s.nodeId, s.msgbus)

	return &pb.NodeRFOffResponse{
		NodeId: nodeID.String(),
	}, nil
}

func (s *NodeServer) TurnNodeOff(ctx context.Context, req *pb.TurnNodeOffRequest) (*pb.TurnNodeOffResponse, error) {
	nodeID, err := ukama.ValidateNodeId(req.GetNodeId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", err.Error())
	}

	utils.PushNodeOffViaREST(s.orgName, s.nodeId, s.msgbus)

	return &pb.TurnNodeOffResponse{
		NodeId: nodeID.String(),
	}, nil
}
