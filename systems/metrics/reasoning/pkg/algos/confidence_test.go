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
	"github.com/ukama/ukama/systems/metrics/reasoning/pkg/metric"
)

func confidenceMakeSamples(values [][]interface{}) []metric.FilteredPrometheusResult {
	return []metric.FilteredPrometheusResult{
		{
			Metric: metric.FilteredMetric{NodeID: "node-1", Metric: "cpu"},
			Values: values,
		},
	}
}

func TestCalculateConfidence(t *testing.T) {
	t.Run("EmptyResults", func(t *testing.T) {
		aggNow := AggregationStats{AggregatedValue: 50, SampleCount: 5, NoiseEstimate: 2}
		aggPrev := AggregationStats{AggregatedValue: 45}
		conf, err := CalculateConfidence(nil, aggNow, aggPrev, "healthy", 10)
		require.NoError(t, err)
		assert.Equal(t, 0.0, conf)
	})

	t.Run("FullCoverageStrongSignalConsistentTrend", func(t *testing.T) {
		// Perfectly increasing: 10, 20, 30, 40, 50 -> strong positive correlation
		samples := confidenceMakeSamples([][]interface{}{
			{1000.0, 10.0}, {2000.0, 20.0}, {3000.0, 30.0}, {4000.0, 40.0}, {5000.0, 50.0},
		})
		aggNow := AggregationStats{
			AggregatedValue: 50,
			SampleCount:     5,
			NoiseEstimate:   2,
		}
		aggPrev := AggregationStats{AggregatedValue: 10} // delta = 40, strong signal
		conf, err := CalculateConfidence(samples, aggNow, aggPrev, "increasing", 5)
		require.NoError(t, err)
		assert.Greater(t, conf, 0.8)
		assert.LessOrEqual(t, conf, 1.0)
	})

	t.Run("UnknownStateReducesConfidence", func(t *testing.T) {
		samples := confidenceMakeSamples([][]interface{}{
			{1000.0, 10.0}, {2000.0, 20.0}, {3000.0, 30.0}, {4000.0, 40.0}, {5000.0, 50.0},
		})
		aggNow := AggregationStats{AggregatedValue: 50, SampleCount: 5, NoiseEstimate: 2}
		aggPrev := AggregationStats{AggregatedValue: 10}
		confKnown, err := CalculateConfidence(samples, aggNow, aggPrev, "healthy", 5)
		require.NoError(t, err)
		confUnknown, err := CalculateConfidence(samples, aggNow, aggPrev, "unknown", 5)
		require.NoError(t, err)
		assert.Less(t, confUnknown, confKnown)
		assert.InDelta(t, confUnknown, confKnown*0.4, 0.01)
	})

	t.Run("NoPreviousDataZeroSignal", func(t *testing.T) {
		samples := confidenceMakeSamples([][]interface{}{
			{1000.0, 10.0}, {2000.0, 20.0}, {3000.0, 30.0}, {4000.0, 40.0}, {5000.0, 50.0},
		})
		aggNow := AggregationStats{AggregatedValue: 30, SampleCount: 5, NoiseEstimate: 5}
		aggPrev := AggregationStats{AggregatedValue: math.NaN(), NoiseEstimate: math.NaN()}
		conf, err := CalculateConfidence(samples, aggNow, aggPrev, "healthy", 5)
		require.NoError(t, err)
		// sig = 0 due to NaN delta; base = 0.5*1 + 0.3*0 + 0.2*con
		assert.Greater(t, conf, 0.0)
		assert.Less(t, conf, 0.8) // no signal contribution
	})

	t.Run("LowCoverage", func(t *testing.T) {
		// 2 samples -> trendConsistency returns 0 (need 3+), coverage: 2/10 = 0.2
		samples := confidenceMakeSamples([][]interface{}{
			{1000.0, 10.0}, {2000.0, 20.0},
		})
		aggNow := AggregationStats{AggregatedValue: 15, SampleCount: 2, NoiseEstimate: 5}
		aggPrev := AggregationStats{AggregatedValue: 10}
		conf, err := CalculateConfidence(samples, aggNow, aggPrev, "healthy", 10)
		require.NoError(t, err)
		assert.Less(t, conf, 0.5)
	})

	t.Run("ConstantValuesZeroConsistency", func(t *testing.T) {
		// Constant values -> varV = 0 -> trendConsistency = 0
		samples := confidenceMakeSamples([][]interface{}{
			{1000.0, 50.0}, {2000.0, 50.0}, {3000.0, 50.0}, {4000.0, 50.0}, {5000.0, 50.0},
		})
		aggNow := AggregationStats{AggregatedValue: 50, SampleCount: 5, NoiseEstimate: 2}
		aggPrev := AggregationStats{AggregatedValue: 45}
		conf, err := CalculateConfidence(samples, aggNow, aggPrev, "healthy", 5)
		require.NoError(t, err)
		assert.Greater(t, conf, 0.0)
		assert.Less(t, conf, 1.0)
	})

	t.Run("ExpectedSamplesZeroUsesMinSamples", func(t *testing.T) {
		samples := confidenceMakeSamples([][]interface{}{
			{1000.0, 10.0}, {2000.0, 20.0}, {3000.0, 30.0}, {4000.0, 40.0}, {5000.0, 50.0},
		})
		aggNow := AggregationStats{AggregatedValue: 30, SampleCount: 5, NoiseEstimate: 5}
		aggPrev := AggregationStats{AggregatedValue: 25}
		conf, err := CalculateConfidence(samples, aggNow, aggPrev, "healthy", 0)
		require.NoError(t, err)
		// exp = minSamplesForCoverage (5), so coverage = 5/5 = 1
		assert.Greater(t, conf, 0.0)
	})

	t.Run("ResultClampedToZeroOne", func(t *testing.T) {
		samples := confidenceMakeSamples([][]interface{}{
			{1000.0, 10.0}, {2000.0, 20.0}, {3000.0, 30.0}, {4000.0, 40.0}, {5000.0, 50.0},
		})
		aggNow := AggregationStats{AggregatedValue: 100, SampleCount: 10, NoiseEstimate: 1}
		aggPrev := AggregationStats{AggregatedValue: 0}
		conf, err := CalculateConfidence(samples, aggNow, aggPrev, "healthy", 5)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, conf, 0.0)
		assert.LessOrEqual(t, conf, 1.0)
	})

	t.Run("NaNNoiseEstimateUsesEpsilon", func(t *testing.T) {
		samples := confidenceMakeSamples([][]interface{}{
			{1000.0, 10.0}, {2000.0, 20.0}, {3000.0, 30.0},
		})
		aggNow := AggregationStats{AggregatedValue: 25, SampleCount: 3, NoiseEstimate: math.NaN()}
		aggPrev := AggregationStats{AggregatedValue: 20}
		conf, err := CalculateConfidence(samples, aggNow, aggPrev, "healthy", 5)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, conf, 0.0)
		assert.LessOrEqual(t, conf, 1.0)
	})
}
