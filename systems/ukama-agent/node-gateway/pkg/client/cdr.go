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
	pb "github.com/ukama/ukama/systems/ukama-agent/cdr/pb/gen"
)

type CDR struct {
	conn    *grpc.ClientConn
	timeout time.Duration `default:"3s"`
	client  pb.CDRServiceClient
	host    string `default:"localhost:9090"`
}

func NewCDR(host string, timeout time.Duration) *CDR {
	conn, err := grpc.NewClient(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to CDR Service: %v", err)
	}

	client := pb.NewCDRServiceClient(conn)

	return &CDR{
		conn:    conn,
		client:  client,
		timeout: timeout,
		host:    host,
	}
}

func NewCdrFromClient(asrClient pb.CDRServiceClient) *CDR {
	return &CDR{
		host:    "localhost",
		timeout: 1 * time.Second,
		conn:    nil,
		client:  asrClient,
	}
}

func (c *CDR) Close() {
	if c.conn != nil {
		err := c.conn.Close()
		if err != nil {
			log.Errorf("Failed to close CDR client connection. Error: %v ", err)
		}
	}

}

func (c *CDR) PostCDR(req *pb.CDR) (*pb.CDRResp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	return c.client.PostCDR(ctx, req)
}

func (c *CDR) GetUsage(req *pb.UsageReq) (*pb.UsageResp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	return c.client.GetUsage(ctx, req)
}

func (c *CDR) GetCDR(req *pb.RecordReq) (*pb.RecordResp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	return c.client.GetCDR(ctx, req)
}
