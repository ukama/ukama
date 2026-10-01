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
	"github.com/ukama/ukama/systems/common/roles"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/notification/event-notify/pkg/db"

	log "github.com/sirupsen/logrus"
	evt "github.com/ukama/ukama/systems/common/events"
	notif "github.com/ukama/ukama/systems/common/notification"
	cpb "github.com/ukama/ukama/systems/common/pb/events"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	csub "github.com/ukama/ukama/systems/common/rest/client/subscriber"
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
		msg, err := cpb.UnmarshalProtoEvent[epb.EventOrgCreate](e.Msg)
		if err != nil {
			log.Errorf("Failed to unmarshal %s event: %v", c.Name, err)
			return nil, err
		}
		return handleEventOrgAdd(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventUserAdd]):
		c := evt.EventToEventConfig[evt.EventUserAdd]
		msg, err := cpb.UnmarshalProtoEvent[epb.EventUserCreate](e.Msg)
		if err != nil {
			log.Errorf("Failed to unmarshal %s event: %v", c.Name, err)
			return nil, err
		}
		return handleEventUserAdd(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventUserDeactivate]):
		c := evt.EventToEventConfig[evt.EventUserDeactivate]
		msg, err := cpb.UnmarshalProtoEvent[epb.EventUserDeactivate](e.Msg)
		if err != nil {
			log.Errorf("Failed to unmarshal %s event: %v", c.Name, err)
			return nil, err
		}
		return handleEventUserDeactivate(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventUserDelete]):
		c := evt.EventToEventConfig[evt.EventUserDelete]
		msg, err := cpb.UnmarshalProtoEvent[epb.EventUserDelete](e.Msg)
		if err != nil {
			log.Errorf("Failed to unmarshal %s event: %v", c.Name, err)
			return nil, err
		}
		return handleEventUserDelete(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventMemberCreate]):
		c := evt.EventToEventConfig[evt.EventMemberCreate]
		msg, err := cpb.UnmarshalProtoEvent[epb.AddMemberEventRequest](e.Msg)
		if err != nil {
			log.Errorf("Failed to unmarshal %s event: %v", c.Name, err)
			return nil, err
		}
		return handleEventMemberCreate(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventMemberDelete]):
		c := evt.EventToEventConfig[evt.EventMemberDelete]
		msg, err := cpb.UnmarshalProtoEvent[epb.DeleteMemberEventRequest](e.Msg)
		if err != nil {
			log.Errorf("Failed to unmarshal %s event: %v", c.Name, err)
			return nil, err
		}
		return handleEventMemberDelete(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventNetworkAdd]):
		c := evt.EventToEventConfig[evt.EventNetworkAdd]
		msg, err := cpb.UnmarshalProtoEvent[epb.EventNetworkCreate](e.Msg)
		if err != nil {
			log.Errorf("Failed to unmarshal %s event: %v", c.Name, err)
			return nil, err
		}
		return handleEventNetworkAdd(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventNetworkDelete]):
		c := evt.EventToEventConfig[evt.EventNetworkDelete]
		msg, err := cpb.UnmarshalProtoEvent[epb.EventNetworkDelete](e.Msg)
		if err != nil {
			log.Errorf("Failed to unmarshal %s event: %v", c.Name, err)
			return nil, err
		}
		return handleEventNetworkDelete(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventNodeCreate]):
		c := evt.EventToEventConfig[evt.EventNodeCreate]
		msg, err := cpb.UnmarshalProtoEvent[epb.EventRegistryNodeCreate](e.Msg)
		if err != nil {
			log.Errorf("Failed to unmarshal %s event: %v", c.Name, err)
			return nil, err
		}
		return handleEventNodeCreate(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventNodeUpdate]):
		c := evt.EventToEventConfig[evt.EventNodeUpdate]
		msg, err := cpb.UnmarshalProtoEvent[epb.EventRegistryNodeUpdate](e.Msg)
		if err != nil {
			log.Errorf("Failed to unmarshal %s event: %v", c.Name, err)
			return nil, err
		}
		return handleEventNodeUpdate(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventNodeStateUpdate]):
		c := evt.EventToEventConfig[evt.EventNodeStateUpdate]
		msg, err := cpb.UnmarshalProtoEvent[epb.EventRegistryNodeStatusUpdate](e.Msg)
		if err != nil {
			log.Errorf("Failed to unmarshal %s event: %v", c.Name, err)
			return nil, err
		}
		return handleEventNodeStateUpdate(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventNodeDelete]):
		c := evt.EventToEventConfig[evt.EventNodeDelete]
		msg, err := cpb.UnmarshalProtoEvent[epb.EventRegistryNodeDelete](e.Msg)
		if err != nil {
			log.Errorf("Failed to unmarshal %s event: %v", c.Name, err)
			return nil, err
		}
		return handleEventNodeDelete(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventNodeAssign]):
		c := evt.EventToEventConfig[evt.EventNodeAssign]
		msg, err := cpb.UnmarshalProtoEvent[epb.EventRegistryNodeAssign](e.Msg)
		if err != nil {
			log.Errorf("Failed to unmarshal %s event: %v", c.Name, err)
			return nil, err
		}
		return handleEventNodeAssign(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventNodeRelease]):
		c := evt.EventToEventConfig[evt.EventNodeRelease]
		msg, err := cpb.UnmarshalProtoEvent[epb.NodeReleasedEvent](e.Msg)
		if err != nil {
			log.Errorf("Failed to unmarshal %s event: %v", c.Name, err)
			return nil, err
		}
		return handleEventNodeRelease(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventInviteCreate]):
		c := evt.EventToEventConfig[evt.EventInviteCreate]
		msg, err := cpb.UnmarshalProtoEvent[epb.EventInvitationCreated](e.Msg)
		if err != nil {
			log.Errorf("Failed to unmarshal %s event: %v", c.Name, err)
			return nil, err
		}
		return handleEventInviteCreate(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventInviteDelete]):
		c := evt.EventToEventConfig[evt.EventInviteDelete]
		msg, err := cpb.UnmarshalProtoEvent[epb.EventInvitationDeleted](e.Msg)
		if err != nil {
			log.Errorf("Failed to unmarshal %s event: %v", c.Name, err)
			return nil, err
		}
		return handleEventInviteDelete(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventInviteUpdate]):
		c := evt.EventToEventConfig[evt.EventInviteUpdate]
		msg, err := cpb.UnmarshalProtoEvent[epb.EventInvitationUpdated](e.Msg)
		if err != nil {
			log.Errorf("Failed to unmarshal %s event: %v", c.Name, err)
			return nil, err
		}
		return handleEventInviteUpdate(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventNodeOnline]):
		c := evt.EventToEventConfig[evt.EventNodeOnline]
		msg, err := cpb.UnmarshalProtoEvent[epb.NodeOnlineEvent](e.Msg)
		if err != nil {
			log.Errorf("Failed to unmarshal %s event: %v", c.Name, err)
			return nil, err
		}
		return handleEventNodeOnline(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventNodeOffline]):
		c := evt.EventToEventConfig[evt.EventNodeOffline]
		msg, err := cpb.UnmarshalProtoEvent[epb.NodeOfflineEvent](e.Msg)
		if err != nil {
			log.Errorf("Failed to unmarshal %s event: %v", c.Name, err)
			return nil, err
		}
		return handleEventNodeOffline(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventSimServiceOn]):
		c := evt.EventToEventConfig[evt.EventSimServiceOn]
		msg, err := cpb.UnmarshalProtoEvent[epb.EventSimServiceOn](e.Msg)
		if err != nil {
			log.Errorf("Failed to unmarshal %s event: %v", c.Name, err)
			return nil, err
		}
		return handleEventSimServiceOn(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventSimAllocate]):
		c := evt.EventToEventConfig[evt.EventSimAllocate]
		msg, err := cpb.UnmarshalProtoEvent[epb.EventSimAllocation](e.Msg)
		if err != nil {
			log.Errorf("Failed to unmarshal %s event: %v", c.Name, err)
			return nil, err
		}
		return handleEventSimAllocate(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventSimDelete]):
		c := evt.EventToEventConfig[evt.EventSimDelete]
		msg, err := cpb.UnmarshalProtoEvent[epb.EventSimTermination](e.Msg)
		if err != nil {
			log.Errorf("Failed to unmarshal %s event: %v", c.Name, err)
			return nil, err
		}
		return handleEventSimDelete(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventSimAddPackage]):
		c := evt.EventToEventConfig[evt.EventSimAddPackage]
		msg, err := cpb.UnmarshalProtoEvent[epb.EventSimAddPackage](e.Msg)
		if err != nil {
			log.Errorf("Failed to unmarshal %s event: %v", c.Name, err)
			return nil, err
		}
		return handleEventSimAddPackage(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventSiteCreate]):
		c := evt.EventToEventConfig[evt.EventSiteCreate]
		msg, err := cpb.UnmarshalProtoEvent[epb.EventAddSite](e.Msg)
		if err != nil {
			log.Errorf("Failed to unmarshal %s event: %v", c.Name, err)
			return nil, err
		}
		return handleEventSiteCreate(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventSiteUpdate]):
		c := evt.EventToEventConfig[evt.EventSiteUpdate]
		msg, err := cpb.UnmarshalProtoEvent[epb.EventUpdateSite](e.Msg)
		if err != nil {
			log.Errorf("Failed to unmarshal %s event: %v", c.Name, err)
			return nil, err
		}
		return handleEventSiteUpdate(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventSimActivePackage]):
		c := evt.EventToEventConfig[evt.EventSimActivePackage]
		msg, err := cpb.UnmarshalProtoEvent[epb.EventSimActivePackage](e.Msg)
		if err != nil {
			log.Errorf("Failed to unmarshal %s event: %v", c.Name, err)
			return nil, err
		}
		return handleEventSimActivePackage(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventSimRemovePackage]):
		c := evt.EventToEventConfig[evt.EventSimRemovePackage]
		msg, err := cpb.UnmarshalProtoEvent[epb.EventSimRemovePackage](e.Msg)
		if err != nil {
			log.Errorf("Failed to unmarshal %s event: %v", c.Name, err)
			return nil, err
		}
		return handleEventSimRemovePackage(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventSubscriberCreate]):
		c := evt.EventToEventConfig[evt.EventSubscriberCreate]
		msg, err := cpb.UnmarshalProtoEvent[epb.EventSubscriberAdded](e.Msg)
		if err != nil {
			log.Errorf("Failed to unmarshal %s event: %v", c.Name, err)
			return nil, err
		}
		return handleEventSubscriberCreate(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventSubscriberUpdate]):
		c := evt.EventToEventConfig[evt.EventSubscriberUpdate]
		msg, err := cpb.UnmarshalProtoEvent[epb.EventSubscriberAdded](e.Msg)
		if err != nil {
			log.Errorf("Failed to unmarshal %s event: %v", c.Name, err)
			return nil, err
		}
		return handleEventSubscriberUpdate(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventSubscriberDelete]):
		c := evt.EventToEventConfig[evt.EventSubscriberDelete]
		msg, err := cpb.UnmarshalProtoEvent[epb.EventSubscriberDeleted](e.Msg)
		if err != nil {
			log.Errorf("Failed to unmarshal %s event: %v", c.Name, err)
			return nil, err
		}
		return handleEventSubscriberDelete(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventSimsUpload]):
		c := evt.EventToEventConfig[evt.EventSimsUpload]
		msg, err := cpb.UnmarshalProtoEvent[epb.EventSimsUploaded](e.Msg)
		if err != nil {
			log.Errorf("Failed to unmarshal %s event: %v", c.Name, err)
			return nil, err
		}
		return handleEventSimsUpload(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventBaserateUpload]):
		c := evt.EventToEventConfig[evt.EventBaserateUpload]
		msg, err := cpb.UnmarshalProtoEvent[epb.EventBaserateUploaded](e.Msg)
		if err != nil {
			log.Errorf("Failed to unmarshal %s event: %v", c.Name, err)
			return nil, err
		}
		return handleEventBaserateUpload(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventPackageCreate]):
		c := evt.EventToEventConfig[evt.EventPackageCreate]
		msg, err := cpb.UnmarshalProtoEvent[epb.CreatePackageEvent](e.Msg)
		if err != nil {
			log.Errorf("Failed to unmarshal %s event: %v", c.Name, err)
			return nil, err
		}
		return handleEventPackageCreate(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventPackageUpdate]):
		c := evt.EventToEventConfig[evt.EventPackageUpdate]
		msg, err := cpb.UnmarshalProtoEvent[epb.UpdatePackageEvent](e.Msg)
		if err != nil {
			log.Errorf("Failed to unmarshal %s event: %v", c.Name, err)
			return nil, err
		}
		return handleEventPackageUpdate(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventPackageDelete]):
		c := evt.EventToEventConfig[evt.EventPackageDelete]
		msg, err := cpb.UnmarshalProtoEvent[epb.DeletePackageEvent](e.Msg)
		if err != nil {
			log.Errorf("Failed to unmarshal %s event: %v", c.Name, err)
			return nil, err
		}
		return handleEventPackageDelete(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventMarkupUpdate]):
		c := evt.EventToEventConfig[evt.EventMarkupUpdate]
		msg, err := cpb.UnmarshalProtoEvent[epb.DefaultMarkupUpdate](e.Msg)
		if err != nil {
			log.Errorf("Failed to unmarshal %s event: %v", c.Name, err)
			return nil, err
		}
		return handleEventMarkupUpdate(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventNodeStateTransition]):
		c := evt.EventToEventConfig[evt.EventNodeStateTransition]
		msg, err := cpb.UnmarshalProtoEvent[epb.NodeStateChangeEvent](e.Msg)
		if err != nil {
			log.Errorf("Failed to unmarshal %s event: %v", c.Name, err)
			return nil, err
		}
		return handleEventNodeStateTransition(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventPaymentSuccess]):
		c := evt.EventToEventConfig[evt.EventPaymentSuccess]
		msg, err := cpb.UnmarshalProtoEvent[epb.Payment](e.Msg)
		if err != nil {
			log.Errorf("Failed to unmarshal %s event: %v", c.Name, err)
			return nil, err
		}
		return handleEventPaymentSuccess(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventPaymentFailed]):
		c := evt.EventToEventConfig[evt.EventPaymentFailed]
		msg, err := cpb.UnmarshalProtoEvent[epb.Payment](e.Msg)
		if err != nil {
			log.Errorf("Failed to unmarshal %s event: %v", c.Name, err)
			return nil, err
		}
		return handleEventPaymentFailed(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventOperationCompleted]):
		c := evt.EventToEventConfig[evt.EventOperationCompleted]
		msg, err := cpb.UnmarshalProtoEvent[epb.OperationCompletedEvent](e.Msg)
		if err != nil {
			log.Errorf("Failed to unmarshal %s event: %v", c.Name, err)
			return nil, err
		}
		return handleEventOperationCompleted(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventOperationFailed]):
		c := evt.EventToEventConfig[evt.EventOperationFailed]
		msg, err := cpb.UnmarshalProtoEvent[epb.OperationFailedEvent](e.Msg)
		if err != nil {
			log.Errorf("Failed to unmarshal %s event: %v", c.Name, err)
			return nil, err
		}
		return handleEventOperationFailed(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventInvoiceGenerate]):
		c := evt.EventToEventConfig[evt.EventInvoiceGenerate]
		msg, err := cpb.UnmarshalProtoEvent[epb.Report](e.Msg)
		if err != nil {
			log.Errorf("Failed to unmarshal %s event: %v", c.Name, err)
			return nil, err
		}
		return handleEventInvoiceGenerate(es, msg, &c)

	default:
		log.Errorf("No handler routing key %s", e.RoutingKey)
		return nil, fmt.Errorf("no handler for routing key %s", e.RoutingKey)
	}
}

