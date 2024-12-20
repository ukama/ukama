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

	"github.com/ukama/ukama/systems/billing/report/pkg/db"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"

	log "github.com/sirupsen/logrus"
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

func TestReportRepo_Add(t *testing.T) {
	t.Run("Addreport", func(t *testing.T) {
		// Arrange
		report := db.Report{
			Id:      uuid.NewV4(),
			OwnerId: uuid.NewV4(),
		}

		mock, gdb := prepare_db(t)

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`INSERT`)).
			WithArgs(report.Id, report.OwnerId, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
				sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		r := db.NewReportRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		err := r.Add(&report, nil)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("AddFaillure", func(t *testing.T) {
		// Arrange
		report := db.Report{
			Id:      uuid.NewV4(),
			OwnerId: uuid.NewV4(),
		}

		mock, gdb := prepare_db(t)

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`INSERT`)).
			WithArgs(report.Id, report.OwnerId, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
				sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnError(sql.ErrNoRows)

		// mock.ExpectCommit()

		r := db.NewReportRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		err := r.Add(&report, nil)

		// Assert
		assert.Error(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestReportRepo_Get(t *testing.T) {
	t.Run("reportFound", func(t *testing.T) {
		// Arrange
		var reportId = uuid.NewV4()
		var ownerId = uuid.NewV4()

		mock, gdb := prepare_db(t)

		rows := sqlmock.NewRows([]string{"id", "owner_id"}).
			AddRow(reportId, ownerId)

		mock.ExpectQuery(`^SELECT.*reports.*`).
			WithArgs(reportId, sqlmock.AnyArg()).
			WillReturnRows(rows)

		r := db.NewReportRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		rep, err := r.Get(reportId)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.NotNil(t, rep)
	})

	t.Run("reportNotFound", func(t *testing.T) {
		// Arrange
		var reportId = uuid.NewV4()

		mock, gdb := prepare_db(t)

		mock.ExpectQuery(`^SELECT.*reports.*`).
			WithArgs(reportId, sqlmock.AnyArg()).
			WillReturnError(sql.ErrNoRows)

		r := db.NewReportRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		rep, err := r.Get(reportId)

		// Assert
		assert.Error(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.Nil(t, rep)
	})
}

func TestReportRepo_List(t *testing.T) {
	t.Run("ListAll", func(t *testing.T) {
		report := &db.Report{
			Id:        uuid.NewV4(),
			OwnerId:   uuid.NewV4(),
			OwnerType: ukama.OwnerTypeOrg,
			IsPaid:    false,
		}

		mock, gdb := prepare_db(t)

		rows := sqlmock.NewRows([]string{"report_id", "owner_id", "owner_type", "is_paid"}).
			AddRow(report.Id, report.OwnerId, report.OwnerType, report.IsPaid)

		mock.ExpectQuery(`^SELECT.*reports.*`).
			WithArgs().
			WillReturnRows(rows)

		r := db.NewReportRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		reps, err := r.List("", ukama.OwnerTypeUnknown, "",
			ukama.ReportTypeUnknown, false, 0, false)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, reps)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ListByownerId", func(t *testing.T) {
		report := &db.Report{
			Id:        uuid.NewV4(),
			OwnerId:   uuid.NewV4(),
			OwnerType: ukama.OwnerTypeOrg,
			IsPaid:    false,
		}

		mock, gdb := prepare_db(t)

		rows := sqlmock.NewRows([]string{"report_id", "owner_id", "owner_type", "is_paid"}).
			AddRow(report.Id, report.OwnerId, report.OwnerType, report.IsPaid)

		mock.ExpectQuery(`^SELECT.*reports.*`).
			WithArgs(report.OwnerId).
			WillReturnRows(rows)

		r := db.NewReportRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		reps, err := r.List(report.OwnerId.String(), ukama.OwnerTypeUnknown, "",
			ukama.ReportTypeUnknown, false, 0, false)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, reps)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ListByOwnerType", func(t *testing.T) {
		isSorted := true

		report := &db.Report{
			Id:        uuid.NewV4(),
			OwnerId:   uuid.NewV4(),
			OwnerType: ukama.OwnerTypeSubscriber,
			NetworkId: uuid.NewV4(),
			IsPaid:    false,
		}

		mock, gdb := prepare_db(t)

		rows := sqlmock.NewRows([]string{"report_id", "owner_id", "owner_type", "network_id", "is_paid"}).
			AddRow(report.Id, report.OwnerId, report.OwnerType, report.NetworkId, report.IsPaid)

		mock.ExpectQuery(`^SELECT.*reports.*`).
			WithArgs(ukama.OwnerTypeSubscriber).
			WillReturnRows(rows)

		r := db.NewReportRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		reps, err := r.List("", ukama.OwnerTypeSubscriber, "",
			ukama.ReportTypeUnknown, false, 0, isSorted)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, reps)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ListNetworkReportsWithCount", func(t *testing.T) {
		var count uint32 = 1

		report := &db.Report{
			Id:        uuid.NewV4(),
			OwnerId:   uuid.NewV4(),
			OwnerType: ukama.OwnerTypeSubscriber,
			NetworkId: uuid.NewV4(),
			IsPaid:    false,
		}

		mock, gdb := prepare_db(t)

		rows := sqlmock.NewRows([]string{"report_id", "owner_id", "owner_type", "network_id", "is_paid"}).
			AddRow(report.Id, report.OwnerId, report.OwnerType, report.NetworkId, report.IsPaid)

		mock.ExpectQuery(`^SELECT.*reports.*`).
			WithArgs(report.NetworkId, count).
			WillReturnRows(rows)

		r := db.NewReportRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		reps, err := r.List("", ukama.OwnerTypeUnknown, report.NetworkId.String(),
			ukama.ReportTypeUnknown, false, count, false)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, reps)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ListPaidInvoiceReports", func(t *testing.T) {
		report := &db.Report{
			Id:        uuid.NewV4(),
			OwnerId:   uuid.NewV4(),
			OwnerType: ukama.OwnerTypeOrg,
			Type:      ukama.ReportTypeInvoice,
			IsPaid:    true,
		}

		mock, gdb := prepare_db(t)

		rows := sqlmock.NewRows([]string{"report_id", "owner_id", "owner_type", "is_paid"}).
			AddRow(report.Id, report.OwnerId, report.OwnerType, report.IsPaid)

		mock.ExpectQuery(`^SELECT.*reports.*`).
			WithArgs(report.Type, report.IsPaid).
			WillReturnRows(rows)

		r := db.NewReportRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		reps, err := r.List("", ukama.OwnerTypeUnknown, "",
			ukama.ReportTypeInvoice, true, 0, false)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, reps)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ListError", func(t *testing.T) {
		mock, gdb := prepare_db(t)

		mock.ExpectQuery(`^SELECT.*reports.*`).
			WithArgs().
			WillReturnError(sql.ErrNoRows)

		r := db.NewReportRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		reps, err := r.List("", ukama.OwnerTypeUnknown, "",
			ukama.ReportTypeUnknown, false, 0, false)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, reps)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestReportRepo_Delete(t *testing.T) {
	t.Run("reportFound", func(t *testing.T) {
		// Arrange
		var reportId = uuid.NewV4()

		mock, gdb := prepare_db(t)

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "reports" SET`)).
			WithArgs(sqlmock.AnyArg(), reportId).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		r := db.NewReportRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		err := r.Delete(reportId, nil)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("reportNotFound", func(t *testing.T) {
		// Arrange
		var reportId = uuid.NewV4()

		mock, gdb := prepare_db(t)

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "reports" SET`)).
			WithArgs(sqlmock.AnyArg(), reportId).
			WillReturnError(sql.ErrNoRows)

		mock.ExpectCommit()

		r := db.NewReportRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		err := r.Delete(reportId, nil)

		// Assert
		assert.Error(t, err)

		err = mock.ExpectationsWereMet()
		assert.Error(t, err)
	})
}

func prepare_db(t *testing.T) (sqlmock.Sqlmock, *gorm.DB) {
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
