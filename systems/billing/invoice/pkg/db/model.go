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

	"github.com/ukama/ukama/systems/common/uuid"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Invoice struct {
	Id           uuid.UUID      `gorm:"primaryKey;type:uuid"`
	SubscriberId uuid.UUID      `gorm:"uniqueIndex:subscriber_id_period,where:deleted_at is null;not null;type:uuid"`
	NetworkId    uuid.UUID      `gorm:"not null;type:uuid"`
	Period       time.Time      `gorm:"uniqueIndex:subscriber_id_period,where:deleted_at is null;not null"`
	RawInvoice   datatypes.JSON `gorm:"not null"`
	IsPaid       bool           `gorm:"default:false"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}
