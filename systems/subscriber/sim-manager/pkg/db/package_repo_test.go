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
	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/subscriber/sim-manager/pkg/db"
)

var (
	validNestedPackageFunc = func(pckg *db.Package, tx *gorm.DB) error {
		return nil
	}
	unvalidNestedPackageFunc = func(pckg *db.Package, tx *gorm.DB) error {
		return errors.New("some errors occurred")
	}
)

func TestPackageRepo_Add(t *testing.T) {
	t.Run("AddPackage", func(t *testing.T) {
		pkg := db.Package{
			Id:    uuid.NewV4(),
			SimId: uuid.NewV4(),
		}

		mock, gdb := prepareDb(t)

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`INSERT`)).
			WithArgs(pkg.Id, pkg.SimId, sqlmock.AnyArg(), sqlmock.AnyArg(),
				sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
				sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		r := db.NewPackageRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		err := r.Add(&pkg, nil)

		// Assert
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("AddPackageError", func(t *testing.T) {
		pkg := db.Package{
			Id:    uuid.NewV4(),
			SimId: uuid.NewV4(),
		}

		mock, gdb := prepareDb(t)

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`INSERT`)).
			WithArgs(pkg.Id, pkg.SimId, sqlmock.AnyArg(), sqlmock.AnyArg(),
				sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
				sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnError(sql.ErrNoRows)

		r := db.NewPackageRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		err := r.Add(&pkg, validNestedPackageFunc)

		// Assert
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("AddPackageNestedFuncError", func(t *testing.T) {
		pkg := db.Package{
			Id:    uuid.NewV4(),
			SimId: uuid.NewV4(),
		}

		mock, gdb := prepareDb(t)

		r := db.NewPackageRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		err := r.Add(&pkg, unvalidNestedPackageFunc)

		// Assert
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestPackageRepo_Get(t *testing.T) {
	t.Run("PackageFound", func(t *testing.T) {
		var (
			simID     = uuid.NewV4()
			packageID = uuid.NewV4()
		)

		mock, gdb := prepareDb(t)

		row := sqlmock.NewRows([]string{"id", "sim_id"}).
			AddRow(packageID, simID)

		mock.ExpectQuery(`^SELECT.*packages.*`).
			WithArgs(packageID, sqlmock.AnyArg()).
			WillReturnRows(row)

		r := db.NewPackageRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		pkg, err := r.Get(packageID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, pkg)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("PackageNotFound", func(t *testing.T) {
		var packageID = uuid.NewV4()

		mock, gdb := prepareDb(t)

		mock.ExpectQuery(`^SELECT.*packages.*`).
			WithArgs(packageID, sqlmock.AnyArg()).
			WillReturnError(sql.ErrNoRows)

		r := db.NewPackageRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		pkg, err := r.Get(packageID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, pkg)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestPackageRepo_List(t *testing.T) {
	const (
		from = "2022-12-01T00:00:00Z"
		to   = "2023-12-01T00:00:00Z"
	)

	t.Run("ListAll", func(t *testing.T) {
		var (
			packageID = uuid.NewV4()
			simID     = uuid.NewV4()
		)

		mock, gdb := prepareDb(t)
		packageRow := sqlmock.NewRows([]string{"id", "sim_id"}).
			AddRow(packageID, simID)

		mock.ExpectQuery(`^SELECT.*packages.*`).
			WithArgs().
			WillReturnRows(packageRow)

		r := db.NewPackageRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		list, err := r.List("", "", "", "", "", "", false, false, 0, false)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, list)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("PackageFound", func(t *testing.T) {
		var (
			packageID  = uuid.NewV4()
			simID      = uuid.NewV4()
			dataplanID = uuid.NewV4()
		)

		mock, gdb := prepareDb(t)
		packageRow := sqlmock.NewRows([]string{"id", "sim_id"}).
			AddRow(packageID, simID)

		mock.ExpectQuery(`^SELECT.*packages.*`).
			WithArgs().
			WillReturnRows(packageRow)

		r := db.NewPackageRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		list, err := r.List(simID.String(), dataplanID.String(), from, to, from, to,
			true, true, 1, true)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, list)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("PackageNotFound", func(t *testing.T) {
		var (
			simID      = uuid.NewV4()
			dataplanID = uuid.NewV4()
		)

		mock, gdb := prepareDb(t)

		mock.ExpectQuery(`^SELECT.*packages.*`).
			WithArgs().
			WillReturnError(sql.ErrNoRows)

		r := db.NewPackageRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		list, err := r.List(simID.String(), dataplanID.String(), from, to, from, to,
			true, true, 1, true)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, list)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestPackageRepo_GetBySim(t *testing.T) {
	t.Run("SimFound", func(t *testing.T) {
		var (
			simID     = uuid.NewV4()
			packageID = uuid.NewV4()
		)

		mock, gdb := prepareDb(t)
		packageRow := sqlmock.NewRows([]string{"id", "sim_id"}).
			AddRow(packageID, simID)

		mock.ExpectQuery(`^SELECT.*packages.*`).
			WithArgs(simID).
			WillReturnRows(packageRow)

		r := db.NewPackageRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		packages, err := r.GetBySim(simID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, packages)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("SimNotFound", func(t *testing.T) {
		var simID = uuid.NewV4()

		mock, gdb := prepareDb(t)

		mock.ExpectQuery(`^SELECT.*packages.*`).
			WithArgs(simID).
			WillReturnError(sql.ErrNoRows)

		r := db.NewPackageRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		packages, err := r.GetBySim(simID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, packages)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestPackageRepo_Update(t *testing.T) {
	t.Run("PackageFound", func(t *testing.T) {
		var (
			packageID = uuid.NewV4()
			simID     = uuid.NewV4()
		)

		pckg := db.Package{
			Id:    packageID,
			SimId: simID,
		}

		mock, gdb := prepareDb(t)

		packageRow := sqlmock.NewRows([]string{"id", "sim_id"}).
			AddRow(packageID, simID)

		mock.ExpectBegin()

		mock.ExpectQuery(`^UPDATE.*packages.*`).
			WithArgs(pckg.SimId, sqlmock.AnyArg(), pckg.Id).
			WillReturnRows(packageRow)

		mock.ExpectCommit()

		r := db.NewPackageRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		err := r.Update(&pckg, nil)

		// Assert
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("PackageNotFound", func(t *testing.T) {
		var (
			packageID = uuid.NewV4()
			simID     = uuid.NewV4()
		)

		pckg := db.Package{
			Id:    packageID,
			SimId: simID,
		}

		mock, gdb := prepareDb(t)
		packageRow := sqlmock.NewRows([]string{"id", "sim_id"})
		mock.ExpectBegin()

		mock.ExpectQuery(`^UPDATE.*packages.*`).
			WithArgs(pckg.SimId, sqlmock.AnyArg(), pckg.Id).
			WillReturnRows(packageRow)

		r := db.NewPackageRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		err := r.Update(&pckg, validNestedPackageFunc)

		// Assert
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("PackageUpdateError", func(t *testing.T) {
		var (
			packageID = uuid.NewV4()
			simID     = uuid.NewV4()
		)

		pckg := db.Package{
			Id:    packageID,
			SimId: simID,
		}

		mock, gdb := prepareDb(t)

		mock.ExpectBegin()

		mock.ExpectQuery(`^UPDATE.*packages.*`).
			WithArgs(pckg.SimId, sqlmock.AnyArg(), pckg.Id).
			WillReturnError(sql.ErrNoRows)

		r := db.NewPackageRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		err := r.Update(&pckg, nil)

		// Assert
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("PackageUpdateNestedFuncError", func(t *testing.T) {
		var (
			packageID = uuid.NewV4()
			simID     = uuid.NewV4()
		)

		pckg := db.Package{
			Id:    packageID,
			SimId: simID,
		}

		mock, gdb := prepareDb(t)

		r := db.NewPackageRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		err := r.Update(&pckg, unvalidNestedPackageFunc)

		// Assert
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestPackageRepo_Delete(t *testing.T) {
	t.Run("PackageFound", func(t *testing.T) {
		var packageID = uuid.NewV4()

		mock, gdb := prepareDb(t)

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`DELETE`)).
			WithArgs(packageID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		r := db.NewPackageRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		err := r.Delete(packageID, nil)

		// Assert
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("PackageDeleteError", func(t *testing.T) {
		var packageID = uuid.NewV4()

		mock, gdb := prepareDb(t)

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`DELETE`)).
			WithArgs(packageID).
			WillReturnError(sql.ErrNoRows)

		r := db.NewPackageRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		err := r.Delete(packageID,
			func(uuid.UUID, *gorm.DB) error { return nil })

		// Assert
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("PackageDeleteNetedFuncError", func(t *testing.T) {
		var packageID = uuid.NewV4()

		mock, gdb := prepareDb(t)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE`)).
			WithArgs(packageID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		r := db.NewPackageRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		err := r.Delete(packageID,
			func(uuid.UUID, *gorm.DB) error {
				return errors.
					New("some error occurred")
			})

		// Assert
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
