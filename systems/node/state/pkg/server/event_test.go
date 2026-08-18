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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	mbmocks "github.com/ukama/ukama/systems/common/mocks"
	"github.com/ukama/ukama/systems/common/msgbus"
	stm "github.com/ukama/ukama/systems/common/stateMachine"

	"github.com/ukama/ukama/systems/common/ukama"
	pb "github.com/ukama/ukama/systems/node/state/pb/gen"
	"github.com/ukama/ukama/systems/node/state/pkg"
)

type MockStateEventServer struct {
	mock.Mock
	latestStateResponse *pb.GetLatestStateResponse
}

type MockStateMachine struct {
	mock.Mock
	stm.StateMachine
}

func NewMockStateMachine() *MockStateMachine {
	return &MockStateMachine{}
}

func (m *MockStateMachine) NewInstance(configPath, instanceID, initialState string) (*stm.StateMachineInstance, error) {
	args := m.Called(configPath, instanceID, initialState)
	if instance, ok := args.Get(0).(*stm.StateMachineInstance); ok {
		return instance, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockStateEventServer) GetLatestState(ctx context.Context, req *pb.GetLatestStateRequest) (*pb.GetLatestStateResponse, error) {
	if m.latestStateResponse != nil {
		return m.latestStateResponse, nil
	}
	return nil, fmt.Errorf("no response set")
}

func (m *MockStateEventServer) SetLatestStateResponse(response *pb.GetLatestStateResponse) {
	m.latestStateResponse = response
}
func (m *MockStateEventServer) handleTransition(event stm.Event) {
    state := event.NewState
    substate := event.NewSubstate
    m.publishStateChangeEvent(state, substate, event.InstanceID)
}

func (m *MockStateEventServer) publishStateChangeEvent(state, substate, nodeID string) {
    m.Called(state, substate, nodeID)
}

func TestStateEventServer_handleTransition(t *testing.T) {
    mockServer := new(MockStateEventServer)
    nodeId := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE).String()

    testEvent := stm.Event{
        NewState:    "unknown",
        NewSubstate: "config",
        InstanceID:  nodeId,
    }

    mockServer.On("publishStateChangeEvent", testEvent.NewState, testEvent.NewSubstate, testEvent.InstanceID).Once()
    mockServer.handleTransition(testEvent)
    mockServer.AssertExpectations(t)
}
func TestStateEventServer_publishStateChangeEvent(t *testing.T) {
	mockMsgBus := new(mbmocks.MsgBusServiceClient)

	server := &StateEventServer{
		msgbus:         mockMsgBus,
		baseRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName("test-org").SetService(pkg.ServiceName),
		orgName:        "test-org",
		orgId:          "test-org-id",
		instances:      make(map[string]*stm.StateMachineInstance),
		eventBuffer:    make(map[string][]string),
	}

	nodeId := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE).String()

	mockMsgBus.On("PublishRequest", mock.Anything, mock.Anything).Return(nil)

	server.publishStateChangeEvent("unknown", "config", nodeId)

	mockMsgBus.AssertExpectations(t)
}

func TestStateEventForNotifyValue(t *testing.T) {
	cases := []struct {
		name  string
		value string
		want  string
	}{
		{"node reports READY as sent on the wire", "READY", NodeStateEventPlatformReady},
		{"node reports FAULTY as sent on the wire", "FAULTY", NodeStateEventFault},
		{"lower case ready still maps", "ready", NodeStateEventPlatformReady},
		{"mixed case faulty still maps", "Faulty", NodeStateEventFault},
		{"unmapped values pass through", "reboot", "reboot"},
		{"empty value passes through", "", ""},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, stateEventForNotifyValue(tc.value))
		})
	}
}

func TestStateEventForNotifyValue_DrivesTheStateMachine(t *testing.T) {
	t.Run("READY moves an initializing node to Ready", func(t *testing.T) {
		instance := newNodeInstance(t, "notify-ready", "Initializing")

		require.NoError(t, instance.Transition(stateEventForNotifyValue("READY")))
		assert.Equal(t, "Ready", instance.CurrentState)
	})

	t.Run("FAULTY moves an operational node to Faulty", func(t *testing.T) {
		instance := newNodeInstance(t, "notify-faulty", "Operational")

		require.NoError(t, instance.Transition(stateEventForNotifyValue("FAULTY")))
		assert.Equal(t, "Faulty", instance.CurrentState)
	})
}
