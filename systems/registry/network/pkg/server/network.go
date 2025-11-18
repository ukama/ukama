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
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/common/validation"
	"github.com/ukama/ukama/systems/registry/network/pkg"
	"github.com/ukama/ukama/systems/registry/network/pkg/db"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"

	log "github.com/sirupsen/logrus"
	metric "github.com/ukama/ukama/systems/common/metrics"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	cnucl "github.com/ukama/ukama/systems/common/rest/client/nucleus"
	pb "github.com/ukama/ukama/systems/registry/network/pb/gen"
)

const uuidParsingError = "Error parsing UUID"

type NetworkServer struct {
	pb.UnimplementedNetworkServiceServer
	orgName        string
	netRepo        db.NetRepo
	orgClient      cnucl.OrgClient
	msgbus         mb.MsgBusServiceClient
	baseRoutingKey msgbus.RoutingKeyBuilder
	pushGateway    string
	country        string
	language       string
	currency       string
	orgId          string
}

func NewNetworkServer(orgName string, netRepo db.NetRepo, orgService cnucl.OrgClient, msgBus mb.MsgBusServiceClient, pushGateway, country, language, currency, orgId string) *NetworkServer {
	return &NetworkServer{
		orgName:        orgName,
		netRepo:        netRepo,
		orgClient:      orgService,
		msgbus:         msgBus,
		baseRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
		pushGateway:    pushGateway,
		country:        country,
		language:       language,
		currency:       currency,
		orgId:          orgId,
	}
}

func (n *NetworkServer) Add(ctx context.Context, req *pb.AddRequest) (*pb.AddResponse, error) {
	// Get the Org locally
	orgName := n.orgName
	networkName := req.GetName()

	log.Infof("lookup for org %s remotely", orgName)
	remoteOrg, err := n.orgClient.Get(orgName)
	if err != nil {
		return nil, err
	}

	if !validation.IsValidDnsLabelName(networkName) {
		return nil, status.Errorf(codes.InvalidArgument, "invalid network name: must be less than 253 "+
			"characters and consist of lowercase characters with hyphens")
	}

	// What should we do if the remote org exists but is deactivated?
	// For now we simply abort.
	if remoteOrg.IsDeactivated {
		return nil, status.Errorf(codes.FailedPrecondition,
			"org is deactivated: cannot add network to it")
	}

	network := &db.Network{
		Name:             networkName,
		AllowedCountries: req.AllowedCountries,
		AllowedNetworks:  req.AllowedNetworks,
		Budget:           req.Budget,
		Overdraft:        req.Overdraft,
		TrafficPolicy:    req.TrafficPolicy,
		PaymentLinks:     req.PaymentLinks,
		SyncStatus:       ukama.StatusTypePending,
	}

	log.Infof("Adding network %s", networkName)
	err = n.netRepo.Add(network, func(*db.Network, *gorm.DB) error {
		network.Id = uuid.NewV4()

		return nil
	})
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "network")
	}

	if n.msgbus != nil {
		route := n.baseRoutingKey.SetAction("add").SetObject("network").MustBuild()
		evt := &epb.EventNetworkCreate{
			Id:               network.Id.String(),
			Name:             network.Name,
			OrgId:            n.orgId,
			AllowedCountries: network.AllowedCountries,
			AllowedNetworks:  network.AllowedNetworks,
			Budget:           network.Budget,
			Overdraft:        network.Overdraft,
			TrafficPolicy:    network.TrafficPolicy,
			PaymentLinks:     network.PaymentLinks,
			IsDeactivated:    network.Deactivated,
		}

		err = n.msgbus.PublishRequest(route, evt)
		if err != nil {
			log.Errorf("Failed to publish message %+v with key %+v. Errors %s",
				evt, route, err.Error())
		}
	}

	n.pushNetworkCount()

	return &pb.AddResponse{
		Network: dbNtwkToPbNtwk(network),
	}, nil
}

func (n *NetworkServer) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	netId, err := uuid.FromString(req.NetworkId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, uuidParsingError)
	}

	nt, err := n.netRepo.Get(netId)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "network")
	}

	return &pb.GetResponse{
		Network: dbNtwkToPbNtwk(nt),
	}, nil
}

