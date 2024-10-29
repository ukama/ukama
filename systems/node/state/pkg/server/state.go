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
	"sort"

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/grpc"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"

	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
	pb "github.com/ukama/ukama/systems/node/state/pb/gen"
	"github.com/ukama/ukama/systems/node/state/pkg"
	"github.com/ukama/ukama/systems/node/state/pkg/db"
)

type StateServer struct {
	pb.UnimplementedStateServiceServer
	orgName         string
	orgId           string
	sRepo           db.StateRepo
	StateRoutingKey msgbus.RoutingKeyBuilder
	msgbus          mb.MsgBusServiceClient
}

func NewStateServer(orgName string, orgId string, sRepo db.StateRepo, msgBus mb.MsgBusServiceClient) *StateServer {

	ns := &StateServer{
		sRepo:           sRepo,
		orgName:         orgName,
		msgbus:          msgBus,
		orgId:           orgId,
		StateRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
	}

	return ns
}

func (s *StateServer) AddNodeState(ctx context.Context, req *pb.AddStateRequest) (*pb.AddStateResponse, error) {
	log.Infof("Adding nodeState for Node ID: %v", req.NodeId)
	nId, err := ukama.ValidateNodeId(req.NodeId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of node id: %s", err.Error())
	}

	events := db.StringArray(req.Events)
	subState := db.StringArray(req.SubState)

	currentState, err := s.sRepo.GetLatestState(nId.String())
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Errorf("Error retrieving current state: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to retrieve current state: %v", err)
	}

	newNodeState := &db.State{
		Id:           uuid.NewV4(),
		NodeId:       nId.String(),
		CurrentState: req.CurrentState,
		SubState:     subState,
		Events:       events,
		NodeType:     req.GetNodeType(),
		NodeIp:       req.NodeIp,
		NodePort:     req.NodePort,
		MeshIp:       req.MeshIp,
		MeshPort:     req.MeshPort,
		MeshHostName: req.MeshHostName,
	}

	if currentState != nil {
		newNodeState.PreviousStateId = &currentState.Id
	}

	err = s.sRepo.AddState(newNodeState, currentState)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to add node state: %v", err)
	}

	return &pb.AddStateResponse{
		Id: newNodeState.Id.String(),
	}, nil
}

func (s *StateServer) GetStates(ctx context.Context, req *pb.GetStatesRequest) (*pb.GetStatesResponse, error) {
	log.Infof("Getting node states for Node ID: %v", req.NodeId)

	nId, err := ukama.ValidateNodeId(req.NodeId)
	if err != nil {
		return &pb.GetStatesResponse{}, status.Errorf(codes.InvalidArgument,
			"invalid format of node id: %s", err.Error())
	}
	history, err := s.sRepo.GetStateHistory(nId.String())
	if err != nil {
		log.Errorf("Failed to get node state history: %v", err)
		return &pb.GetStatesResponse{}, status.Errorf(codes.Internal, "failed to get node state history: %v", err)
	}

	if history == nil {
		log.Warnf("No history found for Node ID: %v", req.NodeId)
		return &pb.GetStatesResponse{States: []*pb.State{}}, nil
	}

	stateMap := make(map[string]*pb.State)
	var latestStates []*pb.State

	for _, nodeState := range history {
		NodeStateRes := &pb.State{
			Id:           nodeState.Id.String(),
			NodeId:       nodeState.NodeId,
			CurrentState: nodeState.CurrentState,
			SubState:     nodeState.SubState,
			Events:       nodeState.Events,
			CreatedAt:    timestamppb.New(nodeState.CreatedAt),
			UpdatedAt:    timestamppb.New(nodeState.UpdatedAt),
		}

		if nodeState.PreviousStateId != nil {
			NodeStateRes.PreviousStateId = nodeState.PreviousStateId.String()
		}

		stateMap[NodeStateRes.Id] = NodeStateRes

		isLatest := true
		for _, s := range history {
			if s.PreviousStateId != nil && s.PreviousStateId.String() == NodeStateRes.Id {
				isLatest = false
				break
			}
		}
		if isLatest {
			latestStates = append(latestStates, NodeStateRes)
		}
	}

	for _, state := range stateMap {
		if state.PreviousStateId != "" {
			if prevState, exists := stateMap[state.PreviousStateId]; exists {
				state.PreviousState = prevState
			}
		}
	}

	sort.Slice(latestStates, func(i, j int) bool {
		return latestStates[i].UpdatedAt.AsTime().After(latestStates[j].UpdatedAt.AsTime())
	})

	return &pb.GetStatesResponse{
		States: latestStates,
	}, nil
}
func (s *StateServer) GetLatestState(ctx context.Context, req *pb.GetLatestStateRequest) (*pb.GetLatestStateResponse, error) {
	log.Infof("Getting latest node state for Node ID: %v", req.NodeId)

	if req.GetNodeId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "node ID cannot be empty")
	}

	nId, err := ukama.ValidateNodeId(req.GetNodeId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid format of node id: %v", err)
	}

	latestState, err := s.sRepo.GetLatestState(nId.String())
	if err != nil {
		log.Errorf("Error getting the node state: %v", err)
		return nil, grpc.SqlErrorToGrpc(err, "node state")
	}

	if latestState == nil {
		return &pb.GetLatestStateResponse{}, nil
	}

	stateRes := &pb.State{
		Id:           latestState.Id.String(),
		NodeId:       latestState.NodeId,
		CurrentState: latestState.CurrentState,
		SubState:     latestState.SubState,
		Events:       latestState.Events,
		CreatedAt:    timestamppb.New(latestState.CreatedAt),
		UpdatedAt:    timestamppb.New(latestState.UpdatedAt),
	}

	if latestState.PreviousStateId != nil {
		stateRes.PreviousStateId = latestState.PreviousStateId.String()
	}

	return &pb.GetLatestStateResponse{
		State: stateRes,
	}, nil
}

func (s *StateServer) UpdateState(ctx context.Context, req *pb.UpdateStateRequest) (*pb.UpdateStateResponse, error) {
	log.Infof("Updating node state for Node ID: %v", req.NodeId)

	nId, err := ukama.ValidateNodeId(req.NodeId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid format of node id: %s", err.Error())
	}

	currentState, err := s.sRepo.GetLatestState(nId.String())
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Errorf(codes.NotFound, "state not found for Node ID: %s", req.NodeId)
		}
		log.Errorf("Error retrieving current state: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to retrieve current state: %v", err)
	}

	if currentState == nil {
		return nil, status.Errorf(codes.NotFound, "state not found for Node ID: %s", req.NodeId)
	}

	updatedSubState := append(currentState.SubState, db.StringArray(req.SubState)...)
	updatedEvents := append(currentState.Events, db.StringArray(req.Events)...)

	if _, err := s.sRepo.UpdateState(nId.String(), updatedSubState, updatedEvents); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update node state: %v", err)
	}

	return &pb.UpdateStateResponse{}, nil
}
