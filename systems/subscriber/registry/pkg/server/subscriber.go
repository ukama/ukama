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
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/grpc"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	upb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
	cnucl "github.com/ukama/ukama/systems/common/rest/client/nucleus"
	creg "github.com/ukama/ukama/systems/common/rest/client/registry"
	uuid "github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/common/validation"
	validate "github.com/ukama/ukama/systems/common/validation"
	pb "github.com/ukama/ukama/systems/subscriber/registry/pb/gen"
	"github.com/ukama/ukama/systems/subscriber/registry/pkg"
	"github.com/ukama/ukama/systems/subscriber/registry/pkg/client"
	"github.com/ukama/ukama/systems/subscriber/registry/pkg/db"
	simMangerPb "github.com/ukama/ukama/systems/subscriber/sim-manager/pb/gen"
)

type SubcriberServer struct {
	orgName              string
	orgId                string
	msgbus               mb.MsgBusServiceClient
	subscriberRepo       db.SubscriberRepo
	subscriberRoutingKey msgbus.RoutingKeyBuilder
	simManagerService    client.SimManagerClientProvider
	orgClient            cnucl.OrgClient
	networkClient        creg.NetworkClient
	pb.UnimplementedRegistryServiceServer
}

func NewSubscriberServer(orgName string, subscriberRepo db.SubscriberRepo, msgBus mb.MsgBusServiceClient, simManagerService client.SimManagerClientProvider, orgId string, orgService cnucl.OrgClient, networkClient creg.NetworkClient) *SubcriberServer {
	return &SubcriberServer{
		orgName:              orgName,
		subscriberRepo:       subscriberRepo,
		msgbus:               msgBus,
		simManagerService:    simManagerService,
		subscriberRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
		orgId:                orgId,
		orgClient:            orgService,
		networkClient:        networkClient,
	}
}

func (s *SubcriberServer) Add(ctx context.Context, req *pb.AddSubscriberRequest) (*pb.AddSubscriberResponse, error) {
	log.Infof("Adding subscriber: %v", req)

	var dob string
	var err error

	if req.GetDob() != "" {
		dob, err = validate.ValidateDate(req.GetDob())
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}
	}
	// remoteOrg, err := s.orgClient.Get(s.orgName)
	// if err != nil {
	// 	return nil, err
	// }

	var networkInfo *creg.NetworkInfo

	if req.GetNetworkId() != "" {
		networkInfo, err = s.networkClient.Get(req.GetNetworkId())
		if err != nil {
			return nil, status.Errorf(codes.NotFound, "network not found: %s", err.Error())
		}
	} else {
		networkInfo, err = s.networkClient.GetDefault()
		if err != nil {
			return nil, status.Errorf(codes.NotFound, "default network not found: %s", err.Error())
		}
		log.Infof("Default network %+v", networkInfo)
	}

	// if s.orgId != networkInfo.OrgId {
	// 	log.Error("Missing network.")

	// 	return nil, fmt.Errorf("Network mismatch")
	// }

	// if remoteOrg.IsDeactivated {
	// 	return nil, status.Errorf(codes.FailedPrecondition,
	// 		"org is deactivated: cannot add network to it")
	// }

	nid, err := uuid.FromString(networkInfo.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of network uuid. Error %s", err.Error())
	}

	subscriber := &db.Subscriber{
		FirstName:             req.GetFirstName(),
		LastName:              req.GetLastName(),
		NetworkId:             nid,
		Email:                 strings.ToLower(req.GetEmail()),
		PhoneNumber:           req.GetPhoneNumber(),
		Gender:                req.GetGender(),
		Address:               req.GetAddress(),
		ProofOfIdentification: req.GetProofOfIdentification(),
		DOB:                   dob,
		IdSerial:              req.GetIdSerial(),
	}

	err = s.subscriberRepo.Add(subscriber, func(*db.Subscriber, *gorm.DB) error {
		subscriber.SubscriberId = uuid.NewV4()

		return nil
	})
	if err != nil {
		log.Error("error while adding subscriber" + err.Error())

		return nil, grpc.SqlErrorToGrpc(err, "subscriber")
	}

	subscriberPb := dbSubscriberToPbSubscriber(subscriber, nil)
	route := s.subscriberRoutingKey.SetAction("create").SetObject("subscriber").MustBuild()
	_ = s.PublishEventMessage(route, &epb.AddSubscriber{
		Subscriber: subscriberPb,
	})

	return &pb.AddSubscriberResponse{
		Subscriber: subscriberPb,
	}, nil
}

