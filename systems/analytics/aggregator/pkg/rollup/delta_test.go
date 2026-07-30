/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain it at http://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package rollup

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpValue_Delta(t *testing.T) {
	// opValue(op, sum, count, min, max, last)
	v, ok := opValue("DELTA", 0, 0, 100, 350, 350)
	assert.True(t, ok)
	assert.Equal(t, float64(250), v)

	v, ok = opValue("DELTA", 0, 0, 100, 40, 40) // reset -> clamp
	assert.True(t, ok)
	assert.Equal(t, float64(0), v)
}
