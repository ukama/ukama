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
	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
)

type Network struct {
	Id               uuid.UUID      `gorm:"primaryKey;type:uuid"`
	Name             string         `gorm:"uniqueIndex:network_name"`
	Deactivated      bool           `gorm:"default:false"`
	AllowedCountries pq.StringArray `gorm:"type:varchar(64)[]" json:"allowed_countries"`
	AllowedNetworks  pq.StringArray `gorm:"type:varchar(64)[]" json:"allowed_networks"`
	Budget           float64
	Overdraft        float64
	TrafficPolicy    uint32
	PaymentLinks     bool
	IsDefault        bool `gorm:"default:false"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt `gorm:"index"`
	SyncStatus       ukama.StatusType
}
