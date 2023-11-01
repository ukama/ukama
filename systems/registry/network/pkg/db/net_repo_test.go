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
	extsql "database/sql"
	"log"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/tj/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/common/uuid"
	net_db "github.com/ukama/ukama/systems/registry/network/pkg/db"
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

func Test_NetRepo_Get(t *testing.T) {
	t.Run("NetworkExist", func(t *testing.T) {
		// Arrange
		const netName = "network1"

		var db *extsql.DB
		var netID = uuid.NewV4()
		var orgID = uuid.NewV4()

		networks := pq.StringArray{"Verizon"}
		countries := pq.StringArray{"USA"}

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"id", "name", "org_id", "allowed_networks",
			"allowed_countries"}).
			AddRow(netID, netName, orgID, networks, countries)

		mock.ExpectQuery(`^SELECT.*networks.*`).
			WithArgs(netID).
			WillReturnRows(rows)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := net_db.NewNetRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		net, err := r.Get(netID)

		// Assert
		assert.NoError(t, err)

		assert.NotNil(t, net)
		assert.Equal(t, net.Id, netID)
		assert.Equal(t, net.OrgId, orgID)
		assert.Equal(t, net.AllowedNetworks, networks)
		assert.Equal(t, net.AllowedCountries, countries)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("NetworkNotFound", func(t *testing.T) {
		// Arrange
		var db *extsql.DB
		var netID = uuid.NewV4()

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		mock.ExpectQuery(`^SELECT.*networks.*`).
			WithArgs(netID).
			WillReturnError(sql.ErrNoRows)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := net_db.NewNetRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		net, err := r.Get(netID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, net)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func Test_NetRepo_GetByName(t *testing.T) {
	t.Run("NetworkExist", func(t *testing.T) {
		// Arrange
		const netName = "network1"
		var netID = uuid.NewV4()
		var orgID = uuid.NewV4()
		const orgName = "org1"

		var db *extsql.DB

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"id", "name", "org_id"}).
			AddRow(netID, netName, orgID)

		mock.ExpectQuery(`^SELECT.*networks.*`).
			WithArgs(orgName, netName).
			WillReturnRows(rows)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := net_db.NewNetRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		network, err := r.GetByName(orgName, netName)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.NotNil(t, network)
	})

	t.Run("NetworkNotFound", func(t *testing.T) {
		// Arrange
		const netName = "network1"
		const orgName = "org1"

		var db *extsql.DB

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		mock.ExpectQuery(`^SELECT.*network.*`).
			WithArgs(orgName, netName).
			WillReturnError(sql.ErrNoRows)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := net_db.NewNetRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		network, err := r.GetByName(orgName, netName)

		// Assert
		assert.Error(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.Nil(t, network)
	})
}

func Test_NetRepo_GetByOrgId(t *testing.T) {
	t.Run("OrgExist", func(t *testing.T) {
		// Arrange
		const netName = "network1"
		var netID = uuid.NewV4()
		var orgID = uuid.NewV4()

		var db *extsql.DB

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"id", "name", "network_id"}).
			AddRow(netID, netName, orgID)

		mock.ExpectQuery(`^SELECT.*networks.*`).
			WithArgs(orgID).
			WillReturnRows(rows)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := net_db.NewNetRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		networks, err := r.GetByOrg(orgID)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.NotNil(t, networks)
	})

	t.Run("NetworkNotFound", func(t *testing.T) {
		// Arrange
		var orgID = uuid.NewV4()

		var db *extsql.DB

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		mock.ExpectQuery(`^SELECT.*networks.*`).
			WithArgs(orgID).
			WillReturnError(sql.ErrNoRows)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := net_db.NewNetRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		networks, err := r.GetByOrg(orgID)

		// Assert
		assert.Error(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.Nil(t, networks)
	})
}

func Test_NetRepo_Add(t *testing.T) {
	t.Run("AddNetwork", func(t *testing.T) {
		// Arrange
		var db *extsql.DB

		network := net_db.Network{
			Id:    uuid.NewV4(),
			Name:  "network1",
			OrgId: uuid.NewV4(),
		}

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`INSERT`)).
			WithArgs(network.Id, network.Name, network.OrgId, sqlmock.AnyArg(), sqlmock.AnyArg(),
				sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
				sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := net_db.NewNetRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		err = r.Add(&network, nil)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func Test_NetRepo_Delete(t *testing.T) {
	t.Run("DeleteNetwork", func(t *testing.T) {
		var db *extsql.DB

		const orgName = "org1"
		const netName = "net1"
		var netID = uuid.NewV4()
		var orgID = uuid.NewV4()

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"id", "name", "network_id"}).
			AddRow(netID, netName, orgID)

		mock.ExpectBegin()

		mock.ExpectQuery(`^SELECT.*networks.*`).
			WithArgs(orgName, netName).
			WillReturnRows(rows)

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "networks" SET`)).
			WithArgs(sqlmock.AnyArg(), netID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := net_db.NewNetRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		err = r.Delete(orgName, netName)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}
