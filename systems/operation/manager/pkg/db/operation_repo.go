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
	"hash/fnv"
	"sort"
	"time"

	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrLockConflict = errors.New("resource_locked")

type OperationRepo interface {
	Start(op *Operation, lockTTL time.Duration) (*Operation, error)
	StartBatch(ops []*Operation, conflictKeys []string, lockTTL time.Duration) ([]*Operation, *Operation, error)
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
	ops, holder, err := r.StartBatch([]*Operation{op}, nil, lockTTL)
	if err != nil {
		return holder, err
	}
	return ops[0], nil
}

// StartBatch checks every conflict key and reserves every target together.
// The existing per-node locks remain authoritative until each operation ends.
func (r *operationRepo) StartBatch(ops []*Operation, conflictKeys []string, lockTTL time.Duration) ([]*Operation, *Operation, error) {
	if len(ops) == 0 {
		return nil, nil, fmt.Errorf("no operation targets")
	}
	keys := make(map[string]bool, len(conflictKeys)+len(ops))
	for _, key := range conflictKeys {
		if key == "" {
			return nil, nil, fmt.Errorf("empty conflict resource key")
		}
		keys[key] = true
	}
	targets := make(map[string]bool, len(ops))
	for _, op := range ops {
		if op == nil || op.ResourceKey == "" || targets[op.ResourceKey] {
			return nil, nil, fmt.Errorf("invalid or duplicate operation target")
		}
		targets[op.ResourceKey] = true
		keys[op.ResourceKey] = true
	}
	ordered := make([]string, 0, len(keys))
	for key := range keys {
		ordered = append(ordered, key)
	}
	sort.Strings(ordered)

	var holdingOp *Operation
	err := r.db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		// Lock admission even when no ResourceLock row exists yet. All Start
		// calls use this path; deterministic order prevents lock-order cycles.
		for _, key := range ordered {
			h := fnv.New64a()
			_, _ = h.Write([]byte("ukama-operation:" + key))
			if err := tx.Exec("SELECT pg_advisory_xact_lock(?)", int64(h.Sum64())).Error; err != nil {
				return err
			}
		}
		var existing ResourceLock
		err := tx.Where("resource_key IN ?", ordered).Order("resource_key").First(&existing).Error
		if err == nil {
			var holder Operation
			if err := tx.Where("id = ?", existing.OperationId).First(&holder).Error; err != nil {
				return err
			}
			holdingOp = &holder
			return ErrLockConflict
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		for _, op := range ops {
			if err := tx.Create(op).Error; err != nil {
				return err
			}
			now := time.Now().UTC()
			lock := &ResourceLock{
				ResourceKey: op.ResourceKey, OperationId: op.Id,
				FencingToken: op.FencingToken, AcquiredAt: now,
				ExpiresAt: now.Add(lockTTL),
			}
			result := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(lock)
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 0 {
				return ErrLockConflict
			}
			if err := tx.Create(&OperationAudit{
				Id: uuid.NewV4(), OperationId: op.Id, ResourceKey: op.ResourceKey,
				Event: "lock_acquired", Actor: op.RequestedBy, At: now,
			}).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if errors.Is(err, ErrLockConflict) {
		return nil, holdingOp, ErrLockConflict
	}
	if err != nil {
		return nil, nil, fmt.Errorf("start operation: %w", err)
	}
	return ops, nil, nil
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
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ? AND fencing_token = ?", id, fencingToken).First(&op).Error; err != nil {
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
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", id).First(&op).Error; err != nil {
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
