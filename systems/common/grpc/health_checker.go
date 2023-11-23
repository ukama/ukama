/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

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
