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

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
)

type Report struct {
	Id            uuid.UUID        `gorm:"primaryKey;type:uuid"`
	OwnerId       uuid.UUID        `gorm:"uniqueIndex:report_id_period,where:deleted_at is null;not null;type:uuid"`
	OwnerType     ukama.OwnerType  `gorm:"not null"`
	NetworkId     uuid.UUID        `gorm:"type:uuid"`
	Type          ukama.ReportType `gorm:"not null"`
	Period        time.Time        `gorm:"uniqueIndex:report_id_period,where:deleted_at is null;not null"`
	RawReport     datatypes.JSON   `gorm:"not null"`
	IsPaid        bool             `gorm:"default:false"`
	TransactionId string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}
