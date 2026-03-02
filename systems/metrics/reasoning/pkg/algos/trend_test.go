/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package algos

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCalculateTrend(t *testing.T) {
	t.Run("Increasing", func(t *testing.T) {
		// delta = 20, noise = 2, sensitivity = 1 -> threshold = 2, |20| > 2
		aggNow := AggregationStats{AggregatedValue: 80, NoiseEstimate: 2}
		aggPrev := AggregationStats{AggregatedValue: 60}
		trend, err := CalculateTrend(aggNow, aggPrev, 1.0)
		require.NoError(t, err)
		assert.Equal(t, "increasing", trend)
	})

	t.Run("Decreasing", func(t *testing.T) {
		aggNow := AggregationStats{AggregatedValue: 60, NoiseEstimate: 2}
		aggPrev := AggregationStats{AggregatedValue: 80}
		trend, err := CalculateTrend(aggNow, aggPrev, 1.0)
		require.NoError(t, err)
		assert.Equal(t, "decreasing", trend)
	})

	t.Run("Stable_DeltaWithinThreshold", func(t *testing.T) {
		// delta = 1, noise = 5, sensitivity = 1 -> threshold = 5, |1| < 5
		aggNow := AggregationStats{AggregatedValue: 51, NoiseEstimate: 5}
		aggPrev := AggregationStats{AggregatedValue: 50}
		trend, err := CalculateTrend(aggNow, aggPrev, 1.0)
		require.NoError(t, err)
		assert.Equal(t, "stable", trend)
	})

	t.Run("Stable_HighSensitivity", func(t *testing.T) {
		// delta = 3, noise = 2, sensitivity = 2 -> threshold = 4, |3| < 4
		aggNow := AggregationStats{AggregatedValue: 53, NoiseEstimate: 2}
		aggPrev := AggregationStats{AggregatedValue: 50}
		trend, err := CalculateTrend(aggNow, aggPrev, 2.0)
		require.NoError(t, err)
		assert.Equal(t, "stable", trend)
	})

	t.Run("PrevAggNaN_ReturnsStable", func(t *testing.T) {
		aggNow := AggregationStats{AggregatedValue: 80, NoiseEstimate: 2}
		aggPrev := AggregationStats{AggregatedValue: math.NaN()}
		trend, err := CalculateTrend(aggNow, aggPrev, 1.0)
		require.NoError(t, err)
		assert.Equal(t, "stable", trend)
	})

	t.Run("PrevAggInf_ReturnsStable", func(t *testing.T) {
		aggNow := AggregationStats{AggregatedValue: 80, NoiseEstimate: 2}
		aggPrev := AggregationStats{AggregatedValue: math.Inf(1)}
		trend, err := CalculateTrend(aggNow, aggPrev, 1.0)
		require.NoError(t, err)
		assert.Equal(t, "stable", trend)
	})

	t.Run("NoiseNaN_ReturnsStable", func(t *testing.T) {
		aggNow := AggregationStats{AggregatedValue: 80, NoiseEstimate: math.NaN()}
		aggPrev := AggregationStats{AggregatedValue: 60}
		trend, err := CalculateTrend(aggNow, aggPrev, 1.0)
		require.NoError(t, err)
		assert.Equal(t, "stable", trend)
	})

	t.Run("NoiseInf_ReturnsStable", func(t *testing.T) {
		aggNow := AggregationStats{AggregatedValue: 80, NoiseEstimate: math.Inf(1)}
		aggPrev := AggregationStats{AggregatedValue: 60}
		trend, err := CalculateTrend(aggNow, aggPrev, 1.0)
		require.NoError(t, err)
		assert.Equal(t, "stable", trend)
	})

	t.Run("Volatile_HighNoiseRelativeToSignal", func(t *testing.T) {
		// aggNow = 100, noise = 60 -> noise/|aggNow| = 0.6 > 0.5
		aggNow := AggregationStats{AggregatedValue: 100, NoiseEstimate: 60}
		aggPrev := AggregationStats{AggregatedValue: 50}
		trend, err := CalculateTrend(aggNow, aggPrev, 1.0)
		require.NoError(t, err)
		assert.Equal(t, "volatile", trend)
	})

	t.Run("NotVolatile_LowNoiseRelativeToSignal", func(t *testing.T) {
		// aggNow = 100, noise = 30 -> noise/|aggNow| = 0.3 < 0.5 (not volatile)
		// delta = 50, threshold = 30 -> increasing
		aggNow := AggregationStats{AggregatedValue: 100, NoiseEstimate: 30}
		aggPrev := AggregationStats{AggregatedValue: 50}
		trend, err := CalculateTrend(aggNow, aggPrev, 1.0)
		require.NoError(t, err)
		assert.Equal(t, "increasing", trend)
	})

	t.Run("SensitivityZero_DefaultsToOne", func(t *testing.T) {
		// With sensitivity 0, should use 1.0 -> same as increasing test
		aggNow := AggregationStats{AggregatedValue: 80, NoiseEstimate: 2}
		aggPrev := AggregationStats{AggregatedValue: 60}
		trend, err := CalculateTrend(aggNow, aggPrev, 0)
		require.NoError(t, err)
		assert.Equal(t, "increasing", trend)
	})

	t.Run("SensitivityNegative_DefaultsToOne", func(t *testing.T) {
		aggNow := AggregationStats{AggregatedValue: 80, NoiseEstimate: 2}
		aggPrev := AggregationStats{AggregatedValue: 60}
		trend, err := CalculateTrend(aggNow, aggPrev, -1)
		require.NoError(t, err)
		assert.Equal(t, "increasing", trend)
	})

	t.Run("DeltaZero_ReturnsStable", func(t *testing.T) {
		aggNow := AggregationStats{AggregatedValue: 50, NoiseEstimate: 2}
		aggPrev := AggregationStats{AggregatedValue: 50}
		trend, err := CalculateTrend(aggNow, aggPrev, 1.0)
		require.NoError(t, err)
		assert.Equal(t, "stable", trend)
	})

	t.Run("Volatile_SmallAggNowValue_SkipsVolatileCheck", func(t *testing.T) {
		// aggNow near zero: |aggNow| <= noiseEpsilon skips volatile check
		// Use aggNow = 0, noise huge - would be volatile if we checked, but |0| <= epsilon
		// Actually: math.Abs(0) > 1e-10 is false, so we skip volatile and go to classifyTrend
		aggNow := AggregationStats{AggregatedValue: 0, NoiseEstimate: 100}
		aggPrev := AggregationStats{AggregatedValue: -10}
		trend, err := CalculateTrend(aggNow, aggPrev, 1.0)
		require.NoError(t, err)
		// delta = 10, threshold = 100, |10| < 100 -> stable
		assert.Equal(t, "stable", trend)
	})
}
