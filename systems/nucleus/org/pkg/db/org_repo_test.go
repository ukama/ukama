/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package db_test

import (
	"database/sql"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/tj/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/common/uuid"

	log "github.com/sirupsen/logrus"
	org_db "github.com/ukama/ukama/systems/nucleus/org/pkg/db"
)

// Test data constants
const (
	testOrgName        = "ukama"
	testOrgCert        = "ukama_certs"
	testOrgCertShort   = "ukamacert"
	testOrgCountry     = "us"
	testOrgCurrency    = "usd"
	testInvalidOrgName = "Invalid Name With Spaces!"
	testOrgNotFound    = "lol"
	testUserId         = 1
	testUserIdNotFound = 999
	testActiveCount    = int64(5)
	testDeactiveCount  = int64(2)
)

// Test data variables
var (
	testOrgId     = uuid.NewV4()
	testOrgId2    = uuid.NewV4()
	testOrgId3    = uuid.NewV4()
	testOwnerId   = uuid.NewV4()
	testOwnerId2  = uuid.NewV4()
	testUserUuid  = uuid.NewV4()
	testUserUuid2 = uuid.NewV4()
	testTime      = time.Now()
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

func setupTestDB(t *testing.T) (sqlmock.Sqlmock, org_db.OrgRepo, func()) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		DSN:                  testDSN,
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})

	gdb, err := gorm.Open(dialector, &gorm.Config{})
	assert.NoError(t, err)

	repo := org_db.NewOrgRepo(&UkamaDbMock{
		GormDb: gdb,
	})

	cleanup := func() {
		if err := db.Close(); err != nil {
			t.Logf("Error closing database: %v", err)
		}
	}

	return mock, repo, cleanup
}

