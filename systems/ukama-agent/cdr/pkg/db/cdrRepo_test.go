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

	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	int_db "github.com/ukama/ukama/systems/ukama-agent/cdr/pkg/db"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var cdr = int_db.CDR{
	Session:       1,
	NodeId:        ukama.NewVirtualHomeNodeId().String(),
	Imsi:          "123456789012345678",
	Policy:        uuid.NewV4().String(),
	ApnName:       "ukama.co",
	Ip:            "192.168.8.2",
	StartTime:     uint64(time.Now().Unix() - 100000),
	EndTime:       uint64(time.Now().Unix() - 50000),
	LastUpdatedAt: uint64(time.Now().Unix() - 50000),
	TxBytes:       2048000,
	RxBytes:       1024000,
	TotalBytes:    3072000,
}

func TestCDRRepo_Add(t *testing.T) {

	t.Run("Add", func(t *testing.T) {
		// Arrange
		var db *extsql.DB
		var err error

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		mock.ExpectBegin()

		mock.ExpectQuery(regexp.QuoteMeta(`INSERT`)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), cdr.Session, cdr.NodeId, cdr.Imsi, cdr.Policy, cdr.ApnName, cdr.Ip, cdr.StartTime, cdr.EndTime, cdr.LastUpdatedAt, cdr.TxBytes, cdr.RxBytes, cdr.TotalBytes).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		mock.ExpectCommit()

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := int_db.NewCDRRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		err = r.Add(&cdr)

		// Assert12
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)

	})

}

