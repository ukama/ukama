/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
package server

import (
	"context"
	"fmt"

	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
	creg "github.com/ukama/ukama/systems/common/rest/client/registry"
	pb "github.com/ukama/ukama/testing/services/dummy/dcontroller/pb/gen"
	"github.com/ukama/ukama/testing/services/dummy/dcontroller/pkg"
	"github.com/ukama/ukama/testing/services/dummy/dcontroller/pkg/client"
	"github.com/ukama/ukama/testing/services/dummy/dcontroller/pkg/metrics"
)

type DControllerServer struct {
	pb.UnimplementedMetricsControllerServer
	orgName          string
	msgbus           mb.MsgBusServiceClient           
	baseRoutingKey   msgbus.RoutingKeyBuilder
	metricsManager   *metrics.MetricsManager  
    dnodeClient      *client.DNodeClient
    nodeClient       creg.NodeClient        
}

func NewControllerServer(orgName string, msgBus mb.MsgBusServiceClient, nodeClient creg.NodeClient, dnodeClient *client.DNodeClient) *DControllerServer {
	return &DControllerServer{
		orgName:        orgName,
		baseRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
		msgbus:         msgBus,
		metricsManager: metrics.NewMetricsManager(),
        nodeClient:    nodeClient,
        dnodeClient:   dnodeClient,
	}
}

func (s *DControllerServer) StartMetrics(ctx context.Context, req *pb.StartMetricsRequest) (*pb.StartMetricsResponse, error) {
	siteId := req.SiteId
	if siteId == "" {
		return &pb.StartMetricsResponse{
			Success: false,
			Message: "Site ID is required",
		}, nil
	}

	config := metrics.SiteConfig{
		AvgBackhaulSpeed: req.SiteConfig.AvgBackhaulSpeed,
		AvgLatency: req.SiteConfig.AvgLatency,
		SolarEfficiency: req.SiteConfig.SolarEfficiency,
	}
	err := s.metricsManager.StartSiteMetrics(siteId, config)
	if err != nil {
		return &pb.StartMetricsResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to start metrics: %v", err),
		}, nil
	}

	return &pb.StartMetricsResponse{
		Success: true,
		Message: fmt.Sprintf("Started metrics for site %s", siteId),
	}, nil
}

func (s *DControllerServer) StopMetrics(ctx context.Context, req *pb.StopMetricsRequest) (*pb.StopMetricsResponse, error) {
	siteId := req.SiteId
	if siteId == "" {
		return &pb.StopMetricsResponse{
			Success: false,
			Message: "Site ID is required",
		}, nil
	}

	if !s.metricsManager.IsMetricsRunning(siteId) {
		return &pb.StopMetricsResponse{
			Success: false,
			Message: "Metrics not running for this site",
		}, nil
	}

	s.metricsManager.StopSiteMetrics(siteId)
	return &pb.StopMetricsResponse{
		Success: true,
		Message: fmt.Sprintf("Stopped metrics for site %s", siteId),
	}, nil
}

func (s *DControllerServer) UpdatePortStatus(ctx context.Context, req *pb.UpdatePortStatusRequest) (*pb.UpdatePortStatusResponse, error) {
	siteId := req.SiteId
	if siteId == "" {
		return &pb.UpdatePortStatusResponse{
			Success: false,
			Message: "Site ID is required",
		}, nil
	}
	
	if !s.metricsManager.IsMetricsRunning(siteId) {
		return &pb.UpdatePortStatusResponse{
			Success: false,
			Message: "Metrics not running for this site",
		}, nil
	}
	
	portNumber := int(req.PortNumber)
	enabled := req.Enabled
	
	err := s.metricsManager.UpdatePortStatus(siteId, portNumber, enabled)
	if err != nil {
		return &pb.UpdatePortStatusResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to update port status: %v", err),
		}, nil
	}
	
	statusStr := "disabled"
	if enabled {
		statusStr = "enabled"
	}
	
	var portName string
	switch portNumber {
	case metrics.PORT_NODE:
		portName = "Node"
	case metrics.PORT_SOLAR:
		portName = "Solar"
	case metrics.PORT_BACKHAUL:
		portName = "Backhaul"
	default:
		portName = fmt.Sprintf("Unknown (%d)", portNumber)
	}
	
	return &pb.UpdatePortStatusResponse{
		Success: true,
		Message: fmt.Sprintf("%s port %s for site %s", portName, statusStr, siteId),
	}, nil
}

func (s *DControllerServer) GetSiteMetrics(ctx context.Context, req *pb.GetSiteMetricsRequest) (*pb.GetSiteMetricsResponse, error) {
    siteId := req.SiteId
    if siteId == "" {
        return nil, fmt.Errorf("site ID is required")
    }

    metrics, err := s.metricsManager.GetSiteMetrics(siteId)
    if err != nil {
        return nil, fmt.Errorf("failed to get metrics: %v", err)
    }

    return &pb.GetSiteMetricsResponse{
        Backhaul: &pb.BackhaulMetrics{
            Latency: metrics["main_backhaul_latency"],
            Speed:   metrics["backhaul_speed"],
            Status:  metrics["backhaul_switch_port_status"],
        },
        Ethernet: &pb.EthernetMetrics{
            PortStatus: metrics["node_switch_port_status"],
            PortSpeed:  metrics["node_switch_port_speed"],
        },
        Power: &pb.PowerMetrics{
            BatteryPower:      metrics["battery_charge_percentage"],
            SolarPanelVoltage: metrics["solar_panel_voltage"],
            SolarPanelCurrent: metrics["solar_panel_current"],
            SolarPanelPower:   metrics["solar_panel_power"],
        },
        UptimeSeconds: metrics["site_uptime_seconds"],
    }, nil
}

