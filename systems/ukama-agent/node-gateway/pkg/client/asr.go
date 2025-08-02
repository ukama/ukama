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
	host    string `default:"localhost:9090"`
}

func NewAsr(host string, timeout time.Duration) *Asr {
	conn, err := grpc.NewClient(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to ASR Service: %v", err)
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
		if err := a.conn.Close(); err != nil {
			log.Errorf("Failed to close ASR client connection. Error: %v ", err)
		}
	}
}

func (a *Asr) UpdateGuti(req *pb.UpdateGutiReq) (*pb.UpdateGutiResp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), a.timeout)
	defer cancel()

	return a.client.UpdateGuti(ctx, req)
}

func (a *Asr) UpdateTai(req *pb.UpdateTaiReq) (*pb.UpdateTaiResp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), a.timeout)
	defer cancel()

	return a.client.UpdateTai(ctx, req)
}

func (a *Asr) Read(req *pb.ReadReq) (*pb.ReadResp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), a.timeout)
	defer cancel()

	return a.client.Read(ctx, req)
}
