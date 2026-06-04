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
	"errors"
	"fmt"
	"strings"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/grpc"

	pb "github.com/ukama/ukama/systems/analytics/network/pb/gen"
	"github.com/ukama/ukama/systems/analytics/network/pkg/db"
)

const (
	recentEventsLimit  = 10
	resourceAlarmLimit = 20
	supportSearchLimit = 10
)

type NetworkServer struct {
	pb.UnimplementedNetworkServiceServer
	orgName                string
	siteRepo               db.SiteRepo
	nodeRepo               db.NodeRepo
	alarmRepo              db.AlarmRepo
	metricRepo             db.MetricRepo
	eventRepo              db.EventRepo
	healthRepo             db.HealthRepo
	latencyThresholdMs     float64
	batteryCriticalPercent float64
	telemetryFreshSeconds  int64
}

func NewNetworkServer(orgName string, siteRepo db.SiteRepo, nodeRepo db.NodeRepo,
	alarmRepo db.AlarmRepo, metricRepo db.MetricRepo, eventRepo db.EventRepo,
	healthRepo db.HealthRepo,
	latencyThresholdMs, batteryCriticalPercent float64, telemetryFreshSeconds int64) *NetworkServer {
	return &NetworkServer{
		orgName:    orgName,
		siteRepo:   siteRepo,
		nodeRepo:   nodeRepo,
		alarmRepo:  alarmRepo,
		metricRepo: metricRepo,
		eventRepo:  eventRepo,
		healthRepo: healthRepo,

		latencyThresholdMs:     latencyThresholdMs,
		batteryCriticalPercent: batteryCriticalPercent,
		telemetryFreshSeconds:  telemetryFreshSeconds,
	}
}

func (n *NetworkServer) GetOverview(ctx context.Context, req *pb.GetOverviewRequest) (*pb.GetOverviewResponse, error) {
	log.Infof("GetOverview for network %s", req.GetNetworkId())

	w := ResolveWindow(req.GetWindow(), time.Now().UTC())

	siteCounts, err := n.siteRepo.StatusCounts(req.GetNetworkId())
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "sites")
	}

	nodeCounts, err := n.nodeRepo.StatusCounts(req.GetNetworkId(), "")
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "nodes")
	}

	alarmCounts, err := n.alarmRepo.Counts(req.GetNetworkId(), "")
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "alarms")
	}

	customersAffected, revenueAtRisk, err := n.alarmRepo.OpenImpact(req.GetNetworkId())
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "alarms")
	}

	activeUes, _, err := n.metricRepo.RadioRollupSums(req.GetNetworkId(), w.To)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "radio")
	}

	events, _, err := n.eventRepo.Recent(req.GetNetworkId(), "", "",
		w.From, w.To, 1, recentEventsLimit)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "events")
	}

	networkStatus := DeriveNetworkStatus(alarmCounts.Critical, alarmCounts.Open,
		siteCounts.Total, siteCounts.Degraded, siteCounts.Offline)

	kpis := []*pb.Kpi{
		ratioKpi("sites_online", siteCounts.Online, siteCounts.Total),
		ratioKpi("nodes_online", nodeCounts.Online, nodeCounts.Total),
		countKpi("active_ues", activeUes),
		countKpi("open_alarms", alarmCounts.Open),
		countKpi("critical_alarms", alarmCounts.Critical),
		countKpi("customers_affected", customersAffected),
		moneyKpi("revenue_at_risk", revenueAtRisk),
	}

	return &pb.GetOverviewResponse{
		NetworkStatus: networkStatus,
		Kpis:          kpis,
		RecentEvents:  pbEvents(events),
	}, nil
}

