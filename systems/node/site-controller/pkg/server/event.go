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
	"time"

	log "github.com/sirupsen/logrus"
	evt "github.com/ukama/ukama/systems/common/events"
	"github.com/ukama/ukama/systems/common/msgbus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	creg "github.com/ukama/ukama/systems/common/rest/client/registry"
	"github.com/ukama/ukama/systems/common/ukama"
	hpb "github.com/ukama/ukama/systems/node/health/pb/gen"
	pb "github.com/ukama/ukama/systems/node/site-controller/pb/gen"
	"github.com/ukama/ukama/systems/node/site-controller/pkg"
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

type SiteControllerEventServer struct {
	s                    *SiteControllerServer
	orgName              string
	sites                db.SiteRepo
	intents              db.IntentRepo
	flights              db.IntentFlightRepo
	states               db.StateRepo
	components           db.ComponentRepo
	componentSyncDelay   time.Duration
	componentSyncTimeout time.Duration
	epb.UnimplementedEventNotificationServiceServer
}

func NewSiteControllerEventServer(
	s *SiteControllerServer,
	sites db.SiteRepo,
	intents db.IntentRepo,
	flights db.IntentFlightRepo,
	states db.StateRepo,
	components db.ComponentRepo,
	config *pkg.Config,
) *SiteControllerEventServer {
	return &SiteControllerEventServer{
		s:                    s,
		orgName:              config.OrgName,
		sites:                sites,
		intents:              intents,
		flights:              flights,
		states:               states,
		components:           components,
		componentSyncDelay:   config.ComponentSyncDelay,
		componentSyncTimeout: config.Timeout,
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
		err = c.handleHealthReport(ctx, msg)
		if err != nil {
			return nil, err
		}

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

func (c *SiteControllerEventServer) handleHealthReport(ctx context.Context, msg *epb.HealthReportEvent) error {
	log.Infof("Received a health report event %+v", msg)
	hClient, err := c.s.healthClient.GetClient()
	if err != nil {
		return err
	}

	nodeId := msg.NodeId
	nodeType := msg.NodeType
	if _, err := c.s.nodeClient.Get(nodeId); err != nil {
		return fmt.Errorf("failed to get node %s: %w", nodeId, err)
	}

	nodes, err := c.getNodesBySite("", nodeId)
	if err != nil {
		return fmt.Errorf("failed to get nodes for node %s: %w", nodeId, err)
	}
	if len(nodes.Nodes) == 0 {
		return fmt.Errorf("no nodes found for node %s", nodeId)
	}
	if nodes.Nodes[0] == nil {
		return fmt.Errorf("node %s is nil", nodeId)
	}
	if nodes.Nodes[0].Site.SiteId == "" {
		return fmt.Errorf("no site found for node %s", nodeId)
	}
	siteId := nodes.Nodes[0].Site.SiteId

	report, err := hClient.ListInterfaces(ctx, &hpb.ListInterfacesRequest{
		NodeId: nodeId,
	})
	if err != nil {
		return fmt.Errorf("failed to list interfaces for node %s: %w", nodeId, err)
	}

	if report.Interfaces == nil {
		return fmt.Errorf("no interfaces found for node %s", nodeId)
	}

	switch nodeType {
	case ukama.NODE_ID_TYPE_TOWERNODE:
		// TODO: Handle tower node health report
		return nil

	case ukama.NODE_ID_TYPE_AMPNODE:
		if report.Interfaces.Radio == nil {
			return fmt.Errorf("no radio found for node %s", nodeId)
		}

		radioState := report.Interfaces.Radio.State
		_, err = c.s.SetRadio(ctx, &pb.SetRadioRequest{
			SiteId: siteId,
			State:  radioState,
		})
		if err != nil {
			return err
		}
		return nil

	case ukama.NODE_ID_TYPE_CNODE:
		// TODO: Handle dnode health report
		return nil

	default:
		return fmt.Errorf("unsupported node type %s", nodeType)
	}
}

func (c *SiteControllerEventServer) handleAddSite(ctx context.Context, msg *epb.EventAddSite) error {
	log.Infof("Adding site %s with network %s", msg.SiteId, msg.NetworkId)

	if err := c.sites.Ensure(msg.SiteId); err != nil {
		return err
	}

	expiresAt := time.Now().UTC().Add(time.Minute)
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

	c.scheduleSetSiteComponents(msg.SiteId)
	return nil
}

func (c *SiteControllerEventServer) scheduleSetSiteComponents(siteID string) {
	delay := c.componentSyncDelay
	go func() {
		log.Infof("site-controller: scheduling component sync for site %s in %s", siteID, delay)
		time.Sleep(delay)

		ctx, cancel := context.WithTimeout(context.Background(), c.componentSyncTimeout)
		defer cancel()

		if err := c.setSiteComponents(ctx, siteID); err != nil {
			log.Warnf("site-controller: failed to set site components for site %s: %v", siteID, err)
			return
		}
		log.Infof("site-controller: updated site components for site %s", siteID)
	}()
}

func (c *SiteControllerEventServer) getNodesBySite(siteID string, nodeId string) (*creg.ListNodesResponse, error) {
	nodes, err := c.s.nodeClient.List(creg.ListNodesRequest{SiteId: siteID, NodeId: nodeId})
	if err != nil {
		return nil, fmt.Errorf("failed to list nodes for site %s: %w", siteID, err)
	}
	return nodes, nil
}

func (c *SiteControllerEventServer) setSiteComponents(ctx context.Context, siteID string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	nodes, err := c.getNodesBySite(siteID, "")
	if err != nil {
		return err
	}
	components := make([]string, 0, len(nodes.Nodes))
	for _, node := range nodes.Nodes {
		if node == nil || node.Id == "" {
			continue
		}
		components = append(components, node.Id)
	}
	return c.components.Upsert(&db.SiteComponent{
		SiteID:     siteID,
		Components: components,
	})
}
