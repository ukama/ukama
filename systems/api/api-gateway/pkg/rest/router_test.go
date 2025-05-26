/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package rest

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/ukama/ukama/systems/api/api-gateway/mocks"
	"github.com/ukama/ukama/systems/common/providers"
	"github.com/ukama/ukama/systems/common/uuid"

	cconfig "github.com/ukama/ukama/systems/common/config"
	crest "github.com/ukama/ukama/systems/common/rest"
	cdplan "github.com/ukama/ukama/systems/common/rest/client/dataplan"
	creg "github.com/ukama/ukama/systems/common/rest/client/registry"
	csub "github.com/ukama/ukama/systems/common/rest/client/subscriber"
)

const (
	netEndpoint  = "/v1/networks"
	pkgEndpoint  = "/v1/packages"
	simEndpoint  = "/v1/sims"
	nodeEndpoint = "/v1/nodes"
)

var (
	netClient     = &mocks.Network{}
	packageClient = &mocks.Package{}
	simClient     = &mocks.Sim{}
	nodeClient    = &mocks.Node{}
)

var defaultCors = cors.Config{
	AllowAllOrigins: true,
}

var routerConfig = &RouterConfig{
	serverConf: &crest.HttpConfig{
		Cors: defaultCors,
	},
	auth: &cconfig.Auth{
		AuthAppUrl:    "http://localhost:4455",
		AuthServerUrl: "http://localhost:4434",
		AuthAPIGW:     "http://localhost:8080",
	},
}

func init() {
	gin.SetMode(gin.TestMode)
}

func TestRouter_PingRoute(t *testing.T) {
	var arc = &providers.AuthRestClient{}

	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)

	r := NewRouter(netClient, packageClient, simClient, nodeClient, routerConfig,
		arc.MockAuthenticateUser).f.Engine()

	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "pong")
}

