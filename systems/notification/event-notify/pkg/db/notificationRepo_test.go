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
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/tj/assert"
	notif "github.com/ukama/ukama/systems/common/notification"
	"github.com/ukama/ukama/systems/common/uuid"
	int_db "github.com/ukama/ukama/systems/notification/event-notify/pkg/db"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var testOrgId = uuid.NewV4().String()
var testUserId = uuid.NewV4().String()

var n = int_db.Notification{
	Id:          uuid.NewV4(),
	Title:       "Title1",
	Description: "Description1",
	Type:        notif.TYPE_INFO,
	Scope:       notif.SCOPE_ORG,
	ResourceId:  uuid.NewV4(),
	OrgId:       testOrgId,
	UserId:      testUserId,
	CreatedAt:   time.Now(),
	UpdatedAt:   time.Now(),
}

var evtMsg = int_db.EventMsg{
	Model: gorm.Model{
		ID: 1,
	},
}

func TestNotificationRepo_Add(t *testing.T) {
	t.Run("Add", func(t *testing.T) {
		// Arrange
		var db *extsql.DB

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`INSERT`)).
			WithArgs(n.Id, n.Title, n.Description, n.Type, n.Scope, n.ResourceId, n.OrgId, n.NetworkId, n.SubscriberId, n.UserId, n.NodeId, evtMsg.ID, evtMsg, n.CreatedAt, n.UpdatedAt, sqlmock.AnyArg()).
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

		r := int_db.NewNotificationRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		err = r.Add(&n)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestNotificationRepo_Get(t *testing.T) {
	t.Run("Get", func(t *testing.T) {
		// Arrange
		var db *extsql.DB

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		row := sqlmock.NewRows([]string{"id", "title", "description", "type", "scope", "resource_id", "org_id", "user_id"}).
			AddRow(n.Id, n.Title, n.Description, uint8(n.Type), uint8(n.Scope), n.ResourceId, n.OrgId, n.UserId)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).
			WithArgs(n.Id, 1).
			WillReturnRows(row)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := int_db.NewNotificationRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		_, err = r.Get(n.Id)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestNotificationRepo_Update(t *testing.T) {
	t.Run("Update", func(t *testing.T) {
		// Arrange
		var db *extsql.DB

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE`)).
			WithArgs(true, sqlmock.AnyArg(), n.Id).
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

		r := int_db.NewNotificationRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		err = r.Update(n.Id, true)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}
