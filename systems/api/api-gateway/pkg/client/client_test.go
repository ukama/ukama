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

	c := client.NewClientsSet(netClient, nil, nil, nil)

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

func TestCient_CreateNetwork(t *testing.T) {
	netClient := &mocks.NetworkClient{}

	netId := uuid.NewV4()
	netName := "net-1"
	orgName := "org-A"
	networks := []string{"Verizon"}
	countries := []string{"USA"}
	budget := float64(0)
	overdraft := float64(0)
	trafficPolicy := uint(0)
	paymentLinks := false

	c := client.NewClientsSet(netClient, nil, nil, nil)

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

		netInfo, err := c.CreateNetwork(orgName, netName, countries, networks, budget,
			overdraft, trafficPolicy, paymentLinks)

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

		netInfo, err := c.CreateNetwork(orgName, netName, countries, networks, budget,
			overdraft, trafficPolicy, paymentLinks)

		assert.Contains(t, err.Error(), "error")
		assert.Nil(t, netInfo)
	})
}

func TestCient_GetPackage(t *testing.T) {
	packageClient := &mocks.PackageClient{}

	packageId := uuid.NewV4()
	pkgName := "Monthly Data"

	c := client.NewClientsSet(nil, packageClient, nil, nil)

	t.Run("PackageFoundAndStatusCompleted", func(t *testing.T) {
		packageClient.On("Get", packageId.String()).
			Return(&client.PackageInfo{
				Id:       packageId,
				Name:     pkgName,
				IsSynced: true,
			}, nil).Once()

		pkgInfo, err := c.GetPackage(packageId.String())

		assert.NoError(t, err)

		assert.NotNil(t, pkgInfo)
		assert.Equal(t, pkgInfo.Id, packageId)
		assert.Equal(t, pkgInfo.Name, pkgName)
	})

	t.Run("PackageFoundAndStatusPending", func(t *testing.T) {
		packageClient.On("Get", packageId.String()).
			Return(&client.PackageInfo{
				Id:       packageId,
				Name:     pkgName,
				IsSynced: false,
			}, nil).Once()

		pkgInfo, err := c.GetPackage(packageId.String())

		assert.Error(t, err)
		assert.IsType(t, err, rest.HttpError{})
		assert.Contains(t, err.Error(), "partial")

		assert.NotNil(t, pkgInfo)
		assert.Equal(t, pkgInfo.Id, packageId)
		assert.Equal(t, pkgInfo.Name, pkgName)
	})

	t.Run("PackageNotFound", func(t *testing.T) {
		packageClient.On("Get", packageId.String()).
			Return(nil,
				fmt.Errorf("GetNetwork failure: %w",
					client.ErrorStatus{StatusCode: 404})).Once()

		pkgInfo, err := c.GetPackage(packageId.String())

		assert.Error(t, err)
		assert.IsType(t, err, rest.HttpError{})
		assert.Contains(t, err.Error(), "404")

		assert.Nil(t, pkgInfo)
	})

	t.Run("PackageGetError", func(t *testing.T) {
		packageClient.On("Get", packageId.String()).
			Return(nil,
				fmt.Errorf("Some unexpected error")).Once()

		pkgInfo, err := c.GetPackage(packageId.String())

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error")

		assert.Nil(t, pkgInfo)
	})
}