func TestRouter_GetNetwork(t *testing.T) {
	arc := &providers.AuthRestClient{}

	netName := "net-1"

	t.Run("NetworkFoundAndStatusCompleted", func(t *testing.T) {
		netId := uuid.NewV4()

		netInfo := &creg.NetworkInfo{
			Id:   netId.String(),
			Name: netName,
		}

		netClient.On("GetNetwork", netId.String()).Return(netInfo, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/%s", netEndpoint, netId), nil)

		r := NewRouter(netClient, nil, nil, nil, routerConfig,
			arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusOK, w.Code)
		netClient.AssertExpectations(t)
	})

	t.Run("NetworkFoundAndStatusPending", func(t *testing.T) {
		netId := uuid.NewV4()

		netInfo := &creg.NetworkInfo{
			Id:   netId.String(),
			Name: netName,
		}

		netClient.On("GetNetwork", netId.String()).Return(netInfo,
			crest.HttpError{
				HttpCode: http.StatusPartialContent,
				Message:  "partial content. request is still ongoing",
			})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/%s", netEndpoint, netId), nil)

		r := NewRouter(netClient, nil, nil, nil, routerConfig,
			arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusPartialContent, w.Code)
		netClient.AssertExpectations(t)
	})

	t.Run("NetworkNotFound", func(t *testing.T) {
		netId := uuid.NewV4()

		netClient.On("GetNetwork", netId.String()).Return(nil,
			crest.HttpError{
				HttpCode: http.StatusNotFound,
				Message:  "GetNetwork failure",
			})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/%s", netEndpoint, netId), nil)

		r := NewRouter(netClient, nil, nil, nil, routerConfig,
			arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusNotFound, w.Code)
		netClient.AssertExpectations(t)
	})

	t.Run("GetNetworkGetError", func(t *testing.T) {
		netId := uuid.NewV4()

		netClient.On("GetNetwork", netId.String()).Return(nil,
			errors.New("some unexpected error"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/%s", netEndpoint, netId), nil)

		r := NewRouter(netClient, nil, nil, nil, routerConfig,
			arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		netClient.AssertExpectations(t)
	})
}

func TestRouter_CreateNetwork(t *testing.T) {
	arc := &providers.AuthRestClient{}

	t.Run("NetworkCreatedAndStatusUpdated", func(t *testing.T) {
		netId := uuid.NewV4()
		netName := "net-1"
		orgName := "org-A"
		networks := []string{"Verizon"}
		countries := []string{"USA"}
		budget := float64(0)
		overdraft := float64(0)
		trafficPolicy := uint32(0)
		paymentLinks := false

		var ntwk = AddNetworkReq{
			OrgName:          orgName,
			NetName:          netName,
			AllowedCountries: countries,
			AllowedNetworks:  networks,
			Budget:           budget,
			Overdraft:        overdraft,
			TrafficPolicy:    trafficPolicy,
			PaymentLinks:     paymentLinks,
		}

		netInfo := &creg.NetworkInfo{
			Id:   netId.String(),
			Name: netName,
		}

		body, err := json.Marshal(ntwk)
		if err != nil {
			t.Errorf("fail to marshal request data: %v. Error: %v", ntwk, err)
		}

		netClient.On("CreateNetwork", orgName, netName, countries, networks, budget, overdraft, trafficPolicy, paymentLinks).
			Return(netInfo, nil)

		r := NewRouter(netClient, nil, nil, nil, routerConfig,
			arc.MockAuthenticateUser).f.Engine()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", netEndpoint, bytes.NewReader(body))

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusPartialContent, w.Code)
		netClient.AssertExpectations(t)
	})

	t.Run("NetworkCreatedAndStatusFailed", func(t *testing.T) {
		netName := "net-2"
		orgName := "org-B"
		networks := []string{"Verizon"}
		countries := []string{"USA"}
		budget := float64(0)
		overdraft := float64(0)
		trafficPolicy := uint32(0)
		paymentLinks := false

		var ntwk = AddNetworkReq{
			OrgName:          orgName,
			NetName:          netName,
			AllowedCountries: countries,
			AllowedNetworks:  networks,
			Budget:           budget,
			Overdraft:        overdraft,
			TrafficPolicy:    trafficPolicy,
			PaymentLinks:     paymentLinks,
		}

		body, err := json.Marshal(ntwk)
		if err != nil {
			t.Errorf("fail to marshal request data: %v. Error: %v", ntwk, err)
		}

		netClient.On("CreateNetwork", orgName, netName, countries, networks, budget, overdraft, trafficPolicy, paymentLinks).
			Return(nil, errors.New("some unexpected error occured"))

		r := NewRouter(netClient, nil, nil, nil, routerConfig,
			arc.MockAuthenticateUser).f.Engine()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", netEndpoint, bytes.NewReader(body))

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		netClient.AssertExpectations(t)
	})
}

func TestRouter_GetPackage(t *testing.T) {
	arc := &providers.AuthRestClient{}

	pkgName := "Monthly Data"

	t.Run("PackageFoundAndStatusCompleted", func(t *testing.T) {
		pkgId := uuid.NewV4()

		pkgInfo := &cdplan.PackageInfo{
			Id:   pkgId.String(),
			Name: pkgName,
		}

		packageClient.On("GetPackage", pkgId.String()).Return(pkgInfo, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/%s", pkgEndpoint, pkgId), nil)

		r := NewRouter(nil, packageClient, nil, nil, routerConfig,
			arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusOK, w.Code)
		packageClient.AssertExpectations(t)
	})

	t.Run("PackageFoundAndStatusPending", func(t *testing.T) {
		pkgId := uuid.NewV4()

		pkgInfo := &cdplan.PackageInfo{
			Id:   pkgId.String(),
			Name: pkgName,
		}

		packageClient.On("GetPackage", pkgId.String()).Return(pkgInfo,
			crest.HttpError{
				HttpCode: http.StatusPartialContent,
				Message:  "partial content. request is still ongoing",
			})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/%s", pkgEndpoint, pkgId), nil)

		r := NewRouter(nil, packageClient, nil, nil, routerConfig,
			arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusPartialContent, w.Code)
		packageClient.AssertExpectations(t)
	})

	t.Run("PackageNotFound", func(t *testing.T) {
		pkgId := uuid.NewV4()

		packageClient.On("GetPackage", pkgId.String()).Return(nil,
			crest.HttpError{
				HttpCode: http.StatusNotFound,
				Message:  "GetNetwork failure",
			})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/%s", pkgEndpoint, pkgId), nil)

		r := NewRouter(nil, packageClient, nil, nil, routerConfig,
			arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusNotFound, w.Code)
		packageClient.AssertExpectations(t)
	})

	t.Run("GetPackageError", func(t *testing.T) {
		pkgId := uuid.NewV4()

		packageClient.On("GetPackage", pkgId.String()).Return(nil,
			errors.New("some unexpected error"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/%s", pkgEndpoint, pkgId), nil)

		r := NewRouter(nil, packageClient, nil, nil, routerConfig,
			arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		packageClient.AssertExpectations(t)
	})
}

func TestRouter_AddPackage(t *testing.T) {
	arc := &providers.AuthRestClient{}

	pkgName := "Monthly Data"
	from := "2023-04-01T00:00:00Z"
	to := "2023-04-01T00:00:00Z"
	voiceVolume := uint64(0)
	isActive := true
	dataVolume := uint64(1024)
	smsVolume := uint64(0)
	dataUnit := "MegaBytes"
	voiceUnit := "seconds"
	simType := "test"
	apn := "ukama.tel"
	markupValue := float64(0)
	markup := cdplan.PackageMarkup{
		Markup: markupValue,
	}
	pType := "postpaid"
	duration := uint64(30)
	flatRate := false
	amount := float64(0)
	overdraft := float64(0)
	trafficPolicy := uint32(0)
	networks := []string{""}

	t.Run("PackageCreatedAndStatusUpdated", func(t *testing.T) {
		pkgId := uuid.NewV4()
		orgId := uuid.NewV4().String()
		ownerId := uuid.NewV4().String()
		baserateId := uuid.NewV4().String()

		var pkg = AddPackageReq{
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
			Markup:        markupValue,
			Type:          pType,
			Duration:      duration,
			Flatrate:      flatRate,
			Amount:        amount,
			Overdraft:     overdraft,
			TrafficPolicy: trafficPolicy,
			Networks:      networks,
		}

		pkgInfo := &cdplan.PackageInfo{
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
		}

		body, err := json.Marshal(pkg)
		if err != nil {
			t.Errorf("fail to marshal request data: %v. Error: %v", pkg, err)
		}

		packageClient.On("AddPackage", pkgName, orgId, ownerId, from, to, baserateId,
			isActive, flatRate, smsVolume, voiceVolume, dataVolume, voiceUnit, dataUnit,
			simType, apn, pType, duration, markupValue, amount, overdraft, trafficPolicy, networks).
			Return(pkgInfo, nil)

		r := NewRouter(nil, packageClient, nil, nil, routerConfig,
			arc.MockAuthenticateUser).f.Engine()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", pkgEndpoint, bytes.NewReader(body))

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusPartialContent, w.Code)
		packageClient.AssertExpectations(t)
	})

	t.Run("PackageCreatedAndStatusFailed", func(t *testing.T) {
		orgId := uuid.NewV4().String()
		ownerId := uuid.NewV4().String()
		baserateId := uuid.NewV4().String()

		var pkg = AddPackageReq{
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
			Markup:        markupValue,
			Type:          pType,
			Duration:      duration,
			Flatrate:      flatRate,
			Amount:        amount,
			Overdraft:     overdraft,
			TrafficPolicy: trafficPolicy,
			Networks:      networks,
		}

		body, err := json.Marshal(pkg)
		if err != nil {
			t.Errorf("fail to marshal request data: %v. Error: %v", pkg, err)
		}

		packageClient.On("AddPackage", pkgName, orgId, ownerId, from, to, baserateId,
			isActive, flatRate, smsVolume, voiceVolume, dataVolume, voiceUnit, dataUnit,
			simType, apn, pType, duration, markupValue, amount, overdraft, trafficPolicy, networks).
			Return(nil, errors.New("some unexpected error occured"))

		r := NewRouter(nil, packageClient, nil, nil, routerConfig,
			arc.MockAuthenticateUser).f.Engine()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", pkgEndpoint, bytes.NewReader(body))

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		packageClient.AssertExpectations(t)
	})
}

func TestRouter_GetSim(t *testing.T) {
	arc := &providers.AuthRestClient{}

	subscriberId := uuid.NewV4()

	t.Run("SimFoundAndStatusCompleted", func(t *testing.T) {
		simId := uuid.NewV4()

		simInfo := &csub.SimInfo{
			Id:           simId.String(),
			SubscriberId: subscriberId.String(),
		}

		simClient.On("GetSim", simId.String()).Return(simInfo, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/%s", simEndpoint, simId), nil)

		r := NewRouter(nil, nil, simClient, nil, routerConfig,
			arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusOK, w.Code)
		simClient.AssertExpectations(t)
	})

	t.Run("SimFoundAndStatusPending", func(t *testing.T) {
		simId := uuid.NewV4()

		simInfo := &csub.SimInfo{
			Id:           simId.String(),
			SubscriberId: subscriberId.String(),
		}

		simClient.On("GetSim", simId.String()).Return(simInfo,
			crest.HttpError{
				HttpCode: http.StatusPartialContent,
				Message:  "partial content. request is still ongoing",
			})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/%s", simEndpoint, simId), nil)

		r := NewRouter(nil, nil, simClient, nil, routerConfig,
			arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusPartialContent, w.Code)
		simClient.AssertExpectations(t)
	})

	t.Run("SimNotFound", func(t *testing.T) {
		simId := uuid.NewV4()

		simClient.On("GetSim", simId.String()).Return(nil,
			crest.HttpError{
				HttpCode: http.StatusNotFound,
				Message:  "GetSim failure",
			})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/%s", simEndpoint, simId), nil)

		r := NewRouter(nil, nil, simClient, nil, routerConfig,
			arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusNotFound, w.Code)
		simClient.AssertExpectations(t)
	})

	t.Run("GetSimError", func(t *testing.T) {
		simId := uuid.NewV4()

		simClient.On("GetSim", simId.String()).Return(nil,
			errors.New("some unexpected error"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/%s", simEndpoint, simId), nil)

		r := NewRouter(nil, nil, simClient, nil, routerConfig,
			arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		simClient.AssertExpectations(t)
	})
}

func TestRouter_ConfigureSim(t *testing.T) {
	arc := &providers.AuthRestClient{}
	networkId := uuid.NewV4()
	packageId := uuid.NewV4()
	simType := "some-sim-type"
	simToken := "some-sim-token"
	trafficPolicy := uint32(0)

	orgId := uuid.NewV4()
	name := "John Doe"
	email := "johndoe@example.com"
	phoneNumber := "0123456789"
	address := "2 Rivers"
	dob := "2023/09/01"
	proofOfID := "passport"
	idSerial := "987654321"

	var sim = AddSimReq{
		// SubscriberId:          subscriberId.String(),
		NetworkId:             networkId.String(),
		PackageId:             packageId.String(),
		OrgId:                 orgId.String(),
		Name:                  name,
		Email:                 email,
		PhoneNumber:           phoneNumber,
		Address:               address,
		Dob:                   dob,
		ProofOfIdentification: proofOfID,
		IdSerial:              idSerial,
		SimType:               simType,
		SimToken:              simToken,
		TrafficPolicy:         trafficPolicy,
	}

	t.Run("SimConfiguredAndStatusUpdated", func(t *testing.T) {
		simId := uuid.NewV4()
		subscriberId := uuid.NewV4()

		sim.SubscriberId = subscriberId.String()

		simInfo := &csub.SimInfo{
			Id:           simId.String(),
			SubscriberId: subscriberId.String(),
		}

		body, err := json.Marshal(sim)
		if err != nil {
			t.Errorf("fail to marshal request data: %v. Error: %v", sim, err)
		}

		simClient.On("ConfigureSim", subscriberId.String(), orgId.String(),
			networkId.String(), name, email, phoneNumber, address,
			dob, proofOfID, idSerial, packageId.String(), simType, simToken, trafficPolicy).
			Return(simInfo, nil)

		r := NewRouter(nil, nil, simClient, nil, routerConfig,
			arc.MockAuthenticateUser).f.Engine()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", simEndpoint, bytes.NewReader(body))

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusPartialContent, w.Code)
		simClient.AssertExpectations(t)
	})

	t.Run("SimconfiguredAndStatusFailed", func(t *testing.T) {
		subscriberId := uuid.NewV4()

		sim.SubscriberId = subscriberId.String()

		body, err := json.Marshal(sim)
		if err != nil {
			t.Errorf("fail to marshal request data: %v. Error: %v", sim, err)
		}

		simClient.On("ConfigureSim", subscriberId.String(), orgId.String(),
			networkId.String(), name, email, phoneNumber, address,
			dob, proofOfID, idSerial, packageId.String(), simType, simToken, trafficPolicy).
			Return(nil, errors.New("some unexpected error occured"))

		r := NewRouter(nil, nil, simClient, nil, routerConfig,
			arc.MockAuthenticateUser).f.Engine()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", simEndpoint, bytes.NewReader(body))

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		simClient.AssertExpectations(t)
	})
}

func TestRouter_GetNode(t *testing.T) {
	arc := &providers.AuthRestClient{}

	nodeName := "node-1"

	t.Run("NodeFound", func(t *testing.T) {
		nodeId := "uk-sa2341-hnode-v0-a1a0"

		nodeInfo := &creg.NodeInfo{
			Id:   nodeId,
			Name: nodeName,
		}

		nodeClient.On("GetNode", nodeId).Return(nodeInfo, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/%s", nodeEndpoint, nodeId), nil)

		r := NewRouter(nil, nil, nil, nodeClient, routerConfig,
			arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusOK, w.Code)
		nodeClient.AssertExpectations(t)
	})

	t.Run("NodeNotFound", func(t *testing.T) {
		nodeId := "uk-sa2341-tnode-v0-a1a0"

		nodeClient.On("GetNode", nodeId).Return(nil,
			crest.HttpError{
				HttpCode: http.StatusNotFound,
				Message:  "GetNode failure",
			})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/%s", nodeEndpoint, nodeId), nil)

		r := NewRouter(nil, nil, nil, nodeClient, routerConfig,
			arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusNotFound, w.Code)
		nodeClient.AssertExpectations(t)
	})

	t.Run("GetNodeError", func(t *testing.T) {
		nodeId := "uk-sa2341-anode-v0-a1a0"

		nodeClient.On("GetNode", nodeId).Return(nil,
			errors.New("some unexpected error"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/%s", nodeEndpoint, nodeId), nil)

		r := NewRouter(nil, nil, nil, nodeClient, routerConfig,
			arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		nodeClient.AssertExpectations(t)
	})
}

func TestRouter_RegisterNode(t *testing.T) {
	arc := &providers.AuthRestClient{}

	t.Run("NodeRegistered", func(t *testing.T) {
		nodeId := "uk-sa2341-hnode-v0-a1a0"
		nodeName := "node-1"
		orgId := uuid.NewV4()
		state := "pending"

		var node = AddNodeRequest{
			NodeId: nodeId,
			Name:   nodeName,
			OrgId:  orgId.String(),
			State:  state,
		}

		nodeInfo := &creg.NodeInfo{
			Id:   nodeId,
			Name: nodeName,
		}

		body, err := json.Marshal(node)
		if err != nil {
			t.Errorf("fail to marshal request data: %v. Error: %v", node, err)
		}

		nodeClient.On("RegisterNode", nodeId, nodeName, orgId.String(), state).
			Return(nodeInfo, nil)

		r := NewRouter(nil, nil, nil, nodeClient, routerConfig,
			arc.MockAuthenticateUser).f.Engine()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", nodeEndpoint, bytes.NewReader(body))

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusCreated, w.Code)
		nodeClient.AssertExpectations(t)
	})

	t.Run("NodeNotRegistered", func(t *testing.T) {
		nodeId := "uk-sa2341-hnode-v0-a1a0"
		nodeName := "node-1"
		orgId := uuid.NewV4()
		state := "pending"

		var node = AddNodeRequest{
			NodeId: nodeId,
			Name:   nodeName,
			OrgId:  orgId.String(),
			State:  state,
		}

		body, err := json.Marshal(node)
		if err != nil {
			t.Errorf("fail to marshal request data: %v. Error: %v", node, err)
		}

		nodeClient.On("RegisterNode", nodeId, nodeName, orgId.String(), state).
			Return(nil, errors.New("some unexpected error occured"))

		r := NewRouter(nil, nil, nil, nodeClient, routerConfig,
			arc.MockAuthenticateUser).f.Engine()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", nodeEndpoint, bytes.NewReader(body))

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		nodeClient.AssertExpectations(t)
	})
}

func TestRouter_AttachNode(t *testing.T) {
	arc := &providers.AuthRestClient{}

	ampNodeL := "uk-sa2341-anode-v0-a1a0"
	ampNodeR := "uk-sa2341-anode-v0-a1a1"

	var nodes = AttachNodesRequest{
		AmpNodeL: ampNodeL,
		AmpNodeR: ampNodeR,
	}

	t.Run("NodeAttached", func(t *testing.T) {
		nodeId := "uk-sa2341-tnode-v0-a1a0"

		body, err := json.Marshal(nodes)
		if err != nil {
			t.Errorf("fail to marshal request data: %v. Error: %v", nodes, err)
		}

		nodeClient.On("AttachNode", nodeId, ampNodeL, ampNodeR).
			Return(nil)

		r := NewRouter(nil, nil, nil, nodeClient, routerConfig,
			arc.MockAuthenticateUser).f.Engine()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", nodeEndpoint+"/"+nodeId+"/attach", bytes.NewReader(body))

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusCreated, w.Code)
		nodeClient.AssertExpectations(t)
	})

	t.Run("NodeNotAttached", func(t *testing.T) {
		nodeId := "uk-sa2341-tnode-v0-a1a1"

		body, err := json.Marshal(nodes)
		if err != nil {
			t.Errorf("fail to marshal request data: %v. Error: %v", nodes, err)
		}

		nodeClient.On("AttachNode", nodeId, ampNodeL, ampNodeR).
			Return(errors.New("some unexpected error occured"))

		r := NewRouter(nil, nil, nil, nodeClient, routerConfig,
			arc.MockAuthenticateUser).f.Engine()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", nodeEndpoint+"/"+nodeId+"/attach", bytes.NewReader(body))

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		nodeClient.AssertExpectations(t)
	})
}

func TestRouter_DetachNode(t *testing.T) {
	arc := &providers.AuthRestClient{}

	t.Run("NodeDetached", func(t *testing.T) {
		nodeId := "uk-sa2341-hnode-v0-a1a0"

		nodeClient.On("DetachNode", nodeId).Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/%s/attach", nodeEndpoint, nodeId), nil)

		r := NewRouter(nil, nil, nil, nodeClient, routerConfig,
			arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusOK, w.Code)
		nodeClient.AssertExpectations(t)
	})

	t.Run("NodeNotFound", func(t *testing.T) {
		nodeId := "uk-sa2341-tnode-v0-a1a0"

		nodeClient.On("DetachNode", nodeId).Return(
			crest.HttpError{
				HttpCode: http.StatusNotFound,
				Message:  "DeleteNode failure",
			})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/%s/attach", nodeEndpoint, nodeId), nil)

		r := NewRouter(nil, nil, nil, nodeClient, routerConfig,
			arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusNotFound, w.Code)
		nodeClient.AssertExpectations(t)
	})

	t.Run("DetachNodeError", func(t *testing.T) {
		nodeId := "uk-sa2341-anode-v0-a1a0"

		nodeClient.On("DetachNode", nodeId).Return(
			errors.New("some unexpected error"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/%s/attach", nodeEndpoint, nodeId), nil)

		r := NewRouter(nil, nil, nil, nodeClient, routerConfig,
			arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		nodeClient.AssertExpectations(t)
	})
}

func TestRouter_AddNodeToSite(t *testing.T) {
	arc := &providers.AuthRestClient{}

	networkId := uuid.NewV4().String()
	siteId := uuid.NewV4().String()

	var req = AddNodeToSiteRequest{
		NetworkId: networkId,
		SiteId:    siteId,
	}

	t.Run("NodeAdded", func(t *testing.T) {
		nodeId := "uk-sa2341-tnode-v0-a1a0"

		body, err := json.Marshal(req)
		if err != nil {
			t.Errorf("fail to marshal request data: %v. Error: %v", req, err)
		}

		nodeClient.On("AddNodeToSite", nodeId, networkId, siteId).
			Return(nil)

		r := NewRouter(nil, nil, nil, nodeClient, routerConfig,
			arc.MockAuthenticateUser).f.Engine()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", nodeEndpoint+"/"+nodeId+"/sites", bytes.NewReader(body))

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusCreated, w.Code)
		nodeClient.AssertExpectations(t)
	})

	t.Run("NodeNotAdded", func(t *testing.T) {
		nodeId := "uk-sa2341-tnode-v0-a1a1"

		body, err := json.Marshal(req)
		if err != nil {
			t.Errorf("fail to marshal request data: %v. Error: %v", req, err)
		}

		nodeClient.On("AddNodeToSite", nodeId, networkId, siteId).
			Return(errors.New("some unexpected error occured"))

		r := NewRouter(nil, nil, nil, nodeClient, routerConfig,
			arc.MockAuthenticateUser).f.Engine()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", nodeEndpoint+"/"+nodeId+"/sites", bytes.NewReader(body))

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		nodeClient.AssertExpectations(t)
	})
}

func TestRouter_RemoveNodeFromSite(t *testing.T) {
	arc := &providers.AuthRestClient{}

	t.Run("NodeRemoved", func(t *testing.T) {
		nodeId := "uk-sa2341-hnode-v0-a1a0"

		nodeClient.On("RemoveNodeFromSite", nodeId).Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/%s/sites", nodeEndpoint, nodeId), nil)

		r := NewRouter(nil, nil, nil, nodeClient, routerConfig,
			arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusOK, w.Code)
		nodeClient.AssertExpectations(t)
	})

	t.Run("NodeNotFound", func(t *testing.T) {
		nodeId := "uk-sa2341-tnode-v0-a1a0"

		nodeClient.On("RemoveNodeFromSite", nodeId).Return(
			crest.HttpError{
				HttpCode: http.StatusNotFound,
				Message:  "DeleteNode failure",
			})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/%s/sites", nodeEndpoint, nodeId), nil)

		r := NewRouter(nil, nil, nil, nodeClient, routerConfig,
			arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusNotFound, w.Code)
		nodeClient.AssertExpectations(t)
	})

	t.Run("RemoveNodeError", func(t *testing.T) {
		nodeId := "uk-sa2341-anode-v0-a1a0"

		nodeClient.On("RemoveNodeFromSite", nodeId).Return(
			errors.New("some unexpected error"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/%s/sites", nodeEndpoint, nodeId), nil)

		r := NewRouter(nil, nil, nil, nodeClient, routerConfig,
			arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		nodeClient.AssertExpectations(t)
	})
}

func TestRouter_DeleteNode(t *testing.T) {
	arc := &providers.AuthRestClient{}

	t.Run("NodeDeleted", func(t *testing.T) {
		nodeId := "uk-sa2341-hnode-v0-a1a0"

		nodeClient.On("DeleteNode", nodeId).Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/%s", nodeEndpoint, nodeId), nil)

		r := NewRouter(nil, nil, nil, nodeClient, routerConfig,
			arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusOK, w.Code)
		nodeClient.AssertExpectations(t)
	})

	t.Run("NodeNotFound", func(t *testing.T) {
		nodeId := "uk-sa2341-tnode-v0-a1a0"

		nodeClient.On("DeleteNode", nodeId).Return(
			crest.HttpError{
				HttpCode: http.StatusNotFound,
				Message:  "DeleteNode failure",
			})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/%s", nodeEndpoint, nodeId), nil)

		r := NewRouter(nil, nil, nil, nodeClient, routerConfig,
			arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusNotFound, w.Code)
		nodeClient.AssertExpectations(t)
	})

	t.Run("NodeDeleteError", func(t *testing.T) {
		nodeId := "uk-sa2341-anode-v0-a1a0"

		nodeClient.On("DeleteNode", nodeId).Return(
			errors.New("some unexpected error"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/%s", nodeEndpoint, nodeId), nil)

		r := NewRouter(nil, nil, nil, nodeClient, routerConfig,
			arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		nodeClient.AssertExpectations(t)
	})
}
