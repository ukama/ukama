/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package server

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	npb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/node/state/pkg/db"
	"github.com/ukama/ukama/systems/node/state/pkg/lifecycle"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Set LIFECYCLE_TEST_DSN to a disposable PostgreSQL database. Each test uses a
// private schema and removes only that schema on completion.
func lifecycleDatabase(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := os.Getenv("LIFECYCLE_TEST_DSN")
	if dsn == "" {
		t.Skip("LIFECYCLE_TEST_DSN is not set")
	}
	connection, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	require.NoError(t, err)
	sqlDB, err := connection.DB()
	require.NoError(t, err)
	sqlDB.SetMaxOpenConns(1)
	schema := "lifecycle_test_" + strings.ReplaceAll(uuid.NewV4().String(), "-", "")
	require.NoError(t, connection.Exec("CREATE SCHEMA "+schema).Error)
	require.NoError(t, connection.Exec("SET search_path TO "+schema).Error)
	t.Cleanup(func() {
		connection.Exec("DROP SCHEMA " + schema + " CASCADE")
		sqlDB.Close()
	})
	require.NoError(t, connection.AutoMigrate(&db.NodeConfig{}, &db.State{}, &db.LifecycleRecord{}, &db.LifecyclePublication{}))
	return connection
}

func TestLifecyclePersistenceAndRecovery(t *testing.T) {
	connection := lifecycleDatabase(t)
	ctx := context.Background()
	nodeID := ukama.NewVirtualHomeNodeId().String()
	repo := db.NewLifecycleRepo(connection)
	server := NewStateEventServer("org", "", nil, nodeStateConfigPath, nil)
	server.SetLifecycleRepo(repo)
	read := func() (db.LifecycleRecord, db.State) {
		var record db.LifecycleRecord
		var state db.State
		require.NoError(t, connection.First(&record, "node_id = ?", nodeID).Error)
		require.NoError(t, connection.Preload("Config").Where("node_id = ?", nodeID).Order("created_at DESC").First(&state).Error)
		return record, state
	}
	observe := func(boot, value string, seq uint64, request string) error {
		return server.processStoredEvent(ctx, value, nodeID, &lifecycle.Event{
			BootID: boot, State: value, Sequence: seq, Time: 100, RequestID: request, ConfigMode: "NOCONFIG", Generation: 1,
		})
	}
	require.NoError(t, observe("a", "INIT", 1, ""))
	require.NoError(t, server.ProcessEvent(ctx, "online", nodeID, &epb.NodeOnlineEvent{NodeId: nodeID, NodeIp: "192.0.2.10", NodePort: 8080}))
	require.NoError(t, observe("a", "READY", 2, ""))
	record, state := read()
	require.Empty(t, record.RequestID)
	require.Equal(t, npb.NodeState_Ready, state.CurrentState)
	sends := []string{}
	send := func(record *db.LifecycleRecord) { sends = append(sends, record.RequestID) }
	require.NoError(t, repo.RetryAssignments(ctx, time.Now().Add(time.Hour), send))
	require.Empty(t, sends, "unassigned READY must not send NOCONFIG")
	assignment := &epb.EventRegistryNodeAssign{NodeId: nodeID, Network: "network", Site: "site"}
	require.NoError(t, server.ProcessEvent(ctx, "assign", nodeID, assignment))
	record, _ = read()
	requestID := record.RequestID
	now := record.RetryAt.Add(time.Millisecond)
	require.NoError(t, repo.RetryAssignments(ctx, now, send))
	require.Len(t, sends, 1)
	require.NoError(t, repo.RetryAssignments(ctx, now.Add(59*time.Second), send))
	require.Len(t, sends, 1)
	// Recreate the backend objects, preserving only the database.
	repo = db.NewLifecycleRepo(connection)
	server = NewStateEventServer("org", "", nil, nodeStateConfigPath, nil)
	server.SetLifecycleRepo(repo)
	require.NoError(t, repo.RetryAssignments(ctx, now.Add(60*time.Second), send))
	require.Equal(t, []string{requestID, requestID}, sends)
	require.NoError(t, server.ProcessEvent(ctx, "assign", nodeID, assignment))
	record, _ = read()
	require.Equal(t, requestID, record.RequestID)
	require.Error(t, observe("a", "OPERATIONAL", 4, "wrong-request"))
	record, _ = read()
	require.Equal(t, uint64(2), record.Cursor.Sequence, "rejected events do not advance the cursor")
	require.NoError(t, observe("a", "CONFIGURING", 3, requestID))
	// Fail the final write after the cursor and assignment updates. All three
	// changes must roll back, permitting the same notification to retry.
	require.NoError(t, connection.Callback().Create().Before("gorm:create").Register("test:reject_operational", func(tx *gorm.DB) {
		if item, ok := tx.Statement.Dest.(*db.LifecyclePublication); ok && item.State == "Operational" {
			tx.AddError(errors.New("injected publication write failure"))
		}
	}))
	require.Error(t, observe("a", "OPERATIONAL", 4, requestID))
	record, state = read()
	require.False(t, record.Completed)
	require.Equal(t, uint64(3), record.Cursor.Sequence)
	require.Equal(t, npb.NodeState_Configuring, state.CurrentState)
	require.NoError(t, connection.Callback().Create().Remove("test:reject_operational"))
	require.NoError(t, observe("a", "OPERATIONAL", 4, requestID))
	record, state = read()
	require.True(t, record.Completed)
	require.False(t, record.AwaitingOperational)
	require.Equal(t, "192.0.2.10", state.Config.NodeIp, "lifecycle must preserve connection metadata")
	require.NoError(t, repo.RetryAssignments(ctx, now.Add(time.Hour), send))
	require.Len(t, sends, 2)
	require.NoError(t, observe("b", "INIT", 1, ""))
	require.NoError(t, observe("b", "READY", 2, ""))
	record, _ = read()
	require.True(t, record.AwaitingOperational)
	require.Equal(t, requestID, record.RequestID)
	require.NoError(t, observe("b", "CONFIGURING", 3, requestID))
	require.NoError(t, observe("b", "OPERATIONAL", 4, requestID))
	require.NoError(t, observe("a", "READY", 100, ""))
	record, state = read()
	require.Equal(t, npb.NodeState_Operational, state.CurrentState)
	require.False(t, record.AwaitingOperational)
	require.NoError(t, repo.RetryAssignments(ctx, now.Add(2*time.Hour), send))
	require.Len(t, sends, 2, "restored completion needs no new command")
	// Failed publication remains durable; a new worker drains the same row.
	var firstID uint64
	require.Error(t, repo.PublishPending(ctx, func(item *db.LifecyclePublication) error {
		firstID = item.ID
		return errors.New("broker unavailable")
	}))
	replayed := false
	repo = db.NewLifecycleRepo(connection)
	require.NoError(t, repo.PublishPending(ctx, func(item *db.LifecyclePublication) error {
		if !replayed {
			require.Equal(t, firstID, item.ID)
			replayed = true
		}
		return nil
	}))
	require.True(t, replayed)
	var count int64
	require.NoError(t, connection.Model(&db.LifecyclePublication{}).Count(&count).Error)
	require.Zero(t, count)
}
