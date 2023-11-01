/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package db

import (
	"time"

	uuid "github.com/ukama/ukama/systems/common/uuid"
)

type Mailing struct {
	MailId       uuid.UUID  `gorm:"primaryKey;type:uuid"`
	Email        string     `gorm:"size:255"`
	TemplateName string     `gorm:"size:255"`
	SentAt       *time.Time `gorm:"index"`
	Status       string     `gorm:"not null"`
	CreatedAt    time.Time  `gorm:"not null"`
	UpdatedAt    time.Time  `gorm:"not null"`
	DeletedAt    *time.Time `gorm:"index"`
}
