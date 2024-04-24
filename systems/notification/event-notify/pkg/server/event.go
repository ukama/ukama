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
	"strings"

	log "github.com/sirupsen/logrus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	"github.com/ukama/ukama/systems/notification/event-notify/pkg"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type EventToNotifyEventServer struct {
	orgName string
	n       *EventToNotifyServer
	epb.UnimplementedEventNotificationServiceServer
}

func NewNotificationEventServer(orgName string, n *EventToNotifyServer) *EventToNotifyEventServer {
	return &EventToNotifyEventServer{
		orgName: "*",
		n:       n,
	}
}

func (es *EventToNotifyEventServer) EventNotification(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
	log.Infof("Received a message with Routing key %s and Message %+v.", e.RoutingKey, e.Msg)

	org := strings.Split(e.RoutingKey, ".")[3]
	newRoutingKey := strings.ReplaceAll(e.RoutingKey, "event.cloud.local."+org+".", "")
	switch newRoutingKey {
	case pkg.EventPackageCreate:
		msg, err := unmarshalMessage(e.Msg, &epb.CreatePackageEvent{})
		if err != nil {
			return nil, err
		}
		log.Infof("Received a message with Routing key %s and Message", msg)
	default:
		log.Errorf("No handler routing key %s", e.RoutingKey)
	}

	return &epb.EventResponse{}, nil
}

func unmarshalMessage(msg *anypb.Any, p proto.Message) (proto.Message, error) {
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal message with : %+v. Error %s.", msg, err.Error())
		return nil, err
	}

	return p, nil
}
