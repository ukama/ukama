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
	pb "github.com/ukama/ukama/systems/inventory/component/pb/gen"
)

type Component interface {
	Get(id string) (*pb.GetResponse, error)
	GetByUser(uid string, c string) (*pb.GetByUserResponse, error)
	SyncComponent() (*pb.SyncComponentsResponse, error)
	List(userId, partNumber, category string) (*pb.ListResponse, error)
	Verify(partNumber string) (*pb.VerifyResponse, error)
	StartScheduler() (*pb.StartSchedulerResponse, error)
	StopScheduler() (*pb.StopSchedulerResponse, error)
}

type ComponentInventory struct {
	conn    *grpc.ClientConn
	client  pb.ComponentServiceClient
	timeout time.Duration
	host    string
}

func NewComponentInventory(componentHost string, timeout time.Duration) *ComponentInventory {
	conn, err := grpc.NewClient(componentHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to Component Service: %v", err)
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

func (c *ComponentInventory) Close() {
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			log.Warnf("Failed to gracefully close Component Service connection: %v", err)
		}
	}
}

func (c *ComponentInventory) Get(id string) (*pb.GetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	return c.client.Get(ctx, &pb.GetRequest{
		Id: id,
	})
}

func (c *ComponentInventory) GetByUser(uid, category string) (*pb.GetByUserResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	return c.client.GetByUser(ctx, &pb.GetByUserRequest{
		UserId:   uid,
		Category: category,
	})
}

func (c *ComponentInventory) SyncComponent() (*pb.SyncComponentsResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	return c.client.SyncComponents(ctx, &pb.SyncComponentsRequest{})
}

func (c *ComponentInventory) List(userId, partNumber, category string) (*pb.ListResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	return c.client.List(ctx, &pb.ListRequest{
		UserId:     userId,
		PartNumber: partNumber,
		Category:   category,
	})
}

func (c *ComponentInventory) Verify(partNumber string) (*pb.VerifyResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	return c.client.Verify(ctx, &pb.VerifyRequest{
		PartNumber: partNumber,
	})
}

func (c *ComponentInventory) StartScheduler() (*pb.StartSchedulerResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	return c.client.StartScheduler(ctx, &pb.StartSchedulerRequest{})
}

func (c *ComponentInventory) StopScheduler() (*pb.StopSchedulerResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	return c.client.StopScheduler(ctx, &pb.StopSchedulerRequest{})
}
