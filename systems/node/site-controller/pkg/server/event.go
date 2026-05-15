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
	pbhealth "github.com/ukama/ukama/systems/node/health/pb/gen"
)

type SiteControllerEventServer struct {
	s       *SiteControllerServer
	orgName string
	epb.UnimplementedEventNotificationServiceServer
}

func NewSiteControllerEventServer(orgName string, s *SiteControllerServer) *SiteControllerEventServer {
	return &SiteControllerEventServer{
		s:       s,
		orgName: orgName,
	}
}

func (c *SiteControllerEventServer) EventNotification(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
	log.Infof("Received a message with Routing key %s and Message %+v", e.RoutingKey, e.Msg)

	switch e.RoutingKey {
	case msgbus.PrepareRoute(c.orgName, "event.cloud.local.{{ .Org}}.node.health.report.store"):
		cfg := evt.EventToEventConfig[evt.EventHealthReportStore]
		msg, err := epb.UnmarshalHealthReportEvent(e.Msg, cfg.Name)
		if err != nil {
			return nil, err
		}

		log.Infof("Received a health report event %+v", msg)

		return &epb.EventResponse{}, nil

	case msgbus.PrepareRoute(c.orgName, "event.cloud.local.{{ .Org}}.registry.site.site.create"):
		cfg := evt.EventToEventConfig[evt.EventSiteCreate]
		msg, err := epb.UnmarshalEventAddSite(e.Msg, cfg.Name)
		if err != nil {
			return nil, err
		}

		err = c.handleAddSite(ctx, msg)
		if err != nil {
			return nil, err
		}

		return &epb.EventResponse{}, nil

	default:
		log.Errorf("No handler routing key %s", e.RoutingKey)
	}

	return &epb.EventResponse{}, nil
}

func (s *SiteControllerEventServer) handleAddSite(ctx context.Context, msg *epb.EventAddSite) error {
	log.Infof("Adding site %s with network %s", msg.SiteId, msg.NetworkId)
	healthClient, err := s.s.healthClient.GetClient()
	if err != nil {
		return fmt.Errorf("failed to get health client: %w", err)
	}

	nodes, err := s.s.nodeClient.GetNodesBySite(msg.SiteId)
	if err != nil {
		return fmt.Errorf("failed to get nodes by site: %w", err)
	}

	interfacesMap := make(map[string]*pbhealth.Interface)

	for _, node := range nodes.Nodes {
		interfaces, err := healthClient.ListInterfaces(ctx, &pbhealth.ListInterfacesRequest{
			NodeId: node.Id,
		})
		if err != nil {
			return fmt.Errorf("failed to list interfaces: %w", err)
		}
		interfacesMap[node.Type] = interfaces.Interfaces
	}

	return nil
}