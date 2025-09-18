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
	pb "github.com/ukama/ukama/testing/services/dummy/dsubscriber/pb/gen"
	"google.golang.org/grpc"
)

type Dsubscriber struct {
	conn    *grpc.ClientConn
	client  pb.DsubscriberServiceClient
	timeout time.Duration
	host    string
}

func NewDsubscriber(healthHost string, timeout time.Duration) *Dsubscriber {
	conn, err := grpc.NewClient(healthHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatalf("did not connect: %v", err)
	}
	client := pb.NewDsubscriberServiceClient(conn)

	return &Dsubscriber{
		conn:    conn,
		client:  client,
		timeout: timeout,
		host:    healthHost,
	}
}

func NewHealthFromClient(mClient pb.DsubscriberServiceClient) *Dsubscriber {
	return &Dsubscriber{
		host:    "localhost",
		timeout: 1 * time.Second,
		conn:    nil,
		client:  mClient,
	}
}

func (r *Dsubscriber) Close() {
	if r.conn != nil {
		if err := r.conn.Close(); err != nil {
			logrus.Errorf("failed to close connection: %v", err)
		}
	}
}

func (h *Dsubscriber) Update(req *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), h.timeout)
	defer cancel()

	resp, err := h.client.Update(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
