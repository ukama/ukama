/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package pkg

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/rest"
)

type NameUpdate struct {
	Required bool   `json:"required" default:"false"`
	Slice    string `json:"slice" default:""`
}

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	Server            rest.HttpConfig
	Services          GrpcEndpoints  `mapstructure:"services"`
	HttpServices      HttpEndpoints  `mapstructure:"httpServices"`
	MetricsServer     config.Metrics `mapstructure:"metrics"`
	MetricsStore      string         `default:"http://localhost:8080"`
	Auth              *config.Auth   `mapstructure:"auth"`
	MetricsConfig     *MetricsConfig
	OrgName           string
	Period            time.Duration `default:"5s"`
	Http              HttpServices
}

type HttpServices struct {
	InitClient string `default:"api-gateway-init:8080"`
}

type Metric struct {
	NeedRate bool   `json:"needRate"`
	Metric   string `json:"metric"`
	// Range vector duration used in Rate func https://prometheus.io/docs/prometheus/latest/querying/basics/#time-durations
	// if NeedRate is false then this field is ignored
	// Example: 1d or 5h, or 30s
	RateInterval string `json:"rateInterval"`

	// consider adding aggregation function as a parameter
}

type MetricsConfig struct {
	Metrics             map[string]Metric  `json:"metrics"`
	MetricsServer       string             `default:"http://localhost:9090"`
	Timeout             time.Duration
	DefaultRateInterval string
}

