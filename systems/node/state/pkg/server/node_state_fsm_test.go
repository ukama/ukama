/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
package server

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	stm "github.com/ukama/ukama/systems/common/stateMachine"
	"github.com/ukama/ukama/systems/node/state/mocks"
)

const nodeStateConfigPath = "../nodeState.json"

func newNodeInstance(t *testing.T, instanceID, initialState string) *stm.StateMachineInstance {
	t.Helper()

	sm := stm.NewStateMachine(func(stm.Event) {})

	instance, err := sm.NewInstance(nodeStateConfigPath, instanceID, initialState)
	require.NoError(t, err)

	return instance
}

func TestNodeStateFsm_HappyPath(t *testing.T) {
	instance := newNodeInstance(t, "happy-path", "Unknown")

	require.NoError(t, instance.Transition("online"))
	assert.Equal(t, "Unknown", instance.CurrentState,
		"connectivity must not move the lifecycle state")
	assert.Equal(t, "on", instance.CurrentSubstate)

	require.NoError(t, instance.Transition("platformready"))
	assert.Equal(t, "Ready", instance.CurrentState)

	require.NoError(t, instance.Transition("assign"))
	assert.Equal(t, "Configuring", instance.CurrentState)

	require.NoError(t, instance.Transition("configapplied"))
	assert.Equal(t, "Operational", instance.CurrentState)
}

func TestNodeStateFsm_AssignWithoutConfig(t *testing.T) {
	instance := newNodeInstance(t, "assign-no-config", "Ready")

	require.NoError(t, instance.Transition("assignnoconfig"))
	assert.Equal(t, "Operational", instance.CurrentState)
}

func TestNodeStateFsm_ReleaseReturnsToReady(t *testing.T) {
	t.Run("from configuring", func(t *testing.T) {
		instance := newNodeInstance(t, "release-configuring", "Configuring")

		require.NoError(t, instance.Transition("release"))
		assert.Equal(t, "Ready", instance.CurrentState)
	})

	t.Run("from operational", func(t *testing.T) {
		instance := newNodeInstance(t, "release-operational", "Operational")

		require.NoError(t, instance.Transition("release"))
		assert.Equal(t, "Ready", instance.CurrentState)
	})
}

func TestNodeStateFsm_ReconfigureFromOperational(t *testing.T) {
	instance := newNodeInstance(t, "reconfigure", "Operational")

	require.NoError(t, instance.Transition("configchange"))
	assert.Equal(t, "Configuring", instance.CurrentState)
}

func TestNodeStateFsm_FaultyRecovery(t *testing.T) {
	t.Run("faulty node recovers to initializing", func(t *testing.T) {
		instance := newNodeInstance(t, "faulty-recover", "Faulty")

		require.NoError(t, instance.Transition("recover"))
		assert.Equal(t, "Initializing", instance.CurrentState)
	})

	t.Run("operational to faulty and back again", func(t *testing.T) {
		instance := newNodeInstance(t, "fault-round-trip", "Operational")

		require.NoError(t, instance.Transition("fault"))
		assert.Equal(t, "Faulty", instance.CurrentState)

		require.NoError(t, instance.Transition("recover"))
		assert.Equal(t, "Initializing", instance.CurrentState)
	})

	t.Run("faulty node is reconfigurable", func(t *testing.T) {
		instance := newNodeInstance(t, "faulty-assign", "Faulty")

		require.NoError(t, instance.Transition("assign"))
		assert.Equal(t, "Configuring", instance.CurrentState)
	})
}

func TestNodeStateFsm_Offboarding(t *testing.T) {
	instance := newNodeInstance(t, "offboard-round-trip", "Operational")

	require.NoError(t, instance.Transition("offboard"))
	assert.Equal(t, "Offboarded", instance.CurrentState)

	require.NoError(t, instance.Transition("onboarding"))
	assert.Equal(t, "Unknown", instance.CurrentState,
		"an offboarded node must be able to come back")
}

