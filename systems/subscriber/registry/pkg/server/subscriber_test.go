package server

import (
	"context"
	"testing"

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

		s := NewSubscriberServer(subscriberRepo, msgBus, simManagerService, network)

		req := &pb.AddSubscriberRequest{
			OrgID:               "7e82c8b1-a746-4f2c-a80e-f4d14d863ea3",
			FirstName:           "John",
			LastName:            "Doe",
			NetworkID:           "9e82c8b1-a746-4f2c-a80e-f4d14d863ea3",
			Email:               "johndoe@example.com",
			PhoneNumber:         "1234567890",
			Gender:              "Male",
			DateOfBirth:"16-04-1995",
			Address:             "1 Main St",
			ProofOfIdentification: "Passport",
			IdSerial: "123456789",
		}

		resp, err := s.Add(context.Background(), req)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "7e82c8b1-a746-4f2c-a80e-f4d14d863ea3", resp.Subscriber.OrgID)
		assert.Equal(t, "John", resp.Subscriber.FirstName)
		assert.Equal(t, "Doe", resp.Subscriber.LastName)
		assert.Equal(t, "9e82c8b1-a746-4f2c-a80e-f4d14d863ea3", resp.Subscriber.NetworkID)
		assert.Equal(t, "johndoe@example.com", resp.Subscriber.Email)
		assert.Equal(t, "1234567890", resp.Subscriber.PhoneNumber)
		assert.Equal(t, "Male", resp.Subscriber.Gender)
	})}
	func TestSubscriberServer_Get(t *testing.T) {
	
		t.Run("SubscriberNotFound", func(t *testing.T) {
			var subscriberID = uuid.Nil
	
			subRepo := &mocks.SubscriberRepo{}
	
			subRepo.On("Get", subscriberID).Return(nil, gorm.ErrRecordNotFound).Once()

			s := NewSubscriberServer(subRepo, nil, nil, nil)
			subResp, err := s.Get(context.TODO(), &pb.GetSubscriberRequest{
				SubscriberID: subscriberID.String()})
	
			assert.Error(t, err)
			assert.Nil(t, subResp)
			subRepo.AssertExpectations(t)
		})
	
		t.Run("SubscriberUUIDInvalid", func(t *testing.T) {
			var subscriberID = "1"
	
			subRepo := &mocks.SubscriberRepo{}
	
			s := NewSubscriberServer(subRepo, nil, nil, nil)
			subResp, err := s.Get(context.TODO(), &pb.GetSubscriberRequest{
				SubscriberID: subscriberID})
	
			assert.Error(t, err)
			assert.Nil(t, subResp)
			subRepo.AssertExpectations(t)
		})
	}
	func TestSubscriberServer_GetbyNetwork(t *testing.T) {
		
		
		  
		t.Run("NetworkNotFound", func(t *testing.T) {
			var networkID = uuid.Nil
	
			subRepo := &mocks.SubscriberRepo{}
	
			subRepo.On("GetByNetwork", networkID).Return(nil, gorm.ErrRecordNotFound).Once()
	
			s := NewSubscriberServer(subRepo, nil, nil, nil)
			subResp, err := s.GetByNetwork(context.TODO(), &pb.GetByNetworkRequest{
				NetworkID: networkID.String()})
	
			assert.Error(t, err)
			assert.Nil(t, subResp)
			subRepo.AssertExpectations(t)
		})
	
		t.Run("NetworkUUIDInvalid", func(t *testing.T) {
			var networkID = "1"
	
			subRepo := &mocks.SubscriberRepo{}
	
			s := NewSubscriberServer(subRepo, nil, nil, nil)
			subResp, err := s.GetByNetwork(context.TODO(), &pb.GetByNetworkRequest{
				NetworkID: networkID})
	
			assert.Error(t, err)
			assert.Nil(t, subResp)
			subRepo.AssertExpectations(t)
		})
	}
