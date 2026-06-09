/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package clients_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tj/assert"

	"github.com/ukama/ukama/systems/analytics/collector/pkg/clients"
)

func testServer() *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/networks", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"networks":[{"id":"n1","name":"Net"}]}`))
	})
	mux.HandleFunc("/v1/sites", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"sites":[{"id":"s1","name":"Site","network_id":"n1"}]}`))
	})
	mux.HandleFunc("/v1/subscribers", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"subscribers":[{"subscriber_id":"sub1"}]}`))
	})
	mux.HandleFunc("/v1/packages", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"packages":[{"uuid":"p1","name":"Starter"}]}`))
	})
	mux.HandleFunc("/v1/metrics", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"metrics":[{"metric":"x","value":1}]}`))
	})
	mux.HandleFunc("/v1/nodes", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"nodes":[{"node_id":"node-1","name":"Node"}]}`))
	})
	mux.HandleFunc("/v1/components", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"components":[{"id":"c1"}]}`))
	})
	mux.HandleFunc("/v1/account", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"account":{"balance_cents":1000}}`))
	})
	return httptest.NewServer(mux)
}

func TestRegistryClient(t *testing.T) {
	s := testServer()
	defer s.Close()
	c := clients.NewRegistryClient(s.URL)

	nets, err := c.GetNetworks()
	assert.NoError(t, err)
	assert.Len(t, nets, 1)

	sites, err := c.GetSites()
	assert.NoError(t, err)
	assert.Len(t, sites, 1)
}

func TestSubscriberClient(t *testing.T) {
	s := testServer()
	defer s.Close()
	out, err := clients.NewSubscriberClient(s.URL).GetSubscribers()
	assert.NoError(t, err)
	assert.Len(t, out, 1)
}

func TestDataplanClient(t *testing.T) {
	s := testServer()
	defer s.Close()
	out, err := clients.NewDataplanClient(s.URL).GetPackages()
	assert.NoError(t, err)
	assert.Len(t, out, 1)
}

func TestMetricsClient(t *testing.T) {
	s := testServer()
	defer s.Close()
	out, err := clients.NewMetricsClient(s.URL).GetLatestMetrics()
	assert.NoError(t, err)
	assert.Len(t, out, 1)
}

func TestNodeClient(t *testing.T) {
	s := testServer()
	defer s.Close()
	out, err := clients.NewNodeClient(s.URL).GetNodes()
	assert.NoError(t, err)
	assert.Len(t, out, 1)
}

func TestInventoryClient(t *testing.T) {
	s := testServer()
	defer s.Close()
	out, err := clients.NewInventoryClient(s.URL).GetComponents()
	assert.NoError(t, err)
	assert.Len(t, out, 1)
}

func TestBillingClient(t *testing.T) {
	s := testServer()
	defer s.Close()
	out, err := clients.NewBillingClient(s.URL).GetAccount()
	assert.NoError(t, err)
	assert.NotNil(t, out)
}

func TestClient_DeserializationError(t *testing.T) {
	bad := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`not-json`))
	}))
	defer bad.Close()

	_, err := clients.NewRegistryClient(bad.URL).GetNetworks()
	assert.Error(t, err)
}

// dead points every client at a closed port so the HTTP GET fails, exercising
// each method's request-failure branch.
const dead = "http://127.0.0.1:1"

func TestClient_RequestFailures(t *testing.T) {
	reg := clients.NewRegistryClient(dead)
	_, err := reg.GetNetworks()
	assert.Error(t, err)
	_, err = reg.GetSites()
	assert.Error(t, err)

	_, err = clients.NewSubscriberClient(dead).GetSubscribers()
	assert.Error(t, err)
	_, err = clients.NewDataplanClient(dead).GetPackages()
	assert.Error(t, err)
	_, err = clients.NewMetricsClient(dead).GetLatestMetrics()
	assert.Error(t, err)
	_, err = clients.NewNodeClient(dead).GetNodes()
	assert.Error(t, err)
	_, err = clients.NewInventoryClient(dead).GetComponents()
	assert.Error(t, err)
	_, err = clients.NewBillingClient(dead).GetAccount()
	assert.Error(t, err)
}
