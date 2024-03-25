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
	"github.com/ukama/ukama/systems/inventory/component/pkg"
	"github.com/ukama/ukama/systems/inventory/component/pkg/db"
	"github.com/ukama/ukama/systems/inventory/component/pkg/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/yaml.v2"

	log "github.com/sirupsen/logrus"

	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	pb "github.com/ukama/ukama/systems/inventory/component/pb/gen"
)

type ComponentServer struct {
	pb.UnimplementedComponentServiceServer
	orgName        string
	gitClient      gitClient.GitClient
	componentRepo  db.ComponentRepo
	msgbus         mb.MsgBusServiceClient
	baseRoutingKey msgbus.RoutingKeyBuilder
	pushGateway    string
	gitDirPath     string
}

func NewComponentServer(orgName string, componentRepo db.ComponentRepo, msgBus mb.MsgBusServiceClient, pushGateway string, gc gitClient.GitClient, path string) *ComponentServer {
	return &ComponentServer{
		gitClient:      gc,
		orgName:        orgName,
		componentRepo:  componentRepo,
		msgbus:         msgBus,
		baseRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
		pushGateway:    pushGateway,
		gitDirPath:     path,
	}
}

func (c *ComponentServer) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	log.Infof("Getting component %v", req)

	cuuid, err := uuid.FromString(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of component uuid. Error %s", err.Error())
	}
	component, err := c.componentRepo.Get(cuuid)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "component")
	}

	return &pb.GetResponse{
		Component: dbComponentToPbComponent(component),
	}, nil
}

func (c *ComponentServer) GetByCompany(ctx context.Context, req *pb.GetByCompanyRequest) (*pb.GetByCompanyResponse, error) {
	log.Infof("Getting components %v", req)

	components, err := c.componentRepo.GetByCompany(req.GetCompany(), int32(req.GetCategory()))
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "component")
	}

	return &pb.GetByCompanyResponse{
		Components: dbComponentsToPbComponents(components),
	}, nil
}

func (c *ComponentServer) SyncComponents(ctx context.Context, req *pb.SyncComponentsRequest) (*pb.SyncComponentsResponse, error) {
	log.Infof("Syncing components %v", req)

	c.gitClient.SetupDir()
	err := c.gitClient.CloneGitRepo()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to clone git repo. Error %s", err.Error())
	}

	rootFileContent, err := c.gitClient.ReadFileJSON(c.gitDirPath + "/root.json")
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to read file. Error %s", err.Error())
	}

	var enviroment gitClient.Environment
	err = json.Unmarshal(rootFileContent, &enviroment)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to unmarshal json. Error %s", err.Error())
	}

	for _, company := range enviroment.Test {
		c.gitClient.BranchCheckout(company.GitBranchName)
		paths, _ := c.gitClient.GetFilesPath("components")
		var components []utils.Component
		for _, path := range paths {
			content, err := c.gitClient.ReadFileYML(path)
			if err != nil {
				return nil, status.Errorf(codes.Internal, "failed to read file. Error %s", err.Error())
			}
			var component utils.Component
			err = yaml.Unmarshal(content, &component)
			if err != nil {
				return nil, status.Errorf(codes.Internal, "failed to unmarshal json. Error %s", err.Error())
			}
			component.Company = company.Company
			components = append(components, component)
		}
		cdb := utilComponentsToDbComponents(components)
		err = c.componentRepo.Add(cdb)
		if err != nil {
			return nil, grpc.SqlErrorToGrpc(err, "component")
		}
	}

	return &pb.SyncComponentsResponse{}, nil
}

func dbComponentToPbComponent(component *db.Component) *pb.Component {
	return &pb.Component{
		Id:            component.Id.String(),
		Company:       component.Company,
		InventoryId:   component.InventoryId,
		Category:      pb.ComponentCategory(component.Category),
		Type:          component.Type,
		Description:   component.Description,
		DatasheetURL:  component.DatasheetURL,
		ImagesURL:     component.ImagesURL,
		PartNumber:    component.PartNumber,
		Manufacturer:  component.Manufacturer,
		Managed:       component.Managed,
		Warranty:      component.Warranty,
		Specification: component.Specification,
	}
}

func dbComponentsToPbComponents(components []*db.Component) []*pb.Component {
	res := []*pb.Component{}

	for _, i := range components {
		res = append(res, dbComponentToPbComponent(i))
	}

	return res
}

func utilComponentsToDbComponents(components []utils.Component) []db.Component {
	res := []db.Component{}

	for _, i := range components {
		res = append(res, db.Component{
			Id:            uuid.NewV4(),
			Company:       i.Company,
			InventoryId:   i.InventoryID,
			Category:      db.ParseType(i.Category),
			Type:          i.Type,
			Description:   i.Description,
			DatasheetURL:  i.DatasheetURL,
			ImagesURL:     i.ImagesURL,
			PartNumber:    i.PartNumber,
			Manufacturer:  i.Manufacturer,
			Managed:       i.Managed,
			Warranty:      i.Warranty,
			Specification: i.Specification,
		})
	}
	return res
}
