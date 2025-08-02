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
	pb "github.com/ukama/ukama/systems/init/lookup/pb/gen"
)

type Lookup struct {
	conn    *grpc.ClientConn
	timeout time.Duration
	client  pb.LookupServiceClient
	host    string
}

func Newlookup(host string, timeout time.Duration) *Lookup {

	conn, err := grpc.NewClient(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to Lookup Service: %v", err)
	}
	client := pb.NewLookupServiceClient(conn)

	return &Lookup{
		conn:    conn,
		client:  client,
		timeout: timeout,
		host:    host,
	}
}

func NewLookupFromClient(lookupClient pb.LookupServiceClient) *Lookup {
	return &Lookup{
		host:    "localhost",
		timeout: 1 * time.Second,
		conn:    nil,
		client:  lookupClient,
	}
}

func (r *Lookup) Close() {
	if r.conn != nil {
		if err := r.conn.Close(); err != nil {
			log.Warnf("Failed to gracefully close Lookup Service connection: %v", err)
		}
	}
}

func (l *Lookup) GetNode(req *pb.GetNodeRequest) (*pb.GetNodeResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), l.timeout)
	defer cancel()

	return l.client.GetNode(ctx, req)
}
