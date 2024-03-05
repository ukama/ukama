package server

import (
	"context"
	"sort"
	"time"

	log "github.com/sirupsen/logrus"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
	dsql "github.com/ukama/ukama/systems/common/sql"
	pb "github.com/ukama/ukama/systems/ukama-agent/cdr/pb/gen"
	"github.com/ukama/ukama/systems/ukama-agent/cdr/pkg"
	"github.com/ukama/ukama/systems/ukama-agent/cdr/pkg/db"
)

type CDRServer struct {
	pb.UnimplementedCDRServiceServer
	cdrRepo        db.CDRRepo
	usageRepo      db.UsageRepo
	msgbus         mb.MsgBusServiceClient
	baseRoutingKey msgbus.RoutingKeyBuilder
	OrgName        string
	OrgId          string
}

func NewCDRServer(cdrRepo db.CDRRepo, usageRepo db.UsageRepo, orgId, orgName string, msgBus mb.MsgBusServiceClient) (*CDRServer, error) {

	cdr := CDRServer{
		cdrRepo:   cdrRepo,
		usageRepo: usageRepo,
		OrgName:   orgName,
		OrgId:     orgId,
		msgbus:    msgBus,
	}

	if msgBus != nil {
		cdr.baseRoutingKey = msgbus.NewRoutingKeyBuilder().SetRequestType().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName)
	}

	log.Infof("CDR server is %+v", cdr)

	return &cdr, nil
}

func (s *CDRServer) PostCDR(c context.Context, req *pb.CDR) (*pb.CDRResp, error) {
	log.Debugf("Received CDR post request %+v", req)

	cdr := pbCDRToDbCDR(req)
	err := s.cdrRepo.Add(cdr)
	if err != nil {
		return nil, err
	}

	err = s.UpdateUsage(req.Imsi, cdr)
	if err != nil {
		log.Errorf("Error updating usage for imsi %s", err)
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
	usage, err := s.usageRepo.Get(req.Imsi)
	if err != nil {
		return nil, err
	}
	return &pb.UsageResp{
		Imsi:  req.Imsi,
		Usage: usage.Usage,
	}, nil
}

func (s *CDRServer) ResetPackageUsage(imsi string) error {

	ou, err := s.usageRepo.Get(imsi)
	if err != nil {
		log.Errorf("Error getting usage for imsi %s. Error %+v", imsi, err)
		return err
	}

	u := db.Usage{
		Imsi:        imsi,
		Usage:       0,
		LastSession: 0,
		Historical:  ou.Historical,
	}

	err = s.usageRepo.Add(&u)
	if err != nil {
		log.Errorf("Error getting usage for imsi %s. Error %+v", imsi, err)
		return err
	}
	log.Infof("Reset package usage for imsi %s  from %+v to %+v", u.Imsi, ou, u)

	return nil
}

func (s *CDRServer) UpdateUsage(imsi string, cdr *db.CDR) error {

	ou, err := s.usageRepo.Get(imsi)
	if err != nil {
		if !dsql.IsNotFoundError(err) {
			log.Errorf("Error getting usage for imsi %s. Error %+v", imsi, err)
			return err
		} else {
			ou = &db.Usage{}
		}
	}

	recs, err := s.cdrRepo.GetByTime(cdr.Imsi, ou.LastCDRUpdatedAt, (uint64)(time.Now().Unix()))
	if err != nil && recs != nil {
		log.Errorf("Error getting CDR for imsi %s. Error %+v", imsi, err)
		return err
	}

	cdrs := *recs
	sort.Slice(cdrs, func(i, j int) bool {
		return cdrs[i].LastUpdatedAt < cdrs[j].LastUpdatedAt
	})

	u := db.Usage{
		Imsi:       imsi,
		Usage:      ou.Usage,
		Historical: ou.Historical,
	}

	var lastUpdatedAt uint64 = 0
	var session uint64 = 0
	var lastSession uint64 = 0
	if len(cdrs) > 1 {
		if cdrs[0].StartTime < ou.LastCDRUpdatedAt {
			// This means it's same session as last recorded session
			session = cdrs[0].Session
			lastSession = cdrs[0].Session

		} else {
			/* new session */
			session = cdrs[0].Session
			lastSession = 0
		}

	}

	for _, cdr := range cdrs {
		if session == lastSession {
			/* if session is continued */
			if cdr.LastUpdatedAt > lastUpdatedAt {
				lastUpdatedAt = cdr.LastUpdatedAt
				u.Historical = (u.Historical - u.Usage) + cdr.TotalBytes
				u.Usage = u.LastSession + cdr.TotalBytes
			}

		} else {
			/* New session */
			if cdr.LastUpdatedAt > lastUpdatedAt {
				lastUpdatedAt = cdr.LastUpdatedAt
				u.LastSession = u.Usage                      /* Usage till last session last cdr */
				u.Historical = u.Historical + cdr.TotalBytes /* usage is hitorical + current */
				u.Usage = u.LastSession + cdr.TotalBytes     /*usage for this package is last session + current */
			}

		}

		/* Session is changed */
		lastSession = session

	}

	err = s.usageRepo.Add(&u)
	if err != nil {
		log.Errorf("Error updating usage for imsi %s. Error %+v", imsi, err)
		return err
	}
	log.Infof("Updated usage for imsi %s to %+v", u.Imsi, u)

	return nil
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
