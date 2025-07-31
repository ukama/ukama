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
	"encoding/json"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/tj/assert"
	int_db "github.com/ukama/ukama/systems/notification/event-notify/pkg/db"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type EventSample struct {
	Name string `json:"name"`
	Id   int64  `json:"id"`
}

func GetEventMsg() *int_db.EventMsg {
	es := &EventSample{
		Name: "Sample",
		Id:   1,
	}

	jdata, err := json.Marshal(es)
	if err != nil {
		return nil
	}
	em := int_db.EventMsg{}
	err = em.Data.Set(jdata)
	if err != nil {
		return nil
	}

	return &em
}
func TestEventMsgRepo_Add(t *testing.T) {
	t.Run("Add", func(t *testing.T) {
		// Arrange
		var db *extsql.DB

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		em := GetEventMsg()
		assert.NotNil(t, em)

		mock.ExpectBegin()

		mock.ExpectQuery(regexp.QuoteMeta(`INSERT`)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), em.Data).
			WillReturnRows(sqlmock.NewRows([]string{"id", "data"}).AddRow(1, em.Data))

		mock.ExpectCommit()

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := int_db.NewEventMsgRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		_, err = r.Add(em)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestEventMsgRepo_Get(t *testing.T) {
	t.Run("Get", func(t *testing.T) {
		// Arrange
		var db *extsql.DB
		var id uint = 1
		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		em := GetEventMsg()
		assert.NotNil(t, em)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).
			WithArgs(id).
			WillReturnRows(sqlmock.NewRows([]string{"id", "data"}).AddRow(1, em.Data))

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := int_db.NewEventMsgRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		_, err = r.Get(id)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}