func (n *NetworkServer) GetTopology(ctx context.Context, req *pb.GetTopologyRequest) (*pb.GetTopologyResponse, error) {
	log.Infof("GetTopology for network %s", req.GetNetworkId())

	sites, _, err := n.siteRepo.List(req.GetNetworkId(), "", 0, 0)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "sites")
	}

	nodes, err := n.nodeRepo.ListAll(req.GetNetworkId())
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "nodes")
	}

	nodesBySite := make(map[string][]*pb.TopologyNode)
	for i := range nodes {
		nd := &nodes[i]
		nodesBySite[nd.SiteId.String()] = append(nodesBySite[nd.SiteId.String()], &pb.TopologyNode{
			NodeId: nd.NodeId,
			Name:   nd.Name,
			Type:   nd.Type,
			Status: nd.Status,
		})
	}

	resp := &pb.GetTopologyResponse{}
	for i := range sites {
		s := &sites[i]
		resp.Sites = append(resp.Sites, &pb.TopologySite{
			SiteId:    s.SiteId.String(),
			Name:      s.Name,
			Status:    s.Status,
			Latitude:  s.Latitude,
			Longitude: s.Longitude,
			Nodes:     nodesBySite[s.SiteId.String()],
		})
	}

	return resp, nil
}

// siteRow builds a SiteRow for a snapshot, enriching it with uptime, customer
// count and threshold-derived flags.
func (n *NetworkServer) siteRow(s *db.SiteSnapshot, w TimeWindow) *pb.SiteRow {
	row := &pb.SiteRow{
		SiteId:    s.SiteId.String(),
		Name:      s.Name,
		Status:    s.Status,
		NodeCount: s.NodeCount,
		Latitude:  s.Latitude,
		Longitude: s.Longitude,
	}

	if customers, err := n.siteRepo.CustomerCount(s.SiteId.String()); err == nil {
		row.Customers = uint32(customers)
	}

	if uptime, err := n.siteRepo.UptimeBetween(s.SiteId.String(), w.From, w.To); err == nil {
		row.Uptime = uptime
	}

	if offline, err := n.siteRepo.OfflineDuration(s.SiteId.String()); err == nil {
		row.OfflineDurationSeconds = offline
	}

	issues := []string{}

	if bh, err := n.metricRepo.LatestBackhaulRollup(s.SiteId.String()); err == nil {
		if bh.LatencyMs > n.latencyThresholdMs {
			row.BackhaulLatencyHigh = true

			issues = append(issues, "backhaul latency high")
		}
	}

	if pw, err := n.metricRepo.LatestPowerRollup(s.SiteId.String()); err == nil {
		if pw.BatteryPercent < n.batteryCriticalPercent {
			row.BatteryCritical = true

			issues = append(issues, "battery critical")
		}
	}

	switch s.Status {
	case "offline":
		issues = append(issues, "site offline")
	case "degraded":
		issues = append(issues, "site degraded")
	}

	row.IssueSummary = strings.Join(issues, "; ")

	return row
}

func (n *NetworkServer) GetSites(ctx context.Context, req *pb.GetSitesRequest) (*pb.GetSitesResponse, error) {
	log.Infof("GetSites for network %s", req.GetNetworkId())

	w := ResolveWindow(req.GetWindow(), time.Now().UTC())

	counts, err := n.siteRepo.StatusCounts(req.GetNetworkId())
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "sites")
	}

	sites, total, err := n.siteRepo.List(req.GetNetworkId(), req.GetStatus(),
		req.GetPage(), req.GetPageSize())
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "sites")
	}

	rows := make([]*pb.SiteRow, 0, len(sites))
	for i := range sites {
		rows = append(rows, n.siteRow(&sites[i], w))
	}

	return &pb.GetSitesResponse{
		Kpis: []*pb.Kpi{
			countKpi("total", counts.Total),
			countKpi("online", counts.Online),
			countKpi("degraded", counts.Degraded),
			countKpi("offline", counts.Offline),
		},
		Sites: rows,
		Meta:  pbMeta(total, req.GetPage(), req.GetPageSize()),
	}, nil
}

