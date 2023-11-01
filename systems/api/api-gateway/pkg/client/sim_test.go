/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package client_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/tj/assert"

	"github.com/ukama/ukama/systems/api/api-gateway/mocks"
	"github.com/ukama/ukama/systems/api/api-gateway/pkg/client"
	"github.com/ukama/ukama/systems/api/api-gateway/pkg/client/rest"
	"github.com/ukama/ukama/systems/common/types"
	"github.com/ukama/ukama/systems/common/uuid"

	crest "github.com/ukama/ukama/systems/common/rest"
)

func TestCient_GetSim(t *testing.T) {
	simClient := &mocks.SimClient{}

	simId := uuid.NewV4()
	subscriberId := uuid.NewV4()

	s := client.NewSimClientSet(simClient, nil)

	t.Run("SimFoundAndStatusCompleted", func(t *testing.T) {
		simClient.On("Get", simId.String()).
			Return(&rest.SimInfo{
				Id:           simId.String(),
				SubscriberId: subscriberId.String(),
				SyncStatus:   types.SyncStatusCompleted.String(),
			}, nil).Once()

		simInfo, err := s.GetSim(simId.String())

		assert.NoError(t, err)

		assert.NotNil(t, simInfo)
		assert.Equal(t, simInfo.Id, simId.String())
		assert.Equal(t, simInfo.SubscriberId, subscriberId.String())
	})

	t.Run("SimFoundAndStatusPending", func(t *testing.T) {
		simClient.On("Get", simId.String()).
			Return(&rest.SimInfo{
				Id:           simId.String(),
				SubscriberId: subscriberId.String(),
				SyncStatus:   types.SyncStatusPending.String(),
			}, nil).Once()

		simInfo, err := s.GetSim(simId.String())

		assert.Error(t, err)
		assert.IsType(t, err, crest.HttpError{})
		assert.Contains(t, err.Error(), "partial")

		assert.NotNil(t, simInfo)
		assert.Equal(t, simInfo.Id, simId.String())
		assert.Equal(t, simInfo.SubscriberId, subscriberId.String())
	})

	t.Run("SimFoundAndStatusFailed", func(t *testing.T) {
		simClient.On("Get", simId.String()).
			Return(&rest.SimInfo{
				Id:           simId.String(),
				SubscriberId: subscriberId.String(),
				SyncStatus:   types.SyncStatusFailed.String(),
			}, nil).Once()

		simInfo, err := s.GetSim(simId.String())

		assert.Error(t, err)
		assert.IsType(t, err, crest.HttpError{})
		assert.Contains(t, err.Error(), "invalid")

		assert.Nil(t, simInfo)
	})

	t.Run("SimNotFound", func(t *testing.T) {
		simClient.On("Get", simId.String()).
			Return(nil,
				fmt.Errorf("GetSim failure: %w",
					rest.ErrorStatus{StatusCode: 404})).Once()

		simInfo, err := s.GetSim(simId.String())

		assert.Error(t, err)
		assert.IsType(t, err, crest.HttpError{})
		assert.Contains(t, err.Error(), "404")

		assert.Nil(t, simInfo)
	})

	t.Run("SimGetError", func(t *testing.T) {
		simClient.On("Get", simId.String()).
			Return(nil,
				fmt.Errorf("Some unexpected error")).Once()

		simInfo, err := s.GetSim(simId.String())

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error")

		assert.Nil(t, simInfo)
	})
}

