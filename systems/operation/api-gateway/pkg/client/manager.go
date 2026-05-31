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

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/ukama/ukama/systems/operation/manager/pb/gen"
)

type Manager struct {
	conn    *grpc.ClientConn
	client  pb.OperationManagerServiceClient
	timeout time.Duration
	host    string
}

func NewManager(host string, timeout time.Duration) *Manager {
	conn, err := grpc.NewClient(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to operation manager: %v", err)
	}
	return &Manager{
		conn:    conn,
		client:  pb.NewOperationManagerServiceClient(conn),
		timeout: timeout,
		host:    host,
	}
}

func NewManagerFromClient(c pb.OperationManagerServiceClient) *Manager {
	return &Manager{
		host:    "localhost",
		timeout: 1 * time.Second,
		conn:    nil,
		client:  c,
	}
}

func (m *Manager) Close() {
	if m.conn != nil {
		if err := m.conn.Close(); err != nil {
			log.Warnf("Failed to close manager connection: %v", err)
		}
	}
}

func (m *Manager) Start(req *pb.StartOperationRequest) (*pb.StartOperationResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()
	return m.client.StartOperation(ctx, req)
}

func (m *Manager) Get(id string) (*pb.GetOperationResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()
	return m.client.GetOperation(ctx, &pb.GetOperationRequest{Id: id})
}

func (m *Manager) GetByResource(resourceKey string) (*pb.GetByResourceResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()
	return m.client.GetByResource(ctx, &pb.GetByResourceRequest{ResourceKey: resourceKey})
}

func (m *Manager) MarkRunning(id string, fencingToken uint64) (*pb.MarkRunningResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()
	return m.client.MarkRunning(ctx, &pb.MarkRunningRequest{Id: id, FencingToken: fencingToken})
}

func (m *Manager) ForceUnlock(id, actor, reason string) (*pb.ForceUnlockResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()
	return m.client.ForceUnlock(ctx, &pb.ForceUnlockRequest{Id: id, Actor: actor, Reason: reason})
}
