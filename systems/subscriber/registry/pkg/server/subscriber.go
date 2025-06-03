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
	"time"

	"github.com/ukama/ukama/systems/common/ukama"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/common/grpc"
	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/validation"
	"github.com/ukama/ukama/systems/subscriber/registry/pkg"
	"github.com/ukama/ukama/systems/subscriber/registry/pkg/client"
	"github.com/ukama/ukama/systems/subscriber/registry/pkg/db"

	log "github.com/sirupsen/logrus"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	upb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
	cnucl "github.com/ukama/ukama/systems/common/rest/client/nucleus"
	creg "github.com/ukama/ukama/systems/common/rest/client/registry"
	uuid "github.com/ukama/ukama/systems/common/uuid"
	validate "github.com/ukama/ukama/systems/common/validation"
	pb "github.com/ukama/ukama/systems/subscriber/registry/pb/gen"
	simMangerPb "github.com/ukama/ukama/systems/subscriber/sim-manager/pb/gen"
)

const MAX_DELETION_RETRIES = 3

type SubscriberServer struct {
	orgName              string
	orgId                string
	msgbus               mb.MsgBusServiceClient
	subscriberRepo       db.SubscriberRepo
	subscriberRoutingKey msgbus.RoutingKeyBuilder
	simManagerService    client.SimManagerClientProvider
	orgClient            cnucl.OrgClient
	networkClient        creg.NetworkClient
	deletionCheckCancel  context.CancelFunc
	pb.UnimplementedRegistryServiceServer
}

func NewSubscriberServer(orgName string, subscriberRepo db.SubscriberRepo, msgBus mb.MsgBusServiceClient, simManagerService client.SimManagerClientProvider, orgId string, orgService cnucl.OrgClient, networkClient creg.NetworkClient) *SubscriberServer {
    server := &SubscriberServer{
        orgName:              orgName,
        subscriberRepo:       subscriberRepo,
        msgbus:               msgBus,
        simManagerService:    simManagerService,
        subscriberRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
        orgId:                orgId,
        orgClient:            orgService,
        networkClient:        networkClient,
    }
    
    go server.startDeletionCheck()
    
    return server
}


func (s *SubscriberServer) Add(ctx context.Context, req *pb.AddSubscriberRequest) (*pb.AddSubscriberResponse, error) {
	log.Infof("Adding subscriber: %v", req)

	var dob string
	var err error

	if req.GetDob() != "" {
		dob, err = validate.ValidateDate(req.GetDob())
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument,
				"invalid format for DoB value %s. Error: %v", req.GetDob(), err)
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
		Name:                  req.GetName(),
		NetworkId:             nid,
		Email:                 strings.ToLower(req.GetEmail()),
		PhoneNumber:           req.GetPhoneNumber(),
		Gender:                req.GetGender(),
		Address:               req.GetAddress(),
		ProofOfIdentification: req.GetProofOfIdentification(),
		SubscriberStatus:      ukama.SubscriberStatusActive, 
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

	route := s.subscriberRoutingKey.SetAction("create").SetObject("subscriber").MustBuild()
	log.Infof("Pushing add subscriber event to %v", route)
	_ = s.PublishEventMessage(route, &epb.EventSubscriberAdded{
		Dob:          req.GetDob(),
		Email:        req.GetEmail(),
		Gender:       req.GetGender(),
		Name:         req.GetName(),
		NetworkId:    req.GetNetworkId(),
		PhoneNumber:  req.GetPhoneNumber(),
		SubscriberId: subscriber.SubscriberId.String(),
	})

	return &pb.AddSubscriberResponse{
		Subscriber: dbSubscriberToPbSubscriber(subscriber, []*upb.Sim{}),
	}, nil
}

