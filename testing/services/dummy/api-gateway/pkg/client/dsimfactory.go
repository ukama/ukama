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
	pb "github.com/ukama/ukama/testing/services/dummy/dsimfactory/pb/gen"
	"google.golang.org/grpc"
)

type Dsimfactory struct {
	conn    *grpc.ClientConn
	client  pb.DsimfactoryServiceClient
	timeout time.Duration
	host    string
}

func NewDsimfactory(healthHost string, timeout time.Duration) *Dsimfactory {
	conn, err := grpc.NewClient(healthHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatalf("did not connect: %v", err)
	}
	client := pb.NewDsimfactoryServiceClient(conn)

	return &Dsimfactory{
		conn:    conn,
		client:  client,
		timeout: timeout,
		host:    healthHost,
	}
}

func NewDsimfactoryFromClient(mClient pb.DsimfactoryServiceClient) *Dsimfactory {
	return &Dsimfactory{
		host:    "localhost",
		timeout: 1 * time.Second,
		conn:    nil,
		client:  mClient,
	}
}

func (r *Dsimfactory) Close() {
	if r.conn != nil {
		if err := r.conn.Close(); err != nil {
			logrus.Errorf("failed to close connection: %v", err)
		}
	}
}

func (sp *Dsimfactory) GetSims() (*pb.GetSimsResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), sp.timeout)
	defer cancel()

	return sp.client.GetSims(ctx, &pb.GetSimsRequest{})
}

func (sp *Dsimfactory) GetSim(iccid string) (*pb.GetByIccidResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), sp.timeout)
	defer cancel()

	return sp.client.GetByIccid(ctx, &pb.GetByIccidRequest{
		Iccid: iccid,
	})
}

func (sp *Dsimfactory) UploadSimsToSimPool(req *pb.UploadRequest) (*pb.UploadResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), sp.timeout)
	defer cancel()

	return sp.client.Upload(ctx, req)
}
