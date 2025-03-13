/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
package main

import (
	"fmt"
	"os"

	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/ukama"
	cconfig "github.com/ukama/ukama/testing/common/config"
	cenums "github.com/ukama/ukama/testing/common/enums"
	"github.com/ukama/ukama/testing/services/dummy/dnode/config"
	"github.com/ukama/ukama/testing/services/dummy/dnode/utils"
)

type Server struct {
	orgName    string
	mu         sync.Mutex
	amqpConfig cconfig.Queue
	coroutines map[string]chan config.WMessage
}

func NewServer() *Server {
	orgname := os.Getenv("ORGNAME")
	amqp := os.Getenv("AMQPCONFIG_URI")
	amqpUsername := os.Getenv("AMQPCONFIG_USERNAME")
	amqpPassword := os.Getenv("AMQPCONFIG_PASSWORD")
	return &Server{
		orgName: orgname,
		amqpConfig: cconfig.Queue{
			Uri:      amqp,
			Vhost:    "%2F",
			Username: amqpUsername,
			Password: amqpPassword,
			Exchange: "amq.topic",
		},
		coroutines: make(map[string]chan config.WMessage),
	}
}

func init() {
	for _, kpi := range config.KPI_CONFIG.KPIs {
		prometheus.MustRegister(kpi.KPI)
	}
}

func main() {
	server := NewServer()

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/update", server.updateHandler)
	http.HandleFunc("/online", server.onlineHandler)
	log.Printf("Server listening on port %d", config.PORT)
	err := http.ListenAndServe(fmt.Sprintf(":%d", config.PORT), nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func (s *Server) onlineHandler(w http.ResponseWriter, r *http.Request) {
	nodeId := r.URL.Query().Get("nodeid")
	nodeID, err := ukama.ValidateNodeId(nodeId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	utils.PushNodeOnlineViaREST(s.amqpConfig, s.orgName, nodeID.String(), nil)

	log.Printf("Online event pushed for NodeId: %s", nodeID.String())

	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.coroutines[nodeID.String()]
	if !exists {
		updateChan := make(chan config.WMessage, 10)
		s.coroutines[nodeID.String()] = updateChan

		log.Printf("Starting coroutine, NodeId: %s, Profile: %d, Scenario: %s", nodeID.String(), cenums.PROFILE_NORMAL, cenums.SCENARIO_DEFAULT)
		go utils.Worker(nodeID.String(), updateChan, config.WMessage{NodeId: nodeID.String(), Profile: cenums.PROFILE_NORMAL, Scenario: cenums.SCENARIO_DEFAULT, Kpis: config.KPI_CONFIG})
	} else {
		log.Printf("Coroutine already exists for NodeId: %s", nodeID.String())
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("NodeId: " + nodeID.String()))
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

	log.Printf("Updating coroutine, NodeId: %s, Profile: %s, Scenario: %s", nodeID.String(), profile, scenario)

	s.mu.Lock()
	defer s.mu.Unlock()

	updateChan, exists := s.coroutines[nodeID.String()]
	if !exists {
		http.Error(w, "Coroutine not found", http.StatusNotFound)
		return
	}

	updateChan <- config.WMessage{
		NodeId:   nodeID.String(),
		Kpis:     config.KPI_CONFIG,
		Profile:  cenums.ParseProfileType(profile),
		Scenario: cenums.ParseScenarioType(scenario),
	}

	if cenums.ParseScenarioType(scenario) == cenums.SCENARIO_BACKHAUL_DOWN ||
		cenums.ParseScenarioType(scenario) == cenums.SCENARIO_NODE_OFF {
		log.Printf("Scenario is: %s, which leads to coroutine shutdown.", scenario)
		updateChan, exists := s.coroutines[nodeID.String()]
		if exists {
			close(updateChan)
			delete(s.coroutines, nodeID.String())
		}
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("NodeId: " + nodeID.String()))
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}
