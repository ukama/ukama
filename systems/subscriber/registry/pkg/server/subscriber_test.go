package server

import (
	"context"
	"testing"

	mbmocks "github.com/ukama/ukama/systems/common/mocks"
	"github.com/ukama/ukama/systems/subscriber/registry/mocks"
	pb "github.com/ukama/ukama/systems/subscriber/registry/pb/gen"

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
