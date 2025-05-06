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
	"encoding/json"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	evt "github.com/ukama/ukama/systems/common/events"
	"github.com/ukama/ukama/systems/common/msgbus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
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
	case msgbus.PrepareRoute(n.orgName, "event.cloud.local.{{ .Org}}.messaging.mesh.node.online"):
		c := evt.EventToEventConfig[evt.EventNodeOnline]
		msg, err := epb.UnmarshalNodeOnlineEvent(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}

		err = n.handleNodeOnlineEvent(e.RoutingKey, msg)
		if err != nil {
			return nil, err
		}
	case msgbus.PrepareRoute(n.orgName, "event.cloud.local.{{ .Org}}.messaging.mesh.node.offline"):
		c := evt.EventToEventConfig[evt.EventNodeOffline]
		msg, err := epb.UnmarshalNodeOfflineEvent(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}

		err = n.handleNodeOfflineEvent(e.RoutingKey, msg)
		if err != nil {
			return nil, err
		}
	case msgbus.PrepareRoute(n.orgName, "event.cloud.local.{{ .Org}}.node.state.node.transition"):
		c := evt.NodeEventToEventConfig[evt.NodeStateTransition]
		msg, err := epb.UnmarshalEventNodeStateTransition(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}
		err = n.handleNodeStateTransitionEvent(e.RoutingKey, msg)
		if err != nil {
			return nil, err
		}
	case msgbus.PrepareRoute(n.orgName, "event.cloud.local.{{ .Org}}.node.notify.notification.store"):
		c := evt.EventToEventConfig[evt.EventPaymentFailed]
		msg, err := epb.UnmarshalNotification(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}
		err = n.handleNotifyEvent(ctx, e.RoutingKey, msg)
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
	default:
		log.Errorf("No handler routing key %s", e.RoutingKey)
	}

	return &epb.EventResponse{}, nil
}

func (n *NodeEventServer) handleNotifyEvent(ctx context.Context, key string, msg *epb.Notification) error {
	log.Infof("Keys %s and Proto is: %+v", key, msg)
	var details map[string]interface{}
	if err := json.Unmarshal(msg.Details, &details); err != nil {
		log.WithError(err).Error("Failed to unmarshal details")
		return err
	}
	lat := details["latitude"]
	lon := details["longitude"]
	if lat == nil || lon == nil {
		log.Errorf("Latitude or Longitude key not found in details")
		return fmt.Errorf("latitude or longitude key not found in details")
	}

	updateRequest := &pb.UpdateNodeRequest{
		NodeId:    msg.NodeId,
		Latitude:  lat.(float64),
		Longitude: lon.(float64),
	}

	_, err := n.s.UpdateNode(ctx, updateRequest)
	if err != nil {
		log.WithError(err).Error("Failed to update node")
		return err
	}
	return nil
}

// so, commenting for compiling.
func (n *NodeEventServer) handleNodeOnlineEvent(key string, msg *epb.NodeOnlineEvent) error {
	log.Infof("Keys %s and Proto is: %+v", key, msg)

	nodeID := strings.ToLower(msg.GetNodeId())
	node, err := n.s.GetNode(context.Background(), &pb.GetNodeRequest{NodeId: nodeID})

	// Handle error cases except NotFound
	if err != nil && status.Code(err) != codes.NotFound {
		log.Errorf("error getting the node: %v", err)
		return err
	}

	// If node doesn't exist, create it
	if node == nil {
		if err := n.createNewNode(nodeID); err != nil {
			return err
		}
	}

	// Update node status
	return n.updateNodeConnectivity(nodeID, ukama.Online.String(), node)
}

func (n *NodeEventServer) handleNodeOfflineEvent(key string, msg *epb.NodeOfflineEvent) error {
	log.Infof("Keys %s and Proto is: %+v", key, msg)

	nodeID := msg.GetNodeId()
	node, err := n.s.GetNode(context.Background(), &pb.GetNodeRequest{NodeId: nodeID})
	if err != nil {
		log.Errorf("error getting the node: %v", err)
		return err
	}

	if node != nil {
		return n.updateNodeConnectivity(nodeID, ukama.Offline.String(), node)
	}
	return nil
}

func (n *NodeEventServer) createNewNode(nodeID string) error {
	req := &pb.AddNodeRequest{
		NodeId: nodeID,
		Name:   nodeID[len(nodeID)-7:],
	}
	_, err := n.s.AddNode(context.Background(), req)
	return err
}

func (n *NodeEventServer) updateNodeConnectivity(nodeID, connectivity string, node *pb.GetNodeResponse) error {
	state := ukama.Unknown.String()
	if node != nil {
		state = node.Node.Status.State.Enum().String()
	}

	_, err := n.s.UpdateNodeStatus(context.Background(), &pb.UpdateNodeStateRequest{
		NodeId:       nodeID,
		Connectivity: connectivity,
		State:        state,
	})
	return err
}
func (n *NodeEventServer) handleNodeStateTransitionEvent(key string, msg *epb.NodeStateChangeEvent) error {
    log.Infof("Keys %s and Proto is: %+v", key, msg)    
	nodeID := strings.ToLower(msg.GetNodeId())

    node, err := n.s.GetNode(context.Background(), &pb.GetNodeRequest{NodeId: nodeID})
    if err != nil {
        log.Errorf("error getting the node: %v", err)
        return err
    }
    
    log.Infof("Updating node %s with current state: %v", node.Node.Id, node.Node.Status.State, msg.State, msg.Substate)
    
    _, err = n.s.UpdateNodeStatus(context.Background(), &pb.UpdateNodeStateRequest{
        NodeId:       nodeID,
        Connectivity: ukama.Online.String(),
        State:       msg.State,
    })
    if err != nil {
        log.Errorf("error updating node status: %v", err)
        return err
    }
    
    return nil
}