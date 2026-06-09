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
	netpb "github.com/ukama/ukama/systems/analytics/network/pb/gen"
)

type NetworkAnalytics struct {
	conn    *grpc.ClientConn
	client  netpb.NetworkServiceClient
	timeout time.Duration
	host    string
}

func NewNetworkAnalytics(networkHost string, timeout time.Duration) *NetworkAnalytics {
	conn, err := grpc.NewClient(networkHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to analytics' network service: %v", err)
	}
	client := netpb.NewNetworkServiceClient(conn)

	return &NetworkAnalytics{
		conn:    conn,
		client:  client,
		timeout: timeout,
		host:    networkHost,
	}
}

func NewNetworkAnalyticsFromClient(networkClient netpb.NetworkServiceClient) *NetworkAnalytics {
	return &NetworkAnalytics{
		host:    "localhost",
		timeout: 1 * time.Second,
		conn:    nil,
		client:  networkClient,
	}
}

func (n *NetworkAnalytics) Close() {
	if n.conn != nil {
		if err := n.conn.Close(); err != nil {
			log.Warnf("Failed to gracefully close Network Service connection: %v", err)
		}
	}
}

func (n *NetworkAnalytics) window(period, from, to, tz string) *netpb.Window {
	w := toWindow(period, from, to, tz)
	if w.empty {
		return nil
	}

	return &netpb.Window{
		Period:   w.period,
		From:     w.from,
		To:       w.to,
		Timezone: w.timezone,
	}
}

func (n *NetworkAnalytics) GetOverview(networkId, period, from, to, tz string) (*netpb.GetOverviewResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), n.timeout)
	defer cancel()

	return n.client.GetOverview(ctx, &netpb.GetOverviewRequest{
		NetworkId: networkId,
		Window:    n.window(period, from, to, tz),
	})
}

func (n *NetworkAnalytics) GetTopology(networkId string) (*netpb.GetTopologyResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), n.timeout)
	defer cancel()

	return n.client.GetTopology(ctx, &netpb.GetTopologyRequest{
		NetworkId: networkId,
	})
}

func (n *NetworkAnalytics) GetSites(networkId, status, period, from, to, tz string, page, pageSize uint32) (*netpb.GetSitesResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), n.timeout)
	defer cancel()

	return n.client.GetSites(ctx, &netpb.GetSitesRequest{
		NetworkId: networkId,
		Status:    status,
		Window:    n.window(period, from, to, tz),
		Page:      page,
		PageSize:  pageSize,
	})
}

func (n *NetworkAnalytics) GetSite(siteId, period, from, to, tz string) (*netpb.GetSiteResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), n.timeout)
	defer cancel()

	return n.client.GetSite(ctx, &netpb.GetSiteRequest{
		SiteId: siteId,
		Window: n.window(period, from, to, tz),
	})
}

func (n *NetworkAnalytics) GetNodes(networkId, siteId, status string, page, pageSize uint32) (*netpb.GetNodesResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), n.timeout)
	defer cancel()

	return n.client.GetNodes(ctx, &netpb.GetNodesRequest{
		NetworkId: networkId,
		SiteId:    siteId,
		Status:    status,
		Page:      page,
		PageSize:  pageSize,
	})
}

func (n *NetworkAnalytics) GetNode(nodeId, period, from, to, tz string) (*netpb.GetNodeResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), n.timeout)
	defer cancel()

	return n.client.GetNode(ctx, &netpb.GetNodeRequest{
		NodeId: nodeId,
		Window: n.window(period, from, to, tz),
	})
}

func (n *NetworkAnalytics) GetNodePool(networkId string) (*netpb.GetNodePoolResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), n.timeout)
	defer cancel()

	return n.client.GetNodePool(ctx, &netpb.GetNodePoolRequest{
		NetworkId: networkId,
	})
}

func (n *NetworkAnalytics) GetRadio(networkId, siteId, nodeId, period, from, to, tz string) (*netpb.GetRadioResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), n.timeout)
	defer cancel()

	return n.client.GetRadio(ctx, &netpb.GetRadioRequest{
		NetworkId: networkId,
		SiteId:    siteId,
		NodeId:    nodeId,
		Window:    n.window(period, from, to, tz),
	})
}

func (n *NetworkAnalytics) GetBackhaul(networkId, siteId, period, from, to, tz string) (*netpb.GetBackhaulResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), n.timeout)
	defer cancel()

	return n.client.GetBackhaul(ctx, &netpb.GetBackhaulRequest{
		NetworkId: networkId,
		SiteId:    siteId,
		Window:    n.window(period, from, to, tz),
	})
}

func (n *NetworkAnalytics) GetPower(networkId, siteId, period, from, to, tz string) (*netpb.GetPowerResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), n.timeout)
	defer cancel()

	return n.client.GetPower(ctx, &netpb.GetPowerRequest{
		NetworkId: networkId,
		SiteId:    siteId,
		Window:    n.window(period, from, to, tz),
	})
}

func (n *NetworkAnalytics) GetAlarms(networkId, siteId, severity, state string, page, pageSize uint32) (*netpb.GetAlarmsResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), n.timeout)
	defer cancel()

	return n.client.GetAlarms(ctx, &netpb.GetAlarmsRequest{
		NetworkId: networkId,
		SiteId:    siteId,
		Severity:  severity,
		State:     state,
		Page:      page,
		PageSize:  pageSize,
	})
}

func (n *NetworkAnalytics) GetMetrics(networkId, siteId, nodeId, metric, period, from, to, tz string) (*netpb.GetMetricsResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), n.timeout)
	defer cancel()

	return n.client.GetMetrics(ctx, &netpb.GetMetricsRequest{
		NetworkId: networkId,
		SiteId:    siteId,
		NodeId:    nodeId,
		Metric:    metric,
		Window:    n.window(period, from, to, tz),
	})
}

func (n *NetworkAnalytics) GetEvents(networkId, siteId, nodeId, period, from, to, tz string, page, pageSize uint32) (*netpb.GetEventsResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), n.timeout)
	defer cancel()

	return n.client.GetEvents(ctx, &netpb.GetEventsRequest{
		NetworkId: networkId,
		SiteId:    siteId,
		NodeId:    nodeId,
		Window:    n.window(period, from, to, tz),
		Page:      page,
		PageSize:  pageSize,
	})
}

func (n *NetworkAnalytics) SupportSearch(query, networkId string) (*netpb.SupportSearchResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), n.timeout)
	defer cancel()

	return n.client.SupportSearch(ctx, &netpb.SupportSearchRequest{
		Query:     query,
		NetworkId: networkId,
	})
}
