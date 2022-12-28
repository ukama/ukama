package server

import (
	"context"

	pb "github.com/ukama/ukama/systems/subscriber/sim-manager/pb/gen"

	"github.com/ukama/ukama/systems/subscriber/sim-manager/pkg/clients"
	"github.com/ukama/ukama/systems/subscriber/sim-manager/pkg/db"
)

type SimManagerServer struct {
	pb.UnimplementedSimManagerServiceServer
	simRepo          db.SimRepo
	testAgentAdapter *clients.TestAgentAdapter
}

func NewSimManagerServer(simRepo db.SimRepo, testAgentAdapter *clients.TestAgentAdapter) *SimManagerServer {
	return &SimManagerServer{
		simRepo:          simRepo,
		testAgentAdapter: testAgentAdapter,
	}
}

func (s *SimManagerServer) ActivateSim(ctx context.Context, req *pb.ActivateSimRequest) (*pb.ActivateSimResponse, error) {
	err := s.testAgentAdapter.ActivateSim(ctx, req.SimID)
	if err != nil {
		return nil, err
	}

	return &pb.ActivateSimResponse{}, nil
}

func (s *SimManagerServer) DeactivateSim(ctx context.Context, req *pb.DeactivateSimRequest) (*pb.DeactivateSimResponse, error) {
	err := s.testAgentAdapter.DeactivateSim(ctx, req.SimID)
	if err != nil {
		return nil, err
	}

	return &pb.DeactivateSimResponse{}, nil
}
