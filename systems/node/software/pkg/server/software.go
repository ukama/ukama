package server

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/grpc"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/uuid"
	pb "github.com/ukama/ukama/systems/node/software/pb/gen"
	"github.com/ukama/ukama/systems/node/software/pkg"
	"github.com/ukama/ukama/systems/node/software/pkg/db"
	providers "github.com/ukama/ukama/systems/node/software/pkg/provider"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SoftwareServer struct {
	pb.UnimplementedSoftwareServiceServer
	sRepo                db.SoftwareRepo
	nodeFeederRoutingKey msgbus.RoutingKeyBuilder
	msgbus               mb.MsgBusServiceClient
	wimsiClient        providers.WimsiClientProvider
	debug                bool
	orgName              string

}

func NewSoftwareServer(orgName string, sRepo db.SoftwareRepo, msgBus mb.MsgBusServiceClient, debug bool, wimsiClient providers.WimsiClientProvider) *SoftwareServer {
	return &SoftwareServer{
		sRepo:                sRepo,
		orgName:              orgName,
		nodeFeederRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
		msgbus:               msgBus,
		debug:                debug,
		wimsiClient:          wimsiClient,
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
		Space :      req.Space,
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

func (s *SoftwareServer) SoftwareUpdate(ctx context.Context, req *pb.UpdateSoftwareRequest) (*pb.UpdateSoftwareResponse, error) {
	log.Infof("Getting software update")

	softwareUpdate, err := s.sRepo.GetLatestSoftwareUpdate()
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "Failed to get software update")
	}

    err =s.wimsiClient.RequestSoftwareUpdate(softwareUpdate.Space,softwareUpdate.Tag,softwareUpdate.Name)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "Failed to request software update")
	}

	return &pb.UpdateSoftwareResponse{
	
	}, nil

}

func dbSoftwareToPbSoftwareUpdate(software *db.Software) *pb.SoftwareUpdate {
	return &pb.SoftwareUpdate{
		Id:          software.Id.String(),
		Name:        software.Name,
		Tag:         software.Tag,
		Space:       software.Space,
		Status:      pb.Status(software.Status),
	}
}
