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
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/msgbus"
	cpb "github.com/ukama/ukama/systems/common/pb/events"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
)

type ControllerEventServer struct {
	s       *ControllerServer
	orgName string
	epb.UnimplementedEventNotificationServiceServer
}

func NewControllerEventServer(orgName string, s *ControllerServer) *ControllerEventServer {
	return &ControllerEventServer{
		s:       s,
		orgName: orgName,
	}
}

func (n *ControllerEventServer) EventNotification(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
	log.Infof("Received a message with Routing key %s and Message %+v", e.RoutingKey, e.Msg)
	switch e.RoutingKey {
	case msgbus.PrepareRoute(n.orgName, "event.cloud.local.{{ .Org}}.registry.node.node.create"):
		msg, err := cpb.UnmarshalProtoEvent[epb.EventRegistryNodeCreate](e.Msg)
		if err != nil {
			log.Errorf("Failed to unmarshal EventRegistryNodeCreate: %v", err)
			return nil, err
		}

		err = n.handleRegistryNodeAddEvent(e.RoutingKey, msg)
		if err != nil {
			return nil, err
		}

	case msgbus.PrepareRoute(n.orgName, "event.cloud.local.{{ .Org}}.registry.node.node.delete"):
		msg, err := cpb.UnmarshalProtoEvent[epb.EventRegistryNodeDelete](e.Msg)
		if err != nil {
			log.Errorf("Failed to unmarshal EventRegistryNodeDelete: %v", err)
			return nil, err
		}

		err = n.handleRegistryNodeDeleteEvent(msg)
		if err != nil {
			return nil, err
		}

	default:
		log.Errorf("No handler routing key %s", e.RoutingKey)
	}

	return &epb.EventResponse{}, nil
}

func (n *ControllerEventServer) handleRegistryNodeAddEvent(key string, msg *epb.EventRegistryNodeCreate) error {
	log.Infof("Keys %s and Proto is: %+v", key, msg)

	err := n.s.nRepo.Add(strings.ToLower(msg.NodeId))
	if err != nil {
		log.Errorf("Error adding node %s to controller/nodeLog repo.Error: %+v", msg.NodeId, err)
		return err
	}
	return nil
}

func (n *ControllerEventServer) handleRegistryNodeDeleteEvent(msg *epb.EventRegistryNodeDelete) error {
	nodeId := strings.ToLower(msg.NodeId)
	if nodeId == "" {
		return fmt.Errorf("node delete event has no node id")
	}

	if err := n.s.nRepo.Delete(nodeId); err != nil {
		log.Errorf("Error deleting node %s from controller/nodeLog repo. Error: %+v", nodeId, err)
		return err
	}
	return nil
}
