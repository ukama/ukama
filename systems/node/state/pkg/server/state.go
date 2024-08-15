package server

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ukama/ukama/systems/common/grpc"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/ukama"
	pb "github.com/ukama/ukama/systems/node/state/pb/gen"
	"github.com/ukama/ukama/systems/node/state/pkg"
	"github.com/ukama/ukama/systems/node/state/pkg/db"
)

const (
	FaultyThresholdDuration = 5 * time.Minute // Adjust as needed
)

type StateServer struct {
	pb.UnimplementedStateServiceServer
	sRepo           db.StateRepo
	stateRoutingKey msgbus.RoutingKeyBuilder
	msgbus          mb.MsgBusServiceClient
	debug           bool
	orgName         string
}

func NewstateServer(orgName string, sRepo db.StateRepo, msgBus mb.MsgBusServiceClient, debug bool) *StateServer {
	return &StateServer{
		sRepo:           sRepo,
		orgName:         orgName,
		stateRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
		msgbus:          msgBus,
		debug:           debug,
	}
}

func (s *StateServer) Create(ctx context.Context, req *pb.CreateStateRequest) (*pb.CreateStateResponse, error) {
	log.Infof("Adding node state  %v", req)

	nId, err := ukama.ValidateNodeId(req.State.NodeId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of node id. Error %s", err.Error())
	}

	now := time.Now()
	state := &db.State{
		NodeId:          nId.StringLowercase(),
		CurrentState:    db.NodeStateEnum(req.State.CurrentState),
		Connectivity:    db.Connectivity(req.State.Connectivity),
		LastHeartbeat:   now,
		LastStateChange: now,
		Type:            req.State.Type,
		Version:         req.State.Version,
	}

	err = s.sRepo.Create(state, nil)
	if err != nil {
		log.Error("Failed to create state: " + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "node")
	}
	return &pb.CreateStateResponse{
		State: convertStateToProto(state),
	}, nil
}

func (s *StateServer) GetByNodeId(ctx context.Context, req *pb.GetByNodeIdRequest) (*pb.GetByNodeIdResponse, error) {
	log.Infof("Getting node state  %v", req)

	nId, err := ukama.ValidateNodeId(req.NodeId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of node id. Error %s", err.Error())
	}

	state, err := s.sRepo.GetByNodeId(nId)
	if err != nil {
		log.Error("State not found: " + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "node")
	}

	return &pb.GetByNodeIdResponse{State: convertStateToProto(state)}, nil
}

func (s *StateServer) Update(ctx context.Context, req *pb.UpdateStateRequest) (*pb.UpdateStateResponse, error) {
	nId, err := ukama.ValidateNodeId(req.State.NodeId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of node id. Error %s", err.Error())
	}

	existingState, err := s.sRepo.GetByNodeId(nId)
	if err != nil {
		log.Error("Failed to get existing state: " + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "state")
	}

	now := time.Now()
	state := &db.State{
		Id:              existingState.Id,
		NodeId:          nId.StringLowercase(),
		CurrentState:    db.NodeStateEnum(req.State.CurrentState),
		Connectivity:    db.Connectivity(req.State.Connectivity),
		LastHeartbeat:   now,
		LastStateChange: existingState.LastStateChange,
		Type:            req.State.Type,
		Version:         req.State.Version,
	}

	if state.CurrentState != existingState.CurrentState {
		state.LastStateChange = now
	}

	err = s.sRepo.Update(state)
	if err != nil {
		log.Error("Failed to update state: " + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "state")
	}

	return &pb.UpdateStateResponse{State: convertStateToProto(state)}, nil
}

func (s *StateServer) Delete(ctx context.Context, req *pb.DeleteStateRequest) (*pb.DeleteStateResponse, error) {
	nId, err := ukama.ValidateNodeId(req.NodeId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of node id. Error %s", err.Error())
	}

	err = s.sRepo.Delete(ukama.NodeID(nId))
	if err != nil {
		log.Error("Failed to delete state: " + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "state")
	}

	return &pb.DeleteStateResponse{}, nil
}