func TestCient_AddPackage(t *testing.T) {
	pkgClient := &mocks.PackageClient{}

	pkgId := uuid.NewV4()
	pkgName := "Monthly Data"
	orgId := uuid.NewV4().String()
	ownerId := uuid.NewV4().String()
	from := "2023-04-01T00:00:00Z"
	to := "2023-04-01T00:00:00Z"
	baserateId := uuid.NewV4().String()
	voiceVolume := int64(0)
	isActive := true
	dataVolume := int64(1024)
	smsVolume := int64(0)
	dataUnit := "MegaBytes"
	voiceUnit := "seconds"
	simType := "test"
	apn := "ukama.tel"
	markup := float64(0)
	pType := "postpaid"
	flatRate := false
	amount := float64(0)
	overdraft := float64(0)
	trafficPolicy := uint(0)
	networks := []string{""}

	c := client.NewClientsSet(nil, pkgClient, nil, nil)

	t.Run("PackageCreatedAndStatusUpdated", func(t *testing.T) {
		pkgClient.On("Add", client.AddPackageRequest{
			Name:          pkgName,
			OrgId:         orgId,
			OwnerId:       ownerId,
			From:          from,
			To:            to,
			BaserateId:    baserateId,
			Active:        isActive,
			SmsVolume:     smsVolume,
			VoiceVolume:   voiceVolume,
			DataVolume:    dataVolume,
			VoiceUnit:     voiceUnit,
			DataUnit:      dataUnit,
			SimType:       simType,
			Apn:           apn,
			Markup:        markup,
			Type:          pType,
			Flatrate:      flatRate,
			Amount:        amount,
			Overdraft:     overdraft,
			TrafficPolicy: trafficPolicy,
			Networks:      networks,
		}).Return(&client.PackageInfo{
			Id:            pkgId,
			Name:          pkgName,
			OrgId:         orgId,
			OwnerId:       ownerId,
			From:          from,
			To:            to,
			BaserateId:    baserateId,
			IsActive:      isActive,
			SmsVolume:     smsVolume,
			VoiceVolume:   voiceVolume,
			DataVolume:    dataVolume,
			VoiceUnit:     voiceUnit,
			DataUnit:      dataUnit,
			SimType:       simType,
			Apn:           apn,
			Markup:        markup,
			Type:          pType,
			Flatrate:      flatRate,
			Amount:        amount,
			Overdraft:     overdraft,
			TrafficPolicy: trafficPolicy,
			Networks:      networks,
			IsSynced:      false,
		}, nil).Once()

		pkgInfo, err := c.AddPackage(pkgName, orgId, ownerId, from, to, baserateId,
			isActive, flatRate, smsVolume, voiceVolume, dataVolume, voiceUnit, dataUnit,
			simType, apn, pType, markup, amount, overdraft, trafficPolicy, networks)

		assert.NoError(t, err)

		assert.Equal(t, pkgInfo.Id, pkgId)
		assert.Equal(t, pkgInfo.Name, pkgName)
	})

	t.Run("PackageNotCreated", func(t *testing.T) {
		pkgClient.On("Add", client.AddPackageRequest{
			Name:          pkgName,
			OrgId:         orgId,
			OwnerId:       ownerId,
			From:          from,
			To:            to,
			BaserateId:    baserateId,
			Active:        isActive,
			SmsVolume:     smsVolume,
			VoiceVolume:   voiceVolume,
			DataVolume:    dataVolume,
			VoiceUnit:     voiceUnit,
			DataUnit:      dataUnit,
			SimType:       simType,
			Apn:           apn,
			Markup:        markup,
			Type:          pType,
			Flatrate:      flatRate,
			Amount:        amount,
			Overdraft:     overdraft,
			TrafficPolicy: trafficPolicy,
			Networks:      networks,
		}).Return(nil, errors.New("some error")).Once()

		pkgInfo, err := c.AddPackage(pkgName, orgId, ownerId, from, to, baserateId,
			isActive, flatRate, smsVolume, voiceVolume, dataVolume, voiceUnit, dataUnit,
			simType, apn, pType, markup, amount, overdraft, trafficPolicy, networks)

		assert.Contains(t, err.Error(), "error")
		assert.Nil(t, pkgInfo)
	})
}

