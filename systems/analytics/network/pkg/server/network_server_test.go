/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package server

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/tj/assert"

	"github.com/ukama/ukama/systems/analytics/network/mocks"
	pb "github.com/ukama/ukama/systems/analytics/network/pb/gen"
	"github.com/ukama/ukama/systems/analytics/network/pkg/db"
	"github.com/ukama/ukama/systems/common/uuid"
)

type netMocks struct {
	site   *mocks.SiteRepo
	node   *mocks.NodeRepo
	alarm  *mocks.AlarmRepo
	metric *mocks.MetricRepo
	event  *mocks.EventRepo
	health *mocks.HealthRepo
}

func anyN(n int) []interface{} {
	a := make([]interface{}, n)
	for i := range a {
		a[i] = mock.Anything
	}
	return a
}

func sampleSiteSnap() db.SiteSnapshot {
	return db.SiteSnapshot{SiteId: uuid.NewV4(), NetworkId: uuid.NewV4(), Name: "Site One", Status: "online", NodeCount: 2, Latitude: 1, Longitude: 2}
}

func sampleNodeSnap() db.NodeSnapshot {
	now := time.Now()
	return db.NodeSnapshot{NodeId: "node-1", Name: "Node One", Type: "tower", Status: "online", SiteId: uuid.NewV4(), NetworkId: uuid.NewV4(), LastTelemetryAt: &now}
}

