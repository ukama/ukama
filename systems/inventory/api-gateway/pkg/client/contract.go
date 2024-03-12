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
	pb "github.com/ukama/ukama/systems/inventory/contract/pb/gen"
	"google.golang.org/grpc"
)

type ContractInventory struct {
	conn    *grpc.ClientConn
	client  pb.ContractServiceClient
	timeout time.Duration
	host    string
}

func NewContractInventory(contractHost string, timeout time.Duration) *ContractInventory {
	// using same context for three connections
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, contractHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatalf("did not connect: %v", err)
	}
	client := pb.NewContractServiceClient(conn)

	return &ContractInventory{
		conn:    conn,
		client:  client,
		timeout: timeout,
		host:    contractHost,
	}
}

func NewNewContractInventoryFromClient(mClient pb.ContractServiceClient) *ContractInventory {
	return &ContractInventory{
		host:    "localhost",
		timeout: 1 * time.Second,
		conn:    nil,
		client:  mClient,
	}
}

func (r *ContractInventory) Close() {
	r.conn.Close()
}

func (r *ContractInventory) GetContracts(c string, a bool) (*pb.GetContractsResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	res, err := r.client.GetContracts(ctx, &pb.GetContractsRequest{
		Company: c,
		Active:  a,
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}
