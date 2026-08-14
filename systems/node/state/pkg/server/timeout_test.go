/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
package server

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	npb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
	stm "github.com/ukama/ukama/systems/common/stateMachine"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/node/state/mocks"
	"github.com/ukama/ukama/systems/node/state/pkg/db"
)

func newTimeoutServer(t *testing.T, repo *mocks.StateRepo) *StateEventServer {
	t.Helper()

	return &StateEventServer{
		stateMachine:  stm.NewStateMachine(func(stm.Event) {}),
		configPath:    nodeStateConfigPath,
		instances:     make(map[string]*stm.StateMachineInstance),
		latchedHealth: make(map[string]string),
		s:             NewStateServer("test-org", "test-org-id", repo, nil),
	}
}

func timeoutStateRow(nodeId string, state npb.NodeState, enteredAt time.Time) db.State {
	return db.State{
		Id:           uuid.NewV4(),
		NodeId:       nodeId,
		CurrentState: state,
		SubState:     db.StringArray{"on"},
		CreatedAt:    enteredAt,
	}
}

func TestRunTimeouts(t *testing.T) {
	nodeId := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE).String()

	t.Run("advances a node that has been in Configuring too long", func(t *testing.T) {
		now := time.Now().UTC()
		enteredAt := now.Add(-90 * time.Second)

		repo := &mocks.StateRepo{}
		repo.On("ListLatestStates").
			Return([]db.State{timeoutStateRow(nodeId, npb.NodeState_Configuring, enteredAt)}, nil).Once()
		repo.On("UpdateState", nodeId, mock.Anything, mock.Anything).
			Return(&db.State{NodeId: nodeId}, nil).Once()
		repo.On("GetLatestState", nodeId).
			Return(&db.State{NodeId: nodeId}, nil).Maybe()
		repo.On("AddState", mock.Anything, mock.Anything).Return(nil).Once()

		srv := newTimeoutServer(t, repo)

		advanced, err := srv.RunTimeouts(context.Background(), now)

		require.NoError(t, err)
		assert.Equal(t, 1, advanced)
		repo.AssertExpectations(t)
	})

	t.Run("leaves a node that is still inside its window", func(t *testing.T) {
		now := time.Now().UTC()
		enteredAt := now.Add(-30 * time.Second)

		repo := &mocks.StateRepo{}
		repo.On("ListLatestStates").
			Return([]db.State{timeoutStateRow(nodeId, npb.NodeState_Configuring, enteredAt)}, nil).Once()

		srv := newTimeoutServer(t, repo)

		advanced, err := srv.RunTimeouts(context.Background(), now)

		require.NoError(t, err)
		assert.Equal(t, 0, advanced)
		repo.AssertNotCalled(t, "AddState", mock.Anything, mock.Anything)
	})

	t.Run("leaves a state that declares no timeout", func(t *testing.T) {
		now := time.Now().UTC()
		enteredAt := now.Add(-24 * time.Hour)

		repo := &mocks.StateRepo{}
		repo.On("ListLatestStates").
			Return([]db.State{timeoutStateRow(nodeId, npb.NodeState_Operational, enteredAt)}, nil).Once()

		srv := newTimeoutServer(t, repo)

		advanced, err := srv.RunTimeouts(context.Background(), now)

		require.NoError(t, err)
		assert.Equal(t, 0, advanced)
		repo.AssertNotCalled(t, "AddState", mock.Anything, mock.Anything)
	})

	t.Run("keeps sweeping when one node fails", func(t *testing.T) {
		now := time.Now().UTC()
		enteredAt := now.Add(-90 * time.Second)
		healthyNode := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE).String()

		repo := &mocks.StateRepo{}
		repo.On("ListLatestStates").Return([]db.State{
			timeoutStateRow("not-a-node-id", npb.NodeState_Configuring, enteredAt),
			timeoutStateRow(healthyNode, npb.NodeState_Configuring, enteredAt),
		}, nil).Once()
		repo.On("UpdateState", healthyNode, mock.Anything, mock.Anything).
			Return(&db.State{NodeId: healthyNode}, nil).Once()
		repo.On("GetLatestState", healthyNode).
			Return(&db.State{NodeId: healthyNode}, nil).Maybe()
		repo.On("AddState", mock.Anything, mock.Anything).Return(nil).Once()

		srv := newTimeoutServer(t, repo)

		advanced, err := srv.RunTimeouts(context.Background(), now)

		require.NoError(t, err)
		assert.Equal(t, 1, advanced,
			"a bad row must not stop the rest of the sweep")
	})

	t.Run("reports a listing failure", func(t *testing.T) {
		repo := &mocks.StateRepo{}
		repo.On("ListLatestStates").Return(nil, assert.AnError).Once()

		srv := newTimeoutServer(t, repo)

		_, err := srv.RunTimeouts(context.Background(), time.Now().UTC())

		require.Error(t, err)
	})
}

func TestStartTimeoutWorker_DisabledWhenIntervalIsZero(t *testing.T) {
	repo := &mocks.StateRepo{}
	srv := newTimeoutServer(t, repo)

	srv.StartTimeoutWorker(context.Background(), 0)

	time.Sleep(50 * time.Millisecond)
	repo.AssertNotCalled(t, "ListLatestStates")
}
