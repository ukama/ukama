/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package db_test

import (
	extsql "database/sql"
	"log"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/tj/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/common/uuid"

	cust_db "github.com/ukama/ukama/systems/analytics/customer/pkg/db"
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

func setupMockDb(t *testing.T) (*extsql.DB, sqlmock.Sqlmock, *gorm.DB) {
	t.Helper()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	assert.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})

	gdb, err := gorm.Open(dialector, &gorm.Config{})
	assert.NoError(t, err)

	return db, mock, gdb
}

func TestCustomerRepo_Get(t *testing.T) {
	t.Run("CustomerFound", func(t *testing.T) {
		db, mock, gdb := setupMockDb(t)
		defer db.Close()

		customerId := uuid.NewV4()

		rows := sqlmock.NewRows([]string{"customer_id", "name", "email", "status"}).
			AddRow(customerId, "John Doe", "john@example.com", "active")

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "analytics_customer_snapshots" WHERE customer_id = $1`)).
			WithArgs(customerId, sqlmock.AnyArg()).
			WillReturnRows(rows)

		r := cust_db.NewCustomerRepo(&UkamaDbMock{GormDb: gdb})

		customer, err := r.Get(customerId)

		assert.NoError(t, err)
		assert.NotNil(t, customer)
		assert.Equal(t, "John Doe", customer.Name)
		assert.Equal(t, "active", customer.Status)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("CustomerNotFound", func(t *testing.T) {
		db, mock, gdb := setupMockDb(t)
		defer db.Close()

		customerId := uuid.NewV4()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "analytics_customer_snapshots"`)).
			WithArgs(customerId, sqlmock.AnyArg()).
			WillReturnError(gorm.ErrRecordNotFound)

		r := cust_db.NewCustomerRepo(&UkamaDbMock{GormDb: gdb})

		customer, err := r.Get(customerId)

		assert.Error(t, err)
		assert.Nil(t, customer)
	})
}

func TestCustomerRepo_Counts(t *testing.T) {
	t.Run("CountsWithoutNetworkFilter", func(t *testing.T) {
		db, mock, gdb := setupMockDb(t)
		defer db.Close()

		from := time.Date(2026, 6, 3, 0, 0, 0, 0, time.UTC)
		to := from.AddDate(0, 0, 1)

		countRow := func(n int64) *sqlmock.Rows {
			return sqlmock.NewRows([]string{"count"}).AddRow(n)
		}

		/* total */
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "analytics_customer_snapshots"`)).
			WillReturnRows(countRow(10))

		/* active */
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "analytics_customer_snapshots" WHERE status = $1`)).
			WithArgs("active").
			WillReturnRows(countRow(7))

		/* expired */
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "analytics_customer_snapshots" WHERE status = $1`)).
			WithArgs("expired").
			WillReturnRows(countRow(2))

		/* new (create events in window) */
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "analytics_customer_events" WHERE kind = $1`)).
			WithArgs("create", from, to).
			WillReturnRows(countRow(3))

		/* failed activations in window */
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "analytics_customer_events" WHERE kind = $1`)).
			WithArgs("activation_failed", from, to).
			WillReturnRows(countRow(1))

		r := cust_db.NewCustomerRepo(&UkamaDbMock{GormDb: gdb})

		total, active, newCount, expired, failed, err := r.Counts(uuid.Nil, from, to)

		assert.NoError(t, err)
		assert.Equal(t, uint32(10), total)
		assert.Equal(t, uint32(7), active)
		assert.Equal(t, uint32(3), newCount)
		assert.Equal(t, uint32(2), expired)
		assert.Equal(t, uint32(1), failed)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCustomerRepo_List(t *testing.T) {
	t.Run("ListByStatus", func(t *testing.T) {
		db, mock, gdb := setupMockDb(t)
		defer db.Close()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "analytics_customer_snapshots" WHERE status = $1`)).
			WithArgs("active").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		rows := sqlmock.NewRows([]string{"customer_id", "name", "email", "status"}).
			AddRow(uuid.NewV4(), "Jane Doe", "jane@example.com", "active")

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "analytics_customer_snapshots" WHERE status = $1`)).
			WithArgs("active").
			WillReturnRows(rows)

		r := cust_db.NewCustomerRepo(&UkamaDbMock{GormDb: gdb})

		customers, count, err := r.List(uuid.Nil, uuid.Nil, "active", 1, 20)

		assert.NoError(t, err)
		assert.Equal(t, int64(1), count)
		assert.Len(t, customers, 1)
		assert.Equal(t, "Jane Doe", customers[0].Name)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
