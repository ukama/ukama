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

}

func NewSubscriberServer(subscriberRepo db.SubscriberRepo, msgBus mb.MsgBusServiceClient,simManagerService clientPkg.SimManagerClientProvider) *SubcriberServer {
	return &SubcriberServer{subscriberRepo: subscriberRepo,
		msgbus:               msgBus,
		simManagerService:simManagerService,
		subscriberRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetContainer(pkg.ServiceName)}
}




func (s *SubcriberServer) Get(ctx context.Context, req *pb.GetSubscriberRequest) (*pb.GetSubscriberResponse, error) {
	subscriberIdReq := req.GetSubscriberID()

	subscriberID, error := uuid.FromString(subscriberIdReq)
	if error != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of subscriber uuid. Error %s", error.Error())
	}

	logrus.Infof("GetSubscriber : %v ", subscriberID)
	subscriber, err := s.subscriberRepo.Get(subscriberID)

	if err != nil {
		logrus.Error("error getting a subscriber" + err.Error())

		return nil, grpc.SqlErrorToGrpc(err, "subscriber")
	}
	smc, err := s.simManagerService.GetSimsBySubscriber()
	if err != nil {
		logrus.Error("Failed to get SimManagerServiceClient. Error: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	simRep, err := smc.GetSimsBySubscriber(ctx, &simMangerPb.GetSimsBySubscriberRequest{
		SubscriberID: subscriberID.String(),
	})

	if err != nil {
		log.Fatalf("Failed to get Sims by subscriber. Error: %v", err)
	}
	fmt.Println(simRep)


	resp := &pb.GetSubscriberResponse{Subscriber: dbSubscriberToPbSubscribers(subscriber)}

	return resp, nil

}




func dbsubscriberToPbSubscribers(subscriber []db.Subscriber) []*pb.Subscriber {
	res := []*pb.Subscriber{}
	for _, u := range subscriber {
		res = append(res, dbSubscriberToPbSubscribers(&u))
	}
	return res
}

func dbSubscriberToPbSubscribers(s *db.Subscriber) *pb.Subscriber {
	pbTimestamp := s.DOB.Format("2006-01-02")

	// Create a slice of mock SIMs
	simList := []*pb.Sim{
		{
			Id:           "12345",
			SubscriberID: s.SubscriberID.String(),
			Iccid:        "9876543210",
			Msisdn:       "123-456-7890",
			IsPhysical:   true,
		},
		{
			Id:           "54321",
			SubscriberID: s.SubscriberID.String(),
			Iccid:        "0123456789",
			Msisdn:       "123-456-7891",
			IsPhysical:   false,
		},
		{
			Id:           "67890",
			SubscriberID: s.SubscriberID.String(),
			Iccid:        "0123456789",
			Msisdn:       "123-456-7891",
			IsPhysical:   false,
		},
	}

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
