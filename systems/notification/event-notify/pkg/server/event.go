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

	log "github.com/sirupsen/logrus"
	evt "github.com/ukama/ukama/systems/common/events"
	"github.com/ukama/ukama/systems/common/msgbus"
	notif "github.com/ukama/ukama/systems/common/notification"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	csub "github.com/ukama/ukama/systems/common/rest/client/subscriber"
	"github.com/ukama/ukama/systems/common/roles"
	"github.com/ukama/ukama/systems/common/ukama"
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
		handleEventOrgAdd(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventUserAdd]):
		c := evt.EventToEventConfig[evt.EventUserAdd]
		msg, err := epb.UnmarshalEventUserCreate(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}
		handleEventUserAdd(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventUserDeactivate]):
		c := evt.EventToEventConfig[evt.EventUserDeactivate]
		msg, err := epb.UnmarshalEventUserDeactivate(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}
		handleEventUserDeactivate(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventUserDelete]):
		c := evt.EventToEventConfig[evt.EventUserDelete]
		msg, err := epb.UnmarshalEventUserDelete(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}
		handleEventUserDelete(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventMemberCreate]):
		c := evt.EventToEventConfig[evt.EventMemberCreate]
		msg, err := epb.UnmarshalAddMemberEventRequest(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}
		handleEventMemberCreate(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventMemberDelete]):
		c := evt.EventToEventConfig[evt.EventMemberDelete]
		msg, err := epb.UnmarshalDeleteMemberEventRequest(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}
		handleEventMemberDelete(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventNetworkAdd]):
		c := evt.EventToEventConfig[evt.EventNetworkAdd]
		msg, err := epb.UnmarshalEventNetworkCreate(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}
		handleEventNetworkAdd(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventNetworkDelete]):
		c := evt.EventToEventConfig[evt.EventNetworkDelete]
		msg, err := epb.UnmarshalEventNetworkDelete(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}
		handleEventNetworkDelete(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventNodeCreate]):
		c := evt.EventToEventConfig[evt.EventNodeCreate]
		msg, err := epb.UnmarshalEventRegistryNodeCreate(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}
		handleEventNodeCreate(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventNodeUpdate]):
		c := evt.EventToEventConfig[evt.EventNodeUpdate]
		msg, err := epb.UnmarshalEventRegistryNodeUpdate(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}
		handleEventNodeUpdate(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventNodeStateUpdate]):
		c := evt.EventToEventConfig[evt.EventNodeStateUpdate]
		msg, err := epb.UnmarshalEventRegistryNodeStatusUpdate(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}
		handleEventNodeStateUpdate(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventNodeDelete]):
		c := evt.EventToEventConfig[evt.EventNodeDelete]
		msg, err := epb.UnmarshalEventRegistryNodeDelete(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}
		handleEventNodeDelete(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventNodeAssign]):
		c := evt.EventToEventConfig[evt.EventNodeAssign]
		msg, err := epb.UnmarshalEventRegistryNodeAssign(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}
		handleEventNodeAssign(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventNodeRelease]):
		c := evt.EventToEventConfig[evt.EventNodeRelease]
		msg, err := epb.UnmarshalEventRegistryNodeRelease(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}
		handleEventNodeRelease(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventInviteCreate]):
		c := evt.EventToEventConfig[evt.EventInviteCreate]
		msg, err := epb.UnmarshalEventInvitationCreated(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}
		handleEventInviteCreate(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventInviteDelete]):
		c := evt.EventToEventConfig[evt.EventInviteDelete]
		msg, err := epb.UnmarshalEventInvitationDeleted(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}
		handleEventInviteDelete(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventInviteUpdate]):
		c := evt.EventToEventConfig[evt.EventInviteUpdate]
		msg, err := epb.UnmarshalEventInvitationUpdated(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}
		handleEventInviteUpdate(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventNodeOnline]):
		c := evt.EventToEventConfig[evt.EventNodeOnline]
		msg, err := epb.UnmarshalNodeOnlineEvent(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}
		handleEventNodeOnline(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventNodeOffline]):
		c := evt.EventToEventConfig[evt.EventNodeOffline]
		msg, err := epb.UnmarshalNodeOfflineEvent(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}
		handleEventNodeOffline(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventSimActivate]):
		c := evt.EventToEventConfig[evt.EventSimActivate]
		msg, err := epb.UnmarshalEventSimActivation(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}
		handleEventSimActivate(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventSimAllocate]):
		c := evt.EventToEventConfig[evt.EventSimAllocate]
		msg, err := epb.UnmarshalEventSimAllocation(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}
		handleEventSimAllocate(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventSimDelete]):
		c := evt.EventToEventConfig[evt.EventSimDelete]
		msg, err := epb.UnmarshalEventSimTermination(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}
		handleEventSimDelete(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventSimAddPackage]):
		c := evt.EventToEventConfig[evt.EventSimAddPackage]
		msg, err := epb.UnmarshalEventSimAddPackage(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}
		handleEventSimAddPackage(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventSiteCreate]):
		c := evt.EventToEventConfig[evt.EventSiteCreate]
		msg, err := epb.UnmarshalEventAddSite(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}
		handleEventSiteCreate(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventSiteUpdate]):
		c := evt.EventToEventConfig[evt.EventSiteUpdate]
		msg, err := epb.UnmarshalEventUpdateSite(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}
		handleEventSiteUpdate(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventSimActivePackage]):
		c := evt.EventToEventConfig[evt.EventSimActivePackage]
		msg, err := epb.UnmarshalEventSimActivePackage(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}
		handleEventSimActivePackage(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventSimRemovePackage]):
		c := evt.EventToEventConfig[evt.EventSimRemovePackage]
		msg, err := epb.UnmarshalEventSimRemovePackage(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}
		handleEventSimRemovePackage(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventSubscriberCreate]):
		c := evt.EventToEventConfig[evt.EventSubscriberCreate]
		msg, err := epb.UnmarshalEventSubscriberAdded(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}
		handleEventSubscriberCreate(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventSubscriberUpdate]):
		c := evt.EventToEventConfig[evt.EventSubscriberUpdate]
		msg, err := epb.UnmarshalEventSubscriberAdded(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}
		handleEventSubscriberUpdate(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventSubscriberDelete]):
		c := evt.EventToEventConfig[evt.EventSubscriberDelete]
		msg, err := epb.UnmarshalEventSubscriberDeleted(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}
		handleEventSubscriberDelete(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventSimsUpload]):
		c := evt.EventToEventConfig[evt.EventSimsUpload]
		msg, err := epb.UnmarshalEventSimsUploaded(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}
		handleEventSimsUpload(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventBaserateUpload]):
		c := evt.EventToEventConfig[evt.EventBaserateUpload]
		msg, err := epb.UnmarshalEventBaserateUploaded(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}
		handleEventBaserateUpload(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventPackageCreate]):
		c := evt.EventToEventConfig[evt.EventPackageCreate]
		msg, err := epb.UnmarshalCreatePackageEvent(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}
		handleEventPackageCreate(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventPackageUpdate]):
		c := evt.EventToEventConfig[evt.EventPackageUpdate]
		msg, err := epb.UnmarshalUpdatePackageEvent(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}
		handleEventPackageUpdate(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventPackageDelete]):
		c := evt.EventToEventConfig[evt.EventPackageDelete]
		msg, err := epb.UnmarshalDeletePackageEvent(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}
		handleEventPackageDelete(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventMarkupUpdate]):
		c := evt.EventToEventConfig[evt.EventMarkupUpdate]
		msg, err := epb.UnmarshalDefaultMarkupUpdate(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}
		handleEventMarkupUpdate(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventNodeStateTransition]):
		c := evt.EventToEventConfig[evt.EventNodeStateTransition]
		msg, err := epb.UnmarshalEventNodeStateTransition(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}
		handleEventNodeStateTransition(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventPaymentSuccess]):
		c := evt.EventToEventConfig[evt.EventPaymentSuccess]
		msg, err := epb.UnmarshalPayment(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}
		handleEventPaymentSuccess(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventPaymentFailed]):
		c := evt.EventToEventConfig[evt.EventPaymentFailed]
		msg, err := epb.UnmarshalPayment(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}
		handleEventPaymentFailed(es, msg, &c)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventInvoiceGenerate]):
		c := evt.EventToEventConfig[evt.EventInvoiceGenerate]
		msg, err := epb.UnmarshalInvoiceGenerated(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}
		handleEventInvoiceGenerate(es, msg, &c)

	default:
		log.Errorf("No handler routing key %s", e.RoutingKey)
	}

	return &epb.EventResponse{}, nil
}

