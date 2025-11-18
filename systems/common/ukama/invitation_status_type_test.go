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

func TestInvitationStatusType(t *testing.T) {
	t.Run("StatusTypeValidString", func(tt *testing.T) {
		status := ukama.ParseInvitationStatusType("pending")

		assert.NotNil(t, status)
		assert.Equal(t, status.String(), ukama.Pending.String())
		assert.Equal(t, uint8(status), uint8(0))
	})

	t.Run("StatusTypeValidNumber", func(tt *testing.T) {
		status := ukama.ParseInvitationStatusType("2")

		assert.NotNil(t, status)
		assert.Equal(t, uint8(status), uint8(2))
		assert.Equal(t, status.String(), ukama.Declined.String())
	})

	t.Run("StatusTypeNonValidString", func(tt *testing.T) {
		status := ukama.ParseInvitationStatusType("failure")

		assert.NotNil(t, status)
		assert.Equal(t, status.String(), ukama.Pending.String())
		assert.Equal(t, uint8(status), uint8(0))
	})

	t.Run("StatusTypeNonValidNumber", func(tt *testing.T) {
		status := ukama.InvitationStatusType(uint8(10))

		assert.NotNil(t, status)
		assert.Equal(t, status.String(), ukama.Pending.String())
		assert.Equal(t, uint8(status), uint8(10))
	})

	t.Run("StatusTypeAcceptedString", func(tt *testing.T) {
		status := ukama.ParseInvitationStatusType("accepted")

		assert.NotNil(t, status)
		assert.Equal(t, status.String(), ukama.Accepted.String())
		assert.Equal(t, uint8(status), uint8(1))
	})

	t.Run("StatusTypeDeclinedString", func(tt *testing.T) {
		status := ukama.ParseInvitationStatusType("declined")

		assert.NotNil(t, status)
		assert.Equal(t, status.String(), ukama.Declined.String())
		assert.Equal(t, uint8(status), uint8(2))
	})

	t.Run("StatusTypeCaseInsensitive", func(tt *testing.T) {
		status := ukama.ParseInvitationStatusType("ACCEPTED")

		assert.NotNil(t, status)
		assert.Equal(t, status.String(), ukama.Accepted.String())
		assert.Equal(t, uint8(status), uint8(1))
	})

	t.Run("StatusTypeValidNumberZero", func(tt *testing.T) {
		status := ukama.ParseInvitationStatusType("0")

		assert.NotNil(t, status)
		assert.Equal(t, status.String(), ukama.Pending.String())
		assert.Equal(t, uint8(status), uint8(0))
	})

	t.Run("StatusTypeValidNumberOne", func(tt *testing.T) {
		status := ukama.ParseInvitationStatusType("1")

		assert.NotNil(t, status)
		assert.Equal(t, status.String(), ukama.Accepted.String())
		assert.Equal(t, uint8(status), uint8(1))
	})
}
