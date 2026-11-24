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
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/tj/assert"
	"github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/gorm"
)

func setupLockRepo(t *testing.T) (sqlmock.Sqlmock, LockRepo, *gorm.DB) {
	mock, udb := setupTestDB(t)
	mock.ExpectBegin()
	return mock, NewLockRepo(udb), udb.GetGormDb().Begin()
}

func Test_LockAcquire(t *testing.T) {
	const (
		tower = "node:uk-sa2450-tnode-v0-4e86"
		amp   = "node:uk-sa2450-anode-v0-4e86"
		ctl   = "node:uk-sa2450-cnode-v0-4e86"
	)
	expectGuards := func(mock sqlmock.Sqlmock, keys ...string) {
		for _, key := range keys {
			mock.ExpectExec(regexp.QuoteMeta("SELECT pg_advisory_xact_lock")).
				WithArgs(key).WillReturnResult(sqlmock.NewResult(0, 0))
		}
	}
	expectHeld := func(mock sqlmock.Sqlmock, rows ...heldLock) {
		out := sqlmock.NewRows([]string{"resource_key", "operation_id", "type", "status"})
		for _, r := range rows {
			out.AddRow(r.ResourceKey, r.OperationId, r.Type, r.Status)
		}
		mock.ExpectQuery(`SELECT resource_locks.resource_key.* FROM "resource_locks"`).
			WillReturnRows(out)
	}
	expectLockInsert := func(mock sqlmock.Sqlmock, op *Operation) {
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "resource_locks"`)).
			WithArgs(op.ResourceKey, op.Id, op.FencingToken, sqlmock.AnyArg(),
				sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))
	}

	t.Run("SingleTargetAcquiresLock", func(t *testing.T) {
		mock, repo, tx := setupLockRepo(t)
		op := newOp(tower)

		expectGuards(mock, tower)
		expectHeld(mock)
		expectLockInsert(mock, op)

		err := repo.Acquire(tx, []*Operation{op}, nil, 5*time.Minute)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GroupLocksEveryTargetAfterGuardingAllKeysInOrder", func(t *testing.T) {
		mock, repo, tx := setupLockRepo(t)
		first, second := newOp(amp), newOp(tower)

		expectGuards(mock, amp, ctl, tower)
		expectHeld(mock)
		expectLockInsert(mock, first)
		expectLockInsert(mock, second)

		err := repo.Acquire(tx, []*Operation{first, second}, []string{ctl, tower},
			5*time.Minute)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("FinishedHolderOnSiblingIsFree", func(t *testing.T) {
		mock, repo, tx := setupLockRepo(t)
		op := newOp(tower)

		expectGuards(mock, ctl, tower)
		expectHeld(mock, heldLock{ResourceKey: ctl, OperationId: uuid.NewV4(),
			Type: "RestartNode", Status: OperationSuccess})
		expectLockInsert(mock, op)

		err := repo.Acquire(tx, []*Operation{op}, []string{ctl}, 5*time.Minute)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	for name, held := range map[string]heldLock{
		"BusyTargetIsConflict": {ResourceKey: tower, Status: OperationRunning},
		"FinishedHolderOnTargetIsConflict": {
			ResourceKey: tower, Status: OperationSuccess,
		},
		"BusySiblingIsConflict": {ResourceKey: ctl, Status: OperationPending},
	} {
		t.Run(name, func(t *testing.T) {
			mock, repo, tx := setupLockRepo(t)
			held.OperationId = uuid.NewV4()
			held.Type = "UpdateSoftware"

			expectGuards(mock, ctl, tower)
			expectHeld(mock, held)

			err := repo.Acquire(tx, []*Operation{newOp(tower)}, []string{ctl},
				5*time.Minute)

			var conflict *ConflictError
			assert.True(t, errors.As(err, &conflict), err)
			assert.True(t, errors.Is(err, ErrLockConflict), err)
			assert.Equal(t, held.ResourceKey, conflict.ResourceKey)
			assert.Equal(t, held.OperationId, conflict.HolderId)
			assert.Contains(t, err.Error(), "UpdateSoftware")
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}

	t.Run("LockRowRaceIsConflict", func(t *testing.T) {
		mock, repo, tx := setupLockRepo(t)
		op := newOp(tower)

		expectGuards(mock, tower)
		expectHeld(mock)
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "resource_locks"`)).
			WillReturnError(&pgconn.PgError{Code: "23505"})

		err := repo.Acquire(tx, []*Operation{op}, nil, 5*time.Minute)

		var conflict *ConflictError
		assert.True(t, errors.As(err, &conflict), err)
		assert.Equal(t, uuid.Nil, conflict.HolderId)
		assert.Contains(t, err.Error(), "locked by an active operation")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("OtherLockInsertErrorIsReturned", func(t *testing.T) {
		mock, repo, tx := setupLockRepo(t)

		expectGuards(mock, tower)
		expectHeld(mock)
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "resource_locks"`)).
			WillReturnError(errors.New("connection reset"))

		err := repo.Acquire(tx, []*Operation{newOp(tower)}, nil, 5*time.Minute)

		assert.Error(t, err)
		assert.False(t, errors.Is(err, ErrLockConflict), err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GuardFailureStops", func(t *testing.T) {
		mock, repo, tx := setupLockRepo(t)

		mock.ExpectExec(regexp.QuoteMeta("SELECT pg_advisory_xact_lock")).
			WillReturnError(&pgconn.PgError{Code: pgErrLockNotAvailable})

		err := repo.Acquire(tx, []*Operation{newOp(tower)}, nil, 5*time.Minute)

		assert.True(t, errors.Is(lockTimeoutErr(err), ErrLockTimeout), err)
		assert.False(t, errors.Is(err, ErrLockConflict), err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("HeldLookupFailure", func(t *testing.T) {
		mock, repo, tx := setupLockRepo(t)

		expectGuards(mock, tower)
		mock.ExpectQuery(`SELECT resource_locks.resource_key.* FROM "resource_locks"`).
			WillReturnError(errors.New("connection reset"))

		err := repo.Acquire(tx, []*Operation{newOp(tower)}, nil, 5*time.Minute)

		assert.Error(t, err)
		assert.False(t, errors.Is(err, ErrLockConflict), err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func Test_LockRelease(t *testing.T) {
	t.Run("DeletesTheOperationsLock", func(t *testing.T) {
		mock, repo, tx := setupLockRepo(t)
		op := newOp("node:a")

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "resource_locks"`)).
			WithArgs(op.ResourceKey, op.Id).
			WillReturnResult(sqlmock.NewResult(0, 1))

		assert.NoError(t, repo.Release(tx, op))
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ErrorIsReturned", func(t *testing.T) {
		mock, repo, tx := setupLockRepo(t)

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "resource_locks"`)).
			WillReturnError(errors.New("connection reset"))

		assert.Error(t, repo.Release(tx, newOp("node:a")))
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func Test_LockHolder(t *testing.T) {
	t.Run("ReturnsNilWhenNotLocked", func(t *testing.T) {
		mock, udb := setupTestDB(t)
		repo := NewLockRepo(udb)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "resource_locks"`)).
			WillReturnError(gorm.ErrRecordNotFound)

		out, err := repo.Holder(ctx, "node:free")

		assert.NoError(t, err)
		assert.Nil(t, out)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ReturnsTheLock", func(t *testing.T) {
		mock, udb := setupTestDB(t)
		repo := NewLockRepo(udb)
		holder := uuid.NewV4()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "resource_locks"`)).
			WithArgs("node:busy", sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{
				"resource_key", "operation_id", "fencing_token",
			}).AddRow("node:busy", holder, 4))

		out, err := repo.Holder(ctx, "node:busy")

		assert.NoError(t, err)
		assert.Equal(t, holder, out.OperationId)
		assert.Equal(t, uint64(4), out.FencingToken)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ErrorIsReturned", func(t *testing.T) {
		mock, udb := setupTestDB(t)
		repo := NewLockRepo(udb)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "resource_locks"`)).
			WillReturnError(errors.New("connection reset"))

		out, err := repo.Holder(ctx, "node:a")

		assert.Error(t, err)
		assert.Nil(t, out)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