func (s *SubcriberServer) Get(ctx context.Context, req *pb.GetSubscriberRequest) (*pb.GetSubscriberResponse, error) {
	log.Infof("Getting subscriber with ID: %v", req)

	subscriberIdReq := req.GetSubscriberId()

	subscriberId, err := uuid.FromString(subscriberIdReq)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid subscriberId format: %v", err.Error())
	}

	subscriber, err := s.subscriberRepo.Get(subscriberId)
	if err != nil {
		log.Errorf("Error while getting subscriber: %s", err.Error())

		return nil, grpc.SqlErrorToGrpc(err, "subscriber")
	}

	smc, err := s.simManagerService.GetSimManagerService()
	if err != nil {
		log.Errorf("Error while calling SimManagerServiceClient: %s", err.Error())

		return nil, err
	}

	simRep, err := smc.GetSimsBySubscriber(ctx, &simMangerPb.GetSimsBySubscriberRequest{
		SubscriberId: subscriberId.String()})
	if err != nil {
		log.Errorf("Error while getting Sims by subscriber: %s", err.Error())

		return nil, err
	}

	resp := &pb.GetSubscriberResponse{
		Subscriber: dbSubscriberToPbSubscriber(subscriber, pbManagerSimsToPbSubscriberSims(simRep.Sims)),
	}

	return resp, nil
}

func (s *SubcriberServer) GetByEmail(ctx context.Context, req *pb.GetSubscriberByEmailRequest) (*pb.GetSubscriberByEmailResponse, error) {
	log.Infof("Getting subscriber with email: %v", req)

	subscriberEmailReq := req.GetEmail()

	subscriber, err := s.subscriberRepo.GetByEmail(subscriberEmailReq)
	if err != nil {
		log.Errorf("Error while getting subscriber: %s", err.Error())

		return nil, grpc.SqlErrorToGrpc(err, "subscriber")
	}

	smc, err := s.simManagerService.GetSimManagerService()
	if err != nil {
		log.Errorf("Error while calling SimManagerServiceClient: %s", err.Error())

		return nil, err
	}

	simRep, err := smc.GetSimsBySubscriber(ctx, &simMangerPb.GetSimsBySubscriberRequest{
		SubscriberId: subscriber.SubscriberId.String()})
	if err != nil {
		log.Errorf("Error while getting Sims by subscriber: %s", err.Error())

		return nil, err
	}

	resp := &pb.GetSubscriberByEmailResponse{
		Subscriber: dbSubscriberToPbSubscriber(subscriber, pbManagerSimsToPbSubscriberSims(simRep.Sims)),
	}

	return resp, nil
}

