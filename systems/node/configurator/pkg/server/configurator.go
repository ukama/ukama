package server

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"

	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
	pb "github.com/ukama/ukama/systems/node/configurator/pb/gen"

	"github.com/ukama/ukama/systems/node/configurator/pkg"
	configstore "github.com/ukama/ukama/systems/node/configurator/pkg/configStore"
	"github.com/ukama/ukama/systems/node/configurator/pkg/db"
	"github.com/ukama/ukama/systems/node/configurator/pkg/providers"
)

type ConfiguratorServer struct {
	pb.UnimplementedConfiguratorServiceServer
	msgbus                 mb.MsgBusServiceClient
	registrySystem         providers.RegistryProvider
	configuratorRoutingKey msgbus.RoutingKeyBuilder
	debug                  bool
	orgName                string
	configStore            *configstore.ConfigStore
	commitRepo             db.CommitRepo
	configRepo             db.ConfigRepo
}

func NewConfiguratorServer(msgBus mb.MsgBusServiceClient, registry providers.RegistryProvider, cfgDb db.ConfigRepo, cmtDb db.CommitRepo, orgName string, url string, user string, pat string, t time.Duration, debug bool) *ConfiguratorServer {
	s, err := providers.NewStoreClient(url, user, pat, t)
	if err != nil {
		return nil
	}
	configStore := configstore.NewConfigStore(msgBus, registry, cfgDb, cmtDb, orgName, s, t)

	log.Infof("Config store created: %+v", configStore)
	return &ConfiguratorServer{
		configuratorRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
		msgbus:                 msgBus,
		registrySystem:         registry,
		debug:                  pkg.IsDebugMode,
		orgName:                orgName,
		configStore:            configStore,
		commitRepo:             cmtDb,
		configRepo:             cfgDb,
	}
}

func (c *ConfiguratorServer) ConfigEvent(ctx context.Context, req *pb.ConfigStoreEvent) (*pb.ConfigStoreEventResponse, error) {
	log.Infof("Received a event from config store %v", req)
	err := c.configStore.HandleConfigStoreEvent(ctx)
	if err != nil {
		log.Errorf("Error while handling config store event.Error: %s", err.Error())
	}
	return &pb.ConfigStoreEventResponse{}, err
}

func (c *ConfiguratorServer) ApplyConfig(ctx context.Context, req *pb.ApplyConfigRequest) (*pb.ApplyConfigResponse, error) {
	log.Infof("Received a request to apply config  %v", req)
	err := c.configStore.HandleConfigCommitReq(ctx, req.Hash)
	if err != nil {
		log.Errorf("Error while handling apply config req commit %s.Error: %s", req.Hash, err.Error())
	}
	return &pb.ApplyConfigResponse{}, err
}

func (c *ConfiguratorServer) GetConfigVersion(ctx context.Context, req *pb.ConfigVersionRequest) (*pb.ConfigVersionResponse, error) {
	log.Infof("Received a request to get config for node  %v", req)
	cfg, err := c.configRepo.Get(req.NodeId)
	if err != nil {
		log.Errorf("Error while reading config for node %s. Error: %s", req.NodeId, err.Error())
	}

	return &pb.ConfigVersionResponse{
		NodeId: req.NodeId,
		Status: cfg.State.String(),
		Commit: cfg.Commit.Hash,
		// LastStatus: cfg.LastStatus.String(),
		// LastCommit: cfg.LastCommit.Hash,
	}, err
}
