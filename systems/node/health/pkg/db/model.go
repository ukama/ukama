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

	uuid "github.com/ukama/ukama/systems/common/uuid"
)

type Health struct {
	Id        uuid.UUID  `gorm:"primaryKey;type:uuid"`
	NodeId    string     `gorm:"not null"`
	TimeStamp string     
	System    []System   `gorm:"foreignKey:HealthID"`
	Capps     []Capp     `gorm:"foreignKey:HealthID"`
	CreatedAt time.Time  `gorm:"not null"`
	UpdatedAt time.Time  `gorm:"not null"`
	DeletedAt *time.Time `gorm:"index"`
}
type System struct {
	Id       uuid.UUID `gorm:"primaryKey;type:uuid"`
	HealthID uuid.UUID `gorm:"type:uuid"`
	Name     string    `gorm:"not null"`
	Value    string    `gorm:"not null"`
}

type Capp struct {
	Id        uuid.UUID  `gorm:"primaryKey;type:uuid"`
	HealthID  uuid.UUID  `gorm:"type:uuid"`
	Space    string     
	Name      string     
	Tag       string     
	Status    Status `gorm:"type:uint;not null;default:3"`
	Resources []Resource `gorm:"foreignKey:CappID"`
}

type Resource struct {
	Id     uuid.UUID `gorm:"primaryKey;type:uuid"`
	CappID uuid.UUID `gorm:"type:uuid"` // Foreign key to associate with Capp
	Name   string    `gorm:"not null"`
	Value  string    `gorm:"not null"` // "value" field from the JSON payload
}

type Status uint8

const (
	Pending Status = 0
	Active  Status = 1
	Done    Status = 2
	Unknown Status = 3
)

func (e *Status) Scan(value interface{}) error {
	*e = Status(uint8(value.(int64)))

	return nil
}

func (e Status) Value() (uint8, error) {
	return uint8(e), nil
}
