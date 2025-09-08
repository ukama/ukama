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

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/hub/distributor/pb/gen"
)

type Chunker struct {
	conn    *grpc.ClientConn
	timeout time.Duration
	client  pb.ChunkerServiceClient
	host    string
}

func NewChunker(host string, maxMsgSize int, timeout time.Duration) *Chunker {
	conn, err := grpc.NewClient(host, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultCallOptions(
		grpc.MaxCallRecvMsgSize(maxMsgSize),
		grpc.MaxCallSendMsgSize(maxMsgSize)))
	if err != nil {
		log.Fatalf("Failed to connect to Chunker Service: %v", err)
	}

	client := pb.NewChunkerServiceClient(conn)

	return &Chunker{
		conn:    conn,
		client:  client,
		timeout: timeout,
		host:    host,
	}
}

func NewChunkerFromClient(c pb.ChunkerServiceClient) *Chunker {
	return &Chunker{
		host:    "localhost",
		timeout: 1 * time.Second,
		conn:    nil,
		client:  c,
	}
}

func (c *Chunker) Close() {
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			log.Warnf("Failed to gracefully close Chunker Service connection: %v", err)
		}
	}
}

func (c *Chunker) CreateChunk(in *pb.CreateChunkRequest) (*pb.CreateChunkResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	log.Infof("Sending chunking request: %+v", in)

	return c.client.CreateChunk(ctx, in)
}
