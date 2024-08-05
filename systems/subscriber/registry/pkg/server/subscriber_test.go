/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package server

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ukama/ukama/systems/subscriber/registry/mocks"
	"github.com/ukama/ukama/systems/subscriber/registry/pkg/db"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	cmocks "github.com/ukama/ukama/systems/common/mocks"
	creg "github.com/ukama/ukama/systems/common/rest/client/registry"
	"github.com/ukama/ukama/systems/common/ukama"
	uuid "github.com/ukama/ukama/systems/common/uuid"
	pb "github.com/ukama/ukama/systems/subscriber/registry/pb/gen"
	simMangerPb "github.com/ukama/ukama/systems/subscriber/sim-manager/pb/gen"
	mocksks "github.com/ukama/ukama/systems/subscriber/sim-manager/pb/gen/mocks"
)

const OrgName = "testorg"
const OrgId = "8c6c2bec-5f90-4fee-8ffd-ee6456abf4fc"


func TestAdd(t *testing.T) {
	t.Run("Add subscriber successfully", func(t *testing.T) {
		subscriberRepo := &mocks.SubscriberRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		simManagerService := &mocks.SimManagerClientProvider{}
		regClient := &cmocks.OrgClient{}
		networkClient := &cmocks.NetworkClient{}
		networkId := uuid.NewV4()
		sub := &db.Subscriber{
			FirstName:             "John",
			LastName:              "Doe",
			Email:                 "johndoe@example.com",
			PhoneNumber:           "1234567890",
			Gender:                "Male",
			Address:               "1 Main St",
			ProofOfIdentification: "Passport",
			IdSerial:              "123456789",
			NetworkId:             networkId,
			DOB:                   time.Now().Add(time.Hour * 24 * 365 * 18).Format(time.RFC3339),
		}

		subscriberRepo.On("Add", sub, mock.Anything).Return(nil).Once()
		networkClient.On("Get", networkId.String()).
			Return(&creg.NetworkInfo{
				Id:         networkId.String(),
				Name:       "net-1",
				SyncStatus: ukama.StatusTypeCompleted.String(),
			}, nil).Once()
		msgBus.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

		s := NewSubscriberServer(OrgName, subscriberRepo, msgBus, simManagerService, OrgId, regClient, networkClient)
		_, err := s.Add(context.TODO(), &pb.AddSubscriberRequest{
			FirstName:             sub.FirstName,
			LastName:              sub.LastName,
			Email:                 sub.Email,
			PhoneNumber:           sub.PhoneNumber,
			Gender:                sub.Gender,
			Dob:                   sub.DOB,
			NetworkId:             networkId.String(),
			Address:               sub.Address,
			ProofOfIdentification: sub.ProofOfIdentification,
			IdSerial:              sub.IdSerial,
		})
		assert.NoError(t, err)

		regClient.AssertExpectations(t)
		subscriberRepo.AssertExpectations(t)
		networkClient.AssertExpectations(t)
		msgBus.AssertExpectations(t)
	})

	t.Run("Add subscriber with default network", func(t *testing.T) {
		subscriberRepo := &mocks.SubscriberRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		simManagerService := &mocks.SimManagerClientProvider{}
		orgClient := &cmocks.OrgClient{}
		networkClient := &cmocks.NetworkClient{}
		networkId := uuid.NewV4()

		sub := &db.Subscriber{
			FirstName:             "John",
			LastName:              "Doe",
			Email:                 "johndoe@example.com",
			PhoneNumber:           "1234567890",
			Gender:                "Male",
			Address:               "1 Main St",
			ProofOfIdentification: "Passport",
			IdSerial:              "123456789",
			NetworkId:             networkId,
			DOB:                   time.Now().Add(time.Hour * 24 * 365 * 18).Format(time.RFC3339),
		}

		subscriberRepo.On("Add", sub, mock.Anything).Return(nil).Once()
		msgBus.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()
		networkClient.On("GetDefault", mock.Anything).Return(
			&creg.NetworkInfo{
				Id: networkId.String(),
			}, nil).Once()

		s := NewSubscriberServer(OrgName, subscriberRepo, msgBus, simManagerService, OrgId, orgClient, networkClient)

		res, err := s.Add(context.TODO(), &pb.AddSubscriberRequest{
			FirstName:             sub.FirstName,
			LastName:              sub.LastName,
			Email:                 sub.Email,
			PhoneNumber:           sub.PhoneNumber,
			Gender:                sub.Gender,
			Dob:                   sub.DOB,
			Address:               sub.Address,
			ProofOfIdentification: sub.ProofOfIdentification,
			IdSerial:              sub.IdSerial,
		})

		assert.NoError(t, err)
		assert.NotNil(t, res)
	})

	t.Run("Add subscriber with network client error", func(t *testing.T) {
		subscriberRepo := &mocks.SubscriberRepo{}
		networkClient := &cmocks.NetworkClient{}
		networkId := uuid.NewV4()

		networkClient.On("Get", networkId.String()).Return(nil, errors.New("network error")).Once()

		s := NewSubscriberServer(OrgName, subscriberRepo, nil, nil, OrgId, nil, networkClient)

		_, err := s.Add(context.TODO(), &pb.AddSubscriberRequest{
			Email:     "valid@email.com",
			NetworkId: networkId.String(),
		})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "network not found")
		networkClient.AssertExpectations(t)
	})

	
}

