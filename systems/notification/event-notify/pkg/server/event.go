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

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/msgbus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/notification/event-notify/pkg"
	"github.com/ukama/ukama/systems/notification/event-notify/pkg/db"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type EventToNotifyEventServer struct {
	orgName string
	orgId   string
	n       *EventToNotifyServer
	epb.UnimplementedEventNotificationServiceServer
}

func NewNotificationEventServer(orgName string, orgId string, n *EventToNotifyServer) *EventToNotifyEventServer {
	return &EventToNotifyEventServer{
		orgName: orgName,
		orgId:   orgId,
		n:       n,
	}
}

func (es *EventToNotifyEventServer) EventNotification(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
	log.Infof("Received a message with Routing key %s and Message %+v.", e.RoutingKey, e.Msg)
	switch e.RoutingKey {
	case msgbus.PrepareRoute(es.orgName, pkg.EventPackageCreate):
		c := pkg.EventsSTMapping["EventPackageCreate"]
		msg, err := unmarshalMessage(e.Msg, &epb.CreatePackageEvent{})
		if err != nil {
			return nil, err
		}
		event := msg.(*epb.CreatePackageEvent)
		notification := &db.Notification{
			Id:           uuid.NewV4(),
			Title:        c.Title,
			Description:  c.Description,
			Type:         db.NotificationType(c.Type),
			Scope:        db.NotificationScope(c.Scope),
			OrgId:        es.orgId,
			UserId:       "",
			NetworkId:    "",
			SubscriberId: "",
		}
		es.n.eventPbToDBNotification(notification)
		log.Infof("Received a message with Routing key %s and Message", event)

	case msgbus.PrepareRoute(es.orgName, pkg.EventMemberCreate):
		c := pkg.EventsSTMapping[pkg.EventMemberCreate]
		msg, err := unmarshalMessage(e.Msg, &epb.AddMemberEventRequest{})
		if err != nil {
			return nil, err
		}
		event := msg.(*epb.AddMemberEventRequest)
		notification := &db.Notification{
			Id:           uuid.NewV4(),
			Title:        c.Title,
			Description:  c.Description,
			Type:         db.NotificationType(c.Type),
			Scope:        db.NotificationScope(c.Scope),
			OrgId:        event.OrgId,
			UserId:       event.UserId,
			NetworkId:    "",
			SubscriberId: "",
		}
		es.n.eventPbToDBNotification(notification)
		user := &db.Users{
			Id:           uuid.NewV4(),
			OrgId:        event.OrgId,
			UserId:       event.UserId,
			Role:         db.RoleType(event.Role),
			NetworkId:    "",
			SubscriberId: "",
		}
		es.n.storeUser(user)
		log.Infof("Received a message with Routing key %s and Message", event)

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
