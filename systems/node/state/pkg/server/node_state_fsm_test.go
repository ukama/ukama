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

func TestNodeStateFsm_FaultyRecovery(t *testing.T) {
	t.Run("faulty node returns to operational on ready", func(t *testing.T) {
		instance := newNodeInstance(t, "faulty-ready", "Faulty")

		require.NoError(t, instance.Transition("ready"))
		assert.Equal(t, "Operational", instance.CurrentState)
	})

	t.Run("operational to faulty and back again", func(t *testing.T) {
		instance := newNodeInstance(t, "fault-round-trip", "Operational")

		require.NoError(t, instance.Transition("fault"))
		assert.Equal(t, "Faulty", instance.CurrentState)

		require.NoError(t, instance.Transition("ready"))
		assert.Equal(t, "Operational", instance.CurrentState)
	})

	t.Run("faulty node still reaches configured on config", func(t *testing.T) {
		instance := newNodeInstance(t, "faulty-config", "Faulty")

		require.NoError(t, instance.Transition("config"))
		assert.Equal(t, "Configured", instance.CurrentState)
	})

	t.Run("faulty node reaches configured on upgrade and downgrade", func(t *testing.T) {
		for _, trigger := range []string{"upgrade", "downgrade"} {
			instance := newNodeInstance(t, "faulty-"+trigger, "Faulty")

			require.NoError(t, instance.Transition(trigger))
			assert.Equal(t, "Configured", instance.CurrentState,
				"%s must recover a faulty node to Configured", trigger)
		}
	})
}

func TestNodeStateFsm_RepeatReadyInOperational(t *testing.T) {
	instance := newNodeInstance(t, "operational-repeat-ready", "Operational")

	require.NoError(t, instance.Transition("online"))
	assert.Equal(t, "on", instance.CurrentSubstate)

	err := instance.Transition("ready")

	assert.NoError(t, err)
	assert.Equal(t, "Operational", instance.CurrentState)
	assert.Equal(t, "on", instance.CurrentSubstate)
}

func TestNodeStateFsm_ConfiguredReadyKeepsSubstate(t *testing.T) {
	instance := newNodeInstance(t, "configured-ready", "Configured")

	require.NoError(t, instance.Transition("online"))
	require.Equal(t, "on", instance.CurrentSubstate)

	require.NoError(t, instance.Transition("ready"))

	assert.Equal(t, "Operational", instance.CurrentState)
	assert.Equal(t, "on", instance.CurrentSubstate,
		"ready must not move the substate into update")
}

func TestNodeStateFsm_OnboardingFromUnknown(t *testing.T) {
	instance := newNodeInstance(t, "unknown-onboarding", "Unknown")

	require.NoError(t, instance.Transition("onboarding"))
	assert.Equal(t, "Configured", instance.CurrentState)

	require.NoError(t, instance.Transition("ready"))
	assert.Equal(t, "Operational", instance.CurrentState)
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

	require.NoError(t, instance.Transition("ready"))
	assert.Equal(t, 1, published, "repeat ready must not publish a state change")

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

	first, err := srv.getOrCreateInstance(nodeId, "Configured", "on")
	require.NoError(t, err)
	require.NoError(t, first.Transition("ready"))
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

		require.NoError(t, resynced.Transition("ready"))
		assert.Equal(t, "Operational", resynced.CurrentState,
			"resynced instance must transition from the stored state")
	})
}

func TestNodeStateFsm_ReadyBeforeOnboardingIsLatched(t *testing.T) {
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

	t.Run("ready in Unknown is latched rather than dropped", func(t *testing.T) {
		instance, err := srv.getOrCreateInstance(nodeId, "Unknown", "on")
		require.NoError(t, err)

		require.NoError(t, instance.Transition(NodeNotifyEventReady))
		assert.Equal(t, "Unknown", instance.CurrentState,
			"ready must not move a node out of Unknown")

		srv.latchHealthEvent(nodeId, NodeNotifyEventReady)

		latched, ok := srv.takeLatchedHealthEvent(nodeId)
		require.True(t, ok)
		assert.Equal(t, NodeNotifyEventReady, latched)
	})

	t.Run("latched ready is consumed only once", func(t *testing.T) {
		_, ok := srv.takeLatchedHealthEvent(nodeId)
		assert.False(t, ok, "latch must be cleared after being taken")
	})

	t.Run("replaying the latched ready after onboarding reaches Operational", func(t *testing.T) {
		instance, err := srv.getOrCreateInstance("test-node-replay", "Unknown", "on")
		require.NoError(t, err)

		require.NoError(t, instance.Transition("onboarding"))
		require.Equal(t, "Configured", instance.CurrentState)

		require.NoError(t, instance.Transition(NodeNotifyEventReady))
		assert.Equal(t, "Operational", instance.CurrentState,
			"the latched ready replayed after onboarding must reach Operational")
	})
}

func TestIsHealthEvent(t *testing.T) {
	assert.True(t, isHealthEvent(NodeNotifyEventReady))
	assert.True(t, isHealthEvent(NodeNotifyEventFault))
	assert.False(t, isHealthEvent("online"))
	assert.False(t, isHealthEvent("onboarding"))
}
