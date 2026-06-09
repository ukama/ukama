/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package client

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/ukama/ukama/systems/node/operation-monitor/pb/gen"
)

type OperationMonitor interface {
	Register(req *pb.RegisterIntentRequest) (*pb.RegisterIntentResponse, error)
	Close()
}

type operationMonitor struct {
	conn    *grpc.ClientConn
	client  pb.OperationMonitorServiceClient
	timeout time.Duration
}

func NewOperationMonitor(host string, timeout time.Duration) OperationMonitor {
	conn, err := grpc.NewClient(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to operation monitor at %s: %v", host, err)
	}
	return &operationMonitor{
		conn:    conn,
		client:  pb.NewOperationMonitorServiceClient(conn),
		timeout: timeout,
	}
}

func (m *operationMonitor) Register(req *pb.RegisterIntentRequest) (*pb.RegisterIntentResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()
	return m.client.RegisterIntent(ctx, req)
}

func (m *operationMonitor) Close() {
	if m.conn != nil {
		_ = m.conn.Close()
	}
}
