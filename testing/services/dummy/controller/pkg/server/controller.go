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
	pb "github.com/ukama/ukama/testing/services/dummy/controller/pb/gen"
	"github.com/ukama/ukama/testing/services/dummy/controller/pkg"
	"github.com/ukama/ukama/testing/services/dummy/controller/pkg/client"

	"github.com/ukama/ukama/testing/services/dummy/controller/pkg/metrics"
)



type SiteMetricsConfig struct {
	ScanInterval int
	Profile      cenums.Profile
	Scenario     cenums.SCENARIOS
	Active       bool
	Exporter     *metrics.PrometheusExporter
	Context      context.Context
	CancelFunc   context.CancelFunc
}

type ControllerServer struct {
	pb.UnimplementedMetricsControllerServer
	orgName          string
	metricsProviders map[string]*metrics.MetricsProvider
	siteConfigs      map[string]*SiteMetricsConfig
	mutex            sync.RWMutex
	nodeClient        creg.NodeClient
	dnodeClient      *client.DNodeClient 
	msgbus         mb.MsgBusServiceClient
	baseRoutingKey msgbus.RoutingKeyBuilder

}

func NewControllerServer(orgName string, nodeClient creg.NodeClient,dnodeHost string, msgBus mb.MsgBusServiceClient) *ControllerServer {
	dnodeClient := client.NewDNodeClient(dnodeHost, 5*time.Second)
	return &ControllerServer{
		orgName:          orgName,
		metricsProviders: make(map[string]*metrics.MetricsProvider),
		siteConfigs:      make(map[string]*SiteMetricsConfig),
		mutex:            sync.RWMutex{},
		nodeClient:       nodeClient,
		dnodeClient:      dnodeClient,
		baseRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
		msgbus: msgBus,
	}
}

