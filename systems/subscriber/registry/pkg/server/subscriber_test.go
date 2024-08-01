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

		s := NewSubscriberServer(OrgName, subscriberRepo, msgBus, simManagerService, OrgId, orgClient, networkClient)

		networkClient.On("GetDefault", mock.Anything).Return(
			&creg.NetworkInfo{
				Id: networkId.String(),
			}, nil).Once()

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

		s := NewSubscriberServer(OrgName, subRepo, nil, nil, OrgId, nil, networkClient)
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

		s := NewSubscriberServer(OrgName, subRepo, nil, nil, OrgId, nil, networkClient)
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
}

// func TestGetByEmail(t *testing.T) {
//     t.Run("EmailNotFound", func(t *testing.T) {
// 		email := "johndoe@example.com"

// 		subRepo := &mocks.SubscriberRepo{}


// 		subRepo.On("GetByEmail", email).Return(nil, gorm.ErrRecordNotFound).Once()

// 		s := NewSubscriberServer(OrgName, subRepo, nil, nil, OrgId, nil, nil)
// 		subResp, err := s.GetByEmail(context.TODO(), &pb.GetSubscriberByEmailRequest{
// 			Email:email,
// 		})

// 		assert.Error(t, err)
// 		assert.Nil(t, subResp)
// 		subRepo.AssertExpectations(t)
// 	})
// 	t.Run("EmailFound", func(t *testing.T) {
//         email := "johndoe@example.com"

//         subRepo := &mocks.SubscriberRepo{}
//         simManagerClient := &mocks.SimManagerClientProvider{}

//         // Simulate a found subscriber with the given email
//         subscriber := &db.Subscriber{
//             Email: email,
//             SubscriberId:uuid.NewV4() ,
//         }

//         subRepo.On("GetByEmail", email).Return(subscriber, nil).Once()

//         simManagerService := &mocks.SimManagerClientProvider{}
//         simManagerClient.On("GetSimManagerService").Return(simManagerService, nil).Once()

//         simManagerService.On("GetSimsBySubscriber", mock.Anything, &simMangerPb.GetSimsBySubscriberRequest{
//             SubscriberId: subscriber.SubscriberId.String(),
//         }).Return(&simMangerPb.GetSimsBySubscriberResponse{
//             Sims: []*simMangerPb.Sim{
//             },
//         }, nil).Once()

//         s := NewSubscriberServer(OrgName, subRepo, nil, simManagerClient, OrgId, nil, nil)

//         subResp, err := s.GetByEmail(context.TODO(), &pb.GetSubscriberByEmailRequest{
//             Email: email,
//         })

//         assert.NoError(t, err)
//         assert.NotNil(t, subResp)
//         assert.Equal(t, email, subResp.Subscriber.Email)
//         // Add more assertions as needed for other fields

//         subRepo.AssertExpectations(t)
//         simManagerClient.AssertExpectations(t)
//         simManagerService.AssertExpectations(t)
//     })
// }

