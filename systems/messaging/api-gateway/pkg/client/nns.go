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
		log.Fatalf("did not connect: %v", err)
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
		err := n.conn.Close()
		if err != nil {
			log.Warnf("Failed to gracefully close NNS server connection: %v", err)
		}
	}
}

func (n *Nns) GetNodeIpRequest(req *pb.GetNodeIPRequest) (*pb.GetNodeIPResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), n.timeout)
	defer cancel()

	return n.client.Get(ctx, req)
}

func (n *Nns) SetNodeIpRequest(req *pb.SetNodeIPRequest) (*pb.SetNodeIPResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), n.timeout)
	defer cancel()

	return n.client.Set(ctx, req)
}

func (n *Nns) DeleteNodeIpRequest(req *pb.DeleteNodeIPRequest) (*pb.DeleteNodeIPResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), n.timeout)
	defer cancel()

	return n.client.Delete(ctx, req)
}

func (n *Nns) ListNodeIpRequest(req *pb.ListNodeIPRequest) (*pb.ListNodeIPResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), n.timeout)
	defer cancel()

	return n.client.List(ctx, req)
}

func (n *Nns) GetNodeOrgMapListRequest(req *pb.NodeOrgMapListRequest) (*pb.NodeOrgMapListResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), n.timeout)
	defer cancel()

	return n.client.GetNodeOrgMapList(ctx, req)
}

func (n *Nns) GetNodeIPMapListRequest(req *pb.NodeIPMapListRequest) (*pb.NodeIPMapListResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), n.timeout)
	defer cancel()

	return n.client.GetNodeIPMapList(ctx, req)
}
