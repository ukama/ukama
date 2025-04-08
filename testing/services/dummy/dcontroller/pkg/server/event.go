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
	"math/rand/v2"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	evt "github.com/ukama/ukama/systems/common/events"
	"github.com/ukama/ukama/systems/common/msgbus"
	cpb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
	cenums "github.com/ukama/ukama/testing/common/enums"
	pb "github.com/ukama/ukama/testing/services/dummy/dcontroller/pb/gen"
	"github.com/ukama/ukama/testing/services/dummy/dcontroller/pkg/metrics"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
)
  
 type DControllerEventServer struct {
	 orgName          string
	 server           *DControllerServer
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

		err = n.handleSiteMonitoring(msg)
		if err != nil {
			return nil, err
		}
	case msgbus.PrepareRoute(n.orgName,"request.cloud.local.{{ .Org}}.node.controller.nodefeeder.publish"):
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
func (n *DControllerEventServer) handleSiteMonitoring(msg *epb.EventAddSite) error {
    log.Infof("Handling node assignment event for site: %s", msg.SiteId)
    
    randomConfig := &pb.SiteConfig{
        AvgBackhaulSpeed: 30 + rand.Float64()*70,    
        AvgLatency:       10 + rand.Float64()*40,    
        SolarEfficiency:  0.7 + rand.Float64()*0.2,  
    }
    
    metricsReq := &pb.StartMetricsRequest{
        SiteId:     msg.SiteId,
        SiteConfig: randomConfig,
    }
    
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    resp, err := n.server.StartMetrics(ctx, metricsReq)
    if err != nil {
        log.Errorf("Failed to start metrics for site %s: %v", msg.SiteId, err)
        return err
    }
    
    if !resp.Success {
        log.Warnf("StartMetrics returned unsuccessful for site %s", msg.SiteId)
        return fmt.Errorf("failed to start metrics for site %s", msg.SiteId)
    }
    
    log.Infof("Successfully started metrics for site %s with config: %+v", msg.SiteId, randomConfig)
    
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
	
	log.Infof("Received toggle event for port %d with status %v for site %s", port, status, siteId)
	
	var component string
	switch port {
	case metrics.PORT_NODE: 
	   component = "Node"
	case metrics.PORT_SOLAR:  
	   component = "Solar Controller"
	case metrics.PORT_BACKHAUL:
	   component = "Backhaul"
	default:
	   component = fmt.Sprintf("Unknown(%d)", port)
	}
	
	log.Infof("Toggling %s port for site %s to %v", component, siteId, status)
	
	req := &pb.UpdatePortStatusRequest{
		SiteId:    siteId,
		PortNumber: int32(port),  
		Enabled:   status,        
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	resp, err := n.server.UpdatePortStatus(ctx, req)
	if err != nil {
		log.Errorf("Failed to update port %d status for site %s: %v", port, siteId, err)
		return err
	}
	
	if !resp.Success {
		log.Warnf("UpdatePortStatus returned unsuccessful for site %s: %s", siteId, resp.Message)
		return fmt.Errorf("failed to update port status: %s", resp.Message)
	}
	
	var scenario cenums.SCENARIOS
	if !status { 
		switch port {
		case metrics.PORT_NODE:
			scenario = cenums.SCENARIO_NODE_OFF
		case metrics.PORT_BACKHAUL:
			scenario = cenums.SCENARIO_BACKHAUL_DOWN
		case metrics.PORT_SOLAR:
			scenario = cenums.SCENARIO_DEFAULT
		default:
			scenario = cenums.SCENARIO_DEFAULT
		}
	} else {
		scenario = cenums.SCENARIO_DEFAULT
	}
	
	nodes, err := n.server.nodeClient.GetNodesBySite(siteId)
	if err != nil {
		log.Errorf("Failed to get nodes for site %s: %v", siteId, err)
		return err
	}
	
	for _, node := range nodes.Nodes {
		if err := n.server.dnodeClient.UpdateNodeScenario(node.Id, scenario); err != nil {
			log.Errorf("Failed to update node %s scenario: %v", node.Id, err)
		} else {
			log.Infof("Updated scenario for node %s to %s", node.Id, scenario)
		}
	}

	log.Infof("Successfully updated %s port status to %v for site %s", component, status, siteId)
	return nil
}
