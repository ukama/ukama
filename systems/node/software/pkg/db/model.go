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

type Node struct {
	NodeId    string         `gorm:"type:string;uniqueIndex:idx_node_node_id,where:deleted_at is null;size:23;not null"`
	NodeType  ukama.NodeType `gorm:"not null"`
	CreatedAt time.Time      `gorm:"not null;default:now()"`
	UpdatedAt time.Time      `gorm:"not null;default:now()"`
	DeletedAt *time.Time     `gorm:"index;default:null"`
}

type Software struct {
	Id             uuid.UUID `gorm:"primaryKey;type:uuid;index"`
	NodeId         string    `gorm:"foreignKey:NodeId;references:NodeId"`
	AppName        string    `gorm:"not null"`
	App            App       `gorm:"foreignKey:AppName;references:Name"`
	ChangeLogs     []string  `gorm:"serializer:json"`
	CurrentVersion string    `gorm:"not null;default:'0.0.1'"`
	DesiredVersion string    `gorm:"not null;default:''"`
	ReleaseDate    time.Time
	CreatedAt      time.Time  `gorm:"not null;default:now()"`
	UpdatedAt      time.Time  `gorm:"not null;default:now()"`
	DeletedAt      *time.Time `gorm:"index;default:null"`
	Status         ukama.SoftwareStatusType
}

// ReleaseCatalog is the node-independent record of what the Hub has published.
// One row per (name, type, version). Populated from Artifact Manager availability
// signals and the periodic Hub reconcile — independent of node discovery.
type ReleaseCatalog struct {
	Id         uuid.UUID `gorm:"primaryKey;type:uuid"`
	Name       string    `gorm:"not null;uniqueIndex:idx_release_name_type_ver,priority:1"`
	Type       string    `gorm:"not null;default:'app';uniqueIndex:idx_release_name_type_ver,priority:2"`
	Version    string    `gorm:"not null;uniqueIndex:idx_release_name_type_ver,priority:3"`
	Digest     string
	SizeBytes  int64
	Chunked    bool
	Available  bool `gorm:"not null;default:true"`
	UploadedAt time.Time
	CreatedAt  time.Time  `gorm:"not null;default:now()"`
	UpdatedAt  time.Time  `gorm:"not null;default:now()"`
	DeletedAt  *time.Time `gorm:"index;default:null"`
}

// AppDesiredRelease is the explicitly-promoted desired version per app, fleet-wide.
// Set only by PromoteRelease — never by upload/chunk events.
type AppDesiredRelease struct {
	Id             uuid.UUID `gorm:"primaryKey;type:uuid"`
	Name           string    `gorm:"not null;uniqueIndex:idx_desired_name_type,priority:1"`
	Type           string    `gorm:"not null;default:'app';uniqueIndex:idx_desired_name_type,priority:2"`
	DesiredVersion string    `gorm:"not null"`
	PromotedAt     time.Time
	PromotedBy     string
	CreatedAt      time.Time `gorm:"not null;default:now()"`
	UpdatedAt      time.Time `gorm:"not null;default:now()"`
}