func TestCient_ConfigureSim(t *testing.T) {
	simClient := &mocks.SimClient{}
	subscriberClient := &mocks.SubscriberClient{}

	simId := uuid.NewV4()
	subscriberId := uuid.NewV4()
	networkId := uuid.NewV4()
	packageId := uuid.NewV4()
	simType := "some-sim-type"
	simToken := "some-sim-token"
	trafficPolicy := uint32(0)

	orgId := uuid.NewV4()
	firstName := "John"
	lastName := "Doe"
	email := "johndoe@example.com"
	phoneNumber := "0123456789"
	address := "2 Rivers"
	dob := "2023/09/01"
	proofOfID := "passport"
	idSerial := "987654321"

	s := client.NewSimClientSet(simClient, subscriberClient)

	t.Run("SimAndSubscriberCreatedAndStatusUpdated", func(t *testing.T) {
		subscriberClient.On("Add", rest.AddSubscriberRequest{
			OrgId:                 orgId.String(),
			NetworkId:             networkId.String(),
			FirstName:             firstName,
			LastName:              lastName,
			Email:                 email,
			PhoneNumber:           phoneNumber,
			Address:               address,
			Dob:                   dob,
			ProofOfIdentification: proofOfID,
			IdSerial:              idSerial,
		}).
			Return(&rest.SubscriberInfo{
				SubscriberId:          subscriberId,
				OrgId:                 orgId,
				NetworkId:             networkId,
				FirstName:             firstName,
				LastName:              lastName,
				Email:                 email,
				PhoneNumber:           phoneNumber,
				Address:               address,
				Dob:                   dob,
				ProofOfIdentification: proofOfID,
				IdSerial:              idSerial,
			}, nil).Once()

		simClient.On("Add", rest.AddSimRequest{
			SubscriberId:  subscriberId.String(),
			NetworkId:     networkId.String(),
			PackageId:     packageId.String(),
			SimType:       simType,
			SimToken:      simToken,
			TrafficPolicy: trafficPolicy}).
			Return(&rest.SimInfo{
				Id:           simId.String(),
				SubscriberId: subscriberId.String(),
				NetworkId:    networkId.String(),
				// PackageId:     packageId,
				SimType: simType,
				// SimToken:      simToken,
				TrafficPolicy: trafficPolicy,
			}, nil).Once()

		simInfo, err := s.ConfigureSim("", orgId.String(),
			networkId.String(), firstName, lastName, email, phoneNumber, address,
			dob, proofOfID, idSerial, packageId.String(), simType, simToken, trafficPolicy)

		assert.NoError(t, err)

		assert.Equal(t, simInfo.Id, simId.String())
		assert.Equal(t, simInfo.SubscriberId, subscriberId.String())
	})

	t.Run("SimCreatedAndStatusUpdated", func(t *testing.T) {
		subscriberClient.On("Get", subscriberId.String()).
			Return(&rest.SubscriberInfo{
				SubscriberId:          subscriberId,
				OrgId:                 orgId,
				NetworkId:             networkId,
				FirstName:             firstName,
				LastName:              lastName,
				Email:                 email,
				PhoneNumber:           phoneNumber,
				Address:               address,
				Dob:                   dob,
				ProofOfIdentification: proofOfID,
				IdSerial:              idSerial,
			}, nil).Once()

		simClient.On("Add", rest.AddSimRequest{
			SubscriberId:  subscriberId.String(),
			NetworkId:     networkId.String(),
			PackageId:     packageId.String(),
			SimType:       simType,
			SimToken:      simToken,
			TrafficPolicy: trafficPolicy}).
			Return(&rest.SimInfo{
				Id:           simId.String(),
				SubscriberId: subscriberId.String(),
				NetworkId:    networkId.String(),
				// PackageId:     packageId,
				SimType: simType,
				// SimToken:      simToken,
				TrafficPolicy: trafficPolicy,
			}, nil).Once()

		simInfo, err := s.ConfigureSim(subscriberId.String(), orgId.String(),
			networkId.String(), firstName, lastName, email, phoneNumber, address,
			dob, proofOfID, idSerial, packageId.String(), simType, simToken, trafficPolicy)

		assert.NoError(t, err)

		assert.Equal(t, simInfo.Id, simId.String())
		assert.Equal(t, simInfo.SubscriberId, subscriberId.String())
	})

	t.Run("SubscriberNotCreated", func(t *testing.T) {
		subscriberClient.On("Add", rest.AddSubscriberRequest{
			OrgId:                 orgId.String(),
			NetworkId:             networkId.String(),
			FirstName:             firstName,
			LastName:              lastName,
			Email:                 email,
			PhoneNumber:           phoneNumber,
			Address:               address,
			Dob:                   dob,
			ProofOfIdentification: proofOfID,
			IdSerial:              idSerial,
		}).Return(nil, errors.New("some error")).Once()

		simInfo, err := s.ConfigureSim("", orgId.String(),
			networkId.String(), firstName, lastName, email, phoneNumber, address,
			dob, proofOfID, idSerial, packageId.String(), simType, simToken, trafficPolicy)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error")
		assert.Nil(t, simInfo)
	})

	t.Run("SimNotCreated", func(t *testing.T) {
		subscriberClient.On("Get", subscriberId.String()).
			Return(nil, nil).Once()

		simClient.On("Add", rest.AddSimRequest{
			SubscriberId: subscriberId.String(),
			NetworkId:    networkId.String(),
			PackageId:    packageId.String(),
			SimType:      simType,
			SimToken:     simToken,
		}).Return(nil, errors.New("some error")).Once()

		simInfo, err := s.ConfigureSim(subscriberId.String(), orgId.String(),
			networkId.String(), firstName, lastName, email, phoneNumber, address,
			dob, proofOfID, idSerial, packageId.String(), simType, simToken, trafficPolicy)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error")
		assert.Nil(t, simInfo)
	})
}
