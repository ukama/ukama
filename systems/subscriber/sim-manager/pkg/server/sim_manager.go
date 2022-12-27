package server

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/ukama/ukama/systems/subscriber/sim-manager/pb/gen"

	testagentpb "github.com/ukama/ukama/systems/subscriber/test-agent/pb/gen"

	"github.com/ukama/ukama/systems/subscriber/sim-manager/pkg/db"
	"github.com/ukama/ukama/systems/subscriber/sim-manager/pkg/providers"
)

type SimManagerServer struct {
	pb.UnimplementedSimManagerServiceServer
	simRepo          db.SimRepo
	testAgentService providers.TestAgentClientProvider
}

func NewSimManagerServer(simRepo db.SimRepo, testAgentService providers.TestAgentClientProvider) *SimManagerServer {
	return &SimManagerServer{
		simRepo:          simRepo,
		testAgentService: testAgentService,
	}
}

func (s *SimManagerServer) ActivateSim(ctx context.Context, req *pb.ActivateSimRequest) (*pb.ActivateSimResponse, error) {
	svc, err := s.testAgentService.GetClient()
	if err != nil {
		return nil, err
	}

	_, err = svc.ActivateSim(ctx, &testagentpb.ActivateSimRequest{SimID: uuid.NewString()})
	if err != nil {
		return nil, err
	}

	return &pb.ActivateSimResponse{}, nil
}

func (s *SimManagerServer) DeactivateSim(ctx context.Context, req *pb.DeactivateSimRequest) (*pb.DeactivateSimResponse, error) {
	svc, err := s.testAgentService.GetClient()
	if err != nil {
		return nil, err
	}

	_, err = svc.DeactivateSim(ctx, &testagentpb.DeactivateSimRequest{SimID: uuid.NewString()})
	if err != nil {
		return nil, err
	}

	return &pb.DeactivateSimResponse{}, nil
}