func TestGet(t *testing.T) {
	t.Run("Get subscriber successfully", func(t *testing.T) {
		subscriberId := uuid.NewV4()
		networkId := uuid.NewV4()
		email:="test@example.com"

		subscriberRepo := &mocks.SubscriberRepo{}
		simManagerService := &mocks.SimManagerClientProvider{}
		mockSimManagerClient := &mocksks.SimManagerServiceClient{}

		
		sub := subscriberRepo.On("Get", subscriberId).
		Return(&db.Subscriber{
			SubscriberId:subscriberId ,
			NetworkId:      networkId,
			Email:     email,
		}, nil).
		Once().
		ReturnArguments.Get(0).(*db.Subscriber)
		simManagerService.On("GetSimManagerService").Return(mockSimManagerClient, nil).Once()
		mockSimManagerClient.On("GetSimsBySubscriber", mock.Anything, &simMangerPb.GetSimsBySubscriberRequest{
			SubscriberId: subscriberId.String(),
		}).Return(&simMangerPb.GetSimsBySubscriberResponse{
			Sims: []*simMangerPb.Sim{},
		}, nil).Once()

		s := NewSubscriberServer(OrgName, subscriberRepo, nil, simManagerService, OrgId, nil, nil)
		resp, err := s.Get(context.TODO(), &pb.GetSubscriberRequest{
			SubscriberId: subscriberId.String(),
		})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, sub.Email, resp.Subscriber.Email)
		subscriberRepo.AssertExpectations(t)
		simManagerService.AssertExpectations(t)
		mockSimManagerClient.AssertExpectations(t)
	})
	t.Run("SusbscriberUUIDInvalid", func(t *testing.T) {
		var subID = "1"
		simManagerService := &mocks.SimManagerClientProvider{}

		subRepo := &mocks.SubscriberRepo{}

			s := NewSubscriberServer(OrgName, subRepo, nil, simManagerService, OrgId, nil, nil)

		subResp, err := s.Get(context.TODO(), &pb.GetSubscriberRequest{
			SubscriberId: subID})

		assert.Error(t, err)
		assert.Nil(t, subResp)
		subRepo.AssertExpectations(t)
	})

	t.Run("Get non-existent subscriber", func(t *testing.T) {
		subscriberId := uuid.NewV4()
		subscriberRepo := &mocks.SubscriberRepo{}

		subscriberRepo.On("Get", subscriberId).Return(nil, gorm.ErrRecordNotFound).Once()

		s := NewSubscriberServer(OrgName, subscriberRepo, nil, nil, OrgId, nil, nil)
		resp, err := s.Get(context.TODO(), &pb.GetSubscriberRequest{
			SubscriberId: subscriberId.String(),
		})

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.True(t, status.Code(err) == codes.NotFound)

		subscriberRepo.AssertExpectations(t)
	})
	t.Run("Database error", func(t *testing.T) {
        subscriberId := uuid.NewV4()
        subscriberRepo := &mocks.SubscriberRepo{}

        subscriberRepo.On("Get", subscriberId).Return(nil, errors.New("database error")).Once()

        s := NewSubscriberServer(OrgName, subscriberRepo, nil, nil, OrgId, nil, nil)
        resp, err := s.Get(context.TODO(), &pb.GetSubscriberRequest{
            SubscriberId: subscriberId.String(),
        })

        assert.Error(t, err)
        assert.Nil(t, resp)
        assert.Contains(t, err.Error(), "database error")

        subscriberRepo.AssertExpectations(t)
    })
	t.Run("Subscriber not found", func(t *testing.T) {
        subscriberId := uuid.NewV4()
        subscriberRepo := &mocks.SubscriberRepo{}

        subscriberRepo.On("Get", subscriberId).Return(nil, gorm.ErrRecordNotFound).Once()

        s := NewSubscriberServer(OrgName, subscriberRepo, nil, nil, OrgId, nil, nil)
        resp, err := s.Get(context.TODO(), &pb.GetSubscriberRequest{
            SubscriberId: subscriberId.String(),
        })

        assert.Error(t, err)
        assert.Nil(t, resp)
        assert.True(t, status.Code(err) == codes.NotFound)

        subscriberRepo.AssertExpectations(t)
    })
	t.Run("Get subscriber with invalid UUID", func(t *testing.T) {
		s := NewSubscriberServer(OrgName, nil, nil, nil, OrgId, nil, nil)
		resp, err := s.Get(context.TODO(), &pb.GetSubscriberRequest{
			SubscriberId: "invalid-uuid",
		})

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.True(t, status.Code(err) == codes.InvalidArgument)
	})
	t.Run("GetSimsBySubscriber error", func(t *testing.T) {
        subscriberId := uuid.NewV4()
        subscriber := &db.Subscriber{
            SubscriberId: subscriberId,
            Email:        "test@example.com",
        }

        subscriberRepo := &mocks.SubscriberRepo{}
        simManagerService := &mocks.SimManagerClientProvider{}
        mockSimManagerClient := &mocksks.SimManagerServiceClient{}

        subscriberRepo.On("Get", subscriberId).Return(subscriber, nil).Once()
        simManagerService.On("GetSimManagerService").Return(mockSimManagerClient, nil).Once()
        mockSimManagerClient.On("GetSimsBySubscriber", mock.Anything, &simMangerPb.GetSimsBySubscriberRequest{
            SubscriberId: subscriberId.String(),
        }).Return(nil, errors.New("Failed to get Sims by subscriber")).Once()

        s := NewSubscriberServer(OrgName, subscriberRepo, nil, simManagerService, OrgId, nil, nil)
        resp, err := s.Get(context.TODO(), &pb.GetSubscriberRequest{
            SubscriberId: subscriberId.String(),
        })

        assert.Error(t, err)
        assert.Nil(t, resp)
        assert.Contains(t, err.Error(), "Failed to get Sims by subscriber")

        subscriberRepo.AssertExpectations(t)
        simManagerService.AssertExpectations(t)
        mockSimManagerClient.AssertExpectations(t)
    })
   
	t.Run("GetSimManagerService error", func(t *testing.T) {
        subscriberId := uuid.NewV4()
        subscriber := &db.Subscriber{
            SubscriberId: subscriberId,
            Email:        "test@example.com",
        }

        subscriberRepo := &mocks.SubscriberRepo{}
        simManagerService := &mocks.SimManagerClientProvider{}

        subscriberRepo.On("Get", subscriberId).Return(subscriber, nil).Once()
        simManagerService.On("GetSimManagerService").Return(nil, errors.New("Failed to get SimManagerService")).Once()

        s := NewSubscriberServer(OrgName, subscriberRepo, nil, simManagerService, OrgId, nil, nil)
        resp, err := s.Get(context.TODO(), &pb.GetSubscriberRequest{
            SubscriberId: subscriberId.String(),
        })

        assert.Error(t, err)
        assert.Nil(t, resp)
        assert.Contains(t, err.Error(), "Failed to get SimManagerService")

        subscriberRepo.AssertExpectations(t)
        simManagerService.AssertExpectations(t)
    })
}
func TestSubscriberServer_GetbyNetwork(t *testing.T) {

	t.Run("NetworkNotFound", func(t *testing.T) {
		var networkId = uuid.Nil
		networkClient := &cmocks.NetworkClient{}

		subRepo := &mocks.SubscriberRepo{}
		simManagerProvider := &mocks.SimManagerClientProvider{}
        mockSimManagerService := &mocksks.SimManagerServiceClient{}

       
		subRepo.On("GetByNetwork", networkId).Return(nil, gorm.ErrRecordNotFound).Once()
		simManagerProvider.On("GetSimManagerService").Return(mockSimManagerService, nil).Once()
		mockSimManagerService.On("GetSimsByNetwork", mock.Anything, &simMangerPb.GetSimsByNetworkRequest{
			NetworkId: networkId.String(),
		}).Return(&simMangerPb.GetSimsBySubscriberResponse{
			Sims: []*simMangerPb.Sim{},
		}, nil).Once()
		s := NewSubscriberServer(OrgName, subRepo, nil, nil, OrgId, nil, networkClient)
		subResp, err := s.GetByNetwork(context.TODO(), &pb.GetByNetworkRequest{
			NetworkId: networkId.String()})

		assert.Error(t, err)
		assert.Nil(t, subResp)
		subRepo.AssertExpectations(t)
	})

	t.Run("NetworkUUIDInvalid", func(t *testing.T) {
		var networkId = "1"
		networkClient := &cmocks.NetworkClient{}

		subRepo := &mocks.SubscriberRepo{}

		s := NewSubscriberServer(OrgName, subRepo, nil, nil, OrgId, nil, networkClient)
		subResp, err := s.GetByNetwork(context.TODO(), &pb.GetByNetworkRequest{
			NetworkId: networkId})

		assert.Error(t, err)
		assert.Nil(t, subResp)
		subRepo.AssertExpectations(t)
	})
	t.Run("subscriberFoundByNetwork ", func(t *testing.T) {
        subscriberId := uuid.NewV4()
        networkId := uuid.NewV4()

        subRepo := &mocks.SubscriberRepo{}
        simManagerProvider := &mocks.SimManagerClientProvider{}
        mockSimManagerService := &mocksks.SimManagerServiceClient{}

       
        subscribers := []db.Subscriber{
            {SubscriberId: subscriberId, Email: "user1@example.com", NetworkId:networkId},
            {SubscriberId: subscriberId, Email: "user2@example.com", NetworkId: networkId},
        }

		sim := &simMangerPb.Sim{
			Id:             subscriberId.String(), 
			SubscriberId:   subscriberId.String(), 
		}
        subRepo.On("GetByNetwork", networkId).Return(subscribers, nil).Once()
        simManagerProvider.On("GetSimManagerService").Return(mockSimManagerService, nil).Once()
        mockSimManagerService.On("GetSimsByNetwork", mock.Anything, &simMangerPb.GetSimsByNetworkRequest{
            NetworkId: networkId.String(),
        }).Return(&simMangerPb.GetSimsByNetworkResponse{
            Sims: []*simMangerPb.Sim{
				sim,
			},
        }, nil).Once()

        s := NewSubscriberServer(OrgName, subRepo, nil, simManagerProvider, OrgId, nil, nil)

        resp, err := s.GetByNetwork(context.TODO(), &pb.GetByNetworkRequest{
            NetworkId: networkId.String(),
        })

        assert.NoError(t, err)
        assert.NotNil(t, resp)
		actualNetworkId, err := uuid.FromString(resp.Subscribers[0].NetworkId)
		assert.NoError(t, err)
	   assert.Equal(t, networkId, actualNetworkId)
        subRepo.AssertExpectations(t)
        simManagerProvider.AssertExpectations(t)
        mockSimManagerService.AssertExpectations(t)
    })
	t.Run("Get subscribers by network with empty result", func(t *testing.T) {
        networkId := uuid.NewV4()
        subscriberRepo := &mocks.SubscriberRepo{}
        simManagerService := &mocks.SimManagerClientProvider{}
        mockSimManagerClient := &mocksks.SimManagerServiceClient{}

        subscriberRepo.On("GetByNetwork", networkId).Return([]db.Subscriber{}, nil).Once()
        simManagerService.On("GetSimManagerService").Return(mockSimManagerClient, nil).Once()
        mockSimManagerClient.On("GetSimsByNetwork", mock.Anything, &simMangerPb.GetSimsByNetworkRequest{
            NetworkId: networkId.String(),
        }).Return(&simMangerPb.GetSimsByNetworkResponse{Sims: []*simMangerPb.Sim{}}, nil).Once()

        s := NewSubscriberServer(OrgName, subscriberRepo, nil, simManagerService, OrgId, nil, nil)
        resp, err := s.GetByNetwork(context.TODO(), &pb.GetByNetworkRequest{
            NetworkId: networkId.String(),
        })

        assert.NoError(t, err)
        assert.NotNil(t, resp)
        assert.Empty(t, resp.Subscribers)
        subscriberRepo.AssertExpectations(t)
        simManagerService.AssertExpectations(t)
        mockSimManagerClient.AssertExpectations(t)
    })
}

