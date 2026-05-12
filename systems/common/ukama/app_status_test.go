/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package ukama_test

import (
	"testing"

	"github.com/tj/assert"

	"github.com/ukama/ukama/systems/common/ukama"
)

func TestAppStatus(t *testing.T) {
	t.Run("ParseValidStringNames", func(tt *testing.T) {
		cases := []struct {
			parsed   string
			constant ukama.AppStatus
		}{
			{"unknown", ukama.AppStatusUnknown},
			{"pending", ukama.AppStatusPending},
			{"active", ukama.AppStatusActive},
			{"inactive", ukama.AppStatusInactive},
			{"running", ukama.AppStatusRunning},
			{"stopped", ukama.AppStatusStopped},
		}
		for _, c := range cases {
			got := ukama.ParseAppStatus(c.parsed)
			assert.Equal(tt, c.constant, got)
			assert.Equal(tt, uint8(c.constant), uint8(got))
		}
	})

	t.Run("ParseValidNumericString", func(tt *testing.T) {
		got := ukama.ParseAppStatus("2")
		assert.Equal(tt, ukama.AppStatusActive, got)
		assert.Equal(tt, "active", got.String())
	})

	t.Run("ParseInvalidStringMapsToUnknown", func(tt *testing.T) {
		got := ukama.ParseAppStatus("not_a_status")
		assert.Equal(tt, ukama.AppStatusUnknown, got)
		assert.Equal(tt, ukama.AppStatusUnknown.String(), got.String())
	})

	t.Run("StringUnknownForUndefinedValue", func(tt *testing.T) {
		s := ukama.AppStatus(99)
		assert.Equal(tt, ukama.AppStatusUnknown.String(), s.String())
		assert.Equal(tt, uint8(99), uint8(s))
	})

	t.Run("ParseCaseSensitiveName", func(tt *testing.T) {
		got := ukama.ParseAppStatus("ACTIVE")
		assert.Equal(tt, ukama.AppStatusUnknown, got)
	})

	t.Run("Value", func(tt *testing.T) {
		val, err := ukama.AppStatusRunning.Value()

		assert.NoError(tt, err)
		assert.Equal(tt, int64(ukama.AppStatusRunning), val)
	})

	t.Run("Scan", func(tt *testing.T) {
		var s ukama.AppStatus
		err := s.Scan(int64(4))

		assert.NoError(tt, err)
		assert.Equal(tt, ukama.AppStatusRunning, s)
	})
}
