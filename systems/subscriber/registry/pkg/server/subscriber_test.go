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

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ukama/ukama/systems/subscriber/registry/mocks"
	"github.com/ukama/ukama/systems/subscriber/registry/pkg/db"
	"gorm.io/gorm"

	cmocks "github.com/ukama/ukama/systems/common/mocks"
	creg "github.com/ukama/ukama/systems/common/rest/client/registry"
	"github.com/ukama/ukama/systems/common/ukama"
	uuid "github.com/ukama/ukama/systems/common/uuid"
	simMangerPb "github.com/ukama/ukama/systems/subscriber/sim-manager/pb/gen"
	simMocks "github.com/ukama/ukama/systems/subscriber/sim-manager/pb/gen/mocks"

	pb "github.com/ukama/ukama/systems/subscriber/registry/pb/gen"
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
			Name:                  "John",
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
			Name:                  sub.Name,
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
			Name:                  "John",
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

		s := NewSubscriberServer(OrgName, subscriberRepo, msgBus, simManagerService, OrgId, orgClient, networkClient)

		networkClient.On("GetDefault", mock.Anything).Return(
			&creg.NetworkInfo{
				Id: networkId.String(),
			}, nil).Once()

		res, err := s.Add(context.TODO(), &pb.AddSubscriberRequest{
			Name:                  sub.Name,
			Email:                 sub.Email,
			PhoneNumber:           sub.PhoneNumber,
			Gender:                sub.Gender,
			Dob:                   sub.DOB,
			Address:               sub.Address,
			ProofOfIdentification: sub.ProofOfIdentification,
			IdSerial:              sub.IdSerial,
		})

		log.Info("res", res)

		assert.NoError(t, err)
	})

	t.Run("Add subscriber with empty DOB", func(t *testing.T) {
		subscriberRepo := &mocks.SubscriberRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		simManagerService := &mocks.SimManagerClientProvider{}
		orgClient := &cmocks.OrgClient{}
		networkClient := &cmocks.NetworkClient{}
		networkId := uuid.NewV4()

		sub := &db.Subscriber{
			Name:                  "John",
			Email:                 "johndoe@example.com",
			PhoneNumber:           "1234567890",
			Gender:                "Male",
			Address:               "1 Main St",
			ProofOfIdentification: "Passport",
			IdSerial:              "123456789",
			NetworkId:             networkId,
			DOB:                   "", // Empty DOB
		}

		subscriberRepo.On("Add", sub, mock.Anything).Return(nil).Once()
		networkClient.On("Get", networkId.String()).
			Return(&creg.NetworkInfo{
				Id:         networkId.String(),
				Name:       "net-1",
				SyncStatus: ukama.StatusTypeCompleted.String(),
			}, nil).Once()
		msgBus.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

		s := NewSubscriberServer(OrgName, subscriberRepo, msgBus, simManagerService, OrgId, orgClient, networkClient)
		_, err := s.Add(context.TODO(), &pb.AddSubscriberRequest{
			Name:                  sub.Name,
			Email:                 sub.Email,
			PhoneNumber:           sub.PhoneNumber,
			Gender:                sub.Gender,
			Dob:                   "", // Empty DOB
			NetworkId:             networkId.String(),
			Address:               sub.Address,
			ProofOfIdentification: sub.ProofOfIdentification,
			IdSerial:              sub.IdSerial,
		})
		assert.NoError(t, err)
	})

	t.Run("Add subscriber with email case conversion", func(t *testing.T) {
		subscriberRepo := &mocks.SubscriberRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		simManagerService := &mocks.SimManagerClientProvider{}
		orgClient := &cmocks.OrgClient{}
		networkClient := &cmocks.NetworkClient{}
		networkId := uuid.NewV4()

		sub := &db.Subscriber{
			Name:                  "John",
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

		s := NewSubscriberServer(OrgName, subscriberRepo, msgBus, simManagerService, OrgId, orgClient, networkClient)
		_, err := s.Add(context.TODO(), &pb.AddSubscriberRequest{
			Name:                  sub.Name,
			Email:                 "JOHNDOE@EXAMPLE.COM", // Uppercase email
			PhoneNumber:           sub.PhoneNumber,
			Gender:                sub.Gender,
			Dob:                   sub.DOB,
			NetworkId:             networkId.String(),
			Address:               sub.Address,
			ProofOfIdentification: sub.ProofOfIdentification,
			IdSerial:              sub.IdSerial,
		})
		assert.NoError(t, err)
	})

	t.Run("Invalid date format", func(t *testing.T) {
		subscriberRepo := &mocks.SubscriberRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		simManagerService := &mocks.SimManagerClientProvider{}
		orgClient := &cmocks.OrgClient{}
		networkClient := &cmocks.NetworkClient{}

		s := NewSubscriberServer(OrgName, subscriberRepo, msgBus, simManagerService, OrgId, orgClient, networkClient)
		_, err := s.Add(context.TODO(), &pb.AddSubscriberRequest{
			Name:                  "John",
			Email:                 "johndoe@example.com",
			PhoneNumber:           "1234567890",
			Gender:                "Male",
			Dob:                   "invalid-date-format", // Invalid date format
			Address:               "1 Main St",
			ProofOfIdentification: "Passport",
			IdSerial:              "123456789",
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid format for DoB value")
	})

	t.Run("Network not found", func(t *testing.T) {
		subscriberRepo := &mocks.SubscriberRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		simManagerService := &mocks.SimManagerClientProvider{}
		orgClient := &cmocks.OrgClient{}
		networkClient := &cmocks.NetworkClient{}
		networkId := uuid.NewV4()

		networkClient.On("Get", networkId.String()).
			Return(nil, errors.New("network not found")).Once()

		s := NewSubscriberServer(OrgName, subscriberRepo, msgBus, simManagerService, OrgId, orgClient, networkClient)
		_, err := s.Add(context.TODO(), &pb.AddSubscriberRequest{
			Name:                  "John",
			Email:                 "johndoe@example.com",
			PhoneNumber:           "1234567890",
			Gender:                "Male",
			NetworkId:             networkId.String(),
			Address:               "1 Main St",
			ProofOfIdentification: "Passport",
			IdSerial:              "123456789",
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "network not found")
	})

	t.Run("Default network not found", func(t *testing.T) {
		subscriberRepo := &mocks.SubscriberRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		simManagerService := &mocks.SimManagerClientProvider{}
		orgClient := &cmocks.OrgClient{}
		networkClient := &cmocks.NetworkClient{}

		networkClient.On("GetDefault").Return(nil, errors.New("default network not found")).Once()

		s := NewSubscriberServer(OrgName, subscriberRepo, msgBus, simManagerService, OrgId, orgClient, networkClient)
		_, err := s.Add(context.TODO(), &pb.AddSubscriberRequest{
			Name:                  "John",
			Email:                 "johndoe@example.com",
			PhoneNumber:           "1234567890",
			Gender:                "Male",
			Address:               "1 Main St",
			ProofOfIdentification: "Passport",
			IdSerial:              "123456789",
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "default network not found")
	})

	t.Run("Invalid network UUID format", func(t *testing.T) {
		subscriberRepo := &mocks.SubscriberRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		simManagerService := &mocks.SimManagerClientProvider{}
		orgClient := &cmocks.OrgClient{}
		networkClient := &cmocks.NetworkClient{}

		networkClient.On("Get", "invalid-uuid").
			Return(&creg.NetworkInfo{
				Id: "invalid-uuid", // Invalid UUID format
			}, nil).Once()

		s := NewSubscriberServer(OrgName, subscriberRepo, msgBus, simManagerService, OrgId, orgClient, networkClient)
		_, err := s.Add(context.TODO(), &pb.AddSubscriberRequest{
			Name:                  "John",
			Email:                 "johndoe@example.com",
			PhoneNumber:           "1234567890",
			Gender:                "Male",
			NetworkId:             "invalid-uuid",
			Address:               "1 Main St",
			ProofOfIdentification: "Passport",
			IdSerial:              "123456789",
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid format of network uuid")
	})

	t.Run("Database error - duplicate key", func(t *testing.T) {
		subscriberRepo := &mocks.SubscriberRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		simManagerService := &mocks.SimManagerClientProvider{}
		orgClient := &cmocks.OrgClient{}
		networkClient := &cmocks.NetworkClient{}
		networkId := uuid.NewV4()

		sub := &db.Subscriber{
			Name:                  "John",
			Email:                 "johndoe@example.com",
			PhoneNumber:           "1234567890",
			Gender:                "Male",
			Address:               "1 Main St",
			ProofOfIdentification: "Passport",
			IdSerial:              "123456789",
			NetworkId:             networkId,
			DOB:                   time.Now().Add(time.Hour * 24 * 365 * 18).Format(time.RFC3339),
		}

		networkClient.On("Get", networkId.String()).
			Return(&creg.NetworkInfo{
				Id:         networkId.String(),
				Name:       "net-1",
				SyncStatus: ukama.StatusTypeCompleted.String(),
			}, nil).Once()

		subscriberRepo.On("Add", sub, mock.Anything).Return(errors.New("duplicate key value violates unique constraint")).Once()

		s := NewSubscriberServer(OrgName, subscriberRepo, msgBus, simManagerService, OrgId, orgClient, networkClient)
		_, err := s.Add(context.TODO(), &pb.AddSubscriberRequest{
			Name:                  sub.Name,
			Email:                 sub.Email,
			PhoneNumber:           sub.PhoneNumber,
			Gender:                sub.Gender,
			Dob:                   sub.DOB,
			NetworkId:             networkId.String(),
			Address:               sub.Address,
			ProofOfIdentification: sub.ProofOfIdentification,
			IdSerial:              sub.IdSerial,
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "duplicate key value violates unique constraint")
	})

	t.Run("Database error - internal error", func(t *testing.T) {
		subscriberRepo := &mocks.SubscriberRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		simManagerService := &mocks.SimManagerClientProvider{}
		orgClient := &cmocks.OrgClient{}
		networkClient := &cmocks.NetworkClient{}
		networkId := uuid.NewV4()

		sub := &db.Subscriber{
			Name:                  "John",
			Email:                 "johndoe@example.com",
			PhoneNumber:           "1234567890",
			Gender:                "Male",
			Address:               "1 Main St",
			ProofOfIdentification: "Passport",
			IdSerial:              "123456789",
			NetworkId:             networkId,
			DOB:                   time.Now().Add(time.Hour * 24 * 365 * 18).Format(time.RFC3339),
		}

		networkClient.On("Get", networkId.String()).
			Return(&creg.NetworkInfo{
				Id:         networkId.String(),
				Name:       "net-1",
				SyncStatus: ukama.StatusTypeCompleted.String(),
			}, nil).Once()

		subscriberRepo.On("Add", sub, mock.Anything).Return(errors.New("connection timeout")).Once()

		s := NewSubscriberServer(OrgName, subscriberRepo, msgBus, simManagerService, OrgId, orgClient, networkClient)
		_, err := s.Add(context.TODO(), &pb.AddSubscriberRequest{
			Name:                  sub.Name,
			Email:                 sub.Email,
			PhoneNumber:           sub.PhoneNumber,
			Gender:                sub.Gender,
			Dob:                   sub.DOB,
			NetworkId:             networkId.String(),
			Address:               sub.Address,
			ProofOfIdentification: sub.ProofOfIdentification,
			IdSerial:              sub.IdSerial,
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "connection timeout")
	})

	t.Run("Message bus publish error", func(t *testing.T) {
		subscriberRepo := &mocks.SubscriberRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		simManagerService := &mocks.SimManagerClientProvider{}
		orgClient := &cmocks.OrgClient{}
		networkClient := &cmocks.NetworkClient{}
		networkId := uuid.NewV4()

		sub := &db.Subscriber{
			Name:                  "John",
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

		msgBus.On("PublishRequest", mock.Anything, mock.Anything).Return(errors.New("publish failed")).Once()

		s := NewSubscriberServer(OrgName, subscriberRepo, msgBus, simManagerService, OrgId, orgClient, networkClient)
		resp, err := s.Add(context.TODO(), &pb.AddSubscriberRequest{
			Name:                  sub.Name,
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
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.Subscriber)
	})
}

func TestSubscriberServer_Get(t *testing.T) {

	t.Run("Get subscriber successfully", func(t *testing.T) {
		subscriberRepo := &mocks.SubscriberRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		simManagerService := &mocks.SimManagerClientProvider{}
		orgClient := &cmocks.OrgClient{}
		networkClient := &cmocks.NetworkClient{}

		subscriberId := uuid.NewV4()
		networkId := uuid.NewV4()

		subscriber := &db.Subscriber{
			SubscriberId:          subscriberId,
			Name:                  "John",
			Email:                 "johndoe@example.com",
			PhoneNumber:           "1234567890",
			Gender:                "Male",
			Address:               "1 Main St",
			ProofOfIdentification: "Passport",
			IdSerial:              "123456789",
			NetworkId:             networkId,
			DOB:                   time.Now().Add(time.Hour * 24 * 365 * 18).Format(time.RFC3339),
			CreatedAt:             time.Now(),
			UpdatedAt:             time.Now(),
		}

		// Mock the sim manager service client
		simManagerClient := &simMocks.SimManagerServiceClient{}

		subscriberRepo.On("Get", subscriberId).Return(subscriber, nil).Once()
		simManagerService.On("GetSimManagerService").Return(simManagerClient, nil).Once()
		simManagerClient.On("ListSims", mock.Anything, &simMangerPb.ListSimsRequest{
			SubscriberId: subscriberId.String(),
		}).Return(&simMangerPb.ListSimsResponse{
			Sims: []*simMangerPb.Sim{
				{
					Id:           "sim-1",
					SubscriberId: subscriberId.String(),
					NetworkId:    networkId.String(),
					Iccid:        "12345678901234567890",
					Msisdn:       "1234567890",
					Type:         "physical",
					Status:       "active",
					IsPhysical:   true,
				},
			},
		}, nil).Once()

		s := NewSubscriberServer(OrgName, subscriberRepo, msgBus, simManagerService, OrgId, orgClient, networkClient)
		resp, err := s.Get(context.TODO(), &pb.GetSubscriberRequest{
			SubscriberId: subscriberId.String(),
		})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.Subscriber)
		assert.Equal(t, subscriber.Name, resp.Subscriber.Name)
		assert.Equal(t, subscriber.Email, resp.Subscriber.Email)
		assert.Equal(t, subscriber.SubscriberId.String(), resp.Subscriber.SubscriberId)
		assert.Len(t, resp.Subscriber.Sim, 1)
		assert.Equal(t, "sim-1", resp.Subscriber.Sim[0].Id)

		subscriberRepo.AssertExpectations(t)
		simManagerService.AssertExpectations(t)
		simManagerClient.AssertExpectations(t)
	})

	t.Run("Get subscriber with no sims", func(t *testing.T) {
		subscriberRepo := &mocks.SubscriberRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		simManagerService := &mocks.SimManagerClientProvider{}
		orgClient := &cmocks.OrgClient{}
		networkClient := &cmocks.NetworkClient{}

		subscriberId := uuid.NewV4()
		networkId := uuid.NewV4()

		subscriber := &db.Subscriber{
			SubscriberId:          subscriberId,
			Name:                  "John",
			Email:                 "johndoe@example.com",
			PhoneNumber:           "1234567890",
			Gender:                "Male",
			Address:               "1 Main St",
			ProofOfIdentification: "Passport",
			IdSerial:              "123456789",
			NetworkId:             networkId,
			DOB:                   time.Now().Add(time.Hour * 24 * 365 * 18).Format(time.RFC3339),
			CreatedAt:             time.Now(),
			UpdatedAt:             time.Now(),
		}

		simManagerClient := &simMocks.SimManagerServiceClient{}

		subscriberRepo.On("Get", subscriberId).Return(subscriber, nil).Once()
		simManagerService.On("GetSimManagerService").Return(simManagerClient, nil).Once()
		simManagerClient.On("ListSims", mock.Anything, &simMangerPb.ListSimsRequest{
			SubscriberId: subscriberId.String(),
		}).Return(&simMangerPb.ListSimsResponse{
			Sims: []*simMangerPb.Sim{},
		}, nil).Once()

		s := NewSubscriberServer(OrgName, subscriberRepo, msgBus, simManagerService, OrgId, orgClient, networkClient)
		resp, err := s.Get(context.TODO(), &pb.GetSubscriberRequest{
			SubscriberId: subscriberId.String(),
		})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.Subscriber)
		assert.Equal(t, subscriber.Name, resp.Subscriber.Name)
		assert.Equal(t, subscriber.Email, resp.Subscriber.Email)
		assert.Equal(t, subscriber.SubscriberId.String(), resp.Subscriber.SubscriberId)
		assert.Len(t, resp.Subscriber.Sim, 0)

		subscriberRepo.AssertExpectations(t)
		simManagerService.AssertExpectations(t)
		simManagerClient.AssertExpectations(t)
	})

	t.Run("Subscriber not found", func(t *testing.T) {
		subscriberRepo := &mocks.SubscriberRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		simManagerService := &mocks.SimManagerClientProvider{}
		orgClient := &cmocks.OrgClient{}
		networkClient := &cmocks.NetworkClient{}

		subscriberId := uuid.NewV4()

		subscriberRepo.On("Get", subscriberId).Return(nil, gorm.ErrRecordNotFound).Once()

		s := NewSubscriberServer(OrgName, subscriberRepo, msgBus, simManagerService, OrgId, orgClient, networkClient)
		resp, err := s.Get(context.TODO(), &pb.GetSubscriberRequest{
			SubscriberId: subscriberId.String(),
		})

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "subscriber record not found")
		subscriberRepo.AssertExpectations(t)
	})

	t.Run("Invalid subscriber UUID format", func(t *testing.T) {
		subscriberRepo := &mocks.SubscriberRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		simManagerService := &mocks.SimManagerClientProvider{}
		orgClient := &cmocks.OrgClient{}
		networkClient := &cmocks.NetworkClient{}

		s := NewSubscriberServer(OrgName, subscriberRepo, msgBus, simManagerService, OrgId, orgClient, networkClient)
		resp, err := s.Get(context.TODO(), &pb.GetSubscriberRequest{
			SubscriberId: "invalid-uuid",
		})

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "Invalid subscriberId format")
	})

	t.Run("Database error", func(t *testing.T) {
		subscriberRepo := &mocks.SubscriberRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		simManagerService := &mocks.SimManagerClientProvider{}
		orgClient := &cmocks.OrgClient{}
		networkClient := &cmocks.NetworkClient{}

		subscriberId := uuid.NewV4()

		subscriberRepo.On("Get", subscriberId).Return(nil, errors.New("database connection error")).Once()

		s := NewSubscriberServer(OrgName, subscriberRepo, msgBus, simManagerService, OrgId, orgClient, networkClient)
		resp, err := s.Get(context.TODO(), &pb.GetSubscriberRequest{
			SubscriberId: subscriberId.String(),
		})

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "database connection error")
		subscriberRepo.AssertExpectations(t)
	})

	t.Run("Sim manager service client error", func(t *testing.T) {
		subscriberRepo := &mocks.SubscriberRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		simManagerService := &mocks.SimManagerClientProvider{}
		orgClient := &cmocks.OrgClient{}
		networkClient := &cmocks.NetworkClient{}

		subscriberId := uuid.NewV4()
		networkId := uuid.NewV4()

		subscriber := &db.Subscriber{
			SubscriberId:          subscriberId,
			Name:                  "John",
			Email:                 "johndoe@example.com",
			PhoneNumber:           "1234567890",
			Gender:                "Male",
			Address:               "1 Main St",
			ProofOfIdentification: "Passport",
			IdSerial:              "123456789",
			NetworkId:             networkId,
			DOB:                   time.Now().Add(time.Hour * 24 * 365 * 18).Format(time.RFC3339),
			CreatedAt:             time.Now(),
			UpdatedAt:             time.Now(),
		}

		subscriberRepo.On("Get", subscriberId).Return(subscriber, nil).Once()
		simManagerService.On("GetSimManagerService").Return(nil, errors.New("sim manager service unavailable")).Once()

		s := NewSubscriberServer(OrgName, subscriberRepo, msgBus, simManagerService, OrgId, orgClient, networkClient)
		resp, err := s.Get(context.TODO(), &pb.GetSubscriberRequest{
			SubscriberId: subscriberId.String(),
		})

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "sim manager service unavailable")
		subscriberRepo.AssertExpectations(t)
		simManagerService.AssertExpectations(t)
	})

	t.Run("List sims error", func(t *testing.T) {
		subscriberRepo := &mocks.SubscriberRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		simManagerService := &mocks.SimManagerClientProvider{}
		orgClient := &cmocks.OrgClient{}
		networkClient := &cmocks.NetworkClient{}

		subscriberId := uuid.NewV4()
		networkId := uuid.NewV4()

		subscriber := &db.Subscriber{
			SubscriberId:          subscriberId,
			Name:                  "John",
			Email:                 "johndoe@example.com",
			PhoneNumber:           "1234567890",
			Gender:                "Male",
			Address:               "1 Main St",
			ProofOfIdentification: "Passport",
			IdSerial:              "123456789",
			NetworkId:             networkId,
			DOB:                   time.Now().Add(time.Hour * 24 * 365 * 18).Format(time.RFC3339),
			CreatedAt:             time.Now(),
			UpdatedAt:             time.Now(),
		}

		simManagerClient := &simMocks.SimManagerServiceClient{}

		subscriberRepo.On("Get", subscriberId).Return(subscriber, nil).Once()
		simManagerService.On("GetSimManagerService").Return(simManagerClient, nil).Once()
		simManagerClient.On("ListSims", mock.Anything, &simMangerPb.ListSimsRequest{
			SubscriberId: subscriberId.String(),
		}).Return(nil, errors.New("failed to list sims")).Once()

		s := NewSubscriberServer(OrgName, subscriberRepo, msgBus, simManagerService, OrgId, orgClient, networkClient)
		resp, err := s.Get(context.TODO(), &pb.GetSubscriberRequest{
			SubscriberId: subscriberId.String(),
		})

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "failed to list sims")
		subscriberRepo.AssertExpectations(t)
		simManagerService.AssertExpectations(t)
		simManagerClient.AssertExpectations(t)
	})

	t.Run("Get subscriber with multiple sims", func(t *testing.T) {
		subscriberRepo := &mocks.SubscriberRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		simManagerService := &mocks.SimManagerClientProvider{}
		orgClient := &cmocks.OrgClient{}
		networkClient := &cmocks.NetworkClient{}

		subscriberId := uuid.NewV4()
		networkId := uuid.NewV4()

		subscriber := &db.Subscriber{
			SubscriberId:          subscriberId,
			Name:                  "John",
			Email:                 "johndoe@example.com",
			PhoneNumber:           "1234567890",
			Gender:                "Male",
			Address:               "1 Main St",
			ProofOfIdentification: "Passport",
			IdSerial:              "123456789",
			NetworkId:             networkId,
			DOB:                   time.Now().Add(time.Hour * 24 * 365 * 18).Format(time.RFC3339),
			CreatedAt:             time.Now(),
			UpdatedAt:             time.Now(),
		}

		simManagerClient := &simMocks.SimManagerServiceClient{}

		subscriberRepo.On("Get", subscriberId).Return(subscriber, nil).Once()
		simManagerService.On("GetSimManagerService").Return(simManagerClient, nil).Once()
		simManagerClient.On("ListSims", mock.Anything, &simMangerPb.ListSimsRequest{
			SubscriberId: subscriberId.String(),
		}).Return(&simMangerPb.ListSimsResponse{
			Sims: []*simMangerPb.Sim{
				{
					Id:           "sim-1",
					SubscriberId: subscriberId.String(),
					NetworkId:    networkId.String(),
					Iccid:        "12345678901234567890",
					Msisdn:       "1234567890",
					Type:         "physical",
					Status:       "active",
					IsPhysical:   true,
				},
				{
					Id:           "sim-2",
					SubscriberId: subscriberId.String(),
					NetworkId:    networkId.String(),
					Iccid:        "09876543210987654321",
					Msisdn:       "0987654321",
					Type:         "virtual",
					Status:       "inactive",
					IsPhysical:   false,
				},
			},
		}, nil).Once()

		s := NewSubscriberServer(OrgName, subscriberRepo, msgBus, simManagerService, OrgId, orgClient, networkClient)
		resp, err := s.Get(context.TODO(), &pb.GetSubscriberRequest{
			SubscriberId: subscriberId.String(),
		})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.Subscriber)
		assert.Equal(t, subscriber.Name, resp.Subscriber.Name)
		assert.Equal(t, subscriber.Email, resp.Subscriber.Email)
		assert.Equal(t, subscriber.SubscriberId.String(), resp.Subscriber.SubscriberId)
		assert.Len(t, resp.Subscriber.Sim, 2)
		assert.Equal(t, "sim-1", resp.Subscriber.Sim[0].Id)
		assert.Equal(t, "sim-2", resp.Subscriber.Sim[1].Id)
		assert.Equal(t, "physical", resp.Subscriber.Sim[0].Type)
		assert.Equal(t, "virtual", resp.Subscriber.Sim[1].Type)

		subscriberRepo.AssertExpectations(t)
		simManagerService.AssertExpectations(t)
		simManagerClient.AssertExpectations(t)
	})
}

func TestSubscriberServer_GetByEmail(t *testing.T) {

	t.Run("Get subscriber by email successfully", func(t *testing.T) {
		subscriberRepo := &mocks.SubscriberRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		simManagerService := &mocks.SimManagerClientProvider{}
		orgClient := &cmocks.OrgClient{}
		networkClient := &cmocks.NetworkClient{}

		subscriberId := uuid.NewV4()
		networkId := uuid.NewV4()
		email := "johndoe@example.com"

		subscriber := &db.Subscriber{
			SubscriberId:          subscriberId,
			Name:                  "John",
			Email:                 email,
			PhoneNumber:           "1234567890",
			Gender:                "Male",
			Address:               "1 Main St",
			ProofOfIdentification: "Passport",
			IdSerial:              "123456789",
			NetworkId:             networkId,
			DOB:                   time.Now().Add(time.Hour * 24 * 365 * 18).Format(time.RFC3339),
			CreatedAt:             time.Now(),
			UpdatedAt:             time.Now(),
		}

		simManagerClient := &simMocks.SimManagerServiceClient{}

		subscriberRepo.On("GetByEmail", email).Return(subscriber, nil).Once()
		simManagerService.On("GetSimManagerService").Return(simManagerClient, nil).Once()
		simManagerClient.On("ListSims", mock.Anything, &simMangerPb.ListSimsRequest{
			SubscriberId: subscriberId.String(),
		}).Return(&simMangerPb.ListSimsResponse{
			Sims: []*simMangerPb.Sim{
				{
					Id:           "sim-1",
					SubscriberId: subscriberId.String(),
					NetworkId:    networkId.String(),
					Iccid:        "12345678901234567890",
					Msisdn:       "1234567890",
					Type:         "physical",
					Status:       "active",
					IsPhysical:   true,
				},
			},
		}, nil).Once()

		s := NewSubscriberServer(OrgName, subscriberRepo, msgBus, simManagerService, OrgId, orgClient, networkClient)
		resp, err := s.GetByEmail(context.TODO(), &pb.GetSubscriberByEmailRequest{
			Email: email,
		})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.Subscriber)
		assert.Equal(t, subscriber.Name, resp.Subscriber.Name)
		assert.Equal(t, subscriber.Email, resp.Subscriber.Email)
		assert.Equal(t, subscriber.SubscriberId.String(), resp.Subscriber.SubscriberId)
		assert.Len(t, resp.Subscriber.Sim, 1)
		assert.Equal(t, "sim-1", resp.Subscriber.Sim[0].Id)

		subscriberRepo.AssertExpectations(t)
		simManagerService.AssertExpectations(t)
		simManagerClient.AssertExpectations(t)
	})

	t.Run("Get subscriber by email with no sims", func(t *testing.T) {
		subscriberRepo := &mocks.SubscriberRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		simManagerService := &mocks.SimManagerClientProvider{}
		orgClient := &cmocks.OrgClient{}
		networkClient := &cmocks.NetworkClient{}

		subscriberId := uuid.NewV4()
		networkId := uuid.NewV4()
		email := "johndoe@example.com"

		subscriber := &db.Subscriber{
			SubscriberId:          subscriberId,
			Name:                  "John",
			Email:                 email,
			PhoneNumber:           "1234567890",
			Gender:                "Male",
			Address:               "1 Main St",
			ProofOfIdentification: "Passport",
			IdSerial:              "123456789",
			NetworkId:             networkId,
			DOB:                   time.Now().Add(time.Hour * 24 * 365 * 18).Format(time.RFC3339),
			CreatedAt:             time.Now(),
			UpdatedAt:             time.Now(),
		}

		simManagerClient := &simMocks.SimManagerServiceClient{}

		subscriberRepo.On("GetByEmail", email).Return(subscriber, nil).Once()
		simManagerService.On("GetSimManagerService").Return(simManagerClient, nil).Once()
		simManagerClient.On("ListSims", mock.Anything, &simMangerPb.ListSimsRequest{
			SubscriberId: subscriberId.String(),
		}).Return(&simMangerPb.ListSimsResponse{
			Sims: []*simMangerPb.Sim{},
		}, nil).Once()

		s := NewSubscriberServer(OrgName, subscriberRepo, msgBus, simManagerService, OrgId, orgClient, networkClient)
		resp, err := s.GetByEmail(context.TODO(), &pb.GetSubscriberByEmailRequest{
			Email: email,
		})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.Subscriber)
		assert.Equal(t, subscriber.Name, resp.Subscriber.Name)
		assert.Equal(t, subscriber.Email, resp.Subscriber.Email)
		assert.Equal(t, subscriber.SubscriberId.String(), resp.Subscriber.SubscriberId)
		assert.Len(t, resp.Subscriber.Sim, 0)

		subscriberRepo.AssertExpectations(t)
		simManagerService.AssertExpectations(t)
		simManagerClient.AssertExpectations(t)
	})

	t.Run("Subscriber not found by email", func(t *testing.T) {
		subscriberRepo := &mocks.SubscriberRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		simManagerService := &mocks.SimManagerClientProvider{}
		orgClient := &cmocks.OrgClient{}
		networkClient := &cmocks.NetworkClient{}

		email := "nonexistent@example.com"

		subscriberRepo.On("GetByEmail", email).Return(nil, gorm.ErrRecordNotFound).Once()

		s := NewSubscriberServer(OrgName, subscriberRepo, msgBus, simManagerService, OrgId, orgClient, networkClient)
		resp, err := s.GetByEmail(context.TODO(), &pb.GetSubscriberByEmailRequest{
			Email: email,
		})

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "subscriber record not found")
		subscriberRepo.AssertExpectations(t)
	})

	t.Run("Database error when getting subscriber by email", func(t *testing.T) {
		subscriberRepo := &mocks.SubscriberRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		simManagerService := &mocks.SimManagerClientProvider{}
		orgClient := &cmocks.OrgClient{}
		networkClient := &cmocks.NetworkClient{}

		email := "johndoe@example.com"

		subscriberRepo.On("GetByEmail", email).Return(nil, errors.New("database connection error")).Once()

		s := NewSubscriberServer(OrgName, subscriberRepo, msgBus, simManagerService, OrgId, orgClient, networkClient)
		resp, err := s.GetByEmail(context.TODO(), &pb.GetSubscriberByEmailRequest{
			Email: email,
		})

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "database connection error")
		subscriberRepo.AssertExpectations(t)
	})

	t.Run("Sim manager service client error", func(t *testing.T) {
		subscriberRepo := &mocks.SubscriberRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		simManagerService := &mocks.SimManagerClientProvider{}
		orgClient := &cmocks.OrgClient{}
		networkClient := &cmocks.NetworkClient{}

		subscriberId := uuid.NewV4()
		networkId := uuid.NewV4()
		email := "johndoe@example.com"

		subscriber := &db.Subscriber{
			SubscriberId:          subscriberId,
			Name:                  "John",
			Email:                 email,
			PhoneNumber:           "1234567890",
			Gender:                "Male",
			Address:               "1 Main St",
			ProofOfIdentification: "Passport",
			IdSerial:              "123456789",
			NetworkId:             networkId,
			DOB:                   time.Now().Add(time.Hour * 24 * 365 * 18).Format(time.RFC3339),
			CreatedAt:             time.Now(),
			UpdatedAt:             time.Now(),
		}

		subscriberRepo.On("GetByEmail", email).Return(subscriber, nil).Once()
		simManagerService.On("GetSimManagerService").Return(nil, errors.New("sim manager service unavailable")).Once()

		s := NewSubscriberServer(OrgName, subscriberRepo, msgBus, simManagerService, OrgId, orgClient, networkClient)
		resp, err := s.GetByEmail(context.TODO(), &pb.GetSubscriberByEmailRequest{
			Email: email,
		})

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "sim manager service unavailable")
		subscriberRepo.AssertExpectations(t)
		simManagerService.AssertExpectations(t)
	})

	t.Run("List sims error", func(t *testing.T) {
		subscriberRepo := &mocks.SubscriberRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		simManagerService := &mocks.SimManagerClientProvider{}
		orgClient := &cmocks.OrgClient{}
		networkClient := &cmocks.NetworkClient{}

		subscriberId := uuid.NewV4()
		networkId := uuid.NewV4()
		email := "johndoe@example.com"

		subscriber := &db.Subscriber{
			SubscriberId:          subscriberId,
			Name:                  "John",
			Email:                 email,
			PhoneNumber:           "1234567890",
			Gender:                "Male",
			Address:               "1 Main St",
			ProofOfIdentification: "Passport",
			IdSerial:              "123456789",
			NetworkId:             networkId,
			DOB:                   time.Now().Add(time.Hour * 24 * 365 * 18).Format(time.RFC3339),
			CreatedAt:             time.Now(),
			UpdatedAt:             time.Now(),
		}

		simManagerClient := &simMocks.SimManagerServiceClient{}

		subscriberRepo.On("GetByEmail", email).Return(subscriber, nil).Once()
		simManagerService.On("GetSimManagerService").Return(simManagerClient, nil).Once()
		simManagerClient.On("ListSims", mock.Anything, &simMangerPb.ListSimsRequest{
			SubscriberId: subscriberId.String(),
		}).Return(nil, errors.New("failed to list sims")).Once()

		s := NewSubscriberServer(OrgName, subscriberRepo, msgBus, simManagerService, OrgId, orgClient, networkClient)
		resp, err := s.GetByEmail(context.TODO(), &pb.GetSubscriberByEmailRequest{
			Email: email,
		})

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "failed to list sims")
		subscriberRepo.AssertExpectations(t)
		simManagerService.AssertExpectations(t)
		simManagerClient.AssertExpectations(t)
	})
}

