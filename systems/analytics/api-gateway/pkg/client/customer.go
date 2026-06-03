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
	custpb "github.com/ukama/ukama/systems/analytics/customer/pb/gen"
)

type CustomerAnalytics struct {
	conn    *grpc.ClientConn
	client  custpb.CustomerServiceClient
	timeout time.Duration
	host    string
}

func NewCustomerAnalytics(customerHost string, timeout time.Duration) *CustomerAnalytics {
	conn, err := grpc.NewClient(customerHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to analytics' customer service: %v", err)
	}
	client := custpb.NewCustomerServiceClient(conn)

	return &CustomerAnalytics{
		conn:    conn,
		client:  client,
		timeout: timeout,
		host:    customerHost,
	}
}

func NewCustomerAnalyticsFromClient(customerClient custpb.CustomerServiceClient) *CustomerAnalytics {
	return &CustomerAnalytics{
		host:    "localhost",
		timeout: 1 * time.Second,
		conn:    nil,
		client:  customerClient,
	}
}

func (c *CustomerAnalytics) Close() {
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			log.Warnf("Failed to gracefully close Customer Service connection: %v", err)
		}
	}
}

func (c *CustomerAnalytics) window(period, from, to, tz string) *custpb.Window {
	w := toWindow(period, from, to, tz)
	if w.empty {
		return nil
	}

	return &custpb.Window{
		Period:   w.period,
		From:     w.from,
		To:       w.to,
		Timezone: w.timezone,
	}
}

func (c *CustomerAnalytics) GetOverview(networkId, period, from, to, tz string) (*custpb.GetOverviewResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	return c.client.GetOverview(ctx, &custpb.GetOverviewRequest{
		NetworkId: networkId,
		Window:    c.window(period, from, to, tz),
	})
}

func (c *CustomerAnalytics) List(networkId, siteId, status string, page, pageSize uint32) (*custpb.ListResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	return c.client.List(ctx, &custpb.ListRequest{
		NetworkId: networkId,
		SiteId:    siteId,
		Status:    status,
		Page:      page,
		PageSize:  pageSize,
	})
}

func (c *CustomerAnalytics) Search(query, networkId string, page, pageSize uint32) (*custpb.SearchResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	return c.client.Search(ctx, &custpb.SearchRequest{
		Query:     query,
		NetworkId: networkId,
		Page:      page,
		PageSize:  pageSize,
	})
}

func (c *CustomerAnalytics) Get(customerId, period, from, to, tz string) (*custpb.GetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	return c.client.Get(ctx, &custpb.GetRequest{
		CustomerId: customerId,
		Window:     c.window(period, from, to, tz),
	})
}

func (c *CustomerAnalytics) GetSupport(customerId string) (*custpb.GetSupportResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	return c.client.GetSupport(ctx, &custpb.GetSupportRequest{
		CustomerId: customerId,
	})
}

func (c *CustomerAnalytics) GetSims(networkId, status string, page, pageSize uint32) (*custpb.GetSimsResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	return c.client.GetSims(ctx, &custpb.GetSimsRequest{
		NetworkId: networkId,
		Status:    status,
		Page:      page,
		PageSize:  pageSize,
	})
}

func (c *CustomerAnalytics) GetSimPool(networkId string) (*custpb.GetSimPoolResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	return c.client.GetSimPool(ctx, &custpb.GetSimPoolRequest{
		NetworkId: networkId,
	})
}