func (n *NetworkServer) GetSite(ctx context.Context, req *pb.GetSiteRequest) (*pb.GetSiteResponse, error) {
	log.Infof("GetSite %s", req.GetSiteId())

	if req.GetSiteId() == "" {
		return nil, status.Error(codes.InvalidArgument, "site_id is required")
	}

	w := ResolveWindow(req.GetWindow(), time.Now().UTC())

	site, err := n.siteRepo.Get(req.GetSiteId())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "site %s not found", req.GetSiteId())
		}

		return nil, grpc.SqlErrorToGrpc(err, "site")
	}

	row := n.siteRow(site, w)

	health, err := n.siteRepo.SiteHealthSeries(req.GetSiteId(), w.From, w.To)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "site health")
	}

	uptimeSeries := &pb.TimeSeries{Key: "uptime"}
	latencySeries := &pb.TimeSeries{Key: "backhaul_latency"}
	batterySeries := &pb.TimeSeries{Key: "battery"}

	for i := range health {
		h := &health[i]
		uptimeSeries.Points = append(uptimeSeries.Points, point(h.Hour, h.UptimePercent))
		latencySeries.Points = append(latencySeries.Points, point(h.Hour, h.BackhaulLatencyMs))
		batterySeries.Points = append(batterySeries.Points, point(h.Hour, h.BatteryPercent))
	}

	alarms, err := n.alarmRepo.ForResource("", req.GetSiteId(), resourceAlarmLimit)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "alarms")
	}

	kpis := []*pb.Kpi{
		percentKpi("uptime", row.Uptime),
		countKpi("customers", int64(row.Customers)),
		countKpi("nodes", int64(row.NodeCount)),
	}

	return &pb.GetSiteResponse{
		Site:   row,
		Kpis:   kpis,
		Series: []*pb.TimeSeries{uptimeSeries, latencySeries, batterySeries},
		Alarms: pbAlarms(alarms),
	}, nil
}

// nodeRow builds a NodeRow for a snapshot, enriching it with uptime, last
// telemetry freshness and configuring duration.
func (n *NetworkServer) nodeRow(nd *db.NodeSnapshot, w TimeWindow, siteNames map[string]string) *pb.NodeRow {
	row := &pb.NodeRow{
		NodeId: nd.NodeId,
		Name:   nd.Name,
		Type:   nd.Type,
		Status: nd.Status,
		SiteId: nd.SiteId.String(),
	}

	if siteNames != nil {
		row.SiteName = siteNames[nd.SiteId.String()]
	}

	if uptime, err := n.nodeRepo.UptimeBetween(nd.NodeId, w.From, w.To); err == nil {
		row.Uptime = uptime
	}

	if nd.LastTelemetryAt != nil {
		row.LastTelemetry = timestamppb.New(*nd.LastTelemetryAt)

		if time.Since(*nd.LastTelemetryAt) > time.Duration(n.telemetryFreshSeconds)*time.Second {
			row.NoTelemetryWarning = true
		}
	} else {
		row.NoTelemetryWarning = true
	}

	if nd.Status == "configuring" {
		if d, err := n.nodeRepo.ConfiguringDuration(nd.NodeId); err == nil {
			row.ConfiguringDurationSeconds = d
		}
	}

	return row
}

// siteNameMap returns a site_id -> name map for the network.
func (n *NetworkServer) siteNameMap(networkId string) map[string]string {
	names := make(map[string]string)

	sites, _, err := n.siteRepo.List(networkId, "", 0, 0)
	if err != nil {
		log.Warnf("failed to list sites for name map: %v", err)

		return names
	}

	for i := range sites {
		names[sites[i].SiteId.String()] = sites[i].Name
	}

	return names
}

