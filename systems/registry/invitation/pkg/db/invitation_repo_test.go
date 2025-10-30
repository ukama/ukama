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
	"time"

	"github.com/ukama/ukama/systems/common/pb/gen/ukama"
	"github.com/ukama/ukama/systems/common/roles"
	"github.com/ukama/ukama/systems/common/uuid"
	db_inv "github.com/ukama/ukama/systems/registry/invitation/pkg/db"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/tj/assert"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Test constants and data
const (
	// Test emails
	testEmail1       = "test@ukama.com"
	testEmail2       = "test1@ukama.com"
	testEmail3       = "test2@ukama.com"
	nonExistentEmail = "nonexistent@ukama.com"

	// Test names
	testName1 = "test"
	testName2 = "test1"
	testName3 = "test2"

	// Test links
	testLinkBase = "https://ukama.com/invitation/accept/"

	// Database configuration
	testDSN        = "sqlmock_db_0"
	testDriverName = "postgres"

	// SQL query patterns
	selectQueryPattern = `^SELECT.*invitations.*`
	insertQueryPattern = `INSERT INTO "invitations"`
	updateQueryPattern = `UPDATE "invitations"`
	updateStatusQuery  = `UPDATE "invitations" SET "status"=$1,"updated_at"=$2 WHERE id = $3 AND "invitations"."deleted_at" IS NULL`
	updateUserIdQuery  = `UPDATE "invitations" SET "user_id"=$1,"updated_at"=$2 WHERE id = $3 AND "invitations"."deleted_at" IS NULL`
)

// Test data variables
var (
	// Test UUIDs
	testInvitationId1 = uuid.NewV4()
	testUserId1       = uuid.NewV4()

	// Test time
	testExpiryTime = time.Date(2023, 8, 25, 17, 59, 43, 831000000, time.UTC)
	testCreatedAt  = time.Now()
	testUpdatedAt  = time.Now()

	// Test roles and statuses
	testRole           = roles.TYPE_ADMIN
	testStatus         = ukama.InvitationStatus_INVITE_PENDING

	// Database column names
	dbColumns = []string{"id", "name", "email", "role", "status", "user_id", "expires_at", "link", "created_at", "updated_at", "deleted_at"}
)

type UkamaDbMock struct {
	GormDb *gorm.DB
}

func (u UkamaDbMock) Init(model ...interface{}) error {
	return nil
}

func (u UkamaDbMock) Connect() error {
	return nil
}

func (u UkamaDbMock) GetGormDb() *gorm.DB {
	return u.GormDb
}

func (u UkamaDbMock) InitDB() error {
	return nil
}

func (u UkamaDbMock) ExecuteInTransaction(dbOperation func(tx *gorm.DB) *gorm.DB,
	nestedFuncs ...func() error) error {
	return nil
}

func (u UkamaDbMock) ExecuteInTransaction2(dbOperation func(tx *gorm.DB) *gorm.DB,
	nestedFuncs ...func(tx *gorm.DB) error) error {
	return nil
}

// Test data helpers
func createTestInvitation(id uuid.UUID, email, name string) *db_inv.Invitation {
	return &db_inv.Invitation{
		Id:        id,
		Name:      name,
		Email:     email,
		Role:      testRole,
		Status:    testStatus,
		UserId:    testUserId1.String(),
		ExpiresAt: testExpiryTime,
		Link:      testLinkBase + uuid.NewV4().String(),
		CreatedAt: testCreatedAt,
		UpdatedAt: testUpdatedAt,
		DeletedAt: gorm.DeletedAt{},
	}
}

func createDefaultTestInvitation() *db_inv.Invitation {
	return createTestInvitation(
		testInvitationId1,
		testEmail1,
		testName1,
	)
}

// Database setup helpers
func setupTestDB(t *testing.T) (sqlmock.Sqlmock, *gorm.DB, db_inv.InvitationRepo) {
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

	repo := db_inv.NewInvitationRepo(&UkamaDbMock{
		GormDb: gdb,
	})

	return mock, gdb, repo
}

func setupTestDBWithInvitation(t *testing.T, invitation *db_inv.Invitation) (sqlmock.Sqlmock, *gorm.DB, db_inv.InvitationRepo) {
	mock, gdb, repo := setupTestDB(t)

	rows := sqlmock.NewRows(dbColumns).
		AddRow(invitation.Id, invitation.Name, invitation.Email, invitation.Role, invitation.Status,
			invitation.UserId, invitation.ExpiresAt, invitation.Link, invitation.CreatedAt, invitation.UpdatedAt, invitation.DeletedAt)

	mock.ExpectQuery(selectQueryPattern).
		WithArgs(invitation.Id, sqlmock.AnyArg()).
		WillReturnRows(rows)

	return mock, gdb, repo
}

