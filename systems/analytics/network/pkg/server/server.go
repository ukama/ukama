/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package server

import (
	"encoding/json"
	"net/http"

	"github.com/ukama/ukama/systems/analytics/network/pkg"
	"github.com/ukama/ukama/systems/analytics/network/pkg/query"
)

type Server struct {
	cfg *pkg.Config
	mux *http.ServeMux
}

func NewServer(cfg *pkg.Config) *Server {
	s := &Server{cfg: cfg, mux: http.NewServeMux()}

	s.mux.HandleFunc("/health", s.health)
	s.mux.HandleFunc("/overview", s.overview)
	s.mux.HandleFunc("/topology", s.topology)
	s.mux.HandleFunc("/sites", s.sites)
	s.mux.HandleFunc("/nodes", s.nodes)
	s.mux.HandleFunc("/radio", s.radio)
	s.mux.HandleFunc("/backhaul", s.backhaul)
	s.mux.HandleFunc("/power", s.power)
	s.mux.HandleFunc("/alarms", s.alarms)
	s.mux.HandleFunc("/metrics", s.metrics)
	s.mux.HandleFunc("/events", s.events)
	s.mux.HandleFunc("/maintenance", s.maintenance)
	return s
}

func (s *Server) Run() error {
	return http.ListenAndServe(":"+s.cfg.ServicePort, s.mux)
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	s.writeJSON(w, map[string]any{"status": "ok", "service": "network"})
}

func (s *Server) overview(w http.ResponseWriter, r *http.Request) {
	s.writeJSON(w, map[string]any{
		"service": "network",
		"name":    "overview",
		"message": "Network overview data",
		"data":    query.Demo("overview"),
	})
}

func (s *Server) topology(w http.ResponseWriter, r *http.Request) {
	s.writeJSON(w, map[string]any{
		"service": "network",
		"name":    "topology",
		"message": "Topology data",
		"data":    query.Demo("topology"),
	})
}

func (s *Server) sites(w http.ResponseWriter, r *http.Request) {
	s.writeJSON(w, map[string]any{
		"service": "network",
		"name":    "sites",
		"message": "Network site health data",
		"data":    query.Demo("sites"),
	})
}

func (s *Server) nodes(w http.ResponseWriter, r *http.Request) {
	s.writeJSON(w, map[string]any{
		"service": "network",
		"name":    "nodes",
		"message": "Node health data",
		"data":    query.Demo("nodes"),
	})
}

func (s *Server) radio(w http.ResponseWriter, r *http.Request) {
	s.writeJSON(w, map[string]any{
		"service": "network",
		"name":    "radio",
		"message": "Radio LTE data",
		"data":    query.Demo("radio"),
	})
}

func (s *Server) backhaul(w http.ResponseWriter, r *http.Request) {
	s.writeJSON(w, map[string]any{
		"service": "network",
		"name":    "backhaul",
		"message": "Backhaul health data",
		"data":    query.Demo("backhaul"),
	})
}

func (s *Server) power(w http.ResponseWriter, r *http.Request) {
	s.writeJSON(w, map[string]any{
		"service": "network",
		"name":    "power",
		"message": "Power health data",
		"data":    query.Demo("power"),
	})
}

func (s *Server) alarms(w http.ResponseWriter, r *http.Request) {
	s.writeJSON(w, map[string]any{
		"service": "network",
		"name":    "alarms",
		"message": "Alarm data",
		"data":    query.Demo("alarms"),
	})
}

func (s *Server) metrics(w http.ResponseWriter, r *http.Request) {
	s.writeJSON(w, map[string]any{
		"service": "network",
		"name":    "metrics",
		"message": "Metric data",
		"data":    query.Demo("metrics"),
	})
}

func (s *Server) events(w http.ResponseWriter, r *http.Request) {
	s.writeJSON(w, map[string]any{
		"service": "network",
		"name":    "events",
		"message": "Network events",
		"data":    query.Demo("events"),
	})
}

func (s *Server) maintenance(w http.ResponseWriter, r *http.Request) {
	s.writeJSON(w, map[string]any{
		"service": "network",
		"name":    "maintenance",
		"message": "Maintenance state",
		"data":    query.Demo("maintenance"),
	})
}

func (s *Server) writeJSON(w http.ResponseWriter, value any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(value)
}
