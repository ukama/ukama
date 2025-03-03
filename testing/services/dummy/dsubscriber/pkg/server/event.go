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

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/msgbus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
)

type DsubEventServer struct {
	orgName string
	server  *DsubscriberServer
	epb.UnimplementedEventNotificationServiceServer
}

func NewDsubEventServer(orgName string, server *DsubscriberServer) *DsubEventServer {
	return &DsubEventServer{
		orgName: orgName,
		server:  server,
	}
}

func (l *DsubEventServer) EventNotification(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
	log.Infof("Received a message with Routing key %s and Message %+v", e.RoutingKey, e.Msg)
	switch e.RoutingKey {
	case msgbus.PrepareRoute(l.orgName, "event.cloud.local.{{ .Org}}.subscriber.simmanager.sim.allocate"):
		msg, err := epb.UnmarshalEventSimAllocation(e.Msg, "EventSimAllocate")
		if err != nil {
			log.Errorf("Failed to unmarshal EventSimAllocate: %v", err)
			return nil, err
		}

		log.Infof("Received a message with Routing key %s and Message %+v", e.RoutingKey, msg)

		l.server.startHandler(msg.Iccid, msg.PackageEndDate.AsTime().Format(time.RFC3339))

	case msgbus.PrepareRoute(l.orgName, "event.cloud.local.{{ .Org}}.subscriber.simmanager.sim.activepackage"):
		msg, err := epb.UnmarshalEventSimActivePackage(e.Msg, "EventSimActivePackage")
		if err != nil {
			log.Errorf("Failed to unmarshal EventSimActivePackage: %v", err)
			return nil, err
		}

		log.Infof("Received a message with Routing key %s and Message %+v", e.RoutingKey, msg)

		l.server.updateHandler(msg.Iccid, msg.PackageEndDate.AsTime().Format(time.RFC3339))

	case msgbus.PrepareRoute(l.orgName, "event.cloud.local.{{ .Org}}.messaging.mesh.node.offline"):
		msg, err := epb.UnmarshalNodeOfflineEvent(e.Msg, "EventNodeOffline")
		if err != nil {
			log.Errorf("Failed to unmarshal NodeOfflineEvent: %v", err)
			return nil, err
		}

		log.Infof("Received a message with Routing key %s and Message %+v", e.RoutingKey, msg)

		l.server.toggleUsageGenerationByNodeId(msg.NodeId, false)
	case msgbus.PrepareRoute(l.orgName, "event.cloud.local.{{ .Org}}.messaging.mesh.node.online"):
		msg, err := epb.UnmarshalNodeOnlineEvent(e.Msg, "EventNodeOnline")
		if err != nil {
			log.Errorf("Failed to unmarshal NodeOnlineEvent: %v", err)
			return nil, err
		}

		log.Infof("Received a message with Routing key %s and Message %+v", e.RoutingKey, msg)

		l.server.toggleUsageGenerationByNodeId(msg.NodeId, true)
	case msgbus.PrepareRoute(l.orgName, "event.cloud.local.{{ .Org}}.subscriber.simmanager.sim.deactivate"):
		msg, err := epb.UnmarshalEventSimDeactivation(e.Msg, "EventSimDeactivation")
		if err != nil {
			log.Errorf("Failed to unmarshal EventSimDeactivation: %v", err)
			return nil, err
		}

		log.Infof("Received a message with Routing key %s and Message %+v", e.RoutingKey, msg)

		l.server.toggleUsageGenerationByIccid(msg.Iccid, false)
	case msgbus.PrepareRoute(l.orgName, "event.cloud.local.{{ .Org}}.subscriber.simmanager.sim.activate"):
		msg, err := epb.UnmarshalEventSimActivation(e.Msg, "EventSimActivation")
		if err != nil {
			log.Errorf("Failed to unmarshal EventSimActivation: %v", err)
			return nil, err
		}

		log.Infof("Received a message with Routing key %s and Message %+v", e.RoutingKey, msg)

		l.server.toggleUsageGenerationByIccid(msg.Iccid, true)
	default:
		log.Errorf("handler not registered for %s", e.RoutingKey)
	}

	return &epb.EventResponse{}, nil
}