func handleEventOrgAdd(es *EventToNotifyEventServer, msg *epb.EventOrgCreate, c *evt.EventConfig) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return
	}
	es.processEvent(c, msg.Id, "", "", "", msg.Owner, jmsg, msg.Id)
}

func handleEventUserAdd(es *EventToNotifyEventServer, msg *epb.EventUserCreate, c *evt.EventConfig) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return
	}
	es.processEvent(c, es.orgId, "", "", "", msg.UserId, jmsg, msg.UserId)
}

func handleEventUserDeactivate(es *EventToNotifyEventServer, msg *epb.EventUserDeactivate, c *evt.EventConfig) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return
	}
	es.processEvent(c, es.orgId, "", "", "", msg.UserId, jmsg, msg.UserId)
}

func handleEventUserDelete(es *EventToNotifyEventServer, msg *epb.EventUserDelete, c *evt.EventConfig) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return
	}
	es.processEvent(c, es.orgId, "", "", "", msg.UserId, jmsg, msg.UserId)
}

func handleEventMemberCreate(es *EventToNotifyEventServer, msg *epb.AddMemberEventRequest, c *evt.EventConfig) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return
	}
	es.processEvent(c, msg.OrgId, "", "", "", msg.UserId, jmsg, msg.MemberId)

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
}

