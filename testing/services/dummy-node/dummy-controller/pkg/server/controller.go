package server

import (
	"context"

	pb "github.com/ukama/ukama/testing/services/dummy-node/dummy-controller/pb/gen"
	"github.com/ukama/ukama/testing/services/dummy-node/dummy-controller/pkg/metrics"
)

type ControllerServer struct {
	pb.UnimplementedMetricsControllerServer
	orgName        string
	metricsProvider *metrics.MetricsProvider
	siteId          string
}

func NewControllerServer(orgName string, siteId string) *ControllerServer {
	return &ControllerServer{
		orgName:         orgName,
		metricsProvider: metrics.NewMetricsProvider(),
		siteId:          siteId,
	}
}


func (s *ControllerServer) GetSiteMetrics(ctx context.Context, req *pb.GetSiteMetricsRequest) (*pb.GetSiteMetricsResponse, error) {
	systemMetrics, err := s.metricsProvider.GetMetrics(s.siteId)
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
			Health:       map[bool]float64{true: 1.0, false: 0.0}[systemMetrics.Battery.Health == "Good"],
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