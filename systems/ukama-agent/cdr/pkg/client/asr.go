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

type asrClient struct {
	conn    *grpc.ClientConn
	timeout time.Duration
	client  pb.AsrRecordServiceClient
}

type AsrService interface {
	GetAsr(imsi string) (*pb.ReadResp, error)
}

func NewAsrClient(asrHost string, timeout time.Duration) (AsrService, error) {
	conn, err := grpc.NewClient(asrHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Errorf("Failed to connect to ASR service at %s. Error %s", asrHost, err.Error())

		return nil, err
	}

	client := pb.NewAsrRecordServiceClient(conn)

	return &asrClient{
		conn:    conn,
		client:  client,
		timeout: timeout,
	}, nil
}

func (c *asrClient) Close() {
	err := c.conn.Close()
	if err != nil {
		log.Errorf("Failed to close ASR client connection. Error: %v ", err)
	}
}

func (c *asrClient) GetAsr(imsi string) (*pb.ReadResp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	req := &pb.ReadReq{
		Id: &pb.ReadReq_Imsi{
			Imsi: imsi,
		},
	}

	return c.client.Read(ctx, req)
}
