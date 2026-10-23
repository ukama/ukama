/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package server

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	mbmocks "github.com/ukama/ukama/systems/common/mocks"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	npb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/node/state/pkg/db"
	"github.com/ukama/ukama/systems/node/state/pkg/lifecycle"
)

func TestAssignmentLifecycle(t *testing.T) {
	n := NewStateEventServer("org", "", nil, nodeStateConfigPath, nil)
	nodeID := ukama.NewVirtualHomeNodeId().String()
	record := &db.LifecycleRecord{NodeID: nodeID}
	state := &db.State{NodeId: nodeID, SubState: db.StringArray{"off"}}
	observe := func(boot, value string, seq uint64, request string) error {
		return n.observeLifecycle(record, state, &lifecycle.Event{BootID: boot, State: value,
			Sequence: seq, Time: 100, RequestID: request, ConfigMode: "NOCONFIG", Generation: 1})
	}
	require.NoError(t, observe("a", "INIT", 1, ""))
	require.Equal(t, npb.NodeState_Initializing, state.CurrentState)
	require.NoError(t, observe("a", "READY", 2, ""))
	require.Equal(t, npb.NodeState_Ready, state.CurrentState)
	require.Empty(t, record.RequestID)
	require.False(t, record.AwaitingOperational)
	msg := &epb.EventRegistryNodeAssign{NodeId: nodeID, Network: "network", Site: "site"}
	require.NoError(t, n.assignNode(record, state, msg))
	request := record.RequestID
	require.NotEmpty(t, request)
	require.True(t, record.AwaitingOperational)
	require.Equal(t, npb.NodeState_Configuring, state.CurrentState)
	require.NoError(t, n.assignNode(record, state, msg))
	require.Equal(t, request, record.RequestID)
	require.NoError(t, observe("a", "CONFIGURING", 3, request))
	require.False(t, record.Completed)
	require.NoError(t, observe("a", "OPERATIONAL", 4, request))
	require.True(t, record.Completed)
	require.False(t, record.AwaitingOperational)
	require.Equal(t, npb.NodeState_Operational, state.CurrentState)
	require.Equal(t, db.StringArray{"off"}, state.SubState, "lifecycle does not imply connectivity")
	require.NoError(t, observe("b", "INIT", 1, ""))
	require.True(t, record.AwaitingOperational)
	require.True(t, record.RetryAt.After(time.Now().Add(59*time.Second)))
	require.NoError(t, observe("b", "READY", 2, ""))
	require.NoError(t, observe("b", "CONFIGURING", 3, request))
	require.NoError(t, observe("b", "OPERATIONAL", 4, request))
	require.NoError(t, observe("a", "READY", 10, ""))
	require.Equal(t, npb.NodeState_Operational, state.CurrentState, "old boot cannot rewind state")
}

func TestReleaseClearsAssignmentInAnyState(t *testing.T) {
	n := NewStateEventServer("org", "", nil, nodeStateConfigPath, nil)
	assigned := func() *db.LifecycleRecord {
		return &db.LifecycleRecord{NodeID: "node-1", Network: "network", Site: "site", RequestID: "request-1",
			ProvisionManaged: true, Completed: true, AwaitingOperational: true}
	}
	release := &epb.NodeReleasedEvent{NodeId: "node-1", Network: "network", Site: "site"}

	for current, want := range map[npb.NodeState]npb.NodeState{
		npb.NodeState_Configuring:  npb.NodeState_Ready,
		npb.NodeState_Operational:  npb.NodeState_Ready,
		npb.NodeState_Initializing: npb.NodeState_Initializing,
		npb.NodeState_Updating:     npb.NodeState_Updating,
		npb.NodeState_Faulty:       npb.NodeState_Faulty,
		npb.NodeState_Ready:        npb.NodeState_Ready,
	} {
		record := assigned()
		state := &db.State{CurrentState: current, SubState: db.StringArray{"on"}}
		require.NoError(t, n.releaseNode(record, state, release), current.String())
		require.Equal(t, want, state.CurrentState, current.String())
		require.Equal(t, db.StringArray{"on"}, state.SubState, current.String())
		require.Equal(t, &db.LifecycleRecord{NodeID: "node-1"}, record, current.String())
	}
}

