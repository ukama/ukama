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
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/prometheus/prometheus/prompb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ukama/ukama/systems/common/rest/client/registry"
	"github.com/ukama/ukama/systems/metrics/sanitizer/pkg"
)

func TestMetricAndNodeId(t *testing.T) {
	t.Run("NodeIdLabel", func(t *testing.T) {
		metricName, nodeId := metricAndNodeId([]prompb.Label{
			{Name: "__name__", Value: pkg.RawComNodeUptimeSeconds},
			{Name: "node_id", Value: "uk-sa2341-tnode-a1-1234"},
			{Name: "env", Value: "dev"},
			{Name: "job", Value: "ukama-nodes"},
		})

		assert.Equal(t, pkg.RawComNodeUptimeSeconds, metricName)
		assert.Equal(t, "uk-sa2341-tnode-a1-1234", nodeId)
	})

	t.Run("LegacyNodeIdLabel", func(t *testing.T) {
		metricName, nodeId := metricAndNodeId([]prompb.Label{
			{Name: "__name__", Value: pkg.RawCtlNodeUptimeSeconds},
			{Name: "nodeid", Value: "uk-sa2341-anode-a1-1234"},
			{Name: "instance", Value: "10.0.0.1:8082"},
		})

		assert.Equal(t, pkg.RawCtlNodeUptimeSeconds, metricName)
		assert.Equal(t, "uk-sa2341-anode-a1-1234", nodeId)
	})

	t.Run("NoNodeLabel", func(t *testing.T) {
		metricName, nodeId := metricAndNodeId([]prompb.Label{
			{Name: "__name__", Value: pkg.RawComNodeUptimeSeconds},
		})

		assert.Equal(t, pkg.RawComNodeUptimeSeconds, metricName)
		assert.Empty(t, nodeId)
	})
}

func TestMetricCacheKey(t *testing.T) {
	nodeId := "uk-sa2341-tnode-a1-1234"

	assert.NotEqual(t, metricCacheKey(pkg.ComNodeUptime, nodeId),
		metricCacheKey(pkg.NodeActiveSubscribers, nodeId))
}

func TestSanitizedMetricsAreDeclared(t *testing.T) {
	for raw, sanitized := range pkg.SanitizedMetrics {
		found := false

		for _, m := range pkg.NodeMetrics {
			if m.Name == sanitized.Name {
				found = true

				assert.Equal(t, map[string]string{
					pkg.NodeIdLabel: "", pkg.SiteIdLabel: "", pkg.NetworkLabel: ""},
					m.Labels, "metric %q must declare exactly the sanitizer label set", m.Name)

				break
			}
		}

		assert.True(t, found, "no metric config declared for %q (raw metric %q)",
			sanitized.Name, raw)
	}
}

func TestLatestSampleValue(t *testing.T) {
	t.Run("SingleSample", func(t *testing.T) {
		assert.Equal(t, 42.0, latestSampleValue([]prompb.Sample{{Value: 42, Timestamp: 1}}))
	})

	t.Run("OutOfOrderSamples", func(t *testing.T) {
		assert.Equal(t, 300.0, latestSampleValue([]prompb.Sample{
			{Value: 100, Timestamp: 1},
			{Value: 300, Timestamp: 3},
			{Value: 200, Timestamp: 2},
		}))
	})
}

func TestRepublishedLabelNames(t *testing.T) {
	assert.Equal(t, "node_id", nodeLabel)
	assert.Equal(t, "site", siteLabel)
	assert.Equal(t, "network", networkLabel)

	assert.Equal(t, "nodeid", altNodeLabel)
}

func TestSanitizedMetricNames(t *testing.T) {
	active, ok := pkg.SanitizedMetrics["trx_lte_core_active_ue"]
	assert.True(t, ok)
	assert.Equal(t, "active_subscribers", active.Name)
	assert.True(t, active.SkipUnchanged)

	for raw, want := range map[string]string{
		"com_generic_system_uptime_seconds": "com_node_uptime",
		"ctl_generic_system_uptime_seconds": "ctl_node_uptime",
	} {
		uptime, ok := pkg.SanitizedMetrics[raw]
		assert.True(t, ok, "missing routing entry for %q", raw)
		assert.Equal(t, want, uptime.Name)
		assert.Equal(t, false, uptime.SkipUnchanged)
	}
}

func TestStaleNodes(t *testing.T) {
	nodeA := NodeMetaData{NodeId: "node-a", NetworkId: "net-1", SiteId: "site-1"}
	nodeB := NodeMetaData{NodeId: "node-b", NetworkId: "net-1", SiteId: "site-1"}

	t.Run("Deleted", func(t *testing.T) {
		stale := staleNodes(map[string]NodeMetaData{"node-a": nodeA, "node-b": nodeB},
			map[string]NodeMetaData{"node-b": nodeB})

		assert.Equal(t, []string{"node-a"}, stale)
	})

	t.Run("MovedSite", func(t *testing.T) {
		moved := nodeA
		moved.SiteId = "site-2"

		stale := staleNodes(map[string]NodeMetaData{"node-a": nodeA},
			map[string]NodeMetaData{"node-a": moved})

		assert.Equal(t, []string{"node-a"}, stale)
	})

	t.Run("Unchanged", func(t *testing.T) {
		stale := staleNodes(map[string]NodeMetaData{"node-a": nodeA},
			map[string]NodeMetaData{"node-a": nodeA, "node-b": nodeB})

		assert.Empty(t, stale)
	})

	t.Run("FirstSync", func(t *testing.T) {
		assert.Empty(t, staleNodes(nil, map[string]NodeMetaData{"node-a": nodeA}))
	})
}

func TestSyncNodeCacheRemovesSeriesOfDeletedNodes(t *testing.T) {
	var mu sync.Mutex
	nodes := []*registry.NodeInfo{
		{Id: "node-a", Site: registry.NodeSiteInfo{SiteId: "site-1", NetworkId: "net-1"}},
		{Id: "node-b", Site: registry.NodeSiteInfo{SiteId: "site-1", NetworkId: "net-1"}},
	}

	registrySrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()

		_ = json.NewEncoder(w).Encode(registry.Nodes{Nodes: nodes})
	}))
	defer registrySrv.Close()

	var lastPush string
	pushSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)

		mu.Lock()
		lastPush = string(body)
		mu.Unlock()
	}))
	defer pushSrv.Close()

	s, err := NewSanitizerServer(registrySrv.URL, pushSrv.URL, "testorg", "testorg-id", nil)
	require.NoError(t, err)

	for _, nodeId := range []string{"node-a", "node-b"} {
		s.updateNodeMetricCache(pkg.NodeActiveSubscribers, nodeId, 3)
		s.pushSanitizedNodeMetric(NodeMetricMetaData{
			MetricName: pkg.ComNodeUptime,
			NodeId:     nodeId,
			Value:      100,
			Labels:     map[string]string{nodeLabel: nodeId, siteLabel: "site-1", networkLabel: "net-1"},
		})
	}

	mu.Lock()
	require.Contains(t, lastPush, "node-a")
	nodes = nodes[1:]
	mu.Unlock()

	require.NoError(t, s.syncNodeCache())

	mu.Lock()
	assert.NotContains(t, lastPush, "node-a")
	assert.Contains(t, lastPush, "node-b")
	mu.Unlock()

	_, ok := s.getNodeMetricFromCache(pkg.NodeActiveSubscribers, "node-a")
	assert.False(t, ok)

	_, ok = s.getNodeMetricFromCache(pkg.NodeActiveSubscribers, "node-b")
	assert.True(t, ok)
}
