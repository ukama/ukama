package db

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/tj/assert"
	"github.com/ukama/ukama/systems/subscriber/subscriber-registry/pb/gen"
	"github.com/ukama/ukama/systems/subscriber/subscriber-registry/pb/gen/mocks"
)

func TestSubscriberRepo_Get(t *testing.T) {
	subscriberServiceServerMock := &mocks.SubscriberRegistryServiceServer{}

	t.Run("SubscriberExist", func(t *testing.T) {

		// Set up mock database
		subscriberServiceServerMock.On("Get", mock.Anything, &gen.GetSubscriberRequest{
			SubscriberID: "12345",
		}).Return(&gen.GetSubscriberResponse{
			Subscriber: &gen.Subscriber{
				SubscriberID: "12345",
				FirstName:    "John",
				LastName:     "Doe",
			},
		}, nil)
		expectedResponse := &gen.GetSubscriberResponse{
			Subscriber: &gen.Subscriber{
				SubscriberID: "12345",
				FirstName:    "John",
				LastName:     "Doe",
			},
		}

		response, err := subscriberServiceServerMock.Get(context.Background(), &gen.GetSubscriberRequest{
			SubscriberID: "12345",
		})
		assert.Equal(t, expectedResponse, response)
		assert.Nil(t, err)
		subscriberServiceServerMock.AssertExpectations(t)
	})
	t.Run("SubscriberDoesNotExist", func(t *testing.T) {
		subscriberServiceServerMock.On("Get", mock.Anything, &gen.GetSubscriberRequest{
			SubscriberID: "123456",
		}).Return(nil, errors.New("subscriber not found"))
		response, err := subscriberServiceServerMock.Get(context.Background(), &gen.GetSubscriberRequest{
			SubscriberID: "123456",
		})

		assert.Nil(t, response)
		assert.EqualError(t, err, "subscriber not found")
		subscriberServiceServerMock.AssertExpectations(t)

	})
}

func TestSubscriberRepo_Add(t *testing.T) {
	subscriberServiceServerMock := &mocks.SubscriberRegistryServiceServer{}
	t.Run("Add a Subscriber", func(t *testing.T) {

		subscriberServiceServerMock.On("Add", mock.Anything, &gen.AddSubscriberRequest{
			FirstName: "John",
			LastName:  "Doe",
		}).Return(&gen.AddSubscriberResponse{
			Subscriber: &gen.Subscriber{
				FirstName:    "John",
				LastName:     "Doe",
				SubscriberID: "12345",
			},
		}, nil)
		response, err := subscriberServiceServerMock.Add(context.Background(), &gen.AddSubscriberRequest{
			FirstName: "John",
			LastName:  "Doe",
		})

		expectedResponse := &gen.AddSubscriberResponse{
			Subscriber: &gen.Subscriber{
				FirstName:    "John",
				LastName:     "Doe",
				SubscriberID: "12345",
			},
		}

		assert.Equal(t, expectedResponse, response)
		assert.Nil(t, err)
	})
}

func TestSubscriberRepo_Update(t *testing.T) {
	subscriberServiceServerMock := &mocks.SubscriberRegistryServiceServer{}
	t.Run("SubscriberExist", func(t *testing.T) {

		subscriberServiceServerMock.On("Update", mock.Anything, &gen.UpdateSubscriberRequest{
			SubscriberID:          "12345",
			IdSerial:              "12345",
			Email:                 "john@gmail.com",
			Address:               "kigali",
			ProofOfIdentification: "0909099",
		}).Return(&gen.UpdateSubscriberResponse{
			IdSerial:              "12345",
			ProofOfIdentification: "0909099",
			Email:                 "john@gmail.com",
			Address:               "kigali",
		}, nil)

		response, err := subscriberServiceServerMock.Update(context.Background(), &gen.UpdateSubscriberRequest{
			SubscriberID: "12345",
			Email:        "john@gmail.com",
			Address:      "kigali",
		})
		expectedResponse := &gen.UpdateSubscriberResponse{
			IdSerial:              "12345",
			ProofOfIdentification: "0909099",
			Email:                 "john@gmail.com",
			Address:               "kigali",
		}

		assert.Equal(t, expectedResponse, response)
		assert.Nil(t, err)

	})
}

func TestSubscriberRepo_Delete(t *testing.T) {
	subscriberServiceServerMock := &mocks.SubscriberRegistryServiceServer{}

	t.Run("SubscriberExist", func(t *testing.T) {

		// Set up mock database
		subscriberServiceServerMock.On("Delete", mock.Anything, &gen.DeleteSubscriberRequest{
			SubscriberID: "12345",
		}).Return(&gen.DeleteSubscriberResponse{}, nil)
		expectedResponse := &gen.DeleteSubscriberResponse{}

		response, err := subscriberServiceServerMock.Delete(context.Background(), &gen.DeleteSubscriberRequest{
			SubscriberID: "12345",
		})
		assert.Equal(t, expectedResponse, response)
		assert.Nil(t, err)
		subscriberServiceServerMock.AssertExpectations(t)
	})
	t.Run("SubscriberAlreadyDeleted", func(t *testing.T) {
		subscriberServiceServerMock.On("delete", mock.Anything, &gen.DeleteSubscriberRequest{
			SubscriberID: "123456",
		}).Return(nil, errors.New("subscriber not found"))
		response, err := subscriberServiceServerMock.Delete(context.Background(), &gen.DeleteSubscriberRequest{
			SubscriberID: "123456",
		})

		assert.Nil(t, response)
		assert.EqualError(t, err, "subscriber not found")
		subscriberServiceServerMock.AssertExpectations(t)

	})
}
