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
	"database/sql/driver"
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

var ctx = context.Background()

func setupTestDB(t *testing.T) (sqlmock.Sqlmock, *UkamaDbMock) {
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

	return mock, &UkamaDbMock{GormDb: gdb}
}

func setupOperationRepo(t *testing.T) (sqlmock.Sqlmock, OperationRepo) {
	mock, udb := setupTestDB(t)
	return mock, NewOperationRepo(udb)
}

func operationRow(op *Operation) *sqlmock.Rows {
	return sqlmock.NewRows([]string{
		"id", "type", "system", "status", "fencing_token", "requested_by",
		"resource_key", "lease_expires_at", "created_at",
	}).AddRow(op.Id, op.Type, op.System, op.Status, op.FencingToken,
		op.RequestedBy, op.ResourceKey, op.LeaseExpiresAt, time.Now())
}

func expectLockTimeout(mock sqlmock.Sqlmock) {
	mock.ExpectExec(regexp.QuoteMeta(lockTimeoutStmt)).
		WillReturnResult(sqlmock.NewResult(0, 0))
}

func expectLockedSelect(mock sqlmock.Sqlmock) *sqlmock.ExpectedQuery {
	expectLockTimeout(mock)
	return mock.ExpectQuery(`SELECT \* FROM "operations" .* FOR UPDATE`)
}

func newOp(key string) *Operation {
	return &Operation{
		Id:             uuid.NewV4(),
		Type:           "RestartNode",
		System:         "node",
		Status:         OperationPending,
		FencingToken:   3,
		ResourceKey:    key,
		LeaseExpiresAt: time.Now().Add(5 * time.Minute),
	}
}

func expectOperationInsert(mock sqlmock.Sqlmock, op *Operation, token int) {
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "operations"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "fencing_token"}).
			AddRow(op.Id, token))
}

