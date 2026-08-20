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
	"fmt"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ukama/ukama/systems/common/grpc"
	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/rest/client/factory"
	"github.com/ukama/ukama/systems/common/rest/client/registry"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/ukama-agent/asr/pkg/client"
	"github.com/ukama/ukama/systems/ukama-agent/asr/pkg/db"

	log "github.com/sirupsen/logrus"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	pb "github.com/ukama/ukama/systems/ukama-agent/asr/pb/gen"
	pm "github.com/ukama/ukama/systems/ukama-agent/asr/pkg/policy"
)

const (
	activeSubscriberEventObject = "activesubscriber"
)

type AsrRecordServer struct {
	pb.UnimplementedAsrRecordServiceServer
	asrRepo    db.AsrRecordRepo
	gutiRepo   db.GutiRepo
	network    registry.NetworkClient
	factory    factory.SimFactoryClient
	cdr        client.CDRService
	msgbus     mb.MsgBusServiceClient
	pc         pm.Controller
	OrgName    string
	OrgId      string
	allowedToS int64
}

func NewAsrRecordServer(asrRepo db.AsrRecordRepo, gutiRepo db.GutiRepo, factory factory.SimFactoryClient, network registry.NetworkClient,
	pc pm.Controller, cdr client.CDRService, orgId, orgName string, msgBus mb.MsgBusServiceClient, aToS int64) (*AsrRecordServer, error) {
	asr := AsrRecordServer{
		asrRepo:    asrRepo,
		gutiRepo:   gutiRepo,
		OrgName:    orgName,
		OrgId:      orgId,
		factory:    factory,
		network:    network,
		msgbus:     msgBus,
		pc:         pc,
		cdr:        cdr,
		allowedToS: aToS,
	}

	log.Infof("Asr is %+v", asr)

	return &asr, nil
}

func (ar *AsrRecordServer) HandePostCDREvent(imsi string, policy string, session uint64) error {
	log.Infof("Handling POST CDR event for imsi %s", imsi)

	sub, err := ar.asrRepo.GetByImsi(imsi)
	if err != nil {
		log.Errorf("Error getting ASR profile for ismi %s.Error: %v", imsi, err)

		return grpc.SqlErrorToGrpc(err, "error getting ASR record for given imsi:")
	}

	r, err := ar.cdr.GetUsage(imsi)
	if err != nil {
		log.Errorf("Failed to get usage: %v for imsi %s", err, imsi)

		return fmt.Errorf("failed to get usage for imsi %s. Error: %w", imsi, err)
	}

	if r.Policy != sub.Policy.Id.String() {
		log.Errorf("Looks like sync failure for the subcriber %s. Policy expected %s is not matching CDR session %d",
			imsi, sub.Policy.Id.String(), session)

		return fmt.Errorf("request policy %s mismatches imsi policy %s", r.Policy, sub.Policy.Id.String())
	}

	//TODO

	return nil
}

func (ar *AsrRecordServer) Read(c context.Context, req *pb.ReadReq) (*pb.ReadResp, error) {
	log.Infof("Reading ASR data for imsi %s", req.GetImsi())

	var sub *db.Asr
	var err error

	switch req.Id.(type) {
	case *pb.ReadReq_Imsi:
		sub, err = ar.asrRepo.GetByImsi(req.GetImsi())
		if err != nil {
			return nil, grpc.SqlErrorToGrpc(err, "error getting ASR record for given imsi:")
		}

	case *pb.ReadReq_Iccid:
		sub, err = ar.asrRepo.GetByIccid(req.GetIccid())
		if err != nil {
			return nil, grpc.SqlErrorToGrpc(err, "error getting ASR record for given iccid:")
		}
	}

	r, err := ar.cdr.GetUsage(sub.Imsi)
	if err != nil {
		log.Errorf("Failed to get usage for imsi %s. Error: %v", req.GetImsi(), err)

		return nil, fmt.Errorf("failed to get usage for imsi %s. Error: %w", req.GetImsi(), err)
	}

	resp := &pb.ReadResp{Record: &pb.Record{
		Imsi:  sub.Imsi,
		Iccid: sub.Iccid,
		Key:   sub.Key,
		Amf:   sub.Amf,
		Op:    sub.Op,
		Apn: &pb.Apn{
			Name: sub.DefaultApnName,
		},
		AlgoType:          sub.AlgoType,
		CsgId:             sub.CsgId,
		CsgIdPrsent:       sub.CsgIdPrsent,
		Sqn:               sub.Sqn,
		UeDlAmbrBps:       sub.UeDlAmbrBps,
		UeUlAmbrBps:       sub.UeDlAmbrBps,
		IsServiceStatusOn: sub.IsServiceStatusOn,
		NetworkId:         sub.NetworkId.String(),
		PackageId:         sub.PackageId.String(),
		SimPackageId:      sub.SimPackageId.String(),
		Policy: &pb.Policy{
			Uuid:         sub.Policy.Id.String(),
			Burst:        sub.Policy.Burst,
			TotalData:    sub.Policy.TotalData,
			ConsumedData: r.Usage,
			Ulbr:         sub.Policy.Ulbr,
			Dlbr:         sub.Policy.Dlbr,
			StartTime:    sub.Policy.StartTime,
			EndTime:      sub.Policy.EndTime,
		},
	}}

	log.Infof("Active subscriber is  %+v", resp)
	return resp, nil
}

