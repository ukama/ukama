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
	"fmt"

	log "github.com/sirupsen/logrus"
	cpb "github.com/ukama/ukama/systems/common/pb/gen/ukama"

	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/node/controller/pkg/providers"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/uuid"
	pb "github.com/ukama/ukama/systems/node/controller/pb/gen"

	"github.com/ukama/ukama/systems/node/controller/pkg"
	"github.com/ukama/ukama/systems/node/controller/pkg/db"
)

type ControllerServer struct {
	pb.UnimplementedControllerServiceServer
	nRepo                db.NodeLogRepo
	nodeFeederRoutingKey msgbus.RoutingKeyBuilder
	msgbus               mb.MsgBusServiceClient
	registrySystem       providers.RegistryProvider
	debug                bool
	orgName              string
}

func NewControllerServer(orgName string, nRepo db.NodeLogRepo, msgBus mb.MsgBusServiceClient, registry providers.RegistryProvider, debug bool) *ControllerServer {
	return &ControllerServer{
		nRepo:                nRepo,
		orgName:              orgName,
		nodeFeederRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
		msgbus:               msgBus,
		debug:                debug,
		registrySystem:       registry,
	}
}

func (c *ControllerServer) RestartSite(ctx context.Context, req *pb.RestartSiteRequest) (*pb.RestartSiteResponse, error) {
	log.Infof("Restarting site %v", req)

	if req.SiteName == "" {
		return nil, status.Errorf(codes.InvalidArgument, "site name cannot be empty")
	}

	if req.NetworkId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "network cannot be empty")
	}

	netId, err := uuid.FromString(req.GetNetworkId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid network ID format: %s", err.Error())
	}

	if err := c.registrySystem.ValidateSite(netId.String(), req.GetSiteName(), c.orgName); err != nil {
		return nil, fmt.Errorf("failed to validate site %s and network %s. Error %s", req.GetSiteName(), netId.String(), err.Error())
	}

	if err := c.registrySystem.ValidateNetwork(netId.String(), c.orgName); err != nil {
		return nil, fmt.Errorf("failed to validate network with network %s. Error %s", netId.String(), err.Error())
	}

	nodes, err := c.registrySystem.GetNodesBySite(req.SiteName)
	if err != nil {
		return nil, fmt.Errorf("failed to get nodes with site %s and network %s. Error %s", req.GetSiteName(), netId.String(), err.Error())

	}
	for _, nodeId := range nodes {

		nId, err := ukama.ValidateNodeId(nodeId)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument,
				"invalid format of node id. Error %s", err.Error())
		}

		_, err = c.nRepo.Get(nId.String())
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "Node has not been registered yet: %s", err.Error())
		}

		msg := &pb.RestartNodeRequest{
			NodeId: nId.String(),
		}
		data, err := proto.Marshal(msg)
		if err != nil {
			return nil, err
		}

		err = c.publishMessage(c.orgName+"."+"."+"."+nId.String(), "/v1/node/"+nId.String()+"/reboot", data)

		if err != nil {
			log.Errorf("Failed to publish message. Errors %s", err.Error())
			return nil, status.Errorf(codes.Internal, "Failed to publish message: %s", err.Error())

		}
	}

	return &pb.RestartSiteResponse{}, nil
}

func (c *ControllerServer) RestartNode(ctx context.Context, req *pb.RestartNodeRequest) (*pb.RestartNodeResponse, error) {
	if req.NodeId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "node ID cannot be empty")
	}

	nId, err := ukama.ValidateNodeId(req.NodeId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of node id. Error %s", err.Error())
	}

	_, err = c.nRepo.Get(nId.String())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Node has not been registered yet: %s", err.Error())
	}

	msg := &pb.RestartNodeRequest{
		NodeId: nId.String(),
	}
	data, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}

	err = c.publishMessage(c.orgName+"."+"."+"."+nId.String(), "/v1/node/"+nId.String()+"/reboot", data)

	if err != nil {
		log.Errorf("Failed to publish message. Errors %s", err.Error())
		return nil, status.Errorf(codes.Internal, "Failed to publish message: %s", err.Error())

	}
	return &pb.RestartNodeResponse{}, nil
}

func (c *ControllerServer) RestartNodes(ctx context.Context, req *pb.RestartNodesRequest) (*pb.RestartNodesResponse, error) {
	if len(req.NodeIds) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "node IDs cannot be empty")
	}

	for _, nodeId := range req.NodeIds {
		nId, err := ukama.ValidateNodeId(string(nodeId))
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument,
				"invalid format of node id. Error %s", err.Error())
		}

		_, err = c.nRepo.Get(nId.String())
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "Node has not been registered yet: %s", err.Error())
		}
		msg := &pb.RestartNodeRequest{
			NodeId: string(nodeId),
		}
		data, err := proto.Marshal(msg)
		if err != nil {
			return nil, err
		}

		err = c.publishMessage(c.orgName+"."+"."+"."+nId.String(), "/v1/node/"+nId.String()+"/reboot", data)

		if err != nil {
			log.Errorf("Failed to publish message . Errors %s", err.Error())
			return nil, status.Errorf(codes.Internal, "Failed to publish message: %s", err.Error())

		}
	}

	return &pb.RestartNodesResponse{}, nil

}
func (c *ControllerServer) ToggleInternetSwitch(ctx context.Context, req *pb.ToggleInternetSwitchRequest) (*pb.ToggleInternetSwitchResponse, error) {
	log.Infof("Toggling internet switch for site %v, port %v to %v", req.SiteId, req.Port, req.Status)

	if req.SiteId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "site ID cannot be empty")
	}
	siteId, err := uuid.FromString(req.GetSiteId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid site ID format: %s", err.Error())
	}

	if err := c.registrySystem.ValidateSite(req.SiteId, "", c.orgName); err != nil {
		return nil, fmt.Errorf("failed to validate site %s. Error %s", req.SiteId, err.Error())
	}

	msg := &pb.ToggleInternetSwitchRequest{
		SiteId: siteId.String(),
		Status: req.Status,
		Port:   req.Port,
	}
	data, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}
	err = c.publishMessage(c.orgName+"."+"."+"."+siteId.String(), "/v1/switch/"+fmt.Sprintf("%d/%t", req.Port, req.Status), data)

	if err != nil {
		log.Errorf("Failed to publish switch port reboot message. Errors: %s", err.Error())
		return nil, status.Errorf(codes.Internal, "Failed to publish switch port reboot message: %s", err.Error())
	}
	return &pb.ToggleInternetSwitchResponse{}, nil
}

func (c *ControllerServer) publishMessage(target string, path string, anyMsg []byte) error {
	route := "request.cloud.local" + "." + c.orgName + "." + pkg.SystemName + "." + pkg.ServiceName + "." + "nodefeeder" + "." + "publish"
	msg := &cpb.NodeFeederMessage{
	  Target: target,
	  HTTPMethod: "POST",
	  Path: path,
	  Msg: anyMsg,
	}
	log.Infof("Published controller %s on route %s on target %s ", anyMsg, route, target)
	err := c.msgbus.PublishRequest(route, msg)
	return err
  }
  