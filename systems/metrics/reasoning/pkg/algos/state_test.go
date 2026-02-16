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
	"github.com/ukama/ukama/systems/metrics/reasoning/pkg"
)

func TestCalculateState(t *testing.T) {
	t.Run("HigherIsWorse_Healthy", func(t *testing.T) {
		thresholds := StateThresholds{Warning: 80, Critical: 95}
		state, err := CalculateState(50, thresholds, DirectionHigherIsWorse)
		require.NoError(t, err)
		assert.Equal(t, "healthy", state)
	})

	t.Run("HigherIsWorse_Warning", func(t *testing.T) {
		thresholds := StateThresholds{Warning: 80, Critical: 95}
		state, err := CalculateState(85, thresholds, DirectionHigherIsWorse)
		require.NoError(t, err)
		assert.Equal(t, "warning", state)
	})

	t.Run("HigherIsWorse_WarningAtBoundary", func(t *testing.T) {
		thresholds := StateThresholds{Warning: 80, Critical: 95}
		state, err := CalculateState(80, thresholds, DirectionHigherIsWorse)
		require.NoError(t, err)
		assert.Equal(t, "warning", state)
	})

	t.Run("HigherIsWorse_Critical", func(t *testing.T) {
		thresholds := StateThresholds{Warning: 80, Critical: 95}
		state, err := CalculateState(100, thresholds, DirectionHigherIsWorse)
		require.NoError(t, err)
		assert.Equal(t, "critical", state)
	})

	t.Run("HigherIsWorse_CriticalAtBoundary", func(t *testing.T) {
		thresholds := StateThresholds{Warning: 80, Critical: 95}
		state, err := CalculateState(95, thresholds, DirectionHigherIsWorse)
		require.NoError(t, err)
		assert.Equal(t, "critical", state)
	})

	t.Run("LowerIsWorse_Healthy", func(t *testing.T) {
		thresholds := StateThresholds{Warning: 80, Critical: 50}
		state, err := CalculateState(95, thresholds, DirectionLowerIsWorse)
		require.NoError(t, err)
		assert.Equal(t, "healthy", state)
	})

	t.Run("LowerIsWorse_Warning", func(t *testing.T) {
		thresholds := StateThresholds{Warning: 80, Critical: 50}
		state, err := CalculateState(75, thresholds, DirectionLowerIsWorse)
		require.NoError(t, err)
		assert.Equal(t, "warning", state)
	})

	t.Run("LowerIsWorse_WarningAtBoundary", func(t *testing.T) {
		thresholds := StateThresholds{Warning: 80, Critical: 50}
		state, err := CalculateState(80, thresholds, DirectionLowerIsWorse)
		require.NoError(t, err)
		assert.Equal(t, "warning", state)
	})

	t.Run("LowerIsWorse_Critical", func(t *testing.T) {
		thresholds := StateThresholds{Warning: 80, Critical: 50}
		state, err := CalculateState(30, thresholds, DirectionLowerIsWorse)
		require.NoError(t, err)
		assert.Equal(t, "critical", state)
	})

	t.Run("LowerIsWorse_CriticalAtBoundary", func(t *testing.T) {
		thresholds := StateThresholds{Warning: 80, Critical: 50}
		state, err := CalculateState(50, thresholds, DirectionLowerIsWorse)
		require.NoError(t, err)
		assert.Equal(t, "critical", state)
	})

	t.Run("Range_Healthy", func(t *testing.T) {
		thresholds := StateThresholds{
			LowWarning: 20, HighWarning: 80,
			LowCritical: 10, HighCritical: 90,
		}
		state, err := CalculateState(50, thresholds, DirectionRange)
		require.NoError(t, err)
		assert.Equal(t, "healthy", state)
	})

	t.Run("Range_WarningBelowBand", func(t *testing.T) {
		thresholds := StateThresholds{
			LowWarning: 20, HighWarning: 80,
			LowCritical: 10, HighCritical: 90,
		}
		state, err := CalculateState(15, thresholds, DirectionRange)
		require.NoError(t, err)
		assert.Equal(t, "warning", state)
	})

	t.Run("Range_WarningAboveBand", func(t *testing.T) {
		thresholds := StateThresholds{
			LowWarning: 20, HighWarning: 80,
			LowCritical: 10, HighCritical: 90,
		}
		state, err := CalculateState(85, thresholds, DirectionRange)
		require.NoError(t, err)
		assert.Equal(t, "warning", state)
	})

	t.Run("Range_CriticalBelowBand", func(t *testing.T) {
		thresholds := StateThresholds{
			LowWarning: 20, HighWarning: 80,
			LowCritical: 10, HighCritical: 90,
		}
		state, err := CalculateState(5, thresholds, DirectionRange)
		require.NoError(t, err)
		assert.Equal(t, "critical", state)
	})

	t.Run("Range_CriticalAboveBand", func(t *testing.T) {
		thresholds := StateThresholds{
			LowWarning: 20, HighWarning: 80,
			LowCritical: 10, HighCritical: 90,
		}
		state, err := CalculateState(95, thresholds, DirectionRange)
		require.NoError(t, err)
		assert.Equal(t, "critical", state)
	})

	t.Run("ValueNaN_ReturnsUnknown", func(t *testing.T) {
		thresholds := StateThresholds{Warning: 80, Critical: 95}
		state, err := CalculateState(math.NaN(), thresholds, DirectionHigherIsWorse)
		require.NoError(t, err)
		assert.Equal(t, "unknown", state)
	})

	t.Run("ValueInf_ReturnsUnknown", func(t *testing.T) {
		thresholds := StateThresholds{Warning: 80, Critical: 95}
		state, err := CalculateState(math.Inf(1), thresholds, DirectionHigherIsWorse)
		require.NoError(t, err)
		assert.Equal(t, "unknown", state)
	})

	t.Run("EmptyDirection_DefaultsToHigherIsWorse", func(t *testing.T) {
		thresholds := StateThresholds{Warning: 80, Critical: 95}
		state, err := CalculateState(90, thresholds, "")
		require.NoError(t, err)
		assert.Equal(t, "warning", state)
	})

	t.Run("UnknownDirection_ReturnsUnknown", func(t *testing.T) {
		thresholds := StateThresholds{Warning: 80, Critical: 95}
		state, err := CalculateState(50, thresholds, "invalid")
		require.NoError(t, err)
		assert.Equal(t, "unknown", state)
	})
}

