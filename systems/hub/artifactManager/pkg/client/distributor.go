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

type Chunker struct {
	conn    *grpc.ClientConn
	timeout time.Duration
	client  pb.ChunkerServiceClient
	host    string
}

func NewChunker(host string, maxMsgSize int, timeout time.Duration) *Chunker {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, host, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultCallOptions(
		grpc.MaxCallRecvMsgSize(maxMsgSize),
		grpc.MaxCallSendMsgSize(maxMsgSize)))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	client := pb.NewChunkerServiceClient(conn)

	return &Chunker{
		conn:    conn,
		client:  client,
		timeout: timeout,
		host:    host,
	}
}

func NewDistributorFromClient(c pb.ChunkerServiceClient) *Chunker {
	return &Chunker{
		host:    "localhost",
		timeout: 1 * time.Second,
		conn:    nil,
		client:  c,
	}
}

func (d *Chunker) Close() {
	d.conn.Close()
}

func (d *Chunker) CreateChunk(in *pb.CreateChunkRequest) (*pb.CreateChunkResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), d.timeout)
	defer cancel()
	log.Infof("Sending chunking request: %+v", in)

	return d.client.CreateChunk(ctx, in)
}
