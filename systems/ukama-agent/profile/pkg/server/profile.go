package server

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/grpc"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	uuid "github.com/ukama/ukama/systems/common/uuid"
	pb "github.com/ukama/ukama/systems/ukama-agent/profile/pb/gen"
	"github.com/ukama/ukama/systems/ukama-agent/profile/pkg"
	"github.com/ukama/ukama/systems/ukama-agent/profile/pkg/db"
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
	var p *db.Profile
	var err error

	switch req.Id.(type) {
	case *pb.ReadReq_Imsi:

		p, err = s.profileRepo.GetByImsi(req.GetImsi())
		if err != nil {
			return nil, grpc.SqlErrorToGrpc(err, "error getting imsi")
		}

	case *pb.ReadReq_Iccid:
		p, err = s.profileRepo.GetByIccid(req.GetIccid())
		if err != nil {
			return nil, grpc.SqlErrorToGrpc(err, "error getting iccid")
		}
	}

	resp := &pb.ReadResp{Profile: &pb.Profile{
		Imsi:    p.Imsi,
		Iccid:   p.Iccid,
		UeDlBps: p.UeDlBps,
		UeUlBps: p.UeUlBps,
		Apn: &pb.Apn{
			Name: p.ApnName,
		},
		NetworkId:            p.NetworkId.String(),
		PackageId:            p.PackageId.String(),
		AllowedTimeOfService: uint64(p.AllowedTimeOfService.Seconds()),
		ConsumedDataBytes:    p.ConsumedDataBytes,
		UpdatedAt:            uint64(p.Model.UpdatedAt.Unix()),
	}}

	logrus.Infof("Subscriber is having %+v", resp)
	return resp, nil
}

func (s *ProfileServer) Add(c context.Context, req *pb.AddReq) (*pb.AddResp, error) {

	/* Send message to PCRF */
	nId, err := uuid.FromString(req.Profile.NetworkId)
	if err != nil {
		logrus.Errorf("NetworkId not valid.")
		return nil, err
	}

	pId, err := uuid.FromString(req.Profile.PackageId)
	if err != nil {
		logrus.Errorf("PackageId not valid.")
	}

	/* Add to Profile */
	p := &db.Profile{
		Iccid:                req.Profile.Iccid,
		Imsi:                 req.Profile.Imsi,
		PackageId:            pId,
		NetworkId:            nId,
		UeDlBps:              req.Profile.UeDlBps,
		UeUlBps:              req.Profile.UeUlBps,
		ApnName:              req.Profile.Apn.Name,
		AllowedTimeOfService: time.Duration(req.Profile.AllowedTimeOfService) * time.Second,
		ConsumedDataBytes:    req.Profile.ConsumedDataBytes,
		TotalDataBytes:       req.Profile.TotalDataBytes,
	}

	err = s.profileRepo.Add(p)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "error updating profile")
	}

	/* Create event */
	e := &epb.ProfileUpdated{
		Profile: &epb.Profile{
			Imsi:                 p.Imsi,
			Iccid:                p.Iccid,
			Network:              p.NetworkId.String(),
			Package:              p.PackageId.String(),
			Org:                  s.Org,
			AllowedTimeOfService: uint64(p.AllowedTimeOfService.Seconds()),
			TotalDataBytes:       p.TotalDataBytes,
		},
	}

	if s.msgbus != nil {
		route := s.baseRoutingKey.SetAction("create").SetObject("profile").MustBuild()
		merr := s.msgbus.PublishRequest(route, e)
		if merr != nil {
			logrus.Errorf("Failed to publish message %+v with key %+v. Errors %s", e, route, err.Error())
		}
	}

	return &pb.AddResp{}, err
}

func (s *ProfileServer) UpdatePackage(c context.Context, req *pb.UpdatePackageReq) (*pb.UpdatePackageResp, error) {
	p, err := s.profileRepo.GetByIccid(req.GetIccid())
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "error getting iccid")
	}

	/* We assum that packageId is validated by subscriber. */
	pId, err := uuid.FromString(req.Package.PackageId)
	if err != nil {
		logrus.Errorf("PackageId not valid.")
		return nil, grpc.SqlErrorToGrpc(err, "error invalid package id")
	}

	pack := db.PackageDetails{
		PackageId:            pId,
		AllowedTimeOfService: time.Duration(req.Package.AllowedTimeOfService) * time.Second,
		TotalDataBytes:       req.Package.TotalDataBytes,
		ConsumedDataBytes:    req.Package.ConsumedDataBytes,
		UeDlBps:              req.Package.UeDlBps,
		UeUlBps:              req.Package.UeUlBps,
		ApnName:              req.Package.Apn.Name,
	}

	err = s.profileRepo.UpdatePackage(p.Imsi, pack)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "error updating asr")
	}

	/* Create event */
	e := &epb.ProfileUpdated{
		Profile: &epb.Profile{
			Imsi:                 p.Imsi,
			Iccid:                p.Iccid,
			Network:              p.NetworkId.String(),
			Package:              req.Package.PackageId,
			Org:                  s.Org,
			AllowedTimeOfService: req.Package.AllowedTimeOfService,
			TotalDataBytes:       p.TotalDataBytes,
		},
	}

	if s.msgbus != nil {
		route := s.baseRoutingKey.SetActionUpdate().SetObject("profile").MustBuild()
		merr := s.msgbus.PublishRequest(route, e)
		if merr != nil {
			logrus.Errorf("Failed to publish message %+v with key %+v. Errors %s", e, route, err.Error())
		}
	}
	return &pb.UpdatePackageResp{}, nil
}

func (s *ProfileServer) Remove(c context.Context, req *pb.RemoveReq) (*pb.RemoveResp, error) {
	var delProfile *db.Profile
	var err error

	switch req.Id.(type) {
	case *pb.RemoveReq_Imsi:

		delProfile, err = s.profileRepo.GetByImsi(req.GetImsi())
		if err != nil {
			return nil, grpc.SqlErrorToGrpc(err, "error getting imsi")
		}

	case *pb.RemoveReq_Iccid:
		delProfile, err = s.profileRepo.GetByIccid(req.GetIccid())
		if err != nil {
			return nil, grpc.SqlErrorToGrpc(err, "error getting iccid")
		}
	}

	err = s.profileRepo.Delete(delProfile.Imsi)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "error updating asr")
	}

	/* Create event */
	e := &epb.ProfileRemoved{
		Profile: &epb.Profile{
			Imsi:                 delProfile.Imsi,
			Iccid:                delProfile.Iccid,
			Network:              delProfile.NetworkId.String(),
			Package:              delProfile.PackageId.String(),
			Org:                  s.Org,
			AllowedTimeOfService: uint64(delProfile.AllowedTimeOfService.Seconds()),
			TotalDataBytes:       delProfile.TotalDataBytes,
		},
	}

	if s.msgbus != nil {
		route := s.baseRoutingKey.SetActionDelete().SetObject("profile").MustBuild()
		merr := s.msgbus.PublishRequest(route, e)
		if merr != nil {
			logrus.Errorf("Failed to publish message %+v with key %+v. Errors %s", e, route, err.Error())
		}
	}

	return &pb.RemoveResp{}, nil

}

func (s *ProfileServer) Sync(c context.Context, req *pb.SyncReq) (*pb.SyncResp, error) {

	return &pb.SyncResp{}, nil
}
