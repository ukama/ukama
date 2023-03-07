package server

import (
	"context"

	uuid "github.com/ukama/ukama/systems/common/uuid"

	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/grpc"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
	pb "github.com/ukama/ukama/systems/subscriber/registry/pb/gen"
	"github.com/ukama/ukama/systems/subscriber/registry/pkg"
	"github.com/ukama/ukama/systems/subscriber/registry/pkg/client"
	"github.com/ukama/ukama/systems/subscriber/registry/pkg/db"
	utils "github.com/ukama/ukama/systems/subscriber/registry/pkg/util"
	simMangerPb "github.com/ukama/ukama/systems/subscriber/sim-manager/pb/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SubcriberServer struct {
	subscriberRepo       db.SubscriberRepo
	msgbus               mb.MsgBusServiceClient
	subscriberRoutingKey msgbus.RoutingKeyBuilder
	pb.UnimplementedRegistryServiceServer
	simManagerService client.SimManagerClientProvider
	network           client.NetworkInfoClient
}

func NewSubscriberServer(subscriberRepo db.SubscriberRepo, msgBus mb.MsgBusServiceClient, simManagerService client.SimManagerClientProvider, network client.NetworkInfoClient) *SubcriberServer {
	return &SubcriberServer{subscriberRepo: subscriberRepo,
		msgbus:               msgBus,
		simManagerService:    simManagerService,
		network:              network,
		subscriberRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetContainer(pkg.ServiceName)}
}

func (s *SubcriberServer) Add(ctx context.Context, req *pb.AddSubscriberRequest) (*pb.AddSubscriberResponse, error) {
	logrus.Infof("Adding subscriber: %v", req)
	networkId, err := uuid.FromString(req.GetNetworkId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of network uuid. Error %s", err.Error())
	}
	orgId, err := uuid.FromString(req.GetOrgId())

	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of org uuid. Error %s", err.Error())
	}

	dob, err := utils.ValidateDOB(req.GetDob())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	subscriberId := uuid.NewV4()
	err = s.network.ValidateNetwork(networkId.String(), orgId.String())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "network not found for that org %s", err.Error())
	}

	subscriber := &db.Subscriber{
		OrgId:                 orgId,
		SubscriberId:          subscriberId,
		FirstName:             req.GetFirstName(),
		LastName:              req.GetLastName(),
		NetworkId:             networkId,
		Email:                 req.GetEmail(),
		PhoneNumber:           req.GetPhoneNumber(),
		Gender:                req.GetGender(),
		Address:               req.GetAddress(),
		ProofOfIdentification: req.GetProofOfIdentification(),
		DOB:                   dob,
		IdSerial:              req.GetIdSerial(),
	}
	err = s.subscriberRepo.Add(subscriber)
	if err != nil {
		logrus.Error("error while adding subscriber" + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "subscriber")
	}

	return &pb.AddSubscriberResponse{
		Subscriber: dbSubscriberToPbSubscriber(subscriber, nil),
	}, nil

}

func (s *SubcriberServer) Get(ctx context.Context, req *pb.GetSubscriberRequest) (*pb.GetSubscriberResponse, error) {

	subscriberIdReq := req.GetSubscriberId()
	subscriberId, err := uuid.FromString(subscriberIdReq)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid subscriberId format: %v", err.Error())
	}

	logrus.Infof("Getting subscriber with ID: %v", subscriberId)
	subscriber, err := s.subscriberRepo.Get(subscriberId)
	if err != nil {
		logrus.Errorf("Error while getting subscriber: %s", err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "subscriber")
	}

	smc, err := s.simManagerService.GetSimManagerService()
	if err != nil {
		logrus.Errorf("Error while calling SimManagerServiceClient: %s", err.Error())
		return nil, err
	}

	simRep, err := smc.GetSimsBySubscriber(ctx, &simMangerPb.GetSimsBySubscriberRequest{
		SubscriberId: subscriberId.String(),
	})
	if err != nil {
		logrus.Errorf("Error while getting Sims by subscriber: %s", err.Error())
		return nil, err
	}

	resp := &pb.GetSubscriberResponse{
		Subscriber: dbSubscriberToPbSubscriber(subscriber, pbManagerSimsToPbSubscriberSims(simRep.Sims)),
	}
	return resp, nil
}

