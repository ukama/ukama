package server

import (
	"context"

	log "github.com/sirupsen/logrus"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
	pb "github.com/ukama/ukama/systems/ukama-agent/cdr/pb/gen"
	"github.com/ukama/ukama/systems/ukama-agent/cdr/pkg"
	"github.com/ukama/ukama/systems/ukama-agent/cdr/pkg/db"
)

type CDRServer struct {
	pb.UnimplementedCDRServiceServer
	cdrRepo        db.CDRRepo
	msgbus         mb.MsgBusServiceClient
	baseRoutingKey msgbus.RoutingKeyBuilder
	OrgName        string
	OrgId          string
}

func NewCDRServer(cdrRepo db.CDRRepo, orgId, orgName string, msgBus mb.MsgBusServiceClient) (*CDRServer, error) {

	cdr := CDRServer{
		cdrRepo: cdrRepo,
		OrgName: orgName,
		OrgId:   orgId,
		msgbus:  msgBus,
	}

	if msgBus != nil {
		cdr.baseRoutingKey = msgbus.NewRoutingKeyBuilder().SetRequestType().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName)
	}

	log.Infof("CDR server is %+v", cdr)

	return &cdr, nil
}

func (s *CDRServer) PostCDR(c context.Context, req *pb.CDR) (*pb.CDRResp, error) {

	return &pb.CDRResp{}, nil
}

func (s *CDRServer) GetCDR(c context.Context, req *pb.RecordReq) (*pb.RecordResp, error) {

	return &pb.RecordResp{}, nil
}

func (s *CDRServer) GetUsage(c context.Context, req *pb.UsageReq) (*pb.UsageResp, error) {

	return &pb.UsageResp{}, nil

}
