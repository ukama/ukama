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

	log "github.com/sirupsen/logrus"
	stm "github.com/ukama/ukama/systems/common/stateMachine"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"

	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
	pb "github.com/ukama/ukama/systems/node/state/pb/gen"
	"github.com/ukama/ukama/systems/node/state/pkg/db"
)
  
 type NodeStateServer struct {
	 pb.UnimplementedNodeStateServiceServer
	 sRepo               db.NodeStateRepo
	 nodeStateRoutingKey msgbus.RoutingKeyBuilder
	 msgbus              mb.MsgBusServiceClient
	 stateMachine        *stm.StateMachine 
	 debug               bool
	 orgName             string
 }
  
 func NewNodeStateServer(orgName string,sRepo db.NodeStateRepo, msgBus mb.MsgBusServiceClient, debug bool, configPath string) *NodeStateServer {
	 ns := &NodeStateServer{
		 sRepo:   sRepo,
		 orgName: orgName,
		 msgbus:  msgBus,
		 debug:   debug,
	 }
 
	 if err := ns.InitializeStateMachine(configPath); err != nil {
		 log.Fatalf("Failed to initialize state machine: %v", err)
	 }
 
	 return ns
 }
  
 func (s *NodeStateServer) AddNodeState(ctx context.Context, req *pb.AddNodeStateRequest) (*pb.AddNodeStateResponse, error) {
    log.Infof("Adding nodeState for Node ID: %v", req.NodeId)
    nId, err := ukama.ValidateNodeId(req.NodeId)
    if err != nil {
        return nil, status.Errorf(codes.InvalidArgument,
            "invalid format of node id: %s", err.Error())
    }

    events := db.StringArray(req.Events)

    nextState, err := s.stateMachine.GetNextState(req.CurrentState, req.Events)
    if err != nil {
        log.Errorf("Error getting next state: %v", err)
        return nil, status.Errorf(codes.Internal, "failed to determine next state: %v", err)
    }
    log.Infof("Next state determined: %s", nextState)

    currentState, err := s.sRepo.GetLatestNodeState(nId.String())
    if err != nil && err != gorm.ErrRecordNotFound {
        log.Errorf("Error retrieving current state: %v", err)
        return nil, status.Errorf(codes.Internal, "failed to retrieve current state: %v", err)
    }

    newNodeState := &db.NodeState{
        Id:           uuid.NewV4(),
        NodeId:       nId.String(),
        CurrentState: nextState,
        SubState:     req.SubState,
        Events:       events,
        Severity:     req.Severity,
    }

    if currentState != nil {
        newNodeState.PreviousStateId = &currentState.Id
    }

    err = s.sRepo.AddNodeState(newNodeState, currentState)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "failed to add node state: %v", err)
    }

    return &pb.AddNodeStateResponse{
		Id:newNodeState.Id.String(),
	}, nil
}
 
func (s *NodeStateServer) GetNodeStates(ctx context.Context, req *pb.GetNodeStatesRequest) (*pb.GetNodeStatesResponse, error) {
    log.Infof("Getting node states for Node ID: %v", req.NodeId)

    if req == nil {
        return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
    }

    nId, err := ukama.ValidateNodeId(req.NodeId)
    if err != nil {
        return nil, status.Errorf(codes.InvalidArgument,
            "invalid format of node id: %s", err.Error())
    }

    history, err := s.sRepo.GetNodeStateHistory(nId.String())
    if err != nil {
        log.Errorf("Failed to get node state history: %v", err)
        return nil, status.Errorf(codes.Internal, "failed to get node state history: %v", err)
    }

    if history == nil {
        log.Warnf("No history found for Node ID: %v", req.NodeId)
        return &pb.GetNodeStatesResponse{NodeStates: []*pb.NodeState{}}, nil
    }

    stateMap := make(map[string]*pb.NodeState)

    for _, nodeState := range history {
        grpcNodeState := &pb.NodeState{
            Id:               nodeState.Id.String(),
            NodeId:           nodeState.NodeId,
            CurrentState:     nodeState.CurrentState,
            SubState:         nodeState.SubState,
            Events:           nodeState.Events,
            Severity:         nodeState.Severity,
            CreatedAt:        timestamppb.New(nodeState.CreatedAt),
            UpdatedAt:        timestamppb.New(nodeState.UpdatedAt),
        }

        if nodeState.PreviousStateId != nil {
            grpcNodeState.PreviousStateId = nodeState.PreviousStateId.String()
        }

        stateMap[grpcNodeState.Id] = grpcNodeState
    }

    for _, state := range stateMap {
        if state.PreviousStateId != "" {
            if prevState, exists := stateMap[state.PreviousStateId]; exists {
                prevStateCopy := &pb.NodeState{
                    Id:               prevState.Id,
                    NodeId:           prevState.NodeId,
                    PreviousStateId:  prevState.PreviousStateId,
                    CurrentState:     prevState.CurrentState,
                    SubState:         prevState.SubState,
                    Events:           prevState.Events,
                    Severity:         prevState.Severity,
                    CreatedAt:        prevState.CreatedAt,
                    UpdatedAt:        prevState.UpdatedAt,
                    DeletedAt:        prevState.DeletedAt,
                }
                state.PreviousState = prevStateCopy
            }
        }
    }

    stateHistory := &pb.GetNodeStatesResponse{
        NodeStates: make([]*pb.NodeState, 0, len(history)),
    }

    for _, state := range stateMap {
        stateHistory.NodeStates = append(stateHistory.NodeStates, state)
    }

    log.Infof("Retrieved %d node states for Node ID: %v", len(stateHistory.NodeStates), req.NodeId)
    return stateHistory, nil
}
func (s *NodeStateServer) GetLatestNodeState(ctx context.Context, req *pb.GetLatestNodeStateRequest) (*pb.GetLatestNodeStateResponse, error) {
    log.Infof("Getting latest node state for Node ID: %v", req.NodeId)


    nId, err := ukama.ValidateNodeId(req.NodeId)
    if err != nil {
        return nil, status.Errorf(codes.InvalidArgument,
            "invalid format of node id: %s", err.Error())
    }

    latestState, err := s.sRepo.GetLatestNodeState(nId.String())
    if err != nil {
        log.Errorf("Failed to get latest node state: %v", err)
        return nil, status.Errorf(codes.Internal, "failed to get latest node state: %v", err)
    }

    if latestState == nil {
        log.Warnf("No state found for Node ID: %v", req.NodeId)
        return nil, status.Error(codes.NotFound, "no state found for the given node ID")
    }

    grpcNodeState := &pb.NodeState{
        Id:           latestState.Id.String(),
        NodeId:       latestState.NodeId,
        CurrentState: latestState.CurrentState,
        SubState:     latestState.SubState,
        Events:       latestState.Events,
        Severity:     latestState.Severity,
        CreatedAt:    timestamppb.New(latestState.CreatedAt),
        UpdatedAt:    timestamppb.New(latestState.UpdatedAt),
    }

    if latestState.PreviousStateId != nil {
        grpcNodeState.PreviousStateId = latestState.PreviousStateId.String()
    }



    return &pb.GetLatestNodeStateResponse{
        NodeState: grpcNodeState,
    }, nil
}