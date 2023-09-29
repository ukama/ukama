package server

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/grpc"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/uuid"
	pb "github.com/ukama/ukama/systems/node/software-manager/pb/gen"
	"github.com/ukama/ukama/systems/node/software-manager/pkg"
	"github.com/ukama/ukama/systems/node/software-manager/pkg/db"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SoftwaManagerServer struct {
	pb.UnimplementedSoftwareManagerServiceServer
	sRepo          db.SoftwareManagerRepo
	msgBus         mb.MsgBusServiceClient
	baseRoutingKey msgbus.RoutingKeyBuilder
	debug          bool
	orgName        string
}

func NewSoftwareManagerServer(msgBus mb.MsgBusServiceClient, debug bool, orgName string, sRepo db.SoftwareManagerRepo) *SoftwaManagerServer {
	return &SoftwaManagerServer{
		sRepo:          sRepo,
		msgBus:         msgBus,
		baseRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
		debug:          debug,
	}
}

func (s *SoftwaManagerServer) CreateSoftware(ctx context.Context, req *pb.CreateSoftwareUpdateRequest) (*pb.CreateSoftwareUpdateResponse, error) {
	if req.Name == "" || req.Version == "" || req.Description == "" || req.Size == 0 || req.ReleaseDate == "" || req.NodeId == "" {
		return nil, status.Errorf(codes.InvalidArgument,
			" Name, Version, Description, Size, ReleaseDate, Status, NodeId are required")
	}

	releaseDate, err := time.Parse("2006-01-02", req.ReleaseDate)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"Invalid format for ReleaseDate. Error: %s", err.Error())
	}
	log.Infof("Creating software update %s", req)

	softwareUpdate := &db.Software{
		Id:          uuid.NewV4(),
		Name:        req.Name,
		Tag:     req.Version,
		Description: req.Description,
		Size:        req.Size,
		ReleaseDate: releaseDate,
		Status:      db.Status(req.Status),
	}

	err = s.sRepo.CreateSoftware(softwareUpdate, nil)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "Failed to create software update")
	}

	route := s.baseRoutingKey.SetActionCreate().SetObject("newUpdate").MustBuild()
	err = s.msgBus.PublishRequest(route, req)
	if err != nil {
		log.Errorf("Failed to publish message %+v with key %+v. Errors %s", req, route, err.Error())
	}

	return &pb.CreateSoftwareUpdateResponse{
		SoftwareUpdate: dbSoftwareToPbSoftwareUpdate(softwareUpdate),
	}, nil

}
func (s *SoftwaManagerServer) ReadSoftware(ctx context.Context, req *pb.ReadSoftwareUpdateRequest) (*pb.ReadSoftwareUpdateResponse, error) {
	log.Infof("Reading software update with id %s", req.Id)

	uuid, err := uuid.FromString(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of software uuid. Error %s", err.Error())
	}

	softwareUpdate, err := s.sRepo.ReadSoftware(uuid)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "Failed to read software update")
	}
	return &pb.ReadSoftwareUpdateResponse{
		SoftwareUpdate: dbSoftwareToPbSoftwareUpdate(softwareUpdate),
	}, nil

}
func (s *SoftwaManagerServer) ListSoftwares(ctx context.Context, req *pb.ListSoftwareUpdatesRequest) (*pb.ListSoftwareUpdatesResponse, error) {
	log.Infof("Listing software updates")

	softwareUpdates, err := s.sRepo.ListSoftwares()
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "Failed to list software updates")
	}

	var pbSoftwareUpdates []*pb.SoftwareUpdate
	for _, software := range softwareUpdates {
		pbSoftwareUpdates = append(pbSoftwareUpdates, dbSoftwareToPbSoftwareUpdate(software))
	}

	return &pb.ListSoftwareUpdatesResponse{
		SoftwareUpdates: pbSoftwareUpdates,
	}, nil

}
func (s *SoftwaManagerServer) GetLatestSoftware(ctx context.Context, req *pb.GetLatestSoftwareUpdateRequest) (*pb.GetLatestSoftwareUpdateResponse, error) {
	log.Infof("Getting latest software update")

	softwareUpdate, err := s.sRepo.GetLatestSoftware()
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "Failed to get latest software update")
	}

	return &pb.GetLatestSoftwareUpdateResponse{
		SoftwareUpdate: dbSoftwareToPbSoftwareUpdate(softwareUpdate),
	}, nil

}
func (s *SoftwaManagerServer) GetSoftwareByVersion(ctx context.Context, req *pb.GetLatestSoftwareUpdateByVersionRequest) (*pb.GetLatestSoftwareUpdateByVersionResponse, error) {
	log.Infof("Getting latest software update by version")

	softwareUpdate, err := s.sRepo.GetSoftwareByVersion(req.Version)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "Failed to get latest software update by version")
	}

	return &pb.GetLatestSoftwareUpdateByVersionResponse{
		SoftwareUpdate: dbSoftwareToPbSoftwareUpdate(softwareUpdate),
	}, nil

}
func (s *SoftwaManagerServer) GetSoftwareByAppName(ctx context.Context, req *pb.GetSoftwareUodateByAppNameRequest) (*pb.GetSoftwareUodateByAppNameResponse, error) {

	log.Infof("Getting software update by app name")

	if req.Name == "" {
		return nil, status.Errorf(codes.InvalidArgument,
			" Name is required")
	}

	softwareUpdate, err := s.sRepo.GetSoftwareByAppName(req.GetName())
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "Failed to get software update by app name")
	}

	return &pb.GetSoftwareUodateByAppNameResponse{
		SoftwareUpdate: dbSoftwareToPbSoftwareUpdate(softwareUpdate),
	}, nil

}

func dbSoftwareToPbSoftwareUpdate(software *db.Software) *pb.SoftwareUpdate {
	return &pb.SoftwareUpdate{
		Id:          software.Id.String(),
		Name:        software.Name,
		Version:     software.Tag,
		Description: software.Description,
		Size:        software.Size,
		ReleaseDate: software.ReleaseDate.String(),
		Status:      pb.Status(software.Status),
	}
}
