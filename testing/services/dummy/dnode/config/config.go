/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
package config

import (
	"github.com/prometheus/client_golang/prometheus"
	cenums "github.com/ukama/ukama/testing/common/enums"
)

const PORT = 8085

type WMessage struct {
	Kpis     NodeKPIs         `json:"kpis"`
	NodeId   string           `json:"nodeId"`
	Profile  cenums.Profile   `json:"profile"`
	Scenario cenums.SCENARIOS `json:"scenario"`
}

type NodeKPI struct {
	Key    string
	Min    float64
	Normal float64
	Max    float64
	KPI    *prometheus.GaugeVec
}

type NodeKPIs struct {
	KPIs []NodeKPI
}

var KPI_CONFIG = NodeKPIs{
	KPIs: []NodeKPI{
		{
			Key: "unit_uptime",
			Min: 0,
			Max: 2678400,
			KPI: prometheus.NewGaugeVec(
				prometheus.GaugeOpts{
					Name: "unit_uptime",
					Help: "Node uptime",
				},
				[]string{"nodeid"},
			),
		},
		{
			Key:    "unit_health",
			Min:    0,
			Normal: 80,
			Max:    100,
			KPI: prometheus.NewGaugeVec(
				prometheus.GaugeOpts{
					Name: "unit_health",
					Help: "Health status of the unit",
				},
				[]string{"nodeid"},
			),
		},
		{
			Key:    "trx_lte_core_active_ue",
			Min:    80,
			Normal: 95,
			Max:    100,
			KPI: prometheus.NewGaugeVec(
				prometheus.GaugeOpts{
					Name: "trx_lte_core_active_ue",
					Help: "Active subscriber within the network",
				},
				[]string{"nodeid"},
			),
		},
		{
			Key:    "node_load",
			Min:    0,
			Normal: 80,
			Max:    100,
			KPI: prometheus.NewGaugeVec(
				prometheus.GaugeOpts{
					Name: "node_load",
					Help: "Load on the node",
				},
				[]string{"nodeid"},
			),
		},
		{
			Key:    "cellular_uplink",
			Min:    0,
			Normal: 5,
			Max:    30,
			KPI: prometheus.NewGaugeVec(
				prometheus.GaugeOpts{
					Name: "cellular_uplink",
					Help: "Cellular uplink",
				},
				[]string{"nodeid"},
			),
		},
		{
			Key:    "cellular_downlink",
			Min:    0,
			Normal: 60,
			Max:    160,
			KPI: prometheus.NewGaugeVec(
				prometheus.GaugeOpts{
					Name: "cellular_downlink",
					Help: "Cellular downlink",
				},
				[]string{"nodeid"},
			),
		},
		{
			Key:    "backhaul_uplink",
			Min:    0,
			Normal: 10,
			Max:    200,
			KPI: prometheus.NewGaugeVec(
				prometheus.GaugeOpts{
					Name: "backhaul_uplink",
					Help: "Backhaul uplink",
				},
				[]string{"nodeid"},
			),
		},
		{
			Key:    "backhaul_downlink",
			Min:    0,
			Normal: 10,
			Max:    200,
			KPI: prometheus.NewGaugeVec(
				prometheus.GaugeOpts{
					Name: "backhaul_downlink",
					Help: "Backhaul downlink",
				},
				[]string{"nodeid"},
			),
		},
		{
			Key:    "backhaul_latency",
			Min:    0,
			Normal: 800,
			Max:    1000,
			KPI: prometheus.NewGaugeVec(
				prometheus.GaugeOpts{
					Name: "backhaul_latency",
					Help: "Backhaul latency",
				},
				[]string{"nodeid"},
			),
		},
		{
			Key:    "hwd_load",
			Min:    0,
			Normal: 80,
			Max:    100,
			KPI: prometheus.NewGaugeVec(
				prometheus.GaugeOpts{
					Name: "hwd_load",
					Help: "Hardware load",
				},
				[]string{"nodeid"},
			),
		},
		{
			Key:    "memory_usage",
			Min:    0,
			Normal: 80,
			Max:    100,
			KPI: prometheus.NewGaugeVec(
				prometheus.GaugeOpts{
					Name: "memory_usage",
					Help: "Memory usage",
				},
				[]string{"nodeid"},
			),
		},
		{
			Key:    "cpu_usage",
			Min:    0,
			Normal: 80,
			Max:    100,
			KPI: prometheus.NewGaugeVec(
				prometheus.GaugeOpts{
					Name: "cpu_usage",
					Help: "CPU usage",
				},
				[]string{"nodeid"},
			),
		},
		{
			Key:    "disk_usage",
			Min:    0,
			Normal: 80,
			Max:    100,
			KPI: prometheus.NewGaugeVec(
				prometheus.GaugeOpts{
					Name: "disk_usage",
					Help: "Disk usage",
				},
				[]string{"nodeid"},
			),
		},
		{
			Key:    "txpower",
			Min:    0,
			Normal: 31,
			Max:    34,
			KPI: prometheus.NewGaugeVec(
				prometheus.GaugeOpts{
					Name: "txpower",
					Help: "Transmit power",
				},
				[]string{"nodeid"},
			),
		},
	},
}
