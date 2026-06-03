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

	col_db "github.com/ukama/ukama/systems/analytics/collector/pkg/db"
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

func setupMockDB(t *testing.T) (sqlmock.Sqlmock, *gorm.DB) {
	db, mock, err := sqlmock.New()
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

func Test_EventRepo_LogEvent(t *testing.T) {
	t.Run("NewEventIsLogged", func(t *testing.T) {
		// Arrange
		mock, gdb := setupMockDB(t)
		repo := col_db.NewEventRepo(&UkamaDbMock{GormDb: gdb})

		mock.ExpectBegin()
		mock.ExpectQuery(`^INSERT INTO "analytics_event_logs".*ON CONFLICT.*DO NOTHING.*`).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
				sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		// Act
		fresh, err := repo.LogEvent(&col_db.EventLog{
			RoutingKey: "event.cloud.local.org.payments.processor.payment.success",
			MsgId:      "payment-1",
			OccurredAt: time.Now(),
		})

		// Assert
		assert.NoError(t, err)
		assert.True(t, fresh)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("DuplicateEventIsSkipped", func(t *testing.T) {
		// Arrange
		mock, gdb := setupMockDB(t)
		repo := col_db.NewEventRepo(&UkamaDbMock{GormDb: gdb})

		mock.ExpectBegin()
		mock.ExpectQuery(`^INSERT INTO "analytics_event_logs".*ON CONFLICT.*DO NOTHING.*`).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
				sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}))
		mock.ExpectCommit()

		// Act
		fresh, err := repo.LogEvent(&col_db.EventLog{
			RoutingKey: "event.cloud.local.org.payments.processor.payment.success",
			MsgId:      "payment-1",
			OccurredAt: time.Now(),
		})

		// Assert
		assert.NoError(t, err)
		assert.False(t, fresh)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func Test_EventRepo_RecordError(t *testing.T) {
	t.Run("ErrorIsRecorded", func(t *testing.T) {
		// Arrange
		mock, gdb := setupMockDB(t)
		repo := col_db.NewEventRepo(&UkamaDbMock{GormDb: gdb})

		mock.ExpectBegin()
		mock.ExpectQuery(`^INSERT INTO "analytics_event_errors".*`).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
				sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		// Act
		err := repo.RecordError(&col_db.EventError{
			RoutingKey: "event.cloud.local.org.payments.processor.payment.success",
			Reason:     "unmarshal failure",
		})

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func Test_EventRepo_GetRecent(t *testing.T) {
	t.Run("RecentEventsAreReturned", func(t *testing.T) {
		// Arrange
		mock, gdb := setupMockDB(t)
		repo := col_db.NewEventRepo(&UkamaDbMock{GormDb: gdb})

		rows := sqlmock.NewRows([]string{"id", "routing_key", "msg_id"}).
			AddRow(1, "event.cloud.local.org.payments.processor.payment.success", "payment-1")

		mock.ExpectQuery(`^SELECT.*analytics_event_logs.*`).
			WillReturnRows(rows)

		// Act
		logs, err := repo.GetRecent(10)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, logs, 1)
		assert.Equal(t, "payment-1", logs[0].MsgId)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}
