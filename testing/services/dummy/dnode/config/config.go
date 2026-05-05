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

type Config struct {
	uconf.BaseConfig `mapstructure:",squash"`
	Port             int `default:"8080"`
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
	Metric   string
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
	UnitUptime         Ranges
	UnitHealth         Ranges
	TrxLteCoreActiveUE Ranges
	NodeLoad           Ranges
	CellularUplink     Ranges
	BackhaulUplink     Ranges
	BackhaulDownlink   Ranges
	BackhaulLatency    Ranges
	HwdLoad            Ranges
	MemoryUsage        Ranges
	CpuUsage           Ranges
	DiskUsage          Ranges
	TxPower            Ranges
}

func NewConfig() *Config {
	return &Config{
		Port: 8080,
		KpiConfig: NodeKPIs{
			KPIs: []NodeKPI{
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
						[]string{"node_id", "metric"},
					),
					Metric: "trx_lte_core_active_ue",
				},
				{
					Key:    "com_generic_system_uptime_seconds",
					Min:    getConfigValue("KPIRANGES_NODELOAD_MIN", 10),
					Normal: getConfigValue("KPIRANGES_NODELOAD_NORMAL", 80),
					Max:    getConfigValue("KPIRANGES_NODELOAD_MAX", 100),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_generic_system_uptime_seconds",
							Help: "Load on the node",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_generic_system_uptime_seconds",
				},
				{
					Key:    "trx_lte_stack_throughput_uplink",
					Min:    getConfigValue("KPIRANGES_CELLULARUP_MIN", 2),
					Normal: getConfigValue("KPIRANGES_CELLULARUP_NORMAL", 5),
					Max:    getConfigValue("KPIRANGES_CELLULARUP_MAX", 30),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "trx_lte_stack_throughput_uplink",
							Help: "Cellular uplink",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "trx_lte_stack_throughput_uplink",
				},
				{
					Key:    "trx_lte_stack_throughput_downlink",
					Min:    getConfigValue("KPIRANGES_CELLULARDOWN_MIN", 8),
					Normal: getConfigValue("KPIRANGES_CELLULARDOWN_NORMAL", 60),
					Max:    getConfigValue("KPIRANGES_CELLULARDOWN_MAX", 160),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "trx_lte_stack_throughput_downlink",
							Help: "Cellular downlink",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "trx_lte_stack_throughput_downlink",
				},
				{
					Key:    "com_backhaul_backhaul_dl_goodput_mbps",
					Min:    getConfigValue("KPIRANGES_BACKHAULUP_MIN", 2),
					Normal: getConfigValue("KPIRANGES_BACKHAULUP_NORMAL", 10),
					Max:    getConfigValue("KPIRANGES_BACKHAULUP_MAX", 200),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_backhaul_backhaul_dl_goodput_mbps",
							Help: "Backhaul uplink",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_backhaul_backhaul_dl_goodput_mbps",
				},
				{
					Key:    "com_backhaul_backhaul_ul_goodput_mbps",
					Min:    getConfigValue("KPIRANGES_BACKHAULDOWN_MIN", 2),
					Normal: getConfigValue("KPIRANGES_BACKHAULDOWN_NORMAL", 10),
					Max:    getConfigValue("KPIRANGES_BACKHAULDOWN_MAX", 200),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_backhaul_backhaul_ul_goodput_mbps",
							Help: "Backhaul downlink",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_backhaul_backhaul_ul_goodput_mbps",
				},
				{
					Key:    "com_network_uplink_latency",
					Min:    getConfigValue("KPIRANGES_BACKHAULLATENCY_MIN", 10),
					Normal: getConfigValue("KPIRANGES_BACKHAULLATENCY_NORMAL", 800),
					Max:    getConfigValue("KPIRANGES_BACKHAULLATENCY_MAX", 1000),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_network_uplink_latency",
							Help: "Backhaul latency",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_network_uplink_latency",
				},
				{
					Key:    "cpu_temperature",
					Min:    getConfigValue("KPIRANGES_HWLOAD_MIN", 10),
					Normal: getConfigValue("KPIRANGES_HWLOAD_NORMAL", 80),
					Max:    getConfigValue("KPIRANGES_HWLOAD_MAX", 100),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "cpu_temperature",
							Help: "Hardware load",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "cpu_temperature",
				},
				{
					Key:    "memory",
					Min:    getConfigValue("KPIRANGES_MEMORYUSAGE_MIN", 10),
					Normal: getConfigValue("KPIRANGES_MEMORYUSAGE_NORMAL", 80),
					Max:    getConfigValue("KPIRANGES_MEMORYUSAGE_MAX", 100),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "memory",
							Help: "Trx memory usage",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "memory",
				},
				{
					Key:    "cpu",
					Min:    getConfigValue("KPIRANGES_CPUUSAGE_MIN", 10),
					Normal: getConfigValue("KPIRANGES_CPUUSAGE_NORMAL", 80),
					Max:    getConfigValue("KPIRANGES_CPUUSAGE_MAX", 100),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "cpu",
							Help: "Trx cpu usage",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "cpu",
				},
				{
					Key:    "disk",
					Min:    getConfigValue("KPIRANGES_DISKUSAGE_MIN", 10),
					Normal: getConfigValue("KPIRANGES_DISKUSAGE_NORMAL", 80),
					Max:    getConfigValue("KPIRANGES_DISKUSAGE_MAX", 100),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "disk",
							Help: "Trx disk usage",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "disk",
				},				
				{
					Key:    "com_power_power_power_total_watts",
					Min:    getConfigValue("KPIRANGES_TXPOWER_MIN", 25),
					Normal: getConfigValue("KPIRANGES_TXPOWER_NORMAL", 31),
					Max:    getConfigValue("KPIRANGES_TXPOWER_MAX", 34),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_power_power_power_total_watts",
							Help: "Total power",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_power_power_power_total_watts",
				},
				{
					Key:    "ctl_sensors_adc_tp_tx_power",
					Min:    getConfigValue("KPIRANGES_TXPOWER_MIN", 25),
					Normal: getConfigValue("KPIRANGES_TXPOWER_NORMAL", 31),
					Max:    getConfigValue("KPIRANGES_TXPOWER_MAX", 34),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "ctl_sensors_adc_tp_tx_power",
							Help: "TX power",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "ctl_sensors_adc_tp_tx_power",
				},
				{
					Key:    "ctl_sensors_adc_rp_rx_power",
					Min:    getConfigValue("KPIRANGES_TXPOWER_MIN", 25),
					Normal: getConfigValue("KPIRANGES_TXPOWER_NORMAL", 31),
					Max:    getConfigValue("KPIRANGES_TXPOWER_MAX", 34),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "ctl_sensors_adc_rp_rx_power",
							Help: "RX power",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "ctl_sensors_adc_rp_rx_power",
				},
				{
					Key:    "ctl_sensors_adc_pa_pa_power",
					Min:    getConfigValue("KPIRANGES_TXPOWER_MIN", 25),
					Normal: getConfigValue("KPIRANGES_TXPOWER_NORMAL", 31),
					Max:    getConfigValue("KPIRANGES_TXPOWER_MAX", 34),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "ctl_sensors_adc_pa_pa_power",
							Help: "PA power",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "ctl_sensors_adc_pa_pa_power",
				},
				{
					Key:    "ctl_sensors_femd_fem1_temperature_c",
					Min:    getConfigValue("KPIRANGES_TXPOWER_MIN", 25),
					Normal: getConfigValue("KPIRANGES_TXPOWER_NORMAL", 31),
					Max:    getConfigValue("KPIRANGES_TXPOWER_MAX", 34),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "ctl_sensors_femd_fem1_temperature_c",
							Help: "FEM1 temperature",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "ctl_sensors_femd_fem1_temperature_c",
				},
				{
					Key:    "ctl_sensors_femd_fem2_temperature_c",
					Min:    getConfigValue("KPIRANGES_TXPOWER_MIN", 25),
					Normal: getConfigValue("KPIRANGES_TXPOWER_NORMAL", 31),
					Max:    getConfigValue("KPIRANGES_TXPOWER_MAX", 34),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "ctl_sensors_femd_fem2_temperature_c",
							Help: "FEM2 temperature",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "ctl_sensors_femd_fem2_temperature_c",
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
