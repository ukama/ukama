package db_test

import (
	"testing"

	"github.com/tj/assert"
	"github.com/ukama/ukama/systems/api/api-gateway/pkg/db"
)

func TestResourceStatus(t *testing.T) {
	t.Run("ResourceStatusValidString", func(tt *testing.T) {
		resourceStatus := db.ParseStatus("pending")

		assert.NotNil(t, resourceStatus)
		assert.Equal(t, resourceStatus.String(), "pending")
		assert.Equal(t, uint8(resourceStatus), uint8(1))
	})

	t.Run("ResourceStatusValidNumber", func(tt *testing.T) {
		resourceStatus := db.ParseStatus("2")

		assert.NotNil(t, resourceStatus)
		assert.Equal(t, uint8(resourceStatus), uint8(2))
		assert.Equal(t, resourceStatus.String(), "completed")
	})

	t.Run("ResourceStatusNonValidString", func(tt *testing.T) {
		resourceStatus := db.ParseStatus("failure")

		assert.NotNil(t, resourceStatus)
		assert.Equal(t, resourceStatus.String(), "unknown")
		assert.Equal(t, uint8(resourceStatus), uint8(0))
	})

	t.Run("ResourceStatusNonValidNumber", func(tt *testing.T) {
		resourceStatus := db.ResourceStatus(uint8(10))

		assert.NotNil(t, resourceStatus)
		assert.Equal(t, resourceStatus.String(), "unknown")
		assert.Equal(t, uint8(resourceStatus), uint8(10))
	})
}