func (n *NetworkServer) GetNodes(ctx context.Context, req *pb.GetNodesRequest) (*pb.GetNodesResponse, error) {
	log.Infof("GetNodes for network %s site %s", req.GetNetworkId(), req.GetSiteId())

	w := ResolveWindow(nil, time.Now().UTC())

	counts, err := n.nodeRepo.StatusCounts(req.GetNetworkId(), req.GetSiteId())
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "nodes")
	}

	nodes, total, err := n.nodeRepo.List(req.GetNetworkId(), req.GetSiteId(),
		req.GetStatus(), req.GetPage(), req.GetPageSize())
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "nodes")
	}

	siteNames := n.siteNameMap(req.GetNetworkId())

	rows := make([]*pb.NodeRow, 0, len(nodes))
	for i := range nodes {
		rows = append(rows, n.nodeRow(&nodes[i], w, siteNames))
	}

	return &pb.GetNodesResponse{
		Kpis: []*pb.Kpi{
			countKpi("total", counts.Total),
			countKpi("online", counts.Online),
			countKpi("offline", counts.Offline),
			countKpi("needs_attention", counts.NeedsAttention),
		},
		Nodes: rows,
		Meta:  pbMeta(total, req.GetPage(), req.GetPageSize()),
	}, nil
}

func (n *NetworkServer) GetNode(ctx context.Context, req *pb.GetNodeRequest) (*pb.GetNodeResponse, error) {
	log.Infof("GetNode %s", req.GetNodeId())

	if req.GetNodeId() == "" {
		return nil, status.Error(codes.InvalidArgument, "node_id is required")
	}

	w := ResolveWindow(req.GetWindow(), time.Now().UTC())

	node, err := n.nodeRepo.Get(req.GetNodeId())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "node %s not found", req.GetNodeId())
		}

		return nil, grpc.SqlErrorToGrpc(err, "node")
	}

	row := n.nodeRow(node, w, n.siteNameMap(node.NetworkId.String()))

	radio, err := n.metricRepo.RadioRollups(req.GetNodeId(), w.From, w.To)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "radio")
	}

	uesSeries := &pb.TimeSeries{Key: "active_ues"}
	dlSeries := &pb.TimeSeries{Key: "dl_throughput"}
	ulSeries := &pb.TimeSeries{Key: "ul_throughput"}
	signalSeries := &pb.TimeSeries{Key: "signal_dbm"}

	for i := range radio {
		rr := &radio[i]
		uesSeries.Points = append(uesSeries.Points, point(rr.Hour, float64(rr.ActiveUes)))
		dlSeries.Points = append(dlSeries.Points, point(rr.Hour, rr.DlThroughputMbps))
		ulSeries.Points = append(ulSeries.Points, point(rr.Hour, rr.UlThroughputMbps))
		signalSeries.Points = append(signalSeries.Points, point(rr.Hour, rr.SignalDbm))
	}

	events, _, err := n.eventRepo.Recent("", "", req.GetNodeId(), w.From, w.To, 1, recentEventsLimit)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "events")
	}

	kpis := []*pb.Kpi{
		percentKpi("uptime", row.Uptime),
	}

	if latest, err := n.metricRepo.LatestRadioRollup(req.GetNodeId()); err == nil {
		kpis = append(kpis,
			countKpi("active_ues", int64(latest.ActiveUes)),
			floatKpi("dl_throughput", latest.DlThroughputMbps, "Mbps"),
			floatKpi("ul_throughput", latest.UlThroughputMbps, "Mbps"),
			floatKpi("signal_quality", latest.SignalDbm, "dBm"),
		)
	}

	return &pb.GetNodeResponse{
		Node:         row,
		Kpis:         kpis,
		Series:       []*pb.TimeSeries{uesSeries, dlSeries, ulSeries, signalSeries},
		RecentEvents: pbEvents(events),
	}, nil
}

