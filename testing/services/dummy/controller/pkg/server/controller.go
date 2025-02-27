package server

import (
	"context"
	"fmt"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	cenums "github.com/ukama/ukama/testing/common/enums"
	pb "github.com/ukama/ukama/testing/services/dummy/controller/pb/gen"
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
}

func NewControllerServer(orgName string) *ControllerServer {
	return &ControllerServer{
		orgName:          orgName,
		metricsProviders: make(map[string]*metrics.MetricsProvider),
		siteConfigs:      make(map[string]*SiteMetricsConfig),
		mutex:            sync.RWMutex{},
	}
}

func (s *ControllerServer) GetSiteMetrics(ctx context.Context, req *pb.GetSiteMetricsRequest) (*pb.GetSiteMetricsResponse, error) {
	var siteId string
	var provider *metrics.MetricsProvider
	
	s.mutex.RLock()
	if len(s.metricsProviders) > 0 {
		for id, p := range s.metricsProviders {
			siteId = id
			provider = p
			break
		}
	} else {
		s.mutex.RUnlock()
		return nil, fmt.Errorf("no site metrics available")
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
	
	if req.Profile != pb.Profile_PROFILE_NORMAL {
		if req.Profile == pb.Profile_PROFILE_MIN {
			config.Profile = cenums.PROFILE_MIN
			provider.SetProfile(cenums.PROFILE_MIN)
		} else if req.Profile == pb.Profile_PROFILE_MAX {
			config.Profile = cenums.PROFILE_MAX
			provider.SetProfile(cenums.PROFILE_MAX)
		}
	}
	
	var scenarioChanged bool
	switch req.Scenario {
	case pb.Scenario_SCENARIO_POWER_DOWN:
		config.Scenario = cenums.SCENARIO_POWER_DOWN
		provider.SetScenario(string(cenums.SCENARIO_POWER_DOWN))
		scenarioChanged = true
	case pb.Scenario_SCENARIO_BACKHAUL_DOWN:
		config.Scenario = cenums.SCENARIO_BACKHAUL_DOWN
		provider.SetScenario(string(cenums.SCENARIO_BACKHAUL_DOWN))
		scenarioChanged = true
	case pb.Scenario_SCENARIO_DEFAULT:
		if config.Scenario != cenums.SCENARIO_DEFAULT {
			config.Scenario = cenums.SCENARIO_DEFAULT
			provider.SetScenario(string(cenums.SCENARIO_DEFAULT))
			scenarioChanged = true
		}
	}
	
	if req.PortUpdates != nil {
		for _, portUpdate := range req.PortUpdates {
			portNumber := int(portUpdate.PortNumber)
			portStatus := portUpdate.Status
			
			err := provider.SetPortStatus(portNumber, portStatus)
			if err != nil {
				log.Infof("Error updating port %d status: %v", portNumber, err)
			} else {
				log.Infof("Updated port %d status to %v for site %s", 
					portNumber, portStatus, siteId)
			}
		}
	}

	statusMessage := "Updated metrics configuration"
	if scenarioChanged {
		statusMessage += fmt.Sprintf(" - Scenario set to %s", config.Scenario)
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