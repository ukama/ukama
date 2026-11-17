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

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/common/store"
	"github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrLockConflict        = errors.New("resource_locked")
	ErrLockTimeout         = errors.New("lock_timeout")
	ErrIdempotencyKeyTaken = errors.New("idempotency_key_taken")
)

const (
	lockTimeoutStmt       = "SET LOCAL lock_timeout = '2s'"
	advisoryLockStmt      = "SELECT pg_advisory_xact_lock(hashtextextended(?, 0))"
	pgErrLockNotAvailable = "55P03"
)

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

type OperationRepo interface {
	Start(ctx context.Context, ops []*Operation, mustBeFree []string,
		lockTTL time.Duration) ([]*Operation, error)
	Get(ctx context.Context, id uuid.UUID) (*Operation, error)
	GetByResource(ctx context.Context, resourceKey string) (*Operation, error)
	GetByIdempotencyKey(ctx context.Context, key string) ([]*Operation, error)
	MarkRunning(ctx context.Context, id uuid.UUID, fencingToken uint64,
		actor string) (*Operation, error)
	Terminate(ctx context.Context, id uuid.UUID, fencingToken uint64,
		status OperationStatus, audit OperationAudit, opErr string) (*Operation, error)
	FindExpired(ctx context.Context, now time.Time, limit int) ([]Operation, error)
	Delete(ctx context.Context, id uuid.UUID, actor string) error
}

type operationRepo struct {
	db sql.Db
}

func NewOperationRepo(db sql.Db) OperationRepo {
	return &operationRepo{db: db}
}

func (r *operationRepo) Start(ctx context.Context, ops []*Operation,
	mustBeFree []string, lockTTL time.Duration) ([]*Operation, error) {
	targets := make([]string, 0, len(ops))
	for _, op := range ops {
		targets = append(targets, op.ResourceKey)
	}

	err := r.db.GetGormDb().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(lockTimeoutStmt).Error; err != nil {
			return err
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
			if err := tx.Create(op).Error; err != nil {
				if sql.IsDuplicateKeyError(err) {
					return ErrIdempotencyKeyTaken
				}
				return err
			}
			if err := tx.Create(&ResourceLock{
				ResourceKey:  op.ResourceKey,
				OperationId:  op.Id,
				FencingToken: op.FencingToken,
				AcquiredAt:   now,
				ExpiresAt:    now.Add(lockTTL),
			}).Error; err != nil {
				if sql.IsDuplicateKeyError(err) {
					return &ConflictError{ResourceKey: op.ResourceKey}
				}
				return err
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

	if errors.Is(err, ErrLockConflict) || errors.Is(err, ErrIdempotencyKeyTaken) {
		return nil, err
	}
	if err != nil {
		return nil, fmt.Errorf("start operation: %w", lockTimeoutErr(err))
	}
	return ops, nil
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

func (r *operationRepo) Get(ctx context.Context, id uuid.UUID) (*Operation, error) {
	var op Operation
	if err := r.db.GetGormDb().WithContext(ctx).Where("id = ?", id).
		First(&op).Error; err != nil {
		return nil, err
	}
	return &op, nil
}

func (r *operationRepo) GetByResource(ctx context.Context,
	resourceKey string) (*Operation, error) {
	var lock ResourceLock
	err := r.db.GetGormDb().WithContext(ctx).
		Where("resource_key = ?", resourceKey).First(&lock).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return r.Get(ctx, lock.OperationId)
}

func (r *operationRepo) GetByIdempotencyKey(ctx context.Context,
	key string) ([]*Operation, error) {
	gdb := r.db.GetGormDb().WithContext(ctx)

	var leader Operation
	err := gdb.Where("idempotency_key = ?", key).First(&leader).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var members []*Operation
	if err := gdb.Where("group_id = ?", leader.Id).Order("resource_key").
		Find(&members).Error; err != nil {
		return nil, err
	}
	return append([]*Operation{&leader}, members...), nil
}

func (r *operationRepo) MarkRunning(ctx context.Context, id uuid.UUID,
	fencingToken uint64, actor string) (*Operation, error) {
	var op Operation
	var rejectErr error

	err := r.db.GetGormDb().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := lockOperation(tx, &op, "id = ?", id); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("operation %s: %w", id, store.ErrNotFound)
			}
			return err
		}

		now := time.Now().UTC()
		writeAudit := func(prefix, reason string, cause error) error {
			rejectErr = cause
			return tx.Create(&OperationAudit{
				Id: uuid.NewV4(), OperationId: op.Id, ResourceKey: op.ResourceKey,
				Event: prefix + "running", Actor: actor, Reason: reason, At: now,
			}).Error
		}
		reject := func(cause error) error {
			return writeAudit("rejected_", cause.Error(), cause)
		}

		switch {
		case op.FencingToken != fencingToken:
			return reject(fmt.Errorf("fencing token mismatch (op=%d, given=%d): %w",
				op.FencingToken, fencingToken, store.ErrPrecondition))
		case op.Status.IsTerminal():
			return reject(fmt.Errorf("operation %s is already %s: %w",
				id, op.Status, store.ErrPrecondition))
		case !op.LeaseExpiresAt.After(now):
			return reject(fmt.Errorf("operation %s lease expired at %s: %w",
				id, op.LeaseExpiresAt.Format(time.RFC3339), store.ErrPrecondition))
		case op.Status == OperationRunning:
			return writeAudit("ignored_", "", nil)
		}

		op.Status = OperationRunning
		op.StartedAt = &now
		if err := tx.Save(&op).Error; err != nil {
			return err
		}
		return writeAudit("", "", nil)
	})
	if err != nil {
		return nil, lockTimeoutErr(err)
	}
	if rejectErr != nil {
		return nil, rejectErr
	}
	return &op, nil
}