func TestNodeStateFsm_UpdateReboots(t *testing.T) {
	instance := newNodeInstance(t, "update-reboot", "Initializing")

	require.NoError(t, instance.Transition("update"))
	assert.Equal(t, "Updating", instance.CurrentState)

	require.NoError(t, instance.Transition("rebooted"))
	assert.Equal(t, "Initializing", instance.CurrentState)
}

func TestNodeStateFsm_Timeouts(t *testing.T) {
	enteredAt := time.Date(2026, time.August, 14, 10, 0, 0, 0, time.UTC)

	cases := []struct {
		from string
		to   string
	}{
		{"Ready", "Configuring"},
		{"Configuring", "Operational"},
		{"Updating", "Initializing"},
	}

	for _, tc := range cases {
		t.Run(tc.from+" advances after 60s", func(t *testing.T) {
			instance := newNodeInstance(t, "timeout-"+tc.from, tc.from)

			_, due := instance.DueTransition(enteredAt, enteredAt.Add(59*time.Second))
			assert.False(t, due, "must not advance before the window closes")

			to, due := instance.DueTransition(enteredAt, enteredAt.Add(60*time.Second))
			require.True(t, due)
			assert.Equal(t, tc.to, to)
		})
	}

	t.Run("Operational never times out", func(t *testing.T) {
		instance := newNodeInstance(t, "timeout-operational", "Operational")

		_, due := instance.DueTransition(enteredAt, enteredAt.Add(time.Hour))
		assert.False(t, due)
	})
}

func TestNodeStateFsm_NoOpTransitionDoesNotPublish(t *testing.T) {
	published := 0

	sm := stm.NewStateMachine(func(e stm.Event) {
		if e.OldState == e.NewState && e.OldSubstate == e.NewSubstate {
			return
		}
		published++
	})

	instance, err := sm.NewInstance(nodeStateConfigPath, "no-op-publish", "Operational")
	require.NoError(t, err)

	require.NoError(t, instance.Transition("online"))
	assert.Equal(t, 1, published, "online changes substate and must publish")

	require.NoError(t, instance.Transition("online"))
	assert.Equal(t, 1, published, "repeat online must not publish a state change")

	require.NoError(t, instance.Transition("fault"))
	assert.Equal(t, 2, published, "fault must publish")
}

func TestStateEventServer_getOrCreateInstance_ResyncsWithStoredState(t *testing.T) {
	srv := &StateEventServer{
		stateMachine: stm.NewStateMachine(func(stm.Event) {}),
		configPath:   nodeStateConfigPath,
		instances:    make(map[string]*stm.StateMachineInstance),
	}

	nodeId := "test-node-resync"

	first, err := srv.getOrCreateInstance(nodeId, "Configuring", "on")
	require.NoError(t, err)
	require.NoError(t, first.Transition("configapplied"))
	require.Equal(t, "Operational", first.CurrentState)

	t.Run("returns the cached instance when it matches stored state", func(t *testing.T) {
		same, err := srv.getOrCreateInstance(nodeId, "Operational", "on")

		require.NoError(t, err)
		assert.Same(t, first, same)
	})

	t.Run("rebuilds from stored state when the cache has drifted", func(t *testing.T) {
		resynced, err := srv.getOrCreateInstance(nodeId, "Faulty", "off")

		require.NoError(t, err)
		assert.NotSame(t, first, resynced)
		assert.Equal(t, "Faulty", resynced.CurrentState)
		assert.Equal(t, "off", resynced.CurrentSubstate,
			"stored substate must survive the rebuild")

		require.NoError(t, resynced.Transition("recover"))
		assert.Equal(t, "Initializing", resynced.CurrentState,
			"resynced instance must transition from the stored state")
	})
}

