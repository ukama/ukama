package server

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/grpc"
	pb "github.com/ukama/ukama/systems/subscriber/subscriber/pb/gen"
	"github.com/ukama/ukama/systems/subscriber/subscriber/pkg/db"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SubcriberServer struct {
	subscriberRepo db.SubscriberRepo
	pb.UnimplementedSubscriberServiceServer
}

func NewSubscriberServer(subscriberRepo db.SubscriberRepo) *SubcriberServer {
	return &SubcriberServer{subscriberRepo: subscriberRepo}
}

func (s *SubcriberServer) Add(ctx context.Context, req *pb.AddSubscriberRequest) (*pb.AddSubscriberResponse, error) {
	logrus.Infof("Adding subscriber: %v", req)
	uuid, err := uuid.NewV4()
	if err != nil {
		logrus.Errorf("Failed to generate UUID: %s", err)
		return nil, err
	}
	parsedDoB, err := time.Parse("2006-1-2", req.GetDateOfBith())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid date format!")
	}

	subscriber := &db.Subscriber{
		SubscriberID:          uuid,
		FirstName:             req.GetFirstName(),
		LastName:              req.GetLastName(),
		Email:                 req.GetEmail(),
		PhoneNumber:           req.GetPhoneNumber(),
		Gender:                req.GetGender(),
		Address:               req.GetAddress(),
		ProofOfIdentification: req.GetProofOfIdentification(),
		DOB:                   parsedDoB.Format("2006-01-02"),
		IdSerial:              req.GetIdSerial(),
	}
	err = s.subscriberRepo.Add(subscriber)
	if err != nil {
		logrus.Error("error while adding subscriber" + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "subscriber")
	}
	return &pb.AddSubscriberResponse{Subscriberid: uuid.String()}, nil

}
func (s *SubcriberServer) Delete(ctx context.Context, req *pb.DeleteSubscriberRequest) (*pb.DeleteSubscriberResponse, error) {
	subscriberID := req.GetSubscriberid()
	logrus.Infof("Delete Subscriber : %v ", subscriberID)
	err := s.subscriberRepo.Delete(subscriberID)
	if err != nil {
		logrus.Error("error while deleting subscriber" + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "subscriber")
	}
	return &pb.DeleteSubscriberResponse{}, nil

}

func (s *SubcriberServer) Get(ctx context.Context, req *pb.GetSubscriberRequest) (*pb.GetSubscriberResponse, error) {
	subscriberID := req.GetSubscriberid()

	logrus.Infof("GetSubscriber : %v ", subscriberID)
	subscriber, err := s.subscriberRepo.Get(subscriberID)

	if err != nil {
		logrus.Error("error getting a subscriber" + err.Error())

		return nil, grpc.SqlErrorToGrpc(err, "subscriber")
	}

	resp := &pb.GetSubscriberResponse{Subscriber: dbSubscriberToPbSubscribers(subscriber)}

	return resp, nil

}
func (s *SubcriberServer) Update(ctx context.Context, req *pb.UpdateSubscriberRequest) (*pb.UpdateSubscriberResponse, error) {
	logrus.Infof("Update Subscriber Id: %v, FullName: %v, Email: %v, Email: %v, ProofOfIdentification: %v, Address: %v, DateOfBith: %v",
		req.FirstName, req.Email, req.Email, req.ProofOfIdentification, req.Address, req.DateOfBith)
	subscriber := db.Subscriber{
		FirstName:             req.GetFirstName(),
		LastName:              req.GetLastName(),
		Email:                 req.GetEmail(),
		PhoneNumber:           req.GetPhoneNumber(),
		SimID:                 req.GetSimRef(),
		Address:               req.GetAddress(),
		ProofOfIdentification: req.GetProofOfIdentification(),
		DOB:                   req.GetDateOfBith(),
		IdSerial:              req.GetIdSerial(),
	}

	sub, err := s.subscriberRepo.Update(req.GetSubscriberId(), subscriber)
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

	return &pb.Subscriber{
		Id:                    uint64(s.ID),
		FirstName:             s.FirstName,
		LastName:              s.LastName,
		Email:                 s.Email,
		SubscriberId:          s.SubscriberID.String(),
		ProofOfIdentification: s.ProofOfIdentification,
		PhoneNumber:           s.PhoneNumber,
		IdSerial:              s.IdSerial,
		Gender:                s.Gender,
		Address:               s.Address,
		CreatedAt:             s.CreatedAt.String(),
		UpdatedAt:             s.UpdatedAt.String(),
		DeletedAt:             s.DeletedAt.Time.String(),
		DateOfBith:            s.DOB,
	}
}