func (r *operationRepo) Terminate(ctx context.Context, id uuid.UUID,
	fencingToken uint64, status OperationStatus, audit OperationAudit,
	opErr string) (*Operation, error) {

	if !status.IsTerminal() {
		return nil, fmt.Errorf("terminate: %s is not a terminal status", status)
	}

	var op Operation
	var rejectErr error

	err := r.db.GetGormDb().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
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
		return nil, lockTimeoutErr(err)
	}
	if rejectErr != nil {
		return nil, rejectErr
	}
	return &op, nil
}

func lockTimeoutErr(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == pgErrLockNotAvailable {
		return fmt.Errorf("%w: %v", ErrLockTimeout, err)
	}
	return err
}

func lockOperation(tx *gorm.DB, op *Operation, query string,
	args ...interface{}) error {
	if err := tx.Exec(lockTimeoutStmt).Error; err != nil {
		return err
	}
	return tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where(query, args...).First(op).Error
}

func (r *operationRepo) FindExpired(ctx context.Context, now time.Time,
	limit int) ([]Operation, error) {
	var ops []Operation
	err := r.db.GetGormDb().WithContext(ctx).
		Where("lease_expires_at < ? AND status IN ?", now, []OperationStatus{OperationPending, OperationRunning}).
		Limit(limit).
		Find(&ops).Error
	return ops, err
}

func (r *operationRepo) Delete(ctx context.Context, id uuid.UUID, actor string) error {
	err := r.db.GetGormDb().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var op Operation
		if err := lockOperation(tx, &op, "id = ?", id); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("operation %s: %w", id, store.ErrNotFound)
			}
			return err
		}
		if !op.Status.IsTerminal() {
			return fmt.Errorf("operation %s is %s and cannot be deleted: %w",
				id, op.Status, store.ErrPrecondition)
		}
		if err := tx.Delete(&op).Error; err != nil {
			return err
		}
		return tx.Create(&OperationAudit{
			Id: uuid.NewV4(), OperationId: op.Id, ResourceKey: op.ResourceKey,
			Event: "deleted", Actor: actor, At: time.Now().UTC(),
		}).Error
	})
	return lockTimeoutErr(err)
}
