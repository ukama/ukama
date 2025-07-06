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

func TestStatusType(t *testing.T) {
	t.Run("StatusTypeValidString", func(tt *testing.T) {
		status := ukama.ParseStatusType("pending")

		assert.NotNil(t, status)
		assert.Equal(t, status.String(), ukama.StatusTypePending.String())
		assert.Equal(t, uint8(status), uint8(1))
	})

	t.Run("StatusTypeValidNumber", func(tt *testing.T) {
		status := ukama.ParseStatusType("3")

		assert.NotNil(t, status)
		assert.Equal(t, uint8(status), uint8(3))
		assert.Equal(t, status.String(), ukama.StatusTypeFailed.String())
	})

	t.Run("StatusTypeNonValidString", func(tt *testing.T) {
		status := ukama.ParseStatusType("failure")

		assert.NotNil(t, status)
		assert.Equal(t, status.String(), ukama.StatusTypeUnknown.String())
		assert.Equal(t, uint8(status), uint8(0))
	})

	t.Run("StatusTypeNonValidNumber", func(tt *testing.T) {
		status := ukama.StatusType(uint8(10))

		assert.NotNil(t, status)
		assert.Equal(t, status.String(), ukama.StatusTypeUnknown.String())
		assert.Equal(t, uint8(status), uint8(10))
	})
}
