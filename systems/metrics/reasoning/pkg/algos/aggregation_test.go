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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/ukama/ukama/systems/metrics/reasoning/pkg/metric"
)

func makeSamples(values [][]interface{}) []metric.FilteredPrometheusResult {
	return []metric.FilteredPrometheusResult{
		{
			Metric: metric.FilteredMetric{NodeID: "node-1", Metric: "cpu"},
			Values: values,
		},
	}
}

func TestAggregateMetricAlgo(t *testing.T) {
	t.Run("EmptySamples", func(t *testing.T) {
		samples := makeSamples(nil)
		stats, err := AggregateMetricAlgo(samples, "mean")
		require.NoError(t, err)
		assert.True(t, math.IsNaN(stats.AggregatedValue))
		assert.Equal(t, "mean", stats.Aggregation)
		assert.Equal(t, float64(0), stats.SampleCount)
	})

	t.Run("Last", func(t *testing.T) {
		samples := makeSamples([][]interface{}{
			{1000.0, 10.0}, {2000.0, 20.0}, {3000.0, 30.0},
		})
		stats, err := AggregateMetricAlgo(samples, "last")
		require.NoError(t, err)
		assert.Equal(t, 30.0, stats.AggregatedValue)
		assert.Equal(t, float64(3), stats.SampleCount)
		assert.Equal(t, 10.0, stats.Min)
		assert.Equal(t, 30.0, stats.Max)
	})

	t.Run("Mean", func(t *testing.T) {
		samples := makeSamples([][]interface{}{
			{1000.0, 10.0}, {2000.0, 20.0}, {3000.0, 30.0},
		})
		stats, err := AggregateMetricAlgo(samples, "mean")
		require.NoError(t, err)
		assert.Equal(t, 20.0, stats.AggregatedValue)
		assert.Equal(t, 20.0, stats.Mean)
	})

	t.Run("Median", func(t *testing.T) {
		samples := makeSamples([][]interface{}{
			{1000.0, 10.0}, {2000.0, 20.0}, {3000.0, 30.0},
		})
		stats, err := AggregateMetricAlgo(samples, "median")
		require.NoError(t, err)
		assert.Equal(t, 20.0, stats.AggregatedValue)
		assert.Equal(t, 20.0, stats.Median)
	})

	t.Run("MedianEvenCount", func(t *testing.T) {
		samples := makeSamples([][]interface{}{
			{1000.0, 10.0}, {2000.0, 20.0}, {3000.0, 30.0}, {4000.0, 40.0},
		})
		stats, err := AggregateMetricAlgo(samples, "median")
		require.NoError(t, err)
		assert.Equal(t, 25.0, stats.AggregatedValue)
	})

	t.Run("P95", func(t *testing.T) {
		samples := makeSamples([][]interface{}{
			{1000.0, 10.0}, {2000.0, 20.0}, {3000.0, 30.0}, {4000.0, 40.0},
		})
		stats, err := AggregateMetricAlgo(samples, "p95")
		require.NoError(t, err)
		// p95: idx = (95/100)*3 = 2.85, lo=2, hi=3; interpolate: 30 + 0.85*(40-30) = 38.5
		assert.InDelta(t, 38.5, stats.AggregatedValue, 0.01)
	})

	t.Run("Sum", func(t *testing.T) {
		samples := makeSamples([][]interface{}{
			{1000.0, 10.0}, {2000.0, 20.0}, {3000.0, 30.0},
		})
		stats, err := AggregateMetricAlgo(samples, "sum")
		require.NoError(t, err)
		assert.Equal(t, 60.0, stats.AggregatedValue)
	})

	t.Run("UnknownMethod", func(t *testing.T) {
		samples := makeSamples([][]interface{}{{1000.0, 10.0}})
		stats, err := AggregateMetricAlgo(samples, "unknown")
		require.NoError(t, err)
		assert.True(t, math.IsNaN(stats.AggregatedValue))
		assert.Equal(t, "unknown", stats.Aggregation)
	})

	t.Run("MultipleResults", func(t *testing.T) {
		samples := []metric.FilteredPrometheusResult{
			{Metric: metric.FilteredMetric{NodeID: "node-1", Metric: "cpu"}, Values: [][]interface{}{{1000.0, 5.0}, {2000.0, 15.0}}},
			{Metric: metric.FilteredMetric{NodeID: "node-2", Metric: "cpu"}, Values: [][]interface{}{{1000.0, 25.0}}},
		}
		stats, err := AggregateMetricAlgo(samples, "mean")
		require.NoError(t, err)
		assert.Equal(t, 15.0, stats.AggregatedValue)
		assert.Equal(t, float64(3), stats.SampleCount)
	})

	t.Run("ExtractNumericValuesTypeConversions", func(t *testing.T) {
		samples := makeSamples([][]interface{}{
			{1000.0, json.Number("10.5")},
			{2000.0, "20.5"},
			{3000.0, 30},
			{4000.0, int64(40)},
			{5000.0, 50.0},
		})
		stats, err := AggregateMetricAlgo(samples, "mean")
		require.NoError(t, err)
		assert.InDelta(t, 30.2, stats.AggregatedValue, 0.01)
		assert.Equal(t, float64(5), stats.SampleCount)
	})

	t.Run("SkipsMalformedPairs", func(t *testing.T) {
		samples := makeSamples([][]interface{}{
			{1000.0, 10.0},
			{2000.0},           // missing value
			{3000.0, "notnum"}, // invalid number
			{4000.0, 40.0},
		})
		stats, err := AggregateMetricAlgo(samples, "mean")
		require.NoError(t, err)
		assert.Equal(t, 25.0, stats.AggregatedValue)
		assert.Equal(t, float64(2), stats.SampleCount)
	})

	t.Run("MedianAbsoluteDeviation", func(t *testing.T) {
		// values: 1, 2, 3, 4, 5 -> median=3; deviations: 2, 1, 0, 1, 2 -> MAD = 1
		samples := makeSamples([][]interface{}{
			{1000.0, 1.0}, {2000.0, 2.0}, {3000.0, 3.0}, {4000.0, 4.0}, {5000.0, 5.0},
		})
		stats, err := AggregateMetricAlgo(samples, "mean")
		require.NoError(t, err)
		assert.InDelta(t, 1.0, stats.NoiseEstimate, 0.001)
	})

	t.Run("MedianAbsoluteDeviationSingleSample", func(t *testing.T) {
		samples := makeSamples([][]interface{}{{1000.0, 10.0}})
		stats, err := AggregateMetricAlgo(samples, "mean")
		require.NoError(t, err)
		assert.True(t, math.IsNaN(stats.NoiseEstimate))
	})
}

