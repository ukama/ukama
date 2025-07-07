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
	pb "github.com/ukama/ukama/systems/inventory/component/pb/gen"
	"google.golang.org/grpc"
)

type ComponentInventory struct {
	conn    *grpc.ClientConn
	client  pb.ComponentServiceClient
	timeout time.Duration
	host    string
}

func NewComponentInventory(componentHost string, timeout time.Duration) *ComponentInventory {

	conn, err := grpc.NewClient(componentHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	client := pb.NewComponentServiceClient(conn)

	return &ComponentInventory{
		conn:    conn,
		client:  client,
		timeout: timeout,
		host:    componentHost,
	}
}

func NewComponentInventoryFromClient(mClient pb.ComponentServiceClient) *ComponentInventory {
	return &ComponentInventory{
		host:    "localhost",
		timeout: 1 * time.Second,
		conn:    nil,
		client:  mClient,
	}
}

func (r *ComponentInventory) Close() {
	r.conn.Close()
}

func (r *ComponentInventory) Get(id string) (*pb.GetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	return r.client.Get(ctx, &pb.GetRequest{
		Id: id,
	})
}

func (r *ComponentInventory) GetByUser(uid string, c string) (*pb.GetByUserResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	return r.client.GetByUser(ctx, &pb.GetByUserRequest{
		UserId:   uid,
		Category: c,
	})
}

func (r *ComponentInventory) SyncComponent() (*pb.SyncComponentsResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	return r.client.SyncComponents(ctx, &pb.SyncComponentsRequest{})
}
