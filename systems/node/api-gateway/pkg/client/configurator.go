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
	pb "github.com/ukama/ukama/systems/node/configurator/pb/gen"
	"google.golang.org/grpc"
)

type Configurator struct {
	conn    *grpc.ClientConn
	client  pb.ConfiguratorServiceClient
	timeout time.Duration
	host    string
}

func NewConfigurator(configuratorHost string, timeout time.Duration) *Configurator {

	conn, err := grpc.NewClient(configuratorHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatalf("did not connect: %v", err)
	}
	client := pb.NewConfiguratorServiceClient(conn)

	return &Configurator{
		conn:    conn,
		client:  client,
		timeout: timeout,
		host:    configuratorHost,
	}
}

func NewConfiguratorFromClient(mClient pb.ConfiguratorServiceClient) *Configurator {
	return &Configurator{
		host:    "localhost",
		timeout: 1 * time.Second,
		conn:    nil,
		client:  mClient,
	}
}

func (r *Configurator) Close() {
	r.conn.Close()
}

func (r *Configurator) ConfigEvent(b []byte) (*pb.ConfigStoreEventResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	res, err := r.client.ConfigEvent(ctx, &pb.ConfigStoreEvent{
		Data: b,
	})

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *Configurator) ApplyConfig(commit string) (*pb.ApplyConfigResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	res, err := r.client.ApplyConfig(ctx, &pb.ApplyConfigRequest{Hash: commit})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *Configurator) GetConfigVersion(nodeId string) (*pb.ConfigVersionResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	res, err := r.client.GetConfigVersion(ctx, &pb.ConfigVersionRequest{NodeId: nodeId})
	if err != nil {
		return nil, err
	}

	return res, nil
}
