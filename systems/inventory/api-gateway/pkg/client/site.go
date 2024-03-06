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
	pb "github.com/ukama/ukama/systems/inventory/site/pb/gen"
	"google.golang.org/grpc"
)

type SiteInventory struct {
	conn    *grpc.ClientConn
	client  pb.SiteServiceClient
	timeout time.Duration
	host    string
}

func NewSiteInventory(siteHost string, timeout time.Duration) *SiteInventory {
	// using same context for three connections
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, siteHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatalf("did not connect: %v", err)
	}
	client := pb.NewSiteServiceClient(conn)

	return &SiteInventory{
		conn:    conn,
		client:  client,
		timeout: timeout,
		host:    siteHost,
	}
}

func NewNewSiteInventoryFromClient(mClient pb.SiteServiceClient) *SiteInventory {
	return &SiteInventory{
		host:    "localhost",
		timeout: 1 * time.Second,
		conn:    nil,
		client:  mClient,
	}
}

func (r *SiteInventory) Close() {
	r.conn.Close()
}

func (r *SiteInventory) Get() (*pb.GetTestResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	res, err := r.client.GetTest(ctx, &pb.GetTestRequest{})
	if err != nil {
		return nil, err
	}

	return res, nil
}
