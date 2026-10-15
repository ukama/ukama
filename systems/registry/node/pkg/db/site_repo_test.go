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
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/tj/assert"
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
			WillReturnError(extsql.ErrNoRows)

		// Act
		nodes, err := r.GetNodes(siteID)

		// Assert
		assert.Error(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.Nil(t, nodes)
	})
}

func TestSiteRepo_AddNode(t *testing.T) {
	var (
		nodeId    = ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)
		siteID    = uuid.NewV4()
		networkID = uuid.NewV4()
	)

	site := &sidedb.Site{
		NodeId:    nodeId.StringLowercase(),
		SiteId:    siteID,
		NetworkId: networkID,
	}

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

	r := sidedb.NewSiteRepo(&UkamaDbMock{
		GormDb: gdb,
	})

	t.Run("AddNodeSuccess", func(t *testing.T) {
		mock.ExpectBegin()

		mock.ExpectExec(`^INSERT INTO "sites"`).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectCommit()

		// Act
		err = r.AddNode(site, nil)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	// Insert conflicts (ON CONFLICT DO NOTHING) report zero rows affected:
	// the repo must treat that as the node already belonging to a site.
	t.Run("NodeAlreadyBelongsToSite", func(t *testing.T) {
		mock.ExpectBegin()

		mock.ExpectExec(`^INSERT INTO "sites"`).
			WillReturnResult(sqlmock.NewResult(0, 0))

		mock.ExpectRollback()

		// Act
		err = r.AddNode(site, nil)

		// Assert
		assert.Error(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("NestedFuncError", func(t *testing.T) {
		mock.ExpectBegin()

		mock.ExpectExec(`^INSERT INTO "sites"`).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectRollback()

		nestedErr := errors.New("nested failure")

		// Act
		err = r.AddNode(site, func(*sidedb.Site, *gorm.DB) error {
			return nestedErr
		})

		// Assert
		assert.Error(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
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
			WillReturnError(extsql.ErrNoRows)

		// Act
		nodes, err := r.GetByNetwork(networkID)

		// Assert
		assert.Error(t, err)

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

func TestSiteRepo_GetFreeNodesForOrg(t *testing.T) {
	var (
		nodeId = ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)
		orgId  = uuid.NewV4()
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

	t.Run("FreeNodesFoundForOrg", func(t *testing.T) {
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
		nodes, err := r.GetFreeNodesForOrg(orgId)

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
		nodes, err := r.GetFreeNodesForOrg(orgId)

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
		assert.Error(t, err)
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
		assert.Error(t, err)
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