func handleEventOrgAdd(es *EventToNotifyEventServer, msg *epb.EventOrgCreate, c *evt.EventConfig) (*epb.EventResponse, error) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return nil, err
	}
	return es.processEvent(c, msg.Id, "", "", "", msg.Owner, jmsg, msg.Id)
}

func handleEventUserAdd(es *EventToNotifyEventServer, msg *epb.EventUserCreate, c *evt.EventConfig) (*epb.EventResponse, error) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return nil, err
	}
	return es.processEvent(c, es.orgId, "", "", "", msg.UserId, jmsg, msg.UserId)
}

func handleEventUserDeactivate(es *EventToNotifyEventServer, msg *epb.EventUserDeactivate, c *evt.EventConfig) (*epb.EventResponse, error) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return nil, err
	}
	return es.processEvent(c, es.orgId, "", "", "", msg.UserId, jmsg, msg.UserId)
}

func handleEventUserDelete(es *EventToNotifyEventServer, msg *epb.EventUserDelete, c *evt.EventConfig) (*epb.EventResponse, error) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return nil, err
	}
	return es.processEvent(c, es.orgId, "", "", "", msg.UserId, jmsg, msg.UserId)
}

func handleEventMemberCreate(es *EventToNotifyEventServer, msg *epb.AddMemberEventRequest, c *evt.EventConfig) (*epb.EventResponse, error) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return nil, err
	}

	response, err := es.processEvent(c, msg.OrgId, "", "", "", msg.UserId, jmsg, msg.MemberId)
	if err != nil {
		return nil, err
	}

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
		return nil, err
	}

	return response, nil
}

