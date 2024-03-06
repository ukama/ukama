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
	"github.com/ukama/ukama/systems/inventory/contract/pkg"
	"github.com/ukama/ukama/systems/inventory/contract/pkg/db"

	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	pb "github.com/ukama/ukama/systems/inventory/contract/pb/gen"
)

type ContractServer struct {
	pb.UnimplementedContractServiceServer
	orgName        string
	contractRepo   db.ContractRepo
	msgbus         mb.MsgBusServiceClient
	baseRoutingKey msgbus.RoutingKeyBuilder
	pushGateway    string
}

func NewContractServer(orgName string, contractRepo db.ContractRepo, msgBus mb.MsgBusServiceClient, pushGateway string) *ContractServer {
	return &ContractServer{
		orgName:        orgName,
		contractRepo:   contractRepo,
		msgbus:         msgBus,
		baseRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
		pushGateway:    pushGateway,
	}
}

func (n *ContractServer) GetTest(ctx context.Context, req *pb.GetTestRequest) (*pb.GetTestResponse, error) {
	return &pb.GetTestResponse{
		Service: "Inventory Contract Service",
	}, nil
}
