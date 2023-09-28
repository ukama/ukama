package client_test

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/tj/assert"

	"github.com/ukama/ukama/systems/api/api-gateway/mocks"
	"github.com/ukama/ukama/systems/api/api-gateway/pkg/client"
	"github.com/ukama/ukama/systems/common/rest"
	"github.com/ukama/ukama/systems/common/types"
	"github.com/ukama/ukama/systems/common/uuid"
)

const (
	testUuid   = "03cb753f-5e03-4c97-8e47-625115476c72"
	testNodeId = "uk-sa2341-hnode-v0-a1a0"
)

type RoundTripFunc func(req *http.Request) *http.Response

func (r RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return r(req), nil
}

func TestCient_GetNetwork(t *testing.T) {
	netClient := &mocks.NetworkClient{}

	netId := uuid.NewV4()
	netName := "net-1"

	c := client.NewClientsSet(netClient, nil, nil, nil, nil)

	t.Run("NetworkFoundAndStatusCompleted", func(t *testing.T) {
		netClient.On("Get", netId.String()).
			Return(&client.NetworkInfo{
				Id:         netId.String(),
				Name:       netName,
				SyncStatus: types.SyncStatusCompleted.String(),
			}, nil).Once()

		netInfo, err := c.GetNetwork(netId.String())

		assert.NoError(t, err)

		assert.NotNil(t, netInfo)
		assert.Equal(t, netInfo.Id, netId.String())
		assert.Equal(t, netInfo.Name, netName)
	})

	t.Run("NetworkFoundAndStatusPending", func(t *testing.T) {
		netClient.On("Get", netId.String()).
			Return(&client.NetworkInfo{
				Id:         netId.String(),
				Name:       netName,
				SyncStatus: types.SyncStatusPending.String(),
			}, nil).Once()

		netInfo, err := c.GetNetwork(netId.String())

		assert.Error(t, err)
		assert.IsType(t, err, rest.HttpError{})
		assert.Contains(t, err.Error(), "partial")

		assert.NotNil(t, netInfo)
		assert.Equal(t, netInfo.Id, netId.String())
		assert.Equal(t, netInfo.Name, netName)
	})

	t.Run("NetworkFoundAndStatusFailed", func(t *testing.T) {
		netClient.On("Get", netId.String()).
			Return(&client.NetworkInfo{
				Id:         netId.String(),
				Name:       netName,
				SyncStatus: types.SyncStatusFailed.String(),
			}, nil).Once()

		netInfo, err := c.GetNetwork(netId.String())

		assert.Error(t, err)
		assert.IsType(t, err, rest.HttpError{})
		assert.Contains(t, err.Error(), "invalid")

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

	c := client.NewClientsSet(netClient, nil, nil, nil, nil)

	t.Run("NetworkCreated", func(t *testing.T) {
		netClient.On("Add", client.AddNetworkRequest{
			OrgName:          orgName,
			NetName:          netName,
			AllowedCountries: countries,
			AllowedNetworks:  networks,
			PaymentLinks:     paymentLinks,
		}).Return(&client.NetworkInfo{
			Id:               netId.String(),
			Name:             netName,
			AllowedCountries: countries,
			AllowedNetworks:  networks,
			PaymentLinks:     paymentLinks,
		}, nil).Once()

		netInfo, err := c.CreateNetwork(orgName, netName, countries, networks, budget,
			overdraft, trafficPolicy, paymentLinks)

		assert.NoError(t, err)

		assert.Equal(t, netInfo.Id, netId.String())
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

	c := client.NewClientsSet(nil, packageClient, nil, nil, nil)

	t.Run("PackageFoundAndStatusCompleted", func(t *testing.T) {
		packageClient.On("Get", packageId.String()).
			Return(&client.PackageInfo{
				Id:         packageId.String(),
				Name:       pkgName,
				SyncStatus: types.SyncStatusCompleted.String(),
			}, nil).Once()

		pkgInfo, err := c.GetPackage(packageId.String())

		assert.NoError(t, err)

		assert.NotNil(t, pkgInfo)
		assert.Equal(t, pkgInfo.Id, packageId.String())
		assert.Equal(t, pkgInfo.Name, pkgName)
	})

	t.Run("PackageFoundAndStatusPending", func(t *testing.T) {
		packageClient.On("Get", packageId.String()).
			Return(&client.PackageInfo{
				Id:         packageId.String(),
				Name:       pkgName,
				SyncStatus: types.SyncStatusPending.String(),
			}, nil).Once()

		pkgInfo, err := c.GetPackage(packageId.String())

		assert.Error(t, err)
		assert.IsType(t, err, rest.HttpError{})
		assert.Contains(t, err.Error(), "partial")

		assert.NotNil(t, pkgInfo)
		assert.Equal(t, pkgInfo.Id, packageId.String())
		assert.Equal(t, pkgInfo.Name, pkgName)
	})

	t.Run("PackageFoundAndStatusFailed", func(t *testing.T) {
		packageClient.On("Get", packageId.String()).
			Return(&client.PackageInfo{
				Id:         packageId.String(),
				Name:       pkgName,
				SyncStatus: types.SyncStatusFailed.String(),
			}, nil).Once()

		pkgInfo, err := c.GetPackage(packageId.String())

		assert.Error(t, err)
		assert.IsType(t, err, rest.HttpError{})
		assert.Contains(t, err.Error(), "invalid")

		assert.Nil(t, pkgInfo)
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
	duration := uint64(0)
	flatRate := false
	amount := float64(0)
	overdraft := float64(0)
	trafficPolicy := uint32(0)
	networks := []string{""}

	c := client.NewClientsSet(nil, pkgClient, nil, nil, nil)

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
			Duration:      duration,
			Flatrate:      flatRate,
			Amount:        amount,
			Overdraft:     overdraft,
			TrafficPolicy: trafficPolicy,
			Networks:      networks,
		}).Return(&client.PackageInfo{
			Id:            pkgId.String(),
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
			Duration:      duration,
			Flatrate:      flatRate,
			Amount:        amount,
			Overdraft:     overdraft,
			TrafficPolicy: trafficPolicy,
			Networks:      networks,
		}, nil).Once()

		pkgInfo, err := c.AddPackage(pkgName, orgId, ownerId, from, to, baserateId,
			isActive, flatRate, smsVolume, voiceVolume, dataVolume, voiceUnit, dataUnit,
			simType, apn, pType, duration, markup, amount, overdraft, trafficPolicy, networks)

		assert.NoError(t, err)

		assert.Equal(t, pkgInfo.Id, pkgId.String())
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
			Duration:      duration,
			Flatrate:      flatRate,
			Amount:        amount,
			Overdraft:     overdraft,
			TrafficPolicy: trafficPolicy,
			Networks:      networks,
		}).Return(nil, errors.New("some error")).Once()

		pkgInfo, err := c.AddPackage(pkgName, orgId, ownerId, from, to, baserateId,
			isActive, flatRate, smsVolume, voiceVolume, dataVolume, voiceUnit, dataUnit,
			simType, apn, pType, duration, markup, amount, overdraft, trafficPolicy, networks)

		assert.Contains(t, err.Error(), "error")
		assert.Nil(t, pkgInfo)
	})
}