var defaultPrometheusMetric = map[string]Metric{
	"cpu":                  {false, "trx_soc_cpu_usage", ""},
	"memory":               {false, "trx_memory_ddr_used", ""},
	"subscribers_active":   {false, "trx_lte_core_active_ue", ""},
	"subscribers_attached": {false, "trx_lte_core_subscribers", ""},

	// New Metrics
	"active_org_users":        {false, "number_of_active_users", ""},
	"inactive_org_users":      {false, "number_of_inactive_users", ""},
	"active_orgs":             {false, "number_of_active_org", ""},
	"inactive_orgs":           {false, "number_of_inactive_org", ""},
	"platform_active_users":   {false, "platform_active_users", ""},
	"platform_inactive_users": {false, "platform_inactive_users", ""},
	"networks":                {false, "number_of_networks", ""},
	"sites":                   {false, "number_of_sites", ""},
	"online_node_count":       {false, "online_node_count", ""},
	"offline_node_count":      {false, "offline_node_count", ""},
	"active_members":          {false, "active_members", ""},
	"inactive_members":        {false, "inactive_members", ""},
	"node_active_subscribers": {false, "active_subscribers_per_node", ""},

	"sims":              {false, "number_of_sims", ""},
	"active_sims":       {false, "active_sim_count", ""},
	"inactive_sims":     {false, "inactive_sim_count", ""},
	"package_sales":     {false, "package_sales_sum", ""},
	"data_usage":        {false, "data_usage", ""},
	"unit_health":       {false, "unit_health", ""},
	"unit_status":       {false, "unit_status", ""},
	"node_load":         {false, "node_load", ""},
	"cellular_uplink":   {false, "cellular_uplink", ""},
	"cellular_downlink": {false, "cellular_downlink", ""},
	"backhaul_uplink":   {false, "backhaul_uplink", ""},
	"backhaul_downlink": {false, "backhaul_downlink", ""},
	"backhaul_latency":  {false, "backhaul_latency", ""},
	"hwd_load":          {false, "hwd_load", ""},
	"memory_usage":      {false, "memory_usage", ""},
	"cpu_usage":         {false, "cpu_usage", ""},
	"disk_usage":        {false, "disk_usage", ""},
	"txpower":           {false, "txpower", ""},

	"unit_uptime":         {false, "unit_uptime", ""},
	"network_sales":       {false, "network_sales", ""},
	"network_data_volume": {false, "network_data_volume", ""},
	"network_active_ue":   {false, "network_active_ue", ""},
	"network_uptime":      {false, "network_uptime", ""},
	//

	//Health metrics
	"temperature_trx": {false, "trx_sensors_tempsensor1_temperature", ""},
	"temperature_com": {false, "com_sensors_tempsensor1_temperature_microprocessor", ""},

	"temperature_ctl": {false, "ctl_sensors_tempsensor_microprocessor", ""},
	"temperature_rfe": {false, "rfe_sensors_tempsensor_pa", ""},

	"temperature_S1_trx_hn": {false, "trx_sensors_tempsensor1_temperature", ""},
	"temperature_S2_trx_hn": {false, "trx_sensors_tempsensor2_temperature", ""},
	"temperature_S1_rfe_hn": {false, "rfe_sensors_tempsensor1_pa1", ""},
	"temperature_S2_rfe_hn": {false, "rfe_sensors_tempsensor2_pa2", ""},

	//Uptime Metrics
	"uptime_trx": {false, "trx_generic_system_uptime_seconds", ""},
	"uptime_com": {false, "com_generic_system_uptime_seconds", ""},
	"uptime_ctl": {false, "ctl_generic_system_uptime_seconds", ""},

	//Radio Metrics
	"tx_power": {false, "rfe_sensor_adc_tx_power", ""},
	"rx_power": {false, "rfe_sensor_adc_rx_power", ""},
	"pa_power": {false, "rfe_sensor_adc_pa_power", ""},

	//Resources Metrics
	"memory_trx_total": {false, "trx_memory_ddr_total", ""},
	"memory_trx_used":  {false, "trx_memory_ddr_used", ""},
	"memory_trx_free":  {false, "trx_memory_ddr_free", ""},

	"memory_com_total": {false, "com_memory_ddr_total", ""},
	"memory_com_used":  {false, "com_memory_ddr_used", ""},
	"memory_com_free":  {false, "com_memory_ddr_free", ""},

	"memory_ctl_total": {false, "ctl_memory_ddr_total", ""},
	"memory_ctl_used":  {false, "ctl_memory_ddr_used", ""},
	"memory_ctl_free":  {false, "ctl_memory_ddr_free", ""},

	//CPU Metrics (TRX, COM, CTL)
	"cpu_trx_usage":    {false, "trx_soc_cpu_usage", ""},
	"cpu_trx_c0_usage": {false, "trx_soc_cpu_core0_usage", ""},
	"cpu_trx_c1_usage": {false, "trx_soc_cpu_core1_usage", ""},
	"cpu_trx_c2_usage": {false, "trx_soc_cpu_core2_usage", ""},
	"cpu_trx_c3_usage": {false, "trx_soc_cpu_core3_usage", ""},

	"cpu_com_usage":    {false, "com_soc_cpu_usage", ""},
	"cpu_com_c0_usage": {false, "com_soc_cpu_core0_usage", ""},
	"cpu_com_c1_usage": {false, "com_soc_cpu_core1_usage", ""},
	"cpu_com_c2_usage": {false, "com_soc_cpu_core2_usage", ""},
	"cpu_com_c3_usage": {false, "com_soc_cpu_core3_usage", ""},

	"cpu_ctl_total": {false, "ctl_soc_cpu_usage", ""},
	"cpu_ctl_used":  {false, "ctl_soc_cpu_core0_usage", ""},

	//DISK Metrics (TRX, COM, CTL)
	"disk_trx_total": {false, "trx_storage_emmc_total", ""},
	"disk_trx_used":  {false, "trx_storage_emmc_used", ""},
	"disk_trx_free":  {false, "trx_storage_emmc_free", ""},

	"disk_com_total": {false, "com_storage_emmc_total", ""},
	"disk_com_used":  {false, "com_storage_emmc_used", ""},
	"disk_com_free":  {false, "com_storage_emmc_free", ""},

	"disk_ctl_total": {false, "ctl_storage_emmc_total", ""},
	"disk_ctl_used":  {false, "ctl_storage_emmc_used", ""},
	"disk_ctl_free":  {false, "ctl_storage_emmc_free", ""},

	//Power Level
	"power_level": {false, "trx_sensors_powermanagement_power", ""},

	// backhaul Metrics
	"network_latency":         {false, "process_open_fds", ""},
	"network_packet_loss":     {false, "process_open_fds", ""},
	"network_overall_status":  {false, "process_open_fds", ""},
	"network_throughput_up":   {false, "trx_lte_stack_throughput_uplink", ""},
	"network_throughput_down": {false, "trx_lte_stack_throughput_downlink", ""},

	// Solar Power Metrics
	"solar_panel_power":         {false, "solar_panel_power", ""},
	"solar_panel_voltage":       {false, "solar_panel_voltage", ""},
	"solar_panel_current":       {false, "solar_panel_current", ""},
	"battery_charge_percentage": {false, "battery_charge_percentage", ""},

	// Internet Switch Metrics
	"switch_port_status": {false, "switch_port_status", ""},
	"switch_port_speed":  {false, "switch_port_speed", ""},
	"switch_port_power":  {false, "switch_port_power", ""},

	//main backhaul
	"backhaul_speed":         {false, "backhaul_speed", ""},
	"main_backhaul_latency":  {false, "main_backhaul_latency", ""},
	"site_uptime_seconds":    {false, "site_uptime_seconds", ""},
	"site_uptime_percentage": {false, "site_uptime_percentage", ""},

	"backhaul_switch_port_status": {false, "backhaul_switch_port_status", ""},
	"backhaul_switch_port_speed":  {false, "backhaul_switch_port_speed", ""},
	"backhaul_switch_port_power":  {false, "backhaul_switch_port_power", ""},

	"solar_switch_port_status": {false, "solar_switch_port_status", ""},
	"solar_switch_port_speed":  {false, "solar_switch_port_speed", ""},
	"solar_switch_port_power":  {false, "solar_switch_port_power", ""},

	"node_switch_port_status": {false, "node_switch_port_status", ""},
	"node_switch_port_speed":  {false, "node_switch_port_speed", ""},
	"node_switch_port_power":  {false, "node_switch_port_power", ""},
}

