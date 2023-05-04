package server

import (
	"context"

	log "github.com/sirupsen/logrus"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
	pb "github.com/ukama/ukama/systems/orchestrator/constructor/pb/gen"
	"github.com/ukama/ukama/systems/orchestrator/constructor/pkg"
	"github.com/ukama/ukama/systems/orchestrator/constructor/pkg/db"
)

type ConstructorServer struct {
	systemRepo     db.SystemRepo
	orgRepo        db.OrgRepo
	nodeRepo       db.NodeRepo
	msgbus         mb.MsgBusServiceClient
	baseRoutingKey msgbus.RoutingKeyBuilder
	pb.UnimplementedConstructorServiceServer
}

func NewConstructorServer(nodeRepo db.NodeRepo, orgRepo db.OrgRepo, systemRepo db.SystemRepo, msgBus mb.MsgBusServiceClient) *ConstructorServer {
	return &ConstructorServer{
		nodeRepo:       nodeRepo,
		orgRepo:        orgRepo,
		systemRepo:     systemRepo,
		msgbus:         msgBus,
		baseRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetContainer(pkg.ServiceName),
	}
}

func (c *ConstructorServer) BuildOrg(ctx context.Context, req *pb.BuildOrgRequest) (*pb.BuildOrgResponse, error) {
	log.Infof("Build System For org %s Id %s", req.GetOrgName(), req.GetOrgId())

	return &pb.BuildOrgResponse{}, nil
}
