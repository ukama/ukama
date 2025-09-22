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

	"google.golang.org/grpc/credentials/insecure"

	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/testing/services/dummy/dcontroller/pb/gen"
	"google.golang.org/grpc"
)

type DController struct {
	conn    *grpc.ClientConn
	client  pb.MetricsControllerClient
	timeout time.Duration
	host    string
}

func NewDController(dcontrollerHost string, timeout time.Duration) *DController {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.NewClient(dcontrollerHost, opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	client := pb.NewMetricsControllerClient(conn)

	return &DController{
		conn:    conn,
		client:  client,
		timeout: timeout,
		host:    dcontrollerHost,
	}
}

func NewDControllerFromClient(mClient pb.MetricsControllerClient) *DController {
	return &DController{
		host:    "localhost",
		timeout: 1 * time.Second,
		conn:    nil,
		client:  mClient,
	}
}

func (r *DController) Close() {
	if r.conn != nil {
		if err := r.conn.Close(); err != nil {
			log.Errorf("failed to close connection: %v", err)
		}
	}
}

func (h *DController) Update(req *pb.UpdateMetricsRequest) (*pb.UpdateMetricsResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), h.timeout)
	defer cancel()

	resp, err := h.client.Update(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (h *DController) Start(req *pb.StartMetricsRequest) (*pb.StartMetricsResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), h.timeout)
	defer cancel()

	resp, err := h.client.Start(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
