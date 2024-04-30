package server

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/grpc"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/uuid"
	pb "github.com/ukama/ukama/systems/ukama-agent/asr/pb/gen"
	"github.com/ukama/ukama/systems/ukama-agent/asr/pkg"
	"github.com/ukama/ukama/systems/ukama-agent/asr/pkg/client"
	"github.com/ukama/ukama/systems/ukama-agent/asr/pkg/db"
	pm "github.com/ukama/ukama/systems/ukama-agent/asr/pkg/policy"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AsrRecordServer struct {
	pb.UnimplementedAsrRecordServiceServer
	asrRepo        db.AsrRecordRepo
	pRepo          db.PolicyRepo
	gutiRepo       db.GutiRepo
	network        client.Network
	factory        client.Factory
	cdr            client.CDRService
	msgbus         mb.MsgBusServiceClient
	baseRoutingKey msgbus.RoutingKeyBuilder
	pc             pm.Controller
	OrgName        string
	OrgId          string
	allowedToS     int64
}

func NewAsrRecordServer(asrRepo db.AsrRecordRepo, gutiRepo db.GutiRepo, pRepo db.PolicyRepo, factory client.Factory, network client.Network, pc pm.Controller, cdr client.CDRService, orgId, orgName string, msgBus mb.MsgBusServiceClient, aToS int64) (*AsrRecordServer, error) {

	asr := AsrRecordServer{
		asrRepo:    asrRepo,
		gutiRepo:   gutiRepo,
		OrgName:    orgName,
		OrgId:      orgId,
		factory:    factory,
		network:    network,
		pRepo:      pRepo,
		msgbus:     msgBus,
		pc:         pc,
		cdr:        cdr,
		allowedToS: aToS,
	}

	if msgBus != nil {
		asr.baseRoutingKey = msgbus.NewRoutingKeyBuilder().SetEventType().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName)
	}

	log.Infof("Asr is %+v", asr)

	return &asr, nil
}

func (s *AsrRecordServer) HandePostCDREvent(imsi string, policy string, session uint64) error {
	sub, err := s.asrRepo.GetByImsi(imsi)
	if err != nil {
		log.Errorf("Error getting ASR profile for ismi %s.Error: %v", imsi, err)
		return grpc.SqlErrorToGrpc(err, "error getting imsi")
	}

	r, err := s.cdr.GetUsage(imsi)
	if err != nil {
		log.Errorf("Failed to get usage: %v for imsi %s", err, imsi)
		return err
	}

	if r.Policy != sub.Policy.Id.String() {
		log.Errorf("Looks like sync failure for the subcriber %s. Policy expected %s is not matching CDR session %d", imsi, sub.Policy.Id.String(), session)
		return fmt.Errorf("Policy mismatch.")
	}

	//TODO

	return nil
}

