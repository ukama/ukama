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

func projectionMakeSamples(values [][]interface{}) []metric.FilteredPrometheusResult {
	if values == nil {
		return nil
	}
	return []metric.FilteredPrometheusResult{
		{Metric: metric.FilteredMetric{NodeID: "node-1", Metric: "cpu"}, Values: values},
	}
}

func TestProjectCrossingTime(t *testing.T) {
	thresholds := StateThresholds{Warning: 80, Critical: 95}

	t.Run("ValueNowNaN_ReturnsEmpty", func(t *testing.T) {
		result := ProjectCrossingTime(math.NaN(), 60, 300, thresholds, DirectionHigherIsWorse)
		assert.Empty(t, result.Type)
		assert.Equal(t, 0.0, result.EtaSec)
	})

	t.Run("ValuePrevNaN_ReturnsEmpty", func(t *testing.T) {
		result := ProjectCrossingTime(70, math.NaN(), 300, thresholds, DirectionHigherIsWorse)
		assert.Empty(t, result.Type)
	})

	t.Run("ValueNowInf_ReturnsEmpty", func(t *testing.T) {
		result := ProjectCrossingTime(math.Inf(1), 60, 300, thresholds, DirectionHigherIsWorse)
		assert.Empty(t, result.Type)
	})

	t.Run("ValuePrevInf_ReturnsEmpty", func(t *testing.T) {
		result := ProjectCrossingTime(70, math.Inf(-1), 300, thresholds, DirectionHigherIsWorse)
		assert.Empty(t, result.Type)
	})

	t.Run("WindowSecZero_ReturnsEmpty", func(t *testing.T) {
		result := ProjectCrossingTime(70, 60, 0, thresholds, DirectionHigherIsWorse)
		assert.Empty(t, result.Type)
	})

	t.Run("WindowSecNegative_ReturnsEmpty", func(t *testing.T) {
		result := ProjectCrossingTime(70, 60, -100, thresholds, DirectionHigherIsWorse)
		assert.Empty(t, result.Type)
	})

	t.Run("HigherIsWorse_IncreasingTowardWarning_ReturnsProjection", func(t *testing.T) {
		// valueNow=70, valuePrev=60, windowSec=300 -> slope = 10/300 = 0.0333/sec
		// target 80, distance 10 -> eta = 10/0.0333 ≈ 300 sec
		result := ProjectCrossingTime(70, 60, 300, thresholds, DirectionHigherIsWorse)
		assert.Equal(t, "to_warning", result.Type)
		assert.InDelta(t, 300, result.EtaSec, 1)
	})

	t.Run("HigherIsWorse_EmptyDirection_DefaultsToHigherIsWorse", func(t *testing.T) {
		result := ProjectCrossingTime(70, 60, 300, thresholds, "")
		assert.Equal(t, "to_warning", result.Type)
		assert.InDelta(t, 300, result.EtaSec, 1)
	})

	t.Run("HigherIsWorse_AlreadyAtOrBeyondWarning_ReturnsEmpty", func(t *testing.T) {
		result := ProjectCrossingTime(80, 70, 300, thresholds, DirectionHigherIsWorse)
		assert.Empty(t, result.Type)
		result = ProjectCrossingTime(90, 85, 300, thresholds, DirectionHigherIsWorse)
		assert.Empty(t, result.Type)
	})

	t.Run("HigherIsWorse_SlopeZeroOrNegative_ReturnsEmpty", func(t *testing.T) {
		result := ProjectCrossingTime(70, 70, 300, thresholds, DirectionHigherIsWorse)
		assert.Empty(t, result.Type)
		result = ProjectCrossingTime(70, 75, 300, thresholds, DirectionHigherIsWorse)
		assert.Empty(t, result.Type)
	})

	t.Run("LowerIsWorse_DecreasingTowardWarning_ReturnsProjection", func(t *testing.T) {
		// success rate: warning at 80%, valueNow=90, valuePrev=95, windowSec=300
		// slope = -5/300 = -0.0167/sec (decreasing)
		// distance to warning = 90-80 = 10, eta = 10/0.0167 ≈ 600 sec
		thresholds := StateThresholds{Warning: 80, Critical: 50}
		result := ProjectCrossingTime(90, 95, 300, thresholds, DirectionLowerIsWorse)
		assert.Equal(t, "to_warning", result.Type)
		assert.InDelta(t, 600, result.EtaSec, 1)
	})

	t.Run("LowerIsWorse_AlreadyAtOrBelowWarning_ReturnsEmpty", func(t *testing.T) {
		thresholds := StateThresholds{Warning: 80, Critical: 50}
		result := ProjectCrossingTime(80, 90, 300, thresholds, DirectionLowerIsWorse)
		assert.Empty(t, result.Type)
		result = ProjectCrossingTime(70, 75, 300, thresholds, DirectionLowerIsWorse)
		assert.Empty(t, result.Type)
	})

	t.Run("LowerIsWorse_SlopeZeroOrPositive_ReturnsEmpty", func(t *testing.T) {
		thresholds := StateThresholds{Warning: 80, Critical: 50}
		result := ProjectCrossingTime(90, 90, 300, thresholds, DirectionLowerIsWorse)
		assert.Empty(t, result.Type)
		result = ProjectCrossingTime(90, 85, 300, thresholds, DirectionLowerIsWorse)
		assert.Empty(t, result.Type)
	})

	t.Run("DirectionRange_ReturnsEmpty", func(t *testing.T) {
		thresholds := StateThresholds{LowWarning: 10, HighWarning: 90}
		result := ProjectCrossingTime(50, 45, 300, thresholds, DirectionRange)
		assert.Empty(t, result.Type)
	})

	t.Run("UnknownDirection_ReturnsEmpty", func(t *testing.T) {
		result := ProjectCrossingTime(70, 60, 300, thresholds, "invalid")
		assert.Empty(t, result.Type)
	})

	t.Run("HigherIsWorse_EtaSecCalculation", func(t *testing.T) {
		// valueNow=50, valuePrev=40, windowSec=100 -> slope = 0.1/sec
		// target 80, distance 30 -> eta = 30/0.1 = 300 sec
		result := ProjectCrossingTime(50, 40, 100, thresholds, DirectionHigherIsWorse)
		assert.Equal(t, "to_warning", result.Type)
		assert.InDelta(t, 300, result.EtaSec, 0.1)
	})
}

func TestProjectMetricResults(t *testing.T) {
	t.Run("ReturnsInputUnchanged", func(t *testing.T) {
		samples := projectionMakeSamples([][]interface{}{
			{1000.0, 10.0}, {2000.0, 20.0},
		})
		result, err := ProjectMetricResults(samples)
		require.NoError(t, err)
		assert.Equal(t, samples, result)
	})

	t.Run("NilInput_ReturnsNil", func(t *testing.T) {
		result, err := ProjectMetricResults(nil)
		require.NoError(t, err)
		assert.Nil(t, result)
	})

	t.Run("EmptySlice_ReturnsEmpty", func(t *testing.T) {
		result, err := ProjectMetricResults([]metric.FilteredPrometheusResult{})
		require.NoError(t, err)
		assert.Empty(t, result)
	})
}
