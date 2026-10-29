/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package metrics

import (
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type pushgatewayStub struct {
	sync.Mutex
	pushes int
	method string
	body   string
}

func newPushgatewayStub(t *testing.T) (*pushgatewayStub, string) {
	stub := &pushgatewayStub{}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)

		stub.Lock()
		stub.pushes++
		stub.method = r.Method
		stub.body = string(body)
		stub.Unlock()

		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(srv.Close)

	return stub, srv.URL
}

func uptimeMetric(nodeId string, value float64) []MetricConfig {
	return []MetricConfig{{
		Name:   "test_node_uptime",
		Type:   "gauge",
		Labels: map[string]string{"node_id": nodeId, "site": "site-" + nodeId},
		Value:  value,
	}}
}

func resetCollectors(t *testing.T) {
	metricCollectors = make(map[string][]*Metrics)
	t.Cleanup(func() { metricCollectors = make(map[string][]*Metrics) })
}

func TestRemoveSeries(t *testing.T) {
	t.Run("DropsMatchingSeriesAndRepushes", func(t *testing.T) {
		resetCollectors(t)
		stub, host := newPushgatewayStub(t)

		PushMetrics(host, uptimeMetric("node-a", 10), "test")
		PushMetrics(host, uptimeMetric("node-b", 20), "test")
		require.Contains(t, stub.body, "node-a")

		removed, err := RemoveSeries(host, map[string]string{"node_id": "node-a"}, "test")

		require.NoError(t, err)
		assert.Equal(t, 1, removed)
		assert.Equal(t, 3, stub.pushes)
		assert.Equal(t, http.MethodPut, stub.method)
		assert.NotContains(t, stub.body, "node-a")
		assert.Contains(t, stub.body, "node-b")

		PushMetrics(host, uptimeMetric("node-b", 30), "test")
		assert.NotContains(t, stub.body, "node-a")
	})

	t.Run("LastSeriesLeavesEmptyGroup", func(t *testing.T) {
		resetCollectors(t)
		stub, host := newPushgatewayStub(t)

		PushMetrics(host, uptimeMetric("node-a", 10), "test")

		removed, err := RemoveSeries(host, map[string]string{"node_id": "node-a"}, "test")

		require.NoError(t, err)
		assert.Equal(t, 1, removed)
		assert.Equal(t, 2, stub.pushes)
		assert.Equal(t, http.MethodPut, stub.method)
		assert.Empty(t, stub.body)
	})

	t.Run("NoMatchSkipsPush", func(t *testing.T) {
		resetCollectors(t)
		stub, host := newPushgatewayStub(t)

		PushMetrics(host, uptimeMetric("node-a", 10), "test")

		removed, err := RemoveSeries(host, map[string]string{"node_id": "node-z"}, "test")

		require.NoError(t, err)
		assert.Equal(t, 0, removed)
		assert.Equal(t, 1, stub.pushes)
	})

	t.Run("NoLabels", func(t *testing.T) {
		resetCollectors(t)
		_, host := newPushgatewayStub(t)

		_, err := RemoveSeries(host, nil, "test")

		assert.Error(t, err)
	})
}