func (s *AsrRecordServer) Read(c context.Context, req *pb.ReadReq) (*pb.ReadResp, error) {
	var sub *db.Asr
	var err error

	switch req.Id.(type) {
	case *pb.ReadReq_Imsi:

		sub, err = s.asrRepo.GetByImsi(req.GetImsi())
		if err != nil {
			return nil, grpc.SqlErrorToGrpc(err, "error getting imsi")
		}

	case *pb.ReadReq_Iccid:
		sub, err = s.asrRepo.GetByIccid(req.GetIccid())
		if err != nil {
			return nil, grpc.SqlErrorToGrpc(err, "error getting iccid")
		}
	}

	r, err := s.cdr.GetUsage(sub.Imsi)
	if err != nil {
		log.Errorf("Failed to get usage: %v for imsi %s", err, req.GetImsi())
		return nil, err
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
		AlgoType:    sub.AlgoType,
		CsgId:       sub.CsgId,
		CsgIdPrsent: sub.CsgIdPrsent,
		Sqn:         sub.Sqn,
		UeDlAmbrBps: sub.UeDlAmbrBps,
		UeUlAmbrBps: sub.UeDlAmbrBps,
		PackageId:   sub.PackageId.String(),
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

func (s *AsrRecordServer) Activate(c context.Context, req *pb.ActivateReq) (*pb.ActivateResp, error) {

	/* Send Request to SIM Factory */
	sim, err := s.factory.ReadSimCardInfo(req.Iccid)
	if err != nil {
		return nil, fmt.Errorf("error reading iccid from factory")
	}

	/* Validate network in Org */
	err = s.network.ValidateNetwork(req.NetworkId, s.OrgId)
	if err != nil {
		return nil, fmt.Errorf("error validating network")
	}

	nId, err := uuid.FromString(req.NetworkId)
	if err != nil {
		log.Errorf("NetworkId not valid.")
		return nil, err
	}

	/* PackageId */
	pId, err := uuid.FromString(req.PackageId)
	if err != nil {
		log.Errorf("PackageId not valid.")
	}

	pcrfData := &pm.SimInfo{
		Imsi:      sim.Imsi,
		Iccid:     sim.Iccid,
		PackageId: pId,
		NetworkId: nId,
		Visitor:   false, // We will using this flag on roaming in VLR
	}

	/* Send message to PCRF */
	policy, err := s.pc.NewPolicy(pcrfData.PackageId)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "error creating policy")
	}

	/* Add to ASR */
	asr := &db.Asr{
		Iccid:                   req.Iccid,
		Imsi:                    sim.Imsi,
		Op:                      sim.Op,
		Key:                     sim.Key,
		Amf:                     sim.Amf,
		AlgoType:                sim.AlgoType,
		UeDlAmbrBps:             sim.UeDlAmbrBps,
		UeUlAmbrBps:             sim.UeUlAmbrBps,
		Sqn:                     uint64(sim.Sqn),
		CsgIdPrsent:             sim.CsgIdPrsent,
		CsgId:                   sim.CsgId,
		DefaultApnName:          sim.DefaultApnName,
		PackageId:               pId,
		NetworkId:               nId,
		Policy:                  *policy,
		LastStatusChangeAt:      time.Now(),
		AllowedTimeOfService:    s.allowedToS,
		LastStatusChangeReasons: db.ACTIVATION,
	}

	err = s.asrRepo.Add(asr)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "error updating asr")
	}

	err, removed := s.pc.RunPolicyControl(asr.Imsi)
	if err != nil {
		log.Errorf("error running policy control for imsi %s. Error %s", asr.Imsi, err.Error())
		return nil, err
	}

	if removed {
		log.Infof("Profile not added to repo as one or more policies were failed for %s", asr.Imsi)
		return nil, fmt.Errorf("policy failure for profile")
	}

	err = s.pc.SyncProfile(pcrfData, asr, msgbus.ACTION_CRUD_CREATE, "activesubscriber")
	if err != nil {
		return nil, err
	}

	log.Debugf("Activated %s imsi with %+v", asr.Imsi, asr)
	return &pb.ActivateResp{}, err
}

func (s *AsrRecordServer) UpdatePackage(c context.Context, req *pb.UpdatePackageReq) (*pb.UpdatePackageResp, error) {
	asrRecord, err := s.asrRepo.GetByIccid(req.GetIccid())
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "error getting iccid")
	}

	/* We assum that packageId is validated by subscriber. */
	pId, err := uuid.FromString(req.PackageId)
	if err != nil {
		log.Errorf("PackageId not valid.")
		return nil, grpc.SqlErrorToGrpc(err, "error invalid package id")
	}

	pcrfData := &pm.SimInfo{
		ID:        asrRecord.ID,
		Imsi:      asrRecord.Imsi,
		Iccid:     asrRecord.Iccid,
		PackageId: pId,
		NetworkId: asrRecord.NetworkId,
	}

	/* Send message to PCRF */
	policy, err := s.pc.NewPolicy(pcrfData.PackageId)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "error creating policy")
	}

	err = s.asrRepo.UpdatePackage(asrRecord.Imsi, pId, policy)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "error updating asr")
	}

	err, removed := s.pc.RunPolicyControl(asrRecord.Imsi)
	if err != nil {
		log.Errorf("error running policy control for imsi %s. Error %s", asrRecord.Imsi, err.Error())
		return nil, err
	}

	if removed {
		log.Infof("Profile removed from repo as one or more policies were failed for %s", asrRecord.Imsi)
		return nil, fmt.Errorf("policy failure for profile")
	}

	/* read the updated profile */
	nRec, err := s.asrRepo.GetByIccid(req.GetIccid())
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "error getting iccid")
	}

	err = s.pc.SyncProfile(pcrfData, nRec, msgbus.ACTION_CRUD_UPDATE, "activesubscriber")
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "error updating pcrf")
	}

	asrRecord.Policy = *policy
	log.Debugf("Updated policy for %s imsi to %+v", asrRecord.Imsi, nRec)
	return &pb.UpdatePackageResp{}, nil
}

