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
	"google.golang.org/protobuf/types/known/anypb"

	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	ukama "github.com/ukama/ukama/systems/common/ukama"
	mocks "github.com/ukama/ukama/systems/node/state/mocks"
	db "github.com/ukama/ukama/systems/node/state/pkg/db"
)

var testNode = ukama.NewVirtualNodeId("HomeNode")

func TestHandleNodeOnlineEvent(t *testing.T) {
	mockRepo := mocks.NewStateRepo(t)
	mockServer := &StateServer{sRepo: mockRepo}
	eventServer := NewControllerEventServer("test-org", mockServer, nil)

	// Setup the mock expectation
	msg, _ := anypb.New(&epb.NodeOnlineEvent{NodeId: testNode.String()})
	mockRepo.On("GetByNodeId", ukama.NodeID(testNode.String())).Return(&db.State{NodeId: testNode.String(), State: ukama.StateUnknown}, nil)
	mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

	err := eventServer.handleNodeOnlineEvent(msg)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestHandleNodeHealthSeverityEvent(t *testing.T) {
	mockRepo := mocks.NewStateRepo(t)
	mockServer := &StateServer{sRepo: mockRepo}
	eventServer := NewControllerEventServer("test-org", mockServer, nil)

	// Simulate a health severity event with medium severity
	msg, _ := anypb.New(&epb.Notification{NodeId: testNode.String(), Type: "event", Severity: "medium"})
	mockRepo.On("GetByNodeId", ukama.NodeID(testNode.String())).Return(&db.State{NodeId: testNode.String(), State: ukama.StateUnknown}, nil)
	mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

	err := eventServer.handleNodeHealthSeverityEvent(msg)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUpdateNodeState(t *testing.T) {
	mockRepo := mocks.NewStateRepo(t)
	mockServer := &StateServer{sRepo: mockRepo}
	eventServer := NewControllerEventServer("test-org", mockServer, nil)

	// Setup existing state
	mockRepo.On("GetByNodeId", ukama.NodeID(testNode.String())).Return(&db.State{NodeId: testNode.String(), State: ukama.StateUnknown}, nil)

	// Mock Create method
	mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

	err := eventServer.updateNodeState(testNode.String(), ukama.StateOperational)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUnmarshalNotification(t *testing.T) {
	eventServer := &NodeStateEventServer{}

	// Create a mock notification message
	msg, _ := anypb.New(&epb.Notification{NodeId: testNode.String(), Type: "event"})
	evt, err := eventServer.unmarshalNotification(msg)

	assert.NoError(t, err)
	assert.NotNil(t, evt)
	assert.Equal(t, testNode.String(), evt.NodeId)
}

func TestUnmarshalNodeCreateEvent(t *testing.T) {
	eventServer := &NodeStateEventServer{}

	// Create a mock node create event message
	msg, _ := anypb.New(&epb.EventRegistryNodeCreate{NodeId: testNode.String()})
	evt, err := eventServer.unmarshalNodeCreateEvent(msg)

	assert.NoError(t, err)
	assert.NotNil(t, evt)
	assert.Equal(t, testNode.String(), evt.NodeId)
}
