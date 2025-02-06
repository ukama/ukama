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

	"github.com/ukama/ukama/testing/services/dummy-node/dummy-controller/internal/backhaul"
	"github.com/ukama/ukama/testing/services/dummy-node/dummy-controller/internal/battery"
	"github.com/ukama/ukama/testing/services/dummy-node/dummy-controller/internal/solar"
	pb "github.com/ukama/ukama/testing/services/dummy-node/dummy-controller/pb/gen"
)

type ControllerServer struct {
	pb.UnimplementedMetricsControllerServer
	orgName        string
	solarProvider    *solar.SolarProvider
	backhaulProvider *backhaul.BackhaulProvider
	batteryProvider *battery.MockBatteryProvider
}

func NewControllerServer(orgName string, solarProvider *solar.SolarProvider, backhaulProvider *backhaul.BackhaulProvider, MockBatteryProvider *battery.MockBatteryProvider) *ControllerServer {
	return &ControllerServer{
		orgName:        orgName,
		solarProvider:    solarProvider,
		backhaulProvider: backhaulProvider,
		batteryProvider: MockBatteryProvider,
	}
}

func (s *ControllerServer) GetMetrics(ctx context.Context, req *pb.GetSiteMetricsRequest) (*pb.GetSiteMetricsResponse, error) {
	// Get battery metrics
	batteryMetrics, err := battery.GetBatteryMetrics()
	if err != nil {
		return nil, err
	}

	// Get solar metrics
	solarMetrics := s.solarProvider.GetMetrics()

	// Get backhaul metrics
	backhaulMetrics := s.backhaulProvider.GetMetrics()

	return &pb.GetSiteMetricsResponse{
		Solar: &pb.SolarMetrics{
			PowerGeneration: solarMetrics.PowerGeneration,
			EnergyTotal:    solarMetrics.EnergyTotal,
			PanelPower:     solarMetrics.PanelPower,
			PanelCurrent:   solarMetrics.PanelCurrent,
			PanelVoltage:   solarMetrics.PanelVoltage,
			InverterStatus: solarMetrics.InverterStatus,
		},
		Battery: &pb.BatteryMetrics{
			ChargeStatus: batteryMetrics.Capacity,
			Voltage:      batteryMetrics.Voltage,
			Health:       map[bool]float64{true: 1.0, false: 0.0}[batteryMetrics.Health == "Good"],
			Current:      batteryMetrics.Current,
			Temperature:  batteryMetrics.Temperature,
		},
		Network: &pb.NetworkMetrics{
			BackhaulLatency: backhaulMetrics.Latency,
			BackhaulStatus:  backhaulMetrics.Status,
			BackhaulSpeed:   backhaulMetrics.Speed,
		},
	}, nil
}