// wireHappy stubs every repo method with success defaults so any RPC can run.
func newNetMocks() *netMocks {
	m := &netMocks{
		site:   &mocks.SiteRepo{},
		node:   &mocks.NodeRepo{},
		alarm:  &mocks.AlarmRepo{},
		metric: &mocks.MetricRepo{},
		event:  &mocks.EventRepo{},
		health: &mocks.HealthRepo{},
	}

	site := sampleSiteSnap()
	node := sampleNodeSnap()
	now := time.Now()

	m.site.On("List", anyN(4)...).Return([]db.SiteSnapshot{site}, int64(1), nil).Maybe()
	m.site.On("StatusCounts", anyN(1)...).Return(&db.SiteStatusCounts{Total: 5, Online: 4, Degraded: 1, Offline: 0}, nil).Maybe()
	m.site.On("Get", anyN(1)...).Return(&site, nil).Maybe()
	m.site.On("CustomerCount", anyN(1)...).Return(int64(12), nil).Maybe()
	m.site.On("UptimeBetween", anyN(3)...).Return(99.0, nil).Maybe()
	m.site.On("LatestSiteHealth", anyN(1)...).Return(&db.SiteHealthRollupHourly{Hour: now, UptimePercent: 99}, nil).Maybe()
	m.site.On("SiteHealthSeries", anyN(3)...).Return([]db.SiteHealthRollupHourly{{Hour: now, UptimePercent: 99, BackhaulLatencyMs: 20, BatteryPercent: 80}}, nil).Maybe()
	m.site.On("Search", anyN(3)...).Return([]db.SiteSnapshot{site}, nil).Maybe()
	m.site.On("OfflineDuration", anyN(1)...).Return(0.0, nil).Maybe()

	m.node.On("List", anyN(5)...).Return([]db.NodeSnapshot{node}, int64(1), nil).Maybe()
	m.node.On("ListAll", anyN(1)...).Return([]db.NodeSnapshot{node}, nil).Maybe()
	m.node.On("StatusCounts", anyN(2)...).Return(&db.NodeStatusCounts{Total: 3, Online: 2, Offline: 1}, nil).Maybe()
	m.node.On("Get", anyN(1)...).Return(&node, nil).Maybe()
	m.node.On("UptimeBetween", anyN(3)...).Return(98.0, nil).Maybe()
	m.node.On("PoolCounts", anyN(1)...).Return(&db.NodePoolCounts{AvailableToInstall: 2, Deployed: 5, InInventory: 1, Rma: 0}, nil).Maybe()
	m.node.On("ConfiguringDuration", anyN(1)...).Return(0.0, nil).Maybe()
	m.node.On("Search", anyN(3)...).Return([]db.NodeSnapshot{node}, nil).Maybe()

	m.alarm.On("List", anyN(1)...).Return([]db.AlarmEvent{{AlarmId: "a1", Severity: "critical", State: "open", ResourceType: "site", OpenedAt: now}}, int64(1), nil).Maybe()
	m.alarm.On("Counts", anyN(2)...).Return(&db.AlarmCounts{Open: 3, Critical: 1, Warning: 2}, nil).Maybe()
	m.alarm.On("ForResource", anyN(3)...).Return([]db.AlarmEvent{{AlarmId: "a1", Severity: "warning", State: "open", OpenedAt: now}}, nil).Maybe()
	m.alarm.On("OpenImpact", anyN(1)...).Return(int64(4), 120.0, nil).Maybe()

	m.metric.On("Rollups", anyN(5)...).Return([]db.MetricRollupHourly{{Hour: now, Avg: 1.0}}, nil).Maybe()
	m.metric.On("LatestSamples", anyN(2)...).Return([]db.MetricSample{{Metric: "x", Value: 1, SampledAt: now}}, nil).Maybe()
	m.metric.On("MetricNames", anyN(1)...).Return([]db.MetricName{{Metric: "x", Unit: "u", LastSampleAt: now}}, nil).Maybe()
	m.metric.On("RadioRollups", anyN(3)...).Return([]db.RadioRollupHourly{{Hour: now, ActiveUes: 10, DlThroughputMbps: 5, UlThroughputMbps: 2, SignalDbm: -90, AttachFailures: 1}}, nil).Maybe()
	m.metric.On("LatestRadioRollup", anyN(1)...).Return(&db.RadioRollupHourly{Hour: now, ActiveUes: 10, SignalDbm: -90}, nil).Maybe()
	m.metric.On("RadioRollupSums", anyN(2)...).Return(int64(10), int64(1), nil).Maybe()
	m.metric.On("BackhaulRollups", anyN(3)...).Return([]db.BackhaulRollupHourly{{Hour: now, LatencyMs: 20, DlMbps: 50, UlMbps: 10, PacketLossPercent: 0.5}}, nil).Maybe()
	m.metric.On("LatestBackhaulRollup", anyN(1)...).Return(&db.BackhaulRollupHourly{Hour: now, LatencyMs: 20}, nil).Maybe()
	m.metric.On("PowerRollups", anyN(3)...).Return([]db.PowerRollupHourly{{Hour: now, BatteryPercent: 80, BatteryVoltage: 12, LoadWatts: 100, SolarWatts: 200, TemperatureC: 25}}, nil).Maybe()
	m.metric.On("LatestPowerRollup", anyN(1)...).Return(&db.PowerRollupHourly{Hour: now, BatteryPercent: 80, BatteryVoltage: 12, LoadWatts: 100}, nil).Maybe()

	m.event.On("Recent", anyN(7)...).Return([]db.EventLog{{RoutingKey: "event.x", OccurredAt: now}}, int64(1), nil).Maybe()
	m.health.On("NetworkHealthLatest", anyN(1)...).Return(&db.NetworkHealthRollupHourly{Hour: now}, nil).Maybe()
	m.health.On("NetworkHealthSeries", anyN(3)...).Return([]db.NetworkHealthRollupHourly{{Hour: now}}, nil).Maybe()

	return m
}

func (m *netMocks) server() *NetworkServer {
	return NewNetworkServer("test-org", m.site, m.node, m.alarm, m.metric, m.event, m.health, 100, 20, 300)
}

func ctx() context.Context { return context.Background() }

func TestNet_GetOverview(t *testing.T) {
	resp, err := newNetMocks().server().GetOverview(ctx(), &pb.GetOverviewRequest{NetworkId: "net-1"})
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.Kpis)
}

func TestNet_GetTopology(t *testing.T) {
	resp, err := newNetMocks().server().GetTopology(ctx(), &pb.GetTopologyRequest{NetworkId: "net-1"})
	assert.NoError(t, err)
	assert.Len(t, resp.Sites, 1)
}

