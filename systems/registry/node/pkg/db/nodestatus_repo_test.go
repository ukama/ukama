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
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/tj/assert"
	"github.com/ukama/ukama/systems/common/ukama"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	nodedb "github.com/ukama/ukama/systems/registry/node/pkg/db"
)

func TestNodeStatusRepo_Update(t *testing.T) {
	var (
		nodeId = ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)
		db     *extsql.DB
	)

	ns := &nodedb.NodeStatus{
		NodeId:       nodeId.StringLowercase(),
		Connectivity: ukama.NodeConnectivityOnline,
		State:        ukama.NodeStateUnknown,
	}

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

	r := nodedb.NewNodeStatusRepo(&UkamaDbMock{
		GormDb: gdb,
	})

	t.Run("UpdateSuccess", func(t *testing.T) {
		mock.ExpectBegin()

		mock.ExpectExec(`^UPDATE "node_statuses" SET`).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectCommit()

		// Act
		err = r.Update(ns)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	// No row matched node_id: the status row should already exist (created
	// alongside the Node), so zero rows affected means it's missing, not a
	// no-op update.
	t.Run("NodeNotFound", func(t *testing.T) {
		mock.ExpectBegin()

		mock.ExpectExec(`^UPDATE "node_statuses" SET`).
			WillReturnResult(sqlmock.NewResult(0, 0))

		mock.ExpectCommit()

		// Act
		err = r.Update(ns)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("UpdateQueryError", func(t *testing.T) {
		mock.ExpectBegin()

		mock.ExpectExec(`^UPDATE "node_statuses" SET`).
			WillReturnError(errors.New("internal"))

		mock.ExpectRollback()

		// Act
		err = r.Update(ns)

		// Assert
		assert.Error(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("PartialUpdateOmitsZeroFields", func(t *testing.T) {
		mock.ExpectBegin()

		// Exactly 4 args (updated_at, node_id, connectivity, then the WHERE
		// node_id) proves State was NOT included in the SET clause.
		mock.ExpectExec(`^UPDATE "node_statuses" SET`).
			WithArgs(sqlmock.AnyArg(), nodeId.StringLowercase(), ukama.NodeConnectivityOffline,
				nodeId.StringLowercase()).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectCommit()

		// Act
		err = r.Update(&nodedb.NodeStatus{
			NodeId:       nodeId.StringLowercase(),
			Connectivity: ukama.NodeConnectivityOffline,
		})

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestNodeStatusRepo_Get(t *testing.T) {
	var (
		nodeId = ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)
		db     *extsql.DB
	)

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

	r := nodedb.NewNodeStatusRepo(&UkamaDbMock{
		GormDb: gdb,
	})

	t.Run("StatusFound", func(t *testing.T) {
		row := sqlmock.NewRows([]string{"id", "node_id"}).
			AddRow(1, nodeId.StringLowercase())

		mock.ExpectQuery(`^SELECT.*node_statuses.*`).
			WithArgs(nodeId.StringLowercase(), sqlmock.AnyArg()).
			WillReturnRows(row)

		// Act
		ns, err := r.Get(nodeId)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, ns)
		assert.Equal(t, nodeId.StringLowercase(), ns.NodeId)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("StatusNotFound", func(t *testing.T) {
		mock.ExpectQuery(`^SELECT.*node_statuses.*`).
			WithArgs(nodeId.StringLowercase(), sqlmock.AnyArg()).
			WillReturnError(extsql.ErrNoRows)

		// Act
		ns, err := r.Get(nodeId)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, ns)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestNodeStatusRepo_Delete(t *testing.T) {
	var (
		nodeId = ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)
		db     *extsql.DB
	)

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

	r := nodedb.NewNodeStatusRepo(&UkamaDbMock{
		GormDb: gdb,
	})

	t.Run("DeleteSuccess", func(t *testing.T) {
		mock.ExpectBegin()

		mock.ExpectExec(`UPDATE "node_statuses" SET "deleted_at"=`).
			WithArgs(sqlmock.AnyArg(), nodeId.StringLowercase()).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectCommit()

		// Act
		err = r.Delete(nodeId)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("DeleteError", func(t *testing.T) {
		mock.ExpectBegin()

		mock.ExpectExec(`UPDATE "node_statuses" SET "deleted_at"=`).
			WithArgs(sqlmock.AnyArg(), nodeId.StringLowercase()).
			WillReturnError(errors.New("internal"))

		mock.ExpectRollback()

		// Act
		err = r.Delete(nodeId)

		// Assert
		assert.Error(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestNodeStatusRepo_GetAll(t *testing.T) {
	var (
		nodeId = ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)
		db     *extsql.DB
	)

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

	r := nodedb.NewNodeStatusRepo(&UkamaDbMock{
		GormDb: gdb,
	})

	t.Run("StatusesFound", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "node_id"}).
			AddRow(1, nodeId.StringLowercase())

		mock.ExpectQuery(`^SELECT.*node_statuses.*`).
			WillReturnRows(rows)

		// Act
		ns, err := r.GetAll()

		// Assert
		assert.NoError(t, err)
		assert.Len(t, ns, 1)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("GetAllError", func(t *testing.T) {
		mock.ExpectQuery(`^SELECT.*node_statuses.*`).
			WillReturnError(errors.New("internal"))

		// Act
		ns, err := r.GetAll()

		// Assert
		assert.Error(t, err)
		assert.Nil(t, ns)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestNodeStatusRepo_GetNodeCount(t *testing.T) {
	var db *extsql.DB

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

	r := nodedb.NewNodeStatusRepo(&UkamaDbMock{
		GormDb: gdb,
	})

	t.Run("CountSuccess", func(t *testing.T) {
		mock.ExpectQuery(`^SELECT count\(\*\) FROM "node_statuses".*`).
			WithArgs(ukama.NodeConnectivityOnline).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

		mock.ExpectQuery(`^SELECT count\(\*\) FROM "node_statuses".*`).
			WithArgs(ukama.NodeConnectivityOffline).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		// Act
		online, offline, err := r.GetNodeCount()

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, int64(2), online)
		assert.Equal(t, int64(1), offline)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("OnlineCountError", func(t *testing.T) {
		mock.ExpectQuery(`^SELECT count\(\*\) FROM "node_statuses".*`).
			WithArgs(ukama.NodeConnectivityOnline).
			WillReturnError(errors.New("internal"))

		// Act
		online, offline, err := r.GetNodeCount()

		// Assert
		assert.Error(t, err)
		assert.Equal(t, int64(0), online)
		assert.Equal(t, int64(0), offline)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("OfflineCountError", func(t *testing.T) {
		mock.ExpectQuery(`^SELECT count\(\*\) FROM "node_statuses".*`).
			WithArgs(ukama.NodeConnectivityOnline).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

		mock.ExpectQuery(`^SELECT count\(\*\) FROM "node_statuses".*`).
			WithArgs(ukama.NodeConnectivityOffline).
			WillReturnError(errors.New("internal"))

		// Act
		online, offline, err := r.GetNodeCount()

		// Assert
		assert.Error(t, err)
		assert.Equal(t, int64(0), online)
		assert.Equal(t, int64(0), offline)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}
