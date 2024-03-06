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
	"github.com/ukama/ukama/systems/inventory/account/pkg"
	"github.com/ukama/ukama/systems/inventory/account/pkg/db"

	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	cnucl "github.com/ukama/ukama/systems/common/rest/client/nucleus"
	pb "github.com/ukama/ukama/systems/inventory/account/pb/gen"
)

const uuidParsingError = "Error parsing UUID"

type AccountServer struct {
	pb.UnimplementedAccountServiceServer
	orgName        string
	accountRepo    db.AccountRepo
	orgClient      cnucl.OrgClient
	msgbus         mb.MsgBusServiceClient
	baseRoutingKey msgbus.RoutingKeyBuilder
	pushGateway    string
}

func NewAccountServer(orgName string, accountRepo db.AccountRepo, msgBus mb.MsgBusServiceClient, pushGateway string) *AccountServer {
	return &AccountServer{
		orgName:        orgName,
		accountRepo:    accountRepo,
		msgbus:         msgBus,
		baseRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
		pushGateway:    pushGateway,
	}
}

func (n *AccountServer) GetTest(ctx context.Context, req *pb.GetTestRequest) (*pb.GetTestResponse, error) {
	return &pb.GetTestResponse{
		Service: "Inventory Account Service",
	}, nil
}
