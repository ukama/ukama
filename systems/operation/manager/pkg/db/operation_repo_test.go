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
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jackc/pgx/v5/pgconn"
	log "github.com/sirupsen/logrus"
	"github.com/tj/assert"
	"github.com/ukama/ukama/systems/common/store"
	"github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type UkamaDbMock struct {
	GormDb *gorm.DB
}

func (u UkamaDbMock) Init(model ...interface{}) error {
	panic("implement me: Init()")
}

func (u UkamaDbMock) Connect() error {
	panic("implement me: Connect()")
}

func (u UkamaDbMock) GetGormDb() *gorm.DB {
	return u.GormDb
}

func (u UkamaDbMock) InitDB() error {
	return nil
}

func (u UkamaDbMock) ExecuteInTransaction(dbOperation func(tx *gorm.DB) *gorm.DB,
	nestedFuncs ...func() error) error {
	log.Fatal("implement me: ExecuteInTransaction()")
	return nil
}

func (u UkamaDbMock) ExecuteInTransaction2(dbOperation func(tx *gorm.DB) *gorm.DB,
	nestedFuncs ...func(tx *gorm.DB) error) error {
	log.Fatal("implement me: ExecuteInTransaction2()")
	return nil
}

func setupTestDB(t *testing.T) (sqlmock.Sqlmock, OperationRepo) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})

	gdb, err := gorm.Open(dialector, &gorm.Config{})
	assert.NoError(t, err)

	return mock, NewOperationRepo(&UkamaDbMock{GormDb: gdb})
}

func operationRow(op *Operation) *sqlmock.Rows {
	return sqlmock.NewRows([]string{
		"id", "type", "system", "status", "fencing_token", "requested_by",
		"resource_key", "lease_expires_at", "created_at",
	}).AddRow(op.Id, op.Type, op.System, op.Status, op.FencingToken,
		op.RequestedBy, op.ResourceKey, op.LeaseExpiresAt, time.Now())
}

