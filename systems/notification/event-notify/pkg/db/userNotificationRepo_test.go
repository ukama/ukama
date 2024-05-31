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

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/uuid"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/tj/assert"
	notif "github.com/ukama/ukama/systems/common/notification"
	int_db "github.com/ukama/ukama/systems/notification/event-notify/pkg/db"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

var un = int_db.UserNotification{
	Id:             uuid.NewV4(),
	NotificationId: uuid.NewV4(),
	UserId:         uuid.NewV4(),
	IsRead:         false,
	CreatedAt:      time.Now(),
	UpdatedAt:      time.Now(),
}

var n = int_db.Notification{
	Id:          uuid.NewV4(),
	Title:       "Title1",
	Description: "Description1",
	Type:        notif.TYPE_INFO,
	Scope:       notif.SCOPE_ORG,
	CreatedAt:   time.Now(),
	UpdatedAt:   time.Now(),
}

func Test_Add(t *testing.T) {
	t.Run("Add", func(t *testing.T) {
		// Arrange
		var db *extsql.DB

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`INSERT`)).
			WithArgs(un.Id, un.NotificationId, un.UserId, un.IsRead, un.CreatedAt, un.UpdatedAt, sqlmock.AnyArg()).
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
		err = r.Add([]*int_db.UserNotification{&un})

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

// func Test_GetNotificationsByUserID(t *testing.T) {
// 	t.Run("GetNotificationsByUserID", func(t *testing.T) {
// 		// Arrange
// 		var db *extsql.DB

// 		un := int_db.UserNotification{
// 			Id:             uuid.NewV4(),
// 			NotificationId: uuid.NewV4(),
// 			UserId:         uuid.NewV4(),
// 			IsRead:         false,
// 			CreatedAt:      time.Now(),
// 			UpdatedAt:      time.Now(),
// 		}

// 		db, mock, err := sqlmock.New() // mock sql.DB
// 		assert.NoError(t, err)

// 		row := sqlmock.NewRows([]string{"user_notifications.is_read", "notifications.title", "notifications.description", "notifications.scope", "notifications.type", "notifications.id", "notifications.created_at", "notifications.updated_at"}).
// 			AddRow(un.IsRead, n.Title, n.Description, uint8(n.Scope), uint8(n.Type), n.Id, n.CreatedAt, n.UpdatedAt)

// 		mock.ExpectBegin()

// 		mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).
// 			WithArgs(un.UserId.String()).
// 			WillReturnRows(row)

// 		mock.ExpectCommit()

// 		dialector := postgres.New(postgres.Config{
// 			DSN:                  "sqlmock_db_0",
// 			DriverName:           "postgres",
// 			Conn:                 db,
// 			PreferSimpleProtocol: true,
// 		})

// 		gdb, err := gorm.Open(dialector, &gorm.Config{})
// 		assert.NoError(t, err)

// 		r := int_db.NewUserNotificationRepo(&UkamaDbMock{
// 			GormDb: gdb,
// 		})

// 		assert.NoError(t, err)

// 		// Act
// 		resp, err := r.GetNotificationsByUserID(un.UserId.String())

// 		// Assert
// 		assert.NoError(t, err)
// 		assert.NotNil(t, resp)

// 		assert.GreaterOrEqual(t, 1, len(resp))

// 		assert.NotNil(t, resp[0])

// 		assert.Equal(t, un.Id.String(), resp[0].Id.String())

// 		err = mock.ExpectationsWereMet()
// 		assert.NoError(t, err)
// 	})
// }