func TestReleaseWithoutAssignmentStillLeavesOperational(t *testing.T) {
	n := NewStateEventServer("org", "", nil, nodeStateConfigPath, nil)
	record := &db.LifecycleRecord{NodeID: "node-1"}
	state := &db.State{CurrentState: npb.NodeState_Operational}
	require.NoError(t, n.releaseNode(record, state,
		&epb.NodeReleasedEvent{NodeId: "node-1", Network: "network", Site: "site"}))
	require.Equal(t, npb.NodeState_Ready, state.CurrentState)
}

func TestReleaseForAnotherSiteIsIgnored(t *testing.T) {
	n := NewStateEventServer("org", "", nil, nodeStateConfigPath, nil)
	record := &db.LifecycleRecord{NodeID: "node-1", Network: "network", Site: "site-2", RequestID: "request-2", Completed: true}
	state := &db.State{CurrentState: npb.NodeState_Operational}
	require.NoError(t, n.releaseNode(record, state,
		&epb.NodeReleasedEvent{NodeId: "node-1", Network: "network", Site: "site-1"}))
	require.Equal(t, npb.NodeState_Operational, state.CurrentState)
	require.Equal(t, "request-2", record.RequestID)
	require.Equal(t, "site-2", record.Site)
}

// registry/node parses these values, so they must stay registry node state names and on/off.
func TestTransitionEventContract(t *testing.T) {
	for _, name := range npb.NodeState_name {
		require.True(t, ukama.IsValidNodeState(name), name)
	}

	var config struct {
		States []struct {
			Substate struct {
				Transition []struct {
					ToState string `json:"to_state"`
				} `json:"transition"`
			} `json:"substate"`
		} `json:"states"`
	}
	data, err := os.ReadFile(nodeStateConfigPath)
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(data, &config))
	for _, state := range config.States {
		for _, transition := range state.Substate.Transition {
			require.Contains(t, []string{"on", "off"}, transition.ToState)
		}
	}

	bus := &mbmocks.MsgBusServiceClient{}
	bus.On("PublishRequest", "event.cloud.local.org.node.state.node.transition",
		mock.MatchedBy(func(value interface{}) bool {
			event, ok := value.(*epb.NodeStateChangeEvent)
			return ok && event.NodeId == "node-1" && event.State == "Operational" &&
				event.Substate == "on" && len(event.Events) == 1 && event.Events[0] == "operational"
		})).Return(nil).Once()
	n := NewStateEventServer("org", "", nil, nodeStateConfigPath, bus)
	require.NoError(t, n.publishLifecycle(&db.LifecyclePublication{
		NodeID: "node-1", State: npb.NodeState_Operational.String(), Substate: "on", Event: "operational",
	}))
	bus.AssertExpectations(t)
}

func TestNoConfigWireContract(t *testing.T) {
	bus := &mbmocks.MsgBusServiceClient{}
	bus.On("PublishRequest", "request.cloud.local.org.node.state.nodefeeder.publish", mock.MatchedBy(func(value interface{}) bool {
		message, ok := value.(*epb.NodeFeederMessage)
		if !ok {
			return false
		}
		var command map[string]string
		err := json.Unmarshal(message.Msg, &command)
		return err == nil && message.Target == "org.network.site.node-1" &&
			message.Path == "configd/v1/config" && message.HttpMethod == "POST" &&
			command["mode"] == "NOCONFIG" && command["requestId"] == "assignment-1"
	})).Return(nil).Once()
	n := NewStateEventServer("org", "", nil, nodeStateConfigPath, bus)
	n.sendNoConfig(&db.LifecycleRecord{NodeID: "node-1", Network: "network", Site: "site", RequestID: "assignment-1"})
	bus.AssertExpectations(t)
}