func TestNet_GetSites(t *testing.T) {
	resp, err := newNetMocks().server().GetSites(ctx(), &pb.GetSitesRequest{NetworkId: "net-1", Page: 1, PageSize: 10})
	assert.NoError(t, err)
	assert.Len(t, resp.Sites, 1)
}

func TestNet_GetSite(t *testing.T) {
	resp, err := newNetMocks().server().GetSite(ctx(), &pb.GetSiteRequest{SiteId: uuid.NewV4().String()})
	assert.NoError(t, err)
	assert.NotNil(t, resp.Site)
}

func TestNet_GetSite_MissingId(t *testing.T) {
	_, err := newNetMocks().server().GetSite(ctx(), &pb.GetSiteRequest{SiteId: ""})
	assert.Error(t, err)
}

func TestNet_GetNodes(t *testing.T) {
	resp, err := newNetMocks().server().GetNodes(ctx(), &pb.GetNodesRequest{NetworkId: "net-1", Page: 1, PageSize: 10})
	assert.NoError(t, err)
	assert.Len(t, resp.Nodes, 1)
}

func TestNet_GetNode(t *testing.T) {
	resp, err := newNetMocks().server().GetNode(ctx(), &pb.GetNodeRequest{NodeId: "node-1"})
	assert.NoError(t, err)
	assert.NotNil(t, resp.Node)
}

func TestNet_GetNode_MissingId(t *testing.T) {
	_, err := newNetMocks().server().GetNode(ctx(), &pb.GetNodeRequest{NodeId: ""})
	assert.Error(t, err)
}

func TestNet_GetNodePool(t *testing.T) {
	resp, err := newNetMocks().server().GetNodePool(ctx(), &pb.GetNodePoolRequest{NetworkId: "net-1"})
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.Kpis)
}

func TestNet_GetRadio_Node(t *testing.T) {
	resp, err := newNetMocks().server().GetRadio(ctx(), &pb.GetRadioRequest{NodeId: "node-1"})
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.Series)
}

func TestNet_GetRadio_Network(t *testing.T) {
	resp, err := newNetMocks().server().GetRadio(ctx(), &pb.GetRadioRequest{NetworkId: "net-1"})
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.Kpis)
}

func TestNet_GetBackhaul(t *testing.T) {
	resp, err := newNetMocks().server().GetBackhaul(ctx(), &pb.GetBackhaulRequest{SiteId: "site-1"})
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.Series)
}

func TestNet_GetPower(t *testing.T) {
	resp, err := newNetMocks().server().GetPower(ctx(), &pb.GetPowerRequest{SiteId: "site-1"})
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.Series)
}

func TestNet_GetAlarms(t *testing.T) {
	resp, err := newNetMocks().server().GetAlarms(ctx(), &pb.GetAlarmsRequest{NetworkId: "net-1", Page: 1, PageSize: 10})
	assert.NoError(t, err)
	assert.Len(t, resp.Alarms, 1)
}

func TestNet_GetMetrics(t *testing.T) {
	resp, err := newNetMocks().server().GetMetrics(ctx(), &pb.GetMetricsRequest{NodeId: "node-1", Metric: "x"})
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.Metrics)
}

func TestNet_GetEvents(t *testing.T) {
	resp, err := newNetMocks().server().GetEvents(ctx(), &pb.GetEventsRequest{NetworkId: "net-1", Page: 1, PageSize: 10})
	assert.NoError(t, err)
	assert.Len(t, resp.Events, 1)
}

func TestNet_SupportSearch(t *testing.T) {
	resp, err := newNetMocks().server().SupportSearch(ctx(), &pb.SupportSearchRequest{Query: "site", NetworkId: "net-1"})
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.Results)
}

func TestNet_SupportSearch_EmptyQuery(t *testing.T) {
	_, err := newNetMocks().server().SupportSearch(ctx(), &pb.SupportSearchRequest{Query: ""})
	assert.Error(t, err)
}

func TestEstimateRuntimeHours(t *testing.T) {
	assert.Equal(t, 0.0, estimateRuntimeHours(80, 0, 100))
	assert.Equal(t, 0.0, estimateRuntimeHours(80, 12, 0))
	assert.InDelta(t, 9.6, estimateRuntimeHours(80, 12, 100), 0.1)
}
