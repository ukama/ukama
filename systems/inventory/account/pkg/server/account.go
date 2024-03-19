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
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/inventory/account/pkg"
	"github.com/ukama/ukama/systems/inventory/account/pkg/db"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	log "github.com/sirupsen/logrus"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	pb "github.com/ukama/ukama/systems/inventory/account/pb/gen"
)

type AccountServer struct {
	pb.UnimplementedAccountServiceServer
	orgName        string
	accountRepo    db.AccountRepo
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

func (a *AccountServer) GetAccount(ctx context.Context, req *pb.GetAcountRequest) (*pb.GetAcountResponse, error) {
	log.Infof("Getting account %v", req)

	auuid, err := uuid.FromString(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of account uuid. Error %s", err.Error())
	}
	account, err := a.accountRepo.Get(auuid)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "account")
	}

	return &pb.GetAcountResponse{
		Account: dbAccountToPbAccount(account),
	}, nil
}

func (a *AccountServer) GetAccounts(ctx context.Context, req *pb.GetAcountsRequest) (*pb.GetAcountsResponse, error) {
	log.Infof("Getting accounts %v", req)

	accounts, err := a.accountRepo.GetByCompany(req.GetCompany(), req.GetCategory().Enum().String())
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "component")
	}

	return &pb.GetAcountsResponse{
		Accounts: dbAccountsToPbAccounts(accounts),
	}, nil
}

func (n *AccountServer) SyncAccounts(ctx context.Context, req *pb.SyncAcountsRequest) (*pb.SyncAcountsResponse, error) {
	return &pb.SyncAcountsResponse{}, nil
}

func dbAccountToPbAccount(component *db.Account) *pb.Account {
	return &pb.Account{
		Id:            component.Id.String(),
		Company:       component.Company,
		Category:      pb.Category(pb.Category_value[component.Category]),
		Item:          component.Item,
		Quantity:      component.Quantity,
		PricePerUnit:  component.PricePerUnit,
		TotalPrice:    component.TotalPrice,
		Description:   component.Description,
		Specification: component.Specification,
		PaymentType:   pb.PaymentType(component.PaymentType),
	}
}

func dbAccountsToPbAccounts(accounts []*db.Account) []*pb.Account {
	res := []*pb.Account{}

	for _, i := range accounts {
		res = append(res, dbAccountToPbAccount(i))
	}

	return res
}
