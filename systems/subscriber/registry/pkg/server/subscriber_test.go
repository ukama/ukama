package server_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	mbmocks "github.com/ukama/ukama/systems/common/mocks"
	"google.golang.org/protobuf/types/known/timestamppb"

	uuid "github.com/ukama/ukama/systems/common/uuid"

	mocks "github.com/ukama/ukama/systems/subscriber/registry/mocks"
	pb "github.com/ukama/ukama/systems/subscriber/registry/pb/gen"
	"github.com/ukama/ukama/systems/subscriber/registry/pkg/db"
	"github.com/ukama/ukama/systems/subscriber/registry/pkg/server"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSubscriberServer_Add(t *testing.T) {
	subRepo := &mocks.SubscriberRepo{}
	netRepo := &mocks.Network{}
	simRepo := &mocks.SimManagerClientProvider{}

	msgbusClient := &mbmocks.MsgBusServiceClient{}
	netID := uuid.NewV4()
	orgID := uuid.NewV4()
	mockDate := &timestamppb.Timestamp{
		Seconds: time.Now().Unix(),
		Nanos:   0,
	}

	sub := &db.Subscriber{
		SubscriberID:          uuid.NewV4(),
		FirstName:             "john",
		LastName:              "Doe",
		Email:                 "john@gmail.com",
		NetworkID:             netID,
		OrgID:                 orgID,
		PhoneNumber:           "0791240041",
		Gender:                "male",
		IdSerial:              "00000",
		DOB:                   mockDate.AsTime(),
		ProofOfIdentification: "passwport",
		Address:               "kigali",
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
		DeletedAt:             nil,
	}

	psub := &pb.AddSubscriberRequest{FirstName: "john", LastName: "Doe", Email: "john@gmail.com", NetworkID: netID.String(), OrgID: orgID.String(), PhoneNumber: "0791240041", Gender: "male", IdSerial: "00000", DateOfBirth: mockDate, ProofOfIdentification: "passwport", Address: "kigali"}

	subRepo.On("Add", sub).Return(nil).Once()
	msgbusClient.On("PublishRequest", mock.Anything, psub).Return(nil).Once()

	s := server.NewSubscriberServer(subRepo, msgbusClient, nil, netRepo)
	simRepo.On("GetClient").Return(nil, nil)

	 netRepo.On("ValidateNetwork", netID.String(), orgID.String()).Return(nil).Once()
	_, err := s.Add(context.TODO(), psub)
	assert.NoError(t, err)
	subRepo.AssertExpectations(t)
}
func TestSubscriberServer_Update(t *testing.T) {
	subRepo := &mocks.SubscriberRepo{}
	msgbusClient := &mbmocks.MsgBusServiceClient{}
	sub := &db.Subscriber{
		SubscriberID: uuid.NewV4(),
		FirstName:    "john",
		LastName:     "Doe",
		NetworkID:    uuid.NewV4(),
		Email:        "john@gmail.com",
		PhoneNumber:  "0791240041",
		Gender:       "male",
	}

	psub := &pb.UpdateSubscriberRequest{Email: "john@gmail.com", PhoneNumber: "0791240041"}

	subRepo.On("GetByNetwork", sub.NetworkID).Return(sub, nil).Once()
	subRepo.On("Update", sub).Return(nil).Once()
	msgbusClient.On("PublishRequest", mock.Anything, psub).Return(nil).Once()

	s := server.NewSubscriberServer(subRepo, msgbusClient, nil, nil)

	_, err := s.Update(context.TODO(), psub)

	assert.NoError(t, err)
	subRepo.AssertExpectations(t)
}

func TestSubscriberServer_Get(t *testing.T) {
	subRepo := &mocks.SubscriberRepo{}
	msgbusClient := &mbmocks.MsgBusServiceClient{}
	subscriberUUID, err := uuid.FromString("dbe9a556-a626-11ed-afa1-0242ac120002")
	if err != nil {
		logrus.Errorf("Invalid format uuid %s", err.Error())
	}
	sub := &db.Subscriber{
		SubscriberID: subscriberUUID,
		FirstName:    "john",
		LastName:     "Doe",
		NetworkID:    uuid.NewV4(),
		Email:        "john@gmail.com",
		PhoneNumber:  "0791240041",
		Gender:       "male",
	}
	subRepo.On("Get", sub.SubscriberID).Return(sub, nil).Once()

	s := server.NewSubscriberServer(subRepo, msgbusClient, nil, nil)
	fmt.Println("BRACKLEY",s)

	resp, err := s.Get(context.TODO(), &pb.GetSubscriberRequest{SubscriberID: subscriberUUID.String()})
	fmt.Println("VANESSA",resp)

	assert.NoError(t, err)
	assert.Equal(t, sub.SubscriberID, resp.GetSubscriber())
	subRepo.AssertExpectations(t)
}
