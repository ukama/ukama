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

	"github.com/lib/pq"
	"github.com/ukama/ukama/systems/common/ukama"

	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/common/uuid"
)

type Package struct {
	gorm.Model
	Uuid           uuid.UUID `gorm:"unique;type:uuid;index"`
	OwnerId        uuid.UUID
	Name           string
	SimType        ukama.SimType
	Active         bool           `gorm:"not null; default:false"`
	Duration       uint64         `gorm:"not null; default:0"`
	SmsVolume      uint64         `gorm:"not null; default:0"`
	DataVolume     uint64         `gorm:"not null; default:0"`
	VoiceVolume    uint64         `gorm:"not null; default:0"`
	PackageRate    PackageRate    `gorm:"foreignKey:PackageID;references:Uuid"`
	PackageMarkup  PackageMarkup  `gorm:"foreignKey:PackageID;references:Uuid"`
	PackageDetails PackageDetails `gorm:"foreignKey:PackageID;references:Uuid"`
	Type           ukama.PackageType
	DataUnits      ukama.DataUnitType
	VoiceUnits     ukama.CallUnitType
	MessageUnits   ukama.MessageUnitType
	Flatrate       bool      `gorm:"not null; default:false"`
	Currency       string    `gorm:"not null; default:Dollar"`
	From           time.Time `gorm:"not null"`
	To             time.Time `gorm:"not null"`
	Country        string    `gorm:"not null;type:string"`
	Provider       string    `gorm:"not null;type:string"`
	Overdraft      float64
	TrafficPolicy  uint32
	Networks       pq.StringArray `gorm:"type:varchar(64)[]" json:"networks"`
	SyncStatus     ukama.StatusType
}

type PackageDetails struct {
	gorm.Model
	PackageID uuid.UUID
	Dlbr      uint64 `gorm:"not null; default:10240000"`
	Ulbr      uint64 `gorm:"not null; default:10240000"`
	Apn       string
}

type PackageRate struct {
	gorm.Model
	PackageID uuid.UUID
	Amount    float64 `gorm:"type:float"`
	SmsMo     float64 `gorm:"type:float"`
	SmsMt     float64 `gorm:"type:float"`
	Data      float64 `gorm:"type:float"`
}

/* View only for owners */
type PackageMarkup struct {
	gorm.Model
	PackageID  uuid.UUID
	BaseRateId uuid.UUID `gorm:"not null;type:uuid"`
	Markup     float64   `gorm:"type:float"`
}
