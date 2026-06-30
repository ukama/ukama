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
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestAllRoutes fires one request at every route so each handler body executes.
// Stub clients return success, so a non-5xx response means the handler ran.
func TestAllRoutes(t *testing.T) {
	const base = "/v1/analytics"

	gets := []string{
		base + "/business/sales/overview",
		base + "/business/sales/packages",
		base + "/business/packages",
		base + "/business/billing",
		base + "/business/sites",
		base + "/business/sites/site-a",
		base + "/business/inventory",
		base + "/customers/list",
		base + "/customers/search?q=john",
		base + "/customers/sims",
		base + "/customers/sim-pool",
		base + "/customers/cust-1/support",
		base + "/network/topology",
		base + "/network/sites",
		base + "/network/sites/site-a",
		base + "/network/nodes",
		base + "/network/nodes/node-1",
		base + "/network/node-pool",
		base + "/network/radio",
		base + "/network/backhaul",
		base + "/network/power",
		base + "/network/alarms",
		base + "/network/metrics?family=all",
		base + "/network/events",
		base + "/network/support/search?q=site-a",
		base + "/collector/state",
	}

	router := newTestRouter(defaultClients())

	for _, path := range gets {
		t.Run("GET "+path, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, path, nil)
			router.ServeHTTP(w, req)
			assert.Less(t, w.Code, http.StatusInternalServerError, "handler should not 5xx: %s", path)
		})
	}

	posts := []struct {
		path string
		body string
	}{
		{base + "/collector/rollups/rebuild", `{"family":"all"}`},
		{base + "/collector/seed-demo", `{"sites":1,"nodes":1,"customers":1,"days":1}`},
	}

	for _, p := range posts {
		t.Run("POST "+p.path, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, p.path, strings.NewReader(p.body))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)
			assert.Less(t, w.Code, http.StatusInternalServerError, "handler should not 5xx: %s", p.path)
		})
	}
}
