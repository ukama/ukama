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

	pb "github.com/ukama/ukama/systems/subscriber/registry/pb/gen"
	"github.com/ukama/ukama/systems/subscriber/registry/pkg"
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
			Name:                 "John",
			Email:                "johndoe@example.com",
			PhoneNumber:          "1234567890",
			Gender:               "Male",
			Address:              "1 Main St",
			ProofOfIdentification: "Passport",
			IdSerial:             "123456789",
			NetworkId:            networkId,
			SubscriberStatus:     ukama.SubscriberStatusActive, 
			DOB:                  time.Now().Add(time.Hour * 24 * 365 * 18).Format(time.RFC3339),
		}

		subscriberRepo.On("Add", sub, mock.AnythingOfType("func(*db.Subscriber, *gorm.DB) error")).Return(nil).Once()
		networkClient.On("Get", networkId.String()).
			Return(&creg.NetworkInfo{
				Id:         networkId.String(),
				Name:       "net-1",
				SyncStatus: ukama.StatusTypeCompleted.String(),
			}, nil).Once()
		msgBus.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

		config := &pkg.Config{
			DeletionWorker: &pkg.DeletionWorkerConfig{
				CheckInterval:  time.Minute,
				DeletionTimeout: time.Hour,
				MaxRetries:     3,
			},
		}

		s := NewSubscriberServer(OrgName, subscriberRepo, msgBus, simManagerService, OrgId, regClient, networkClient, config)
		
		defer s.Shutdown()
		
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
			SubscriberStatus:      ukama.SubscriberStatusActive,
			DOB:                   time.Now().Add(time.Hour * 24 * 365 * 18).Format(time.RFC3339),
		}

		subscriberRepo.On("Add", sub, mock.AnythingOfType("func(*db.Subscriber, *gorm.DB) error")).Return(nil).Once()

		msgBus.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

		config := &pkg.Config{
			DeletionWorker: &pkg.DeletionWorkerConfig{
				CheckInterval:  time.Minute,
				DeletionTimeout: time.Hour,
				MaxRetries:     3,
			},
		}

		s := NewSubscriberServer(OrgName, subscriberRepo, msgBus, simManagerService, OrgId, orgClient, networkClient, config)
		
		defer s.Shutdown()

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
}

func TestSubscriberServer_Get(t *testing.T) {

	t.Run("SubscriberNotFound", func(t *testing.T) {
		var subscriberId = uuid.Nil
		networkClient := &cmocks.NetworkClient{}

		subRepo := &mocks.SubscriberRepo{}

		subRepo.On("Get", subscriberId).Return(nil, gorm.ErrRecordNotFound).Once()

		config := &pkg.Config{
			DeletionWorker: &pkg.DeletionWorkerConfig{
				CheckInterval:  time.Minute,
				DeletionTimeout: time.Hour,
				MaxRetries:     3,
			},
		}

		s := NewSubscriberServer(OrgName, subRepo, nil, nil, OrgId, nil, networkClient, config)
		defer s.Shutdown()
		
		subResp, err := s.Get(context.TODO(), &pb.GetSubscriberRequest{
			SubscriberId: subscriberId.String()})

		assert.Error(t, err)
		assert.Nil(t, subResp)
		subRepo.AssertExpectations(t)
	})

	t.Run("SubscriberUUIDInvalid", func(t *testing.T) {
		var subscriberId = "1"
		networkClient := &cmocks.NetworkClient{}

		subRepo := &mocks.SubscriberRepo{}

		config := &pkg.Config{
			DeletionWorker: &pkg.DeletionWorkerConfig{
				CheckInterval:  time.Minute,
				DeletionTimeout: time.Hour,
				MaxRetries:     3,
			},
		}

		s := NewSubscriberServer(OrgName, subRepo, nil, nil, OrgId, nil, networkClient, config)
		defer s.Shutdown()
		
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
		networkClient := &cmocks.NetworkClient{}

		subRepo := &mocks.SubscriberRepo{}

		subRepo.On("GetByNetwork", networkId).Return(nil, gorm.ErrRecordNotFound).Once()

		config := &pkg.Config{
			DeletionWorker: &pkg.DeletionWorkerConfig{
				CheckInterval:  time.Minute,
				DeletionTimeout: time.Hour,
				MaxRetries:     3,
			},
		}

		s := NewSubscriberServer(OrgName, subRepo, nil, nil, OrgId, nil, networkClient, config)
		defer s.Shutdown()
		
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

		config := &pkg.Config{
			DeletionWorker: &pkg.DeletionWorkerConfig{
				CheckInterval:  time.Minute,
				DeletionTimeout: time.Hour,
				MaxRetries:     3,
			},
		}

		s := NewSubscriberServer(OrgName, subRepo, nil, nil, OrgId, nil, networkClient, config)
		defer s.Shutdown()
		
		subResp, err := s.GetByNetwork(context.TODO(), &pb.GetByNetworkRequest{
			NetworkId: networkId})

		assert.Error(t, err)
		assert.Nil(t, subResp)
		subRepo.AssertExpectations(t)
	})
}