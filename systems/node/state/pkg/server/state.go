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
 
 type StateServer struct {
	 pb.UnimplementedStateServiceServer
	 sRepo               db.StateRepo
	 nodeStateRoutingKey msgbus.RoutingKeyBuilder
	 msgbus              mb.MsgBusServiceClient
	 debug               bool
	 orgName             string
 }
 
 func NewStateServer(orgName string, sRepo db.StateRepo, msgBus mb.MsgBusServiceClient, debug bool) *StateServer {
	 ns := &StateServer{
		 sRepo:   sRepo,
		 orgName: orgName,
		 msgbus:  msgBus,
		 debug:   debug,
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
 
	 
	 currentState, err := s.sRepo.GetLatestState(nId.String())
	 if err != nil && err != gorm.ErrRecordNotFound {
		 log.Errorf("Error retrieving current state: %v", err)
		 return nil, status.Errorf(codes.Internal, "failed to retrieve current state: %v", err)
	 }
 
	 newNodeState := &db.State{
		 Id:           uuid.NewV4(),
		 NodeId:       nId.String(),
		 CurrentState: req.CurrentState,
		 SubState:     req.SubState,
		 Events:       events,
	 }
 
	 if currentState != nil {
		 newNodeState.PreviousStateId = &currentState.Id
	 }
 
	 err = s.sRepo.AddState(newNodeState, currentState)
	 if err != nil {
		 return nil, status.Errorf(codes.Internal, "failed to add node state: %v", err)
	 }
	 if s.msgbus != nil {
		route := s.nodeStateRoutingKey.SetAction("add").SetObject("network").MustBuild()
		evt := &pb.AddStateRequest {
			NodeId: newNodeState.NodeId,
			CurrentState: newNodeState.CurrentState,
			SubState: newNodeState.SubState,
			Events: newNodeState.Events,
		}

		err = s.msgbus.PublishRequest(route, evt)
		if err != nil {
			log.Errorf("Failed to publish message %+v with key %+v. Errors %s",
				evt, route, err.Error())
		}
	}
 
	 return &pb.AddStateResponse{
		 Id: newNodeState.Id.String(),
	 }, nil
 }
 
 func (s *StateServer) GetStates(ctx context.Context, req *pb.GetStatesRequest) (*pb.GetStatesResponse, error) {
	 log.Infof("Getting node states for Node ID: %v", req.NodeId)
 
	 nId, err := ukama.ValidateNodeId(req.NodeId)
	 if err != nil {
		 return nil, status.Errorf(codes.InvalidArgument,
			 "invalid format of node id: %s", err.Error())
	 }
	 history, err := s.sRepo.GetStateHistory(nId.String())
	 if err != nil {
		 log.Errorf("Failed to get node state history: %v", err)
		 return nil, status.Errorf(codes.Internal, "failed to get node state history: %v", err)
	 }
 
	 if history == nil {
		 log.Warnf("No history found for Node ID: %v", req.NodeId)
		 return &pb.GetStatesResponse{States: []*pb.State{}}, nil
	 }
 
	 stateMap := make(map[string]*pb.State)
 
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
	 }
 
	 for _, state := range stateMap {
		 if state.PreviousStateId != "" {
			 if prevState, exists := stateMap[state.PreviousStateId]; exists {
				 prevStateCopy := &pb.State{
					 Id:              prevState.Id,
					 NodeId:          prevState.NodeId,
					 PreviousStateId: prevState.PreviousStateId,
					 CurrentState:    prevState.CurrentState,
					 SubState:        prevState.SubState,
					 Events:          prevState.Events,
					 Severity:        prevState.Severity,
					 CreatedAt:       prevState.CreatedAt,
					 UpdatedAt:       prevState.UpdatedAt,
					 DeletedAt:       prevState.DeletedAt,
				 }
				 state.PreviousState = prevStateCopy
			 }
		 }
	 }
 
	 stateHistory := &pb.GetStatesResponse{
		 States: make([]*pb.State, 0, len(history)),
	 }
 
	 for _, state := range stateMap {
		 stateHistory.States = append(stateHistory.States, state)
	 }
 
	 return stateHistory, nil
 }
 func (s *StateServer) GetLatestState(ctx context.Context, req *pb.GetLatestStateRequest) (*pb.GetLatestStateResponse, error) {
	 log.Infof("Getting latest node state for Node ID: %v", req.NodeId)
 
	 nId, err := ukama.ValidateNodeId(req.NodeId)
	 if err != nil {
		 return nil, status.Errorf(codes.InvalidArgument,
			 "invalid format of node id: %s", err.Error())
	 }
 
	 latestState, err := s.sRepo.GetLatestState(nId.String())
	 if err != nil {
		 log.Errorf("Failed to get latest node state: %v", err)
		 return nil, status.Errorf(codes.Internal, "failed to get latest node state: %v", err)
	 }
 
	 StateRes := &pb.State{
		 Id:           latestState.Id.String(),
		 NodeId:       latestState.NodeId,
		 CurrentState: latestState.CurrentState,
		 SubState:     latestState.SubState,
		 Events:       latestState.Events,
		 CreatedAt:    timestamppb.New(latestState.CreatedAt),
		 UpdatedAt:    timestamppb.New(latestState.UpdatedAt),
	 }
 
	 if latestState.PreviousStateId != nil {
		 StateRes.PreviousStateId = latestState.PreviousStateId.String()
	 }
 
	 return &pb.GetLatestStateResponse{
		 State: StateRes,
	 }, nil
 }
 