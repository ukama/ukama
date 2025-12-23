/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2025-present, Ukama Inc.
 */

package server

import (
	"context"

	log "github.com/sirupsen/logrus"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
	factory "github.com/ukama/ukama/systems/common/rest/client/factory"
	lookupclient "github.com/ukama/ukama/systems/common/rest/client/initclient"
	pb "github.com/ukama/ukama/systems/init/bootstrap/pb/gen"
	"github.com/ukama/ukama/systems/init/bootstrap/pkg"
)

const MessagingSystem = "messaging"

type BootstrapServer struct {
	pb.UnimplementedBootstrapServiceServer
	bootstrapRoutingKey msgbus.RoutingKeyBuilder
	msgbus              mb.MsgBusServiceClient
	debug               bool
	orgName             string
	lookupClient        lookupclient.InitClient
	factoryClient       factory.NodeFactoryClient
}

func NewBootstrapServer(orgName string, msgBus mb.MsgBusServiceClient, debug bool, lookupClient lookupclient.InitClient, factoryClient factory.NodeFactoryClient) *BootstrapServer {
	return &BootstrapServer{
		orgName:             orgName,
		bootstrapRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
		msgbus:              msgBus,
		debug:               debug,
		lookupClient:        lookupClient,
		factoryClient:       factoryClient,
	}
}

func (s *BootstrapServer) GetNodeCredentials(ctx context.Context, req *pb.GetNodeCredentialsRequest) (*pb.GetNodeCredentialsResponse, error) {
	node, err := s.factoryClient.Get(req.Id)
	if err != nil {
		log.Errorf("Failed to get node from factory: %v", err)
		return nil, err
	}

	systemInfo, err := s.lookupClient.GetSystem(node.Node.OrgName, "messaging")
	if err != nil {
		log.Errorf("Failed to get messaging system: %v", err)
		return nil, err
	}

	return &pb.GetNodeCredentialsResponse{
		Id:          node.Node.Id,
		Ip:          systemInfo.Ip,
		OrgName:     node.Node.OrgName,
		Certificate: systemInfo.Certificate,
	}, nil
}
