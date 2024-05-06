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

	"github.com/ukama/ukama/systems/subscriber/registry/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"

	cmocks "github.com/ukama/ukama/systems/common/mocks"
	cnucl "github.com/ukama/ukama/systems/common/rest/client/nucleus"
	creg "github.com/ukama/ukama/systems/common/rest/client/registry"
	"github.com/ukama/ukama/systems/common/ukama"

	uuid "github.com/ukama/ukama/systems/common/uuid"
	pb "github.com/ukama/ukama/systems/subscriber/registry/pb/gen"
)

const OrgName = "testorg"
const orgId = "8c6c2bec-5f90-4fee-8ffd-ee6456abf4fc"

func TestAdd(t *testing.T) {
	t.Run("Add subscriber successfully", func(t *testing.T) {

		subscriberRepo := &mocks.SubscriberRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		simManagerService := &mocks.SimManagerClientProvider{}
		regClient := &cmocks.OrgClient{}
		networkClient := &cmocks.NetworkClient{}
		netId := "9e82c8b1-a746-4f2c-a80e-f4d14d863ea3"
		orgId:="7e82c8b1-a746-4f2c-a80e-f4d14d863ea3"

		firstName := "John"
		lastName := "Doe"
		email := "johndoe@example.com"
		phoneNumber := "1234567890"
		gender := "Male"
		address := "1 Main St"
		proofOfIdentification := "Passport"
		idSerial := "123456789"

		regClient.On("Get", OrgName).Return(
			&cnucl.OrgInfo{
				Id:            orgId,
				Name:          OrgName,
				IsDeactivated: false,
			}, nil).Once()
		subscriberRepo.On("Add", mock.AnythingOfType("*db.Subscriber")).Return(nil)
		networkClient.On("Get", netId).
			Return(&creg.NetworkInfo{
				Id:         netId,
				OrgId:      orgId,
				Name:       "net-1",
				SyncStatus: ukama.StatusTypeCompleted.String(),
			}, nil).Once()
		msgBus.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

		s := NewSubscriberServer(OrgName, subscriberRepo, msgBus, simManagerService, orgId, regClient, networkClient)
		_, err := s.Add(context.TODO(), &pb.AddSubscriberRequest{
			FirstName:             firstName,
			NetworkId:             netId,
			LastName:              lastName,
			Email:                 email,
			PhoneNumber:           phoneNumber,
			Gender:                gender,
			Dob:                   time.Now().Add(time.Hour * 24 * 365 * 18).Format(time.RFC3339),
			Address:               address,
			ProofOfIdentification: proofOfIdentification,
			IdSerial:              idSerial,
		})
		assert.NoError(t, err)
	})
}
func TestSubscriberServer_Get(t *testing.T) {

	t.Run("SubscriberNotFound", func(t *testing.T) {
		var subscriberId = uuid.Nil
		networkClient := &cmocks.NetworkClient{}

		subRepo := &mocks.SubscriberRepo{}

		subRepo.On("Get", subscriberId).Return(nil, gorm.ErrRecordNotFound).Once()

		s := NewSubscriberServer(OrgName, subRepo, nil, nil, orgId, nil, networkClient)
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

		s := NewSubscriberServer(OrgName, subRepo, nil, nil, orgId, nil, networkClient)
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

		s := NewSubscriberServer(OrgName, subRepo, nil, nil, orgId, nil, networkClient)
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

		s := NewSubscriberServer(OrgName, subRepo, nil, nil, orgId, nil, networkClient)
		subResp, err := s.GetByNetwork(context.TODO(), &pb.GetByNetworkRequest{
			NetworkId: networkId})

		assert.Error(t, err)
		assert.Nil(t, subResp)
		subRepo.AssertExpectations(t)
	})
}
