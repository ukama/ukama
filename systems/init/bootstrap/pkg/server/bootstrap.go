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
	"github.com/ukama/ukama/systems/init/bootstrap/client"
	pb "github.com/ukama/ukama/systems/init/bootstrap/pb/gen"
	"github.com/ukama/ukama/systems/init/bootstrap/pkg"
	lpb "github.com/ukama/ukama/systems/init/lookup/pb/gen"
)

const MessagingSystem = "messaging"

type BootstrapServer struct {
	pb.UnimplementedBootstrapServiceServer
	bootstrapRoutingKey msgbus.RoutingKeyBuilder
	msgbus              mb.MsgBusServiceClient
	debug               bool
	orgName             string
	lookupClient        client.LookupClientProvider
	factoryClient       factory.NodeFactoryClient
}

func NewBootstrapServer(orgName string, msgBus mb.MsgBusServiceClient, debug bool, lookupClient client.LookupClientProvider, factoryClient factory.NodeFactoryClient) *BootstrapServer {
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

	lookupSvc, err := s.lookupClient.GetClient()
	if err != nil {
		return nil, err
	}

	var ip string
	var certificate string

	if node.Node.OrgName != "" {
		msgSystem, err := lookupSvc.GetSystemForOrg(ctx, &lpb.GetSystemRequest{
			OrgName:    node.Node.OrgName,
			SystemName: MessagingSystem,
		})
		if err != nil {
			log.Errorf("Failed to get messaging system: %v", err)
		} else {
			ip = msgSystem.Ip
			certificate = msgSystem.Certificate
		}
	}

	return &pb.GetNodeCredentialsResponse{
		Id:          node.Node.Id,
		OrgName:     node.Node.OrgName,
		Ip:          ip,
		Certificate: certificate,
	}, nil
}
