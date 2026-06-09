/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package db

import (
	"database/sql/driver"
	"time"

	"github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/gorm"
)

type IntentStatus uint8

const (
	IntentWatching IntentStatus = iota
	IntentCompleted
	IntentFailed
	IntentExpired
	IntentCancelled
)

func (s IntentStatus) String() string {
	return map[IntentStatus]string{
		0: "watching", 1: "completed", 2: "failed", 3: "expired", 4: "cancelled",
	}[s]
}

func (s *IntentStatus) Scan(v interface{}) error {
	if v == nil {
		*s = IntentWatching
		return nil
	}
	*s = IntentStatus(v.(int64))
	return nil
}

func (s IntentStatus) Value() (driver.Value, error) {
	return int64(s), nil
}

func (s IntentStatus) IsTerminal() bool {
	return s == IntentCompleted || s == IntentFailed ||
		s == IntentExpired || s == IntentCancelled
}

type MonitoredIntent struct {
	Id             uuid.UUID      `gorm:"primaryKey;type:uuid" json:"id"`
	OperationId    uuid.UUID      `gorm:"type:uuid;not null;uniqueIndex" json:"operationId"`
	ResourceKey    string         `gorm:"not null;index" json:"resourceKey"`
	ActionType     string         `gorm:"not null;index" json:"actionType"`
	FencingToken   uint64         `gorm:"not null" json:"fencingToken"`
	CompletionRule string         `gorm:"not null" json:"completionRule"`
	Status         IntentStatus   `gorm:"not null;index" json:"status"`
	Deadline       time.Time      `gorm:"not null;index" json:"deadline"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}
