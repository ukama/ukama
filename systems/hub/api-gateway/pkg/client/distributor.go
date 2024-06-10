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
	dpb "github.com/ukama/ukama/systems/hub/distributor/pb/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Distributor struct {
	conn    *grpc.ClientConn
	timeout time.Duration
	client  dpb.DistributorServiceClient
	host    string
}

func NewDistributor(host string, timeout time.Duration) *Distributor {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	client := dpb.NewDistributorServiceClient(conn)

	return &Distributor{
		conn:    conn,
		client:  client,
		timeout: timeout,
		host:    host,
	}
}

func NewDistributorFromClient(c dpb.DistributorServiceClient) *Distributor {
	return &Distributor{
		host:    "localhost",
		timeout: 1 * time.Second,
		conn:    nil,
		client:  c,
	}
}

func (r *Distributor) Close() {
	r.conn.Close()
}

func (d *Distributor) Store(in *dpb.Request) (*dpb.Response, error) {
	return nil, nil
}

func (d *Distributor) Get(in *dpb.ChunkRequest) (*dpb.ChunkResponse, error) {
	return nil, nil
}