func TestCient_GetSim(t *testing.T) {
	simClient := &mocks.SimClient{}

	simId := uuid.NewV4()
	subscriberId := uuid.NewV4()

	c := client.NewClientsSet(nil, nil, nil, simClient, nil)

	t.Run("SimFoundAndStatusCompleted", func(t *testing.T) {
		simClient.On("Get", simId.String()).
			Return(&client.SimInfo{
				Id:           simId.String(),
				SubscriberId: subscriberId.String(),
				SyncStatus:   types.SyncStatusCompleted.String(),
			}, nil).Once()

		simInfo, err := c.GetSim(simId.String())

		assert.NoError(t, err)

		assert.NotNil(t, simInfo)
		assert.Equal(t, simInfo.Id, simId.String())
		assert.Equal(t, simInfo.SubscriberId, subscriberId.String())
	})

	t.Run("SimFoundAndStatusPending", func(t *testing.T) {
		simClient.On("Get", simId.String()).
			Return(&client.SimInfo{
				Id:           simId.String(),
				SubscriberId: subscriberId.String(),
				SyncStatus:   types.SyncStatusPending.String(),
			}, nil).Once()

		simInfo, err := c.GetSim(simId.String())

		assert.Error(t, err)
		assert.IsType(t, err, rest.HttpError{})
		assert.Contains(t, err.Error(), "partial")

		assert.NotNil(t, simInfo)
		assert.Equal(t, simInfo.Id, simId.String())
		assert.Equal(t, simInfo.SubscriberId, subscriberId.String())
	})

	t.Run("SimFoundAndStatusFailed", func(t *testing.T) {
		simClient.On("Get", simId.String()).
			Return(&client.SimInfo{
				Id:           simId.String(),
				SubscriberId: subscriberId.String(),
				SyncStatus:   types.SyncStatusFailed.String(),
			}, nil).Once()

		simInfo, err := c.GetSim(simId.String())

		assert.Error(t, err)
		assert.IsType(t, err, rest.HttpError{})
		assert.Contains(t, err.Error(), "invalid")

		assert.Nil(t, simInfo)
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

	c := client.NewClientsSet(nil, nil, subscriberClient, simClient, nil)

	t.Run("SimAndSubscriberCreatedAndStatusUpdated", func(t *testing.T) {
		subscriberClient.On("Add", client.AddSubscriberRequest{
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
			Return(&client.SubscriberInfo{
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

		simClient.On("Add", client.AddSimRequest{
			SubscriberId:  subscriberId.String(),
			NetworkId:     networkId.String(),
			PackageId:     packageId.String(),
			SimType:       simType,
			SimToken:      simToken,
			TrafficPolicy: trafficPolicy}).
			Return(&client.SimInfo{
				Id:           simId.String(),
				SubscriberId: subscriberId.String(),
				NetworkId:    networkId.String(),
				// PackageId:     packageId,
				SimType: simType,
				// SimToken:      simToken,
				TrafficPolicy: trafficPolicy,
			}, nil).Once()

		simInfo, err := c.ConfigureSim("", orgId.String(),
			networkId.String(), firstName, lastName, email, phoneNumber, address,
			dob, proofOfID, idSerial, packageId.String(), simType, simToken, trafficPolicy)

		assert.NoError(t, err)

		assert.Equal(t, simInfo.Id, simId.String())
		assert.Equal(t, simInfo.SubscriberId, subscriberId.String())
	})

	t.Run("SimCreatedAndStatusUpdated", func(t *testing.T) {
		subscriberClient.On("Get", subscriberId.String()).
			Return(&client.SubscriberInfo{
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

		simClient.On("Add", client.AddSimRequest{
			SubscriberId:  subscriberId.String(),
			NetworkId:     networkId.String(),
			PackageId:     packageId.String(),
			SimType:       simType,
			SimToken:      simToken,
			TrafficPolicy: trafficPolicy}).
			Return(&client.SimInfo{
				Id:           simId.String(),
				SubscriberId: subscriberId.String(),
				NetworkId:    networkId.String(),
				// PackageId:     packageId,
				SimType: simType,
				// SimToken:      simToken,
				TrafficPolicy: trafficPolicy,
			}, nil).Once()

		simInfo, err := c.ConfigureSim(subscriberId.String(), orgId.String(),
			networkId.String(), firstName, lastName, email, phoneNumber, address,
			dob, proofOfID, idSerial, packageId.String(), simType, simToken, trafficPolicy)

		assert.NoError(t, err)

		assert.Equal(t, simInfo.Id, simId.String())
		assert.Equal(t, simInfo.SubscriberId, subscriberId.String())
	})

	t.Run("SubscriberNotCreated", func(t *testing.T) {
		subscriberClient.On("Add", client.AddSubscriberRequest{
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

		simInfo, err := c.ConfigureSim("", orgId.String(),
			networkId.String(), firstName, lastName, email, phoneNumber, address,
			dob, proofOfID, idSerial, packageId.String(), simType, simToken, trafficPolicy)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error")
		assert.Nil(t, simInfo)
	})

	t.Run("SimNotCreated", func(t *testing.T) {
		subscriberClient.On("Get", subscriberId.String()).
			Return(nil, nil).Once()

		simClient.On("Add", client.AddSimRequest{
			SubscriberId: subscriberId.String(),
			NetworkId:    networkId.String(),
			PackageId:    packageId.String(),
			SimType:      simType,
			SimToken:     simToken,
		}).Return(nil, errors.New("some error")).Once()

		simInfo, err := c.ConfigureSim(subscriberId.String(), orgId.String(),
			networkId.String(), firstName, lastName, email, phoneNumber, address,
			dob, proofOfID, idSerial, packageId.String(), simType, simToken, trafficPolicy)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error")
		assert.Nil(t, simInfo)
	})
}

func TestCient_GetNode(t *testing.T) {
	nodeClient := &mocks.NodeClient{}

	nodeId := "uk-sa2341-hnode-v0-a1a0"
	nodeName := "node-1"

	c := client.NewClientsSet(nil, nil, nil, nil, nodeClient)

	t.Run("NodeFound", func(t *testing.T) {
		nodeClient.On("Get", nodeId).
			Return(&client.NodeInfo{
				Id:   nodeId,
				Name: nodeName,
			}, nil).Once()

		nodeInfo, err := c.GetNode(nodeId)

		assert.NoError(t, err)

		assert.NotNil(t, nodeInfo)
		assert.Equal(t, nodeInfo.Id, nodeId)
		assert.Equal(t, nodeInfo.Name, nodeName)
	})

	t.Run("NodeNotFound", func(t *testing.T) {
		nodeClient.On("Get", nodeId).
			Return(nil,
				fmt.Errorf("GetNode failure: %w",
					client.ErrorStatus{StatusCode: 404})).Once()

		nodeInfo, err := c.GetNode(nodeId)

		assert.Error(t, err)
		assert.IsType(t, err, rest.HttpError{})
		assert.Contains(t, err.Error(), "404")

		assert.Nil(t, nodeInfo)
	})

	t.Run("NodeGetError", func(t *testing.T) {
		nodeClient.On("Get", nodeId).
			Return(nil,
				fmt.Errorf("Some unexpected error")).Once()

		nodeInfo, err := c.GetNode(nodeId)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error")

		assert.Nil(t, nodeInfo)
	})
}

func TestCient_RegisterNode(t *testing.T) {
	nodeClient := &mocks.NodeClient{}

	nodeId := "uk-sa2341-hnode-v0-a1a0"
	nodeName := "node-1"
	orgId := uuid.NewV4()
	state := "pending"

	c := client.NewClientsSet(nil, nil, nil, nil, nodeClient)

	t.Run("NodeRegistered", func(t *testing.T) {
		nodeClient.On("Add", client.AddNodeRequest{
			NodeId: nodeId,
			Name:   nodeName,
			OrgId:  orgId.String(),
			State:  state,
		}).Return(&client.NodeInfo{
			Id:    nodeId,
			Name:  nodeName,
			OrgId: orgId.String(),
			State: state,
		}, nil).Once()

		nodeInfo, err := c.RegisterNode(nodeId, nodeName, orgId.String(), state)

		assert.NoError(t, err)

		assert.Equal(t, nodeInfo.Id, nodeId)
		assert.Equal(t, nodeInfo.Name, nodeName)
	})

	t.Run("NodeNotRegistered", func(t *testing.T) {
		nodeClient.On("Add", client.AddNodeRequest{
			NodeId: nodeId,
			Name:   nodeName,
			OrgId:  orgId.String(),
			State:  state,
		}).Return(nil, errors.New("some error")).Once()

		nodeInfo, err := c.RegisterNode(nodeId, nodeName, orgId.String(), state)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error")
		assert.Nil(t, nodeInfo)
	})
}

func TestCient_AttachNode(t *testing.T) {
	nodeClient := &mocks.NodeClient{}

	ampNodeL := "uk-sa2341-anode-v0-a1a0"
	ampNodeR := "uk-sa2341-anode-v0-a1a1"

	c := client.NewClientsSet(nil, nil, nil, nil, nodeClient)

	t.Run("NodeAttached", func(t *testing.T) {
		nodeId := "uk-sa2341-tnode-v0-a1a0"

		nodeClient.On("Attach", nodeId, client.AttachNodesRequest{
			AmpNodeL: ampNodeL,
			AmpNodeR: ampNodeR,
		}).Return(nil).Once()

		err := c.AttachNode(nodeId, ampNodeL, ampNodeR)

		assert.NoError(t, err)
	})

	t.Run("NodeNotAttached", func(t *testing.T) {
		nodeId := "uk-sa2341-tnode-v0-a1a1"

		nodeClient.On("Attach", nodeId, client.AttachNodesRequest{
			AmpNodeL: ampNodeL,
			AmpNodeR: ampNodeR,
		}).Return(errors.New("some error")).Once()

		err := c.AttachNode(nodeId, ampNodeL, ampNodeR)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error")
	})
}

func TestCient_DetachNode(t *testing.T) {
	nodeClient := &mocks.NodeClient{}

	nodeId := "uk-sa2341-hnode-v0-a1a0"

	c := client.NewClientsSet(nil, nil, nil, nil, nodeClient)

	t.Run("NodeDetached", func(t *testing.T) {
		nodeClient.On("Detach", nodeId).
			Return(nil).Once()

		err := c.DetachNode(nodeId)

		assert.NoError(t, err)
	})

	t.Run("NodeNotFound", func(t *testing.T) {
		nodeClient.On("Detach", nodeId).
			Return(fmt.Errorf("DetachNode failure: %w",
				client.ErrorStatus{StatusCode: 404})).Once()

		err := c.DetachNode(nodeId)

		assert.Error(t, err)
		assert.IsType(t, err, rest.HttpError{})
		assert.Contains(t, err.Error(), "404")
	})

	t.Run("NodeDetachError", func(t *testing.T) {
		nodeClient.On("Detach", nodeId).
			Return(fmt.Errorf("Some unexpected error")).Once()

		err := c.DetachNode(nodeId)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error")
	})
}

func TestCient_AddToSite(t *testing.T) {
	nodeClient := &mocks.NodeClient{}

	networkId := uuid.NewV4().String()
	siteId := uuid.NewV4().String()

	c := client.NewClientsSet(nil, nil, nil, nil, nodeClient)

	t.Run("NodeAdded", func(t *testing.T) {
		nodeId := "uk-sa2341-tnode-v0-a1a0"

		nodeClient.On("AddToSite", nodeId, client.AddToSiteRequest{
			NetworkId: networkId,
			SiteId:    siteId,
		}).Return(nil).Once()

		err := c.AddNodeToSite(nodeId, networkId, siteId)

		assert.NoError(t, err)
	})

	t.Run("NodeNotAdded", func(t *testing.T) {
		nodeId := "uk-sa2341-tnode-v0-a1a1"

		nodeClient.On("AddToSite", nodeId, client.AddToSiteRequest{
			NetworkId: networkId,
			SiteId:    siteId,
		}).Return(errors.New("some error")).Once()

		err := c.AddNodeToSite(nodeId, networkId, siteId)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error")
	})
}

func TestCient_RemoveNodeFromSite(t *testing.T) {
	nodeClient := &mocks.NodeClient{}

	nodeId := "uk-sa2341-hnode-v0-a1a0"

	c := client.NewClientsSet(nil, nil, nil, nil, nodeClient)

	t.Run("NodeRemoved", func(t *testing.T) {
		nodeClient.On("RemoveFromSite", nodeId).
			Return(nil).Once()

		err := c.RemoveNodeFromSite(nodeId)

		assert.NoError(t, err)
	})

	t.Run("NodeNotFound", func(t *testing.T) {
		nodeClient.On("RemoveFromSite", nodeId).
			Return(fmt.Errorf("DetachNode failure: %w",
				client.ErrorStatus{StatusCode: 404})).Once()

		err := c.RemoveNodeFromSite(nodeId)

		assert.Error(t, err)
		assert.IsType(t, err, rest.HttpError{})
		assert.Contains(t, err.Error(), "404")
	})

	t.Run("NodeRemoveError", func(t *testing.T) {
		nodeClient.On("RemoveFromSite", nodeId).
			Return(fmt.Errorf("Some unexpected error")).Once()

		err := c.RemoveNodeFromSite(nodeId)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error")
	})
}

func TestCient_DeleteNode(t *testing.T) {
	nodeClient := &mocks.NodeClient{}

	nodeId := "uk-sa2341-hnode-v0-a1a0"

	c := client.NewClientsSet(nil, nil, nil, nil, nodeClient)

	t.Run("NodeDeleted", func(t *testing.T) {
		nodeClient.On("Delete", nodeId).
			Return(nil).Once()

		err := c.DeleteNode(nodeId)

		assert.NoError(t, err)
	})

	t.Run("NodeNotFound", func(t *testing.T) {
		nodeClient.On("Delete", nodeId).
			Return(fmt.Errorf("DeleteNode failure: %w",
				client.ErrorStatus{StatusCode: 404})).Once()

		err := c.DeleteNode(nodeId)

		assert.Error(t, err)
		assert.IsType(t, err, rest.HttpError{})
		assert.Contains(t, err.Error(), "404")
	})

	t.Run("NodeDeleteError", func(t *testing.T) {
		nodeClient.On("Delete", nodeId).
			Return(fmt.Errorf("Some unexpected error")).Once()

		err := c.DeleteNode(nodeId)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error")
	})
}
