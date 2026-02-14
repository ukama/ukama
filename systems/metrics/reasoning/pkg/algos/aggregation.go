/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package algos

import (
	"encoding/json"
	"math"
	"sort"
	"strconv"

	"github.com/ukama/ukama/systems/metrics/reasoning/pkg/metric"
	"github.com/ukama/ukama/systems/metrics/reasoning/utils"
)

type AggregationStats struct {
	AggregatedValue float64
	Min             float64
	Max             float64
	P95             float64
	Mean            float64
	Median          float64
	SampleCount     float64
	Aggregation     string
	NoiseEstimate   float64
}

// jsonSafeFloat returns the value for JSON marshal; NaN/Inf become nil (marshals as null).
func jsonSafeFloat(f float64) any {
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return nil
	}
	return f
}

func UnmarshalAggStats(data []byte) (AggregationStats, error) {
	var aggStats AggregationStats
	err := json.Unmarshal(data, &aggStats)
	if err != nil {
		return AggregationStats{}, err
	}
	return aggStats, nil
}

// MarshalJSON marshals AggregationStats, converting NaN and Inf to null.
func (s AggregationStats) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		AggregatedValue any `json:"AggregatedValue"`
		Min             any `json:"Min"`
		Max             any `json:"Max"`
		P95             any `json:"P95"`
		Mean            any `json:"Mean"`
		Median          any `json:"Median"`
		SampleCount     any `json:"SampleCount"`
		Aggregation     string `json:"Aggregation"`
		NoiseEstimate   any    `json:"NoiseEstimate"`
	}{
		AggregatedValue: jsonSafeFloat(s.AggregatedValue),
		Min:             jsonSafeFloat(s.Min),
		Max:             jsonSafeFloat(s.Max),
		P95:             jsonSafeFloat(s.P95),
		Mean:            jsonSafeFloat(s.Mean),
		Median:          jsonSafeFloat(s.Median),
		SampleCount:     jsonSafeFloat(s.SampleCount),
		Aggregation:     s.Aggregation,
		NoiseEstimate:   jsonSafeFloat(s.NoiseEstimate),
	})
}

