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
	"log"
	"regexp"
	"testing"
	"time"

	uuid "github.com/ukama/ukama/systems/common/uuid"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/tj/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Test constants
const (
	defaultMarkupRate = 10.0
	testDBDSN         = "sqlmock_db_0"
	testDriverName    = "postgres"
)

// Test data structures
type testData struct {
	userId uuid.UUID
	markup float64
}

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

// Helper functions for test setup
func setupMarkupTestDB(t *testing.T) (sqlmock.Sqlmock, *gorm.DB) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		DSN:                  testDBDSN,
		DriverName:           testDriverName,
		Conn:                 db,
		PreferSimpleProtocol: true,
	})
	gdb, err := gorm.Open(dialector, &gorm.Config{})
	assert.NoError(t, err)

	return mock, gdb
}

func createMarkupsRepo(gdb *gorm.DB) *markupsRepo {
	return NewMarkupsRepo(&UkamaDbMock{
		GormDb: gdb,
	})
}

func createTestData() testData {
	return testData{
		userId: uuid.NewV4(),
		markup: defaultMarkupRate,
	}
}

func createMockRows(userId uuid.UUID, markup float64) *sqlmock.Rows {
	return sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "owner_id", "markup"}).
		AddRow(1, time.Now(), time.Now(), nil, userId, markup)
}