func (s *SubcriberServer) ListSubscribers(ctx context.Context, req *pb.ListSubscribersRequest) (*pb.ListSubscribersResponse, error) {
	logrus.Infof("List all subscribers")

	subscribers, err := s.subscriberRepo.ListSubscribers()
	if err != nil {
		logrus.WithError(err).Error("error while getting all subscribers")
		return nil, grpc.SqlErrorToGrpc(err, "subscribers")
	}

	simManagerClient, err := s.simManagerService.GetSimManagerService()
	if err != nil {
		logrus.Errorf("Failed to call SimManagerServiceClient. Error: %s", err.Error())
		return nil, err
	}

	simRep, err := simManagerClient.ListSims(ctx, &simMangerPb.ListSimsRequest{})
	if err != nil {
		logrus.Errorf("Failed to get Sims by subscriber. Error: %s", err.Error())
		return nil, err
	}

	allSims := simRep.Sims

	// Store Sims by their SubscriberId
	simMap := make(map[string][]*pb.Sim)
	for _, sim := range allSims {
		simMap[sim.SubscriberId] = append(simMap[sim.SubscriberId], &pb.Sim{
			Id:           sim.Id,
			SubscriberId: sim.SubscriberId,
			NetworkId:    sim.NetworkId,
			OrgId:        sim.OrgId,
			Iccid:        sim.Iccid,
			Msisdn:       sim.Msisdn,
			Package: &pb.Package{
				Id:        sim.Package.Id,
				StartDate: sim.Package.StartDate,
				EndDate:   sim.Package.EndDate,
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

	var res []*pb.Subscriber
	for _, sub := range subscribers {
		res = append(res, dbSubscriberToPbSubscriber(&sub, simMap[sub.SubscriberId.String()]))
	}
	subscriberList := &pb.ListSubscribersResponse{
		Subscribers: res,
	}

	return subscriberList, nil
}
func (s *SubcriberServer) GetByNetwork(ctx context.Context, req *pb.GetByNetworkRequest) (*pb.GetByNetworkResponse, error) {
	networkIdReq := req.GetNetworkId()
	logrus.Infof("Get subscribers by network: %v ", networkIdReq)
	networkId, err := uuid.FromString(networkIdReq)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid networkId: %s", err.Error())
	}

	subscribers, err := s.subscriberRepo.GetByNetwork(networkId)
	if err != nil {
		logrus.WithError(err).Error("error while getting subscribers by network")
		return nil, grpc.SqlErrorToGrpc(err, "subscribers")
	}

	smc, err := s.simManagerService.GetSimManagerService()
	if err != nil {
		logrus.Errorf("Failed to get SimManagerServiceClient. Error: %s", err.Error())
		return nil, err
	}

	simRep, err := smc.GetSimsByNetwork(ctx, &simMangerPb.GetSimsByNetworkRequest{NetworkId: networkIdReq})
	if err != nil {
		logrus.Errorf("Failed to get Sims by network. Error: %s", err.Error())
		return nil, err
	}

	subscriberSims := pbManagerSimsToPbSubscriberSims(simRep.Sims)
	subscriberList := &pb.GetByNetworkResponse{
		Subscribers: dbSubScribersToPbSubscribers(subscribers, subscriberSims),
	}

	return subscriberList, nil
}
func (s *SubcriberServer) Update(ctx context.Context, req *pb.UpdateSubscriberRequest) (*pb.UpdateSubscriberResponse, error) {
	logrus.Infof("Updating subscriber: %v", req)
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
	}

	err = s.subscriberRepo.Update(subscriberId, *subscriber)
	if err != nil {
		logrus.Errorf("error while updating subscriber" + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "subscriber")
	}

	return &pb.UpdateSubscriberResponse{}, nil
}
func (s *SubcriberServer) Delete(ctx context.Context, req *pb.DeleteSubscriberRequest) (*pb.DeleteSubscriberResponse, error) {
	subscriberIdReq := req.GetSubscriberId()
	subscriberId, err := uuid.FromString(subscriberIdReq)

	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of subscriber uuid. Error %s", err.Error())
	}
	logrus.Infof("Delete Subscriber : %v ", subscriberId)
	err = s.subscriberRepo.Delete(subscriberId)
	if err != nil {
		logrus.WithError(err).Error("error while deleting subscriber")
		return nil, grpc.SqlErrorToGrpc(err, "subscriber")
	}
	route := s.subscriberRoutingKey.SetAction("delete").SetObject("subscriber").MustBuild()
	err = s.msgbus.PublishRequest(route, req)
	if err != nil {
		logrus.Errorf("Failed to publish message %+v with key %+v. Errors %s", req, route, err.Error())
	}
	return &pb.DeleteSubscriberResponse{}, nil
}

func dbSubScribersToPbSubscribers(subscriber []db.Subscriber, sims []*pb.Sim) []*pb.Subscriber {
	res := []*pb.Subscriber{}

	for _, u := range subscriber {
		subscriberSims := []*pb.Sim{}
		for _, sim := range sims {
			if sim.SubscriberId == u.SubscriberId.String() {
				subscriberSims = append(subscriberSims, sim)
			}
		}
		res = append(res, dbSubscriberToPbSubscriber(&u, subscriberSims))
	}
	return res
}
func pbManagerSimsToPbSubscriberSims(s []*simMangerPb.Sim) []*pb.Sim {
	res := []*pb.Sim{}
	for _, u := range s {
		ss := &pb.Sim{
			Id:                 u.Id,
			SubscriberId:       u.SubscriberId,
			NetworkId:          u.NetworkId,
			OrgId:              u.OrgId,
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

func dbSubscriberToPbSubscriber(s *db.Subscriber, simList []*pb.Sim) *pb.Subscriber {

	return &pb.Subscriber{
		FirstName:             s.FirstName,
		LastName:              s.LastName,
		Email:                 s.Email,
		SubscriberId:          s.SubscriberId.String(),
		ProofOfIdentification: s.ProofOfIdentification,
		Sim:                   simList,
		PhoneNumber:           s.PhoneNumber,
		IdSerial:              s.IdSerial,
		NetworkId:             s.NetworkId.String(),
		OrgId:                 s.OrgId.String(),
		Gender:                s.Gender,
		Address:               s.Address,
		CreatedAt:             s.CreatedAt.String(),
		UpdatedAt:             s.UpdatedAt.String(),
		Dob:           s.DOB,
	}

}