func handleEventMemberDelete(es *EventToNotifyEventServer, msg *epb.DeleteMemberEventRequest, c *evt.EventConfig) (*epb.EventResponse, error) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return nil, err
	}
	return es.processEvent(c, msg.OrgId, "", "", "", msg.UserId, jmsg, msg.MemberId)
}

func handleEventNetworkAdd(es *EventToNotifyEventServer, msg *epb.EventNetworkCreate, c *evt.EventConfig) (*epb.EventResponse, error) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return nil, err
	}
	return es.processEvent(c, msg.OrgId, msg.Id, "", "", "", jmsg, msg.Id)
}

func handleEventNetworkDelete(es *EventToNotifyEventServer, msg *epb.EventNetworkDelete, c *evt.EventConfig) (*epb.EventResponse, error) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return nil, err
	}
	return es.processEvent(c, msg.OrgId, msg.Id, "", "", "", jmsg, msg.Id)
}

func handleEventNodeCreate(es *EventToNotifyEventServer, msg *epb.EventRegistryNodeCreate, c *evt.EventConfig) (*epb.EventResponse, error) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return nil, err
	}
	return es.processEvent(c, es.orgId, "", msg.NodeId, "", "", jmsg, msg.NodeId)
}

func handleEventNodeUpdate(es *EventToNotifyEventServer, msg *epb.EventRegistryNodeUpdate, c *evt.EventConfig) (*epb.EventResponse, error) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return nil, err
	}
	return es.processEvent(c, es.orgId, "", msg.NodeId, "", "", jmsg, msg.NodeId)
}