func (s *StateServer) ListAll(ctx context.Context, req *pb.ListAllRequest) (*pb.ListAllResponse, error) {
	states, err := s.sRepo.ListAll()
	if err != nil {
		log.Error("Failed to list states: " + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "state")
	}

	pbStates := make([]*pb.State, len(states))
	for i, state := range states {
		pbStates[i] = convertStateToProto(&state)
	}

	return &pb.ListAllResponse{States: pbStates}, nil
}

func (s *StateServer) UpdateConnectivity(ctx context.Context, req *pb.UpdateConnectivityRequest) (*pb.UpdateConnectivityResponse, error) {
	nId, err := ukama.ValidateNodeId(req.NodeId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of node id. Error %s", err.Error())
	}

	state, err := s.sRepo.GetByNodeId(nId)
	if err != nil {
		log.Error("Failed to get existing state: " + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "state")
	}

	newConnectivity := db.Connectivity(req.Connectivity)
	now := time.Now()

	if state.Connectivity != newConnectivity {
		if state.CurrentState == db.StateFaulty {
			if now.Sub(state.LastStateChange) > FaultyThresholdDuration {
				state.CurrentState = db.StateActive
				state.LastStateChange = now
			}
		} else if state.CurrentState == db.StateActive {
			state.CurrentState = db.StateFaulty
			state.LastStateChange = now
		}
	}

	state.Connectivity = newConnectivity
	state.LastHeartbeat = now

	err = s.sRepo.Update(state)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to update connectivity: %v", err)
	}

	return &pb.UpdateConnectivityResponse{}, nil
}

func (s *StateServer) UpdateCurrentState(ctx context.Context, req *pb.UpdateCurrentStateRequest) (*pb.UpdateCurrentStateResponse, error) {
	nId, err := ukama.ValidateNodeId(req.NodeId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of node id. Error %s", err.Error())
	}

	state, err := s.sRepo.GetByNodeId(nId)
	if err != nil {
		log.Error("Failed to get existing state: " + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "state")
	}

	newState := db.NodeStateEnum(req.CurrentState)
	now := time.Now()

	if state.CurrentState != newState {
		state.CurrentState = newState
		state.LastStateChange = now
	}

	state.LastHeartbeat = now

	err = s.sRepo.Update(state)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to update current state: %v", err)
	}

	return &pb.UpdateCurrentStateResponse{}, nil
}

func (s *StateServer) GetStateHistoryByTimeRange(ctx context.Context, req *pb.GetStateHistoryByTimeRangeRequest) (*pb.GetStateHistoryByTimeRangeResponse, error) {
	nId, err := ukama.ValidateNodeId(req.NodeId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of node id. Error %s", err.Error())
	}

	fromTime := req.From.AsTime()
	toTime := req.To.AsTime()

	history, err := s.sRepo.GetStateHistoryByTimeRange(nId, fromTime, toTime)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get state history: %v", err)
	}

	pbHistory := make([]*pb.StateHistory, len(history))
	for i, h := range history {
		pbHistory[i] = &pb.StateHistory{
			Id:            h.Id.String(),
			NodeStateId:   h.NodeStateId,
			PreviousState: pb.NodeStateEnum(h.PreviousState),
			NewState:      pb.NodeStateEnum(h.NewState),
			TimeStamp:     timestamppb.New(h.Timestamp),
		}
	}

	return &pb.GetStateHistoryByTimeRangeResponse{History: pbHistory}, nil
}

func convertStateToProto(state *db.State) *pb.State {
	return &pb.State{
		Id:              state.Id.String(),
		NodeId:          state.NodeId,
		CurrentState:    pb.NodeStateEnum(state.CurrentState),
		Connectivity:    pb.Connectivity(state.Connectivity),
		LastHeartbeat:   timestamppb.New(state.LastHeartbeat),
		LastStateChange: timestamppb.New(state.LastStateChange),
		Type:            state.Type,
		Version:         state.Version,
	}
}
