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
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/common/store"
	"github.com/ukama/ukama/systems/common/uuid"
)

var (
	ErrLockTimeout         = errors.New("lock_timeout")
	ErrIdempotencyKeyTaken = errors.New("idempotency_key_taken")
)

const (
	lockTimeoutStmt       = "SET LOCAL lock_timeout = '2s'"
	pgErrLockNotAvailable = "55P03"
)

type OperationOrder string

const (
	OrderByResourceKey OperationOrder = "resource_key"
	OrderByLeaseExpiry OperationOrder = "lease_expires_at"
)

type OperationFilter struct {
	IdempotencyKey     string
	GroupId            *uuid.UUID
	ResourceKey        string
	Type               string
	System             string
	Statuses           []OperationStatus
	LeaseExpiresBefore *time.Time
	OrderBy            OperationOrder
	Limit              int
}

type OperationRepo interface {
	Add(ctx context.Context, ops []*Operation, nestedFunc func(tx *gorm.DB) error) error
	Get(ctx context.Context, id uuid.UUID) (*Operation, error)
	List(ctx context.Context, filter OperationFilter) ([]*Operation, error)
	Update(ctx context.Context, id uuid.UUID,
		nestedFunc func(op *Operation, tx *gorm.DB) error) (*Operation, error)
	Delete(ctx context.Context, id uuid.UUID,
		nestedFunc func(op *Operation, tx *gorm.DB) error) error
}

type operationRepo struct {
	db sql.Db
}

func NewOperationRepo(db sql.Db) OperationRepo {
	return &operationRepo{db: db}
}

func (r *operationRepo) Add(ctx context.Context, ops []*Operation,
	nestedFunc func(tx *gorm.DB) error) error {
	err := r.db.GetGormDb().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(lockTimeoutStmt).Error; err != nil {
			return err
		}
		for _, op := range ops {
			if err := tx.Create(op).Error; err != nil {
				if sql.IsDuplicateKeyError(err) {
					return ErrIdempotencyKeyTaken
				}
				return err
			}
		}
		if nestedFunc != nil {
			return nestedFunc(tx)
		}
		return nil
	})

	if errors.Is(err, ErrLockConflict) || errors.Is(err, ErrIdempotencyKeyTaken) {
		return err
	}
	if err != nil {
		return fmt.Errorf("add operations: %w", lockTimeoutErr(err))
	}
	return nil
}

func (r *operationRepo) Get(ctx context.Context, id uuid.UUID) (*Operation, error) {
	var op Operation
	if err := r.db.GetGormDb().WithContext(ctx).Where("id = ?", id).
		First(&op).Error; err != nil {
		return nil, err
	}
	return &op, nil
}

func (r *operationRepo) List(ctx context.Context,
	filter OperationFilter) ([]*Operation, error) {
	q := r.db.GetGormDb().WithContext(ctx)
	if filter.IdempotencyKey != "" {
		q = q.Where("idempotency_key = ?", filter.IdempotencyKey)
	}
	if filter.GroupId != nil {
		q = q.Where("group_id = ?", *filter.GroupId)
	}
	if filter.ResourceKey != "" {
		q = q.Where("resource_key = ?", filter.ResourceKey)
	}
	if filter.Type != "" {
		q = q.Where("type = ?", filter.Type)
	}
	if filter.System != "" {
		q = q.Where("system = ?", filter.System)
	}
	if len(filter.Statuses) > 0 {
		q = q.Where("status IN ?", filter.Statuses)
	}
	if filter.LeaseExpiresBefore != nil {
		q = q.Where("lease_expires_at < ?", *filter.LeaseExpiresBefore)
	}
	switch filter.OrderBy {
	case "":
	case OrderByResourceKey, OrderByLeaseExpiry:
		q = q.Order(string(filter.OrderBy))
	default:
		return nil, fmt.Errorf("list operations: unknown order %q: %w",
			filter.OrderBy, store.ErrInvalid)
	}
	if filter.Limit > 0 {
		q = q.Limit(filter.Limit)
	}

	var ops []*Operation
	if err := q.Find(&ops).Error; err != nil {
		return nil, err
	}
	return ops, nil
}

func (r *operationRepo) Update(ctx context.Context, id uuid.UUID,
	nestedFunc func(op *Operation, tx *gorm.DB) error) (*Operation, error) {
	var op Operation
	err := r.db.GetGormDb().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := lockOperation(tx, &op, id); err != nil {
			return err
		}
		if nestedFunc != nil {
			if err := nestedFunc(&op, tx); err != nil {
				return err
			}
		}
		return tx.Save(&op).Error
	})
	if err != nil {
		return nil, lockTimeoutErr(err)
	}
	return &op, nil
}

func (r *operationRepo) Delete(ctx context.Context, id uuid.UUID,
	nestedFunc func(op *Operation, tx *gorm.DB) error) error {
	err := r.db.GetGormDb().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var op Operation
		if err := lockOperation(tx, &op, id); err != nil {
			return err
		}
		if nestedFunc != nil {
			if err := nestedFunc(&op, tx); err != nil {
				return err
			}
		}
		return tx.Delete(&op).Error
	})
	return lockTimeoutErr(err)
}

func lockOperation(tx *gorm.DB, op *Operation, id uuid.UUID) error {
	if err := tx.Exec(lockTimeoutStmt).Error; err != nil {
		return err
	}
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", id).First(op).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("operation %s: %w", id, store.ErrNotFound)
	}
	return err
}

func lockTimeoutErr(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == pgErrLockNotAvailable {
		return fmt.Errorf("%w: %v", ErrLockTimeout, err)
	}
	return err
}
