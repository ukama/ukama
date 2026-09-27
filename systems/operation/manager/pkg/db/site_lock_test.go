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
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func newSiteTestOperation(key string) *Operation {
	return &Operation{Id: uuid.NewV4(), Type: "RestartNode", System: "node",
		Status: OperationPending, ResourceKey: key, LeaseExpiresAt: time.Now().Add(time.Minute)}
}

func TestStartBatchRejectsBusyControllerBeforeAnyInsert(t *testing.T) {
	mock, repo := setupTestDB(t)
	holder := newSiteTestOperation("node:c")
	holder.Type, holder.Status = "UpdateSoftware", OperationRunning
	mock.ExpectBegin()
	for i := 0; i < 3; i++ {
		mock.ExpectExec("SELECT pg_advisory_xact_lock").WillReturnResult(sqlmock.NewResult(0, 1))
	}
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "resource_locks"`)).
		WithArgs("node:a", "node:c", "node:t", 1).
		WillReturnRows(sqlmock.NewRows([]string{"resource_key", "operation_id"}).AddRow(holder.ResourceKey, holder.Id))
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "operations"`)).WillReturnRows(operationRow(holder))
	mock.ExpectRollback()
	ops, blockedBy, err := repo.StartBatch([]*Operation{newSiteTestOperation("node:t"), newSiteTestOperation("node:a")},
		[]string{"node:t", "node:c", "node:a"}, time.Minute)
	require.ErrorIs(t, err, ErrLockConflict)
	require.Empty(t, ops)
	require.Equal(t, holder.Id, blockedBy.Id)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStartBatchRollsBackAllTargetsOnInsertFailure(t *testing.T) {
	mock, repo := setupTestDB(t)
	tower, amp := newSiteTestOperation("node:t"), newSiteTestOperation("node:a")
	mock.ExpectBegin()
	for i := 0; i < 3; i++ {
		mock.ExpectExec("SELECT pg_advisory_xact_lock").WillReturnResult(sqlmock.NewResult(0, 1))
	}
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "resource_locks"`)).WillReturnError(gorm.ErrRecordNotFound)
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "operations"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "fencing_token"}).AddRow(tower.Id, 1))
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "resource_locks"`)).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "operation_audits"`)).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "operations"`)).WillReturnError(errors.New("insert failed"))
	mock.ExpectRollback()
	ops, _, err := repo.StartBatch([]*Operation{tower, amp}, []string{"node:c", "node:t", "node:a"}, time.Minute)
	require.Error(t, err)
	require.Empty(t, ops)
	require.NoError(t, mock.ExpectationsWereMet())
}

// Set UKAMA_OPERATION_TEST_DSN to a PostgreSQL test database to exercise real
// concurrent transactions. Every run creates and removes its own schema.
func TestSiteAdmissionPostgres(t *testing.T) {
	dsn := os.Getenv("UKAMA_OPERATION_TEST_DSN")
	if dsn == "" { t.Skip("UKAMA_OPERATION_TEST_DSN is not set") }
	admin, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	require.NoError(t, err)
	schema := "operation_test_" + strings.ReplaceAll(uuid.NewV4().String(), "-", "")
	require.NoError(t, admin.Exec("CREATE SCHEMA " + schema).Error)
	t.Cleanup(func() {
		_ = admin.Exec("DROP SCHEMA " + schema + " CASCADE").Error
		if db, err := admin.DB(); err == nil { _ = db.Close() }
	})
	if strings.HasPrefix(dsn, "postgres://") || strings.HasPrefix(dsn, "postgresql://") {
		u, err := url.Parse(dsn)
		require.NoError(t, err)
		query := u.Query(); query.Set("search_path", schema); u.RawQuery = query.Encode(); dsn = u.String()
	} else { dsn += " search_path=" + schema }
	gdb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	require.NoError(t, err)
	t.Cleanup(func() { if db, err := gdb.DB(); err == nil { _ = db.Close() } })
	require.NoError(t, gdb.AutoMigrate(&Operation{}, &ResourceLock{}, &OperationAudit{}))
	repo := NewOperationRepo(&UkamaDbMock{GormDb: gdb})
	keys := []string{"node:c", "node:t", "node:a"}

	// Simultaneous updates/controls for different nodes cannot both win.
	var wg sync.WaitGroup
	start := make(chan struct{})
	results := make(chan error, 2)
	for _, key := range []string{"node:c", "node:t"} {
		key := key
		wg.Add(1)
		go func() {
			defer wg.Done(); <-start
			_, _, err := repo.StartBatch([]*Operation{newSiteTestOperation(key)}, keys, time.Minute)
			results <- err
		}()
	}
	close(start); wg.Wait(); close(results)
	won, blocked := 0, 0
	for err := range results {
		if err == nil { won++ } else if errors.Is(err, ErrLockConflict) { blocked++ } else { t.Fatal(err) }
	}
	require.Equal(t, 1, won); require.Equal(t, 1, blocked)
	for _, key := range keys {
		op, err := repo.GetByResource(key); require.NoError(t, err)
		if op != nil {
			_, err = repo.Terminate(op.Id, op.FencingToken, OperationSuccess, OperationAudit{Event: "completed"}, "")
			require.NoError(t, err)
		}
	}

	// Reserve only tower/amplifier; one completed child must not unlock site.
	ops, _, err := repo.StartBatch([]*Operation{newSiteTestOperation("node:t"), newSiteTestOperation("node:a")}, keys, time.Minute)
	require.NoError(t, err)
	cnode, err := repo.GetByResource("node:c"); require.NoError(t, err); require.Nil(t, cnode)
	_, err = repo.Terminate(ops[0].Id, ops[0].FencingToken, OperationSuccess, OperationAudit{Event: "completed"}, "")
	require.NoError(t, err)
	_, _, err = repo.StartBatch([]*Operation{newSiteTestOperation("node:c")}, keys, time.Minute)
	require.ErrorIs(t, err, ErrLockConflict)
	_, err = repo.Terminate(ops[1].Id, ops[1].FencingToken, OperationSuccess, OperationAudit{Event: "completed"}, "")
	require.NoError(t, err)
	_, _, err = repo.StartBatch([]*Operation{newSiteTestOperation("node:c")}, keys, time.Minute)
	require.NoError(t, err)
}
