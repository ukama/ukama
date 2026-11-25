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

	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/common/uuid"
)

type OperationStatus uint8

const (
	OperationPending OperationStatus = iota
	OperationRunning
	OperationSuccess
	OperationFailed
	OperationTimeout
	OperationCancelled
)

func (s OperationStatus) String() string {
	return map[OperationStatus]string{
		0: "pending", 1: "running", 2: "success",
		3: "failed", 4: "timeout", 5: "cancelled",
	}[s]
}

func (s *OperationStatus) Scan(v interface{}) error {
	if v == nil {
		*s = OperationPending
		return nil
	}
	*s = OperationStatus(v.(int64))
	return nil
}

func (s OperationStatus) Value() (driver.Value, error) {
	return int64(s), nil
}

func (s OperationStatus) IsTerminal() bool {
	return s == OperationSuccess || s == OperationFailed ||
		s == OperationTimeout || s == OperationCancelled
}

// Operation is one locked action on one resource key.
//
// One StartOperation request creates a group of operations, one per target key.
// The group convention keeps everything in this table:
//   - The leader is the operation on the first target in sorted key order. Its
//     GroupId is nil and it carries the request's IdempotencyKey, if any.
//   - Every other member points GroupId at the leader's Id and has a nil
//     IdempotencyKey, so idx_op_idem stays one key per request.
//   - A single-target request is a leader with no members.
//
// A group is loaded with "id = leader OR group_id = leader". Locks and audits stay
// per operation.
type Operation struct {
	Id             uuid.UUID       `gorm:"primaryKey;type:uuid" json:"id"`
	GroupId        *uuid.UUID      `gorm:"type:uuid;index" json:"groupId,omitempty"`
	Type           string          `gorm:"not null;index" json:"type"`
	System         string          `gorm:"not null;index" json:"system"`
	Status         OperationStatus `gorm:"not null;index" json:"status"`
	FencingToken   uint64          `gorm:"autoIncrement;uniqueIndex" json:"fencingToken"`
	RequestedBy    string          `gorm:"index" json:"requestedBy,omitempty"`
	IdempotencyKey *string         `gorm:"uniqueIndex:idx_op_idem,where:deleted_at is null" json:"idempotencyKey,omitempty"`
	ResourceKey    string          `gorm:"not null;index" json:"resourceKey"`
	LeaseExpiresAt time.Time       `gorm:"not null;index" json:"leaseExpiresAt"`
	Error          string          `gorm:"" json:"error,omitempty"`
	StartedAt      *time.Time      `json:"startedAt,omitempty"`
	TerminalAt     *time.Time      `json:"terminalAt,omitempty"`
	CreatedAt      time.Time       `json:"createdAt"`
	UpdatedAt      time.Time       `json:"updatedAt"`
	DeletedAt      gorm.DeletedAt  `gorm:"index" json:"deletedAt,omitempty"`
}

type ResourceLock struct {
	ResourceKey string     `gorm:"primaryKey" json:"resourceKey"`
	OperationId uuid.UUID  `gorm:"type:uuid;not null;index" json:"operationId"`
	Operation   *Operation `gorm:"foreignKey:OperationId" json:"operation,omitempty"`
	AcquiredAt  time.Time  `gorm:"not null" json:"acquiredAt"`
}

type OperationAudit struct {
	Id          uuid.UUID `gorm:"primaryKey;type:uuid" json:"id"`
	OperationId uuid.UUID `gorm:"type:uuid;index" json:"operationId"`
	ResourceKey string    `gorm:"index" json:"resourceKey,omitempty"`
	Event       string    `gorm:"not null" json:"event"`
	Actor       string    `gorm:"" json:"actor,omitempty"`
	Reason      string    `gorm:"" json:"reason,omitempty"`
	At          time.Time `gorm:"not null;index" json:"at"`
}