func (s *SubcriberServer) ListSubscribers(ctx context.Context, req *pb.ListSubscribersRequest) (*pb.ListSubscribersResponse, error) {
	log.Infof("List all subscribers")

	subscribers, err := s.subscriberRepo.ListSubscribers()
	if err != nil {
		log.WithError(err).Error("error while getting all subscribers")

		return nil, grpc.SqlErrorToGrpc(err, "subscribers")
	}

	simManagerClient, err := s.simManagerService.GetSimManagerService()
	if err != nil {
		log.Errorf("Failed to call SimManagerServiceClient. Error: %s", err.Error())

		return nil, err
	}

	simRep, err := simManagerClient.ListSims(ctx, &simMangerPb.ListSimsRequest{})
	if err != nil {
		log.Errorf("Failed to get Sims by subscriber. Error: %s", err.Error())

		return nil, err
	}

	allSims := simRep.Sims

	// Store Sims by their SubscriberId
	simMap := make(map[string][]*upb.Sim)
	for _, sim := range allSims {
		start, err := validation.FromString(sim.Package.StartDate)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}

		end, err := validation.FromString(sim.Package.EndDate)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}
		simMap[sim.SubscriberId] = append(simMap[sim.SubscriberId], &upb.Sim{
			Id:           sim.Id,
			SubscriberId: sim.SubscriberId,
			NetworkId:    sim.NetworkId,
			Iccid:        sim.Iccid,
			Msisdn:       sim.Msisdn,
			Package: &upb.Package{
				Id:        sim.Package.Id,
				StartDate: timestamppb.New(start),
				EndDate:   timestamppb.New(end),
			},
			Type:               sim.Type,
			Status:             sim.Status,
			IsPhysical:         sim.IsPhysical,
			FirstActivatedOn:   sim.FirstActivatedOn,
			LastActivatedOn:    sim.LastActivatedOn,
			ActivationsCount:   sim.ActivationsCount,
			DeactivationsCount: sim.DeactivationsCount,
			AllocatedAt:        sim.AllocatedAt,
		})
	}

	var res []*upb.Subscriber
	for _, sub := range subscribers {
		res = append(res, dbSubscriberToPbSubscriber(&sub, simMap[sub.SubscriberId.String()]))
	}
	subscriberList := &pb.ListSubscribersResponse{
		Subscribers: res,
	}

	return subscriberList, nil
}

func (s *SubcriberServer) GetByNetwork(ctx context.Context, req *pb.GetByNetworkRequest) (*pb.GetByNetworkResponse, error) {
	log.Infof("Get subscribers by network: %v ", req)

	networkIdReq := req.GetNetworkId()

	networkId, err := uuid.FromString(networkIdReq)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid networkId: %s", err.Error())
	}

	subscribers, err := s.subscriberRepo.GetByNetwork(networkId)
	if err != nil {
		log.WithError(err).Error("error while getting subscribers by network")

		return nil, grpc.SqlErrorToGrpc(err, "subscribers")
	}

	smc, err := s.simManagerService.GetSimManagerService()
	if err != nil {
		log.Errorf("Failed to get SimManagerServiceClient. Error: %s", err.Error())

		return nil, err
	}

	simRep, err := smc.GetSimsByNetwork(ctx, &simMangerPb.GetSimsByNetworkRequest{NetworkId: networkIdReq})
	if err != nil {
		log.Errorf("Failed to get Sims by network. Error: %s", err.Error())

		return nil, err
	}

	subscriberSims := pbManagerSimsToPbSubscriberSims(simRep.Sims)
	subscriberList := &pb.GetByNetworkResponse{
		Subscribers: dbSubScribersToPbSubscribers(subscribers, subscriberSims),
	}

	return subscriberList, nil
}

func (s *SubcriberServer) Update(ctx context.Context, req *pb.UpdateSubscriberRequest) (*pb.UpdateSubscriberResponse, error) {
	log.Infof("Updating subscriber: %v", req)

	subscriberId, err := uuid.FromString(req.GetSubscriberId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of subscriber uuid. Error %s", err.Error())
	}

	subscriber := &db.Subscriber{
		Email:                 req.GetEmail(),
		PhoneNumber:           req.GetPhoneNumber(),
		Address:               req.GetAddress(),
		ProofOfIdentification: req.GetProofOfIdentification(),
		IdSerial:              req.GetIdSerial(),
		SubscriberId:          subscriberId,
	}

	err = s.subscriberRepo.Update(subscriberId, *subscriber)
	if err != nil {
		log.Errorf("error while updating subscriber" + err.Error())

		return nil, grpc.SqlErrorToGrpc(err, "subscriber")
	}

	subscriberPb := dbSubscriberToPbSubscriber(subscriber, nil)

	route := s.subscriberRoutingKey.SetAction("update").SetObject("subscriber").MustBuild()
	_ = s.PublishEventMessage(route, &epb.UpdateSubscriber{
		Subscriber: subscriberPb,
	})

	return &pb.UpdateSubscriberResponse{}, nil
}

