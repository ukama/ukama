package server

import (
	"context"

	pb "github.com/ukama/ukama/systems/subscriber/sim-manager/pb/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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