func handleEventNodeStateUpdate(es *EventToNotifyEventServer, msg *epb.EventRegistryNodeStatusUpdate, c *evt.EventConfig) (*epb.EventResponse, error) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return nil, err
	}
	return es.processEvent(c, es.orgId, "", msg.NodeId, "", "", jmsg, msg.NodeId)
}

func handleEventNodeDelete(es *EventToNotifyEventServer, msg *epb.EventRegistryNodeDelete, c *evt.EventConfig) (*epb.EventResponse, error) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return nil, err
	}
	return es.processEvent(c, es.orgId, "", msg.NodeId, "", "", jmsg, msg.NodeId)
}

func handleEventNodeAssign(es *EventToNotifyEventServer, msg *epb.EventRegistryNodeAssign, c *evt.EventConfig) (*epb.EventResponse, error) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return nil, err
	}
	return es.processEvent(c, es.orgId, "", msg.NodeId, "", "", jmsg, msg.NodeId)
}

func handleEventNodeRelease(es *EventToNotifyEventServer, msg *epb.NodeReleasedEvent, c *evt.EventConfig) (*epb.EventResponse, error) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return nil, err
	}
	return es.processEvent(c, es.orgId, "", msg.NodeId, "", "", jmsg, msg.NodeId)
}

func handleEventInviteCreate(es *EventToNotifyEventServer, msg *epb.EventInvitationCreated, c *evt.EventConfig) (*epb.EventResponse, error) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return nil, err
	}
	return es.processEvent(c, es.orgId, "", "", "", msg.UserId, jmsg, msg.Id)
}

