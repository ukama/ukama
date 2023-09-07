package grpc

import (
	"context"

	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/common/pb/gen/health"
)

type HealthCheckerInterface interface {
	Check(ctx context.Context, request *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error)
	Watch(request *pb.HealthCheckRequest, server pb.Health_WatchServer) error
}

type HealthChecker struct {
	pb.UnimplementedHealthServer
}

func NewDefaultHealthChecker() pb.HealthServer {
	return &HealthChecker{}
}

func (s *HealthChecker) Check(ctx context.Context, request *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	log.Info("Serving the Check request for health check")

	return &pb.HealthCheckResponse{
		Status: pb.HealthCheckResponse_SERVING,
	}, nil
}

func (s *HealthChecker) Watch(request *pb.HealthCheckRequest, server pb.Health_WatchServer) error {
	panic("Watch method is not supported")
}
