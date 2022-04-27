package grpc

import (
	"fmt"
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/services/common/config"
	"google.golang.org/grpc"
	pbhealth "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

// Basic GrpcServer with the set of middlewares
type UkamaGrpcServer struct {
	// replace with custom implementation if needed
	healthChecker           pbhealth.HealthServer
	config                  config.Grpc
	serviceRegistrar        func(s *grpc.Server)
	ExtraUnaryInterceptors  []grpc.UnaryServerInterceptor
	ExtraStreamInterceptors []grpc.StreamServerInterceptor
}

func NewGrpcServer(config config.Grpc, serviceRegistrar func(s *grpc.Server)) *UkamaGrpcServer {
	return &UkamaGrpcServer{healthChecker: NewDefaultHealthChecker(), config: config,
		serviceRegistrar: serviceRegistrar}
}

func NewGrpcServerWithCustomHealthcheck(healthChecker *HealthChecker, config config.Grpc, serviceRegistrator func(s *grpc.Server)) *UkamaGrpcServer {
	return &UkamaGrpcServer{healthChecker: healthChecker, config: config, serviceRegistrar: serviceRegistrator}
}

func (g *UkamaGrpcServer) StartServer() {
	log.Infof("Starting gRpc on port " + fmt.Sprintf(":%d", g.config.Port))
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", g.config.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	g.startServerInternal(lis)
}

func (g *UkamaGrpcServer) startServerInternal(listener net.Listener) {
	logrusEntry := log.NewEntry(log.New())

	sInterc := []grpc.StreamServerInterceptor{
		grpc_logrus.StreamServerInterceptor(logrusEntry),
		grpc_prometheus.StreamServerInterceptor,
		grpc_validator.StreamServerInterceptor(),
	}
	sInterc = append(sInterc, g.ExtraStreamInterceptors...)

	uInterc := []grpc.UnaryServerInterceptor{
		grpc_logrus.UnaryServerInterceptor(logrusEntry),
		grpc_prometheus.UnaryServerInterceptor,
		grpc_validator.UnaryServerInterceptor(),
	}
	uInterc = append(uInterc, g.ExtraUnaryInterceptors...)

	s := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(sInterc...)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(uInterc...)),
	)

	g.serviceRegistrar(s)

	pbhealth.RegisterHealthServer(s, g.healthChecker)
	reflection.Register(s)
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
