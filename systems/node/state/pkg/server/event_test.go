package server

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ukama/ukama/systems/common/msgbus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	eCfgPb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
	"github.com/ukama/ukama/systems/node/state/mocks"
	"github.com/ukama/ukama/systems/node/state/pkg/db"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

func TestEventNotification(t *testing.T) {
	mockRepo := mocks.NewStateRepo(t)
	server := &StateServer{sRepo: mockRepo}
	eventServer := NewControllerEventServer("testOrg", server)

	tests := []struct {
		name       string
		routingKey string
		msg        proto.Message
		setupMock  func()
		wantErr    bool
	}{
		{
			name:       "Handle Registry Node Add Event",
			routingKey: msgbus.PrepareRoute("testOrg", "event.cloud.local.{{ .Org}}.registry.node.node.create"),
			msg: &epb.NodeCreatedEvent{
				NodeId: "uk-sa2434-hnode-v0-cdb2",
				Type:   "test-type",
			},
			setupMock: func() {
				mockRepo.On("Create", mock.AnythingOfType("*db.State"), mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name:       "Handle Node Online Event",
			routingKey: msgbus.PrepareRoute("testOrg", "event.cloud.local.{{ .Org}}.messaging.mesh.node.online"),
			msg: &epb.NodeOnlineEvent{
				NodeId: "uk-sa2434-hnode-v0-cdb2",
			},
			setupMock: func() {
				mockRepo.On("GetByNodeId", mock.AnythingOfType("ukama.NodeID")).Return(&db.State{}, nil)
				mockRepo.On("Update", mock.AnythingOfType("*db.State")).Return(nil)
			},
			wantErr: false,
		},
		// Add more test cases for other event types
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			anyMsg, err := anypb.New(tt.msg)
			assert.NoError(t, err)

			event := &epb.Event{
				RoutingKey: tt.routingKey,
				Msg:        anyMsg,
			}

			_, err = eventServer.EventNotification(context.Background(), event)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestHandleRegistryNodeAddEvent(t *testing.T) {
	mockRepo := mocks.NewStateRepo(t)
	server := &StateServer{sRepo: mockRepo}
	eventServer := NewControllerEventServer("testOrg", server)

	msg := &epb.NodeCreatedEvent{
		NodeId: "uk-sa2434-hnode-v0-cdb2",
		Type:   "test-type",
	}

	mockRepo.On("Create", mock.AnythingOfType("*db.State"), mock.Anything).Return(nil)

	err := eventServer.handleRegistryNodeAddEvent("test-key", msg)
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestHandleNodeOnlineEvent(t *testing.T) {
	mockRepo := mocks.NewStateRepo(t)
	server := &StateServer{sRepo: mockRepo}
	eventServer := NewControllerEventServer("testOrg", server)

	msg := &epb.NodeOnlineEvent{
		NodeId: "uk-sa2434-hnode-v0-cdb2",
	}

	mockRepo.On("GetByNodeId", mock.AnythingOfType("ukama.NodeID")).Return(&db.State{}, nil)
	mockRepo.On("Update", mock.AnythingOfType("*db.State")).Return(nil)

	err := eventServer.handleNodeOnlineEvent("test-key", msg)
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestHandleNodeOfflineEvent(t *testing.T) {
	mockRepo := mocks.NewStateRepo(t)
	server := &StateServer{sRepo: mockRepo}
	eventServer := NewControllerEventServer("testOrg", server)

	msg := &epb.NodeOfflineEvent{
		NodeId: "uk-sa2434-hnode-v0-cdb2",
	}

	mockRepo.On("GetByNodeId", mock.AnythingOfType("ukama.NodeID")).Return(&db.State{}, nil)
	mockRepo.On("Update", mock.AnythingOfType("*db.State")).Return(nil)

	err := eventServer.handleNodeOfflineEvent("test-key", msg)
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestHandleNodeConfigUpdateEvent(t *testing.T) {
	mockRepo := mocks.NewStateRepo(t)
	server := &StateServer{sRepo: mockRepo}
	eventServer := NewControllerEventServer("testOrg", server)

	msg := &eCfgPb.NodeConfigUpdateEvent{
		NodeId: "uk-sa2434-hnode-v0-cdb2",
		Commit: "test-commit",
	}

	mockRepo.On("GetByNodeId", mock.AnythingOfType("ukama.NodeID")).Return(&db.State{CurrentState: db.StateOnboarded}, nil)
	mockRepo.On("Update", mock.AnythingOfType("*db.State")).Return(nil)

	err := eventServer.handleNodeConfigUpdateEvent("test-key", msg)
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}