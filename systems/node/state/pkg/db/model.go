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

	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"

	"gorm.io/gorm"
)

type State struct {
    Id              uuid.UUID           `gorm:"primaryKey;type:string;size:23;not null"`
    NodeId          string              `gorm:"type:string;size:23;not null"` 
    State           ukama.NodeStateEnum `gorm:"type:uint;not null"`
    LastHeartbeat   time.Time
    LastStateChange time.Time
    Connectivity    ukama.Connectivity `gorm:"type:uint;not null"`
    Type            string             `gorm:"type:string;not null"`
    Version         string             `gorm:"type:string"`
    CreatedAt       time.Time
    UpdatedAt       time.Time
    DeletedAt       gorm.DeletedAt `gorm:"index"`
}
