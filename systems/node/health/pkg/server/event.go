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
	"strings"

	log "github.com/sirupsen/logrus"
	evt "github.com/ukama/ukama/systems/common/events"
	"github.com/ukama/ukama/systems/common/msgbus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
)

type HealthEventServer struct {
	s       *HealthServer
	orgName string
	epb.UnimplementedEventNotificationServiceServer
}

func NewHealthEventServer(orgName string, s *HealthServer) *HealthEventServer {
	return &HealthEventServer{
		s:       s,
		orgName: orgName,
	}
}

func (h *HealthEventServer) EventNotification(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
	log.Infof("Received a message with Routing key %s and Message %+v", e.RoutingKey, e.Msg)
	switch e.RoutingKey {
	case msgbus.PrepareRoute(h.orgName, evt.EventRoutingKey[evt.EventNodeDelete]):
		c := evt.EventToEventConfig[evt.EventNodeDelete]
		msg, err := epb.UnmarshalEventRegistryNodeDelete(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}

		err = h.handleNodeDeleteEvent(msg)
		if err != nil {
			return nil, err
		}

	default:
		log.Errorf("No handler routing key %s", e.RoutingKey)
	}

	return &epb.EventResponse{}, nil
}

func (h *HealthEventServer) handleNodeDeleteEvent(msg *epb.EventRegistryNodeDelete) error {
	nodeId := strings.ToLower(msg.NodeId)
	if nodeId == "" {
		return fmt.Errorf("node delete event has no node id")
	}

	if err := h.s.sRepo.DeleteNode(nodeId); err != nil {
		log.Errorf("Failed to delete health data of node %s: %v", nodeId, err)
		return err
	}
	log.Infof("Deleted health data of node %s", nodeId)

	return nil
}
