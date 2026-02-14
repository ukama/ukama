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
	"github.com/ukama/ukama/systems/metrics/reasoning/pkg/metric"
)

// EvaluationPolicy holds the policy for metric evaluation.
type EvaluationPolicy struct {
	WindowSec                 float64        // time window in seconds
	ExpectedScrapeIntervalSec float64        // expected interval between samples
	MinCoveragePct           float64        // minimum coverage [0,1] to evaluate
	Aggregation              string         // "mean", "median", "last", "p95", "sum"
	TrendSensitivity         float64        // multiplier for trend noise threshold
	EnableVolatility         bool           // override trend with "volatile" when noise is high
	VolatilityRatio          float64        // noise/|value| threshold for volatile (e.g. 0.5)
	Thresholds               StateThresholds
	Direction                string // "higher_is_worse", "lower_is_worse", "range"
}

// MetricEvaluation is the result of evaluating a metric.
type MetricEvaluation struct {
	MetricID    string          `json:"metric_id"`
	EvaluatedAt int64           `json:"evaluated_at"`
	AggNow      float64         `json:"agg_now"`
	StatsNow    AggregationStats `json:"stats_now"`
	State       string          `json:"state"`
	Trend       string          `json:"trend"`
	Conclusion  string          `json:"conclusion"`
	Confidence  float64         `json:"confidence"`
	Projection  ProjectionStats `json:"projection"`
}

// CombineStateAndTrend defines the "meaning vocabulary" for conclusions.
func CombineStateAndTrend(state, trend string) string {
	switch state {
	case "critical":
		if trend == "decreasing" {
			return "recovering"
		}
		return "degrading"
	case "warning":
		if trend == "increasing" {
			return "risk_rising"
		}
		if trend == "decreasing" {
			return "recovering"
		}
		return "persistent"
	case "healthy":
		if trend == "increasing" {
			return "risk_rising"
		}
		return "ok"
	}
	return "unknown"
}

// maybeOverrideWithVolatility returns "volatile" when noise is high relative to signal.
func maybeOverrideWithVolatility(noise, aggNow float64, volatilityRatio float64) string {
	if math.Abs(aggNow) < noiseEpsilon {
		return "volatile"
	}
	if noise/math.Abs(aggNow) > volatilityRatio {
		return "volatile"
	}
	return ""
}

func estimateCoverage(samples []metric.FilteredPrometheusResult, windowSec, expectedIntervalSec float64) float64 {
	if expectedIntervalSec <= 0 {
		expectedIntervalSec = 1
	}
	expectedCount := windowSec / expectedIntervalSec
	if expectedCount < 1 {
		expectedCount = 1
	}
	actualCount := float64(countSamples(samples))
	return clamp(actualCount/expectedCount, 0, 1)
}

func countSamples(samples []metric.FilteredPrometheusResult) int {
	n := 0
	for _, s := range samples {
		for _, pair := range s.Values {
			if len(pair) >= 2 {
				n++
			}
		}
	}
	return n
}

// filterSamplesByTime returns only [ts, value] pairs where ts is in [startSec, endSec].
func filterSamplesByTime(samples []metric.FilteredPrometheusResult, startSec, endSec float64) []metric.FilteredPrometheusResult {
	result := make([]metric.FilteredPrometheusResult, 0, len(samples))
	for _, s := range samples {
		var filtered [][]interface{}
		for _, pair := range s.Values {
			if len(pair) < 2 {
				continue
			}
			ts, ok := extractTimestamp(pair)
			if !ok || ts < startSec || ts > endSec {
				continue
			}
			filtered = append(filtered, pair)
		}
		if len(filtered) > 0 {
			result = append(result, metric.FilteredPrometheusResult{
				Metric: s.Metric,
				Values: filtered,
			})
		}
	}
	return result
}

// computeConfidenceFromComponents computes confidence from coverage, signal strength, consistency, and state.
func computeConfidenceFromComponents(coveragePct, signalStrength, consistencyScore float64, state string) float64 {
	cov := clamp(coveragePct, 0, 1)
	sig := clamp(signalStrength/(signalStrength+1), 0, 1)
	con := clamp(consistencyScore, 0, 1)
	base := 0.5*cov + 0.3*sig + 0.2*con
	if state == "unknown" {
		base *= 0.4
	}
	return clamp(base, 0, 1)
}

