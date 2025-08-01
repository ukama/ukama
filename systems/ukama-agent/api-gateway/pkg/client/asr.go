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
	pb "github.com/ukama/ukama/systems/ukama-agent/asr/pb/gen"
)

type Asr struct {
	conn    *grpc.ClientConn
	timeout time.Duration `default:"3s"`
	client  pb.AsrRecordServiceClient
	host    string `deafault:"localhost:9090"`
}

func NewAsr(host string, timeout time.Duration) *Asr {

	conn, err := grpc.NewClient(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	client := pb.NewAsrRecordServiceClient(conn)

	return &Asr{
		conn:    conn,
		client:  client,
		timeout: timeout,
		host:    host,
	}
}

func NewAsrFromClient(asrClient pb.AsrRecordServiceClient) *Asr {
	return &Asr{
		host:    "localhost",
		timeout: 1 * time.Second,
		conn:    nil,
		client:  asrClient,
	}
}

func (a *Asr) Close() {
	if a.conn != nil {
		err := a.conn.Close()
		if err != nil {
			log.Errorf("Failed to close ASR client connection. Error: %v ", err)
		}
	}
}

func (a *Asr) Activate(req *pb.ActivateReq) (*pb.ActivateResp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), a.timeout)
	defer cancel()

	return a.client.Activate(ctx, req)
}

func (a *Asr) Inactivate(req *pb.InactivateReq) (*pb.InactivateResp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), a.timeout)
	defer cancel()

	return a.client.Inactivate(ctx, req)
}

func (a *Asr) UpdatePackage(req *pb.UpdatePackageReq) (*pb.UpdatePackageResp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), a.timeout)
	defer cancel()

	return a.client.UpdatePackage(ctx, req)
}

func (a *Asr) Read(req *pb.ReadReq) (*pb.ReadResp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), a.timeout)
	defer cancel()

	return a.client.Read(ctx, req)
}

func (a *Asr) GetUsage(req *pb.UsageReq) (*pb.UsageResp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), a.timeout)
	defer cancel()

	return a.client.GetUsage(ctx, req)
}

func (a *Asr) GetUsageForPeriod(req *pb.UsageForPeriodReq) (*pb.UsageResp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), a.timeout)
	defer cancel()

	return a.client.GetUsageForPeriod(ctx, req)
}

func (a *Asr) QueryUsage(req *pb.QueryUsageReq) (*pb.QueryUsageResp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), a.timeout)
	defer cancel()

	return a.client.QueryUsage(ctx, req)
}