func TestSubscriberServer_ListSubscribers(t *testing.T) {

	t.Run("List subscribers successfully", func(t *testing.T) {
		subscriberRepo := &mocks.SubscriberRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		simManagerService := &mocks.SimManagerClientProvider{}
		orgClient := &cmocks.OrgClient{}
		networkClient := &cmocks.NetworkClient{}

		subscriberId1 := uuid.NewV4()
		subscriberId2 := uuid.NewV4()
		networkId := uuid.NewV4()

		subscribers := []db.Subscriber{
			{
				SubscriberId:          subscriberId1,
				Name:                  "John",
				Email:                 "john@example.com",
				PhoneNumber:           "1234567890",
				Gender:                "Male",
				Address:               "1 Main St",
				ProofOfIdentification: "Passport",
				IdSerial:              "123456789",
				NetworkId:             networkId,
				DOB:                   time.Now().Add(time.Hour * 24 * 365 * 18).Format(time.RFC3339),
				CreatedAt:             time.Now(),
				UpdatedAt:             time.Now(),
			},
			{
				SubscriberId:          subscriberId2,
				Name:                  "Jane",
				Email:                 "jane@example.com",
				PhoneNumber:           "0987654321",
				Gender:                "Female",
				Address:               "2 Oak Ave",
				ProofOfIdentification: "Driver License",
				IdSerial:              "987654321",
				NetworkId:             networkId,
				DOB:                   time.Now().Add(time.Hour * 24 * 365 * 20).Format(time.RFC3339),
				CreatedAt:             time.Now(),
				UpdatedAt:             time.Now(),
			},
		}

		simManagerClient := &simMocks.SimManagerServiceClient{}

		subscriberRepo.On("ListSubscribers").Return(subscribers, nil).Once()
		simManagerService.On("GetSimManagerService").Return(simManagerClient, nil).Once()
		simManagerClient.On("ListSims", mock.Anything, &simMangerPb.ListSimsRequest{}).Return(&simMangerPb.ListSimsResponse{
			Sims: []*simMangerPb.Sim{
				{
					Id:           "sim-1",
					SubscriberId: subscriberId1.String(),
					NetworkId:    networkId.String(),
					Iccid:        "12345678901234567890",
					Msisdn:       "1234567890",
					Type:         "physical",
					Status:       "active",
					IsPhysical:   true,
					Package: &simMangerPb.Package{
						Id:        "pkg-1",
						StartDate: time.Now().Format(time.RFC3339),
						EndDate:   time.Now().Add(time.Hour * 24 * 30).Format(time.RFC3339),
					},
				},
				{
					Id:           "sim-2",
					SubscriberId: subscriberId2.String(),
					NetworkId:    networkId.String(),
					Iccid:        "09876543210987654321",
					Msisdn:       "0987654321",
					Type:         "virtual",
					Status:       "inactive",
					IsPhysical:   false,
					Package: &simMangerPb.Package{
						Id:        "pkg-2",
						StartDate: time.Now().Add(time.Hour * 24).Format(time.RFC3339),
						EndDate:   time.Now().Add(time.Hour * 24 * 60).Format(time.RFC3339),
					},
				},
			},
		}, nil).Once()

		s := NewSubscriberServer(OrgName, subscriberRepo, msgBus, simManagerService, OrgId, orgClient, networkClient)
		resp, err := s.ListSubscribers(context.TODO(), &pb.ListSubscribersRequest{})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Len(t, resp.Subscribers, 2)

		assert.Equal(t, "John", resp.Subscribers[0].Name)
		assert.Equal(t, "john@example.com", resp.Subscribers[0].Email)
		assert.Equal(t, subscriberId1.String(), resp.Subscribers[0].SubscriberId)
		assert.Len(t, resp.Subscribers[0].Sim, 1)
		assert.Equal(t, "sim-1", resp.Subscribers[0].Sim[0].Id)

		// Check second subscriber
		assert.Equal(t, "Jane", resp.Subscribers[1].Name)
		assert.Equal(t, "jane@example.com", resp.Subscribers[1].Email)
		assert.Equal(t, subscriberId2.String(), resp.Subscribers[1].SubscriberId)
		assert.Len(t, resp.Subscribers[1].Sim, 1)
		assert.Equal(t, "sim-2", resp.Subscribers[1].Sim[0].Id)

		subscriberRepo.AssertExpectations(t)
		simManagerService.AssertExpectations(t)
		simManagerClient.AssertExpectations(t)
	})

	t.Run("List subscribers with no sims", func(t *testing.T) {
		subscriberRepo := &mocks.SubscriberRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		simManagerService := &mocks.SimManagerClientProvider{}
		orgClient := &cmocks.OrgClient{}
		networkClient := &cmocks.NetworkClient{}

		subscriberId := uuid.NewV4()
		networkId := uuid.NewV4()

		subscribers := []db.Subscriber{
			{
				SubscriberId:          subscriberId,
				Name:                  "John",
				Email:                 "john@example.com",
				PhoneNumber:           "1234567890",
				Gender:                "Male",
				Address:               "1 Main St",
				ProofOfIdentification: "Passport",
				IdSerial:              "123456789",
				NetworkId:             networkId,
				DOB:                   time.Now().Add(time.Hour * 24 * 365 * 18).Format(time.RFC3339),
				CreatedAt:             time.Now(),
				UpdatedAt:             time.Now(),
			},
		}

		simManagerClient := &simMocks.SimManagerServiceClient{}

		subscriberRepo.On("ListSubscribers").Return(subscribers, nil).Once()
		simManagerService.On("GetSimManagerService").Return(simManagerClient, nil).Once()
		simManagerClient.On("ListSims", mock.Anything, &simMangerPb.ListSimsRequest{}).Return(&simMangerPb.ListSimsResponse{
			Sims: []*simMangerPb.Sim{},
		}, nil).Once()

		s := NewSubscriberServer(OrgName, subscriberRepo, msgBus, simManagerService, OrgId, orgClient, networkClient)
		resp, err := s.ListSubscribers(context.TODO(), &pb.ListSubscribersRequest{})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Len(t, resp.Subscribers, 1)
		assert.Equal(t, "John", resp.Subscribers[0].Name)
		assert.Equal(t, "john@example.com", resp.Subscribers[0].Email)
		assert.Equal(t, subscriberId.String(), resp.Subscribers[0].SubscriberId)
		assert.Len(t, resp.Subscribers[0].Sim, 0)

		subscriberRepo.AssertExpectations(t)
		simManagerService.AssertExpectations(t)
		simManagerClient.AssertExpectations(t)
	})

	t.Run("List subscribers with empty list", func(t *testing.T) {
		subscriberRepo := &mocks.SubscriberRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		simManagerService := &mocks.SimManagerClientProvider{}
		orgClient := &cmocks.OrgClient{}
		networkClient := &cmocks.NetworkClient{}

		subscribers := []db.Subscriber{}

		simManagerClient := &simMocks.SimManagerServiceClient{}

		subscriberRepo.On("ListSubscribers").Return(subscribers, nil).Once()
		simManagerService.On("GetSimManagerService").Return(simManagerClient, nil).Once()
		simManagerClient.On("ListSims", mock.Anything, &simMangerPb.ListSimsRequest{}).Return(&simMangerPb.ListSimsResponse{
			Sims: []*simMangerPb.Sim{},
		}, nil).Once()

		s := NewSubscriberServer(OrgName, subscriberRepo, msgBus, simManagerService, OrgId, orgClient, networkClient)
		resp, err := s.ListSubscribers(context.TODO(), &pb.ListSubscribersRequest{})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Len(t, resp.Subscribers, 0)

		subscriberRepo.AssertExpectations(t)
		simManagerService.AssertExpectations(t)
		simManagerClient.AssertExpectations(t)
	})

	t.Run("Database error when listing subscribers", func(t *testing.T) {
		subscriberRepo := &mocks.SubscriberRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		simManagerService := &mocks.SimManagerClientProvider{}
		orgClient := &cmocks.OrgClient{}
		networkClient := &cmocks.NetworkClient{}

		subscriberRepo.On("ListSubscribers").Return(nil, errors.New("database connection error")).Once()

		s := NewSubscriberServer(OrgName, subscriberRepo, msgBus, simManagerService, OrgId, orgClient, networkClient)
		resp, err := s.ListSubscribers(context.TODO(), &pb.ListSubscribersRequest{})

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "database connection error")
		subscriberRepo.AssertExpectations(t)
	})

	t.Run("Sim manager service client error", func(t *testing.T) {
		subscriberRepo := &mocks.SubscriberRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		simManagerService := &mocks.SimManagerClientProvider{}
		orgClient := &cmocks.OrgClient{}
		networkClient := &cmocks.NetworkClient{}

		subscriberId := uuid.NewV4()
		networkId := uuid.NewV4()

		subscribers := []db.Subscriber{
			{
				SubscriberId:          subscriberId,
				Name:                  "John",
				Email:                 "john@example.com",
				PhoneNumber:           "1234567890",
				Gender:                "Male",
				Address:               "1 Main St",
				ProofOfIdentification: "Passport",
				IdSerial:              "123456789",
				NetworkId:             networkId,
				DOB:                   time.Now().Add(time.Hour * 24 * 365 * 18).Format(time.RFC3339),
				CreatedAt:             time.Now(),
				UpdatedAt:             time.Now(),
			},
		}

		subscriberRepo.On("ListSubscribers").Return(subscribers, nil).Once()
		simManagerService.On("GetSimManagerService").Return(nil, errors.New("sim manager service unavailable")).Once()

		s := NewSubscriberServer(OrgName, subscriberRepo, msgBus, simManagerService, OrgId, orgClient, networkClient)
		resp, err := s.ListSubscribers(context.TODO(), &pb.ListSubscribersRequest{})

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "sim manager service unavailable")
		subscriberRepo.AssertExpectations(t)
		simManagerService.AssertExpectations(t)
	})

	t.Run("List sims error", func(t *testing.T) {
		subscriberRepo := &mocks.SubscriberRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		simManagerService := &mocks.SimManagerClientProvider{}
		orgClient := &cmocks.OrgClient{}
		networkClient := &cmocks.NetworkClient{}

		subscriberId := uuid.NewV4()
		networkId := uuid.NewV4()

		subscribers := []db.Subscriber{
			{
				SubscriberId:          subscriberId,
				Name:                  "John",
				Email:                 "john@example.com",
				PhoneNumber:           "1234567890",
				Gender:                "Male",
				Address:               "1 Main St",
				ProofOfIdentification: "Passport",
				IdSerial:              "123456789",
				NetworkId:             networkId,
				DOB:                   time.Now().Add(time.Hour * 24 * 365 * 18).Format(time.RFC3339),
				CreatedAt:             time.Now(),
				UpdatedAt:             time.Now(),
			},
		}

		simManagerClient := &simMocks.SimManagerServiceClient{}

		subscriberRepo.On("ListSubscribers").Return(subscribers, nil).Once()
		simManagerService.On("GetSimManagerService").Return(simManagerClient, nil).Once()
		simManagerClient.On("ListSims", mock.Anything, &simMangerPb.ListSimsRequest{}).Return(nil, errors.New("failed to list sims")).Once()

		s := NewSubscriberServer(OrgName, subscriberRepo, msgBus, simManagerService, OrgId, orgClient, networkClient)
		resp, err := s.ListSubscribers(context.TODO(), &pb.ListSubscribersRequest{})

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "failed to list sims")
		subscriberRepo.AssertExpectations(t)
		simManagerService.AssertExpectations(t)
		simManagerClient.AssertExpectations(t)
	})

	t.Run("Invalid package start date format", func(t *testing.T) {
		subscriberRepo := &mocks.SubscriberRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		simManagerService := &mocks.SimManagerClientProvider{}
		orgClient := &cmocks.OrgClient{}
		networkClient := &cmocks.NetworkClient{}

		subscriberId := uuid.NewV4()
		networkId := uuid.NewV4()

		subscribers := []db.Subscriber{
			{
				SubscriberId:          subscriberId,
				Name:                  "John",
				Email:                 "john@example.com",
				PhoneNumber:           "1234567890",
				Gender:                "Male",
				Address:               "1 Main St",
				ProofOfIdentification: "Passport",
				IdSerial:              "123456789",
				NetworkId:             networkId,
				DOB:                   time.Now().Add(time.Hour * 24 * 365 * 18).Format(time.RFC3339),
				CreatedAt:             time.Now(),
				UpdatedAt:             time.Now(),
			},
		}

		simManagerClient := &simMocks.SimManagerServiceClient{}

		subscriberRepo.On("ListSubscribers").Return(subscribers, nil).Once()
		simManagerService.On("GetSimManagerService").Return(simManagerClient, nil).Once()
		simManagerClient.On("ListSims", mock.Anything, &simMangerPb.ListSimsRequest{}).Return(&simMangerPb.ListSimsResponse{
			Sims: []*simMangerPb.Sim{
				{
					Id:           "sim-1",
					SubscriberId: subscriberId.String(),
					NetworkId:    networkId.String(),
					Iccid:        "12345678901234567890",
					Msisdn:       "1234567890",
					Type:         "physical",
					Status:       "active",
					IsPhysical:   true,
					Package: &simMangerPb.Package{
						Id:        "pkg-1",
						StartDate: "invalid-date-format", // Invalid date format
						EndDate:   time.Now().Add(time.Hour * 24 * 30).Format(time.RFC3339),
					},
				},
			},
		}, nil).Once()

		s := NewSubscriberServer(OrgName, subscriberRepo, msgBus, simManagerService, OrgId, orgClient, networkClient)
		resp, err := s.ListSubscribers(context.TODO(), &pb.ListSubscribersRequest{})

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "invalid format for Package.StartDate value")
		subscriberRepo.AssertExpectations(t)
		simManagerService.AssertExpectations(t)
		simManagerClient.AssertExpectations(t)
	})

	t.Run("Invalid package end date format", func(t *testing.T) {
		subscriberRepo := &mocks.SubscriberRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		simManagerService := &mocks.SimManagerClientProvider{}
		orgClient := &cmocks.OrgClient{}
		networkClient := &cmocks.NetworkClient{}

		subscriberId := uuid.NewV4()
		networkId := uuid.NewV4()

		subscribers := []db.Subscriber{
			{
				SubscriberId:          subscriberId,
				Name:                  "John",
				Email:                 "john@example.com",
				PhoneNumber:           "1234567890",
				Gender:                "Male",
				Address:               "1 Main St",
				ProofOfIdentification: "Passport",
				IdSerial:              "123456789",
				NetworkId:             networkId,
				DOB:                   time.Now().Add(time.Hour * 24 * 365 * 18).Format(time.RFC3339),
				CreatedAt:             time.Now(),
				UpdatedAt:             time.Now(),
			},
		}

		simManagerClient := &simMocks.SimManagerServiceClient{}

		subscriberRepo.On("ListSubscribers").Return(subscribers, nil).Once()
		simManagerService.On("GetSimManagerService").Return(simManagerClient, nil).Once()
		simManagerClient.On("ListSims", mock.Anything, &simMangerPb.ListSimsRequest{}).Return(&simMangerPb.ListSimsResponse{
			Sims: []*simMangerPb.Sim{
				{
					Id:           "sim-1",
					SubscriberId: subscriberId.String(),
					NetworkId:    networkId.String(),
					Iccid:        "12345678901234567890",
					Msisdn:       "1234567890",
					Type:         "physical",
					Status:       "active",
					IsPhysical:   true,
					Package: &simMangerPb.Package{
						Id:        "pkg-1",
						StartDate: time.Now().Format(time.RFC3339),
						EndDate:   "invalid-date-format", // Invalid date format
					},
				},
			},
		}, nil).Once()

		s := NewSubscriberServer(OrgName, subscriberRepo, msgBus, simManagerService, OrgId, orgClient, networkClient)
		resp, err := s.ListSubscribers(context.TODO(), &pb.ListSubscribersRequest{})

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "invalid format for Package.EndDate value")
		subscriberRepo.AssertExpectations(t)
		simManagerService.AssertExpectations(t)
		simManagerClient.AssertExpectations(t)
	})

	t.Run("List subscribers with multiple sims per subscriber", func(t *testing.T) {
		subscriberRepo := &mocks.SubscriberRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		simManagerService := &mocks.SimManagerClientProvider{}
		orgClient := &cmocks.OrgClient{}
		networkClient := &cmocks.NetworkClient{}

		subscriberId := uuid.NewV4()
		networkId := uuid.NewV4()

		subscribers := []db.Subscriber{
			{
				SubscriberId:          subscriberId,
				Name:                  "John",
				Email:                 "john@example.com",
				PhoneNumber:           "1234567890",
				Gender:                "Male",
				Address:               "1 Main St",
				ProofOfIdentification: "Passport",
				IdSerial:              "123456789",
				NetworkId:             networkId,
				DOB:                   time.Now().Add(time.Hour * 24 * 365 * 18).Format(time.RFC3339),
				CreatedAt:             time.Now(),
				UpdatedAt:             time.Now(),
			},
		}

		simManagerClient := &simMocks.SimManagerServiceClient{}

		subscriberRepo.On("ListSubscribers").Return(subscribers, nil).Once()
		simManagerService.On("GetSimManagerService").Return(simManagerClient, nil).Once()
		simManagerClient.On("ListSims", mock.Anything, &simMangerPb.ListSimsRequest{}).Return(&simMangerPb.ListSimsResponse{
			Sims: []*simMangerPb.Sim{
				{
					Id:           "sim-1",
					SubscriberId: subscriberId.String(),
					NetworkId:    networkId.String(),
					Iccid:        "12345678901234567890",
					Msisdn:       "1234567890",
					Type:         "physical",
					Status:       "active",
					IsPhysical:   true,
					Package: &simMangerPb.Package{
						Id:        "pkg-1",
						StartDate: time.Now().Format(time.RFC3339),
						EndDate:   time.Now().Add(time.Hour * 24 * 30).Format(time.RFC3339),
					},
				},
				{
					Id:           "sim-2",
					SubscriberId: subscriberId.String(),
					NetworkId:    networkId.String(),
					Iccid:        "09876543210987654321",
					Msisdn:       "0987654321",
					Type:         "virtual",
					Status:       "inactive",
					IsPhysical:   false,
					Package: &simMangerPb.Package{
						Id:        "pkg-2",
						StartDate: time.Now().Add(time.Hour * 24).Format(time.RFC3339),
						EndDate:   time.Now().Add(time.Hour * 24 * 60).Format(time.RFC3339),
					},
				},
			},
		}, nil).Once()

		s := NewSubscriberServer(OrgName, subscriberRepo, msgBus, simManagerService, OrgId, orgClient, networkClient)
		resp, err := s.ListSubscribers(context.TODO(), &pb.ListSubscribersRequest{})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Len(t, resp.Subscribers, 1)
		assert.Equal(t, "John", resp.Subscribers[0].Name)
		assert.Equal(t, "john@example.com", resp.Subscribers[0].Email)
		assert.Equal(t, subscriberId.String(), resp.Subscribers[0].SubscriberId)
		assert.Len(t, resp.Subscribers[0].Sim, 2)
		assert.Equal(t, "sim-1", resp.Subscribers[0].Sim[0].Id)
		assert.Equal(t, "sim-2", resp.Subscribers[0].Sim[1].Id)
		assert.Equal(t, "physical", resp.Subscribers[0].Sim[0].Type)
		assert.Equal(t, "virtual", resp.Subscribers[0].Sim[1].Type)

		subscriberRepo.AssertExpectations(t)
		simManagerService.AssertExpectations(t)
		simManagerClient.AssertExpectations(t)
	})
}