func (ar *AsrRecordServer) CreateProfile(ctx context.Context, req *pb.CreateProfileReq) (*pb.CreateProfileResp, error) {
	return createProfile(req.Iccid, req.Imsi, req.SimPackageId, req.PackageId, req.NetworkId, ar.network,
		ar.factory, ar.asrRepo, ar.pc, ar.allowedToS)
}

func (ar *AsrRecordServer) UpdatePackage(c context.Context, req *pb.UpdatePackageReq) (*pb.UpdatePackageResp, error) {
	log.Infof("Updating ASR profile package for imsi %s", req.GetImsi())

	asrRecord, err := ar.asrRepo.GetByIccid(req.GetIccid())
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "error getting ASR record for given iccid:")
	}

	/* We assume that packageId is validated by subscriber. */
	pId, err := uuid.FromString(req.PackageId)
	if err != nil {
		log.Errorf("PackageId not valid.")

		return nil, grpc.SqlErrorToGrpc(err, "error invalid package id:")
	}

	pcrfData := &pm.SimInfo{
		ID:        asrRecord.ID,
		Imsi:      asrRecord.Imsi,
		Iccid:     asrRecord.Iccid,
		PackageId: pId,
		NetworkId: asrRecord.NetworkId,
	}

	/* Create policy and send message to PCRF */
	policy, err := ar.pc.NewPolicy(pcrfData.PackageId)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "error creating policy:")
	}

	//Update policy
	err = ar.asrRepo.UpdatePackage(asrRecord.Imsi, pId, policy)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "error updating asr:")
	}

	err, removed := ar.pc.RunPolicyControl(asrRecord.Imsi, false)
	if err != nil {
		log.Errorf("Error running policy control for imsi %s. Error %v",
			asrRecord.Imsi, err)

		return nil, fmt.Errorf("error running policy control for imsi %s. Error %w",
			asrRecord.Imsi, err)
	}

	if removed {
		log.Infof("Profile removed from repo as one or more policies were failed for %s",
			asrRecord.Imsi)

		return nil, fmt.Errorf("profile removed from repo as one or more policies were failed for imsi %s",
			asrRecord.Imsi)
	}

	/* read the updated profile */
	nRec, err := ar.asrRepo.GetByIccid(req.GetIccid())
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "error getting ASR record for given iccid:")
	}

	err = ar.pc.SyncProfile(pcrfData, nRec, msgbus.ACTION_CRUD_UPDATE, activeSubscriberEventObject, true)
	if err != nil {
		log.Errorf("Failure to sync imsi %s pcrf profile for ASR package update. Error: %v",
			asrRecord.Imsi, err)

		//TODO: We need some kind of retry mecanism here

		return nil, fmt.Errorf("failure to sync imsi %s pcrf profile for ASR package update. Error: %w",
			asrRecord.Imsi, err)
	}

	asrRecord.Policy = *policy
	log.Debugf("Updated policy for %s imsi to %+v", asrRecord.Imsi, nRec)
	return &pb.UpdatePackageResp{}, nil
}

func (ar *AsrRecordServer) DeleteProfile(c context.Context, req *pb.DeleteProfileReq) (*pb.DeleteProfileResp, error) {
	return deleteProfile(req.Iccid, ar.asrRepo, ar.pc)
}

