/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
package config

import (
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

type AmqpConfig struct {
	Uri      string `default:"http://rabbitmq:15672"`
	Username string `default:"guest"`
	Password string `default:"guest"`
	Exchange string `default:"amq.topic"`
	Vhost    string `default:"%2F"`
}

type Profile uint8

const (
	PROFILE_NORMAL Profile = 0
	PROFILE_MIN    Profile = 1
	PROFILE_MAX    Profile = 2
)

func ParseProfileType(value string) Profile {
	i, err := strconv.Atoi(value)
	if err == nil {
		return Profile(i)
	}

	t := map[string]Profile{"normal": 0, "min": 1, "max": 2}

	v, ok := t[strings.ToLower(value)]
	if !ok {
		return Profile(0)
	}

	return Profile(v)
}

type WMessage struct {
	NodeId   string    `json:"nodeId"`
	Profile  Profile   `json:"profile"`
	Scenario SCENARIOS `json:"scenario"`
	Kpis     NodeKPIs  `json:"kpis"`
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

const PORT = 8085

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
			Min:    50,
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
			Min:    50,
			Normal: 75,
			Max:    90,
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
			Min:    1024,
			Normal: 5120,
			Max:    10240,
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
			Min:    1024,
			Normal: 8192,
			Max:    10240,
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
			Min:    1024,
			Normal: 5120,
			Max:    10240,
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
			Min:    1024,
			Normal: 8192,
			Max:    10240,
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
			Min:    30,
			Normal: 50,
			Max:    80,
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
			Min:    50,
			Normal: 70,
			Max:    80,
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
			Min:    40,
			Normal: 70,
			Max:    80,
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
			Min:    40,
			Normal: 70,
			Max:    80,
			KPI: prometheus.NewGaugeVec(
				prometheus.GaugeOpts{
					Name: "cpu_usage",
					Help: "Cpu usage",
				},
				[]string{"nodeid"},
			),
		},
		{
			Key:    "disk_usage",
			Min:    40,
			Normal: 70,
			Max:    80,
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
			Min:    30,
			Normal: 60,
			Max:    95,
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

type SCENARIOS string

const (
	SCENARIO_DEFAULT                SCENARIOS = "default"
	SCENARIO_BACKHAUL_DOWN          SCENARIOS = "backhaul_down"
	SCENARIO_BACKHAUL_DOWNLINK_DOWN SCENARIOS = "backhaul_downlink_down"
	SCENARIO_SOLAR_DOWN             SCENARIOS = "solar_down"
	SCENARIO_SWITCH_OFF             SCENARIOS = "switch_off"
	SCENARIO_SITE_RESTART           SCENARIOS = "site_restart"
	SCENARIO_NODE_OFF               SCENARIOS = "node_off"
)

func ParseScenarioType(value string) SCENARIOS {
	t := map[string]SCENARIOS{
		"default":                SCENARIO_DEFAULT,
		"backhaul_down":          SCENARIO_BACKHAUL_DOWN,
		"backhaul_downlink_down": SCENARIO_BACKHAUL_DOWNLINK_DOWN,
		"solar_down":             SCENARIO_SOLAR_DOWN,
		"switch_off":             SCENARIO_SWITCH_OFF,
		"site_restart":           SCENARIO_SITE_RESTART,
		"node_off":               SCENARIO_NODE_OFF,
	}

	v, ok := t[strings.ToLower(value)]
	if !ok {
		return SCENARIO_DEFAULT
	}

	return SCENARIOS(v)
}
