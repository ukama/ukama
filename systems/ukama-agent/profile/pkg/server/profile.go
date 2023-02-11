package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
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
	"github.com/ukama/ukama/systems/ukama-agent/profile/pkg/policy"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type ProfileServer struct {
	pb.UnimplementedProfileServiceServer
	profileRepo db.ProfileRepo

	msgbus           mb.MsgBusServiceClient
	baseRoutingKey   msgbus.RoutingKeyBuilder
	Org              string
	PolicyController *policy.PolicyController
	nodePolicyPath   string
}

func NewProfileServer(pRepo db.ProfileRepo, org string, msgBus mb.MsgBusServiceClient, nodePath string, period time.Duration) (*ProfileServer, error) {

	ps := &ProfileServer{
		profileRepo:    pRepo,
		Org:            org,
		msgbus:         msgBus,
		nodePolicyPath: nodePath,
	}

	if msgBus != nil {
		ps.baseRoutingKey = msgbus.NewRoutingKeyBuilder().SetCloudSource().SetContainer(pkg.ServiceName)
	}

	ps.PolicyController = policy.NewPolicyController(pRepo, org, msgBus, nodePath, period)

	return ps, nil
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
		AllowedTimeOfService: p.AllowedTimeOfService,
		ConsumedDataBytes:    p.ConsumedDataBytes,
		UpdatedAt:            int64(p.Model.UpdatedAt.Unix()),
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
		Iccid:                   req.Profile.Iccid,
		Imsi:                    req.Profile.Imsi,
		PackageId:               pId,
		NetworkId:               nId,
		UeDlBps:                 req.Profile.UeDlBps,
		UeUlBps:                 req.Profile.UeUlBps,
		ApnName:                 req.Profile.Apn.Name,
		AllowedTimeOfService:    req.Profile.AllowedTimeOfService,
		ConsumedDataBytes:       req.Profile.ConsumedDataBytes,
		TotalDataBytes:          req.Profile.TotalDataBytes,
		LastStatusChangeReasons: db.ACTIVATION,
		LastStatusChangeAt:      time.Now(),
	}

	err = s.profileRepo.Add(p)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "error updating profile")
	}

	err, pState := s.PolicyController.RunPolicyControl(p.Imsi)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "error checking policies")
	}

	if pState {
		logrus.Errorf("Policy controller rejceted profile.")
		return nil, fmt.Errorf("policy control rejected profile")
	}

	/* Create event */
	e := &epb.ProfileUpdated{
		Profile: &epb.Profile{
			Imsi:                 p.Imsi,
			Iccid:                p.Iccid,
			Network:              p.NetworkId.String(),
			Package:              p.PackageId.String(),
			Org:                  s.Org,
			AllowedTimeOfService: p.AllowedTimeOfService,
			TotalDataBytes:       p.TotalDataBytes,
			LastStatusChangeAt:   p.LastStatusChangeAt.Unix(),
		},
	}

	_ = s.publishEvent(msgbus.ACTION_CRUD_CREATE, "profile", e)

	s.syncProfile(http.MethodPut, e.Profile.Iccid)

	return &pb.AddResp{}, err
}

func (s *ProfileServer) UpdatePackage(c context.Context, req *pb.UpdatePackageReq) (*pb.UpdatePackageResp, error) {
	p, err := s.profileRepo.GetByIccid(req.GetIccid())
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "error getting iccid")
	}

	/* We assume that packageId is validated by subscriber. */
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
		LastStatusChangeAt:   time.Now(),
	}

	err = s.profileRepo.UpdatePackage(p.Imsi, pack)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "error updating asr")
	}

	err, pState := s.PolicyController.RunPolicyControl(p.Imsi)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "error checking policies")
	}

	if pState {
		logrus.Errorf("Policy controller rejceted profile.")
		return nil, fmt.Errorf("policy control rejected profile")
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

	_ = s.publishEvent(msgbus.ACTION_CRUD_UPDATE, "profile", e)

	s.syncProfile(http.MethodPut, e.Profile.Iccid)

	return &pb.UpdatePackageResp{}, nil
}

