package server

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	mbmocks "github.com/ukama/ukama/systems/common/mocks"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/node/notify/mocks" // Replace with your mock package path
	"google.golang.org/protobuf/types/known/anypb"
)

func TestEventNotification(t *testing.T) {
	msgbusClient := &mbmocks.MsgBusServiceClient{}

	repo := new(mocks.NotificationRepo)

	orgName := "testorg"
	notifyEventServer := NewNotifyEventServer(orgName, repo, msgbusClient)

	nodeId := ukama.NewVirtualHomeNodeId().String()

	onlineEvent := &epb.NodeOnlineEvent{
		NodeId: nodeId,
	}
	onlineMsg, _ := anypb.New(onlineEvent)

	repo.On("Add", mock.Anything).Return(nil).Once()
	msgbusClient.On("PublishRequest", mock.Anything, mock.AnythingOfType("*events.Notification")).
		Run(func(args mock.Arguments) {
			notification := args.Get(1).(*epb.Notification)
			assert.NotEmpty(t, notification.Id)
			assert.Equal(t, nodeId, notification.NodeId)
			assert.Equal(t, "hnode", notification.NodeType)
			assert.Equal(t, "low", notification.Severity)
			assert.Equal(t, "event", notification.Type)
			assert.Equal(t, "mesh", notification.ServiceName)
		}).
		Return(nil).Once()

	eventOnline := &epb.Event{
		RoutingKey: "event.cloud.local.testorg.messaging.mesh.node.online",
		Msg:        onlineMsg,
	}
	_, err := notifyEventServer.EventNotification(context.Background(), eventOnline)
	assert.NoError(t, err)

	offlineEvent := &epb.NodeOfflineEvent{
		NodeId: nodeId,
	}
	offlineMsg, _ := anypb.New(offlineEvent)

	repo.On("Add", mock.Anything).Return(nil).Once()

	msgbusClient.On("PublishRequest", mock.Anything, mock.AnythingOfType("*events.Notification")).
		Run(func(args mock.Arguments) {
			notification := args.Get(1).(*epb.Notification)
			assert.NotEmpty(t, notification.Id)
			assert.Equal(t, nodeId, notification.NodeId)
			assert.Equal(t, "hnode", notification.NodeType)
			assert.Equal(t, "low", notification.Severity)
			assert.Equal(t, "event", notification.Type)
			assert.Equal(t, "mesh", notification.ServiceName)
		}).
		Return(nil).Once()

	eventOffline := &epb.Event{
		RoutingKey: "event.cloud.local.testorg.messaging.mesh.node.offline",
		Msg:        offlineMsg,
	}
	_, err = notifyEventServer.EventNotification(context.Background(), eventOffline)
	assert.NoError(t, err)

	createdEvent := &epb.EventRegistryNodeCreate{
		NodeId: nodeId,
	}
	createdMsg, _ := anypb.New(createdEvent)

	repo.On("Add", mock.Anything).Return(nil).Once()

	msgbusClient.On("PublishRequest", mock.Anything, mock.AnythingOfType("*events.Notification")).
		Run(func(args mock.Arguments) {
			notification := args.Get(1).(*epb.Notification)
			assert.NotEmpty(t, notification.Id)
			assert.Equal(t, nodeId, notification.NodeId)
			assert.Equal(t, "hnode", notification.NodeType)
			assert.Equal(t, "low", notification.Severity)
			assert.Equal(t, "event", notification.Type)
			assert.Equal(t, "registry", notification.ServiceName)
		}).
		Return(nil).Once()

	eventCreated := &epb.Event{
		RoutingKey: "event.cloud.local.testorg.registry.node.node.create",
		Msg:        createdMsg,
	}
	_, err = notifyEventServer.EventNotification(context.Background(), eventCreated)
	assert.NoError(t, err)

	repo.AssertExpectations(t)
	msgbusClient.AssertExpectations(t)
}
