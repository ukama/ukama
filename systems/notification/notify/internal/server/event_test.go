package server_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	mbmocks "github.com/ukama/ukama/systems/common/mocks"
	"github.com/ukama/ukama/systems/common/msgbus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/notification/notify/internal/server"
	"github.com/ukama/ukama/systems/notification/notify/mocks"
	"google.golang.org/protobuf/types/known/anypb"
)

const OrgName = "testorg"

func TestNotifyEventServer_HandleNotificationSentEvent(t *testing.T) {
	msgbusClient := &mbmocks.MsgBusServiceClient{}
	msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()
	repo := mocks.NotificationRepo{}

	node := ukama.NewVirtualHomeNodeId().String()
	nt := NewTestDbNotification(node, "alert")

	listenerRoutes := msgbus.PrepareRoutes(OrgName, []string{"event.cloud.local.{{ .Org}}.subscriber.cdr.sim.usage",
		"event.cloud.local.{{ .Org}}.registry.node.notification.sent"})

	t.Run("NotificationEventSent", func(t *testing.T) {
		routingKey := msgbus.PrepareRoute(OrgName, "event.cloud.local.{{ .Org}}.registry.node.notification.sent")

		repo.On("Add", mock.Anything).Return(nil)

		evt := &epb.Notification{
			Id:          nt.Id.String(),
			NodeId:      nt.NodeId,
			NodeType:    nt.NodeType,
			Severity:    nt.Severity.String(),
			Type:        nt.Type.String(),
			ServiceName: nt.ServiceName,
			Status:      nt.Status,
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

		s := server.NewNotifyEventServer(OrgName, &repo, msgbusClient, listenerRoutes)

		_, err = s.EventNotification(context.TODO(), msg)

		assert.NoError(t, err)
	})

	t.Run("InvalidNotificationEventSent", func(t *testing.T) {
		routingKey := msgbus.PrepareRoute(OrgName, "event.cloud.local.{{ .Org}}.registry.node.notification.sent")

		evt := &epb.Notification{
			Id:          nt.Id.String(),
			NodeId:      "foo",
			NodeType:    nt.NodeType,
			Severity:    nt.Severity.String(),
			Type:        "bar",
			ServiceName: nt.ServiceName,
			Status:      nt.Status,
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

		s := server.NewNotifyEventServer(OrgName, &repo, msgbusClient, listenerRoutes)

		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})

	t.Run("NotificationEventNotSent", func(t *testing.T) {
		routingKey := msgbus.PrepareRoute(OrgName, "event.cloud.local.{{ .Org}}.subscriber.cdr.sim.usage")

		evt := epb.SimUsage{}

		anyE, err := anypb.New(&evt)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s := server.NewNotifyEventServer(OrgName, &repo, msgbusClient, listenerRoutes)

		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})
}
