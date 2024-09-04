package server

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/types/known/anypb"

	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	ukama "github.com/ukama/ukama/systems/common/ukama"
	mocks "github.com/ukama/ukama/systems/node/state/mocks"
	utils "github.com/ukama/ukama/systems/node/state/pkg/utils"

	db "github.com/ukama/ukama/systems/node/state/pkg/db"
)

var testNode = ukama.NewVirtualNodeId("HomeNode")

func TestHandleNodeOnlineEvent(t *testing.T) {
	mockRepo := mocks.NewStateRepo(t)
	mockServer := &StateServer{sRepo: mockRepo}
	eventServer := NewControllerEventServer("test-org", mockServer, nil)

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

	msg, _ := anypb.New(&epb.Notification{
		NodeId:   testNode.String(),
		Type:     string(utils.Event),
		Severity: string(utils.Medium),
		Details:  []byte(`{"config": "ready"}`),
	})

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

	mockRepo.On("GetByNodeId", ukama.NodeID(testNode.String())).Return(&db.State{NodeId: testNode.String(), State: ukama.StateUnknown}, nil)
	mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

	err := eventServer.updateNodeState(testNode.String(), ukama.StateOperational, utils.Medium, utils.Event)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUnmarshalNotification(t *testing.T) {
	eventServer := &NodeStateEventServer{}

	msg, _ := anypb.New(&epb.Notification{NodeId: testNode.String(), Type: string(utils.Event)})
	evt, err := eventServer.unmarshalNotification(msg)

	assert.NoError(t, err)
	assert.NotNil(t, evt)
	assert.Equal(t, testNode.String(), evt.NodeId)
}

func TestUnmarshalNodeCreateEvent(t *testing.T) {
	eventServer := &NodeStateEventServer{}

	msg, _ := anypb.New(&epb.EventRegistryNodeCreate{NodeId: testNode.String()})
	evt, err := eventServer.unmarshalNodeCreateEvent(msg)

	assert.NoError(t, err)
	assert.NotNil(t, evt)
	assert.Equal(t, testNode.String(), evt.NodeId)
}

func TestDetermineNodeState(t *testing.T) {
	eventServer := &NodeStateEventServer{}

	tests := []struct {
		name            string
		systemStatuses  map[string]string
		bootstrapStatus string
		expectedState   ukama.NodeStateEnum
	}{
		{
			name: "All systems operational",
			systemStatuses: map[string]string{
				"CPU":    "Fine",
				"Memory": "Fine",
				"Radio":  "On",
			},
			bootstrapStatus: "Running",
			expectedState:   ukama.StateOperational,
		},
		{
			name: "Faulty CPU",
			systemStatuses: map[string]string{
				"CPU":    "Error",
				"Memory": "Fine",
				"Radio":  "On",
			},
			bootstrapStatus: "Running",
			expectedState:   ukama.StateFaulty,
		},
		{
			name: "Bootstrap not running",
			systemStatuses: map[string]string{
				"CPU":    "Fine",
				"Memory": "Fine",
				"Radio":  "On",
			},
			bootstrapStatus: "Not Found",
			expectedState:   ukama.StateFaulty,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := eventServer.determineNodeState(tt.systemStatuses, tt.bootstrapStatus)
			assert.Equal(t, tt.expectedState, state)
		})
	}
}

func TestHandleConfigState(t *testing.T) {
	eventServer := &NodeStateEventServer{}

	tests := []struct {
		name          string
		configKey     string
		currentState  ukama.NodeStateEnum
		expectedState ukama.NodeStateEnum
	}{
		{
			name:          "Ready config",
			configKey:     "ready",
			currentState:  ukama.StateUnknown,
			expectedState: ukama.StateUnknown, // Remains unchanged until threshold is met
		},
		{
			name:          "Faulty config",
			configKey:     "faulty",
			currentState:  ukama.StateOperational,
			expectedState: ukama.StateFaulty,
		},
		{
			name:          "Other config",
			configKey:     "other",
			currentState:  ukama.StateOperational,
			expectedState: ukama.StateConfigure,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := eventServer.handleConfigState(tt.configKey, testNode.String(), tt.currentState, time.Now())
			assert.Equal(t, tt.expectedState, state)
		})
	}

	// Test for multiple "ready" configs within the timeout
	t.Run("Multiple ready configs", func(t *testing.T) {
		configReadyCount = 0
		now := time.Now()
		for i := 0; i < configReadyThreshold; i++ {
			state := eventServer.handleConfigState("ready", testNode.String(), ukama.StateUnknown, now)
			if i < configReadyThreshold-1 {
				assert.Equal(t, ukama.StateUnknown, state)
			} else {
				assert.Equal(t, ukama.StateOperational, state)
			}
		}
	})
}