package server

import (
	"context"
	"sort"
	"time"

	log "github.com/sirupsen/logrus"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	"github.com/ukama/ukama/systems/common/sql"
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
		cdr.baseRoutingKey = msgbus.NewRoutingKeyBuilder().SetEventType().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName)
	}

	log.Infof("CDR server is %+v", cdr)

	return &cdr, nil
}

func (s *CDRServer) InitUsage(imsi string, policy string) error {
	ou, err := s.usageRepo.Get(imsi)
	if err != nil {
		if !sql.IsNotFoundError(err) {
			log.Errorf("Error getting usage for imsi %s. Error %+v", imsi, err)
			return err
		}
	}

	u := db.Usage{
		Imsi:        imsi,
		Policy:      policy,
		Usage:       0,
		LastSession: 0,
		Historical:  ou.Historical,
	}

	err = s.usageRepo.Add(&u)
	if err != nil {
		log.Errorf("Error initalizing usage for imsi %s. Error %+v", imsi, err)
		return err
	}
	log.Infof("Initilaize package usage for imsi %s to %+v", u.Imsi, u)

	return nil
}

func (s *CDRServer) PostCDR(c context.Context, req *pb.CDR) (*pb.CDRResp, error) {
	log.Debugf("Received CDR post request %+v", req)

	cdr := pbCDRToDbCDR(req)
	err := s.cdrRepo.Add(cdr)
	if err != nil {
		return nil, err
	}

	/* Publish event for new CDR */
	e := dbCDRToepbCDR(*cdr)
	if s.msgbus != nil {
		route := s.baseRoutingKey.SetActionCreate().SetObject("cdr").MustBuild()
		merr := s.msgbus.PublishRequest(route, e)
		if merr != nil {
			log.Errorf("Failed to publish message %+v with key %+v. Errors %s", e, route, err.Error())
		}
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
		Imsi:   req.Imsi,
		Usage:  usage.Usage,
		Policy: usage.Policy,
	}, nil
}

func (s *CDRServer) GetUsageDetails(c context.Context, req *pb.CycleUsageReq) (*pb.CycleUsageResp, error) {
	log.Debugf("Received get cycle usage request %+v", req)
	usage, err := s.usageRepo.Get(req.Imsi)
	if err != nil {
		return nil, err
	}
	return &pb.CycleUsageResp{
		Imsi:             req.Imsi,
		Usage:            usage.Usage,
		Historical:       usage.Historical,
		LastSessions:     usage.LastSession,
		LastNodeId:       usage.LastNodeId,
		LastCDRUpdatedAt: usage.LastCDRUpdatedAt,
		Policy:           usage.Policy,
	}, nil
}

