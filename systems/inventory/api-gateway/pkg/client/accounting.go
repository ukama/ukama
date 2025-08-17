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

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/inventory/accounting/pb/gen"
)

type Accounting interface {
	Get(id string) (*pb.GetResponse, error)
	GetByUser(uid string) (*pb.GetByUserResponse, error)
	SyncAccounts() (*pb.SyncAcountingResponse, error)
}

type AccountingInventory struct {
	conn    *grpc.ClientConn
	client  pb.AccountingServiceClient
	timeout time.Duration
	host    string
}

func NewAccountingInventory(accountHost string, timeout time.Duration) *AccountingInventory {
	conn, err := grpc.NewClient(accountHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to Accounting Service: %v", err)
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

func (a *AccountingInventory) Close() {
	if a.conn != nil {
		if err := a.conn.Close(); err != nil {
			log.Warnf("Failed to gracefully close Accounting Service connection: %v", err)
		}
	}
}

func (a *AccountingInventory) Get(id string) (*pb.GetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), a.timeout)
	defer cancel()

	return a.client.Get(ctx, &pb.GetRequest{
		Id: id,
	})
}

func (a *AccountingInventory) GetByUser(uid string) (*pb.GetByUserResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), a.timeout)
	defer cancel()

	return a.client.GetByUser(ctx, &pb.GetByUserRequest{
		UserId: uid,
	})
}

func (a *AccountingInventory) SyncAccounts() (*pb.SyncAcountingResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), a.timeout)
	defer cancel()

	return a.client.SyncAccounting(ctx, &pb.SyncAcountingRequest{})
}
