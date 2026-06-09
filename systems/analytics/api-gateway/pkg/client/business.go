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
	bizpb "github.com/ukama/ukama/systems/analytics/business/pb/gen"
)

type BusinessAnalytics struct {
	conn    *grpc.ClientConn
	client  bizpb.BusinessServiceClient
	timeout time.Duration
	host    string
}

func NewBusinessAnalytics(businessHost string, timeout time.Duration) *BusinessAnalytics {
	conn, err := grpc.NewClient(businessHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to analytics' business service: %v", err)
	}
	client := bizpb.NewBusinessServiceClient(conn)

	return &BusinessAnalytics{
		conn:    conn,
		client:  client,
		timeout: timeout,
		host:    businessHost,
	}
}

func NewBusinessAnalyticsFromClient(businessClient bizpb.BusinessServiceClient) *BusinessAnalytics {
	return &BusinessAnalytics{
		host:    "localhost",
		timeout: 1 * time.Second,
		conn:    nil,
		client:  businessClient,
	}
}

func (b *BusinessAnalytics) Close() {
	if b.conn != nil {
		if err := b.conn.Close(); err != nil {
			log.Warnf("Failed to gracefully close Business Service connection: %v", err)
		}
	}
}

func (b *BusinessAnalytics) window(period, from, to, tz string) *bizpb.Window {
	w := toWindow(period, from, to, tz)
	if w.empty {
		return nil
	}

	return &bizpb.Window{
		Period:   w.period,
		From:     w.from,
		To:       w.to,
		Timezone: w.timezone,
	}
}

func (b *BusinessAnalytics) GetHome(networkId, period, from, to, tz string) (*bizpb.GetHomeResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), b.timeout)
	defer cancel()

	return b.client.GetHome(ctx, &bizpb.GetHomeRequest{
		NetworkId: networkId,
		Window:    b.window(period, from, to, tz),
	})
}

func (b *BusinessAnalytics) GetSalesOverview(networkId, period, from, to, tz string) (*bizpb.GetSalesOverviewResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), b.timeout)
	defer cancel()

	return b.client.GetSalesOverview(ctx, &bizpb.GetSalesOverviewRequest{
		NetworkId: networkId,
		Window:    b.window(period, from, to, tz),
	})
}

func (b *BusinessAnalytics) GetPackagePerformance(networkId, period, from, to, tz string, page, pageSize uint32) (*bizpb.GetPackagePerformanceResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), b.timeout)
	defer cancel()

	return b.client.GetPackagePerformance(ctx, &bizpb.GetPackagePerformanceRequest{
		NetworkId: networkId,
		Window:    b.window(period, from, to, tz),
		Page:      page,
		PageSize:  pageSize,
	})
}

func (b *BusinessAnalytics) GetBillingSummary(period, from, to, tz string, page, pageSize uint32) (*bizpb.GetBillingSummaryResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), b.timeout)
	defer cancel()

	return b.client.GetBillingSummary(ctx, &bizpb.GetBillingSummaryRequest{
		Window:   b.window(period, from, to, tz),
		Page:     page,
		PageSize: pageSize,
	})
}

func (b *BusinessAnalytics) GetSites(networkId, period, from, to, tz string, page, pageSize uint32) (*bizpb.GetSitesResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), b.timeout)
	defer cancel()

	return b.client.GetSites(ctx, &bizpb.GetSitesRequest{
		NetworkId: networkId,
		Window:    b.window(period, from, to, tz),
		Page:      page,
		PageSize:  pageSize,
	})
}

func (b *BusinessAnalytics) GetSite(siteId, period, from, to, tz string) (*bizpb.GetSiteResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), b.timeout)
	defer cancel()

	return b.client.GetSite(ctx, &bizpb.GetSiteRequest{
		SiteId: siteId,
		Window: b.window(period, from, to, tz),
	})
}

func (b *BusinessAnalytics) GetInventoryReadiness(networkId string) (*bizpb.GetInventoryReadinessResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), b.timeout)
	defer cancel()

	return b.client.GetInventoryReadiness(ctx, &bizpb.GetInventoryReadinessRequest{
		NetworkId: networkId,
	})
}