func handleEventMemberDelete(es *EventToNotifyEventServer, msg *epb.DeleteMemberEventRequest, c *evt.EventConfig) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return
	}
	es.processEvent(c, msg.OrgId, "", "", "", msg.UserId, jmsg, msg.MemberId)
}

func handleEventNetworkAdd(es *EventToNotifyEventServer, msg *epb.EventNetworkCreate, c *evt.EventConfig) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return
	}
	es.processEvent(c, msg.OrgId, msg.Id, "", "", "", jmsg, msg.Id)
}

func handleEventNetworkDelete(es *EventToNotifyEventServer, msg *epb.EventNetworkDelete, c *evt.EventConfig) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return
	}
	es.processEvent(c, msg.OrgId, msg.Id, "", "", "", jmsg, msg.Id)
}

func handleEventNodeCreate(es *EventToNotifyEventServer, msg *epb.EventRegistryNodeCreate, c *evt.EventConfig) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return
	}
	es.processEvent(c, es.orgId, "", msg.NodeId, "", "", jmsg, msg.NodeId)
}

func handleEventNodeUpdate(es *EventToNotifyEventServer, msg *epb.EventRegistryNodeUpdate, c *evt.EventConfig) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return
	}
	es.processEvent(c, es.orgId, "", msg.NodeId, "", "", jmsg, msg.NodeId)
}

func handleEventNodeStateUpdate(es *EventToNotifyEventServer, msg *epb.EventRegistryNodeStatusUpdate, c *evt.EventConfig) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return
	}
	es.processEvent(c, es.orgId, "", msg.NodeId, "", "", jmsg, msg.NodeId)
}

func handleEventNodeDelete(es *EventToNotifyEventServer, msg *epb.EventRegistryNodeDelete, c *evt.EventConfig) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return
	}
	es.processEvent(c, es.orgId, "", msg.NodeId, "", "", jmsg, msg.NodeId)
}

func handleEventNodeAssign(es *EventToNotifyEventServer, msg *epb.EventRegistryNodeAssign, c *evt.EventConfig) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return
	}
	es.processEvent(c, es.orgId, "", msg.NodeId, "", "", jmsg, msg.NodeId)
}

func handleEventNodeRelease(es *EventToNotifyEventServer, msg *epb.EventRegistryNodeRelease, c *evt.EventConfig) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return
	}
	es.processEvent(c, es.orgId, "", msg.NodeId, "", "", jmsg, msg.NodeId)
}

func handleEventInviteCreate(es *EventToNotifyEventServer, msg *epb.EventInvitationCreated, c *evt.EventConfig) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return
	}
	es.processEvent(c, es.orgId, "", "", "", msg.UserId, jmsg, msg.Id)
}

func handleEventInviteDelete(es *EventToNotifyEventServer, msg *epb.EventInvitationDeleted, c *evt.EventConfig) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return
	}
	es.processEvent(c, es.orgId, "", "", "", msg.UserId, jmsg, msg.Id)
}

