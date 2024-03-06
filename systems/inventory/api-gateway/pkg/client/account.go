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
	pb "github.com/ukama/ukama/systems/inventory/account/pb/gen"
	"google.golang.org/grpc"
)

type AccountInventory struct {
	conn    *grpc.ClientConn
	client  pb.AccountServiceClient
	timeout time.Duration
	host    string
}

func NewAccountInventory(accountHost string, timeout time.Duration) *AccountInventory {
	// using same context for three connections
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, accountHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatalf("did not connect: %v", err)
	}
	client := pb.NewAccountServiceClient(conn)

	return &AccountInventory{
		conn:    conn,
		client:  client,
		timeout: timeout,
		host:    accountHost,
	}
}

func NewNewAccountInventoryFromClient(mClient pb.AccountServiceClient) *AccountInventory {
	return &AccountInventory{
		host:    "localhost",
		timeout: 1 * time.Second,
		conn:    nil,
		client:  mClient,
	}
}

func (r *AccountInventory) Close() {
	r.conn.Close()
}

func (r *AccountInventory) Get() (*pb.GetTestResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	res, err := r.client.GetTest(ctx, &pb.GetTestRequest{})
	if err != nil {
		return nil, err
	}

	return res, nil
}
