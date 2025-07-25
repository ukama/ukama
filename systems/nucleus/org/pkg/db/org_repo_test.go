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
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/tj/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/common/uuid"

	log "github.com/sirupsen/logrus"
	org_db "github.com/ukama/ukama/systems/nucleus/org/pkg/db"
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

func TestOrgRepo_Add(t *testing.T) {
	var db *sql.DB

	org := org_db.Org{
		Id:          uuid.NewV4(),
		Name:        "ukama",
		Owner:       uuid.NewV4(),
		Certificate: "ukama_certs",
		Country:     "us",
		Currency:    "usd",
	}

	db, mock, err := sqlmock.New() // mock sql.DB
	assert.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})

	gdb, err := gorm.Open(dialector, &gorm.Config{})
	assert.NoError(t, err)

	r := org_db.NewOrgRepo(&UkamaDbMock{
		GormDb: gdb,
	})

	t.Run("AddValidOrg", func(t *testing.T) {
		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`INSERT`)).
			WithArgs(org.Id, org.Name, org.Owner, org.Country, org.Currency,
				sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), org.Certificate, sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		// Act
		err = r.Add(&org, nil)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestOrgRepo_Get(t *testing.T) {
	var db *sql.DB

	db, mock, err := sqlmock.New() // mock sql.DB
	assert.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})

	gdb, err := gorm.Open(dialector, &gorm.Config{})
	assert.NoError(t, err)

	r := org_db.NewOrgRepo(&UkamaDbMock{
		GormDb: gdb,
	})

	t.Run("OrgFound", func(t *testing.T) {
		// Arrange
		const orgName = "ukama"
		const orgCert = "ukamacert"

		var orgId = uuid.NewV4()
		var orgOwner = uuid.NewV4()

		rows := sqlmock.NewRows([]string{"id", "name", "owner", "certificate"}).
			AddRow(orgId, orgName, orgOwner, orgCert)

		mock.ExpectQuery(`^SELECT.*orgs.*`).
			WithArgs(orgId, sqlmock.AnyArg()).
			WillReturnRows(rows)

		// Act
		org, err := r.Get(orgId)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, org)

		assert.Equal(t, orgId, org.Id)
		assert.Equal(t, orgName, org.Name)
		assert.Equal(t, orgOwner, org.Owner)
		assert.Equal(t, orgCert, org.Certificate)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("OrgNotFound", func(t *testing.T) {
		// Arrange

		var orgId = uuid.NewV4()

		mock.ExpectQuery(`^SELECT.*orgs.*`).
			WithArgs(orgId, sqlmock.AnyArg()).
			WillReturnError(sql.ErrNoRows)

		// Act
		org, err := r.Get(orgId)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, org)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestOrgRepo_GetByName(t *testing.T) {
	var db *sql.DB

	db, mock, err := sqlmock.New() // mock sql.DB
	assert.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})

	gdb, err := gorm.Open(dialector, &gorm.Config{})
	assert.NoError(t, err)

	r := org_db.NewOrgRepo(&UkamaDbMock{
		GormDb: gdb,
	})

	t.Run("OrgFound", func(t *testing.T) {
		// Arrange
		const orgName = "ukama"
		const orgCert = "ukamacert"

		var orgId = uuid.NewV4()
		var orgOwner = uuid.NewV4()

		rows := sqlmock.NewRows([]string{"id", "name", "owner", "certificate"}).
			AddRow(orgId, orgName, orgOwner, orgCert)

		mock.ExpectQuery(`^SELECT.*orgs.*`).
			WithArgs(orgName, sqlmock.AnyArg()).
			WillReturnRows(rows)

		// Act
		org, err := r.GetByName(orgName)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, org)

		assert.Equal(t, orgId, org.Id)
		assert.Equal(t, orgName, org.Name)
		assert.Equal(t, orgOwner, org.Owner)
		assert.Equal(t, orgCert, org.Certificate)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("OrgNotFound", func(t *testing.T) {
		// Arrange

		var orgName = "lol"

		mock.ExpectQuery(`^SELECT.*orgs.*`).
			WithArgs(orgName, sqlmock.AnyArg()).
			WillReturnError(sql.ErrNoRows)

		// Act
		org, err := r.GetByName(orgName)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, org)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestOrgRepo_GetByOwner(t *testing.T) {
	var db *sql.DB

	db, mock, err := sqlmock.New() // mock sql.DB
	assert.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})

	gdb, err := gorm.Open(dialector, &gorm.Config{})
	assert.NoError(t, err)

	r := org_db.NewOrgRepo(&UkamaDbMock{
		GormDb: gdb,
	})

	t.Run("OwnerFound", func(t *testing.T) {
		// Arrange
		const orgName = "ukama"
		const orgCert = "ukamacert"

		var orgId = uuid.NewV4()
		var orgOwner = uuid.NewV4()

		rows := sqlmock.NewRows([]string{"id", "name", "owner", "certificate"}).
			AddRow(orgId, orgName, orgOwner, orgCert)

		mock.ExpectQuery(`^SELECT.*orgs.*`).
			WithArgs(orgOwner).
			WillReturnRows(rows)

		// Act
		orgs, err := r.GetByOwner(orgOwner)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, orgs)

		assert.Equal(t, orgId, orgs[0].Id)
		assert.Equal(t, orgName, orgs[0].Name)
		assert.Equal(t, orgOwner, orgs[0].Owner)
		assert.Equal(t, orgCert, orgs[0].Certificate)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("OwnerNotFound", func(t *testing.T) {
		// Arrange

		var orgOwner = uuid.NewV4()

		mock.ExpectQuery(`^SELECT.*orgs.*`).
			WithArgs(orgOwner).
			WillReturnError(sql.ErrNoRows)

		// Act
		orgs, err := r.GetByOwner(orgOwner)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, orgs)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

// func TestOrgRepo_GetByMember(t *testing.T) {
// 	var db *sql.DB

// 	db, mock, err := sqlmock.New() // mock sql.DB
// 	assert.NoError(t, err)

// 	dialector := postgres.New(postgres.Config{
// 		DSN:                  "sqlmock_db_0",
// 		DriverName:           "postgres",
// 		Conn:                 db,
// 		PreferSimpleProtocol: true,
// 	})

// 	gdb, err := gorm.Open(dialector, &gorm.Config{})
// 	assert.NoError(t, err)

// 	r := org_db.NewOrgRepo(&UkamaDbMock{
// 		GormDb: gdb,
// 	})

// 	t.Run("MemberFound", func(t *testing.T) {
// 		// Arrange
// 		const userId = 1
// 		const deactivated = false

// 		var orgId = uuid.NewV4()
// 		var uuid = uuid.NewV4()

// 		rows := sqlmock.NewRows([]string{"id","org_id", "user_id", "uuid", "deactivated"}).
// 			AddRow(orgId, userId, uuid, deactivated)

// 		mock.ExpectQuery(`^SELECT.*org_users.*`).
// 			WithArgs(userId).
// 			WillReturnRows(rows)

// 		// Act
// 		members, err := r.GetByMember(userId)

// 		// Assert
// 		assert.NoError(t, err)
// 		assert.NotNil(t, members)

// 		assert.Equal(t, orgId, members[0].OrgId)
// 		assert.Equal(t, uuid, members[0].Uuid)
// 		assert.Equal(t, deactivated, members[0].Deactivated)

// 		err = mock.ExpectationsWereMet()
// 		assert.NoError(t, err)
// 	})

// 	t.Run("MemberNotFound", func(t *testing.T) {
// 		// Arrange
// 		var uuid = uuid.NewV4()

// 		mock.ExpectQuery(`^SELECT.*org_users.*`).
// 			WithArgs(uuid).
// 			WillReturnError(sql.ErrNoRows)

// 		// Act
// 		members, err := r.GetByMember(uuid)

// 		// Assert
// 		assert.Error(t, err)
// 		assert.Nil(t, members)

// 		err = mock.ExpectationsWereMet()
// 		assert.NoError(t, err)
// 	})
//}

func TestOrgRepo_GetAll(t *testing.T) {
	var db *sql.DB

	db, mock, err := sqlmock.New() // mock sql.DB
	assert.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})

	gdb, err := gorm.Open(dialector, &gorm.Config{})
	assert.NoError(t, err)

	r := org_db.NewOrgRepo(&UkamaDbMock{
		GormDb: gdb,
	})

	t.Run("GetAll", func(t *testing.T) {
		// Arrange
		const orgName = "ukama"
		const orgCert = "ukamacert"

		var orgId = uuid.NewV4()
		var orgOwner = uuid.NewV4()

		rows := sqlmock.NewRows([]string{"id", "name", "owner", "certificate"}).
			AddRow(orgId, orgName, orgOwner, orgCert)

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
}