func TestCient_GetSim(t *testing.T) {
	simClient := &mocks.SimClient{}

	simId := uuid.NewV4()
	subscriberId := uuid.NewV4()

	c := client.NewClientsSet(nil, nil, nil, simClient)

	t.Run("SimFoundAndStatusCompleted", func(t *testing.T) {
		simClient.On("Get", simId.String()).
			Return(&client.SimInfo{
				Id:           simId,
				SubscriberId: subscriberId,
				IsSynced:     true,
			}, nil).Once()

		simInfo, err := c.GetSim(simId.String())

		assert.NoError(t, err)

		assert.NotNil(t, simInfo)
		assert.Equal(t, simInfo.Id, simId)
		assert.Equal(t, simInfo.SubscriberId, subscriberId)
	})

	t.Run("SimFoundAndStatusPending", func(t *testing.T) {
		simClient.On("Get", simId.String()).
			Return(&client.SimInfo{
				Id:           simId,
				SubscriberId: subscriberId,
				IsSynced:     false,
			}, nil).Once()

		simInfo, err := c.GetSim(simId.String())

		assert.Error(t, err)
		assert.IsType(t, err, rest.HttpError{})
		assert.Contains(t, err.Error(), "partial")

		assert.NotNil(t, simInfo)
		assert.Equal(t, simInfo.Id, simId)
		assert.Equal(t, simInfo.SubscriberId, subscriberId)
	})

	t.Run("SimNotFound", func(t *testing.T) {
		simClient.On("Get", simId.String()).
			Return(nil,
				fmt.Errorf("GetSim failure: %w",
					client.ErrorStatus{StatusCode: 404})).Once()

		simInfo, err := c.GetSim(simId.String())

		assert.Error(t, err)
		assert.IsType(t, err, rest.HttpError{})
		assert.Contains(t, err.Error(), "404")

		assert.Nil(t, simInfo)
	})

	t.Run("SimGetError", func(t *testing.T) {
		simClient.On("Get", simId.String()).
			Return(nil,
				fmt.Errorf("Some unexpected error")).Once()

		simInfo, err := c.GetSim(simId.String())

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error")

		assert.Nil(t, simInfo)
	})
}

func TestCient_ConfigureSim(t *testing.T) {
	simClient := &mocks.SimClient{}

	simId := uuid.NewV4()
	subscriberId := uuid.NewV4()
	networkId := uuid.NewV4()
	packageId := uuid.NewV4()
	simType := "some-sim-type"
	simToken := "some-sim-token"
	trafficPolicy := uint(0)

	c := client.NewClientsSet(nil, nil, nil, simClient)

	t.Run("SimCreatedAndStatusUpdated", func(t *testing.T) {
		simClient.On("Add", client.AddSimRequest{
			SubscriberId:  subscriberId.String(),
			NetworkId:     networkId.String(),
			PackageId:     packageId.String(),
			SimType:       simType,
			SimToken:      simToken,
			TrafficPolicy: trafficPolicy,
		}).Return(&client.SimInfo{
			Id:           simId,
			SubscriberId: subscriberId,
			NetworkId:    networkId,
			// PackageId:     packageId,
			SimType: simType,
			// SimToken:      simToken,
			TrafficPolicy: trafficPolicy,
			IsSynced:      false,
		}, nil).Once()

		simInfo, err := c.ConfigureSim(subscriberId.String(),
			networkId.String(), packageId.String(), simType, simToken, trafficPolicy)

		assert.NoError(t, err)

		assert.Equal(t, simInfo.Id, simId)
		assert.Equal(t, simInfo.SubscriberId, subscriberId)
	})

	t.Run("SimNotCreated", func(t *testing.T) {
		simClient.On("Add", client.AddSimRequest{
			SubscriberId: subscriberId.String(),
			NetworkId:    networkId.String(),
			PackageId:    packageId.String(),
			SimType:      simType,
			SimToken:     simToken,
		}).Return(nil, errors.New("some error")).Once()

		simInfo, err := c.ConfigureSim(subscriberId.String(),
			networkId.String(), packageId.String(), simType, simToken, trafficPolicy)

		assert.Contains(t, err.Error(), "error")
		assert.Nil(t, simInfo)
	})
}
