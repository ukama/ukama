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
	"fmt"
	"time"

	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/ukama/ukama/testing/services/dummy/controller/pb/gen"
	"google.golang.org/grpc"
)
 
 type Controller struct {
	 conn    *grpc.ClientConn
	 client  pb.MetricsControllerClient
	 timeout time.Duration
	 host    string
 }
 
 func NewController(controllerHost string, timeout time.Duration) (*Controller, error) {
	 conn, err := grpc.Dial(controllerHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	 if err != nil {
		 return nil, fmt.Errorf("failed to connect to controller: %w", err)
	 }
	 client := pb.NewMetricsControllerClient(conn)
 
	 return &Controller{
		 conn:    conn,
		 client:  client,
		 timeout: timeout,
		 host:    controllerHost,
	 }, nil
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

 func (h *Controller) Start(req *pb.StartMetricsRequest) (*pb.StartMetricsResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), h.timeout)
	defer cancel()

	resp, err := h.client.StartMetrics(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
 }