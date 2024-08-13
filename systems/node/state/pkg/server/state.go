package server

import (
	"context"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ukama/ukama/systems/common/grpc"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	pb "github.com/ukama/ukama/systems/node/state/pb/gen"
	"github.com/ukama/ukama/systems/node/state/pkg"
	"github.com/ukama/ukama/systems/node/state/pkg/db"
)

type StateServer struct {
    pb.UnimplementedStateServiceServer
    sRepo            db.StateRepo
    stateRoutingKey  msgbus.RoutingKeyBuilder
    msgbus           mb.MsgBusServiceClient
    debug            bool
    orgName          string
}

func NewstateServer(orgName string, sRepo db.StateRepo, msgBus mb.MsgBusServiceClient, debug bool) *StateServer {
    return &StateServer{
        sRepo:            sRepo,
        orgName:          orgName,
        stateRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
        msgbus:           msgBus,
        debug:            debug,
    }
}

func (s *StateServer) Create(ctx context.Context, req *pb.CreateStateRequest) (*pb.State, error) {
	log.Infof("Adding node state  %v", req)

	nId, err := ukama.ValidateNodeId(req.State.NodeId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of node id. Error %s", err.Error())
	}

    state := &db.State{
        NodeId:        nId.StringLowercase(),
        CurrentState:  db.NodeStateEnum(req.State.CurrentState),
        Connectivity:  db.Connectivity(req.State.Connectivity),
        LastHeartbeat: req.State.LastHeartbeat.AsTime(),
        Type:          req.State.Type,
        Version:       req.State.Version,
    }

    err = s.sRepo.Create(state, nil)
	if err != nil {
		log.Error("Failed to create state: " + err.Error())

		return nil, grpc.SqlErrorToGrpc(err, "node")
	}

    return convertStateToProto(state), nil
}

func (s *StateServer) GetByNodeId(ctx context.Context, req *pb.GetByNodeIdRequest) (*pb.State, error) {
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


    return convertStateToProto(state), nil
}

func (s *StateServer) Update(ctx context.Context, req *pb.UpdateStateRequest) (*pb.State, error) {
	nId, err := ukama.ValidateNodeId(req.State.NodeId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of node id. Error %s", err.Error())
	}
	

    state := &db.State{
        Id:            uuid.NewV4(),
        NodeId:        nId.StringLowercase(),
        CurrentState:  db.NodeStateEnum(req.State.CurrentState),
        Connectivity:  db.Connectivity(req.State.Connectivity),
        LastHeartbeat: req.State.LastHeartbeat.AsTime(),
        Type:          req.State.Type,
        Version:       req.State.Version,
    }

    err = s.sRepo.Update(state)
   
	if err != nil {
		log.Error("Failed to update state: " + err.Error())

		return nil, grpc.SqlErrorToGrpc(err, "state")
	}

    return convertStateToProto(state), nil
}

func (s *StateServer) Delete(ctx context.Context, req *pb.DeleteStateRequest) (*emptypb.Empty, error) {
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

    return &emptypb.Empty{}, nil
}

func (s *StateServer) ListAll(ctx context.Context,req *pb.ListAllRequest) (*pb.ListAllResponse, error) {
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

func (s *StateServer) UpdateConnectivity(ctx context.Context, req *pb.UpdateConnectivityRequest) (*emptypb.Empty, error) {
	nId, err := ukama.ValidateNodeId(req.NodeId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of node id. Error %s", err.Error())
	}
    err = s.sRepo.UpdateConnectivity(nId, db.Connectivity(req.Connectivity))
    if err != nil {
        return nil, status.Errorf(codes.Internal, "Failed to update connectivity: %v", err)
    }

    return &emptypb.Empty{}, nil
}

func (s *StateServer) UpdateCurrentState(ctx context.Context, req *pb.UpdateCurrentStateRequest) (*emptypb.Empty, error) {
    err := s.sRepo.UpdateCurrentState(ukama.NodeID(req.NodeId), db.NodeStateEnum(req.CurrentState))
    if err != nil {
        return nil, status.Errorf(codes.Internal, "Failed to update current state: %v", err)
    }

    return &emptypb.Empty{}, nil
}

func (s *StateServer) GetStateHistory(ctx context.Context, req *pb.GetStateHistoryRequest) (*pb.GetStateHistoryResponse, error) {
    history, err := s.sRepo.GetStateHistory(ukama.NodeID(req.NodeId))
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

    return &pb.GetStateHistoryResponse{History: pbHistory}, nil
}

func convertStateToProto(state *db.State) *pb.State {
    return &pb.State{
        Id:            state.Id.String(),
        NodeId:        state.NodeId,
        CurrentState:  pb.NodeStateEnum(state.CurrentState),
        Connectivity:  pb.Connectivity(state.Connectivity),
        LastHeartbeat: timestamppb.New(state.LastHeartbeat),
        Type:          state.Type,
        Version:       state.Version,
    }
}