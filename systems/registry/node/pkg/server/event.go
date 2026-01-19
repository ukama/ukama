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

	log "github.com/sirupsen/logrus"
	evt "github.com/ukama/ukama/systems/common/events"
	"github.com/ukama/ukama/systems/common/msgbus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	npb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
	"github.com/ukama/ukama/systems/common/ukama"
	pb "github.com/ukama/ukama/systems/registry/node/pb/gen"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type NodeEventServer struct {
	s       *NodeServer
	orgName string
	epb.UnimplementedEventNotificationServiceServer
}

func NewNodeEventServer(orgName string, s *NodeServer) *NodeEventServer {
	return &NodeEventServer{
		s:       s,
		orgName: orgName,
	}
}

func (n *NodeEventServer) EventNotification(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
	log.Infof("Received a message with Routing key %s and Message %+v", e.RoutingKey, e.Msg)

	switch e.RoutingKey {
	case msgbus.PrepareRoute(n.orgName, "event.cloud.local.{{ .Org}}.node.state.node.transition"):
		c := evt.NodeEventToEventConfig[evt.NodeStateTransition]
		msg, err := epb.UnmarshalNodeStateChangeEvent(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}
		err = n.handleNodeStateTransitionEvent(ctx, e.RoutingKey, msg)
		if err != nil {
			return nil, err
		}
	case msgbus.PrepareRoute(n.orgName, "event.cloud.local.{{ .Org}}.node.health.capps.store"):
		c := evt.EventToEventConfig[evt.EventHealthCappStore]
		msg, err := epb.UnmarshalStoreRunningAppsInfoEvent(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}
		err = n.handleStoreRunningAppsInfoEvent(ctx, e.RoutingKey, msg)
		if err != nil {
			return nil, err
		}
	case msgbus.PrepareRoute(n.orgName, "event.cloud.local.{{ .Org}}.registry.site.site.create"):
		c := evt.EventToEventConfig[evt.EventPaymentFailed]
		msg, err := epb.UnmarshalEventAddSite(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}

		err = n.s.addNodeToSite(msg.AccessId, msg.SiteId, msg.NetworkId)
		if err != nil {
			return nil, err
		}
	case msgbus.PrepareRoute(n.orgName, "event.cloud.local.{{ .Org}}.inventory.component.node.added"):
		msg, err := epb.UnmarshalEventInventoryNodeComponentAdd(e.Msg, "EventInventoryComponentNodeAdded")
		if err != nil {
			return nil, err
		}
		err = n.handleAddNode(ctx, e.RoutingKey, msg)
		if err != nil {
			return nil, err
		}
	default:
		log.Errorf("No handler routing key %s", e.RoutingKey)
	}

	return &epb.EventResponse{}, nil
}

func (n *NodeEventServer) handleStoreRunningAppsInfoEvent(ctx context.Context, key string, msg *epb.StoreRunningAppsInfoEvent) error {
	log.Infof("Processing store running apps info event: %s, nodeID: %s, timestamp: %s",
		key, msg.NodeId, msg.Timestamp)

	node, err := n.s.GetNode(ctx, &pb.GetNodeRequest{NodeId: msg.NodeId})
	if err != nil {
		log.Errorf("Failed to get node: %v", err)
		return fmt.Errorf("failed to get node: %w", err)
	}
	if node == nil || node.Node == nil {
		log.Errorf("Node not found")
		return fmt.Errorf("node not found")
	}
	
	coordinates := ""
	for _, system := range msg.System {
		if system.Name == "coordinates" {
			coordinates = system.Value
			break
		}
	}
	if coordinates == "" {
		log.Errorf("Coordinates not found")
		return fmt.Errorf("coordinates not found")
	}
	lat, lon, err := parseCoordinates(coordinates)
	if err != nil {
		log.Errorf("Failed to parse coordinates: %v", err)
		return fmt.Errorf("failed to parse coordinates: %w", err)
	}
	if node.Node.Latitude == lat && node.Node.Longitude == lon {
		log.Infof("Node %s already has latitude=%f, longitude=%f",
			msg.NodeId, lat, lon)
		return nil
	}
	log.Infof("Updating node %s: latitude=%f, longitude=%f",
		msg.NodeId, lat, lon)
	updateRequest := &pb.UpdateNodeRequest{
		NodeId:    msg.NodeId,
		Latitude:  lat,
		Longitude: lon,
	}
	_, err = n.s.UpdateNode(ctx, updateRequest)
	if err != nil {
		log.WithError(err).Error("Failed to update node")
		return fmt.Errorf("failed to update node: %w", err)
	}
	return nil
}

