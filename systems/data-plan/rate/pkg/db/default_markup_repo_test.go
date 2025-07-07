/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package db

import (
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/tj/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	defaultMarkupValue = 10.0
	testID             = 1
)

// Test data structures
type testSetup struct {
	mock   sqlmock.Sqlmock
	gormDB *gorm.DB
	repo   *defaultMarkupRepo
}

// Helper functions
func setupTestDB(t *testing.T) *testSetup {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})
	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	assert.NoError(t, err)

	repo := NewDefaultMarkupRepo(&UkamaDbMock{
		GormDb: gormDB,
	})

	return &testSetup{
		mock:   mock,
		gormDB: gormDB,
		repo:   repo,
	}
}

func createDefaultMarkupRow(markup float64) *sqlmock.Rows {
	return sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "markup"}).
		AddRow(testID, time.Now(), time.Now(), nil, markup)
}

func expectTransactionBegin(mock sqlmock.Sqlmock) {
	mock.ExpectBegin()
}

func expectTransactionCommit(mock sqlmock.Sqlmock) {
	mock.ExpectCommit()
}

func expectTransactionRollback(mock sqlmock.Sqlmock) {
	mock.ExpectRollback()
}

func expectInsertQuery(mock sqlmock.Sqlmock, markup float64) {
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), markup).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(testID))
}

func expectInsertQueryError(mock sqlmock.Sqlmock, markup float64, errMsg string) {
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), markup).
		WillReturnError(errors.New(errMsg))
}

func expectSelectQuery(mock sqlmock.Sqlmock, rows *sqlmock.Rows) {
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).
		WillReturnRows(rows)
}

func expectSelectQueryError(mock sqlmock.Sqlmock, err error) {
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).
		WillReturnError(err)
}

func expectUpdateQuery(mock sqlmock.Sqlmock) {
	mock.ExpectExec(regexp.QuoteMeta("UPDATE")).WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
}

func expectUpdateQueryError(mock sqlmock.Sqlmock, errMsg string) {
	mock.ExpectExec(regexp.QuoteMeta("UPDATE")).WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(errors.New(errMsg))
}

func expectDeleteQuery(mock sqlmock.Sqlmock) {
	mock.ExpectExec(regexp.QuoteMeta("UPDATE")).WithArgs(sqlmock.AnyArg(), nil).
		WillReturnResult(sqlmock.NewResult(1, 1))
}

func expectDeleteQueryError(mock sqlmock.Sqlmock, errMsg string) {
	mock.ExpectExec(regexp.QuoteMeta("UPDATE")).WithArgs(sqlmock.AnyArg(), nil).
		WillReturnError(errors.New(errMsg))
}

func expectTransactionBeginError(mock sqlmock.Sqlmock, errMsg string) {
	mock.ExpectBegin().WillReturnError(errors.New(errMsg))
}

func TestDefaultMarkupRepo_Create(t *testing.T) {

	t.Run("Create", func(t *testing.T) {
		// Arrange
		setup := setupTestDB(t)

		expectTransactionBegin(setup.mock)
		expectInsertQuery(setup.mock, defaultMarkupValue)
		expectTransactionCommit(setup.mock)

		// Act
		err := setup.repo.CreateDefaultMarkupRate(defaultMarkupValue)

		// Assert
		assert.NoError(t, err)
		assert.NoError(t, setup.mock.ExpectationsWereMet())
	})

	t.Run("Create_DatabaseError", func(t *testing.T) {
		// Arrange
		setup := setupTestDB(t)

		expectTransactionBegin(setup.mock)
		expectInsertQueryError(setup.mock, defaultMarkupValue, "database error")
		expectTransactionRollback(setup.mock)

		// Act
		err := setup.repo.CreateDefaultMarkupRate(defaultMarkupValue)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")
		assert.NoError(t, setup.mock.ExpectationsWereMet())
	})
}

func TestDefaultMarkupRepo_Get(t *testing.T) {

	t.Run("Get_Success", func(t *testing.T) {
		// Arrange
		setup := setupTestDB(t)

		rows := createDefaultMarkupRow(defaultMarkupValue)
		expectSelectQuery(setup.mock, rows)

		// Act
		m, err := setup.repo.GetDefaultMarkupRate()

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, m)
		assert.EqualValues(t, defaultMarkupValue, m.Markup)
		assert.NoError(t, setup.mock.ExpectationsWereMet())
	})

	t.Run("Get_NotFound", func(t *testing.T) {
		// Arrange
		setup := setupTestDB(t)

		expectSelectQueryError(setup.mock, gorm.ErrRecordNotFound)

		// Act
		m, err := setup.repo.GetDefaultMarkupRate()

		// Assert
		assert.Error(t, err)
		assert.Nil(t, m)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
		assert.NoError(t, setup.mock.ExpectationsWereMet())
	})

	t.Run("Get_DatabaseError", func(t *testing.T) {
		// Arrange
		setup := setupTestDB(t)

		expectSelectQueryError(setup.mock, errors.New("database error"))

		// Act
		m, err := setup.repo.GetDefaultMarkupRate()

		// Assert
		assert.Error(t, err)
		assert.Nil(t, m)
		assert.Contains(t, err.Error(), "database error")
		assert.NoError(t, setup.mock.ExpectationsWereMet())
	})
}

