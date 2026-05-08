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

	log "github.com/sirupsen/logrus"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
)

type SiteControllerEventServer struct {
	orgName string
	s       *SiteControllerServer
	msgbus  mb.MsgBusServiceClient

	epb.UnimplementedEventNotificationServiceServer
}

func NewSiteControllerEventServer(orgName string,
	s *SiteControllerServer,
	msgbus mb.MsgBusServiceClient) *SiteControllerEventServer {

	return &SiteControllerEventServer{
		orgName: orgName,
		s:       s,
		msgbus:  msgbus,
	}
}

func (n *SiteControllerEventServer) EventNotification(ctx context.Context,
	e *epb.Event) (*epb.EventResponse, error) {

	if e == nil {
		return &epb.EventResponse{}, nil
	}

	log.Infof("site-controller: received event routingKey=%s", e.RoutingKey)

	/*
	 *
	 * CNode online / health events will later trigger:
	 *
	 *   1. resolve site from registry
	 *   2. GET /switch/v1/status
	 *   3. GET /switch/v1/ports/policy if policy hash changed
	 *   4. validate/cache policy
	 *   5. derive site state
	 *
	 * For now, site-controller is API driven and event-ready.
	 */

	return &epb.EventResponse{}, nil
}
