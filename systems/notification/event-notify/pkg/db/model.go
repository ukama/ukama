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

	"github.com/jackc/pgtype"
	notif "github.com/ukama/ukama/systems/common/notification"
	"github.com/ukama/ukama/systems/common/roles"
	"github.com/ukama/ukama/systems/common/uuid"
)

type JSONB map[string]interface{}

type Notification struct {
	Id           uuid.UUID `gorm:"primaryKey;type:uuid"`
	Title        string
	Description  string
	Type         notif.NotificationType  `gorm:"type:uint;not null;default:0"`
	Scope        notif.NotificationScope `gorm:"type:uint;not null;default:0"`
	ResourceId   uuid.UUID
	OrgId        string
	NetworkId    string
	SubscriberId string
	UserId       string
	NodeId       string
	EventMsgID   uint
	EventMsg     EventMsg
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

type Users struct {
	Id           uuid.UUID `gorm:"primaryKey;type:uuid"`
	OrgId        string
	NetworkId    string
	SubscriberId string
	UserId       string
	Role         roles.RoleType `gorm:"type:uint;not null;default:0"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

type UserNotification struct {
	Id             uuid.UUID `gorm:"primaryKey;type:uuid"`
	NotificationId uuid.UUID `gorm:"type:uuid"`
	UserId         uuid.UUID
	IsRead         bool `gorm:"default:false"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

type Notifications struct {
	Id          uuid.UUID `gorm:"type:uuid"`
	Title       string
	Description string
	Type        notif.NotificationType  `gorm:"type:uint;not null;default:0"`
	Scope       notif.NotificationScope `gorm:"type:uint;not null;default:0"`
	IsRead      bool                    `gorm:"default:false"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type EventMsg struct {
	gorm.Model
	Data pgtype.JSONB `gorm:"type:jsonb;default:'[]';not null"`
}
