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
		syncStatus := types.ParseStatus("2")

		assert.NotNil(t, syncStatus)
		assert.Equal(t, uint8(syncStatus), uint8(2))
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
