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

type Ranges struct {
	Min    float64 `json:"min"`
	Normal float64 `json:"normal"`
	Max    float64 `json:"max"`
}

type Config struct {
	UnitUptime         Ranges `json:"{min:0,normal:2678400,max:2678400}"`
	UnitHealth         Ranges `json:"{min:10,normal:80,max:100}"`
	TrxLteCoreActiveUE Ranges `json:"{min:80,normal:95,max:100}"`
	NodeLoad           Ranges `json:"{min:10,normal:80,max:100}"`
	CellularUplink     Ranges `json:"{min:2,normal:5,max:30}"`
	CellularDownlink   Ranges `json:"{min:8,normal:60,max:160}"`
	BackhaulUplink     Ranges `json:"{min:1,normal:10,max:200}"`
	BackhaulDownlink   Ranges `json:"{min:1,normal:10,max:200}"`
	BackhaulLatency    Ranges `json:"{min:10,normal:800,max:1000}"`
	HwdLoad            Ranges `json:"{min:10,normal:80,max:100}"`
	MemoryUsage        Ranges `json:"{min:10,normal:80,max:100}"`
	CpuUsage           Ranges `json:"{min:10,normal:80,max:100}"`
	DiskUsage          Ranges `json:"{min:10,normal:80,max:100}"`
	TxPower            Ranges `json:"{min:25,normal:31,max:34}"`
}

type WMessage struct {
	Kpis     NodeKPIs         `json:"kpis"`
	NodeId   string           `json:"nodeId"`
	Profile  cenums.Profile   `json:"profile"`
	Scenario cenums.SCENARIOS `json:"scenario"`
}

type NodeKPI struct {
	Id     string
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
			Id:  "UnitUptime",
			Key: "unit_uptime",
			KPI: prometheus.NewGaugeVec(
				prometheus.GaugeOpts{
					Name: "unit_uptime",
					Help: "Node uptime",
				},
				[]string{"nodeid"},
			),
		},
		{
			Id:  "UnitHealth",
			Key: "unit_health",
			KPI: prometheus.NewGaugeVec(
				prometheus.GaugeOpts{
					Name: "unit_health",
					Help: "Health status of the unit",
				},
				[]string{"nodeid"},
			),
		},
		{
			Id:  "TrxLteCoreActiveUE",
			Key: "trx_lte_core_active_ue",
			KPI: prometheus.NewGaugeVec(
				prometheus.GaugeOpts{
					Name: "trx_lte_core_active_ue",
					Help: "Active subscriber within the network",
				},
				[]string{"nodeid"},
			),
		},
		{
			Id:  "NodeLoad",
			Key: "node_load",
			KPI: prometheus.NewGaugeVec(
				prometheus.GaugeOpts{
					Name: "node_load",
					Help: "Load on the node",
				},
				[]string{"nodeid"},
			),
		},
		{
			Id:  "CellularUplink",
			Key: "cellular_uplink",
			KPI: prometheus.NewGaugeVec(
				prometheus.GaugeOpts{
					Name: "cellular_uplink",
					Help: "Cellular uplink",
				},
				[]string{"nodeid"},
			),
		},
		{
			Id:  "CellularDownlink",
			Key: "cellular_downlink",
			KPI: prometheus.NewGaugeVec(
				prometheus.GaugeOpts{
					Name: "cellular_downlink",
					Help: "Cellular downlink",
				},
				[]string{"nodeid"},
			),
		},
		{
			Id:  "BackhaulUplink",
			Key: "backhaul_uplink",
			KPI: prometheus.NewGaugeVec(
				prometheus.GaugeOpts{
					Name: "backhaul_uplink",
					Help: "Backhaul uplink",
				},
				[]string{"nodeid"},
			),
		},
		{
			Id:  "BackhaulDownlink",
			Key: "backhaul_downlink",
			KPI: prometheus.NewGaugeVec(
				prometheus.GaugeOpts{
					Name: "backhaul_downlink",
					Help: "Backhaul downlink",
				},
				[]string{"nodeid"},
			),
		},
		{
			Id:  "BackhaulLatency",
			Key: "backhaul_latency",
			KPI: prometheus.NewGaugeVec(
				prometheus.GaugeOpts{
					Name: "backhaul_latency",
					Help: "Backhaul latency",
				},
				[]string{"nodeid"},
			),
		},
		{
			Id:  "HwdLoad",
			Key: "hwd_load",
			KPI: prometheus.NewGaugeVec(
				prometheus.GaugeOpts{
					Name: "hwd_load",
					Help: "Hardware load",
				},
				[]string{"nodeid"},
			),
		},
		{
			Id:  "MemoryUsage",
			Key: "memory_usage",
			KPI: prometheus.NewGaugeVec(
				prometheus.GaugeOpts{
					Name: "memory_usage",
					Help: "Memory usage",
				},
				[]string{"nodeid"},
			),
		},
		{
			Id:  "CpuUsage",
			Key: "cpu_usage",
			KPI: prometheus.NewGaugeVec(
				prometheus.GaugeOpts{
					Name: "cpu_usage",
					Help: "CPU usage",
				},
				[]string{"nodeid"},
			),
		},
		{
			Id:  "DiskUsage",
			Key: "disk_usage",
			KPI: prometheus.NewGaugeVec(
				prometheus.GaugeOpts{
					Name: "disk_usage",
					Help: "Disk usage",
				},
				[]string{"nodeid"},
			),
		},
		{
			Id:  "TxPower",
			Key: "txpower",
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
