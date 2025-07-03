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

	ukama "github.com/ukama/ukama/systems/common/ukama"
	uuid "github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/gorm"
)

type BaseRate struct {
	gorm.Model
	Uuid        uuid.UUID `gorm:"uniqueIndex:uuid_idx,where:deleted_at is null;not null;type:uuid"`
	Country     string    `gorm:"uniqueIndex:baserate_idx,priority:1,where:deleted_at is null;not null;type:string"`
	Provider    string    `gorm:"uniqueIndex:baserate_idx,priority:4,where:deleted_at is null;not null;type:string"`
	Vpmn        string
	Imsi        int64
	SmsMo       float64 `gorm:"type:float"`
	SmsMt       float64 `gorm:"type:float"`
	Data        float64 `gorm:"type:float"`
	X2g         bool    `gorm:"type:bool; default:false"`
	X3g         bool    `gorm:"type:bool; default:false"`
	X5g         bool    `gorm:"type:bool; default:false"`
	Lte         bool    `gorm:"type:bool; default:false"`
	LteM        bool    `gorm:"type:bool; default:false"`
	Apn         string
	EffectiveAt time.Time `gorm:"uniqueIndex:baserate_idx,priority:3,where:deleted_at is null;not null"`
	EndAt       time.Time
	SimType     ukama.SimType `gorm:"uniqueIndex:baserate_idx,priority:2,where:deleted_at is null;not null"`
	Currency    string        `gorm:"not null; default:Dollar"`
}
