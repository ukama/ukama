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

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"

	log "github.com/sirupsen/logrus"
	cpb "github.com/ukama/ukama/systems/common/pb/events"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	npb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
	cinvent "github.com/ukama/ukama/systems/common/rest/client/inventory"
	pb "github.com/ukama/ukama/systems/registry/node/pb/gen"
)

const errFailedAddNodeToSiteFmt = "failed to add node to site: %v"
const errFailedUpdateNodeStatusFmt = "failed to update node status: %v"

type NodeEventServer struct {
	s         *NodeServer
	orgName   string
	invClient cinvent.ComponentClient
	epb.UnimplementedEventNotificationServiceServer
}

func NewNodeEventServer(orgName string, s *NodeServer,
	invClient cinvent.ComponentClient) *NodeEventServer {
	return &NodeEventServer{
		s:         s,
		orgName:   orgName,
		invClient: invClient,
	}
}

func (n *NodeEventServer) EventNotification(ctx context.Context,
	e *epb.Event) (*epb.EventResponse, error) {
	log.Infof("Received a message with Routing key %s and Message %+v", e.RoutingKey, e.Msg)

	switch e.RoutingKey {
	case msgbus.PrepareRoute(n.orgName,
		"event.cloud.local.{{ .Org}}.node.state.node.transition"):
		msg, err := cpb.UnmarshalProtoEvent[epb.NodeStateChangeEvent](e.Msg)
		if err != nil {
			return nil, err
		}
		err = n.handleNodeStateTransitionEvent(ctx, e.RoutingKey, msg)
		if err != nil {
			return nil, err
		}
	case msgbus.PrepareRoute(n.orgName,
		"event.cloud.local.{{ .Org}}.node.health.report.store"):
		msg, err := cpb.UnmarshalProtoEvent[epb.HealthReportEvent](e.Msg)
		if err != nil {
			return nil, err
		}
		err = n.handleHealthReportEvent(ctx, e.RoutingKey, msg)
		if err != nil {
			return nil, err
		}
	case msgbus.PrepareRoute(n.orgName,
		"event.cloud.local.{{ .Org}}.registry.site.site.create"):
		msg, err := cpb.UnmarshalProtoEvent[epb.EventAddSite](e.Msg)
		if err != nil {
			return nil, err
		}

		err = n.handleAddNodeToSite(ctx, msg.AccessId, msg.SiteId, msg.NetworkId)
		if err != nil {
			return nil, err
		}

		return &epb.EventResponse{}, nil

	case msgbus.PrepareRoute(n.orgName,
		"event.cloud.local.{{ .Org}}.inventory.component.node.added"):
		msg, err := cpb.UnmarshalProtoEvent[epb.EventInventoryNodeComponentAdd](e.Msg)
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

func (n *NodeEventServer) handleHealthReportEvent(ctx context.Context, key string,
	msg *epb.HealthReportEvent) error {
	log.Infof("Processing health report event: %s, nodeID: %s",
		key, msg.NodeId)

	node, err := n.s.GetNode(ctx, &pb.GetNodeRequest{NodeId: msg.NodeId})
	if status.Code(err) == codes.NotFound {
		log.Warnf("Ignoring health report for unknown node %s", msg.NodeId)
		return nil
	}
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

	interfaces, err := n.s.healthClient.GetInterfaces(msg.NodeId, msg.Id)
	if err != nil {
		log.Errorf("Failed to get interfaces: %v", err)
		return fmt.Errorf("failed to get interfaces: %w", err)
	}

	if interfaces.Gps == nil {
		log.Errorf("GPS not found: %+v", interfaces.Gps)
		return fmt.Errorf("GPS not found: %+v", interfaces.Gps)
	}

	if !interfaces.Gps.Lock {
		log.Infof("Node %s has no GPS lock yet, skipping location update",
			msg.NodeId)
		return nil
	}

	if !interfaces.Gps.Available {
		log.Infof("Node %s has GPS not available, skipping location update",
			msg.NodeId)
		return nil
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

func (n *NodeEventServer) handleNodeStateTransitionEvent(ctx context.Context, key string,
	msg *epb.NodeStateChangeEvent) error {
	log.Infof("Processing state transition event: %s, nodeID: %s, state: %s, substate: %s",
		key, msg.GetNodeId(), msg.State, msg.Substate)

	nodeID := strings.ToLower(msg.GetNodeId())

	// Without an on/off substate the event says nothing about connectivity, so leave it
	// unchanged.
	var connectivity string

	switch strings.ToLower(msg.Substate) {
	case "on":
		connectivity = npb.NodeConnectivity_Online.String()
	case "off":
		connectivity = npb.NodeConnectivity_Offline.String()
	}

	log.Infof("Updating node %s: connectivity=%s, state=%s",
		nodeID, connectivity, msg.State)

	return n.updateNodeStatus(ctx, msg.State+" transition", &pb.UpdateNodeStateRequest{
		NodeId:       nodeID,
		Connectivity: connectivity,
		State:        msg.State,
	})
}

// updateNodeStatus writes an event's status; NotFound means the node isn't registered, so
// the event is skipped.
func (n *NodeEventServer) updateNodeStatus(ctx context.Context, source string,
	req *pb.UpdateNodeStateRequest) error {
	_, err := n.s.UpdateNodeStatus(ctx, req)
	if status.Code(err) == codes.NotFound {
		log.Warnf("Ignoring %s for unregistered node %s", source, req.NodeId)
		return nil
	}
	if err != nil {
		log.Errorf(errFailedUpdateNodeStatusFmt, err)
		return fmt.Errorf(errFailedUpdateNodeStatusFmt, err)
	}

	return nil
}

func (n *NodeEventServer) handleAddNode(ctx context.Context, key string,
	msg *epb.EventInventoryNodeComponentAdd) error {
	log.Infof("Processing add node event: %s, nodeID: %s, nodeType: %s",
		key, msg.PartNumber, msg.Type)

	nodeID, err := ukama.ValidateNodeId(msg.PartNumber)
	if err != nil {
		return fmt.Errorf("invalid node id: %w", err)
	}

	_, err = n.s.AddNode(ctx, &pb.AddNodeRequest{
		NodeId: nodeID.StringLowercase(),
		Name:   ukama.GetPlaceholderNameByType(nodeID.GetNodeType()),
	})
	if status.Code(err) == codes.AlreadyExists {
		log.Infof("Node %s is already registered", nodeID.StringLowercase())
		return nil
	}
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
		return "", "", fmt.Errorf(
			"invalid coordinates format: expected 'lat,lon', got %d parts", len(parts))
	}

	latStr := strings.TrimSpace(parts[0])
	lonStr := strings.TrimSpace(parts[1])

	if latStr == "" || lonStr == "" {
		return "", "", fmt.Errorf("invalid coordinates: latitude or longitude is empty")
	}

	return latStr, lonStr, nil
}

func (n *NodeEventServer) handleAddNodeToSite(ctx context.Context, accessId string,
	siteID string, networkID string) error {
	log.Infof("Adding node with access id %s to site %s with network %s", accessId, siteID,
		networkID)

	component, err := n.invClient.Get(accessId)
	if err != nil {
		return fmt.Errorf("failed to get component: %w", err)
	}

	tNodeId, err := ukama.ValidateNodeId(component.PartNumber)
	if err != nil {
		return fmt.Errorf("invalid access node id %s: %w", component.PartNumber, err)
	}

	if tNodeId.GetNodeType() != ukama.NODE_ID_TYPE_TOWERNODE {
		return fmt.Errorf("access node %s is not a tower node", tNodeId.StringLowercase())
	}

	nodeID := tNodeId.StringLowercase()

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

	siteUUID, err := uuid.FromString(siteID)
	if err != nil {
		return fmt.Errorf("invalid site id %s: %w", siteID, err)
	}

	networkUUID, err := uuid.FromString(networkID)
	if err != nil {
		return fmt.Errorf("invalid network id %s: %w", networkID, err)
	}

	nodeIds := []string{tNodeId.StringLowercase(), aNodeId.StringLowercase(),
		cNodeId.StringLowercase()}

	log.Infof("Adding nodes %v to site %s with network %s", nodeIds, siteID, networkID)

	err = n.s.assignNodesToSite(siteUUID, networkUUID, nodeIds)
	if err != nil {
		log.Errorf("No node of %v was added to site %s: %v", nodeIds, siteID, err)
		return fmt.Errorf(errFailedAddNodeToSiteFmt, err)
	}

	return nil
}
