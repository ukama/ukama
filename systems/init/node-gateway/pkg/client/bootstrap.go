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
	pb "github.com/ukama/ukama/systems/init/bootstrap/pb/gen"
)

type BootstrapEP interface {
	GetNodeCredentials(req *pb.GetNodeCredentialsRequest) (*pb.GetNodeCredentialsResponse, error)
	GetNodeMeshInfo(req *pb.GetNodeMeshInfoRequest) (*pb.GetNodeMeshInfoResponse, error)
}

type Bootstrap struct {
	conn    *grpc.ClientConn
	client  pb.BootstrapServiceClient
	timeout time.Duration
	host    string
}

func NewBootstrap(bootstrapHost string, timeout time.Duration) *Bootstrap {

	conn, err := grpc.NewClient(bootstrapHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to Bootstrap Service:  %v", err)
	}
	client := pb.NewBootstrapServiceClient(conn)

	return &Bootstrap{
		conn:    conn,
		client:  client,
		timeout: timeout,
		host:    bootstrapHost,
	}
}

func NewBootstrapFromClient(mClient pb.BootstrapServiceClient) *Bootstrap {
	return &Bootstrap{
		host:    "localhost",
		timeout: 1 * time.Second,
		conn:    nil,
		client:  mClient,
	}
}

func (r *Bootstrap) Close() {
	if r.conn != nil {
		if err := r.conn.Close(); err != nil {
			log.Warnf("Failed to gracefully close Bootstrap Service connection: %v", err)
		}
	}
}

func (r *Bootstrap) GetNodeCredentials(req *pb.GetNodeCredentialsRequest) (*pb.GetNodeCredentialsResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	return r.client.GetNodeCredentials(ctx, req)
}

func (r *Bootstrap) GetNodeMeshInfo(req *pb.GetNodeMeshInfoRequest) (*pb.GetNodeMeshInfoResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	return r.client.GetNodeMeshInfo(ctx, req)
}