/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package rest

import "net/http"

type RouterConfig struct {
	Port string
}

type Router struct {
	mux *http.ServeMux
	cfg RouterConfig
}

func NewRouter(clients Clients, cfg RouterConfig) *Router {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/analytics/business/", clients.ProxyBusiness)
	mux.HandleFunc("/v1/analytics/network/", clients.ProxyNetwork)
	mux.HandleFunc("/v1/analytics/customers/", clients.ProxyCustomers)
	mux.HandleFunc("/v1/analytics/collector/", clients.ProxyCollector)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok","service":"analytics-api-gateway"}`))
	})
	return &Router{mux: mux, cfg: cfg}
}

func (r *Router) Run() error {
	return http.ListenAndServe(":"+r.cfg.Port, r.mux)
}
