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
	pb "github.com/ukama/ukama/systems/messaging/nns/pb/gen"
)

type Nns struct {
	conn    *grpc.ClientConn
	timeout time.Duration
	client  pb.NnsClient
	host    string
}

func NewNns(host string, timeout time.Duration) *Nns {
	conn, err := grpc.NewClient(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to NNS Service: %v", err)
	}
	client := pb.NewNnsClient(conn)

	return &Nns{
		conn:    conn,
		client:  client,
		timeout: timeout,
		host:    host,
	}
}

func NewNnsFromClient(NnsClient pb.NnsClient) *Nns {
	return &Nns{
		host:    "localhost",
		timeout: 1 * time.Second,
		conn:    nil,
		client:  NnsClient,
	}
}

func (n *Nns) Close() {
	if n.conn != nil {
		if err := n.conn.Close(); err != nil {
			log.Warnf("Failed to gracefully close NNS server connection: %v", err)
		}
	}
}

func (n *Nns) GetNodeRequest(req *pb.GetNodeRequest) (*pb.GetNodeResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), n.timeout)
	defer cancel()

	return n.client.GetNode(ctx, req)
}

func (n *Nns) GetMeshRequest(req *pb.GetMeshRequest) (*pb.GetMeshResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), n.timeout)
	defer cancel()

	return n.client.GetMesh(ctx, req)
}

func (n *Nns) SetRequest(req *pb.SetRequest) (*pb.SetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), n.timeout)
	defer cancel()

	return n.client.Set(ctx, req)
}

func (n *Nns) UpdateMeshRequest(req *pb.UpdateMeshRequest) (*pb.UpdateMeshResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), n.timeout)
	defer cancel()

	return n.client.UpdateMesh(ctx, req)
}

func (n *Nns) UpdateNodeRequest(req *pb.UpdateNodeRequest) (*pb.UpdateNodeResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), n.timeout)
	defer cancel()

	return n.client.UpdateNode(ctx, req)
}

func (n *Nns) DeleteRequest(req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), n.timeout)
	defer cancel()

	return n.client.Delete(ctx, req)
}

func (n *Nns) ListRequest(req *pb.ListRequest) (*pb.ListResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), n.timeout)
	defer cancel()

	return n.client.List(ctx, req)
}