func TestIsHealthEvent(t *testing.T) {
	assert.True(t, isHealthEvent(NodeStateEventPlatformReady))
	assert.True(t, isHealthEvent(NodeStateEventFault))
	assert.False(t, isHealthEvent("online"))
	assert.False(t, isHealthEvent("assign"))
	assert.False(t, isHealthEvent(""))
}

func TestShouldLatch(t *testing.T) {
	t.Run("latches health events only while the node is Unknown", func(t *testing.T) {
		assert.True(t, shouldLatch("Unknown", NodeStateEventPlatformReady))
		assert.True(t, shouldLatch("Unknown", NodeStateEventFault))
	})

	t.Run("never latches non health events", func(t *testing.T) {
		assert.False(t, shouldLatch("Unknown", "online"))
		assert.False(t, shouldLatch("Unknown", "assign"))
	})

	t.Run("never latches once the node has come up", func(t *testing.T) {
		for _, state := range []string{
			"Initializing", "Ready", "Configuring",
			"Operational", "Updating", "Faulty", "Offboarded",
		} {
			assert.False(t, shouldLatch(state, NodeStateEventPlatformReady),
				"a duplicate platformready in %s must not be latched", state)
			assert.False(t, shouldLatch(state, NodeStateEventFault),
				"a duplicate fault in %s must not be latched", state)
		}
	})
}

func TestNodeStateFsm_PlatformReadyBeforeNodeIsUpIsLatched(t *testing.T) {
	repo := &mocks.StateRepo{}
	repo.On("SetLatchedEvent", mock.Anything, mock.Anything).Return(nil).Maybe()
	repo.On("TakeLatchedEvent", mock.Anything).Return("", nil).Maybe()

	srv := &StateEventServer{
		stateMachine:  stm.NewStateMachine(func(stm.Event) {}),
		configPath:    nodeStateConfigPath,
		instances:     make(map[string]*stm.StateMachineInstance),
		latchedHealth: make(map[string]string),
		s:             NewStateServer("test-org", "test-org-id", repo, nil),
	}

	nodeId := "test-node-ready-race"

	t.Run("a latched event is stored and handed back once", func(t *testing.T) {
		srv.latchHealthEvent(nodeId, NodeStateEventPlatformReady)

		latched, ok := srv.takeLatchedHealthEvent(nodeId)
		require.True(t, ok)
		assert.Equal(t, NodeStateEventPlatformReady, latched)
	})

	t.Run("latched event is consumed only once", func(t *testing.T) {
		_, ok := srv.takeLatchedHealthEvent(nodeId)
		assert.False(t, ok, "latch must be cleared after being taken")
	})

	t.Run("replaying the latched platformready reaches Ready", func(t *testing.T) {
		instance, err := srv.getOrCreateInstance("test-node-replay", "Unknown", "on")
		require.NoError(t, err)

		require.NoError(t, instance.Transition("online"))
		require.Equal(t, "Unknown", instance.CurrentState)

		require.NoError(t, instance.Transition(NodeStateEventPlatformReady))
		assert.Equal(t, "Ready", instance.CurrentState,
			"the latched platformready replayed after the node came up must reach Ready")
	})
}

func TestNodeStateFsm_LatchSurvivesRestart(t *testing.T) {
	repo := &mocks.StateRepo{}
	repo.On("TakeLatchedEvent", "restarted-node").Return(NodeStateEventPlatformReady, nil).Once()

	srv := &StateEventServer{
		stateMachine:  stm.NewStateMachine(func(stm.Event) {}),
		configPath:    nodeStateConfigPath,
		instances:     make(map[string]*stm.StateMachineInstance),
		latchedHealth: make(map[string]string),
		s:             NewStateServer("test-org", "test-org-id", repo, nil),
	}

	latched, ok := srv.takeLatchedHealthEvent("restarted-node")

	require.True(t, ok, "a latch persisted before the restart must still be found")
	assert.Equal(t, NodeStateEventPlatformReady, latched)
	repo.AssertExpectations(t)
}
