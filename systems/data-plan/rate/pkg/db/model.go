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

type DefaultMarkup struct {
	gorm.Model
	Markup float64 `gorm:"type:float; default:0"`
}

type Markups struct {
	gorm.Model
	OwnerId uuid.UUID `gorm:"uniqueIndex:owner_id_unique,where:deleted_at is null;not null;type:uuid"`
	Markup  float64   `gorm:"type:float; default:0"`
}
