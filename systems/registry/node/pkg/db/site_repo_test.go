/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package db_test

import (
	extsql "database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/tj/assert"
	"github.com/ukama/ukama/systems/common/store"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	sidedb "github.com/ukama/ukama/systems/registry/node/pkg/db"
)

func TestSiteRepo_GetNodes(t *testing.T) {
	// Arrange
	var (
		siteID  = uuid.NewV4()
		nodeIDa = ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)
		nodeIDb = ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)
		db      *extsql.DB
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

	r := sidedb.NewSiteRepo(&UkamaDbMock{
		GormDb: gdb,
	})

	t.Run("SiteFound", func(t *testing.T) {
		siteRows := sqlmock.NewRows([]string{"node_id", "site_id"}).
			AddRow(nodeIDa, siteID).
			AddRow(nodeIDb, siteID)

		mock.ExpectQuery(`^SELECT.*sites.*`).
			WithArgs(siteID).
			WillReturnRows(siteRows)

		// Act
		nodes, err := r.GetNodes(siteID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, nodes)
		assert.Equal(t, 2, len(nodes))

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("SiteNotFound", func(t *testing.T) {
		// Arrange

		mock.ExpectQuery(`^SELECT.*sites.*`).
			WithArgs(siteID).
			WillReturnRows(sqlmock.NewRows([]string{"node_id", "site_id"}))

		// Act
		nodes, err := r.GetNodes(siteID)

		// Assert
		assert.True(t, errors.Is(err, store.ErrNotFound), err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.Nil(t, nodes)
	})
}

func TestSiteRepo_AssignNodes(t *testing.T) {
	siteId := uuid.NewV4()
	networkId := uuid.NewV4()
	otherSiteId := uuid.NewV4()

	// Node IDs sort as anode < tnode, which is the lock order the repo must use.
	tNode := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_TOWERNODE).StringLowercase()
	aNode := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_AMPNODE).StringLowercase()

	lockQuery := regexp.QuoteMeta(`SELECT "id" FROM "nodes" WHERE id = $1 AND "nodes"."deleted_at" IS NULL ORDER BY "nodes"."id" LIMIT $2 FOR UPDATE`)
	siteQuery := regexp.QuoteMeta(`SELECT * FROM "sites" WHERE lower(node_id) = $1 AND "sites"."deleted_at" IS NULL LIMIT $2`)
	insertQuery := regexp.QuoteMeta(`INSERT INTO "sites"`)
	siteCols := []string{"node_id", "site_id", "network_id"}

	newRepo := func(t *testing.T) (sidedb.SiteRepo, sqlmock.Sqlmock) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		gdb, err := gorm.Open(postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		}), &gorm.Config{})
		assert.NoError(t, err)

		return sidedb.NewSiteRepo(&UkamaDbMock{GormDb: gdb}), mock
	}

	t.Run("AssignsInLockOrder", func(t *testing.T) {
		r, mock := newRepo(t)

		mock.ExpectBegin()
		for _, id := range []string{aNode, tNode} {
			mock.ExpectQuery(lockQuery).WithArgs(id, 1).
				WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(id))
			mock.ExpectQuery(siteQuery).WithArgs(id, 1).
				WillReturnRows(sqlmock.NewRows(siteCols))
			mock.ExpectExec(insertQuery).
				WithArgs(id, siteId, networkId, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
				WillReturnResult(sqlmock.NewResult(0, 1))
		}
		mock.ExpectCommit()

		assigned, err := r.AssignNodes(siteId, networkId, []string{tNode, aNode})

		assert.NoError(t, err)
		assert.Len(t, assigned, 2)
		assert.Equal(t, aNode, assigned[0].NodeId)
		assert.Equal(t, tNode, assigned[1].NodeId)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("SameSiteIsSkipped", func(t *testing.T) {
		r, mock := newRepo(t)

		mock.ExpectBegin()
		mock.ExpectQuery(lockQuery).WithArgs(tNode, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(tNode))
		mock.ExpectQuery(siteQuery).WithArgs(tNode, 1).
			WillReturnRows(sqlmock.NewRows(siteCols).AddRow(tNode, siteId, networkId))
		mock.ExpectCommit()

		assigned, err := r.AssignNodes(siteId, networkId, []string{tNode})

		assert.NoError(t, err)
		assert.Empty(t, assigned)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("OtherSiteFailsWholeCall", func(t *testing.T) {
		r, mock := newRepo(t)

		mock.ExpectBegin()
		mock.ExpectQuery(lockQuery).WithArgs(aNode, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(aNode))
		mock.ExpectQuery(siteQuery).WithArgs(aNode, 1).
			WillReturnRows(sqlmock.NewRows(siteCols))
		mock.ExpectExec(insertQuery).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectQuery(lockQuery).WithArgs(tNode, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(tNode))
		mock.ExpectQuery(siteQuery).WithArgs(tNode, 1).
			WillReturnRows(sqlmock.NewRows(siteCols).AddRow(tNode, otherSiteId, networkId))
		mock.ExpectRollback()

		assigned, err := r.AssignNodes(siteId, networkId, []string{tNode, aNode})

		assert.True(t, errors.Is(err, store.ErrAlreadyExists), err)
		assert.Nil(t, assigned)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("InsertConflictIsAlreadyExists", func(t *testing.T) {
		r, mock := newRepo(t)

		mock.ExpectBegin()
		mock.ExpectQuery(lockQuery).WithArgs(tNode, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(tNode))
		mock.ExpectQuery(siteQuery).WithArgs(tNode, 1).
			WillReturnRows(sqlmock.NewRows(siteCols))
		mock.ExpectExec(insertQuery).
			WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectRollback()

		assigned, err := r.AssignNodes(siteId, networkId, []string{tNode})

		assert.True(t, errors.Is(err, store.ErrAlreadyExists), err)
		assert.Contains(t, err.Error(), tNode)
		assert.Nil(t, assigned)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("MissingNodeIsNotFound", func(t *testing.T) {
		r, mock := newRepo(t)

		mock.ExpectBegin()
		mock.ExpectQuery(lockQuery).WithArgs(tNode, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}))
		mock.ExpectRollback()

		assigned, err := r.AssignNodes(siteId, networkId, []string{tNode})

		assert.True(t, errors.Is(err, store.ErrNotFound), err)
		assert.Contains(t, err.Error(), tNode)
		assert.Nil(t, assigned)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSiteRepo_GetByNetwork(t *testing.T) {
	// Arrange
	var (
		networkID = uuid.NewV4()
		nodeIDa   = ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)
		db        *extsql.DB
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

	r := sidedb.NewSiteRepo(&UkamaDbMock{
		GormDb: gdb,
	})

	t.Run("NetworkFound", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"node_id", "network_id"}).
			AddRow(nodeIDa, networkID)

		mock.ExpectQuery(`^SELECT.*sites.*`).
			WithArgs(networkID).
			WillReturnRows(rows)

		// Act
		nodes, err := r.GetByNetwork(networkID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, nodes)
		assert.Equal(t, 1, len(nodes))

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("NetworkNotFound", func(t *testing.T) {
		mock.ExpectQuery(`^SELECT.*sites.*`).
			WithArgs(networkID).
			WillReturnRows(sqlmock.NewRows([]string{"node_id", "site_id"}))

		// Act
		nodes, err := r.GetByNetwork(networkID)

		// Assert
		assert.True(t, errors.Is(err, store.ErrNotFound), err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.Nil(t, nodes)
	})
}

func TestSiteRepo_GetFreeNodes(t *testing.T) {
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

	r := sidedb.NewSiteRepo(&UkamaDbMock{
		GormDb: gdb,
	})

	t.Run("FreeNodesFound", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name"}).
			AddRow(nodeId.StringLowercase(), "node-1")

		mock.ExpectQuery(`^SELECT.*nodes.*`).
			WillReturnRows(rows)

		mock.ExpectQuery(`^SELECT.*parent_node_id.*`).
			WithArgs(nodeId.StringLowercase()).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name"}))

		mock.ExpectQuery(`^SELECT.*sites.*`).
			WithArgs(nodeId.StringLowercase()).
			WillReturnRows(sqlmock.NewRows([]string{"node_id"}))

		mock.ExpectQuery(`^SELECT.*node_statuses.*`).
			WithArgs(nodeId.StringLowercase()).
			WillReturnRows(sqlmock.NewRows([]string{"node_id"}))

		// Act
		nodes, err := r.GetFreeNodes()

		// Assert
		assert.NoError(t, err)
		assert.Len(t, nodes, 1)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("QueryError", func(t *testing.T) {
		mock.ExpectQuery(`^SELECT.*nodes.*`).
			WillReturnError(errors.New("internal"))

		// Act
		nodes, err := r.GetFreeNodes()

		// Assert
		assert.Error(t, err)
		assert.Nil(t, nodes)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestSiteRepo_IsAllocated(t *testing.T) {
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

	r := sidedb.NewSiteRepo(&UkamaDbMock{
		GormDb: gdb,
	})

	t.Run("Allocated", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"node_id"}).
			AddRow(nodeId.StringLowercase())

		mock.ExpectQuery(`^SELECT.*sites.*`).
			WillReturnRows(rows)

		// Act
		ok, site := r.IsAllocated(nodeId)

		// Assert
		assert.True(t, ok)
		assert.NotNil(t, site)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("NotAllocated", func(t *testing.T) {
		mock.ExpectQuery(`^SELECT.*sites.*`).
			WillReturnError(extsql.ErrNoRows)

		// Act
		ok, site := r.IsAllocated(nodeId)

		// Assert
		assert.False(t, ok)
		assert.NotNil(t, site)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestSiteRepo_RemoveNode(t *testing.T) {
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

	r := sidedb.NewSiteRepo(&UkamaDbMock{
		GormDb: gdb,
	})

	t.Run("NotAllocated", func(t *testing.T) {
		mock.ExpectQuery(`^SELECT.*sites.*`).
			WillReturnError(extsql.ErrNoRows)

		// Act
		site, err := r.RemoveNode(nodeId)

		// Assert
		assert.True(t, errors.Is(err, store.ErrPrecondition), err)
		assert.Nil(t, site)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("StillGrouped", func(t *testing.T) {
		allocatedRows := sqlmock.NewRows([]string{"node_id"}).
			AddRow(nodeId.StringLowercase())

		mock.ExpectQuery(`^SELECT.*sites.*`).
			WillReturnRows(allocatedRows)

		mock.ExpectExec(regexp.QuoteMeta(`select * from nodes where parent_node_id=`+
			` $1  OR (id= $2 AND parent_node_id is NOT NULL)`)).
			WithArgs(nodeId.StringLowercase(), nodeId.StringLowercase()).
			WillReturnResult(sqlmock.NewResult(0, 1))

		// Act
		site, err := r.RemoveNode(nodeId)

		// Assert
		assert.True(t, errors.Is(err, store.ErrPrecondition), err)
		assert.Nil(t, site)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("RemoveSuccess", func(t *testing.T) {
		allocatedRows := sqlmock.NewRows([]string{"node_id"}).
			AddRow(nodeId.StringLowercase())

		mock.ExpectQuery(`^SELECT.*sites.*`).
			WillReturnRows(allocatedRows)

		mock.ExpectExec(regexp.QuoteMeta(`select * from nodes where parent_node_id=`+
			` $1  OR (id= $2 AND parent_node_id is NOT NULL)`)).
			WithArgs(nodeId.StringLowercase(), nodeId.StringLowercase()).
			WillReturnResult(sqlmock.NewResult(0, 0))

		mock.ExpectBegin()

		mock.ExpectExec(`^UPDATE "sites" SET "deleted_at"=`).
			WithArgs(sqlmock.AnyArg(), nodeId.StringLowercase()).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectCommit()

		// Act
		site, err := r.RemoveNode(nodeId)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, site)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestSiteRepo_ReleaseSite(t *testing.T) {
	siteId := uuid.NewV4()
	networkId := uuid.NewV4()
	tNode := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_TOWERNODE).StringLowercase()
	aNode := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_AMPNODE).StringLowercase()

	selectQuery := regexp.QuoteMeta(`SELECT * FROM "sites" WHERE site_id = $1 AND "sites"."deleted_at" IS NULL FOR UPDATE`)

	newRepo := func(t *testing.T) (sidedb.SiteRepo, sqlmock.Sqlmock) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		gdb, err := gorm.Open(postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		}), &gorm.Config{})
		assert.NoError(t, err)

		return sidedb.NewSiteRepo(&UkamaDbMock{GormDb: gdb}), mock
	}

	t.Run("ReleasesAllNodes", func(t *testing.T) {
		r, mock := newRepo(t)

		mock.ExpectBegin()
		mock.ExpectQuery(selectQuery).
			WithArgs(siteId).
			WillReturnRows(sqlmock.NewRows([]string{"node_id", "site_id", "network_id"}).
				AddRow(tNode, siteId, networkId).
				AddRow(aNode, siteId, networkId))
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "sites" SET "deleted_at"=$1 WHERE site_id = $2 AND "sites"."deleted_at" IS NULL`)).
			WithArgs(sqlmock.AnyArg(), siteId).
			WillReturnResult(sqlmock.NewResult(0, 2))
		mock.ExpectCommit()

		released, err := r.ReleaseSite(siteId)

		assert.NoError(t, err)
		assert.Len(t, released, 2)
		assert.Equal(t, tNode, released[0].NodeId)
		assert.Equal(t, aNode, released[1].NodeId)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("NoAssignedNodesIsNoop", func(t *testing.T) {
		r, mock := newRepo(t)

		mock.ExpectBegin()
		mock.ExpectQuery(selectQuery).
			WithArgs(siteId).
			WillReturnRows(sqlmock.NewRows([]string{"node_id", "site_id", "network_id"}))
		mock.ExpectCommit()

		released, err := r.ReleaseSite(siteId)

		assert.NoError(t, err)
		assert.Empty(t, released)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("SelectError", func(t *testing.T) {
		r, mock := newRepo(t)

		mock.ExpectBegin()
		mock.ExpectQuery(selectQuery).
			WithArgs(siteId).
			WillReturnError(errors.New("internal"))
		mock.ExpectRollback()

		released, err := r.ReleaseSite(siteId)

		assert.Error(t, err)
		assert.Nil(t, released)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
