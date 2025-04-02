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
             Speed:   metrics.Backhaul.Speed,
         },
         Ethernet: &pb.EthernetMetrics{
             PortStatus: metrics.Backhaul.SwitchStatus,
             PortSpeed:  metrics.Backhaul.SwitchBandwidth,
         },
         Power: &pb.PowerMetrics{
             BatteryPower:      metrics.Battery.Voltage * metrics.Battery.Current,
             SolarPanelVoltage: metrics.Solar.PanelVoltage,
             SolarPanelCurrent: metrics.Solar.PanelCurrent,
             SolarPanelPower:   metrics.Solar.PanelPower,
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
     } else if config.Profile == 0 {
         config.Profile = cenums.PROFILE_NORMAL
         provider.SetProfile(cenums.PROFILE_NORMAL)
         log.Infof("Set default profile to PROFILE_NORMAL for site %s", siteId)
     }
 
     shouldCheckStatus := false
 
     for _, portUpdate := range req.PortUpdates {
         portNumber := int(portUpdate.PortNumber)
         if err := provider.SetPortStatus(portNumber, portUpdate.Status); err != nil {
             log.Warnf("Error updating port %d status for site %s: %v", portNumber, siteId, err)
         } else {
             log.Infof("Updated port %d status to %v for site %s", portNumber, portUpdate.Status, siteId)
             shouldCheckStatus = true
         }
     }
 
     backhaulUpdated := false
     for _, pu := range req.PortUpdates {
         if pu.PortNumber == int32(metrics.PORT_BACKHAUL) {
             backhaulUpdated = true
             break
         }
     }
     if !backhaulUpdated {
         if err := provider.SetPortStatus(metrics.PORT_BACKHAUL, true); err != nil {
             log.Warnf("Failed to ensure backhaul port enabled for site %s: %v", siteId, err)
         } else {
             log.Infof("Ensured backhaul port enabled for site %s", siteId)
         }
     }
 
     metrics, err := provider.GetMetrics(siteId)
     if err != nil {
         log.Errorf("Failed to get metrics for site %s: %v", siteId, err)
         return &pb.UpdateMetricsResponse{
             Success: false,
             Message: fmt.Sprintf("Failed to retrieve metrics: %v", err),
         }, nil
     }
 
     voltage := metrics.Battery.Voltage
     var percentage float64
     switch config.Profile {
     case cenums.PROFILE_MIN:
         if voltage <= 10.5 {
             percentage = 0
         } else if voltage >= 12.0 {
             percentage = 100
         } else {
             percentage = (voltage - 10.5) / (12.0 - 10.5) * 100
         }
     case cenums.PROFILE_NORMAL:
         if voltage <= 10.5 {
             percentage = 0
         } else if voltage >= 12.7 {
             percentage = 100
         } else {
             percentage = (voltage - 10.5) / (12.7 - 10.5) * 100
         }
     case cenums.PROFILE_MAX:
         if voltage <= 11.0 {
             percentage = 0
         } else if voltage >= 13.5 {
             percentage = 100
         } else {
             percentage = (voltage - 11.0) / (13.5 - 11.0) * 100
         }
     default:
         percentage = (voltage - 10.5) / (12.7 - 10.5) * 100
     }
 
     isUp := metrics.Backhaul.Speed > 0
 
     log.Debugf("Site %s - Profile: %v, Voltage: %.2fV, Battery percentage: %.2f%%, Backhaul speed: %.2f, Switch port power: %.2fW, Site up: %v",
         siteId, config.Profile, voltage, percentage, metrics.Backhaul.Speed, metrics.Backhaul.SwitchPortPower, isUp)
 
     var needToCheckStatus bool
     var siteToCheck string
 
     if shouldCheckStatus {
         monConfig, monExists := s.monitoringStatus[siteId]
         if monExists && monConfig.Active {
             needToCheckStatus = true
             siteToCheck = siteId
         }
     }
 
     resp := &pb.UpdateMetricsResponse{
         Success: true,
         Message: "Metrics updated",
     }
 
     if needToCheckStatus {
         defer func(site string) {
             s.checkAndUpdateStatus(site, "")
         }(siteToCheck)
     }
 
     return resp, nil
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
 
     config, exists := s.siteConfigs[siteId]
     if !exists || !config.Active {
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
 
     nodes, err := s.nodeClient.GetNodesBySite(siteId)
     if err != nil || len(nodes.Nodes) == 0 {
         return &pb.MonitorSiteResponse{
             Success: false,
             Message: "Failed to validate site nodes or no nodes found",
         }, nil
     }
 
     if config.Exporter != nil {
         config.Exporter.ResetUptimeCounter()
         log.Infof("Initialized uptime counter for site %s", siteId)
     }
 
     monitorCtx, cancelFunc := context.WithCancel(context.Background())
     monConfig := &MonitoringConfig{
         NodeId:     "",
         Active:     true,
         Context:    monitorCtx,
         CancelFunc: cancelFunc,
         LastStatus: cenums.SCENARIO_DEFAULT, 
     }
     s.monitoringStatus[siteId] = monConfig
     s.mutex.RLock()
     provider, exists := s.metricsProviders[siteId]
     siteConfig, siteExists := s.siteConfigs[siteId]
     s.mutex.RUnlock()
 
     if exists && siteExists {
         metrics, err := provider.GetMetrics(siteId)
         if err == nil {
             percentage := s.calculateBatteryPercentage(siteConfig.Profile, metrics.Battery.Voltage)
             currentScenario := s.determineScenario(percentage, metrics.Backhaul.Speed)
             monConfig.LastStatus = currentScenario
             log.Infof("Initialized LastStatus to %s for site %s", currentScenario, siteId)
         } else {
             log.Warnf("Failed to get initial metrics for site %s: %v", siteId, err)
         }
     }
 
     go s.monitorSiteStatusWorker(siteId, "")
 
 
     return &pb.MonitorSiteResponse{
         Success: true,
         Message: fmt.Sprintf("Started status monitoring for site %s with %d nodes", siteId, len(nodes.Nodes)),
     }, nil
 }
 
 func (s *DControllerServer) monitorSiteStatusWorker(siteId, _ string) {
     log.Infof("Starting status monitoring for site %s", siteId)
 
     uptimeTicker := time.NewTicker(1 * time.Second)
     defer uptimeTicker.Stop()
 
     statusTicker := time.NewTicker(monitorInterval)
     defer statusTicker.Stop()
 
     s.checkAndUpdateStatus(siteId, "")
 
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
 
         case <-statusTicker.C:
             s.checkAndUpdateStatus(siteId, "")
 
         case <-uptimeTicker.C:
             s.incrementUptimeIfUp(siteId)
         }
     }
 }
 
 func (s *DControllerServer) incrementUptimeIfUp(siteId string) {
    s.mutex.RLock()
    provider, providerExists := s.metricsProviders[siteId]
    siteConfig, siteExists := s.siteConfigs[siteId]
    s.mutex.RUnlock()

    if !providerExists {
        log.Warnf("Metrics provider not found for site %s, skipping uptime update", siteId)
        return
    }
    if !siteExists {
        log.Warnf("Site configuration not found for site %s, skipping uptime update", siteId)
        return
    }
    if siteConfig.Exporter == nil {
        log.Warnf("Exporter is nil for site %s, skipping uptime update", siteId)
        return
    }

    backhaulPortOn := provider.GetPortStatus(metrics.PORT_BACKHAUL)
    log.Debugf("Checking uptime for site %s: backhaul port status = %v", siteId, backhaulPortOn)

    if backhaulPortOn {
        siteConfig.Exporter.IncrementUptimeCounter(1.0)
        log.Debugf("Incremented uptime counter for site %s", siteId)
    } else {
        log.Infof("Site %s is down due to backhaul port being off. Resetting uptime counter", siteId)
        siteConfig.Exporter.ResetUptimeCounter()
    }
}
 
 func (s *DControllerServer) resetUptimeCounter(siteId string) {
     s.mutex.RLock()
     defer s.mutex.RUnlock()
 
     siteConfig, exists := s.siteConfigs[siteId]
     if !exists || siteConfig.Exporter == nil {
         return
     }
 
     siteConfig.Exporter.ResetUptimeCounter()
 }
 
 func (s *DControllerServer) calculateBatteryPercentage(profile cenums.Profile, voltage float64) float64 {
     switch profile {
     case cenums.PROFILE_MIN:
         if voltage <= 10.0 {
             return 0
         } else if voltage >= 12.0 {
             return 100
         }
         return (voltage - 10.0) / (12.0 - 10.0) * 100
     case cenums.PROFILE_MAX:
         if voltage <= 12.0 {
             return 70
         } else if voltage >= 13.0 {
             return 100
         }
         return 70 + (voltage-12.0)/(13.0-12.0)*30
     default:
         if voltage <= 10.5 {
             return 0
         } else if voltage >= 12.7 {
             return 100
         }
         return (voltage - 10.5) / (12.7 - 10.5) * 100
     }
 }
 
 func (s *DControllerServer) determineScenario(percentage float64, backhaulSpeed float64) cenums.SCENARIOS {
     if percentage < 50 {
         return cenums.SCENARIO_NODE_OFF
     } else if backhaulSpeed <= 0 {
         return cenums.SCENARIO_BACKHAUL_DOWN
     }
     return cenums.SCENARIO_DEFAULT
 }
 
 func (s *DControllerServer) checkAndUpdateStatus(siteId, _ string) {
     s.mutex.RLock()
     provider, providerExists := s.metricsProviders[siteId]
     monConfig, monExists := s.monitoringStatus[siteId]
     siteConfig, siteExists := s.siteConfigs[siteId]
     if !providerExists || !monExists || !siteExists {
         s.mutex.RUnlock()
         log.Warnf("Missing configuration for site %s", siteId)
         return
     }
 
     lastStatus := monConfig.LastStatus
     currentProfile := siteConfig.Profile
     s.mutex.RUnlock()
 
     metrics, err := provider.GetMetrics(siteId)
     if err != nil {
         log.Errorf("Failed to get metrics for site %s: %v", siteId, err)
         return
     }
 
     voltage := metrics.Battery.Voltage
     percentage := s.calculateBatteryPercentage(currentProfile, voltage)
     currentScenario := s.determineScenario(percentage, metrics.Backhaul.Speed)
 
     switch currentScenario {
     case cenums.SCENARIO_NODE_OFF:
         log.Infof("Battery low for site %s: percentage=%.1f%%, voltage=%.1fV, setting scenario to %s",
             siteId, percentage, voltage, currentScenario)
     case cenums.SCENARIO_BACKHAUL_DOWN:
         log.Infof("Backhaul down detected for site %s: speed=%.1f, setting scenario to %s",
             siteId, metrics.Backhaul.Speed, currentScenario)
     case cenums.SCENARIO_DEFAULT:
         log.Infof("Setting scenario to default for site %s", siteId)
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
 
     if siteConfig, exists := s.siteConfigs[siteId]; exists && siteConfig.Exporter != nil {
         siteConfig.Exporter.ResetUptimeCounter()
         log.Infof("Reset uptime counter for site %s due to monitoring stop", siteId)
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