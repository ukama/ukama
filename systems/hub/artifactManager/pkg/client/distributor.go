/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2024-present, Ukama Inc.
 */

package client

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/hub/distributor/pb/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Distributor struct {
	conn    *grpc.ClientConn
	timeout time.Duration
	client  pb.DistributorServiceClient
	host    string
}

func NewDistributor(host string, timeout time.Duration) *Distributor {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	client := pb.NewDistributorServiceClient(conn)

	return &Distributor{
		conn:    conn,
		client:  client,
		timeout: timeout,
		host:    host,
	}
}

func NewDistributorFromClient(c pb.DistributorServiceClient) *Distributor {
	return &Distributor{
		host:    "localhost",
		timeout: 1 * time.Second,
		conn:    nil,
		client:  c,
	}
}

func (d *Distributor) Close() {
	d.conn.Close()
}

func (d *Distributor) CreateChunk(in *pb.CreateChunkRequest) (*pb.CreateChunkResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), d.timeout)
	defer cancel()
	log.Infof("Sending chunking request: %+v", in)

	return d.client.CreateChunk(ctx, in)
}

func (d *Distributor) GetChunk(in *pb.GetChunkRequest) (*pb.GetChunkResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), d.timeout)
	defer cancel()

	return d.client.GetChunk(ctx, in)
}
