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
	pb "github.com/ukama/ukama/systems/inventory/accounting/pb/gen"
	"google.golang.org/grpc"
)

type AccountingInventory struct {
	conn    *grpc.ClientConn
	client  pb.AccountingServiceClient
	timeout time.Duration
	host    string
}

func NewAccountingInventory(accountHost string, timeout time.Duration) *AccountingInventory {

	conn, err := grpc.NewClient(accountHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	client := pb.NewAccountingServiceClient(conn)

	return &AccountingInventory{
		conn:    conn,
		client:  client,
		timeout: timeout,
		host:    accountHost,
	}
}

func NewAccountingInventoryFromClient(mClient pb.AccountingServiceClient) *AccountingInventory {
	return &AccountingInventory{
		host:    "localhost",
		timeout: 1 * time.Second,
		conn:    nil,
		client:  mClient,
	}
}

func (r *AccountingInventory) Close() {
	r.conn.Close()
}

func (r *AccountingInventory) Get(id string) (*pb.GetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	res, err := r.client.Get(ctx, &pb.GetRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *AccountingInventory) GetByUser(uid string) (*pb.GetByUserResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	res, err := r.client.GetByUser(ctx, &pb.GetByUserRequest{
		UserId: uid,
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *AccountingInventory) SyncAccounts() (*pb.SyncAcountingResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	res, err := r.client.SyncAccounting(ctx, &pb.SyncAcountingRequest{})
	if err != nil {
		return nil, err
	}

	return res, nil
}
