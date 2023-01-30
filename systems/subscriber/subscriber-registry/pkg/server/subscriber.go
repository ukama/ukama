package server

import (
	"context"

	uuid "github.com/ukama/ukama/systems/common/uuid"

	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/grpc"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
	pb "github.com/ukama/ukama/systems/subscriber/subscriber-registry/pb/gen"
	"github.com/ukama/ukama/systems/subscriber/subscriber-registry/pkg"
	"github.com/ukama/ukama/systems/subscriber/subscriber-registry/pkg/db"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type SubcriberServer struct {
	subscriberRepo       db.SubscriberRepo
	msgbus               mb.MsgBusServiceClient
	subscriberRoutingKey msgbus.RoutingKeyBuilder
	pb.UnimplementedSubscriberRegistryServiceServer
}



func NewSubscriberServer(subscriberRepo db.SubscriberRepo, msgBus mb.MsgBusServiceClient) *SubcriberServer {
	return &SubcriberServer{subscriberRepo: subscriberRepo,
		msgbus:         msgBus,
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
	timestamp := &timestamppb.Timestamp{
		Seconds: req.DateOfBirth.GetSeconds(),
		Nanos:   req.DateOfBirth.GetNanos(),
	}

	birthday := timestamp.AsTime()

	subscriber := &db.Subscriber{
		OrgID:orgID,
		SubscriberID:          subscriberID,
		FirstName:             req.GetFirstName(),
		LastName:              req.GetLastName(),
		NetworkID:             networkID,
		Email:                 req.GetEmail(),
		PhoneNumber:           req.GetPhoneNumber(),
		Gender:                req.GetGender(),
		Address:               req.GetAddress(),
		ProofOfIdentification: req.GetProofOfIdentification(),
		DOB:                   birthday,
		IdSerial:              req.GetIdSerial(),
	}
	err := s.subscriberRepo.Add(subscriber)
	if err != nil {
		logrus.Error("error while adding subscriber" + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "subscriber")
	}

	return &pb.AddSubscriberResponse{
		Subscriber: &pb.Subscriber{
			OrgID: orgID.String(),
			SubscriberID:          subscriberID.String(),
			FirstName:             req.GetFirstName(),
			LastName:              req.GetLastName(),
			NetworkID:             networkID.String(),
			Email:                 req.GetEmail(),
			PhoneNumber:           req.GetPhoneNumber(),
			Gender:                req.GetGender(),
			Address:               req.GetAddress(),
			ProofOfIdentification: req.GetProofOfIdentification(),
			DateOfBirth:           birthday.GoString(),
			IdSerial:              req.GetIdSerial()},
	}, nil

}

func (s *SubcriberServer) Delete(ctx context.Context, req *pb.DeleteSubscriberRequest) (*pb.DeleteSubscriberResponse, error) {
	subscriberIdReq := req.GetSubscriberID()
	if subscriberIdReq == "" {
		return nil, status.Errorf(codes.InvalidArgument, "subscriberID must not be empty")
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}
	subscriberID, err := uuid.FromString(subscriberIdReq)

	logrus.Infof("Delete Subscriber : %v ", subscriberID)
	err = s.subscriberRepo.Delete(subscriberID)
	if err != nil {
		logrus.WithError(err).Error("error while deleting subscriber")
		return nil, grpc.SqlErrorToGrpc(err, "subscriber")
	}
	route := s.subscriberRoutingKey.SetAction("delete").SetObject("subscriber").MustBuild()
	err = s.msgbus.PublishRequest(route, req)
	if err != nil {
		logrus.Errorf("Failed to publish message %+v with key %+v. Errors %s", req, route, err.Error())
	}
	if err != nil {
		logrus.Errorf("Failed to publish message %+v with key %+v. Errors %s", req, route, err.Error())
	}
	return &pb.DeleteSubscriberResponse{
		SubscriberID: subscriberID.String(),
	}, nil
}

func (s *SubcriberServer) Get(ctx context.Context, req *pb.GetSubscriberRequest) (*pb.GetSubscriberResponse, error) {

	subscriberIdReq := req.GetSubscriberID()
	if subscriberIdReq == "" {
		return nil, status.Errorf(codes.InvalidArgument, "subscriberID must not be empty")
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	subscriberID, err := uuid.FromString(subscriberIdReq)
	logrus.Infof("GetSubscriber : %v ", subscriberID)
	subscriber, err := s.subscriberRepo.Get(subscriberID)

	if err != nil {
		logrus.Error("error getting a subscriber" + err.Error())

		return nil, grpc.SqlErrorToGrpc(err, "subscriber")
	}

	resp := &pb.GetSubscriberResponse{Subscriber: dbSubscriberToPbSubscribers(subscriber)}

	return resp, nil

}

func (s *SubcriberServer) ListSubscribers(ctx context.Context, req *pb.ListSubscribersRequest) (*pb.ListSubscribersResponse, error) {

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	logrus.Infof("List all subscribers")

	subscribers, err := s.subscriberRepo.ListSubscribers()
	if err != nil {
		logrus.WithError(err).Error("error while getting all subscribers")
		return nil, grpc.SqlErrorToGrpc(err, "subscribers")
	}

	subscriberList := &pb.ListSubscribersResponse{
		Subscribers: dbsubscriberToPbSubscribers(subscribers),
	}

	return subscriberList, nil
}
func (s *SubcriberServer) GetByNetwork(ctx context.Context, req *pb.GetByNetworkRequest) (*pb.GetByNetworkResponse, error) {
	networkIdReq := req.GetNetworkID()
	if networkIdReq == "" {
		return nil, status.Errorf(codes.InvalidArgument, "networkID must not be empty")
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}

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

	subscriberList := &pb.GetByNetworkResponse{
		Subscribers: dbsubscriberToPbSubscribers(subscribers),
	}

	return subscriberList, nil
}
func (s *SubcriberServer) Update(ctx context.Context, req *pb.UpdateSubscriberRequest) (*pb.UpdateSubscriberResponse, error) {
	subscriberIdReq := req.GetSubscriberID()
	if subscriberIdReq == "" {
		return nil, status.Errorf(codes.InvalidArgument, "subscriberID must not be empty")
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}
	subscribeID, err := uuid.FromString(subscriberIdReq)

	logrus.Infof("Update Subscriber Id: %v, Email: %v, ProofOfIdentification: %v, Address: %v",
		req.SubscriberID, req.Email, req.ProofOfIdentification, req.Address)

	updateSubscriber := db.Subscriber{
		Email:                 req.GetEmail(),
		PhoneNumber:           req.GetPhoneNumber(),
		Address:               req.GetAddress(),
		ProofOfIdentification: req.GetProofOfIdentification(),
		IdSerial:              req.GetIdSerial(),
	}

	updatedSubscriberRes, err := s.subscriberRepo.Update(subscribeID, updateSubscriber)
	if err != nil {
		logrus.Errorf("error while updating a subscriber: %v", err)
		return nil, grpc.SqlErrorToGrpc(err, "subscriber")
	}

	return &pb.UpdateSubscriberResponse{
		Email:                 updatedSubscriberRes.Email,
		PhoneNumber:           updatedSubscriberRes.PhoneNumber,
		Address:               updatedSubscriberRes.Address,
		IdSerial:              updatedSubscriberRes.IdSerial,
		ProofOfIdentification: updatedSubscriberRes.ProofOfIdentification,
	}, nil
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
		OrgID: s.OrgID.String(),
	}

}