func (n *NetworkServer) GetNodePool(ctx context.Context, req *pb.GetNodePoolRequest) (*pb.GetNodePoolResponse, error) {
	log.Infof("GetNodePool for network %s", req.GetNetworkId())

	w := ResolveWindow(nil, time.Now().UTC())

	pool, err := n.nodeRepo.PoolCounts(req.GetNetworkId())
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "node pool")
	}

	nodes, err := n.nodeRepo.ListAll(req.GetNetworkId())
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "nodes")
	}

	siteNames := n.siteNameMap(req.GetNetworkId())

	rows := make([]*pb.NodeRow, 0, len(nodes))
	for i := range nodes {
		rows = append(rows, n.nodeRow(&nodes[i], w, siteNames))
	}

	return &pb.GetNodePoolResponse{
		Kpis: []*pb.Kpi{
			countKpi("available_to_install", pool.AvailableToInstall),
			countKpi("deployed", pool.Deployed),
			countKpi("in_inventory", pool.InInventory),
			countKpi("rma", pool.Rma),
		},
		Nodes: rows,
	}, nil
}

func (n *NetworkServer) GetRadio(ctx context.Context, req *pb.GetRadioRequest) (*pb.GetRadioResponse, error) {
	log.Infof("GetRadio network %s site %s node %s", req.GetNetworkId(), req.GetSiteId(), req.GetNodeId())

	w := ResolveWindow(req.GetWindow(), time.Now().UTC())

	rollups, err := n.metricRepo.RadioRollups(req.GetNodeId(), w.From, w.To)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "radio")
	}

	uesSeries := &pb.TimeSeries{Key: "active_ues"}
	failSeries := &pb.TimeSeries{Key: "attach_failures"}
	dlSeries := &pb.TimeSeries{Key: "dl_throughput"}
	ulSeries := &pb.TimeSeries{Key: "ul_throughput"}
	signalSeries := &pb.TimeSeries{Key: "signal_quality"}

	for i := range rollups {
		rr := &rollups[i]
		uesSeries.Points = append(uesSeries.Points, point(rr.Hour, float64(rr.ActiveUes)))
		failSeries.Points = append(failSeries.Points, point(rr.Hour, float64(rr.AttachFailures)))
		dlSeries.Points = append(dlSeries.Points, point(rr.Hour, rr.DlThroughputMbps))
		ulSeries.Points = append(ulSeries.Points, point(rr.Hour, rr.UlThroughputMbps))
		signalSeries.Points = append(signalSeries.Points, point(rr.Hour, rr.SignalDbm))
	}

	kpis := []*pb.Kpi{}

	if req.GetNodeId() != "" {
		if latest, err := n.metricRepo.LatestRadioRollup(req.GetNodeId()); err == nil {
			kpis = append(kpis,
				countKpi("active_ues", int64(latest.ActiveUes)),
				countKpi("attach_failures", int64(latest.AttachFailures)),
				floatKpi("dl_throughput", latest.DlThroughputMbps, "Mbps"),
				floatKpi("ul_throughput", latest.UlThroughputMbps, "Mbps"),
				floatKpi("signal_quality", latest.SignalDbm, "dBm"),
			)
		}
	} else {
		activeUes, attachFailures, err := n.metricRepo.RadioRollupSums(req.GetNetworkId(), w.To)
		if err != nil {
			return nil, grpc.SqlErrorToGrpc(err, "radio")
		}

		kpis = append(kpis,
			countKpi("active_ues", activeUes),
			countKpi("attach_failures", attachFailures),
		)
	}

	alarms, err := n.alarmRepo.ForResource("radio", "", resourceAlarmLimit)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "alarms")
	}

	return &pb.GetRadioResponse{
		Kpis:   kpis,
		Series: []*pb.TimeSeries{uesSeries, failSeries, dlSeries, ulSeries, signalSeries},
		Alarms: pbAlarms(alarms),
	}, nil
}

