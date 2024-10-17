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
	"encoding/json"
	"time"

	"github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/gorm"
)

type StringArray []string

func (a StringArray) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *StringArray) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &a)
}

type State struct {
	Id              uuid.UUID       `gorm:"primaryKey;type:uuid" json:"id"`
	NodeId          string          `gorm:"not null;index" json:"nodeId"`
	PreviousStateId *uuid.UUID      `gorm:"column:previous_state_id;index" json:"previousStateId,omitempty"`
	PreviousState   *State          `gorm:"-" json:"previousState,omitempty"`
	CurrentState    string `gorm:"not null" json:"currentState"`
	SubState        StringArray     `gorm:"type:jsonb" json:"subState"`
	Events          StringArray     `gorm:"type:jsonb" json:"events"`
	Version         string          `gorm:"" json:"version,omitempty"`
	NodeType        string          `gorm:"" json:"nodeType,omitempty"`
	NodeIp          string          `gorm:"" json:"nodeIp,omitempty"`
	NodePort        int32           `gorm:"" json:"nodePort,omitempty"`
	MeshIp          string          `gorm:"" json:"meshIp,omitempty"`
	MeshPort        int32           `gorm:"" json:"meshPort,omitempty"`
	MeshHostName    string          `gorm:"" json:"meshHostName,omitempty"`
	CreatedAt       time.Time       `json:"createdAt"`
	UpdatedAt       time.Time       `json:"updatedAt"`
	DeletedAt       gorm.DeletedAt  `gorm:"index" json:"deletedAt,omitempty"`
}