func handleEventInviteDelete(es *EventToNotifyEventServer, msg *epb.EventInvitationDeleted, c *evt.EventConfig) (*epb.EventResponse, error) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return nil, err
	}
	return es.processEvent(c, es.orgId, "", "", "", msg.UserId, jmsg, msg.Id)
}

func handleEventInviteUpdate(es *EventToNotifyEventServer, msg *epb.EventInvitationUpdated, c *evt.EventConfig) (*epb.EventResponse, error) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return nil, err
	}
	return es.processEvent(c, es.orgId, "", "", "", msg.UserId, jmsg, msg.Id)
}

func handleEventNodeOnline(es *EventToNotifyEventServer, msg *epb.NodeOnlineEvent, c *evt.EventConfig) (*epb.EventResponse, error) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return nil, err
	}
	return es.processEvent(c, es.orgId, "", msg.NodeId, "", "", jmsg, msg.NodeId)
}

func handleEventNodeOffline(es *EventToNotifyEventServer, msg *epb.NodeOfflineEvent, c *evt.EventConfig) (*epb.EventResponse, error) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return nil, err
	}
	return es.processEvent(c, es.orgId, "", msg.NodeId, "", "", jmsg, msg.NodeId)
}

func handleEventSimServiceOn(es *EventToNotifyEventServer, msg *epb.EventSimServiceOn, c *evt.EventConfig) (*epb.EventResponse, error) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return nil, err
	}
	return es.processEvent(c, es.orgId, "", "", msg.SubscriberId, "", jmsg, msg.Id)
}