type GrpcEndpoints struct {
	Timeout   time.Duration
	Exporter  string
	Sanitizer string
}

type HttpEndpoints struct {
	Timeout     time.Duration
	NodeMetrics string
}

func NewConfig() *Config {
	defaultCors := cors.DefaultConfig()
	defaultCors.AllowWildcard = true
	defaultCors.AllowOrigins = []string{"http://localhost", "https://localhost"}

	return &Config{
		BaseConfig: config.BaseConfig{
			DebugMode: false,
		},

		Services: GrpcEndpoints{
			Timeout:   3 * time.Second,
			Exporter:  "0.0.0.0:9090",
			Sanitizer: "sanitizer:9090",
		},
		HttpServices: HttpEndpoints{
			Timeout:     3 * time.Second,
			NodeMetrics: "http://localhost",
		},

		Server: rest.HttpConfig{
			Port: 8080,
			Cors: defaultCors,
		},
		MetricsServer: *config.DefaultMetrics(),

		MetricsConfig: &MetricsConfig{
			Metrics:             defaultPrometheusMetric,
			Timeout:             time.Second * 5,
			DefaultRateInterval: "5m",
		},
		Auth:   config.LoadAuthHostConfig("auth"),
		Period: time.Second * 5,
	}
}

func MergeMetricsConfigDefaults(c *Config) {
	if c == nil || c.MetricsConfig == nil {
		return
	}
	if c.MetricsConfig.Metrics == nil {
		c.MetricsConfig.Metrics = make(map[string]Metric)
	}
	for k, v := range defaultPrometheusMetric {
		if _, ok := c.MetricsConfig.Metrics[k]; !ok {
			c.MetricsConfig.Metrics[k] = v
		}
	}
}
