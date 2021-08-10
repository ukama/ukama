package server

import (
	"context"
	"github.com/sirupsen/logrus"
	pb "ukamaX/registry/pb/generated"
)

type HealthChecker struct{}

func NewHealthChecker() *HealthChecker {
	return &HealthChecker{}
}

func (s *HealthChecker) Check(ctx context.Context, request *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	logrus.Info("Serving the Check request for health check")

	return &pb.HealthCheckResponse{
		Status: pb.HealthCheckResponse_SERVING,
	}, nil
}

func (s *HealthChecker) Watch(request *pb.HealthCheckRequest, server pb.Health_WatchServer) error {
	panic("Watch method is not supported")
}