func handleEventSimAllocate(es *EventToNotifyEventServer, msg *epb.EventSimAllocation, c *evt.EventConfig) (*epb.EventResponse, error) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return nil, err
	}
	return es.processEvent(c, es.orgId, "", "", msg.SubscriberId, "", jmsg, msg.Id)
}

func handleEventSimDelete(es *EventToNotifyEventServer, msg *epb.EventSimTermination, c *evt.EventConfig) (*epb.EventResponse, error) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return nil, err
	}
	return es.processEvent(c, es.orgId, "", "", "", msg.SubscriberId, jmsg, msg.Id)
}

func handleEventSimAddPackage(es *EventToNotifyEventServer, msg *epb.EventSimAddPackage, c *evt.EventConfig) (*epb.EventResponse, error) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return nil, err
	}
	return es.processEvent(c, es.orgId, "", "", "", msg.SubscriberId, jmsg, msg.Id)
}

func handleEventSiteCreate(es *EventToNotifyEventServer, msg *epb.EventAddSite, c *evt.EventConfig) (*epb.EventResponse, error) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return nil, err
	}
	return es.processEvent(c, es.orgId, msg.NetworkId, "", "", "", jmsg, msg.SiteId)
}

func handleEventSiteUpdate(es *EventToNotifyEventServer, msg *epb.EventUpdateSite, c *evt.EventConfig) (*epb.EventResponse, error) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return nil, err
	}
	return es.processEvent(c, es.orgId, msg.NetworkId, "", "", "", jmsg, msg.SiteId)
}

func handleEventSimActivePackage(es *EventToNotifyEventServer, msg *epb.EventSimActivePackage, c *evt.EventConfig) (*epb.EventResponse, error) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return nil, err
	}
	return es.processEvent(c, es.orgId, "", "", "", msg.SubscriberId, jmsg, msg.Id)
}

