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
	"github.com/ukama/ukama/systems/messaging/nns/pkg"
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
	case msgbus.PrepareRoute(l.orgName, "event.cloud.local.{{ .Org}}.registry.node.node.assign"):
		msg, err := l.unmarshalNodeAssignedEvent(e.Msg)
		if err != nil {
			return nil, err
		}

		err = l.handleNodeAssignedEvent(e.RoutingKey, msg)
		if err != nil {
			return nil, err
		}

	case msgbus.PrepareRoute(l.orgName, "event.cloud.local.{{ .Org}}.registry.node.node.release"):
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

	network, site := l.resolveNodeLineage(msg.GetNodeId())

	_, err := l.Nns.Set(context.Background(), &pb.SetRequest{
		NodeId:       msg.GetNodeId(),
		NodeIp:       msg.GetNodeIp(),
		MeshIp:       msg.GetMeshIp(),
		NodePort:     msg.GetNodePort(),
		MeshPort:     msg.GetMeshPort(),
		Network:      network,
		Site:         site,
		MeshHostName: msg.GetMeshHostName(),
	})

	if err != nil {
		log.Errorf("Failed to set node IP. Error: %+v", err)

		return err
	}

	log.Infof("Node %s IP set to %s", msg.GetNodeId(), msg.GetMeshIp())

	return nil
}

// resolveNodeLineage returns the network and site to store for a node. Registry
// is authoritative, including when it reports no site: an unattached node
// legitimately has empty lineage. When registry cannot be reached the currently
// stored lineage is kept, so a registry outage does not blank out the mapping
// of every node that reconnects during it.
func (l *NnsEventServer) resolveNodeLineage(nodeId string) (network, site string) {
	log.Infof("Getting org and network for %s", nodeId)

	nodeInfo, err := l.NodeClient.Get(nodeId)
	if err == nil && nodeInfo != nil {
		return nodeInfo.Site.NetworkId, nodeInfo.Site.SiteId
	}

	if err != nil {
		log.Errorf("Failed to get org and network for %s. Error: %+v", nodeId, err)
	} else {
		log.Errorf("Registry returned no node info for %s", nodeId)
	}

	stored, sErr := l.Nns.nns.Get(context.Background(), nodeId)
	if sErr != nil || stored == nil {
		log.Warningf("Node id %s won't have org and network info", nodeId)

		return "", ""
	}

	log.Warningf("Keeping stored network %q and site %q for node %s", stored.Network, stored.Site, nodeId)

	return stored.Network, stored.Site
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

func (l *NnsEventServer) unmarshalNodeAssignedEvent(msg *anypb.Any) (*epb.EventRegistryNodeAssign, error) {
	p := &epb.EventRegistryNodeAssign{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to unmarshal EventRegistryNodeAssign message with : %+v. Error %s.", msg, err.Error())
		return nil, err
	}
	return p, nil
}

func (l *NnsEventServer) unmarshalNodeReleaseEvent(msg *anypb.Any) (*epb.NodeReleasedEvent, error) {
	p := &epb.NodeReleasedEvent{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to unmarshal NodeReleasedEvent message with : %+v. Error %s.", msg, err.Error())
		return nil, err
	}
	return p, nil
}

func (l *NnsEventServer) handleNodeAssignedEvent(key string, msg *epb.EventRegistryNodeAssign) error {
	log.Infof("Keys %s and Proto is: %+v", key, msg)

	return l.updateNodeLineage(msg.GetNodeId(), msg.GetNetwork(), msg.GetSite())
}

func (l *NnsEventServer) handleNodeReleaseEvent(key string, msg *epb.NodeReleasedEvent) error {
	log.Infof("Keys %s and Proto is: %+v", key, msg)

	return l.updateNodeLineage(msg.GetNodeId(), "", "")
}

// updateNodeLineage rewrites the network/site of an already known node, leaving
// its mesh and node addressing untouched. A node that has never been online is
// not an error: it has no mapping yet, and handleNodeOnlineEvent resolves its
// lineage from registry when it first connects.
func (l *NnsEventServer) updateNodeLineage(nodeId, network, site string) error {
	orgNet, err := l.Nns.nns.Get(context.Background(), nodeId)
	if err != nil {
		log.Warningf("Skipping lineage update for node %s: no mesh mapping yet. Error %v", nodeId, err)

		return nil
	}

	if orgNet.Network == network && orgNet.Site == site {
		return nil
	}

	obj := pkg.NodeMeshMap{
		NodeId:       nodeId,
		NodeIp:       orgNet.NodeIp,
		NodePort:     orgNet.NodePort,
		MeshIp:       orgNet.MeshIp,
		MeshHostName: orgNet.MeshHostName,
		MeshPort:     orgNet.MeshPort,
		Org:          l.orgName,
		Network:      network,
		Site:         site,
	}

	if err := l.Nns.nns.Add(context.Background(), obj); err != nil {
		log.Errorf("failed to update labels for %s. Error %v", nodeId, err)

		return err
	}

	return nil
}
