/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
package main

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Profile uint8

const (
	PROFILE_NORMAL Profile = 0
	PROFILE_MIN    Profile = 1
	PROFILE_MAX    Profile = 2
)

type NodeKPIs struct {
	Key    string
	Min    float64
	Normal float64
	Max    float64
}

var (
	network_sales = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "network_sales",
			Help: "Overall network sales",
		},
		[]string{"nodeid"},
	)
	network_data_volume = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "network_data_volume",
			Help: "Network data volume",
		},
		[]string{"nodeid"},
	)
	network_active_ue = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "network_active_ue",
			Help: "Active subscriber within the network",
		},
		[]string{"nodeid"},
	)
	network_uptime = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "network_uptime",
			Help: "Network uptime",
		},
		[]string{"nodeid"},
	)
	unit_uptime = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "unit_uptime",
			Help: "Node uptime",
		},
		[]string{"nodeid"},
	)
	unit_status = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "unit_status",
			Help: "Unit status",
		},
		[]string{"nodeid"},
	)
	unit_health = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "unit_health",
			Help: "Health status of the unit",
		},
		[]string{"nodeid"},
	)
	node_load = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "node_load",
			Help: "Load on the node",
		},
		[]string{"nodeid"},
	)
	subscriber_active = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "trx_lte_core_active_ue",
			Help: "Subscriber active",
		},
		[]string{"nodeid"},
	)
	cellular_uplink = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cellular_uplink",
			Help: "Cellular uplink",
		},
		[]string{"nodeid"},
	)
	cellular_downlink = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cellular_downlink",
			Help: "Cellular downlink",
		},
		[]string{"nodeid"},
	)
	backhaul_uplink = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "backhaul_uplink",
			Help: "Backhaul downlink",
		},
		[]string{"nodeid"},
	)
	backhaul_downlink = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "backhaul_downlink",
			Help: "Backhaul downlink",
		},
		[]string{"nodeid"},
	)
	backhaul_latency = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "backhaul_latency",
			Help: "Backhaul latency",
		},
		[]string{"nodeid"},
	)
	hwd_load = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "hwd_load",
			Help: "Hardware load",
		},
		[]string{"nodeid"},
	)
	memory_usage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "memory_usage",
			Help: "Memory usage",
		},
		[]string{"nodeid"},
	)
	cpu_usage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cpu_usage",
			Help: "Cpu usage",
		},
		[]string{"nodeid"},
	)
	disk_usage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "disk_usage",
			Help: "Disk usage",
		},
		[]string{"nodeid"},
	)
	txpower = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "txpower",
			Help: "TX power",
		},
		[]string{"nodeid"},
	)
)

func init() {
	prometheus.MustRegister(network_sales)
	prometheus.MustRegister(network_data_volume)
	prometheus.MustRegister(network_active_ue)
	prometheus.MustRegister(network_uptime)
	prometheus.MustRegister(unit_uptime)
	prometheus.MustRegister(unit_status)
	prometheus.MustRegister(unit_health)
	prometheus.MustRegister(node_load)
	prometheus.MustRegister(subscriber_active)
	prometheus.MustRegister(cellular_uplink)
	prometheus.MustRegister(cellular_downlink)
	prometheus.MustRegister(backhaul_uplink)
	prometheus.MustRegister(backhaul_downlink)
	prometheus.MustRegister(backhaul_latency)
	prometheus.MustRegister(hwd_load)
	prometheus.MustRegister(memory_usage)
	prometheus.MustRegister(cpu_usage)
	prometheus.MustRegister(disk_usage)
	prometheus.MustRegister(txpower)
}

func generateRandomData(ctx context.Context, nodeId string, profile Profile) {
	config := []NodeKPIs{
		// 0-Min, Min-Normal, Normal-Max
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
	}
	count := 1.0
	increment := 10.0
	for {
		select {
		case <-ctx.Done():
			return
		default:
			labels := prometheus.Labels{"nodeid": nodeId}
			values := make(map[string]float64)

			for _, kpi := range config {
				if kpi.Key == "network_uptime" || kpi.Key == "unit_uptime" {
					values[kpi.Key] = count
				} else if kpi.Key == "network_sales" {
					values[kpi.Key] = increment / 12
				} else if kpi.Key == "network_data_volume" {
					values[kpi.Key] = increment
				} else {
					switch profile {
					case PROFILE_MIN:
						values[kpi.Key] = kpi.Min + rand.Float64()*(kpi.Normal-kpi.Min)*0.1
					case PROFILE_MAX:
						values[kpi.Key] = kpi.Normal + rand.Float64()*(kpi.Max-kpi.Normal)*0.1
					default:
						values[kpi.Key] = kpi.Min + rand.Float64()*(kpi.Normal-kpi.Min)*0.1
					}
				}
			}

			network_sales.With(labels).Set(values["network_sales"])
			network_data_volume.With(labels).Set(values["network_data_volume"])
			network_active_ue.With(labels).Set(values["network_active_ue"])
			network_uptime.With(labels).Set(values["network_uptime"])
			unit_uptime.With(labels).Set(values["unit_uptime"])
			unit_status.With(labels).Set(1)
			unit_health.With(labels).Set(values["unit_health"])
			node_load.With(labels).Set(values["node_load"])
			subscriber_active.With(labels).Set(values["trx_lte_core_active_ue"])
			cellular_uplink.With(labels).Set(values["cellular_uplink"])
			cellular_downlink.With(labels).Set(values["cellular_downlink"])
			backhaul_uplink.With(labels).Set(values["backhaul_uplink"])
			backhaul_downlink.With(labels).Set(values["backhaul_downlink"])
			backhaul_latency.With(labels).Set(values["backhaul_latency"])
			hwd_load.With(labels).Set(values["hwd_load"])
			memory_usage.With(labels).Set(values["memory_usage"])
			cpu_usage.With(labels).Set(values["cpu_usage"])
			disk_usage.With(labels).Set(values["disk_usage"])
			txpower.With(labels).Set(values["txpower"])

			time.Sleep(1 * time.Second)
			count++
			increment += 10
		}
	}
}

type Server struct {
	cancelFuncs map[string]context.CancelFunc
	mu          sync.Mutex
}

func NewServer() *Server {
	return &Server{
		cancelFuncs: make(map[string]context.CancelFunc),
	}
}

func (s *Server) metricsHandler(w http.ResponseWriter, r *http.Request) {
	nodeId := r.URL.Query().Get("nodeid")
	profileStr := r.URL.Query().Get("profile")

	if nodeId == "" || profileStr == "" {
		http.Error(w, "Missing nodeid or profile parameter", http.StatusBadRequest)
		return
	}

	var profile Profile
	switch profileStr {
	case "normal":
		profile = PROFILE_NORMAL
	case "min":
		profile = PROFILE_MIN
	case "max":
		profile = PROFILE_MAX
	default:
		http.Error(w, "Invalid profile parameter", http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	if cancelFunc, exists := s.cancelFuncs[nodeId]; exists {
		cancelFunc()
	}
	newCtx, cancelFunc := context.WithCancel(context.Background())
	s.cancelFuncs[nodeId] = cancelFunc
	s.mu.Unlock()

	go generateRandomData(newCtx, nodeId, profile)
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("Metrics generation started"))
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

func main() {
	server := NewServer()

	http.HandleFunc("/start-metrics", server.metricsHandler)
	http.Handle("/metrics", promhttp.Handler())

	port := "8085"
	log.Printf("Server listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