func (ar *AsrRecordServer) GetUsage(c context.Context, req *pb.UsageReq) (*pb.UsageResp, error) {
	log.Debugf("Received a usage request %+v", req)
	var sub *db.Asr
	var err error

	switch req.Id.(type) {
	case *pb.UsageReq_Imsi:

		sub, err = ar.asrRepo.GetByImsi(req.GetImsi())
		if err != nil {
			return nil, grpc.SqlErrorToGrpc(err, "error getting ASR record for given imsi:")
		}

	case *pb.UsageReq_Iccid:
		sub, err = ar.asrRepo.GetByIccid(req.GetIccid())
		if err != nil {
			return nil, grpc.SqlErrorToGrpc(err, "error getting ASR record for given iccid:")
		}
	}

	r, err := ar.cdr.GetUsage(sub.Imsi)
	if err != nil {
		log.Errorf("Failed to get usage: %v for imsi %s", err, req.GetImsi())

		return nil, fmt.Errorf("failed to get usage for imsi %s. Error: %w", req.GetImsi(), err)
	}

	return &pb.UsageResp{
		Usage: r.Usage,
	}, nil
}

func (ar *AsrRecordServer) GetUsageForPeriod(c context.Context, req *pb.UsageForPeriodReq) (*pb.UsageResp, error) {
	log.Debugf("Received a usage request for period %+v", req)

	var sub *db.Asr
	var err error

	switch req.Id.(type) {
	case *pb.UsageForPeriodReq_Imsi:

		sub, err = ar.asrRepo.GetByImsi(req.GetImsi())
		if err != nil {
			return nil, grpc.SqlErrorToGrpc(err, "error getting ASR record for given imsi:")
		}

	case *pb.UsageForPeriodReq_Iccid:
		sub, err = ar.asrRepo.GetByIccid(req.GetIccid())
		if err != nil {
			return nil, grpc.SqlErrorToGrpc(err, "error getting ASR record for given iccid:")
		}
	}

	r, err := ar.cdr.GetUsageForPeriod(sub.Imsi, req.StartTime, req.EndTime)
	if err != nil {
		log.Errorf("Failed to get usage for imsi %s. Error: %v", req.GetImsi(), err)

		return nil, fmt.Errorf("failed to get usage for imsi %s. Error: %w", req.GetImsi(), err)
	}

	return &pb.UsageResp{
		Usage: r.Usage,
	}, nil

}
func (ar *AsrRecordServer) QueryUsage(c context.Context, req *pb.QueryUsageReq) (*pb.QueryUsageResp, error) {
	log.Debugf("Received a query usage request: %+v", req)

	var sub *db.Asr
	var err error
	var policies []string

	sub, err = ar.asrRepo.GetByIccid(req.GetIccid())
	if err != nil {
		log.Errorf("Failed to query ASR record for given iccid : %s. Error: %v", req.Iccid, err)

		return nil, grpc.SqlErrorToGrpc(err, "query usage failure: Error getting ASR record for given iccid:")
	}

	policies = []string{sub.Policy.Id.String()}

	r, err := ar.cdr.QueryUsage(sub.Imsi, req.NodeId, req.Session, req.From, req.To, policies, req.Count, req.Sort)
	if err != nil {
		log.Errorf("Failed to query usage: for imsi %s. Error: %v", sub.Imsi, err)

		return nil, fmt.Errorf("failed to query usage for imsi %s. Error: %w", sub.Imsi, err)
	}

	return &pb.QueryUsageResp{
		Usage: r.Usage,
	}, nil
}

func (ar *AsrRecordServer) UpdateGuti(c context.Context, req *pb.UpdateGutiReq) (*pb.UpdateGutiResp, error) {
	log.Infof("Updating DUTI for imsi %s", req.GetImsi())

	_, err := ar.asrRepo.GetByImsi(req.Imsi)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "imsi")
	}

	err = ar.gutiRepo.Update(&db.Guti{
		Imsi:            req.Imsi,
		PlmnId:          req.Guti.PlmnId,
		Mmegi:           req.Guti.Mmegi,
		Mmec:            req.Guti.Mmec,
		MTmsi:           req.Guti.Mtmsi,
		DeviceUpdatedAt: time.Unix(int64(req.UpdatedAt), 0),
	})
	if err != nil {
		log.Errorf("Failed to update GUTI: %v", err)

		if err.Error() == db.GutiNotUpdatedErr {
			return nil, status.Errorf(codes.AlreadyExists, "%v", err)
		}

		return nil, grpc.SqlErrorToGrpc(err, "guti")
	}

	return &pb.UpdateGutiResp{}, nil
}

