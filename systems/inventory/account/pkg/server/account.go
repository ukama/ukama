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
	"encoding/json"

	"github.com/ukama/ukama/systems/common/gitClient"
	"github.com/ukama/ukama/systems/common/grpc"
	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/inventory/account/pkg"
	"github.com/ukama/ukama/systems/inventory/account/pkg/db"
	"github.com/ukama/ukama/systems/inventory/account/pkg/utils"
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
	gitClient      gitClient.GitClient
	gitDirPath     string
}

func NewAccountServer(orgName string, accountRepo db.AccountRepo, msgBus mb.MsgBusServiceClient, pushGateway string, gc gitClient.GitClient, path string) *AccountServer {
	return &AccountServer{
		orgName:        orgName,
		accountRepo:    accountRepo,
		msgbus:         msgBus,
		baseRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
		pushGateway:    pushGateway,
		gitClient:      gc,
		gitDirPath:     path,
	}
}

func (a *AccountServer) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
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

	return &pb.GetResponse{
		Account: dbAccountToPbAccount(account),
	}, nil
}

func (a *AccountServer) GetByCompany(ctx context.Context, req *pb.GetByCompanmyRequest) (*pb.GetByCompanmyResponse, error) {
	log.Infof("Getting accounts %v", req)

	accounts, err := a.accountRepo.GetByCompany(req.GetCompany())
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "component")
	}

	return &pb.GetByCompanmyResponse{
		Accounts: dbAccountsToPbAccounts(accounts),
	}, nil
}

func (a *AccountServer) SyncAccounts(ctx context.Context, req *pb.SyncAcountsRequest) (*pb.SyncAcountsResponse, error) {
	log.Infof("Syncing accounts %v", req)

	a.gitClient.SetupDir()
	err := a.gitClient.CloneGitRepo()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to clone git repo. Error %s", err.Error())
	}

	rootFileContent, err := a.gitClient.ReadFileJSON(a.gitDirPath + "/root.json")
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to read file. Error %s", err.Error())
	}

	var enviroment gitClient.Environment
	err = json.Unmarshal(rootFileContent, &enviroment)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to unmarshal json. Error %s", err.Error())
	}

	for _, company := range enviroment.Test {
		a.gitClient.BranchCheckout(company.GitBranchName)
		paths, _ := a.gitClient.GetFilesPath("accounts")
		var accounts []utils.Account
		for _, path := range paths {
			content, err := a.gitClient.ReadFileYML(path)
			if err != nil {
				return nil, status.Errorf(codes.Internal, "failed to read file. Error %s", err.Error())
			}
			var account utils.Account
			err = json.Unmarshal(content, &account)
			if err != nil {
				return nil, status.Errorf(codes.Internal, "failed to unmarshal json. Error %s", err.Error())
			}
			account.Company = company.Company
			accounts = append(accounts, account)
		}
		adb := utilAccountsToDbAccounts(accounts)
		err = a.accountRepo.Add(adb)
		if err != nil {
			return nil, grpc.SqlErrorToGrpc(err, "account")
		}
	}

	return &pb.SyncAcountsResponse{}, nil
}

func dbAccountToPbAccount(component *db.Account) *pb.Account {
	return &pb.Account{
		Id:            component.Id.String(),
		Company:       component.Company,
		Item:          component.Item,
		Description:   component.Description,
		Inventory:     component.Inventory,
		OpexFee:       component.OpexFee,
		Vat:           component.Vat,
		EffectiveDate: component.EffectiveDate,
	}
}

func dbAccountsToPbAccounts(accounts []*db.Account) []*pb.Account {
	res := []*pb.Account{}

	for _, i := range accounts {
		res = append(res, dbAccountToPbAccount(i))
	}

	return res
}

func utilAccountsToDbAccounts(accounts []utils.Account) []db.Account {
	res := []db.Account{}

	for _, i := range accounts {
		res = append(res, db.Account{
			Id:            uuid.NewV4(),
			Company:       i.Company,
			Description:   i.Description,
			Item:          i.Item,
			Inventory:     i.Inventory,
			EffectiveDate: i.EffectiveDate,
			OpexFee:       i.OpexFee,
			Vat:           i.Vat,
		})
	}
	return res
}