func (s *SubscriberServer) Get(ctx context.Context, req *pb.GetSubscriberRequest) (*pb.GetSubscriberResponse, error) {
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

	simRep, err := smc.ListSims(ctx, &simMangerPb.ListSimsRequest{
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

func (s *SubscriberServer) GetByEmail(ctx context.Context, req *pb.GetSubscriberByEmailRequest) (*pb.GetSubscriberByEmailResponse, error) {
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

	simRep, err := smc.ListSims(ctx, &simMangerPb.ListSimsRequest{
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

func (s *SubscriberServer) ListSubscribers(ctx context.Context, req *pb.ListSubscribersRequest) (*pb.ListSubscribersResponse, error) {
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
			return nil, status.Errorf(codes.InvalidArgument,
				"invalid format for Package.StartDate value %s. Error: %v", sim.Package.StartDate, err)
		}

		end, err := validation.FromString(sim.Package.EndDate)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument,
				"invalid format for Package.EndDate value %s. Error: %v", sim.Package.EndDate, err)
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

func (s *SubscriberServer) GetByNetwork(ctx context.Context, req *pb.GetByNetworkRequest) (*pb.GetByNetworkResponse, error) {
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

	simRep, err := smc.ListSims(ctx, &simMangerPb.ListSimsRequest{NetworkId: networkIdReq})
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

func (s *SubscriberServer) Update(ctx context.Context, req *pb.UpdateSubscriberRequest) (*pb.UpdateSubscriberResponse, error) {
	log.Infof("Updating subscriber: %v", req)

	subscriberId, err := uuid.FromString(req.GetSubscriberId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of subscriber uuid. Error %s", err.Error())
	}

	subscriber := &db.Subscriber{
		Name:                  req.GetName(),
		PhoneNumber:           req.GetPhoneNumber(),
		Address:               req.GetAddress(),
		ProofOfIdentification: req.GetProofOfIdentification(),
		IdSerial:              req.GetIdSerial(),
		SubscriberId:          subscriberId,
	}

	err = s.subscriberRepo.Update(subscriberId, *subscriber)
	if err != nil {
		log.Errorf("error while updating subscriber %s. Error: %v", subscriberId, err)

		return nil, grpc.SqlErrorToGrpc(err, "subscriber")
	}

	route := s.subscriberRoutingKey.SetAction("update").SetObject("subscriber").MustBuild()
	log.Infof("Pushing update subscriber event to %v", route)
	_ = s.PublishEventMessage(route, &epb.EventSubscriberUpdate{
		Email:                 subscriber.Email,
		Address:               subscriber.Address,
		IdSerial:              subscriber.IdSerial,
		PhoneNumber:           subscriber.PhoneNumber,
		SubscriberId:          subscriber.SubscriberId.String(),
		ProofOfIdentification: subscriber.ProofOfIdentification,
	})

	return &pb.UpdateSubscriberResponse{}, nil
}

func (s *SubscriberServer) Delete(ctx context.Context, req *pb.DeleteSubscriberRequest) (*pb.DeleteSubscriberResponse, error) {
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

    if subscriber.SubscriberStatus == ukama.SubscriberStatusPendingDeletion {
        return &pb.DeleteSubscriberResponse{
        }, nil
    }

    err = s.subscriberRepo.MarkAsPendingDeletion(subscriberId)
    if err != nil {
        log.Errorf("Error marking subscriber as pending deletion: %s", err.Error())
        return nil, grpc.SqlErrorToGrpc(err, "subscriber")
    }

    log.Infof("Initiating subscriber deletion: %v", subscriberId)

    simManagerClient, err := s.simManagerService.GetSimManagerService()
    if err != nil {
        log.Errorf("Failed to get SimManagerServiceClient. Error: %s", err.Error())
        return nil, err
    }

    _, err = simManagerClient.TerminateSimsForSubscriber(ctx, &simMangerPb.TerminateSimsForSubscriberRequest{
        SubscriberId: subscriber.SubscriberId.String(),
    })
    if err != nil {
        log.Errorf("Failed to terminate SIMs for subscriber %s: %v", subscriberId, err)
        return nil, status.Errorf(codes.Internal, "Failed to terminate SIMs: %v", err)
    }
    
    log.Infof("Successfully initiated deletion for subscriber: %v. SIM Manager will handle coordination.", subscriberId)
    return &pb.DeleteSubscriberResponse{
    }, nil
}


func (s *SubscriberServer) PublishEventMessage(route string, msg protoreflect.ProtoMessage) error {

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
		Name:                  s.Name,
		Email:                 s.Email,
		SubscriberId:          s.SubscriberId.String(),
		ProofOfIdentification: s.ProofOfIdentification,
		Sim:                   simList,
		PhoneNumber:           s.PhoneNumber,
		IdSerial:              s.IdSerial,
		NetworkId:             s.NetworkId.String(),
		Gender:                s.Gender,
		Address:               s.Address,
		SubscriberStatus:                s.SubscriberStatus.String(),
		CreatedAt:             s.CreatedAt.String(),
		UpdatedAt:             s.UpdatedAt.String(),
		Dob:                   s.DOB,
	}
}
func (s *SubscriberServer) checkStuckDeletions() {
    threshold := time.Now().Add(-15 * time.Minute)
    
    stuckSubscribers, err := s.subscriberRepo.FindPendingDeletionBefore(threshold)
    if err != nil {
        log.Errorf("Error checking for stuck deletions: %v", err)
        return
    }
    
    if len(stuckSubscribers) == 0 {
        return
    }
    
    log.Infof("Found %d subscribers stuck in pending deletion state", len(stuckSubscribers))
    
    for _, subscriber := range stuckSubscribers {
        if subscriber.DeletionRetryCount >= MAX_DELETION_RETRIES {
            log.Errorf("Subscriber %s has exceeded maximum retry attempts (%d). Manual intervention required.", 
                subscriber.SubscriberId, MAX_DELETION_RETRIES)
            continue
        }
        
        log.Infof("Retrying deletion for subscriber %s (attempt %d/%d)", 
            subscriber.SubscriberId, subscriber.DeletionRetryCount+1, MAX_DELETION_RETRIES)
        
        ctx := context.Background()
        go s.retrySubscriberDeletion(ctx, subscriber)
    }
}



func (s *SubscriberServer) retrySubscriberDeletion(ctx context.Context, subscriber db.Subscriber) {
    err := s.subscriberRepo.IncrementDeletionRetry(subscriber.SubscriberId)
    if err != nil {
        log.Errorf("Failed to increment retry count for subscriber %s: %v", 
            subscriber.SubscriberId, err)
        return
    }
    
    simManagerClient, err := s.simManagerService.GetSimManagerService()
    if err != nil {
        log.Errorf("Failed to get SimManagerClient for retry: %v", err)
        return
    }
    
    _, err = simManagerClient.TerminateSimsForSubscriber(ctx, &simMangerPb.TerminateSimsForSubscriberRequest{
        SubscriberId: subscriber.SubscriberId.String(),
    })
    
    if err != nil {
        log.Errorf("Retry failed for subscriber %s: %v", subscriber.SubscriberId, err)
        
        if subscriber.DeletionRetryCount+1 >= MAX_DELETION_RETRIES {
            log.Errorf("Subscriber %s deletion failed after %d attempts. Manual intervention required.", 
                subscriber.SubscriberId, MAX_DELETION_RETRIES)
        }
    }
}
func (s *SubscriberServer) startDeletionCheck() {
    ticker := time.NewTicker(10 * time.Minute)
    defer ticker.Stop()
    
    ctx, cancel := context.WithCancel(context.Background())
    s.deletionCheckCancel = cancel
    defer cancel()
    
    for {
        select {
        case <-ticker.C:
            s.checkStuckDeletions()
        case <-ctx.Done():
            return
        }
    }
}
func (s *SubscriberServer) Shutdown() {
    if s.deletionCheckCancel != nil {
        s.deletionCheckCancel()
    }
}