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

const (
	minSamplesForConsistency = 3
	minSamplesForCoverage    = 5
)

// trendConsistency computes how consistently values move over time using
// simple correlation between time indices and values. Returns [0,1].
// 1 means strong consistent direction, 0 means no pattern or too few samples.
func trendConsistency(values []float64) float64 {
	if len(values) < minSamplesForConsistency {
		return 0.0
	}
	n := float64(len(values))
	meanT := (n - 1) / 2
	meanV := 0.0
	for _, v := range values {
		meanV += v
	}
	meanV /= n

	var cov, varT, varV float64
	for i, v := range values {
		ti := float64(i)
		dt := ti - meanT
		dv := v - meanV
		cov += dt * dv
		varT += dt * dt
		varV += dv * dv
	}
	cov /= n
	varT /= n
	varV /= n

	if varT < noiseEpsilon || varV < noiseEpsilon {
		return 0.0
	}
	corr := cov / (math.Sqrt(varT) * math.Sqrt(varV))
	return clamp(math.Abs(corr), 0, 1)
}

func clamp(x, lo, hi float64) float64 {
	if x < lo {
		return lo
	}
	if x > hi {
		return hi
	}
	return x
}

// CalculateConfidence answers "how sure is this evaluation?" based on evidence quality.
//
// Inputs:
//   - coverage_pct: sample coverage (did we see enough samples in the window?)
//   - signal_strength: abs(delta)/noise (is the delta stronger than normal variation?)
//   - consistency_score: trend consistency (did values move consistently, or bounce?)
//   - state: if "unknown", confidence is reduced
//
// Uses outputs from AggregateMetricAlgo (NoiseEstimate, SampleCount), CalculateTrend (delta),
// and CalculateState.
func CalculateConfidence(results []metric.FilteredPrometheusResult, aggNow, aggPrev AggregationStats, state string, expectedSamples int) (float64, error) {
	values := extractNumericValues(results)
	if len(values) == 0 {
		return 0.0, nil
	}

	// coverage_pct: [0,1] - did we see enough samples?
	exp := expectedSamples
	if exp <= 0 {
		exp = minSamplesForCoverage
	}
	coveragePct := clamp(aggNow.SampleCount/float64(exp), 0, 1)

	// signal_strength: abs(delta)/noise, mapped to [0,1] via sig/(sig+1)
	delta := aggNow.AggregatedValue - aggPrev.AggregatedValue
	noise := aggNow.NoiseEstimate
	if math.IsNaN(noise) || math.IsInf(noise, 1) || noise < 0 {
		noise = noiseEpsilon
	}
	noise += noiseEpsilon
	signalStrength := math.Abs(delta) / noise
	sig := clamp(signalStrength/(signalStrength+1), 0, 1)

	// consistency_score: [0,1] from trend consistency
	con := trendConsistency(values)

	// base = 0.5*cov + 0.3*sig + 0.2*con
	base := 0.5*coveragePct + 0.3*sig + 0.2*con

	if state == "unknown" {
		base *= 0.4
	}

	return clamp(base, 0, 1), nil
}
