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

	"github.com/stretchr/testify/mock"
	"github.com/tj/assert"
	mbmocks "github.com/ukama/ukama/systems/common/mocks"
	"github.com/ukama/ukama/systems/common/msgbus"
	stm "github.com/ukama/ukama/systems/common/stateMachine"

	"github.com/ukama/ukama/systems/common/ukama"
	pb "github.com/ukama/ukama/systems/node/state/pb/gen"
	"github.com/ukama/ukama/systems/node/state/pkg"
)

type MockStateEventServer struct {
	mock.Mock
	msgbus *mbmocks.MsgBusServiceClient
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
	var state, substate string
	if event.IsSubstate {
		state = event.NewState
		substate = event.NewSubstate
	} else {
		state = event.NewState
		substate = event.NewSubstate
	}

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
		IsSubstate:  true,
	}

	mockServer.On("publishStateChangeEvent", "unknown", "config", nodeId).Once()
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
func TestStateEventServer_EventBuffer(t *testing.T) {
	server := &StateEventServer{
		eventBuffer: make(map[string][]string),
	}
	
	nodeId := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE).String()
	
	t.Run("Empty Buffer", func(t *testing.T) {
		events := server.getEventsForNode(nodeId)
		assert.Empty(t, events, "Event buffer should be empty initially")
	})
	
	t.Run("Add Events", func(t *testing.T) {
		// Add single event
		server.addEventToBuffer(nodeId, "event1")
		events := server.getEventsForNode(nodeId)
		assert.Len(t, events, 1, "Should have one event")
		assert.Equal(t, "event1", events[0])
		
		// Add another event
		server.addEventToBuffer(nodeId, "event2")
		events = server.getEventsForNode(nodeId)
		assert.Len(t, events, 2, "Should have two events")
		assert.Equal(t, []string{"event1", "event2"}, events)
	})
	
	t.Run("Clear Events", func(t *testing.T) {
		server.clearEventsForNode(nodeId)
		events := server.getEventsForNode(nodeId)
		assert.Empty(t, events, "Event buffer should be empty after clearing")
	})
	
	t.Run("Multiple Nodes", func(t *testing.T) {
		nodeId2 := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE).String()
		
		server.addEventToBuffer(nodeId, "event1")
		server.addEventToBuffer(nodeId2, "event2")
		
		events1 := server.getEventsForNode(nodeId)
		events2 := server.getEventsForNode(nodeId2)
		
		assert.Len(t, events1, 1, "Node 1 should have one event")
		assert.Len(t, events2, 1, "Node 2 should have one event")
		assert.Equal(t, "event1", events1[0])
		assert.Equal(t, "event2", events2[0])
		
		server.clearEventsForNode(nodeId)
		assert.Empty(t, server.getEventsForNode(nodeId), "Node 1 events should be cleared")
		assert.NotEmpty(t, server.getEventsForNode(nodeId2), "Node 2 events should remain")
	})
}


