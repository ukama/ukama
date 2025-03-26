/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
package config

import (
	"os"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	cenums "github.com/ukama/ukama/testing/common/enums"

	uconf "github.com/ukama/ukama/systems/common/config"
)

const PORT = 8085

type Config struct {
	uconf.BaseConfig `mapstructure:",squash"`
	Port             int `default:"8085"`
	KpiConfig        NodeKPIs
	KpiRanges        KPIRanges
}

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

type Ranges struct {
	Min    float64 `json:"min"`
	Normal float64 `json:"normal"`
	Max    float64 `json:"max"`
}

type KPIRanges struct {
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
	TxPower            Ranges
}

func NewConfig() *Config {
	return &Config{
		KpiConfig: NodeKPIs{
			KPIs: []NodeKPI{
				{
					Key:    "unit_uptime",
					Min:    getConfigValue("KPIRANGES_UNITUPTIME_MIN", 0),
					Normal: getConfigValue("KPIRANGES_UNITUPTIME_NORMAL", 2678400),
					Max:    getConfigValue("KPIRANGES_UNITUPTIME_MAX", 3678400),
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
					Min:    getConfigValue("KPIRANGES_UNITHEALTH_MIN", 10),
					Normal: getConfigValue("KPIRANGES_UNITHEALTH_NORMAL", 80),
					Max:    getConfigValue("KPIRANGES_UNITHEALTH_MAX", 100),
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
					Min:    getConfigValue("KPIRANGES_TRXLTECOREACTIVEUE_MIN", 80),
					Normal: getConfigValue("KPIRANGES_TRXLTECOREACTIVEUE_NORMAL", 95),
					Max:    getConfigValue("KPIRANGES_TRXLTECOREACTIVEUE_MAX", 100),
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
					Min:    getConfigValue("KPIRANGES_NODELOAD_MIN", 10),
					Normal: getConfigValue("KPIRANGES_NODELOAD_NORMAL", 80),
					Max:    getConfigValue("KPIRANGES_NODELOAD_MAX", 100),
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
					Min:    getConfigValue("KPIRANGES_CELLULARUP_MIN", 2),
					Normal: getConfigValue("KPIRANGES_CELLULARUP_NORMAL", 5),
					Max:    getConfigValue("KPIRANGES_CELLULARUP_MAX", 30),
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
					Min:    getConfigValue("KPIRANGES_CELLULARDOWN_MIN", 8),
					Normal: getConfigValue("KPIRANGES_CELLULARDOWN_NORMAL", 60),
					Max:    getConfigValue("KPIRANGES_CELLULARDOWN_MAX", 160),
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
					Min:    getConfigValue("KPIRANGES_BACKHAULUP_MIN", 2),
					Normal: getConfigValue("KPIRANGES_BACKHAULUP_NORMAL", 10),
					Max:    getConfigValue("KPIRANGES_BACKHAULUP_MAX", 200),
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
					Min:    getConfigValue("KPIRANGES_BACKHAULDOWN_MIN", 2),
					Normal: getConfigValue("KPIRANGES_BACKHAULDOWN_NORMAL", 10),
					Max:    getConfigValue("KPIRANGES_BACKHAULDOWN_MAX", 200),
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
					Min:    getConfigValue("KPIRANGES_BACKHAULLATENCY_MIN", 10),
					Normal: getConfigValue("KPIRANGES_BACKHAULLATENCY_NORMAL", 800),
					Max:    getConfigValue("KPIRANGES_BACKHAULLATENCY_MAX", 1000),
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
					Min:    getConfigValue("KPIRANGES_HWLOAD_MIN", 10),
					Normal: getConfigValue("KPIRANGES_HWLOAD_NORMAL", 80),
					Max:    getConfigValue("KPIRANGES_HWLOAD_MAX", 100),
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
					Min:    getConfigValue("KPIRANGES_MEMORYUSAGE_MIN", 10),
					Normal: getConfigValue("KPIRANGES_MEMORYUSAGE_NORMAL", 80),
					Max:    getConfigValue("KPIRANGES_MEMORYUSAGE_MAX", 100),
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
					Min:    getConfigValue("KPIRANGES_CPUUSAGE_MIN", 10),
					Normal: getConfigValue("KPIRANGES_CPUUSAGE_NORMAL", 80),
					Max:    getConfigValue("KPIRANGES_CPUUSAGE_MAX", 100),
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
					Min:    getConfigValue("KPIRANGES_DISKUSAGE_MIN", 10),
					Normal: getConfigValue("KPIRANGES_DISKUSAGE_NORMAL", 80),
					Max:    getConfigValue("KPIRANGES_DISKUSAGE_MAX", 100),
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
					Min:    getConfigValue("KPIRANGES_TXPOWER_MIN", 25),
					Normal: getConfigValue("KPIRANGES_TXPOWER_NORMAL", 31),
					Max:    getConfigValue("KPIRANGES_TXPOWER_MAX", 34),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "txpower",
							Help: "Transmit power",
						},
						[]string{"nodeid"},
					),
				},
			},
		},
	}
}

func getConfigValue(key string, d float64) float64 {
	v := os.Getenv(key)
	if v == "" {
		return d
	}
	if val, err := strconv.ParseFloat(v, 64); err == nil {
		return val
	}
	return d
}
