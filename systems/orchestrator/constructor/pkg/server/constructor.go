package server

import (
	"context"

	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
	pb "github.com/ukama/ukama/systems/orchestrator/constructor/pb/gen"
	"github.com/ukama/ukama/systems/orchestrator/constructor/pkg"
	"github.com/ukama/ukama/systems/orchestrator/constructor/pkg/db"
)

type ConstructorServer struct {
	oRepo          db.OrgRepo
	dRepo          db.DeploymentRepo
	cRepo          db.ConfigRepo
	msgbus         mb.MsgBusServiceClient
	baseRoutingKey msgbus.RoutingKeyBuilder
	pb.UnimplementedConstructorServiceServer
}

func NewConstructorServer(o db.OrgRepo, d db.DeploymentRepo, c db.ConfigRepo, msgBus mb.MsgBusServiceClient) *ConstructorServer {
	return &ConstructorServer{
		dRepo:          d,
		oRepo:          o,
		cRepo:          c,
		msgbus:         msgBus,
		baseRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetContainer(pkg.ServiceName),
	}
}

func ConstructOrg(ctx context.Context, in *pb.ConstructOrgRequest) (*pb.ConstructOrgResponse, error) {
	return &pb.ConstructOrgResponse{}, nil
}

func DistructOrg(ctx context.Context, in *pb.DistructOrgRequest) (*pb.DistructOrgResponse, error) {
	return &pb.DistructOrgResponse{}, nil
}

func Deployment(ctx context.Context, in *pb.DeploymentRequest) (*pb.DeploymentResponse, error) {
	return &pb.DeploymentResponse{}, nil
}

func GetDeployment(ctx context.Context, in *pb.GetDeploymentRequests) (*pb.GetDeploymentResponse, error) {
	return &pb.GetDeploymentResponse{}, nil
}

func RemoveDeployment(ctx context.Context, in *pb.RemoveDeploymentRequest) (*pb.RemoveDeploymentResponse, error) {
	return &pb.RemoveDeploymentResponse{}, nil
}

func AddConfig(ctx context.Context, in *pb.AddConfigRequest) (*pb.AddConfigResponse, error) {
	return &pb.AddConfigResponse{}, nil
}

func GetConfig(ctx context.Context, in *pb.GetConfigRequest) (*pb.GetConfigResponse, error) {
	return &pb.GetConfigResponse{}, nil
}

func GetDeploymentHistory(ctx context.Context, in *pb.GetDeploymentRequests) (*pb.GetDeploymentResponse, error) {
	return &pb.GetDeploymentResponse{}, nil
}
