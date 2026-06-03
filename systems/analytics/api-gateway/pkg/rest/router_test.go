/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package rest

import (
	"net/http"
	"testing"
)

type mockClients struct{}

func (m mockClients) ProxyBusiness(http.ResponseWriter, *http.Request)  {}
func (m mockClients) ProxyNetwork(http.ResponseWriter, *http.Request)   {}
func (m mockClients) ProxyCustomers(http.ResponseWriter, *http.Request) {}
func (m mockClients) ProxyCollector(http.ResponseWriter, *http.Request) {}

func TestNewRouter(t *testing.T) {
	if NewRouter(mockClients{}, RouterConfig{Port: "8080"}) == nil {
		t.Fatalf("expected router")
	}
}
