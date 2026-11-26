/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package db

import (
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/node/operation-monitor/pkg"
)

type testDb struct {
	gormDb *gorm.DB
}

func (d testDb) Init(...interface{}) error { return nil }
func (d testDb) Connect() error            { return nil }
func (d testDb) GetGormDb() *gorm.DB       { return d.gormDb }
func (d testDb) ExecuteInTransaction(func(tx *gorm.DB) *gorm.DB,
	...func() error) error {
	return nil
}
func (d testDb) ExecuteInTransaction2(func(tx *gorm.DB) *gorm.DB,
	...func(tx *gorm.DB) error) error {
	return nil
}

var intentColumns = []string{"id", "operation_id", "resource_key", "action_type",
	"fencing_token", "completion_rule", "status", "armed", "dispatched", "deadline"}

func setupIntentRepo(t *testing.T) (sqlmock.Sqlmock, IntentRepo) {
	t.Helper()
	sqlDb, mock, err := sqlmock.New()
	require.NoError(t, err)
	gdb, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 sqlDb,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	require.NoError(t, err)
	return mock, NewIntentRepo(testDb{gormDb: gdb})
}

func intentRow(opId uuid.UUID, st IntentStatus) *sqlmock.Rows {
	return sqlmock.NewRows(intentColumns).AddRow(uuid.NewV4(), opId, "node:abc",
		"RestartNode", uint64(3), pkg.RuleSubstateOn, int64(st), false, false, time.Now())
}

