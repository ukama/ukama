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
	pb "github.com/ukama/ukama/systems/data-plan/rate/pb/gen"
)

type Rate struct {
	conn    *grpc.ClientConn
	timeout time.Duration
	client  pb.RateServiceClient
}

type RateService interface {
	GetRates(req *pb.GetRatesRequest) (*pb.GetRatesResponse, error)
	GetRateById(req *pb.GetRateByIdRequest) (*pb.GetRateByIdResponse, error)
}

func NewRate(rate string, timeout time.Duration) (*Rate, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, rate, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Errorf("Failed to connect to rate service at %s. Error %s", rate, err.Error())
		return nil, err
	}
	client := pb.NewRateServiceClient(conn)

	return &Rate{
		conn:    conn,
		client:  client,
		timeout: timeout,
	}, nil
}

func (c *Rate) Close() {
	c.conn.Close()
}

func (c *Rate) GetRates(req *pb.GetRatesRequest) (*pb.GetRatesResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	return c.client.GetRates(ctx, req)
}

func (c *Rate) GetRateById(req *pb.GetRateByIdRequest) (*pb.GetRateByIdResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	return c.client.GetRateById(ctx, req)
}