func Test_OperationAdd(t *testing.T) {
	t.Run("InsertsEveryOperationThenRunsNested", func(t *testing.T) {
		mock, repo := setupOperationRepo(t)
		first, second := newOp("node:a"), newOp("node:t")

		mock.ExpectBegin()
		expectLockTimeout(mock)
		expectOperationInsert(mock, first, 1)
		expectOperationInsert(mock, second, 2)
		mock.ExpectExec(regexp.QuoteMeta("SELECT pg_advisory_xact_lock")).
			WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		var seen *gorm.DB
		err := repo.Add(ctx, []*Operation{first, second}, func(tx *gorm.DB) error {
			seen = tx
			return tx.Exec(advisoryLockStmt, "node:a").Error
		})

		assert.NoError(t, err)
		assert.NotNil(t, seen)
		assert.Equal(t, uint64(1), first.FencingToken)
		assert.Equal(t, uint64(2), second.FencingToken)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("NilNestedOnlyInserts", func(t *testing.T) {
		mock, repo := setupOperationRepo(t)
		op := newOp("node:a")

		mock.ExpectBegin()
		expectLockTimeout(mock)
		expectOperationInsert(mock, op, 1)
		mock.ExpectCommit()

		assert.NoError(t, repo.Add(ctx, []*Operation{op}, nil))
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("TakenIdempotencyKey", func(t *testing.T) {
		mock, repo := setupOperationRepo(t)

		mock.ExpectBegin()
		expectLockTimeout(mock)
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "operations"`)).
			WillReturnError(&pgconn.PgError{Code: "23505"})
		mock.ExpectRollback()

		err := repo.Add(ctx, []*Operation{newOp("node:a")}, nil)

		assert.True(t, errors.Is(err, ErrIdempotencyKeyTaken), err)
		assert.False(t, errors.Is(err, ErrLockConflict), err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("NestedConflictIsReturnedAsIs", func(t *testing.T) {
		mock, repo := setupOperationRepo(t)
		op := newOp("node:a")
		conflict := &ConflictError{ResourceKey: "node:a", HolderId: uuid.NewV4()}

		mock.ExpectBegin()
		expectLockTimeout(mock)
		expectOperationInsert(mock, op, 1)
		mock.ExpectRollback()

		err := repo.Add(ctx, []*Operation{op}, func(*gorm.DB) error { return conflict })

		var got *ConflictError
		assert.True(t, errors.As(err, &got), err)
		assert.Equal(t, conflict.HolderId, got.HolderId)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("NestedLockTimeoutIsLockTimeout", func(t *testing.T) {
		mock, repo := setupOperationRepo(t)
		op := newOp("node:a")

		mock.ExpectBegin()
		expectLockTimeout(mock)
		expectOperationInsert(mock, op, 1)
		mock.ExpectRollback()

		err := repo.Add(ctx, []*Operation{op}, func(*gorm.DB) error {
			return &pgconn.PgError{Code: pgErrLockNotAvailable}
		})

		assert.True(t, errors.Is(err, ErrLockTimeout), err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("OtherInsertErrorIsWrapped", func(t *testing.T) {
		mock, repo := setupOperationRepo(t)

		mock.ExpectBegin()
		expectLockTimeout(mock)
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "operations"`)).
			WillReturnError(errors.New("connection reset"))
		mock.ExpectRollback()

		err := repo.Add(ctx, []*Operation{newOp("node:a")}, nil)

		assert.Error(t, err)
		assert.False(t, errors.Is(err, ErrIdempotencyKeyTaken), err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("LockTimeoutStatementFailure", func(t *testing.T) {
		mock, repo := setupOperationRepo(t)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(lockTimeoutStmt)).
			WillReturnError(errors.New("connection reset"))
		mock.ExpectRollback()

		err := repo.Add(ctx, []*Operation{newOp("node:a")}, nil)

		assert.Error(t, err)
		assert.False(t, errors.Is(err, ErrLockConflict), err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func Test_OperationGet(t *testing.T) {
	t.Run("ReturnsOperation", func(t *testing.T) {
		mock, repo := setupOperationRepo(t)
		op := newOp("node:a")

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "operations"`)).
			WillReturnRows(operationRow(op))

		out, err := repo.Get(ctx, op.Id)

		assert.NoError(t, err)
		assert.Equal(t, op.Id, out.Id)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("MissingOperationIsNotFound", func(t *testing.T) {
		mock, repo := setupOperationRepo(t)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "operations"`)).
			WillReturnError(gorm.ErrRecordNotFound)

		out, err := repo.Get(ctx, uuid.NewV4())

		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound), err)
		assert.Nil(t, out)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func Test_OperationList(t *testing.T) {
	groupId := uuid.NewV4()
	before := time.Now()

	for name, tc := range map[string]struct {
		filter OperationFilter
		query  string
		args   []interface{}
	}{
		"NoFilterListsAll": {
			query: `SELECT \* FROM "operations" WHERE "operations"."deleted_at" IS NULL$`,
		},
		"ByIdempotencyKey": {
			filter: OperationFilter{IdempotencyKey: "key"},
			query:  `WHERE idempotency_key = \$1 AND`,
			args:   []interface{}{"key"},
		},
		"GroupMembersByResourceKey": {
			filter: OperationFilter{GroupId: &groupId, OrderBy: OrderByResourceKey},
			query:  `WHERE group_id = \$1 AND .* ORDER BY resource_key$`,
			args:   []interface{}{groupId},
		},
		"ExpiredOldestFirst": {
			filter: OperationFilter{
				Statuses:           []OperationStatus{OperationPending, OperationRunning},
				LeaseExpiresBefore: &before, OrderBy: OrderByLeaseExpiry, Limit: 100,
			},
			query: `WHERE status IN \(\$1,\$2\) AND lease_expires_at < \$3 AND .* ` +
				`ORDER BY lease_expires_at LIMIT \$4$`,
			args: []interface{}{OperationPending, OperationRunning, before, 100},
		},
		"ByResourceTypeAndSystem": {
			filter: OperationFilter{
				ResourceKey: "node:a", Type: "RestartNode", System: "node",
			},
			query: `WHERE resource_key = \$1 AND type = \$2 AND system = \$3 AND`,
			args:  []interface{}{"node:a", "RestartNode", "node"},
		},
	} {
		t.Run(name, func(t *testing.T) {
			mock, repo := setupOperationRepo(t)
			op := newOp("node:a")

			q := mock.ExpectQuery(tc.query)
			if tc.args != nil {
				values := make([]driver.Value, 0, len(tc.args))
				for _, a := range tc.args {
					values = append(values, a)
				}
				q.WithArgs(values...)
			}
			q.WillReturnRows(operationRow(op))

			out, err := repo.List(ctx, tc.filter)

			assert.NoError(t, err)
			assert.Len(t, out, 1)
			assert.Equal(t, op.Id, out[0].Id)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}

	t.Run("NoMatchIsEmpty", func(t *testing.T) {
		mock, repo := setupOperationRepo(t)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "operations"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}))

		out, err := repo.List(ctx, OperationFilter{IdempotencyKey: "unknown"})

		assert.NoError(t, err)
		assert.Empty(t, out)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("UnknownOrderIsInvalid", func(t *testing.T) {
		mock, repo := setupOperationRepo(t)

		out, err := repo.List(ctx, OperationFilter{OrderBy: "id; DROP TABLE operations"})

		assert.True(t, errors.Is(err, store.ErrInvalid), err)
		assert.Nil(t, out)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("QueryErrorIsReturned", func(t *testing.T) {
		mock, repo := setupOperationRepo(t)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "operations"`)).
			WillReturnError(errors.New("connection reset"))

		out, err := repo.List(ctx, OperationFilter{})

		assert.Error(t, err)
		assert.Nil(t, out)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func Test_OperationUpdate(t *testing.T) {
	t.Run("SavesWhatNestedChanged", func(t *testing.T) {
		mock, repo := setupOperationRepo(t)
		op := newOp("node:a")

		mock.ExpectBegin()
		expectLockedSelect(mock).WillReturnRows(operationRow(op))
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "operation_audits"`)).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "operations"`)).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		out, err := repo.Update(ctx, op.Id, func(locked *Operation, tx *gorm.DB) error {
			assert.Equal(t, op.Id, locked.Id)
			assert.Equal(t, uint64(3), locked.FencingToken)
			locked.Status = OperationRunning
			return tx.Create(&OperationAudit{Id: uuid.NewV4(), OperationId: locked.Id,
				Event: "running", At: time.Now()}).Error
		})

		assert.NoError(t, err)
		assert.Equal(t, OperationRunning, out.Status)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("NilNestedSaves", func(t *testing.T) {
		mock, repo := setupOperationRepo(t)
		op := newOp("node:a")

		mock.ExpectBegin()
		expectLockedSelect(mock).WillReturnRows(operationRow(op))
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "operations"`)).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		out, err := repo.Update(ctx, op.Id, nil)

		assert.NoError(t, err)
		assert.Equal(t, op.Id, out.Id)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("NestedErrorRollsBackAndIsReturned", func(t *testing.T) {
		mock, repo := setupOperationRepo(t)
		op := newOp("node:a")
		refused := errors.New("refused")

		mock.ExpectBegin()
		expectLockedSelect(mock).WillReturnRows(operationRow(op))
		mock.ExpectRollback()

		out, err := repo.Update(ctx, op.Id, func(*Operation, *gorm.DB) error {
			return refused
		})

		assert.True(t, errors.Is(err, refused), err)
		assert.Nil(t, out)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("SaveFailureRollsBack", func(t *testing.T) {
		mock, repo := setupOperationRepo(t)
		op := newOp("node:a")

		mock.ExpectBegin()
		expectLockedSelect(mock).WillReturnRows(operationRow(op))
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "operations"`)).
			WillReturnError(errors.New("connection reset"))
		mock.ExpectRollback()

		out, err := repo.Update(ctx, op.Id, nil)

		assert.Error(t, err)
		assert.Nil(t, out)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("UnknownOperationIsNotFound", func(t *testing.T) {
		mock, repo := setupOperationRepo(t)

		mock.ExpectBegin()
		expectLockedSelect(mock).WillReturnError(gorm.ErrRecordNotFound)
		mock.ExpectRollback()

		out, err := repo.Update(ctx, uuid.NewV4(), nil)

		assert.True(t, errors.Is(err, store.ErrNotFound), err)
		assert.Nil(t, out)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ReturnsLockTimeout", func(t *testing.T) {
		mock, repo := setupOperationRepo(t)

		mock.ExpectBegin()
		expectLockedSelect(mock).
			WillReturnError(&pgconn.PgError{Code: pgErrLockNotAvailable})
		mock.ExpectRollback()

		out, err := repo.Update(ctx, uuid.NewV4(), nil)

		assert.True(t, errors.Is(err, ErrLockTimeout), err)
		assert.False(t, errors.Is(err, store.ErrPrecondition))
		assert.Nil(t, out)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("LockTimeoutStatementFailure", func(t *testing.T) {
		mock, repo := setupOperationRepo(t)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(lockTimeoutStmt)).
			WillReturnError(errors.New("connection reset"))
		mock.ExpectRollback()

		out, err := repo.Update(ctx, uuid.NewV4(), nil)

		assert.Error(t, err)
		assert.Nil(t, out)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("CancelledContextChangesNothing", func(t *testing.T) {
		mock, repo := setupOperationRepo(t)
		cancelled, cancel := context.WithCancel(context.Background())
		cancel()

		out, err := repo.Update(cancelled, uuid.NewV4(), nil)

		assert.True(t, errors.Is(err, context.Canceled), err)
		assert.Nil(t, out)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func Test_OperationDelete(t *testing.T) {
	t.Run("SoftDeletesAfterNested", func(t *testing.T) {
		mock, repo := setupOperationRepo(t)
		op := newOp("node:a")
		op.Status = OperationSuccess

		mock.ExpectBegin()
		expectLockedSelect(mock).WillReturnRows(operationRow(op))
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "operation_audits"`)).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "operations" SET "deleted_at"`)).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		err := repo.Delete(ctx, op.Id, func(locked *Operation, tx *gorm.DB) error {
			assert.Equal(t, OperationSuccess, locked.Status)
			return tx.Create(&OperationAudit{Id: uuid.NewV4(), OperationId: locked.Id,
				Event: "deleted", At: time.Now()}).Error
		})

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("NilNestedDeletes", func(t *testing.T) {
		mock, repo := setupOperationRepo(t)
		op := newOp("node:a")

		mock.ExpectBegin()
		expectLockedSelect(mock).WillReturnRows(operationRow(op))
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "operations" SET "deleted_at"`)).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		assert.NoError(t, repo.Delete(ctx, op.Id, nil))
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("NestedErrorRollsBack", func(t *testing.T) {
		mock, repo := setupOperationRepo(t)
		op := newOp("node:a")

		mock.ExpectBegin()
		expectLockedSelect(mock).WillReturnRows(operationRow(op))
		mock.ExpectRollback()

		err := repo.Delete(ctx, op.Id, func(*Operation, *gorm.DB) error {
			return store.ErrPrecondition
		})

		assert.True(t, errors.Is(err, store.ErrPrecondition), err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("UnknownOperationIsNotFound", func(t *testing.T) {
		mock, repo := setupOperationRepo(t)

		mock.ExpectBegin()
		expectLockedSelect(mock).WillReturnError(gorm.ErrRecordNotFound)
		mock.ExpectRollback()

		err := repo.Delete(ctx, uuid.NewV4(), nil)

		assert.True(t, errors.Is(err, store.ErrNotFound), err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ReturnsLockTimeout", func(t *testing.T) {
		mock, repo := setupOperationRepo(t)

		mock.ExpectBegin()
		expectLockedSelect(mock).
			WillReturnError(&pgconn.PgError{Code: pgErrLockNotAvailable})
		mock.ExpectRollback()

		err := repo.Delete(ctx, uuid.NewV4(), nil)

		assert.True(t, errors.Is(err, ErrLockTimeout), err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
