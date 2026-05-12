/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package client

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/node/health/pb/gen"
)

type Health struct {
	conn    *grpc.ClientConn
	client  pb.HealthServiceClient
	timeout time.Duration
	host    string
}

func NewHealth(healthHost string, timeout time.Duration) *Health {
	conn, err := grpc.NewClient(healthHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to Health Service host: %v", err)
	}
	client := pb.NewHealthServiceClient(conn)

	return &Health{
		conn:    conn,
		client:  client,
		timeout: timeout,
		host:    healthHost,
	}
}

func NewHealthFromClient(mClient pb.HealthServiceClient) *Health {
	return &Health{
		host:    "localhost",
		timeout: 1 * time.Second,
		conn:    nil,
		client:  mClient,
	}
}

func (h *Health) Close() {
	if h.conn != nil {
		if err := h.conn.Close(); err != nil {
			log.Warnf("Failed to gracefully close Health service connection: %v", err)
		}
	}
}

func (h *Health) StoreHealthReport(request *pb.StoreHealthReportRequest) (*pb.StoreHealthReportResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), h.timeout)
	defer cancel()

	return h.client.StoreHealthReport(ctx, request)
}

func (h *Health) List(request *pb.ListRequest) (*pb.ListResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), h.timeout)
	defer cancel()

	return h.client.List(ctx, request)
}
