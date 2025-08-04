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

type RateClient struct {
	conn    *grpc.ClientConn
	timeout time.Duration
	host    string
	client  pb.RateServiceClient
}

func NewRateClient(rateHost string, timeout time.Duration) *RateClient {
	conn, err := grpc.NewClient(rateHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to Rate Service: %v", err)
	}
	client := pb.NewRateServiceClient(conn)

	return &RateClient{
		conn:    conn,
		client:  client,
		timeout: timeout,
		host:    rateHost,
	}
}

func NewRateClientFromClient(client pb.RateServiceClient) *RateClient {
	return &RateClient{
		host:    "localhost",
		timeout: 1 * time.Second,
		conn:    nil,
		client:  client,
	}
}

func (r *RateClient) Close() {
	if r.conn != nil {
		if err := r.conn.Close(); err != nil {
			log.Warnf("Failed to gracefully close Rate Service connection: %v", err)
		}
	}
}

func (r *RateClient) GetRate(req *pb.GetRateRequest) (*pb.GetRateResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	return r.client.GetRate(ctx, req)
}

func (r *RateClient) UpdateDefaultMarkup(req *pb.UpdateDefaultMarkupRequest) (*pb.UpdateDefaultMarkupResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	return r.client.UpdateDefaultMarkup(ctx, req)
}

func (r *RateClient) GetDefaultMarkup(req *pb.GetDefaultMarkupRequest) (*pb.GetDefaultMarkupResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	return r.client.GetDefaultMarkup(ctx, req)
}

func (r *RateClient) GetDefaultMarkupHistory(req *pb.GetDefaultMarkupHistoryRequest) (*pb.GetDefaultMarkupHistoryResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	return r.client.GetDefaultMarkupHistory(ctx, req)
}

func (r *RateClient) UpdateMarkup(req *pb.UpdateMarkupRequest) (*pb.UpdateMarkupResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	return r.client.UpdateMarkup(ctx, req)
}

func (r *RateClient) DeleteMarkup(req *pb.DeleteMarkupRequest) (*pb.DeleteMarkupResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	return r.client.DeleteMarkup(ctx, req)
}

func (r *RateClient) GetMarkup(req *pb.GetMarkupRequest) (*pb.GetMarkupResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	return r.client.GetMarkup(ctx, req)
}

func (r *RateClient) GetMarkupHistory(req *pb.GetMarkupHistoryRequest) (*pb.GetMarkupHistoryResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	return r.client.GetMarkupHistory(ctx, req)
}
