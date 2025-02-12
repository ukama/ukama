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
	pb "github.com/ukama/ukama/testing/services/dummy-node/dmetrics/pb/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Dmetrics struct {
	conn    *grpc.ClientConn
	client  pb.DmetricsServiceClient
	timeout time.Duration
	host    string
}

func NewDmetricsService(accountHost string, timeout time.Duration) *Dmetrics {

	conn, err := grpc.NewClient(accountHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	client := pb.NewDmetricsServiceClient(conn)

	return &Dmetrics{
		conn:    conn,
		client:  client,
		timeout: timeout,
		host:    accountHost,
	}
}

func NewDmetricsServiceFromClient(mClient pb.DmetricsServiceClient) *Dmetrics {
	return &Dmetrics{
		host:    "localhost",
		timeout: 1 * time.Second,
		conn:    nil,
		client:  mClient,
	}
}

func (r *Dmetrics) Close() {
	r.conn.Close()
}

func (r *Dmetrics) NodeMetrics(id string, profile string) (*pb.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	res, err := r.client.NodeMetrics(ctx, &pb.Request{NodeId: id,
		Profile: pb.Profile(pb.Profile_value[profile])})
	if err != nil {
		return nil, err
	}

	return res, nil
}