func TestDefaultMarkupRepo_Delete(t *testing.T) {

	t.Run("Delete_Success", func(t *testing.T) {
		// Arrange
		setup := setupTestDB(t)

		expectTransactionBegin(setup.mock)
		expectDeleteQuery(setup.mock)
		expectTransactionCommit(setup.mock)

		// Act
		err := setup.repo.DeleteDefaultMarkupRate()

		// Assert
		assert.NoError(t, err)
		assert.NoError(t, setup.mock.ExpectationsWereMet())
	})

	t.Run("Delete_DatabaseError", func(t *testing.T) {
		// Arrange
		setup := setupTestDB(t)

		expectTransactionBegin(setup.mock)
		expectDeleteQueryError(setup.mock, "database error")
		expectTransactionRollback(setup.mock)

		// Act
		err := setup.repo.DeleteDefaultMarkupRate()

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")
		assert.NoError(t, setup.mock.ExpectationsWereMet())
	})

	t.Run("Delete_TransactionBeginError", func(t *testing.T) {
		// Arrange
		setup := setupTestDB(t)

		expectTransactionBeginError(setup.mock, "transaction begin error")

		// Act
		err := setup.repo.DeleteDefaultMarkupRate()

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "transaction begin error")
		assert.NoError(t, setup.mock.ExpectationsWereMet())
	})
}

func TestDefaultMarkupRepo_Update(t *testing.T) {

	t.Run("Update_Success", func(t *testing.T) {
		// Arrange
		setup := setupTestDB(t)

		expectTransactionBegin(setup.mock)
		expectUpdateQuery(setup.mock)
		expectInsertQuery(setup.mock, defaultMarkupValue)
		expectTransactionCommit(setup.mock)

		// Act
		err := setup.repo.UpdateDefaultMarkupRate(defaultMarkupValue)

		// Assert
		assert.NoError(t, err)
		assert.NoError(t, setup.mock.ExpectationsWereMet())
	})

	t.Run("Update_DeleteError", func(t *testing.T) {
		// Arrange
		setup := setupTestDB(t)

		expectTransactionBegin(setup.mock)
		expectUpdateQueryError(setup.mock, "delete error")
		expectTransactionRollback(setup.mock)

		// Act
		err := setup.repo.UpdateDefaultMarkupRate(defaultMarkupValue)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "delete error")
		assert.NoError(t, setup.mock.ExpectationsWereMet())
	})

	t.Run("Update_CreateError", func(t *testing.T) {
		// Arrange
		setup := setupTestDB(t)

		expectTransactionBegin(setup.mock)
		expectUpdateQuery(setup.mock)
		expectInsertQueryError(setup.mock, defaultMarkupValue, "create error")
		expectTransactionRollback(setup.mock)

		// Act
		err := setup.repo.UpdateDefaultMarkupRate(defaultMarkupValue)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "create error")
		assert.NoError(t, setup.mock.ExpectationsWereMet())
	})
}

func TestDefaultMarkupRepo_GetHistory(t *testing.T) {

	t.Run("GetHistory_Success", func(t *testing.T) {
		// Arrange
		setup := setupTestDB(t)

		rows := createDefaultMarkupRow(defaultMarkupValue)
		expectSelectQuery(setup.mock, rows)

		// Act
		m, err := setup.repo.GetDefaultMarkupRateHistory()

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, m)
		assert.EqualValues(t, defaultMarkupValue, m[0].Markup)
		assert.NoError(t, setup.mock.ExpectationsWereMet())
	})

	t.Run("GetHistory_DatabaseError", func(t *testing.T) {
		// Arrange
		setup := setupTestDB(t)

		expectSelectQueryError(setup.mock, errors.New("db error"))

		// Act
		m, err := setup.repo.GetDefaultMarkupRateHistory()

		// Assert
		assert.Error(t, err)
		assert.Nil(t, m)
		assert.Contains(t, err.Error(), "db error")
		assert.NoError(t, setup.mock.ExpectationsWereMet())
	})
}
