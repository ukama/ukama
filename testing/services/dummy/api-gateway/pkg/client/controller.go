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

	"github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/testing/services/dummy/controller/pb/gen"
	"google.golang.org/grpc"
)

type Controller struct {
	conn    *grpc.ClientConn
	client  pb.MetricsControllerClient
	timeout time.Duration
	host    string
}

func NewController(controllerHost string, timeout time.Duration) *Controller {
	conn, err := grpc.NewClient(controllerHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatalf("did not connect: %v", err)
	}
	client := pb.NewMetricsControllerClient(conn)

	return &Controller{
		conn:    conn,
		client:  client,
		timeout: timeout,
		host:    controllerHost,
	}
}

func NewHealthFromClient(mClient pb.MetricsControllerClient) *Controller {
	return &Controller{
		host:    "localhost",
		timeout: 1 * time.Second,
		conn:    nil,
		client:  mClient,
	}
}

func (r *Controller) Close() {
	r.conn.Close()
}

func (h *Controller) Update(req *pb.UpdateMetricsRequest) (*pb.UpdateMetricsResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), h.timeout)
	defer cancel()

	resp, err := h.client.UpdateMetrics(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
