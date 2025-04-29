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

	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
)

type Sim struct {
	Id                 uuid.UUID `gorm:"primaryKey;type:uuid"`
	SubscriberId       uuid.UUID `gorm:"not null;type:uuid"`
	NetworkId          uuid.UUID `gorm:"not null;type:uuid"`
	Package            Package
	Iccid              string `gorm:"index:idx_iccid,unique"`
	Msisdn             string
	Imsi               string
	Type               ukama.SimType
	Status             ukama.SimStatus
	IsPhysical         bool
	ActivationsCount   uint64 `gorm:"default:0"`
	DeactivationsCount uint64 `gorm:"default:0"`
	FirstActivatedOn   time.Time
	LastActivatedOn    time.Time
	TrafficPolicy      uint32
	AllocatedAt        int64 `gorm:"autoCreateTime"`
	UpdatedAt          time.Time
	TerminatedAt       time.Time
	DetetedAt          gorm.DeletedAt `gorm:"index"`
	SyncStatus         ukama.StatusType
}

type Package struct {
	Id              uuid.UUID `gorm:"primaryKey;type:uuid"`
	SimId           uuid.UUID `gorm:"uniqueIndex:unique_sim_package_is_active,where:is_active is true;not null;type:uuid"`
	StartDate       time.Time
	EndDate         time.Time
	DefaultDuration uint64
	PackageId       uuid.UUID `gorm:"not null;type:uuid"`
	IsActive        bool      `gorm:"uniqueIndex:unique_sim_package_is_active,where:is_active is true;default:false"`
	AsExpired       bool      `gorm:"default:false"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
