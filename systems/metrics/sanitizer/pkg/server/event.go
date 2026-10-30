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
	"fmt"

	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/metrics/sanitizer/pkg"

	log "github.com/sirupsen/logrus"
	evt "github.com/ukama/ukama/systems/common/events"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	cpb "github.com/ukama/ukama/systems/common/pb/events"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
)

type SanitizerEventServer struct {
	orgName        string
	ss             *SanitizerServer
	msgbus         mb.MsgBusServiceClient
	baseRoutingKey msgbus.RoutingKeyBuilder
	epb.UnimplementedEventNotificationServiceServer
}

func NewSanitizerEventServer(orgName string, ss *SanitizerServer, msgBus mb.MsgBusServiceClient) *SanitizerEventServer {
	return &SanitizerEventServer{
		ss:      ss,
		orgName: orgName,
		msgbus:  msgBus,
		baseRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().
			SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
	}
}

func (se *SanitizerEventServer) EventNotification(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
	log.Infof("Received a message with Routing key %s and Message %+v", e.RoutingKey, e.Msg)

	switch e.RoutingKey {
	case msgbus.PrepareRoute(se.orgName, evt.EventRoutingKey[evt.EventNodeAssign]):
		msg, err := cpb.UnmarshalProtoEvent[epb.EventRegistryNodeAssign](e.Msg)
		if err != nil {
			log.Errorf("Failed to unmarshal node assign event: %v", err)

			return nil, err
		}

		err = se.handleNodeAssignEvent(e.RoutingKey, msg)
		if err != nil {
			return nil, err
		}

	case msgbus.PrepareRoute(se.orgName, evt.EventRoutingKey[evt.EventNodeRelease]):
		msg, err := cpb.UnmarshalProtoEvent[epb.NodeReleasedEvent](e.Msg)
		if err != nil {
			log.Errorf("Failed to unmarshal node release event: %v", err)

			return nil, err
		}

		err = se.handleNodeReleaseEvent(e.RoutingKey, msg)
		if err != nil {
			return nil, err
		}

	case msgbus.PrepareRoute(se.orgName, evt.EventRoutingKey[evt.EventNodeDelete]):
		msg, err := cpb.UnmarshalProtoEvent[epb.EventRegistryNodeDelete](e.Msg)
		if err != nil {
			log.Errorf("Failed to unmarshal node delete event: %v", err)

			return nil, err
		}

		err = se.handleNodeDeleteEvent(e.RoutingKey, msg)
		if err != nil {
			return nil, err
		}

	default:
		log.Errorf("No handler routing key %s", e.RoutingKey)

		return nil, fmt.Errorf("no handler routing key %s", e.RoutingKey)
	}

	return &epb.EventResponse{}, nil
}

func (se *SanitizerEventServer) handleNodeAssignEvent(key string, msg *epb.EventRegistryNodeAssign) error {
	return se.ss.syncNodeCache()
}

func (se *SanitizerEventServer) handleNodeReleaseEvent(key string, msg *epb.NodeReleasedEvent) error {
	return se.ss.syncNodeCache()
}

func (se *SanitizerEventServer) handleNodeDeleteEvent(key string, msg *epb.EventRegistryNodeDelete) error {
	return se.ss.syncNodeCache()
}
