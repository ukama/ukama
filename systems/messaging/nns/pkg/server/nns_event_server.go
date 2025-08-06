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

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/ukama/ukama/systems/common/msgbus"

	log "github.com/sirupsen/logrus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	creg "github.com/ukama/ukama/systems/common/rest/client/registry"
	pb "github.com/ukama/ukama/systems/messaging/nns/pb/gen"
)

type NnsEventServer struct {
	orgName    string
	Nns        *NnsServer
	NodeClient creg.NodeClient
	Org        string
	epb.UnimplementedEventNotificationServiceServer
}

func NewNnsEventServer(orgName string, c creg.NodeClient, s *NnsServer, o string) *NnsEventServer {

	return &NnsEventServer{
		orgName:    orgName,
		NodeClient: c,
		Nns:        s,
		Org:        o,
	}
}

func (l *NnsEventServer) EventNotification(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
	log.Infof("Received a message with Routing key %s and Message %+v", e.RoutingKey, e.Msg)
	switch e.RoutingKey {
	case msgbus.PrepareRoute(l.orgName, "event.cloud.local.{{ .Org}}.messaging.mesh.node.online"):
		msg, err := l.unmarshalNodeOnlineEvent(e.Msg)
		if err != nil {
			return nil, err
		}

		err = l.handleNodeOnlineEvent(e.RoutingKey, msg)
		if err != nil {
			return nil, err
		}

	case msgbus.PrepareRoute(l.orgName, "event.cloud.local.{{ .Org}}.messaging.mesh.node.offline"):
		msg, err := l.unmarshalNodeOfflineEvent(e.Msg)
		if err != nil {
			return nil, err
		}

		err = l.handleNodeOfflineEvent(e.RoutingKey, msg)
		if err != nil {
			return nil, err
		}
	case msgbus.PrepareRoute(l.orgName, "event.cloud.local.{{ .Org}}.registry.node.node.assigned"):
		msg, err := l.unmarshalNodeAssignedEvent(e.Msg)
		if err != nil {
			return nil, err
		}

		err = l.handleNodeAssignedEvent(e.RoutingKey, msg)
		if err != nil {
			return nil, err
		}

	case msgbus.PrepareRoute(l.orgName, "event.cloud.local.{{ .Org}}.registry.node.node.released"):
		msg, err := l.unmarshalNodeReleaseEvent(e.Msg)
		if err != nil {
			return nil, err
		}

		err = l.handleNodeReleaseEvent(e.RoutingKey, msg)
		if err != nil {
			return nil, err
		}

	case msgbus.PrepareRoute(l.orgName, "event.cloud.global.{{ .Org}}.messaging.mesh.ip.update"):
		msg, err := l.unmarshalNodeReleaseEvent(e.Msg)
		if err != nil {
			return nil, err
		}

		err = l.handleNodeReleaseEvent(e.RoutingKey, msg)
		if err != nil {
			return nil, err
		}

	default:
		log.Errorf("No handler routing key %s", e.RoutingKey)
	}

	return &epb.EventResponse{}, nil
}

func (l *NnsEventServer) unmarshalNodeOnlineEvent(msg *anypb.Any) (*epb.NodeOnlineEvent, error) {
	p := &epb.NodeOnlineEvent{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal AddOrgRequest message with : %+v. Error %s.", msg, err.Error())
		return nil, err
	}
	return p, nil
}