func TestMarkupRepo_Create(t *testing.T) {

	t.Run("Create", func(t *testing.T) {
		// Arrange
		testData := createTestData()
		mock, gdb := setupMarkupTestDB(t)
		r := createMarkupsRepo(gdb)

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT`)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), testData.userId, testData.markup).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		// Act
		err := r.CreateMarkupRate(testData.userId, testData.markup)

		// Assert
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Create_DatabaseError", func(t *testing.T) {
		// Arrange
		testData := createTestData()
		mock, gdb := setupMarkupTestDB(t)
		r := createMarkupsRepo(gdb)

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT`)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), testData.userId, testData.markup).
			WillReturnError(errors.New("database error"))
		mock.ExpectRollback()

		// Act
		err := r.CreateMarkupRate(testData.userId, testData.markup)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Create_TransactionBeginError", func(t *testing.T) {
		// Arrange
		testData := createTestData()
		mock, gdb := setupMarkupTestDB(t)
		r := createMarkupsRepo(gdb)

		mock.ExpectBegin().WillReturnError(errors.New("transaction begin error"))

		// Act
		err := r.CreateMarkupRate(testData.userId, testData.markup)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "transaction begin error")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestMarkupRepo_Get(t *testing.T) {

	t.Run("Get_Success", func(t *testing.T) {
		// Arrange
		testData := createTestData()
		mock, gdb := setupMarkupTestDB(t)
		r := createMarkupsRepo(gdb)

		row := createMockRows(testData.userId, testData.markup)
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).
			WithArgs(testData.userId, sqlmock.AnyArg()).
			WillReturnRows(row)

		// Act
		m, err := r.GetMarkupRate(testData.userId)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, m)
		assert.EqualValues(t, testData.markup, m.Markup)
		assert.Equal(t, testData.userId.String(), m.OwnerId.String())
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Get_NotFound", func(t *testing.T) {
		// Arrange
		testData := createTestData()
		mock, gdb := setupMarkupTestDB(t)
		r := createMarkupsRepo(gdb)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).
			WithArgs(testData.userId, sqlmock.AnyArg()).
			WillReturnError(gorm.ErrRecordNotFound)

		// Act
		m, err := r.GetMarkupRate(testData.userId)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, m)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Get_DatabaseError", func(t *testing.T) {
		// Arrange
		testData := createTestData()
		mock, gdb := setupMarkupTestDB(t)
		r := createMarkupsRepo(gdb)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).
			WithArgs(testData.userId, sqlmock.AnyArg()).
			WillReturnError(errors.New("database error"))

		// Act
		m, err := r.GetMarkupRate(testData.userId)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, m)
		assert.Contains(t, err.Error(), "database error")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestMarkupRepo_Delete(t *testing.T) {

	t.Run("Delete_Success", func(t *testing.T) {
		// Arrange
		testData := createTestData()
		mock, gdb := setupMarkupTestDB(t)
		r := createMarkupsRepo(gdb)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("UPDATE")).WithArgs(sqlmock.AnyArg(), testData.userId).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		// Act
		err := r.DeleteMarkupRate(testData.userId)

		// Assert
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Delete_DatabaseError", func(t *testing.T) {
		// Arrange
		testData := createTestData()
		mock, gdb := setupMarkupTestDB(t)
		r := createMarkupsRepo(gdb)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("UPDATE")).WithArgs(sqlmock.AnyArg(), testData.userId).
			WillReturnError(errors.New("database error"))
		mock.ExpectRollback()

		// Act
		err := r.DeleteMarkupRate(testData.userId)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Delete_TransactionBeginError", func(t *testing.T) {
		// Arrange
		testData := createTestData()
		mock, gdb := setupMarkupTestDB(t)
		r := createMarkupsRepo(gdb)

		mock.ExpectBegin().WillReturnError(errors.New("transaction begin error"))

		// Act
		err := r.DeleteMarkupRate(testData.userId)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "transaction begin error")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestMarkupRepo_Update(t *testing.T) {

	t.Run("Update_Success", func(t *testing.T) {
		// Arrange
		testData := createTestData()
		mock, gdb := setupMarkupTestDB(t)
		r := createMarkupsRepo(gdb)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("UPDATE")).WithArgs(sqlmock.AnyArg(), testData.userId).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT`)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), testData.userId, testData.markup).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		// Act
		err := r.UpdateMarkupRate(testData.userId, testData.markup)

		// Assert
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Update_DeleteError", func(t *testing.T) {
		// Arrange
		testData := createTestData()
		mock, gdb := setupMarkupTestDB(t)
		r := createMarkupsRepo(gdb)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("UPDATE")).WithArgs(sqlmock.AnyArg(), testData.userId).
			WillReturnError(errors.New("delete error"))
		mock.ExpectRollback()

		// Act
		err := r.UpdateMarkupRate(testData.userId, testData.markup)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "delete error")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Update_CreateError", func(t *testing.T) {
		// Arrange
		testData := createTestData()
		mock, gdb := setupMarkupTestDB(t)
		r := createMarkupsRepo(gdb)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("UPDATE")).WithArgs(sqlmock.AnyArg(), testData.userId).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT`)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), testData.userId, testData.markup).
			WillReturnError(errors.New("create error"))
		mock.ExpectRollback()

		// Act
		err := r.UpdateMarkupRate(testData.userId, testData.markup)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "create error")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestMarkupRepo_GetHistory(t *testing.T) {
	t.Run("GetHistory_Success", func(t *testing.T) {
		// Arrange
		testData := createTestData()
		mock, gdb := setupMarkupTestDB(t)
		r := createMarkupsRepo(gdb)

		row := createMockRows(testData.userId, testData.markup)
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).
			WithArgs(testData.userId).
			WillReturnRows(row)

		// Act
		m, err := r.GetMarkupRateHistory(testData.userId)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, m)
		assert.EqualValues(t, testData.markup, m[0].Markup)
		assert.Equal(t, testData.userId.String(), m[0].OwnerId.String())
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GetHistory_DatabaseError", func(t *testing.T) {
		// Arrange
		testData := createTestData()
		mock, gdb := setupMarkupTestDB(t)
		r := createMarkupsRepo(gdb)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).
			WithArgs(testData.userId).
			WillReturnError(errors.New("db error"))

		// Act
		m, err := r.GetMarkupRateHistory(testData.userId)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, m)
		assert.Contains(t, err.Error(), "db error")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
