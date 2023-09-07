package grpc

import (
	"fmt"
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/config"
	pbhealth "github.com/ukama/ukama/systems/common/pb/gen/health"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Basic GrpcServer with the set of middlewares
type UkamaGrpcServer struct {
	// replace with custom implementation if needed
	server                  *grpc.Server
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

	server := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(sInterc...)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(uInterc...)),
	)

	g.serviceRegistrar(server)

	pbhealth.RegisterHealthServer(server, g.healthChecker)
	reflection.Register(server)
	g.server = server
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (g *UkamaGrpcServer) StopServer() {
	log.Infof("Stoping gRpc server.")

	g.server.Stop()
}