func TestSubscriberServer_GetbyNetwork(t *testing.T) {

	t.Run("Get subscribers by network successfully", func(t *testing.T) {
		subscriberRepo := &mocks.SubscriberRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		simManagerService := &mocks.SimManagerClientProvider{}
		orgClient := &cmocks.OrgClient{}
		networkClient := &cmocks.NetworkClient{}

		subscriberId1 := uuid.NewV4()
		subscriberId2 := uuid.NewV4()
		networkId := uuid.NewV4()

		subscribers := []db.Subscriber{
			{
				SubscriberId:          subscriberId1,
				Name:                  "John",
				Email:                 "john@example.com",
				PhoneNumber:           "1234567890",
				Gender:                "Male",
				Address:               "1 Main St",
				ProofOfIdentification: "Passport",
				IdSerial:              "123456789",
				NetworkId:             networkId,
				DOB:                   time.Now().Add(time.Hour * 24 * 365 * 18).Format(time.RFC3339),
				CreatedAt:             time.Now(),
				UpdatedAt:             time.Now(),
			},
			{
				SubscriberId:          subscriberId2,
				Name:                  "Jane",
				Email:                 "jane@example.com",
				PhoneNumber:           "0987654321",
				Gender:                "Female",
				Address:               "2 Oak Ave",
				ProofOfIdentification: "Driver License",
				IdSerial:              "987654321",
				NetworkId:             networkId,
				DOB:                   time.Now().Add(time.Hour * 24 * 365 * 20).Format(time.RFC3339),
				CreatedAt:             time.Now(),
				UpdatedAt:             time.Now(),
			},
		}

		simManagerClient := &simMocks.SimManagerServiceClient{}

		subscriberRepo.On("GetByNetwork", networkId).Return(subscribers, nil).Once()
		simManagerService.On("GetSimManagerService").Return(simManagerClient, nil).Once()
		simManagerClient.On("ListSims", mock.Anything, &simMangerPb.ListSimsRequest{NetworkId: networkId.String()}).Return(&simMangerPb.ListSimsResponse{
			Sims: []*simMangerPb.Sim{
				{
					Id:           "sim-1",
					SubscriberId: subscriberId1.String(),
					NetworkId:    networkId.String(),
					Iccid:        "12345678901234567890",
					Msisdn:       "1234567890",
					Type:         "physical",
					Status:       "active",
					IsPhysical:   true,
					Package: &simMangerPb.Package{
						Id:        "pkg-1",
						StartDate: time.Now().Format(time.RFC3339),
						EndDate:   time.Now().Add(time.Hour * 24 * 30).Format(time.RFC3339),
					},
				},
				{
					Id:           "sim-2",
					SubscriberId: subscriberId2.String(),
					NetworkId:    networkId.String(),
					Iccid:        "09876543210987654321",
					Msisdn:       "0987654321",
					Type:         "virtual",
					Status:       "inactive",
					IsPhysical:   false,
					Package: &simMangerPb.Package{
						Id:        "pkg-2",
						StartDate: time.Now().Add(time.Hour * 24).Format(time.RFC3339),
						EndDate:   time.Now().Add(time.Hour * 24 * 60).Format(time.RFC3339),
					},
				},
			},
		}, nil).Once()

		s := NewSubscriberServer(OrgName, subscriberRepo, msgBus, simManagerService, OrgId, orgClient, networkClient)
		resp, err := s.GetByNetwork(context.TODO(), &pb.GetByNetworkRequest{NetworkId: networkId.String()})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Len(t, resp.Subscribers, 2)

		// Check first subscriber
		assert.Equal(t, "John", resp.Subscribers[0].Name)
		assert.Equal(t, "john@example.com", resp.Subscribers[0].Email)
		assert.Equal(t, subscriberId1.String(), resp.Subscribers[0].SubscriberId)
		assert.Len(t, resp.Subscribers[0].Sim, 1)
		assert.Equal(t, "sim-1", resp.Subscribers[0].Sim[0].Id)

		// Check second subscriber
		assert.Equal(t, "Jane", resp.Subscribers[1].Name)
		assert.Equal(t, "jane@example.com", resp.Subscribers[1].Email)
		assert.Equal(t, subscriberId2.String(), resp.Subscribers[1].SubscriberId)
		assert.Len(t, resp.Subscribers[1].Sim, 1)
		assert.Equal(t, "sim-2", resp.Subscribers[1].Sim[0].Id)

		subscriberRepo.AssertExpectations(t)
		simManagerService.AssertExpectations(t)
		simManagerClient.AssertExpectations(t)
	})

	t.Run("Get subscribers by network with no sims", func(t *testing.T) {
		subscriberRepo := &mocks.SubscriberRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		simManagerService := &mocks.SimManagerClientProvider{}
		orgClient := &cmocks.OrgClient{}
		networkClient := &cmocks.NetworkClient{}

		subscriberId := uuid.NewV4()
		networkId := uuid.NewV4()

		subscribers := []db.Subscriber{
			{
				SubscriberId:          subscriberId,
				Name:                  "John",
				Email:                 "john@example.com",
				PhoneNumber:           "1234567890",
				Gender:                "Male",
				Address:               "1 Main St",
				ProofOfIdentification: "Passport",
				IdSerial:              "123456789",
				NetworkId:             networkId,
				DOB:                   time.Now().Add(time.Hour * 24 * 365 * 18).Format(time.RFC3339),
				CreatedAt:             time.Now(),
				UpdatedAt:             time.Now(),
			},
		}

		simManagerClient := &simMocks.SimManagerServiceClient{}

		subscriberRepo.On("GetByNetwork", networkId).Return(subscribers, nil).Once()
		simManagerService.On("GetSimManagerService").Return(simManagerClient, nil).Once()
		simManagerClient.On("ListSims", mock.Anything, &simMangerPb.ListSimsRequest{NetworkId: networkId.String()}).Return(&simMangerPb.ListSimsResponse{
			Sims: []*simMangerPb.Sim{},
		}, nil).Once()

		s := NewSubscriberServer(OrgName, subscriberRepo, msgBus, simManagerService, OrgId, orgClient, networkClient)
		resp, err := s.GetByNetwork(context.TODO(), &pb.GetByNetworkRequest{NetworkId: networkId.String()})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Len(t, resp.Subscribers, 1)
		assert.Equal(t, "John", resp.Subscribers[0].Name)
		assert.Equal(t, "john@example.com", resp.Subscribers[0].Email)
		assert.Equal(t, subscriberId.String(), resp.Subscribers[0].SubscriberId)
		assert.Len(t, resp.Subscribers[0].Sim, 0)

		subscriberRepo.AssertExpectations(t)
		simManagerService.AssertExpectations(t)
		simManagerClient.AssertExpectations(t)
	})

	t.Run("Get subscribers by network with empty subscriber list", func(t *testing.T) {
		subscriberRepo := &mocks.SubscriberRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		simManagerService := &mocks.SimManagerClientProvider{}
		orgClient := &cmocks.OrgClient{}
		networkClient := &cmocks.NetworkClient{}

		networkId := uuid.NewV4()
		subscribers := []db.Subscriber{}

		simManagerClient := &simMocks.SimManagerServiceClient{}

		subscriberRepo.On("GetByNetwork", networkId).Return(subscribers, nil).Once()
		simManagerService.On("GetSimManagerService").Return(simManagerClient, nil).Once()
		simManagerClient.On("ListSims", mock.Anything, &simMangerPb.ListSimsRequest{NetworkId: networkId.String()}).Return(&simMangerPb.ListSimsResponse{
			Sims: []*simMangerPb.Sim{},
		}, nil).Once()

		s := NewSubscriberServer(OrgName, subscriberRepo, msgBus, simManagerService, OrgId, orgClient, networkClient)
		resp, err := s.GetByNetwork(context.TODO(), &pb.GetByNetworkRequest{NetworkId: networkId.String()})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Len(t, resp.Subscribers, 0)

		subscriberRepo.AssertExpectations(t)
		simManagerService.AssertExpectations(t)
		simManagerClient.AssertExpectations(t)
	})

	t.Run("Invalid network UUID", func(t *testing.T) {
		subscriberRepo := &mocks.SubscriberRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		simManagerService := &mocks.SimManagerClientProvider{}
		orgClient := &cmocks.OrgClient{}
		networkClient := &cmocks.NetworkClient{}

		s := NewSubscriberServer(OrgName, subscriberRepo, msgBus, simManagerService, OrgId, orgClient, networkClient)
		resp, err := s.GetByNetwork(context.TODO(), &pb.GetByNetworkRequest{NetworkId: "invalid-uuid"})

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "invalid networkId")
	})

	t.Run("Database error when getting subscribers by network", func(t *testing.T) {
		subscriberRepo := &mocks.SubscriberRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		simManagerService := &mocks.SimManagerClientProvider{}
		orgClient := &cmocks.OrgClient{}
		networkClient := &cmocks.NetworkClient{}

		networkId := uuid.NewV4()

		subscriberRepo.On("GetByNetwork", networkId).Return(nil, errors.New("database connection error")).Once()

		s := NewSubscriberServer(OrgName, subscriberRepo, msgBus, simManagerService, OrgId, orgClient, networkClient)
		resp, err := s.GetByNetwork(context.TODO(), &pb.GetByNetworkRequest{NetworkId: networkId.String()})

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "database connection error")
		subscriberRepo.AssertExpectations(t)
	})

	t.Run("Network not found error", func(t *testing.T) {
		subscriberRepo := &mocks.SubscriberRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		simManagerService := &mocks.SimManagerClientProvider{}
		orgClient := &cmocks.OrgClient{}
		networkClient := &cmocks.NetworkClient{}

		networkId := uuid.NewV4()

		subscriberRepo.On("GetByNetwork", networkId).Return(nil, gorm.ErrRecordNotFound).Once()

		s := NewSubscriberServer(OrgName, subscriberRepo, msgBus, simManagerService, OrgId, orgClient, networkClient)
		resp, err := s.GetByNetwork(context.TODO(), &pb.GetByNetworkRequest{NetworkId: networkId.String()})

		assert.Error(t, err)
		assert.Nil(t, resp)
		subscriberRepo.AssertExpectations(t)
	})

	t.Run("Sim manager service client error", func(t *testing.T) {
		subscriberRepo := &mocks.SubscriberRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		simManagerService := &mocks.SimManagerClientProvider{}
		orgClient := &cmocks.OrgClient{}
		networkClient := &cmocks.NetworkClient{}

		subscriberId := uuid.NewV4()
		networkId := uuid.NewV4()

		subscribers := []db.Subscriber{
			{
				SubscriberId:          subscriberId,
				Name:                  "John",
				Email:                 "john@example.com",
				PhoneNumber:           "1234567890",
				Gender:                "Male",
				Address:               "1 Main St",
				ProofOfIdentification: "Passport",
				IdSerial:              "123456789",
				NetworkId:             networkId,
				DOB:                   time.Now().Add(time.Hour * 24 * 365 * 18).Format(time.RFC3339),
				CreatedAt:             time.Now(),
				UpdatedAt:             time.Now(),
			},
		}

		subscriberRepo.On("GetByNetwork", networkId).Return(subscribers, nil).Once()
		simManagerService.On("GetSimManagerService").Return(nil, errors.New("sim manager service unavailable")).Once()

		s := NewSubscriberServer(OrgName, subscriberRepo, msgBus, simManagerService, OrgId, orgClient, networkClient)
		resp, err := s.GetByNetwork(context.TODO(), &pb.GetByNetworkRequest{NetworkId: networkId.String()})

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "sim manager service unavailable")
		subscriberRepo.AssertExpectations(t)
		simManagerService.AssertExpectations(t)
	})

	t.Run("List sims by network error", func(t *testing.T) {
		subscriberRepo := &mocks.SubscriberRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		simManagerService := &mocks.SimManagerClientProvider{}
		orgClient := &cmocks.OrgClient{}
		networkClient := &cmocks.NetworkClient{}

		subscriberId := uuid.NewV4()
		networkId := uuid.NewV4()

		subscribers := []db.Subscriber{
			{
				SubscriberId:          subscriberId,
				Name:                  "John",
				Email:                 "john@example.com",
				PhoneNumber:           "1234567890",
				Gender:                "Male",
				Address:               "1 Main St",
				ProofOfIdentification: "Passport",
				IdSerial:              "123456789",
				NetworkId:             networkId,
				DOB:                   time.Now().Add(time.Hour * 24 * 365 * 18).Format(time.RFC3339),
				CreatedAt:             time.Now(),
				UpdatedAt:             time.Now(),
			},
		}

		simManagerClient := &simMocks.SimManagerServiceClient{}

		subscriberRepo.On("GetByNetwork", networkId).Return(subscribers, nil).Once()
		simManagerService.On("GetSimManagerService").Return(simManagerClient, nil).Once()
		simManagerClient.On("ListSims", mock.Anything, &simMangerPb.ListSimsRequest{NetworkId: networkId.String()}).Return(nil, errors.New("failed to list sims by network")).Once()

		s := NewSubscriberServer(OrgName, subscriberRepo, msgBus, simManagerService, OrgId, orgClient, networkClient)
		resp, err := s.GetByNetwork(context.TODO(), &pb.GetByNetworkRequest{NetworkId: networkId.String()})

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "failed to list sims by network")
		subscriberRepo.AssertExpectations(t)
		simManagerService.AssertExpectations(t)
		simManagerClient.AssertExpectations(t)
	})
}

