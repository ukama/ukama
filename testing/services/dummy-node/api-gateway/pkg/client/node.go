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
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type NodeDummy struct {
	conn *grpc.ClientConn
	// client  pb.NodeServiceClient
	timeout time.Duration
	host    string
}

func NewNodeService(accountHost string, timeout time.Duration) *NodeDummy {

	conn, err := grpc.NewClient(accountHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	// client := pb.NewNodeServiceClient(conn)

	return &NodeDummy{
		conn: conn,
		// client:  client,
		timeout: timeout,
		host:    accountHost,
	}
}

// func NewNodeServiceFromClient(mClient pb.NodeServiceClient) *NodeDummy {
// 	return &NodeDummy{
// 		host:    "localhost",
// 		timeout: 1 * time.Second,
// 		conn:    nil,
// 		// client:  mClient,
// 	}
// }

func (r *NodeDummy) Close() {
	r.conn.Close()
}

func (r *NodeDummy) Get() (string, error) {
	_, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	return "uk-sa2450-tnode-v0-4e86", nil
}
