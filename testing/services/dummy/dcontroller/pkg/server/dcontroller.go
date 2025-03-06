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
	cenums "github.com/ukama/ukama/testing/common/enums"
	pb "github.com/ukama/ukama/testing/services/dummy/dcontroller/pb/gen"
	"github.com/ukama/ukama/testing/services/dummy/dcontroller/pkg"

	"github.com/ukama/ukama/testing/services/dummy/dcontroller/pkg/metrics"
)
 
 type SiteMetricsConfig struct {
	 ScanInterval int
	 Profile      cenums.Profile
	 Active       bool
	 Exporter     *metrics.PrometheusExporter
	 Context      context.Context
	 CancelFunc   context.CancelFunc
 }
 
 type DControllerServer struct {
	 pb.UnimplementedMetricsControllerServer
	 orgName          string
	 metricsProviders map[string]*metrics.MetricsProvider
	 siteConfigs      map[string]*SiteMetricsConfig
	 mutex            sync.RWMutex
	 msgbus           mb.MsgBusServiceClient
	 baseRoutingKey   msgbus.RoutingKeyBuilder
 }
 
 func NewControllerServer(orgName string, msgBus mb.MsgBusServiceClient) *DControllerServer {
	 return &DControllerServer{
		 orgName:          orgName,
		 metricsProviders: make(map[string]*metrics.MetricsProvider),
		 siteConfigs:      make(map[string]*SiteMetricsConfig),
		 mutex:            sync.RWMutex{},
		 baseRoutingKey:   msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
		 msgbus:           msgBus,
	 }
 }
 
 func (s *DControllerServer) GetSiteMetrics(ctx context.Context, req *pb.GetSiteMetricsRequest) (*pb.GetSiteMetricsResponse, error) {
	 siteId := req.SiteId
	 if siteId == "" {
		 return nil, fmt.Errorf("site ID is required")
	 }
	 
	 s.mutex.RLock()
	 provider, exists := s.metricsProviders[siteId]
	 if !exists {
		 s.mutex.RUnlock()
		 return nil, fmt.Errorf("no metrics available for site %s", siteId)
	 }
	 s.mutex.RUnlock()
 
	 systemMetrics, err := provider.GetMetrics(siteId)
	 if err != nil {
		 return nil, err
	 }
 
	 return &pb.GetSiteMetricsResponse{
		 Solar: &pb.SolarMetrics{
			 PowerGeneration: systemMetrics.Solar.PowerGeneration,
			 EnergyTotal:    systemMetrics.Solar.EnergyTotal,
			 PanelPower:     systemMetrics.Solar.PanelPower,
			 PanelCurrent:   systemMetrics.Solar.PanelCurrent,
			 PanelVoltage:   systemMetrics.Solar.PanelVoltage,
			 InverterStatus: systemMetrics.Solar.InverterStatus,
		 },
		 Battery: &pb.BatteryMetrics{
			 ChargeStatus: systemMetrics.Battery.Capacity,
			 Voltage:      systemMetrics.Battery.Voltage,
			 Health:       map[string]float64{
				 "Good": 1.0,
				 "Fair": 0.5,
				 "Poor": 0.0,
			 }[systemMetrics.Battery.Health],
			 Current:      systemMetrics.Battery.Current,
			 Temperature:  systemMetrics.Battery.Temperature,
		 },
		 Network: &pb.NetworkMetrics{
			 BackhaulLatency:      systemMetrics.Backhaul.Latency,
			 BackhaulStatus:       systemMetrics.Backhaul.Status,
			 BackhaulSpeed:        systemMetrics.Backhaul.Speed,
			 SwitchPortStatus:     systemMetrics.Backhaul.SwitchStatus,
			 SwitchPortBandwidth:  systemMetrics.Backhaul.SwitchBandwidth,
		 },
	 }, nil
 }
 
 func (s *DControllerServer) StartMetrics(ctx context.Context, req *pb.StartMetricsRequest) (*pb.StartMetricsResponse, error) {
	siteId := req.SiteId
	
	log.Infof("Starting metrics for site ID: %s", siteId)
	
	scanInterval := 3
	log.Infof("Starting metrics collection goroutine for site %s with scan interval %d seconds", 
	siteId, scanInterval)
	
	profile := cenums.Profile(req.Profile)
	
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
		ScanInterval: scanInterval,
		Profile:      profile,
		Active:       true,
		Exporter:     exporter,
		Context:      siteCtx,
		CancelFunc:   cancelFunc,
	}
	
	go func() {
		scanIntervalDuration := time.Duration(scanInterval) * time.Second
		log.Infof("Inside goroutine: Starting metrics collection for site %s", siteId)
		err := exporter.StartMetricsCollection(siteCtx, scanIntervalDuration)
		if err != nil && err != context.Canceled {
			log.Infof("ERROR collecting metrics for site %s: %v\n", siteId, err)
		}
	}()
	
	return &pb.StartMetricsResponse{
		Success: true,
		Message: "Started metrics collection",
	}, nil
}
func (s *DControllerServer) UpdateMetrics(ctx context.Context, req *pb.UpdateMetricsRequest) (*pb.UpdateMetricsResponse, error) {
	siteId := req.SiteId

	log.Infof("Updating metrics for site ID: %s", siteId)
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

	if len(req.PortUpdates) > 0 {
		for _, portUpdate := range req.PortUpdates {
			portNumber := int(portUpdate.PortNumber)
			portStatus := portUpdate.Status

			err := provider.SetPortStatus(portNumber, portStatus)
			if err != nil {
				log.Infof("Error updating port %d status: %v", portNumber, err)
			} else {
				log.Infof("Updated port %d status to %v for site %s", portNumber, portStatus, siteId)
			}
		}
	}

	return &pb.UpdateMetricsResponse{
		Success: true,
		Message: "metrics updated",
	}, nil
}
 func (s *DControllerServer) StopMetricsCollection(siteId string) bool {
	 s.mutex.Lock()
	 defer s.mutex.Unlock()
	 
	 config, exists := s.siteConfigs[siteId]
	 if !exists || !config.Active {
		 return false
	 }
	 
	 config.CancelFunc()
	 config.Exporter.Shutdown()
	 config.Active = false
	 
	 return true
 }
 
 func (s *DControllerServer) Cleanup() {
	 s.mutex.Lock()
	 defer s.mutex.Unlock()
	 
	 for _, config := range s.siteConfigs {
		 if config.Active {
			 config.CancelFunc()
			 config.Exporter.Shutdown()
			 config.Active = false
		 }
	 }
 }
