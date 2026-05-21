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
	evt "github.com/ukama/ukama/systems/common/events"
	"github.com/ukama/ukama/systems/common/msgbus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	npb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
	cinvent "github.com/ukama/ukama/systems/common/rest/client/inventory"
	"github.com/ukama/ukama/systems/common/ukama"
	pb "github.com/ukama/ukama/systems/registry/node/pb/gen"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const errFailedAddNodeToSiteFmt = "failed to add node to site: %v"
const errFailedGetNodeFmt = "failed to get node: %v"
const errFailedUpdateNodeStatusFmt = "failed to update node status: %v"
const errNodeNotFoundFmt = "node %s not found"

type NodeEventServer struct {
	s       *NodeServer
	orgName string
	invClient cinvent.ComponentClient
	epb.UnimplementedEventNotificationServiceServer
}

func NewNodeEventServer(orgName string, s *NodeServer, invClient cinvent.ComponentClient) *NodeEventServer {
	return &NodeEventServer{
		s:       s,
		orgName: orgName,
		invClient: invClient,
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
		c := evt.EventToEventConfig[evt.EventHealthReportStore]
		msg, err := epb.UnmarshalHealthReportEvent(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}
		err = n.handleHealthReportEvent(ctx, e.RoutingKey, msg)
		if err != nil {
			return nil, err
		}
	case msgbus.PrepareRoute(n.orgName, "event.cloud.local.{{ .Org}}.registry.site.site.create"):
		c := evt.EventToEventConfig[evt.EventSiteCreate]
		msg, err := epb.UnmarshalEventAddSite(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}

		err = n.handleAddNodeToSite(ctx, msg.AccessId, msg.SiteId, msg.NetworkId)
		if err != nil {
			return nil, err
		}

		return &epb.EventResponse{}, nil

	case msgbus.PrepareRoute(n.orgName, "event.cloud.local.{{ .Org}}.inventory.component.node.added"):
		msg, err := epb.UnmarshalEventInventoryNodeComponentAdd(e.Msg, "EventInventoryComponentNodeAdded")
		if err != nil {
			return nil, err
		}
		err = n.handleAddNode(ctx, e.RoutingKey, msg)
		if err != nil {
			return nil, err
		}

	case msgbus.PrepareRoute(n.orgName, evt.NodeStateEventRoutingKey[evt.NodeStateEventOnline]):
		c := evt.NodeEventToEventConfig[evt.NodeStateEventOnline]
		msg, err := epb.UnmarshalNodeOnlineEvent(e.Msg, c.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal node online event: %w", err)
		}
		err = n.handleNodeOnlineEvent(ctx, e.RoutingKey, msg)
		if err != nil {
			return nil, err
		}
	case msgbus.PrepareRoute(n.orgName, evt.NodeStateEventRoutingKey[evt.NodeStateEventOffline]):
		c := evt.NodeEventToEventConfig[evt.NodeStateEventOffline]
		msg, err := epb.UnmarshalNodeOfflineEvent(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}
		err = n.handleNodeOfflineEvent(ctx, e.RoutingKey, msg)
		if err != nil {
			return nil, err
		}
	default:
		log.Errorf("No handler routing key %s", e.RoutingKey)
	}

	return &epb.EventResponse{}, nil
}

func (n *NodeEventServer) handleNodeOnlineEvent(ctx context.Context, key string, msg *epb.NodeOnlineEvent) error {
	log.Infof("Processing node online event: %s, nodeID: %s", key, msg.NodeId)

	node, err := n.s.GetNode(ctx, &pb.GetNodeRequest{NodeId: msg.NodeId})
	if err != nil {
		log.Errorf(errFailedGetNodeFmt, err)
		return fmt.Errorf(errFailedGetNodeFmt, err)
	}
	if node == nil || node.Node == nil {
		log.Errorf(errNodeNotFoundFmt, msg.NodeId)
		return fmt.Errorf(errNodeNotFoundFmt, msg.NodeId)
	}	

	_, err = n.s.UpdateNodeStatus(ctx, &pb.UpdateNodeStateRequest{
		NodeId:       msg.NodeId,
		Connectivity: npb.NodeConnectivity_Online.String(),
		State:        node.Node.Status.State.String(),
	})
	if err != nil {
		log.Errorf(errFailedUpdateNodeStatusFmt, err)
		return fmt.Errorf(errFailedUpdateNodeStatusFmt, err)
	}
	return nil
}

func (n *NodeEventServer) handleNodeOfflineEvent(ctx context.Context, key string, msg *epb.NodeOfflineEvent) error {
	log.Infof("Processing node offline event: %s, nodeID: %s", key, msg.NodeId)

	node, err := n.s.GetNode(ctx, &pb.GetNodeRequest{NodeId: msg.NodeId})
	if err != nil {
		log.Errorf(errFailedGetNodeFmt, err)
		return fmt.Errorf(errFailedGetNodeFmt, err)
	}
	if node == nil || node.Node == nil {
		log.Errorf(errNodeNotFoundFmt, msg.NodeId)
		return fmt.Errorf(errNodeNotFoundFmt, msg.NodeId)
	}

	_, err = n.s.UpdateNodeStatus(ctx, &pb.UpdateNodeStateRequest{
		NodeId:       msg.NodeId,
		Connectivity: npb.NodeConnectivity_Offline.String(),
		State:        node.Node.Status.State.String(),
	})

	if err != nil {
		log.Errorf(errFailedUpdateNodeStatusFmt, err)
		return fmt.Errorf(errFailedUpdateNodeStatusFmt, err)
	}
	return nil
}

