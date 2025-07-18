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
	"gorm.io/gorm"
)

type Org struct {
	Id          uuid.UUID      `gorm:"primaryKey;type:uuid"`
	Name        string         `gorm:"uniqueIndex"`
	Owner       uuid.UUID      `gorm:"type:uuid"`
	Users       []User         `gorm:"many2many:org_users"`
	Country     string         `gorm:"default:us"`
	Currency    string         `gorm:"default:usd"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Certificate string
	Deactivated bool
}

/*
TODO: Check if this is user table is still required if not then instead of calling org we should call member directly from user service
*/
type User struct {
	Id          uint      `gorm:"primaryKey"`
	Uuid        uuid.UUID `gorm:"uniqueIndex:uuid_unique,where:deleted_at is null;not null;type:uuid"`
	Deactivated bool
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
