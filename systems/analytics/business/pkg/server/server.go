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

	"github.com/ukama/ukama/systems/analytics/business/pkg"
	"github.com/ukama/ukama/systems/analytics/business/pkg/query"
)

type Server struct {
	cfg *pkg.Config
	mux *http.ServeMux
}

func NewServer(cfg *pkg.Config) *Server {
	s := &Server{cfg: cfg, mux: http.NewServeMux()}

	s.mux.HandleFunc("/health", s.health)
	s.mux.HandleFunc("/home", s.home)
	s.mux.HandleFunc("/sales/overview", s.salesOverview)
	s.mux.HandleFunc("/sales/packages", s.salesPackages)
	s.mux.HandleFunc("/sites", s.sites)
	s.mux.HandleFunc("/inventory", s.inventory)
	return s
}

func (s *Server) Run() error {
	return http.ListenAndServe(":"+s.cfg.ServicePort, s.mux)
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	s.writeJSON(w, map[string]any{"status": "ok", "service": "business"})
}

func (s *Server) home(w http.ResponseWriter, r *http.Request) {
	s.writeJSON(w, map[string]any{
		"service": "business",
		"name":    "home",
		"message": "Business home data",
		"data":    query.Demo("home"),
	})
}

func (s *Server) salesOverview(w http.ResponseWriter, r *http.Request) {
	s.writeJSON(w, map[string]any{
		"service": "business",
		"name":    "salesOverview",
		"message": "Sales overview data",
		"data":    query.Demo("salesOverview"),
	})
}

func (s *Server) salesPackages(w http.ResponseWriter, r *http.Request) {
	s.writeJSON(w, map[string]any{
		"service": "business",
		"name":    "salesPackages",
		"message": "Package performance data",
		"data":    query.Demo("salesPackages"),
	})
}

func (s *Server) sites(w http.ResponseWriter, r *http.Request) {
	s.writeJSON(w, map[string]any{
		"service": "business",
		"name":    "sites",
		"message": "Business site data",
		"data":    query.Demo("sites"),
	})
}

func (s *Server) inventory(w http.ResponseWriter, r *http.Request) {
	s.writeJSON(w, map[string]any{
		"service": "business",
		"name":    "inventory",
		"message": "Inventory readiness data",
		"data":    query.Demo("inventory"),
	})
}

func (s *Server) writeJSON(w http.ResponseWriter, value any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(value)
}
