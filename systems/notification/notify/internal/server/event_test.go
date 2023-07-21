package server_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/types/known/anypb"

	mbmocks "github.com/ukama/ukama/systems/common/mocks"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/notification/notify/internal/server"
	"github.com/ukama/ukama/systems/notification/notify/mocks"
)

func TestNotifyEventServer_HandleNotificationSentEvent(t *testing.T) {
	msgbusClient := &mbmocks.MsgBusServiceClient{}
	msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()
	repo := mocks.NotificationRepo{}

	node := ukama.NewVirtualHomeNodeId().String()
	nt := NewTestDbNotification(node, "alert")

	listenerRoutes := []string{"event.cloud.org.notification.sent",
		"event.cloud.node.notification.sent"}

	t.Run("NotificationEventSent", func(t *testing.T) {
		routingKey := "event.cloud.node.notification.sent"

		repo.On("Add", mock.Anything).Return(nil)

		evt := &epb.Notification{
			Id:          nt.Id.String(),
			NodeId:      nt.NodeId,
			NodeType:    nt.NodeType,
			Severity:    nt.Severity.String(),
			Type:        nt.Type.String(),
			ServiceName: nt.ServiceName,
			EpochTime:   nt.Time,
			Description: nt.Description,
			Details:     nt.Details.String(),
		}

		anyE, err := anypb.New(evt)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s := server.NewNotifyEventServer(&repo, msgbusClient, listenerRoutes)

		_, err = s.EventNotification(context.TODO(), msg)

		assert.NoError(t, err)
	})

	t.Run("InvalidNotificationEventSent", func(t *testing.T) {
		routingKey := "event.cloud.node.notification.sent"

		evt := &epb.Notification{
			Id:          nt.Id.String(),
			NodeId:      "foo",
			NodeType:    nt.NodeType,
			Severity:    nt.Severity.String(),
			Type:        "bar",
			ServiceName: nt.ServiceName,
			EpochTime:   nt.Time,
			Description: nt.Description,
			Details:     nt.Details.String(),
		}

		anyE, err := anypb.New(evt)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s := server.NewNotifyEventServer(&repo, msgbusClient, listenerRoutes)

		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})

	t.Run("NotificationEventNotSent", func(t *testing.T) {
		routingKey := "event.cloud.cdr.sim.usage"

		evt := epb.SimUsage{}

		anyE, err := anypb.New(&evt)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s := server.NewNotifyEventServer(&repo, msgbusClient, listenerRoutes)

		_, err = s.EventNotification(context.TODO(), msg)

		assert.NoError(t, err)
	})
}
