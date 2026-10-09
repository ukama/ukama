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
	"github.com/ukama/ukama/systems/common/uuid"
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

func (u UkamaDbMock) ExecuteInTransaction(
	dbOperation func(tx *gorm.DB) *gorm.DB,
	nestedFuncs ...func() error) error {
	log.Fatal("implement me: ExecuteInTransaction()")
	return nil
}

func (u UkamaDbMock) ExecuteInTransaction2(
	dbOperation func(tx *gorm.DB) *gorm.DB,
	nestedFuncs ...func(tx *gorm.DB) error) error {
	log.Fatal("implement me: ExecuteInTransaction2()")
	return nil
}

func TestNodeRepo_Add(t *testing.T) {
	var (
		nodeId = ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)
		db     *extsql.DB
	)

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

		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "nodes"`+
			` ("id","name","type","parent_node_id","created_at","updated_at","deleted_`+
			`at") VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "latitude","longitude"`)).
			WithArgs(node.Id, node.Name, node.Type, node.ParentNodeId, sqlmock.AnyArg(),
				sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"latitude", "longitude"}).
				AddRow(0.0, 0.0))
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
	var (
		nodeId = ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)
		name   = "node-1"
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
	var (
		nodeId = ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)
		name   = "node-1"
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
	var (
		nodeId = ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)
		name   = "node-1"
		db     *extsql.DB
	)

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
		siteRow := sqlmock.NewRows([]string{"id", "name"}).
			AddRow(nodeId, name)

		mock.ExpectQuery(`^SELECT.*nodes.*`).
			WithArgs(nodeId, sqlmock.AnyArg()).
			WillReturnRows(siteRow)

		mock.ExpectQuery(`^SELECT.*parent_node_id.*`).
			WithArgs(nodeId).
			WillReturnRows(siteRow)

		mock.ExpectQuery(`^SELECT.*sites.*`).
			WithArgs(nodeId).
			WillReturnRows(siteRow)

		mock.ExpectQuery(`^SELECT.*node_statuses.*`).
			WithArgs(nodeId).
			WillReturnRows(siteRow)

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE`)).
			WithArgs(sqlmock.AnyArg(), nodeId).
			WillReturnResult(sqlmock.NewResult(1, 1))

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

	t.Run("NodeErrorGrouped", func(t *testing.T) {
		mock.ExpectQuery(`^SELECT.*nodes.*`).
			WithArgs(nodeId, sqlmock.AnyArg()).
			WillReturnError(extsql.ErrNoRows)

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

		// Act
		err = r.Delete(nodeId, nil)

		// Assert
		assert.Error(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestNodeRepo_List(t *testing.T) {
	var (
		nodeId       = ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)
		siteId       = uuid.NewV4()
		networkId    = uuid.NewV4()
		ntype        = ukama.NODE_ID_TYPE_HOMENODE
		connectivity = uint8(ukama.NodeConnectivityOnline)
		state        = uint8(ukama.NodeStateUnknown)
		db           *extsql.DB
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

	r := nodedb.NewNodeRepo(&UkamaDbMock{
		GormDb: gdb,
	})

	t.Run("ListWithAllFilters", func(t *testing.T) {
		// Arrange
		rows := sqlmock.NewRows([]string{
			"id", "name", "type", "connectivity", "state", "site_id", "network_id",
		}).
			AddRow(nodeId.String(), "node-1", ntype,
				ukama.NodeConnectivity(connectivity), ukama.NodeState(state), siteId,
				networkId)

		// Mock the main query first
		mock.ExpectQuery(`^SELECT nodes.* FROM "nodes" INNER JOIN node_statuses ON`+
			` nodes.id = node_statuses.node_id LEFT JOIN sites ON nodes.id =`+
			` sites.node_id AND sites.deleted_at IS NULL WHERE`+
			` node_statuses.deleted_at IS NULL AND nodes.id = \$1 AND sites.site_id =`+
			` \$2 AND sites.network_id = \$3 AND node_statuses.connectivity = \$4 AND`+
			` node_statuses.state = \$5 AND nodes.type = \$6 AND "nodes"."deleted_at"`+
			` IS NULL$`).
			WithArgs(nodeId.String(), siteId, networkId, connectivity, state, ntype).
			WillReturnRows(rows)

		// Mock the attached nodes query
		mock.ExpectQuery(`^SELECT \* FROM "nodes" WHERE "nodes"."parent_node_id" =` +
			` \$1 AND "nodes"."deleted_at" IS NULL$`).
			WithArgs(nodeId.String()).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "type"}))

		// Mock the sites query
		mock.ExpectQuery(`^SELECT \* FROM "sites" WHERE "sites"."node_id" = \$1 AND` +
			` "sites"."deleted_at" IS NULL$`).
			WithArgs(nodeId.String()).
			WillReturnRows(sqlmock.NewRows([]string{"node_id", "site_id", "network_id"}).
				AddRow(nodeId.String(), siteId, networkId))

		// Mock the node_statuses query
		mock.ExpectQuery(`^SELECT \* FROM "node_statuses" WHERE` +
			` "node_statuses"."node_id" = \$1 AND "node_statuses"."deleted_at" IS` +
			` NULL$`).
			WithArgs(nodeId.String()).
			WillReturnRows(sqlmock.NewRows([]string{"node_id", "connectivity", "state"}).
				AddRow(nodeId.String(), ukama.NodeConnectivity(connectivity),
					ukama.NodeState(state)))

		// Act
		nodes, err := r.List(nodeId.String(), siteId.String(), networkId.String(),
			ntype, &connectivity, &state)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, nodes)
		assert.Len(t, nodes, 1)
		assert.Equal(t, nodeId.String(), nodes[0].Id)
		assert.Equal(t, "node-1", nodes[0].Name)
		assert.Equal(t, ukama.NodeType(ntype), nodes[0].Type)
		assert.Equal(t, ukama.NodeConnectivity(connectivity),
			nodes[0].Status.Connectivity)
		assert.Equal(t, ukama.NodeState(state), nodes[0].Status.State)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	// Regression test: the network-scoped list must report the same connectivity as
	// GetNode. Status is populated by Preload, never by the joined columns, so a
	// stale value on the join row must not leak into the result.
	t.Run("ListReportsPreloadedConnectivityNotJoinedColumns", func(t *testing.T) {
		staleOnline := uint8(ukama.NodeConnectivityOnline)
		freshOffline := uint8(ukama.NodeConnectivityOffline)

		rows := sqlmock.NewRows([]string{
			"id", "name", "type", "connectivity", "state", "site_id", "network_id",
		}).
			AddRow(nodeId.String(), "node-1", ntype, staleOnline, state, siteId,
				networkId)

		mock.ExpectQuery(`^SELECT nodes.* FROM "nodes" INNER JOIN node_statuses ON` +
			` nodes.id = node_statuses.node_id LEFT JOIN sites ON nodes.id =` +
			` sites.node_id AND sites.deleted_at IS NULL WHERE` +
			` node_statuses.deleted_at IS NULL AND "nodes"."deleted_at" IS NULL$`).
			WithArgs().
			WillReturnRows(rows)

		mock.ExpectQuery(`^SELECT \* FROM "nodes" WHERE "nodes"."parent_node_id" =` +
			` \$1 AND "nodes"."deleted_at" IS NULL$`).
			WithArgs(nodeId.String()).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "type"}))

		mock.ExpectQuery(`^SELECT \* FROM "sites" WHERE "sites"."node_id" = \$1 AND` +
			` "sites"."deleted_at" IS NULL$`).
			WithArgs(nodeId.String()).
			WillReturnRows(sqlmock.NewRows([]string{"node_id", "site_id", "network_id"}).
				AddRow(nodeId.String(), siteId, networkId))

		mock.ExpectQuery(`^SELECT \* FROM "node_statuses" WHERE` +
			` "node_statuses"."node_id" = \$1 AND "node_statuses"."deleted_at" IS` +
			` NULL$`).
			WithArgs(nodeId.String()).
			WillReturnRows(sqlmock.NewRows([]string{"node_id", "connectivity", "state"}).
				AddRow(nodeId.String(), freshOffline, state))

		nodes, err := r.List("", "", "", "", nil, nil)

		assert.NoError(t, err)
		assert.Len(t, nodes, 1)
		assert.Equal(t, ukama.NodeConnectivityOffline, nodes[0].Status.Connectivity,
			"list must report the current status row, not a stale joined value")

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ListWithNoFilters", func(t *testing.T) {
		// Arrange
		rows := sqlmock.NewRows([]string{
			"id", "name", "type", "connectivity", "state", "site_id", "network_id",
		}).
			AddRow(nodeId.String(), "node-1", ntype, connectivity, state, siteId,
				networkId)

		// Mock the main query first
		mock.ExpectQuery(`^SELECT nodes.* FROM "nodes" INNER JOIN node_statuses ON` +
			` nodes.id = node_statuses.node_id LEFT JOIN sites ON nodes.id =` +
			` sites.node_id AND sites.deleted_at IS NULL WHERE` +
			` node_statuses.deleted_at IS NULL AND "nodes"."deleted_at" IS NULL$`).
			WithArgs().
			WillReturnRows(rows)

		// Mock the attached nodes query
		mock.ExpectQuery(`^SELECT \* FROM "nodes" WHERE "nodes"."parent_node_id" =` +
			` \$1 AND "nodes"."deleted_at" IS NULL$`).
			WithArgs(nodeId.String()).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "type"}))

		// Mock the sites query
		mock.ExpectQuery(`^SELECT \* FROM "sites" WHERE "sites"."node_id" = \$1 AND` +
			` "sites"."deleted_at" IS NULL$`).
			WithArgs(nodeId.String()).
			WillReturnRows(sqlmock.NewRows([]string{"node_id", "site_id", "network_id"}).
				AddRow(nodeId.String(), siteId, networkId))

		// Mock the node_statuses query
		mock.ExpectQuery(`^SELECT \* FROM "node_statuses" WHERE` +
			` "node_statuses"."node_id" = \$1 AND "node_statuses"."deleted_at" IS` +
			` NULL$`).
			WithArgs(nodeId.String()).
			WillReturnRows(sqlmock.NewRows([]string{"node_id", "connectivity", "state"}).
				AddRow(nodeId.String(), connectivity, state))

		// Act
		nodes, err := r.List("", "", "", "", nil, nil)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, nodes)
		assert.Len(t, nodes, 1)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	// Unlike siteRepo.GetByNetwork, List must not turn zero rows into
	// gorm.ErrRecordNotFound — an empty network is not a missing one.
	t.Run("ListEmptyNetworkReturnsNoError", func(t *testing.T) {
		// Arrange
		mock.ExpectQuery(`^SELECT nodes.* FROM "nodes" INNER JOIN node_statuses ON` +
			` nodes.id = node_statuses.node_id LEFT JOIN sites ON nodes.id =` +
			` sites.node_id AND sites.deleted_at IS NULL WHERE` +
			` node_statuses.deleted_at IS NULL AND sites.network_id = \$1 AND` +
			` "nodes"."deleted_at" IS NULL$`).
			WithArgs(networkId).
			WillReturnRows(sqlmock.NewRows([]string{
				"id", "name", "type", "connectivity", "state", "site_id", "network_id",
			}))

		// Act
		nodes, err := r.List("", "", networkId.String(), "", nil, nil)

		// Assert
		assert.NoError(t, err)
		assert.Empty(t, nodes)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("ListExcludesSoftDeletedSiteRows", func(t *testing.T) {
		// Arrange
		mock.ExpectQuery(`LEFT JOIN sites ON nodes\.id = sites\.node_id AND` +
			` sites\.deleted_at IS NULL`).
			WithArgs(siteId).
			WillReturnRows(sqlmock.NewRows([]string{
				"id", "name", "type", "connectivity", "state", "site_id", "network_id",
			}))

		// Act
		nodes, err := r.List("", siteId.String(), "", "", nil, nil)

		// Assert
		assert.NoError(t, err)
		assert.Empty(t, nodes)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

}

func TestNodeRepo_GetNodesByState(t *testing.T) {
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

	t.Run("NoNodesInState", func(t *testing.T) {
		mock.ExpectQuery(`^SELECT.*nodes.*`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name"}))

		// Act
		nodes, err := r.GetNodesByState(uint8(ukama.NodeConnectivityOnline),
			uint8(ukama.NodeStateUnknown))

		// Assert
		assert.NoError(t, err)
		assert.Empty(t, nodes)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("QueryError", func(t *testing.T) {
		mock.ExpectQuery(`^SELECT.*nodes.*`).
			WillReturnError(errors.New("internal"))

		// Act
		nodes, err := r.GetNodesByState(uint8(ukama.NodeConnectivityOnline),
			uint8(ukama.NodeStateUnknown))

		// Assert
		assert.Error(t, err)
		assert.Nil(t, nodes)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestNodeRepo_Update(t *testing.T) {
	var (
		nodeId = ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)
		db     *extsql.DB
	)

	node := &nodedb.Node{
		Id:   nodeId.StringLowercase(),
		Name: "node-1",
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

	t.Run("UpdateSuccess", func(t *testing.T) {
		mock.ExpectBegin()

		mock.ExpectQuery(`^UPDATE "nodes" SET`).
			WillReturnRows(sqlmock.NewRows([]string{"latitude", "longitude"}).
				AddRow("0", "0"))

		mock.ExpectCommit()

		// Act
		err = r.Update(node, nil)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	// Updates reports zero rows affected when the node id doesn't exist:
	// the repo must surface that as gorm.ErrRecordNotFound, not silently succeed.
	t.Run("UpdateNotFound", func(t *testing.T) {
		mock.ExpectBegin()

		mock.ExpectQuery(`^UPDATE "nodes" SET`).
			WillReturnRows(sqlmock.NewRows([]string{"latitude", "longitude"}))

		mock.ExpectRollback()

		// Act
		err = r.Update(node, nil)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("UpdateQueryError", func(t *testing.T) {
		mock.ExpectBegin()

		mock.ExpectQuery(`^UPDATE "nodes" SET`).
			WillReturnError(errors.New("internal"))

		mock.ExpectRollback()

		// Act
		err = r.Update(node, nil)

		// Assert
		assert.Error(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("UpdateNestedFuncError", func(t *testing.T) {
		mock.ExpectBegin()

		mock.ExpectQuery(`^UPDATE "nodes" SET`).
			WillReturnRows(sqlmock.NewRows([]string{"latitude", "longitude"}).
				AddRow("0", "0"))

		mock.ExpectRollback()

		nestedErr := errors.New("nested failure")

		// Act
		err = r.Update(node, func(*nodedb.Node, *gorm.DB) error {
			return nestedErr
		})

		// Assert
		assert.Error(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestNodeRepo_AttachNodes(t *testing.T) {
	var (
		parentId = ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_TOWERNODE)
		ampId    = ukama.NewVirtualAmplifierNodeId()
		siteId   = uuid.NewV4()
		db       *extsql.DB
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

	r := nodedb.NewNodeRepo(&UkamaDbMock{
		GormDb: gdb,
	})

	expectParentGet := func(nodeType ukama.NodeType, withSite bool) {
		parentRows := sqlmock.NewRows([]string{"id", "name", "type"}).
			AddRow(parentId.StringLowercase(), "parent", nodeType)

		mock.ExpectQuery(`^SELECT.*nodes.*`).
			WithArgs(parentId.StringLowercase(), sqlmock.AnyArg()).
			WillReturnRows(parentRows)

		mock.ExpectQuery(`^SELECT.*parent_node_id.*`).
			WithArgs(parentId.StringLowercase()).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name"}))

		siteRows := sqlmock.NewRows([]string{"node_id", "site_id"})
		if withSite {
			siteRows.AddRow(parentId.StringLowercase(), siteId)
		}

		mock.ExpectQuery(`^SELECT.*sites.*`).
			WithArgs(parentId.StringLowercase()).
			WillReturnRows(siteRows)

		mock.ExpectQuery(`^SELECT.*node_statuses.*`).
			WithArgs(parentId.StringLowercase()).
			WillReturnRows(sqlmock.NewRows([]string{"node_id"}))
	}

	t.Run("InvalidNodeCount", func(t *testing.T) {
		// Act
		err := r.AttachNodes(parentId, []string{})

		// Assert
		assert.Error(t, err)
	})

	t.Run("ParentNotFound", func(t *testing.T) {
		mock.ExpectQuery(`^SELECT.*nodes.*`).
			WithArgs(parentId.StringLowercase(), sqlmock.AnyArg()).
			WillReturnError(extsql.ErrNoRows)

		// Act
		err := r.AttachNodes(parentId, []string{ampId.StringLowercase()})

		// Assert
		assert.Error(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("ParentNotTowerNode", func(t *testing.T) {
		expectParentGet(ukama.NODE_ID_TYPE_HOMENODE, true)

		// Act
		err := r.AttachNodes(parentId, []string{ampId.StringLowercase()})

		// Assert
		assert.Error(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("ParentHasNoSite", func(t *testing.T) {
		expectParentGet(ukama.NODE_ID_TYPE_TOWERNODE, false)

		// Act
		err := r.AttachNodes(parentId, []string{ampId.StringLowercase()})

		// Assert
		assert.Error(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("BatchGetFails", func(t *testing.T) {
		expectParentGet(ukama.NODE_ID_TYPE_TOWERNODE, true)

		mock.ExpectQuery(`^SELECT.*nodes.*`).
			WithArgs(ampId.StringLowercase()).
			WillReturnError(errors.New("internal"))

		// Act
		err := r.AttachNodes(parentId, []string{ampId.StringLowercase()})

		// Assert
		assert.Error(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("AttachedNodeNotFound", func(t *testing.T) {
		expectParentGet(ukama.NODE_ID_TYPE_TOWERNODE, true)

		mock.ExpectQuery(`^SELECT.*nodes.*`).
			WithArgs(ampId.StringLowercase()).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name"}))

		// Act
		err := r.AttachNodes(parentId, []string{ampId.StringLowercase()})

		// Assert
		assert.Error(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("AttachSuccess", func(t *testing.T) {
		expectParentGet(ukama.NODE_ID_TYPE_TOWERNODE, true)

		ampRows := sqlmock.NewRows([]string{"id", "name", "type", "site_id"}).
			AddRow(ampId.StringLowercase(), "amp", ukama.NODE_ID_TYPE_AMPNODE, siteId)

		mock.ExpectQuery(`^SELECT.*nodes.*`).
			WithArgs(ampId.StringLowercase()).
			WillReturnRows(ampRows)

		mock.ExpectQuery(`^SELECT.*parent_node_id.*`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name"}))

		mock.ExpectQuery(`^SELECT.*sites.*`).
			WillReturnRows(sqlmock.NewRows([]string{"node_id", "site_id"}).
				AddRow(ampId.StringLowercase(), siteId))

		mock.ExpectQuery(`^SELECT.*node_statuses.*`).
			WillReturnRows(sqlmock.NewRows([]string{"node_id"}))

		mock.ExpectBegin()

		// Lock the parent row, then re-count its currently-attached nodes
		// inside the transaction (the TOCTOU-safe re-check).
		mock.ExpectQuery(`^SELECT.*nodes.*FOR UPDATE`).
			WithArgs(parentId.StringLowercase(), sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(parentId.StringLowercase()))

		mock.ExpectQuery(`^SELECT count\(\*\) FROM "nodes" WHERE parent_node_id`).
			WithArgs(parentId.StringLowercase()).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		// Re-fetch and lock the attached node fresh inside the transaction.
		ampRowsLocked := sqlmock.NewRows([]string{"id", "name", "type", "site_id"}).
			AddRow(ampId.StringLowercase(), "amp", ukama.NODE_ID_TYPE_AMPNODE, siteId)

		mock.ExpectQuery(`^SELECT.*nodes.*FOR UPDATE`).
			WithArgs(ampId.StringLowercase(), sqlmock.AnyArg()).
			WillReturnRows(ampRowsLocked)

		mock.ExpectQuery(`^SELECT.*sites.*`).
			WithArgs(ampId.StringLowercase()).
			WillReturnRows(sqlmock.NewRows([]string{"node_id", "site_id"}).
				AddRow(ampId.StringLowercase(), siteId))

		mock.ExpectExec(`^UPDATE "nodes" SET`).
			WillReturnResult(sqlmock.NewResult(0, 1))

		// Save() cascades to the Site association, which has its own
		// BeforeSave upsert (ON CONFLICT DO NOTHING).
		mock.ExpectExec(`^INSERT INTO "sites"`).
			WillReturnResult(sqlmock.NewResult(0, 0))

		mock.ExpectCommit()

		// Act
		err := r.AttachNodes(parentId, []string{ampId.StringLowercase()})

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	// Regression test for the TOCTOU race
	t.Run("ConcurrentAttachRejectedInsideTx", func(t *testing.T) {
		expectParentGet(ukama.NODE_ID_TYPE_TOWERNODE, true)

		ampRows := sqlmock.NewRows([]string{"id", "name", "type", "site_id"}).
			AddRow(ampId.StringLowercase(), "amp", ukama.NODE_ID_TYPE_AMPNODE, siteId)

		mock.ExpectQuery(`^SELECT.*nodes.*`).
			WithArgs(ampId.StringLowercase()).
			WillReturnRows(ampRows)

		mock.ExpectQuery(`^SELECT.*parent_node_id.*`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name"}))

		mock.ExpectQuery(`^SELECT.*sites.*`).
			WillReturnRows(sqlmock.NewRows([]string{"node_id", "site_id"}).
				AddRow(ampId.StringLowercase(), siteId))

		mock.ExpectQuery(`^SELECT.*node_statuses.*`).
			WillReturnRows(sqlmock.NewRows([]string{"node_id"}))

		mock.ExpectBegin()

		mock.ExpectQuery(`^SELECT.*nodes.*FOR UPDATE`).
			WithArgs(parentId.StringLowercase(), sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(parentId.StringLowercase()))

		// Another concurrent AttachNodes call already filled both slots
		// between our pre-transaction Get and this lock.
		mock.ExpectQuery(`^SELECT count\(\*\) FROM "nodes" WHERE parent_node_id`).
			WithArgs(parentId.StringLowercase()).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

		mock.ExpectRollback()

		// Act
		err := r.AttachNodes(parentId, []string{ampId.StringLowercase()})

		// Assert
		assert.Error(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestNodeRepo_DetachNode(t *testing.T) {
	var (
		nodeId   = ukama.NewVirtualAmplifierNodeId()
		parentId = "uk-sa2641-tnode-v0-aaaa"
		db       *extsql.DB
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

	r := nodedb.NewNodeRepo(&UkamaDbMock{
		GormDb: gdb,
	})

	expectNodeGet := func(parentNodeId *string) {
		rows := sqlmock.NewRows([]string{"id", "name", "parent_node_id"}).
			AddRow(nodeId.StringLowercase(), "amp", parentNodeId)

		mock.ExpectQuery(`^SELECT.*nodes.*`).
			WithArgs(nodeId.StringLowercase(), sqlmock.AnyArg()).
			WillReturnRows(rows)

		mock.ExpectQuery(`^SELECT.*parent_node_id.*`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name"}))

		mock.ExpectQuery(`^SELECT.*sites.*`).
			WillReturnRows(sqlmock.NewRows([]string{"node_id"}))

		mock.ExpectQuery(`^SELECT.*node_statuses.*`).
			WillReturnRows(sqlmock.NewRows([]string{"node_id"}))
	}

	t.Run("NodeNotFound", func(t *testing.T) {
		mock.ExpectBegin()

		mock.ExpectQuery(`^SELECT.*nodes.*`).
			WithArgs(nodeId.StringLowercase(), sqlmock.AnyArg()).
			WillReturnError(extsql.ErrNoRows)

		mock.ExpectRollback()

		// Act
		err := r.DetachNode(nodeId)

		// Assert
		assert.Error(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("NotAttachedToParent", func(t *testing.T) {
		mock.ExpectBegin()

		expectNodeGet(nil)

		mock.ExpectRollback()

		// Act
		err := r.DetachNode(nodeId)

		// Assert
		assert.Error(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("DetachSuccess", func(t *testing.T) {
		mock.ExpectBegin()

		expectNodeGet(&parentId)

		mock.ExpectExec(`^UPDATE "nodes" SET`).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectCommit()

		// Act
		err := r.DetachNode(nodeId)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("SaveFails", func(t *testing.T) {
		mock.ExpectBegin()

		expectNodeGet(&parentId)

		mock.ExpectExec(`^UPDATE "nodes" SET`).
			WillReturnError(errors.New("internal"))

		mock.ExpectRollback()

		// Act
		err := r.DetachNode(nodeId)

		// Assert
		assert.Error(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestNodeRepo_GetNodeCount(t *testing.T) {
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

	t.Run("CountSuccess", func(t *testing.T) {
		mock.ExpectQuery(`^SELECT count\(\*\) FROM "nodes"`).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(3))

		mock.ExpectQuery(`^SELECT count\(\*\) FROM "nodes" JOIN` +
			` node_statuses.*connectivity = \$1`).
			WithArgs(ukama.NodeConnectivityOnline).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

		mock.ExpectQuery(`^SELECT count\(\*\) FROM "nodes" JOIN` +
			` node_statuses.*connectivity = \$1`).
			WithArgs(ukama.NodeConnectivityOffline).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		// Act
		total, online, offline, err := r.GetNodeCount()

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, int64(3), total)
		assert.Equal(t, int64(2), online)
		assert.Equal(t, int64(1), offline)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("TotalCountError", func(t *testing.T) {
		mock.ExpectQuery(`^SELECT count\(\*\) FROM "nodes"`).
			WillReturnError(errors.New("internal"))

		// Act
		total, online, offline, err := r.GetNodeCount()

		// Assert
		assert.Error(t, err)
		assert.Equal(t, int64(0), total)
		assert.Equal(t, int64(0), online)
		assert.Equal(t, int64(0), offline)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("OnlineCountError", func(t *testing.T) {
		mock.ExpectQuery(`^SELECT count\(\*\) FROM "nodes"`).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(3))

		mock.ExpectQuery(`^SELECT count\(\*\) FROM "nodes" JOIN` +
			` node_statuses.*connectivity = \$1`).
			WithArgs(ukama.NodeConnectivityOnline).
			WillReturnError(errors.New("internal"))

		// Act
		total, online, offline, err := r.GetNodeCount()

		// Assert
		assert.Error(t, err)
		assert.Equal(t, int64(0), total)
		assert.Equal(t, int64(0), online)
		assert.Equal(t, int64(0), offline)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("OfflineCountError", func(t *testing.T) {
		mock.ExpectQuery(`^SELECT count\(\*\) FROM "nodes"`).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(3))

		mock.ExpectQuery(`^SELECT count\(\*\) FROM "nodes" JOIN` +
			` node_statuses.*connectivity = \$1`).
			WithArgs(ukama.NodeConnectivityOnline).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

		mock.ExpectQuery(`^SELECT count\(\*\) FROM "nodes" JOIN` +
			` node_statuses.*connectivity = \$1`).
			WithArgs(ukama.NodeConnectivityOffline).
			WillReturnError(errors.New("internal"))

		// Act
		total, online, offline, err := r.GetNodeCount()

		// Assert
		assert.Error(t, err)
		assert.Equal(t, int64(0), total)
		assert.Equal(t, int64(0), online)
		assert.Equal(t, int64(0), offline)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}
