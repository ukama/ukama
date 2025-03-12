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
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
	creg "github.com/ukama/ukama/systems/common/rest/client/registry"
	cenums "github.com/ukama/ukama/testing/common/enums"
	pb "github.com/ukama/ukama/testing/services/dummy/dcontroller/pb/gen"
	"github.com/ukama/ukama/testing/services/dummy/dcontroller/pkg"
	"github.com/ukama/ukama/testing/services/dummy/dcontroller/pkg/client"
	"github.com/ukama/ukama/testing/services/dummy/dcontroller/pkg/metrics"
)
  
 const (
	 defaultScanInterval = 3
	 monitorInterval     = 5 * time.Second
 )
  
 type SiteMetricsConfig struct {
	 Profile    cenums.Profile
	 Active     bool
	 Exporter   *metrics.PrometheusExporter
	 Context    context.Context
	 CancelFunc context.CancelFunc
 }
  
 type MonitoringConfig struct {
	 NodeId     string
	 Active     bool
	 Context    context.Context
	 CancelFunc context.CancelFunc
	 LastStatus cenums.SCENARIOS
 }
  
 type DControllerServer struct {
	 pb.UnimplementedMetricsControllerServer
	 orgName          string
	 metricsProviders map[string]*metrics.MetricsProvider
	 siteConfigs      map[string]*SiteMetricsConfig
	 monitoringStatus map[string]*MonitoringConfig
	 mutex            sync.RWMutex
	 msgbus           mb.MsgBusServiceClient
	 baseRoutingKey   msgbus.RoutingKeyBuilder
	 dnodeClient      *client.DNodeClient
	 nodeClient       creg.NodeClient
 }
  
 func NewControllerServer(orgName string, msgBus mb.MsgBusServiceClient, nodeClient creg.NodeClient) *DControllerServer {
	 return &DControllerServer{
		 orgName:          orgName,
		 metricsProviders: make(map[string]*metrics.MetricsProvider),
		 siteConfigs:      make(map[string]*SiteMetricsConfig),
		 monitoringStatus: make(map[string]*MonitoringConfig),
		 mutex:            sync.RWMutex{},
		 baseRoutingKey:   msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
		 msgbus:           msgBus,
		 nodeClient:       nodeClient,
	 }
 }
  
 func (s *DControllerServer) AddScenarioMonitoring(dnodeBaseURL string) {
	 s.dnodeClient = client.NewDNodeClient(dnodeBaseURL, 10*time.Second)
	 log.Infof("DNode client initialized with base URL: %s", dnodeBaseURL)
 }
  
 func (s *DControllerServer) GetSiteMetrics(ctx context.Context, req *pb.GetSiteMetricsRequest) (*pb.GetSiteMetricsResponse, error) {
	 siteId := req.SiteId
	 if siteId == "" {
		 return nil, fmt.Errorf("site ID is required")
	 }
 
	 s.mutex.RLock()
	 provider, exists := s.metricsProviders[siteId]
	 s.mutex.RUnlock()
	 
	 if !exists {
		 return nil, fmt.Errorf("no metrics available for site %s", siteId)
	 }
 
	 metrics, err := provider.GetMetrics(siteId)
	 if err != nil {
		 return nil, err
	 }
 
	 return &pb.GetSiteMetricsResponse{
		 Backhaul: &pb.BackhaulMetrics{
			 Latency: metrics.Backhaul.Latency,
			 Status:  metrics.Backhaul.Status,
			 Speed:   metrics.Backhaul.Speed,
		 },
		 Ethernet: &pb.EthernetMetrics{
			 PortStatus: metrics.Backhaul.SwitchStatus,
			 PortSpeed:  metrics.Backhaul.SwitchBandwidth,
		 },
		 Power: &pb.PowerMetrics{
			 BatteryPower:            metrics.Battery.Voltage * metrics.Battery.Current,
			 SolarPanelVoltage:       metrics.Solar.PanelVoltage,
			 SolarPanelCurrent:       metrics.Solar.PanelCurrent,
			 SolarPanelPower:         metrics.Solar.PanelPower,
			 ChargeControllerStatus:  metrics.Solar.ControllerStatus,
			 ChargeControllerMode:    float64(metrics.Solar.ControllerModeValue),
			 ChargeControllerCurrent: metrics.Solar.ControllerCurrent,
			 ChargeControllerVoltage: metrics.Solar.ControllerVoltage,
		 },
	 }, nil
 }
  
 func (s *DControllerServer) StartMetrics(ctx context.Context, req *pb.StartMetricsRequest) (*pb.StartMetricsResponse, error) {
	 siteId := req.SiteId
	 profile := cenums.Profile(req.Profile)
	 
	 log.Infof("Starting metrics for site ID: %s with profile: %v", siteId, profile)
	 
	 s.mutex.Lock()
	 defer s.mutex.Unlock()
	 
	 if config, exists := s.siteConfigs[siteId]; exists && config.Active {
		 return &pb.StartMetricsResponse{
			 Success: false,
			 Message: "Site metrics already active",
		 }, nil
	 }
	 
	 if _, exists := s.metricsProviders[siteId]; !exists {
		 s.metricsProviders[siteId] = metrics.NewMetricsProvider()
	 }
	 
	 s.metricsProviders[siteId].SetProfile(profile)
	 
	 siteCtx, cancelFunc := context.WithCancel(context.Background())
	 exporter := metrics.NewPrometheusExporter(s.metricsProviders[siteId], siteId)
	 
	 s.siteConfigs[siteId] = &SiteMetricsConfig{
		 Profile:    profile,
		 Active:     true,
		 Exporter:   exporter,
		 Context:    siteCtx,
		 CancelFunc: cancelFunc,
	 }
	 
	 go s.collectMetrics(siteId, exporter, siteCtx)
	 
	 return &pb.StartMetricsResponse{
		 Success: true,
		 Message: "Started metrics collection",
	 }, nil
 }
  
 func (s *DControllerServer) collectMetrics(siteId string, exporter *metrics.PrometheusExporter, ctx context.Context) {
	 interval := time.Duration(defaultScanInterval) * time.Second
	 log.Infof("Starting metrics collection for site %s with interval %v", siteId, interval)
	 
	 err := exporter.StartMetricsCollection(ctx, interval)
	 if err != nil && err != context.Canceled {
		 log.Errorf("ERROR collecting metrics for site %s: %v", siteId, err)
	 }
 }
 func (s *DControllerServer) UpdateMetrics(ctx context.Context, req *pb.UpdateMetricsRequest) (*pb.UpdateMetricsResponse, error) {
	siteId := req.SiteId
	if siteId == "" {
		return &pb.UpdateMetricsResponse{
			Success: false,
			Message: "Site ID is required",
		}, nil
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	config, exists := s.siteConfigs[siteId]
	if !exists || !config.Active {
		return &pb.UpdateMetricsResponse{
			Success: false,
			Message: "Site metrics not active",
		}, nil
	}

	provider, exists := s.metricsProviders[siteId]
	if !exists {
		return &pb.UpdateMetricsResponse{
			Success: false,
			Message: "Metrics provider not found",
		}, nil
	}

	if req.Profile > 0 {
		profile := cenums.Profile(req.Profile)
		config.Profile = profile
		provider.SetProfile(profile)
		log.Infof("Updated profile to %v for site %s", profile, siteId)
	}

	shouldCheckStatus := false
	
	for _, portUpdate := range req.PortUpdates {
		portNumber := int(portUpdate.PortNumber)
		if err := provider.SetPortStatus(portNumber, portUpdate.Status); err != nil {
			log.Warnf("Error updating port %d status: %v", portNumber, err)
		} else {
			log.Infof("Updated port %d status to %v for site %s", portNumber, portUpdate.Status, siteId)
			shouldCheckStatus = true
		}
	}
	if shouldCheckStatus {
		monConfig, monExists := s.monitoringStatus[siteId]
		if monExists && monConfig.Active {
			s.mutex.Unlock()
			s.checkAndUpdateStatus(siteId, "")
			s.mutex.Lock()
		}
	}

	return &pb.UpdateMetricsResponse{
		Success: true,
		Message: "Metrics updated",
	}, nil
}
  
 func (s *DControllerServer) MonitorSite(ctx context.Context, req *pb.MonitorSiteRequest) (*pb.MonitorSiteResponse, error) {
	 siteId := req.SiteId
	 
	 if siteId == "" {
		 return &pb.MonitorSiteResponse{
			 Success: false,
			 Message: "Site ID is required",
		 }, nil
	 }
	 
	 s.mutex.Lock()
	 defer s.mutex.Unlock()
	 
	 if config, exists := s.siteConfigs[siteId]; !exists || !config.Active {
		 return &pb.MonitorSiteResponse{
			 Success: false,
			 Message: "Site metrics not active",
		 }, nil
	 }
	 
	 if _, exists := s.monitoringStatus[siteId]; exists {
		 return &pb.MonitorSiteResponse{
			 Success: false,
			 Message: "Site is already being monitored",
		 }, nil
	 }
	 
	 // Validate that the site has nodes
	 nodes, err := s.nodeClient.GetNodesBySite(siteId)
	 if err != nil {
		 log.Errorf("Failed to get nodes for site %s: %v", siteId, err)
		 return &pb.MonitorSiteResponse{
			 Success: false,
			 Message: fmt.Sprintf("Failed to validate site nodes: %v", err),
		 }, nil
	 }
	 
	 if len(nodes.Nodes) == 0 {
		 return &pb.MonitorSiteResponse{
			 Success: false,
			 Message: "Site has no associated nodes",
		 }, nil
	 }
	 
	 monitorCtx, cancelFunc := context.WithCancel(context.Background())
	 s.monitoringStatus[siteId] = &MonitoringConfig{
		 NodeId:     "",  
		 Active:     true,
		 Context:    monitorCtx,
		 CancelFunc: cancelFunc,
		 LastStatus: cenums.SCENARIO_BACKHAUL_DOWN,
	 }
	 
	 go s.monitorSiteStatusWorker(siteId, "")
	 
	 return &pb.MonitorSiteResponse{
		 Success: true,
		 Message: fmt.Sprintf("Started status monitoring for site %s with %d nodes", siteId, len(nodes.Nodes)),
	 }, nil
 }
  
 func (s *DControllerServer) monitorSiteStatusWorker(siteId, _ string) {
	 log.Infof("Starting status monitoring for site %s", siteId)
	 
	 ticker := time.NewTicker(monitorInterval)
	 defer ticker.Stop()
	 
	 for {
		 s.mutex.RLock()
		 monConfig, exists := s.monitoringStatus[siteId]
		 if !exists || !monConfig.Active {
			 s.mutex.RUnlock()
			 log.Infof("Monitoring for site %s has been stopped", siteId)
			 return
		 }
		 monCtx := monConfig.Context
		 s.mutex.RUnlock()
		 
		 select {
		 case <-monCtx.Done():
			 log.Infof("Monitoring context for site %s has been cancelled", siteId)
			 return
		 case <-ticker.C:
			 s.checkAndUpdateStatus(siteId, "")
		 }
	 }
 }
  
 func (s *DControllerServer) checkAndUpdateStatus(siteId, _ string) {
	s.mutex.RLock()
	provider, providerExists := s.metricsProviders[siteId]
	monConfig, monExists := s.monitoringStatus[siteId]
	_, siteExists := s.siteConfigs[siteId]
	
	if !providerExists || !monExists || !siteExists {
		s.mutex.RUnlock()
		log.Warnf("Missing configuration for site %s", siteId)
		return
	}
	
	lastStatus := monConfig.LastStatus
	s.mutex.RUnlock()
	
	metrics, err := provider.GetMetrics(siteId)
	if err != nil {
		log.Errorf("Failed to get metrics for site %s: %v", siteId, err)
		return
	}
	
	var currentScenario cenums.SCENARIOS
	
	powerStatus, err := provider.GetPowerStatus()
	switch {
	case err != nil || !powerStatus:
		currentScenario = cenums.SCENARIO_NODE_OFF
	case metrics.Backhaul.Status < 0.5 || metrics.Backhaul.Speed <= 0:
		currentScenario = cenums.SCENARIO_BACKHAUL_DOWN
		log.Infof("Backhaul down detected for site %s: status=%.1f, speed=%.1f", 
			siteId, metrics.Backhaul.Status, metrics.Backhaul.Speed)

	case metrics.Battery.Power <= 0:
		currentScenario = cenums.SCENARIO_NODE_OFF
		log.Infof("Battery down detected for site %s: voltage=%.1f, current=%.1f",
			siteId, metrics.Battery.Voltage, metrics.Battery.Current)
	default:
		currentScenario = cenums.SCENARIO_DEFAULT
	}
	
	if currentScenario != lastStatus && s.dnodeClient != nil {
		log.Infof("Status changed for site %s from %s to %s", siteId, lastStatus, currentScenario)
		
		nodes, err := s.nodeClient.GetNodesBySite(siteId)
		if err != nil {
			log.Errorf("Failed to get nodes for site %s: %v", siteId, err)
			return
		}
		
		for _, node := range nodes.Nodes {
			if err := s.dnodeClient.UpdateNodeScenario(node.Id, currentScenario); err != nil {
				log.Errorf("Failed to update node %s scenario: %v", node.Id, err)
			} else {
				log.Infof("Updated scenario for node %s to %s", node.Id, currentScenario)
			}
		}
		
		s.mutex.Lock()
		if config, exists := s.monitoringStatus[siteId]; exists {
			config.LastStatus = currentScenario
		}
		s.mutex.Unlock()
	}
}
 func (s *DControllerServer) StopMonitoring(ctx context.Context, req *pb.StopMonitoringRequest) (*pb.StopMonitoringResponse, error) {
	 siteId := req.SiteId
	 
	 s.mutex.Lock()
	 defer s.mutex.Unlock()
	 
	 config, exists := s.monitoringStatus[siteId]
	 if !exists || !config.Active {
		 return &pb.StopMonitoringResponse{
			 Success: false,
			 Message: "Site is not being monitored",
		 }, nil
	 }
	 
	 config.CancelFunc()
	 config.Active = false
	 
	 if s.dnodeClient != nil{
		 nodes, err := s.nodeClient.GetNodesBySite(siteId)
		 if err != nil {
			 log.Warnf("Failed to get nodes for site %s: %v", siteId, err)
		 } else {
			 for _, node := range nodes.Nodes {
				 if err := s.dnodeClient.SetDefault(node.Id); err != nil {
					 log.Warnf("Failed to reset node %s status to default: %v", node.Id, err)
				 } else {
					 log.Infof("Reset node %s to default status", node.Id)
				 }
			 }
		 }
	 }
	 
	 delete(s.monitoringStatus, siteId)
	 
	 return &pb.StopMonitoringResponse{
		 Success: true,
		 Message: "Stopped site status monitoring",
	 }, nil
 }
  

  
 func (s *DControllerServer) Cleanup() {
	 s.mutex.Lock()
	 
	 for siteId, config := range s.siteConfigs {
		 if config.Active {
			 config.CancelFunc()
			 config.Exporter.Shutdown()
			 config.Active = false
			 
			 if monConfig, exists := s.monitoringStatus[siteId]; exists && monConfig.Active {
				 monConfig.CancelFunc()
				 monConfig.Active = false
			 }
		 }
	 }
	 
	 s.mutex.Unlock()
	 log.Info("All metrics collection and monitoring stopped")
 }