func (n *NetworkServer) GetBackhaul(ctx context.Context, req *pb.GetBackhaulRequest) (*pb.GetBackhaulResponse, error) {
	log.Infof("GetBackhaul network %s site %s", req.GetNetworkId(), req.GetSiteId())

	w := ResolveWindow(req.GetWindow(), time.Now().UTC())

	rollups, err := n.metricRepo.BackhaulRollups(req.GetSiteId(), w.From, w.To)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "backhaul")
	}

	latencySeries := &pb.TimeSeries{Key: "latency"}
	dlSeries := &pb.TimeSeries{Key: "dl_throughput"}
	ulSeries := &pb.TimeSeries{Key: "ul_throughput"}
	lossSeries := &pb.TimeSeries{Key: "packet_loss"}

	for i := range rollups {
		rr := &rollups[i]
		latencySeries.Points = append(latencySeries.Points, point(rr.Hour, rr.LatencyMs))
		dlSeries.Points = append(dlSeries.Points, point(rr.Hour, rr.DlMbps))
		ulSeries.Points = append(ulSeries.Points, point(rr.Hour, rr.UlMbps))
		lossSeries.Points = append(lossSeries.Points, point(rr.Hour, rr.PacketLossPercent))
	}

	kpis := []*pb.Kpi{}

	if req.GetSiteId() != "" {
		if latest, err := n.metricRepo.LatestBackhaulRollup(req.GetSiteId()); err == nil {
			kpis = append(kpis,
				floatKpi("latency", latest.LatencyMs, "ms"),
				floatKpi("dl_throughput", latest.DlMbps, "Mbps"),
				floatKpi("ul_throughput", latest.UlMbps, "Mbps"),
				percentKpi("packet_loss", latest.PacketLossPercent),
			)
		}
	}

	alarms, err := n.alarmRepo.ForResource("backhaul", "", resourceAlarmLimit)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "alarms")
	}

	kpis = append(kpis, countKpi("recent_failures", int64(len(alarms))))

	return &pb.GetBackhaulResponse{
		Kpis:   kpis,
		Series: []*pb.TimeSeries{latencySeries, dlSeries, ulSeries, lossSeries},
		Alarms: pbAlarms(alarms),
	}, nil
}

func (n *NetworkServer) GetPower(ctx context.Context, req *pb.GetPowerRequest) (*pb.GetPowerResponse, error) {
	log.Infof("GetPower network %s site %s", req.GetNetworkId(), req.GetSiteId())

	w := ResolveWindow(req.GetWindow(), time.Now().UTC())

	rollups, err := n.metricRepo.PowerRollups(req.GetSiteId(), w.From, w.To)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "power")
	}

	batterySeries := &pb.TimeSeries{Key: "battery_percent"}
	voltageSeries := &pb.TimeSeries{Key: "battery_voltage"}
	loadSeries := &pb.TimeSeries{Key: "load_watts"}
	solarSeries := &pb.TimeSeries{Key: "solar_watts"}
	tempSeries := &pb.TimeSeries{Key: "temperature"}

	for i := range rollups {
		rr := &rollups[i]
		batterySeries.Points = append(batterySeries.Points, point(rr.Hour, rr.BatteryPercent))
		voltageSeries.Points = append(voltageSeries.Points, point(rr.Hour, rr.BatteryVoltage))
		loadSeries.Points = append(loadSeries.Points, point(rr.Hour, rr.LoadWatts))
		solarSeries.Points = append(solarSeries.Points, point(rr.Hour, rr.SolarWatts))
		tempSeries.Points = append(tempSeries.Points, point(rr.Hour, rr.TemperatureC))
	}

	kpis := []*pb.Kpi{}

	if req.GetSiteId() != "" {
		if latest, err := n.metricRepo.LatestPowerRollup(req.GetSiteId()); err == nil {
			kpis = append(kpis,
				percentKpi("battery_percent", latest.BatteryPercent),
				floatKpi("battery_voltage", latest.BatteryVoltage, "V"),
				floatKpi("load_watts", latest.LoadWatts, "W"),
				floatKpi("solar_watts", latest.SolarWatts, "W"),
				floatKpi("temperature", latest.TemperatureC, "C"),
				floatKpi("runtime_estimate", estimateRuntimeHours(latest.BatteryPercent, latest.BatteryVoltage, latest.LoadWatts), "h"),
			)
		}
	}

	alarms, err := n.alarmRepo.ForResource("power", "", resourceAlarmLimit)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "alarms")
	}

	return &pb.GetPowerResponse{
		Kpis:   kpis,
		Series: []*pb.TimeSeries{batterySeries, voltageSeries, loadSeries, solarSeries, tempSeries},
		Alarms: pbAlarms(alarms),
	}, nil
}

