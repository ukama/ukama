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

	"github.com/ukama/ukama/systems/common/roles"
	uuid "github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/gorm"
)

type Invitation struct {
	Id        uuid.UUID `gorm:"primaryKey;type:uuid"`
	Link      string
	Email     string
	Name      string
	ExpiresAt time.Time
	Role      roles.RoleType   `gorm:"type:uint;not null;default:4"`
	Status    InvitationStatus `gorm:"type:uint;not null;default:0"`
	UserId    string           `gorm:"type:uuid"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type InvitationStatus uint8

const (
	Pending  InvitationStatus = 0
	Accepted InvitationStatus = 1
	Declined InvitationStatus = 2
)

func (e *InvitationStatus) Scan(value interface{}) error {
	*e = InvitationStatus(uint8(value.(int64)))
	return nil
}

func (e InvitationStatus) Value() (uint8, error) {
	return uint8(e), nil
}
