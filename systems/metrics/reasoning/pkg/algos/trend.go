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
)

const noiseEpsilon = 1e-10 // avoid division by zero

// classifyTrend returns "stable", "increasing", or "decreasing" based on delta vs threshold.
func classifyTrend(delta, noise, sensitivity float64) string {
	threshold := sensitivity * noise
	if math.Abs(delta) < threshold {
		return "stable"
	}
	if delta > 0 {
		return "increasing"
	}
	return "decreasing"
}

// CalculateTrend measures direction of change by comparing the current window to the previous window.
// Uses pre-computed AggregationStats from AggregateMetricAlgo
// Trend is only "real" if delta exceeds normal noise.
//
// Parameters:
//   - aggNow, aggPrev: aggregation results from AggregateMetricAlgo for now and previous windows
//   - sensitivity: multiplier for noise threshold; higher = harder to call a trend (more stable results)
//
// Returns: "stable", "increasing", "decreasing", or "volatile" (when window has high variation)
func CalculateTrend(aggNow, aggPrev AggregationStats, sensitivity float64) (string, error) {
	if sensitivity <= 0 {
		sensitivity = 1.0
	}

	delta := aggNow.AggregatedValue - aggPrev.AggregatedValue
	noise := aggNow.NoiseEstimate

	if math.IsNaN(noise) || math.IsInf(noise, 1) {
		return "stable", nil // too few samples for reliable trend
	}
	noise += noiseEpsilon

	// Volatile: when noise is very high relative to signal, trend is unreliable
	if math.Abs(aggNow.AggregatedValue) > noiseEpsilon && noise/math.Abs(aggNow.AggregatedValue) > 0.5 {
		return "volatile", nil
	}

	return classifyTrend(delta, noise, sensitivity), nil
}