func (n *NodeEventServer) handleNodeStateTransitionEvent(ctx context.Context, key string, msg *epb.NodeStateChangeEvent) error {
	log.Infof("Processing state transition event: %s, nodeID: %s, state: %s, substate: %s",
		key, msg.GetNodeId(), msg.State, msg.Substate)

	nodeID := strings.ToLower(msg.GetNodeId())

	node, err := n.s.GetNode(ctx, &pb.GetNodeRequest{NodeId: nodeID})

	if err != nil && status.Code(err) == codes.NotFound {
		log.Infof("Node %s not found, creating new node", nodeID)
		req := &pb.AddNodeRequest{
			NodeId: nodeID,
		}
		_, err = n.s.AddNode(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to create node: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("error retrieving node: %w", err)
	}

	var connectivity string

	switch strings.ToLower(msg.Substate) {
	case "on":
		connectivity = npb.NodeConnectivity_Online.String()
	case "off":
		connectivity = npb.NodeConnectivity_Offline.String()
	default:
		if node != nil && node.Node != nil && node.Node.Status != nil {
			connectivity = node.Node.Status.Connectivity.String()
		} else {
			connectivity = npb.NodeConnectivity_Undefined.String()
		}
	}

	log.Infof("Updating node %s: connectivity=%s, state=%s",
		nodeID, connectivity, msg.State)

	_, err = n.s.UpdateNodeStatus(ctx, &pb.UpdateNodeStateRequest{
		NodeId:       nodeID,
		Connectivity: connectivity,
		State:        msg.State,
	})
	if err != nil {
		return fmt.Errorf("failed to update node status: %w", err)
	}

	return nil
}

func (n *NodeEventServer) handleAddNode(ctx context.Context, key string, msg *epb.EventInventoryNodeComponentAdd) error {
	log.Infof("Processing add node event: %s, nodeID: %s, nodeType: %s",
		key, msg.PartNumber, msg.Type)

	nodeID, err := ukama.ValidateNodeId(msg.PartNumber)
	if err != nil {
		return fmt.Errorf("invalid node id: %w", err)
	}

	_, err = n.s.AddNode(ctx, &pb.AddNodeRequest{
		NodeId:    nodeID.StringLowercase(),
		Name:      ukama.GetPlaceholderNameByType(nodeID.GetNodeType()),
		Latitude:  0,
		Longitude: 0,
	})
	if err != nil {
		return fmt.Errorf("failed to add node: %w", err)
	}
	return nil
}

func parseCoordinates(coordinates string) (float64, float64, error) {
	if coordinates == "" {
		return 0, 0, fmt.Errorf("coordinates string is empty")
	}

	parts := strings.Split(coordinates, ",")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid coordinates format: expected 'lat,lon', got %d parts", len(parts))
	}

	latStr := strings.TrimSpace(parts[0])
	lonStr := strings.TrimSpace(parts[1])

	if latStr == "" || lonStr == "" {
		return 0, 0, fmt.Errorf("invalid coordinates: latitude or longitude is empty")
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse latitude '%s': %w", latStr, err)
	}

	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse longitude '%s': %w", lonStr, err)
	}

	return lat, lon, nil
}