func handleEventSimRemovePackage(es *EventToNotifyEventServer, msg *epb.EventSimRemovePackage, c *evt.EventConfig) (*epb.EventResponse, error) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return nil, err
	}
	return es.processEvent(c, es.orgId, "", "", "", msg.SubscriberId, jmsg, msg.Id)
}

func handleEventSubscriberCreate(es *EventToNotifyEventServer, msg *epb.EventSubscriberAdded, c *evt.EventConfig) (*epb.EventResponse, error) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return nil, err
	}

	response, err := es.processEvent(c, es.orgId, "", "", msg.SubscriberId, "", jmsg, msg.SubscriberId)
	if err != nil {
		return nil, err
	}

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
		return nil, err
	}

	return response, nil
}

func handleEventSubscriberUpdate(es *EventToNotifyEventServer, msg *epb.EventSubscriberAdded, c *evt.EventConfig) (*epb.EventResponse, error) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return nil, err
	}
	return es.processEvent(c, es.orgId, "", "", msg.SubscriberId, "", jmsg, msg.SubscriberId)
}

func handleEventSubscriberDelete(es *EventToNotifyEventServer, msg *epb.EventSubscriberDeleted, c *evt.EventConfig) (*epb.EventResponse, error) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return nil, err
	}
	return es.processEvent(c, es.orgId, "", "", msg.SubscriberId, "", jmsg, msg.SubscriberId)
}

func handleEventSimsUpload(es *EventToNotifyEventServer, msg *epb.EventSimsUploaded, c *evt.EventConfig) (*epb.EventResponse, error) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return nil, err
	}
	return es.processEvent(c, es.orgId, "", "", "", "", jmsg, "")
}

func handleEventBaserateUpload(es *EventToNotifyEventServer, msg *epb.EventBaserateUploaded, c *evt.EventConfig) (*epb.EventResponse, error) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return nil, err
	}
	return es.processEvent(c, es.orgId, "", "", "", "", jmsg, "")
}

func handleEventPackageCreate(es *EventToNotifyEventServer, msg *epb.CreatePackageEvent, c *evt.EventConfig) (*epb.EventResponse, error) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return nil, err
	}
	return es.processEvent(c, es.orgId, "", "", "", "", jmsg, msg.Uuid)
}

func handleEventPackageUpdate(es *EventToNotifyEventServer, msg *epb.UpdatePackageEvent, c *evt.EventConfig) (*epb.EventResponse, error) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return nil, err
	}
	return es.processEvent(c, es.orgId, "", "", "", "", jmsg, msg.Uuid)
}

func handleEventPackageDelete(es *EventToNotifyEventServer, msg *epb.DeletePackageEvent, c *evt.EventConfig) (*epb.EventResponse, error) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return nil, err
	}
	return es.processEvent(c, es.orgId, "", "", "", "", jmsg, msg.Uuid)
}

func handleEventMarkupUpdate(es *EventToNotifyEventServer, msg *epb.DefaultMarkupUpdate, c *evt.EventConfig) (*epb.EventResponse, error) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return nil, err
	}
	return es.processEvent(c, es.orgId, "", "", "", "", jmsg, "")
}

func handleEventNodeStateTransition(es *EventToNotifyEventServer, msg *epb.NodeStateChangeEvent, c *evt.EventConfig) (*epb.EventResponse, error) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return nil, err
	}

	dynamicConfig := *c
	shortNodeId := msg.NodeId
	if len(msg.NodeId) > 6 {
		shortNodeId = msg.NodeId[len(msg.NodeId)-6:]
	}

	dynamicConfig.Title = fmt.Sprintf("Node %s: %s", shortNodeId, msg.State)
	dynamicConfig.Description = fmt.Sprintf("Status: %s", msg.Substate)

	notificationType := notif.TYPE_INFO

	switch msg.State {
	case "Faulty":
		notificationType = notif.TYPE_CRITICAL
	case "Unknown":
		notificationType = notif.TYPE_ACTIONABLE_WARNING
	}

	if notificationType == notif.TYPE_INFO {
		switch msg.Substate {
		case "off":
			notificationType = notif.TYPE_WARNING
		case "reboot", "update", "upgrade", "downgrade":
			notificationType = notif.TYPE_WARNING
		}
	}

	dynamicConfig.Type = notificationType

	return es.processEvent(&dynamicConfig, es.orgId, "", msg.NodeId, "", "", jmsg, msg.NodeId)
}

