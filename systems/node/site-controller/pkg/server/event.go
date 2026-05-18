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
	"time"

	log "github.com/sirupsen/logrus"
	evt "github.com/ukama/ukama/systems/common/events"
	"github.com/ukama/ukama/systems/common/msgbus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	"github.com/ukama/ukama/systems/node/site-controller/pkg/db"
)

var defaultSiteIntent = db.SiteIntent{
	DesiredService: "off",
	DesiredRadio:   "off",
	RequestedBy:    "system",
}

var defaultSiteState = db.SiteState{
	PowerState:   "unknown",
	ServiceState: "unknown",
	RadioState:   "unknown",
	AccessState:  "unavailable",
	Reason:       "site_created",
}

var defaultSiteComponent = db.SiteComponent{
	Components: []string{"unknown"},
}

type SiteControllerEventServer struct {
	s          *SiteControllerServer
	orgName    string
	sites      db.SiteRepo
	intents    db.IntentRepo
	flights    db.IntentFlightRepo
	states     db.StateRepo
	components db.ComponentRepo
	epb.UnimplementedEventNotificationServiceServer
}

func NewSiteControllerEventServer(
	orgName string,
	s *SiteControllerServer,
	sites db.SiteRepo,
	intents db.IntentRepo,
	flights db.IntentFlightRepo,
	states db.StateRepo,
	components db.ComponentRepo,
) *SiteControllerEventServer {
	return &SiteControllerEventServer{
		s:          s,
		orgName:    orgName,
		sites:      sites,
		intents:    intents,
		flights:    flights,
		states:     states,
		components: components,
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

func (c *SiteControllerEventServer) handleAddSite(ctx context.Context, msg *epb.EventAddSite) error {
	log.Infof("Adding site %s with network %s", msg.SiteId, msg.NetworkId)

	if err := c.sites.Ensure(msg.SiteId); err != nil {
		return err
	}

	expiresAt := time.Now().UTC().Add(time.Hour)
	intent := &db.SiteIntent{
		SiteID:         msg.SiteId,
		DesiredService: defaultSiteIntent.DesiredService,
		DesiredRadio:   defaultSiteIntent.DesiredRadio,
		Reason:         "site_created",
		RequestedBy:    defaultSiteIntent.RequestedBy,
	}
	if err := c.intents.Upsert(intent); err != nil {
		return err
	}

	flight := &db.SiteIntentFlight{
		SiteIntentID: intent.ID,
		Status:       db.IntentFlightStatusPending,
		RetryCount:   0,
		ExpiresAt:    expiresAt,
	}
	if err := c.flights.Upsert(flight); err != nil {
		return err
	}

	state := &db.SiteState{
		SiteID:       msg.SiteId,
		PowerState:   defaultSiteState.PowerState,
		ServiceState: defaultSiteState.ServiceState,
		RadioState:   defaultSiteState.RadioState,
		AccessState:  defaultSiteState.AccessState,
		Reason:       defaultSiteState.Reason,
	}
	if err := c.states.Upsert(state); err != nil {
		return err
	}

	component := &db.SiteComponent{
		SiteID:     msg.SiteId,
		Components: defaultSiteComponent.Components,
	}
	return c.components.Upsert(component)
}