func TestGetByEmail(t *testing.T) {
    t.Run("Get subscriber by existing email", func(t *testing.T) {
        email := "johndoe@example.com"
        subscriberId := uuid.NewV4()

        subRepo := &mocks.SubscriberRepo{}
        simManagerProvider := &mocks.SimManagerClientProvider{}
        mockSimManagerService := &mocksks.SimManagerServiceClient{}

        subscriber := &db.Subscriber{
            Email:        email,
            SubscriberId: subscriberId,
        }

		sim := &simMangerPb.Sim{
			Id:             subscriberId.String(), 
			SubscriberId:   subscriberId.String(), 
		}
        subRepo.On("GetByEmail", email).Return(subscriber, nil).Once()
        simManagerProvider.On("GetSimManagerService").Return(mockSimManagerService, nil).Once()
        mockSimManagerService.On("GetSimsBySubscriber", mock.Anything, &simMangerPb.GetSimsBySubscriberRequest{
            SubscriberId: subscriberId.String(),
        }).Return(&simMangerPb.GetSimsBySubscriberResponse{
            Sims: []*simMangerPb.Sim{sim},
        }, nil).Once()

        s := NewSubscriberServer(OrgName, subRepo, nil, simManagerProvider, OrgId, nil, nil)

        resp, err := s.GetByEmail(context.TODO(), &pb.GetSubscriberByEmailRequest{
            Email: email,
        })

        assert.NoError(t, err)
        assert.NotNil(t, resp)
        assert.Equal(t, email, resp.Subscriber.Email)

        subRepo.AssertExpectations(t)
        simManagerProvider.AssertExpectations(t)
        mockSimManagerService.AssertExpectations(t)
    })

    t.Run("Get subscriber by non-existent email", func(t *testing.T) {
        email := "nonexistent@example.com"

        subRepo := &mocks.SubscriberRepo{}
        subRepo.On("GetByEmail", email).Return(nil, gorm.ErrRecordNotFound).Once()

        s := NewSubscriberServer(OrgName, subRepo, nil, nil, OrgId, nil, nil)
        
        resp, err := s.GetByEmail(context.TODO(), &pb.GetSubscriberByEmailRequest{
            Email: email,
        })

        assert.Error(t, err)
        assert.Nil(t, resp)
        assert.True(t, status.Code(err) == codes.NotFound)
        subRepo.AssertExpectations(t)
    })

}
func TestUpdateSubscriber(t *testing.T) {
	t.Run("Update subscriber successfully", func(t *testing.T) {
        subscriberId := uuid.NewV4()
        subscriberRepo := &mocks.SubscriberRepo{}
        msgBus := &cmocks.MsgBusServiceClient{}

        subscriberRepo.On("Update", subscriberId, mock.AnythingOfType("db.Subscriber")).Return(nil).Once()
        msgBus.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

        s := NewSubscriberServer(OrgName, subscriberRepo, msgBus, nil, OrgId, nil, nil)
        _, err := s.Update(context.TODO(), &pb.UpdateSubscriberRequest{
            SubscriberId: subscriberId.String(),
            Email:        "newemail@example.com",
        })

        assert.NoError(t, err)
        subscriberRepo.AssertExpectations(t)
        msgBus.AssertExpectations(t)
    })
  

    t.Run("InvalidSubscriberId", func(t *testing.T) {
        s := NewSubscriberServer(OrgName, nil, nil, nil, OrgId, nil, nil)

        updateReq := &pb.UpdateSubscriberRequest{
            SubscriberId: "invalid-uuid",
        }

        resp, err := s.Update(context.TODO(), updateReq)

        assert.Error(t, err)
        assert.Nil(t, resp)
        assert.Contains(t, err.Error(), "invalid format of subscriber uuid")
    })
}
func TestListSubscribers(t *testing.T) {
    t.Run("SuccessfulList", func(t *testing.T) {
        subRepo := &mocks.SubscriberRepo{}
        simManagerProvider := &mocks.SimManagerClientProvider{}
        mockSimManagerService := &mocksks.SimManagerServiceClient{} // New mock for SimManagerServiceClient
        msgBus := &cmocks.MsgBusServiceClient{}
	
        subscriberId := uuid.NewV4()
        subscribers := []db.Subscriber{
            {SubscriberId: subscriberId, Email: "user1@example.com", NetworkId: uuid.NewV4()},
            {SubscriberId: subscriberId, Email: "user2@example.com", NetworkId: uuid.NewV4()},
        }
		
        subRepo.On("ListSubscribers").Return(subscribers, nil).Once()
        simManagerProvider.On("GetSimManagerService").Return(mockSimManagerService, nil).Once()


        mockSimManagerService.On("ListSims", mock.Anything, &simMangerPb.ListSimsRequest{
        }).Return(&simMangerPb.ListSimsResponse{
            Sims: []*simMangerPb.Sim{
			},
        }, nil).Once()

        s := NewSubscriberServer(OrgName, subRepo, msgBus, simManagerProvider, OrgId, nil, nil)

        subResp, err := s.ListSubscribers(context.TODO(), &pb.ListSubscribersRequest{})

        assert.NoError(t, err)
        assert.NotNil(t, subResp)

        subRepo.AssertExpectations(t)
        simManagerProvider.AssertExpectations(t)
        mockSimManagerService.AssertExpectations(t)
    })

    t.Run("List subscribers with empty result", func(t *testing.T) {
        subscriberRepo := &mocks.SubscriberRepo{}
        simManagerService := &mocks.SimManagerClientProvider{}
        mockSimManagerClient := &mocksks.SimManagerServiceClient{}

        subscriberRepo.On("ListSubscribers").Return([]db.Subscriber{}, nil).Once()
        simManagerService.On("GetSimManagerService").Return(mockSimManagerClient, nil).Once()
        mockSimManagerClient.On("ListSims", mock.Anything, &simMangerPb.ListSimsRequest{}).
            Return(&simMangerPb.ListSimsResponse{Sims: []*simMangerPb.Sim{}}, nil).Once()

        s := NewSubscriberServer(OrgName, subscriberRepo, nil, simManagerService, OrgId, nil, nil)
        resp, err := s.ListSubscribers(context.TODO(), &pb.ListSubscribersRequest{})

        assert.NoError(t, err)
        assert.NotNil(t, resp)
        assert.Empty(t, resp.Subscribers)
        subscriberRepo.AssertExpectations(t)
        simManagerService.AssertExpectations(t)
        mockSimManagerClient.AssertExpectations(t)
    })

    t.Run("List subscribers with database error", func(t *testing.T) {
        subscriberRepo := &mocks.SubscriberRepo{}

        subscriberRepo.On("ListSubscribers").Return(nil, errors.New("database error")).Once()

        s := NewSubscriberServer(OrgName, subscriberRepo, nil, nil, OrgId, nil, nil)
        _, err := s.ListSubscribers(context.TODO(), &pb.ListSubscribersRequest{})

        assert.Error(t, err)
        assert.Contains(t, err.Error(), "database error")
        subscriberRepo.AssertExpectations(t)
    })

}

func TestDeleteSubscriber(t *testing.T) {
    t.Run("Delete subscriber successfully", func(t *testing.T) {
        subscriberId := uuid.NewV4()
        subscriberRepo := &mocks.SubscriberRepo{}
        msgBus := &cmocks.MsgBusServiceClient{}
        simManagerService := &mocks.SimManagerClientProvider{}

        subscriberRepo.On("Get", subscriberId).Return(&db.Subscriber{SubscriberId: subscriberId}, nil).Once()
        subscriberRepo.On("Delete", subscriberId).Return(nil).Once()
        msgBus.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

        s := NewSubscriberServer(OrgName, subscriberRepo, msgBus, simManagerService, OrgId, nil, nil)

        _, err := s.Delete(context.TODO(), &pb.DeleteSubscriberRequest{
            SubscriberId: subscriberId.String(),
        })

        assert.NoError(t, err)
        subscriberRepo.AssertExpectations(t)
        msgBus.AssertExpectations(t)
    })

   
}


