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

	log "github.com/sirupsen/logrus"
	evt "github.com/ukama/ukama/systems/common/events"
	"github.com/ukama/ukama/systems/common/msgbus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
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

	return nil
}