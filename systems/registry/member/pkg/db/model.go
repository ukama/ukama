/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package db

import (
	"github.com/ukama/ukama/systems/common/roles"
	uuid "github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/gorm"
)

type Member struct {
	gorm.Model
	MemberId    uuid.UUID      `gorm:"primaryKey;type:uuid"`
	UserId      uuid.UUID      `gorm:"uniqueIndex:user_id_idx,where:deleted_at is null;not null;type:uuid"`
	Deactivated bool           `gorm:"default:false"`
	Role        roles.RoleType `gorm:"type:uint;not null"` // Set the default value to Member
}
