/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
package server

import (
	"context"
	"errors"
	"sort"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/common/grpc"
	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/ukama-agent/cdr/pkg"
	"github.com/ukama/ukama/systems/ukama-agent/cdr/pkg/client"
	"github.com/ukama/ukama/systems/ukama-agent/cdr/pkg/db"

	log "github.com/sirupsen/logrus"
	pmetric "github.com/ukama/ukama/systems/common/metrics"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	pb "github.com/ukama/ukama/systems/ukama-agent/cdr/pb/gen"
)

type CDRServer struct {
	pb.UnimplementedCDRServiceServer
	cdrRepo         db.CDRRepo
	usageRepo       db.UsageRepo
	asrClient       client.AsrService
	msgbus          mb.MsgBusServiceClient
	baseRoutingKey  msgbus.RoutingKeyBuilder
	OrgName         string
	OrgId           string
	pushGatewayHost string
}

func NewCDRServer(cdrRepo db.CDRRepo, usageRepo db.UsageRepo, orgId, orgName, pushGatewayHost string, asrClient client.AsrService, msgBus mb.MsgBusServiceClient) (*CDRServer, error) {

	cdr := CDRServer{
		cdrRepo:         cdrRepo,
		usageRepo:       usageRepo,
		asrClient:       asrClient,
		OrgName:         orgName,
		OrgId:           orgId,
		pushGatewayHost: pushGatewayHost,
		msgbus:          msgBus,
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
		ou = &db.Usage{} 
	}

	u := db.Usage{
		Imsi:             imsi,
		Policy:           policy,
		Usage:            ou.Usage,            
		LastSessionId:    ou.LastSessionId,    
		LastSessionUsage: ou.LastSessionUsage, 
		Historical:       ou.Historical,      
		LastNodeId:       ou.LastNodeId,       
		LastCDRUpdatedAt: ou.LastCDRUpdatedAt, 
	}

	err = s.usageRepo.Add(&u)
	if err != nil {
		log.Errorf("Error initialization usage for imsi %s. Error %+v", imsi, err)
		return err
	}

	log.Infof("Initialize/reactivate usage for imsi %s, preserved usage: %d bytes", u.Imsi, u.Usage)
	return nil
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

	usage, err := s.cdrRepo.QueryUsage(req.Imsi, "", 0, 0, 0, []string{cdr.Policy}, 0, false)
	if err != nil {
		return nil, err
	}

	asr, err := s.asrClient.GetAsr(cdr.Imsi)
	if err == nil && asr.Record != nil && asr.Record.Policy != nil && asr.Record.Policy.Uuid == cdr.Policy {
		labels := map[string]string{
			"package":  asr.Record.SimPackageId,
			"dataplan": asr.Record.PackageId,
			"network":  asr.Record.NetworkId,
			"iccid":    asr.Record.Iccid,
		}

		pushDataUsageMetrics(float64(usage), labels, s.pushGatewayHost)
	} else {
		log.Errorf("Failure while processing  ASR for policy %s : Skipping data usage metric push.",
			cdr.Policy)
	}

	/* Publish event for new CDR */
	e := dbCDRToepbCDR(*cdr)
	if s.msgbus != nil {
		route := s.baseRoutingKey.SetActionCreate().SetObject("cdr").MustBuild()
		merr := s.msgbus.PublishRequest(route, e)
		if merr != nil {
			log.Errorf("Failed to publish message %+v with key %+v. Errors %s", e, route, merr.Error())
		}
	}

	return &pb.CDRResp{}, nil
}

func (s *CDRServer) GetCDR(c context.Context, req *pb.RecordReq) (*pb.RecordResp, error) {
	log.Debugf("Received CDR get request %+v", req)
	cdrs, err := s.cdrRepo.GetByFilters(req.Imsi, req.SessionId, req.Policy, req.StartTime, req.EndTime)
	if err != nil {
		return nil, err
	}
	log.Debugf("CDR read: %+v", cdrs)
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
		LastSessionId:    usage.LastSessionId,
		LastSessionUsage: usage.LastSessionUsage,
		LastNodeId:       usage.LastNodeId,
		LastCDRUpdatedAt: usage.LastCDRUpdatedAt,
		Policy:           usage.Policy,
	}, nil
}

func (s *CDRServer) GetUsageForPeriod(c context.Context, req *pb.UsageForPeriodReq) (*pb.UsageForPeriodResp, error) {
	log.Debugf("Received request for usage during package %+v", req)
	usage, err := s.GetPeriodUsage(req.Imsi, req.StartTime, req.EndTime)
	if err != nil {
		return nil, err
	}
	return &pb.UsageForPeriodResp{
		Usage: usage,
	}, nil
}

