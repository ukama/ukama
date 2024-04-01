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

	"github.com/ukama/ukama/systems/billing/invoice/pkg/db"
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

func TestInvoiceRepo_Add(t *testing.T) {
	t.Run("AddINvoice", func(t *testing.T) {
		// Arrange
		invoice := db.Invoice{
			Id:         uuid.NewV4(),
			InvoiceeId: uuid.NewV4(),
		}

		mock, gdb := prepare_db(t)

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`INSERT`)).
			WithArgs(invoice.Id, invoice.InvoiceeId, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
				sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		r := db.NewInvoiceRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		err := r.Add(&invoice, nil)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestInvoiceRepo_Get(t *testing.T) {
	t.Run("InvoiceFound", func(t *testing.T) {
		// Arrange
		var invoiceId = uuid.NewV4()
		var invoiceeId = uuid.NewV4()

		mock, gdb := prepare_db(t)

		rows := sqlmock.NewRows([]string{"id", "invoicee_id"}).
			AddRow(invoiceId, invoiceeId)

		mock.ExpectQuery(`^SELECT.*invoices.*`).
			WithArgs(invoiceId, sqlmock.AnyArg()).
			WillReturnRows(rows)

		r := db.NewInvoiceRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		inv, err := r.Get(invoiceId)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.NotNil(t, inv)
	})

	t.Run("InvoiceNotFound", func(t *testing.T) {
		// Arrange
		var invoiceId = uuid.NewV4()

		mock, gdb := prepare_db(t)

		mock.ExpectQuery(`^SELECT.*invoices.*`).
			WithArgs(invoiceId, sqlmock.AnyArg()).
			WillReturnError(sql.ErrNoRows)

		r := db.NewInvoiceRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		inv, err := r.Get(invoiceId)

		// Assert
		assert.Error(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.Nil(t, inv)
	})
}

func TestInvoiceRepo_List(t *testing.T) {
	t.Run("ListAll", func(t *testing.T) {
		i := &db.Invoice{
			Id:           uuid.NewV4(),
			InvoiceeId:   uuid.NewV4(),
			InvoiceeType: db.InvoiceeTypeOrg,
			IsPaid:       false,
		}

		mock, gdb := prepare_db(t)

		rows := sqlmock.NewRows([]string{"invoice_id", "invoicee_id", "invoicee_type", "is_paid"}).
			AddRow(i.Id, i.InvoiceeId, i.InvoiceeType, i.IsPaid)

		mock.ExpectQuery(`^SELECT.*invoices.*`).
			WithArgs().
			WillReturnRows(rows)

		r := db.NewInvoiceRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		list, err := r.List("", db.InvoiceeTypeUnknown, "",
			false, 0, false)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, list)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ListByInvoiceeId", func(t *testing.T) {
		i := &db.Invoice{
			Id:           uuid.NewV4(),
			InvoiceeId:   uuid.NewV4(),
			InvoiceeType: db.InvoiceeTypeOrg,
			IsPaid:       false,
		}

		mock, gdb := prepare_db(t)

		rows := sqlmock.NewRows([]string{"invoice_id", "invoicee_id", "invoicee_type", "is_paid"}).
			AddRow(i.Id, i.InvoiceeId, i.InvoiceeType, i.IsPaid)

		mock.ExpectQuery(`^SELECT.*invoices.*`).
			WithArgs(i.InvoiceeId).
			WillReturnRows(rows)

		r := db.NewInvoiceRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		list, err := r.List(i.InvoiceeId.String(), db.InvoiceeTypeUnknown, "",
			false, 0, false)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, list)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ListSubscriberInvoices", func(t *testing.T) {
		i := &db.Invoice{
			Id:           uuid.NewV4(),
			InvoiceeId:   uuid.NewV4(),
			InvoiceeType: db.InvoiceeTypeSubscriber,
			NetworkId:    uuid.NewV4(),
			IsPaid:       false,
		}

		mock, gdb := prepare_db(t)

		rows := sqlmock.NewRows([]string{"invoice_id", "invoicee_id", "invoicee_type", "network_id", "is_paid"}).
			AddRow(i.Id, i.InvoiceeId, i.InvoiceeType, i.NetworkId, i.IsPaid)

		mock.ExpectQuery(`^SELECT.*invoices.*`).
			WithArgs(db.InvoiceeTypeSubscriber).
			WillReturnRows(rows)

		r := db.NewInvoiceRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		list, err := r.List("", db.InvoiceeTypeSubscriber, "",
			false, 0, false)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, list)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ListNetworkInvoices", func(t *testing.T) {
		i := &db.Invoice{
			Id:           uuid.NewV4(),
			InvoiceeId:   uuid.NewV4(),
			InvoiceeType: db.InvoiceeTypeSubscriber,
			NetworkId:    uuid.NewV4(),
			IsPaid:       false,
		}

		mock, gdb := prepare_db(t)

		rows := sqlmock.NewRows([]string{"invoice_id", "invoicee_id", "invoicee_type", "network_id", "is_paid"}).
			AddRow(i.Id, i.InvoiceeId, i.InvoiceeType, i.NetworkId, i.IsPaid)

		mock.ExpectQuery(`^SELECT.*invoices.*`).
			WithArgs(i.NetworkId).
			WillReturnRows(rows)

		r := db.NewInvoiceRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		list, err := r.List("", db.InvoiceeTypeUnknown, i.NetworkId.String(),
			false, 0, false)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, list)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ListPaidInvoices", func(t *testing.T) {
		i := &db.Invoice{
			Id:           uuid.NewV4(),
			InvoiceeId:   uuid.NewV4(),
			InvoiceeType: db.InvoiceeTypeOrg,
			IsPaid:       true,
		}

		mock, gdb := prepare_db(t)

		rows := sqlmock.NewRows([]string{"invoice_id", "invoicee_id", "invoicee_type", "is_paid"}).
			AddRow(i.Id, i.InvoiceeId, i.InvoiceeType, i.IsPaid)

		mock.ExpectQuery(`^SELECT.*invoices.*`).
			WithArgs(i.IsPaid).
			WillReturnRows(rows)

		r := db.NewInvoiceRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		list, err := r.List("", db.InvoiceeTypeUnknown, "", true, 1, true)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, list)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ListError", func(t *testing.T) {
		mock, gdb := prepare_db(t)

		mock.ExpectQuery(`^SELECT.*invoices.*`).
			WithArgs().
			WillReturnError(sql.ErrNoRows)

		r := db.NewInvoiceRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		list, err := r.List("", db.InvoiceeTypeUnknown, "",
			false, 0, false)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, list)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestInvoiceRepo_Delete(t *testing.T) {
	t.Run("InvoiceFound", func(t *testing.T) {
		// Arrange
		var invoiceId = uuid.NewV4()

		mock, gdb := prepare_db(t)

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "invoices" SET`)).
			WithArgs(sqlmock.AnyArg(), invoiceId).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		r := db.NewInvoiceRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		err := r.Delete(invoiceId, nil)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("InvoiceNotFound", func(t *testing.T) {
		// Arrange
		var invoiceId = uuid.NewV4()

		mock, gdb := prepare_db(t)

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "invoices" SET`)).
			WithArgs(sqlmock.AnyArg(), invoiceId).
			WillReturnError(sql.ErrNoRows)

		mock.ExpectCommit()

		r := db.NewInvoiceRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		err := r.Delete(invoiceId, nil)

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