func (s *ProfileServer) UpdateUsage(c context.Context, req *pb.UpdateUsageReq) (*pb.UpdateUsageResp, error) {
	p, err := s.profileRepo.GetByImsi(req.GetImsi())
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "error getting imsi")
	}

	err = s.profileRepo.UpdateUsage(p.Imsi, req.ConsumedDataBytes)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "error updating asr")
	}

	/* Create event */
	e := &epb.ProfileUpdated{
		Profile: &epb.Profile{
			Imsi:                 p.Imsi,
			Iccid:                p.Iccid,
			Network:              p.NetworkId.String(),
			Package:              p.PackageId.String(),
			Org:                  s.Org,
			AllowedTimeOfService: p.AllowedTimeOfService,
			TotalDataBytes:       p.TotalDataBytes,
		},
	}

	/* Check for the data */
	if p.TotalDataBytes <= req.ConsumedDataBytes {
		err = s.profileRepo.Delete(p.Imsi, db.NO_DATA_AVAILABLE)
		if err != nil {
			return nil, grpc.SqlErrorToGrpc(err, "error updating asr")
		}

		_ = s.publishEvent(msgbus.ACTION_CRUD_DELETE, "profile", e)

		s.syncProfile(http.MethodDelete, e.Profile.Iccid)
	}

	return &pb.UpdateUsageResp{}, nil
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

	err = s.profileRepo.Delete(delProfile.Imsi, db.DEACTIVATION)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "error updating profile db")
	}

	/* Create event */
	e := &epb.ProfileRemoved{
		Profile: &epb.Profile{
			Imsi:                 delProfile.Imsi,
			Iccid:                delProfile.Iccid,
			Network:              delProfile.NetworkId.String(),
			Package:              delProfile.PackageId.String(),
			Org:                  s.Org,
			AllowedTimeOfService: delProfile.AllowedTimeOfService,
			TotalDataBytes:       delProfile.TotalDataBytes,
		},
	}

	_ = s.publishEvent(msgbus.ACTION_CRUD_DELETE, "profile", e)

	s.syncProfile(http.MethodDelete, e.Profile.Iccid)

	return &pb.RemoveResp{}, nil

}

func (s *ProfileServer) Sync(c context.Context, req *pb.SyncReq) (*pb.SyncResp, error) {

	for _, iccid := range req.Iccid {
		s.syncProfile(http.MethodPut, iccid)
	}

	return &pb.SyncResp{}, nil
}

func (s *ProfileServer) syncProfile(method string, iccid string) {
	p, err := s.Read(context.Background(), &pb.ReadReq{
		Id: &pb.ReadReq_Iccid{
			Iccid: iccid,
		},
	})
	if err != nil {
		logrus.Errorf("error syncing %s: %s", iccid, err.Error())
		return
	}

	body, err := json.Marshal(p.Profile)
	if err != nil {
		logrus.Errorf("error marshaling profile: %s", err.Error())
		return
	}

	if s.msgbus != nil {
		route := s.baseRoutingKey.SetAction("node-feed").SetObject("profile").MustBuild()
		err = s.msgbus.PublishToNodeFeeder(route, s.Org, "*", s.nodePolicyPath, method, body)
		if err != nil {
			logrus.Errorf("Failed to publish message %+v with key %+v. Errors %s", body, route, err.Error())
		}
	}

}

func (s *ProfileServer) publishEvent(action string, object string, msg protoreflect.ProtoMessage) error {
	var err error
	if s.msgbus != nil {
		route := s.baseRoutingKey.SetAction(action).SetObject(object).MustBuild()
		err = s.msgbus.PublishRequest(route, msg)
		if err != nil {
			logrus.Errorf("Failed to publish message %+v with key %+v. Errors %s", msg, route, err.Error())
		}
	}

	return err
}
