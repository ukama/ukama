package server

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/ukama/ukama/systems/common/grpc"
	pb "github.com/ukama/ukama/systems/subscriber/sim-manager/pb/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ukama/ukama/systems/subscriber/sim-manager/pkg/clients"
	"github.com/ukama/ukama/systems/subscriber/sim-manager/pkg/db"
)

type SimManagerServer struct {
	pb.UnimplementedSimManagerServiceServer
	simRepo      db.SimRepo
	agentFactory *clients.AgentFactory
}

func NewSimManagerServer(simRepo db.SimRepo, agentFactory *clients.AgentFactory) *SimManagerServer {
	return &SimManagerServer{
		simRepo:      simRepo,
		agentFactory: agentFactory,
	}
}

func (s *SimManagerServer) GetBySubscriber(ctx context.Context, req *pb.GetBySubscriberRequest) (*pb.GetBySubscriberResponse, error) {
	subID, err := uuid.Parse(req.GetSubscriberID())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid format of subscriber uuid. Error %s", err.Error())
	}

	sims, err := s.simRepo.GetBySubscriber(subID)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "sims")
	}

	resp := &pb.GetBySubscriberResponse{
		SubscriberID: req.GetSubscriberID(),
		Sims:         dbSimsToPbSims(sims),
	}

	return resp, nil
}

func (s *SimManagerServer) ActivateSim(ctx context.Context, req *pb.ActivateSimRequest) (*pb.ActivateSimResponse, error) {
	simAgent, ok := s.agentFactory.GetAgentAdapter(req.SimType)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "invalid sim type:%s", req.SimID)
	}

	err := simAgent.ActivateSim(ctx, req.SimID)
	if err != nil {
		return nil, err
	}

	return &pb.ActivateSimResponse{}, nil
}

func (s *SimManagerServer) DeactivateSim(ctx context.Context, req *pb.DeactivateSimRequest) (*pb.DeactivateSimResponse, error) {
	simAgent, ok := s.agentFactory.GetAgentAdapter(req.SimType)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "invalid sim type:%s", req.SimID)
	}

	err := simAgent.DeactivateSim(ctx, req.SimID)
	if err != nil {
		return nil, err
	}

	return &pb.DeactivateSimResponse{}, nil
}

func dbSimToPbSim(sim *db.Sim) *pb.Sim {
	return &pb.Sim{
		Id:           sim.ID.String(),
		SubscriberID: sim.SubscriberID.String(),
		Iccid:        sim.Iccid,
		Msisdn:       sim.Msisdn,
		IsPhysical:   sim.IsPhysical,
		AllocatedAt:  timestamppb.New(time.Unix(sim.AllocatedAt, 0)),
	}
}

func dbSimsToPbSims(sims []db.Sim) []*pb.Sim {
	res := []*pb.Sim{}

	for _, s := range sims {
		res = append(res, dbSimToPbSim(&s))
	}

	return res
}