func handleEventInviteUpdate(es *EventToNotifyEventServer, msg *epb.EventInvitationUpdated, c *evt.EventConfig) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return
	}
	es.processEvent(c, es.orgId, "", "", "", msg.UserId, jmsg, msg.Id)
}

func handleEventNodeOnline(es *EventToNotifyEventServer, msg *epb.NodeOnlineEvent, c *evt.EventConfig) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return
	}
	es.processEvent(c, es.orgId, "", msg.NodeId, "", "", jmsg, msg.NodeId)
}

func handleEventNodeOffline(es *EventToNotifyEventServer, msg *epb.NodeOfflineEvent, c *evt.EventConfig) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return
	}
	es.processEvent(c, es.orgId, "", msg.NodeId, "", "", jmsg, msg.NodeId)
}

func handleEventSimActivate(es *EventToNotifyEventServer, msg *epb.EventSimActivation, c *evt.EventConfig) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return
	}
	es.processEvent(c, es.orgId, "", "", msg.SubscriberId, "", jmsg, msg.Id)
}

func handleEventSimAllocate(es *EventToNotifyEventServer, msg *epb.EventSimAllocation, c *evt.EventConfig) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return
	}
	es.processEvent(c, es.orgId, "", "", msg.SubscriberId, "", jmsg, msg.Id)
}

func handleEventSimDelete(es *EventToNotifyEventServer, msg *epb.EventSimTermination, c *evt.EventConfig) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return
	}
	es.processEvent(c, es.orgId, "", "", "", msg.SubscriberId, jmsg, msg.Id)
}

func handleEventSimAddPackage(es *EventToNotifyEventServer, msg *epb.EventSimAddPackage, c *evt.EventConfig) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return
	}
	es.processEvent(c, es.orgId, "", "", "", msg.SubscriberId, jmsg, msg.Id)
}

func handleEventSiteCreate(es *EventToNotifyEventServer, msg *epb.EventAddSite, c *evt.EventConfig) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return
	}
	es.processEvent(c, es.orgId, msg.NetworkId, "", "", "", jmsg, msg.SiteId)
}

func handleEventSiteUpdate(es *EventToNotifyEventServer, msg *epb.EventUpdateSite, c *evt.EventConfig) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return
	}
	es.processEvent(c, es.orgId, msg.NetworkId, "", "", "", jmsg, msg.SiteId)
}

func handleEventSimActivePackage(es *EventToNotifyEventServer, msg *epb.EventSimActivePackage, c *evt.EventConfig) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return
	}
	es.processEvent(c, es.orgId, "", "", "", msg.SubscriberId, jmsg, msg.Id)
}

func handleEventSimRemovePackage(es *EventToNotifyEventServer, msg *epb.EventSimRemovePackage, c *evt.EventConfig) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return
	}
	es.processEvent(c, es.orgId, "", "", "", msg.SubscriberId, jmsg, msg.Id)
}

func handleEventSubscriberCreate(es *EventToNotifyEventServer, msg *epb.EventSubscriberAdded, c *evt.EventConfig) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return
	}
	es.processEvent(c, es.orgId, "", "", msg.SubscriberId, "", jmsg, msg.SubscriberId)

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
}

func handleEventSubscriberUpdate(es *EventToNotifyEventServer, msg *epb.EventSubscriberAdded, c *evt.EventConfig) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return
	}
	es.processEvent(c, es.orgId, "", "", msg.SubscriberId, "", jmsg, msg.SubscriberId)
}

func handleEventSubscriberDelete(es *EventToNotifyEventServer, msg *epb.EventSubscriberDeleted, c *evt.EventConfig) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return
	}
	es.processEvent(c, es.orgId, "", "", msg.SubscriberId, "", jmsg, msg.SubscriberId)
}

func handleEventSimsUpload(es *EventToNotifyEventServer, msg *epb.EventSimsUploaded, c *evt.EventConfig) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return
	}
	es.processEvent(c, es.orgId, "", "", "", "", jmsg, "")
}

func handleEventBaserateUpload(es *EventToNotifyEventServer, msg *epb.EventBaserateUploaded, c *evt.EventConfig) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return
	}
	es.processEvent(c, es.orgId, "", "", "", "", jmsg, "")
}

