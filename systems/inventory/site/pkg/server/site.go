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

	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/inventory/site/pkg"
	"github.com/ukama/ukama/systems/inventory/site/pkg/db"

	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	cnucl "github.com/ukama/ukama/systems/common/rest/client/nucleus"
	pb "github.com/ukama/ukama/systems/inventory/site/pb/gen"
)

const uuidParsingError = "Error parsing UUID"

type SiteServer struct {
	pb.UnimplementedSiteServiceServer
	orgName        string
	siteRepo       db.SiteRepo
	orgClient      cnucl.OrgClient
	msgbus         mb.MsgBusServiceClient
	baseRoutingKey msgbus.RoutingKeyBuilder
	pushGateway    string
}

func NewSiteServer(orgName string, siteRepo db.SiteRepo,
	orgService cnucl.OrgClient, msgBus mb.MsgBusServiceClient, pushGateway string) *SiteServer {
	return &SiteServer{
		orgName:        orgName,
		siteRepo:       siteRepo,
		orgClient:      orgService,
		msgbus:         msgBus,
		baseRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
		pushGateway:    pushGateway,
	}
}

func (n *SiteServer) GetTest(ctx context.Context, req *pb.GetTestRequest) (*pb.GetTestResponse, error) {
	return &pb.GetTestResponse{
		Service: "Inventory Site Service",
	}, nil
}
