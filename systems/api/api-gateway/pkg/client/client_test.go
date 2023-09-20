package client_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/tj/assert"
	"github.com/ukama/ukama/systems/api/api-gateway/mocks"
	"github.com/ukama/ukama/systems/api/api-gateway/pkg/client"
	"github.com/ukama/ukama/systems/common/rest"
	"github.com/ukama/ukama/systems/common/uuid"
)

func TestCient_GetNetwork(t *testing.T) {
	netClient := &mocks.NetworkClient{}

	netId := uuid.NewV4()
	netName := "net-1"

	c := client.NewClientsSet(netClient)

	t.Run("NetworkFoundAndStatusCompleted", func(t *testing.T) {
		netClient.On("Get", netId.String()).
			Return(&client.NetworkInfo{
				Id:       netId,
				Name:     netName,
				IsSynced: true,
			}, nil).Once()

		netInfo, err := c.GetNetwork(netId.String())

		assert.NoError(t, err)

		assert.NotNil(t, netInfo)
		assert.Equal(t, netInfo.Id, netId)
		assert.Equal(t, netInfo.Name, netName)
	})

	t.Run("NetworkFoundAndStatusPending", func(t *testing.T) {
		netClient.On("Get", netId.String()).
			Return(&client.NetworkInfo{
				Id:       netId,
				Name:     netName,
				IsSynced: false,
			}, nil).Once()

		netInfo, err := c.GetNetwork(netId.String())

		assert.Error(t, err)
		assert.IsType(t, err, rest.HttpError{})
		assert.Contains(t, err.Error(), "partial")

		assert.NotNil(t, netInfo)
		assert.Equal(t, netInfo.Id, netId)
		assert.Equal(t, netInfo.Name, netName)
	})

	t.Run("NetworkNotFound", func(t *testing.T) {
		netClient.On("Get", netId.String()).
			Return(nil,
				fmt.Errorf("GetNetwork failure: %w",
					client.ErrorStatus{StatusCode: 404})).Once()

		netInfo, err := c.GetNetwork(netId.String())

		assert.Error(t, err)
		assert.IsType(t, err, rest.HttpError{})
		assert.Contains(t, err.Error(), "404")

		assert.Nil(t, netInfo)
	})

	t.Run("NetworkGetError", func(t *testing.T) {
		netClient.On("Get", netId.String()).
			Return(nil,
				fmt.Errorf("Some unexpected error")).Once()

		netInfo, err := c.GetNetwork(netId.String())

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error")

		assert.Nil(t, netInfo)
	})

}

func TestCient_AddNetwork(t *testing.T) {
	netClient := &mocks.NetworkClient{}

	netId := uuid.NewV4()
	netName := "net-1"
	orgName := "org-A"
	networks := []string{"Verizon"}
	countries := []string{"USA"}
	paymentLinks := false

	c := client.NewClientsSet(netClient)

	t.Run("NetworkCreatedAndStatusUpdated", func(t *testing.T) {
		netClient.On("Add", client.AddNetworkRequest{
			OrgName:          orgName,
			NetName:          netName,
			AllowedCountries: countries,
			AllowedNetworks:  networks,
			PaymentLinks:     paymentLinks,
		}).Return(&client.NetworkInfo{
			Id:               netId,
			Name:             netName,
			AllowedCountries: countries,
			AllowedNetworks:  networks,
			PaymentLinks:     paymentLinks,
			IsSynced:         false,
		}, nil).Once()

		netInfo, err := c.CreateNetwork(orgName, netName, countries, networks, paymentLinks)

		assert.NoError(t, err)

		assert.Equal(t, netInfo.Id, netId)
		assert.Equal(t, netInfo.Name, netName)
	})

	t.Run("NetworkNotCreated", func(t *testing.T) {
		netClient.On("Add", client.AddNetworkRequest{
			OrgName:          orgName,
			NetName:          netName,
			AllowedCountries: countries,
			AllowedNetworks:  networks,
			PaymentLinks:     paymentLinks,
		}).Return(nil, errors.New("some error")).Once()

		netInfo, err := c.CreateNetwork(orgName, netName, countries, networks, paymentLinks)

		assert.Contains(t, err.Error(), "error")
		assert.Nil(t, netInfo)
	})
}
