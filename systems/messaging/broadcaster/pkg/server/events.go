/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package server

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	evt "github.com/ukama/ukama/systems/common/events"
	"github.com/ukama/ukama/systems/common/msgbus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
)

const logMsgRoutingKey = "Received a message with Routing key %s and Message %+v"

type BroadcasterEventServer struct {
	s       *BroadcasterServer
	orgName string
	epb.UnimplementedEventNotificationServiceServer
}

func NewBroadcasterEventServer(orgName string, s *BroadcasterServer) *BroadcasterEventServer {
	return &BroadcasterEventServer{
		s:       s,
		orgName: orgName,
	}
}

func (n *BroadcasterEventServer) EventNotification(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
	log.Infof(logMsgRoutingKey, e.RoutingKey, e.Msg)
	switch e.RoutingKey {
		case msgbus.PrepareRoute(n.orgName, "event.cloud.local.{{ .Org}}.messaging.broadcast.nodefeeder"):
			c := evt.EventToEventConfig[evt.EventSiteCreate]
			msg, err := epb.UnmarshalBroadcasterEvent(e.Msg, c.Name)
			if err != nil {
				log.Errorf("Failed to unmarshal broadcaster event: %+v", err)
				return nil, err
			}
	
			err = n.handleBroadcastNodeFeeder(ctx, msg)
			if err != nil {
				log.Errorf("Failed to handle broadcast node feeder event: %+v", err)
				return nil, err
			}
	
			return &epb.EventResponse{}, nil
		default:
			log.Errorf("No handler routing key %s", e.RoutingKey)
	}

	return &epb.EventResponse{}, nil
}


func (n *BroadcasterEventServer) handleBroadcastNodeFeeder(ctx context.Context, msg *epb.BroadcasterEvent) error {
	log.Infof("Broadcasting node feeder event: %+v", msg)

	switch msg.Scope {
		case epb.BroadcastScope_ORGANIZATIONAL_SCOPE:
			err := n.s.BroadcastOrganizational(ctx, msg)
			if err != nil {
				log.Errorf("Failed to broadcast organizational event: %+v", err)
				return err
		}
		default:
			log.Errorf("No handler scope %s", msg.Scope)
			return fmt.Errorf("no handler scope %s", msg.Scope)
	}
	return nil
}