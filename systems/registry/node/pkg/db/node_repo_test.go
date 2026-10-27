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
	"slices"
	"sort"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/tj/assert"
	"github.com/ukama/ukama/systems/common/store"
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
			` ("id","name","type","created_at","updated_at","deleted_at")`+
			` VALUES ($1,$2,$3,$4,$5,$6) RETURNING "parent_node_id","latitude","longitude"`)).
			WithArgs(node.Id, node.Name, node.Type, sqlmock.AnyArg(),
				sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"parent_node_id", "latitude", "longitude"}).
				AddRow(nil, 0.0, 0.0))
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
	var (
		nodeId = ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)
		id     = nodeId.StringLowercase()
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

	lockQuery := `^SELECT "id","parent_node_id" FROM "nodes" WHERE id = \$1 .*FOR UPDATE`
	countQuery := `^SELECT count\(\*\) FROM "nodes" WHERE parent_node_id = \$1`

	expectLock := func(parent *string) {
		mock.ExpectQuery(lockQuery).WithArgs(id, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "parent_node_id"}).
				AddRow(id, parent))
	}

	expectChildren := func(count int) {
		mock.ExpectQuery(countQuery).WithArgs(id).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(count))
	}

	t.Run("NodeFound", func(t *testing.T) {
		mock.ExpectBegin()
		expectLock(nil)
		expectChildren(0)

		for _, table := range []string{"sites", "node_statuses", "nodes"} {
			mock.ExpectExec(`^UPDATE "`+table+`" SET "deleted_at"=`).
				WithArgs(sqlmock.AnyArg(), id).
				WillReturnResult(sqlmock.NewResult(1, 1))
		}

		mock.ExpectCommit()

		err := r.Delete(nodeId, nil)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("NodeNotFound", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(lockQuery).WithArgs(id, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}))
		mock.ExpectRollback()

		err := r.Delete(nodeId, nil)

		assert.True(t, errors.Is(err, store.ErrNotFound), err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("AttachedToParent", func(t *testing.T) {
		parent := "uk-sa2641-tnode-v0-aaaa"

		mock.ExpectBegin()
		expectLock(&parent)
		mock.ExpectRollback()

		err := r.Delete(nodeId, nil)

		assert.True(t, errors.Is(err, store.ErrPrecondition), err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("HasChildren", func(t *testing.T) {
		mock.ExpectBegin()
		expectLock(nil)
		expectChildren(1)
		mock.ExpectRollback()

		err := r.Delete(nodeId, nil)

		assert.True(t, errors.Is(err, store.ErrPrecondition), err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	dbErr := errors.New("db failed")

	for _, tc := range []struct {
		name   string
		expect func()
	}{
		{"LockFails", func() {
			mock.ExpectQuery(lockQuery).WithArgs(id, 1).WillReturnError(dbErr)
		}},
		{"CountChildrenFails", func() {
			expectLock(nil)
			mock.ExpectQuery(countQuery).WithArgs(id).WillReturnError(dbErr)
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			mock.ExpectBegin()
			tc.expect()
			mock.ExpectRollback()

			err := r.Delete(nodeId, nil)

			assert.True(t, errors.Is(err, dbErr), err)
			assert.False(t, errors.Is(err, store.ErrPrecondition))
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
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
		assert.True(t, errors.Is(err, store.ErrNotFound), err)

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

	pId := parentId.StringLowercase()
	aId := ampId.StringLowercase()

	expectLocks := func(ids ...string) {
		sorted := slices.Clone(ids)
		sort.Strings(sorted)

		for _, id := range sorted {
			mock.ExpectQuery(`^SELECT "id" FROM "nodes" WHERE id = \$1 .*FOR UPDATE`).
				WithArgs(id, 1).
				WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(id))
		}
	}

	expectParentSiteInTx := func(site uuid.UUID) {
		rows := sqlmock.NewRows([]string{"node_id", "site_id"})
		if site != uuid.Nil {
			rows.AddRow(pId, site)
		}

		mock.ExpectQuery(`^SELECT \* FROM "sites" WHERE node_id = `).
			WithArgs(pId, 1).
			WillReturnRows(rows)
	}

	expectCount := func(count int) {
		mock.ExpectQuery(`^SELECT count\(\*\) FROM "nodes" WHERE parent_node_id`).
			WithArgs(pId).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(count))
	}

	expectChildGet := func(id string, site uuid.UUID, parent *string) {
		mock.ExpectQuery(`^SELECT \* FROM "nodes" WHERE id = \$1`).
			WithArgs(id, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "type", "parent_node_id"}).
				AddRow(id, "amp", ukama.NODE_ID_TYPE_AMPNODE, parent))

		mock.ExpectQuery(`^SELECT \* FROM "sites" WHERE "sites"."node_id" = \$1`).
			WithArgs(id).
			WillReturnRows(sqlmock.NewRows([]string{"node_id", "site_id"}).AddRow(id, site))
	}

	updateQuery := `^UPDATE "nodes" SET "parent_node_id"=\$1,"updated_at"=\$2 WHERE id = \$3`

	t.Run("InvalidNodeCount", func(t *testing.T) {
		err := r.AttachNodes(parentId, []string{})

		assert.True(t, errors.Is(err, store.ErrInvalid), err)
	})

	t.Run("ParentNotFound", func(t *testing.T) {
		mock.ExpectQuery(`^SELECT.*nodes.*`).
			WithArgs(pId, sqlmock.AnyArg()).
			WillReturnError(extsql.ErrNoRows)

		err := r.AttachNodes(parentId, []string{aId})

		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ParentNotTowerNode", func(t *testing.T) {
		expectParentGet(ukama.NODE_ID_TYPE_HOMENODE, true)

		err := r.AttachNodes(parentId, []string{aId})

		assert.True(t, errors.Is(err, store.ErrInvalid), err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("AttachSuccess", func(t *testing.T) {
		expectParentGet(ukama.NODE_ID_TYPE_TOWERNODE, true)
		mock.ExpectBegin()
		expectLocks(pId, aId)
		expectParentSiteInTx(siteId)
		expectCount(0)
		expectChildGet(aId, siteId, nil)
		mock.ExpectExec(updateQuery).
			WithArgs(pId, sqlmock.AnyArg(), aId).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		err := r.AttachNodes(parentId, []string{aId})

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("LocksParentAndChildrenInSortedOrder", func(t *testing.T) {
		lo := ukama.NewVirtualAmplifierNodeId().StringLowercase()
		hi := ukama.NewVirtualAmplifierNodeId().StringLowercase()
		if hi < lo {
			lo, hi = hi, lo
		}

		lockIds := []string{pId, lo, hi}
		sort.Strings(lockIds)

		expectParentGet(ukama.NODE_ID_TYPE_TOWERNODE, true)
		mock.ExpectBegin()
		expectLocks(lockIds[:2]...)

		lockErr := errors.New("lock failed")
		mock.ExpectQuery(`^SELECT "id" FROM "nodes" WHERE id = \$1 .*FOR UPDATE`).
			WithArgs(lockIds[2], 1).
			WillReturnError(lockErr)
		mock.ExpectRollback()

		err := r.AttachNodes(parentId, []string{hi, lo})

		assert.True(t, errors.Is(err, lockErr), err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("AttachedNodeNotFound", func(t *testing.T) {
		expectParentGet(ukama.NODE_ID_TYPE_TOWERNODE, true)
		mock.ExpectBegin()
		mock.ExpectQuery(`^SELECT "id" FROM "nodes" WHERE id = \$1 .*FOR UPDATE`).
			WithArgs(aId, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}))
		mock.ExpectRollback()

		err := r.AttachNodes(parentId, []string{aId})

		assert.True(t, errors.Is(err, store.ErrNotFound), err)
		assert.Contains(t, err.Error(), aId)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	rejectCases := []struct {
		name   string
		expect func()
		target error
	}{
		{"ParentHasNoSite", func() {
			expectParentSiteInTx(uuid.Nil)
		}, store.ErrPrecondition},
		{"ParentSiteReadInsideTx", func() {
			expectParentSiteInTx(uuid.NewV4())
			expectCount(0)
			expectChildGet(aId, siteId, nil)
		}, store.ErrInvalid},
		{"ConcurrentAttachRejectedInsideTx", func() {
			expectParentSiteInTx(siteId)
			expectCount(2)
		}, store.ErrPrecondition},
		{"AlreadyAttached", func() {
			other := "uk-sa2641-tnode-v0-aaaa"
			expectParentSiteInTx(siteId)
			expectCount(0)
			expectChildGet(aId, siteId, &other)
		}, store.ErrPrecondition},
	}

	for _, tc := range rejectCases {
		t.Run(tc.name, func(t *testing.T) {
			expectParentGet(ukama.NODE_ID_TYPE_TOWERNODE, true)
			mock.ExpectBegin()
			expectLocks(pId, aId)
			tc.expect()
			mock.ExpectRollback()

			err := r.AttachNodes(parentId, []string{aId})

			assert.True(t, errors.Is(err, tc.target), err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}

	txErr := errors.New("tx query failed")

	failCases := []struct {
		name   string
		expect func()
	}{
		{"LockFails", func() {
			mock.ExpectQuery(`^SELECT "id" FROM "nodes" WHERE id = \$1 .*FOR UPDATE`).
				WithArgs(aId, 1).
				WillReturnError(txErr)
		}},
		{"ReadParentSiteFails", func() {
			expectLocks(pId, aId)
			mock.ExpectQuery(`^SELECT \* FROM "sites" WHERE node_id = `).
				WithArgs(pId, 1).
				WillReturnError(txErr)
		}},
		{"CountAttachedFails", func() {
			expectLocks(pId, aId)
			expectParentSiteInTx(siteId)
			mock.ExpectQuery(`^SELECT count\(\*\) FROM "nodes" WHERE parent_node_id`).
				WithArgs(pId).
				WillReturnError(txErr)
		}},
		{"UpdateFails", func() {
			expectLocks(pId, aId)
			expectParentSiteInTx(siteId)
			expectCount(0)
			expectChildGet(aId, siteId, nil)
			mock.ExpectExec(updateQuery).
				WithArgs(pId, sqlmock.AnyArg(), aId).
				WillReturnError(txErr)
		}},
	}

	for _, tc := range failCases {
		t.Run(tc.name, func(t *testing.T) {
			expectParentGet(ukama.NODE_ID_TYPE_TOWERNODE, true)
			mock.ExpectBegin()
			tc.expect()
			mock.ExpectRollback()

			err := r.AttachNodes(parentId, []string{aId})

			assert.True(t, errors.Is(err, txErr), err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
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
			WillReturnRows(sqlmock.NewRows([]string{"node_id", "site_id"}).
				AddRow(nodeId.StringLowercase(), uuid.NewV4()))

		mock.ExpectQuery(`^SELECT.*node_statuses.*`).
			WillReturnRows(sqlmock.NewRows([]string{"node_id"}))
	}

	t.Run("NodeNotFound", func(t *testing.T) {
		mock.ExpectBegin()

		mock.ExpectQuery(`^SELECT.*nodes.*`).
			WithArgs(nodeId.StringLowercase(), sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}))

		mock.ExpectRollback()

		// Act
		err := r.DetachNode(nodeId)

		// Assert
		assert.True(t, errors.Is(err, store.ErrNotFound), err)

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
		assert.True(t, errors.Is(err, store.ErrPrecondition), err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("DetachSuccess", func(t *testing.T) {
		mock.ExpectBegin()

		expectNodeGet(&parentId)

		mock.ExpectExec(`^UPDATE "nodes" SET "parent_node_id"=\$1,"updated_at"=\$2`+
			` WHERE id = \$3`).
			WithArgs(nil, sqlmock.AnyArg(), nodeId.StringLowercase()).
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
