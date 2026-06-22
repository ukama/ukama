/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package server

import (
	"context"

	log "github.com/sirupsen/logrus"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	creg "github.com/ukama/ukama/systems/common/rest/client/registry"
	"github.com/ukama/ukama/systems/messaging/broadcaster/pkg"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type BroadcasterServer struct {
	broadcasterRoutingKey msgbus.RoutingKeyBuilder
	msgbus                mb.MsgBusServiceClient
	debug                 bool
	orgName               string
	nodeClient            creg.NodeClient
}

func NewBroadcasterServer(orgName string, msgBus mb.MsgBusServiceClient, nodeClient creg.NodeClient, debug bool) *BroadcasterServer {
	return &BroadcasterServer{
		debug:                 debug,
		msgbus:                msgBus,
		orgName:               orgName,
		broadcasterRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
		nodeClient:            nodeClient,
	}
}

func (b *BroadcasterServer) NodeFeederBroadcast(ctx context.Context, msg *epb.BroadcasterEvent) error {
	log.Infof("Broadcasting node feeder event: %+v", msg)

	nfMsg := &epb.NodeFeederMessage{}
	err := proto.Unmarshal(msg.Msg, nfMsg)
	if err != nil {
		log.Errorf("Failed to unmarshal broadcaster event: %+v", err)
		return err
	}

	listReq := creg.ListNodesRequest{
		NetworkId:    "",
		SiteId:       "",
		Connectivity: "",
		State:        "",
		NodeId:       "",
		Type:         "",
	}

	switch msg.Scope {
	case epb.BroadcastScope_UNKNOWN_SCOPE:
		log.Errorf("Unknown broadcast scope: %s, Publishing msg as it is.", msg.Scope)
		return b.publishMessage(msg.RoutingKey, nfMsg)
	case epb.BroadcastScope_NETWORK_SCOPE:
		listReq.NetworkId = msg.TargetId
	case epb.BroadcastScope_SITE_SCOPE:
		listReq.SiteId = msg.TargetId
	}

	nodes, err := b.nodeClient.List(listReq)
	if err != nil {
		log.Errorf("Failed to get nodes: %+v", err)
		return err
	}

	for _, node := range nodes.Nodes {
		nfMsg.NodeId = node.Id
		nfMsg.Target = b.orgName + "." + "*" + "." + "*" + "." + node.Id
		err = b.publishMessage(msg.RoutingKey, nfMsg)
		if err != nil {
			log.Errorf("Failed to publish message: %+v for node %s", err, node.Id)
			return err
		}
	}

	return nil
}

func (b *BroadcasterServer) publishMessage(routeKey string, msg protoreflect.ProtoMessage) error {
	log.Infof("Published message on route %s ", routeKey)
	err := b.msgbus.PublishRequest(routeKey, msg)
	return err
}
