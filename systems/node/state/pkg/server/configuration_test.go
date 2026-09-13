/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 * Copyright (c) 2026-present, Ukama Inc.
 */
package server

import (
	"context"
	"github.com/stretchr/testify/mock"
	mbmocks "github.com/ukama/ukama/systems/common/mocks"
	"github.com/ukama/ukama/systems/common/ukama"
	pb "github.com/ukama/ukama/systems/node/state/pb/gen"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"

	"github.com/stretchr/testify/require"
	npb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
	"github.com/ukama/ukama/systems/node/state/pkg/db"
	"github.com/ukama/ukama/systems/node/state/pkg/lifecycle"
	"testing"
)

func TestCancelledAttemptCannotBecomeOperational(t *testing.T) {
	n := NewStateEventServer("org", "", nil, nodeStateConfigPath, nil)
	record := &db.LifecycleRecord{RequestID: "attempt-1", ProvisionManaged: true,
		Cursor:   lifecycle.Cursor{BootID: "boot", Sequence: 1},
		Attempts: map[string]db.ProvisionAttempt{"attempt-1": {Cancelled: true}}}
	state := &db.State{CurrentState: npb.NodeState_Configuring}
	require.NoError(t, n.observeLifecycle(record, state, &lifecycle.Event{BootID: "boot", Sequence: 2, State: "OPERATIONAL", RequestID: "attempt-1", ConfigMode: "NOCONFIG", Generation: 1}))
	require.False(t, record.Completed)
	require.Equal(t, npb.NodeState_Configuring, state.CurrentState)
	require.NoError(t, n.observeLifecycle(record, state, &lifecycle.Event{BootID: "boot", Sequence: 3, State: "READY", RequestID: "attempt-1", ConfigMode: "NONE", Generation: 2}))
	require.True(t, record.Attempts["attempt-1"].Cleared)
	require.Equal(t, npb.NodeState_Ready, state.CurrentState)
	require.NoError(t, n.observeLifecycle(record, state, &lifecycle.Event{BootID: "boot", Sequence: 4, State: "OPERATIONAL", RequestID: "attempt-1", ConfigMode: "NOCONFIG", Generation: 1}))
	require.Equal(t, npb.NodeState_Ready, state.CurrentState)
}

func TestProvisioningBootDoesNotEnableAssignmentRetry(t *testing.T) {
	n := NewStateEventServer("org", "", nil, nodeStateConfigPath, nil)
	record := &db.LifecycleRecord{RequestID: "attempt-1", ProvisionManaged: true}
	state := &db.State{CurrentState: npb.NodeState_Ready}
	require.NoError(t, n.observeLifecycle(record, state, &lifecycle.Event{BootID: "boot", Sequence: 1, State: "INIT"}))
	require.False(t, record.AwaitingOperational)
}

func TestConfigurationRecoveryDatabase(t *testing.T) {
	dsn := os.Getenv("PROVISION_TEST_DSN")
	if dsn == "" {
		t.Skip("PROVISION_TEST_DSN is not set")
	}
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	require.NoError(t, err)
	pool, poolErr := database.DB()
	require.NoError(t, poolErr)
	if os.Getenv("PROVISION_TEST_SINGLE_CONNECTION") != "" {
		pool.SetMaxOpenConns(1)
	}
	require.NoError(t, database.AutoMigrate(&db.State{}, &db.LifecycleRecord{}, &db.LifecyclePublication{}))
	bus := &mbmocks.MsgBusServiceClient{}
	bus.On("PublishRequest", mock.Anything, mock.Anything).Return(nil)
	events := NewStateEventServer("org", "", nil, nodeStateConfigPath, bus)
	events.SetLifecycleRepo(db.NewLifecycleRepo(database))
	server := &StateServer{configurationEvents: events}
	nodeID := ukama.NewVirtualTowerNodeId().String()
	req := &pb.RecordConfigurationRequest{NodeId: nodeID, RequestId: "attempt-1", SiteId: "site-1", NetworkId: "net-1", Cancelled: true}
	ctx := context.Background()
	out, err := server.RecordConfiguration(ctx, req)
	require.NoError(t, err)
	require.True(t, out.Cancelled)
	require.False(t, out.Cleared)
	req.Cancelled = false
	_, err = server.RecordConfiguration(ctx, req)
	require.Error(t, err)
	require.NoError(t, events.processStoredEvent(ctx, "init", nodeID, &lifecycle.Event{State: "INIT", BootID: "boot", Sequence: 1}))
	require.NoError(t, events.processStoredEvent(ctx, "platformready", nodeID, &lifecycle.Event{State: "READY", BootID: "boot", Sequence: 2, RequestID: "attempt-1", ConfigMode: "NONE", Generation: 1}))
	// Recreate the server/repository as after a process crash.
	events.SetLifecycleRepo(db.NewLifecycleRepo(database))
	server = &StateServer{configurationEvents: events}
	out, err = server.configurationStatus(ctx, nodeID, req.RequestId)
	require.NoError(t, err)
	require.True(t, out.Cleared)
	req.RequestId = "attempt-2"
	req.Cancelled = false
	_, err = server.RecordConfiguration(ctx, req)
	require.NoError(t, err)
	require.NoError(t, events.processStoredEvent(ctx, "operational", nodeID, &lifecycle.Event{State: "OPERATIONAL", BootID: "boot", Sequence: 3, RequestID: "attempt-1", ConfigMode: "NOCONFIG", Generation: 1}))
	out, err = server.configurationStatus(ctx, nodeID, req.RequestId)
	require.NoError(t, err)
	require.False(t, out.Completed)
	require.NoError(t, events.processStoredEvent(ctx, "operational", nodeID, &lifecycle.Event{State: "OPERATIONAL", BootID: "boot", Sequence: 4, RequestID: "attempt-2", ConfigMode: "NOCONFIG", Generation: 2}))
	out, err = server.configurationStatus(ctx, nodeID, req.RequestId)
	require.NoError(t, err)
	require.True(t, out.Completed)
}
