/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package ukama_test

import (
	"testing"

	"github.com/tj/assert"

	"github.com/ukama/ukama/systems/common/ukama"
)

func TestSubscriberStatus(t *testing.T) {
	t.Run("ParseValidString_Active", func(tt *testing.T) {
		status := ukama.ParseSubscriberStatus("active")

		assert.NotNil(t, status)
		assert.Equal(t, status.String(), "active")
		assert.Equal(t, uint8(status), uint8(1))
	})

	t.Run("ParseValidString_PendingDeletion", func(tt *testing.T) {
		status := ukama.ParseSubscriberStatus("pending_deletion")

		assert.NotNil(t, status)
		assert.Equal(t, status.String(), "pending_deletion")
		assert.Equal(t, uint8(status), uint8(2))
	})

	t.Run("ParseValidNumber", func(tt *testing.T) {
		status := ukama.ParseSubscriberStatus("2")

		assert.NotNil(t, status)
		assert.Equal(t, uint8(status), uint8(2))
		assert.Equal(t, status.String(), "pending_deletion")
	})

	t.Run("ParseInvalidString", func(tt *testing.T) {
		status := ukama.ParseSubscriberStatus("deleted")

		assert.NotNil(t, status)
		assert.Equal(t, status.String(), "unknown")
		assert.Equal(t, uint8(status), uint8(0))
	})

	t.Run("ParseInvalidNumber", func(tt *testing.T) {
		status := ukama.ParseSubscriberStatus("10")

		assert.NotNil(t, status)
		assert.Equal(t, status.String(), "unknown")
		assert.Equal(t, uint8(status), uint8(0))
	})

	t.Run("StringConversion_ValidValue", func(tt *testing.T) {
		status := ukama.SubscriberStatusActive
		assert.Equal(t, status.String(), "active")
	})

	t.Run("StringConversion_InvalidValue", func(tt *testing.T) {
		status := ukama.SubscriberStatus(10)
		assert.Equal(t, status.String(), "unknown")
	})

	t.Run("CaseInsensitivity", func(tt *testing.T) {
		status1 := ukama.ParseSubscriberStatus("ACTIVE")
		status2 := ukama.ParseSubscriberStatus("active")
		
		assert.Equal(t, status1, status2)
		assert.Equal(t, uint8(status1), uint8(1))
	})

	t.Run("ConstantValues", func(tt *testing.T) {
		assert.Equal(t, uint8(ukama.SubscriberStatusUnknown), uint8(0))
		assert.Equal(t, uint8(ukama.SubscriberStatusActive), uint8(1))
		assert.Equal(t, uint8(ukama.SubscriberStatusPendingDeletion), uint8(2))
	})
}