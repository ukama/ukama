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
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/tj/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/subscriber/sim-manager/pkg/db"

	log "github.com/sirupsen/logrus"
)

var (
	validNestedSimFunc = func(sim *db.Sim, tx *gorm.DB) error {
		return nil
	}
	unvalidNestedSimFunc = func(sim *db.Sim, tx *gorm.DB) error {
		return errors.New("some errors occurred")
	}
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

func TestSimRepo_Add(t *testing.T) {
	t.Run("AddSim", func(t *testing.T) {
		sim := db.Sim{
			Id:           uuid.NewV4(),
			SubscriberId: uuid.NewV4(),
		}

		mock, gdb := prepareDb(t)

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`INSERT`)).
			WithArgs(sim.Id, sim.SubscriberId, sqlmock.AnyArg(), sqlmock.AnyArg(),
				sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
				sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
				sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
				sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		r := db.NewSimRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		err := r.Add(&sim, nil)

		// Assert
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("AddSimError", func(t *testing.T) {
		sim := db.Sim{
			Id:           uuid.NewV4(),
			SubscriberId: uuid.NewV4(),
		}

		mock, gdb := prepareDb(t)

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`INSERT`)).
			WithArgs(sim.Id, sim.SubscriberId, sqlmock.AnyArg(), sqlmock.AnyArg(),
				sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
				sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
				sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
				sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnError(sql.ErrNoRows)

		r := db.NewSimRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		err := r.Add(&sim, validNestedSimFunc)

		// Assert
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("AddSimNestedFuncError", func(t *testing.T) {
		sim := db.Sim{
			Id:           uuid.NewV4(),
			SubscriberId: uuid.NewV4(),
		}

		mock, gdb := prepareDb(t)

		mock.ExpectBegin()

		r := db.NewSimRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		err := r.Add(&sim, unvalidNestedSimFunc)

		// Assert
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSimRepo_Get(t *testing.T) {
	t.Run("SimFound", func(t *testing.T) {
		var (
			simID     = uuid.NewV4()
			netID     = uuid.NewV4()
			subID     = uuid.NewV4()
			packageID = uuid.NewV4()
		)

		mock, gdb := prepareDb(t)
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

		r := db.NewSimRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		sim, err := r.Get(simID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, sim)
		assert.Equal(t, packageID, sim.Package.Id)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("SimNotFound", func(t *testing.T) {
		var simID = uuid.NewV4()

		mock, gdb := prepareDb(t)

		mock.ExpectQuery(`^SELECT.*sims.*`).
			WithArgs(simID, sqlmock.AnyArg()).
			WillReturnError(sql.ErrNoRows)

		r := db.NewSimRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		sim, err := r.Get(simID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, sim)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSimRepo_List(t *testing.T) {
	const (
		testIccid = "890000-this-is-a-test-iccid"
		testImsi  = "890000-this-is-a-test-imsi"
	)

	t.Run("ListAll", func(t *testing.T) {
		var (
			simID     = uuid.NewV4()
			netID     = uuid.NewV4()
			subID     = uuid.NewV4()
			packageID = uuid.NewV4()
		)

		mock, gdb := prepareDb(t)
		simRow := sqlmock.NewRows([]string{"id", "network_id", "subscriber_id", "iccid"}).
			AddRow(simID, netID, subID, testIccid)

		packageRow := sqlmock.NewRows([]string{"id", "sim_id"}).
			AddRow(packageID, simID)

		mock.ExpectQuery(`^SELECT.*sims.*`).
			WithArgs().
			WillReturnRows(simRow)

		mock.ExpectQuery(`^SELECT.*packages.*`).
			WithArgs(simID).
			WillReturnRows(packageRow)

		r := db.NewSimRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		list, err := r.List("", "", "", "", ukama.SimTypeUnknown, ukama.SimStatusUnknown,
			0, false, 0, false)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, list)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("SimFound", func(t *testing.T) {
		var (
			simID     = uuid.NewV4()
			netID     = uuid.NewV4()
			subID     = uuid.NewV4()
			packageID = uuid.NewV4()
		)

		mock, gdb := prepareDb(t)
		simRow := sqlmock.NewRows([]string{"id", "network_id", "subscriber_id", "iccid"}).
			AddRow(simID, netID, subID, testIccid)

		packageRow := sqlmock.NewRows([]string{"id", "sim_id"}).
			AddRow(packageID, simID)

		mock.ExpectQuery(`^SELECT.*sims.*`).
			WithArgs().
			WillReturnRows(simRow)

		mock.ExpectQuery(`^SELECT.*packages.*`).
			WithArgs(simID).
			WillReturnRows(packageRow)

		r := db.NewSimRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		list, err := r.List(testIccid, testImsi, subID.String(), netID.String(),
			ukama.SimTypeUkamaData, ukama.SimStatusActive, 22, true, 1, true)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, list)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("SimNotFound", func(t *testing.T) {
		var (
			netID = uuid.NewV4()
			subID = uuid.NewV4()
		)

		mock, gdb := prepareDb(t)

		mock.ExpectQuery(`^SELECT.*sims.*`).
			WithArgs().
			WillReturnError(sql.ErrNoRows)

		r := db.NewSimRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		list, err := r.List(testIccid, testImsi, subID.String(), netID.String(),
			ukama.SimTypeUkamaData, ukama.SimStatusActive, 22, true, 1, true)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, list)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSimRepo_GetByIccid(t *testing.T) {
	const testIccid = "890000-this-is-a-test-iccid"

	t.Run("IccidFound", func(t *testing.T) {
		var (
			simID     = uuid.NewV4()
			netID     = uuid.NewV4()
			subID     = uuid.NewV4()
			packageID = uuid.NewV4()
		)

		mock, gdb := prepareDb(t)
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

		r := db.NewSimRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		sim, err := r.GetByIccid(testIccid)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, sim)
		assert.Equal(t, testIccid, sim.Iccid)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("IccidNotFound", func(t *testing.T) {
		mock, gdb := prepareDb(t)

		mock.ExpectQuery(`^SELECT.*sims.*`).
			WithArgs(testIccid, 1).
			WillReturnError(sql.ErrNoRows)

		r := db.NewSimRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		sim, err := r.GetByIccid(testIccid)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, sim)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSimRepo_GetBySubscriber(t *testing.T) {
	t.Run("SubscriberFound", func(t *testing.T) {
		var (
			simID     = uuid.NewV4()
			netID     = uuid.NewV4()
			subID     = uuid.NewV4()
			packageID = uuid.NewV4()
		)

		mock, gdb := prepareDb(t)
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

		r := db.NewSimRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		sims, err := r.GetBySubscriber(subID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, sims)
		assert.Equal(t, packageID, sims[0].Package.Id)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("SubscriberNotFound", func(t *testing.T) {
		var subID = uuid.NewV4()

		mock, gdb := prepareDb(t)

		mock.ExpectQuery(`^SELECT.*sims.*`).
			WithArgs(subID).
			WillReturnError(sql.ErrNoRows)

		r := db.NewSimRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		sims, err := r.GetBySubscriber(subID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, sims)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSimRepo_GetByNetwork(t *testing.T) {
	t.Run("NetworkFound", func(t *testing.T) {
		var (
			simID     = uuid.NewV4()
			netID     = uuid.NewV4()
			subID     = uuid.NewV4()
			packageID = uuid.NewV4()
		)

		mock, gdb := prepareDb(t)
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

		r := db.NewSimRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		sims, err := r.GetByNetwork(netID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, sims)
		assert.Equal(t, packageID, sims[0].Package.Id)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("NetworkNotFound", func(t *testing.T) {
		var netID = uuid.NewV4()

		mock, gdb := prepareDb(t)

		mock.ExpectQuery(`^SELECT.*sims.*`).
			WithArgs(netID).
			WillReturnError(sql.ErrNoRows)

		r := db.NewSimRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		sims, err := r.GetByNetwork(netID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, sims)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSimRepo_Update(t *testing.T) {
	t.Run("SimFound", func(t *testing.T) {
		var (
			simID = uuid.NewV4()
			netID = uuid.NewV4()
			subID = uuid.NewV4()
		)

		sim := db.Sim{
			Id:           simID,
			SubscriberId: subID,
			NetworkId:    netID,
		}

		mock, gdb := prepareDb(t)
		simRow := sqlmock.NewRows([]string{"id", "network_id", "subscriber_id"}).
			AddRow(simID, netID, subID)

		mock.ExpectBegin()

		mock.ExpectQuery(`^UPDATE.*sims.*`).
			WithArgs(sim.SubscriberId, sim.NetworkId, sqlmock.AnyArg(), sim.Id).
			WillReturnRows(simRow)

		mock.ExpectCommit()

		r := db.NewSimRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		err := r.Update(&sim, nil)

		// Assert
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("SimNotFound", func(t *testing.T) {
		var (
			simID = uuid.NewV4()
			netID = uuid.NewV4()
			subID = uuid.NewV4()
		)

		sim := db.Sim{
			Id:           simID,
			SubscriberId: subID,
			NetworkId:    netID,
		}

		mock, gdb := prepareDb(t)
		simRow := sqlmock.NewRows([]string{"id", "network_id", "subscriber_id"})
		mock.ExpectBegin()

		mock.ExpectQuery(`^UPDATE.*sims.*`).
			WithArgs(sim.SubscriberId, sim.NetworkId, sqlmock.AnyArg(), sim.Id).
			WillReturnRows(simRow)

		r := db.NewSimRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		err := r.Update(&sim, validNestedSimFunc)

		// Assert
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("SimUpdateError", func(t *testing.T) {
		var (
			simID = uuid.NewV4()
			netID = uuid.NewV4()
			subID = uuid.NewV4()
		)

		sim := db.Sim{
			Id:           simID,
			SubscriberId: subID,
			NetworkId:    netID,
		}

		mock, gdb := prepareDb(t)

		mock.ExpectBegin()

		mock.ExpectQuery(`^UPDATE.*sims.*`).
			WithArgs(sim.SubscriberId, sim.NetworkId, sqlmock.AnyArg(), sim.Id).
			WillReturnError(sql.ErrNoRows)

		r := db.NewSimRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		err := r.Update(&sim, nil)

		// Assert
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("SimUpdateNestedFuncError", func(t *testing.T) {
		var (
			simID = uuid.NewV4()
			netID = uuid.NewV4()
			subID = uuid.NewV4()
		)

		sim := db.Sim{
			Id:           simID,
			SubscriberId: subID,
			NetworkId:    netID,
		}

		mock, gdb := prepareDb(t)

		r := db.NewSimRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		err := r.Update(&sim, unvalidNestedSimFunc)

		// Assert
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSimRepo_Delete(t *testing.T) {
	t.Run("SimFound", func(t *testing.T) {
		var simID = uuid.NewV4()

		mock, gdb := prepareDb(t)

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "sims" SET`)).
			WithArgs(sqlmock.AnyArg(), simID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		r := db.NewSimRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		err := r.Delete(simID, nil)

		// Assert
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("SimDeleteError", func(t *testing.T) {
		var simID = uuid.NewV4()

		mock, gdb := prepareDb(t)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "sims" SET`)).
			WithArgs(sqlmock.AnyArg(), simID).
			WillReturnError(sql.ErrNoRows)

		r := db.NewSimRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		err := r.Delete(simID,
			func(uuid.UUID, *gorm.DB) error { return nil })

		// Assert
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("SimDeleteNetedFuncError", func(t *testing.T) {
		var simID = uuid.NewV4()

		mock, gdb := prepareDb(t)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "sims" SET`)).
			WithArgs(sqlmock.AnyArg(), simID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		r := db.NewSimRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		err := r.Delete(simID,
			func(uuid.UUID, *gorm.DB) error {
				return errors.
					New("some error occurred")
			})

		// Assert
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func prepareDb(t *testing.T) (sqlmock.Sqlmock, *gorm.DB) {
	var db *sql.DB
	var err error

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

	return mock, gdb
}
