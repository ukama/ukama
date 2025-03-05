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
	"time"

	log "github.com/sirupsen/logrus"
	evt "github.com/ukama/ukama/systems/common/events"
	"github.com/ukama/ukama/systems/common/msgbus"
	pb "github.com/ukama/ukama/testing/services/dummy/controller/pb/gen"

	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
)
 
 type DControllerEventServer struct {
	 orgName        string
	 server  *ControllerServer
	 controllerClient pb.MetricsControllerClient 
	 epb.UnimplementedEventNotificationServiceServer
 }
 
 func NewEventServer(orgName string, server *ControllerServer) *DControllerEventServer {
	 return &DControllerEventServer{
		 orgName:        orgName,
		 server:  server,
	 }
 }
 
 func (n *DControllerEventServer) EventNotification(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
	 log.Infof("Received a message with Routing key %s and Message %+v", e.RoutingKey, e.Msg)
	 switch e.RoutingKey {
 
	 case msgbus.PrepareRoute(n.orgName, evt.EventRoutingKey[evt.EventSiteCreate]):
		 c := evt.EventToEventConfig[evt.EventSiteCreate]
		 msg, err := epb.UnmarshalEventAddSite(e.Msg, c.Name)
		 if err != nil {
			 return nil, err
		 }
 
		 err = n.handleSiteCreateEvent(msg, c.Title)
		 if err != nil {
			 return nil, err
		 }
	 default:
		 log.Errorf("No handler routing key %s", e.RoutingKey)
	 }
 
	 return &epb.EventResponse{}, nil
 }
 
 func (n *DControllerEventServer) handleSiteCreateEvent(msg *epb.EventAddSite, name string) error {
    log.Infof("Handling site create event for site ID: %s", msg.SiteId)

    req := &pb.StartMetricsRequest{
        SiteId:   msg.SiteId,
        Profile:  pb.Profile_PROFILE_NORMAL, 
        Scenario: pb.Scenario_SCENARIO_DEFAULT, 
    }
    
  
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    resp, err := n.controllerClient.StartMetrics(ctx, req)
    if err != nil {
        log.Errorf("Failed to start metrics for site %s: %v", msg.SiteId, err)
        return err
    }
    
    if !resp.Success {
        log.Warnf("StartMetrics returned unsuccessful for site %s", msg.SiteId)
    } else {
        log.Infof("Successfully started metrics for site %s", msg.SiteId)
    }
    
    return nil
}