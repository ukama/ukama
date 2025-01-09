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
	"log"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/tj/assert"
	"github.com/ukama/ukama/systems/common/ukama"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	nodedb "github.com/ukama/ukama/systems/registry/node/pkg/db"
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

func TestNodeRepo_Add(t *testing.T) {
	var nodeId = ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)

	var db *extsql.DB

	node := nodedb.Node{
		Id:   nodeId.String(),
		Name: "node-1",
		Type: "hnode",
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

	r := nodedb.NewNodeRepo(&UkamaDbMock{
		GormDb: gdb,
	})

	t.Run("AddNode", func(t *testing.T) {
		// Arrange
		mock.ExpectBegin()

		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "nodes" ("id","name","type","parent_node_id","created_at","updated_at","deleted_at") VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "latitude","longitude"`)).
			WithArgs(node.Id, node.Name, node.Type, node.ParentNodeId, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"latitude", "longitude"}).AddRow(0.0, 0.0))
		mock.ExpectCommit()

		// Act
		err = r.Add(&node, nil)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestNodeRepo_Get(t *testing.T) {
	var nodeId = ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)
	var name = "node-1"

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

	r := nodedb.NewNodeRepo(&UkamaDbMock{
		GormDb: gdb,
	})

	t.Run("NodeFound", func(t *testing.T) {
		// Arrange
		row := sqlmock.NewRows([]string{"id", "name"}).
			AddRow(nodeId, name)

		mock.ExpectQuery(`^SELECT.*nodes.*`).
			WithArgs(nodeId, sqlmock.AnyArg()).
			WillReturnRows(row)

		mock.ExpectQuery(`^SELECT.*parent_node_id.*`).
			WithArgs(nodeId).
			WillReturnRows(row)

		mock.ExpectQuery(`^SELECT.*sites.*`).
			WithArgs(nodeId).
			WillReturnRows(row)

		mock.ExpectQuery(`^SELECT.*node_statuses.*`).
			WithArgs(nodeId).
			WillReturnRows(row)

		// Act
		node, err := r.Get(nodeId)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, node)

		assert.Equal(t, nodeId.String(), node.Id)
		assert.Equal(t, name, node.Name)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("NodeNotFound", func(t *testing.T) {
		mock.ExpectQuery(`^SELECT.*nodes.*`).
			WithArgs(nodeId, sqlmock.AnyArg()).
			WillReturnError(extsql.ErrNoRows)

		// Act
		node, err := r.Get(nodeId)

		// Assert
		assert.Error(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.Nil(t, node)
	})
}

func TestNodeRepo_GetAll(t *testing.T) {
	var nodeId = ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)
	var name = "node-1"

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

	r := nodedb.NewNodeRepo(&UkamaDbMock{
		GormDb: gdb,
	})

	t.Run("NodeFound", func(t *testing.T) {
		// Arrange
		row := sqlmock.NewRows([]string{"id", "name"}).
			AddRow(nodeId, name)

		mock.ExpectQuery(`^SELECT.*nodes.*`).
			WillReturnRows(row)

		mock.ExpectQuery(`^SELECT.*parent_node_id.*`).
			WithArgs(nodeId).
			WillReturnRows(row)

		mock.ExpectQuery(`^SELECT.*sites.*`).
			WithArgs(nodeId).
			WillReturnRows(row)

		mock.ExpectQuery(`^SELECT.*node_statuses.*`).
			WithArgs(nodeId).
			WillReturnRows(row)

		// Act
		nodes, err := r.GetAll()

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, nodes)

		assert.Equal(t, nodeId.String(), nodes[0].Id)
		assert.Equal(t, name, nodes[0].Name)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("UnknownError", func(t *testing.T) {
		mock.ExpectQuery(`^SELECT.*nodes.*`).
			WillReturnError(errors.New("internal"))

		// Act
		node, err := r.GetAll()

		// Assert
		assert.Error(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.Nil(t, node)
	})
}

func TestNodeRepo_Delete(t *testing.T) {
	// Arrange
	var nodeId = ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)
	var name = "node-1"

	var db *extsql.DB

	row := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(nodeId, name)

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

	r := nodedb.NewNodeRepo(&UkamaDbMock{
		GormDb: gdb,
	})

	assert.NoError(t, err)

	t.Run("NodeFound", func(t *testing.T) {
		mock.ExpectQuery(`^SELECT.*nodes.*`).
			WithArgs(nodeId, sqlmock.AnyArg()).
			WillReturnRows(row)

		mock.ExpectQuery(`^SELECT.*parent_node_id.*`).
			WithArgs(nodeId).
			WillReturnRows(row)

		mock.ExpectQuery(`^SELECT.*sites.*`).
			WithArgs(nodeId).
			WillReturnRows(row)

		mock.ExpectQuery(`^SELECT.*node_statuses.*`).
			WithArgs(nodeId).
			WillReturnRows(row)

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE`)).
			WithArgs(sqlmock.AnyArg(), nodeId).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE`)).
			WithArgs(sqlmock.AnyArg(), nodeId).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		// Act
		err = r.Delete(nodeId, nil)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("NodeOnSite", func(t *testing.T) {
		mock.ExpectQuery(`^SELECT.*nodes.*`).
			WithArgs(nodeId, sqlmock.AnyArg()).
			WillReturnRows(row)

		// Act
		err = r.Delete(nodeId, nil)

		// Assert
		assert.Error(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("NodeErrorGrouped", func(t *testing.T) {
		mock.ExpectQuery(`^SELECT.*nodes.*`).
			WithArgs(nodeId, sqlmock.AnyArg()).
			WillReturnError(extsql.ErrNoRows)

		// mock.ExpectExec(regexp.QuoteMeta(`select * from attached_nodes where attached_id= $1 OR node_id= $2`)).
		// WithArgs(nodeId, nodeId).
		// WillReturnError(extsql.ErrNoRows)

		// Act
		err = r.Delete(nodeId, nil)

		// Assert
		assert.Error(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("NodeStillGrouped", func(t *testing.T) {
		mock.ExpectQuery(`^SELECT.*nodes.*`).
			WithArgs(nodeId, sqlmock.AnyArg()).
			WillReturnError(extsql.ErrNoRows)

		// mock.ExpectExec(regexp.QuoteMeta(`select * from attached_nodes where attached_id= $1 OR node_id= $2`)).
		// WithArgs(nodeId, nodeId).
		// WillReturnResult(sqlmock.NewResult(1, 1))

		// Act
		err = r.Delete(nodeId, nil)

		// Assert
		assert.Error(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}
