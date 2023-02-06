package server

import (
	"context"
	"fmt"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/grpc"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	pb "github.com/ukama/ukama/systems/ukama-agent/profile/pb/gen"
	"github.com/ukama/ukama/systems/ukama-agent/profile/pkg"
	"github.com/ukama/ukama/systems/ukama-agent/profile/pkg/client"
	"github.com/ukama/ukama/systems/ukama-agent/profile/pkg/db"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProfileServer struct {
	pb.UnimplementedProfileServiceServer
	profileRepo db.ProfileRepo

	msgbus         mb.MsgBusServiceClient
	baseRoutingKey msgbus.RoutingKeyBuilder
	Org            string
}

func NewProfileServer(pRepo db.ProfileRepo, org string, msgBus mb.MsgBusServiceClient) (*ProfileServer, error) {

	asr := ProfileServer{
		profileRepo: pRepo,
		Org:         org,
		msgbus:      msgBus,
	}

	if msgBus != nil {
		asr.baseRoutingKey = msgbus.NewRoutingKeyBuilder().SetCloudSource().SetContainer(pkg.ServiceName)
	}

	return &asr, nil
}

func (s *ProfileServer) Read(c context.Context, req *pb.ReadReq) (*pb.ReadResp, error) {
	var sub *db.Profile
	var err error

	switch req.Id.(type) {
	case *pb.ReadReq_Imsi:

		sub, err = s.profileRepo.GetByImsi(req.GetImsi())
		if err != nil {
			return nil, grpc.SqlErrorToGrpc(err, "error getting imsi")
		}

	case *pb.ReadReq_Iccid:
		sub, err = s.profileRepo.GetByIccid(req.GetIccid())
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

	logrus.Infof("Subscriber is having %+v", resp)
	return resp, nil
}

func (s *ProfileServer) Activate(c context.Context, req *pb.ActivateReq) (*pb.ActivateResp, error) {

	/* Validate network in Org */
	err := s.network.ValidateNetwork(req.Network, s.Org)
	if err != nil {
		return nil, fmt.Errorf("error validating network")
	}

	/* Send Request to SIM Factory */
	sim, err := s.factory.ReadSimCardInfo(req.Iccid)
	if err != nil {
		return nil, fmt.Errorf("error reading iccid from factory")
	}

	/* Send message to PCRF */
	nId, err := uuid.FromString(req.Network)
	if err != nil {
		logrus.Errorf("NetworkId not valid.")
		return nil, err
	}

	pId, err := uuid.FromString(req.PackageId)
	if err != nil {
		logrus.Errorf("PackageId not valid.")
	}

	pcrfData := client.PolicyControlSimInfo{
		Imsi:      sim.Imsi,
		Iccid:     sim.Iccid,
		PackageId: pId,
		NetworkId: nId,
		Visitor:   false, // We will using this flag on roaming in VLR
	}

	err = s.pcrf.AddSim(pcrfData)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "error adding to pcrf")
	}

	/* Add to ASR */
	asr := &db.Profile{
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
	}

	err = s.profileRepo.Add(asr)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "error updating asr")
	}

	/* Create event */
	e := &epb.ProfileActivated{
		Subscriber: &epb.Subscriber{
			Imsi:    asr.Imsi,
			Iccid:   asr.Iccid,
			Network: asr.NetworkID.String(),
			Package: asr.PackageId.String(),
			Org:     s.Org,
		},
	}

	if s.msgbus != nil {
		route := s.baseRoutingKey.SetAction("create").SetObject("activesubscriber").MustBuild()
		merr := s.msgbus.PublishRequest(route, e)
		if merr != nil {
			logrus.Errorf("Failed to publish message %+v with key %+v. Errors %s", e, route, err.Error())
		}
	}

	return &pb.ActivateResp{}, err
}