func TestInvitationRepo_AddInvitation(t *testing.T) {
	invitation := createDefaultTestInvitation()
	mock, _, r := setupTestDB(t)

	t.Run("AddValidInvitation", func(t *testing.T) {
		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(insertQueryPattern)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
				sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		// Act
		err := r.Add(invitation, nil)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("AddInvitationWithNestedFuncSuccess", func(t *testing.T) {
		// Arrange
		invitation := createDefaultTestInvitation()
		nestedFuncCalled := false
		nestedFunc := func(inv *db_inv.Invitation, tx *gorm.DB) error {
			nestedFuncCalled = true
			assert.Equal(t, invitation.Id, inv.Id)
			assert.NotNil(t, tx)
			return nil
		}

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(insertQueryPattern)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
				sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		// Act
		err := r.Add(invitation, nestedFunc)

		// Assert
		assert.NoError(t, err)
		assert.True(t, nestedFuncCalled)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("AddInvitationWithNestedFuncError", func(t *testing.T) {
		// Arrange
		invitation := createDefaultTestInvitation()
		expectedError := gorm.ErrInvalidTransaction
		nestedFunc := func(inv *db_inv.Invitation, tx *gorm.DB) error {
			return expectedError
		}

		mock.ExpectBegin()
		mock.ExpectRollback()

		// Act
		err := r.Add(invitation, nestedFunc)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("AddInvitationDatabaseError", func(t *testing.T) {
		// Arrange
		invitation := createDefaultTestInvitation()
		expectedError := gorm.ErrDuplicatedKey

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "invitations"`)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
				sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnError(expectedError)

		mock.ExpectRollback()

		// Act
		err := r.Add(invitation, nil)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestInvitationRepo_GetInvitation(t *testing.T) {
	t.Run("InvitationExist", func(t *testing.T) {
		// Arrange
		invitation := createDefaultTestInvitation()
		mock, _, r := setupTestDBWithInvitation(t, invitation)

		// Act
		rm, err := r.Get(invitation.Id)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, rm)
		assert.Equal(t, invitation.Id, rm.Id)
		assert.Equal(t, invitation.Email, rm.Email)
		assert.Equal(t, invitation.Name, rm.Name)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("InvitationNotFound", func(t *testing.T) {
		// Arrange
		nonExistentId := uuid.NewV4()
		mock, _, r := setupTestDB(t)

		mock.ExpectQuery(`^SELECT.*invitations.*`).
			WithArgs(nonExistentId, sqlmock.AnyArg()).
			WillReturnError(gorm.ErrRecordNotFound)

		// Act
		rm, err := r.Get(nonExistentId)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, rm)
		assert.Equal(t, gorm.ErrRecordNotFound, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestInvitationRepo_GetByEmail(t *testing.T) {
	t.Run("GetByEmailSuccess", func(t *testing.T) {
		// Arrange
		invitation := createDefaultTestInvitation()
		mock, _, r := setupTestDB(t)

		rows := sqlmock.NewRows([]string{"id", "name", "email", "role", "status", "user_id", "expires_at", "link", "created_at", "updated_at", "deleted_at"}).
			AddRow(invitation.Id, invitation.Name, invitation.Email, invitation.Role, invitation.Status,
				invitation.UserId, invitation.ExpiresAt, invitation.Link, invitation.CreatedAt, invitation.UpdatedAt, invitation.DeletedAt)

		mock.ExpectQuery(`^SELECT.*invitations.*`).
			WithArgs(invitation.Email, sqlmock.AnyArg()).
			WillReturnRows(rows)

		// Act
		result, err := r.GetByEmail(invitation.Email)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, invitation.Id, result.Id)
		assert.Equal(t, invitation.Email, result.Email)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("GetByEmailNotFound", func(t *testing.T) {
		// Arrange
		nonExistentEmail := "nonexistent@ukama.com"
		mock, _, r := setupTestDB(t)

		mock.ExpectQuery(`^SELECT.*invitations.*`).
			WithArgs(nonExistentEmail, sqlmock.AnyArg()).
			WillReturnError(gorm.ErrRecordNotFound)

		// Act
		result, err := r.GetByEmail(nonExistentEmail)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, gorm.ErrRecordNotFound, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestInvitationRepo_GetAll(t *testing.T) {
	t.Run("GetAllInvitations", func(t *testing.T) {
		// Arrange
		invitation := createDefaultTestInvitation()
		mock, _, r := setupTestDB(t)

		rows := sqlmock.NewRows([]string{"id", "name", "email", "role", "status", "user_id", "expires_at", "link", "created_at", "updated_at", "deleted_at"}).
			AddRow(invitation.Id, invitation.Name, invitation.Email, invitation.Role, invitation.Status,
				invitation.UserId, invitation.ExpiresAt, invitation.Link, invitation.CreatedAt, invitation.UpdatedAt, invitation.DeletedAt)

		mock.ExpectQuery(`^SELECT.*invitations.*`).
			WithArgs().
			WillReturnRows(rows)

		// Act
		rm, err := r.GetAll()

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, rm)
		assert.Len(t, rm, 1)
		assert.Equal(t, invitation.Id, rm[0].Id)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("GetAllInvitationsEmpty", func(t *testing.T) {
		// Arrange
		mock, _, r := setupTestDB(t)

		rows := sqlmock.NewRows([]string{"id", "name", "email", "role", "status", "user_id", "expires_at", "link", "created_at", "updated_at", "deleted_at"})

		mock.ExpectQuery(`^SELECT.*invitations.*`).
			WithArgs().
			WillReturnRows(rows)

		// Act
		rm, err := r.GetAll()

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, rm)
		assert.Len(t, rm, 0)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("GetAllInvitationsMultiple", func(t *testing.T) {
		// Arrange
		invitation1 := createTestInvitation(uuid.NewV4(), "test1@ukama.com", "test1")
		invitation2 := createTestInvitation(uuid.NewV4(), "test2@ukama.com", "test2")
		mock, _, r := setupTestDB(t)

		rows := sqlmock.NewRows([]string{"id", "name", "email", "role", "status", "user_id", "expires_at", "link", "created_at", "updated_at", "deleted_at"}).
			AddRow(invitation1.Id, invitation1.Name, invitation1.Email, invitation1.Role, invitation1.Status,
				invitation1.UserId, invitation1.ExpiresAt, invitation1.Link, invitation1.CreatedAt, invitation1.UpdatedAt, invitation1.DeletedAt).
			AddRow(invitation2.Id, invitation2.Name, invitation2.Email, invitation2.Role, invitation2.Status,
				invitation2.UserId, invitation2.ExpiresAt, invitation2.Link, invitation2.CreatedAt, invitation2.UpdatedAt, invitation2.DeletedAt)

		mock.ExpectQuery(`^SELECT.*invitations.*`).
			WithArgs().
			WillReturnRows(rows)

		// Act
		rm, err := r.GetAll()

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, rm)
		assert.Len(t, rm, 2)
		assert.Equal(t, invitation1.Id, rm[0].Id)
		assert.Equal(t, invitation2.Id, rm[1].Id)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("GetAllInvitationsDatabaseError", func(t *testing.T) {
		// Arrange
		expectedError := gorm.ErrInvalidDB
		mock, _, r := setupTestDB(t)

		mock.ExpectQuery(`^SELECT.*invitations.*`).
			WithArgs().
			WillReturnError(expectedError)

		// Act
		rm, err := r.GetAll()

		// Assert
		assert.Error(t, err)
		assert.Nil(t, rm)
		assert.Equal(t, expectedError, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestInvitationRepo_Delete(t *testing.T) {
	t.Run("DeleteInvitation", func(t *testing.T) {
		// Arrange
		invitation := createDefaultTestInvitation()
		mock, _, r := setupTestDB(t)

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "invitations"`)).
			WithArgs(sqlmock.AnyArg(), invitation.Id).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		// Act
		err := r.Delete(invitation.Id, nil)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("DeleteInvitationWithNestedFuncSuccess", func(t *testing.T) {
		// Arrange
		invitation := createDefaultTestInvitation()
		mock, _, r := setupTestDB(t)
		nestedFuncCalled := false
		nestedFunc := func(string, string) error {
			nestedFuncCalled = true
			return nil
		}

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "invitations"`)).
			WithArgs(sqlmock.AnyArg(), invitation.Id).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		// Act
		err := r.Delete(invitation.Id, nestedFunc)

		// Assert
		assert.NoError(t, err)
		assert.True(t, nestedFuncCalled)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("DeleteInvitationWithNestedFuncError", func(t *testing.T) {
		// Arrange
		invitation := createDefaultTestInvitation()
		mock, _, r := setupTestDB(t)
		expectedError := gorm.ErrInvalidTransaction
		nestedFunc := func(string, string) error {
			return expectedError
		}

		mock.ExpectBegin()
		mock.ExpectRollback()

		// Act
		err := r.Delete(invitation.Id, nestedFunc)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("DeleteInvitationDatabaseError", func(t *testing.T) {
		// Arrange
		invitation := createDefaultTestInvitation()
		mock, _, r := setupTestDB(t)
		expectedError := gorm.ErrInvalidDB

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "invitations"`)).
			WithArgs(sqlmock.AnyArg(), invitation.Id).
			WillReturnError(expectedError)

		mock.ExpectRollback()

		// Act
		err := r.Delete(invitation.Id, nil)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

}

func TestInvitationRepo_UpdateStatus(t *testing.T) {
	t.Run("UpdateStatusSuccess", func(t *testing.T) {
		// Arrange
		invitation := createDefaultTestInvitation()
		mock, _, r := setupTestDB(t)

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "invitations" SET "status"=$1,"updated_at"=$2 WHERE id = $3 AND "invitations"."deleted_at" IS NULL`)).
			WithArgs(int32(ukama.InvitationStatus_INVITE_ACCEPTED), sqlmock.AnyArg(), invitation.Id).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		// Act
		err := r.UpdateStatus(invitation.Id, uint8(ukama.InvitationStatus_INVITE_ACCEPTED))

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("UpdateStatusNotFound", func(t *testing.T) {
		// Arrange
		nonExistentId := uuid.NewV4()
		mock, _, r := setupTestDB(t)

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "invitations" SET "status"=$1,"updated_at"=$2 WHERE id = $3 AND "invitations"."deleted_at" IS NULL`)).
			WithArgs(int32(ukama.InvitationStatus_INVITE_ACCEPTED), sqlmock.AnyArg(), nonExistentId).
			WillReturnResult(sqlmock.NewResult(0, 0))

		mock.ExpectCommit()

		// Act
		err := r.UpdateStatus(nonExistentId, uint8(ukama.InvitationStatus_INVITE_ACCEPTED))

		// Assert
		assert.NoError(t, err) // GORM doesn't return error for 0 rows affected

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("UpdateStatusDatabaseError", func(t *testing.T) {
		// Arrange
		invitation := createDefaultTestInvitation()
		mock, _, r := setupTestDB(t)
		expectedError := gorm.ErrInvalidDB

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "invitations" SET "status"=$1,"updated_at"=$2 WHERE id = $3 AND "invitations"."deleted_at" IS NULL`)).
			WithArgs(int32(ukama.InvitationStatus_INVITE_ACCEPTED), sqlmock.AnyArg(), invitation.Id).
			WillReturnError(expectedError)

		mock.ExpectRollback()

		// Act
		err := r.UpdateStatus(invitation.Id, uint8(ukama.InvitationStatus_INVITE_ACCEPTED))

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("UpdateStatusTransactionBeginError", func(t *testing.T) {
		// Arrange
		invitation := createDefaultTestInvitation()
		mock, _, r := setupTestDB(t)
		expectedError := gorm.ErrInvalidTransaction

		mock.ExpectBegin().WillReturnError(expectedError)

		// Act
		err := r.UpdateStatus(invitation.Id, uint8(ukama.InvitationStatus_INVITE_ACCEPTED))

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestInvitationRepo_UpdateUserId(t *testing.T) {
	t.Run("UpdateUserIdSuccess", func(t *testing.T) {
		// Arrange
		invitation := createDefaultTestInvitation()
		newUserId := uuid.NewV4()
		mock, _, r := setupTestDB(t)

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "invitations" SET "user_id"=$1,"updated_at"=$2 WHERE id = $3 AND "invitations"."deleted_at" IS NULL`)).
			WithArgs(newUserId, sqlmock.AnyArg(), invitation.Id).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		// Act
		err := r.UpdateUserId(invitation.Id, newUserId)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("UpdateUserIdNotFound", func(t *testing.T) {
		// Arrange
		nonExistentId := uuid.NewV4()
		newUserId := uuid.NewV4()
		mock, _, r := setupTestDB(t)

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "invitations" SET "user_id"=$1,"updated_at"=$2 WHERE id = $3 AND "invitations"."deleted_at" IS NULL`)).
			WithArgs(newUserId, sqlmock.AnyArg(), nonExistentId).
			WillReturnResult(sqlmock.NewResult(0, 0))

		mock.ExpectCommit()

		// Act
		err := r.UpdateUserId(nonExistentId, newUserId)

		// Assert
		assert.NoError(t, err) // GORM doesn't return error for 0 rows affected

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("UpdateUserIdDatabaseError", func(t *testing.T) {
		// Arrange
		invitation := createDefaultTestInvitation()
		newUserId := uuid.NewV4()
		mock, _, r := setupTestDB(t)
		expectedError := gorm.ErrInvalidDB

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "invitations" SET "user_id"=$1,"updated_at"=$2 WHERE id = $3 AND "invitations"."deleted_at" IS NULL`)).
			WithArgs(newUserId, sqlmock.AnyArg(), invitation.Id).
			WillReturnError(expectedError)

		mock.ExpectRollback()

		// Act
		err := r.UpdateUserId(invitation.Id, newUserId)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("UpdateUserIdTransactionBeginError", func(t *testing.T) {
		// Arrange
		invitation := createDefaultTestInvitation()
		newUserId := uuid.NewV4()
		mock, _, r := setupTestDB(t)
		expectedError := gorm.ErrInvalidTransaction

		mock.ExpectBegin().WillReturnError(expectedError)

		// Act
		err := r.UpdateUserId(invitation.Id, newUserId)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}
