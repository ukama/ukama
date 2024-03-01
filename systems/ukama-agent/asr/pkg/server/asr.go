package server

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/grpc"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	"github.com/ukama/ukama/systems/common/uuid"
	pb "github.com/ukama/ukama/systems/ukama-agent/asr/pb/gen"
	"github.com/ukama/ukama/systems/ukama-agent/asr/pkg"
	"github.com/ukama/ukama/systems/ukama-agent/asr/pkg/client"
	"github.com/ukama/ukama/systems/ukama-agent/asr/pkg/db"
	"github.com/ukama/ukama/systems/ukama-agent/asr/pkg/pcrf"

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
	msgbus         mb.MsgBusServiceClient
	baseRoutingKey msgbus.RoutingKeyBuilder
	pcrf           pcrf.PCRFController
	OrgName        string
	OrgId          string
}

func NewAsrRecordServer(asrRepo db.AsrRecordRepo, gutiRepo db.GutiRepo, pRepo db.PolicyRepo, factory client.Factory, network client.Network, pcrf pcrf.PCRFController, orgId, orgName string, msgBus mb.MsgBusServiceClient) (*AsrRecordServer, error) {

	asr := AsrRecordServer{
		asrRepo:  asrRepo,
		gutiRepo: gutiRepo,
		OrgName:  orgName,
		OrgId:    orgId,
		factory:  factory,
		network:  network,
		pRepo:    pRepo,
		msgbus:   msgBus,
		pcrf:     pcrf,
	}

	if msgBus != nil {
		asr.baseRoutingKey = msgbus.NewRoutingKeyBuilder().SetRequestType().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName)
	}

	log.Infof("Asr is %+v", asr)

	return &asr, nil
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
	}}

	log.Infof("Subscriber is having %+v", resp)
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

	pcrfData := &pcrf.SimInfo{
		Imsi:      sim.Imsi,
		Iccid:     sim.Iccid,
		PackageId: pId,
		NetworkId: nId,
		Visitor:   false, // We will using this flag on roaming in VLR
	}

	/* Send message to PCRF */
	policy, err := s.pcrf.NewPolicy(pcrfData.PackageId)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "error creating policy")
	}

	/* Add to ASR */
	asr := &db.Asr{
		Iccid:          req.Iccid,
		Imsi:           sim.Imsi,
		Op:             sim.Op,
		Key:            sim.Key,
		Amf:            sim.Amf,
		AlgoType:       sim.AlgoType,
		UeDlAmbrBps:    sim.UeDlAmbrBps,
		UeUlAmbrBps:    sim.UeUlAmbrBps,
		Sqn:            uint64(sim.Sqn),
		CsgIdPrsent:    sim.CsgIdPrsent,
		CsgId:          sim.CsgId,
		DefaultApnName: sim.DefaultApnName,
		PackageId:      pId,
		NetworkID:      nId,
		Policy:         *policy,
	}

	err = s.asrRepo.Add(asr)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "error updating asr")
	}

	err = s.pcrf.AddPolicy(pcrfData, policy)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "error adding policy")
	}

	/* Create event */
	e := &epb.AsrActivated{
		Subscriber: &epb.Subscriber{
			Imsi:    asr.Imsi,
			Iccid:   asr.Iccid,
			Network: asr.NetworkID.String(),
			Package: asr.PackageId.String(),
			Org:     s.OrgId,
			Policy:  policy.Id.String(),
		},
	}

	if s.msgbus != nil {
		route := s.baseRoutingKey.SetAction("create").SetObject("activesubscriber").MustBuild()
		merr := s.msgbus.PublishRequest(route, e)
		if merr != nil {
			log.Errorf("Failed to publish message %+v with key %+v. Errors %s", e, route, err.Error())
		}
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

	pcrfData := &pcrf.SimInfo{
		ID:        asrRecord.ID,
		Imsi:      asrRecord.Imsi,
		Iccid:     asrRecord.Iccid,
		PackageId: pId,
		NetworkId: asrRecord.NetworkID,
	}

	/* Send message to PCRF */
	policy, err := s.pcrf.NewPolicy(pcrfData.PackageId)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "error creating policy")
	}

	err = s.asrRepo.UpdatePackage(asrRecord.Imsi, pId, policy)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "error updating asr")
	}

	err = s.pcrf.UpdatePolicy(pcrfData, policy)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "error updating pcrf")
	}

	/* Create event */
	e := &epb.AsrUpdated{
		Subscriber: &epb.Subscriber{
			Imsi:    asrRecord.Imsi,
			Iccid:   asrRecord.Iccid,
			Network: asrRecord.NetworkID.String(),
			Package: req.PackageId,
			Org:     s.OrgId,
			Policy:  policy.Id.String(),
		},
	}

	if s.msgbus != nil {
		route := s.baseRoutingKey.SetActionUpdate().SetObject("updateactivesubscriber").MustBuild()
		merr := s.msgbus.PublishRequest(route, e)
		if merr != nil {
			log.Errorf("Failed to publish message %+v with key %+v. Errors %s", e, route, err.Error())
		}
	}

	asrRecord.Policy = *policy
	log.Debugf("Updated policy for %s imsi to %+v", asrRecord.Imsi, asrRecord)
	return &pb.UpdatePackageResp{}, nil
}

func (s *AsrRecordServer) Inactivate(c context.Context, req *pb.InactivateReq) (*pb.InactivateResp, error) {

	delAsrRecord, err := s.asrRepo.GetByIccid(req.GetIccid())
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "error getting iccid")
	}

	pcrfData := &pcrf.SimInfo{
		ID:        delAsrRecord.ID,
		Imsi:      delAsrRecord.Imsi,
		Iccid:     delAsrRecord.Iccid,
		NetworkId: delAsrRecord.NetworkID,
	}

	err = s.pcrf.DeletePolicy(pcrfData)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "error updating pcrf")
	}

	err = s.asrRepo.Delete(delAsrRecord.Imsi)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "error updating asr")
	}

	/* Create event */
	e := &epb.AsrInactivated{
		Subscriber: &epb.Subscriber{
			Imsi:    delAsrRecord.Imsi,
			Iccid:   delAsrRecord.Iccid,
			Network: delAsrRecord.NetworkID.String(),
			Package: delAsrRecord.PackageId.String(),
			Org:     s.OrgId,
		},
	}

	if s.msgbus != nil {
		route := s.baseRoutingKey.SetActionDelete().SetObject("activesubscriber").MustBuild()
		merr := s.msgbus.PublishRequest(route, e)
		if merr != nil {
			log.Errorf("Failed to publish message %+v with key %+v. Errors %s", e, route, err.Error())
		}
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
