package server

import (
	"context"

	pb "github.com/ukama/ukama/systems/subscriber/test-agent/pb/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TestAgentServer struct {
	pb.UnimplementedTestAgentServiceServer
}

func NewTestAgentServer() *TestAgentServer {
	return &TestAgentServer{}
}

func (s *TestAgentServer) ActivateSim(ctx context.Context, req *pb.ActivateSimRequest) (*pb.ActivateSimResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "cannot activate sim %s: method TestAgent.ActivateSim not implemented", req.Iccid)
}
func (s *TestAgentServer) DeactivateSim(ctx context.Context, req *pb.DeactivateSimRequest) (*pb.DeactivateSimResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "cannot deactivate sim %s: method TestAgent.DeactivateSim not implemented", req.Iccid)
}
