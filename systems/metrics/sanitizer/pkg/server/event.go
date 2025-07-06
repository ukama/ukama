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

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/metrics/sanitizer/pkg"

	log "github.com/sirupsen/logrus"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
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
	case msgbus.PrepareRoute(se.orgName, "event.cloud.local.{{ .Org}}.registry.node.node.assign"):
		msg, err := unmarshalNodeAssignEvent(e.Msg)
		if err != nil {
			return nil, err
		}

		err = se.handleNodeAssignEvent(e.RoutingKey, msg)
		if err != nil {
			return nil, err
		}

	case msgbus.PrepareRoute(se.orgName, "event.cloud.local.{{ .Org}}.registry.node.node.release"):
		msg, err := unmarshalNodeReleaseEvent(e.Msg)
		if err != nil {
			return nil, err
		}

		err = se.handleNodeReleaseEvent(e.RoutingKey, msg)
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

func unmarshalNodeAssignEvent(msg *anypb.Any) (*epb.EventRegistryNodeAssign, error) {
	p := &epb.EventRegistryNodeAssign{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal registry node assign message with : %+v. Error %s.", msg, err.Error())

		return nil, err
	}
	return p, nil
}

func unmarshalNodeReleaseEvent(msg *anypb.Any) (*epb.NodeReleasedEvent, error) {
	p := &epb.NodeReleasedEvent{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal node release message with : %+v. Error %s.", msg, err.Error())

		return nil, err
	}
	return p, nil
}