func (n *NetworkServer) GetAlarms(ctx context.Context, req *pb.GetAlarmsRequest) (*pb.GetAlarmsResponse, error) {
	log.Infof("GetAlarms network %s site %s", req.GetNetworkId(), req.GetSiteId())

	counts, err := n.alarmRepo.Counts(req.GetNetworkId(), req.GetSiteId())
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "alarms")
	}

	alarms, total, err := n.alarmRepo.List(db.AlarmFilter{
		NetworkId: req.GetNetworkId(),
		SiteId:    req.GetSiteId(),
		Severity:  req.GetSeverity(),
		State:     req.GetState(),
		Page:      req.GetPage(),
		PageSize:  req.GetPageSize(),
	})
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "alarms")
	}

	return &pb.GetAlarmsResponse{
		Kpis: []*pb.Kpi{
			countKpi("open", counts.Open),
			countKpi("critical", counts.Critical),
			countKpi("warning", counts.Warning),
		},
		Alarms: pbAlarms(alarms),
		Meta:   pbMeta(total, req.GetPage(), req.GetPageSize()),
	}, nil
}

func (n *NetworkServer) GetMetrics(ctx context.Context, req *pb.GetMetricsRequest) (*pb.GetMetricsResponse, error) {
	log.Infof("GetMetrics network %s site %s node %s metric %s",
		req.GetNetworkId(), req.GetSiteId(), req.GetNodeId(), req.GetMetric())

	w := ResolveWindow(req.GetWindow(), time.Now().UTC())

	resourceId := req.GetNodeId()
	if resourceId == "" {
		resourceId = req.GetSiteId()
	}

	names, err := n.metricRepo.MetricNames(resourceId)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "metrics")
	}

	fresh := time.Duration(n.telemetryFreshSeconds) * time.Second

	infos := make([]*pb.MetricInfo, 0, len(names))
	for i := range names {
		nm := &names[i]
		infos = append(infos, &pb.MetricInfo{
			Name:         nm.Metric,
			Unit:         nm.Unit,
			LastSampleAt: timestamppb.New(nm.LastSampleAt),
			Stale:        time.Since(nm.LastSampleAt) > fresh,
		})
	}

	series := []*pb.TimeSeries{}

	if req.GetMetric() != "" {
		rollups, err := n.metricRepo.Rollups(req.GetMetric(), "", resourceId, w.From, w.To)
		if err != nil {
			return nil, grpc.SqlErrorToGrpc(err, "metrics")
		}

		ts := &pb.TimeSeries{Key: req.GetMetric()}
		for i := range rollups {
			ts.Points = append(ts.Points, point(rollups[i].Hour, rollups[i].Avg))
		}

		series = append(series, ts)
	}

	return &pb.GetMetricsResponse{
		Metrics: infos,
		Series:  series,
	}, nil
}

func (n *NetworkServer) GetEvents(ctx context.Context, req *pb.GetEventsRequest) (*pb.GetEventsResponse, error) {
	log.Infof("GetEvents network %s site %s node %s", req.GetNetworkId(), req.GetSiteId(), req.GetNodeId())

	w := ResolveWindow(req.GetWindow(), time.Now().UTC())

	events, total, err := n.eventRepo.Recent(req.GetNetworkId(), req.GetSiteId(),
		req.GetNodeId(), w.From, w.To, req.GetPage(), req.GetPageSize())
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "events")
	}

	return &pb.GetEventsResponse{
		Events: pbEvents(events),
		Meta:   pbMeta(total, req.GetPage(), req.GetPageSize()),
	}, nil
}