func (ar *AsrRecordServer) UpdateTai(c context.Context, req *pb.UpdateTaiReq) (*pb.UpdateTaiResp, error) {
	log.Infof("Updating TAI for imsi %s", req.GetImsi())

	_, err := ar.asrRepo.GetByImsi(req.Imsi)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "imsi")
	}

	err = ar.asrRepo.UpdateTai(req.Imsi, db.Tai{
		PlmnId:          req.PlmnId,
		Tac:             req.Tac,
		DeviceUpdatedAt: time.Unix(int64(req.UpdatedAt), 0),
	})

	if err != nil {
		log.Errorf("Failed to update TAI: %v", err)

		if err.Error() == db.TaiNotUpdatedErr {
			return nil, status.Errorf(codes.AlreadyExists, "%v", err)
		}

		return nil, grpc.SqlErrorToGrpc(err, "tai")
	}

	return &pb.UpdateTaiResp{}, nil
}

func (ar *AsrRecordServer) UpdateAndSyncAsrProfileFromCdr(imsi string) error {
	log.Infof("Updating and syncing ASR profile for imsi %s", imsi)

	sub, err := ar.asrRepo.GetByImsi(imsi)
	if err != nil {
		log.Errorf("ASR record not found for imsi %s. Error: %v", imsi, err)

		return grpc.SqlErrorToGrpc(err, "error getting ASR record by imsi:")
	}

	r, err := ar.cdr.GetUsage(imsi)
	if err != nil {
		log.Errorf("Failed to get usage: %v for imsi %s", err, imsi)

		return fmt.Errorf("failed to get usage for imsi %s. Error: %w", imsi, err)
	}

	sub.Policy.ConsumedData = r.Usage

	err = ar.asrRepo.Update(imsi, sub)
	if err != nil {
		log.Errorf("Failed to update usage: %v for imsi %s. Error: %v", r, imsi, err)

		return fmt.Errorf("failed to update usage for imsi %s. Error: %w", imsi, err)
	}

	err, removed := ar.pc.RunPolicyControl(imsi, false)
	if err != nil {
		log.Errorf("Error running policy control for imsi %s. Error: %v", sub.Imsi, err)

		return fmt.Errorf("error running policy control for imsi %s. Error: %w", sub.Imsi, err)
	}

	if removed {
		log.Infof("Profile removed from repo as one or more policies were failed for imsi %s", sub.Imsi)

		return fmt.Errorf("profile removed from repo as one or more policies were failed for imsi %s", sub.Imsi)
	}

	pcrfData := &pm.SimInfo{
		ID:        sub.ID,
		Imsi:      sub.Imsi,
		Iccid:     sub.Iccid,
		PackageId: sub.PackageId,
		NetworkId: sub.NetworkId,
	}

	err = ar.pc.SyncProfile(pcrfData, sub, msgbus.ACTION_CRUD_UPDATE, activeSubscriberEventObject, false)
	if err != nil {
		log.Errorf("Failure to sync imsi %s pcrf profile for ASR update. Error: %v", sub.Imsi, err)

		return fmt.Errorf("failure to sync imsi %s pcrf profile for ASR update. Error: %w", sub.Imsi, err)
	}

	return nil
}

