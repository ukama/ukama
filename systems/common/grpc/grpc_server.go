/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

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
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

// Basic GrpcServer with the set of middlewares
type UkamaGrpcServer struct {
	server                  *grpc.Server
	config                  config.Grpc
	serviceRegistrar        func(s *grpc.Server)
	ExtraUnaryInterceptors  []grpc.UnaryServerInterceptor
	ExtraStreamInterceptors []grpc.StreamServerInterceptor

	// dependencyState carries opt-in dependency checks registered via
	// RegisterDependency and the background monitor driving them.
	dependencyState

	// GrpcHealth is the standard grpc.health.v1 health service
	// (google.golang.org/grpc/health). It is registered on every server so
	// that grpc_health_probe and Kubernetes native gRPC probes work out of
	// the box. Services can flip their status (e.g. when a dependency goes
	// down) via GrpcHealth.SetServingStatus.
	GrpcHealth *health.Server
}

func NewGrpcServer(config config.Grpc, serviceRegistrar func(s *grpc.Server)) *UkamaGrpcServer {
	return &UkamaGrpcServer{config: config,
		serviceRegistrar: serviceRegistrar, GrpcHealth: health.NewServer()}
}

func (g *UkamaGrpcServer) StartServer() {
	log.Infof("Starting gRpc on port %s", fmt.Sprintf(":%d", g.config.Port))
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", g.config.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	g.startServerInternal(lis)
}

func (g *UkamaGrpcServer) startServerInternal(listener net.Listener) {
	logrusEntry := log.NewEntry(log.New())
	if g.config.MaxMsgSize == 0 {
		g.config.MaxMsgSize = 4194304
	}
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
		grpc.MaxRecvMsgSize(g.config.MaxMsgSize),
		grpc.MaxSendMsgSize(g.config.MaxMsgSize),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(sInterc...)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(uInterc...)),
	)

	g.serviceRegistrar(server)

	// Standard health service (grpc.health.v1), used by Kubernetes native
	// gRPC probes and grpc_health_probe.
	if g.GrpcHealth == nil {
		g.GrpcHealth = health.NewServer()
	}
	healthpb.RegisterHealthServer(server, g.GrpcHealth)
	g.GrpcHealth.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	// Liveness target: SERVING while the process is up. Liveness probes
	// (grpc service: "live") do not react to short dependency failures —
	// only the default "" status above does. The dependency monitor
	// escalates "live" to NOT_SERVING when a critical dependency stays
	// down past DependencyRestartThreshold, so kubelet restarts the
	// container (see dependency.go).
	g.GrpcHealth.SetServingStatus(LivenessService, healthpb.HealthCheckResponse_SERVING)

	// Monitors registered dependency checks (no-op when none registered)
	// and keeps the "" status in sync with critical dependencies.
	g.startDependencyMonitor()

	reflection.Register(server)
	g.server = server
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (g *UkamaGrpcServer) StopServer() {
	log.Infof("Stoping gRpc server.")

	g.stopDependencyMonitor()

	if g.GrpcHealth != nil {
		// Flip all services to NOT_SERVING so readiness probes fail fast
		// while the server drains.
		g.GrpcHealth.Shutdown()
	}

	g.server.Stop()
}
