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

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/grpc"
	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/inventory/component/pkg"
	"github.com/ukama/ukama/systems/inventory/component/pkg/db"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	pb "github.com/ukama/ukama/systems/inventory/component/pb/gen"
	pkgP "github.com/ukama/ukama/systems/inventory/component/pkg/providers"
	gpb "github.com/ukama/ukama/systems/services/gitClient/pb/gen"
)

type ComponentServer struct {
	pb.UnimplementedComponentServiceServer
	orgName        string
	gitService     pkgP.GitClientProvider
	componentRepo  db.ComponentRepo
	msgbus         mb.MsgBusServiceClient
	baseRoutingKey msgbus.RoutingKeyBuilder
	pushGateway    string
}

func NewComponentServer(orgName string, componentRepo db.ComponentRepo,
	msgBus mb.MsgBusServiceClient, pushGateway string, gitService pkgP.GitClientProvider) *ComponentServer {
	return &ComponentServer{
		msgbus:         msgBus,
		orgName:        orgName,
		gitService:     gitService,
		pushGateway:    pushGateway,
		componentRepo:  componentRepo,
		baseRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
	}
}

func (s *ComponentServer) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	log.Infof("Getting component %v", req)

	cuuid, err := uuid.FromString(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of component uuid. Error %s", err.Error())
	}
	component, err := s.componentRepo.Get(cuuid)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "component")
	}

	return &pb.GetResponse{
		Component: dbComponentToPbComponent(component),
	}, nil
}

func (s *ComponentServer) GetByCompany(ctx context.Context, req *pb.GetByCompanyRequest) (*pb.GetByCompanyResponse, error) {
	log.Infof("Getting components %v", req)

	components, err := s.componentRepo.GetByCompany(req.GetCompany(), req.GetType().Enum().String())
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "component")
	}

	return &pb.GetByCompanyResponse{
		Components: dbComponentsToPbComponents(components),
	}, nil
}

func (s *ComponentServer) Add(ctx context.Context, req *pb.AddRequest) (*pb.AddResponse, error) {
	log.Infof("Adding component %v", req)

	cuuid := uuid.NewV4()

	component := &db.Component{
		Id:            cuuid,
		Company:       req.GetCompany(),
		InventoryId:   req.GetInventoryId(),
		Category:      req.GetCategory(),
		Type:          db.ComponentType(req.GetType()),
		Description:   req.GetDescription(),
		DatasheetURL:  req.GetDatasheetURL(),
		ImagesURL:     req.GetImagesURL(),
		PartNumber:    req.GetPartNumber(),
		Manufacturer:  req.GetManufacturer(),
		Managed:       req.GetManaged(),
		Warranty:      req.GetWarranty(),
		Specification: req.GetSpecification(),
	}

	err := s.componentRepo.Add(component, nil)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "component")
	}

	return &pb.AddResponse{
		Id: cuuid.String(),
	}, nil
}

func (s *ComponentServer) SyncComponents(ctx context.Context, req *pb.SyncComponentsRequest) (*pb.SyncComponentsResponse, error) {
	log.Infof("Syncing components %v", req)

	gc, err := s.gitService.GetClient()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get git client. Error: %v", err)
	}

	fcr, err := gc.FetchComponents(ctx, &gpb.FetchComponentsRequest{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch components from git client. Error: %v", err)
	}

	for _, c := range fcr.Component {
		cuuid := uuid.NewV4()

		component := &db.Component{
			Id:            cuuid,
			Company:       c.GetCompany(),
			InventoryId:   c.GetInventoryId(),
			Category:      c.GetCategory(),
			Type:          db.ComponentType(c.GetType()),
			Description:   c.GetDescription(),
			DatasheetURL:  c.GetDatasheetURL(),
			ImagesURL:     c.GetImagesURL(),
			PartNumber:    c.GetPartNumber(),
			Manufacturer:  c.GetManufacturer(),
			Managed:       c.GetManaged(),
			Warranty:      c.GetWarranty(),
			Specification: c.GetSpecification(),
		}

		err = s.componentRepo.Add(component, nil)
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
		Category:      component.Category,
		Type:          pb.ComponentType(component.Type),
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