func (n *NetworkServer) SetDefault(ctx context.Context, req *pb.SetDefaultRequest) (*pb.SetDefaultResponse, error) {
	netId, err := uuid.FromString(req.NetworkId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, uuidParsingError)
	}

	_, err = n.netRepo.SetDefault(netId, true)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "network")
	}

	return &pb.SetDefaultResponse{}, nil
}

func (n *NetworkServer) GetByName(ctx context.Context, req *pb.GetByNameRequest) (*pb.GetByNameResponse, error) {
	nt, err := n.netRepo.GetByName(req.GetName())
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "network")
	}

	return &pb.GetByNameResponse{
		Network: dbNtwkToPbNtwk(nt),
	}, nil
}

func (n *NetworkServer) GetDefault(ctx context.Context, req *pb.GetDefaultRequest) (*pb.GetDefaultResponse, error) {
	nt, err := n.netRepo.GetDefault()
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "network")
	}

	return &pb.GetDefaultResponse{
		Network: dbNtwkToPbNtwk(nt),
	}, nil
}

func (n *NetworkServer) GetAll(ctx context.Context, req *pb.GetNetworksRequest) (*pb.GetNetworksResponse, error) {

	ntwks, err := n.netRepo.GetAll()
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "networks")
	}

	resp := &pb.GetNetworksResponse{
		Networks: dbNtwksToPbNtwks(ntwks),
	}

	return resp, nil
}

func (n *NetworkServer) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	log.Infof("Deleting network %s", req.NetworkId)
	netId, err := uuid.FromString(req.NetworkId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, uuidParsingError)
	}

	err = n.netRepo.Delete(netId)
	if err != nil {
		log.Error(err)

		return nil, grpc.SqlErrorToGrpc(err, "network")
	}

	if n.msgbus != nil {
		route := n.baseRoutingKey.SetAction("delete").SetObject("network").MustBuild()
		evt := &epb.EventNetworkDelete{
			Id:    req.NetworkId,
			OrgId: n.orgId,
		}

		err = n.msgbus.PublishRequest(route, evt)
		if err != nil {
			log.Errorf("Failed to publish message %+v with key %+v. Errors %s",
				evt, route, err.Error())
		}
	}

	n.pushNetworkCount()

	return &pb.DeleteResponse{}, nil
}

func dbNtwkToPbNtwk(ntwk *db.Network) *pb.Network {
	return &pb.Network{
		Id:               ntwk.Id.String(),
		Name:             ntwk.Name,
		AllowedCountries: ntwk.AllowedCountries,
		AllowedNetworks:  ntwk.AllowedNetworks,
		Budget:           ntwk.Budget,
		Overdraft:        ntwk.Overdraft,
		TrafficPolicy:    ntwk.TrafficPolicy,
		PaymentLinks:     ntwk.PaymentLinks,
		IsDeactivated:    ntwk.Deactivated,
		SyncStatus:       ntwk.SyncStatus.String(),
		IsDefault:        ntwk.IsDefault,
		CreatedAt:        timestamppb.New(ntwk.CreatedAt),
	}
}

func dbNtwksToPbNtwks(ntwks []db.Network) []*pb.Network {
	res := []*pb.Network{}

	for _, n := range ntwks {
		res = append(res, dbNtwkToPbNtwk(&n))
	}

	return res
}

func (n *NetworkServer) pushNetworkCount() {
	networkCount, err := n.netRepo.GetNetworkCount()
	if err != nil {
		log.Errorf("failed to get network counts: %s", err.Error())
	}

	err = metric.CollectAndPushSimMetrics(n.pushGateway, pkg.NetworkMetric, pkg.NumberOfNetworks,
		float64(networkCount), map[string]string{"org": n.orgId}, pkg.SystemName+"-"+pkg.ServiceName)
	if err != nil {
		log.Errorf("Error while pushing network count metric to pushgateway %s", err.Error())
	}
}

func (n *NetworkServer) PushMetrics() error {
	n.pushNetworkCount()
	return nil
}