func TestOrgRepo_Add(t *testing.T) {
	mock, r, cleanup := setupTestDB(t)
	defer cleanup()

	t.Run("AddValidOrg", func(t *testing.T) {
		// Arrange
		org := org_db.Org{
			Id:          testOrgId,
			Name:        testOrgName,
			Owner:       testOwnerId,
			Certificate: testOrgCert,
			Country:     testOrgCountry,
			Currency:    testOrgCurrency,
		}

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`INSERT`)).
			WithArgs(org.Id, org.Name, org.Owner, org.Country, org.Currency,
				sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), org.Certificate, sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		// Act
		err := r.Add(&org, nil)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("AddOrgWithInvalidName", func(t *testing.T) {
		// Arrange
		org := org_db.Org{
			Id:          testOrgId2,
			Name:        testInvalidOrgName, // Invalid DNS label
			Owner:       testOwnerId2,
			Certificate: testOrgCert,
			Country:     testOrgCountry,
			Currency:    testOrgCurrency,
		}

		// Act
		err := r.Add(&org, nil)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid name must be less then 253")

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("AddOrgWithNestedFuncSuccess", func(t *testing.T) {
		// Arrange
		org := org_db.Org{
			Id:          testOrgId3,
			Name:        testOrgName,
			Owner:       testOwnerId,
			Certificate: testOrgCert,
			Country:     testOrgCountry,
			Currency:    testOrgCurrency,
		}

		nestedFunc := func(org *org_db.Org, tx *gorm.DB) error {
			// Simulate some nested operation
			return nil
		}

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`INSERT`)).
			WithArgs(org.Id, org.Name, org.Owner, org.Country, org.Currency,
				sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), org.Certificate, sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		// Act
		err := r.Add(&org, nestedFunc)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("AddOrgWithNestedFuncError", func(t *testing.T) {
		// Arrange
		org := org_db.Org{
			Id:          testOrgId,
			Name:        testOrgName,
			Owner:       testOwnerId,
			Certificate: testOrgCert,
			Country:     testOrgCountry,
			Currency:    testOrgCurrency,
		}

		nestedFunc := func(org *org_db.Org, tx *gorm.DB) error {
			return fmt.Errorf("nested function error")
		}

		mock.ExpectBegin()
		mock.ExpectRollback()

		// Act
		err := r.Add(&org, nestedFunc)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "nested function error")

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("AddOrgWithDatabaseError", func(t *testing.T) {
		// Arrange
		org := org_db.Org{
			Id:          testOrgId2,
			Name:        testOrgName,
			Owner:       testOwnerId2,
			Certificate: testOrgCert,
			Country:     testOrgCountry,
			Currency:    testOrgCurrency,
		}

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`INSERT`)).
			WithArgs(org.Id, org.Name, org.Owner, org.Country, org.Currency,
				sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), org.Certificate, sqlmock.AnyArg()).
			WillReturnError(sql.ErrConnDone)

		mock.ExpectRollback()

		// Act
		err := r.Add(&org, nil)

		// Assert
		assert.Error(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestOrgRepo_Get(t *testing.T) {
	mock, r, cleanup := setupTestDB(t)
	defer cleanup()

	t.Run("OrgFound", func(t *testing.T) {
		// Arrange
		rows := sqlmock.NewRows([]string{"id", "name", "owner", "certificate"}).
			AddRow(testOrgId, testOrgName, testOwnerId, testOrgCertShort)

		mock.ExpectQuery(`^SELECT.*orgs.*`).
			WithArgs(testOrgId, sqlmock.AnyArg()).
			WillReturnRows(rows)

		// Act
		org, err := r.Get(testOrgId)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, org)

		assert.Equal(t, testOrgId, org.Id)
		assert.Equal(t, testOrgName, org.Name)
		assert.Equal(t, testOwnerId, org.Owner)
		assert.Equal(t, testOrgCertShort, org.Certificate)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("OrgNotFound", func(t *testing.T) {
		// Arrange
		mock.ExpectQuery(`^SELECT.*orgs.*`).
			WithArgs(testOrgId2, sqlmock.AnyArg()).
			WillReturnError(sql.ErrNoRows)

		// Act
		org, err := r.Get(testOrgId2)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, org)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestOrgRepo_GetByName(t *testing.T) {
	mock, r, cleanup := setupTestDB(t)
	defer cleanup()

	t.Run("OrgFound", func(t *testing.T) {
		// Arrange
		rows := sqlmock.NewRows([]string{"id", "name", "owner", "certificate"}).
			AddRow(testOrgId, testOrgName, testOwnerId, testOrgCertShort)

		mock.ExpectQuery(`^SELECT.*orgs.*`).
			WithArgs(testOrgName, sqlmock.AnyArg()).
			WillReturnRows(rows)

		// Act
		org, err := r.GetByName(testOrgName)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, org)

		assert.Equal(t, testOrgId, org.Id)
		assert.Equal(t, testOrgName, org.Name)
		assert.Equal(t, testOwnerId, org.Owner)
		assert.Equal(t, testOrgCertShort, org.Certificate)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("OrgNotFound", func(t *testing.T) {
		// Arrange
		mock.ExpectQuery(`^SELECT.*orgs.*`).
			WithArgs(testOrgNotFound, sqlmock.AnyArg()).
			WillReturnError(sql.ErrNoRows)

		// Act
		org, err := r.GetByName(testOrgNotFound)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, org)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestOrgRepo_GetByOwner(t *testing.T) {
	mock, r, cleanup := setupTestDB(t)
	defer cleanup()

	t.Run("OwnerFound", func(t *testing.T) {
		// Arrange
		rows := sqlmock.NewRows([]string{"id", "name", "owner", "certificate"}).
			AddRow(testOrgId, testOrgName, testOwnerId, testOrgCertShort)

		mock.ExpectQuery(`^SELECT.*orgs.*`).
			WithArgs(testOwnerId).
			WillReturnRows(rows)

		// Act
		orgs, err := r.GetByOwner(testOwnerId)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, orgs)

		assert.Equal(t, testOrgId, orgs[0].Id)
		assert.Equal(t, testOrgName, orgs[0].Name)
		assert.Equal(t, testOwnerId, orgs[0].Owner)
		assert.Equal(t, testOrgCertShort, orgs[0].Certificate)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("OwnerNotFound", func(t *testing.T) {
		// Arrange
		mock.ExpectQuery(`^SELECT.*orgs.*`).
			WithArgs(testOwnerId2).
			WillReturnError(sql.ErrNoRows)

		// Act
		orgs, err := r.GetByOwner(testOwnerId2)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, orgs)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestOrgRepo_GetByMember(t *testing.T) {
	mock, r, cleanup := setupTestDB(t)
	defer cleanup()

	t.Run("MemberFound", func(t *testing.T) {
		// Arrange
		// Mock the main query for orgs
		rows := sqlmock.NewRows([]string{"id", "name", "owner", "certificate", "country", "currency", "deleted_at", "created_at", "updated_at", "deactivated"}).
			AddRow(testOrgId, testOrgName, testOwnerId, testOrgCertShort, testOrgCountry, testOrgCurrency, nil, testTime, testTime, false)

		mock.ExpectQuery(`^SELECT.*orgs.*`).
			WillReturnRows(rows)

		// Mock the preload query for users - this is a separate query for the association
		userRows := sqlmock.NewRows([]string{"id", "uuid", "deactivated", "deleted_at"}).
			AddRow(testUserId, testUserUuid, false, nil)

		mock.ExpectQuery(`^SELECT.*org_users.*`).
			WithArgs(testOrgId).
			WillReturnRows(userRows)

		// Act
		orgs, err := r.GetByMember(testUserId)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, orgs)
		assert.Equal(t, 1, len(orgs))
		assert.Equal(t, testOrgId, orgs[0].Id)
		assert.Equal(t, testOrgName, orgs[0].Name)
		assert.Equal(t, testOwnerId, orgs[0].Owner)
		assert.Equal(t, testOrgCertShort, orgs[0].Certificate)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("MemberNotFound", func(t *testing.T) {
		// Arrange
		// Mock the main query for orgs (empty result)
		rows := sqlmock.NewRows([]string{"id", "name", "owner", "certificate", "country", "currency", "deleted_at", "created_at", "updated_at", "deactivated"})

		mock.ExpectQuery(`^SELECT.*orgs.*`).
			WillReturnRows(rows)

		// Act
		orgs, err := r.GetByMember(testUserIdNotFound)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, orgs)
		assert.Equal(t, 0, len(orgs))

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("GetByMemberDatabaseError", func(t *testing.T) {
		// Arrange
		// Mock the main query for orgs to return an error
		mock.ExpectQuery(`^SELECT.*orgs.*`).
			WillReturnError(sql.ErrConnDone)

		// Act
		orgs, err := r.GetByMember(testUserId)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, orgs)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestOrgRepo_AddUser(t *testing.T) {
	mock, r, cleanup := setupTestDB(t)
	defer cleanup()

	t.Run("AddUserSuccess", func(t *testing.T) {
		// Arrange
		org := &org_db.Org{
			Id:   testOrgId,
			Name: testOrgName,
		}

		user := &org_db.User{
			Id:   testUserId,
			Uuid: testUserUuid,
		}

		// Mock transaction begin
		mock.ExpectBegin()

		// Mock the org update (GORM updates updated_at field first)
		mock.ExpectExec(`^UPDATE "orgs" SET "updated_at"`).
			WithArgs(sqlmock.AnyArg(), org.Id).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Mock the user insert (GORM tries to insert user)
		mock.ExpectQuery(`^INSERT INTO "users"`).
			WithArgs(user.Uuid, false, nil, user.Id).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(user.Id))

		// Mock the association append operation
		mock.ExpectExec(`^INSERT INTO "org_users"`).
			WithArgs(org.Id, user.Id).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Mock transaction commit
		mock.ExpectCommit()

		// Act
		err := r.AddUser(org, user)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("AddUserError", func(t *testing.T) {
		// Arrange
		org := &org_db.Org{
			Id:   testOrgId2,
			Name: testOrgName,
		}

		user := &org_db.User{
			Id:   testUserId,
			Uuid: testUserUuid2,
		}

		// Mock transaction begin
		mock.ExpectBegin()

		// Mock the org update (GORM updates updated_at field first)
		mock.ExpectExec(`^UPDATE "orgs" SET "updated_at"`).
			WithArgs(sqlmock.AnyArg(), org.Id).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Mock the user insert (GORM tries to insert user)
		mock.ExpectQuery(`^INSERT INTO "users"`).
			WithArgs(user.Uuid, false, nil, user.Id).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(user.Id))

		// Mock the association append operation to return an error
		mock.ExpectExec(`^INSERT INTO "org_users"`).
			WithArgs(org.Id, user.Id).
			WillReturnError(sql.ErrConnDone)

		// Mock transaction rollback
		mock.ExpectRollback()

		// Act
		err := r.AddUser(org, user)

		// Assert
		assert.Error(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestOrgRepo_RemoveUser(t *testing.T) {
	mock, r, cleanup := setupTestDB(t)
	defer cleanup()

	t.Run("RemoveUserSuccess", func(t *testing.T) {
		// Arrange
		org := &org_db.Org{
			Id:   testOrgId3,
			Name: testOrgName,
		}

		user := &org_db.User{
			Id:   testUserId,
			Uuid: testUserUuid,
		}

		// Mock transaction begin
		mock.ExpectBegin()

		// Mock the association delete operation
		mock.ExpectExec(`^DELETE FROM "org_users"`).
			WithArgs(org.Id, user.Id).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Mock transaction commit
		mock.ExpectCommit()

		// Act
		err := r.RemoveUser(org, user)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("RemoveUserError", func(t *testing.T) {
		// Arrange
		org := &org_db.Org{
			Id:   testOrgId,
			Name: testOrgName,
		}

		user := &org_db.User{
			Id:   testUserId,
			Uuid: testUserUuid2,
		}

		// Mock transaction begin
		mock.ExpectBegin()

		// Mock the association delete operation to return an error
		mock.ExpectExec(`^DELETE FROM "org_users"`).
			WithArgs(org.Id, user.Id).
			WillReturnError(sql.ErrConnDone)

		// Mock transaction rollback
		mock.ExpectRollback()

		// Act
		err := r.RemoveUser(org, user)

		// Assert
		assert.Error(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestOrgRepo_GetOrgCount(t *testing.T) {
	mock, r, cleanup := setupTestDB(t)
	defer cleanup()

	t.Run("GetOrgCountSuccess", func(t *testing.T) {
		// Arrange
		// Mock the active org count query
		activeRows := sqlmock.NewRows([]string{"count"}).AddRow(testActiveCount)
		mock.ExpectQuery(`^SELECT count\(\*\) FROM "orgs"`).
			WithArgs(false).
			WillReturnRows(activeRows)

		// Mock the deactive org count query
		deactiveRows := sqlmock.NewRows([]string{"count"}).AddRow(testDeactiveCount)
		mock.ExpectQuery(`^SELECT count\(\*\) FROM "orgs"`).
			WithArgs(true).
			WillReturnRows(deactiveRows)

		// Act
		active, deactive, err := r.GetOrgCount()

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, testActiveCount, active)
		assert.Equal(t, testDeactiveCount, deactive)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("GetOrgCountActiveError", func(t *testing.T) {
		// Arrange
		// Mock the active org count query to return an error
		mock.ExpectQuery(`^SELECT count\(\*\) FROM "orgs"`).
			WithArgs(false).
			WillReturnError(sql.ErrConnDone)

		// Act
		active, deactive, err := r.GetOrgCount()

		// Assert
		assert.Error(t, err)
		assert.Equal(t, int64(0), active)
		assert.Equal(t, int64(0), deactive)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("GetOrgCountDeactiveError", func(t *testing.T) {
		// Arrange
		// Mock the active org count query
		activeRows := sqlmock.NewRows([]string{"count"}).AddRow(testActiveCount)
		mock.ExpectQuery(`^SELECT count\(\*\) FROM "orgs"`).
			WithArgs(false).
			WillReturnRows(activeRows)

		// Mock the deactive org count query to return an error
		mock.ExpectQuery(`^SELECT count\(\*\) FROM "orgs"`).
			WithArgs(true).
			WillReturnError(sql.ErrConnDone)

		// Act
		active, deactive, err := r.GetOrgCount()

		// Assert
		assert.Error(t, err)
		assert.Equal(t, int64(0), active)
		assert.Equal(t, int64(0), deactive)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestOrgRepo_GetAll(t *testing.T) {
	mock, r, cleanup := setupTestDB(t)
	defer cleanup()

	t.Run("GetAll", func(t *testing.T) {
		// Arrange
		rows := sqlmock.NewRows([]string{"id", "name", "owner", "certificate"}).
			AddRow(testOrgId, testOrgName, testOwnerId, testOrgCertShort)

		mock.ExpectQuery(`^SELECT.*orgs.*`).
			WillReturnRows(rows)

		// Act
		orgs, err := r.GetAll()

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, orgs)

		assert.Equal(t, 1, len(orgs))

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("GetAllEmptyResult", func(t *testing.T) {
		// Arrange
		rows := sqlmock.NewRows([]string{"id", "name", "owner", "certificate"})

		mock.ExpectQuery(`^SELECT.*orgs.*`).
			WillReturnRows(rows)

		// Act
		orgs, err := r.GetAll()

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, orgs)
		assert.Equal(t, 0, len(orgs))

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("GetAllDatabaseError", func(t *testing.T) {
		// Arrange
		mock.ExpectQuery(`^SELECT.*orgs.*`).
			WillReturnError(sql.ErrConnDone)

		// Act
		orgs, err := r.GetAll()

		// Assert
		assert.Error(t, err)
		assert.Nil(t, orgs)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}
