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
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/ukama/ukama/systems/common/grpc"
	metric "github.com/ukama/ukama/systems/common/metrics"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
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
	orgName        string
	siteRepo       db.SiteRepo
	msgbus         mb.MsgBusServiceClient
	baseRoutingKey msgbus.RoutingKeyBuilder
	networkService providers.NetworkClientProvider
	pushGateway    string
}

func NewsiteServer(orgName string, siteRepo db.SiteRepo, msgBus mb.MsgBusServiceClient, networkService providers.NetworkClientProvider, pushGateway string) *SiteServer {
	return &SiteServer{
		orgName:        orgName,
		siteRepo:       siteRepo,
		msgbus:         msgBus,
		baseRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
		networkService: networkService,
		pushGateway:    pushGateway,
	}
}

func (s *SiteServer) Add(ctx context.Context, req *pb.AddRequest) (*pb.AddResponse, error) {

	log.Infof("Adding site %s", req.Name)

	backhaulID, err := uuid.FromString(req.BackhaulId)

	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, uuidParsingError)
	}
	powerID, err := uuid.FromString(req.PowerId)

	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, uuidParsingError)
	}

	accessID, err := uuid.FromString(req.AccessId)

	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, uuidParsingError)
	}
	switchID, err := uuid.FromString(req.SwitchId)

	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, uuidParsingError)
	}
	netID, err := uuid.FromString(req.NetworkId)

	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, uuidParsingError)
	}

	svc, err := s.networkService.GetClient()
	if err != nil {
		return nil, err
	}
	log.Infof("checking if network exist %s", req.NetworkId)

	_, err = svc.Get(ctx, &npb.GetRequest{
		NetworkId: netID.String(),
	})
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "network")

	}

	site := &db.Site{
		NetworkID:     netID,
		Name:          req.Name,
		BackhaulID:    backhaulID,
		PowerID:       powerID,
		AccessID:      accessID,
		SwitchID:      switchID,
		IsDeactivated: req.IsDeactivated,
		Latitude:      req.Latitude,
		Longitude:     req.Longitude,
		InstallDate:   req.InstallDate.AsTime(),
	}

	err = s.siteRepo.Add(site, func(*db.Site, *gorm.DB) error {
		site.ID = uuid.NewV4()
		return nil
	})

	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "site")
	}
	route := s.baseRoutingKey.SetAction("add").SetObject("site").MustBuild()

	err = s.msgbus.PublishRequest(route, req)
	if err != nil {
		log.Errorf("Failed to publish message %+v with key %+v. Errors %s", req, route, err.Error())
	}
	s.pushSiteCount(netID)

	return &pb.AddResponse{
		Site: dbSiteToPbSite(site)}, nil
}

func (s *SiteServer) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	log.Infof("Getting site %s-%s", req.NetworkId, req.SiteId)

	netID, err := uuid.FromString(req.NetworkId)

	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, uuidParsingError)
	}
	siteID, err := uuid.FromString(req.SiteId)

	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, uuidParsingError)
	}

	svc, err := s.networkService.GetClient()
	if err != nil {
		return nil, err
	}
	log.Infof("checking if network exist %s", req.NetworkId)

	_, err = svc.Get(ctx, &npb.GetRequest{
		NetworkId: netID.String(),
	})
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "network")

	}
	site, err := s.siteRepo.Get(netID, siteID)

	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "site")
	}
	return &pb.GetResponse{
		Site: dbSiteToPbSite(site)}, nil
}

func (s *SiteServer) GetSites(ctx context.Context, req *pb.GetSitesRequest) (*pb.GetSitesResponse, error) {

	log.Infof("Getting sites %s", req.NetworkId)

	netID, err := uuid.FromString(req.NetworkId)

	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, uuidParsingError)
	}

	svc, err := s.networkService.GetClient()
	if err != nil {
		return nil, err
	}
	log.Infof("checking if network exist %s", req.NetworkId)
	
	_, err = svc.Get(ctx, &npb.GetRequest{
		NetworkId: netID.String(),
	})
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "network")

	}

	sites, err := s.siteRepo.GetSites(netID)

	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "site")
	}
	resp := &pb.GetSitesResponse{
		Sites: dbSitesToPbSites(sites),
	}

	return resp, nil
}

func (s *SiteServer) Update(ctx context.Context, req *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	log.Infof("Updating site %s-%s", req.NetworkId, req.SiteId)

	netID, err := uuid.FromString(req.NetworkId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, uuidParsingError)
	}

	siteID, err := uuid.FromString(req.SiteId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, uuidParsingError)
	}

	site, err := s.siteRepo.Get(netID, siteID)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "site")
	}

	// Update the site fields
	if req.Name != "" {
		site.Name = req.Name
	}
	if req.BackhaulId != "" {
		backhaulID, err := uuid.FromString(req.BackhaulId)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, uuidParsingError)
		}
		site.BackhaulID = backhaulID
	}
	// Similarly, update other fields as required
	site, err = s.siteRepo.Get(netID, siteID)

	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "site")
	}

	// Save the updated site to the database
	err = s.siteRepo.Update(site)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "site")
	}

	return &pb.UpdateResponse{
		Site: dbSiteToPbSite(site),
	}, nil
}

func dbSiteToPbSite(site *db.Site) *pb.Site {
	return &pb.Site{
		Id:            site.ID.String(),
		Name:          site.Name,
		NetworkId:     site.NetworkID.String(),
		BackhaulId:    site.BackhaulID.String(),
		PowerId:       site.PowerID.String(),
		AccessId:      site.AccessID.String(),
		SwitchId:      site.SwitchID.String(),
		IsDeactivated: site.IsDeactivated,
		Latitude:      site.Latitude,
		Longitude:     site.Longitude,
		InstallDate:   &timestamp.Timestamp{Seconds: site.InstallDate.Unix()}, // Convert time.Time to Timestamp
	}
}

func dbSitesToPbSites(sites []db.Site) []*pb.Site {
	res := []*pb.Site{}

	for _, s := range sites {
		res = append(res, dbSiteToPbSite(&s))
	}

	return res
}
func (s *SiteServer) pushSiteCount(netId uuid.UUID) {
	siteCount, err := s.siteRepo.GetSiteCount(netId)
	if err != nil {
		log.Errorf("failed to get site count: %s", err.Error())
	}

	err = metric.CollectAndPushSimMetrics(s.pushGateway, pkg.SiteMetric, pkg.NumberOfSites, float64(siteCount), map[string]string{"network": netId.String()}, pkg.SystemName+"-"+pkg.ServiceName)
	if err != nil {
		log.Errorf("Error while pushing site count metric to pushgateway %s", err.Error())
	}
}

