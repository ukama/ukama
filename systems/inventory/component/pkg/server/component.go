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

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/yaml.v2"

	"github.com/ukama/ukama/systems/common/gitClient"
	"github.com/ukama/ukama/systems/common/grpc"
	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/inventory/component/pkg"
	"github.com/ukama/ukama/systems/inventory/component/pkg/db"
	"github.com/ukama/ukama/systems/inventory/component/pkg/utils"

	log "github.com/sirupsen/logrus"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	pb "github.com/ukama/ukama/systems/inventory/component/pb/gen"
)

const uuidParsingError = "Error parsing UUID"

type ComponentServer struct {
	pb.UnimplementedComponentServiceServer
	orgName              string
	gitClient            gitClient.GitClient
	componentRepo        db.ComponentRepo
	msgbus               mb.MsgBusServiceClient
	baseRoutingKey       msgbus.RoutingKeyBuilder
	pushGateway          string
	gitDirPath           string
	componentEnvironment string
	testUserId           string
}

func NewComponentServer(orgName string, componentRepo db.ComponentRepo, msgBus mb.MsgBusServiceClient, pushGateway string, gc gitClient.GitClient, path string, componentEnvironment string, testUserId string) *ComponentServer {
	return &ComponentServer{
		gitClient:            gc,
		gitDirPath:           path,
		msgbus:               msgBus,
		orgName:              orgName,
		testUserId:           testUserId,
		pushGateway:          pushGateway,
		componentRepo:        componentRepo,
		componentEnvironment: componentEnvironment,
		baseRoutingKey:       msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
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

func (c *ComponentServer) GetByUser(ctx context.Context, req *pb.GetByUserRequest) (*pb.GetByUserResponse, error) {
	log.Infof("Getting components by user %v", req)

	components, err := c.componentRepo.GetByUser(req.GetUserId(), int32(ukama.ParseComponentCategory(req.GetCategory())))
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "component")
	}

	return &pb.GetByUserResponse{
		Components: dbComponentsToPbComponents(components),
	}, nil
}

func (c *ComponentServer) List(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error) {
	log.Infof("Listing components %v", req)

	components, err := c.componentRepo.List(req.Id, req.GetUserId(), req.GetPartNumber(), int32(ukama.ParseComponentCategory(req.GetCategory())))
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "component")
	}

	return &pb.ListResponse{
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

	var env gitClient.Environment
	err = json.Unmarshal(rootFileContent, &env)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to unmarshal json. Error %s", err.Error())
	}

	var components []utils.Component
	var environment []gitClient.Company
	if c.componentEnvironment == "test" {
		environment = env.Test
	} else {
		environment = env.Production
	}

	for _, company := range environment {
		err := c.gitClient.BranchCheckout(company.GitBranchName)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to checkout branch. Error %s", err.Error())
		}

		paths, _ := c.gitClient.GetFilesPath("components")
		for _, path := range paths {
			content, err := c.gitClient.ReadFileYML(path)
			if err != nil {
				return nil, status.Errorf(codes.Internal, "failed to read file. Error %s", err.Error())
			}
			var component utils.Component
			err = yaml.Unmarshal(content, &component)
			if err != nil {
				return nil, status.Errorf(codes.Internal, "failed to unmarshal yaml. Error %s", err.Error())
			}

			components = append(components, component)
		}
		var userId string

		if c.componentEnvironment == "test" {
			userId = c.testUserId
		} else {
			userId = company.UserId
		}
		company.UserId = userId
		userIdUUID, err := uuid.FromString(userId)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, uuidParsingError)
		}
		log.Infof("User ID: %s, %s", userIdUUID, company.UserId)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, uuidParsingError)
		}
		cdb := utilComponentsToDbComponents(components, userIdUUID)

		err = c.componentRepo.Delete()
		if err != nil {
			return nil, grpc.SqlErrorToGrpc(err, "component")
		}
		log.Info("Deleted all component records")

		err = c.componentRepo.Add(cdb)
		if err != nil {
			return nil, grpc.SqlErrorToGrpc(err, "component")
		}

		route := c.baseRoutingKey.SetAction("sync").SetObject("components").MustBuild()
		err = c.msgbus.PublishRequest(route, req)
		if err != nil {
			log.Errorf("Failed to publish message %+v with key %+v. Errors %s", req, route, err.Error())
		}
	}

	return &pb.SyncComponentsResponse{}, nil
}

func dbComponentToPbComponent(component *db.Component) *pb.Component {
	return &pb.Component{
		Id:            component.Id.String(),
		Inventory:     component.Inventory,
		UserId:        component.UserId.String(),
		Category:      ukama.ComponentCategory(component.Category).String(),
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

func utilComponentsToDbComponents(components []utils.Component, uId uuid.UUID) []*db.Component {
	res := []*db.Component{}

	for _, i := range components {
		res = append(res, &db.Component{
			Id:            uuid.NewV4(),
			Inventory:     i.InventoryID,
			Category:      ukama.ParseComponentCategory(i.Category),
			UserId:        uId,
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
