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
	uuid "github.com/ukama/ukama/systems/common/uuid"
)

type App struct {
	Id          uuid.UUID `gorm:"primaryKey;type:uuid;index"`
	Name        string    `gorm:"not null;index:idx_app_name,unique"`
	Space       string
	Notes       string
	MetricsKeys []string `gorm:"serializer:json"`
}

type Software struct {
	Id             uuid.UUID `gorm:"primaryKey;type:uuid;index"`
	NodeId         string    `gorm:"not null;uniqueIndex:idx_software_node_app"`
	AppName        string    `gorm:"not null;uniqueIndex:idx_software_node_app"`
	App            App       `gorm:"foreignKey:AppName;references:Name"`
	ChangeLogs     []string  `gorm:"serializer:json"`
	CurrentVersion string    `gorm:"not null;default:'0.0.1'"`
	DesiredVersion string    `gorm:"not null;default:''"`
	ReleaseDate    time.Time
	CreatedAt      time.Time                `gorm:"not null;default:now()"`
	UpdatedAt      time.Time                `gorm:"not null;default:now()"`
	DeletedAt      *time.Time               `gorm:"index;default:null"`
	Status         ukama.SoftwareStatusType
}
