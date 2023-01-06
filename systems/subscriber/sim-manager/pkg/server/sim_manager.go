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

	sims "github.com/ukama/ukama/systems/subscriber/sim-manager/pkg/db"
)

type SimManagerServer struct {
	pb.UnimplementedSimManagerServiceServer
	simRepo      sims.SimRepo
	agentFactory *clients.AgentFactory
}

func NewSimManagerServer(simRepo sims.SimRepo, agentFactory *clients.AgentFactory) *SimManagerServer {
	return &SimManagerServer{
		simRepo:      simRepo,
		agentFactory: agentFactory,
	}
}

func (s *SimManagerServer) GetSimsBySubscriber(ctx context.Context, req *pb.GetSimsBySubscriberRequest) (*pb.GetSimsBySubscriberResponse, error) {
	subID, err := uuid.Parse(req.GetSubscriberID())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid format of subscriber uuid. Error %s", err.Error())
	}

	sims, err := s.simRepo.GetBySubscriber(subID)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "sims")
	}

	resp := &pb.GetSimsBySubscriberResponse{
		SubscriberID: req.GetSubscriberID(),
		Sims:         dbSimsToPbSims(sims),
	}

	return resp, nil
}

func (s *SimManagerServer) ActivateSim(ctx context.Context, req *pb.ActivateSimRequest) (*pb.ActivateSimResponse, error) {
	simID, err := uuid.Parse(req.GetSimID())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid format of sim uuid. Error %s", err.Error())
	}

	sim, err := s.simRepo.Get(simID)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "sim")
	}

	simAgent, ok := s.agentFactory.GetAgentAdapter(sim.Type)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "invalid sim type %q for sim ID %q", sim.Type, req.SimID)
	}

	err = simAgent.ActivateSim(ctx, req.SimID)
	if err != nil {
		return nil, err
	}

	return &pb.ActivateSimResponse{}, nil
}

func (s *SimManagerServer) DeactivateSim(ctx context.Context, req *pb.DeactivateSimRequest) (*pb.DeactivateSimResponse, error) {
	simID, err := uuid.Parse(req.GetSimID())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid format of sim uuid. Error %s", err.Error())
	}

	sim, err := s.simRepo.Get(simID)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "sim")
	}

	simAgent, ok := s.agentFactory.GetAgentAdapter(sim.Type)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "invalid sim type:%q for sim ID %q", sim.Type, req.SimID)
	}

	err = simAgent.DeactivateSim(ctx, req.SimID)
	if err != nil {
		return nil, err
	}

	return &pb.DeactivateSimResponse{}, nil
}

func dbSimToPbSim(sim *sims.Sim) *pb.Sim {
	return &pb.Sim{
		Id:           sim.ID.String(),
		SubscriberID: sim.SubscriberID.String(),
		Iccid:        sim.Iccid,
		Msisdn:       sim.Msisdn,
		Type:         sim.Type.String(),
		Status:       sim.Status.String(),
		IsPhysical:   sim.IsPhysical,
		AllocatedAt:  timestamppb.New(time.Unix(sim.AllocatedAt, 0)),
	}
}

func dbSimsToPbSims(sims []sims.Sim) []*pb.Sim {
	res := []*pb.Sim{}

	for _, s := range sims {
		res = append(res, dbSimToPbSim(&s))
	}

	return res
}
