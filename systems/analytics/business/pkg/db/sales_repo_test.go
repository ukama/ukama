/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package db_test

import (
	"log"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/tj/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/common/uuid"

	biz_db "github.com/ukama/ukama/systems/analytics/business/pkg/db"
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

func setupMockDB(t *testing.T) (sqlmock.Sqlmock, biz_db.SalesRepo) {
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

	repo := biz_db.NewSalesRepo(&UkamaDbMock{
		GormDb: gdb,
	})

	return mock, repo
}

func Test_SalesRepo_RevenueBetween(t *testing.T) {
	from := time.Date(2026, time.June, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, time.June, 3, 0, 0, 0, 0, time.UTC)

	t.Run("RevenueWithNetworkFilter", func(t *testing.T) {
		// Arrange
		mock, repo := setupMockDB(t)

		networkId := uuid.NewV4().String()

		rows := sqlmock.NewRows([]string{"coalesce"}).AddRow(150.5)

		mock.ExpectQuery(`SELECT COALESCE\(SUM\(amount_cents\), 0\) / 100.0 FROM "analytics_payment_events"`).
			WithArgs("success", from, to, networkId).
			WillReturnRows(rows)

		// Act
		revenue, err := repo.RevenueBetween(networkId, from, to)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, 150.5, revenue)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("RevenueWithoutNetworkFilter", func(t *testing.T) {
		// Arrange
		mock, repo := setupMockDB(t)

		rows := sqlmock.NewRows([]string{"coalesce"}).AddRow(0.0)

		mock.ExpectQuery(`SELECT COALESCE\(SUM\(amount_cents\), 0\) / 100.0 FROM "analytics_payment_events"`).
			WithArgs("success", from, to).
			WillReturnRows(rows)

		// Act
		revenue, err := repo.RevenueBetween("", from, to)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, 0.0, revenue)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("QueryError", func(t *testing.T) {
		// Arrange
		mock, repo := setupMockDB(t)

		mock.ExpectQuery(`SELECT COALESCE\(SUM\(amount_cents\), 0\) / 100.0 FROM "analytics_payment_events"`).
			WillReturnError(gorm.ErrInvalidDB)

		// Act
		_, err := repo.RevenueBetween("", from, to)

		// Assert
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
