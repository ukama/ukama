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

	"github.com/ukama/ukama/systems/metrics/reasoning/pkg/metric"
)

// ProjectionStats holds the result of projecting when a metric will cross a threshold.
// It is not a precise prediction; it is an early warning with a rough time-to-threshold.
type ProjectionStats struct {
	Type   string  `json:"type"`    // e.g. "to_warning", "to_critical"
	EtaSec float64 `json:"eta_sec"` // estimated seconds until crossing
}

// ProjectCrossingTime answers: "if this continues, will we cross a boundary soon?"
// Uses the delta between windows as a slope.
//
// Returns nil when:
//   - value_now or value_prev is unknown (NaN/Inf)
//   - already at or beyond the target threshold
//   - slope does not move toward the threshold (e.g. improving when we care about worsening)
//
// Direction "range" is not yet supported; returns nil.
func ProjectCrossingTime(valueNow, valuePrev, windowSec float64, thresholds StateThresholds, direction string) ProjectionStats {
	if math.IsNaN(valueNow) || math.IsInf(valueNow, 0) ||
		math.IsNaN(valuePrev) || math.IsInf(valuePrev, 0) {
		return ProjectionStats{}
	}
	if windowSec <= 0 {
		return ProjectionStats{}
	}

	slopePerSec := (valueNow - valuePrev) / windowSec

	switch direction {
	case DirectionHigherIsWorse, "":
		target := thresholds.Warning
		if valueNow >= target {
			return ProjectionStats{}
		}
		if slopePerSec <= 0 {
			return ProjectionStats{}
		}
		seconds := (target - valueNow) / slopePerSec
		return ProjectionStats{Type: "to_warning", EtaSec: seconds}

	case DirectionLowerIsWorse:
		target := thresholds.Warning
		if valueNow <= target {
			return ProjectionStats{}
		}
		if slopePerSec >= 0 {
			return ProjectionStats{}
		}
		seconds := (valueNow - target) / math.Abs(slopePerSec)
		return ProjectionStats{Type: "to_warning", EtaSec: seconds}

	case DirectionRange:
		// Range direction not yet supported; would need to project toward LowWarning or HighWarning
		return ProjectionStats{}

	default:
		return ProjectionStats{}
	}
}

// ProjectMetricResults returns the input results unchanged.
// Projection analysis is done via ProjectCrossingTime using aggregated values.
func ProjectMetricResults(results []metric.FilteredPrometheusResult) ([]metric.FilteredPrometheusResult, error) {
	return results, nil
}