func TestIntentRepoAdd(t *testing.T) {
	mock, repo := setupIntentRepo(t)
	intent := &Intent{Id: uuid.NewV4(), OperationId: uuid.NewV4(),
		ResourceKey: "node:abc", ActionType: "RestartNode", FencingToken: 3,
		CompletionRule: pkg.RuleSubstateOn, Deadline: time.Now().Add(time.Minute)}

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO "intents"`).
		WithArgs(intent.Id, intent.OperationId, "node:abc", "RestartNode", 3,
			pkg.RuleSubstateOn, IntentWatching, false, false, intent.Deadline,
			sqlmock.AnyArg(), sqlmock.AnyArg(), nil).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	assert.NoError(t, repo.Add(intent))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestIntentRepoGet(t *testing.T) {
	t.Run("Found", func(t *testing.T) {
		mock, repo := setupIntentRepo(t)
		opId := uuid.NewV4()

		mock.ExpectQuery(`SELECT \* FROM "intents" WHERE operation_id = \$1`).
			WithArgs(opId, 1).WillReturnRows(intentRow(opId, IntentWatching))

		out, err := repo.Get(opId)

		require.NoError(t, err)
		assert.Equal(t, opId, out.OperationId)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("NotFound", func(t *testing.T) {
		mock, repo := setupIntentRepo(t)
		opId := uuid.NewV4()

		mock.ExpectQuery(`SELECT \* FROM "intents" WHERE operation_id = \$1`).
			WithArgs(opId, 1).WillReturnRows(sqlmock.NewRows(intentColumns))

		out, err := repo.Get(opId)

		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
		assert.Nil(t, out)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestIntentRepoList(t *testing.T) {
	watching := IntentWatching

	t.Run("NoFilter", func(t *testing.T) {
		mock, repo := setupIntentRepo(t)

		mock.ExpectQuery(`SELECT \* FROM "intents" WHERE "intents"."deleted_at" IS NULL ` +
			`ORDER BY deadline ASC$`).
			WillReturnRows(intentRow(uuid.NewV4(), IntentCompleted))

		out, err := repo.List(IntentFilter{})

		require.NoError(t, err)
		assert.Len(t, out, 1)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("WatchingByResource", func(t *testing.T) {
		mock, repo := setupIntentRepo(t)
		opId := uuid.NewV4()

		mock.ExpectQuery(`SELECT \* FROM "intents" WHERE resource_key = \$1 `+
			`AND status = \$2 AND "intents"."deleted_at" IS NULL ORDER BY deadline ASC$`).
			WithArgs("node:abc", IntentWatching).
			WillReturnRows(intentRow(opId, IntentWatching))

		out, err := repo.List(IntentFilter{ResourceKey: "node:abc", Status: &watching})

		require.NoError(t, err)
		require.Len(t, out, 1)
		assert.Equal(t, opId, out[0].OperationId)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ExpiredOldestFirst", func(t *testing.T) {
		mock, repo := setupIntentRepo(t)
		now := time.Now()

		mock.ExpectQuery(`SELECT \* FROM "intents" WHERE status = \$1 AND deadline < \$2 `+
			`AND "intents"."deleted_at" IS NULL ORDER BY deadline ASC LIMIT \$3$`).
			WithArgs(IntentWatching, now, 10).
			WillReturnRows(intentRow(uuid.NewV4(), IntentWatching))

		out, err := repo.List(IntentFilter{Status: &watching, DeadlineBefore: now,
			Limit: 10})

		require.NoError(t, err)
		assert.Len(t, out, 1)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("PassesThroughDbError", func(t *testing.T) {
		mock, repo := setupIntentRepo(t)

		mock.ExpectQuery(`SELECT \* FROM "intents"`).WillReturnError(errors.New("db down"))

		_, err := repo.List(IntentFilter{Status: &watching})

		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestIntentRepoUpdate(t *testing.T) {
	updateWatching := `UPDATE "intents" SET (.+) WHERE \(operation_id = \$\d+ ` +
		`AND status = \$\d+\) AND "intents"."deleted_at" IS NULL RETURNING \*`
	selectIntent := `SELECT \* FROM "intents" WHERE operation_id = \$1`
	failed, watching, armed := IntentFailed, IntentWatching, true

	t.Run("ClosesWatchingIntent", func(t *testing.T) {
		mock, repo := setupIntentRepo(t)
		opId := uuid.NewV4()

		mock.ExpectBegin()
		mock.ExpectQuery(updateWatching).
			WithArgs(IntentFailed, sqlmock.AnyArg(), opId, IntentWatching).
			WillReturnRows(intentRow(opId, IntentFailed))
		mock.ExpectCommit()

		out, err := repo.Update(opId, IntentUpdate{Status: &failed})

		require.NoError(t, err)
		assert.Equal(t, IntentFailed, out.Status)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ArmsWatchingIntent", func(t *testing.T) {
		mock, repo := setupIntentRepo(t)
		opId := uuid.NewV4()

		mock.ExpectBegin()
		mock.ExpectQuery(updateWatching).
			WithArgs(true, sqlmock.AnyArg(), opId, IntentWatching).
			WillReturnRows(intentRow(opId, IntentWatching))
		mock.ExpectCommit()

		out, err := repo.Update(opId, IntentUpdate{Armed: &armed})

		require.NoError(t, err)
		assert.Equal(t, opId, out.OperationId)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DispatchesWatchingIntent", func(t *testing.T) {
		mock, repo := setupIntentRepo(t)
		opId := uuid.NewV4()
		dispatched := true

		mock.ExpectBegin()
		mock.ExpectQuery(updateWatching).
			WithArgs(true, sqlmock.AnyArg(), opId, IntentWatching).
			WillReturnRows(intentRow(opId, IntentWatching))
		mock.ExpectCommit()

		out, err := repo.Update(opId, IntentUpdate{Dispatched: &dispatched})

		require.NoError(t, err)
		assert.Equal(t, opId, out.OperationId)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("AlreadyTerminalIsKept", func(t *testing.T) {
		mock, repo := setupIntentRepo(t)
		opId := uuid.NewV4()

		mock.ExpectBegin()
		mock.ExpectQuery(updateWatching).WillReturnRows(sqlmock.NewRows(intentColumns))
		mock.ExpectCommit()
		mock.ExpectQuery(selectIntent).WithArgs(opId, 1).
			WillReturnRows(intentRow(opId, IntentCompleted))

		out, err := repo.Update(opId, IntentUpdate{Status: &failed})

		require.NoError(t, err)
		assert.Equal(t, IntentCompleted, out.Status)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("UnknownIntentIsNotFound", func(t *testing.T) {
		mock, repo := setupIntentRepo(t)
		opId := uuid.NewV4()

		mock.ExpectBegin()
		mock.ExpectQuery(updateWatching).WillReturnRows(sqlmock.NewRows(intentColumns))
		mock.ExpectCommit()
		mock.ExpectQuery(selectIntent).WithArgs(opId, 1).
			WillReturnRows(sqlmock.NewRows(intentColumns))

		out, err := repo.Update(opId, IntentUpdate{Armed: &armed})

		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
		assert.Nil(t, out)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("PassesThroughDbError", func(t *testing.T) {
		mock, repo := setupIntentRepo(t)

		mock.ExpectBegin()
		mock.ExpectQuery(updateWatching).WillReturnError(errors.New("db down"))
		mock.ExpectRollback()

		out, err := repo.Update(uuid.NewV4(), IntentUpdate{Status: &failed})

		assert.Error(t, err)
		assert.Nil(t, out)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("NonTerminalStatusIsRejected", func(t *testing.T) {
		mock, repo := setupIntentRepo(t)

		out, err := repo.Update(uuid.NewV4(), IntentUpdate{Status: &watching})

		assert.Error(t, err)
		assert.Nil(t, out)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("EmptyUpdateIsRejected", func(t *testing.T) {
		mock, repo := setupIntentRepo(t)

		out, err := repo.Update(uuid.NewV4(), IntentUpdate{})

		assert.Error(t, err)
		assert.Nil(t, out)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestIntentStatus(t *testing.T) {
	assert.Equal(t, "failed", IntentFailed.String())
	assert.True(t, IntentExpired.IsTerminal())
	assert.False(t, IntentWatching.IsTerminal())

	var st IntentStatus
	require.NoError(t, st.Scan(nil))
	assert.Equal(t, IntentWatching, st)
	require.NoError(t, st.Scan(int64(2)))
	assert.Equal(t, IntentFailed, st)
	v, err := IntentCancelled.Value()
	require.NoError(t, err)
	assert.Equal(t, int64(4), v)
}
