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
	Name        string `gorm:"not null;primaryKey"`
	Space 	 	string
	Notes       string
	MetricsKeys []string
}

type Software struct {
	Id          	uuid.UUID `gorm:"primaryKey;type:uuid"`
	ReleaseDate 	time.Time `gorm:"not null"`
	NodeId      	string		`gorm:"not null"`
	ChangeLog   	[]string
	CurrentVersion  string
	DesiredVersion  string
	CreatedAt   	time.Time  `gorm:"not null;default:now()"`
	UpdatedAt   	time.Time  `gorm:"not null;default:now()"`
	DeletedAt   	*time.Time `gorm:"index;default:null"`
	App             App 		`gorm:"foreignKey:Name;references:Name"`
	Status 	 		ukama.SoftwareStatusType	`gorm:"not null;default:unknown"`
}
