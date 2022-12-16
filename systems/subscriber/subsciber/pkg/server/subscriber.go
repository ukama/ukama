package server

import (
	"context"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"github.com/golang/protobuf/ptypes"
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
	  logrus.Infof("Adding subscriber: %v", req)
	      uuid, err := uuid.NewV4()
    if err != nil {
        logrus.Errorf("Failed to generate UUID: %s", err)
        return nil, err
    }
	    dateOfBirth, err := ptypes.Timestamp(req.GetDateOfBith())
    if err != nil {
        logrus.Errorf("Error converting timestamp: %s", err)
        return nil, err
    }
   
	subscriber := &db.Subscriber{
		SubscriberID: uuid,
		FullName:     req.GetFullName(),
		Email:        req.GetEmail(),
		PhoneNumber:  req.GetPhoneNumber(),
		Address:      req.GetAddress(),
		ProofOfIdentification:req.GetProofOfIdentification(),
		DOB:  &dateOfBirth,
		IdSerial:req.GetIdSerial(),
	}
	err = s.subscriberRepo.Add(subscriber)
	if err != nil {
		logrus.Error("error while adding subscriber" + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "subscriber")
	}
	return &pb.AddSubscriberResponse{Subscriberid: uuid.String()}, nil

}
func (s *SunscriberServer) Delete(ctx context.Context, req *pb.DeleteSubscriberRequest) (*pb.DeleteSubscriberResponse, error) {
	subscriberID := req.GetSubscriberid()
	logrus.Infof("Delete Subscriber : %v ", subscriberID)
	err := s.subscriberRepo.Delete(subscriberID)
	if err != nil {
		logrus.Error("error while deleting subscriber" + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "subscriber")
	}
	return &pb.DeleteSubscriberResponse{}, nil

}

func (s *SunscriberServer) Get(ctx context.Context, req *pb.GetSubscriberRequest) (*pb.GetSubscriberResponse, error) {
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
func dbsubscriberToPbSubscribers(packages []db.Subscriber) []*pb.Subscriber {
	res := []*pb.Subscriber{}
	for _, u := range packages {
		res = append(res, dbSubscriberToPbSubscribers(&u))
	}
	return res
}

func dbSubscriberToPbSubscribers(s *db.Subscriber) *pb.Subscriber {
	time, err := time.Parse("2006-01-02", s.DOB.String())
if err != nil {
	logrus.Error("error while parsing date of birth" + err.Error())
}
timestamp, err := ptypes.TimestampProto(time)
if err != nil {
	logrus.Error("error while parsing date of birth" + err.Error())
}
	pbSims := make([]*pb.Sims, len(s.Sims))
	for i, sim := range s.Sims {
		pbSims[i] = &pb.Sims{
			SimId:     sim.SimID.String(),
			Iccid:     sim.ICCID,
			Imsi:      sim.IMSI,
			Msisdn:    sim.MSISDN,
			CreatedAt: sim.CreatedAt.String(),
			UpdatedAt: sim.UpdatedAt.String(),
			DeletedAt: sim.DeletedAt.Time.String(),
		}
	}
fmt.Println("SIMS",pbSims)
	return &pb.Subscriber{
		Id:             uint64(s.ID),
		FullName:           s.FullName,
		Email:          s.Email,
		SubscriberId:   s.SubscriberID.String(),
		ProofOfIdentification: s.ProofOfIdentification,
		PhoneNumber:    s.PhoneNumber,
		IdSerial:s.IdSerial,
		Sims:           pbSims,
		Address:        s.Address,
		CreatedAt:      s.CreatedAt.String(),
		UpdatedAt:      s.UpdatedAt.String(),
		DeletedAt:      s.DeletedAt.Time.String(),
		DateOfBith:	 timestamp,
	}
}
