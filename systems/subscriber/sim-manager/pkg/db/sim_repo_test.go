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

	"github.com/ukama/ukama/systems/common/uuid"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/tj/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	log "github.com/sirupsen/logrus"
	simdb "github.com/ukama/ukama/systems/subscriber/sim-manager/pkg/db"
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

func TestSimRepo_Get(t *testing.T) {
	t.Run("SimFound", func(t *testing.T) {
		// Arrange
		var simID = uuid.NewV4()
		var netID = uuid.NewV4()
		var subID = uuid.NewV4()

		var packageID = uuid.NewV4()

		var db *sql.DB

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		simRow := sqlmock.NewRows([]string{"id", "network_id", "subscriber_id"}).
			AddRow(simID, netID, subID)

		packageRow := sqlmock.NewRows([]string{"id", "sim_id"}).
			AddRow(packageID, simID)

		mock.ExpectQuery(`^SELECT.*sims.*`).
			WithArgs(simID, sqlmock.AnyArg()).
			WillReturnRows(simRow)

		mock.ExpectQuery(`^SELECT.*packages.*`).
			WithArgs(simID).
			WillReturnRows(packageRow)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := simdb.NewSimRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		sim, err := r.Get(simID)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.NotNil(t, sim)
		assert.NotNil(t, sim)
		assert.Equal(t, packageID, sim.Package.Id)
	})

	t.Run("SimNotFound", func(t *testing.T) {
		// Arrange
		var simID = uuid.NewV4()

		var db *sql.DB

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		mock.ExpectQuery(`^SELECT.*sims.*`).
			WithArgs(simID, sqlmock.AnyArg()).
			WillReturnError(sql.ErrNoRows)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := simdb.NewSimRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		sim, err := r.Get(simID)

		// Assert
		assert.Error(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.Nil(t, sim)
	})
}

func TestSimRepo_GetByIccid(t *testing.T) {
	const testIccid = "890000-this-is-a-test-iccid"

	t.Run("IccidFound", func(t *testing.T) {
		// Arrange
		var simID = uuid.NewV4()
		var netID = uuid.NewV4()
		var subID = uuid.NewV4()

		var packageID = uuid.NewV4()

		var db *sql.DB

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		simRow := sqlmock.NewRows([]string{"id", "network_id", "subscriber_id", "iccid"}).
			AddRow(simID, netID, subID, testIccid)

		packageRow := sqlmock.NewRows([]string{"id", "sim_id"}).
			AddRow(packageID, simID)

		mock.ExpectQuery(`^SELECT.*sims.*`).
			WithArgs(testIccid, 1).
			WillReturnRows(simRow)

		mock.ExpectQuery(`^SELECT.*packages.*`).
			WithArgs(simID).
			WillReturnRows(packageRow)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := simdb.NewSimRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		sim, err := r.GetByIccid(testIccid)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()

		assert.NoError(t, err)
		assert.NotNil(t, sim)

		assert.Equal(t, testIccid, sim.Iccid)
	})

	t.Run("IccidNotFound", func(t *testing.T) {
		var db *sql.DB

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		mock.ExpectQuery(`^SELECT.*sims.*`).
			WithArgs(testIccid, 1).
			WillReturnError(sql.ErrNoRows)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := simdb.NewSimRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		sim, err := r.GetByIccid(testIccid)

		// Assert
		assert.Error(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.Nil(t, sim)
	})
}

func TestSimRepo_GetBySubscriber(t *testing.T) {
	t.Run("SubscriberFound", func(t *testing.T) {
		// Arrange
		var simID = uuid.NewV4()
		var netID = uuid.NewV4()
		var subID = uuid.NewV4()

		var packageID = uuid.NewV4()

		var db *sql.DB

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		simRow := sqlmock.NewRows([]string{"id", "network_id", "subscriber_id"}).
			AddRow(simID, netID, subID)

		packageRow := sqlmock.NewRows([]string{"id", "sim_id"}).
			AddRow(packageID, simID)

		mock.ExpectQuery(`^SELECT.*sims.*`).
			WithArgs(subID).
			WillReturnRows(simRow)

		mock.ExpectQuery(`^SELECT.*packages.*`).
			WithArgs(simID).
			WillReturnRows(packageRow)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := simdb.NewSimRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		sims, err := r.GetBySubscriber(subID)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.NotNil(t, sims)
		assert.Equal(t, packageID, sims[0].Package.Id)
	})

	t.Run("SubscriberNotFound", func(t *testing.T) {
		// Arrange
		var subID = uuid.NewV4()

		var db *sql.DB

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		mock.ExpectQuery(`^SELECT.*sims.*`).
			WithArgs(subID).
			WillReturnError(sql.ErrNoRows)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := simdb.NewSimRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		sims, err := r.GetBySubscriber(subID)

		// Assert
		assert.Error(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.Nil(t, sims)
	})
}

func TestSimRepo_GetByNetwork(t *testing.T) {
	t.Run("NetworkFound", func(t *testing.T) {
		// Arrange
		var simID = uuid.NewV4()
		var netID = uuid.NewV4()
		var subID = uuid.NewV4()

		var packageID = uuid.NewV4()

		var db *sql.DB

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		simRow := sqlmock.NewRows([]string{"id", "network_id", "subscriber_id"}).
			AddRow(simID, netID, subID)

		packageRow := sqlmock.NewRows([]string{"id", "sim_id"}).
			AddRow(packageID, simID)

		mock.ExpectQuery(`^SELECT.*sims.*`).
			WithArgs(netID).
			WillReturnRows(simRow)

		mock.ExpectQuery(`^SELECT.*packages.*`).
			WithArgs(simID).
			WillReturnRows(packageRow)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := simdb.NewSimRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		sims, err := r.GetByNetwork(netID)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.NotNil(t, sims)
		assert.Equal(t, packageID, sims[0].Package.Id)
	})

	t.Run("NetworkNotFound", func(t *testing.T) {
		// Arrange
		var netID = uuid.NewV4()

		var db *sql.DB

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		mock.ExpectQuery(`^SELECT.*sims.*`).
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

		r := simdb.NewSimRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		sims, err := r.GetByNetwork(netID)

		// Assert
		assert.Error(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.Nil(t, sims)
	})
}

func TestSimRepo_Add(t *testing.T) {
	t.Run("AddSim", func(t *testing.T) {
		// Arrange
		var db *sql.DB

		sim := simdb.Sim{
			Id:           uuid.NewV4(),
			SubscriberId: uuid.NewV4(),
		}

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`INSERT`)).
			WithArgs(sim.Id, sim.SubscriberId, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
				sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
				sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
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

		r := simdb.NewSimRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		err = r.Add(&sim, nil)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestSimRepo_Delete(t *testing.T) {
	t.Run("SimFound", func(t *testing.T) {
		var db *sql.DB

		// Arrange
		var simID = uuid.NewV4()

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "sims" SET`)).
			WithArgs(sqlmock.AnyArg(), simID).
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

		r := simdb.NewSimRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		err = r.Delete(simID, nil)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}
