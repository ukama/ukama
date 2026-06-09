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
	"github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/gorm"
)

var ErrLockConflict = errors.New("resource_locked")

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
		if err := tx.Where("id = ? AND fencing_token = ?", id, fencingToken).First(&op).Error; err != nil {
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
	err := r.db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", id).First(&op).Error; err != nil {
			return err
		}
		if op.FencingToken != fencingToken {
			return fmt.Errorf("fencing token mismatch (op=%d, given=%d)", op.FencingToken, fencingToken)
		}
		if op.Status.IsTerminal() {
			return nil
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

		audit.Id = uuid.NewV4()
		audit.OperationId = op.Id
		audit.ResourceKey = op.ResourceKey
		audit.At = now
		return tx.Create(&audit).Error
	})
	if err != nil {
		return nil, err
	}
	return &op, nil
}

func (r *operationRepo) FindExpired(now time.Time, limit int) ([]Operation, error) {
	var ops []Operation
	err := r.db.GetGormDb().
		Where("lease_expires_at < ? AND status IN ?", now, []OperationStatus{OperationPending, OperationRunning}).
		Limit(limit).
		Find(&ops).Error
	return ops, err
}
