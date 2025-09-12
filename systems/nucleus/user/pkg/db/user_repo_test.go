/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package db_test

import (
	"regexp"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/tj/assert"
	"github.com/ukama/ukama/systems/common/uuid"

	extsql "database/sql"

	log "github.com/sirupsen/logrus"
	userdb "github.com/ukama/ukama/systems/nucleus/user/pkg/db"
)

// Test data constants
var (
	// User data
	testUserName     = "John Doe"
	testUserEmail    = "johndoe@example.com"
	testUserPhone    = "00100000000"
	testUpdatedName  = "Fox Doe"
	testUpdatedEmail = "foxdoe@example.com"
	testUpdatedPhone = "00200000000"

	// Additional test user data
	testUserEmail2 = "janedoe@example.com"

	// Database configuration
	testDSN        = "sqlmock_db_0"
	testDriverName = "postgres"

	// Count data
	testActiveUserCount   = int64(2)
	testInactiveUserCount = int64(1)

	// SQL query patterns
	testSelectQueryPattern     = `^SELECT.*users.*`
	testCountQueryPattern      = `^SELECT count(\\*).*users.*`
	testCountWhereQueryPattern = `^SELECT count(\\*).*users.*WHERE.*`
	testUpdateQueryPattern     = `UPDATE "users" SET`
	testInsertQueryPattern     = `INSERT`

	// Database column names
	testUserColumns = []string{"id", "name", "email", "phone", "auth_id"}
	testCountColumn = []string{"count"}

	// Mock result values
	testMockResultRowsAffected = int64(1)
	testMockResultLastInsertId = int64(1)
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

// Helper functions for test setup
func setupTestDB(t *testing.T) (sqlmock.Sqlmock, userdb.UserRepo, func()) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		DSN:                  testDSN,
		DriverName:           testDriverName,
		Conn:                 db,
		PreferSimpleProtocol: true,
	})

	gdb, err := gorm.Open(dialector, &gorm.Config{})
	assert.NoError(t, err)

	repo := userdb.NewUserRepo(&UkamaDbMock{
		GormDb: gdb,
	})

	cleanup := func() {
		if err := db.Close(); err != nil {
			t.Logf("Error closing database: %v", err)
		}
	}

	return mock, repo, cleanup
}

func createTestUser(id, authId uuid.UUID) userdb.User {
	return userdb.User{
		Id:     id,
		Name:   testUserName,
		Email:  testUserEmail,
		Phone:  testUserPhone,
		AuthId: authId,
	}
}

func createTestUserWithData(id, authId uuid.UUID, name, email, phone string) userdb.User {
	return userdb.User{
		Id:     id,
		Name:   name,
		Email:  email,
		Phone:  phone,
		AuthId: authId,
	}
}

