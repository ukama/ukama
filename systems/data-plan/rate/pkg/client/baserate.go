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

	"github.com/sirupsen/logrus"
	bpb "github.com/ukama/ukama/systems/data-plan/base-rate/pb/gen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type BaseRate struct {
	conn    *grpc.ClientConn
	timeout time.Duration
	client  bpb.BaseRatesServiceClient
}

type BaseRateSrvc interface {
	GetBaseRates(req *bpb.GetBaseRatesByPeriodRequest) (*bpb.GetBaseRatesResponse, error)
	GetBaseRate(req *bpb.GetBaseRatesByIdRequest) (*bpb.GetBaseRatesByIdResponse, error)
}

func NewBaseRate(baseRate string, timeout time.Duration) (*BaseRate, error) {

	conn, err := grpc.NewClient(baseRate, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Errorf("Failed to connect to base rate service at %s. Error %s", baseRate, err.Error())
		return nil, err
	}
	client := bpb.NewBaseRatesServiceClient(conn)

	return &BaseRate{
		conn:    conn,
		client:  client,
		timeout: timeout,
	}, nil
}

func (c *BaseRate) Close() {
	c.conn.Close()
}

func (c *BaseRate) GetBaseRates(req *bpb.GetBaseRatesByPeriodRequest) (*bpb.GetBaseRatesResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	return c.client.GetBaseRatesForPeriod(ctx, req)
}

func (c *BaseRate) GetBaseRate(req *bpb.GetBaseRatesByIdRequest) (*bpb.GetBaseRatesByIdResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	return c.client.GetBaseRatesById(ctx, req)
}