func Test_Start(t *testing.T) {
	const (
		tower = "node:uk-sa2450-tnode-v0-4e86"
		amp   = "node:uk-sa2450-anode-v0-4e86"
		ctl   = "node:uk-sa2450-cnode-v0-4e86"
	)
	newOp := func(key string) *Operation {
		return &Operation{
			Id:             uuid.NewV4(),
			Type:           "RestartNode",
			System:         "node",
			Status:         OperationPending,
			ResourceKey:    key,
			LeaseExpiresAt: time.Now().Add(5 * time.Minute),
		}
	}
	expectGuards := func(mock sqlmock.Sqlmock, keys ...string) {
		expectLockTimeout(mock)
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
	expectInserts := func(mock sqlmock.Sqlmock, op *Operation, token int) {
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "operations"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "fencing_token"}).
				AddRow(op.Id, token))
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "resource_locks"`)).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "operation_audits"`)).
			WithArgs(sqlmock.AnyArg(), op.Id, op.ResourceKey, "lock_acquired",
				sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))
	}

	t.Run("SingleTargetAcquiresLock", func(t *testing.T) {
		mock, repo := setupTestDB(t)
		op := newOp(tower)

		mock.ExpectBegin()
		expectGuards(mock, tower)
		expectHeld(mock)
		expectInserts(mock, op, 1)
		mock.ExpectCommit()

		out, err := repo.Start(ctx, []*Operation{op}, nil, 5*time.Minute)

		assert.NoError(t, err)
		assert.Len(t, out, 1)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GroupLocksEveryTargetAfterGuardingAllKeysInOrder", func(t *testing.T) {
		mock, repo := setupTestDB(t)
		first, second := newOp(amp), newOp(tower)

		mock.ExpectBegin()
		expectGuards(mock, amp, ctl, tower)
		expectHeld(mock)
		expectInserts(mock, first, 1)
		expectInserts(mock, second, 2)
		mock.ExpectCommit()

		out, err := repo.Start(ctx, []*Operation{first, second}, []string{ctl, tower},
			5*time.Minute)

		assert.NoError(t, err)
		assert.Len(t, out, 2)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("FinishedHolderOnSiblingIsFree", func(t *testing.T) {
		mock, repo := setupTestDB(t)
		op := newOp(tower)

		mock.ExpectBegin()
		expectGuards(mock, ctl, tower)
		expectHeld(mock, heldLock{ResourceKey: ctl, OperationId: uuid.NewV4(),
			Type: "RestartNode", Status: OperationSuccess})
		expectInserts(mock, op, 1)
		mock.ExpectCommit()

		_, err := repo.Start(ctx, []*Operation{op}, []string{ctl}, 5*time.Minute)

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
			mock, repo := setupTestDB(t)
			held.OperationId = uuid.NewV4()
			held.Type = "UpdateSoftware"

			mock.ExpectBegin()
			expectGuards(mock, ctl, tower)
			expectHeld(mock, held)
			mock.ExpectRollback()

			out, err := repo.Start(ctx, []*Operation{newOp(tower)}, []string{ctl},
				5*time.Minute)

			var conflict *ConflictError
			assert.True(t, errors.As(err, &conflict), err)
			assert.True(t, errors.Is(err, ErrLockConflict), err)
			assert.Equal(t, held.ResourceKey, conflict.ResourceKey)
			assert.Equal(t, held.OperationId, conflict.HolderId)
			assert.Contains(t, err.Error(), "UpdateSoftware")
			assert.Nil(t, out)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}

	t.Run("LockRowRaceIsConflict", func(t *testing.T) {
		mock, repo := setupTestDB(t)
		op := newOp(tower)

		mock.ExpectBegin()
		expectGuards(mock, tower)
		expectHeld(mock)
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "operations"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "fencing_token"}).AddRow(op.Id, 8))
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "resource_locks"`)).
			WillReturnError(&pgconn.PgError{Code: "23505"})
		mock.ExpectRollback()

		out, err := repo.Start(ctx, []*Operation{op}, nil, 5*time.Minute)

		assert.True(t, errors.Is(err, ErrLockConflict), err)
		assert.Nil(t, out)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("TakenIdempotencyKey", func(t *testing.T) {
		mock, repo := setupTestDB(t)

		mock.ExpectBegin()
		expectGuards(mock, tower)
		expectHeld(mock)
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "operations"`)).
			WillReturnError(&pgconn.PgError{Code: "23505"})
		mock.ExpectRollback()

		out, err := repo.Start(ctx, []*Operation{newOp(tower)}, nil, 5*time.Minute)

		assert.True(t, errors.Is(err, ErrIdempotencyKeyTaken), err)
		assert.False(t, errors.Is(err, ErrLockConflict), err)
		assert.Nil(t, out)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GuardTimeoutIsLockTimeout", func(t *testing.T) {
		mock, repo := setupTestDB(t)

		mock.ExpectBegin()
		expectLockTimeout(mock)
		mock.ExpectExec(regexp.QuoteMeta("SELECT pg_advisory_xact_lock")).
			WillReturnError(&pgconn.PgError{Code: pgErrLockNotAvailable})
		mock.ExpectRollback()

		out, err := repo.Start(ctx, []*Operation{newOp(tower)}, nil, 5*time.Minute)

		assert.True(t, errors.Is(err, ErrLockTimeout), err)
		assert.False(t, errors.Is(err, ErrLockConflict), err)
		assert.Nil(t, out)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("HeldLookupFailure", func(t *testing.T) {
		mock, repo := setupTestDB(t)

		mock.ExpectBegin()
		expectGuards(mock, tower)
		mock.ExpectQuery(`SELECT resource_locks.resource_key.* FROM "resource_locks"`).
			WillReturnError(errors.New("connection reset"))
		mock.ExpectRollback()

		out, err := repo.Start(ctx, []*Operation{newOp(tower)}, nil, 5*time.Minute)

		assert.Error(t, err)
		assert.False(t, errors.Is(err, ErrLockConflict), err)
		assert.Nil(t, out)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("OtherOperationInsertErrorIsWrapped", func(t *testing.T) {
		mock, repo := setupTestDB(t)

		mock.ExpectBegin()
		expectGuards(mock, tower)
		expectHeld(mock)
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "operations"`)).
			WillReturnError(errors.New("connection reset"))
		mock.ExpectRollback()

		out, err := repo.Start(ctx, []*Operation{newOp(tower)}, nil, 5*time.Minute)

		assert.Error(t, err)
		assert.False(t, errors.Is(err, ErrIdempotencyKeyTaken), err)
		assert.Nil(t, out)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("LockTimeoutStatementFailure", func(t *testing.T) {
		mock, repo := setupTestDB(t)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(lockTimeoutStmt)).
			WillReturnError(errors.New("connection reset"))
		mock.ExpectRollback()

		out, err := repo.Start(ctx, []*Operation{newOp(tower)}, nil, 5*time.Minute)

		assert.Error(t, err)
		assert.False(t, errors.Is(err, ErrLockConflict), err)
		assert.Nil(t, out)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func Test_MarkRunning(t *testing.T) {
	pending := func() *Operation {
		return &Operation{
			Id:             uuid.NewV4(),
			Type:           "RestartNode",
			System:         "node",
			Status:         OperationPending,
			FencingToken:   3,
			ResourceKey:    "node:abc",
			LeaseExpiresAt: time.Now().Add(time.Minute),
		}
	}
	expectAudit := func(mock sqlmock.Sqlmock, op *Operation, event string) {
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "operation_audits"`)).
			WithArgs(sqlmock.AnyArg(), op.Id, op.ResourceKey, event, "controller",
				sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))
	}

	t.Run("TransitionsPendingToRunning", func(t *testing.T) {
		mock, repo := setupTestDB(t)
		op := pending()

		mock.ExpectBegin()
		expectLockedSelect(mock).WillReturnRows(operationRow(op))
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "operations"`)).
			WillReturnResult(sqlmock.NewResult(1, 1))
		expectAudit(mock, op, "running")
		mock.ExpectCommit()

		out, err := repo.MarkRunning(ctx, op.Id, op.FencingToken, "controller")

		assert.NoError(t, err)
		assert.Equal(t, OperationRunning, out.Status)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("RepeatIsAuditedAndSucceeds", func(t *testing.T) {
		mock, repo := setupTestDB(t)
		op := pending()
		op.Status = OperationRunning

		mock.ExpectBegin()
		expectLockedSelect(mock).WillReturnRows(operationRow(op))
		expectAudit(mock, op, "ignored_running")
		mock.ExpectCommit()

		out, err := repo.MarkRunning(ctx, op.Id, op.FencingToken, "controller")

		assert.NoError(t, err)
		assert.Equal(t, OperationRunning, out.Status)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	for name, change := range map[string]func(op *Operation) uint64{
		"RejectsStaleFencingToken": func(op *Operation) uint64 {
			return op.FencingToken - 1
		},
		"RejectsFinishedOperation": func(op *Operation) uint64 {
			op.Status = OperationTimeout
			return op.FencingToken
		},
		"RejectsExpiredLease": func(op *Operation) uint64 {
			op.LeaseExpiresAt = time.Now().Add(-time.Second)
			return op.FencingToken
		},
		"RejectsExpiredLeaseWhileRunning": func(op *Operation) uint64 {
			op.Status = OperationRunning
			op.LeaseExpiresAt = time.Now().Add(-time.Second)
			return op.FencingToken
		},
	} {
		t.Run(name, func(t *testing.T) {
			mock, repo := setupTestDB(t)
			op := pending()
			token := change(op)

			mock.ExpectBegin()
			expectLockedSelect(mock).WillReturnRows(operationRow(op))
			expectAudit(mock, op, "rejected_running")
			mock.ExpectCommit()

			out, err := repo.MarkRunning(ctx, op.Id, token, "controller")

			assert.True(t, errors.Is(err, store.ErrPrecondition), err)
			assert.Nil(t, out)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}

	t.Run("UnknownOperationIsNotFound", func(t *testing.T) {
		mock, repo := setupTestDB(t)

		mock.ExpectBegin()
		expectLockedSelect(mock).WillReturnError(gorm.ErrRecordNotFound)
		mock.ExpectRollback()

		out, err := repo.MarkRunning(ctx, uuid.NewV4(), 3, "controller")

		assert.True(t, errors.Is(err, store.ErrNotFound), err)
		assert.Nil(t, out)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ReturnsLockTimeout", func(t *testing.T) {
		mock, repo := setupTestDB(t)

		mock.ExpectBegin()
		expectLockedSelect(mock).
			WillReturnError(&pgconn.PgError{Code: pgErrLockNotAvailable})
		mock.ExpectRollback()

		out, err := repo.MarkRunning(ctx, uuid.NewV4(), 3, "controller")

		assert.True(t, errors.Is(err, ErrLockTimeout), err)
		assert.Nil(t, out)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func Test_Terminate(t *testing.T) {
	t.Run("ReleasesLockOnSuccess", func(t *testing.T) {
		mock, repo := setupTestDB(t)

		op := &Operation{
			Id:           uuid.NewV4(),
			Type:         "RestartNode",
			System:       "node",
			Status:       OperationRunning,
			FencingToken: 5,
			ResourceKey:  "node:abc",
		}

		mock.ExpectBegin()
		expectLockedSelect(mock).
			WillReturnRows(operationRow(op))
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "operations"`)).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "resource_locks"`)).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "operation_audits"`)).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		out, err := repo.Terminate(ctx, op.Id, op.FencingToken, OperationSuccess,
			OperationAudit{Event: "completed"}, "")

		assert.NoError(t, err)
		assert.Equal(t, OperationSuccess, out.Status)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("RejectsStaleFencingToken", func(t *testing.T) {
		mock, repo := setupTestDB(t)

		op := &Operation{
			Id:           uuid.NewV4(),
			Status:       OperationRunning,
			FencingToken: 9, // current
			ResourceKey:  "node:abc",
		}

		mock.ExpectBegin()
		expectLockedSelect(mock).
			WillReturnRows(operationRow(op))
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "operation_audits"`)).
			WithArgs(sqlmock.AnyArg(), op.Id, op.ResourceKey, "rejected_completed",
				sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		// zombie carries an older token (4 < 9) → rejected
		out, err := repo.Terminate(ctx, op.Id, 4, OperationSuccess,
			OperationAudit{Event: "completed"}, "")

		assert.True(t, errors.Is(err, store.ErrPrecondition), err)
		assert.Nil(t, out)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DuplicateWithSameStatusIsAuditedAndSucceeds", func(t *testing.T) {
		mock, repo := setupTestDB(t)

		op := &Operation{
			Id:           uuid.NewV4(),
			Status:       OperationSuccess,
			FencingToken: 5,
			ResourceKey:  "node:abc",
		}

		mock.ExpectBegin()
		expectLockedSelect(mock).
			WillReturnRows(operationRow(op))
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "operation_audits"`)).
			WithArgs(sqlmock.AnyArg(), op.Id, op.ResourceKey, "ignored_completed",
				sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		out, err := repo.Terminate(ctx, op.Id, op.FencingToken, OperationSuccess,
			OperationAudit{Event: "completed"}, "")

		assert.NoError(t, err)
		assert.Equal(t, OperationSuccess, out.Status)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DifferentStatusIsAuditedAndRejected", func(t *testing.T) {
		mock, repo := setupTestDB(t)

		op := &Operation{
			Id:           uuid.NewV4(),
			Status:       OperationTimeout,
			FencingToken: 5,
			ResourceKey:  "node:abc",
		}

		mock.ExpectBegin()
		expectLockedSelect(mock).
			WillReturnRows(operationRow(op))
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "operation_audits"`)).
			WithArgs(sqlmock.AnyArg(), op.Id, op.ResourceKey, "ignored_completed",
				sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		out, err := repo.Terminate(ctx, op.Id, op.FencingToken, OperationSuccess,
			OperationAudit{Event: "completed"}, "")

		assert.True(t, errors.Is(err, store.ErrPrecondition), err)
		assert.Nil(t, out)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("RejectsCompleteOnPending", func(t *testing.T) {
		mock, repo := setupTestDB(t)

		op := &Operation{
			Id:           uuid.NewV4(),
			Status:       OperationPending,
			FencingToken: 5,
			ResourceKey:  "node:abc",
		}

		mock.ExpectBegin()
		expectLockedSelect(mock).
			WillReturnRows(operationRow(op))
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "operation_audits"`)).
			WithArgs(sqlmock.AnyArg(), op.Id, op.ResourceKey, "rejected_completed",
				sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		out, err := repo.Terminate(ctx, op.Id, op.FencingToken, OperationSuccess,
			OperationAudit{Event: "completed"}, "")

		assert.True(t, errors.Is(err, store.ErrPrecondition), err)
		assert.Nil(t, out)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("AllowsFailOnPending", func(t *testing.T) {
		mock, repo := setupTestDB(t)

		op := &Operation{
			Id:           uuid.NewV4(),
			Status:       OperationPending,
			FencingToken: 5,
			ResourceKey:  "node:abc",
		}

		mock.ExpectBegin()
		expectLockedSelect(mock).
			WillReturnRows(operationRow(op))
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "operations"`)).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "resource_locks"`)).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "operation_audits"`)).
			WithArgs(sqlmock.AnyArg(), op.Id, op.ResourceKey, "failed",
				sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		out, err := repo.Terminate(ctx, op.Id, op.FencingToken, OperationFailed,
			OperationAudit{Event: "failed"}, "boom")

		assert.NoError(t, err)
		assert.Equal(t, OperationFailed, out.Status)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ReturnsLockTimeout", func(t *testing.T) {
		mock, repo := setupTestDB(t)

		mock.ExpectBegin()
		expectLockedSelect(mock).
			WillReturnError(&pgconn.PgError{Code: pgErrLockNotAvailable})
		mock.ExpectRollback()

		out, err := repo.Terminate(ctx, uuid.NewV4(), 5, OperationSuccess,
			OperationAudit{Event: "completed"}, "")

		assert.True(t, errors.Is(err, ErrLockTimeout), err)
		assert.False(t, errors.Is(err, store.ErrPrecondition))
		assert.Nil(t, out)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("CancelledContextChangesNothing", func(t *testing.T) {
		mock, repo := setupTestDB(t)

		cancelled, cancel := context.WithCancel(context.Background())
		cancel()

		out, err := repo.Terminate(cancelled, uuid.NewV4(), 5, OperationSuccess,
			OperationAudit{Event: "completed"}, "")

		assert.True(t, errors.Is(err, context.Canceled), err)
		assert.Nil(t, out)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("RejectsNonTerminalStatus", func(t *testing.T) {
		_, repo := setupTestDB(t)

		out, err := repo.Terminate(ctx, uuid.NewV4(), 1, OperationRunning,
			OperationAudit{Event: "running"}, "")

		assert.Error(t, err)
		assert.Nil(t, out)
	})
}

func Test_GetByResource(t *testing.T) {
	t.Run("ReturnsNilWhenNotLocked", func(t *testing.T) {
		mock, repo := setupTestDB(t)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "resource_locks"`)).
			WillReturnError(gorm.ErrRecordNotFound)

		out, err := repo.GetByResource(ctx, "node:free")

		assert.NoError(t, err)
		assert.Nil(t, out)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("ReturnsLockHolder", func(t *testing.T) {
		mock, repo := setupTestDB(t)

		holder := &Operation{
			Id:           uuid.NewV4(),
			Status:       OperationRunning,
			FencingToken: 4,
			ResourceKey:  "node:busy",
		}
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "resource_locks"`)).
			WillReturnRows(sqlmock.NewRows([]string{
				"resource_key", "operation_id", "fencing_token",
			}).AddRow(holder.ResourceKey, holder.Id, holder.FencingToken))
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "operations"`)).
			WillReturnRows(operationRow(holder))

		out, err := repo.GetByResource(ctx, holder.ResourceKey)

		assert.NoError(t, err)
		assert.Equal(t, holder.Id, out.Id)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func Test_Get(t *testing.T) {
	t.Run("ReturnsOperation", func(t *testing.T) {
		mock, repo := setupTestDB(t)

		op := &Operation{Id: uuid.NewV4(), Status: OperationPending, ResourceKey: "node:a"}
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "operations"`)).
			WillReturnRows(operationRow(op))

		out, err := repo.Get(ctx, op.Id)

		assert.NoError(t, err)
		assert.Equal(t, op.Id, out.Id)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("MissingOperationIsNotFound", func(t *testing.T) {
		mock, repo := setupTestDB(t)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "operations"`)).
			WillReturnError(gorm.ErrRecordNotFound)

		out, err := repo.Get(ctx, uuid.NewV4())

		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound), err)
		assert.Nil(t, out)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func Test_GetByIdempotencyKey(t *testing.T) {
	t.Run("ReturnsLeaderThenMembers", func(t *testing.T) {
		mock, repo := setupTestDB(t)

		leader := &Operation{Id: uuid.NewV4(), Status: OperationPending,
			ResourceKey: "node:a"}
		member := &Operation{Id: uuid.NewV4(), Status: OperationPending,
			ResourceKey: "node:t"}
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "operations"`)).
			WithArgs("key", sqlmock.AnyArg()).
			WillReturnRows(operationRow(leader))
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "operations"`)).
			WithArgs(leader.Id).
			WillReturnRows(operationRow(member))

		out, err := repo.GetByIdempotencyKey(ctx, "key")

		assert.NoError(t, err)
		assert.Len(t, out, 2)
		assert.Equal(t, leader.Id, out[0].Id)
		assert.Equal(t, member.Id, out[1].Id)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("UnknownKeyReturnsNil", func(t *testing.T) {
		mock, repo := setupTestDB(t)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "operations"`)).
			WillReturnError(gorm.ErrRecordNotFound)

		out, err := repo.GetByIdempotencyKey(ctx, "unknown")

		assert.NoError(t, err)
		assert.Nil(t, out)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("QueryErrorIsReturned", func(t *testing.T) {
		mock, repo := setupTestDB(t)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "operations"`)).
			WillReturnError(errors.New("connection reset"))

		out, err := repo.GetByIdempotencyKey(ctx, "key")

		assert.Error(t, err)
		assert.Nil(t, out)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("MemberQueryErrorIsReturned", func(t *testing.T) {
		mock, repo := setupTestDB(t)

		leader := &Operation{Id: uuid.NewV4(), ResourceKey: "node:a"}
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "operations"`)).
			WillReturnRows(operationRow(leader))
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "operations"`)).
			WillReturnError(errors.New("connection reset"))

		out, err := repo.GetByIdempotencyKey(ctx, "key")

		assert.Error(t, err)
		assert.Nil(t, out)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func Test_FindExpired(t *testing.T) {
	t.Run("ReturnsExpiredNonTerminalOperations", func(t *testing.T) {
		mock, repo := setupTestDB(t)

		op := &Operation{
			Id:           uuid.NewV4(),
			Status:       OperationRunning,
			FencingToken: 2,
			ResourceKey:  "node:abc",
		}

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "operations"`)).
			WillReturnRows(operationRow(op))

		out, err := repo.FindExpired(ctx, time.Now(), 100)

		assert.NoError(t, err)
		assert.Len(t, out, 1)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func expectLockTimeout(mock sqlmock.Sqlmock) {
	mock.ExpectExec(regexp.QuoteMeta(lockTimeoutStmt)).
		WillReturnResult(sqlmock.NewResult(0, 0))
}

func expectLockedSelect(mock sqlmock.Sqlmock) *sqlmock.ExpectedQuery {
	expectLockTimeout(mock)
	return mock.ExpectQuery(`SELECT \* FROM "operations" .* FOR UPDATE`)
}

var ctx = context.Background()
