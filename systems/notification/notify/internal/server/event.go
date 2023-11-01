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

	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/notification/notify/internal"
	"github.com/ukama/ukama/systems/notification/notify/internal/db"
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

	msg, err := n.unmarshalNotificationSentEvent(e.Msg)
	if err != nil {
		return nil, err
	}

	err = n.handleNotificationSentEvent(e.RoutingKey, msg)
	if err != nil {
		return nil, err
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
