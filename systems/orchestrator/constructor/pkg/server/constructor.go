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
	oRepo          db.OrgsRepo
	dRepo          db.DeploymentsRepo
	sRepo          db.SystemsRepo
	msgbus         mb.MsgBusServiceClient
	baseRoutingKey msgbus.RoutingKeyBuilder
	pb.UnimplementedConstructorServiceServer
}

func NewConstructorServer(o db.OrgsRepo, d db.DeploymentsRepo, s db.SystemsRepo, msgBus mb.MsgBusServiceClient) *ConstructorServer {
	return &ConstructorServer{
		dRepo:          d,
		oRepo:          o,
		sRepo:          s,
		msgbus:         msgBus,
		baseRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetContainer(pkg.ServiceName),
	}
}

func (c *ConstructorServer) BuildOrg(ctx context.Context, req *pb.BuildOrgRequest) (*pb.BuildOrgResponse, error) {
	log.Infof("Build Org Id %s", req.GetOrgId())

	return &pb.BuildOrgResponse{}, nil
}

func (c *ConstructorServer) removeOrg(ctx context.Context, req *pb.RemoveOrgRequest) (*pb.RemoveOrgResponse, error) {
	log.Infof("Remove Org %s Id %s", req.GetOrgId())

	return &pb.RemoveOrgResponse{}, nil
}
