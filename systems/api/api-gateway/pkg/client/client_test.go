package client_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/tj/assert"
	"github.com/ukama/ukama/systems/api/api-gateway/mocks"
	"github.com/ukama/ukama/systems/api/api-gateway/pkg/client"
	"github.com/ukama/ukama/systems/api/api-gateway/pkg/db"
	"github.com/ukama/ukama/systems/common/rest"
	"github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/gorm"
)

func TestCient_GetNetwork(t *testing.T) {
	resRepo := &mocks.ResourceRepo{}
	netClient := &mocks.NetworkClient{}

	netId := uuid.NewV4()
	netName := "net-1"

	c := client.NewClientsSet(resRepo, netClient)

	t.Run("NetworkFoundAndStatusCompleted", func(t *testing.T) {
		netClient.On("Get", netId.String()).
			Return(&client.NetworkInfo{
				Id:   netId,
				Name: netName,
			}, nil).Once()

		resRepo.On("Get", netId).
			Return(&db.Resource{
				Id:     netId,
				Status: db.ResourceStatusCompleted}, nil).Once()

		netInfo, err := c.GetNetwork(netId.String())

		assert.NoError(t, err)

		assert.NotNil(t, netInfo)
		assert.Equal(t, netInfo.Id, netId)
		assert.Equal(t, netInfo.Name, netName)
	})

	t.Run("NetworkFoundAndStatusPending", func(t *testing.T) {
		netClient.On("Get", netId.String()).
			Return(&client.NetworkInfo{
				Id:   netId,
				Name: netName,
			}, nil).Once()

		resRepo.On("Get", netId).
			Return(&db.Resource{
				Id:     netId,
				Status: db.ResourceStatusPending}, nil).Once()

		netInfo, err := c.GetNetwork(netId.String())

		assert.Error(t, err)
		assert.IsType(t, err, rest.HttpError{})
		assert.Contains(t, err.Error(), "partial")

		assert.NotNil(t, netInfo)
		assert.Equal(t, netInfo.Id, netId)
		assert.Equal(t, netInfo.Name, netName)
	})

	t.Run("NetworkFoundAndStatusFailed", func(t *testing.T) {
		netClient.On("Get", netId.String()).
			Return(&client.NetworkInfo{
				Id:   netId,
				Name: netName,
			}, nil).Once()

		resRepo.On("Get", netId).
			Return(&db.Resource{
				Id:     netId,
				Status: db.ResourceStatusFailed}, nil).Once()

		netInfo, err := c.GetNetwork(netId.String())

		assert.Error(t, err)
		assert.IsType(t, err, rest.HttpError{})
		assert.Contains(t, err.Error(), "inconsistent")

		assert.Nil(t, netInfo)
	})

	t.Run("NetworkFoundAndStatusError", func(t *testing.T) {
		netClient.On("Get", netId.String()).
			Return(&client.NetworkInfo{
				Id:   netId,
				Name: netName,
			}, nil).Once()

		resRepo.On("Get", netId).
			Return(nil, gorm.ErrRecordNotFound).Once()

		netInfo, err := c.GetNetwork(netId.String())

		assert.Error(t, err)
		assert.IsType(t, err, rest.HttpError{})
		assert.Contains(t, err.Error(), "inconsistent")

		assert.Nil(t, netInfo)
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
	resRepo := &mocks.ResourceRepo{}
	netClient := &mocks.NetworkClient{}

	netId := uuid.NewV4()
	netName := "net-1"
	orgName := "org-A"

	c := client.NewClientsSet(resRepo, netClient)

	t.Run("NetworkCreatedAndStatusUpdated", func(t *testing.T) {
		netClient.On("Add", client.AddNetworkRequest{
			OrgName: orgName,
			NetName: netName}).Return(&client.NetworkInfo{
			Id:   netId,
			Name: netName,
		}, nil).Once()

		res := &db.Resource{
			Id:     netId,
			Status: db.ParseStatus("pending"),
		}

		resRepo.On("Add", res).
			Return(nil).Once()

		netInfo, err := c.CreateNetwork(orgName, netName)

		assert.NoError(t, err)

		assert.Equal(t, netInfo.Id, netId)
		assert.Equal(t, netInfo.Name, netName)
	})

	t.Run("NetworkCreatedAndStatusNotUpdated", func(t *testing.T) {
		netClient.On("Add", client.AddNetworkRequest{
			OrgName: orgName,
			NetName: netName}).Return(&client.NetworkInfo{
			Id:   netId,
			Name: netName,
		}, nil).Once()

		res := &db.Resource{
			Id:     netId,
			Status: db.ParseStatus("pending"),
		}

		resRepo.On("Add", res).
			Return(errors.New("some unexpected error")).Once()

		netInfo, err := c.CreateNetwork(orgName, netName)

		assert.Contains(t, err.Error(), "error")
		assert.Nil(t, netInfo)
	})

	t.Run("NetworkNotCreated", func(t *testing.T) {
		netClient.On("Add", client.AddNetworkRequest{
			OrgName: orgName,
			NetName: netName}).Return(nil, errors.New("some error")).Once()

		netInfo, err := c.CreateNetwork(orgName, netName)

		assert.Contains(t, err.Error(), "error")
		assert.Nil(t, netInfo)
	})
}
