/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package client

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	log "github.com/sirupsen/logrus"
	pbr "github.com/ukama/ukama/systems/metrics/reasoning/pb/gen"
)
 
 type Reasoning interface {
	GetAlgoStatsForMetric(nodeID string, metric string) (*pbr.GetAlgoStatsForMetricResponse, error)
 }
 
 type reasoning struct {
	conn    *grpc.ClientConn
	client  pbr.ReasoningServiceClient
	timeout time.Duration
	host    string
 }
 
 func NewReasoning(reasoningHost string, timeout time.Duration) *reasoning {
	conn, err := grpc.NewClient(reasoningHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to Reasoning Service: %v", err)
	}
	client := pbr.NewReasoningServiceClient(conn)

	return &reasoning{
		conn:    conn,
		client:  client,
		timeout: timeout,
		host:    reasoningHost,
	}
 }
 
 func NewReasoningFromClient(reasoningClient pbr.ReasoningServiceClient) *reasoning {
	return &reasoning{
		host:    "localhost",
		timeout: 1 * time.Second,
		conn:    nil,
		client:  reasoningClient,
	}
 }
 
 func (s *reasoning) Close() {
	if s.conn != nil {
		if err := s.conn.Close(); err != nil {
			log.Warnf("failed to properly close reasoning client. Error: %v", err)
		}
	}
 }
 
 func (s *reasoning) GetAlgoStatsForMetric(nodeID string, metric string) (*pbr.GetAlgoStatsForMetricResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	return s.client.GetAlgoStatsForMetric(ctx, &pbr.GetAlgoStatsForMetricRequest{
		NodeId: nodeID,
		Metric: metric,
	})
 }
 