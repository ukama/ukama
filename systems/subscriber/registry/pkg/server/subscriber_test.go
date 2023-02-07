package server_test

import (
	"context"
	"testing"

	"github.com/sirupsen/logrus"
	mbmocks "github.com/ukama/ukama/systems/common/mocks"

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

	psub := &pb.AddSubscriberRequest{FirstName: "john", LastName: "Doe", Email: "joe@gmail.com", PhoneNumber: "0791240041"}

	subRepo.On("Add", sub).Return(nil).Once()
	msgbusClient.On("PublishRequest", mock.Anything, psub).Return(nil).Once()

	s := server.NewSubscriberServer(subRepo, msgbusClient, nil, nil)
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
	subscriberUUID, error := uuid.FromString("dbe9a556-a626-11ed-afa1-0242ac120002")
	if error != nil {
		logrus.Error("Invalid format uuid %s ", error.Error())
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
	resp, err := s.Get(context.TODO(), &pb.GetSubscriberRequest{SubscriberID: subscriberUUID.String()})

	assert.NoError(t, err)
	assert.Equal(t, sub.SubscriberID, resp.GetSubscriber())
	subRepo.AssertExpectations(t)
}
