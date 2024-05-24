/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package db

import (
	uuid "github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/gorm"
)

type Member struct {
	gorm.Model
	MemberId    uuid.UUID `gorm:"primaryKey;type:uuid"`
	UserId      uuid.UUID `gorm:"uniqueIndex:user_id_idx,where:deleted_at is null;not null;type:uuid"`
	Deactivated bool      `gorm:"default:false"`
	Role        RoleType  `gorm:"type:uint;not null"` // Set the default value to Member
}

type RoleType uint8

const (
	Owner    RoleType = 0
	Admin    RoleType = 1
	Employee RoleType = 2
	Vendor   RoleType = 3
	Users    RoleType = 4
)

func (e *RoleType) Scan(value interface{}) error {
	*e = RoleType(uint8(value.(int64)))

	return nil
}

func (e RoleType) Value() (uint8, error) {
	return uint8(e), nil
}
