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

	"github.com/ukama/ukama/systems/analytics/customers/pkg"
	"github.com/ukama/ukama/systems/analytics/customers/pkg/query"
)

type Server struct {
	cfg *pkg.Config
	mux *http.ServeMux
}

func NewServer(cfg *pkg.Config) *Server {
	s := &Server{cfg: cfg, mux: http.NewServeMux()}

	s.mux.HandleFunc("/health", s.health)
	s.mux.HandleFunc("/overview", s.overview)
	s.mux.HandleFunc("/list", s.list)
	s.mux.HandleFunc("/search", s.search)
	s.mux.HandleFunc("/support", s.support)
	return s
}

func (s *Server) Run() error {
	return http.ListenAndServe(":"+s.cfg.ServicePort, s.mux)
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	s.writeJSON(w, map[string]any{"status": "ok", "service": "customers"})
}

func (s *Server) overview(w http.ResponseWriter, r *http.Request) {
	s.writeJSON(w, map[string]any{
		"service": "customers",
		"name":    "overview",
		"message": "Customers overview data",
		"data":    query.Demo("overview"),
	})
}

func (s *Server) list(w http.ResponseWriter, r *http.Request) {
	s.writeJSON(w, map[string]any{
		"service": "customers",
		"name":    "list",
		"message": "Customer list data",
		"data":    query.Demo("list"),
	})
}

func (s *Server) search(w http.ResponseWriter, r *http.Request) {
	s.writeJSON(w, map[string]any{
		"service": "customers",
		"name":    "search",
		"message": "Customer search data",
		"data":    query.Demo("search"),
	})
}

func (s *Server) support(w http.ResponseWriter, r *http.Request) {
	s.writeJSON(w, map[string]any{
		"service": "customers",
		"name":    "support",
		"message": "Customer support diagnosis data",
		"data":    query.Demo("support"),
	})
}

func (s *Server) writeJSON(w http.ResponseWriter, value any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(value)
}
