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
	"github.com/ukama/ukama/systems/inventory/accounting/pkg"
	"github.com/ukama/ukama/systems/inventory/accounting/pkg/db"
	"github.com/ukama/ukama/systems/inventory/accounting/pkg/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	log "github.com/sirupsen/logrus"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	pb "github.com/ukama/ukama/systems/inventory/accounting/pb/gen"
)

type AccountingServer struct {
	pb.UnimplementedAccountingServiceServer
	orgName        string
	accountingRepo db.AccountingRepo
	msgbus         mb.MsgBusServiceClient
	baseRoutingKey msgbus.RoutingKeyBuilder
	pushGateway    string
	gitClient      gitClient.GitClient
	gitDirPath     string
}

func NewAccountingServer(orgName string, accountingRepo db.AccountingRepo, msgBus mb.MsgBusServiceClient, pushGateway string, gc gitClient.GitClient, path string) *AccountingServer {
	return &AccountingServer{
		orgName:        orgName,
		accountingRepo: accountingRepo,
		msgbus:         msgBus,
		baseRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
		pushGateway:    pushGateway,
		gitClient:      gc,
		gitDirPath:     path,
	}
}

func (a *AccountingServer) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	log.Infof("Getting accounting %v", req)

	auuid, err := uuid.FromString(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of accounting uuid. Error %s", err.Error())
	}
	accounting, err := a.accountingRepo.Get(auuid)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "accounting")
	}
	return &pb.GetResponse{
		Accounting: dbAccountingToPbAccounting(accounting),
	}, nil
}

func (a *AccountingServer) GetByCompany(ctx context.Context, req *pb.GetByCompanmyRequest) (*pb.GetByCompanmyResponse, error) {
	log.Infof("Getting accountings %v", req)

	accountings, err := a.accountingRepo.GetByCompany(req.GetCompany())
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "component")
	}

	return &pb.GetByCompanmyResponse{
		Accounting: dbAccountingsToPbAccountings(accountings),
	}, nil
}

func (a *AccountingServer) GetByUser(ctx context.Context, req *pb.GetByUserRequest) (*pb.GetByUserResponse, error) {
	log.Infof("Getting accountings by user %v", req)

	accountings, err := a.accountingRepo.GetByUser(req.GetUserId())
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "component")
	}

	return &pb.GetByUserResponse{
		Accounting: dbAccountingsToPbAccountings(accountings),
	}, nil
}

func (a *AccountingServer) SyncAccounting(ctx context.Context, req *pb.SyncAcountingRequest) (*pb.SyncAcountingResponse, error) {
	log.Infof("Syncing accountings %v", req)

	a.gitClient.SetupDir()
	err := a.gitClient.CloneGitRepo()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to clone git repo. Error %s", err.Error())
	}

	rootFileContent, err := a.gitClient.ReadFileJSON(a.gitDirPath + "/root.json")
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to read root file. Error %s", err.Error())
	}

	var enviroment gitClient.Environment
	err = json.Unmarshal(rootFileContent, &enviroment)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to unmarshal root json file. Error %s", err.Error())
	}

	for _, company := range enviroment.Test {
		err := a.gitClient.BranchCheckout(company.GitBranchName)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to checkout branch. Error %s", err.Error())
		}
		manifestFileContent, err := a.gitClient.ReadFileJSON(a.gitDirPath + "/manifest.json")
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to read manifest file. Error %s", err.Error())
		}

		var accounting utils.Accounting
		err = json.Unmarshal(manifestFileContent, &accounting)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to unmarshal manifest json file. Error %s", err.Error())
		}
		adb := utilAccountsToDbAccounts(accounting, company.Company, company.UserId)
		inventoryIds := utils.UniqueInventoryIds(adb)
		if len(inventoryIds) > 0 {
			err = a.accountingRepo.Delete(inventoryIds)
			if err != nil {
				return nil, grpc.SqlErrorToGrpc(err, "accounting")
			}
			log.Info("Deleted accountings with inventory ids: ", inventoryIds)
		}

		err = a.accountingRepo.Add(adb)
		if err != nil {
			return nil, grpc.SqlErrorToGrpc(err, "accounting")
		}
		log.Info("Added accountings: ", adb)

		route := a.baseRoutingKey.SetAction("sync").SetObject("accounting").MustBuild()
		err = a.msgbus.PublishRequest(route, req)
		if err != nil {
			log.Errorf("Failed to publish message %+v with key %+v. Errors %s", req, route, err.Error())
		}
	}

	return &pb.SyncAcountingResponse{}, nil
}

func dbAccountingToPbAccounting(component *db.Accounting) *pb.Accounting {
	return &pb.Accounting{
		Id:            component.Id.String(),
		Company:       component.Company,
		Item:          component.Item,
		UserId:        component.UserId,
		Description:   component.Description,
		Inventory:     component.Inventory,
		OpexFee:       component.OpexFee,
		Vat:           component.Vat,
		EffectiveDate: component.EffectiveDate,
	}
}

func dbAccountingsToPbAccountings(accountings []*db.Accounting) []*pb.Accounting {
	res := []*pb.Accounting{}

	for _, i := range accountings {
		res = append(res, dbAccountingToPbAccounting(i))
	}

	return res
}

func utilAccountsToDbAccounts(accounting utils.Accounting, company string, userId string) []*db.Accounting {
	res := []*db.Accounting{}

	for _, i := range accounting.Ukama {
		res = append(res, &db.Accounting{
			Id:            uuid.NewV4(),
			Company:       company,
			UserId:        userId,
			Description:   i.Description,
			Item:          i.Item,
			Inventory:     i.Inventory,
			EffectiveDate: i.EffectiveDate,
			OpexFee:       i.OpexFee,
			Vat:           i.Vat,
		})
	}
	for _, i := range accounting.Backhaul {
		res = append(res, &db.Accounting{
			Id:            uuid.NewV4(),
			Company:       company,
			UserId:        userId,
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
