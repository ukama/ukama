/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package db

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/common/uuid"
)

var ErrLockConflict = errors.New("resource_locked")

const advisoryLockStmt = "SELECT pg_advisory_xact_lock(hashtextextended(?, 0))"

type ConflictError struct {
	ResourceKey string
	HolderId    uuid.UUID
	HolderType  string
}

func (e *ConflictError) Error() string {
	if e.HolderId == uuid.Nil {
		return fmt.Sprintf("%s is locked by an active operation", e.ResourceKey)
	}
	return fmt.Sprintf("%s is busy with %s operation %s", e.ResourceKey,
		e.HolderType, e.HolderId)
}

func (e *ConflictError) Is(target error) bool {
	return target == ErrLockConflict
}

type heldLock struct {
	ResourceKey string
	OperationId uuid.UUID
	Type        string
	Status      OperationStatus
}

type LockRepo interface {
	Acquire(tx *gorm.DB, ops []*Operation, mustBeFree []string) error
	Release(tx *gorm.DB, op *Operation) error
	Holder(ctx context.Context, resourceKey string) (*ResourceLock, error)
}

type lockRepo struct {
	db sql.Db
}

func NewLockRepo(db sql.Db) LockRepo {
	return &lockRepo{db: db}
}

func (r *lockRepo) Acquire(tx *gorm.DB, ops []*Operation, mustBeFree []string) error {
	targets := make([]string, 0, len(ops))
	for _, op := range ops {
		targets = append(targets, op.ResourceKey)
	}

	for _, key := range sortedUnique(targets, mustBeFree) {
		if err := tx.Exec(advisoryLockStmt, key).Error; err != nil {
			return err
		}
	}
	if err := checkFree(tx, targets, mustBeFree); err != nil {
		return err
	}

	now := time.Now().UTC()
	for _, op := range ops {
		if err := tx.Create(&ResourceLock{
			ResourceKey: op.ResourceKey,
			OperationId: op.Id,
			AcquiredAt:  now,
		}).Error; err != nil {
			if sql.IsDuplicateKeyError(err) {
				return &ConflictError{ResourceKey: op.ResourceKey}
			}
			return err
		}
	}
	return nil
}

func (r *lockRepo) Release(tx *gorm.DB, op *Operation) error {
	return tx.Where("resource_key = ? AND operation_id = ?", op.ResourceKey, op.Id).
		Delete(&ResourceLock{}).Error
}

func (r *lockRepo) Holder(ctx context.Context,
	resourceKey string) (*ResourceLock, error) {
	var lock ResourceLock
	err := r.db.GetGormDb().WithContext(ctx).
		Where("resource_key = ?", resourceKey).First(&lock).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &lock, nil
}

func checkFree(tx *gorm.DB, targets, mustBeFree []string) error {
	var held []heldLock
	if err := tx.Table("resource_locks").
		Select("resource_locks.resource_key, resource_locks.operation_id, "+
			"operations.type, operations.status").
		Joins("LEFT JOIN operations ON operations.id = resource_locks.operation_id").
		Where("resource_locks.resource_key IN ?", sortedUnique(targets, mustBeFree)).
		Order("resource_locks.resource_key").
		Scan(&held).Error; err != nil {
		return err
	}

	isTarget := make(map[string]bool, len(targets))
	for _, key := range targets {
		isTarget[key] = true
	}
	for _, h := range held {
		if isTarget[h.ResourceKey] || !h.Status.IsTerminal() {
			return &ConflictError{ResourceKey: h.ResourceKey,
				HolderId: h.OperationId, HolderType: h.Type}
		}
	}
	return nil
}

func sortedUnique(lists ...[]string) []string {
	seen := map[string]bool{}
	out := []string{}
	for _, list := range lists {
		for _, key := range list {
			if !seen[key] {
				seen[key] = true
				out = append(out, key)
			}
		}
	}
	sort.Strings(out)
	return out
}