// BuildEvaluationPolicy creates an EvaluationPolicy from pkg.Metric and window params.
func BuildEvaluationPolicy(m pkg.Metric, windowSec, expectedScrapeIntervalSec float64, minCoveragePct float64) EvaluationPolicy {
	policy := EvaluationPolicy{
		WindowSec:                 windowSec,
		ExpectedScrapeIntervalSec:  expectedScrapeIntervalSec,
		MinCoveragePct:            minCoveragePct,
		Aggregation:              "mean",
		TrendSensitivity:         m.TrendSensitivity,
		EnableVolatility:         true,
		VolatilityRatio:          0.5,
		Thresholds:               BuildStateThresholds(m),
		Direction:                m.StateDirection,
	}
	if policy.TrendSensitivity <= 0 {
		policy.TrendSensitivity = 1.0
	}
	if policy.ExpectedScrapeIntervalSec <= 0 {
		policy.ExpectedScrapeIntervalSec = float64(m.Step)
	}
	if policy.ExpectedScrapeIntervalSec <= 0 {
		policy.ExpectedScrapeIntervalSec = 1
	}
	return policy
}

// EvaluateMetric ties the 5 algorithms together into a full evaluation pipeline.
//
// samples: raw Prometheus results (should cover [now - 2*window_sec, now] for both windows)
// now: evaluation timestamp (Unix seconds)
func EvaluateMetric(metricID string, samples []metric.FilteredPrometheusResult, policy EvaluationPolicy, now int64) (MetricEvaluation, error) {
	nowF := float64(now)
	nowWindow := filterSamplesByTime(samples, nowF-policy.WindowSec, nowF)
	prevWindow := filterSamplesByTime(samples, nowF-2*policy.WindowSec, nowF-policy.WindowSec)

	coverageNow := estimateCoverage(nowWindow, policy.WindowSec, policy.ExpectedScrapeIntervalSec)
	if coverageNow < policy.MinCoveragePct {
		return MetricEvaluation{
			MetricID:    metricID,
			EvaluatedAt: now,
			State:       "unknown",
			Trend:       "unknown",
			Conclusion:  "unknown",
			Confidence:  0.2,
		}, nil
	}

	aggMethod := policy.Aggregation
	if aggMethod == "" {
		aggMethod = "mean"
	}
	aggNow, err := AggregateMetricAlgo(nowWindow, aggMethod)
	if err != nil {
		return MetricEvaluation{}, err
	}
	aggPrev, err := AggregateMetricAlgo(prevWindow, aggMethod)
	if err != nil {
		return MetricEvaluation{}, err
	}

	noise := aggNow.NoiseEstimate
	if math.IsNaN(noise) || math.IsInf(noise, 1) || noise < 0 {
		noise = noiseEpsilon
	}
	noiseClamped := noise + noiseEpsilon

	delta := aggNow.AggregatedValue - aggPrev.AggregatedValue

	trend, err := CalculateTrend(aggNow, aggPrev, policy.TrendSensitivity)
	if err != nil {
		return MetricEvaluation{}, err
	}
	if policy.EnableVolatility {
		if vol := maybeOverrideWithVolatility(noiseClamped, aggNow.AggregatedValue, policy.VolatilityRatio); vol != "" {
			trend = vol
		}
	}

	direction := policy.Direction
	if direction == "" {
		direction = DirectionHigherIsWorse
	}
	state, err := CalculateState(aggNow.AggregatedValue, policy.Thresholds, direction)
	if err != nil {
		return MetricEvaluation{}, err
	}

	conclusion := CombineStateAndTrend(state, trend)

	signalStrength := math.Abs(delta) / math.Max(noiseClamped, noiseEpsilon)
	if math.IsNaN(delta) || math.IsInf(delta, 0) {
		signalStrength = 0
	}
	valuesNow := extractNumericValues(nowWindow)
	consistency := trendConsistency(valuesNow)

	confidence := computeConfidenceFromComponents(coverageNow, signalStrength, consistency, state)

	projection := ProjectCrossingTime(
		aggNow.AggregatedValue,
		aggPrev.AggregatedValue,
		policy.WindowSec,
		policy.Thresholds,
		direction,
	)

	return MetricEvaluation{
		MetricID:    metricID,
		EvaluatedAt: now,
		AggNow:      aggNow.AggregatedValue,
		StatsNow:    aggNow,
		State:       state,
		Trend:       trend,
		Conclusion:  conclusion,
		Confidence:  confidence,
		Projection:  projection,
	}, nil
}
