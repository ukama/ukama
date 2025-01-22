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

	"github.com/google/uuid"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
	pb "github.com/ukama/ukama/testing/services/dummy-node/node/pb/gen"
	"github.com/ukama/ukama/testing/services/dummy-node/node/pkg"
	"github.com/ukama/ukama/testing/services/dummy-node/node/pkg/db"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const uuidParsingError = "Error parsing UUID"

type NodeServer struct {
	pb.UnimplementedNodeServiceServer
	orgName        string
	nodeRepo       db.NodeRepo
	msgbus         mb.MsgBusServiceClient
	baseRoutingKey msgbus.RoutingKeyBuilder
	pushGateway    string
}

func NewNodeServer(orgName string, nodeRepo db.NodeRepo, msgBus mb.MsgBusServiceClient, pushGateway string) *NodeServer {
	return &NodeServer{
		orgName:        orgName,
		nodeRepo:       nodeRepo,
		msgbus:         msgBus,
		pushGateway:    pushGateway,
		baseRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
	}
}

func (s *NodeServer) ResetNode(ctx context.Context, req *pb.ResetRequest) (*pb.ResetResponse, error) {
	nodeID, err := uuid.Parse(req.NodeId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, uuidParsingError)
	}

	return &pb.ResetResponse{
		NodeId: nodeID.String(),
	}, nil
}

func (s *NodeServer) ToggleNodeRFOn(ctx context.Context, req *pb.NodeRFOnRequest) (*pb.NodeRFOnResponse, error) {
	nodeID, err := uuid.Parse(req.NodeId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, uuidParsingError)
	}

	return &pb.NodeRFOnResponse{
		NodeId: nodeID.String(),
	}, nil
}

func (s *NodeServer) NodeRFOff(ctx context.Context, req *pb.NodeRFOffRequest) (*pb.NodeRFOffResponse, error) {
	nodeID, err := uuid.Parse(req.NodeId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, uuidParsingError)
	}

	return &pb.NodeRFOffResponse{
		NodeId: nodeID.String(),
	}, nil
}

func (s *NodeServer) TurnNodeOff(ctx context.Context, req *pb.TurnNodeOffRequest) (*pb.TurnNodeOffResponse, error) {
	nodeID, err := uuid.Parse(req.NodeId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, uuidParsingError)
	}

	return &pb.TurnNodeOffResponse{
		NodeId: nodeID.String(),
	}, nil
}
