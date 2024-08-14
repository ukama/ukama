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
	"time"

	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/node/notify/internal"
	"github.com/ukama/ukama/systems/node/notify/internal/db"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	log "github.com/sirupsen/logrus"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
)

type NotifiyEventServer struct {
	orgName        string
	notifyRepo     db.NotificationRepo
	listenerRoutes map[string]struct{}
	msgbus         mb.MsgBusServiceClient
	baseRoutingKey msgbus.RoutingKeyBuilder
	epb.UnimplementedEventNotificationServiceServer
}

func NewNotifyEventServer(orgName string, nRepo db.NotificationRepo, msgBus mb.MsgBusServiceClient, routes []string) *NotifiyEventServer {

    pRoutes := msgbus.PrepareRoutes(orgName, routes)
    r := make(map[string]struct{}, len(routes))

    for _, route := range pRoutes {
        r[route] = struct{}{}
    }

    nodeRoutes := []string{
        "event.cloud.local." + orgName + ".registry.node.node.create",
        "event.cloud.local." + orgName + ".messaging.mesh.node.online",
        "event.cloud.local." + orgName + ".messaging.mesh.node.offline",
    }

    for _, route := range nodeRoutes {
        r[route] = struct{}{}
    }

    return &NotifiyEventServer{
        orgName:        orgName,
        notifyRepo:     nRepo,
        listenerRoutes: r,
        msgbus:         msgBus,
        baseRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(internal.SystemName).SetOrgName(orgName).SetService(internal.ServiceName),
    }
}


func (n *NotifiyEventServer) EventNotification(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
    log.Infof("Received a message with Routing key %s and Message %+v", e.RoutingKey, e.Msg)

    if _, ok := n.listenerRoutes[e.RoutingKey]; !ok {
        log.Errorf("No handler routing key %s", e.RoutingKey)
        return nil, nil
    }

    switch e.RoutingKey {
    case "event.cloud.local." + n.orgName + ".registry.node.node.create":
        log.Infof("Handling node.create event")
		msg, err := n.unmarshalNodeCreateEvent(e.Msg)
		if err != nil {
			return nil, err
		}
		err = n.handleNodeCreate(msg)
		if err != nil {
			return nil, err
		}

    case "event.cloud.local." + n.orgName + ".messaging.mesh.node.online":
        log.Infof("Handling node.online event")
		msg, err := n.unmarshalNotificationSentEvent(e.Msg)
        if err != nil {
            return nil, err
        }
        err = n.handleNotificationSentEvent(e.RoutingKey, msg)
        if err != nil {
            return nil, err
        }

    case "event.cloud.local." + n.orgName + ".messaging.mesh.node.offline":
        log.Infof("Handling node.offline event")
		msg, err := n.unmarshalNotificationSentEvent(e.Msg)
        if err != nil {
            return nil, err
        }
        err = n.handleNotificationSentEvent(e.RoutingKey, msg)
        if err != nil {
            return nil, err
        }

    default:
        msg, err := n.unmarshalNotificationSentEvent(e.Msg)
        if err != nil {
            return nil, err
        }
        err = n.handleNotificationSentEvent(e.RoutingKey, msg)
        if err != nil {
            return nil, err
        }
    }

    return &epb.EventResponse{}, nil
}


func (n *NotifiyEventServer) unmarshalNotificationSentEvent(msg *anypb.Any) (*epb.Notification, error) {
	p := &epb.Notification{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal notification created message with : %+v. Error %s.", msg, err.Error())

		return nil, err
	}
	return p, nil
}

func (n *NotifiyEventServer) handleNotificationSentEvent(key string, msg *epb.Notification) error {
	return add(msg.NodeId, msg.Severity, msg.Type, msg.ServiceName, msg.Description,
		msg.Details, msg.Status, msg.EpochTime, n.notifyRepo, n.msgbus, n.baseRoutingKey)
}

func (n *NotifiyEventServer) unmarshalNodeCreateEvent(msg *anypb.Any) (*epb.EventRegistryNodeCreate, error) {
	evt := &epb.EventRegistryNodeCreate{}
	err := anypb.UnmarshalTo(msg, evt, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal NodeCreate event with : %+v. Error %s.", msg, err.Error())
		return nil, err
	}
	return evt, nil
}



func (n *NotifiyEventServer) handleNodeCreate(evt *epb.EventRegistryNodeCreate) error {
	log.Infof("NodeCreate event received: NodeId: %s, Name: %s, Type: %s", evt.NodeId, evt.Name, evt.Type)
	
	return add(evt.NodeId, string(db.Medium), evt.Type, "node", "node creation",
		"", 1, uint32(time.Now().Unix()), n.notifyRepo, n.msgbus, n.baseRoutingKey)
}
