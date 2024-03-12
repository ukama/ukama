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

	"github.com/ukama/ukama/systems/common/grpc"
	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/inventory/contract/pkg"
	"github.com/ukama/ukama/systems/inventory/contract/pkg/db"

	log "github.com/sirupsen/logrus"

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

func (c *ContractServer) GetContracts(ctx context.Context, req *pb.GetContractsRequest) (*pb.GetContractsResponse, error) {

	log.Infof("Getting contracts %v", req)

	contracts, err := c.contractRepo.GetContracts(req.GetCompany(), req.GetActive())
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "component")
	}

	return &pb.GetContractsResponse{
		Contracts: dbContractsToPbContracts(contracts),
		Name:      "test",
		Phone:     "1234567890",
		Email:     "test@ukama.com",
		Address:   "1234 Test St",
		Company:   "Ukama",
	}, nil
}

func dbContractToPbContract(component *db.Contract) *pb.Contract {
	return &pb.Contract{
		Vat:           component.VAT,
		Name:          component.Name,
		Company:       component.Company,
		OpexFee:       component.OpexFee,
		Id:            component.Id.String(),
		Description:   component.Description,
		EffectiveDate: component.EffectiveDate,
		Type:          pb.ContractType(component.Type),
	}
}

func dbContractsToPbContracts(contracts []*db.Contract) []*pb.Contract {
	res := []*pb.Contract{}

	for _, i := range contracts {
		res = append(res, dbContractToPbContract(i))
	}

	return res
}
