package server

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	uuid "github.com/ukama/ukama/systems/common/uuid"
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
		Id:              uuid.NewV4(),
		NodeId:          nId.StringLowercase(),
		State:           db.NodeStateEnum(req.State.State),
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

func (s *StateServer) GetStateHistory(ctx context.Context, req *pb.GetStateHistoryRequest) (*pb.GetStateHistoryResponse, error) {
	log.Infof("Getting state history for node ID %v", req.NodeId)

	nId, err := ukama.ValidateNodeId(req.NodeId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of node id. Error %s", err.Error())
	}

	states, err := s.sRepo.GetStateHistory(ukama.NodeID(nId))
	if err != nil {
		log.Error("Failed to get state history: " + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "node")
	}

	pbStates := make([]*pb.State, len(states))
	for i, state := range states {
		pbStates[i] = convertStateToProto(&state)
	}

	return &pb.GetStateHistoryResponse{StateHistory: pbStates}, nil
}

func convertStateToProto(state *db.State) *pb.State {
	return &pb.State{
		Id:              state.Id.String(),
		NodeId:          state.NodeId,
		State:           pb.NodeStateEnum(state.State),
		LastHeartbeat:   timestamppb.New(state.LastHeartbeat),
		LastStateChange: timestamppb.New(state.LastStateChange),
		Type:            state.Type,
		Version:         state.Version,
	}
}
