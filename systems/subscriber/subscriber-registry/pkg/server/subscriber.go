package server

import (
	"context"
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
	smc, err := s.simManagerService.GetSimManagerService()
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
	subscriberSims := pbManagerSimsToPbSubscriberSims(simRep.Sims)


	resp := &pb.GetSubscriberResponse{Subscriber: dbSubscriberToPbSubscriber(subscriber,subscriberSims)}

	return resp, nil

}




func dbSubScribersToPbSubscribers(subscriber []db.Subscriber,simMaping simMangerPb.GetPackagesBySimResponse ) []*pb.Subscriber {
	res := []*pb.Subscriber{}
	for _, u := range subscriber {
	
		//get list of sims from simManager for that subscriber
		//convert that list of sim manager sims to subscriber sims
		//add the list of subscriber sim to that subscriber
		res = append(res, dbSubscriberToPbSubscriber(&u))
	}
	return res
}
func pbManagerSimsToPbSubscriberSims (s []*simMangerPb.Sim) []*pb.Sim{
	res := []*pb.Sim{}
	for _, u := range s {
		ss := &pb.Sim{
			Id:u.Id,
			SubscriberID: u.SubscriberID,
			NetworkID:u.NetworkID,
		}

		res = append(res, ss)
	}
	return res
}

func dbSubscriberToPbSubscriber(s *db.Subscriber,sims []*pb.Sim) *pb.Subscriber {
	pbTimestamp := s.DOB.Format("2006-01-02")
	
	return &pb.Subscriber{
		FirstName:             s.FirstName,
		LastName:              s.LastName,
		Email:                 s.Email,
		SubscriberID:          s.SubscriberID.String(),
		ProofOfIdentification: s.ProofOfIdentification,
		Sim:                   sims,
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
