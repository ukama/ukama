/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package db_test

import (
	extsql "database/sql"
	"fmt"
	"log"
	"regexp"
	"testing"
	"time"

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

		networks := pq.StringArray{"Verizon"}
		countries := pq.StringArray{"USA"}

		db, mock, err := sqlmock.New() // mock extsql.DB
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"id", "name", "allowed_networks",
			"allowed_countries"}).
			AddRow(netID, netName, networks, countries)

		mock.ExpectQuery(`^SELECT.*networks.*`).
			WithArgs(netID, sqlmock.AnyArg()).
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
		assert.Equal(t, net.AllowedNetworks, networks)
		assert.Equal(t, net.AllowedCountries, countries)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("NetworkNotFound", func(t *testing.T) {
		// Arrange
		var db *extsql.DB
		var netID = uuid.NewV4()

		db, mock, err := sqlmock.New() // mock extsql.DB
		assert.NoError(t, err)

		mock.ExpectQuery(`^SELECT.*networks.*`).
			WithArgs(netID, sqlmock.AnyArg()).
			WillReturnError(extsql.ErrNoRows)

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

		var db *extsql.DB

		db, mock, err := sqlmock.New() // mock extsql.DB
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"id", "name"}).
			AddRow(netID, netName)

		mock.ExpectQuery(`^SELECT.*networks.*`).
			WithArgs(netName, sqlmock.AnyArg()).
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
		network, err := r.GetByName(netName)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.NotNil(t, network)
	})

	t.Run("NetworkNotFound", func(t *testing.T) {
		// Arrange
		const netName = "network1"

		var db *extsql.DB

		db, mock, err := sqlmock.New() // mock extsql.DB
		assert.NoError(t, err)

		mock.ExpectQuery(`^SELECT.*network.*`).
			WithArgs(netName, sqlmock.AnyArg()).
			WillReturnError(extsql.ErrNoRows)

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
		network, err := r.GetByName(netName)

		// Assert
		assert.Error(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.Nil(t, network)
	})
}

func Test_NetRepo_GetAll(t *testing.T) {
	t.Run("networks exist", func(t *testing.T) {
		// Arrange
		const netName = "network1"
		var netID = uuid.NewV4()

		var db *extsql.DB

		db, mock, err := sqlmock.New() // mock extsql.DB
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"id", "name"}).
			AddRow(netID, netName)

		mock.ExpectQuery(`^SELECT.*networks.*`).
			WithArgs().
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
		networks, err := r.GetAll()

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.NotNil(t, networks)
	})

	t.Run("NetworkNotFound", func(t *testing.T) {
		// Arrange

		var db *extsql.DB

		db, mock, err := sqlmock.New() // mock extsql.DB
		assert.NoError(t, err)

		mock.ExpectQuery(`^SELECT.*networks.*`).
			WithArgs().
			WillReturnError(extsql.ErrNoRows)

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
		networks, err := r.GetAll()

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
			Id:   uuid.NewV4(),
			Name: "network1",
		}

		db, mock, err := sqlmock.New() // mock extsql.DB
		assert.NoError(t, err)

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`INSERT`)).
			WithArgs(network.Id, network.Name, sqlmock.AnyArg(), sqlmock.AnyArg(),
				sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
				sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
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
		fmt.Println("NET", network)
		// Act
		err = r.Add(&network, nil)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func Test_NetRepo_Delete(t *testing.T) {
	t.Run("NetworkExist", func(t *testing.T) {
		// Arrange
		net := net_db.Network{
			Id:        uuid.NewV4(),
			Name:      "test",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: gorm.DeletedAt{},
		}
		var db *extsql.DB

		db, mock, err := sqlmock.New() // mock extsql.DB
		assert.NoError(t, err)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "networks"`)).
			WithArgs(sqlmock.AnyArg(), net.Id).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		dialector := postgres.New(postgres.Config{
			DSN:        "sqlmock_db_0",
			DriverName: "postgres",
			Conn:       db,

			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})

		assert.NoError(t, err)

		r := net_db.NewNetRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		err = r.Delete(net.Id)

		// Assert

		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}
