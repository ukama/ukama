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
	log.Debugf("Received CDR post request %+v", req)
	err := s.cdrRepo.Add(pbCDRToDbCDR(req))
	if err != nil {
		return nil, err
	}
	return &pb.CDRResp{}, nil
}

func (s *CDRServer) GetCDR(c context.Context, req *pb.RecordReq) (*pb.RecordResp, error) {
	log.Debugf("Received CDR get request %+v", req)
	cdrs, err := s.cdrRepo.GetByFilters(req.Imsi, req.SessionId, req.Policy, req.StartTime, req.EndTime)
	if err != nil {
		return nil, err
	}
	return dbCDRToRecordResp(cdrs), nil
}

func (s *CDRServer) GetUsage(c context.Context, req *pb.UsageReq) (*pb.UsageResp, error) {
	log.Debugf("Received get usage request %+v", req)
	return &pb.UsageResp{}, nil
}

func dbCDRToRecordResp(cdrs *[]db.CDR) *pb.RecordResp {
	if len(*cdrs) > 0 {
		return &pb.RecordResp{}
	}
	pcdrs := make([]*pb.CDR, 0, len(*cdrs))

	for i, cdr := range *cdrs {
		pcdrs[i] = dbCDRTopbCDR(cdr)
	}

	return &pb.RecordResp{Cdr: pcdrs}
}

func dbCDRTopbCDR(req db.CDR) *pb.CDR {

	pcdr := &pb.CDR{
		Session: req.Session,
		Imsi:    req.Imsi,
		//NodeId:        req.NodeId, TBU
		Policy:        req.Policy,
		ApnName:       req.ApnName,
		Ip:            req.Ip,
		StartTime:     req.StartTime,
		EndTime:       req.EndTime,
		LastUpdatedAt: req.LastUpdatedAt,
		TxBytes:       req.TotalBytes,
		RxBytes:       req.RxBytes,
		TotalBytes:    req.TotalBytes,
	}

	return pcdr
}

func pbCDRToDbCDR(req *pb.CDR) *db.CDR {
	return &db.CDR{
		Session:       req.Session,
		Imsi:          req.Imsi,
		NodeId:        "TBU",
		Policy:        req.Policy,
		ApnName:       req.ApnName,
		Ip:            req.Ip,
		StartTime:     req.StartTime,
		EndTime:       req.EndTime,
		LastUpdatedAt: req.LastUpdatedAt,
		TxBytes:       req.TotalBytes,
		RxBytes:       req.RxBytes,
		TotalBytes:    req.TotalBytes,
	}
}
