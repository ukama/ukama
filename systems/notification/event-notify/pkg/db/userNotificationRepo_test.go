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
	"fmt"
	"regexp"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/uuid"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/tj/assert"
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

func TestUserNotificationRepo_Update(t *testing.T) {
	t.Run("Success_UpdateToRead", func(t *testing.T) {
		// Arrange
		var db *extsql.DB
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		notificationId := uuid.NewV4()

		// Expect transaction begin
		mock.ExpectBegin()

		// Expect the UPDATE query with 3 arguments: is_read, updated_at, notification_id
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE`)).
			WithArgs(true, sqlmock.AnyArg(), notificationId).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Expect transaction commit
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

		// Act
		err = r.Update(notificationId, true)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("Success_UpdateToUnread", func(t *testing.T) {
		// Arrange
		var db *extsql.DB
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		notificationId := uuid.NewV4()

		// Expect transaction begin
		mock.ExpectBegin()

		// Expect the UPDATE query with 3 arguments: is_read, updated_at, notification_id
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE`)).
			WithArgs(false, sqlmock.AnyArg(), notificationId).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Expect transaction commit
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

		// Act
		err = r.Update(notificationId, false)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("Error_UpdateFails", func(t *testing.T) {
		// Arrange
		var db *extsql.DB
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		notificationId := uuid.NewV4()

		// Expect transaction begin
		mock.ExpectBegin()

		// Expect the UPDATE query to fail with 3 arguments
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE`)).
			WithArgs(true, sqlmock.AnyArg(), notificationId).
			WillReturnError(gorm.ErrRecordNotFound)

		// Expect transaction rollback
		mock.ExpectRollback()

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

		// Act
		err = r.Update(notificationId, true)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("Error_NoRowsAffected", func(t *testing.T) {
		// Arrange
		var db *extsql.DB
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		notificationId := uuid.NewV4()

		// Expect transaction begin
		mock.ExpectBegin()

		// Expect the UPDATE query to return no rows affected with 3 arguments
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE`)).
			WithArgs(true, sqlmock.AnyArg(), notificationId).
			WillReturnResult(sqlmock.NewResult(0, 0))

		// Expect transaction commit
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

		// Act
		err = r.Update(notificationId, true)

		// Assert
		assert.NoError(t, err) // GORM doesn't consider 0 rows affected as an error

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestUserNotificationRepo_GetNotificationsByUserID(t *testing.T) {
	t.Run("Success_WithMultipleNotifications", func(t *testing.T) {
		// Arrange
		var db *extsql.DB
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		userId := uuid.NewV4().String()
		notification1 := int_db.Notifications{
			Id:          uuid.NewV4(),
			Title:       "Test Notification 1",
			Description: "Test Description 1",
			Type:        1, // TYPE_INFO
			Scope:       1, // SCOPE_ORG
			IsRead:      false,
			EventKey:    "test.event.1",
			ResourceId:  "resource-1",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		notification2 := int_db.Notifications{
			Id:          uuid.NewV4(),
			Title:       "Test Notification 2",
			Description: "Test Description 2",
			Type:        2, // TYPE_WARNING
			Scope:       2, // SCOPE_NETWORK
			IsRead:      true,
			EventKey:    "test.event.2",
			ResourceId:  "resource-2",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		// Create the expected SQL query
		expectedQuery := fmt.Sprintf("SELECT event_msgs.key AS event_key, user_notifications.is_read, notifications.title, notifications.description, notifications.scope, notifications.type, notifications.id, notifications.created_at, notifications.updated_at, notifications.resource_id FROM user_notifications INNER JOIN notifications ON user_notifications.notification_id = notifications.id INNER JOIN event_msgs ON event_msgs.id = notifications.event_msg_id WHERE user_notifications.user_id = '%s' ORDER BY notifications.created_at DESC;", userId)

		rows := sqlmock.NewRows([]string{
			"event_key", "is_read", "title", "description", "scope", "type", "id", "created_at", "updated_at", "resource_id",
		}).
			AddRow(notification1.EventKey, notification1.IsRead, notification1.Title, notification1.Description, uint8(notification1.Scope), uint8(notification1.Type), notification1.Id, notification1.CreatedAt, notification1.UpdatedAt, notification1.ResourceId).
			AddRow(notification2.EventKey, notification2.IsRead, notification2.Title, notification2.Description, uint8(notification2.Scope), uint8(notification2.Type), notification2.Id, notification2.CreatedAt, notification2.UpdatedAt, notification2.ResourceId)

		mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).
			WillReturnRows(rows)

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

		// Act
		notifications, err := r.GetNotificationsByUserID(userId)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, notifications)
		assert.Len(t, notifications, 2)
		assert.Equal(t, notification1.Title, notifications[0].Title)
		assert.Equal(t, notification2.Title, notifications[1].Title)
		assert.Equal(t, notification1.EventKey, notifications[0].EventKey)
		assert.Equal(t, notification2.EventKey, notifications[1].EventKey)
		assert.Equal(t, notification1.IsRead, notifications[0].IsRead)
		assert.Equal(t, notification2.IsRead, notifications[1].IsRead)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("Success_WithSingleNotification", func(t *testing.T) {
		// Arrange
		var db *extsql.DB
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		userId := uuid.NewV4().String()
		notification := int_db.Notifications{
			Id:          uuid.NewV4(),
			Title:       "Single Notification",
			Description: "Single Description",
			Type:        1, // TYPE_INFO
			Scope:       1, // SCOPE_ORG
			IsRead:      false,
			EventKey:    "single.event",
			ResourceId:  "single-resource",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		expectedQuery := fmt.Sprintf("SELECT event_msgs.key AS event_key, user_notifications.is_read, notifications.title, notifications.description, notifications.scope, notifications.type, notifications.id, notifications.created_at, notifications.updated_at, notifications.resource_id FROM user_notifications INNER JOIN notifications ON user_notifications.notification_id = notifications.id INNER JOIN event_msgs ON event_msgs.id = notifications.event_msg_id WHERE user_notifications.user_id = '%s' ORDER BY notifications.created_at DESC;", userId)

		rows := sqlmock.NewRows([]string{
			"event_key", "is_read", "title", "description", "scope", "type", "id", "created_at", "updated_at", "resource_id",
		}).
			AddRow(notification.EventKey, notification.IsRead, notification.Title, notification.Description, uint8(notification.Scope), uint8(notification.Type), notification.Id, notification.CreatedAt, notification.UpdatedAt, notification.ResourceId)

		mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).
			WillReturnRows(rows)

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

		// Act
		notifications, err := r.GetNotificationsByUserID(userId)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, notifications)
		assert.Len(t, notifications, 1)
		assert.Equal(t, notification.Title, notifications[0].Title)
		assert.Equal(t, notification.EventKey, notifications[0].EventKey)
		assert.Equal(t, notification.IsRead, notifications[0].IsRead)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("Success_WithNoNotifications", func(t *testing.T) {
		// Arrange
		var db *extsql.DB
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		userId := uuid.NewV4().String()

		expectedQuery := fmt.Sprintf("SELECT event_msgs.key AS event_key, user_notifications.is_read, notifications.title, notifications.description, notifications.scope, notifications.type, notifications.id, notifications.created_at, notifications.updated_at, notifications.resource_id FROM user_notifications INNER JOIN notifications ON user_notifications.notification_id = notifications.id INNER JOIN event_msgs ON event_msgs.id = notifications.event_msg_id WHERE user_notifications.user_id = '%s' ORDER BY notifications.created_at DESC;", userId)

		// Return empty result set
		rows := sqlmock.NewRows([]string{
			"event_key", "is_read", "title", "description", "scope", "type", "id", "created_at", "updated_at", "resource_id",
		})

		mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).
			WillReturnRows(rows)

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

		// Act
		notifications, err := r.GetNotificationsByUserID(userId)

		// Assert
		assert.NoError(t, err)
		// Accept both nil and empty slice as valid responses
		assert.True(t, notifications == nil || len(notifications) == 0, "Expected nil or empty slice, got %v", notifications)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("Success_WithEmptyUserId", func(t *testing.T) {
		// Arrange
		var db *extsql.DB
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		userId := ""

		expectedQuery := "SELECT event_msgs.key AS event_key, user_notifications.is_read, notifications.title, notifications.description, notifications.scope, notifications.type, notifications.id, notifications.created_at, notifications.updated_at, notifications.resource_id FROM user_notifications INNER JOIN notifications ON user_notifications.notification_id = notifications.id INNER JOIN event_msgs ON event_msgs.id = notifications.event_msg_id WHERE user_notifications.user_id = '' ORDER BY notifications.created_at DESC;"

		// Return empty result set
		rows := sqlmock.NewRows([]string{
			"event_key", "is_read", "title", "description", "scope", "type", "id", "created_at", "updated_at", "resource_id",
		})

		mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).
			WillReturnRows(rows)

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

		// Act
		notifications, err := r.GetNotificationsByUserID(userId)

		// Assert
		assert.NoError(t, err)
		// Accept both nil and empty slice as valid responses
		assert.True(t, notifications == nil || len(notifications) == 0, "Expected nil or empty slice, got %v", notifications)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}
