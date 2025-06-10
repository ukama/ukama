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
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/ukama"
	cconfig "github.com/ukama/ukama/testing/common/config"
	cenums "github.com/ukama/ukama/testing/common/enums"
	"github.com/ukama/ukama/testing/services/dummy/dnode/config"
	"github.com/ukama/ukama/testing/services/dummy/dnode/utils"
)

var (
	serviceConfig = config.NewConfig()
	logger        = log.New()
)

type Server struct {
	orgName    string
	mu         sync.RWMutex
	amqpConfig cconfig.Queue
	coroutines map[string]chan config.WMessage
	server     *http.Server
}

func NewServer() (*Server, error) {
	orgname := os.Getenv("ORGNAME")
	if orgname == "" {
		return nil, fmt.Errorf("ORGNAME environment variable is required")
	}

	amqp := os.Getenv("AMQPCONFIG_URI")
	amqpUsername := os.Getenv("AMQPCONFIG_USERNAME")
	amqpPassword := os.Getenv("AMQPCONFIG_PASSWORD")

	if amqp == "" || amqpUsername == "" || amqpPassword == "" {
		return nil, fmt.Errorf("AMQP configuration environment variables are required")
	}

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
		server: &http.Server{
			Addr:         fmt.Sprintf(":%d", serviceConfig.Port),
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	}, nil
}

func init() {
	logger.SetFormatter(&log.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(log.InfoLevel)

	for _, kpi := range serviceConfig.KpiConfig.KPIs {
		prometheus.MustRegister(kpi.KPI)
	}
}

func (s *Server) Start() error {
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/update", s.updateHandler)
	http.HandleFunc("/online", s.onlineHandler)

	logger.Infof("Server starting on port %d", config.NewConfig().Port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	<-quit
	logger.Info("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		logger.Errorf("Server forced to shutdown: %v", err)
		return err
	}

	logger.Info("Server exited properly")
	return nil
}

func (s *Server) onlineHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	nodeId := r.URL.Query().Get("nodeid")
	nodeID, err := ukama.ValidateNodeId(nodeId)
	if err != nil {
		logger.Errorf("Invalid node ID: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	utils.PushNodeOnlineViaREST(s.amqpConfig, s.orgName, nodeID.String(), nil)
	logger.Infof("Online event pushed for NodeId: %s", nodeID.String())

	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.handleNodeOn(nodeID, "normal"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("NodeId: " + nodeID.String())); err != nil {
		logger.Errorf("Error writing response: %v", err)
	}
}

func (s *Server) updateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	nodeId := r.URL.Query().Get("nodeid")
	profile := r.URL.Query().Get("profile")
	scenario := r.URL.Query().Get("scenario")

	nodeID, err := ukama.ValidateNodeId(nodeId)
	if err != nil {
		logger.Errorf("Invalid node ID: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	logger.Infof("Updating coroutine, NodeId: %s, Profile: %s, Scenario: %s",
		nodeID.String(), profile, scenario)

	s.mu.Lock()
	defer s.mu.Unlock()

	scenarioType := cenums.ParseScenarioType(scenario)
	profileType := cenums.ParseProfileType(profile)

	var handlerErr error
	switch scenarioType {
	case cenums.SCENARIO_BACKHAUL_DOWN:
	case cenums.SCENARIO_NODE_OFF:
		handlerErr = s.handleNodeOff(nodeID, profile, scenarioType)

	case cenums.SCENARIO_NODE_ON:
		handlerErr = s.handleNodeOn(nodeID, profile)

	case cenums.SCENARIO_NODE_RESTART:
		handlerErr = s.handleNodeRestart(nodeID, profile, scenarioType)
	default:
		s.coroutines[nodeID.String()] <- config.WMessage{
			Profile:  profileType,
			Scenario: scenarioType,
		}
	}

	if handlerErr != nil {
		status := http.StatusInternalServerError
		if scenarioType == cenums.SCENARIO_NODE_OFF {
			status = http.StatusNotFound
		}
		http.Error(w, handlerErr.Error(), status)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("NodeId: " + nodeID.String())); err != nil {
		logger.Errorf("Error writing response: %v", err)
	}
}

func (s *Server) handleNodeOn(nodeID ukama.NodeID, profile string) error {
	if _, exists := s.coroutines[nodeID.String()]; !exists {
		updateChan := make(chan config.WMessage, 10)
		s.coroutines[nodeID.String()] = updateChan

		logger.Infof("Starting new coroutine for node %s", nodeID.String())
		utils.PushNodeOnlineViaREST(s.amqpConfig, s.orgName, nodeID.String(), nil)
		go utils.Worker(nodeID.String(), updateChan, config.WMessage{
			NodeId:   nodeID.String(),
			Profile:  cenums.ParseProfileType(profile),
			Scenario: cenums.SCENARIO_DEFAULT,
			Kpis:     serviceConfig.KpiConfig,
		})
	} else {
		logger.Infof("Coroutine already exists for NodeId: %s", nodeID.String())
	}
	return nil
}

func (s *Server) handleNodeOff(nodeID ukama.NodeID, profile string, scenarioType cenums.SCENARIOS) error {
	logger.Infof("Shutting down coroutine for node %s", nodeID.String())
	utils.PushNodeOffViaREST(s.amqpConfig, s.orgName, nodeID.String(), nil)
	updateChan, exists := s.coroutines[nodeID.String()]
	if !exists {
		return fmt.Errorf("coroutine not found for node %s", nodeID.String())
	}

	updateChan <- config.WMessage{
		NodeId:   nodeID.String(),
		Kpis:     serviceConfig.KpiConfig,
		Profile:  cenums.ParseProfileType(profile),
		Scenario: scenarioType,
	}

	close(updateChan)
	delete(s.coroutines, nodeID.String())
	return nil
}

func (s *Server) handleNodeRestart(nodeID ukama.NodeID, profile string, scenarioType cenums.SCENARIOS) error {
	logger.Infof("Restarting node %s", nodeID.String())

	err := s.handleNodeOff(nodeID, profile, scenarioType)
	if err != nil {
		return err
	}
	time.Sleep(15 * time.Second)
	return s.handleNodeOn(nodeID, profile)
}

func main() {
	server, err := NewServer()
	if err != nil {
		logger.Fatalf("Failed to create server: %v", err)
	}

	if err := server.Start(); err != nil {
		logger.Fatalf("Server error: %v", err)
	}
}
