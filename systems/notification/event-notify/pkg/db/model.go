/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package db

import (
	"database/sql/driver"
	"time"

	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/common/uuid"
)

type Notification struct {
	Id           uuid.UUID `gorm:"primaryKey;type:uuid"`
	Title        string    `gorm:"uniqueIndex"`
	Description  string
	Type         NotificationType  `gorm:"type:uint;not null;default:0"`
	Scope        NotificationScope `gorm:"type:uint;not null;default:4"`
	OrgId        string
	NetworkId    string
	SubscriberId string
	UserId       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

type User struct {
	Id           uuid.UUID `gorm:"primaryKey;type:uuid"`
	OrgId        string
	NetworkId    string
	SubscriberId string
	UserId       string
	Role         RoleType
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

type UserNotification struct {
	Id             uuid.UUID `gorm:"primaryKey;type:uuid"`
	NotificationId uuid.UUID
	UserId         uuid.UUID
	IsRead         bool `gorm:"default:false"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

type RoleType uint8

const (
	Owner    RoleType = 0
	Admin    RoleType = 1
	Employee RoleType = 2
	Vendor   RoleType = 3
	Users    RoleType = 4
)

func (e *RoleType) Scan(value interface{}) error {
	*e = RoleType(uint8(value.(int64)))

	return nil
}

func (e RoleType) Value() (uint8, error) {
	return uint8(e), nil
}

type NotificationType uint8

const (
	INFO     NotificationType = 0
	WARNING  NotificationType = 1
	ERROR    NotificationType = 2
	CRITICAL NotificationType = 3
)

func (l *NotificationType) Scan(value interface{}) error {
	*l = NotificationType(uint8(value.(int64)))
	return nil
}

func (l NotificationType) Value() (driver.Value, error) {
	return int64(l), nil
}

type NotificationScope uint8

const (
	ORG        NotificationScope = 0
	NETWORK    NotificationScope = 1
	SITE       NotificationScope = 2
	SUBSCRIBER NotificationScope = 3
	USER       NotificationScope = 4
)

func (l *NotificationScope) Scan(value interface{}) error {
	*l = NotificationScope(uint8(value.(int64)))
	return nil
}

func (l NotificationScope) Value() (driver.Value, error) {
	return int64(l), nil
}
