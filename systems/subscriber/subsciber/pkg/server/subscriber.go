package server

import (
	"context"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/grpc"
	pb "github.com/ukama/ukama/systems/subscriber/subscriber/pb/gen"
	"github.com/ukama/ukama/systems/subscriber/subscriber/pkg/db"
)

type SunscriberServer struct {
	subscriberRepo db.SubscriberRepo
	pb.UnimplementedSubscriberServiceServer
}

func NewSubscriberServer(subscriberRepo db.SubscriberRepo) *SunscriberServer {
	return &SunscriberServer{subscriberRepo: subscriberRepo}
}
func (s *SunscriberServer) Add(ctx context.Context, req *pb.AddSubscriberRequest) (*pb.AddSubscriberResponse, error) {
	logrus.Infof("Add a subscriber : %v ")
	subId := uuid.New()
	subscriber := &db.Subscriber{
		SubscriberId: subId.String(),
		Name:         req.GetName(),
		Email:        req.GetEmail(),
		Phone:        req.GetPhone(),
		Address:      req.GetAddress(),
	}
	err := s.subscriberRepo.Add(subscriber)
	if err != nil {
		logrus.Error("error while adding subscriber" + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "subscriber")
	}
	return &pb.AddSubscriberResponse{Subscriberid: subId.String()}, nil

}
func (s *SunscriberServer) Delete(ctx context.Context, req *pb.DeleteSubscriberRequest) (*pb.DeleteSubscriberResponse, error) {
	logrus.Infof("Delete Subscriber : %v ", req.GetSubscriberid())
	err := s.subscriberRepo.Delete(req.GetSubscriberid())
	if err != nil {
		logrus.Error("error while deleting subscriber" + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "subscriber")
	}
	return &pb.DeleteSubscriberResponse{}, nil

}

func (s *SunscriberServer) Get(ctx context.Context, req *pb.GetSubscriberRequest) (*pb.GetSubscriberResponse, error) {
	logrus.Infof("GetSubscriber : %v ", req.GetSubscriberid())
	subscriber, err := s.subscriberRepo.Get(req.GetSubscriberid())

	if err != nil {
		logrus.Error("error getting a subscriber" + err.Error())

		return nil, grpc.SqlErrorToGrpc(err, "subscriber")
	}

	resp := &pb.GetSubscriberResponse{Subscriber: dbSubscriberToPbSubscribers(subscriber)}

	return resp, nil

}
func dbsubscriberToPbSubscribers(packages []db.Subscriber) []*pb.Subscriber {
	res := []*pb.Subscriber{}
	for _, u := range packages {
		res = append(res, dbSubscriberToPbSubscribers(&u))
	}
	return res
}

func dbSubscriberToPbSubscribers(s *db.Subscriber) *pb.Subscriber {
	return &pb.Subscriber{
		Id:           uint64(s.ID),
		Name:         s.Name,
		Email:        s.Email,
		SubscriberId: s.SubscriberId,
		Phone:        s.Phone,
		Address:      s.Address,
		CreatedAt:    s.CreatedAt.String(),
		UpdatedAt:    s.UpdatedAt.String(),
		DeletedAt:    s.DeletedAt.Time.String(),
	}
}
