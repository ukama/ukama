package server

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/grpc"
	pb "github.com/ukama/ukama/systems/subscriber/subscriber/pb/gen"
	"github.com/ukama/ukama/systems/subscriber/subscriber/pkg/db"
)

type SubcriberServer struct {
	subscriberRepo db.SubscriberRepo
	pb.UnimplementedSubscriberServiceServer
}
type Date struct {
	Year  int32
	Month int32
	Day   int32
  }

func NewSubscriberServer(subscriberRepo db.SubscriberRepo) *SubcriberServer {
	return &SubcriberServer{subscriberRepo: subscriberRepo}
}

func (s *SubcriberServer) Add(ctx context.Context, req *pb.AddSubscriberRequest) (*pb.AddSubscriberResponse, error) {
	logrus.Infof("Adding subscriber: %v", req)
	networkID_uuid := uuid.FromStringOrNil(req.GetNetworkID())
	uuid, err := uuid.NewV4()
	if err != nil {
		logrus.Errorf("Failed to generate UUID: %s", err)
		return nil, err
	}
	timestamp := &timestamppb.Timestamp{
		Seconds: req.DateOfBirth.Seconds,
		Nanos:   req.DateOfBirth.Nanos,
	}

	birthday := timestamp.AsTime()
	

	subscriber := &db.Subscriber{
		SubscriberID:          uuid,
		FirstName:             req.GetFirstName(),
		LastName:              req.GetLastName(),
		NetworkID:             networkID_uuid,
		Email:                 req.GetEmail(),
		PhoneNumber:           req.GetPhoneNumber(),
		Gender:                req.GetGender(),
		Address:               req.GetAddress(),
		ProofOfIdentification: req.GetProofOfIdentification(),
		DOB:             birthday ,
		IdSerial:              req.GetIdSerial(),
	}
	err = s.subscriberRepo.Add(subscriber)
	if err != nil {
		logrus.Error("error while adding subscriber" + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "subscriber")
	}
	return &pb.AddSubscriberResponse{SubscriberID: uuid.String()}, nil

}
func (s *SubcriberServer) Delete(ctx context.Context, req *pb.DeleteSubscriberRequest) (*pb.DeleteSubscriberResponse, error) {
	subscriberID := req.GetSubscriberID()
	logrus.Infof("Delete Subscriber : %v ", subscriberID)
	err := s.subscriberRepo.Delete(subscriberID)
	if err != nil {
		logrus.Error("error while deleting subscriber" + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "subscriber")
	}
	return &pb.DeleteSubscriberResponse{}, nil

}

func (s *SubcriberServer) Get(ctx context.Context, req *pb.GetSubscriberRequest) (*pb.GetSubscriberResponse, error) {
	subscriberID := req.GetSubscriberID()

	logrus.Infof("GetSubscriber : %v ", subscriberID)
	subscriber, err := s.subscriberRepo.Get(subscriberID)

	if err != nil {
		logrus.Error("error getting a subscriber" + err.Error())

		return nil, grpc.SqlErrorToGrpc(err, "subscriber")
	}

	resp := &pb.GetSubscriberResponse{Subscriber: dbSubscriberToPbSubscribers(subscriber)}

	return resp, nil

}
func (s *SubcriberServer) GetByNetwork(ctx context.Context, req *pb.GetByNetworkRequest) (*pb.GetByNetworkResponse, error) {
	logrus.Infof("Get subscribers by network: %v ", req.GetNetworkID())
	networkID_uuid := uuid.FromStringOrNil(req.GetNetworkID())

	subscribers, err := s.subscriberRepo.GetByNetwork(networkID_uuid)
	if err != nil {
		logrus.Error("error while getting subscribers by network" + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "subscribers")
	}

	subscriberList := &pb.GetByNetworkResponse{
		Subscribers: dbsubscriberToPbSubscribers(subscribers),
	}

	return subscriberList, nil
}
func (s *SubcriberServer) Update(ctx context.Context, req *pb.UpdateSubscriberRequest) (*pb.UpdateSubscriberResponse, error) {
	logrus.Infof("Update Subscriber Id: %v, FullName: %v, Email: %v, Email: %v, ProofOfIdentification: %v, Address: %v, DateOfBith: %v",
		req.Email, req.ProofOfIdentification, req.Address)

	subscriber := db.Subscriber{
		Email:                 req.GetEmail(),
		PhoneNumber:           req.GetPhoneNumber(),
		Address:               req.GetAddress(),
		ProofOfIdentification: req.GetProofOfIdentification(),
		IdSerial:              req.GetIdSerial(),
	}

	sub, err := s.subscriberRepo.Update(req.GetSubscriberID(), subscriber)
	if err != nil {
		logrus.Error("error while updating a subscriber" + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "subscriber")
	}

	return &pb.UpdateSubscriberResponse{
		Subscriber: dbSubscriberToPbSubscribers(sub),
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
	fmt.Println("DATE",s.DOB)
	dateString := s.DOB.Format("2006-01-02")
	fmt.Println("DATE STRING",dateString)

	return &pb.Subscriber{
		FirstName:             s.FirstName,
		LastName:              s.LastName,
		Email:                 s.Email,
		SubscriberID:          s.SubscriberID.String(),
		ProofOfIdentification: s.ProofOfIdentification,
		PhoneNumber:           s.PhoneNumber,
		IdSerial:              s.IdSerial,
		NetworkID: s.NetworkID.String(),
		Gender:                s.Gender,
		Address:               s.Address,
		CreatedAt:              s.CreatedAt.String(),
		UpdatedAt:             s.UpdatedAt.String(),
		DateOfBirth:           dateString,
	}

}
