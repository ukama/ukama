/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package db_test

import (
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	softwaredb "github.com/ukama/ukama/systems/node/software/pkg/db"
)

func TestNodeRepoDelete(t *testing.T) {
	const nodeID = "uk-sa2643-tnode-v0-a60f"

	setup := func(t *testing.T) (sqlmock.Sqlmock, softwaredb.NodeRepo) {
		t.Helper()
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		gdb, err := gorm.Open(postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		}), &gorm.Config{})
		require.NoError(t, err)
		return mock, softwaredb.NewNodeRepo(ukamaDbMock{gormDb: gdb})
	}

	t.Run("deletes software rows then the node in one transaction", func(t *testing.T) {
		mock, repo := setup(t)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "softwares" WHERE node_id = $1`)).
			WithArgs(nodeID).
			WillReturnResult(sqlmock.NewResult(0, 18))
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "nodes" WHERE node_id = $1`)).
			WithArgs(nodeID).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		assert.NoError(t, repo.Delete(nodeID))
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("rolls back when the node delete fails", func(t *testing.T) {
		mock, repo := setup(t)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "softwares" WHERE node_id = $1`)).
			WithArgs(nodeID).
			WillReturnResult(sqlmock.NewResult(0, 18))
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "nodes" WHERE node_id = $1`)).
			WithArgs(nodeID).
			WillReturnError(errors.New("delete failed"))
		mock.ExpectRollback()

		assert.Error(t, repo.Delete(nodeID))
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
