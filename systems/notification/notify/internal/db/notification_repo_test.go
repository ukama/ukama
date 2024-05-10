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
	"time"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/tj/assert"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	intdb "github.com/ukama/ukama/systems/notification/notify/internal/db"
	jdb "gorm.io/datatypes"
)

type UkamaDbMock struct {
	GormDb *gorm.DB
}

func (u UkamaDbMock) Init(model ...interface{}) error {
	panic("implement me")
}

func (u UkamaDbMock) Connect() error {
	panic("implement me")
}

func (u UkamaDbMock) GetGormDb() *gorm.DB {
	return u.GormDb
}

func (u UkamaDbMock) InitDB() error {
	return nil
}

func (u UkamaDbMock) ExecuteInTransaction(dbOperation func(tx *gorm.DB) *gorm.DB,
	nestedFuncs ...func() error) error {
	panic("implement me")
}

func (u UkamaDbMock) ExecuteInTransaction2(dbOperation func(tx *gorm.DB) *gorm.DB,
	nestedFuncs ...func(tx *gorm.DB) error) (err error) {
	panic("implement me")
}

func TestNotificationRepo_Insert(t *testing.T) {
	node := ukama.NewVirtualHomeNodeId()
	nt := NewTestDbNotification(node.String(), "alert")
	var err error

	mock, gdb := prepare_db(t)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`INSERT`)).
		WithArgs(nt.Id, nt.NodeId, nt.NodeType, nt.Severity, nt.Type,
			nt.ServiceName, nt.Status, nt.Time, nt.Description, nt.Details,
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	r := intdb.NewNotificationRepo(&UkamaDbMock{
		GormDb: gdb,
	})

	assert.NoError(t, err)

	// Act
	err = r.Add(&nt)

	// Assert
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestNotificationRepo_Get(t *testing.T) {
	t.Run("NotificationFound", func(t *testing.T) {
		// Arrange
		var notificationId = uuid.NewV4()
		var node = ukama.NewVirtualHomeNodeId()

		mock, gdb := prepare_db(t)

		rows := sqlmock.NewRows([]string{"id", "node_id"}).
			AddRow(notificationId, node.String())

		mock.ExpectQuery(`^SELECT.*notifications.*`).
			WithArgs(notificationId, sqlmock.AnyArg()).
			WillReturnRows(rows)

		r := intdb.NewNotificationRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		notification, err := r.Get(notificationId)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, notification)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("NotificationNotFound", func(t *testing.T) {
		// Arrange
		var notificationId = uuid.NewV4()

		mock, gdb := prepare_db(t)

		mock.ExpectQuery(`^SELECT.*notifications.*`).
			WithArgs(notificationId, sqlmock.AnyArg()).
			WillReturnError(sql.ErrNoRows)

		r := intdb.NewNotificationRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		notification, err := r.Get(notificationId)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, notification)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestNotificationRepo_List(t *testing.T) {
	t.Run("ListAll", func(t *testing.T) {
		node := ukama.NewVirtualHomeNodeId()
		nt := NewTestDbNotification(node.String(), "alert")
		var err error

		mock, gdb := prepare_db(t)

		rows := sqlmock.NewRows([]string{"notification_id", "node_id", "node_type",
			"severity", "service_name", "status", "time", "description", "details"}).
			AddRow(nt.Id, nt.NodeId, nt.NodeType, nt.Severity, nt.ServiceName,
				nt.Status, nt.Time, nt.Description, nt.Details)

		mock.ExpectQuery(`^SELECT.*notifications.*`).
			WithArgs().
			WillReturnRows(rows)

		r := intdb.NewNotificationRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		list, err := r.List("", "", "", 0, false)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, list)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ListAlertsForService", func(t *testing.T) {
		node := ukama.NewVirtualHomeNodeId()
		nt := NewTestDbNotification(node.String(), "alert")
		var err error

		mock, gdb := prepare_db(t)

		rows := sqlmock.NewRows([]string{"notification_id", "node_id", "node_type",
			"severity", "service_name", "status", "time", "description", "details"}).
			AddRow(nt.Id, nt.NodeId, nt.NodeType, nt.Severity, nt.ServiceName,
				nt.Status, nt.Time, nt.Description, nt.Details)

		mock.ExpectQuery(`^SELECT.*notifications.*`).
			WithArgs(nt.ServiceName, string(nt.Type)).
			WillReturnRows(rows)

		r := intdb.NewNotificationRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		list, err := r.List("", nt.ServiceName, nt.Type.String(), 0, false)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, list)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ListEventsForNode", func(t *testing.T) {
		node := ukama.NewVirtualHomeNodeId()
		nt := NewTestDbNotification(node.String(), "event")
		var err error

		mock, gdb := prepare_db(t)

		rows := sqlmock.NewRows([]string{"notification_id", "node_id", "node_type",
			"severity", "service_name", "status", "time", "description", "details"}).
			AddRow(nt.Id, nt.NodeId, nt.NodeType, nt.Severity, nt.ServiceName,
				nt.Status, nt.Time, nt.Description, nt.Details)

		mock.ExpectQuery(`^SELECT.*notifications.*`).
			WithArgs(nt.NodeId, string(nt.Type)).
			WillReturnRows(rows)

		r := intdb.NewNotificationRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		list, err := r.List(nt.NodeId, "", nt.Type.String(), 0, false)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, list)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ListSortedEventsForService", func(t *testing.T) {
		node := ukama.NewVirtualHomeNodeId()
		nt := NewTestDbNotification(node.String(), "event")
		var err error

		mock, gdb := prepare_db(t)

		rows := sqlmock.NewRows([]string{"notification_id", "node_id", "node_type",
			"severity", "service_name", "status", "time", "description", "details"}).
			AddRow(nt.Id, nt.NodeId, nt.NodeType, nt.Severity, nt.ServiceName,
				nt.Status, nt.Time, nt.Description, nt.Details)

		mock.ExpectQuery(`^SELECT.*notifications.*`).
			WithArgs(nt.ServiceName, nt.Type, sqlmock.AnyArg()).
			WillReturnRows(rows)

		r := intdb.NewNotificationRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		list, err := r.List("", nt.ServiceName, nt.Type.String(), 1, true)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, list)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ListSortedAlertsForNode", func(t *testing.T) {
		node := ukama.NewVirtualHomeNodeId()
		nt := NewTestDbNotification(node.String(), "alert")
		var err error

		mock, gdb := prepare_db(t)

		rows := sqlmock.NewRows([]string{"notification_id", "node_id", "node_type",
			"severity", "service_name", "status", "time", "description", "details"}).
			AddRow(nt.Id, nt.NodeId, nt.NodeType, nt.Severity, nt.ServiceName,
				nt.Status, nt.Time, nt.Description, nt.Details)

		mock.ExpectQuery(`^SELECT.*notifications.*`).
			WithArgs(nt.NodeId, nt.Type, sqlmock.AnyArg()).
			WillReturnRows(rows)

		r := intdb.NewNotificationRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		list, err := r.List(nt.NodeId, "", nt.Type.String(), 1, true)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, list)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestNotificationRepo_Delete(t *testing.T) {
	t.Run("NotificationFound", func(t *testing.T) {
		// Arrange
		var notificationId = uuid.NewV4()

		mock, gdb := prepare_db(t)

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "notifications" SET`)).
			WithArgs(sqlmock.AnyArg(), notificationId).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		r := intdb.NewNotificationRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		err := r.Delete(notificationId)

		// Assert
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("NotificationNotFound", func(t *testing.T) {
		// Arrange
		var notificationId = uuid.NewV4()

		mock, gdb := prepare_db(t)

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "notifications" SET`)).
			WithArgs(sqlmock.AnyArg(), notificationId).
			WillReturnError(sql.ErrNoRows)

		// mock.ExpectCommit()

		r := intdb.NewNotificationRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		err := r.Delete(notificationId)

		// Assert
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

}

func TestNotificationRepo_Purge(t *testing.T) {
	t.Run("DeleteAll", func(t *testing.T) {
		var err error

		mock, gdb := prepare_db(t)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE`)).
			WithArgs(sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		r := intdb.NewNotificationRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		_, err = r.Purge("", "", "")

		// Assert
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DeleteAlertsForNode", func(t *testing.T) {
		node := ukama.NewVirtualHomeNodeId()
		nt := NewTestDbNotification(node.String(), "alert")
		var err error

		mock, gdb := prepare_db(t)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("UPDATE")).
			WithArgs(sqlmock.AnyArg(), nt.NodeId, nt.Type.String()).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		r := intdb.NewNotificationRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		_, err = r.Purge(nt.NodeId, "", nt.Type.String())

		// Assert
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DeleteEventsForService", func(t *testing.T) {
		node := ukama.NewVirtualHomeNodeId()
		nt := NewTestDbNotification(node.String(), "alert")
		var err error

		mock, gdb := prepare_db(t)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("UPDATE")).
			WithArgs(sqlmock.AnyArg(), nt.ServiceName, nt.Type.String()).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		r := intdb.NewNotificationRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		_, err = r.Purge("", nt.ServiceName, nt.Type.String())

		// Assert
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func NewTestDbNotification(nodeId string, ntype string) intdb.Notification {
	return intdb.Notification{
		Id:          uuid.NewV4(),
		NodeId:      nodeId,
		NodeType:    *ukama.GetNodeType(nodeId),
		Severity:    intdb.SeverityType("high"),
		Type:        intdb.NotificationType(ntype),
		ServiceName: "noded",
		Status:      8200,
		Time:        uint32(time.Now().Unix()),
		Description: "Some random alert",
		Details:     jdb.JSON(`{"reason": "testing", "component":"router_test"}`),
	}
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
