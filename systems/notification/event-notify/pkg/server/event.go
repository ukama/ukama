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

	log "github.com/sirupsen/logrus"
	evt "github.com/ukama/ukama/systems/common/events"
	"github.com/ukama/ukama/systems/common/grpc"
	"github.com/ukama/ukama/systems/common/msgbus"
	notif "github.com/ukama/ukama/systems/common/notification"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	csub "github.com/ukama/ukama/systems/common/rest/client/subscriber"
	"github.com/ukama/ukama/systems/common/roles"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/notification/event-notify/pkg/db"
)

type EventToNotifyEventServer struct {
	orgName string
	orgId   string
	n       *EventToNotifyServer
	sc      csub.SubscriberClient
	epb.UnimplementedEventNotificationServiceServer
}

func NewNotificationEventServer(orgName string, orgId string, subscriberClient csub.SubscriberClient, n *EventToNotifyServer) *EventToNotifyEventServer {
	return &EventToNotifyEventServer{
		orgName: orgName,
		orgId:   orgId,
		sc:      subscriberClient,
		n:       n,
	}
}

func (es *EventToNotifyEventServer) EventNotification(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
	log.Infof("Received a message with Routing key %s and Message %+v.", e.RoutingKey, e.Msg)
	switch e.RoutingKey {

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventOrgAdd]):
		c := evt.EventToEventConfig[evt.EventOrgAdd]
		msg, err := epb.UnmarshalEventOrgCreate(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}
		// Handle Org Add event
		jmsg, err := json.Marshal(msg)
		if err != nil {
			log.Errorf("Failed to store raw message for %s to db. Error %+v", c.Name, err)
		}

		_ = es.ProcessEvent(&c, msg.Id, "", "", "", msg.Owner, jmsg)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventUserAdd]):
		c := evt.EventToEventConfig[evt.EventUserAdd]
		msg, err := epb.UnmarshalEventUserCreate(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}

		// Handle Org Add event
		jmsg, err := json.Marshal(msg)
		if err != nil {
			log.Errorf("Failed to store raw message for %s to db. Error %+v", c.Name, err)
		}

		_ = es.ProcessEvent(&c, es.orgId, "", "", "", msg.UserId, jmsg)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventUserDeactivate]):
		c := evt.EventToEventConfig[evt.EventUserDeactivate]
		msg, err := epb.UnmarshalEventUserDeactivate(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}

		// Handle Org Add event
		jmsg, err := json.Marshal(msg)
		if err != nil {
			log.Errorf("Failed to store raw message for %s to db. Error %+v", c.Name, err)
		}

		_ = es.ProcessEvent(&c, es.orgId, "", "", "", msg.UserId, jmsg)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventUserDelete]):
		c := evt.EventToEventConfig[evt.EventUserDelete]
		msg, err := epb.UnmarshalEventUserDelete(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}

		// Handle Org Add event
		jmsg, err := json.Marshal(msg)
		if err != nil {
			log.Errorf("Failed to store raw message for %s to db. Error %+v", c.Name, err)
		}

		_ = es.ProcessEvent(&c, es.orgId, "", "", "", msg.UserId, jmsg)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventMemberCreate]):
		c := evt.EventToEventConfig[evt.EventMemberCreate]
		msg, err := epb.UnmarshalAddMemberEventRequest(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}

		jmsg, err := json.Marshal(msg)
		if err != nil {
			log.Errorf("Failed to store raw message for %s to db. Error %+v", c.Name, err)
		}

		_ = es.ProcessEvent(&c, msg.OrgId, "", "", "", msg.UserId, jmsg)

		user := &db.Users{
			Id:           uuid.NewV4(),
			OrgId:        msg.OrgId,
			UserId:       msg.UserId,
			Role:         roles.RoleType(msg.Role),
			NetworkId:    "",
			SubscriberId: "",
		}

		err = es.n.storeUser(user)
		if err != nil {
			log.Errorf("Error storing user: %v", err)
		}

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventMemberDelete]):
		c := evt.EventToEventConfig[evt.EventMemberDelete]
		msg, err := epb.UnmarshalDeleteMemberEventRequest(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}

		jmsg, err := json.Marshal(msg)
		if err != nil {
			log.Errorf("Failed to store raw message for %s to db. Error %+v", c.Name, err)
		}

		_ = es.ProcessEvent(&c, msg.OrgId, "", "", "", msg.UserId, jmsg)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventNetworkAdd]):
		c := evt.EventToEventConfig[evt.EventNetworkAdd]
		msg, err := epb.UnmarshalEventNetworkCreate(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}

		jmsg, err := json.Marshal(msg)
		if err != nil {
			log.Errorf("Failed to store raw message for %s to db. Error %+v", c.Name, err)
		}

		_ = es.ProcessEvent(&c, msg.OrgId, msg.Id, "", "", "", jmsg)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventNetworkDelete]):
		c := evt.EventToEventConfig[evt.EventNetworkDelete]
		msg, err := epb.UnmarshalEventNetworkDelete(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}
		// Handle Network Delete event
		jmsg, err := json.Marshal(msg)
		if err != nil {
			log.Errorf("Failed to store raw message for %s to db. Error %+v", c.Name, err)
		}

		_ = es.ProcessEvent(&c, msg.OrgId, msg.Id, "", "", "", jmsg)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventNodeCreate]):
		c := evt.EventToEventConfig[evt.EventNodeCreate]
		msg, err := epb.UnmarshalEventRegistryNodeCreate(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}

		jmsg, err := json.Marshal(msg)
		if err != nil {
			log.Errorf("Failed to store raw message for %s to db. Error %+v", c.Name, err)
		}

		_ = es.ProcessEvent(&c, es.orgId, "", msg.NodeId, "", "", jmsg)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventNodeUpdate]):
		c := evt.EventToEventConfig[evt.EventNodeUpdate]
		msg, err := epb.UnmarshalEventRegistryNodeUpdate(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}

		jmsg, err := json.Marshal(msg)
		if err != nil {
			log.Errorf("Failed to store raw message for %s to db. Error %+v", c.Name, err)
		}

		_ = es.ProcessEvent(&c, es.orgId, "", msg.NodeId, "", "", jmsg)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventNodeStateUpdate]):
		c := evt.EventToEventConfig[evt.EventNodeStateUpdate]
		msg, err := epb.UnmarshalEventRegistryNodeStatusUpdate(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}

		jmsg, err := json.Marshal(msg)
		if err != nil {
			log.Errorf("Failed to store raw message for %s to db. Error %+v", c.Name, err)
		}

		_ = es.ProcessEvent(&c, es.orgId, "", msg.NodeId, "", "", jmsg)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventNodeDelete]):
		c := evt.EventToEventConfig[evt.EventNodeDelete]
		msg, err := epb.UnmarshalEventRegistryNodeDelete(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}

		jmsg, err := json.Marshal(msg)
		if err != nil {
			log.Errorf("Failed to store raw message for %s to db. Error %+v", c.Name, err)
		}

		_ = es.ProcessEvent(&c, es.orgId, "", msg.NodeId, "", "", jmsg)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventNodeAssign]):
		c := evt.EventToEventConfig[evt.EventNodeAssign]
		msg, err := epb.UnmarshalEventRegistryNodeAssign(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}

		jmsg, err := json.Marshal(msg)
		if err != nil {
			log.Errorf("Failed to store raw message for %s to db. Error %+v", c.Name, err)
		}

		_ = es.ProcessEvent(&c, es.orgId, "", msg.NodeId, "", "", jmsg)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventNodeRelease]):
		c := evt.EventToEventConfig[evt.EventNodeRelease]
		msg, err := epb.UnmarshalEventRegistryNodeRelease(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}

		jmsg, err := json.Marshal(msg)
		if err != nil {
			log.Errorf("Failed to store raw message for %s to db. Error %+v", c.Name, err)
		}

		_ = es.ProcessEvent(&c, es.orgId, "", msg.NodeId, "", "", jmsg)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventInviteCreate]):
		c := evt.EventToEventConfig[evt.EventInviteCreate]
		msg, err := epb.UnmarshalEventInvitationCreated(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}
		jmsg, err := json.Marshal(msg)
		if err != nil {
			log.Errorf("Failed to store raw message for %s to db. Error %+v", c.Name, err)
		}

		_ = es.ProcessEvent(&c, es.orgId, "", "", "", msg.UserId, jmsg)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventInviteDelete]):
		c := evt.EventToEventConfig[evt.EventInviteDelete]
		msg, err := epb.UnmarshalEventInvitationDeleted(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}

		// Handle Invite Delete event
		jmsg, err := json.Marshal(msg)
		if err != nil {
			log.Errorf("Failed to store raw message for %s to db. Error %+v", c.Name, err)
		}

		_ = es.ProcessEvent(&c, es.orgId, "", "", "", msg.UserId, jmsg)
	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventInviteUpdate]):
		c := evt.EventToEventConfig[evt.EventInviteUpdate]
		msg, err := epb.UnmarshalEventInvitationUpdated(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}
		// Handle Invite Update event
		jmsg, err := json.Marshal(msg)
		if err != nil {
			log.Errorf("Failed to store raw message for %s to db. Error %+v", c.Name, err)
		}
		_ = es.ProcessEvent(&c, es.orgId, "", "", "", msg.UserId, jmsg)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventMeshNodeOnline]):
		c := evt.EventToEventConfig[evt.EventMeshNodeOnline]
		msg, err := epb.UnmarshalNodeOnlineEvent(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}

		jmsg, err := json.Marshal(msg)
		if err != nil {
			log.Errorf("Failed to store raw message for %s to db. Error %+v", c.Name, err)
		}

		_ = es.ProcessEvent(&c, es.orgId, "", msg.NodeId, "", "", jmsg)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventMeshNodeOffline]):
		c := evt.EventToEventConfig[evt.EventMeshNodeOffline]
		msg, err := epb.UnmarshalNodeOfflineEvent(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}

		jmsg, err := json.Marshal(msg)
		if err != nil {
			log.Errorf("Failed to store raw message for %s to db. Error %+v", c.Name, err)
		}

		_ = es.ProcessEvent(&c, es.orgId, "", msg.NodeId, "", "", jmsg)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventSimActivate]):
		c := evt.EventToEventConfig[evt.EventSimActivate]
		msg, err := epb.UnmarshalEventSimActivation(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}

		jmsg, err := json.Marshal(msg)
		if err != nil {
			log.Errorf("Failed to store raw message for %s to db. Error %+v", c.Name, err)
		}

		_ = es.ProcessEvent(&c, es.orgId, "", "", msg.SubscriberId, "", jmsg)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventSimAllocate]):
		c := evt.EventToEventConfig[evt.EventSimAllocate]
		msg, err := epb.UnmarshalEventSimAllocation(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}

		jmsg, err := json.Marshal(msg)
		if err != nil {
			log.Errorf("Failed to store raw message for %s to db. Error %+v", c.Name, err)
		}

		_ = es.ProcessEvent(&c, es.orgId, "", "", msg.SubscriberId, "", jmsg)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventSimDelete]):
		c := evt.EventToEventConfig[evt.EventSimDelete]
		msg, err := epb.UnmarshalEventSimTermination(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}

		jmsg, err := json.Marshal(msg)
		if err != nil {
			log.Errorf("Failed to store raw message for %s to db. Error %+v", c.Name, err)
		}

		_ = es.ProcessEvent(&c, es.orgId, "", "", "", msg.SubscriberId, jmsg)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventSimAddPackage]):
		c := evt.EventToEventConfig[evt.EventSimAddPackage]
		msg, err := epb.UnmarshalEventSimAddPackage(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}
		jmsg, err := json.Marshal(msg)
		if err != nil {
			log.Errorf("Failed to store raw message for %s to db. Error %+v", c.Name, err)
		}

		_ = es.ProcessEvent(&c, es.orgId, "", "", "", msg.SubscriberId, jmsg)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventSimActivePackage]):
		c := evt.EventToEventConfig[evt.EventSimActivePackage]
		msg, err := epb.UnmarshalEventSimActivePackage(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}

		jmsg, err := json.Marshal(msg)
		if err != nil {
			log.Errorf("Failed to store raw message for %s to db. Error %+v", c.Name, err)
		}

		_ = es.ProcessEvent(&c, es.orgId, "", "", "", msg.SubscriberId, jmsg)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventSimRemovePackage]):
		c := evt.EventToEventConfig[evt.EventSimRemovePackage]
		msg, err := epb.UnmarshalEventSimRemovePackage(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}

		jmsg, err := json.Marshal(msg)
		if err != nil {
			log.Errorf("Failed to store raw message for %s to db. Error %+v", c.Name, err)
		}

		_ = es.ProcessEvent(&c, es.orgId, "", "", "", msg.SubscriberId, jmsg)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventSubscriberCreate]):
		c := evt.EventToEventConfig[evt.EventSubscriberCreate]
		msg, err := epb.UnmarshalEventSubscriberAdded(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}

		jmsg, err := json.Marshal(msg)
		if err != nil {
			log.Errorf("Failed to store raw message for %s to db. Error %+v", c.Name, err)
		}

		_ = es.ProcessEvent(&c, es.orgId, "", "", msg.SubscriberId, "", jmsg)

		user := &db.Users{
			Id:           uuid.NewV4(),
			OrgId:        es.orgId,
			UserId:       msg.SubscriberId,
			Role:         roles.TYPE_SUBSCRIBER,
			NetworkId:    msg.NetworkId,
			SubscriberId: msg.SubscriberId,
		}

		err = es.n.storeUser(user)
		if err != nil {
			log.Errorf("Error storing user: %v", err)
		}

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventSubscriberUpdate]):
		c := evt.EventToEventConfig[evt.EventSubscriberUpdate]
		msg, err := epb.UnmarshalEventSubscriberAdded(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}

		jmsg, err := json.Marshal(msg)
		if err != nil {
			log.Errorf("Failed to store raw message for %s to db. Error %+v", c.Name, err)
		}

		_ = es.ProcessEvent(&c, es.orgId, "", "", msg.SubscriberId, "", jmsg)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventSubscriberDelete]):
		c := evt.EventToEventConfig[evt.EventSubscriberDelete]
		msg, err := epb.UnmarshalEventSubscriberDeleted(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}

		jmsg, err := json.Marshal(msg)
		if err != nil {
			log.Errorf("Failed to store raw message for %s to db. Error %+v", c.Name, err)
		}

		_ = es.ProcessEvent(&c, es.orgId, "", "", msg.SubscriberId, "", jmsg)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventSimsUpload]):
		c := evt.EventToEventConfig[evt.EventSimsUpload]
		msg, err := epb.UnmarshalEventSimsUploaded(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}

		jmsg, err := json.Marshal(msg)
		if err != nil {
			log.Errorf("Failed to store raw message for %s to db. Error %+v", c.Name, err)
		}

		_ = es.ProcessEvent(&c, es.orgId, "", "", "", "", jmsg)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventBaserateUpload]):
		c := evt.EventToEventConfig[evt.EventBaserateUpload]
		msg, err := epb.UnmarshalEventBaserateUploaded(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}

		jmsg, err := json.Marshal(msg)
		if err != nil {
			log.Errorf("Failed to store raw message for %s to db. Error %+v", c.Name, err)
		}

		_ = es.ProcessEvent(&c, es.orgId, "", "", "", "", jmsg)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventPackageCreate]):
		c := evt.EventToEventConfig[evt.EventPackageCreate]
		msg, err := epb.UnmarshalCreatePackageEvent(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}

		jmsg, err := json.Marshal(msg)
		if err != nil {
			log.Errorf("Failed to store raw message for %s to db. Error %+v", c.Name, err)
		}

		_ = es.ProcessEvent(&c, es.orgId, "", "", "", "", jmsg)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventPackageUpdate]):
		c := evt.EventToEventConfig[evt.EventPackageUpdate]
		msg, err := epb.UnmarshalUpdatePackageEvent(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}

		jmsg, err := json.Marshal(msg)
		if err != nil {
			log.Errorf("Failed to store raw message for %s to db. Error %+v", c.Name, err)
		}

		_ = es.ProcessEvent(&c, es.orgId, "", "", "", "", jmsg)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventPackageDelete]):
		c := evt.EventToEventConfig[evt.EventPackageDelete]
		msg, err := epb.UnmarshalDeletePackageEvent(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}

		jmsg, err := json.Marshal(msg)
		if err != nil {
			log.Errorf("Failed to store raw message for %s to db. Error %+v", c.Name, err)
		}

		_ = es.ProcessEvent(&c, es.orgId, "", "", "", "", jmsg)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventMarkupUpdate]):
		c := evt.EventToEventConfig[evt.EventMarkupUpdate]
		msg, err := epb.UnmarshalDefaultMarkupUpdate(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}

		jmsg, err := json.Marshal(msg)
		if err != nil {
			log.Errorf("Failed to store raw message for %s to db. Error %+v", c.Name, err)
		}

		_ = es.ProcessEvent(&c, es.orgId, "", "", "", "", jmsg)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventPaymentProcessUpdate]):
		c := evt.EventToEventConfig[evt.EventPaymentProcessUpdate]
		msg, err := epb.UnmarshalPayment(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}

		jmsg, err := json.Marshal(msg)
		if err != nil {
			log.Errorf("Failed to store raw message for %s to db. Error %+v", c.Name, err)
		}

		sub, err := es.sc.GetByEmail(msg.PayerEmail)
		if err != nil {
			return nil, grpc.SqlErrorToGrpc(err, "eventnotify")
		}

		_ = es.ProcessEvent(&c, es.orgId, sub.NetworkId.String(), sub.SubscriberId.String(), sub.SubscriberId.String(), "", jmsg)

	default:
		log.Errorf("No handler routing key %s", e.RoutingKey)
	}

	return &epb.EventResponse{}, nil
}

