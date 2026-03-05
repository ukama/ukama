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

func TestSoftwareStatusType(t *testing.T) {
	t.Run("ValidString", func(tt *testing.T) {
		status := ukama.ParseSoftwareStatusType("update_available")

		assert.NotNil(tt, status)
		assert.Equal(tt, status.String(), ukama.SoftwareStatusTypeUpdateAvailable.String())
		assert.Equal(tt, uint8(status), uint8(1))
	})

	t.Run("ValidStringCaseInsensitive", func(tt *testing.T) {
		status := ukama.ParseSoftwareStatusType("UPDATE_IN_PROGRESS")

		assert.NotNil(tt, status)
		assert.Equal(tt, status.String(), ukama.SoftwareStatusTypeUpdateInProgress.String())
		assert.Equal(tt, uint8(status), uint8(3))
	})

	t.Run("ValidNumber", func(tt *testing.T) {
		status := ukama.ParseSoftwareStatusType("4")

		assert.NotNil(tt, status)
		assert.Equal(tt, uint8(status), uint8(4))
		assert.Equal(tt, status.String(), ukama.SoftwareStatusTypeUpdateFailed.String())
	})

	t.Run("NonValidString", func(tt *testing.T) {
		status := ukama.ParseSoftwareStatusType("invalid_status")

		assert.NotNil(tt, status)
		assert.Equal(tt, status.String(), ukama.SoftwareStatusTypeUnknown.String())
		assert.Equal(tt, uint8(status), uint8(0))
	})

	t.Run("NonValidNumber", func(tt *testing.T) {
		status := ukama.SoftwareStatusType(uint8(10))

		assert.NotNil(tt, status)
		assert.Equal(tt, status.String(), ukama.SoftwareStatusTypeUnknown.String())
		assert.Equal(tt, uint8(status), uint8(10))
	})

	t.Run("AllConstantsString", func(tt *testing.T) {
		cases := []struct {
			parsed string
			constant ukama.SoftwareStatusType
			expectStr string
		}{
			{"unknown", ukama.SoftwareStatusTypeUnknown, "unknown"},
			{"update_available", ukama.SoftwareStatusTypeUpdateAvailable, "update_available"},
			{"up_to_date", ukama.SoftwareStatusTypeUpToDate, "up_to_date"},
			{"update_in_progress", ukama.SoftwareStatusTypeUpdateInProgress, "update_in_progress"},
			{"update_failed", ukama.SoftwareStatusTypeUpdateFailed, "update_failed"},
		}
		for _, c := range cases {
			status := ukama.ParseSoftwareStatusType(c.parsed)
			assert.Equal(tt, status, c.constant)
			assert.Equal(tt, status.String(), c.expectStr)
		}
	})

	t.Run("Value", func(tt *testing.T) {
		status := ukama.SoftwareStatusTypeUpToDate
		val, err := status.Value()

		assert.NoError(tt, err)
		assert.Equal(tt, uint8(2), val)
	})

	t.Run("Scan", func(tt *testing.T) {
		var status ukama.SoftwareStatusType
		err := status.Scan(int64(3))

		assert.NoError(tt, err)
		assert.Equal(tt, ukama.SoftwareStatusTypeUpdateInProgress, status)
	})
}