func (s *CDRServer) ResetPackageUsage(imsi string, policy string) error {
	ou, err := s.usageRepo.Get(imsi)
	if err != nil {
		log.Errorf("Error getting usage for imsi %s. Error %+v", imsi, err)
		return err
	}

	var newUsage, newLastSessionUsage uint64
	if ou.Policy != policy {
		newUsage = 0
		newLastSessionUsage = 0
		log.Infof("Policy changed from %s to %s for imsi %s - resetting usage", ou.Policy, policy, imsi)
	} else {
		newUsage = ou.Usage
		newLastSessionUsage = ou.LastSessionUsage
		log.Infof("Same policy %s for imsi %s - preserving usage: %d bytes", policy, imsi, newUsage)
	}

	u := db.Usage{
		Imsi:             imsi,
		Policy:           policy,
		Usage:            newUsage,
		LastSessionId:    0, 
		LastSessionUsage: newLastSessionUsage,
		Historical:       ou.Historical,
		LastNodeId:       ou.LastNodeId,
		LastCDRUpdatedAt: ou.LastCDRUpdatedAt,
	}

	err = s.usageRepo.Add(&u)
	if err != nil {
		log.Errorf("Error updating usage for imsi %s. Error %+v", imsi, err)
		return err
	}

	log.Infof("Updated package usage for imsi %s from %+v to %+v", u.Imsi, ou, u)

	return nil
}
/* This API to be used for period start date to any time till end date */
func (s *CDRServer) GetPeriodUsage(imsi string, startTime uint64, endTime uint64) (uint64, error) {
	var lastSessionId uint64
	var lastNodeId string
	var usage uint64
	var lastUpdatedAt uint64
	var usageTillLastSession uint64
	recs, err := s.cdrRepo.GetByTime(imsi, startTime, endTime)
	if err != nil && recs != nil {
		log.Errorf("Error getting CDR for imsi %s. Error %+v", imsi, err)
		return 0, err
	}

	log.Debugf("Found %d CDR for imsi %s. CDR: %+v", len(*recs), imsi, recs)
	cdrs := *recs
	sort.Slice(cdrs, func(i, j int) bool {
		return cdrs[i].LastUpdatedAt < cdrs[j].LastUpdatedAt
	})

	if len(cdrs) > 0 {
		lastSessionId = cdrs[0].Session
		lastNodeId = cdrs[0].NodeId
		usage = cdrs[0].TotalBytes
		lastUpdatedAt = cdrs[0].LastUpdatedAt
		usageTillLastSession = 0
	}

	for _, cdr := range cdrs {
		log.Debugf("Handling CDR %+v", cdr)
		if lastNodeId == cdr.NodeId {
			if lastSessionId == cdr.Session {
				/* same as previous session */
				if cdr.LastUpdatedAt > lastUpdatedAt {
					usage = usageTillLastSession + cdr.TotalBytes
					lastUpdatedAt = cdr.LastUpdatedAt
				} else {
					log.Infof("Ignoring CDR %+v because last used CDR for usage was with LastUpdatedAt %d", cdr, lastUpdatedAt)
					continue
				}
				/* Handle end of session CDR */
				if cdr.EndTime != 0 {
					usageTillLastSession = usage
				}
			} else {
				/* New session */
				usageTillLastSession = usage
				if cdr.LastUpdatedAt > lastUpdatedAt {
					usage = usageTillLastSession + cdr.TotalBytes
					lastUpdatedAt = cdr.LastUpdatedAt
					lastSessionId = cdr.Session
				} else {
					log.Infof("Ignoring CDR %+v because last used CDR for usage was with LastUpdatedAt %d", cdr, lastUpdatedAt)
					continue
				}
				/* Handle end of session CDR */
				if cdr.EndTime != 0 {
					usageTillLastSession = usage
				}
			}
		} else {
			/* Node is changed this means new session */
			if cdr.LastUpdatedAt > lastUpdatedAt {
				usageTillLastSession = usage
				lastNodeId = cdr.NodeId
				usage = usageTillLastSession + cdr.TotalBytes
				lastUpdatedAt = cdr.LastUpdatedAt
				lastSessionId = cdr.Session
			} else {
				log.Infof("Ignoring CDR %+v because last used CDR for usage was with LastUpdatedAt %d", cdr, lastUpdatedAt)
				continue
			}
		}
	}

	log.Infof("Usage for imsi %s in time interval starting at %d and ending at %d is %d bytes", imsi, startTime, endTime, usage)
	return usage, nil
}

