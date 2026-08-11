/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain it at http://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// DELTA on a cumulative counter is max-min over the span; a mid-span reset
// (max < min after folding) clamps to 0 rather than reporting negative usage.
func TestRollingOpValue_Delta(t *testing.T) {
	v, ok := rollingOpValue("DELTA", &windowAgg{min: 100, max: 350})
	assert.True(t, ok)
	assert.Equal(t, float64(250), v, "usage consumed across the span = max-min")

	// single observation -> no movement -> 0
	v, ok = rollingOpValue("DELTA", &windowAgg{min: 100, max: 100})
	assert.True(t, ok)
	assert.Equal(t, float64(0), v)

	// reset: max < min -> clamp to 0, never negative
	v, ok = rollingOpValue("DELTA", &windowAgg{min: 100, max: 40})
	assert.True(t, ok)
	assert.Equal(t, float64(0), v)
}