func (s *CDRServer) ResetPackageUsage(imsi string, policy string) error {

	ou, err := s.usageRepo.Get(imsi)
	if err != nil {
		log.Errorf("Error getting usage for imsi %s. Error %+v", imsi, err)
		return err
	}

	u := db.Usage{
		Imsi:        imsi,
		Policy:      policy,
		Usage:       0,
		LastSession: 0,
		Historical:  ou.Historical,
	}

	err = s.usageRepo.Add(&u)
	if err != nil {
		log.Errorf("Error updating usage for imsi %s. Error %+v", imsi, err)
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

	log.Infof("Usage for imsi %s before CDR update is %+v and CDR is %+v ", imsi, ou, cdr)

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
		Imsi:             imsi,
		Usage:            ou.Usage,
		Historical:       ou.Historical,
		LastNodeId:       ou.LastNodeId,
		Policy:           ou.Policy,
		LastCDRUpdatedAt: ou.LastCDRUpdatedAt,
		LastSession:      ou.LastSession,
	}

	var lastUpdatedAt uint64 = 0
	var sessionId uint64 = 0
	var lastSessionId uint64 = 0
	tempUsage := db.Usage{}
	lastCDRNodeId := ou.LastNodeId
	var nodeChangedFlag bool = false
	var newSessionFlag bool = false

	//var policy string
	if len(cdrs) > 0 {
		if lastCDRNodeId != cdrs[0].NodeId {
			/* since new node is reporting the session usage this means new session is created */
			sessionId = cdrs[0].Session
			lastSessionId = ou.LastSession
		} else if cdrs[0].StartTime < ou.LastCDRUpdatedAt {
			// This means it's same session as last recorded session
			sessionId = cdrs[0].Session
			lastSessionId = cdrs[0].Session
		} else {
			/* new session */
			sessionId = cdrs[0].Session
			lastSessionId = ou.LastSession
		}
	} else {
		log.Infof("No usage update for imsi %s.Current usage stays: %+v", u.Imsi, u)
		return nil
	}

	for _, cdr := range cdrs {

		if ou.Policy != cdr.Policy {
			log.Errorf("CDR policy is %s not matching usage policy %s for imsi %s. Ignoring CDR record %+v", cdr.Policy, ou.Policy, imsi, cdr)
			continue
		} else {
			log.Infof("Handling CDR %v for imsi %s. Usage: %+v", cdr, imsi, u)
		}

		if cdr.NodeId == lastCDRNodeId {
			tempUsage = u
			if sessionId == lastSessionId {
				/* if session is continued */
				/* check to avoid duplicates updates */
				if cdr.LastUpdatedAt > lastUpdatedAt {
					lastUpdatedAt = cdr.LastUpdatedAt
					u.Historical = (u.Historical - u.Usage) + cdr.TotalBytes
					u.Usage = u.LastSession + cdr.TotalBytes
					u.LastNodeId = cdr.NodeId
					u.LastCDRUpdatedAt = cdr.LastUpdatedAt
					u.Policy = cdr.Policy

				}

			} else {
				/* New session */
				if cdr.LastUpdatedAt > lastUpdatedAt {
					lastUpdatedAt = cdr.LastUpdatedAt
					u.LastSession = u.Usage                      /* Usage till last session last cdr */
					u.Historical = u.Historical + cdr.TotalBytes /* usage is hitorical + current */
					u.Usage = u.LastSession + cdr.TotalBytes     /*usage for this package is last session + current */
					u.LastNodeId = cdr.NodeId
					u.LastCDRUpdatedAt = cdr.LastUpdatedAt
					u.Policy = cdr.Policy
					newSessionFlag = true
				}

			}

		} else {
			/* This will always be new session as new node is reporting cdr now */
			lastUpdatedAt = cdr.LastUpdatedAt
			u.LastSession = u.Usage                      /* Usage till last session last cdr */
			u.Historical = u.Historical + cdr.TotalBytes /* usage is hitorical + current */
			u.Usage = u.LastSession + cdr.TotalBytes
			u.LastNodeId = cdr.NodeId
			u.LastCDRUpdatedAt = cdr.LastUpdatedAt
			u.Policy = cdr.Policy
			newSessionFlag = true
			nodeChangedFlag = true
			//policy = cdr.Policy
		}

		/* If new session created send a event regading last session */
		if newSessionFlag == true && lastSessionId != 0 {
			log.Infof("Session %d is terminated for imsi %s. Session details are: %+v", lastSessionId, imsi, tempUsage)
			// Create session destroyed event
			e := &epb.SessionDestroyed{
				Imsi:         imsi,
				SessionUsage: tempUsage.Usage - tempUsage.LastSession,
				NodeId:       tempUsage.LastNodeId,
				Policy:       tempUsage.Policy,
				SessionId:    tempUsage.LastSession,
				TotalUsage:   tempUsage.Usage,
			}

			if s.msgbus != nil {
				route := s.baseRoutingKey.SetAction("terminated").SetObject("session").MustBuild()
				merr := s.msgbus.PublishRequest(route, e)
				if merr != nil {
					log.Errorf("Failed to publish message %+v with key %+v. Errors %s", e, route, err.Error())
				}
			}
			newSessionFlag = false
		}
		/* Session is changed */
		lastSessionId = sessionId
	}

	if nodeChangedFlag == true && lastCDRNodeId != "" {
		log.Infof("Imsi %s performed a node handover from %s to %s node", imsi, lastCDRNodeId, cdr.NodeId)
		/* subscriber node handover */
		e := &epb.NodeChanged{
			Imsi:              imsi,
			Policy:            cdr.Policy,
			NodeId:            cdr.NodeId,
			OldNodeId:         lastCDRNodeId,
			UsageTillLastNode: tempUsage.Usage,
			TotalUsage:        u.Usage,
		}

		if s.msgbus != nil {
			route := s.baseRoutingKey.SetActionCreate().SetObject("nodehandover").MustBuild()
			merr := s.msgbus.PublishRequest(route, e)
			if merr != nil {
				log.Errorf("Failed to publish message %+v with key %+v. Errors %s", e, route, err.Error())
			}
		}
		nodeChangedFlag = false
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
		Session:       req.Session,
		Imsi:          req.Imsi,
		NodeId:        req.NodeId,
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

func dbCDRToepbCDR(req db.CDR) *epb.CDRReported {

	pcdr := &epb.CDRReported{
		Session:       req.Session,
		Imsi:          req.Imsi,
		NodeId:        req.NodeId,
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
		NodeId:        req.NodeId,
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