func TestAggregationStats(t *testing.T) {
	t.Run("RoundOfDecimalPoints", func(t *testing.T) {
		stats := AggregationStats{
			AggregatedValue: 3.14159,
			Min:             1.12345,
			Max:             5.98765,
			P95:             4.56789,
			Mean:            3.33333,
			Median:          3.2,
			SampleCount:     10.7, // rounds to 11
			NoiseEstimate:   0.55555,
		}
		stats.RoundOfDecimalPoints(2)
		assert.Equal(t, 3.14, stats.AggregatedValue)
		assert.Equal(t, 1.12, stats.Min)
		assert.Equal(t, 5.99, stats.Max)
		assert.Equal(t, 4.57, stats.P95)
		assert.Equal(t, 3.33, stats.Mean)
		assert.Equal(t, 3.2, stats.Median)
		assert.Equal(t, 11.0, stats.SampleCount)
		assert.Equal(t, 0.56, stats.NoiseEstimate)
	})

	t.Run("MarshalJSONNaNInf", func(t *testing.T) {
		stats := AggregationStats{
			AggregatedValue: math.NaN(),
			Min:             math.Inf(1),
			Max:             math.Inf(-1),
			Mean:            10.5,
		}
		data, err := json.Marshal(stats)
		require.NoError(t, err)
		var m map[string]interface{}
		err = json.Unmarshal(data, &m)
		require.NoError(t, err)
		assert.Nil(t, m["AggregatedValue"])
		assert.Nil(t, m["Min"])
		assert.Nil(t, m["Max"])
		assert.Equal(t, 10.5, m["Mean"])
	})

	t.Run("UnmarshalJSONNullToNaN", func(t *testing.T) {
		data := []byte(`{"AggregatedValue":null,"Min":null,"Max":5.0,"Mean":null,"Median":null,"P95":null,"SampleCount":3,"Aggregation":"mean","NoiseEstimate":null}`)
		var stats AggregationStats
		err := json.Unmarshal(data, &stats)
		require.NoError(t, err)
		assert.True(t, math.IsNaN(stats.AggregatedValue))
		assert.True(t, math.IsNaN(stats.Min))
		assert.Equal(t, 5.0, stats.Max)
		assert.True(t, math.IsNaN(stats.Mean))
		assert.True(t, math.IsNaN(stats.Median))
		assert.True(t, math.IsNaN(stats.P95))
		assert.Equal(t, float64(3), stats.SampleCount)
		assert.Equal(t, "mean", stats.Aggregation)
		assert.True(t, math.IsNaN(stats.NoiseEstimate))
	})

	t.Run("MarshalUnmarshalRoundtrip", func(t *testing.T) {
		original := AggregationStats{
			AggregatedValue: 42.0,
			Min:             10.0,
			Max:             100.0,
			P95:             95.0,
			Mean:             50.0,
			Median:           45.0,
			SampleCount:      10,
			Aggregation:      "mean",
			NoiseEstimate:    5.5,
		}
		data, err := json.Marshal(original)
		require.NoError(t, err)
		var decoded AggregationStats
		err = json.Unmarshal(data, &decoded)
		require.NoError(t, err)
		assert.Equal(t, original, decoded)
	})
}

func TestUnmarshalAggStats(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		data := []byte(`{"AggregatedValue":42.5,"Min":10,"Max":100,"P95":95,"Mean":50,"Median":45,"SampleCount":10,"Aggregation":"median","NoiseEstimate":5.2}`)
		stats, err := UnmarshalAggStats(data)
		require.NoError(t, err)
		assert.Equal(t, 42.5, stats.AggregatedValue)
		assert.Equal(t, 10.0, stats.Min)
		assert.Equal(t, 100.0, stats.Max)
		assert.Equal(t, 95.0, stats.P95)
		assert.Equal(t, 50.0, stats.Mean)
		assert.Equal(t, 45.0, stats.Median)
		assert.Equal(t, float64(10), stats.SampleCount)
		assert.Equal(t, "median", stats.Aggregation)
		assert.Equal(t, 5.2, stats.NoiseEstimate)
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		_, err := UnmarshalAggStats([]byte(`invalid json`))
		assert.Error(t, err)
	})
}
