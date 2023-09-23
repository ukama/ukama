package server

import (
	"context"

	log "github.com/sirupsen/logrus"

	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
	pb "github.com/ukama/ukama/systems/node/configurator/pb/gen"

	"github.com/ukama/ukama/systems/node/configurator/pkg"
	"github.com/ukama/ukama/systems/node/configurator/pkg/providers"
)

type ConfiguratorServer struct {
	pb.UnimplementedConfiguratorServiceServer
	msgbus                 mb.MsgBusServiceClient
	registrySystem         providers.RegistryProvider
	configuratorRoutingKey msgbus.RoutingKeyBuilder
	debug                  bool
	orgName                string
}

func NewConfiguratorServer(msgBus mb.MsgBusServiceClient, registry providers.RegistryProvider, debug bool, orgName string) *ConfiguratorServer {
	return &ConfiguratorServer{
		configuratorRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
		msgbus:                 msgBus,
		registrySystem:         registry,
		debug:                  pkg.IsDebugMode,
		orgName:                orgName,
	}
}

func (c *ConfiguratorServer) ConfigEvent(ctx context.Context, req *pb.ConfigStoreEvent) (*pb.ConfigStoreEventResponse, error) {
	log.Infof("Received a event from config store %v", req)

	return &pb.ConfigStoreEventResponse{}, nil
}
