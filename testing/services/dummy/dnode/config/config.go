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
	NodeId   string   `json:"nodeId"`
	Profile  Profile  `json:"profile"`
	Scenario string   `json:"scenario"`
	Kpis     NodeKPIs `json:"kpis"`
}

type NodeKPI struct {
	Key    string
	Min    float64
	Normal float64
	Max    float64
}

type NodeKPIs struct {
	KPIs []NodeKPI
}

const PORT = 8085

var KPI_CONFIG = NodeKPIs{
	KPIs: []NodeKPI{
		{
			Key:    "network_sales",
			Min:    0,
			Normal: 10000,
			Max:    50000,
		},
		{
			Key:    "network_data_volume",
			Min:    0,
			Normal: 512000,
			Max:    1024000,
		},
		{
			Key:    "network_active_ue",
			Min:    0,
			Normal: 500,
			Max:    10000,
		},
		{
			Key: "network_uptime",
			Min: 0,
			Max: 2678400,
		},
		{
			Key: "unit_uptime",
			Min: 0,
			Max: 2678400,
		},
		{
			Key:    "unit_health",
			Min:    50,
			Normal: 80,
			Max:    100,
		},
		{
			Key:    "trx_lte_core_active_ue",
			Min:    80,
			Normal: 95,
			Max:    100,
		},
		{
			Key:    "node_load",
			Min:    50,
			Normal: 75,
			Max:    90,
		},
		{
			Key:    "cellular_uplink",
			Min:    1024,
			Normal: 5120,
			Max:    10240,
		},
		{
			Key:    "cellular_downlink",
			Min:    1024,
			Normal: 8192,
			Max:    10240,
		},
		{
			Key:    "backhaul_uplink",
			Min:    1024,
			Normal: 5120,
			Max:    10240,
		},
		{
			Key:    "backhaul_downlink",
			Min:    1024,
			Normal: 8192,
			Max:    10240,
		},
		{
			Key:    "backhaul_latency",
			Min:    30,
			Normal: 50,
			Max:    80,
		},
		{
			Key:    "hwd_load",
			Min:    50,
			Normal: 70,
			Max:    80,
		},
		{
			Key:    "memory_usage",
			Min:    40,
			Normal: 70,
			Max:    80,
		},
		{
			Key:    "cpu_usage",
			Min:    40,
			Normal: 70,
			Max:    80,
		},
		{
			Key:    "disk_usage",
			Min:    40,
			Normal: 70,
			Max:    80,
		},
		{
			Key:    "txpower",
			Min:    30,
			Normal: 60,
			Max:    95,
		},
	},
}