func TestCDRRepo_Get(t *testing.T) {

	t.Run("ImsiFound", func(t *testing.T) {
		var ID uint = 1
		// Arrange
		var db *extsql.DB
		var err error

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		crow := sqlmock.NewRows([]string{"ID", "session", "node_id", "imsi", "policy", "apn_name", "ip", "start_time", "end_time", "last_updated_at", "tx_bytes", "rx_bytes", "total_bytes"}).
			AddRow(ID, cdr.Session, cdr.NodeId, cdr.Imsi, cdr.Policy, cdr.ApnName, cdr.Ip, cdr.StartTime, cdr.EndTime, cdr.LastUpdatedAt, cdr.TxBytes, cdr.RxBytes, cdr.TotalBytes)

		mock.ExpectQuery(`^SELECT.*cdrs.*`).
			WithArgs(cdr.Imsi).
			WillReturnRows(crow)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := int_db.NewCDRRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		tu, err := r.GetByImsi(cdr.Imsi)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)

		if assert.NotNil(t, tu) {
			c := *tu
			assert.EqualValues(t, ID, c[0].ID)
		}

	})

	t.Run("ImsiNotFound", func(t *testing.T) {

		// Arrange
		var db *extsql.DB
		var err error

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		// urow := sqlmock.NewRows([]string{"ID", "usage", "historical", "last_session_id", "last_session_usage", "last_node_id", "last_cdr_updated_at", "policy"}).
		// 	AddRow(ID, usage.Imsi, usage.Historical, usage.Usage, usage.LastSessionUsage, usage.LastSessionId, usage.LastNodeId, usage.LastCDRUpdatedAt, usage.Policy)

		mock.ExpectQuery(`^SELECT.*cdrs.*`).
			WithArgs(cdr.Imsi).
			WillReturnError(gorm.ErrRecordNotFound)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := int_db.NewCDRRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		_, err = r.GetByImsi(cdr.Imsi)

		// Assert
		assert.Error(t, err)
		if assert.Error(t, err) {
			assert.Equal(t, gorm.ErrRecordNotFound, err)
		}

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)

	})

	t.Run("Session", func(t *testing.T) {
		var ID uint = 1
		// Arrange
		var db *extsql.DB
		var err error

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		crow := sqlmock.NewRows([]string{"ID", "session", "node_id", "imsi", "policy", "apn_name", "ip", "start_time", "end_time", "last_updated_at", "tx_bytes", "rx_bytes", "total_bytes"}).
			AddRow(ID, cdr.Session, cdr.NodeId, cdr.Imsi, cdr.Policy, cdr.ApnName, cdr.Ip, cdr.StartTime, cdr.EndTime, cdr.LastUpdatedAt, cdr.TxBytes, cdr.RxBytes, cdr.TotalBytes)

		mock.ExpectQuery(`^SELECT.*cdrs.*`).
			WithArgs(cdr.Imsi, cdr.Session).
			WillReturnRows(crow)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := int_db.NewCDRRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		tu, err := r.GetBySession(cdr.Imsi, cdr.Session)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)

		if assert.NotNil(t, tu) {
			c := *tu
			assert.EqualValues(t, ID, c[0].ID)
		}

	})

	t.Run("ByTime", func(t *testing.T) {
		var ID uint = 1
		// Arrange
		var db *extsql.DB
		var err error

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		crow := sqlmock.NewRows([]string{"ID", "session", "node_id", "imsi", "policy", "apn_name", "ip", "start_time", "end_time", "last_updated_at", "tx_bytes", "rx_bytes", "total_bytes"}).
			AddRow(ID, cdr.Session, cdr.NodeId, cdr.Imsi, cdr.Policy, cdr.ApnName, cdr.Ip, cdr.StartTime, cdr.EndTime, cdr.LastUpdatedAt, cdr.TxBytes, cdr.RxBytes, cdr.TotalBytes)

		mock.ExpectQuery(`^SELECT.*cdrs.*`).
			WithArgs(cdr.Imsi, cdr.StartTime, cdr.EndTime, cdr.EndTime).
			WillReturnRows(crow)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := int_db.NewCDRRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		tu, err := r.GetByTime(cdr.Imsi, cdr.StartTime, cdr.EndTime)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)

		if assert.NotNil(t, tu) {
			c := *tu
			assert.EqualValues(t, ID, c[0].ID)
		}

	})

	t.Run("ByTimeAndNodeId", func(t *testing.T) {
		var ID uint = 1
		// Arrange
		var db *extsql.DB
		var err error

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		crow := sqlmock.NewRows([]string{"ID", "session", "node_id", "imsi", "policy", "apn_name", "ip", "start_time", "end_time", "last_updated_at", "tx_bytes", "rx_bytes", "total_bytes"}).
			AddRow(ID, cdr.Session, cdr.NodeId, cdr.Imsi, cdr.Policy, cdr.ApnName, cdr.Ip, cdr.StartTime, cdr.EndTime, cdr.LastUpdatedAt, cdr.TxBytes, cdr.RxBytes, cdr.TotalBytes)

		mock.ExpectQuery(`^SELECT.*cdrs.*`).
			WithArgs(cdr.Imsi, cdr.StartTime, cdr.EndTime, cdr.EndTime, cdr.NodeId).
			WillReturnRows(crow)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := int_db.NewCDRRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		tu, err := r.GetByTimeAndNodeId(cdr.Imsi, cdr.StartTime, cdr.EndTime, cdr.NodeId)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)

		if assert.NotNil(t, tu) {
			c := *tu
			assert.EqualValues(t, ID, c[0].ID)
		}

	})

	t.Run("ByPolicy", func(t *testing.T) {
		var ID uint = 1
		// Arrange
		var db *extsql.DB
		var err error

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		crow := sqlmock.NewRows([]string{"ID", "session", "node_id", "imsi", "policy", "apn_name", "ip", "start_time", "end_time", "last_updated_at", "tx_bytes", "rx_bytes", "total_bytes"}).
			AddRow(ID, cdr.Session, cdr.NodeId, cdr.Imsi, cdr.Policy, cdr.ApnName, cdr.Ip, cdr.StartTime, cdr.EndTime, cdr.LastUpdatedAt, cdr.TxBytes, cdr.RxBytes, cdr.TotalBytes)

		mock.ExpectQuery(`^SELECT.*cdrs.*`).
			WithArgs(cdr.Imsi, cdr.Policy).
			WillReturnRows(crow)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := int_db.NewCDRRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		tu, err := r.GetByPolicy(cdr.Imsi, cdr.Policy)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)

		if assert.NotNil(t, tu) {
			c := *tu
			assert.EqualValues(t, ID, c[0].ID)
		}

	})

	t.Run("ByFiltersWithPolicy", func(t *testing.T) {
		var ID uint = 1
		// Arrange
		var db *extsql.DB
		var err error

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		crow := sqlmock.NewRows([]string{"ID", "session", "node_id", "imsi", "policy", "apn_name", "ip", "start_time", "end_time", "last_updated_at", "tx_bytes", "rx_bytes", "total_bytes"}).
			AddRow(ID, cdr.Session, cdr.NodeId, cdr.Imsi, cdr.Policy, cdr.ApnName, cdr.Ip, cdr.StartTime, cdr.EndTime, cdr.LastUpdatedAt, cdr.TxBytes, cdr.RxBytes, cdr.TotalBytes)

		mock.ExpectQuery(`^SELECT.*cdrs.*`).
			WithArgs(cdr.Imsi, cdr.Session, cdr.Policy, cdr.StartTime, cdr.EndTime, cdr.EndTime).
			WillReturnRows(crow)

		mock.ExpectQuery(`^SELECT.*cdrs.*`).
			WithArgs(cdr.Imsi, cdr.Session, cdr.Policy, cdr.StartTime, cdr.EndTime, cdr.EndTime).
			WillReturnRows(crow)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := int_db.NewCDRRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		tu, err := r.GetByFilters(cdr.Imsi, cdr.Session, cdr.Policy, cdr.StartTime, cdr.EndTime)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)

		if assert.NotNil(t, tu) {
			c := *tu
			assert.EqualValues(t, ID, c[0].ID)
		}

	})

	t.Run("ByFiltersWithoutPolicy", func(t *testing.T) {
		var ID uint = 1
		// Arrange
		var db *extsql.DB
		var err error

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		crow := sqlmock.NewRows([]string{"ID", "session", "node_id", "imsi", "policy", "apn_name", "ip", "start_time", "end_time", "last_updated_at", "tx_bytes", "rx_bytes", "total_bytes"}).
			AddRow(ID, cdr.Session, cdr.NodeId, cdr.Imsi, cdr.Policy, cdr.ApnName, cdr.Ip, cdr.StartTime, cdr.EndTime, cdr.LastUpdatedAt, cdr.TxBytes, cdr.RxBytes, cdr.TotalBytes)

		mock.ExpectQuery(`^SELECT.*cdrs.*`).
			WithArgs(cdr.Imsi, cdr.Session, cdr.StartTime, cdr.EndTime, cdr.EndTime).
			WillReturnRows(crow)

		mock.ExpectQuery(`^SELECT.*cdrs.*`).
			WithArgs(cdr.Imsi, cdr.Session, cdr.StartTime, cdr.EndTime, cdr.EndTime).
			WillReturnRows(crow)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := int_db.NewCDRRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		tu, err := r.GetByFilters(cdr.Imsi, cdr.Session, "", cdr.StartTime, cdr.EndTime)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)

		if assert.NotNil(t, tu) {
			c := *tu
			assert.EqualValues(t, ID, c[0].ID)
		}

	})

	t.Run("ByFiltersWithSession", func(t *testing.T) {
		var ID uint = 1
		var session uint64 = 0
		// Arrange
		var db *extsql.DB
		var err error

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		crow := sqlmock.NewRows([]string{"ID", "session", "node_id", "imsi", "policy", "apn_name", "ip", "start_time", "end_time", "last_updated_at", "tx_bytes", "rx_bytes", "total_bytes"}).
			AddRow(ID, cdr.Session, cdr.NodeId, cdr.Imsi, cdr.Policy, cdr.ApnName, cdr.Ip, cdr.StartTime, cdr.EndTime, cdr.LastUpdatedAt, cdr.TxBytes, cdr.RxBytes, cdr.TotalBytes)

		mock.ExpectQuery(`^SELECT.*cdrs.*`).
			WithArgs(cdr.Imsi, session, cdr.Policy, cdr.StartTime, cdr.EndTime, cdr.EndTime).
			WillReturnRows(crow)

		mock.ExpectQuery(`^SELECT.*cdrs.*`).
			WithArgs(cdr.Imsi, session, cdr.Policy, cdr.StartTime, cdr.EndTime, cdr.EndTime).
			WillReturnRows(crow)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := int_db.NewCDRRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		tu, err := r.GetByFilters(cdr.Imsi, session, cdr.Policy, cdr.StartTime, cdr.EndTime)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)

		if assert.NotNil(t, tu) {
			c := *tu
			assert.EqualValues(t, ID, c[0].ID)
		}

	})

	t.Run("ByFiltersWithoutPolicyAndSession", func(t *testing.T) {
		var ID uint = 1
		var session uint64 = 0
		// Arrange
		var db *extsql.DB
		var err error

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		crow := sqlmock.NewRows([]string{"ID", "session", "node_id", "imsi", "policy", "apn_name", "ip", "start_time", "end_time", "last_updated_at", "tx_bytes", "rx_bytes", "total_bytes"}).
			AddRow(ID, cdr.Session, cdr.NodeId, cdr.Imsi, cdr.Policy, cdr.ApnName, cdr.Ip, cdr.StartTime, cdr.EndTime, cdr.LastUpdatedAt, cdr.TxBytes, cdr.RxBytes, cdr.TotalBytes)

		mock.ExpectQuery(`^SELECT.*cdrs.*`).
			WithArgs(cdr.Imsi, session, cdr.StartTime, cdr.EndTime, cdr.EndTime).
			WillReturnRows(crow)

		mock.ExpectQuery(`^SELECT.*cdrs.*`).
			WithArgs(cdr.Imsi, session, cdr.StartTime, cdr.EndTime, cdr.EndTime).
			WillReturnRows(crow)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := int_db.NewCDRRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		tu, err := r.GetByFilters(cdr.Imsi, session, "", cdr.StartTime, cdr.EndTime)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)

		if assert.NotNil(t, tu) {
			c := *tu
			assert.EqualValues(t, ID, c[0].ID)
		}

	})
}
