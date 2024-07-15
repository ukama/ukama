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
	"google.golang.org/grpc/codes"

	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	ukama "github.com/ukama/ukama/systems/common/validation"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/common/grpc"
	metric "github.com/ukama/ukama/systems/common/metrics"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
	cinvent "github.com/ukama/ukama/systems/common/rest/client/inventory"
	"github.com/ukama/ukama/systems/common/uuid"
	npb "github.com/ukama/ukama/systems/registry/network/pb/gen"
	pb "github.com/ukama/ukama/systems/registry/site/pb/gen"
	"github.com/ukama/ukama/systems/registry/site/pkg"
	"github.com/ukama/ukama/systems/registry/site/pkg/db"
	providers "github.com/ukama/ukama/systems/registry/site/pkg/provider"
)

const uuidParsingError = "Error parsing UUID"

type SiteServer struct {
	pb.UnimplementedSiteServiceServer
	orgName         string
	siteRepo        db.SiteRepo
	msgbus          mb.MsgBusServiceClient
	baseRoutingKey  msgbus.RoutingKeyBuilder
	networkService  providers.NetworkClientProvider
	inventoryClient cinvent.ComponentClient
	pushGateway     string
}

func NewSiteServer(orgName string, siteRepo db.SiteRepo, msgBus mb.MsgBusServiceClient, networkService providers.NetworkClientProvider, pushGateway string, inventoryClientProvider cinvent.ComponentClient) *SiteServer {
	return &SiteServer{
		orgName:         orgName,
		siteRepo:        siteRepo,
		msgbus:          msgBus,
		baseRoutingKey:  msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
		networkService:  networkService,
		pushGateway:     pushGateway,
		inventoryClient: inventoryClientProvider,
	}
}

func (s *SiteServer) Add(ctx context.Context, req *pb.AddRequest) (*pb.AddResponse, error) {
	log.Infof("Adding site %v", req)
	spectrumId, err := uuid.FromString(req.SpectrumId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, uuidParsingError)
	}
	backhaulId, err := uuid.FromString(req.BackhaulId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, uuidParsingError)
	}
	powerId, err := uuid.FromString(req.PowerId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, uuidParsingError)
	}

	accessId, err := uuid.FromString(req.AccessId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, uuidParsingError)
	}

	switchId, err := uuid.FromString(req.SwitchId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, uuidParsingError)
	}

	networkId, err := uuid.FromString(req.NetworkId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, uuidParsingError)
	}

	instDate, err := ukama.ValidateDate(req.GetInstallDate())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	for _, componentIdStr := range []string{
		backhaulId.String(),
		powerId.String(),
		accessId.String(),
		switchId.String(),
		spectrumId.String(),
	} {
		// Validate the parsed UUID using s.inventoryClient
		_, err := s.inventoryClient.Get(componentIdStr)
		if err != nil {
			return nil, err
		}
	}
	svc, err := s.networkService.GetClient()
	if err != nil {
		return nil, err
	}

	log.Infof("checking if network exist %s", req.NetworkId)
	_, err = svc.Get(ctx, &npb.GetRequest{
		NetworkId: networkId.String(),
	})
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "network")
	}

	site := &db.Site{
		NetworkId:     networkId,
		Name:          req.Name,
		Location:      req.Location,
		BackhaulId:    backhaulId,
		PowerId:       powerId,
		AccessId:      accessId,
		SwitchId:      switchId,
		SpectrumId:    spectrumId,
		IsDeactivated: req.IsDeactivated,
		Latitude:      req.Latitude,
		Longitude:     req.Longitude,
		InstallDate:   instDate,
	}

	err = s.siteRepo.Add(site, func(*db.Site, *gorm.DB) error {
		site.Id = uuid.NewV4()
		return nil
	})
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "site")
	}

	evt := &epb.EventAddSite{
		SiteId:        site.Id.String(),
		Name:          site.Name,
		NetworkId:     site.NetworkId.String(),
		IsDeactivated: site.IsDeactivated,
		BackhaulId:    site.BackhaulId.String(),
		PowerId:       site.PowerId.String(),
		AccessId:      site.AccessId.String(),
		SwitchId:      site.SwitchId.String(),
		Latitude:      site.Latitude,
		Longitude:     site.Longitude,
		InstallDate:   site.InstallDate,
	}

	route := s.baseRoutingKey.SetAction("add").SetObject("site").MustBuild()

	err = s.msgbus.PublishRequest(route, req)
	if err != nil {
		log.Errorf("Failed to publish message %+v with key %+v. Errors %s", evt, route, err.Error())
	}

	s.pushSiteCount(networkId)

	return &pb.AddResponse{
		Site: dbSiteToPbSite(site),
	}, nil
}