func (s *ProfileServer) UpdatePackage(c context.Context, req *pb.UpdatePackageReq) (*pb.UpdatePackageResp, error) {
	asrRecord, err := s.profileRepo.GetByIccid(req.GetIccid())
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "error getting iccid")
	}

	/* We assum that packageId is validated by subscriber. */
	pId, err := uuid.FromString(req.PackageId)
	if err != nil {
		logrus.Errorf("PackageId not valid.")
		return nil, grpc.SqlErrorToGrpc(err, "error invalid package id")
	}

	pD := client.PolicyControlSimPackageUpdate{
		Imsi:      asrRecord.Imsi,
		PackageId: pId,
	}

	err = s.pcrf.UpdateSim(pD)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "error updating pcrf")
	}

	err = s.profileRepo.UpdatePackage(asrRecord.Imsi, pId)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "error updating asr")
	}

	/* Create event */
	e := &epb.ProfileUpdated{
		Subscriber: &epb.Subscriber{
			Imsi:    asrRecord.Imsi,
			Iccid:   asrRecord.Iccid,
			Network: asrRecord.NetworkID.String(),
			Package: req.PackageId,
			Org:     s.Org,
		},
	}

	if s.msgbus != nil {
		route := s.baseRoutingKey.SetActionUpdate().SetObject("updateactivesubscriber").MustBuild()
		merr := s.msgbus.PublishRequest(route, e)
		if merr != nil {
			logrus.Errorf("Failed to publish message %+v with key %+v. Errors %s", e, route, err.Error())
		}
	}
	return &pb.UpdatePackageResp{}, nil
}

func (s *ProfileServer) Inactivate(c context.Context, req *pb.InactivateReq) (*pb.InactivateResp, error) {
	var delProfile *db.Profile
	var err error

	switch req.Id.(type) {
	case *pb.InactivateReq_Imsi:

		delProfile, err = s.profileRepo.GetByImsi(req.GetImsi())
		if err != nil {
			return nil, grpc.SqlErrorToGrpc(err, "error getting imsi")
		}

	case *pb.InactivateReq_Iccid:
		delProfile, err = s.profileRepo.GetByIccid(req.GetIccid())
		if err != nil {
			return nil, grpc.SqlErrorToGrpc(err, "error getting iccid")
		}
	}

	err = s.pcrf.DeleteSim(delProfile.Imsi)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "error updating pcrf")
	}

	err = s.profileRepo.Delete(delProfile.Imsi)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "error updating asr")
	}

	/* Create event */
	e := &epb.ProfileInactivated{
		Subscriber: &epb.Subscriber{
			Imsi:    delProfile.Imsi,
			Iccid:   delProfile.Iccid,
			Network: delProfile.NetworkID.String(),
			Package: delProfile.PackageId.String(),
			Org:     s.Org,
		},
	}

	if s.msgbus != nil {
		route := s.baseRoutingKey.SetActionDelete().SetObject("activesubscriber").MustBuild()
		merr := s.msgbus.PublishRequest(route, e)
		if merr != nil {
			logrus.Errorf("Failed to publish message %+v with key %+v. Errors %s", e, route, err.Error())
		}
	}

	return &pb.InactivateResp{}, nil

}

func (s *ProfileServer) UpdateGuti(c context.Context, req *pb.UpdateGutiReq) (*pb.UpdateGutiResp, error) {
	_, err := s.profileRepo.GetByImsi(req.Imsi)
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
		logrus.Errorf("Failed to update GUTI: %s", err.Error())
		if err.Error() == db.GutiNotUpdatedErr {
			return nil, status.Errorf(codes.AlreadyExists, err.Error())
		}
		return nil, grpc.SqlErrorToGrpc(err, "guti")
	}

	return &pb.UpdateGutiResp{}, nil
}

func (s *ProfileServer) UpdateTai(c context.Context, req *pb.UpdateTaiReq) (*pb.UpdateTaiResp, error) {
	_, err := s.profileRepo.GetByImsi(req.Imsi)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "imsi")
	}

	err = s.profileRepo.UpdateTai(req.Imsi, db.Tai{
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
