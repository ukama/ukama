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
	"encoding/json"
	"fmt"

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
	msgbus         mb.MsgBusServiceClient
	baseRoutingKey msgbus.RoutingKeyBuilder
	epb.UnimplementedEventNotificationServiceServer
}

func NewNotifyEventServer(orgName string, nRepo db.NotificationRepo, msgBus mb.MsgBusServiceClient) *NotifiyEventServer {
	return &NotifiyEventServer{
		orgName:        orgName,
		notifyRepo:     nRepo,
		msgbus:         msgBus,
		baseRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(internal.SystemName).SetOrgName(orgName).SetService(internal.ServiceName),
	}
}

func (n *NotifiyEventServer) EventNotification(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
	log.Infof("Received a message with Routing key %s and Message %+v", e.RoutingKey, e.Msg)
	switch e.RoutingKey {
	case msgbus.PrepareRoute(n.orgName, "event.cloud.local.{{ .Org}}.messaging.mesh.node.online"):
		msg, err := n.unmarshalNodeOnlineEvent(e.Msg)
		if err != nil {
			return nil, err
		}

		err = n.handleNodeOnlineEvent(msg)
		if err != nil {
			return nil, err
		}
	case msgbus.PrepareRoute(n.orgName, "event.cloud.local.{{ .Org}}.messaging.mesh.node.offline"):
		msg, err := n.unmarshalNodeOfflineEvent(e.Msg)
		if err != nil {
			return nil, err
		}

		err = n.handleNodeOfflineEvent(msg)
		if err != nil {
			return nil, err
		}
	case msgbus.PrepareRoute(n.orgName, "event.cloud.local.{{ .Org}}.registry.node.node.create"):
		msg, err := n.unmarshalNodeCreateEvent(e.Msg)
		if err != nil {
			return nil, err
		}

		err = n.handleNodeCreateEvent(msg)
		if err != nil {
			return nil, err
		}
	default:
		log.Errorf("No handler routing key %s", e.RoutingKey)
	}

	return &epb.EventResponse{}, nil
}

func (n *NotifiyEventServer) handleNodeOnlineEvent(msg *epb.NodeOnlineEvent) error {
	eventData := map[string]interface{}{
		"value": "online",
	}
	data, err := json.Marshal(eventData)
	if err != nil {
		return fmt.Errorf("failed to marshal event data: %v", err)
	}

	return add(
		msg.NodeId,
		string(db.Low),
		db.NotificationType("event").String(),
		"mesh",
		data,
		1,
		1,
		n.notifyRepo,
		n.msgbus,
		n.baseRoutingKey,
	)

}

func (n *NotifiyEventServer) handleNodeCreateEvent(msg *epb.NodeCreatedEvent) error {
	eventData := map[string]interface{}{
		"value": "created",
	}

	data, err := json.Marshal(eventData)
	if err != nil {
		return fmt.Errorf("failed to marshal event data: %v", err)
	}

	return add(
		msg.NodeId,
		string(db.Low),
		db.NotificationType("event").String(),
		"registry",
		data,
		1,
		1,
		n.notifyRepo,
		n.msgbus,
		n.baseRoutingKey,
	)
}

func (n *NotifiyEventServer) handleNodeOfflineEvent(msg *epb.NodeOfflineEvent) error {
	eventData := map[string]interface{}{
		"value": "offline",
	}

	data, err := json.Marshal(eventData)
	if err != nil {
		return fmt.Errorf("failed to marshal event data: %v", err)
	}

	return add(
		msg.NodeId,
		string(db.Low),
		db.NotificationType("event").String(),
		"mesh",
		data,
		1,
		1,
		n.notifyRepo,
		n.msgbus,
		n.baseRoutingKey,
	)
}

func (n *NotifiyEventServer) unmarshalNodeCreateEvent(msg *anypb.Any) (*epb.NodeCreatedEvent, error) {
	p := &epb.NodeCreatedEvent{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal NodeCreate message with : %+v. Error %s.", msg, err.Error())
		return nil, err
	}

	return p, nil
}

func (n *NotifiyEventServer) unmarshalNodeOnlineEvent(msg *anypb.Any) (*epb.NodeOnlineEvent, error) {
	p := &epb.NodeOnlineEvent{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal NodeOnline  message with : %+v. Error %s.", msg, err.Error())
		return nil, err
	}
	return p, nil
}
func (n *NotifiyEventServer) unmarshalNodeOfflineEvent(msg *anypb.Any) (*epb.NodeOfflineEvent, error) {
	p := &epb.NodeOfflineEvent{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal NodeOffline message with : %+v. Error %s.", msg, err.Error())
		return nil, err
	}
	return p, nil
}
