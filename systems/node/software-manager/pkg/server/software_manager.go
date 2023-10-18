package server

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/grpc"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
	cpb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	pb "github.com/ukama/ukama/systems/node/software-manager/pb/gen"
	"github.com/ukama/ukama/systems/node/software-manager/pkg"
	"github.com/ukama/ukama/systems/node/software-manager/pkg/db"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
)

type SoftwaManagerServer struct {
	pb.UnimplementedSoftwareManagerServiceServer
	sRepo                db.SoftwareManagerRepo
	NodeFeederRoutingKey msgbus.RoutingKeyBuilder
	msgbus               mb.MsgBusServiceClient
	debug                bool
	orgName              string
}

func NewSoftwareManagerServer(orgName string, sRepo db.SoftwareManagerRepo, msgBus mb.MsgBusServiceClient, debug bool) *SoftwaManagerServer {
	return &SoftwaManagerServer{
		sRepo:                sRepo,
		orgName:              orgName,
		NodeFeederRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
		msgbus:               msgBus,
		debug:                debug,
	}
}

func (s *SoftwaManagerServer) CreateSoftwareUpdate(ctx context.Context, req *pb.CreateSoftwareUpdateRequest) (*pb.CreateSoftwareUpdateResponse, error) {
	if req.Name == "" || req.Tag == "" {
		return nil, status.Errorf(codes.InvalidArgument,
			" Name, Tag, Description, ReleaseDate, Status")
	}

	log.Infof("Creating software update %s", req)
	//realesase date should be the current date time.time.now()
	softwareUpdate := &db.Software{
		Id:          uuid.NewV4(),
		Name:        req.Name,
		Tag:         req.Tag,
		ReleaseDate: time.Now(),
		Status:      db.Beta,
	}

	err := s.sRepo.CreateSoftwareUpdate(softwareUpdate, nil)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "Failed to create software update")
	}

	capps := &pb.CreateSoftwareUpdateRequest{
		Name: req.Name,
		Tag:  req.Tag,
	}
	route := s.NodeFeederRoutingKey.SetObject("node").SetAction("update").MustBuild()

	anyMsg, err := anypb.New(capps)

	msg := &cpb.NodeFeederMessage{
		Target:     s.orgName + "." + capps.Name + "." + capps.Tag,
		HTTPMethod: "POST",
		Path:       "/v1/node/update",
		Msg:        anyMsg,
	}
	err = s.msgbus.PublishRequest(route, msg)
	if err != nil {
		log.Errorf("Failed to publish message %+v with key %+v. Errors %s", req, route, err.Error())
	}

	return &pb.CreateSoftwareUpdateResponse{
		SoftwareUpdate: dbSoftwareToPbSoftwareUpdate(softwareUpdate),
	}, nil

}

func (s *SoftwaManagerServer) GetLatestSoftwareUpdate(ctx context.Context, req *pb.GetLatestSoftwareUpdateRequest) (*pb.GetLatestSoftwareUpdateResponse, error) {
	log.Infof("Getting latest software update")

	softwareUpdate, err := s.sRepo.GetLatestSoftwareUpdate()
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "Failed to get latest software update")
	}

	return &pb.GetLatestSoftwareUpdateResponse{
		SoftwareUpdate: dbSoftwareToPbSoftwareUpdate(softwareUpdate),
	}, nil

}

func dbSoftwareToPbSoftwareUpdate(software *db.Software) *pb.SoftwareUpdate {
	return &pb.SoftwareUpdate{
		Id:          software.Id.String(),
		Name:        software.Name,
		Tag:         software.Tag,
		ReleaseDate: software.ReleaseDate.String(),
		Status:      pb.Status(software.Status),
	}
}
