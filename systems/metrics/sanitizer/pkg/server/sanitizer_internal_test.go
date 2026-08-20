/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package server

import (
	"testing"

	"github.com/prometheus/prometheus/prompb"
	"github.com/stretchr/testify/assert"

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