func createProfile(iccid, imsi, packageId, dataPlanId, networkId string, networkClient registry.NetworkClient,
	factoryClient factory.SimFactoryClient, asrRepo db.AsrRecordRepo, policyController pm.Controller, allowedToS int64) (*pb.CreateProfileResp, error) {
	log.Infof("Adding ASR profile for iccid %s", iccid)

	/* Package DataPlan Id */
	pId, err := uuid.FromString(dataPlanId)
	if err != nil {
		log.Errorf("PackageId not valid: %s", dataPlanId)
		return nil, fmt.Errorf("packageId %s not valid.`Error:%w ", dataPlanId, err)
	}

	/* Sim Package Id */
	spId, err := uuid.FromString(packageId)
	if err != nil {
		log.Errorf("SimPackageId not valid: %s", packageId)

		return nil, fmt.Errorf("sim packageId %s not valid.`Error:%w ", packageId, err)
	}

	/* NetworkId */
	nId, err := uuid.FromString(networkId)
	if err != nil {
		log.Errorf("NetworkId not valid: %s", networkId)

		return nil, fmt.Errorf("networkId %s not valid.`Error:%w ", networkId, err)
	}

	// Fetch network details from registry
	_, err = networkClient.Get(networkId)
	if err != nil {
		return nil, fmt.Errorf("error while fetching network %s info: %w", networkId, err)
	}

	// network-org validation is no longer needed since we are using initClient to fetch
	// the correct registry system that matches with the current running org.

	/* Send Request to SIM Factory */
	sim, err := factoryClient.ReadSimCardInfo(iccid)
	if err != nil {
		return nil, fmt.Errorf("error reading sim (iccid: %s) info from factory. Error: %w", iccid, err)
	}

	pcrfData := &pm.SimInfo{
		Imsi:      sim.Imsi,
		Iccid:     sim.Iccid,
		PackageId: pId,
		NetworkId: nId,
		Visitor:   false, // We will using this flag on roaming in VLR
	}

	/* Send message to PCRF */
	policy, err := policyController.NewPolicy(pcrfData.PackageId)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "error creating policy:")
	}

	/* Add to ASR */
	asr := &db.Asr{
		Iccid:                   iccid,
		Imsi:                    sim.Imsi,
		Op:                      sim.Op,
		Key:                     sim.Key,
		Amf:                     sim.Amf,
		AlgoType:                sim.AlgoType,
		UeDlAmbrBps:             sim.UeDlAmbrBps,
		UeUlAmbrBps:             sim.UeUlAmbrBps,
		Sqn:                     sim.Sqn,
		CsgIdPrsent:             sim.CsgIdPrsent,
		CsgId:                   sim.CsgId,
		DefaultApnName:          sim.DefaultApnName,
		PackageId:               pId,
		SimPackageId:            spId,
		NetworkId:               nId,
		Policy:                  *policy,
		LastStatusChangeAt:      time.Now(),
		AllowedTimeOfService:    allowedToS,
		IsServiceStatusOn:       true,
		LastStatusChangeReasons: db.PROFILE_CREATION,
	}

	//Create profile
	err = asrRepo.Add(asr)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "error updating asr:")
	}

	err, removed := policyController.RunPolicyControl(asr.Imsi, false)
	if err != nil {
		log.Errorf("Error running policy control for imsi %s. Error %v", asr.Imsi, err)

		return nil, fmt.Errorf("error running policy control for imsi %s. Error: %w", asr.Imsi, err)
	}

	if removed {
		log.Infof("Profile not added to repo as one or more policies were failed for imsi %s", asr.Imsi)

		return nil, fmt.Errorf("profile not added to repo as one or more policies were failed for imsi %s", asr.Imsi)
	}

	err = policyController.SyncProfile(pcrfData, asr, msgbus.ACTION_CRUD_CREATE, activeSubscriberEventObject, true)
	if err != nil {
		log.Warnf("Failure to sync imsi %s pcrf profile for ASR activation. Error: %v", asr.Imsi, err)
		log.Warnf("We'll have to rely on pcrf subscriber reverse lookup!")
	}

	log.Debugf("Activated %s imsi with %+v", asr.Imsi, asr)

	return &pb.CreateProfileResp{}, err
}

func deleteProfile(iccid string, asrRepo db.AsrRecordRepo, policyController pm.Controller) (*pb.DeleteProfileResp, error) {
	log.Infof("Removing ASR profile for Iccid %s", iccid)

	delAsrRecord, err := asrRepo.GetByIccid(iccid)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "error getting ASR record for given iccid:")
	}

	pcrfData := &pm.SimInfo{
		ID:        delAsrRecord.ID,
		Imsi:      delAsrRecord.Imsi,
		Iccid:     delAsrRecord.Iccid,
		NetworkId: delAsrRecord.NetworkId,
	}

	err = asrRepo.Delete(delAsrRecord.Imsi, db.PROFILE_DELETION)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "error updating asr:")
	}

	err = policyController.SyncProfile(pcrfData, delAsrRecord, msgbus.ACTION_CRUD_DELETE, activeSubscriberEventObject, true)
	if err != nil {
		log.Errorf("Failure to sync imsi %s pcrf profile for ASR deactivation. Error: %v",
			delAsrRecord.Imsi, err)

		//TODO: We need some kind of retry mecanism here

		return nil, fmt.Errorf("failure to sync imsi %s pcrf profile for ASR deactivation. Error: %w",
			delAsrRecord.Imsi, err)
	}

	log.Debugf("Deleted subscriber %+v", delAsrRecord)

	return &pb.DeleteProfileResp{}, nil
}
