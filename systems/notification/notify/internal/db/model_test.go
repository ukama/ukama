/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

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