// UnmarshalJSON unmarshals AggregationStats, converting null to NaN for float fields.
func (s *AggregationStats) UnmarshalJSON(data []byte) error {
	var aux struct {
		AggregatedValue *float64 `json:"AggregatedValue"`
		Min             *float64 `json:"Min"`
		Max             *float64 `json:"Max"`
		P95             *float64 `json:"P95"`
		Mean            *float64 `json:"Mean"`
		Median          *float64 `json:"Median"`
		SampleCount     *float64 `json:"SampleCount"`
		Aggregation     string   `json:"Aggregation"`
		NoiseEstimate   *float64 `json:"NoiseEstimate"`
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	s.AggregatedValue = ptrToFloat(aux.AggregatedValue)
	s.Min = ptrToFloat(aux.Min)
	s.Max = ptrToFloat(aux.Max)
	s.P95 = ptrToFloat(aux.P95)
	s.Mean = ptrToFloat(aux.Mean)
	s.Median = ptrToFloat(aux.Median)
	s.SampleCount = ptrToFloat(aux.SampleCount)
	s.Aggregation = aux.Aggregation
	s.NoiseEstimate = ptrToFloat(aux.NoiseEstimate)
	return nil
}

func (s *AggregationStats) RoundOfDecimalPoints(decimalPoints int) {
	s.AggregatedValue = utils.RoundToDecimalPoints(s.AggregatedValue, decimalPoints)
	s.Min = utils.RoundToDecimalPoints(s.Min, decimalPoints)
	s.Max = utils.RoundToDecimalPoints(s.Max, decimalPoints)
	s.P95 = utils.RoundToDecimalPoints(s.P95, decimalPoints)
	s.Mean = utils.RoundToDecimalPoints(s.Mean, decimalPoints)
	s.Median = utils.RoundToDecimalPoints(s.Median, decimalPoints)
	s.NoiseEstimate = utils.RoundToDecimalPoints(s.NoiseEstimate, decimalPoints)
	s.SampleCount = utils.RoundToDecimalPoints(s.SampleCount, 0) // count is always whole number
}

func ptrToFloat(p *float64) float64 {
	if p == nil {
		return math.NaN()
	}
	return *p
}

func AggregateMetricAlgo(windowSamples []metric.FilteredPrometheusResult, method string) (AggregationStats, error) {
	stats := AggregationStats{Aggregation: method}

	values := extractNumericValues(windowSamples)
	if len(values) == 0 {
		stats.AggregatedValue = math.NaN()
		return stats, nil
	}

	stats.SampleCount = float64(len(values))
	stats.Min = min(values)
	stats.Max = max(values)
	stats.Median = median(values)
	stats.Mean = mean(values)
	stats.P95 = percentile(values, 95)
	stats.NoiseEstimate = medianAbsoluteDeviation(values)

	switch method {
	case "last":
		stats.AggregatedValue = values[len(values)-1]
	case "mean":
		stats.AggregatedValue = stats.Mean
	case "median":
		stats.AggregatedValue = stats.Median
	case "p95":
		stats.AggregatedValue = stats.P95
	case "sum":
		stats.AggregatedValue = sum(values)
	default:
		stats.AggregatedValue = math.NaN()
	}

	return stats, nil
}

func extractNumericValues(samples []metric.FilteredPrometheusResult) []float64 {
	var values []float64
	for _, s := range samples {
		for _, pair := range s.Values {
			if len(pair) < 2 {
				continue
			}
			if v, ok := toFloat64(pair[1]); ok {
				values = append(values, v)
			}
		}
	}
	return values
}

func toFloat64(v any) (float64, bool) {
	switch x := v.(type) {
	case float64:
		return x, true
	case json.Number:
		f, err := x.Float64()
		return f, err == nil
	case string:
		f, err := strconv.ParseFloat(x, 64)
		return f, err == nil
	case int:
		return float64(x), true
	case int64:
		return float64(x), true
	default:
		return 0, false
	}
}

// extractTimestamp returns the timestamp from a [ts, value] pair for time-range filtering.
func extractTimestamp(pair []interface{}) (float64, bool) {
	if len(pair) < 1 {
		return 0, false
	}
	return toFloat64(pair[0])
}

func min(values []float64) float64 {
	if len(values) == 0 {
		return math.NaN()
	}
	m := values[0]
	for _, v := range values[1:] {
		if v < m {
			m = v
		}
	}
	return m
}

func max(values []float64) float64 {
	if len(values) == 0 {
		return math.NaN()
	}
	m := values[0]
	for _, v := range values[1:] {
		if v > m {
			m = v
		}
	}
	return m
}

func sum(values []float64) float64 {
	var s float64
	for _, v := range values {
		s += v
	}
	return s
}

func mean(values []float64) float64 {
	if len(values) == 0 {
		return math.NaN()
	}
	return sum(values) / float64(len(values))
}

func median(values []float64) float64 {
	if len(values) == 0 {
		return math.NaN()
	}
	sorted := make([]float64, len(values))
	copy(sorted, values)
	sort.Float64s(sorted)
	n := len(sorted)
	if n%2 == 1 {
		return sorted[n/2]
	}
	return (sorted[n/2-1] + sorted[n/2]) / 2
}

func percentile(values []float64, p float64) float64 {
	if len(values) == 0 {
		return math.NaN()
	}
	sorted := make([]float64, len(values))
	copy(sorted, values)
	sort.Float64s(sorted)
	idx := (p / 100.0) * float64(len(sorted)-1)
	lo := int(idx)
	hi := lo + 1
	if hi >= len(sorted) {
		return sorted[len(sorted)-1]
	}
	return sorted[lo] + (idx-float64(lo))*(sorted[hi]-sorted[lo])
}

// medianAbsoluteDeviation returns median of absolute deviations from median (MAD).
func medianAbsoluteDeviation(values []float64) float64 {
	if len(values) < 2 {
		return math.NaN()
	}
	med := median(values)
	absDev := make([]float64, len(values))
	for i, v := range values {
		absDev[i] = math.Abs(v - med)
	}
	return median(absDev)
}