func (s *ControllerServer) GetSiteMetrics(ctx context.Context, req *pb.GetSiteMetricsRequest) (*pb.GetSiteMetricsResponse, error) {
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

func (s *ControllerServer) StartMetrics(ctx context.Context, req *pb.StartMetricsRequest) (*pb.StartMetricsResponse, error) {
	siteId := req.SiteId
    
    log.Infof("Starting metrics for site ID: %s", siteId)
    
	scanInterval := 3
	log.Infof("Starting metrics collection goroutine for site %s with scan interval %d seconds", 
	siteId, scanInterval)
	profile := cenums.PROFILE_NORMAL
	if req.Profile == pb.Profile_PROFILE_MIN {
		profile = cenums.PROFILE_MIN
	} else if req.Profile == pb.Profile_PROFILE_MAX {
		profile = cenums.PROFILE_MAX
	}
	
	var scenario cenums.SCENARIOS
	switch req.Scenario {
	case pb.Scenario_SCENARIO_POWER_DOWN:
		scenario = cenums.SCENARIO_POWER_DOWN
	case pb.Scenario_SCENARIO_BACKHAUL_DOWN:
		scenario = cenums.SCENARIO_BACKHAUL_DOWN
	default:
		scenario = cenums.SCENARIO_DEFAULT
	}
	
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
	
	s.metricsProviders[siteId].SetScenario(string(scenario))
	
    s.metricsProviders[siteId].SetProfile(profile)
	
	siteCtx, cancelFunc := context.WithCancel(context.Background())
	
	exporter := metrics.NewPrometheusExporter(s.metricsProviders[siteId], siteId)
	
	s.siteConfigs[siteId] = &SiteMetricsConfig{
		ScanInterval: scanInterval,
		Profile:      profile,
		Scenario:     scenario,
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

func (s *ControllerServer) UpdateMetrics(ctx context.Context, req *pb.UpdateMetricsRequest) (*pb.UpdateMetricsResponse, error) {
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
        }, nil
    }

    // Handle profile changes
    if req.Profile != pb.Profile_PROFILE_NORMAL {
        if req.Profile == pb.Profile_PROFILE_MIN {
            config.Profile = cenums.PROFILE_MIN
            provider.SetProfile(cenums.PROFILE_MIN)
        } else if req.Profile == pb.Profile_PROFILE_MAX {
            config.Profile = cenums.PROFILE_MAX
            provider.SetProfile(cenums.PROFILE_MAX)
        }
    }

    // Handle scenario changes
    var scenarioChanged bool
    var newScenario cenums.SCENARIOS = config.Scenario

    switch req.Scenario {
    case pb.Scenario_SCENARIO_POWER_DOWN:
        newScenario = cenums.SCENARIO_POWER_DOWN
        config.Scenario = newScenario
        provider.SetScenario(string(newScenario))
        scenarioChanged = true
    case pb.Scenario_SCENARIO_BACKHAUL_DOWN:
        newScenario = cenums.SCENARIO_BACKHAUL_DOWN
        config.Scenario = newScenario
        provider.SetScenario(string(newScenario))
        scenarioChanged = true
    case pb.Scenario_SCENARIO_DEFAULT:
        if config.Scenario != cenums.SCENARIO_DEFAULT {
            newScenario = cenums.SCENARIO_DEFAULT
            config.Scenario = newScenario
            provider.SetScenario(string(newScenario))
            scenarioChanged = true
        }
    }

    // Track port updates and backhaul ports set to down
    var portUpdatesApplied bool
    var backhaulPortsDown []int
    if req.PortUpdates != nil {
        for _, portUpdate := range req.PortUpdates {
            portNumber := int(portUpdate.PortNumber)
            portStatus := portUpdate.Status

            err := provider.SetPortStatus(portNumber, portStatus)
            if err != nil {
                log.Infof("Error updating port %d status: %v", portNumber, err)
            } else {
                portUpdatesApplied = true
                log.Infof("Updated port %d status to %v for site %s", portNumber, portStatus, siteId)

                // Track backhaul ports set to down
                if !portStatus && isBackhaulPort(portNumber) {
                    backhaulPortsDown = append(backhaulPortsDown, portNumber)
                }
            }
        }
    }

	// Fetch nodes only if scenario changed or backhaul ports were updated
	var nodes *creg.NodesBySite
	needNodes := scenarioChanged || len(backhaulPortsDown) > 0
    if needNodes {
        var err error
        nodes, err = s.nodeClient.GetNodesBySite(siteId)
        if err != nil {
            return &pb.UpdateMetricsResponse{
                Success: false,
                Message: "Site not found",
            }, nil
        }
    }

    // Update nodes if scenario changed
    if scenarioChanged {
        nodeIds := make([]string, len(nodes.Nodes))
        for i, node := range nodes.Nodes {
            nodeIds[i] = node.Id
        }
        go func(nodeList []string, scenario cenums.SCENARIOS, profile cenums.Profile) {
            for _, nodeID := range nodeList {
                err := s.dnodeClient.UpdateNodeScenario(nodeID, scenario, profile)
                if err != nil {
                    log.Errorf("Failed to update node %s scenario: %v", nodeID, err)
                }
            }
        }(nodeIds, newScenario, config.Profile)
    }

    // Notify nodes of backhaul port down events
    if len(backhaulPortsDown) > 0 {
        nodeIDs := make([]string, len(nodes.Nodes))
        for i, node := range nodes.Nodes {
            nodeIDs[i] = node.Id
        }
        for _, portNumber := range backhaulPortsDown {
            go func(portNum int, nodeList []string) {
                for _, nodeID := range nodeList {
                    err := s.dnodeClient.NotifyNodeBackhaulDown(nodeID)
                    if err != nil {
                        log.Errorf("Failed to notify node %s of backhaul down (port %d): %v", nodeID, portNum, err)
                    }
                }
            }(portNumber, nodeIDs)
        }
    }

    // Build response message
    statusMessage := "Updated metrics configuration"
    if scenarioChanged {
        statusMessage += fmt.Sprintf(" - Scenario set to %s", config.Scenario)
    }
    if portUpdatesApplied {
        statusMessage += " - Port updates applied"
    }

    return &pb.UpdateMetricsResponse{
        Success: true,
        Message: statusMessage,
    }, nil
}
func (s *ControllerServer) StopMetricsCollection(siteId string) bool {
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

func (s *ControllerServer) Cleanup() {
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
func isBackhaulPort(portNumber int) bool {

	backhaulPorts := map[int]bool{
		1: true, 
		2: true, 
	}
	
	return backhaulPorts[portNumber]
}
func (s *ControllerServer) MonitorPowerStatus(ctx context.Context, siteID string) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.mutex.RLock()
			provider, exists := s.metricsProviders[siteID]
			s.mutex.RUnlock()
			
			if !exists {
				continue
			}
			
			powerStatus, err := provider.GetPowerStatus()
			if err != nil {
				log.Errorf("Failed to get power status for site %s: %v", siteID, err)
				continue
			}
			
			// If power is down, notify all nodes
			if !powerStatus {
				nodes, err := s.nodeClient.GetNodesBySite(siteID)
				if err != nil {
					log.Errorf("Failed to get nodes for site %s: %v", siteID, err)
					continue
				}
				
				for _, node := range nodes.Nodes {
					err := s.dnodeClient.NotifyNodePowerDown(node.Id)
					if err != nil {
						log.Errorf("Failed to notify node %s of power down: %v", node.Id, err)
					}
				}
			}
		}
	}
}