func TestSubscriberServer_Update(t *testing.T) {

	t.Run("Update subscriber successfully", func(t *testing.T) {
		subscriberRepo := &mocks.SubscriberRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		simManagerService := &mocks.SimManagerClientProvider{}
		orgClient := &cmocks.OrgClient{}
		networkClient := &cmocks.NetworkClient{}

		subscriberId := uuid.NewV4()

		updateRequest := &pb.UpdateSubscriberRequest{
			SubscriberId:          subscriberId.String(),
			Name:                  "John Updated",
			PhoneNumber:           "9876543210",
			Address:               "123 Updated St",
			ProofOfIdentification: "Driver License",
			IdSerial:              "987654321",
		}

		expectedSubscriber := &db.Subscriber{
			Name:                  updateRequest.Name,
			PhoneNumber:           updateRequest.PhoneNumber,
			Address:               updateRequest.Address,
			ProofOfIdentification: updateRequest.ProofOfIdentification,
			IdSerial:              updateRequest.IdSerial,
			SubscriberId:          subscriberId,
		}

		subscriberRepo.On("Update", subscriberId, *expectedSubscriber).Return(nil).Once()
		msgBus.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

		s := NewSubscriberServer(OrgName, subscriberRepo, msgBus, simManagerService, OrgId, orgClient, networkClient)
		resp, err := s.Update(context.TODO(), updateRequest)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		subscriberRepo.AssertExpectations(t)
		msgBus.AssertExpectations(t)
	})

	t.Run("Update subscriber with partial fields", func(t *testing.T) {
		subscriberRepo := &mocks.SubscriberRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		simManagerService := &mocks.SimManagerClientProvider{}
		orgClient := &cmocks.OrgClient{}
		networkClient := &cmocks.NetworkClient{}

		subscriberId := uuid.NewV4()

		updateRequest := &pb.UpdateSubscriberRequest{
			SubscriberId: subscriberId.String(),
			Name:         "John Updated",
			PhoneNumber:  "9876543210",
		}

		expectedSubscriber := &db.Subscriber{
			Name:                  updateRequest.Name,
			PhoneNumber:           updateRequest.PhoneNumber,
			Address:               "",
			ProofOfIdentification: "",
			IdSerial:              "",
			SubscriberId:          subscriberId,
		}

		subscriberRepo.On("Update", subscriberId, *expectedSubscriber).Return(nil).Once()
		msgBus.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

		s := NewSubscriberServer(OrgName, subscriberRepo, msgBus, simManagerService, OrgId, orgClient, networkClient)
		resp, err := s.Update(context.TODO(), updateRequest)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		subscriberRepo.AssertExpectations(t)
		msgBus.AssertExpectations(t)
	})

	t.Run("Invalid subscriber UUID", func(t *testing.T) {
		subscriberRepo := &mocks.SubscriberRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		simManagerService := &mocks.SimManagerClientProvider{}
		orgClient := &cmocks.OrgClient{}
		networkClient := &cmocks.NetworkClient{}

		updateRequest := &pb.UpdateSubscriberRequest{
			SubscriberId:          "invalid-uuid",
			Name:                  "John Updated",
			PhoneNumber:           "9876543210",
			Address:               "123 Updated St",
			ProofOfIdentification: "Driver License",
			IdSerial:              "987654321",
		}

		s := NewSubscriberServer(OrgName, subscriberRepo, msgBus, simManagerService, OrgId, orgClient, networkClient)
		resp, err := s.Update(context.TODO(), updateRequest)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "invalid format of subscriber uuid")
	})

	t.Run("Empty subscriber UUID", func(t *testing.T) {
		subscriberRepo := &mocks.SubscriberRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		simManagerService := &mocks.SimManagerClientProvider{}
		orgClient := &cmocks.OrgClient{}
		networkClient := &cmocks.NetworkClient{}

		updateRequest := &pb.UpdateSubscriberRequest{
			SubscriberId:          "",
			Name:                  "John Updated",
			PhoneNumber:           "9876543210",
			Address:               "123 Updated St",
			ProofOfIdentification: "Driver License",
			IdSerial:              "987654321",
		}

		s := NewSubscriberServer(OrgName, subscriberRepo, msgBus, simManagerService, OrgId, orgClient, networkClient)
		resp, err := s.Update(context.TODO(), updateRequest)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "invalid format of subscriber uuid")
	})

	t.Run("Database error when updating subscriber", func(t *testing.T) {
		subscriberRepo := &mocks.SubscriberRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		simManagerService := &mocks.SimManagerClientProvider{}
		orgClient := &cmocks.OrgClient{}
		networkClient := &cmocks.NetworkClient{}

		subscriberId := uuid.NewV4()

		updateRequest := &pb.UpdateSubscriberRequest{
			SubscriberId:          subscriberId.String(),
			Name:                  "John Updated",
			PhoneNumber:           "9876543210",
			Address:               "123 Updated St",
			ProofOfIdentification: "Driver License",
			IdSerial:              "987654321",
		}

		expectedSubscriber := &db.Subscriber{
			Name:                  updateRequest.Name,
			PhoneNumber:           updateRequest.PhoneNumber,
			Address:               updateRequest.Address,
			ProofOfIdentification: updateRequest.ProofOfIdentification,
			IdSerial:              updateRequest.IdSerial,
			SubscriberId:          subscriberId,
		}

		subscriberRepo.On("Update", subscriberId, *expectedSubscriber).Return(errors.New("database connection error")).Once()

		s := NewSubscriberServer(OrgName, subscriberRepo, msgBus, simManagerService, OrgId, orgClient, networkClient)
		resp, err := s.Update(context.TODO(), updateRequest)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "database connection error")
		subscriberRepo.AssertExpectations(t)
	})

	t.Run("Subscriber not found error", func(t *testing.T) {
		subscriberRepo := &mocks.SubscriberRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		simManagerService := &mocks.SimManagerClientProvider{}
		orgClient := &cmocks.OrgClient{}
		networkClient := &cmocks.NetworkClient{}

		subscriberId := uuid.NewV4()

		updateRequest := &pb.UpdateSubscriberRequest{
			SubscriberId:          subscriberId.String(),
			Name:                  "John Updated",
			PhoneNumber:           "9876543210",
			Address:               "123 Updated St",
			ProofOfIdentification: "Driver License",
			IdSerial:              "987654321",
		}

		expectedSubscriber := &db.Subscriber{
			Name:                  updateRequest.Name,
			PhoneNumber:           updateRequest.PhoneNumber,
			Address:               updateRequest.Address,
			ProofOfIdentification: updateRequest.ProofOfIdentification,
			IdSerial:              updateRequest.IdSerial,
			SubscriberId:          subscriberId,
		}

		subscriberRepo.On("Update", subscriberId, *expectedSubscriber).Return(gorm.ErrRecordNotFound).Once()

		s := NewSubscriberServer(OrgName, subscriberRepo, msgBus, simManagerService, OrgId, orgClient, networkClient)
		resp, err := s.Update(context.TODO(), updateRequest)

		assert.Error(t, err)
		assert.Nil(t, resp)
		subscriberRepo.AssertExpectations(t)
	})

	t.Run("Update subscriber with all fields empty except ID", func(t *testing.T) {
		subscriberRepo := &mocks.SubscriberRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		simManagerService := &mocks.SimManagerClientProvider{}
		orgClient := &cmocks.OrgClient{}
		networkClient := &cmocks.NetworkClient{}

		subscriberId := uuid.NewV4()

		updateRequest := &pb.UpdateSubscriberRequest{
			SubscriberId: subscriberId.String(),
		}

		expectedSubscriber := &db.Subscriber{
			Name:                  "",
			PhoneNumber:           "",
			Address:               "",
			ProofOfIdentification: "",
			IdSerial:              "",
			SubscriberId:          subscriberId,
		}

		subscriberRepo.On("Update", subscriberId, *expectedSubscriber).Return(nil).Once()
		msgBus.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

		s := NewSubscriberServer(OrgName, subscriberRepo, msgBus, simManagerService, OrgId, orgClient, networkClient)
		resp, err := s.Update(context.TODO(), updateRequest)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		subscriberRepo.AssertExpectations(t)
		msgBus.AssertExpectations(t)
	})

}

