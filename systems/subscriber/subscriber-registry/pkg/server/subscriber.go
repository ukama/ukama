package server

import (
	"context"
	"fmt"
	"log"
	"time"

	uuid "github.com/ukama/ukama/systems/common/uuid"

	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/grpc"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
	simMangerPb "github.com/ukama/ukama/systems/subscriber/sim-manager/pb/gen"
	pb "github.com/ukama/ukama/systems/subscriber/subscriber-registry/pb/gen"
	"github.com/ukama/ukama/systems/subscriber/subscriber-registry/pkg"
	client "github.com/ukama/ukama/systems/subscriber/subscriber-registry/pkg/client"
	clientPkg "github.com/ukama/ukama/systems/subscriber/subscriber-registry/pkg/client"
	"github.com/ukama/ukama/systems/subscriber/subscriber-registry/pkg/db"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SubcriberServer struct {
	subscriberRepo       db.SubscriberRepo
	msgbus               mb.MsgBusServiceClient
	subscriberRoutingKey msgbus.RoutingKeyBuilder
	pb.UnimplementedSubscriberRegistryServiceServer
	simManagerService clientPkg.SimManagerClientProvider
	network        client.Network

}

func NewSubscriberServer(subscriberRepo db.SubscriberRepo, msgBus mb.MsgBusServiceClient, simManagerService clientPkg.SimManagerClientProvider, network client.Network) *SubcriberServer {
	return &SubcriberServer{subscriberRepo: subscriberRepo,
		msgbus:               msgBus,
		simManagerService:    simManagerService,
		network:  network,
		subscriberRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetContainer(pkg.ServiceName)}
}

func (s *SubcriberServer) Add(ctx context.Context, req *pb.AddSubscriberRequest) (*pb.AddSubscriberResponse, error) {
	logrus.Infof("Adding subscriber: %v", req)
	networkID := uuid.FromStringOrNil(req.GetNetworkID())
	orgID := uuid.FromStringOrNil(req.GetOrgID())
	if networkID == uuid.Nil {
		return nil, status.Errorf(codes.InvalidArgument, "NetworkID must not be empty")
	}
	if orgID == uuid.Nil {
		return nil, status.Errorf(codes.InvalidArgument, "OrgID must not be empty")
	}
	subscriberID := uuid.NewV4()
	err := s.network.ValidateNetwork(networkID.String(),orgID.String() )
	if err != nil {
		return nil, fmt.Errorf("error validating network")
	}
	
	subscriber := &db.Subscriber{
		OrgID:                 orgID,
		SubscriberID:          subscriberID,
		FirstName:             req.GetFirstName(),
		LastName:              req.GetLastName(),
		NetworkID:             networkID,
		Email:                 req.GetEmail(),
		PhoneNumber:           req.GetPhoneNumber(),
		Gender:                req.GetGender(),
		Address:               req.GetAddress(),
		ProofOfIdentification: req.GetProofOfIdentification(),
		DOB:                   req.DateOfBirth.AsTime(),
		IdSerial:              req.GetIdSerial(),
	}
	err = s.subscriberRepo.Add(subscriber)
	if err != nil {
		logrus.Error("error while adding subscriber" + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "subscriber")
	}

	return &pb.AddSubscriberResponse{
		Subscriber: &pb.Subscriber{
			OrgID:                 orgID.String(),
			SubscriberID:          subscriberID.String(),
			FirstName:             req.GetFirstName(),
			LastName:              req.GetLastName(),
			NetworkID:             networkID.String(),
			Email:                 req.GetEmail(),
			PhoneNumber:           req.GetPhoneNumber(),
			Gender:                req.GetGender(),
			Address:               req.GetAddress(),
			ProofOfIdentification: req.GetProofOfIdentification(),
			DateOfBirth:           req.GetDateOfBirth().String(),
			IdSerial:              req.GetIdSerial()},
	}, nil

}

func (s *SubcriberServer) Get(ctx context.Context, req *pb.GetSubscriberRequest) (*pb.GetSubscriberResponse, error) {
	subscriberIdReq := req.GetSubscriberID()
	subscriberID, err := uuid.FromString(subscriberIdReq)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid subscriberID format: %v", err)
	}

	logrus.Infof("Getting subscriber with ID: %v", subscriberID)
	subscriber, err := s.subscriberRepo.Get(subscriberID)
	if err != nil {
		logrus.Error("Error while getting subscriber: %v", err)
		return nil, grpc.SqlErrorToGrpc(err, "subscriber")
	}

	smc, err := s.simManagerService.GetSimManagerService()
	if err != nil {
		logrus.Error("Error while getting SimManagerServiceClient: %v", err)
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	simRep, err := smc.GetSimsBySubscriber(ctx, &simMangerPb.GetSimsBySubscriberRequest{
		SubscriberID: subscriberID.String(),
	})
	if err != nil {
		logrus.Error("Error while getting Sims by subscriber: %v", err)
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
		logrus.Error("Failed to get SimManagerServiceClient. Error: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	simRep, err := simManagerClient.ListSims(ctx, &simMangerPb.ListSimsRequest{})
	if err != nil {
		log.Fatalf("Failed to get Sims by subscriber. Error: %v", err)
	}
	allSims := simRep.Sims

	var res []*pb.Subscriber
	for _, sub := range subscribers {
		simList := []*pb.Sim{}
		for _, sim := range allSims {
			if sim.SubscriberID == sub.SubscriberID.String() {
				simList = append(simList, &pb.Sim{
					Id:           sim.Id,
					SubscriberID: sim.SubscriberID,
					NetworkID:    sim.NetworkID,
				})
			}
		}
		res = append(res, dbSubscriberToPbSubscriber(&sub, simList))
	}

	subscriberList := &pb.ListSubscribersResponse{
		Subscribers: res,
	}

	return subscriberList, nil
}

func (s *SubcriberServer) GetByNetwork(ctx context.Context, req *pb.GetByNetworkRequest) (*pb.GetByNetworkResponse, error) {
	networkIdReq := req.GetNetworkID()
	logrus.Infof("Get subscribers by network: %v ", networkIdReq)
	networkID, err := uuid.FromString(networkIdReq)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid networkID: %v", err)
	}

	subscribers, err := s.subscriberRepo.GetByNetwork(networkID)
	if err != nil {
		logrus.WithError(err).Error("error while getting subscribers by network")
		return nil, grpc.SqlErrorToGrpc(err, "subscribers")
	}

	smc, err := s.simManagerService.GetSimManagerService()
	if err != nil {
		logrus.Error("Failed to get SimManagerServiceClient. Error: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	simRep, err := smc.GetSimsByNetwork(ctx, &simMangerPb.GetSimsByNetworkRequest{NetworkID: networkIdReq})
	if err != nil {
		log.Fatalf("Failed to get Sims by network. Error: %v", err)
	}

	subscriberSims := pbManagerSimsToPbSubscriberSims(simRep.Sims)
	subscriberList := &pb.GetByNetworkResponse{
		Subscribers: dbSubScribersToPbSubscribers(subscribers, subscriberSims),
	}

	return subscriberList, nil
}

func (s *SubcriberServer) Delete(ctx context.Context, req *pb.DeleteSubscriberRequest) (*pb.DeleteSubscriberResponse, error) {
	subscriberIdReq := req.GetSubscriberID()
	subscriberID, error := uuid.FromString(subscriberIdReq)

	if error != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of subscriber uuid. Error %s", error.Error())
	}
	logrus.Infof("Delete Subscriber : %v ", subscriberID)
	er := s.subscriberRepo.Delete(subscriberID)
	if er != nil {
		logrus.WithError(er).Error("error while deleting subscriber")
		return nil, grpc.SqlErrorToGrpc(er, "subscriber")
	}
	route := s.subscriberRoutingKey.SetAction("delete").SetObject("subscriber").MustBuild()
	err := s.msgbus.PublishRequest(route, req)
	if err != nil {
		logrus.Errorf("Failed to publish message %+v with key %+v. Errors %s", req, route, err.Error())
	}
	return &pb.DeleteSubscriberResponse{}, nil
}

func getSimListBySubscriber(subscriberID string) []*simMangerPb.Sim {
	// Call simManager to get list of sims for the given subscriberID
	var s *SubcriberServer
	smc, err := s.simManagerService.GetSimManagerService()
	if err != nil {
		logrus.Error("Failed to get SimManagerServiceClient. Error: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	simRep, err := smc.GetSimsBySubscriber(ctx, &simMangerPb.GetSimsBySubscriberRequest{
		SubscriberID: subscriberID,
	})

	if err != nil {
		log.Fatalf("Failed to get Sims by subscriber. Error: %v", err)
	}

	return simRep.Sims
}

func dbSubScribersToPbSubscribers(subscriber []db.Subscriber, sims []*pb.Sim) []*pb.Subscriber {
	res := []*pb.Subscriber{}

	for _, u := range subscriber {
		subscriberSims := []*pb.Sim{}
		for _, sim := range sims {
			if sim.SubscriberID == u.SubscriberID.String() {
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
			SubscriberID:       u.SubscriberID,
			NetworkID:          u.NetworkID,
			OrgID:              u.OrgID,
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
	pbTimestamp := s.DOB.Format("2006-01-02")

	// simList := []*pb.Sim{}
	return &pb.Subscriber{
		FirstName:             s.FirstName,
		LastName:              s.LastName,
		Email:                 s.Email,
		SubscriberID:          s.SubscriberID.String(),
		ProofOfIdentification: s.ProofOfIdentification,
		Sim:                   simList,
		PhoneNumber:           s.PhoneNumber,
		IdSerial:              s.IdSerial,
		NetworkID:             s.NetworkID.String(),
		Gender:                s.Gender,
		Address:               s.Address,
		CreatedAt:             s.CreatedAt.String(),
		UpdatedAt:             s.UpdatedAt.String(),
		DateOfBirth:           pbTimestamp,
		OrgID:                 s.OrgID.String(),
	}

}