func (es *EventToNotifyEventServer) ProcessEvent(ec *evt.EventConfig, orgId, networkId, nodeId, subscriberId, userId string, msg []byte) *db.Notification {
	log.Debugf("Processing event OrgId %s NetworkId %s nodeId %s subscriberId %s userId %s", orgId, networkId, nodeId, subscriberId, userId)

	/* Store raw event */
	event := &db.EventMsg{}
	var id uint = 0
	err := event.Data.Set(msg)
	if err != nil {
		log.Errorf("failed to assing event: %v", err)
	} else {
		id, err = es.n.storeEvent(event)
		if err != nil {
			log.Errorf("failed to store event: %v", err)
		}
	}

	dn := &db.Notification{
		Id:           uuid.NewV4(),
		Title:        ec.Title,
		Description:  ec.Description,
		Type:         notif.NotificationType(ec.Type),
		Scope:        notif.NotificationScope(ec.Scope),
		OrgId:        orgId,
		UserId:       userId,
		NetworkId:    networkId,
		NodeId:       nodeId,
		SubscriberId: subscriberId,
	}

	if id != 0 {
		dn.EventMsgID = id
	}

	err = es.n.storeNotification(dn)
	if err != nil {
		log.Errorf("failed to store notification: %v", err)
	}

	return dn
}
