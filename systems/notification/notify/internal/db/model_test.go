package db_test

import (
	"testing"

	"github.com/tj/assert"
	"github.com/ukama/ukama/systems/notification/notify/internal/db"
)

func TestSeverityType(t *testing.T) {
	t.Run("SeverityTypeValid", func(tt *testing.T) {
		severityType, err := db.GetSeverityType("high")

		assert.NoError(t, err)
		assert.Equal(t, severityType.String(), "high")
	})

	t.Run("SeverityTypeNotValid", func(tt *testing.T) {
		severityType, err := db.GetSeverityType("Super-high")

		assert.Error(t, err)
		assert.Nil(t, severityType)
	})
}

func TestNotificationType(t *testing.T) {
	t.Run("NotificationTypeValid", func(tt *testing.T) {
		notificationType, err := db.GetNotificationType("event")

		assert.NoError(t, err)
		assert.Equal(t, notificationType.String(), "event")
	})

	t.Run("NotificationTypeNotValid", func(tt *testing.T) {
		notificationType, err := db.GetNotificationType("failure")

		assert.Error(t, err)
		assert.Nil(t, notificationType)
	})
}
