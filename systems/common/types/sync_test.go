/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package types_test

import (
	"testing"

	"github.com/tj/assert"
	"github.com/ukama/ukama/systems/common/types"
)

func TestSyncStatus(t *testing.T) {
	t.Run("SyncStatusValidString", func(tt *testing.T) {
		syncStatus := types.ParseStatus("pending")

		assert.NotNil(t, syncStatus)
		assert.Equal(t, syncStatus.String(), "pending")
		assert.Equal(t, uint8(syncStatus), uint8(1))
	})

	t.Run("SyncStatusValidNumber", func(tt *testing.T) {
		syncStatus := types.ParseStatus("3")

		assert.NotNil(t, syncStatus)
		assert.Equal(t, uint8(syncStatus), uint8(3))
		assert.Equal(t, syncStatus.String(), "completed")
	})

	t.Run("SyncStatusNonValidString", func(tt *testing.T) {
		syncStatus := types.ParseStatus("failure")

		assert.NotNil(t, syncStatus)
		assert.Equal(t, syncStatus.String(), "unknown")
		assert.Equal(t, uint8(syncStatus), uint8(0))
	})

	t.Run("SyncStatusNonValidNumber", func(tt *testing.T) {
		syncStatus := types.SyncStatus(uint8(10))

		assert.NotNil(t, syncStatus)
		assert.Equal(t, syncStatus.String(), "unknown")
		assert.Equal(t, uint8(syncStatus), uint8(10))
	})
}
