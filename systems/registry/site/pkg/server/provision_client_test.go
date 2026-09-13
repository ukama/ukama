/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 * Copyright (c) 2026-present, Ukama Inc.
 */
package server

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProvisionClientExistingEndpoints(t *testing.T) {
	var calls []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls = append(calls, r.Method+" "+r.URL.Path)
		if r.Method == http.MethodGet {
			require.Equal(t, "attempt-1", r.URL.Query().Get("request_id"))
			_, _ = w.Write([]byte(`{"configuration":{"request_id":"attempt-1","completed":true,"cancelled":true,"cleared":true}}`))
			return
		}
		var body map[string]string
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		require.Equal(t, map[string]string{"request_id": "attempt-1", "site_id": "site-1", "network_id": "network-1"}, body)
		_, _ = w.Write([]byte(`{"status":"DISPATCHED"}`))
	}))
	defer server.Close()
	client := &nodeProvisionClient{url: server.URL, http: server.Client()}
	node := provisionNode{NodeID: "node-1", RequestID: "attempt-1", SiteID: "site-1", NetworkID: "network-1", Action: "configure"}
	result, err := client.Reconcile(context.Background(), node)
	require.NoError(t, err)
	require.False(t, result.Completed, "dispatch is not configuration completion")
	node.Action = "status"
	result, err = client.Reconcile(context.Background(), node)
	require.NoError(t, err)
	require.True(t, result.Completed)
	node.Action = "cancel"
	result, err = client.Reconcile(context.Background(), node)
	require.NoError(t, err)
	require.True(t, result.Cleared)
	require.Equal(t, []string{"POST /v1/controller/nodes/node-1/config", "GET /v1/state/node-1/latest", "DELETE /v1/controller/nodes/node-1/config", "GET /v1/state/node-1/latest"}, calls)
}

func TestProvisionClientRejectsUncorrelatedStatus(t *testing.T) {
	for _, body := range []string{`{"state":{}}`, `{"configuration":{"request_id":"old","completed":true}}`} {
		t.Run(body, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { _, _ = w.Write([]byte(body)) }))
			defer server.Close()
			client := &nodeProvisionClient{url: server.URL, http: server.Client()}
			_, err := client.Reconcile(context.Background(), provisionNode{NodeID: "node-1", RequestID: "current", Action: "status"})
			require.Error(t, err)
		})
	}
}