func TestBuildStateThresholds(t *testing.T) {
	t.Run("HigherIsWorse", func(t *testing.T) {
		m := pkg.Metric{
			StateDirection: DirectionHigherIsWorse,
			Thresholds:     pkg.Thresholds{Min: 0, Medium: 80, Max: 100},
		}
		st := BuildStateThresholds(m)
		assert.Equal(t, 80.0, st.Warning)
		assert.Equal(t, 100.0, st.Critical)
	})

	t.Run("LowerIsWorse", func(t *testing.T) {
		m := pkg.Metric{
			StateDirection: DirectionLowerIsWorse,
			Thresholds:     pkg.Thresholds{Min: 50, Medium: 80, Max: 100},
		}
		st := BuildStateThresholds(m)
		assert.Equal(t, 80.0, st.Warning)
		assert.Equal(t, 50.0, st.Critical)
	})

	t.Run("Range", func(t *testing.T) {
		m := pkg.Metric{
			StateDirection: DirectionRange,
			Thresholds: pkg.Thresholds{
				LowWarning: 20, HighWarning: 80,
				LowCritical: 10, HighCritical: 90,
			},
		}
		st := BuildStateThresholds(m)
		assert.Equal(t, 20.0, st.LowWarning)
		assert.Equal(t, 80.0, st.HighWarning)
		assert.Equal(t, 10.0, st.LowCritical)
		assert.Equal(t, 90.0, st.HighCritical)
	})

	t.Run("EmptyDirection_DefaultsToHigherIsWorse", func(t *testing.T) {
		m := pkg.Metric{
			StateDirection: "",
			Thresholds:     pkg.Thresholds{Min: 0, Medium: 75, Max: 95},
		}
		st := BuildStateThresholds(m)
		assert.Equal(t, 75.0, st.Warning)
		assert.Equal(t, 95.0, st.Critical)
	})
}
