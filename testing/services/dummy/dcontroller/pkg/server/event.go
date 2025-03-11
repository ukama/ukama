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
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	evt "github.com/ukama/ukama/systems/common/events"
	"github.com/ukama/ukama/systems/common/msgbus"
	cpb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
	pb "github.com/ukama/ukama/testing/services/dummy/dcontroller/pb/gen"
	"github.com/ukama/ukama/testing/services/dummy/dcontroller/pkg/metrics"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
)
  
 type DControllerEventServer struct {
	 orgName          string
	 server           *DControllerServer
	 controllerClient pb.MetricsControllerClient 
	 epb.UnimplementedEventNotificationServiceServer
 }
  
 func NewEventServer(orgName string, server *DControllerServer) *DControllerEventServer {
	 return &DControllerEventServer{
		 orgName:        orgName,
		 server:         server,
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
 
	 case "request.cloud.local.{{ .Org}}.node.controller.nodefeeder.publish":
		 nodeMsg := &cpb.NodeFeederMessage{}
		 if err := anypb.UnmarshalTo(e.Msg, nodeMsg, proto.UnmarshalOptions{}); err != nil {
			 log.Errorf("Failed to unmarshal to NodeFeederMessage: %v", err)
			 return nil, err
		 }
		 
		 err := n.handleToggleSwitchEventDirect(nodeMsg)
		 if err != nil {
			 log.Errorf("Error handling toggle switch event: %v", err)
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
		 SiteId:  msg.SiteId,
		 Profile: pb.Profile_PROFILE_NORMAL, 
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
 
 func (n *DControllerEventServer) handleToggleSwitchEventDirect(msg *cpb.NodeFeederMessage) error {
	 log.Infof("Handling toggle switch event: target=%s, path=%s", msg.Target, msg.Path)
	 
	 targetParts := strings.Split(msg.Target, ".")
	 siteId := targetParts[len(targetParts)-1]
	 
	 path := strings.TrimPrefix(msg.Path, "/v1/switch/")
	 pathParts := strings.Split(path, "/")
	 
	 if len(pathParts) != 2 {
		 return fmt.Errorf("invalid path format: %s", msg.Path)
	 }
	 
	 port, err := strconv.Atoi(pathParts[0])
	 if err != nil {
		 return fmt.Errorf("invalid port number: %w", err)
	 }
	 
	 status, err := strconv.ParseBool(pathParts[1])
	 if err != nil {
		 return fmt.Errorf("invalid status value: %w", err)
	 }
	 
	 var component string
	 switch port {
	 case metrics.PORT_AMPLIFIER:
		component = "Amplifier"
	 case metrics.PORT_TOWER:
		component = "Tower"
	 case metrics.PORT_SOLAR:
		component = "Solar"
	 case metrics.PORT_BACKHAUL:
		component = "Backhaul"
	 default:
		component = fmt.Sprintf("Unknown(%d)", port)
	 }
	 
	 log.Infof("Toggling %s port for site %s to %v", component, siteId, status)
	 
	 req := &pb.UpdateMetricsRequest{
		 SiteId: siteId,
		 PortUpdates: []*pb.PortUpdate{
			 {
				 PortNumber: int32(port),
				 Status:     status,
			 },
		 },
	 }
	 
	 ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	 defer cancel()
	 
	 resp, err := n.server.UpdateMetrics(ctx, req)
	 if err != nil {
		 log.Errorf("Failed to update port %d status for site %s: %v", port, siteId, err)
		 return err
	 }
	 
	 if !resp.Success {
		 log.Warnf("UpdateMetrics returned unsuccessful for site %s: %s", siteId, resp.Message)
		 return fmt.Errorf("failed to update port status: %s", resp.Message)
	 }
	 
	 log.Infof("Successfully updated %s port status to %v for site %s", component, status, siteId)
	 return nil
 }
 
 func (n *DControllerEventServer) handleToggleSwitchEvent(eventMsg []byte) error {
	 msg := &cpb.NodeFeederMessage{}
	 if err := proto.Unmarshal(eventMsg, msg); err != nil {
		 return fmt.Errorf("failed to unmarshal NodeFeederMessage: %w", err)
	 }
	 
	 return n.handleToggleSwitchEventDirect(msg)
 }