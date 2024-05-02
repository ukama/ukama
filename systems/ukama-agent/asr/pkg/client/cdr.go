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
	pb "github.com/ukama/ukama/systems/ukama-agent/cdr/pb/gen"
)

type CDR struct {
	conn    *grpc.ClientConn
	timeout time.Duration
	client  pb.CDRServiceClient
}

type CDRService interface {
	GetUsage(req string) (*pb.UsageResp, error)
	GetUsageForPeriod(imsi string, startTime uint64, endTime uint64) (*pb.UsageForPeriodResp, error)
}

func NewCDR(cdr string, timeout time.Duration) (*CDR, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, cdr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Errorf("Failed to connect to cdr service at %s. Error %s", cdr, err.Error())
		return nil, err
	}
	client := pb.NewCDRServiceClient(conn)

	return &CDR{
		conn:    conn,
		client:  client,
		timeout: timeout,
	}, nil
}

func (c *CDR) Close() {
	c.conn.Close()
}

func (c *CDR) GetUsage(imsi string) (*pb.UsageResp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	req := &pb.UsageReq{
		Imsi: imsi,
	}

	return c.client.GetUsage(ctx, req)
}

func (c *CDR) GetUsageForPeriod(imsi string, startTime uint64, endTime uint64) (*pb.UsageForPeriodResp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	req := &pb.UsageForPeriodReq{
		Imsi:      imsi,
		StartTime: startTime,
		EndTime:   endTime,
	}

	return c.client.GetUsageForPeriod(ctx, req)
}
