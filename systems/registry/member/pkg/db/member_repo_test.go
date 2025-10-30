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

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/roles"
	"github.com/ukama/ukama/systems/common/uuid"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/tj/assert"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
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

var orgId = uuid.NewV4()

// Test fixtures and helper functions
type TestMember struct {
	UserId      uuid.UUID
	Role        roles.RoleType
	Deactivated bool
}

// Helper function to create a test member
func createTestMember(role roles.RoleType, deactivated bool) *Member {
	return &Member{
		MemberId:    uuid.NewV4(),
		UserId:      uuid.NewV4(),
		Role:        role,
		Deactivated: deactivated,
	}
}

// Helper function to setup database and mock
func setupTestDB(t *testing.T) (sqlmock.Sqlmock, *gorm.DB, MemberRepo) {
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

	repo := NewMemberRepo(&UkamaDbMock{
		GormDb: gdb,
	})

	return mock, gdb, repo
}

// Helper function to create mock rows for member queries
func createMemberRows(member *Member) *sqlmock.Rows {
	return sqlmock.NewRows([]string{"id", "member_id", "user_id", "role", "deactivated", "created_at", "updated_at", "deleted_at"}).
		AddRow(1, member.MemberId, member.UserId, member.Role, member.Deactivated, nil, nil, nil)
}