func handleEventPackageCreate(es *EventToNotifyEventServer, msg *epb.CreatePackageEvent, c *evt.EventConfig) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return
	}
	es.processEvent(c, es.orgId, "", "", "", "", jmsg, msg.Uuid)
}

func handleEventPackageUpdate(es *EventToNotifyEventServer, msg *epb.UpdatePackageEvent, c *evt.EventConfig) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return
	}
	es.processEvent(c, es.orgId, "", "", "", "", jmsg, msg.Uuid)
}

func handleEventPackageDelete(es *EventToNotifyEventServer, msg *epb.DeletePackageEvent, c *evt.EventConfig) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return
	}
	es.processEvent(c, es.orgId, "", "", "", "", jmsg, msg.Uuid)
}

func handleEventMarkupUpdate(es *EventToNotifyEventServer, msg *epb.DefaultMarkupUpdate, c *evt.EventConfig) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return
	}
	es.processEvent(c, es.orgId, "", "", "", "", jmsg, "")
}

func handleEventNodeStateTransition(es *EventToNotifyEventServer, msg *epb.NodeStateChangeEvent, c *evt.EventConfig) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return
	}

	dynamicConfig := *c
	shortNodeId := msg.NodeId
	if len(msg.NodeId) > 6 {
		shortNodeId = msg.NodeId[len(msg.NodeId)-6:]
	}

	dynamicConfig.Title = fmt.Sprintf("Node %s: %s", shortNodeId, msg.State)
	dynamicConfig.Description = fmt.Sprintf("Status: %s", msg.Substate)

	notificationType := notif.TYPE_INFO

	if msg.State == "Faulty" {
		notificationType = notif.TYPE_CRITICAL
	} else if msg.State == "Unknown" {
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

	es.processEvent(&dynamicConfig, es.orgId, "", msg.NodeId, "", "", jmsg, msg.NodeId)
}

func handleEventPaymentSuccess(es *EventToNotifyEventServer, msg *epb.Payment, c *evt.EventConfig) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return
	}

	if msg.ItemType != ukama.ItemTypeInvoice.String() {
		log.Errorf("unexpected item type for failed payment: %s", msg.ItemType)
		return
	}

	metadata := map[string]string{}

	err = json.Unmarshal(msg.Metadata, &metadata)
	if err != nil {
		log.Errorf("failed to Unmarshal payment metadata as map[string]string: %v", err)
	}

	targetId, ok := metadata["targetId"]
	if !ok {
		log.Errorf("missing targetId metadata for failed package payment: %s", msg.ItemId)
		return
	}

	es.processEvent(c, es.orgId, "", "", targetId, targetId, jmsg, msg.Id)
}

func handleEventPaymentFailed(es *EventToNotifyEventServer, msg *epb.Payment, c *evt.EventConfig) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return
	}

	if msg.ItemType != ukama.ItemTypePackage.String() {
		log.Errorf("unexpected item type for successful payment: %s", msg.ItemType)
		return
	}

	metadata := map[string]string{}

	err = json.Unmarshal(msg.Metadata, &metadata)
	if err != nil {
		log.Errorf("failed to Unmarshal payment metadata as map[string]string: %v", err)
	}

	targetId, ok := metadata["targetId"]
	if !ok {
		log.Errorf("missing targetId metadata for successful package payment: %s", msg.ItemId)
		return
	}

	es.processEvent(c, es.orgId, "", "", targetId, targetId, jmsg, msg.Id)
}

func handleEventInvoiceGenerate(es *EventToNotifyEventServer, msg *epb.Report, c *evt.EventConfig) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal message for %s to JSON. Error %+v", c.Name, err)
		return
	}
	es.processEvent(c, es.orgId, msg.NetworkId, "", "", "", jmsg, msg.Id)
}

func (es *EventToNotifyEventServer) processEvent(ec *evt.EventConfig, orgId, networkId, nodeId, subscriberId, userId string, msg []byte, rid string) *db.Notification {
	log.Debugf("Processing event OrgId %s NetworkId %s nodeId %s subscriberId %s userId %s", orgId, networkId, nodeId, subscriberId, userId)

	/* Store raw event */
	event := &db.EventMsg{}
	var id uint = 0
	event.Key = ec.Name
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
		ResourceId:   rid,
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
