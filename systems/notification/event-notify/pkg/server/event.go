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
		_, err := unmarshalMessage(e.Msg, &epb.CreatePackageEvent{})
		if err != nil {
			return nil, err
		}
		// event := msg.(*epb.CreatePackageEvent)
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

	case msgbus.PrepareRoute(es.orgName, pkg.EventPackageUpdate):
		c := pkg.EventsSTMapping["EventPackageUpdate"]
		_, err := unmarshalMessage(e.Msg, &epb.UpdatePackageEvent{})
		if err != nil {
			return nil, err
		}
		// event := msg.(*epb.UpdatePackageEvent)
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

	case msgbus.PrepareRoute(es.orgName, pkg.EventNetworkAdd):
		c := pkg.EventsSTMapping["EventNetworkAdd"]
		_, err := unmarshalMessage(e.Msg, &epb.NetworkCreatedEvent{})
		if err != nil {
			return nil, err
		}
		// event := msg.(*epb.NetworkCreatedEvent)
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

	case msgbus.PrepareRoute(es.orgName, pkg.EventNodeCreate):
		c := pkg.EventsSTMapping["EventNodeCreate"]
		_, err := unmarshalMessage(e.Msg, &epb.NodeCreatedEvent{})
		if err != nil {
			return nil, err
		}
		// event := msg.(*epb.NodeCreatedEvent)
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

	case msgbus.PrepareRoute(es.orgName, pkg.EventNodeUpdate):
		c := pkg.EventsSTMapping["EventNodeUpdate"]
		_, err := unmarshalMessage(e.Msg, &epb.NodeUpdatedEvent{})
		if err != nil {
			return nil, err
		}
		// event := msg.(*epb.NodeUpdatedEvent)
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

	case msgbus.PrepareRoute(es.orgName, pkg.EventNodeStateUpdate):
		c := pkg.EventsSTMapping["EventNodeStateUpdate"]
		_, err := unmarshalMessage(e.Msg, &epb.NodeStateUpdatedEvent{})
		if err != nil {
			return nil, err
		}
		// event := msg.(*epb.NodeStateUpdatedEvent)
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

	case msgbus.PrepareRoute(es.orgName, pkg.EventNodeDelete):
		c := pkg.EventsSTMapping["EventNodeDelete"]
		_, err := unmarshalMessage(e.Msg, &epb.NodeDeletedEvent{})
		if err != nil {
			return nil, err
		}
		// event := msg.(*epb.NodeDeletedEvent)
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

	case msgbus.PrepareRoute(es.orgName, pkg.EventNodeAssign):
		c := pkg.EventsSTMapping["EventNodeAssign"]
		_, err := unmarshalMessage(e.Msg, &epb.NodeAssignedEvent{})
		if err != nil {
			return nil, err
		}
		// event := msg.(*epb.NodeAssignedEvent)
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

	case msgbus.PrepareRoute(es.orgName, pkg.EventNodeRelease):
		c := pkg.EventsSTMapping["EventNodeRelease"]
		_, err := unmarshalMessage(e.Msg, &epb.NodeReleasedEvent{})
		if err != nil {
			return nil, err
		}
		// event := msg.(*epb.NodeReleasedEvent)
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

	case msgbus.PrepareRoute(es.orgName, pkg.EventMeshNodeOnline):
		c := pkg.EventsSTMapping["EventMeshNodeOnline"]
		_, err := unmarshalMessage(e.Msg, &epb.NodeOnlineEvent{})
		if err != nil {
			return nil, err
		}
		// event := msg.(*epb.NodeOnlineEvent)
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

	case msgbus.PrepareRoute(es.orgName, pkg.EventMeshNodeOffline):
		c := pkg.EventsSTMapping["EventMeshNodeOffline"]
		_, err := unmarshalMessage(e.Msg, &epb.NodeOfflineEvent{})
		if err != nil {
			return nil, err
		}
		// event := msg.(*epb.NodeOfflineEvent)
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
