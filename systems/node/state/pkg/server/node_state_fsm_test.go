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
	"github.com/stretchr/testify/require"

	stm "github.com/ukama/ukama/systems/common/stateMachine"
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
