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

func (l *Lookup) Close() {
	if l.conn != nil {
		err := l.conn.Close()
		if err != nil {
			log.Warnf("Failed to gracefully close Lookup Service connection :%v", err)
		}
	}
}

func (l *Lookup) AddOrg(req *pb.AddOrgRequest) (*pb.AddOrgResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), l.timeout)
	defer cancel()

	return l.client.AddOrg(ctx, req)
}

func (l *Lookup) UpdateOrg(req *pb.UpdateOrgRequest) (*pb.UpdateOrgResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), l.timeout)
	defer cancel()

	return l.client.UpdateOrg(ctx, req)
}

func (l *Lookup) GetOrg(req *pb.GetOrgRequest) (*pb.GetOrgResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), l.timeout)
	defer cancel()

	return l.client.GetOrg(ctx, req)
}

func (l *Lookup) GetOrgs(req *pb.GetOrgsRequest) (*pb.GetOrgsResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), l.timeout)
	defer cancel()

	return l.client.GetOrgs(ctx, req)
}

func (l *Lookup) AddNodeForOrg(req *pb.AddNodeRequest) (*pb.AddNodeResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), l.timeout)
	defer cancel()

	return l.client.AddNodeForOrg(ctx, req)
}

func (l *Lookup) GetNodeForOrg(req *pb.GetNodeForOrgRequest) (*pb.GetNodeResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), l.timeout)
	defer cancel()

	return l.client.GetNodeForOrg(ctx, req)
}

func (l *Lookup) DeleteNodeForOrg(req *pb.DeleteNodeRequest) (*pb.DeleteNodeResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), l.timeout)
	defer cancel()

	return l.client.DeleteNodeForOrg(ctx, req)
}

func (l *Lookup) AddSystemForOrg(req *pb.AddSystemRequest) (*pb.AddSystemResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), l.timeout)
	defer cancel()

	return l.client.AddSystemForOrg(ctx, req)
}

func (l *Lookup) UpdateSystemForOrg(req *pb.UpdateSystemRequest) (*pb.UpdateSystemResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), l.timeout)
	defer cancel()

	return l.client.UpdateSystemForOrg(ctx, req)
}

func (l *Lookup) GetSystemForOrg(req *pb.GetSystemRequest) (*pb.GetSystemResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), l.timeout)
	defer cancel()

	return l.client.GetSystemForOrg(ctx, req)
}

func (l *Lookup) DeleteSystemForOrg(req *pb.DeleteSystemRequest) (*pb.DeleteSystemResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), l.timeout)
	defer cancel()

	return l.client.DeleteSystemForOrg(ctx, req)
}
