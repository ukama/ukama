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

	"github.com/ukama/ukama/systems/common/pb/gen/ukama"
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
	CreatedAt time.Time
	UpdatedAt time.Time
	Role      roles.RoleType         `gorm:"type:uint;not null;default:4"`
	Status    ukama.InvitationStatus `gorm:"type:uint;not null;default:0"`
	UserId    string                 `gorm:"type:uuid"`
	DeletedAt gorm.DeletedAt         `gorm:"index"`
}