func handleEventPaymentSuccess(es *EventToNotifyEventServer, msg *epb.Payment, c *evt.EventConfig) (*epb.EventResponse, error) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return nil, err
	}

	if msg.ItemType != ukama.ItemTypeInvoice.String() {
		log.Infof("Ignoring successful payment for item type %s", msg.ItemType)
		return &epb.EventResponse{}, nil
	}

	metadata := map[string]string{}

	err = json.Unmarshal(msg.Metadata, &metadata)
	if err != nil {
		log.Errorf("failed to Unmarshal payment metadata as map[string]string: %v", err)
	}

	targetId, ok := metadata["targetId"]
	if !ok {
		log.Errorf("missing targetId metadata for successful package payment: %s", msg.ItemId)
		return nil, fmt.Errorf("missing targetId metadata for successful package payment: %s", msg.ItemId)
	}

	return es.processEvent(c, es.orgId, "", "", targetId, targetId, jmsg, msg.Id)
}

func handleEventOperationCompleted(es *EventToNotifyEventServer, msg *epb.OperationCompletedEvent, c *evt.EventConfig) (*epb.EventResponse, error) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return nil, err
	}
	return es.processEvent(c, es.orgId, "", "", "", "", jmsg, msg.OperationId)
}

func handleEventOperationFailed(es *EventToNotifyEventServer, msg *epb.OperationFailedEvent, c *evt.EventConfig) (*epb.EventResponse, error) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return nil, err
	}
	return es.processEvent(c, es.orgId, "", "", "", "", jmsg, msg.OperationId)
}

func handleEventPaymentFailed(es *EventToNotifyEventServer, msg *epb.Payment, c *evt.EventConfig) (*epb.EventResponse, error) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return nil, err
	}

	if msg.ItemType != ukama.ItemTypePackage.String() {
		log.Errorf("unexpected item type for successful payment: %s", msg.ItemType)
		return nil, fmt.Errorf("unexpected item type for successful payment: %s", msg.ItemType)
	}

	metadata := map[string]string{}

	err = json.Unmarshal(msg.Metadata, &metadata)
	if err != nil {
		log.Errorf("failed to Unmarshal payment metadata as map[string]string: %v", err)
	}

	targetId, ok := metadata["targetId"]
	if !ok {
		log.Errorf("missing targetId metadata for successful package payment: %s", msg.ItemId)
		return nil, fmt.Errorf("missing targetId metadata for successful package payment: %s", msg.ItemId)
	}

	return es.processEvent(c, es.orgId, "", "", targetId, targetId, jmsg, msg.Id)
}

func handleEventInvoiceGenerate(es *EventToNotifyEventServer, msg *epb.Report, c *evt.EventConfig) (*epb.EventResponse, error) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return nil, err
	}
	return es.processEvent(c, es.orgId, msg.NetworkId, "", "", "", jmsg, msg.Id)
}

func (es *EventToNotifyEventServer) processEvent(ec *evt.EventConfig, orgId, networkId, nodeId, subscriberId, userId string, msg []byte, rid string) (*epb.EventResponse, error) {
	log.Debugf("Processing event OrgId %s NetworkId %s nodeId %s subscriberId %s userId %s", orgId, networkId, nodeId, subscriberId, userId)

	/* Store raw event */
	event := &db.EventMsg{}
	event.Key = ec.Name
	err := event.Data.Set(msg)
	if err != nil {
		log.Errorf("failed to assing event: %v", err)
		return nil, err
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
		ResourceId:   rid,
		SubscriberId: subscriberId,
	}

	err = es.n.storeNotification(event, dn)
	if err != nil {
		log.Errorf("failed to store notification: %v", err)
		return nil, err
	}

	return &epb.EventResponse{}, nil
}