func TestUserRepo_Add(t *testing.T) {
	// Arrange
	userId := uuid.NewV4()
	authId := uuid.NewV4()
	user := createTestUser(userId, authId)

	mock, r, cleanup := setupTestDB(t)
	defer cleanup()

	t.Run("AddUser", func(t *testing.T) {
		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(testInsertQueryPattern)).
			WithArgs(user.Id, user.Name, user.Email, user.Phone, sqlmock.AnyArg(),
				user.AuthId, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(testMockResultLastInsertId, testMockResultRowsAffected))

		mock.ExpectCommit()

		// Act
		err := r.Add(&user, nil)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("AddUserDatabaseError", func(t *testing.T) {
		// Arrange
		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(testInsertQueryPattern)).
			WithArgs(user.Id, user.Name, user.Email, user.Phone, sqlmock.AnyArg(),
				user.AuthId, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnError(extsql.ErrConnDone)

		mock.ExpectRollback()

		// Act
		err := r.Add(&user, nil)

		// Assert
		assert.Error(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestUserRepo_Get(t *testing.T) {
	// Arrange
	mock, r, cleanup := setupTestDB(t)
	defer cleanup()

	t.Run("UserFound", func(t *testing.T) {
		// Arrange
		userId := uuid.NewV4()
		authId := uuid.NewV4()

		rows := sqlmock.NewRows(testUserColumns).
			AddRow(userId, testUserName, testUserEmail, testUserPhone, authId)

		mock.ExpectQuery(testSelectQueryPattern).
			WithArgs(userId, sqlmock.AnyArg()).
			WillReturnRows(rows)

		// Act
		usr, err := r.Get(userId)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, usr)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)

		assert.Equal(t, testUserName, usr.Name)
		assert.Equal(t, testUserEmail, usr.Email)
		assert.Equal(t, testUserPhone, usr.Phone)
		assert.Equal(t, authId, usr.AuthId)
	})

	t.Run("userNotFound", func(t *testing.T) {
		// Arrange
		userId := uuid.NewV4()

		mock.ExpectQuery(testSelectQueryPattern).
			WithArgs(userId, sqlmock.AnyArg()).
			WillReturnError(extsql.ErrNoRows)

		// Act
		usr, err := r.Get(userId)

		// Assert
		assert.Error(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.Nil(t, usr)
	})
}

func TestUserRepo_GetByAuthId(t *testing.T) {
	// Arrange
	mock, r, cleanup := setupTestDB(t)
	defer cleanup()

	t.Run("UserFound", func(t *testing.T) {
		// Arrange
		userId := uuid.NewV4()
		authId := uuid.NewV4()

		rows := sqlmock.NewRows(testUserColumns).
			AddRow(userId, testUserName, testUserEmail, testUserPhone, authId)

		mock.ExpectQuery(testSelectQueryPattern).
			WithArgs(authId, sqlmock.AnyArg()).
			WillReturnRows(rows)

		// Act
		usr, err := r.GetByAuthId(authId)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.NotNil(t, usr)

		assert.Equal(t, testUserName, usr.Name)
		assert.Equal(t, testUserEmail, usr.Email)
		assert.Equal(t, testUserPhone, usr.Phone)
		assert.Equal(t, authId, usr.AuthId)
	})

	t.Run("userNotFound", func(t *testing.T) {
		// Arrange
		authId := uuid.NewV4()

		mock.ExpectQuery(testSelectQueryPattern).
			WithArgs(authId, sqlmock.AnyArg()).
			WillReturnError(extsql.ErrNoRows)

		// Act
		usr, err := r.GetByAuthId(authId)

		// Assert
		assert.Error(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.Nil(t, usr)
	})
}

func TestUserRepo_Delete(t *testing.T) {
	// Arrange
	mock, r, cleanup := setupTestDB(t)
	defer cleanup()

	t.Run("UserFound", func(t *testing.T) {
		// Arrange
		userId := uuid.NewV4()

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(testUpdateQueryPattern)).
			WithArgs(sqlmock.AnyArg(), userId).
			WillReturnResult(sqlmock.NewResult(testMockResultLastInsertId, testMockResultRowsAffected))

		mock.ExpectCommit()

		// Act
		err := r.Delete(userId, nil)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		// Arrange
		userId := uuid.NewV4()

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(testUpdateQueryPattern)).
			WithArgs(sqlmock.AnyArg(), userId).
			WillReturnError(extsql.ErrNoRows)

		// Act
		err := r.Delete(userId, nil)

		// Assert
		assert.Error(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("DeleteUserDatabaseError", func(t *testing.T) {
		// Arrange
		userId := uuid.NewV4()

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(testUpdateQueryPattern)).
			WithArgs(sqlmock.AnyArg(), userId).
			WillReturnError(extsql.ErrConnDone)

		mock.ExpectRollback()

		// Act
		err := r.Delete(userId, nil)

		// Assert
		assert.Error(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestUserRepo_GetUserCount(t *testing.T) {
	// Arrange
	mock, r, cleanup := setupTestDB(t)
	defer cleanup()

	t.Run("UserFound", func(t *testing.T) {
		// Arrange
		rowsCount1 := sqlmock.NewRows(testCountColumn).
			AddRow(testActiveUserCount)

		rowsCount2 := sqlmock.NewRows(testCountColumn).
			AddRow(testInactiveUserCount)

		mock.ExpectQuery(testCountQueryPattern).
			WillReturnRows(rowsCount1)

		mock.ExpectQuery(testCountWhereQueryPattern).
			WillReturnRows(rowsCount2)

		// Act
		activeUsr, inactiveUsr, err := r.GetUserCount()
		assert.NoError(t, err)

		// Assert
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)

		assert.Equal(t, testActiveUserCount, activeUsr)
		assert.Equal(t, testInactiveUserCount, inactiveUsr)
	})

	t.Run("GetUserCountWithZeroUsers", func(t *testing.T) {
		// Arrange
		zeroCount := int64(0)
		rowsCount1 := sqlmock.NewRows(testCountColumn).
			AddRow(zeroCount)

		rowsCount2 := sqlmock.NewRows(testCountColumn).
			AddRow(zeroCount)

		mock.ExpectQuery(testCountQueryPattern).
			WillReturnRows(rowsCount1)

		mock.ExpectQuery(testCountWhereQueryPattern).
			WillReturnRows(rowsCount2)

		// Act
		activeUsr, inactiveUsr, err := r.GetUserCount()
		assert.NoError(t, err)

		// Assert
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)

		assert.Equal(t, zeroCount, activeUsr)
		assert.Equal(t, zeroCount, inactiveUsr)
	})

	t.Run("GetUserCountDatabaseError", func(t *testing.T) {
		// Arrange
		mock.ExpectQuery(testCountQueryPattern).
			WillReturnError(extsql.ErrConnDone)

		// Act
		activeUsr, inactiveUsr, err := r.GetUserCount()

		// Assert
		assert.Error(t, err)
		assert.Equal(t, int64(0), activeUsr)
		assert.Equal(t, int64(0), inactiveUsr)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("GetUserCountSecondQueryError", func(t *testing.T) {
		// Arrange
		rowsCount1 := sqlmock.NewRows(testCountColumn).
			AddRow(testActiveUserCount)

		mock.ExpectQuery(testCountQueryPattern).
			WillReturnRows(rowsCount1)

		mock.ExpectQuery(testCountWhereQueryPattern).
			WillReturnError(extsql.ErrConnDone)

		// Act
		activeUsr, inactiveUsr, err := r.GetUserCount()

		// Assert
		assert.Error(t, err)
		assert.Equal(t, int64(0), activeUsr)
		assert.Equal(t, int64(0), inactiveUsr)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestUserRepo_GetByEmail(t *testing.T) {
	// Arrange
	mock, r, cleanup := setupTestDB(t)
	defer cleanup()

	t.Run("UserFound", func(t *testing.T) {
		// Arrange
		userId := uuid.NewV4()
		authId := uuid.NewV4()

		rows := sqlmock.NewRows(testUserColumns).
			AddRow(userId, testUserName, testUserEmail, testUserPhone, authId)

		mock.ExpectQuery(testSelectQueryPattern).
			WithArgs(testUserEmail, sqlmock.AnyArg()).
			WillReturnRows(rows)

		// Act
		usr, err := r.GetByEmail(testUserEmail)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, usr)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)

		assert.Equal(t, testUserName, usr.Name)
		assert.Equal(t, testUserEmail, usr.Email)
		assert.Equal(t, testUserPhone, usr.Phone)
		assert.Equal(t, authId, usr.AuthId)
	})

	t.Run("userNotFound", func(t *testing.T) {
		// Arrange
		mock.ExpectQuery(testSelectQueryPattern).
			WithArgs(testUserEmail2, sqlmock.AnyArg()).
			WillReturnError(extsql.ErrNoRows)

		// Act
		usr, err := r.GetByEmail(testUserEmail2)

		// Assert
		assert.Error(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.Nil(t, usr)
	})
}

func TestUserRepo_Update(t *testing.T) {
	// Arrange
	mock, r, cleanup := setupTestDB(t)
	defer cleanup()

	t.Run("UserUpdated", func(t *testing.T) {
		// Arrange
		userId := uuid.NewV4()
		authId := uuid.NewV4()
		user := createTestUserWithData(userId, authId, testUpdatedName, testUpdatedEmail, testUpdatedPhone)

		mock.ExpectBegin()

		mock.ExpectQuery(regexp.QuoteMeta(`UPDATE "users" SET`)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), userId, userId).
			WillReturnRows(sqlmock.NewRows(testUserColumns).
				AddRow(userId, testUpdatedName, testUpdatedEmail, testUpdatedPhone, authId))

		mock.ExpectCommit()

		// Act
		err := r.Update(&user, nil)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		// Arrange
		userId := uuid.NewV4()
		authId := uuid.NewV4()
		user := createTestUserWithData(userId, authId, testUpdatedName, testUpdatedEmail, testUpdatedPhone)

		mock.ExpectBegin()

		mock.ExpectQuery(regexp.QuoteMeta(`UPDATE "users" SET`)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), userId, userId).
			WillReturnRows(sqlmock.NewRows(testUserColumns)) // Empty rows

		mock.ExpectRollback()
		// Act
		err := r.Update(&user, nil)

		// Assert
		assert.Error(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

}
