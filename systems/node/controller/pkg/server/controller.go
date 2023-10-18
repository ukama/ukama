package server

import (
	"context"

	log "github.com/sirupsen/logrus"
	cpb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/node/controller/pkg/providers"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"

	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/uuid"
	pb "github.com/ukama/ukama/systems/node/controller/pb/gen"

	"github.com/ukama/ukama/systems/node/controller/pkg"
	"github.com/ukama/ukama/systems/node/controller/pkg/db"
)

type ControllerServer struct {
	pb.UnimplementedControllerServiceServer
	sRepo          db.NodeLogRepo
	nodeFeederRoutingKey msgbus.RoutingKeyBuilder
	msgbus              mb.MsgBusServiceClient
	registrySystem       providers.RegistryProvider
	debug          bool
	orgName        string
}


func NewControllerServer(orgName string, sRepo db.NodeLogRepo,  msgBus mb.MsgBusServiceClient, debug bool,registry providers.RegistryProvider) *ControllerServer {
	return &ControllerServer{
		sRepo:             sRepo,
		orgName:             orgName,
		nodeFeederRoutingKey:      msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
		msgbus:              msgBus,
		debug:               debug,
		registrySystem:      registry,
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

	if err := c.registrySystem.ValidateNetwork(netId.String(), c.orgName); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid network ID: %s", err.Error())
	}

	route := c.nodeFeederRoutingKey.SetAction("restart").SetObject("site").MustBuild()
	anyMsg, err := anypb.New(req)
	if err != nil {
		return nil,err
	}
	msg := &cpb.NodeFeederMsg{
		Target:     c.orgName + "." + netId.String() + "." + req.SiteName ,
		HTTPMethod: "POST",
		Path:       "/v1/node/site/restart",
		Msg:        anyMsg,
	}

	err = c.msgbus.PublishRequest(route, msg)
	if err != nil {
		log.Errorf("Failed to publish message with key %+v. Errors %s", route, err.Error())
		return nil, err
	}
	log.Infof("Published controller on route %s for node %s ", msg, route)
	return &pb.RestartSiteResponse{
		Status: pb.RestartStatus_ACCEPTED,
	}, nil
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

	if err := c.registrySystem.ValidateNode(nId.String(), c.orgName); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid node ID: %s", err.Error())
	}
	_, err = c.sRepo.Get(nId.String())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Node has not been registered yet: %s", err.Error())
	}
	anyMsg,err:= anypb.New(req)
	if err != nil {
		return nil,err
	}
	route := c.nodeFeederRoutingKey.SetAction("restart").SetObject("node").MustBuild()
	msg:= &cpb.NodeFeederMsg{
		Target:     c.orgName + "." + nId.String(),
		HTTPMethod: "POST",
		Path:       "/v1/controller/restart/node",
		Msg:        anyMsg,
	}

	err = c.msgbus.PublishRequest(route, msg)
	if err != nil {
		log.Errorf("Failed to publish message with key %+v. Errors %s", route, err.Error())
		return nil, err
	}
	log.Infof("Published controller %s on route %s for node %s ", anyMsg, route, nId.String())
	return &pb.RestartNodeResponse{
		Status: pb.RestartStatus_ACCEPTED,
	}, nil
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

		if err := c.registrySystem.ValidateNode(nId.String(), c.orgName); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid node ID: %s", err.Error())
		}

		_, err = c.sRepo.Get(nId.String())
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "Node has not been registered yet: %s", err.Error())
		}
	}

	anyMsg,err:= anypb.New(req)
	if err != nil {
		return nil,err
	}
	route := c.nodeFeederRoutingKey.SetAction("restart").SetObject("nodes").MustBuild()
	msg:= &cpb.NodeFeederMsg{
		Target:     c.orgName,
		HTTPMethod: "POST",
		Path:       "/v1/controller/restart/nodes",
		Msg:        anyMsg,
	}

	err = c.msgbus.PublishRequest(route, msg)
	if err != nil {
		log.Errorf("Failed to publish message with key %+v. Errors %s", route, err.Error())
		return nil, err
	}

	log.Infof("Published controller %s on route %s for nodes %s ", anyMsg, route, req.NodeIds)
	
	return &pb.RestartNodesResponse{
		Status: pb.RestartStatus_ACCEPTED,
	}, nil
}
