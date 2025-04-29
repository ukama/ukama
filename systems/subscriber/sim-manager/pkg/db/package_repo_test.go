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

	simdb "github.com/ukama/ukama/systems/subscriber/sim-manager/pkg/db"
)

func TestPackageRepo_Add(t *testing.T) {
	t.Run("AddPackage", func(t *testing.T) {
		// Arrange
		var db *sql.DB

		pkg := simdb.Package{
			Id:    uuid.NewV4(),
			SimId: uuid.NewV4(),
		}

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`INSERT`)).
			WithArgs(pkg.Id, pkg.SimId, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
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

		r := simdb.NewPackageRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		err = r.Add(&pkg, nil)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestPackageRepo_Get(t *testing.T) {
	t.Run("PackageFound", func(t *testing.T) {
		// Arrange
		var packageID = uuid.NewV4()
		var simID = uuid.NewV4()

		var db *sql.DB

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		row := sqlmock.NewRows([]string{"id", "sim_id"}).
			AddRow(packageID, simID)

		mock.ExpectQuery(`^SELECT.*packages.*`).
			WithArgs(packageID, sqlmock.AnyArg()).
			WillReturnRows(row)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := simdb.NewPackageRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		pkg, err := r.Get(packageID)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.NotNil(t, pkg)
		assert.NotNil(t, pkg)
	})

	t.Run("PackageNotFound", func(t *testing.T) {
		// Arrange
		var packageID = uuid.NewV4()

		var db *sql.DB

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		mock.ExpectQuery(`^SELECT.*packages.*`).
			WithArgs(packageID, sqlmock.AnyArg()).
			WillReturnError(sql.ErrNoRows)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := simdb.NewPackageRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		pkg, err := r.Get(packageID)

		// Assert
		assert.Error(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.Nil(t, pkg)
	})
}

func TestPackageRepo_GetBySim(t *testing.T) {
	t.Run("SimFound", func(t *testing.T) {
		// Arrange
		var simID = uuid.NewV4()
		var packageID = uuid.NewV4()

		var db *sql.DB

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		packageRow := sqlmock.NewRows([]string{"id", "sim_id"}).
			AddRow(packageID, simID)

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

		r := simdb.NewPackageRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		packages, err := r.GetBySim(simID)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.NotNil(t, packages)
	})

	t.Run("SimNotFound", func(t *testing.T) {
		// Arrange
		var simID = uuid.NewV4()

		var db *sql.DB

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		mock.ExpectQuery(`^SELECT.*packages.*`).
			WithArgs(simID).
			WillReturnError(sql.ErrNoRows)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := simdb.NewPackageRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		packages, err := r.GetBySim(simID)

		// Assert
		assert.Error(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.Nil(t, packages)
	})
}

func TestPackageRepo_Delete(t *testing.T) {
	t.Run("PackageFound", func(t *testing.T) {
		var db *sql.DB

		// Arrange
		var packageID = uuid.NewV4()

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`DELETE`)).
			WithArgs(packageID).
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

		r := simdb.NewPackageRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		err = r.Delete(packageID, nil)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}
