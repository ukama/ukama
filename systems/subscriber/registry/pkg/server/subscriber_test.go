package server

import (
	"context"
	"testing"
	"time"

	mbmocks "github.com/ukama/ukama/systems/common/mocks"
	uuid "github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/subscriber/registry/mocks"
	pb "github.com/ukama/ukama/systems/subscriber/registry/pb/gen"
	"gorm.io/gorm"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAdd(t *testing.T) {
	// Test case 1: Add subscriber successfully
	t.Run("Add subscriber successfully", func(t *testing.T) {
		subscriberRepo := &mocks.SubscriberRepo{}
		msgBus := &mbmocks.MsgBusServiceClient{}
		simManagerService := &mocks.SimManagerClientProvider{}
		network := &mocks.NetworkInfoClient{}

		subscriberRepo.On("Add", mock.AnythingOfType("*db.Subscriber")).Return(nil)
		network.On("ValidateNetwork", mock.Anything, mock.Anything).Return(nil)

		msgBus.On("PublishRequest", mock.Anything, mock.Anything).Return(nil)

		s := NewSubscriberServer(subscriberRepo, msgBus, simManagerService, network)

		req := &pb.AddSubscriberRequest{
			OrgId:                 "7e82c8b1-a746-4f2c-a80e-f4d14d863ea3",
			FirstName:             "John",
			LastName:              "Doe",
			NetworkId:             "9e82c8b1-a746-4f2c-a80e-f4d14d863ea3",
			Email:                 "johndoe@example.com",
			PhoneNumber:           "1234567890",
			Gender:                "Male",
			Dob:                   time.Now().Add(time.Hour * 24 * 365 * 18).Format(time.RFC1123),
			Address:               "1 Main St",
			ProofOfIdentification: "Passport",
			IdSerial:              "123456789",
		}

		resp, err := s.Add(context.Background(), req)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "7e82c8b1-a746-4f2c-a80e-f4d14d863ea3", resp.Subscriber.OrgId)
		assert.Equal(t, "John", resp.Subscriber.FirstName)
		assert.Equal(t, "Doe", resp.Subscriber.LastName)
		assert.Equal(t, "9e82c8b1-a746-4f2c-a80e-f4d14d863ea3", resp.Subscriber.NetworkId)
		assert.Equal(t, "johndoe@example.com", resp.Subscriber.Email)
		assert.Equal(t, "1234567890", resp.Subscriber.PhoneNumber)
		assert.Equal(t, "Male", resp.Subscriber.Gender)
	})
}
func TestSubscriberServer_Get(t *testing.T) {

	t.Run("SubscriberNotFound", func(t *testing.T) {
		var subscriberId = uuid.Nil

		subRepo := &mocks.SubscriberRepo{}

		subRepo.On("Get", subscriberId).Return(nil, gorm.ErrRecordNotFound).Once()

		s := NewSubscriberServer(subRepo, nil, nil, nil)
		subResp, err := s.Get(context.TODO(), &pb.GetSubscriberRequest{
			SubscriberId: subscriberId.String()})

		assert.Error(t, err)
		assert.Nil(t, subResp)
		subRepo.AssertExpectations(t)
	})

	t.Run("SubscriberUUIDInvalid", func(t *testing.T) {
		var subscriberId = "1"

		subRepo := &mocks.SubscriberRepo{}

		s := NewSubscriberServer(subRepo, nil, nil, nil)
		subResp, err := s.Get(context.TODO(), &pb.GetSubscriberRequest{
			SubscriberId: subscriberId})

		assert.Error(t, err)
		assert.Nil(t, subResp)
		subRepo.AssertExpectations(t)
	})
}
func TestSubscriberServer_GetbyNetwork(t *testing.T) {

	t.Run("NetworkNotFound", func(t *testing.T) {
		var networkId = uuid.Nil

		subRepo := &mocks.SubscriberRepo{}

		subRepo.On("GetByNetwork", networkId).Return(nil, gorm.ErrRecordNotFound).Once()

		s := NewSubscriberServer(subRepo, nil, nil, nil)
		subResp, err := s.GetByNetwork(context.TODO(), &pb.GetByNetworkRequest{
			NetworkId: networkId.String()})

		assert.Error(t, err)
		assert.Nil(t, subResp)
		subRepo.AssertExpectations(t)
	})

	t.Run("NetworkUUIDInvalid", func(t *testing.T) {
		var networkId = "1"

		subRepo := &mocks.SubscriberRepo{}

		s := NewSubscriberServer(subRepo, nil, nil, nil)
		subResp, err := s.GetByNetwork(context.TODO(), &pb.GetByNetworkRequest{
			NetworkId: networkId})

		assert.Error(t, err)
		assert.Nil(t, subResp)
		subRepo.AssertExpectations(t)
	})
}