func (n *NetworkServer) SupportSearch(ctx context.Context, req *pb.SupportSearchRequest) (*pb.SupportSearchResponse, error) {
	log.Infof("SupportSearch %q network %s", req.GetQuery(), req.GetNetworkId())

	if req.GetQuery() == "" {
		return nil, status.Error(codes.InvalidArgument, "query is required")
	}

	now := time.Now().UTC()
	from30d := now.AddDate(0, 0, -30)

	results := []*pb.SupportResult{}

	sites, err := n.siteRepo.Search(req.GetQuery(), req.GetNetworkId(), supportSearchLimit)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "sites")
	}

	for i := range sites {
		s := &sites[i]

		res := &pb.SupportResult{
			ResourceType: "site",
			ResourceId:   s.SiteId.String(),
			Name:         s.Name,
			Status:       s.Status,
		}

		if customers, err := n.siteRepo.CustomerCount(s.SiteId.String()); err == nil {
			res.Customers = uint32(customers)
		}

		if uptime, err := n.siteRepo.UptimeBetween(s.SiteId.String(), from30d, now); err == nil {
			res.Uptime_30D = uptime
		}

		if pw, err := n.metricRepo.LatestPowerRollup(s.SiteId.String()); err == nil {
			res.BatteryPercent = pw.BatteryPercent
		}

		offline := float64(0)
		if d, err := n.siteRepo.OfflineDuration(s.SiteId.String()); err == nil {
			offline = d
		}

		res.StatusSummary = fmt.Sprintf("site %s is %s; uptime over last 30 days %.1f%%",
			s.Name, s.Status, res.Uptime_30D)
		res.Recommendation = DeriveRecommendation(s.Status, offline)

		results = append(results, res)
	}

	nodes, err := n.nodeRepo.Search(req.GetQuery(), req.GetNetworkId(), supportSearchLimit)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "nodes")
	}

	for i := range nodes {
		nd := &nodes[i]

		res := &pb.SupportResult{
			ResourceType: "node",
			ResourceId:   nd.NodeId,
			Name:         nd.Name,
			Status:       nd.Status,
		}

		if customers, err := n.siteRepo.CustomerCount(nd.SiteId.String()); err == nil {
			res.Customers = uint32(customers)
		}

		if uptime, err := n.nodeRepo.UptimeBetween(nd.NodeId, from30d, now); err == nil {
			res.Uptime_30D = uptime
		}

		if radio, err := n.metricRepo.LatestRadioRollup(nd.NodeId); err == nil {
			res.SignalDbm = radio.SignalDbm
		}

		if pw, err := n.metricRepo.LatestPowerRollup(nd.SiteId.String()); err == nil {
			res.BatteryPercent = pw.BatteryPercent
		}

		res.StatusSummary = fmt.Sprintf("node %s is %s; uptime over last 30 days %.1f%%",
			nd.Name, nd.Status, res.Uptime_30D)
		res.Recommendation = DeriveRecommendation(nd.Status, 0)

		results = append(results, res)
	}

	return &pb.SupportSearchResponse{
		Results: results,
	}, nil
}

// estimateRuntimeHours gives a rough battery runtime estimate. With only
// battery percent, voltage and load available, assume a nominal 100Ah pack:
// energy left = percent/100 * voltage * 100Ah; runtime = energy / load.
func estimateRuntimeHours(batteryPercent, batteryVoltage, loadWatts float64) float64 {
	if loadWatts <= 0 || batteryVoltage <= 0 {
		return 0
	}

	const nominalAmpHours = 100.0

	energyWh := batteryPercent / 100 * batteryVoltage * nominalAmpHours

	return energyWh / loadWatts
}