func Test_AddMember(t *testing.T) {
	t.Run("AddMember", func(t *testing.T) {
		// Arrange
		member := createTestMember(roles.TYPE_USERS, false)
		mock, _, repo := setupTestDB(t)

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT`)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), member.MemberId, member.UserId, false, member.Role).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		// Act
		err := repo.AddMember(member, orgId.String(), nil)

		// Assert
		assert.NoError(t, err)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("AddMemberWithNestedFunc", func(t *testing.T) {
		// Arrange
		member := createTestMember(roles.TYPE_USERS, false)
		mock, _, repo := setupTestDB(t)

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT`)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), member.MemberId, member.UserId, false, member.Role).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		// Act - with nested function
		nestedFunc := func(orgId, userId string) error {
			return nil // Simulate successful nested operation
		}
		err := repo.AddMember(member, orgId.String(), nestedFunc)

		// Assert
		assert.NoError(t, err)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("AddMemberWithNestedFuncError", func(t *testing.T) {
		// Arrange
		member := createTestMember(roles.TYPE_USERS, false)
		mock, _, repo := setupTestDB(t)

		mock.ExpectBegin()
		mock.ExpectRollback()

		// Act - with nested function that returns error
		nestedFunc := func(orgId, userId string) error {
			return errors.New("test error") // Simulate nested operation failure
		}
		err := repo.AddMember(member, orgId.String(), nestedFunc)

		// Assert
		assert.Error(t, err)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("AddMemberWithDatabaseError", func(t *testing.T) {
		// Arrange
		member := createTestMember(roles.TYPE_USERS, false)
		mock, _, repo := setupTestDB(t)

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT`)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), member.MemberId, member.UserId, false, member.Role).
			WillReturnError(errors.New("database constraint violation"))
		mock.ExpectRollback()

		// Act
		err := repo.AddMember(member, orgId.String(), nil)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database constraint violation")
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func Test_GetMember(t *testing.T) {
	t.Run("MemberExist", func(t *testing.T) {
		// Arrange
		member := createTestMember(roles.TYPE_USERS, false)
		mock, _, repo := setupTestDB(t)

		rows := createMemberRows(member)
		mock.ExpectQuery(`^SELECT.*members.*`).
			WithArgs(member.MemberId, sqlmock.AnyArg()).
			WillReturnRows(rows)

		// Act
		rm, err := repo.GetMember(member.MemberId)

		// Assert
		assert.NoError(t, err)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.NotNil(t, rm)
	})

	t.Run("MemberNotFound", func(t *testing.T) {
		// Arrange
		memberId := uuid.NewV4()
		mock, _, repo := setupTestDB(t)

		mock.ExpectQuery(`^SELECT.*members.*`).
			WithArgs(memberId, sqlmock.AnyArg()).
			WillReturnError(gorm.ErrRecordNotFound)

		// Act
		rm, err := repo.GetMember(memberId)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, rm)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func Test_GetMembers(t *testing.T) {
	t.Run("MembersOfAnOrg", func(t *testing.T) {
		// Arrange
		member := createTestMember(roles.TYPE_USERS, false)
		mock, _, repo := setupTestDB(t)

		rows := createMemberRows(member)
		mock.ExpectQuery(`^SELECT.*members.*`).
			WillReturnRows(rows)

		// Act
		members, err := repo.GetMembers()

		// Assert
		assert.NoError(t, err)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.NotNil(t, members)
	})

	t.Run("GetMembersError", func(t *testing.T) {
		// Arrange
		mock, _, repo := setupTestDB(t)

		mock.ExpectQuery(`^SELECT.*members.*`).
			WillReturnError(errors.New("test error"))

		// Act
		members, err := repo.GetMembers()

		// Assert
		assert.Error(t, err)
		assert.Nil(t, members)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func Test_RemoveMember(t *testing.T) {
	t.Run("RemoveMemberOfAnOrg", func(t *testing.T) {
		// Arrange
		member := createTestMember(roles.TYPE_USERS, false)
		mock, _, repo := setupTestDB(t)

		mock.ExpectBegin()
		mock.ExpectExec(`^UPDATE.*members.*`).
			WithArgs(sqlmock.AnyArg(), member.MemberId).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		// Act
		err := repo.RemoveMember(member.MemberId, orgId.String(), nil)

		// Assert
		assert.NoError(t, err)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("RemoveMemberNotFound", func(t *testing.T) {
		// Arrange
		memberId := uuid.NewV4()
		mock, _, repo := setupTestDB(t)

		mock.ExpectBegin()
		mock.ExpectExec(`^UPDATE.*members.*`).
			WithArgs(sqlmock.AnyArg(), memberId).
			WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectRollback()

		// Act
		err := repo.RemoveMember(memberId, orgId.String(), nil)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("RemoveMemberWithError", func(t *testing.T) {
		// Arrange
		memberId := uuid.NewV4()
		mock, _, repo := setupTestDB(t)

		mock.ExpectBegin()
		mock.ExpectExec(`^UPDATE.*members.*`).
			WithArgs(sqlmock.AnyArg(), memberId).
			WillReturnError(errors.New("test error"))
		mock.ExpectRollback()

		// Act
		err := repo.RemoveMember(memberId, orgId.String(), nil)

		// Assert
		assert.Error(t, err)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func Test_GetMemberByUserId(t *testing.T) {
	t.Run("MemberExists", func(t *testing.T) {
		// Arrange
		member := createTestMember(roles.TYPE_USERS, false)
		mock, _, repo := setupTestDB(t)

		rows := createMemberRows(member)
		mock.ExpectQuery(`^SELECT.*members.*`).
			WithArgs(member.UserId, sqlmock.AnyArg()).
			WillReturnRows(rows)

		// Act
		rm, err := repo.GetMemberByUserId(member.UserId)

		// Assert
		assert.NoError(t, err)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.NotNil(t, rm)
		assert.Equal(t, member.UserId, rm.UserId)
	})

	t.Run("MemberNotFound", func(t *testing.T) {
		// Arrange
		userId := uuid.NewV4()
		mock, _, repo := setupTestDB(t)

		mock.ExpectQuery(`^SELECT.*members.*`).
			WithArgs(userId, sqlmock.AnyArg()).
			WillReturnError(gorm.ErrRecordNotFound)

		// Act
		rm, err := repo.GetMemberByUserId(userId)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, rm)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func Test_GetMemberCount(t *testing.T) {
	t.Run("GetMemberCountSuccess", func(t *testing.T) {
		// Arrange
		mock, _, repo := setupTestDB(t)

		// Mock active member count query
		mock.ExpectQuery(`^SELECT.*count.*members.*`).
			WithArgs(false).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(5))

		// Mock deactivated member count query
		mock.ExpectQuery(`^SELECT.*count.*members.*`).
			WithArgs(true).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

		// Act
		activeCount, deactiveCount, err := repo.GetMemberCount()

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, int64(5), activeCount)
		assert.Equal(t, int64(2), deactiveCount)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("GetMemberCountActiveError", func(t *testing.T) {
		// Arrange
		mock, _, repo := setupTestDB(t)

		// Mock active member count query with error
		mock.ExpectQuery(`^SELECT.*count.*members.*`).
			WithArgs(false).
			WillReturnError(errors.New("test error"))

		// Act
		activeCount, deactiveCount, err := repo.GetMemberCount()

		// Assert
		assert.Error(t, err)
		assert.Equal(t, int64(0), activeCount)
		assert.Equal(t, int64(0), deactiveCount)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("GetMemberCountDeactiveError", func(t *testing.T) {
		// Arrange
		mock, _, repo := setupTestDB(t)

		// Mock active member count query success
		mock.ExpectQuery(`^SELECT.*count.*members.*`).
			WithArgs(false).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(5))

		// Mock deactivated member count query with error
		mock.ExpectQuery(`^SELECT.*count.*members.*`).
			WithArgs(true).
			WillReturnError(errors.New("test error"))

		// Act
		activeCount, deactiveCount, err := repo.GetMemberCount()

		// Assert
		assert.Error(t, err)
		assert.Equal(t, int64(0), activeCount)
		assert.Equal(t, int64(0), deactiveCount)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func Test_UpdateMember(t *testing.T) {
	t.Run("UpdateMemberSuccess", func(t *testing.T) {
		// Arrange - Use a real in-memory SQLite database for integration test
		db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		// Auto-migrate the Member table
		err = db.AutoMigrate(&Member{})
		assert.NoError(t, err)

		repo := NewMemberRepo(&UkamaDbMock{GormDb: db})

		// Create a member first
		member := createTestMember(roles.TYPE_USERS, false)
		err = repo.AddMember(member, orgId.String(), nil)
		assert.NoError(t, err)

		// Act - Update the member
		member.Role = roles.TYPE_ADMIN
		member.Deactivated = true
		err = repo.UpdateMember(member)

		// Assert
		assert.NoError(t, err)

		// Verify the update worked
		updatedMember, err := repo.GetMember(member.MemberId)
		assert.NoError(t, err)
		assert.Equal(t, roles.TYPE_ADMIN, updatedMember.Role)
		assert.Equal(t, true, updatedMember.Deactivated)
	})

	t.Run("UpdateMemberNotFound", func(t *testing.T) {
		// Arrange - Use a real in-memory SQLite database for integration test
		db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		// Auto-migrate the Member table
		err = db.AutoMigrate(&Member{})
		assert.NoError(t, err)

		repo := NewMemberRepo(&UkamaDbMock{GormDb: db})

		// Create a member that doesn't exist in the database
		member := createTestMember(roles.TYPE_USERS, false)

		// Act - Try to update a non-existent member
		err = repo.UpdateMember(member)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("UpdateMemberRoleChange", func(t *testing.T) {
		// Arrange - Use a real in-memory SQLite database for integration test
		db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		// Auto-migrate the Member table
		err = db.AutoMigrate(&Member{})
		assert.NoError(t, err)

		repo := NewMemberRepo(&UkamaDbMock{GormDb: db})

		// Create a member first
		member := createTestMember(roles.TYPE_USERS, false)
		err = repo.AddMember(member, orgId.String(), nil)
		assert.NoError(t, err)

		// Act - Update only the role
		member.Role = roles.TYPE_ADMIN
		err = repo.UpdateMember(member)

		// Assert
		assert.NoError(t, err)

		// Verify the role update worked
		updatedMember, err := repo.GetMember(member.MemberId)
		assert.NoError(t, err)
		assert.Equal(t, roles.TYPE_ADMIN, updatedMember.Role)
		assert.Equal(t, false, updatedMember.Deactivated) // Should remain unchanged
	})

	t.Run("UpdateMemberDeactivatedStatus", func(t *testing.T) {
		// Arrange - Use a real in-memory SQLite database for integration test
		db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		// Auto-migrate the Member table
		err = db.AutoMigrate(&Member{})
		assert.NoError(t, err)

		repo := NewMemberRepo(&UkamaDbMock{GormDb: db})

		// Create a member first
		member := createTestMember(roles.TYPE_USERS, false)
		err = repo.AddMember(member, orgId.String(), nil)
		assert.NoError(t, err)

		// Act - Update only the deactivated status
		member.Deactivated = true
		err = repo.UpdateMember(member)

		// Assert
		assert.NoError(t, err)

		// Verify the deactivated status update worked
		updatedMember, err := repo.GetMember(member.MemberId)
		assert.NoError(t, err)
		assert.Equal(t, true, updatedMember.Deactivated)
		assert.Equal(t, roles.TYPE_USERS, updatedMember.Role) // Should remain unchanged
	})

	t.Run("UpdateMemberMultipleFields", func(t *testing.T) {
		// Arrange - Use a real in-memory SQLite database for integration test
		db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		// Auto-migrate the Member table
		err = db.AutoMigrate(&Member{})
		assert.NoError(t, err)

		repo := NewMemberRepo(&UkamaDbMock{GormDb: db})

		// Create a member first
		member := createTestMember(roles.TYPE_USERS, false)
		err = repo.AddMember(member, orgId.String(), nil)
		assert.NoError(t, err)

		// Act - Update multiple fields
		member.Role = roles.TYPE_ADMIN
		member.Deactivated = true
		err = repo.UpdateMember(member)

		// Assert
		assert.NoError(t, err)

		// Verify both fields were updated
		updatedMember, err := repo.GetMember(member.MemberId)
		assert.NoError(t, err)
		assert.Equal(t, roles.TYPE_ADMIN, updatedMember.Role)
		assert.Equal(t, true, updatedMember.Deactivated)
	})

}
