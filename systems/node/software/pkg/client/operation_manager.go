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

	pb "github.com/ukama/ukama/systems/operation/manager/pb/gen"
)

type OperationManager interface {
	Start(req *pb.StartOperationRequest) (*pb.StartOperationResponse, error)
	MarkRunning(id string, fencingToken uint64) (*pb.MarkRunningResponse, error)
	FailOperation(id, actor, reason string) (*pb.ForceUnlockResponse, error)
	Close()
}

type operationManager struct {
	conn    *grpc.ClientConn
	client  pb.OperationManagerServiceClient
	timeout time.Duration
}

func NewOperationManager(host string, timeout time.Duration) OperationManager {
	conn, err := grpc.NewClient(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to operation manager at %s: %v", host, err)
	}
	return &operationManager{
		conn:    conn,
		client:  pb.NewOperationManagerServiceClient(conn),
		timeout: timeout,
	}
}

func (m *operationManager) Start(req *pb.StartOperationRequest) (*pb.StartOperationResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()
	return m.client.StartOperation(ctx, req)
}

func (m *operationManager) MarkRunning(id string, fencingToken uint64) (*pb.MarkRunningResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()
	return m.client.MarkRunning(ctx, &pb.MarkRunningRequest{Id: id, FencingToken: fencingToken})
}

func (m *operationManager) FailOperation(id, actor, reason string) (*pb.ForceUnlockResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()
	return m.client.FailOperation(ctx, &pb.ForceUnlockRequest{Id: id, Actor: actor, Reason: reason})
}

func (m *operationManager) Close() {
	if m.conn != nil {
		_ = m.conn.Close()
	}
}
