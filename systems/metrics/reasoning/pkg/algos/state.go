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

	"github.com/ukama/ukama/systems/metrics/reasoning/pkg"
)

// State direction constants for threshold classification
const (
	DirectionHigherIsWorse = "higher_is_worse" // temperature, memory %, latency
	DirectionLowerIsWorse = "lower_is_worse"   // success rate, signal strength
	DirectionRange        = "range"             // voltage within band
)

// StateThresholds holds threshold values for state classification.
// For higher_is_worse / lower_is_worse: use Warning and Critical.
// For range: use LowWarning, HighWarning, LowCritical, HighCritical.
type StateThresholds struct {
	Warning  float64
	Critical float64
	// Range direction
	LowWarning   float64
	HighWarning  float64
	LowCritical  float64
	HighCritical float64
}

func classifyHigherIsWorse(value float64, t StateThresholds) string {
	if value >= t.Critical {
		return "critical"
	}
	if value >= t.Warning {
		return "warning"
	}
	return "healthy"
}

func classifyLowerIsWorse(value float64, t StateThresholds) string {
	if value <= t.Critical {
		return "critical"
	}
	if value <= t.Warning {
		return "warning"
	}
	return "healthy"
}

func classifyRange(value float64, t StateThresholds) string {
	if value < t.LowCritical || value > t.HighCritical {
		return "critical"
	}
	if value < t.LowWarning || value > t.HighWarning {
		return "warning"
	}
	return "healthy"
}

func classifyState(value float64, thresholds StateThresholds, direction string) string {
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return "unknown"
	}
	switch direction {
	case DirectionHigherIsWorse:
		return classifyHigherIsWorse(value, thresholds)
	case DirectionLowerIsWorse:
		return classifyLowerIsWorse(value, thresholds)
	case DirectionRange:
		return classifyRange(value, thresholds)
	default:
		return "unknown"
	}
}

func BuildStateThresholds(m pkg.Metric) StateThresholds {
	st := StateThresholds{}
	switch m.StateDirection {
	case DirectionLowerIsWorse:
		st.Warning = m.Thresholds.Medium
		st.Critical = m.Thresholds.Min
	case DirectionRange:
		st.LowWarning = m.Thresholds.LowWarning
		st.HighWarning = m.Thresholds.HighWarning
		st.LowCritical = m.Thresholds.LowCritical
		st.HighCritical = m.Thresholds.HighCritical
	case DirectionHigherIsWorse, "":
		// default higher_is_worse for gauge metrics (CPU, memory)
		st.Warning = m.Thresholds.Medium
		st.Critical = m.Thresholds.Max
	}
	return st
}

// CalculateState maps the aggregated value to healthy/warning/critical using thresholds from policy.
// Uses pre-computed AggregatedValue from AggregationStats.
func CalculateState(aggregatedValue float64, thresholds StateThresholds, direction string) (string, error) {
	if direction == "" {
		direction = DirectionHigherIsWorse // default for gauge metrics
	}
	return classifyState(aggregatedValue, thresholds, direction), nil
}
