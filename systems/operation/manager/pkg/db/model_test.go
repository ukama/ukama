/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package db

import (
	"testing"

	"github.com/tj/assert"
)

func TestOperationStatus_String(t *testing.T) {
	assert.Equal(t, "pending", OperationPending.String())
	assert.Equal(t, "running", OperationRunning.String())
	assert.Equal(t, "success", OperationSuccess.String())
	assert.Equal(t, "failed", OperationFailed.String())
	assert.Equal(t, "timeout", OperationTimeout.String())
	assert.Equal(t, "cancelled", OperationCancelled.String())
}

func TestOperationStatus_IsTerminal(t *testing.T) {
	assert.False(t, OperationPending.IsTerminal())
	assert.False(t, OperationRunning.IsTerminal())
	assert.True(t, OperationSuccess.IsTerminal())
	assert.True(t, OperationFailed.IsTerminal())
	assert.True(t, OperationTimeout.IsTerminal())
	assert.True(t, OperationCancelled.IsTerminal())
}

func TestOperationStatus_ScanValue(t *testing.T) {
	var s OperationStatus

	assert.NoError(t, s.Scan(nil))
	assert.Equal(t, OperationPending, s)

	assert.NoError(t, s.Scan(int64(3)))
	assert.Equal(t, OperationFailed, s)

	v, err := OperationSuccess.Value()
	assert.NoError(t, err)
	assert.Equal(t, int64(2), v)
}

func TestJSONMap_ValueAndScan(t *testing.T) {
	t.Run("NilValue", func(t *testing.T) {
		var m JSONMap
		v, err := m.Value()
		assert.NoError(t, err)
		assert.Nil(t, v)
	})

	t.Run("RoundTrip", func(t *testing.T) {
		in := JSONMap{"node": "uk-1234", "count": float64(3)}
		v, err := in.Value()
		assert.NoError(t, err)

		var out JSONMap
		assert.NoError(t, out.Scan(v.([]byte)))
		assert.Equal(t, "uk-1234", out["node"])
		assert.Equal(t, float64(3), out["count"])
	})

	t.Run("ScanNil", func(t *testing.T) {
		var m JSONMap
		assert.NoError(t, m.Scan(nil))
		assert.Nil(t, m)
	})

	t.Run("ScanNonBytesIsNoOp", func(t *testing.T) {
		var m JSONMap
		assert.NoError(t, m.Scan("not-bytes"))
	})
}
