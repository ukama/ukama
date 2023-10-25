package server

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/grpc"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
	cpb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
	"github.com/ukama/ukama/systems/common/ukama"

	"github.com/ukama/ukama/systems/common/uuid"
	hpb "github.com/ukama/ukama/systems/node/health/pb/gen"
	pb "github.com/ukama/ukama/systems/node/software/pb/gen"
	"github.com/ukama/ukama/systems/node/software/pkg"
	"github.com/ukama/ukama/systems/node/software/pkg/db"
	providers "github.com/ukama/ukama/systems/node/software/pkg/provider"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

type SoftwareServer struct {
	pb.UnimplementedSoftwareServiceServer
	sRepo                db.SoftwareRepo
	nodeFeederRoutingKey msgbus.RoutingKeyBuilder
	msgbus               mb.MsgBusServiceClient
	debug                bool
	orgName              string
	healthService        providers.HealthClientProvider
}

func NewSoftwareServer(orgName string, sRepo db.SoftwareRepo, msgBus mb.MsgBusServiceClient, debug bool, healthService providers.HealthClientProvider) *SoftwareServer {
	return &SoftwareServer{
		sRepo:                sRepo,
		orgName:              orgName,
		nodeFeederRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
		msgbus:               msgBus,
		debug:                debug,
		healthService:        healthService,
	}
}

func (s *SoftwareServer) CreateSoftwareUpdate(ctx context.Context, req *pb.CreateSoftwareUpdateRequest) (*pb.CreateSoftwareUpdateResponse, error) {
	if req.Name == "" || req.Tag == "" {
		return nil, status.Errorf(codes.InvalidArgument,
			" Name, Tag, Description, ReleaseDate, Status")
	}

	log.Infof("Creating software update %s", req)
	//realesase date should be the current date time.time.now()
	softwareUpdate := &db.Software{
		Id:          uuid.NewV4(),
		Name:        req.Name,
		Space:       req.Space,
		Tag:         req.Tag,
		ReleaseDate: time.Now(),
		Status:      db.Beta,
	}

	err := s.sRepo.CreateSoftwareUpdate(softwareUpdate, nil)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "Failed to create software update")
	}

	return &pb.CreateSoftwareUpdateResponse{
		SoftwareUpdate: dbSoftwareToPbSoftwareUpdate(softwareUpdate),
	}, nil

}

func (s *SoftwareServer) GetLatestSoftwareUpdate(ctx context.Context, req *pb.GetLatestSoftwareUpdateRequest) (*pb.GetLatestSoftwareUpdateResponse, error) {
	log.Infof("Getting latest software update")

	softwareUpdate, err := s.sRepo.GetLatestSoftwareUpdate()
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "Failed to get latest software update")
	}

	return &pb.GetLatestSoftwareUpdateResponse{
		SoftwareUpdate: dbSoftwareToPbSoftwareUpdate(softwareUpdate),
	}, nil

}

func (s *SoftwareServer) UpdateSoftware(ctx context.Context, req *pb.UpdateSoftwareRequest) (*pb.UpdateSoftwareResponse, error) {
	log.Infof("Getting software update")

	nId, err := ukama.ValidateNodeId(req.NodeId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of node id. Error %s", err.Error())
	}

	softwareUpdate, err := s.sRepo.GetLatestSoftwareUpdate()
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "Failed to get software update")
	}

	svc, err := s.healthService.GetClient()
	if err != nil {
		return nil, err
	}
	runningApps, err := svc.GetRunningApps(ctx, &hpb.GetRunningAppsRequest{
		NodeId: nId.String(),
	})

	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "Failed to get running apps")
	}

	for _, capp := range runningApps.RunningApps.Capps {
		log.Infof("Running app %s", capp.Name)
		softReq := &pb.UpdateSoftwareRequest{
			NodeId: runningApps.RunningApps.NodeId,
			Tag:    capp.Tag,
			Name:   capp.Name,
			Space:  capp.Space,
		}
		if capp.Tag == softwareUpdate.Tag {
			log.Infof("App %s is already running and tag %s", capp.Name, capp.Tag)
			msg := fmt.Sprintf("Capp %s is already running and tag %s", capp.Name, capp.Tag)
			return &pb.UpdateSoftwareResponse{
				Message: msg,
			}, nil
		}

		data, err := proto.Marshal(softReq)
		if err != nil {
			log.Fatalf("Failed to marshal message: %v", err)
		}
		err = s.publishMessage(s.orgName+"."+"."+"."+runningApps.RunningApps.NodeId, data, capp.Space, capp.Name, capp.Tag)
		if err != nil {
			log.Errorf("Failed to publish message. Errors %s", err.Error())
			return nil, status.Errorf(codes.Internal, "Failed to publish message: %s", err.Error())

		}
	}
	return &pb.UpdateSoftwareResponse{
		Message: "Software updated successfully",
	}, nil

}

func dbSoftwareToPbSoftwareUpdate(software *db.Software) *pb.SoftwareUpdate {
	return &pb.SoftwareUpdate{
		Id:     software.Id.String(),
		Name:   software.Name,
		Tag:    software.Tag,
		Space:  software.Space,
		Status: pb.Status(software.Status),
	}
}

func (c *SoftwareServer) publishMessage(target string, anyMsg []byte, space string, name string, tag string) error {
	route := "request.cloud.local" + "." + c.orgName + "." + pkg.SystemName + "." + pkg.ServiceName + "." + "nodefeeder" + "." + "publish"

	msg := &cpb.NodeFeederMessage{
		Target:     target,
		HTTPMethod: "POST",
		Path:       "wimc/v1/ping",
		Msg:        anyMsg,
	}
	log.Infof("Published controller %s on route %s on target %s ", anyMsg, route, target)

	err := c.msgbus.PublishRequest(route, msg)
	return err
}