func (s *CDRServer) QueryUsage(c context.Context, req *pb.QueryUsageReq) (*pb.QueryUsageResp, error) {
	log.Debugf("Received Usage query request %+v", req)

	usage, err := s.cdrRepo.QueryUsage(req.Imsi, req.NodeId, req.Session, req.From, req.To, req.Policies, req.Count, req.Sort)
	if err != nil {
		if errors.Is(err, gorm.ErrInvalidData) {
			log.Errorf("Query usage failure: Inconsistent CDR in DB for request: %v", req)
			log.Warnf("Query Usage failure: You should check nodes correctly reported CDR for request: %v", req)

			return nil, status.Errorf(codes.OutOfRange,
				"query usage failure: inconsistent CDR(s) in DB for request: %v", req)
		}

		log.Errorf("Query usage failure: Error getting usage matching request %v. Error: %v", req, err)

		return nil, grpc.SqlErrorToGrpc(err, "query usage failure: Error getting usage matiching request:")
	}

	log.Debugf("usage query success: %+v", usage)

	return &pb.QueryUsageResp{
		Usage: usage,
	}, nil
}

/* If this function is getting really complex just drop this and use GetPeriodUsage which will read all the CDR from starttime to end time and report the usage */
func (s *CDRServer) UpdateUsage(imsi string, cdrMsg *db.CDR) error {
	ou, err := s.usageRepo.Get(imsi)
	if err != nil {
		if !sql.IsNotFoundError(err) {
			log.Errorf("Error getting usage for imsi %s. Error %+v", imsi, err)
			return err
		} else {
			ou = &db.Usage{}
		}
	}

	log.Infof("Usage for imsi %s before CDR update is %+v and CDR is %+v ", imsi, ou, cdrMsg)

	/* Assumption: Make sure older CDR is sent from the node first
	TODO: Handle the case if node A is not able to update CDR to the backend but
	node B on which subscriber latches after node A was able to publish CDR on backend
	In this case node A CDR will be rejected as of now
	*/
	recs, err := s.cdrRepo.GetByTimeAndNodeId(cdrMsg.Imsi, cdrMsg.StartTime, (uint64)(time.Now().Unix()), cdrMsg.NodeId)
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
		LastSessionId:    ou.LastSessionId,
		LastSessionUsage: ou.LastSessionUsage,
	}

	var lastUpdatedAt uint64 = 0
	// TODO: Removing default value of sessionId and lastSessionId "= 0" to address linting issue "ineffectual assignment to nodeChangedFlag (ineffassign)"
	var sessionId uint64
	var lastSessionId uint64
	tempUsage := db.Usage{}
	lastCDRNodeId := ou.LastNodeId
	var nodeChangedFlag bool
	var newSessionFlag = false

	//var policy string
	if len(cdrs) > 0 {
		//sessionId = cdrs[0].Session
		lastSessionId = ou.LastSessionId
		// if lastCDRNodeId != cdrs[0].NodeId {
		// 	/* since new node is reporting the session usage this means new session is created */
		// 	sessionId = cdrs[0].Session
		// 	lastSessionId = ou.LastSessionId
		// } else if cdrs[0].StartTime < ou.LastCDRUpdatedAt {
		// 	// This means it's same session as last recorded session
		// 	sessionId = cdrs[0].Session
		// 	lastSessionId = cdrs[0].Session
		// } else {
		// 	/* new session */
		// 	sessionId = cdrs[0].Session
		// 	lastSessionId = ou.LastSessionId
		// }
	} else {
		log.Infof("No usage update for imsi %s.Current usage stays: %+v", u.Imsi, u)
		return nil
	}

	for _, cdr := range cdrs {
		sessionId = cdr.Session

		if ou.Policy != cdr.Policy {
			log.Errorf("CDR policy is %s not matching usage policy %s for imsi %s. Ignoring CDR record %+v", cdr.Policy, ou.Policy, imsi, cdr)
			continue
		} else {
			log.Infof("Handling CDR %+v for imsi %s. Usage: %+v", cdr, imsi, u)
		}

		/* Ignoring old CDR
		This helps with multiple CDR from same session where start time is same
		but lastUpdatedAt is incremented in new CDR's
		*/
		if cdr.LastUpdatedAt <= ou.LastCDRUpdatedAt {
			log.Errorf("Last CDR updated was generated at %+v. Ignoring older CDR record %+v", ou.LastCDRUpdatedAt, cdr)
			continue
		}
		if cdr.NodeId == lastCDRNodeId {
			tempUsage = u
			if sessionId == lastSessionId {
				log.Infof("Session %d continued for imsi %s", sessionId, cdr.Imsi)
				/* if session is continued */
				/* check to avoid duplicates updates */
				if cdr.LastUpdatedAt > lastUpdatedAt {
					lastUpdatedAt = cdr.LastUpdatedAt
					u.Historical = (u.Historical - (u.Usage - u.LastSessionUsage)) + cdr.TotalBytes
					u.Usage = u.LastSessionUsage + cdr.TotalBytes
					u.LastNodeId = cdr.NodeId
					u.LastCDRUpdatedAt = cdr.LastUpdatedAt
					u.Policy = cdr.Policy
					/* If this report says session is ending too */
					if cdr.EndTime != 0 {
						u.LastSessionUsage = u.Usage
					}
				}

			} else {
				/* New session
				Assumption: We only allow the CDR which are generated(updated) after the last updated CDR in backend db
				We might have to check it in future if we miss any CDR
				We can still compile the report from CDR table as it contain all received CDR
				*/
				log.Infof("End session %d and create new session %d for imsi %s", ou.LastSessionId, sessionId, cdr.Imsi)
				if cdr.LastUpdatedAt > lastUpdatedAt {
					lastUpdatedAt = cdr.LastUpdatedAt
					u.LastSessionUsage = u.Usage                  /* Usage till last session last CDR */
					u.Historical = u.Historical + cdr.TotalBytes  /* usage is historical + current */
					u.Usage = u.LastSessionUsage + cdr.TotalBytes /*usage for this package is last session + current */
					u.LastNodeId = cdr.NodeId
					u.LastCDRUpdatedAt = cdr.LastUpdatedAt
					u.Policy = cdr.Policy
					u.LastSessionId = cdr.Session
					newSessionFlag = true
					/* If this report says session is ending too */
					if cdr.EndTime != 0 {
						u.LastSessionUsage = u.Usage
					}
				}

			}

		} else {
			/* This will always be new session as new node is reporting CDR now */
			log.Infof("End session %d and create new session %d for imsi %s because of node handover from %s to %s", ou.LastSessionId, sessionId, cdr.Imsi, lastCDRNodeId, cdr.NodeId)
			lastUpdatedAt = cdr.LastUpdatedAt
			u.LastSessionUsage = u.Usage                 /* Usage till last session last CDR */
			u.Historical = u.Historical + cdr.TotalBytes /* usage is historical + current */
			u.Usage = u.LastSessionUsage + cdr.TotalBytes
			u.LastNodeId = cdr.NodeId
			u.LastCDRUpdatedAt = cdr.LastUpdatedAt
			u.Policy = cdr.Policy
			u.LastSessionId = cdr.Session
			newSessionFlag = true
			nodeChangedFlag = true
			/* If this report says session is ending too */
			if cdr.EndTime != 0 {
				u.LastSessionUsage = u.Usage
			}
			//policy = cdr.Policy
		}

		/* If new session created send a event regading last session */
		if newSessionFlag && lastSessionId != 0 {
			log.Infof("Session %d is terminated for imsi %s. Session details are: %+v", lastSessionId, imsi, tempUsage)
			// Create session destroyed event
			e := &epb.SessionDestroyed{
				Imsi:         imsi,
				SessionUsage: tempUsage.Usage - tempUsage.LastSessionUsage,
				NodeId:       tempUsage.LastNodeId,
				Policy:       tempUsage.Policy,
				SessionId:    tempUsage.LastSessionId,
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

	if nodeChangedFlag && lastCDRNodeId != "" {
		log.Infof("Imsi %s performed a node handover from %s to %s node", imsi, lastCDRNodeId, cdrMsg.NodeId)
		/* subscriber node handover */
		e := &epb.NodeChanged{
			Imsi:              imsi,
			Policy:            cdrMsg.Policy,
			NodeId:            cdrMsg.NodeId,
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

		// TODO: Commented this out to address linting issue "ineffectual assignment to nodeChangedFlag (ineffassign)"
		// nodeChangedFlag = false
	}

	err = s.usageRepo.Add(&u)
	if err != nil {
		log.Errorf("Error updating usage for imsi %s. Error %+v", imsi, err)
		return err
	}

	log.Infof("Updated usage for imsi %s to %+v", u.Imsi, u)

	return nil
}

func pushDataUsageMetrics(value float64, labels map[string]string, pushGatewayHost string) {
	log.Infof("Collecting and pushing data usage metric to push gateway host: %s", pushGatewayHost)

	err := pmetric.CollectAndPushSimMetrics(pushGatewayHost, pkg.UsageMetrics,
		pkg.DataUsage, float64(value), labels, pkg.SystemName)
	if err != nil {
		log.Errorf("Error while pushing data usage  metric to push gateway %s", err.Error())
	}
}

func dbCDRToRecordResp(cdrs *[]db.CDR) *pb.RecordResp {
	if len(*cdrs) <= 0 {
		return &pb.RecordResp{}
	}
	pcdrs := make([]*pb.CDR, len(*cdrs))

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
		TxBytes:       req.TxBytes,
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
		TxBytes:       req.TxBytes,
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
		TxBytes:       req.TxBytes,
		RxBytes:       req.RxBytes,
		TotalBytes:    req.TotalBytes,
	}
}
