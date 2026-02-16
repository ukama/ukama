/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package utils

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/ukama/ukama/systems/metrics/reasoning/pkg"
	"github.com/ukama/ukama/systems/metrics/reasoning/pkg/store"
)

func TestRoundToDecimalPoints(t *testing.T) {
	t.Run("NormalValue", func(t *testing.T) {
		assert.Equal(t, 3.14, RoundToDecimalPoints(3.14159, 2))
		assert.Equal(t, 3.142, RoundToDecimalPoints(3.14159, 3))
	})

	t.Run("NaN_Preserved", func(t *testing.T) {
		result := RoundToDecimalPoints(math.NaN(), 2)
		assert.True(t, math.IsNaN(result))
	})

	t.Run("Inf_Preserved", func(t *testing.T) {
		result := RoundToDecimalPoints(math.Inf(1), 2)
		assert.True(t, math.IsInf(result, 1))
	})

	t.Run("NegativeDecimalPoints_ReturnsUnchanged", func(t *testing.T) {
		assert.Equal(t, 3.14159, RoundToDecimalPoints(3.14159, -1))
	})

	t.Run("ZeroDecimals", func(t *testing.T) {
		assert.Equal(t, 42.0, RoundToDecimalPoints(42.3, 0))
	})
}

func TestGetAlgoStatsStoreKey(t *testing.T) {
	key := GetAlgoStatsStoreKey("node-123", "cpu")
	assert.Equal(t, "node-123/cpu/algo_stats", key)
}

func TestSortNodeIds(t *testing.T) {
	t.Run("ValidTowerNode", func(t *testing.T) {
		nodes, err := SortNodeIds("UK-SA2156-TNODE-A1-XXXX")
		require.NoError(t, err)
		assert.Equal(t, "uk-sa2156-tnode-a1-xxxx", nodes.TNode)
		assert.Equal(t, "uk-sa2156-anode-a1-xxxx", nodes.ANode)
	})

	t.Run("InvalidNodeID", func(t *testing.T) {
		_, err := SortNodeIds("invalid")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "validate node id")
	})

	t.Run("NonTowerNode_ReturnsError", func(t *testing.T) {
		_, err := SortNodeIds("UK-SA2156-HNODE-A1-XXXX")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "expected tower node")
	})
}

func TestGetStartEndFromStore(t *testing.T) {
	s := store.NewInMemoryStore()

	t.Run("FirstCall_NoStoredValue_UsesCurrentWindow", func(t *testing.T) {
		start, end, err := GetStartEndFromStore(s, "node-1", 300)
		require.NoError(t, err)
		require.NotEmpty(t, start)
		require.NotEmpty(t, end)
	})

	t.Run("SecondCall_AdvancesWindow", func(t *testing.T) {
		_ = s.Put("node-2/start_end", "1000:1300")
		start, end, err := GetStartEndFromStore(s, "node-2", 300)
		require.NoError(t, err)
		assert.Equal(t, "1300", start)
		assert.Equal(t, "1600", end)
	})

	t.Run("InvalidStoredValue_ReturnsError", func(t *testing.T) {
		_ = s.Put("node-3/start_end", "invalid")
		_, _, err := GetStartEndFromStore(s, "node-3", 300)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid stored value")
	})

	t.Run("InvalidEndTimestamp_ReturnsError", func(t *testing.T) {
		_ = s.Put("node-4/start_end", "1000:notanumber")
		_, _, err := GetStartEndFromStore(s, "node-4", 300)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid end timestamp")
	})
}

func TestValidateMetricKey(t *testing.T) {
	metricsCfg := pkg.Metrics{
		Metrics: []pkg.Metric{
			{Key: "cpu", MetricKey: "cpu_usage_percent"},
			{Key: "memory", MetricKey: "memory_usage_percent"},
		},
	}

	t.Run("ValidKey", func(t *testing.T) {
		metricKey, err := ValidateMetricKey("cpu", metricsCfg, "tnode")
		require.NoError(t, err)
		assert.Equal(t, "cpu_usage_percent", metricKey)
	})

	t.Run("InvalidKey", func(t *testing.T) {
		_, err := ValidateMetricKey("disk", metricsCfg, "tnode")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not valid")
	})
}