func (l *NnsEventServer) handleNodeOnlineEvent(key string, msg *epb.NodeOnlineEvent) error {
	log.Infof("Keys %s and Proto is: %+v", key, msg)

	log.Infof("Getting org and network for %s", msg.GetNodeId())

	nodeInfo, err := l.NodeClient.Get(msg.GetNodeId())
	if err != nil {
		log.Errorf("Failed to get org and network. Error: %+v", err)
		log.Warningf("Node id %s won't have org and network info", msg.GetNodeId())

		nodeInfo = &creg.NodeInfo{
			Id: msg.GetNodeId(),
		}

		nodeInfo.Site = creg.NodeSiteInfo{}
		nodeInfo.Site.NodeId = msg.GetNodeId()
		nodeInfo.Site.SiteId = ""
		nodeInfo.Site.NetworkId = ""
	}

	_, err = l.Nns.Set(context.Background(), &pb.SetNodeIPRequest{
		NodeId:       msg.GetNodeId(),
		NodeIp:       msg.GetMeshIp(),
		MeshIp:       msg.GetMeshIp(),
		NodePort:     msg.GetNodePort(),
		MeshPort:     msg.GetMeshPort(),
		Org:          l.Org,
		Network:      nodeInfo.Site.NetworkId,
		Site:         nodeInfo.Site.SiteId,
		MeshHostName: msg.GetMeshHostName(),
	})

	if err != nil {
		log.Errorf("Failed to set node IP. Error: %+v", err)

		return err
	}

	log.Infof("Node %s IP set to %s", msg.GetNodeId(), msg.GetMeshIp())

	return nil
}

func (l *NnsEventServer) unmarshalNodeOfflineEvent(msg *anypb.Any) (*epb.NodeOfflineEvent, error) {
	p := &epb.NodeOfflineEvent{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal AddOrgRequest message with : %+v. Error %s.", msg, err.Error())
		return nil, err
	}
	return p, nil
}

func (l *NnsEventServer) handleNodeOfflineEvent(key string, msg *epb.NodeOfflineEvent) error {
	log.Infof("Keys %s and Proto is: %+v", key, msg)
	return nil
}

func (l *NnsEventServer) unmarshalNodeAssignedEvent(msg *anypb.Any) (*epb.NodeAssignedEvent, error) {
	p := &epb.NodeAssignedEvent{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal AddOrgRequest message with : %+v. Error %s.", msg, err.Error())
		return nil, err
	}
	return p, nil
}

func (l *NnsEventServer) handleNodeAssignedEvent(key string, msg *epb.NodeAssignedEvent) error {
	log.Infof("Keys %s and Proto is: %+v", key, msg)

	orgNet, err := l.Nns.nodeOrgMapping.Get(context.Background(), msg.GetNodeId())
	if err != nil {
		log.Errorf("node %s doesn't exist. Error %v", msg.GetNodeId(), err)
		return err
	}

	err = l.Nns.nodeOrgMapping.Add(context.Background(), msg.GetNodeId(), l.Org, msg.Network, msg.Site, orgNet.NodeIp, orgNet.MeshHostName, orgNet.NodePort, orgNet.MeshPort)
	if err != nil {
		log.Errorf("failed to update labels for %s. Error %v", msg.GetNodeId(), err)
		return err
	}

	return nil
}

func (l *NnsEventServer) unmarshalNodeReleaseEvent(msg *anypb.Any) (*epb.NodeReleasedEvent, error) {
	p := &epb.NodeReleasedEvent{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal AddOrgRequest message with : %+v. Error %s.", msg, err.Error())
		return nil, err
	}
	return p, nil
}

func (l *NnsEventServer) handleNodeReleaseEvent(key string, msg *epb.NodeReleasedEvent) error {
	log.Infof("Keys %s and Proto is: %+v", key, msg)

	orgNet, err := l.Nns.nodeOrgMapping.Get(context.Background(), msg.GetNodeId())
	if err != nil {
		log.Errorf("node %s doesn't exist. Error %v", msg.GetNodeId(), err)
		return err
	}

	err = l.Nns.nodeOrgMapping.Add(context.Background(), msg.GetNodeId(), l.Org, "", "", orgNet.NodeIp, orgNet.MeshHostName, orgNet.NodePort, orgNet.MeshPort)
	if err != nil {
		log.Errorf("failed to update labels for %s. Error %v", msg.GetNodeId(), err)
		return err
	}

	return nil
}
