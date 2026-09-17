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
	ugrpc "github.com/ukama/ukama/systems/common/grpc"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/subscriber/sim-pool/pb/gen"
)

type SimPool struct {
	conn    *grpc.ClientConn
	timeout time.Duration
	client  pb.SimServiceClient
	host    string
}

func NewSimPool(host string, timeout time.Duration) *SimPool {
	conn, err := grpc.NewClient(host, grpc.WithTransportCredentials(insecure.NewCredentials()),
		ugrpc.TracingDialOption())
	if err != nil {
		log.Fatalf("Failed to connect to Sim Pool Service: %v", err)
	}
	client := pb.NewSimServiceClient(conn)

	return &SimPool{
		conn:    conn,
		client:  client,
		timeout: timeout,
		host:    host,
	}
}

func NewSimPoolFromClient(SimPoolClient pb.SimServiceClient) *SimPool {
	return &SimPool{
		host:    "localhost",
		timeout: 1 * time.Second,
		conn:    nil,
		client:  SimPoolClient,
	}
}

func (sp *SimPool) Close() {
	if sp.conn != nil {
		if err := sp.conn.Close(); err != nil {
			log.Warnf("Failed to gracefully close Sim Pool Service connection: %v", err)
		}
	}
}

func (sp *SimPool) Get(ctx context.Context, iccid string) (*pb.GetByIccidResponse, error) {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), sp.timeout)
	defer cancel()

	return sp.client.GetByIccid(ctx, &pb.GetByIccidRequest{Iccid: iccid})
}

func (sp *SimPool) GetSims(ctx context.Context, simType string) (*pb.GetSimsResponse, error) {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), sp.timeout)
	defer cancel()

	return sp.client.GetSims(ctx, &pb.GetSimsRequest{SimType: simType})
}

func (sp *SimPool) GetStats(ctx context.Context, simType string) (*pb.GetStatsResponse, error) {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), sp.timeout)
	defer cancel()

	return sp.client.GetStats(ctx, &pb.GetStatsRequest{SimType: simType})
}

func (sp *SimPool) AddSimsToSimPool(ctx context.Context, req *pb.AddRequest) (*pb.AddResponse, error) {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), sp.timeout)
	defer cancel()

	return sp.client.Add(ctx, req)
}

func (sp *SimPool) UploadSimsToSimPool(ctx context.Context, req *pb.UploadRequest) (*pb.UploadResponse, error) {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), sp.timeout)
	defer cancel()

	return sp.client.Upload(ctx, req)
}

func (sp *SimPool) DeleteSimFromSimPool(ctx context.Context, id []uint64) (*pb.DeleteResponse, error) {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), sp.timeout)
	defer cancel()

	return sp.client.Delete(ctx, &pb.DeleteRequest{Id: id})
}