func (s *SubcriberServer) Delete(ctx context.Context, req *pb.DeleteSubscriberRequest) (*pb.DeleteSubscriberResponse, error) {
	subscriberIdReq := req.GetSubscriberId()

	subscriberId, err := uuid.FromString(subscriberIdReq)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of subscriber uuid. Error %s", err.Error())
	}

	subscriber, err := s.subscriberRepo.Get(subscriberId)
	if err != nil {
		log.Errorf("Error while getting subscriber: %s", err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "subscriber")
	}
	log.Infof("Delete Subscriber : %v ", subscriberId)

	err = s.subscriberRepo.Delete(subscriberId)
	if err != nil {
		log.WithError(err).Error("error while deleting subscriber")

		return nil, grpc.SqlErrorToGrpc(err, "subscriber")
	}

	subscriberPb := dbSubscriberToPbSubscriber(subscriber, nil)

	route := s.subscriberRoutingKey.SetAction("delete").SetObject("subscriber").MustBuild()
	_ = s.PublishEventMessage(route, &epb.RemoveSubscriber{
		Subscriber: subscriberPb,
	})

	return &pb.DeleteSubscriberResponse{}, nil
}

func (s *SubcriberServer) PublishEventMessage(route string, msg protoreflect.ProtoMessage) error {

	err := s.msgbus.PublishRequest(route, msg)
	if err != nil {
		log.Errorf("Failed to publish message %+v with key %+v. Errors %s", msg, route, err.Error())
	}
	return err

}

func dbSubScribersToPbSubscribers(subscriber []db.Subscriber, sims []*upb.Sim) []*upb.Subscriber {
	res := []*upb.Subscriber{}

	for _, u := range subscriber {
		subscriberSims := []*upb.Sim{}
		for _, sim := range sims {
			if sim.SubscriberId == u.SubscriberId.String() {
				subscriberSims = append(subscriberSims, sim)
			}
		}
		res = append(res, dbSubscriberToPbSubscriber(&u, subscriberSims))
	}

	return res
}
func pbManagerSimsToPbSubscriberSims(s []*simMangerPb.Sim) []*upb.Sim {
	res := []*upb.Sim{}
	for _, u := range s {
		ss := &upb.Sim{
			Id:                 u.Id,
			SubscriberId:       u.SubscriberId,
			NetworkId:          u.NetworkId,
			Iccid:              u.Iccid,
			Msisdn:             u.Msisdn,
			Type:               u.Type,
			Status:             u.Status,
			IsPhysical:         u.IsPhysical,
			FirstActivatedOn:   u.FirstActivatedOn,
			LastActivatedOn:    u.LastActivatedOn,
			DeactivationsCount: u.DeactivationsCount,
			AllocatedAt:        u.AllocatedAt,
		}

		res = append(res, ss)
	}

	return res
}

func dbSubscriberToPbSubscriber(s *db.Subscriber, simList []*upb.Sim) *upb.Subscriber {

	return &upb.Subscriber{
		FirstName:             s.FirstName,
		LastName:              s.LastName,
		Email:                 s.Email,
		SubscriberId:          s.SubscriberId.String(),
		ProofOfIdentification: s.ProofOfIdentification,
		Sim:                   simList,
		PhoneNumber:           s.PhoneNumber,
		IdSerial:              s.IdSerial,
		NetworkId:             s.NetworkId.String(),
		Gender:                s.Gender,
		Address:               s.Address,
		CreatedAt:             s.CreatedAt.String(),
		UpdatedAt:             s.UpdatedAt.String(),
		Dob:                   s.DOB,
	}
}