func (n *NodeEventServer) handleHealthReportEvent(ctx context.Context, key string, msg *epb.HealthReportEvent) error {
	log.Infof("Processing health report event: %s, nodeID: %s",
		key, msg.NodeId)

	node, err := n.s.GetNode(ctx, &pb.GetNodeRequest{NodeId: msg.NodeId})
	if err != nil {
		log.Errorf("Failed to get node: %v", err)
		return fmt.Errorf("failed to get node: %w", err)
	}
	if node == nil || node.Node == nil {
		log.Errorf("Node not found")
		return fmt.Errorf("node not found")
	}

	if node.Node.Type != ukama.NODE_ID_TYPE_TOWERNODE {
		log.Infof("Node %s is not a tower node", msg.NodeId)
		return nil
	}

	interfaces, err := n.s.healthClient.GetInterfaces("", msg.NodeId, msg.Id)
	if err != nil {
		log.Errorf("Failed to get interfaces: %v", err)
		return fmt.Errorf("failed to get interfaces: %w", err)
	}

	coordinates := interfaces.Gps.Coordinates
	
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
		log.Infof("Node %s already has latitude=%s, longitude=%s",
			msg.NodeId, lat, lon)
		return nil
	}
	log.Infof("Updating node %s: latitude=%s, longitude=%s",
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
	})
	if err != nil {
		return fmt.Errorf("failed to add node: %w", err)
	}
	return nil
}

func parseCoordinates(coordinates string) (string, string, error) {
	if coordinates == "" {
		return "", "", fmt.Errorf("coordinates string is empty")
	}

	parts := strings.Split(coordinates, ",")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid coordinates format: expected 'lat,lon', got %d parts", len(parts))
	}

	latStr := strings.TrimSpace(parts[0])
	lonStr := strings.TrimSpace(parts[1])

	if latStr == "" || lonStr == "" {
		return "", "", fmt.Errorf("invalid coordinates: latitude or longitude is empty")
	}

	return latStr, lonStr, nil
}

func (n *NodeEventServer) handleAddNodeToSite(ctx context.Context, accessId string, siteID string, networkID string) error {
	log.Infof("Adding node with access id %s to site %s with network %s", accessId, siteID, networkID)

	component, err := n.invClient.Get(accessId)
	if err != nil {
		return fmt.Errorf("failed to get component: %w", err)
	}
	nodeID := component.PartNumber
	
	log.Infof("Node ID is: %s", nodeID)
	
	aId, err := ukama.GetANodeIdFromTNodeId(nodeID)
	if err != nil {
		return fmt.Errorf("failed to get A Node ID: %w", err)
	}

	aNodeId, err := ukama.ValidateNodeId(aId.StringLowercase())
	if err != nil {
		return fmt.Errorf("failed to validate A Node ID: %w", err)
	}
	
	cId, err := ukama.GetCNodeIdFromTNodeId(nodeID)
	if err != nil {
		return fmt.Errorf("failed to get C Node ID: %w", err)
	}

	cNodeId, err := ukama.ValidateNodeId(cId.StringLowercase())
	if err != nil {
		return fmt.Errorf("failed to validate C Node ID: %w", err)
	}

	// Add Tower Node to Site
	log.Infof("Adding Tower Node %s to Site %s with Network %s", nodeID, siteID, networkID)
	err = n.addNodeToSite(ctx, nodeID, siteID, networkID)
	if err != nil {
		return fmt.Errorf(errFailedAddNodeToSiteFmt, err)
	}

	// Add Amplifier Node to Site
	log.Infof("Adding Amplifier Node %s to Site %s with Network %s", aNodeId.StringLowercase(), siteID, networkID)
	err = n.addNodeToSite(ctx, aNodeId.StringLowercase(), siteID, networkID)
	if err != nil {
		return fmt.Errorf(errFailedAddNodeToSiteFmt, err)
	}

	// Add Controller Node to Site
	log.Infof("Adding Controller Node %s to Site %s with Network %s", cNodeId.StringLowercase(), siteID, networkID)
	err = n.addNodeToSite(ctx, cNodeId.StringLowercase(), siteID, networkID)
	if err != nil {
		return fmt.Errorf(errFailedAddNodeToSiteFmt, err)
	}
	return nil
}

func (n *NodeEventServer) addNodeToSite(ctx context.Context, nodeID string, siteID string, networkID string) error {
	log.Infof("Adding node %s to site %s with network %s", nodeID, siteID, networkID)
	node, err := n.s.GetNode(ctx, &pb.GetNodeRequest{NodeId: nodeID})
	if err != nil {
		return fmt.Errorf("failed to get node: %w", err)
	}
	if node == nil || node.Node == nil {
		return fmt.Errorf("node not found")
	}
	err = n.s.addNodeToSiteServer(nodeID, siteID, networkID)
	if err != nil {
		return fmt.Errorf(errFailedAddNodeToSiteFmt, err)
	}
	return nil
}
