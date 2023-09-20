package server

import (
	"context"

	"github.com/cloudflare/cfssl/log"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/uuid"
	pb "github.com/ukama/ukama/systems/node/health/pb/gen"
	"github.com/ukama/ukama/systems/node/health/pkg"
	"github.com/ukama/ukama/systems/node/health/pkg/db"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type HealthServer struct {
	pb.UnimplementedHealhtServiceServer
	sRepo          db.HealthRepo
	msgBus         mb.MsgBusServiceClient
	baseRoutingKey msgbus.RoutingKeyBuilder
	debug          bool
	orgName        string
}

func NewHealthServer(msgBus mb.MsgBusServiceClient, debug bool, orgName string, sRepo db.HealthRepo) *HealthServer {
	return &HealthServer{
		sRepo:          sRepo,
		msgBus:         msgBus,
		baseRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
		debug:          debug,
	}
}

func (h *HealthServer) StoreRunningAppsInfo(ctx context.Context, req *pb.StoreRunningAppsInfoRequest) (*pb.StoreRunningAppsInfoResponse, error) {
	log.Infof("StoreRunningAppsInfo: %v", req)

	nodeUUID, err := uuid.FromString(req.GetNodeId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of node uuid. Error %s", err.Error())
	}

	for _, app := range req.GetRunningApps() {
		health := db.Health{
			NodeId:    nodeUUID,
			Name:      app.Name,
			Version:   app.Version,
			Status:    db.Status(app.Status),
			Timestamp: app.Timestamp,
		}

		err := h.sRepo.StoreRunningAppsInfo(&health, nil)
		if err != nil {
			return nil, err
		}
	}

	return &pb.StoreRunningAppsInfoResponse{}, nil
}

func (h *HealthServer) GetRunningAppsInfo(ctx context.Context, req *pb.GetRunningAppsRequest) (*pb.GetRunningAppsResponse, error) {
	log.Infof("GetRunningAppsInfo: %v", req)

	healths, err := h.sRepo.GetRunningAppsInfo()
	if err != nil {
		return nil, err
	}

	apps := make([]*pb.App, 0)
	for _, health := range healths {
		app := pb.App{
			Name:      health.Name,
			Version:   health.Version,
			Status:    pb.Status(health.Status),
			Timestamp: health.Timestamp,
		}
		apps = append(apps, &app)
	}

	return &pb.GetRunningAppsResponse{
		RunningApps: apps,
	}, nil
}