func (s *AsrRecordServer) Inactivate(c context.Context, req *pb.InactivateReq) (*pb.InactivateResp, error) {

	delAsrRecord, err := s.asrRepo.GetByIccid(req.GetIccid())
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "error getting iccid")
	}

	pcrfData := &pm.SimInfo{
		ID:        delAsrRecord.ID,
		Imsi:      delAsrRecord.Imsi,
		Iccid:     delAsrRecord.Iccid,
		NetworkId: delAsrRecord.NetworkId,
	}

	err = s.asrRepo.Delete(delAsrRecord.Imsi, db.DEACTIVATION)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "error updating asr")
	}

	err = s.pc.SyncProfile(pcrfData, delAsrRecord, msgbus.ACTION_CRUD_DELETE, "activesubscriber")
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "error updating pcrf")
	}

	log.Debugf("Deleted subscriber %+v", delAsrRecord)

	return &pb.InactivateResp{}, nil

}

func (s *AsrRecordServer) UpdateGuti(c context.Context, req *pb.UpdateGutiReq) (*pb.UpdateGutiResp, error) {
	_, err := s.asrRepo.GetByImsi(req.Imsi)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "imsi")
	}

	err = s.gutiRepo.Update(&db.Guti{
		Imsi:            req.Imsi,
		PlmnId:          req.Guti.PlmnId,
		Mmegi:           req.Guti.Mmegi,
		Mmec:            req.Guti.Mmec,
		MTmsi:           req.Guti.Mtmsi,
		DeviceUpdatedAt: time.Unix(int64(req.UpdatedAt), 0),
	})
	if err != nil {
		log.Errorf("Failed to update GUTI: %s", err.Error())
		if err.Error() == db.GutiNotUpdatedErr {
			return nil, status.Errorf(codes.AlreadyExists, err.Error())
		}
		return nil, grpc.SqlErrorToGrpc(err, "guti")
	}

	return &pb.UpdateGutiResp{}, nil
}

func (s *AsrRecordServer) UpdateTai(c context.Context, req *pb.UpdateTaiReq) (*pb.UpdateTaiResp, error) {
	_, err := s.asrRepo.GetByImsi(req.Imsi)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "imsi")
	}

	err = s.asrRepo.UpdateTai(req.Imsi, db.Tai{
		PlmnId:          req.PlmnId,
		Tac:             req.Tac,
		DeviceUpdatedAt: time.Unix(int64(req.UpdatedAt), 0),
	})

	if err != nil {
		if err.Error() == db.TaiNotUpdatedErr {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, grpc.SqlErrorToGrpc(err, "tai")
	}

	return &pb.UpdateTaiResp{}, nil
}

func (s *AsrRecordServer) UpdateandSyncAsrProfile(imsi string) error {

	sub, err := s.asrRepo.GetByImsi(imsi)
	if err != nil {
		return grpc.SqlErrorToGrpc(err, "error getting imsi")
	}

	r, err := s.cdr.GetUsage(imsi)
	if err != nil {
		log.Errorf("Failed to get usage: %v for imsi %s", err, imsi)
		return err
	}

	sub.Policy.ConsumedData = r.Usage

	err = s.asrRepo.Update(imsi, sub)
	if err != nil {
		log.Errorf("Failed to update usage: %v for imsi %s.Error%s", r, imsi, err.Error())
		return err
	}

	err, removed := s.pc.RunPolicyControl(imsi)
	if err != nil {
		log.Errorf("error running policy control for imsi %s. Error %s", sub.Imsi, err.Error())
		return err
	}

	if removed {
		log.Infof("Profile removed from repo as one or more policies were failed for %s", sub.Imsi)
		return fmt.Errorf("policy failure for profile")
	}

	pcrfData := &pm.SimInfo{
		ID:        sub.ID,
		Imsi:      sub.Imsi,
		Iccid:     sub.Iccid,
		PackageId: sub.PackageId,
		NetworkId: sub.NetworkId,
	}

	err = s.pc.SyncProfile(pcrfData, sub, msgbus.ACTION_CRUD_UPDATE, "activesubscriber")
	if err != nil {
		return err
	}

	return nil
}
