/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package db

import (
	"errors"
	"fmt"
	"time"

	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/common/store"
	"github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrLockConflict = errors.New("resource_locked")

const lockTimeoutStmt = "SET LOCAL lock_timeout = '5s'"

type OperationRepo interface {
	Start(op *Operation, lockTTL time.Duration) (*Operation, error)
	Get(id uuid.UUID) (*Operation, error)
	GetByResource(resourceKey string) (*Operation, error)
	GetByIdempotencyKey(key string) (*Operation, error)
	MarkRunning(id uuid.UUID, fencingToken uint64) (*Operation, error)
	Terminate(id uuid.UUID, fencingToken uint64, status OperationStatus,
		audit OperationAudit, opErr string) (*Operation, error)
	FindExpired(now time.Time, limit int) ([]Operation, error)
}

type operationRepo struct {
	db sql.Db
}

func NewOperationRepo(db sql.Db) OperationRepo {
	return &operationRepo{db: db}
}

func (r *operationRepo) Start(op *Operation, lockTTL time.Duration) (*Operation, error) {
	var holdingOp *Operation
	err := r.db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(op).Error; err != nil {
			return err
		}

		lock := &ResourceLock{
			ResourceKey:  op.ResourceKey,
			OperationId:  op.Id,
			FencingToken: op.FencingToken,
			AcquiredAt:   time.Now().UTC(),
			ExpiresAt:    time.Now().UTC().Add(lockTTL),
		}
		if err := tx.Create(lock).Error; err != nil {
			var existing ResourceLock
			if findErr := tx.Where("resource_key = ?", op.ResourceKey).First(&existing).Error; findErr == nil {
				var loadOp Operation
				if loadErr := tx.Where("id = ?", existing.OperationId).First(&loadOp).Error; loadErr == nil {
					holdingOp = &loadOp
				}
			}
			return ErrLockConflict
		}

		return tx.Create(&OperationAudit{
			Id:          uuid.NewV4(),
			OperationId: op.Id,
			ResourceKey: op.ResourceKey,
			Event:       "lock_acquired",
			Actor:       op.RequestedBy,
			At:          time.Now().UTC(),
		}).Error
	})

	if errors.Is(err, ErrLockConflict) {
		return holdingOp, ErrLockConflict
	}
	if err != nil {
		return nil, fmt.Errorf("start operation: %w", err)
	}
	return op, nil
}

func (r *operationRepo) Get(id uuid.UUID) (*Operation, error) {
	var op Operation
	if err := r.db.GetGormDb().Where("id = ?", id).First(&op).Error; err != nil {
		return nil, err
	}
	return &op, nil
}

func (r *operationRepo) GetByResource(resourceKey string) (*Operation, error) {
	var lock ResourceLock
	err := r.db.GetGormDb().Where("resource_key = ?", resourceKey).First(&lock).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return r.Get(lock.OperationId)
}

func (r *operationRepo) GetByIdempotencyKey(key string) (*Operation, error) {
	var op Operation
	err := r.db.GetGormDb().Where("idempotency_key = ?", key).First(&op).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &op, nil
}

func (r *operationRepo) MarkRunning(id uuid.UUID, fencingToken uint64) (*Operation, error) {
	var op Operation
	err := r.db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		if err := lockOperation(tx, &op, "id = ? AND fencing_token = ?", id,
			fencingToken); err != nil {
			return err
		}
		if op.Status != OperationPending {
			return fmt.Errorf("operation %s not in pending state (was %s)", id, op.Status)
		}
		now := time.Now().UTC()
		op.Status = OperationRunning
		op.StartedAt = &now
		if err := tx.Save(&op).Error; err != nil {
			return err
		}
		return tx.Create(&OperationAudit{
			Id: uuid.NewV4(), OperationId: op.Id, ResourceKey: op.ResourceKey,
			Event: "running", At: now,
		}).Error
	})
	if err != nil {
		return nil, err
	}
	return &op, nil
}

func (r *operationRepo) Terminate(id uuid.UUID, fencingToken uint64,
	status OperationStatus, audit OperationAudit, opErr string) (*Operation, error) {

	if !status.IsTerminal() {
		return nil, fmt.Errorf("terminate: %s is not a terminal status", status)
	}

	var op Operation
	var rejectErr error

	err := r.db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		if err := lockOperation(tx, &op, "id = ?", id); err != nil {
			return err
		}

		writeAudit := func(prefix string, cause error) error {
			rejectErr = cause
			audit.Id = uuid.NewV4()
			audit.OperationId = op.Id
			audit.ResourceKey = op.ResourceKey
			audit.Event = prefix + audit.Event
			audit.At = time.Now().UTC()
			return tx.Create(&audit).Error
		}

		if op.FencingToken != fencingToken {
			return writeAudit("rejected_", fmt.Errorf(
				"fencing token mismatch (op=%d, given=%d): %w",
				op.FencingToken, fencingToken, store.ErrPrecondition))
		}
		if op.Status.IsTerminal() {
			var cause error
			if op.Status != status {
				cause = fmt.Errorf("operation %s is already %s, cannot move to %s: %w",
					id, op.Status, status, store.ErrPrecondition)
			}
			return writeAudit("ignored_", cause)
		}
		if status == OperationSuccess && op.Status == OperationPending {
			return writeAudit("rejected_", fmt.Errorf(
				"operation %s is pending, cannot move to %s: %w",
				id, status, store.ErrPrecondition))
		}

		now := time.Now().UTC()
		op.Status = status
		op.TerminalAt = &now
		if opErr != "" {
			op.Error = opErr
		}
		if err := tx.Save(&op).Error; err != nil {
			return err
		}

		if err := tx.Where("resource_key = ? AND operation_id = ?",
			op.ResourceKey, op.Id).Delete(&ResourceLock{}).Error; err != nil {
			return err
		}

		return writeAudit("", nil)
	})
	if err != nil {
		return nil, err
	}
	if rejectErr != nil {
		return nil, rejectErr
	}
	return &op, nil
}

func lockOperation(tx *gorm.DB, op *Operation, query string,
	args ...interface{}) error {
	if err := tx.Exec(lockTimeoutStmt).Error; err != nil {
		return err
	}
	return tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where(query, args...).First(op).Error
}

func (r *operationRepo) FindExpired(now time.Time, limit int) ([]Operation, error) {
	var ops []Operation
	err := r.db.GetGormDb().
		Where("lease_expires_at < ? AND status IN ?", now, []OperationStatus{OperationPending, OperationRunning}).
		Limit(limit).
		Find(&ops).Error
	return ops, err
}
