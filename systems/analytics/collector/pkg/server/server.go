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

	"github.com/ukama/ukama/systems/analytics/collector/pkg"
	"github.com/ukama/ukama/systems/analytics/collector/pkg/query"
)

type Server struct {
	cfg *pkg.Config
	mux *http.ServeMux
}

func NewServer(cfg *pkg.Config) *Server {
	s := &Server{cfg: cfg, mux: http.NewServeMux()}

	s.mux.HandleFunc("/health", s.health)
	s.mux.HandleFunc("/refresh", s.refresh)
	s.mux.HandleFunc("/events", s.events)
	s.mux.HandleFunc("/routes", s.routes)
	return s
}

func (s *Server) Run() error {
	return http.ListenAndServe(":"+s.cfg.ServicePort, s.mux)
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	s.writeJSON(w, map[string]any{"status": "ok", "service": "collector"})
}

func (s *Server) refresh(w http.ResponseWriter, r *http.Request) {
	s.writeJSON(w, map[string]any{
		"service": "collector",
		"name":    "refresh",
		"message": "Refresh queued",
		"data":    query.Demo("refresh"),
	})
}

func (s *Server) events(w http.ResponseWriter, r *http.Request) {
	s.writeJSON(w, map[string]any{
		"service": "collector",
		"name":    "events",
		"message": "Event ingestion endpoint",
		"data":    query.Demo("events"),
	})
}

func (s *Server) routes(w http.ResponseWriter, r *http.Request) {
	s.writeJSON(w, map[string]any{
		"service": "collector",
		"name":    "routes",
		"message": "Event subscription routes",
		"data":    query.Demo("routes"),
	})
}

func (s *Server) writeJSON(w http.ResponseWriter, value any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(value)
}
