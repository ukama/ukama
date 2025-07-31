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

	"github.com/ukama/ukama/systems/api/api-gateway/pkg/client"
	"github.com/ukama/ukama/systems/common/mocks"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"

	crest "github.com/ukama/ukama/systems/common/rest"
	cclient "github.com/ukama/ukama/systems/common/rest/client"
	creg "github.com/ukama/ukama/systems/common/rest/client/registry"
)

func TestCient_GetNetwork(t *testing.T) {
	netClient := &mocks.NetworkClient{}

	netId := uuid.NewV4()
	netName := "net-1"

	n := client.NewNetworkClientSet(netClient)

	t.Run("NetworkFoundAndStatusCompleted", func(t *testing.T) {
		netClient.On("Get", netId.String()).
			Return(&creg.NetworkInfo{
				Id:         netId.String(),
				Name:       netName,
				SyncStatus: ukama.StatusTypeCompleted.String(),
			}, nil).Once()

		netInfo, err := n.GetNetwork(netId.String())

		assert.NoError(t, err)

		assert.NotNil(t, netInfo)
		assert.Equal(t, netInfo.Id, netId.String())
		assert.Equal(t, netInfo.Name, netName)
	})

	t.Run("NetworkFoundAndStatusPending", func(t *testing.T) {
		netClient.On("Get", netId.String()).
			Return(&creg.NetworkInfo{
				Id:         netId.String(),
				Name:       netName,
				SyncStatus: ukama.StatusTypePending.String(),
			}, nil).Once()

		netInfo, err := n.GetNetwork(netId.String())

		assert.Error(t, err)
		assert.IsType(t, err, crest.HttpError{})
		assert.Contains(t, err.Error(), "partial")

		assert.NotNil(t, netInfo)
		assert.Equal(t, netInfo.Id, netId.String())
		assert.Equal(t, netInfo.Name, netName)
	})

	t.Run("NetworkFoundAndStatusFailed", func(t *testing.T) {
		netClient.On("Get", netId.String()).
			Return(&creg.NetworkInfo{
				Id:         netId.String(),
				Name:       netName,
				SyncStatus: ukama.StatusTypeFailed.String(),
			}, nil).Once()

		netInfo, err := n.GetNetwork(netId.String())

		assert.Error(t, err)
		assert.IsType(t, err, crest.HttpError{})
		assert.Contains(t, err.Error(), "invalid")

		assert.Nil(t, netInfo)
	})

	t.Run("NetworkNotFound", func(t *testing.T) {
		netClient.On("Get", netId.String()).
			Return(nil,
				fmt.Errorf("getNetwork failure: %w",
					&cclient.ErrorStatus{
						StatusCode: 404,
						RawError:   crest.ErrorResponse{Err: "not found"},
					})).Once()

		netInfo, err := n.GetNetwork(netId.String())

		assert.Error(t, err)
		assert.IsType(t, err, crest.HttpError{})
		assert.Contains(t, err.Error(), "404")

		assert.Nil(t, netInfo)
	})

	t.Run("NetworkGetError", func(t *testing.T) {
		netClient.On("Get", netId.String()).
			Return(nil,
				fmt.Errorf("Some unexpected error")).Once()

		netInfo, err := n.GetNetwork(netId.String())

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error")

		assert.Nil(t, netInfo)
	})
}

func TestCient_CreateNetwork(t *testing.T) {
	netClient := &mocks.NetworkClient{}

	netId := uuid.NewV4()
	netName := "net-1"
	orgName := "org-A"
	networks := []string{"Verizon"}
	countries := []string{"USA"}
	budget := float64(0)
	overdraft := float64(0)
	trafficPolicy := uint32(0)
	paymentLinks := false

	n := client.NewNetworkClientSet(netClient)

	t.Run("NetworkCreated", func(t *testing.T) {
		netClient.On("Add", creg.AddNetworkRequest{
			OrgName:          orgName,
			NetName:          netName,
			AllowedCountries: countries,
			AllowedNetworks:  networks,
			PaymentLinks:     paymentLinks,
		}).Return(&creg.NetworkInfo{
			Id:               netId.String(),
			Name:             netName,
			AllowedCountries: countries,
			AllowedNetworks:  networks,
			PaymentLinks:     paymentLinks,
		}, nil).Once()

		netInfo, err := n.CreateNetwork(orgName, netName, countries, networks, budget,
			overdraft, trafficPolicy, paymentLinks)

		assert.NoError(t, err)

		assert.Equal(t, netInfo.Id, netId.String())
		assert.Equal(t, netInfo.Name, netName)
	})

	t.Run("NetworkNotCreated", func(t *testing.T) {
		netClient.On("Add", creg.AddNetworkRequest{
			OrgName:          orgName,
			NetName:          netName,
			AllowedCountries: countries,
			AllowedNetworks:  networks,
			PaymentLinks:     paymentLinks,
		}).Return(nil, errors.New("some error")).Once()

		netInfo, err := n.CreateNetwork(orgName, netName, countries, networks, budget,
			overdraft, trafficPolicy, paymentLinks)

		assert.Contains(t, err.Error(), "error")
		assert.Nil(t, netInfo)
	})
}