func TestSubscriberServer_Delete(t *testing.T) {
	t.Run("Delete subscriber successfully", func(t *testing.T) {
		subscriberRepo := &mocks.SubscriberRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		simManagerService := &mocks.SimManagerClientProvider{}
		orgClient := &cmocks.OrgClient{}
		networkClient := &cmocks.NetworkClient{}

		subscriberId := uuid.NewV4()
		dbSubscriber := &db.Subscriber{
			SubscriberId: subscriberId,
			Name:         "John",
		}

		subscriberRepo.On("Get", subscriberId).Return(dbSubscriber, nil).Once()
		subscriberRepo.On("Delete", subscriberId).Return(nil).Once()
		msgBus.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

		s := NewSubscriberServer(OrgName, subscriberRepo, msgBus, simManagerService, OrgId, orgClient, networkClient)
		resp, err := s.Delete(context.TODO(), &pb.DeleteSubscriberRequest{SubscriberId: subscriberId.String()})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		subscriberRepo.AssertExpectations(t)
		msgBus.AssertExpectations(t)
	})

	t.Run("Invalid subscriber UUID", func(t *testing.T) {
		subscriberRepo := &mocks.SubscriberRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		simManagerService := &mocks.SimManagerClientProvider{}
		orgClient := &cmocks.OrgClient{}
		networkClient := &cmocks.NetworkClient{}

		s := NewSubscriberServer(OrgName, subscriberRepo, msgBus, simManagerService, OrgId, orgClient, networkClient)
		resp, err := s.Delete(context.TODO(), &pb.DeleteSubscriberRequest{SubscriberId: "invalid-uuid"})

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "invalid format of subscriber uuid")
	})

	t.Run("Empty subscriber UUID", func(t *testing.T) {
		subscriberRepo := &mocks.SubscriberRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		simManagerService := &mocks.SimManagerClientProvider{}
		orgClient := &cmocks.OrgClient{}
		networkClient := &cmocks.NetworkClient{}

		s := NewSubscriberServer(OrgName, subscriberRepo, msgBus, simManagerService, OrgId, orgClient, networkClient)
		resp, err := s.Delete(context.TODO(), &pb.DeleteSubscriberRequest{SubscriberId: ""})

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "invalid format of subscriber uuid")
	})

	t.Run("Subscriber not found", func(t *testing.T) {
		subscriberRepo := &mocks.SubscriberRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		simManagerService := &mocks.SimManagerClientProvider{}
		orgClient := &cmocks.OrgClient{}
		networkClient := &cmocks.NetworkClient{}

		subscriberId := uuid.NewV4()
		subscriberRepo.On("Get", subscriberId).Return(nil, gorm.ErrRecordNotFound).Once()

		s := NewSubscriberServer(OrgName, subscriberRepo, msgBus, simManagerService, OrgId, orgClient, networkClient)
		resp, err := s.Delete(context.TODO(), &pb.DeleteSubscriberRequest{SubscriberId: subscriberId.String()})

		assert.Error(t, err)
		assert.Nil(t, resp)
		subscriberRepo.AssertExpectations(t)
	})

	t.Run("Database error on Get", func(t *testing.T) {
		subscriberRepo := &mocks.SubscriberRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		simManagerService := &mocks.SimManagerClientProvider{}
		orgClient := &cmocks.OrgClient{}
		networkClient := &cmocks.NetworkClient{}

		subscriberId := uuid.NewV4()
		subscriberRepo.On("Get", subscriberId).Return(nil, errors.New("db error")).Once()

		s := NewSubscriberServer(OrgName, subscriberRepo, msgBus, simManagerService, OrgId, orgClient, networkClient)
		resp, err := s.Delete(context.TODO(), &pb.DeleteSubscriberRequest{SubscriberId: subscriberId.String()})

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "db error")
		subscriberRepo.AssertExpectations(t)
	})

	t.Run("Database error on Delete", func(t *testing.T) {
		subscriberRepo := &mocks.SubscriberRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		simManagerService := &mocks.SimManagerClientProvider{}
		orgClient := &cmocks.OrgClient{}
		networkClient := &cmocks.NetworkClient{}

		subscriberId := uuid.NewV4()
		dbSubscriber := &db.Subscriber{
			SubscriberId: subscriberId,
			Name:         "John",
		}

		subscriberRepo.On("Get", subscriberId).Return(dbSubscriber, nil).Once()
		subscriberRepo.On("Delete", subscriberId).Return(errors.New("delete error")).Once()

		s := NewSubscriberServer(OrgName, subscriberRepo, msgBus, simManagerService, OrgId, orgClient, networkClient)
		resp, err := s.Delete(context.TODO(), &pb.DeleteSubscriberRequest{SubscriberId: subscriberId.String()})

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "delete error")
		subscriberRepo.AssertExpectations(t)
	})
}
