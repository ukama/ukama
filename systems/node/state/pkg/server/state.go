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
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/grpc"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
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
	log.Infof("Adding nodeState for Node ID: %v with state: %v, subState: %v, events: %v",
		req.NodeId, req.CurrentState, req.SubState, req.Events)

	nId, err := ukama.ValidateNodeId(req.NodeId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of node id: %s", err.Error())
	}

	events := db.StringArray(req.Events)
	subState := db.StringArray(req.SubState)

	currentState, err := s.sRepo.GetLatestState(nId.String())
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, status.Errorf(codes.Internal, "failed to retrieve current state: %v", err)
	}

	config := &db.NodeConfig{
		Id:           uuid.NewV4(),
		NodeId:       nId.String(),
		NodeIp:       req.NodeIp,
		NodePort:     req.NodePort,
		MeshIp:       req.MeshIp,
		MeshPort:     req.MeshPort,
		MeshHostName: req.MeshHostName,
	}

	newNodeState := &db.State{
		Id:           uuid.NewV4(),
		NodeId:       nId.String(),
		CurrentState: req.CurrentState,
		SubState:     subState,
		Events:       events,
		NodeType:     req.GetNodeType(),
		ConfigId:     config.Id,
		Config:       config,
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

	states := make([]*pb.State, 0, len(history))
	for _, nodeState := range history {
		state := &pb.State{
			Id:           nodeState.Id.String(),
			NodeId:       nodeState.NodeId,
			CurrentState: nodeState.CurrentState,
			SubState:     nodeState.SubState,
			Events:       nodeState.Events,
			CreatedAt:    timestamppb.New(nodeState.CreatedAt),
			UpdatedAt:    timestamppb.New(nodeState.UpdatedAt),
		}

		if nodeState.PreviousStateId != nil {
			state.PreviousStateId = nodeState.PreviousStateId.String()
		}

		states = append(states, state)
	}

	sort.Slice(states, func(i, j int) bool {
		return states[i].UpdatedAt.AsTime().After(states[j].UpdatedAt.AsTime())
	})

	nodeConfig, err := s.sRepo.GetNodeConfig(nId.String())
	if err != nil {
		log.Errorf("Failed to get node configuration: %v", err)
		return &pb.GetStatesResponse{}, status.Errorf(codes.Internal, "failed to get node configuration: %v", err)
	}

	return &pb.GetStatesResponse{
		States:     states,
		NodeConfig: convertToGenNodeConfig(nodeConfig),
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
	log.Infof("Updating node state for Node ID: %v with subState: %v, events: %v",
		req.NodeId, req.SubState, req.Events)

	nId, err := ukama.ValidateNodeId(req.NodeId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid format of node id: %s", err.Error())
	}

	currentState, err := s.sRepo.GetLatestState(nId.String())
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Errorf(codes.NotFound, "state not found for Node ID: %s", req.NodeId)
		}
		return nil, status.Errorf(codes.Internal, "failed to retrieve current state: %v", err)
	}

	if currentState == nil {
		return nil, status.Errorf(codes.NotFound, "state not found for Node ID: %s", req.NodeId)
	}

	var updatedSubState = currentState.SubState
	var updatedEvents = currentState.Events

	for _, substate := range req.SubState {
		if substate != "" {
			updatedSubState = append(updatedSubState, substate)
		}
	}

	for _, event := range req.Events {
		if event != "" {
			updatedEvents = append(updatedEvents, event)
		}
	}

	updatedState, err := s.sRepo.UpdateState(nId.String(), updatedSubState, updatedEvents)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update node state: %v", err)
	}

	return &pb.UpdateStateResponse{
		UpdatedState: &pb.State{
			Id:           updatedState.Id.String(),
			NodeId:       updatedState.NodeId,
			CurrentState: updatedState.CurrentState,
			SubState:     updatedState.SubState,
			Events:       updatedState.Events,
			CreatedAt:    timestamppb.New(updatedState.CreatedAt),
			UpdatedAt:    timestamppb.New(updatedState.UpdatedAt),
		},
	}, nil
}
func (s *StateServer) GetStatesHistory(ctx context.Context, req *pb.GetStatesHistoryRequest) (*pb.GetStatesHistoryResponse, error) {
	log.Infof("Getting node state history for Node ID: %v", req.NodeId)

	nId, err := ukama.ValidateNodeId(req.NodeId)
	if err != nil {
		return &pb.GetStatesHistoryResponse{}, status.Errorf(codes.InvalidArgument,
			"invalid format of node id: %s", err.Error())
	}

	var from, to time.Time

	if req.GetStartTime() != "" {
		from, err = time.Parse(time.RFC3339, req.GetStartTime())
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid time format for start_time. Error: %s", err.Error())
		}
	}

	if req.GetEndTime() != "" {
		to, err = time.Parse(time.RFC3339, req.GetEndTime())
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid time format for end_time. Error: %s", err.Error())
		}
	}

	history, err := s.sRepo.GetStateHistoryWithFilter(nId.String(), int(req.PageSize), int(req.PageNumber), from, to)
	if err != nil {
		log.Errorf("Failed to get node state history: %v", err)
		return &pb.GetStatesHistoryResponse{}, status.Errorf(codes.Internal, "failed to get node state history: %v", err)
	}

	if history == nil {
		log.Warnf("No history found for Node ID: %v", req.NodeId)
		return &pb.GetStatesHistoryResponse{States: []*pb.State{}}, nil
	}

	states := make([]*pb.State, 0, len(history))
	for _, nodeState := range history {
		state := &pb.State{
			Id:           nodeState.Id.String(),
			NodeId:       nodeState.NodeId,
			CurrentState: nodeState.CurrentState,
			SubState:     nodeState.SubState,
			Events:       nodeState.Events,
			CreatedAt:    timestamppb.New(nodeState.CreatedAt),
			UpdatedAt:    timestamppb.New(nodeState.UpdatedAt),
		}

		if nodeState.PreviousStateId != nil {
			state.PreviousStateId = nodeState.PreviousStateId.String()
		}

		states = append(states, state)
	}

	sort.Slice(states, func(i, j int) bool {
		return states[i].UpdatedAt.AsTime().After(states[j].UpdatedAt.AsTime())
	})

	nodeConfig, err := s.sRepo.GetNodeConfig(nId.String())
	if err != nil {
		log.Errorf("Failed to get node configuration: %v", err)
		return &pb.GetStatesHistoryResponse{}, status.Errorf(codes.Internal, "failed to get node configuration: %v", err)
	}

	return &pb.GetStatesHistoryResponse{
		States:     states,
		NodeConfig: convertToGenNodeConfig(nodeConfig),
	}, nil
}
func (s *StateServer) EnforceStateTransition(ctx context.Context, req *pb.EnforceStateTransitionRequest) (*pb.EnforceStateTransitionResponse, error) {
	log.Infof("Enforcing state transition for Node ID: %v", req.NodeId)

	nId, err := ukama.ValidateNodeId(req.NodeId)
	if err != nil {
		return &pb.EnforceStateTransitionResponse{}, status.Errorf(codes.InvalidArgument,
			"invalid format of node id: %s", err.Error())
	}
	if s.msgbus != nil {
		route := s.StateRoutingKey.SetAction("force").SetObject("node").MustBuild()

		evt := &epb.EnforceNodeStateEvent{
			NodeId: nId.StringLowercase(),
			Event:  req.Event,
		}

		err = s.msgbus.PublishRequest(route, evt)
		if err != nil {
			log.Errorf("Failed to publish message %+v with key %+v. Errors %s", evt, route, err.Error())
		}
	}

	return &pb.EnforceStateTransitionResponse{}, nil
}

func convertToGenNodeConfig(dbConfig *db.NodeConfig) *pb.NodeConfig {
	if dbConfig == nil {
		return nil
	}
	return &pb.NodeConfig{
		Id:           dbConfig.Id.String(),
		NodeId:       dbConfig.NodeId,
		NodeIp:       dbConfig.NodeIp,
		NodePort:     dbConfig.NodePort,
		MeshIp:       dbConfig.MeshIp,
		MeshPort:     dbConfig.MeshPort,
		MeshHostName: dbConfig.MeshHostName,
		CreatedAt:    timestamppb.New(dbConfig.CreatedAt),
		UpdatedAt:    timestamppb.New(dbConfig.UpdatedAt),
	}
}
