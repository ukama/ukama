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

	"github.com/ukama/ukama/systems/common/ukama"
	int_db "github.com/ukama/ukama/systems/node/configurator/pkg/db"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestConfigRepo_Get(t *testing.T) {

	t.Run("NodeExist", func(t *testing.T) {
		// Arrange
		nid := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)

		var db *extsql.DB
		var err error

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"node_id"}).
			AddRow(nid.String())

		mock.ExpectQuery(`^SELECT.*configurations.*`).
			WithArgs(nid.String()).
			WillReturnRows(rows)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := int_db.NewConfigRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		cfg, err := r.Get(nid.String())

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		if assert.NotNil(t, cfg) {
			assert.Equal(t, nid.String(), cfg.NodeId)
		}
	})

	t.Run("Node Doesn't Exist", func(t *testing.T) {
		// Arrange
		nid := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)
		var db *extsql.DB
		var err error

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		mock.ExpectQuery(`^SELECT.*configurations.*`).
			WithArgs(nid.String()).
			WillReturnError(gorm.ErrRecordNotFound)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := int_db.NewConfigRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		_, err = r.Get(nid.String())

		// Assert
		if assert.Error(t, err) {
			assert.Equal(t, true, errors.Is(gorm.ErrRecordNotFound, err))
		}

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)

	})

	t.Run("GetAll", func(t *testing.T) {
		// Arrange

		nid0 := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)
		nid1 := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)

		var db *extsql.DB
		var err error

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"node_id"}).
			AddRow(nid0).AddRow(nid1)

		mock.ExpectQuery(`^SELECT.*configurations.*`).
			WillReturnRows(rows)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := int_db.NewConfigRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		cfg, err := r.GetAll()

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		if assert.NotNil(t, cfg) {
			assert.Equal(t, nid0.String(), cfg[0].NodeId)
			assert.Equal(t, nid1.String(), cfg[1].NodeId)
		}
	})
}

func TestConfigRepo_Add(t *testing.T) {

	t.Run("Add", func(t *testing.T) {
		// Arrange

		nid := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)
		hash := ""
		id := 1
		var db *extsql.DB
		var err error

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		mock.ExpectBegin()

		mock.ExpectQuery(regexp.QuoteMeta(`INSERT`)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), hash, id).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		mock.ExpectQuery(regexp.QuoteMeta(`INSERT`)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), hash, id).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		mock.ExpectQuery(regexp.QuoteMeta(`INSERT`)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), nid.String(), int_db.Default, 1, 1, 0).
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

		r := int_db.NewConfigRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		err = r.Add(nid.String())

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

}

func TestConfigRepo_Update(t *testing.T) {

	t.Run("UpdateLastCommit", func(t *testing.T) {
		// Arrange

		nid := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)
		state := int_db.Published
		c := int_db.Configuration{
			Model:  gorm.Model{ID: 1},
			NodeId: nid.String(),
			LastCommit: int_db.Commit{
				Model: gorm.Model{ID: 1},
				Hash:  "6b0a48e3d06ae7708b2257321d17b36bd930f670",
			},
			LastCommitId: 1,
		}

		var db *extsql.DB
		var err error

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE`)).
			WithArgs(sqlmock.AnyArg(), c.LastCommit.ID, c.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE`)).
			WithArgs(state, sqlmock.AnyArg(), nid.String()).
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

		r := int_db.NewConfigRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		err = r.UpdateLastCommit(c, &state)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("UpdateCurrentCommit", func(t *testing.T) {
		// Arrange

		nid := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)
		state := int_db.Published
		c := int_db.Configuration{
			Model:  gorm.Model{ID: 1},
			NodeId: nid.String(),
			Commit: int_db.Commit{
				Model: gorm.Model{ID: 1},
				Hash:  "6b0a48e3d06ae7708b2257321d17b36bd930f670",
			},
			LastCommitId: 1,
		}

		var db *extsql.DB
		var err error

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT`)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), c.Commit.Hash, c.Commit.ID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE`)).
			WithArgs(sqlmock.AnyArg(), c.Commit.ID, c.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE`)).
			WithArgs(state, sqlmock.AnyArg(), nid.String()).
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

		r := int_db.NewConfigRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		err = r.UpdateCurrentCommit(c, &state)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

}
