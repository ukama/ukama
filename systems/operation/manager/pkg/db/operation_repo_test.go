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
	log "github.com/sirupsen/logrus"
	"github.com/tj/assert"
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
	t.Run("AcquiresLockSuccessfully", func(t *testing.T) {
		mock, repo := setupTestDB(t)

		op := &Operation{
			Id:             uuid.NewV4(),
			Type:           "RestartNode",
			System:         "node",
			Status:         OperationPending,
			ResourceKey:    "node:uk-sa2450-tnode-v0-4e86",
			LeaseExpiresAt: time.Now().Add(5 * time.Minute),
		}

		mock.ExpectBegin()
		// insert operation (returns generated id + fencing_token)
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "operations"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "fencing_token"}).AddRow(op.Id, 1))
		// insert lock
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "resource_locks"`)).
			WillReturnResult(sqlmock.NewResult(1, 1))
		// insert audit
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "operation_audits"`)).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		out, err := repo.Start(op, 5*time.Minute)

		assert.NoError(t, err)
		assert.NotNil(t, out)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ReturnsConflictWhenResourceLocked", func(t *testing.T) {
		mock, repo := setupTestDB(t)

		op := &Operation{
			Id:             uuid.NewV4(),
			Type:           "RestartNode",
			System:         "node",
			Status:         OperationPending,
			ResourceKey:    "node:uk-sa2450-tnode-v0-4e86",
			LeaseExpiresAt: time.Now().Add(5 * time.Minute),
		}
		holder := &Operation{
			Id:           uuid.NewV4(),
			Type:         "RestartNode",
			System:       "node",
			Status:       OperationRunning,
			FencingToken: 7,
			ResourceKey:  op.ResourceKey,
		}

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "operations"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "fencing_token"}).AddRow(op.Id, 8))
		// lock insert fails on PK collision
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "resource_locks"`)).
			WillReturnError(gorm.ErrDuplicatedKey)
		// look up existing lock holder
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "resource_locks"`)).
			WithArgs(op.ResourceKey, sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"resource_key", "operation_id", "fencing_token"}).
				AddRow(op.ResourceKey, holder.Id, holder.FencingToken))
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "operations"`)).
			WithArgs(holder.Id, sqlmock.AnyArg()).
			WillReturnRows(operationRow(holder))
		mock.ExpectRollback()

		out, err := repo.Start(op, 5*time.Minute)

		assert.True(t, errors.Is(err, ErrLockConflict))
		assert.NotNil(t, out)
		assert.Equal(t, holder.Id, out.Id)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func Test_MarkRunning(t *testing.T) {
	t.Run("TransitionsPendingToRunning", func(t *testing.T) {
		mock, repo := setupTestDB(t)

		op := &Operation{
			Id:           uuid.NewV4(),
			Type:         "RestartNode",
			System:       "node",
			Status:       OperationPending,
			FencingToken: 3,
			ResourceKey:  "node:abc",
		}

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "operations"`)).
			WillReturnRows(operationRow(op))
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "operations"`)).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "operation_audits"`)).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		out, err := repo.MarkRunning(op.Id, op.FencingToken)

		assert.NoError(t, err)
		assert.Equal(t, OperationRunning, out.Status)
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
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "operations"`)).
			WillReturnRows(operationRow(op))
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "operations"`)).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "resource_locks"`)).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "operation_audits"`)).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		out, err := repo.Terminate(op.Id, op.FencingToken, OperationSuccess,
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
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "operations"`)).
			WillReturnRows(operationRow(op))
		mock.ExpectRollback()

		// zombie carries an older token (4 < 9) → rejected
		out, err := repo.Terminate(op.Id, 4, OperationSuccess,
			OperationAudit{Event: "completed"}, "")

		assert.Error(t, err)
		assert.Nil(t, out)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("RejectsNonTerminalStatus", func(t *testing.T) {
		_, repo := setupTestDB(t)

		out, err := repo.Terminate(uuid.NewV4(), 1, OperationRunning,
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

		out, err := repo.GetByResource("node:free")

		assert.NoError(t, err)
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

		out, err := repo.FindExpired(time.Now(), 100)

		assert.NoError(t, err)
		assert.Len(t, out, 1)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
