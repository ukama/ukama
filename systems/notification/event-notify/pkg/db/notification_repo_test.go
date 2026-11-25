/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package db_test

import (
	"database/sql/driver"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/tj/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"

	extsql "database/sql"
	notif "github.com/ukama/ukama/systems/common/notification"
	int_db "github.com/ukama/ukama/systems/notification/event-notify/pkg/db"
)

var testOrgId = uuid.NewV4().String()
var testUserId = uuid.NewV4().String()

var n = int_db.Notification{
	Id:          uuid.NewV4(),
	Title:       "Title1",
	Description: "Description1",
	Type:        notif.TYPE_INFO,
	Scope:       notif.SCOPE_ORG,
	ResourceId:  ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_TOWERNODE).String(),
	OrgId:       testOrgId,
	UserId:      testUserId,
	CreatedAt:   time.Now(),
	UpdatedAt:   time.Now(),
}

func TestNotificationRepo_Get(t *testing.T) {
	t.Run("Get", func(t *testing.T) {
		// Arrange
		var db *extsql.DB

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		row := sqlmock.NewRows([]string{"id", "title", "description", "type", "scope",
			"resource_id", "org_id", "user_id"}).
			AddRow(n.Id, n.Title, n.Description, uint8(n.Type), uint8(n.Scope), n.ResourceId,
				n.OrgId, n.UserId)

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

		r := int_db.NewUserNotificationRepo(&UkamaDbMock{
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

type uintConverter struct{}

func (uintConverter) ConvertValue(v any) (driver.Value, error) {
	if vr, ok := v.(driver.Valuer); ok {
		val, err := vr.Value()
		if err != nil {
			return nil, err
		}
		v = val
	}
	if u, ok := v.(uint8); ok {
		return int64(u), nil
	}
	return driver.DefaultParameterConverter.ConvertValue(v)
}

func newMockRepoDb(t *testing.T) (sqlmock.Sqlmock, *UkamaDbMock) {
	t.Helper()
	db, mock, err := sqlmock.New(sqlmock.ValueConverterOption(uintConverter{}))
	assert.NoError(t, err)
	gdb, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	assert.NoError(t, err)
	return mock, &UkamaDbMock{GormDb: gdb}
}

func TestNotificationRepo_Add(t *testing.T) {
	mock, udb := newMockRepoDb(t)
	notification := n

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "notifications"`)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	assert.NoError(t, int_db.NewNotificationRepo(udb).Add(&notification))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestNotificationRepo_AddWithEvent(t *testing.T) {
	insertEvent := regexp.QuoteMeta(`INSERT INTO "event_msgs"`)
	insertNotification := regexp.QuoteMeta(`INSERT INTO "notifications"`)
	insertUsers := regexp.QuoteMeta(`INSERT INTO "user_notifications"`)
	eventId := func() *sqlmock.Rows { return sqlmock.NewRows([]string{"id"}).AddRow(7) }
	users := func() []*int_db.UserNotification {
		return []*int_db.UserNotification{{Id: uuid.NewV4(), NotificationId: n.Id,
			UserId: uuid.NewV4()}}
	}

	t.Run("StoresEventNotificationAndUsers", func(t *testing.T) {
		mock, udb := newMockRepoDb(t)
		notification := n

		mock.ExpectBegin()
		mock.ExpectQuery(insertEvent).WillReturnRows(eventId())
		mock.ExpectExec(insertNotification).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(insertUsers).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := int_db.NewNotificationRepo(udb).AddWithEvent(&int_db.EventMsg{Key: "k"},
			&notification, users())

		assert.NoError(t, err)
		assert.Equal(t, uint(7), notification.EventMsgID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("NoUsersSkipsUserRows", func(t *testing.T) {
		mock, udb := newMockRepoDb(t)
		notification := n

		mock.ExpectBegin()
		mock.ExpectQuery(insertEvent).WillReturnRows(eventId())
		mock.ExpectExec(insertNotification).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := int_db.NewNotificationRepo(udb).AddWithEvent(&int_db.EventMsg{Key: "k"},
			&notification, nil)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("EventInsertErrorRollsBack", func(t *testing.T) {
		mock, udb := newMockRepoDb(t)
		notification := n

		mock.ExpectBegin()
		mock.ExpectQuery(insertEvent).WillReturnError(errors.New("db down"))
		mock.ExpectRollback()

		err := int_db.NewNotificationRepo(udb).AddWithEvent(&int_db.EventMsg{Key: "k"},
			&notification, users())

		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("NotificationInsertErrorRollsBack", func(t *testing.T) {
		mock, udb := newMockRepoDb(t)
		notification := n

		mock.ExpectBegin()
		mock.ExpectQuery(insertEvent).WillReturnRows(eventId())
		mock.ExpectExec(insertNotification).WillReturnError(errors.New("db down"))
		mock.ExpectRollback()

		err := int_db.NewNotificationRepo(udb).AddWithEvent(&int_db.EventMsg{Key: "k"},
			&notification, users())

		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestNotificationRepo_GetNotFound(t *testing.T) {
	mock, udb := newMockRepoDb(t)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).WithArgs(n.Id, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	out, err := int_db.NewNotificationRepo(udb).Get(n.Id)

	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
	assert.Nil(t, out)
	assert.NoError(t, mock.ExpectationsWereMet())
}
