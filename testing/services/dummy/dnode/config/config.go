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
	Metric string
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

// KPIs are generated from systems/metrics/api-gateway/pkg/default-metrics.yaml:
// one entry per unique Prometheus `metric:` series, with thresholds taken from
// that file (env-overridable via KPIRANGES_<SERIES>_{MIN,NORMAL,MAX}).
func NewConfig() *Config {
	return &Config{
		Port: 8080,
		KpiConfig: NodeKPIs{
			KPIs: []NodeKPI{
				{
					Key:    "com_soc_cpu_usage",
					Min:    getConfigValue("KPIRANGES_COM_SOC_CPU_USAGE_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_SOC_CPU_USAGE_NORMAL", 80),
					Max:    getConfigValue("KPIRANGES_COM_SOC_CPU_USAGE_MAX", 100),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_soc_cpu_usage",
							Help: "cpu (tnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_soc_cpu_usage",
				},
				{
					Key:    "com_power_power_power_board_temp_c_celsius",
					Min:    getConfigValue("KPIRANGES_COM_POWER_POWER_POWER_BOARD_TEMP_C_CELSIUS_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_POWER_POWER_POWER_BOARD_TEMP_C_CELSIUS_NORMAL", 65),
					Max:    getConfigValue("KPIRANGES_COM_POWER_POWER_POWER_BOARD_TEMP_C_CELSIUS_MAX", 80),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_power_power_power_board_temp_c_celsius",
							Help: "cpu_temperature (tnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_power_power_power_board_temp_c_celsius",
				},
				{
					Key:    "com_sensors_thermal_zone0_temp_temperature_celsius",
					Min:    getConfigValue("KPIRANGES_COM_SENSORS_THERMAL_ZONE0_TEMP_TEMPERATURE_CELSIUS_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_SENSORS_THERMAL_ZONE0_TEMP_TEMPERATURE_CELSIUS_NORMAL", 65),
					Max:    getConfigValue("KPIRANGES_COM_SENSORS_THERMAL_ZONE0_TEMP_TEMPERATURE_CELSIUS_MAX", 80),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_sensors_thermal_zone0_temp_temperature_celsius",
							Help: "temperature_z0 (tnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_sensors_thermal_zone0_temp_temperature_celsius",
				},
				{
					Key:    "com_sensors_thermal_zone1_temp_temperature_celsius",
					Min:    getConfigValue("KPIRANGES_COM_SENSORS_THERMAL_ZONE1_TEMP_TEMPERATURE_CELSIUS_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_SENSORS_THERMAL_ZONE1_TEMP_TEMPERATURE_CELSIUS_NORMAL", 65),
					Max:    getConfigValue("KPIRANGES_COM_SENSORS_THERMAL_ZONE1_TEMP_TEMPERATURE_CELSIUS_MAX", 80),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_sensors_thermal_zone1_temp_temperature_celsius",
							Help: "temperature_z1 (tnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_sensors_thermal_zone1_temp_temperature_celsius",
				},
				{
					Key:    "com_sensors_thermal_zone2_temp_temperature_celsius",
					Min:    getConfigValue("KPIRANGES_COM_SENSORS_THERMAL_ZONE2_TEMP_TEMPERATURE_CELSIUS_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_SENSORS_THERMAL_ZONE2_TEMP_TEMPERATURE_CELSIUS_NORMAL", 65),
					Max:    getConfigValue("KPIRANGES_COM_SENSORS_THERMAL_ZONE2_TEMP_TEMPERATURE_CELSIUS_MAX", 80),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_sensors_thermal_zone2_temp_temperature_celsius",
							Help: "temperature_z2 (tnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_sensors_thermal_zone2_temp_temperature_celsius",
				},
				{
					Key:    "com_sensors_thermal_zone3_temp_temperature_celsius",
					Min:    getConfigValue("KPIRANGES_COM_SENSORS_THERMAL_ZONE3_TEMP_TEMPERATURE_CELSIUS_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_SENSORS_THERMAL_ZONE3_TEMP_TEMPERATURE_CELSIUS_NORMAL", 65),
					Max:    getConfigValue("KPIRANGES_COM_SENSORS_THERMAL_ZONE3_TEMP_TEMPERATURE_CELSIUS_MAX", 80),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_sensors_thermal_zone3_temp_temperature_celsius",
							Help: "temperature_z3 (tnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_sensors_thermal_zone3_temp_temperature_celsius",
				},
				{
					Key:    "com_generic_system_uptime_seconds",
					Min:    getConfigValue("KPIRANGES_COM_GENERIC_SYSTEM_UPTIME_SECONDS_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_GENERIC_SYSTEM_UPTIME_SECONDS_NORMAL", 0),
					Max:    getConfigValue("KPIRANGES_COM_GENERIC_SYSTEM_UPTIME_SECONDS_MAX", 0),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_generic_system_uptime_seconds",
							Help: "uptime (tnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_generic_system_uptime_seconds",
				},
				{
					Key:    "com_memory_ddr_used",
					Min:    getConfigValue("KPIRANGES_COM_MEMORY_DDR_USED_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_MEMORY_DDR_USED_NORMAL", 80),
					Max:    getConfigValue("KPIRANGES_COM_MEMORY_DDR_USED_MAX", 100),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_memory_ddr_used",
							Help: "memory (tnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_memory_ddr_used",
				},
				{
					Key:    "com_storage_emmc_used",
					Min:    getConfigValue("KPIRANGES_COM_STORAGE_EMMC_USED_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_STORAGE_EMMC_USED_NORMAL", 80),
					Max:    getConfigValue("KPIRANGES_COM_STORAGE_EMMC_USED_MAX", 100),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_storage_emmc_used",
							Help: "disk (tnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_storage_emmc_used",
				},
				{
					Key:    "com_power_power_power_total_watts",
					Min:    getConfigValue("KPIRANGES_COM_POWER_POWER_POWER_TOTAL_WATTS_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_POWER_POWER_POWER_TOTAL_WATTS_NORMAL", 60),
					Max:    getConfigValue("KPIRANGES_COM_POWER_POWER_POWER_TOTAL_WATTS_MAX", 100),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_power_power_power_total_watts",
							Help: "power (tnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_power_power_power_total_watts",
				},
				{
					Key:    "com_backhaul_backhaul_ul_goodput_mbps",
					Min:    getConfigValue("KPIRANGES_COM_BACKHAUL_BACKHAUL_UL_GOODPUT_MBPS_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_BACKHAUL_BACKHAUL_UL_GOODPUT_MBPS_NORMAL", 100),
					Max:    getConfigValue("KPIRANGES_COM_BACKHAUL_BACKHAUL_UL_GOODPUT_MBPS_MAX", 200),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_backhaul_backhaul_ul_goodput_mbps",
							Help: "backhaul_uplink (tnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_backhaul_backhaul_ul_goodput_mbps",
				},
				{
					Key:    "com_backhaul_backhaul_dl_goodput_mbps",
					Min:    getConfigValue("KPIRANGES_COM_BACKHAUL_BACKHAUL_DL_GOODPUT_MBPS_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_BACKHAUL_BACKHAUL_DL_GOODPUT_MBPS_NORMAL", 100),
					Max:    getConfigValue("KPIRANGES_COM_BACKHAUL_BACKHAUL_DL_GOODPUT_MBPS_MAX", 200),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_backhaul_backhaul_dl_goodput_mbps",
							Help: "backhaul_downlink (tnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_backhaul_backhaul_dl_goodput_mbps",
				},
				{
					Key:    "com_network_uplink_latency",
					Min:    getConfigValue("KPIRANGES_COM_NETWORK_UPLINK_LATENCY_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_NETWORK_UPLINK_LATENCY_NORMAL", 800),
					Max:    getConfigValue("KPIRANGES_COM_NETWORK_UPLINK_LATENCY_MAX", 1000),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_network_uplink_latency",
							Help: "backhaul_latency (tnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_network_uplink_latency",
				},
				{
					Key:    "trx_lte_stack_throughput_uplink",
					Min:    getConfigValue("KPIRANGES_TRX_LTE_STACK_THROUGHPUT_UPLINK_MIN", 0),
					Normal: getConfigValue("KPIRANGES_TRX_LTE_STACK_THROUGHPUT_UPLINK_NORMAL", 30),
					Max:    getConfigValue("KPIRANGES_TRX_LTE_STACK_THROUGHPUT_UPLINK_MAX", 100),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "trx_lte_stack_throughput_uplink",
							Help: "cellular_uplink (tnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "trx_lte_stack_throughput_uplink",
				},
				{
					Key:    "trx_lte_stack_throughput_downlink",
					Min:    getConfigValue("KPIRANGES_TRX_LTE_STACK_THROUGHPUT_DOWNLINK_MIN", 0),
					Normal: getConfigValue("KPIRANGES_TRX_LTE_STACK_THROUGHPUT_DOWNLINK_NORMAL", 60),
					Max:    getConfigValue("KPIRANGES_TRX_LTE_STACK_THROUGHPUT_DOWNLINK_MAX", 160),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "trx_lte_stack_throughput_downlink",
							Help: "cellular_downlink (tnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "trx_lte_stack_throughput_downlink",
				},
				{
					Key:    "trx_lte_stack_signal_strength",
					Min:    getConfigValue("KPIRANGES_TRX_LTE_STACK_SIGNAL_STRENGTH_MIN", -120),
					Normal: getConfigValue("KPIRANGES_TRX_LTE_STACK_SIGNAL_STRENGTH_NORMAL", -95),
					Max:    getConfigValue("KPIRANGES_TRX_LTE_STACK_SIGNAL_STRENGTH_MAX", -70),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "trx_lte_stack_signal_strength",
							Help: "lte_signal_strength (tnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "trx_lte_stack_signal_strength",
				},
				{
					Key:    "trx_lte_core_active_ue",
					Min:    getConfigValue("KPIRANGES_TRX_LTE_CORE_ACTIVE_UE_MIN", 0),
					Normal: getConfigValue("KPIRANGES_TRX_LTE_CORE_ACTIVE_UE_NORMAL", 100),
					Max:    getConfigValue("KPIRANGES_TRX_LTE_CORE_ACTIVE_UE_MAX", 1000),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "trx_lte_core_active_ue",
							Help: "lte_active_ue (tnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "trx_lte_core_active_ue",
				},
				{
					Key:    "trx_lte_core_subscribers",
					Min:    getConfigValue("KPIRANGES_TRX_LTE_CORE_SUBSCRIBERS_MIN", 0),
					Normal: getConfigValue("KPIRANGES_TRX_LTE_CORE_SUBSCRIBERS_NORMAL", 100),
					Max:    getConfigValue("KPIRANGES_TRX_LTE_CORE_SUBSCRIBERS_MAX", 1000),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "trx_lte_core_subscribers",
							Help: "lte_subscribers (tnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "trx_lte_core_subscribers",
				},
				{
					Key:    "trx_lte_cell_id",
					Min:    getConfigValue("KPIRANGES_TRX_LTE_CELL_ID_MIN", 0),
					Normal: getConfigValue("KPIRANGES_TRX_LTE_CELL_ID_NORMAL", 0),
					Max:    getConfigValue("KPIRANGES_TRX_LTE_CELL_ID_MAX", 0),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "trx_lte_cell_id",
							Help: "lte_cell_id (tnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "trx_lte_cell_id",
				},
				{
					Key:    "ctl_soc_cpu_usage",
					Min:    getConfigValue("KPIRANGES_CTL_SOC_CPU_USAGE_MIN", 0),
					Normal: getConfigValue("KPIRANGES_CTL_SOC_CPU_USAGE_NORMAL", 80),
					Max:    getConfigValue("KPIRANGES_CTL_SOC_CPU_USAGE_MAX", 100),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "ctl_soc_cpu_usage",
							Help: "cpu (anode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "ctl_soc_cpu_usage",
				},
				{
					Key:    "ctl_generic_system_uptime_seconds",
					Min:    getConfigValue("KPIRANGES_CTL_GENERIC_SYSTEM_UPTIME_SECONDS_MIN", 0),
					Normal: getConfigValue("KPIRANGES_CTL_GENERIC_SYSTEM_UPTIME_SECONDS_NORMAL", 0),
					Max:    getConfigValue("KPIRANGES_CTL_GENERIC_SYSTEM_UPTIME_SECONDS_MAX", 0),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "ctl_generic_system_uptime_seconds",
							Help: "uptime (anode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "ctl_generic_system_uptime_seconds",
				},
				{
					Key:    "ctl_sensors_tempsensor_proc_microprocessor",
					Min:    getConfigValue("KPIRANGES_CTL_SENSORS_TEMPSENSOR_PROC_MICROPROCESSOR_MIN", 0),
					Normal: getConfigValue("KPIRANGES_CTL_SENSORS_TEMPSENSOR_PROC_MICROPROCESSOR_NORMAL", 70),
					Max:    getConfigValue("KPIRANGES_CTL_SENSORS_TEMPSENSOR_PROC_MICROPROCESSOR_MAX", 85),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "ctl_sensors_tempsensor_proc_microprocessor",
							Help: "processor_temperature (anode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "ctl_sensors_tempsensor_proc_microprocessor",
				},
				{
					Key:    "ctl_memory_ddr_used",
					Min:    getConfigValue("KPIRANGES_CTL_MEMORY_DDR_USED_MIN", 0),
					Normal: getConfigValue("KPIRANGES_CTL_MEMORY_DDR_USED_NORMAL", 80),
					Max:    getConfigValue("KPIRANGES_CTL_MEMORY_DDR_USED_MAX", 100),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "ctl_memory_ddr_used",
							Help: "memory (anode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "ctl_memory_ddr_used",
				},
				{
					Key:    "ctl_storage_emmc_used",
					Min:    getConfigValue("KPIRANGES_CTL_STORAGE_EMMC_USED_MIN", 0),
					Normal: getConfigValue("KPIRANGES_CTL_STORAGE_EMMC_USED_NORMAL", 80),
					Max:    getConfigValue("KPIRANGES_CTL_STORAGE_EMMC_USED_MAX", 100),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "ctl_storage_emmc_used",
							Help: "disk (anode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "ctl_storage_emmc_used",
				},
				{
					Key:    "ctl_sensors_adc_pa_pa_power",
					Min:    getConfigValue("KPIRANGES_CTL_SENSORS_ADC_PA_PA_POWER_MIN", 0),
					Normal: getConfigValue("KPIRANGES_CTL_SENSORS_ADC_PA_PA_POWER_NORMAL", 30),
					Max:    getConfigValue("KPIRANGES_CTL_SENSORS_ADC_PA_PA_POWER_MAX", 60),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "ctl_sensors_adc_pa_pa_power",
							Help: "pa_power (anode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "ctl_sensors_adc_pa_pa_power",
				},
				{
					Key:    "ctl_sensors_adc_rp_rx_power",
					Min:    getConfigValue("KPIRANGES_CTL_SENSORS_ADC_RP_RX_POWER_MIN", -120),
					Normal: getConfigValue("KPIRANGES_CTL_SENSORS_ADC_RP_RX_POWER_NORMAL", -90),
					Max:    getConfigValue("KPIRANGES_CTL_SENSORS_ADC_RP_RX_POWER_MAX", -40),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "ctl_sensors_adc_rp_rx_power",
							Help: "rx_power (anode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "ctl_sensors_adc_rp_rx_power",
				},
				{
					Key:    "ctl_sensors_adc_tp_tx_power",
					Min:    getConfigValue("KPIRANGES_CTL_SENSORS_ADC_TP_TX_POWER_MIN", 0),
					Normal: getConfigValue("KPIRANGES_CTL_SENSORS_ADC_TP_TX_POWER_NORMAL", 31),
					Max:    getConfigValue("KPIRANGES_CTL_SENSORS_ADC_TP_TX_POWER_MAX", 34),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "ctl_sensors_adc_tp_tx_power",
							Help: "tx_power (anode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "ctl_sensors_adc_tp_tx_power",
				},
				{
					Key:    "ctl_sensors_femd_fem1_temperature_c",
					Min:    getConfigValue("KPIRANGES_CTL_SENSORS_FEMD_FEM1_TEMPERATURE_C_MIN", 0),
					Normal: getConfigValue("KPIRANGES_CTL_SENSORS_FEMD_FEM1_TEMPERATURE_C_NORMAL", 70),
					Max:    getConfigValue("KPIRANGES_CTL_SENSORS_FEMD_FEM1_TEMPERATURE_C_MAX", 90),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "ctl_sensors_femd_fem1_temperature_c",
							Help: "fem1_temperature (anode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "ctl_sensors_femd_fem1_temperature_c",
				},
				{
					Key:    "ctl_sensors_femd_fem1_forward_power",
					Min:    getConfigValue("KPIRANGES_CTL_SENSORS_FEMD_FEM1_FORWARD_POWER_MIN", 0),
					Normal: getConfigValue("KPIRANGES_CTL_SENSORS_FEMD_FEM1_FORWARD_POWER_NORMAL", 30),
					Max:    getConfigValue("KPIRANGES_CTL_SENSORS_FEMD_FEM1_FORWARD_POWER_MAX", 60),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "ctl_sensors_femd_fem1_forward_power",
							Help: "fem1_forward_power (anode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "ctl_sensors_femd_fem1_forward_power",
				},
				{
					Key:    "ctl_sensors_femd_fem1_reverse_power",
					Min:    getConfigValue("KPIRANGES_CTL_SENSORS_FEMD_FEM1_REVERSE_POWER_MIN", 0),
					Normal: getConfigValue("KPIRANGES_CTL_SENSORS_FEMD_FEM1_REVERSE_POWER_NORMAL", 20),
					Max:    getConfigValue("KPIRANGES_CTL_SENSORS_FEMD_FEM1_REVERSE_POWER_MAX", 40),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "ctl_sensors_femd_fem1_reverse_power",
							Help: "fem1_reverse_power (anode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "ctl_sensors_femd_fem1_reverse_power",
				},
				{
					Key:    "ctl_sensors_femd_fem1_pa_current",
					Min:    getConfigValue("KPIRANGES_CTL_SENSORS_FEMD_FEM1_PA_CURRENT_MIN", 0),
					Normal: getConfigValue("KPIRANGES_CTL_SENSORS_FEMD_FEM1_PA_CURRENT_NORMAL", 10),
					Max:    getConfigValue("KPIRANGES_CTL_SENSORS_FEMD_FEM1_PA_CURRENT_MAX", 20),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "ctl_sensors_femd_fem1_pa_current",
							Help: "fem1_pa_current (anode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "ctl_sensors_femd_fem1_pa_current",
				},
				{
					Key:    "ctl_sensors_femd_fem1_ok",
					Min:    getConfigValue("KPIRANGES_CTL_SENSORS_FEMD_FEM1_OK_MIN", 0),
					Normal: getConfigValue("KPIRANGES_CTL_SENSORS_FEMD_FEM1_OK_NORMAL", 1),
					Max:    getConfigValue("KPIRANGES_CTL_SENSORS_FEMD_FEM1_OK_MAX", 1),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "ctl_sensors_femd_fem1_ok",
							Help: "fem1_ok (anode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "ctl_sensors_femd_fem1_ok",
				},
				{
					Key:    "ctl_sensors_femd_fem2_temperature_c",
					Min:    getConfigValue("KPIRANGES_CTL_SENSORS_FEMD_FEM2_TEMPERATURE_C_MIN", 0),
					Normal: getConfigValue("KPIRANGES_CTL_SENSORS_FEMD_FEM2_TEMPERATURE_C_NORMAL", 70),
					Max:    getConfigValue("KPIRANGES_CTL_SENSORS_FEMD_FEM2_TEMPERATURE_C_MAX", 90),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "ctl_sensors_femd_fem2_temperature_c",
							Help: "fem2_temperature (anode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "ctl_sensors_femd_fem2_temperature_c",
				},
				{
					Key:    "ctl_sensors_femd_fem2_forward_power",
					Min:    getConfigValue("KPIRANGES_CTL_SENSORS_FEMD_FEM2_FORWARD_POWER_MIN", 0),
					Normal: getConfigValue("KPIRANGES_CTL_SENSORS_FEMD_FEM2_FORWARD_POWER_NORMAL", 30),
					Max:    getConfigValue("KPIRANGES_CTL_SENSORS_FEMD_FEM2_FORWARD_POWER_MAX", 60),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "ctl_sensors_femd_fem2_forward_power",
							Help: "fem2_forward_power (anode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "ctl_sensors_femd_fem2_forward_power",
				},
				{
					Key:    "ctl_sensors_femd_fem2_reverse_power",
					Min:    getConfigValue("KPIRANGES_CTL_SENSORS_FEMD_FEM2_REVERSE_POWER_MIN", 0),
					Normal: getConfigValue("KPIRANGES_CTL_SENSORS_FEMD_FEM2_REVERSE_POWER_NORMAL", 20),
					Max:    getConfigValue("KPIRANGES_CTL_SENSORS_FEMD_FEM2_REVERSE_POWER_MAX", 40),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "ctl_sensors_femd_fem2_reverse_power",
							Help: "fem2_reverse_power (anode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "ctl_sensors_femd_fem2_reverse_power",
				},
				{
					Key:    "ctl_sensors_femd_fem2_pa_current",
					Min:    getConfigValue("KPIRANGES_CTL_SENSORS_FEMD_FEM2_PA_CURRENT_MIN", 0),
					Normal: getConfigValue("KPIRANGES_CTL_SENSORS_FEMD_FEM2_PA_CURRENT_NORMAL", 10),
					Max:    getConfigValue("KPIRANGES_CTL_SENSORS_FEMD_FEM2_PA_CURRENT_MAX", 20),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "ctl_sensors_femd_fem2_pa_current",
							Help: "fem2_pa_current (anode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "ctl_sensors_femd_fem2_pa_current",
				},
				{
					Key:    "ctl_sensors_femd_fem2_ok",
					Min:    getConfigValue("KPIRANGES_CTL_SENSORS_FEMD_FEM2_OK_MIN", 0),
					Normal: getConfigValue("KPIRANGES_CTL_SENSORS_FEMD_FEM2_OK_NORMAL", 1),
					Max:    getConfigValue("KPIRANGES_CTL_SENSORS_FEMD_FEM2_OK_MAX", 1),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "ctl_sensors_femd_fem2_ok",
							Help: "fem2_ok (anode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "ctl_sensors_femd_fem2_ok",
				},
				{
					Key:    "com_backhaul_backhaul_state",
					Min:    getConfigValue("KPIRANGES_COM_BACKHAUL_BACKHAUL_STATE_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_BACKHAUL_BACKHAUL_STATE_NORMAL", 1),
					Max:    getConfigValue("KPIRANGES_COM_BACKHAUL_BACKHAUL_STATE_MAX", 1),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_backhaul_backhaul_state",
							Help: "backhaul_state (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_backhaul_backhaul_state",
				},
				{
					Key:    "com_backhaul_backhaul_link_guess",
					Min:    getConfigValue("KPIRANGES_COM_BACKHAUL_BACKHAUL_LINK_GUESS_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_BACKHAUL_BACKHAUL_LINK_GUESS_NORMAL", 1),
					Max:    getConfigValue("KPIRANGES_COM_BACKHAUL_BACKHAUL_LINK_GUESS_MAX", 1),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_backhaul_backhaul_link_guess",
							Help: "backhaul_link_guess (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_backhaul_backhaul_link_guess",
				},
				{
					Key:    "com_backhaul_backhaul_confidence",
					Min:    getConfigValue("KPIRANGES_COM_BACKHAUL_BACKHAUL_CONFIDENCE_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_BACKHAUL_BACKHAUL_CONFIDENCE_NORMAL", 70),
					Max:    getConfigValue("KPIRANGES_COM_BACKHAUL_BACKHAUL_CONFIDENCE_MAX", 100),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_backhaul_backhaul_confidence",
							Help: "backhaul_confidence (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_backhaul_backhaul_confidence",
				},
				{
					Key:    "com_controller_controller_solar_panel_voltage",
					Min:    getConfigValue("KPIRANGES_COM_CONTROLLER_CONTROLLER_SOLAR_PANEL_VOLTAGE_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_CONTROLLER_CONTROLLER_SOLAR_PANEL_VOLTAGE_NORMAL", 75),
					Max:    getConfigValue("KPIRANGES_COM_CONTROLLER_CONTROLLER_SOLAR_PANEL_VOLTAGE_MAX", 100),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_controller_controller_solar_panel_voltage",
							Help: "solar_panel_voltage (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_controller_controller_solar_panel_voltage",
				},
				{
					Key:    "com_controller_controller_solar_panel_current",
					Min:    getConfigValue("KPIRANGES_COM_CONTROLLER_CONTROLLER_SOLAR_PANEL_CURRENT_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_CONTROLLER_CONTROLLER_SOLAR_PANEL_CURRENT_NORMAL", 5),
					Max:    getConfigValue("KPIRANGES_COM_CONTROLLER_CONTROLLER_SOLAR_PANEL_CURRENT_MAX", 12),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_controller_controller_solar_panel_current",
							Help: "solar_panel_current (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_controller_controller_solar_panel_current",
				},
				{
					Key:    "com_controller_controller_solar_panel_power",
					Min:    getConfigValue("KPIRANGES_COM_CONTROLLER_CONTROLLER_SOLAR_PANEL_POWER_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_CONTROLLER_CONTROLLER_SOLAR_PANEL_POWER_NORMAL", 300),
					Max:    getConfigValue("KPIRANGES_COM_CONTROLLER_CONTROLLER_SOLAR_PANEL_POWER_MAX", 600),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_controller_controller_solar_panel_power",
							Help: "solar_panel_power (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_controller_controller_solar_panel_power",
				},
				{
					Key:    "com_controller_controller_battery_voltage",
					Min:    getConfigValue("KPIRANGES_COM_CONTROLLER_CONTROLLER_BATTERY_VOLTAGE_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_CONTROLLER_CONTROLLER_BATTERY_VOLTAGE_NORMAL", 48),
					Max:    getConfigValue("KPIRANGES_COM_CONTROLLER_CONTROLLER_BATTERY_VOLTAGE_MAX", 60),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_controller_controller_battery_voltage",
							Help: "battery_voltage (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_controller_controller_battery_voltage",
				},
				{
					Key:    "com_controller_controller_battery_current",
					Min:    getConfigValue("KPIRANGES_COM_CONTROLLER_CONTROLLER_BATTERY_CURRENT_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_CONTROLLER_CONTROLLER_BATTERY_CURRENT_NORMAL", 10),
					Max:    getConfigValue("KPIRANGES_COM_CONTROLLER_CONTROLLER_BATTERY_CURRENT_MAX", 20),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_controller_controller_battery_current",
							Help: "battery_current (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_controller_controller_battery_current",
				},
				{
					Key:    "com_controller_controller_battery_charge_percentage",
					Min:    getConfigValue("KPIRANGES_COM_CONTROLLER_CONTROLLER_BATTERY_CHARGE_PERCENTAGE_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_CONTROLLER_CONTROLLER_BATTERY_CHARGE_PERCENTAGE_NORMAL", 70),
					Max:    getConfigValue("KPIRANGES_COM_CONTROLLER_CONTROLLER_BATTERY_CHARGE_PERCENTAGE_MAX", 100),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_controller_controller_battery_charge_percentage",
							Help: "battery_charge (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_controller_controller_battery_charge_percentage",
				},
				{
					Key:    "com_switch_switch_poe_total_power_watts",
					Min:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_POE_TOTAL_POWER_WATTS_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_POE_TOTAL_POWER_WATTS_NORMAL", 60),
					Max:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_POE_TOTAL_POWER_WATTS_MAX", 120),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_switch_switch_poe_total_power_watts",
							Help: "poe_total_power (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_switch_switch_poe_total_power_watts",
				},
				{
					Key:    "com_switch_switch_poe_max_power_watts",
					Min:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_POE_MAX_POWER_WATTS_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_POE_MAX_POWER_WATTS_NORMAL", 90),
					Max:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_POE_MAX_POWER_WATTS_MAX", 120),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_switch_switch_poe_max_power_watts",
							Help: "poe_max_power (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_switch_switch_poe_max_power_watts",
				},
				{
					Key:    "com_switch_switch_system_temperature_c",
					Min:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_SYSTEM_TEMPERATURE_C_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_SYSTEM_TEMPERATURE_C_NORMAL", 65),
					Max:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_SYSTEM_TEMPERATURE_C_MAX", 80),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_switch_switch_system_temperature_c",
							Help: "switch_temperature (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_switch_switch_system_temperature_c",
				},
				{
					Key:    "com_switch_switch_system_power_watts",
					Min:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_SYSTEM_POWER_WATTS_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_SYSTEM_POWER_WATTS_NORMAL", 80),
					Max:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_SYSTEM_POWER_WATTS_MAX", 240),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_switch_switch_system_power_watts",
							Help: "switch_system_power (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_switch_switch_system_power_watts",
				},
				{
					Key:    "com_switch_switch_system_current_amps",
					Min:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_SYSTEM_CURRENT_AMPS_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_SYSTEM_CURRENT_AMPS_NORMAL", 2),
					Max:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_SYSTEM_CURRENT_AMPS_MAX", 10),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_switch_switch_system_current_amps",
							Help: "switch_system_current (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_switch_switch_system_current_amps",
				},
				{
					Key:    "com_switch_switch_input_voltage",
					Min:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_INPUT_VOLTAGE_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_INPUT_VOLTAGE_NORMAL", 54),
					Max:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_INPUT_VOLTAGE_MAX", 60),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_switch_switch_input_voltage",
							Help: "switch_input_voltage (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_switch_switch_input_voltage",
				},
				{
					Key:    "com_switch_switch_ambient_temperature_c",
					Min:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_AMBIENT_TEMPERATURE_C_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_AMBIENT_TEMPERATURE_C_NORMAL", 45),
					Max:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_AMBIENT_TEMPERATURE_C_MAX", 70),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_switch_switch_ambient_temperature_c",
							Help: "switch_ambient_temperature (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_switch_switch_ambient_temperature_c",
				},
				{
					Key:    "com_network_uplink_link",
					Min:    getConfigValue("KPIRANGES_COM_NETWORK_UPLINK_LINK_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_NETWORK_UPLINK_LINK_NORMAL", 1),
					Max:    getConfigValue("KPIRANGES_COM_NETWORK_UPLINK_LINK_MAX", 1),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_network_uplink_link",
							Help: "network_uplink_link (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_network_uplink_link",
				},
				{
					Key:    "com_network_uplink_linkspeed",
					Min:    getConfigValue("KPIRANGES_COM_NETWORK_UPLINK_LINKSPEED_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_NETWORK_UPLINK_LINKSPEED_NORMAL", 1000),
					Max:    getConfigValue("KPIRANGES_COM_NETWORK_UPLINK_LINKSPEED_MAX", 10000),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_network_uplink_linkspeed",
							Help: "network_uplink_linkspeed (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_network_uplink_linkspeed",
				},
				{
					Key:    "com_network_uplink_rx_bytes",
					Min:    getConfigValue("KPIRANGES_COM_NETWORK_UPLINK_RX_BYTES_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_NETWORK_UPLINK_RX_BYTES_NORMAL", 1000000),
					Max:    getConfigValue("KPIRANGES_COM_NETWORK_UPLINK_RX_BYTES_MAX", 1000000000),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_network_uplink_rx_bytes",
							Help: "network_uplink_rx_bytes (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_network_uplink_rx_bytes",
				},
				{
					Key:    "com_network_uplink_tx_bytes",
					Min:    getConfigValue("KPIRANGES_COM_NETWORK_UPLINK_TX_BYTES_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_NETWORK_UPLINK_TX_BYTES_NORMAL", 1000000),
					Max:    getConfigValue("KPIRANGES_COM_NETWORK_UPLINK_TX_BYTES_MAX", 1000000000),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_network_uplink_tx_bytes",
							Help: "network_uplink_tx_bytes (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_network_uplink_tx_bytes",
				},
				{
					Key:    "com_network_uplink_rx_packets",
					Min:    getConfigValue("KPIRANGES_COM_NETWORK_UPLINK_RX_PACKETS_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_NETWORK_UPLINK_RX_PACKETS_NORMAL", 10000),
					Max:    getConfigValue("KPIRANGES_COM_NETWORK_UPLINK_RX_PACKETS_MAX", 1000000),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_network_uplink_rx_packets",
							Help: "network_uplink_rx_packets (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_network_uplink_rx_packets",
				},
				{
					Key:    "com_network_uplink_tx_packets",
					Min:    getConfigValue("KPIRANGES_COM_NETWORK_UPLINK_TX_PACKETS_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_NETWORK_UPLINK_TX_PACKETS_NORMAL", 10000),
					Max:    getConfigValue("KPIRANGES_COM_NETWORK_UPLINK_TX_PACKETS_MAX", 1000000),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_network_uplink_tx_packets",
							Help: "network_uplink_tx_packets (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_network_uplink_tx_packets",
				},
				{
					Key:    "com_network_uplink_rx_error",
					Min:    getConfigValue("KPIRANGES_COM_NETWORK_UPLINK_RX_ERROR_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_NETWORK_UPLINK_RX_ERROR_NORMAL", 0),
					Max:    getConfigValue("KPIRANGES_COM_NETWORK_UPLINK_RX_ERROR_MAX", 100),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_network_uplink_rx_error",
							Help: "network_uplink_rx_error (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_network_uplink_rx_error",
				},
				{
					Key:    "com_network_uplink_tx_error",
					Min:    getConfigValue("KPIRANGES_COM_NETWORK_UPLINK_TX_ERROR_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_NETWORK_UPLINK_TX_ERROR_NORMAL", 0),
					Max:    getConfigValue("KPIRANGES_COM_NETWORK_UPLINK_TX_ERROR_MAX", 100),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_network_uplink_tx_error",
							Help: "network_uplink_tx_error (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_network_uplink_tx_error",
				},
				{
					Key:    "com_network_uplink_rx_dropped",
					Min:    getConfigValue("KPIRANGES_COM_NETWORK_UPLINK_RX_DROPPED_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_NETWORK_UPLINK_RX_DROPPED_NORMAL", 0),
					Max:    getConfigValue("KPIRANGES_COM_NETWORK_UPLINK_RX_DROPPED_MAX", 100),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_network_uplink_rx_dropped",
							Help: "network_uplink_rx_dropped (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_network_uplink_rx_dropped",
				},
				{
					Key:    "com_network_uplink_tx_dropped",
					Min:    getConfigValue("KPIRANGES_COM_NETWORK_UPLINK_TX_DROPPED_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_NETWORK_UPLINK_TX_DROPPED_NORMAL", 0),
					Max:    getConfigValue("KPIRANGES_COM_NETWORK_UPLINK_TX_DROPPED_MAX", 100),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_network_uplink_tx_dropped",
							Help: "network_uplink_tx_dropped (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_network_uplink_tx_dropped",
				},
				{
					Key:    "com_backhaul_backhaul_near_ttfb_median_ms",
					Min:    getConfigValue("KPIRANGES_COM_BACKHAUL_BACKHAUL_NEAR_TTFB_MEDIAN_MS_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_BACKHAUL_BACKHAUL_NEAR_TTFB_MEDIAN_MS_NORMAL", 100),
					Max:    getConfigValue("KPIRANGES_COM_BACKHAUL_BACKHAUL_NEAR_TTFB_MEDIAN_MS_MAX", 500),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_backhaul_backhaul_near_ttfb_median_ms",
							Help: "backhaul_near_ttfb_median (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_backhaul_backhaul_near_ttfb_median_ms",
				},
				{
					Key:    "com_backhaul_backhaul_near_ttfb_p95_ms",
					Min:    getConfigValue("KPIRANGES_COM_BACKHAUL_BACKHAUL_NEAR_TTFB_P95_MS_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_BACKHAUL_BACKHAUL_NEAR_TTFB_P95_MS_NORMAL", 200),
					Max:    getConfigValue("KPIRANGES_COM_BACKHAUL_BACKHAUL_NEAR_TTFB_P95_MS_MAX", 1000),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_backhaul_backhaul_near_ttfb_p95_ms",
							Help: "backhaul_near_ttfb_p95 (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_backhaul_backhaul_near_ttfb_p95_ms",
				},
				{
					Key:    "com_backhaul_backhaul_near_ttfb_p99_ms",
					Min:    getConfigValue("KPIRANGES_COM_BACKHAUL_BACKHAUL_NEAR_TTFB_P99_MS_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_BACKHAUL_BACKHAUL_NEAR_TTFB_P99_MS_NORMAL", 300),
					Max:    getConfigValue("KPIRANGES_COM_BACKHAUL_BACKHAUL_NEAR_TTFB_P99_MS_MAX", 1500),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_backhaul_backhaul_near_ttfb_p99_ms",
							Help: "backhaul_near_ttfb_p99 (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_backhaul_backhaul_near_ttfb_p99_ms",
				},
				{
					Key:    "com_backhaul_backhaul_far_ttfb_median_ms",
					Min:    getConfigValue("KPIRANGES_COM_BACKHAUL_BACKHAUL_FAR_TTFB_MEDIAN_MS_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_BACKHAUL_BACKHAUL_FAR_TTFB_MEDIAN_MS_NORMAL", 100),
					Max:    getConfigValue("KPIRANGES_COM_BACKHAUL_BACKHAUL_FAR_TTFB_MEDIAN_MS_MAX", 500),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_backhaul_backhaul_far_ttfb_median_ms",
							Help: "backhaul_far_ttfb_median (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_backhaul_backhaul_far_ttfb_median_ms",
				},
				{
					Key:    "com_backhaul_backhaul_far_ttfb_p95_ms",
					Min:    getConfigValue("KPIRANGES_COM_BACKHAUL_BACKHAUL_FAR_TTFB_P95_MS_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_BACKHAUL_BACKHAUL_FAR_TTFB_P95_MS_NORMAL", 200),
					Max:    getConfigValue("KPIRANGES_COM_BACKHAUL_BACKHAUL_FAR_TTFB_P95_MS_MAX", 1000),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_backhaul_backhaul_far_ttfb_p95_ms",
							Help: "backhaul_far_ttfb_p95 (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_backhaul_backhaul_far_ttfb_p95_ms",
				},
				{
					Key:    "com_backhaul_backhaul_far_ttfb_p99_ms",
					Min:    getConfigValue("KPIRANGES_COM_BACKHAUL_BACKHAUL_FAR_TTFB_P99_MS_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_BACKHAUL_BACKHAUL_FAR_TTFB_P99_MS_NORMAL", 300),
					Max:    getConfigValue("KPIRANGES_COM_BACKHAUL_BACKHAUL_FAR_TTFB_P99_MS_MAX", 1500),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_backhaul_backhaul_far_ttfb_p99_ms",
							Help: "backhaul_far_ttfb_p99 (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_backhaul_backhaul_far_ttfb_p99_ms",
				},
				{
					Key:    "com_backhaul_backhaul_probe_success_rate_pct",
					Min:    getConfigValue("KPIRANGES_COM_BACKHAUL_BACKHAUL_PROBE_SUCCESS_RATE_PCT_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_BACKHAUL_BACKHAUL_PROBE_SUCCESS_RATE_PCT_NORMAL", 90),
					Max:    getConfigValue("KPIRANGES_COM_BACKHAUL_BACKHAUL_PROBE_SUCCESS_RATE_PCT_MAX", 100),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_backhaul_backhaul_probe_success_rate_pct",
							Help: "backhaul_probe_success_rate (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_backhaul_backhaul_probe_success_rate_pct",
				},
				{
					Key:    "com_backhaul_backhaul_stall_rate_pct",
					Min:    getConfigValue("KPIRANGES_COM_BACKHAUL_BACKHAUL_STALL_RATE_PCT_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_BACKHAUL_BACKHAUL_STALL_RATE_PCT_NORMAL", 10),
					Max:    getConfigValue("KPIRANGES_COM_BACKHAUL_BACKHAUL_STALL_RATE_PCT_MAX", 100),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_backhaul_backhaul_stall_rate_pct",
							Help: "backhaul_stall_rate (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_backhaul_backhaul_stall_rate_pct",
				},
				{
					Key:    "com_backhaul_backhaul_cap_detected_mbps",
					Min:    getConfigValue("KPIRANGES_COM_BACKHAUL_BACKHAUL_CAP_DETECTED_MBPS_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_BACKHAUL_BACKHAUL_CAP_DETECTED_MBPS_NORMAL", 200),
					Max:    getConfigValue("KPIRANGES_COM_BACKHAUL_BACKHAUL_CAP_DETECTED_MBPS_MAX", 1000),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_backhaul_backhaul_cap_detected_mbps",
							Help: "backhaul_cap_detected (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_backhaul_backhaul_cap_detected_mbps",
				},
				{
					Key:    "com_backhaul_backhaul_bufferbloat_inflation_factor",
					Min:    getConfigValue("KPIRANGES_COM_BACKHAUL_BACKHAUL_BUFFERBLOAT_INFLATION_FACTOR_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_BACKHAUL_BACKHAUL_BUFFERBLOAT_INFLATION_FACTOR_NORMAL", 1),
					Max:    getConfigValue("KPIRANGES_COM_BACKHAUL_BACKHAUL_BUFFERBLOAT_INFLATION_FACTOR_MAX", 5),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_backhaul_backhaul_bufferbloat_inflation_factor",
							Help: "backhaul_bufferbloat_inflation (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_backhaul_backhaul_bufferbloat_inflation_factor",
				},
				{
					Key:    "com_backhaul_backhaul_diag_present",
					Min:    getConfigValue("KPIRANGES_COM_BACKHAUL_BACKHAUL_DIAG_PRESENT_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_BACKHAUL_BACKHAUL_DIAG_PRESENT_NORMAL", 1),
					Max:    getConfigValue("KPIRANGES_COM_BACKHAUL_BACKHAUL_DIAG_PRESENT_MAX", 1),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_backhaul_backhaul_diag_present",
							Help: "backhaul_diag_present (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_backhaul_backhaul_diag_present",
				},
				{
					Key:    "com_controller_controller_solar_yield_today",
					Min:    getConfigValue("KPIRANGES_COM_CONTROLLER_CONTROLLER_SOLAR_YIELD_TODAY_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_CONTROLLER_CONTROLLER_SOLAR_YIELD_TODAY_NORMAL", 20),
					Max:    getConfigValue("KPIRANGES_COM_CONTROLLER_CONTROLLER_SOLAR_YIELD_TODAY_MAX", 100),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_controller_controller_solar_yield_today",
							Help: "solar_yield_today (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_controller_controller_solar_yield_today",
				},
				{
					Key:    "com_controller_controller_solar_yield_total",
					Min:    getConfigValue("KPIRANGES_COM_CONTROLLER_CONTROLLER_SOLAR_YIELD_TOTAL_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_CONTROLLER_CONTROLLER_SOLAR_YIELD_TOTAL_NORMAL", 500),
					Max:    getConfigValue("KPIRANGES_COM_CONTROLLER_CONTROLLER_SOLAR_YIELD_TOTAL_MAX", 50000),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_controller_controller_solar_yield_total",
							Help: "solar_yield_total (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_controller_controller_solar_yield_total",
				},
				{
					Key:    "com_controller_controller_controller_temperature",
					Min:    getConfigValue("KPIRANGES_COM_CONTROLLER_CONTROLLER_CONTROLLER_TEMPERATURE_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_CONTROLLER_CONTROLLER_CONTROLLER_TEMPERATURE_NORMAL", 60),
					Max:    getConfigValue("KPIRANGES_COM_CONTROLLER_CONTROLLER_CONTROLLER_TEMPERATURE_MAX", 80),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_controller_controller_controller_temperature",
							Help: "controller_temperature (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_controller_controller_controller_temperature",
				},
				{
					Key:    "com_controller_controller_mppt_efficiency",
					Min:    getConfigValue("KPIRANGES_COM_CONTROLLER_CONTROLLER_MPPT_EFFICIENCY_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_CONTROLLER_CONTROLLER_MPPT_EFFICIENCY_NORMAL", 90),
					Max:    getConfigValue("KPIRANGES_COM_CONTROLLER_CONTROLLER_MPPT_EFFICIENCY_MAX", 100),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_controller_controller_mppt_efficiency",
							Help: "mppt_efficiency (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_controller_controller_mppt_efficiency",
				},
				{
					Key:    "com_controller_controller_load_current",
					Min:    getConfigValue("KPIRANGES_COM_CONTROLLER_CONTROLLER_LOAD_CURRENT_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_CONTROLLER_CONTROLLER_LOAD_CURRENT_NORMAL", 10),
					Max:    getConfigValue("KPIRANGES_COM_CONTROLLER_CONTROLLER_LOAD_CURRENT_MAX", 30),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_controller_controller_load_current",
							Help: "load_current (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_controller_controller_load_current",
				},
				{
					Key:    "com_memory_ddr_total",
					Min:    getConfigValue("KPIRANGES_COM_MEMORY_DDR_TOTAL_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_MEMORY_DDR_TOTAL_NORMAL", 16000),
					Max:    getConfigValue("KPIRANGES_COM_MEMORY_DDR_TOTAL_MAX", 64000),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_memory_ddr_total",
							Help: "memory_total (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_memory_ddr_total",
				},
				{
					Key:    "com_memory_ddr_free",
					Min:    getConfigValue("KPIRANGES_COM_MEMORY_DDR_FREE_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_MEMORY_DDR_FREE_NORMAL", 2000),
					Max:    getConfigValue("KPIRANGES_COM_MEMORY_DDR_FREE_MAX", 32000),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_memory_ddr_free",
							Help: "memory_free (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_memory_ddr_free",
				},
				{
					Key:    "com_storage_emmc_total",
					Min:    getConfigValue("KPIRANGES_COM_STORAGE_EMMC_TOTAL_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_STORAGE_EMMC_TOTAL_NORMAL", 64000),
					Max:    getConfigValue("KPIRANGES_COM_STORAGE_EMMC_TOTAL_MAX", 256000),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_storage_emmc_total",
							Help: "storage_total (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_storage_emmc_total",
				},
				{
					Key:    "com_storage_emmc_free",
					Min:    getConfigValue("KPIRANGES_COM_STORAGE_EMMC_FREE_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_STORAGE_EMMC_FREE_NORMAL", 10000),
					Max:    getConfigValue("KPIRANGES_COM_STORAGE_EMMC_FREE_MAX", 200000),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_storage_emmc_free",
							Help: "storage_free (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_storage_emmc_free",
				},
				{
					Key:    "com_switch_switch_port_1_tnode_poe_present",
					Min:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_1_TNODE_POE_PRESENT_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_1_TNODE_POE_PRESENT_NORMAL", 1),
					Max:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_1_TNODE_POE_PRESENT_MAX", 1),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_switch_switch_port_1_tnode_poe_present",
							Help: "switch_port_1_present (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_switch_switch_port_1_tnode_poe_present",
				},
				{
					Key:    "com_switch_switch_port_1_tnode_poe_admin_up",
					Min:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_1_TNODE_POE_ADMIN_UP_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_1_TNODE_POE_ADMIN_UP_NORMAL", 1),
					Max:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_1_TNODE_POE_ADMIN_UP_MAX", 1),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_switch_switch_port_1_tnode_poe_admin_up",
							Help: "switch_port_1_admin_up (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_switch_switch_port_1_tnode_poe_admin_up",
				},
				{
					Key:    "com_switch_switch_port_1_tnode_poe_link_up",
					Min:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_1_TNODE_POE_LINK_UP_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_1_TNODE_POE_LINK_UP_NORMAL", 1),
					Max:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_1_TNODE_POE_LINK_UP_MAX", 1),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_switch_switch_port_1_tnode_poe_link_up",
							Help: "switch_port_1_link_up (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_switch_switch_port_1_tnode_poe_link_up",
				},
				{
					Key:    "com_switch_switch_port_1_tnode_poe_speed_bps",
					Min:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_1_TNODE_POE_SPEED_BPS_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_1_TNODE_POE_SPEED_BPS_NORMAL", 1000000000),
					Max:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_1_TNODE_POE_SPEED_BPS_MAX", 10000000000),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_switch_switch_port_1_tnode_poe_speed_bps",
							Help: "switch_port_1_speed (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_switch_switch_port_1_tnode_poe_speed_bps",
				},
				{
					Key:    "com_switch_switch_port_1_tnode_poe_poe_power_watts",
					Min:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_1_TNODE_POE_POE_POWER_WATTS_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_1_TNODE_POE_POE_POWER_WATTS_NORMAL", 15),
					Max:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_1_TNODE_POE_POE_POWER_WATTS_MAX", 60),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_switch_switch_port_1_tnode_poe_poe_power_watts",
							Help: "switch_port_1_power (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_switch_switch_port_1_tnode_poe_poe_power_watts",
				},
				{
					Key:    "com_switch_switch_port_2_cnode_poe_present",
					Min:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_2_CNODE_POE_PRESENT_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_2_CNODE_POE_PRESENT_NORMAL", 1),
					Max:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_2_CNODE_POE_PRESENT_MAX", 1),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_switch_switch_port_2_cnode_poe_present",
							Help: "switch_port_2_present (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_switch_switch_port_2_cnode_poe_present",
				},
				{
					Key:    "com_switch_switch_port_2_cnode_poe_admin_up",
					Min:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_2_CNODE_POE_ADMIN_UP_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_2_CNODE_POE_ADMIN_UP_NORMAL", 1),
					Max:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_2_CNODE_POE_ADMIN_UP_MAX", 1),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_switch_switch_port_2_cnode_poe_admin_up",
							Help: "switch_port_2_admin_up (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_switch_switch_port_2_cnode_poe_admin_up",
				},
				{
					Key:    "com_switch_switch_port_2_cnode_poe_link_up",
					Min:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_2_CNODE_POE_LINK_UP_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_2_CNODE_POE_LINK_UP_NORMAL", 1),
					Max:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_2_CNODE_POE_LINK_UP_MAX", 1),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_switch_switch_port_2_cnode_poe_link_up",
							Help: "switch_port_2_link_up (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_switch_switch_port_2_cnode_poe_link_up",
				},
				{
					Key:    "com_switch_switch_port_2_cnode_poe_speed_bps",
					Min:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_2_CNODE_POE_SPEED_BPS_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_2_CNODE_POE_SPEED_BPS_NORMAL", 1000000000),
					Max:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_2_CNODE_POE_SPEED_BPS_MAX", 10000000000),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_switch_switch_port_2_cnode_poe_speed_bps",
							Help: "switch_port_2_speed (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_switch_switch_port_2_cnode_poe_speed_bps",
				},
				{
					Key:    "com_switch_switch_port_2_cnode_poe_poe_power_watts",
					Min:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_2_CNODE_POE_POE_POWER_WATTS_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_2_CNODE_POE_POE_POWER_WATTS_NORMAL", 15),
					Max:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_2_CNODE_POE_POE_POWER_WATTS_MAX", 60),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_switch_switch_port_2_cnode_poe_poe_power_watts",
							Help: "switch_port_2_power (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_switch_switch_port_2_cnode_poe_poe_power_watts",
				},
				{
					Key:    "com_switch_switch_port_3_anode_poe_present",
					Min:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_3_ANODE_POE_PRESENT_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_3_ANODE_POE_PRESENT_NORMAL", 1),
					Max:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_3_ANODE_POE_PRESENT_MAX", 1),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_switch_switch_port_3_anode_poe_present",
							Help: "switch_port_3_present (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_switch_switch_port_3_anode_poe_present",
				},
				{
					Key:    "com_switch_switch_port_3_anode_poe_admin_up",
					Min:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_3_ANODE_POE_ADMIN_UP_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_3_ANODE_POE_ADMIN_UP_NORMAL", 1),
					Max:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_3_ANODE_POE_ADMIN_UP_MAX", 1),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_switch_switch_port_3_anode_poe_admin_up",
							Help: "switch_port_3_admin_up (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_switch_switch_port_3_anode_poe_admin_up",
				},
				{
					Key:    "com_switch_switch_port_3_anode_poe_link_up",
					Min:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_3_ANODE_POE_LINK_UP_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_3_ANODE_POE_LINK_UP_NORMAL", 1),
					Max:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_3_ANODE_POE_LINK_UP_MAX", 1),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_switch_switch_port_3_anode_poe_link_up",
							Help: "switch_port_3_link_up (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_switch_switch_port_3_anode_poe_link_up",
				},
				{
					Key:    "com_switch_switch_port_3_anode_poe_speed_bps",
					Min:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_3_ANODE_POE_SPEED_BPS_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_3_ANODE_POE_SPEED_BPS_NORMAL", 1000000000),
					Max:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_3_ANODE_POE_SPEED_BPS_MAX", 10000000000),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_switch_switch_port_3_anode_poe_speed_bps",
							Help: "switch_port_3_speed (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_switch_switch_port_3_anode_poe_speed_bps",
				},
				{
					Key:    "com_switch_switch_port_3_anode_poe_poe_power_watts",
					Min:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_3_ANODE_POE_POE_POWER_WATTS_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_3_ANODE_POE_POE_POWER_WATTS_NORMAL", 15),
					Max:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_3_ANODE_POE_POE_POWER_WATTS_MAX", 60),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_switch_switch_port_3_anode_poe_poe_power_watts",
							Help: "switch_port_3_power (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_switch_switch_port_3_anode_poe_poe_power_watts",
				},
				{
					Key:    "com_switch_switch_port_4_spare_poe_present",
					Min:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_4_SPARE_POE_PRESENT_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_4_SPARE_POE_PRESENT_NORMAL", 1),
					Max:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_4_SPARE_POE_PRESENT_MAX", 1),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_switch_switch_port_4_spare_poe_present",
							Help: "switch_port_4_present (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_switch_switch_port_4_spare_poe_present",
				},
				{
					Key:    "com_switch_switch_port_4_spare_poe_admin_up",
					Min:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_4_SPARE_POE_ADMIN_UP_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_4_SPARE_POE_ADMIN_UP_NORMAL", 1),
					Max:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_4_SPARE_POE_ADMIN_UP_MAX", 1),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_switch_switch_port_4_spare_poe_admin_up",
							Help: "switch_port_4_admin_up (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_switch_switch_port_4_spare_poe_admin_up",
				},
				{
					Key:    "com_switch_switch_port_4_spare_poe_link_up",
					Min:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_4_SPARE_POE_LINK_UP_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_4_SPARE_POE_LINK_UP_NORMAL", 1),
					Max:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_4_SPARE_POE_LINK_UP_MAX", 1),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_switch_switch_port_4_spare_poe_link_up",
							Help: "switch_port_4_link_up (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_switch_switch_port_4_spare_poe_link_up",
				},
				{
					Key:    "com_switch_switch_port_4_spare_poe_speed_bps",
					Min:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_4_SPARE_POE_SPEED_BPS_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_4_SPARE_POE_SPEED_BPS_NORMAL", 1000000000),
					Max:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_4_SPARE_POE_SPEED_BPS_MAX", 10000000000),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_switch_switch_port_4_spare_poe_speed_bps",
							Help: "switch_port_4_speed (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_switch_switch_port_4_spare_poe_speed_bps",
				},
				{
					Key:    "com_switch_switch_port_4_spare_poe_poe_power_watts",
					Min:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_4_SPARE_POE_POE_POWER_WATTS_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_4_SPARE_POE_POE_POWER_WATTS_NORMAL", 15),
					Max:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_4_SPARE_POE_POE_POWER_WATTS_MAX", 60),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_switch_switch_port_4_spare_poe_poe_power_watts",
							Help: "switch_port_4_power (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_switch_switch_port_4_spare_poe_poe_power_watts",
				},
				{
					Key:    "com_switch_switch_port_5_port5_present",
					Min:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_5_PORT5_PRESENT_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_5_PORT5_PRESENT_NORMAL", 1),
					Max:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_5_PORT5_PRESENT_MAX", 1),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_switch_switch_port_5_port5_present",
							Help: "switch_port_5_present (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_switch_switch_port_5_port5_present",
				},
				{
					Key:    "com_switch_switch_port_5_port5_admin_up",
					Min:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_5_PORT5_ADMIN_UP_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_5_PORT5_ADMIN_UP_NORMAL", 1),
					Max:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_5_PORT5_ADMIN_UP_MAX", 1),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_switch_switch_port_5_port5_admin_up",
							Help: "switch_port_5_admin_up (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_switch_switch_port_5_port5_admin_up",
				},
				{
					Key:    "com_switch_switch_port_5_port5_link_up",
					Min:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_5_PORT5_LINK_UP_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_5_PORT5_LINK_UP_NORMAL", 1),
					Max:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_5_PORT5_LINK_UP_MAX", 1),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_switch_switch_port_5_port5_link_up",
							Help: "switch_port_5_link_up (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_switch_switch_port_5_port5_link_up",
				},
				{
					Key:    "com_switch_switch_port_5_port5_speed_bps",
					Min:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_5_PORT5_SPEED_BPS_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_5_PORT5_SPEED_BPS_NORMAL", 1000000000),
					Max:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_5_PORT5_SPEED_BPS_MAX", 10000000000),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_switch_switch_port_5_port5_speed_bps",
							Help: "switch_port_5_speed (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_switch_switch_port_5_port5_speed_bps",
				},
				{
					Key:    "com_switch_switch_port_5_port5_poe_power_watts",
					Min:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_5_PORT5_POE_POWER_WATTS_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_5_PORT5_POE_POWER_WATTS_NORMAL", 15),
					Max:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_5_PORT5_POE_POWER_WATTS_MAX", 60),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_switch_switch_port_5_port5_poe_power_watts",
							Help: "switch_port_5_power (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_switch_switch_port_5_port5_poe_power_watts",
				},
				{
					Key:    "com_switch_switch_port_6_port6_present",
					Min:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_6_PORT6_PRESENT_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_6_PORT6_PRESENT_NORMAL", 1),
					Max:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_6_PORT6_PRESENT_MAX", 1),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_switch_switch_port_6_port6_present",
							Help: "switch_port_6_present (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_switch_switch_port_6_port6_present",
				},
				{
					Key:    "com_switch_switch_port_6_port6_admin_up",
					Min:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_6_PORT6_ADMIN_UP_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_6_PORT6_ADMIN_UP_NORMAL", 1),
					Max:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_6_PORT6_ADMIN_UP_MAX", 1),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_switch_switch_port_6_port6_admin_up",
							Help: "switch_port_6_admin_up (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_switch_switch_port_6_port6_admin_up",
				},
				{
					Key:    "com_switch_switch_port_6_port6_link_up",
					Min:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_6_PORT6_LINK_UP_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_6_PORT6_LINK_UP_NORMAL", 1),
					Max:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_6_PORT6_LINK_UP_MAX", 1),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_switch_switch_port_6_port6_link_up",
							Help: "switch_port_6_link_up (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_switch_switch_port_6_port6_link_up",
				},
				{
					Key:    "com_switch_switch_port_6_port6_speed_bps",
					Min:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_6_PORT6_SPEED_BPS_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_6_PORT6_SPEED_BPS_NORMAL", 1000000000),
					Max:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_6_PORT6_SPEED_BPS_MAX", 10000000000),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_switch_switch_port_6_port6_speed_bps",
							Help: "switch_port_6_speed (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_switch_switch_port_6_port6_speed_bps",
				},
				{
					Key:    "com_switch_switch_port_6_port6_poe_power_watts",
					Min:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_6_PORT6_POE_POWER_WATTS_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_6_PORT6_POE_POWER_WATTS_NORMAL", 15),
					Max:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_6_PORT6_POE_POWER_WATTS_MAX", 60),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_switch_switch_port_6_port6_poe_power_watts",
							Help: "switch_port_6_power (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_switch_switch_port_6_port6_poe_power_watts",
				},
				{
					Key:    "com_switch_switch_port_7_port7_present",
					Min:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_7_PORT7_PRESENT_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_7_PORT7_PRESENT_NORMAL", 1),
					Max:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_7_PORT7_PRESENT_MAX", 1),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_switch_switch_port_7_port7_present",
							Help: "switch_port_7_present (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_switch_switch_port_7_port7_present",
				},
				{
					Key:    "com_switch_switch_port_7_port7_admin_up",
					Min:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_7_PORT7_ADMIN_UP_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_7_PORT7_ADMIN_UP_NORMAL", 1),
					Max:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_7_PORT7_ADMIN_UP_MAX", 1),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_switch_switch_port_7_port7_admin_up",
							Help: "switch_port_7_admin_up (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_switch_switch_port_7_port7_admin_up",
				},
				{
					Key:    "com_switch_switch_port_7_port7_link_up",
					Min:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_7_PORT7_LINK_UP_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_7_PORT7_LINK_UP_NORMAL", 1),
					Max:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_7_PORT7_LINK_UP_MAX", 1),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_switch_switch_port_7_port7_link_up",
							Help: "switch_port_7_link_up (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_switch_switch_port_7_port7_link_up",
				},
				{
					Key:    "com_switch_switch_port_7_port7_speed_bps",
					Min:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_7_PORT7_SPEED_BPS_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_7_PORT7_SPEED_BPS_NORMAL", 1000000000),
					Max:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_7_PORT7_SPEED_BPS_MAX", 10000000000),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_switch_switch_port_7_port7_speed_bps",
							Help: "switch_port_7_speed (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_switch_switch_port_7_port7_speed_bps",
				},
				{
					Key:    "com_switch_switch_port_7_port7_poe_power_watts",
					Min:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_7_PORT7_POE_POWER_WATTS_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_7_PORT7_POE_POWER_WATTS_NORMAL", 15),
					Max:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_7_PORT7_POE_POWER_WATTS_MAX", 60),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_switch_switch_port_7_port7_poe_power_watts",
							Help: "switch_port_7_power (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_switch_switch_port_7_port7_poe_power_watts",
				},
				{
					Key:    "com_switch_switch_port_8_port8_present",
					Min:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_8_PORT8_PRESENT_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_8_PORT8_PRESENT_NORMAL", 1),
					Max:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_8_PORT8_PRESENT_MAX", 1),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_switch_switch_port_8_port8_present",
							Help: "switch_port_8_present (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_switch_switch_port_8_port8_present",
				},
				{
					Key:    "com_switch_switch_port_8_port8_admin_up",
					Min:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_8_PORT8_ADMIN_UP_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_8_PORT8_ADMIN_UP_NORMAL", 1),
					Max:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_8_PORT8_ADMIN_UP_MAX", 1),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_switch_switch_port_8_port8_admin_up",
							Help: "switch_port_8_admin_up (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_switch_switch_port_8_port8_admin_up",
				},
				{
					Key:    "com_switch_switch_port_8_port8_link_up",
					Min:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_8_PORT8_LINK_UP_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_8_PORT8_LINK_UP_NORMAL", 1),
					Max:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_8_PORT8_LINK_UP_MAX", 1),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_switch_switch_port_8_port8_link_up",
							Help: "switch_port_8_link_up (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_switch_switch_port_8_port8_link_up",
				},
				{
					Key:    "com_switch_switch_port_8_port8_speed_bps",
					Min:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_8_PORT8_SPEED_BPS_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_8_PORT8_SPEED_BPS_NORMAL", 1000000000),
					Max:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_8_PORT8_SPEED_BPS_MAX", 10000000000),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_switch_switch_port_8_port8_speed_bps",
							Help: "switch_port_8_speed (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_switch_switch_port_8_port8_speed_bps",
				},
				{
					Key:    "com_switch_switch_port_8_port8_poe_power_watts",
					Min:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_8_PORT8_POE_POWER_WATTS_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_8_PORT8_POE_POWER_WATTS_NORMAL", 15),
					Max:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_8_PORT8_POE_POWER_WATTS_MAX", 60),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_switch_switch_port_8_port8_poe_power_watts",
							Help: "switch_port_8_power (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_switch_switch_port_8_port8_poe_power_watts",
				},
				{
					Key:    "com_switch_switch_port_9_uplink_sfp_present",
					Min:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_9_UPLINK_SFP_PRESENT_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_9_UPLINK_SFP_PRESENT_NORMAL", 1),
					Max:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_9_UPLINK_SFP_PRESENT_MAX", 1),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_switch_switch_port_9_uplink_sfp_present",
							Help: "switch_port_9_present (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_switch_switch_port_9_uplink_sfp_present",
				},
				{
					Key:    "com_switch_switch_port_9_uplink_sfp_admin_up",
					Min:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_9_UPLINK_SFP_ADMIN_UP_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_9_UPLINK_SFP_ADMIN_UP_NORMAL", 1),
					Max:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_9_UPLINK_SFP_ADMIN_UP_MAX", 1),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_switch_switch_port_9_uplink_sfp_admin_up",
							Help: "switch_port_9_admin_up (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_switch_switch_port_9_uplink_sfp_admin_up",
				},
				{
					Key:    "com_switch_switch_port_9_uplink_sfp_link_up",
					Min:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_9_UPLINK_SFP_LINK_UP_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_9_UPLINK_SFP_LINK_UP_NORMAL", 1),
					Max:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_9_UPLINK_SFP_LINK_UP_MAX", 1),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_switch_switch_port_9_uplink_sfp_link_up",
							Help: "switch_port_9_link_up (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_switch_switch_port_9_uplink_sfp_link_up",
				},
				{
					Key:    "com_switch_switch_port_9_uplink_sfp_speed_bps",
					Min:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_9_UPLINK_SFP_SPEED_BPS_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_9_UPLINK_SFP_SPEED_BPS_NORMAL", 1000000000),
					Max:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_9_UPLINK_SFP_SPEED_BPS_MAX", 10000000000),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_switch_switch_port_9_uplink_sfp_speed_bps",
							Help: "switch_port_9_speed (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_switch_switch_port_9_uplink_sfp_speed_bps",
				},
				{
					Key:    "com_switch_switch_port_9_uplink_sfp_poe_power_watts",
					Min:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_9_UPLINK_SFP_POE_POWER_WATTS_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_9_UPLINK_SFP_POE_POWER_WATTS_NORMAL", 15),
					Max:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_9_UPLINK_SFP_POE_POWER_WATTS_MAX", 60),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_switch_switch_port_9_uplink_sfp_poe_power_watts",
							Help: "switch_port_9_power (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_switch_switch_port_9_uplink_sfp_poe_power_watts",
				},
				{
					Key:    "com_switch_switch_port_10_port10_present",
					Min:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_10_PORT10_PRESENT_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_10_PORT10_PRESENT_NORMAL", 1),
					Max:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_10_PORT10_PRESENT_MAX", 1),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_switch_switch_port_10_port10_present",
							Help: "switch_port_10_present (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_switch_switch_port_10_port10_present",
				},
				{
					Key:    "com_switch_switch_port_10_port10_admin_up",
					Min:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_10_PORT10_ADMIN_UP_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_10_PORT10_ADMIN_UP_NORMAL", 1),
					Max:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_10_PORT10_ADMIN_UP_MAX", 1),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_switch_switch_port_10_port10_admin_up",
							Help: "switch_port_10_admin_up (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_switch_switch_port_10_port10_admin_up",
				},
				{
					Key:    "com_switch_switch_port_10_port10_link_up",
					Min:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_10_PORT10_LINK_UP_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_10_PORT10_LINK_UP_NORMAL", 1),
					Max:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_10_PORT10_LINK_UP_MAX", 1),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_switch_switch_port_10_port10_link_up",
							Help: "switch_port_10_link_up (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_switch_switch_port_10_port10_link_up",
				},
				{
					Key:    "com_switch_switch_port_10_port10_speed_bps",
					Min:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_10_PORT10_SPEED_BPS_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_10_PORT10_SPEED_BPS_NORMAL", 1000000000),
					Max:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_10_PORT10_SPEED_BPS_MAX", 10000000000),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_switch_switch_port_10_port10_speed_bps",
							Help: "switch_port_10_speed (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_switch_switch_port_10_port10_speed_bps",
				},
				{
					Key:    "com_switch_switch_port_10_port10_poe_power_watts",
					Min:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_10_PORT10_POE_POWER_WATTS_MIN", 0),
					Normal: getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_10_PORT10_POE_POWER_WATTS_NORMAL", 15),
					Max:    getConfigValue("KPIRANGES_COM_SWITCH_SWITCH_PORT_10_PORT10_POE_POWER_WATTS_MAX", 60),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "com_switch_switch_port_10_port10_poe_power_watts",
							Help: "switch_port_10_power (cnode)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "com_switch_switch_port_10_port10_poe_power_watts",
				},
				{
					Key:    "active_subscribers_per_node",
					Min:    getConfigValue("KPIRANGES_ACTIVE_SUBSCRIBERS_PER_NODE_MIN", 0),
					Normal: getConfigValue("KPIRANGES_ACTIVE_SUBSCRIBERS_PER_NODE_NORMAL", 100),
					Max:    getConfigValue("KPIRANGES_ACTIVE_SUBSCRIBERS_PER_NODE_MAX", 1000),
					KPI: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "active_subscribers_per_node",
							Help: "node_active_subscribers (system)",
						},
						[]string{"node_id", "metric"},
					),
					Metric: "active_subscribers_per_node",
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
