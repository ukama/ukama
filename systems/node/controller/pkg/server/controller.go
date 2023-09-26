package server

import (
	"context"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/uuid"
	pb "github.com/ukama/ukama/systems/node/controller/pb/gen"

	"github.com/ukama/ukama/systems/node/controller/pkg"
	"github.com/ukama/ukama/systems/node/controller/pkg/providers"
)

type ControllerServer struct {
	pb.UnimplementedControllerServiceServer
	msgbus               mb.MsgBusServiceClient
	registrySystem       providers.RegistryProvider
	controllerRoutingKey msgbus.RoutingKeyBuilder
	debug                bool
	orgName              string
}

func NewControllerServer(msgBus mb.MsgBusServiceClient, registry providers.RegistryProvider, debug bool, orgName string) *ControllerServer {
	return &ControllerServer{
		controllerRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
		msgbus:               msgBus,
		registrySystem:       registry,
		debug:                pkg.IsDebugMode,
		orgName:              orgName,
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

	if err := c.registrySystem.ValidateSite(netId.String(), req.SiteName, c.orgName); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid site or network ID: %s", err.Error())
	}

	route := c.controllerRoutingKey.SetAction("restart").SetObject("site").MustBuild()
	err = c.msgbus.PublishRequest(route, req)
	if err != nil {
		log.Errorf("Failed to publish message %+v with key %+v. Errors %s", req, route, err.Error())
	}
	return &pb.RestartSiteResponse{
		Status: pb.RestartStatus_ACCEPTED,
	}, nil
}

func (c *ControllerServer) RestartNode(ctx context.Context, req *pb.RestartNodeRequest) (*pb.RestartNodeResponse, error) {
	if req.NodeId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "node ID cannot be empty")
	}

	nodeId, err := uuid.FromString(req.GetNodeId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid node ID format: %s", err.Error())
	}

	if err := c.registrySystem.ValidateNode(nodeId.String(), c.orgName); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid node ID: %s", err.Error())
	}

	route := c.controllerRoutingKey.SetAction("restart").SetObject("node").MustBuild()
	err = c.msgbus.PublishRequest(route, req)
	if err != nil {
		log.Errorf("Failed to publish message %+v with key %+v. Errors %s", req, route, err.Error())
	}
	return &pb.RestartNodeResponse{
		Status: pb.RestartStatus_ACCEPTED,
	}, nil
}
