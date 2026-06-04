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
	colpb "github.com/ukama/ukama/systems/analytics/collector/pb/gen"
)

type CollectorAnalytics struct {
	conn    *grpc.ClientConn
	client  colpb.CollectorServiceClient
	timeout time.Duration
	host    string
}

func NewCollectorAnalytics(collectorHost string, timeout time.Duration) *CollectorAnalytics {
	conn, err := grpc.NewClient(collectorHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to analytics' collector service: %v", err)
	}
	client := colpb.NewCollectorServiceClient(conn)

	return &CollectorAnalytics{
		conn:    conn,
		client:  client,
		timeout: timeout,
		host:    collectorHost,
	}
}

func NewCollectorAnalyticsFromClient(collectorClient colpb.CollectorServiceClient) *CollectorAnalytics {
	return &CollectorAnalytics{
		host:    "localhost",
		timeout: 1 * time.Second,
		conn:    nil,
		client:  collectorClient,
	}
}

func (c *CollectorAnalytics) Close() {
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			log.Warnf("Failed to gracefully close Collector Service connection: %v", err)
		}
	}
}

func (c *CollectorAnalytics) Refresh(source string) (*colpb.RefreshResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	return c.client.Refresh(ctx, &colpb.RefreshRequest{
		Source: source,
	})
}

func (c *CollectorAnalytics) GetRefreshState() (*colpb.GetRefreshStateResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	return c.client.GetRefreshState(ctx, &colpb.GetRefreshStateRequest{})
}

func (c *CollectorAnalytics) RebuildRollups(family, from, to string) (*colpb.RebuildRollupsResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	return c.client.RebuildRollups(ctx, &colpb.RebuildRollupsRequest{
		Family: family,
		From:   parseTime(from),
		To:     parseTime(to),
	})
}

func (c *CollectorAnalytics) SeedDemo(sites, nodes, customers, days uint32) (*colpb.SeedDemoResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	return c.client.SeedDemo(ctx, &colpb.SeedDemoRequest{
		Sites:     sites,
		Nodes:     nodes,
		Customers: customers,
		Days:      days,
	})
}
