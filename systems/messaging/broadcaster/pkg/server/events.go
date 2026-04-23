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
		case msgbus.PrepareRoute(n.orgName, "event.cloud.local.{{ .Org}}.ukamaagent.asr.policies.publish"):
			c := evt.EventToEventConfig[evt.EventSiteCreate]
			msg, err := epb.UnmarshalBroadcasterEvent(e.Msg, c.Name)
			if err != nil {
				log.Errorf("Failed to unmarshal broadcaster event: %+v", err)
				return nil, err
			}
			
			if msg.Type == epb.BroadcastType_NODE_BROADCAST {
				err = n.s.NodeFeederBroadcast(ctx, msg)
				if err != nil {
					log.Errorf("Failed to handle broadcast node feeder event: %+v", err)
					return nil, err
				}
			} else {
				log.Errorf("No handler broadcast type %s", msg.Type)
				return nil, fmt.Errorf("no handler broadcast type %s", msg.Type)
			}
			return &epb.EventResponse{}, nil
		default:
			log.Errorf("No handler routing key %s", e.RoutingKey)
			return nil, fmt.Errorf("no handler routing key %s", e.RoutingKey)
		}

}
