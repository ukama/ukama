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

	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/testing/services/dummy-node/dnode/pb/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type NodeDummy struct {
	conn    *grpc.ClientConn
	client  pb.NodeServiceClient
	timeout time.Duration
	host    string
}

func NewNodeService(accountHost string, timeout time.Duration) *NodeDummy {

	conn, err := grpc.NewClient(accountHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	client := pb.NewNodeServiceClient(conn)

	return &NodeDummy{
		conn:    conn,
		client:  client,
		timeout: timeout,
		host:    accountHost,
	}
}

func NewNodeServiceFromClient(mClient pb.NodeServiceClient) *NodeDummy {
	return &NodeDummy{
		host:    "localhost",
		timeout: 1 * time.Second,
		conn:    nil,
		client:  mClient,
	}
}

func (r *NodeDummy) Close() {
	r.conn.Close()
}

func (r *NodeDummy) ResetNode(id string) (*pb.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	res, err := r.client.ResetNode(ctx, &pb.Request{NodeId: id})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *NodeDummy) TurnNodeOff(id string) (*pb.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	res, err := r.client.TurnNodeOff(ctx, &pb.Request{NodeId: id})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *NodeDummy) TurnRFOn(id string) (*pb.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	res, err := r.client.NodeRFOn(ctx, &pb.Request{NodeId: id})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *NodeDummy) TurnRFOff(id string) (*pb.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	res, err := r.client.NodeRFOff(ctx, &pb.Request{NodeId: id})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *NodeDummy) TurnNodeOnline(id string) (*pb.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	res, err := r.client.TurnNodeOnline(ctx, &pb.Request{NodeId: id})
	if err != nil {
		return nil, err
	}

	return res, nil
}
