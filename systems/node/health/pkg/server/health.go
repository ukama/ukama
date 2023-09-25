package server

import (
	"context"

	"github.com/cloudflare/cfssl/log"
	"github.com/ukama/ukama/systems/common/grpc"
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
	healthID :=uuid.NewV4()
	cappID := uuid.NewV4()

    // Create a Health instance
    health := db.Health{
        Id:        healthID,
        NodeID:    nodeUUID,
        Timestamp: req.GetTimestamp(),
    }

    // Populate the System array from the request
    for _, sys := range req.GetSystem() {
        health.System = append(health.System, db.System{
            Id:    uuid.NewV4(),
            HealthID: healthID,
            Name:  sys.GetName(),
            Value: sys.GetValue(),
        })
    }

    for _, capp := range req.GetCapps() {
        health.Capps = append(health.Capps, db.Capp{
            Id:     cappID,
            HealthID: healthID,
            Name:   capp.GetName(),
            Tag:    capp.GetTag(),
            Status: db.Status(capp.GetStatus()),
        })

        for _, resource := range capp.GetResources() {
            health.Capps[len(health.Capps)-1].Resources = append(health.Capps[len(health.Capps)-1].Resources, db.Resource{
                Id:    uuid.NewV4(),
                CappID: cappID,
                Name:  resource.GetName(),
                Value: resource.GetValue(),
            })
        }
    }

    err = h.sRepo.StoreRunningAppsInfo(&health, nil)
    if err != nil {
        return nil, err
    }

    return &pb.StoreRunningAppsInfoResponse{}, nil
}



func (h *HealthServer) GetRunningApps(ctx context.Context, req *pb.GetRunningAppsRequest) (*pb.GetRunningAppsResponse, error) {
    log.Infof("GetRunningAppsInfo: %v", req)
    nodeUUID, err := uuid.FromString(req.GetNodeId())

    if err != nil {
        return nil, status.Errorf(codes.InvalidArgument,
            "invalid format of node uuid. Error %s", err.Error())
    }
    health, err := h.sRepo.GetRunningAppsInfo(nodeUUID)
    if err != nil {
        return nil, grpc.SqlErrorToGrpc(err, "health")

    }

    app := &pb.App{
        Id:        health.Id.String(),
        NodeId:    health.NodeID.String(),
        Timestamp: health.Timestamp,
        System:    []*pb.System{}, // Initialize System and Capps slices
        Capps:     []*pb.Capps{},
    }

    for _, sys := range health.System {
        system := &pb.System{
            Id:    sys.Id.String(),
			HealthId: health.Id.String(),
            Name:  sys.Name,
            Value: sys.Value,
        }
        app.System = append(app.System, system)
    }

    for _, capp := range health.Capps {
        capps := &pb.Capps{
            Id:     capp.Id.String(),
            Name:   capp.Name,
            Tag:    capp.Tag,
            Status: pb.Status(capp.Status), // Convert Status enum to string
            Resources: []*pb.Resource{},   // Initialize Resources slice
        }

        // Extract and format Resource data from Capps
        for _, resource := range capp.Resources {
            res := &pb.Resource{

                Id:    resource.Id.String(),
                Name:  resource.Name,
                Value: resource.Value,
				CappId: capp.Id.String(),
            }
            capps.Resources = append(capps.Resources, res)
        }

        app.Capps = append(app.Capps, capps)
    }

    return &pb.GetRunningAppsResponse{
        RunningApps: app, // Include the populated app in the response
    }, nil
}
