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
	"github.com/ukama/ukama/testing/services/dummy/dnode/config"
	"github.com/ukama/ukama/testing/services/dummy/dnode/utils"
)

type Server struct {
	orgName    string
	mu         sync.Mutex
	amqpConfig config.AmqpConfig
	coroutines map[string]chan config.WMessage
}

func NewServer() *Server {
	orgname := os.Getenv("ORGNAME")
	amqp := os.Getenv("AMQPCONFIG_URI")
	amqpUsername := os.Getenv("AMQPCONFIG_USERNAME")
	amqpPassword := os.Getenv("AMQPCONFIG_PASSWORD")
	return &Server{
		orgName: orgname,
		amqpConfig: config.AmqpConfig{
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
	http.ListenAndServe(fmt.Sprintf(":%d", config.PORT), nil)

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

	_, exists := s.coroutines[nodeId]
	if !exists {
		updateChan := make(chan config.WMessage)
		s.coroutines[nodeId] = updateChan

		log.Printf("Starting coroutine, NodeId: %s, Profile: %d, Scenario: %s", nodeID, config.PROFILE_NORMAL, "DEFAULT")
		go utils.Worker(nodeId, updateChan, config.WMessage{NodeId: nodeId, Profile: config.PROFILE_NORMAL, Scenario: "DEFAULT", Kpis: config.KPI_CONFIG})
	} else {
		log.Printf("Coroutine already exists for NodeId: %s. Please use /update.", nodeID)
	}

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

	updateChan <- config.WMessage{
		NodeId:   nodeId,
		Scenario: scenario,
		Kpis:     config.KPI_CONFIG,
		Profile:  config.ParseProfileType(profile),
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("NodeId: " + nodeId))
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}
