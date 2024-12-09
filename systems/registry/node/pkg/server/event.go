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

	log "github.com/sirupsen/logrus"
	evt "github.com/ukama/ukama/systems/common/events"
	"github.com/ukama/ukama/systems/common/msgbus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	"github.com/ukama/ukama/systems/common/ukama"
	pb "github.com/ukama/ukama/systems/registry/node/pb/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
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
		msg, err := n.unmarshalNodeOnlineEvent(e.Msg)
		if err != nil {
			return nil, err
		}

		err = n.handleNodeOnlineEvent(e.RoutingKey, msg)
		if err != nil {
			return nil, err
		}
	case msgbus.PrepareRoute(n.orgName, "event.cloud.local.{{ .Org}}.messaging.mesh.node.offline"):
		msg, err := n.unmarshalNodeOfflineEvent(e.Msg)
		if err != nil {
			return nil, err
		}

		err = n.handleNodeOfflineEvent(e.RoutingKey, msg)
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

func (n *NodeEventServer) unmarshalNodeOnlineEvent(msg *anypb.Any) (*epb.NodeOnlineEvent, error) {
	p := &epb.NodeOnlineEvent{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal NodeOnline  message with : %+v. Error %s.", msg, err.Error())
		return nil, err
	}
	return p, nil
}

// so, commenting for compiling.
func (n *NodeEventServer) handleNodeOnlineEvent(key string, msg *epb.NodeOnlineEvent) error {
	log.Infof("Keys %s and Proto is: %+v", key, msg)

	node, err := n.s.GetNode(context.Background(), &pb.GetNodeRequest{NodeId: msg.GetNodeId()})
	if err != nil {
		log.Warnf("Error code %v for getting node %s", status.Code(err), msg.GetNodeId())
		if status.Code(err) != codes.NotFound {
			log.Error("error getting the node" + err.Error())
			return err
		}
	}

	/* Add node if you can't find a node */
	if node == nil {
		req := &pb.AddNodeRequest{
			NodeId: msg.GetNodeId(),
		}
		_, err = n.s.AddNode(context.Background(), req)
		if err != nil {
			return err
		}
	}

	/* Update node status */
	_, err = n.s.UpdateNodeStatus(context.Background(), &pb.UpdateNodeStateRequest{
		NodeId:       msg.GetNodeId(),
		Connectivity: ukama.Online.String(),
	})
	if err != nil {
		return err
	}

	return nil
}

func (n *NodeEventServer) unmarshalNodeOfflineEvent(msg *anypb.Any) (*epb.NodeOfflineEvent, error) {
	p := &epb.NodeOfflineEvent{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal NodeOffline message with : %+v. Error %s.", msg, err.Error())
		return nil, err
	}
	return p, nil
}

func (n *NodeEventServer) handleNodeOfflineEvent(key string, msg *epb.NodeOfflineEvent) error {
	log.Infof("Keys %s and Proto is: %+v", key, msg)

	node, err := n.s.GetNode(context.Background(), &pb.GetNodeRequest{NodeId: msg.GetNodeId()})
	if err != nil {
		log.Error("error getting the node" + err.Error())
		return err
	}

	if node != nil {
		/* Update node status */
		_, err = n.s.UpdateNodeStatus(context.Background(), &pb.UpdateNodeStateRequest{
			NodeId:       node.Node.Id,
			Connectivity: ukama.Offline.String(),
		})
		if err != nil {
			return err
		}
	}

	return nil
}
