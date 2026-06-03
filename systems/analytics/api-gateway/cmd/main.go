/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ukama/ukama/systems/analytics/api-gateway/cmd/version"
	"github.com/ukama/ukama/systems/analytics/api-gateway/pkg"
	"github.com/ukama/ukama/systems/analytics/api-gateway/pkg/client"
	"github.com/ukama/ukama/systems/analytics/api-gateway/pkg/rest"
)

type clientsSet struct {
	business  *client.BusinessClient
	network   *client.NetworkClient
	customers *client.CustomersClient
	collector *client.CollectorClient
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "version" {
		log.Println(version.Version)
		return
	}
	cfg := pkg.NewConfig()
	clients := &clientsSet{
		business:  client.NewBusinessClient(cfg.BusinessService, cfg.Timeout),
		network:   client.NewNetworkClient(cfg.NetworkService, cfg.Timeout),
		customers: client.NewCustomersClient(cfg.CustomersService, cfg.Timeout),
		collector: client.NewCollectorClient(cfg.CollectorService, cfg.Timeout),
	}
	r := rest.NewRouter(clients, rest.RouterConfig{Port: cfg.ServerPort})
	log.Printf("analytics api-gateway listening on :%s", cfg.ServerPort)
	log.Fatal(r.Run())
}

func (c *clientsSet) ProxyBusiness(w http.ResponseWriter, r *http.Request) {
	c.business.Proxy(w, r, "/v1/analytics/business")
}

func (c *clientsSet) ProxyNetwork(w http.ResponseWriter, r *http.Request) {
	c.network.Proxy(w, r, "/v1/analytics/network")
}

func (c *clientsSet) ProxyCustomers(w http.ResponseWriter, r *http.Request) {
	c.customers.Proxy(w, r, "/v1/analytics/customers")
}

func (c *clientsSet) ProxyCollector(w http.ResponseWriter, r *http.Request) {
	c.collector.Proxy(w, r, "/v1/analytics/collector")
}
