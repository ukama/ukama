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
	"strings"

	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/node/notify/internal"
	"github.com/ukama/ukama/systems/node/notify/internal/db"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	log "github.com/sirupsen/logrus"
	evt "github.com/ukama/ukama/systems/common/events"
	cpb "github.com/ukama/ukama/systems/common/pb/gen/ukama"

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

	case msgbus.PrepareRoute(n.orgName, evt.EventRoutingKey[evt.EventNodeOnline]):
		c := evt.EventToEventConfig[evt.EventNodeOnline]
		msg, err := epb.UnmarshalNodeOnlineEvent(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}

		err = n.handleNodeOnlineEvent(msg, c.Title)
		if err != nil {
			return nil, err
		}
	case msgbus.PrepareRoute(n.orgName, evt.EventRoutingKey[evt.EventNodeOffline]):
		c := evt.EventToEventConfig[evt.EventNodeOffline]
		msg, err := epb.UnmarshalNodeOfflineEvent(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}

		err = n.handleNodeOfflineEvent(msg, c.Title)
		if err != nil {
			return nil, err
		}

	case msgbus.PrepareRoute(n.orgName,"request.cloud.local.{{ .Org}}.node.controller.nodefeeder.publish"):
		nodeMsg := &cpb.NodeFeederMessage{}
		if err := anypb.UnmarshalTo(e.Msg, nodeMsg, proto.UnmarshalOptions{}); err != nil {
			log.Errorf("Failed to unmarshal to NodeFeederMessage: %v", err)
			return nil, err
		}
		err := n.handleNodeRestartEvent(nodeMsg)
		if err != nil {
			log.Errorf("Failed to handle node restart event: %v", err)
			return nil, err
		}
	case msgbus.PrepareRoute(n.orgName, evt.EventRoutingKey[evt.EventNodeCreate]):
		c := evt.EventToEventConfig[evt.EventNodeCreate]
		msg, err := epb.UnmarshalEventRegistryNodeCreate(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}

		err = n.handleNodeCreateEvent(msg, c.Title)
		if err != nil {
			return nil, err
		}
	default:
		log.Errorf("No handler routing key %s", e.RoutingKey)
	}

	return &epb.EventResponse{}, nil
}

func (n *NotifiyEventServer) handleNodeOnlineEvent(msg *epb.NodeOnlineEvent, name string) error {
	eventData := map[string]interface{}{
		"value": name,
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
func (n *NotifiyEventServer) handleNodeRestartEvent(msg *cpb.NodeFeederMessage) error {
	log.Infof("Handling node restart event: target=%s, path=%s", msg.Target, msg.Path)
	
	targetParts := strings.Split(msg.Target, ".")
	nodeId := targetParts[len(targetParts)-1]
	
	if !strings.Contains(msg.Path, "/v1/reboot/") {
		return fmt.Errorf("not a reboot path: %s", msg.Path)
	}
	
	eventData := map[string]interface{}{
		"value": "reboot",
	}
	
	data, err := json.Marshal(eventData)
	if err != nil {
		return fmt.Errorf("failed to marshal event data: %v", err)
	}
	
	return add(
		nodeId,
		string(db.Medium), 
		db.NotificationType("event").String(),
		"controller",
		data,
		1,
		1,
		n.notifyRepo,
		n.msgbus,
		n.baseRoutingKey,
	)
}

func (n *NotifiyEventServer) handleNodeCreateEvent(msg *epb.EventRegistryNodeCreate, name string) error {
	eventData := map[string]interface{}{
		"value": name,
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

func (n *NotifiyEventServer) handleNodeOfflineEvent(msg *epb.NodeOfflineEvent, name string) error {
	eventData := map[string]interface{}{
		"value": name,
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