func (s *SiteServer) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	log.Infof("Getting site %s", req.SiteId)

	siteId, err := uuid.FromString(req.SiteId)

	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, uuidParsingError)
	}

	site, err := s.siteRepo.Get(siteId)

	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "site")
	}
	return &pb.GetResponse{
		Site: dbSiteToPbSite(site)}, nil
}

func (s *SiteServer) GetSites(ctx context.Context, req *pb.GetSitesRequest) (*pb.GetSitesResponse, error) {

	log.Infof("Getting sites %s", req.NetworkId)

	networkId, err := uuid.FromString(req.NetworkId)

	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, uuidParsingError)
	}

	sites, err := s.siteRepo.GetSites(networkId)

	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "site")
	}
	resp := &pb.GetSitesResponse{
		Sites: dbSitesToPbSites(sites),
	}

	return resp, nil
}

func (s *SiteServer) Update(ctx context.Context, req *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	log.Infof("Updating site %s", req.SiteId)

	siteId, err := uuid.FromString(req.SiteId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, uuidParsingError)
	}

	backhaulId, err := uuid.FromString(req.BackhaulId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, uuidParsingError)
	}

	accessId, err := uuid.FromString(req.AccessId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, uuidParsingError)
	}

	powerId, err := uuid.FromString(req.PowerId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, uuidParsingError)
	}

	switchId, err := uuid.FromString(req.SwitchId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, uuidParsingError)
	}

	instDate, err := ukama.ValidateDate(req.GetInstallDate())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	for _, componentIdStr := range []string{
		backhaulId.String(),
		powerId.String(),
		accessId.String(),
		switchId.String(),
	} {
		_, err := s.inventoryClient.Get(componentIdStr)
		if err != nil {
			return nil, err
		}
	}
	site := &db.Site{
		Id:            siteId,
		Name:          req.Name,
		BackhaulId:    backhaulId,
		PowerId:       powerId,
		AccessId:      accessId,
		SwitchId:      switchId,
		IsDeactivated: req.IsDeactivated,
		Latitude:      req.Latitude,
		Longitude:     req.Longitude,
		InstallDate:   instDate,
	}

	err = s.siteRepo.Update(site)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "site")
	}
	evt := &epb.EventUpdateSite{
		SiteId:        site.Id.String(),
		Name:          site.Name,
		IsDeactivated: site.IsDeactivated,
		BackhaulId:    site.BackhaulId.String(),
		PowerId:       site.PowerId.String(),
		AccessId:      site.AccessId.String(),
		SwitchId:      site.SwitchId.String(),
		Latitude:      site.Latitude,
		Longitude:     site.Longitude,
		InstallDate:   site.InstallDate,
	}

	route := s.baseRoutingKey.SetAction("update").SetObject("site").MustBuild()

	err = s.msgbus.PublishRequest(route, req)
	if err != nil {
		log.Errorf("Failed to publish message %+v with key %+v. Errors %s", evt, route, err.Error())
	}
	return &pb.UpdateResponse{
		Site: dbSiteToPbSite(site),
	}, nil
}

func dbSiteToPbSite(site *db.Site) *pb.Site {

	return &pb.Site{
		Id:            site.Id.String(),
		Name:          site.Name,
		Location:      site.Location,
		NetworkId:     site.NetworkId.String(),
		IsDeactivated: site.IsDeactivated,
		BackhaulId:    site.BackhaulId.String(),
		PowerId:       site.PowerId.String(),
		AccessId:      site.AccessId.String(),
		SwitchId:      site.SwitchId.String(),
		SpectrumId:    site.SpectrumId.String(),
		Latitude:      site.Latitude,
		Longitude:     site.Longitude,
		InstallDate:   site.InstallDate,
		CreatedAt:     site.CreatedAt.String(),
	}
}

func dbSitesToPbSites(sites []db.Site) []*pb.Site {
	res := []*pb.Site{}

	for _, s := range sites {
		res = append(res, dbSiteToPbSite(&s))
	}

	return res
}
func (s *SiteServer) pushSiteCount(networkId uuid.UUID) {
	siteCount, err := s.siteRepo.GetSiteCount(networkId)
	if err != nil {
		log.Errorf("failed to get site count: %s", err.Error())
	}

	err = metric.CollectAndPushSimMetrics(s.pushGateway, pkg.SiteMetric, pkg.NumberOfSites, float64(siteCount), map[string]string{"network": networkId.String()}, pkg.SystemName+"-"+pkg.ServiceName)
	if err != nil {
		log.Errorf("Error while pushing site count metric to pushgateway %s", err.Error())
	}
}
