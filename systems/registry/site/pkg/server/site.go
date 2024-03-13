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

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/registry/site/pkg"
	"github.com/ukama/ukama/systems/registry/site/pkg/db"

	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	pb "github.com/ukama/ukama/systems/registry/site/pb/gen"
)


type SiteServer struct {
	pb.UnimplementedSiteServiceServer
	orgName        string
	siteRepo       db.SiteRepo
	msgbus         mb.MsgBusServiceClient
	baseRoutingKey msgbus.RoutingKeyBuilder
	pushGateway    string
}

func NewsiteServer(orgName string,siteRepo db.SiteRepo, msgBus mb.MsgBusServiceClient, pushGateway string) *SiteServer {
	return &SiteServer{
		orgName:        orgName,
		siteRepo:       siteRepo,
		msgbus:         msgBus,
		baseRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
		pushGateway:    pushGateway,
	}
}

func (n *SiteServer) Add(ctx context.Context, req *pb.AddRequest) (*pb.AddResponse, error) {	

	
	return nil, nil
	
}

func (n *SiteServer) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	
	return nil, nil

}



func (n *SiteServer) AddSite(ctx context.Context, req *pb.AddRequest) (*pb.AddResponse, error) {
	
	return nil, nil

}

func (n *SiteServer) GetSite(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	

	return nil, nil

}


func (n *SiteServer) GetSites(ctx context.Context, req *pb.GetSitesRequest) (*pb.GetSitesResponse, error) {
	
	return nil, nil

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
