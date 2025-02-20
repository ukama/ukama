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
	"os"

	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/testing/services/dummy/dnode/config"
	"github.com/ukama/ukama/testing/services/dummy/dnode/utils"
)

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

type Server struct {
	orgName     string
	mu          sync.Mutex
	amqpConfig  config.AmqpConfig
	cancelFuncs map[string]context.CancelFunc
	coroutines  map[string]chan config.WMessage
}

func NewServer(orgName string) *Server {
	return &Server{
		orgName: orgName,
		amqpConfig: config.AmqpConfig{
			Uri:      "http://rabbitmq:15672",
			Username: "guest",
			Password: "guest",
			Exchange: "amq.topic",
			Vhost:    "%2F",
		},
		cancelFuncs: make(map[string]context.CancelFunc),
		coroutines:  make(map[string]chan config.WMessage),
	}
}

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

func main() {
	ORGNAME := os.Getenv("ORGNAME")
	server := NewServer(ORGNAME)

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/update", server.updateHandler)
	http.HandleFunc("/online", server.onlineHandler)
}

func (s *Server) onlineHandler(w http.ResponseWriter, r *http.Request) {
	nodeId := r.URL.Query().Get("nodeid")
	nodeID, err := ukama.ValidateNodeId(nodeId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	utils.PushNodeResetViaREST(s.amqpConfig, s.orgName, nodeID.String(), nil)

	log.Printf("Online event pushed for NodeId: %s", nodeID)

	s.mu.Lock()
	defer s.mu.Unlock()

	updateChan := make(chan config.WMessage)

	s.coroutines[nodeId] = updateChan

	log.Printf("Starting coroutine, NodeId: %s, Profile: %s, Scenario: %s", nodeID, config.PROFILE_NORMAL, "DEFAULT")

	go utils.Worker(nodeId, updateChan, config.WMessage{nodeId, config.PROFILE_NORMAL, "DEFAULT", config.KPI_CONFIG})

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("NodeId: " + nodeId))
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

func (s *Server) updateHandler(w http.ResponseWriter, r *http.Request) {
	nodeId := r.URL.Query().Get("nodeid")
	profile := r.URL.Query().Get("profile")
	scenario := r.URL.Query().Get("scenario")
	nodeID, err := ukama.ValidateNodeId(nodeId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Updating coroutine, NodeId: %s, Profile: %s, Scenario: %s", nodeID, profile, scenario)

	updateChan, exists := s.coroutines[nodeId]
	if !exists {
		http.Error(w, "Coroutine not found", http.StatusNotFound)
		return
	}

	updateChan <- config.WMessage{nodeId, config.ParseProfileType(profile), scenario, config.KPI_CONFIG}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("NodeId: " + nodeId))
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}
