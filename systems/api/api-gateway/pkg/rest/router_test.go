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
	"github.com/ukama/ukama/systems/api/api-gateway/pkg/client"
	"github.com/ukama/ukama/systems/common/providers"
	"github.com/ukama/ukama/systems/common/rest"
	"github.com/ukama/ukama/systems/common/uuid"

	cconfig "github.com/ukama/ukama/systems/common/config"
)

const netEndpoint = "/v1/networks"
const pkgEndpoint = "/v1/packages"
const simEndpoint = "/v1/sims"
const nodeEndpoint = "/v1/nodes"

var defaultCors = cors.Config{
	AllowAllOrigins: true,
}

var routerConfig = &RouterConfig{
	serverConf: &rest.HttpConfig{
		Cors: defaultCors,
	},
	auth: &cconfig.Auth{
		AuthAppUrl:    "http://localhost:4455",
		AuthServerUrl: "http://localhost:4434",
		AuthAPIGW:     "http://localhost:8080",
	},
}

var testClientSet client.Client

func init() {
	netClient := &mocks.NetworkClient{}
	packageClient := &mocks.PackageClient{}
	subscriberClient := &mocks.SubscriberClient{}
	simClient := &mocks.SimClient{}
	nodeClient := &mocks.NodeClient{}

	gin.SetMode(gin.TestMode)
	testClientSet = client.NewClientsSet(netClient, packageClient, subscriberClient, simClient, nodeClient)
}

func TestRouter_PingRoute(t *testing.T) {
	var c = &mocks.Client{}
	var arc = &providers.AuthRestClient{}

	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)

	r := NewRouter(c, routerConfig, arc.MockAuthenticateUser).f.Engine()

	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "pong")
}

func TestRouter_GetNetwork(t *testing.T) {
	c := &mocks.Client{}
	arc := &providers.AuthRestClient{}

	netName := "net-1"

	t.Run("NetworkFoundAndStatusCompleted", func(t *testing.T) {
		netId := uuid.NewV4()

		netInfo := &client.NetworkInfo{
			Id:   netId.String(),
			Name: netName,
		}

		c.On("GetNetwork", netId.String()).Return(netInfo, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/%s", netEndpoint, netId), nil)

		r := NewRouter(c, routerConfig, arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusOK, w.Code)
		c.AssertExpectations(t)
	})

	t.Run("NetworkFoundAndStatusPending", func(t *testing.T) {
		netId := uuid.NewV4()

		netInfo := &client.NetworkInfo{
			Id:   netId.String(),
			Name: netName,
		}

		c.On("GetNetwork", netId.String()).Return(netInfo,
			rest.HttpError{
				HttpCode: http.StatusPartialContent,
				Message:  "partial content. request is still ongoing",
			})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/%s", netEndpoint, netId), nil)

		r := NewRouter(c, routerConfig, arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusPartialContent, w.Code)
		c.AssertExpectations(t)
	})

	t.Run("NetworkNotFound", func(t *testing.T) {
		netId := uuid.NewV4()

		c.On("GetNetwork", netId.String()).Return(nil,
			rest.HttpError{
				HttpCode: http.StatusNotFound,
				Message:  "GetNetwork failure",
			})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/%s", netEndpoint, netId), nil)

		r := NewRouter(c, routerConfig, arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusNotFound, w.Code)
		c.AssertExpectations(t)
	})

	t.Run("GetNetworkGetError", func(t *testing.T) {
		netId := uuid.NewV4()

		c.On("GetNetwork", netId.String()).Return(nil,
			errors.New("some unexpected error"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/%s", netEndpoint, netId), nil)

		r := NewRouter(c, routerConfig, arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		c.AssertExpectations(t)
	})
}

func TestRouter_CreateNetwork(t *testing.T) {
	c := &mocks.Client{}
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

		netInfo := &client.NetworkInfo{
			Id:   netId.String(),
			Name: netName,
		}

		body, err := json.Marshal(ntwk)
		if err != nil {
			t.Errorf("fail to marshal request data: %v. Error: %v", ntwk, err)
		}

		c.On("CreateNetwork", orgName, netName, countries, networks, budget, overdraft, trafficPolicy, paymentLinks).
			Return(netInfo, nil)

		r := NewRouter(c, routerConfig, arc.MockAuthenticateUser).f.Engine()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", netEndpoint, bytes.NewReader(body))

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusPartialContent, w.Code)
		c.AssertExpectations(t)
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

		c.On("CreateNetwork", orgName, netName, countries, networks, budget, overdraft, trafficPolicy, paymentLinks).
			Return(nil, errors.New("some unexpected error occured"))

		r := NewRouter(c, routerConfig, arc.MockAuthenticateUser).f.Engine()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", netEndpoint, bytes.NewReader(body))

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		c.AssertExpectations(t)
	})
}

func TestRouter_GetPackage(t *testing.T) {
	c := &mocks.Client{}
	arc := &providers.AuthRestClient{}

	pkgName := "Monthly Data"

	t.Run("PackageFoundAndStatusCompleted", func(t *testing.T) {
		pkgId := uuid.NewV4()

		pkgInfo := &client.PackageInfo{
			Id:   pkgId.String(),
			Name: pkgName,
		}

		c.On("GetPackage", pkgId.String()).Return(pkgInfo, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/%s", pkgEndpoint, pkgId), nil)

		r := NewRouter(c, routerConfig, arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusOK, w.Code)
		c.AssertExpectations(t)
	})

	t.Run("PackageFoundAndStatusPending", func(t *testing.T) {
		pkgId := uuid.NewV4()

		pkgInfo := &client.PackageInfo{
			Id:   pkgId.String(),
			Name: pkgName,
		}

		c.On("GetPackage", pkgId.String()).Return(pkgInfo,
			rest.HttpError{
				HttpCode: http.StatusPartialContent,
				Message:  "partial content. request is still ongoing",
			})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/%s", pkgEndpoint, pkgId), nil)

		r := NewRouter(c, routerConfig, arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusPartialContent, w.Code)
		c.AssertExpectations(t)
	})

	t.Run("PackageNotFound", func(t *testing.T) {
		pkgId := uuid.NewV4()

		c.On("GetPackage", pkgId.String()).Return(nil,
			rest.HttpError{
				HttpCode: http.StatusNotFound,
				Message:  "GetNetwork failure",
			})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/%s", pkgEndpoint, pkgId), nil)

		r := NewRouter(c, routerConfig, arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusNotFound, w.Code)
		c.AssertExpectations(t)
	})

	t.Run("GetPackageError", func(t *testing.T) {
		pkgId := uuid.NewV4()

		c.On("GetPackage", pkgId.String()).Return(nil,
			errors.New("some unexpected error"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/%s", pkgEndpoint, pkgId), nil)

		r := NewRouter(c, routerConfig, arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		c.AssertExpectations(t)
	})
}

func TestRouter_AddPackage(t *testing.T) {
	c := &mocks.Client{}
	arc := &providers.AuthRestClient{}

	pkgName := "Monthly Data"
	from := "2023-04-01T00:00:00Z"
	to := "2023-04-01T00:00:00Z"
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
			Markup:        markup,
			Type:          pType,
			Duration:      duration,
			Flatrate:      flatRate,
			Amount:        amount,
			Overdraft:     overdraft,
			TrafficPolicy: trafficPolicy,
			Networks:      networks,
		}

		pkgInfo := &client.PackageInfo{
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

		c.On("AddPackage", pkgName, orgId, ownerId, from, to, baserateId,
			isActive, flatRate, smsVolume, voiceVolume, dataVolume, voiceUnit, dataUnit,
			simType, apn, pType, duration, markup, amount, overdraft, trafficPolicy, networks).
			Return(pkgInfo, nil)

		r := NewRouter(c, routerConfig, arc.MockAuthenticateUser).f.Engine()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", pkgEndpoint, bytes.NewReader(body))

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusPartialContent, w.Code)
		c.AssertExpectations(t)
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

		c.On("AddPackage", pkgName, orgId, ownerId, from, to, baserateId,
			isActive, flatRate, smsVolume, voiceVolume, dataVolume, voiceUnit, dataUnit,
			simType, apn, pType, duration, markup, amount, overdraft, trafficPolicy, networks).
			Return(nil, errors.New("some unexpected error occured"))

		r := NewRouter(c, routerConfig, arc.MockAuthenticateUser).f.Engine()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", pkgEndpoint, bytes.NewReader(body))

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		c.AssertExpectations(t)
	})
}

func TestRouter_GetSim(t *testing.T) {
	c := &mocks.Client{}
	arc := &providers.AuthRestClient{}

	subscriberId := uuid.NewV4()

	t.Run("SimFoundAndStatusCompleted", func(t *testing.T) {
		simId := uuid.NewV4()

		simInfo := &client.SimInfo{
			Id:           simId.String(),
			SubscriberId: subscriberId.String(),
		}

		c.On("GetSim", simId.String()).Return(simInfo, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/%s", simEndpoint, simId), nil)

		r := NewRouter(c, routerConfig, arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusOK, w.Code)
		c.AssertExpectations(t)
	})

	t.Run("SimFoundAndStatusPending", func(t *testing.T) {
		simId := uuid.NewV4()

		simInfo := &client.SimInfo{
			Id:           simId.String(),
			SubscriberId: subscriberId.String(),
		}

		c.On("GetSim", simId.String()).Return(simInfo,
			rest.HttpError{
				HttpCode: http.StatusPartialContent,
				Message:  "partial content. request is still ongoing",
			})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/%s", simEndpoint, simId), nil)

		r := NewRouter(c, routerConfig, arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusPartialContent, w.Code)
		c.AssertExpectations(t)
	})

	t.Run("SimNotFound", func(t *testing.T) {
		simId := uuid.NewV4()

		c.On("GetSim", simId.String()).Return(nil,
			rest.HttpError{
				HttpCode: http.StatusNotFound,
				Message:  "GetSim failure",
			})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/%s", simEndpoint, simId), nil)

		r := NewRouter(c, routerConfig, arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusNotFound, w.Code)
		c.AssertExpectations(t)
	})

	t.Run("GetSimError", func(t *testing.T) {
		simId := uuid.NewV4()

		c.On("GetSim", simId.String()).Return(nil,
			errors.New("some unexpected error"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/%s", simEndpoint, simId), nil)

		r := NewRouter(c, routerConfig, arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		c.AssertExpectations(t)
	})
}

func TestRouter_ConfigureSim(t *testing.T) {
	c := &mocks.Client{}
	arc := &providers.AuthRestClient{}
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

	var sim = AddSimReq{
		// SubscriberId:          subscriberId.String(),
		NetworkId:             networkId.String(),
		PackageId:             packageId.String(),
		OrgId:                 orgId.String(),
		FirstName:             firstName,
		LastName:              lastName,
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

		simInfo := &client.SimInfo{
			Id:           simId.String(),
			SubscriberId: subscriberId.String(),
		}

		body, err := json.Marshal(sim)
		if err != nil {
			t.Errorf("fail to marshal request data: %v. Error: %v", sim, err)
		}

		c.On("ConfigureSim", subscriberId.String(), orgId.String(),
			networkId.String(), firstName, lastName, email, phoneNumber, address,
			dob, proofOfID, idSerial, packageId.String(), simType, simToken, trafficPolicy).
			Return(simInfo, nil)

		r := NewRouter(c, routerConfig, arc.MockAuthenticateUser).f.Engine()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", simEndpoint, bytes.NewReader(body))

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusPartialContent, w.Code)
		c.AssertExpectations(t)
	})

	t.Run("SimconfiguredAndStatusFailed", func(t *testing.T) {
		subscriberId := uuid.NewV4()

		sim.SubscriberId = subscriberId.String()

		body, err := json.Marshal(sim)
		if err != nil {
			t.Errorf("fail to marshal request data: %v. Error: %v", sim, err)
		}

		c.On("ConfigureSim", subscriberId.String(), orgId.String(),
			networkId.String(), firstName, lastName, email, phoneNumber, address,
			dob, proofOfID, idSerial, packageId.String(), simType, simToken, trafficPolicy).
			Return(nil, errors.New("some unexpected error occured"))

		r := NewRouter(c, routerConfig, arc.MockAuthenticateUser).f.Engine()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", simEndpoint, bytes.NewReader(body))

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		c.AssertExpectations(t)
	})
}

func TestRouter_GetNode(t *testing.T) {
	c := &mocks.Client{}
	arc := &providers.AuthRestClient{}

	nodeName := "node-1"

	t.Run("NodeFound", func(t *testing.T) {
		nodeId := "uk-sa2341-hnode-v0-a1a0"

		nodeInfo := &client.NodeInfo{
			Id:   nodeId,
			Name: nodeName,
		}

		c.On("GetNode", nodeId).Return(nodeInfo, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/%s", nodeEndpoint, nodeId), nil)

		r := NewRouter(c, routerConfig, arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusOK, w.Code)
		c.AssertExpectations(t)
	})

	t.Run("NodeNotFound", func(t *testing.T) {
		nodeId := "uk-sa2341-tnode-v0-a1a0"

		c.On("GetNode", nodeId).Return(nil,
			rest.HttpError{
				HttpCode: http.StatusNotFound,
				Message:  "GetNode failure",
			})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/%s", nodeEndpoint, nodeId), nil)

		r := NewRouter(c, routerConfig, arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusNotFound, w.Code)
		c.AssertExpectations(t)
	})

	t.Run("GetNodeError", func(t *testing.T) {
		nodeId := "uk-sa2341-anode-v0-a1a0"

		c.On("GetNode", nodeId).Return(nil,
			errors.New("some unexpected error"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/%s", nodeEndpoint, nodeId), nil)

		r := NewRouter(c, routerConfig, arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		c.AssertExpectations(t)
	})
}

func TestRouter_RegisterNode(t *testing.T) {
	c := &mocks.Client{}
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

		nodeInfo := &client.NodeInfo{
			Id:    nodeId,
			Name:  nodeName,
			OrgId: orgId.String(),
			State: state,
		}

		body, err := json.Marshal(node)
		if err != nil {
			t.Errorf("fail to marshal request data: %v. Error: %v", node, err)
		}

		c.On("RegisterNode", nodeId, nodeName, orgId.String(), state).
			Return(nodeInfo, nil)

		r := NewRouter(c, routerConfig, arc.MockAuthenticateUser).f.Engine()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", nodeEndpoint, bytes.NewReader(body))

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusCreated, w.Code)
		c.AssertExpectations(t)
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

		c.On("RegisterNode", nodeId, nodeName, orgId.String(), state).
			Return(nil, errors.New("some unexpected error occured"))

		r := NewRouter(c, routerConfig, arc.MockAuthenticateUser).f.Engine()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", nodeEndpoint, bytes.NewReader(body))

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		c.AssertExpectations(t)
	})
}

func TestRouter_AttachNode(t *testing.T) {
	c := &mocks.Client{}
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

		c.On("AttachNode", nodeId, ampNodeL, ampNodeR).
			Return(nil)

		r := NewRouter(c, routerConfig, arc.MockAuthenticateUser).f.Engine()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", nodeEndpoint+"/"+nodeId+"/attach", bytes.NewReader(body))

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusCreated, w.Code)
		c.AssertExpectations(t)
	})

	t.Run("NodeNotAttached", func(t *testing.T) {
		nodeId := "uk-sa2341-tnode-v0-a1a1"

		body, err := json.Marshal(nodes)
		if err != nil {
			t.Errorf("fail to marshal request data: %v. Error: %v", nodes, err)
		}

		c.On("AttachNode", nodeId, ampNodeL, ampNodeR).
			Return(errors.New("some unexpected error occured"))

		r := NewRouter(c, routerConfig, arc.MockAuthenticateUser).f.Engine()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", nodeEndpoint+"/"+nodeId+"/attach", bytes.NewReader(body))

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		c.AssertExpectations(t)
	})
}

func TestRouter_DetachNode(t *testing.T) {
	c := &mocks.Client{}
	arc := &providers.AuthRestClient{}

	t.Run("NodeDetached", func(t *testing.T) {
		nodeId := "uk-sa2341-hnode-v0-a1a0"

		c.On("DetachNode", nodeId).Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/%s/attach", nodeEndpoint, nodeId), nil)

		r := NewRouter(c, routerConfig, arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusOK, w.Code)
		c.AssertExpectations(t)
	})

	t.Run("NodeNotFound", func(t *testing.T) {
		nodeId := "uk-sa2341-tnode-v0-a1a0"

		c.On("DetachNode", nodeId).Return(
			rest.HttpError{
				HttpCode: http.StatusNotFound,
				Message:  "DeleteNode failure",
			})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/%s/attach", nodeEndpoint, nodeId), nil)

		r := NewRouter(c, routerConfig, arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusNotFound, w.Code)
		c.AssertExpectations(t)
	})

	t.Run("DetachNodeError", func(t *testing.T) {
		nodeId := "uk-sa2341-anode-v0-a1a0"

		c.On("DetachNode", nodeId).Return(
			errors.New("some unexpected error"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/%s/attach", nodeEndpoint, nodeId), nil)

		r := NewRouter(c, routerConfig, arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		c.AssertExpectations(t)
	})
}

func TestRouter_AddNodeToSite(t *testing.T) {
	c := &mocks.Client{}
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

		c.On("AddNodeToSite", nodeId, networkId, siteId).
			Return(nil)

		r := NewRouter(c, routerConfig, arc.MockAuthenticateUser).f.Engine()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", nodeEndpoint+"/"+nodeId+"/sites", bytes.NewReader(body))

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusCreated, w.Code)
		c.AssertExpectations(t)
	})

	t.Run("NodeNotAdded", func(t *testing.T) {
		nodeId := "uk-sa2341-tnode-v0-a1a1"

		body, err := json.Marshal(req)
		if err != nil {
			t.Errorf("fail to marshal request data: %v. Error: %v", req, err)
		}

		c.On("AddNodeToSite", nodeId, networkId, siteId).
			Return(errors.New("some unexpected error occured"))

		r := NewRouter(c, routerConfig, arc.MockAuthenticateUser).f.Engine()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", nodeEndpoint+"/"+nodeId+"/sites", bytes.NewReader(body))

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		c.AssertExpectations(t)
	})
}

func TestRouter_RemoveNodeFromSite(t *testing.T) {
	c := &mocks.Client{}
	arc := &providers.AuthRestClient{}

	t.Run("NodeRemoved", func(t *testing.T) {
		nodeId := "uk-sa2341-hnode-v0-a1a0"

		c.On("RemoveNodeFromSite", nodeId).Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/%s/sites", nodeEndpoint, nodeId), nil)

		r := NewRouter(c, routerConfig, arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusOK, w.Code)
		c.AssertExpectations(t)
	})

	t.Run("NodeNotFound", func(t *testing.T) {
		nodeId := "uk-sa2341-tnode-v0-a1a0"

		c.On("RemoveNodeFromSite", nodeId).Return(
			rest.HttpError{
				HttpCode: http.StatusNotFound,
				Message:  "DeleteNode failure",
			})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/%s/sites", nodeEndpoint, nodeId), nil)

		r := NewRouter(c, routerConfig, arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusNotFound, w.Code)
		c.AssertExpectations(t)
	})

	t.Run("RemoveNodeError", func(t *testing.T) {
		nodeId := "uk-sa2341-anode-v0-a1a0"

		c.On("RemoveNodeFromSite", nodeId).Return(
			errors.New("some unexpected error"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/%s/sites", nodeEndpoint, nodeId), nil)

		r := NewRouter(c, routerConfig, arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		c.AssertExpectations(t)
	})
}

func TestRouter_DeleteNode(t *testing.T) {
	c := &mocks.Client{}
	arc := &providers.AuthRestClient{}

	t.Run("NodeDeleted", func(t *testing.T) {
		nodeId := "uk-sa2341-hnode-v0-a1a0"

		c.On("DeleteNode", nodeId).Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/%s", nodeEndpoint, nodeId), nil)

		r := NewRouter(c, routerConfig, arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusOK, w.Code)
		c.AssertExpectations(t)
	})

	t.Run("NodeNotFound", func(t *testing.T) {
		nodeId := "uk-sa2341-tnode-v0-a1a0"

		c.On("DeleteNode", nodeId).Return(
			rest.HttpError{
				HttpCode: http.StatusNotFound,
				Message:  "DeleteNode failure",
			})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/%s", nodeEndpoint, nodeId), nil)

		r := NewRouter(c, routerConfig, arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusNotFound, w.Code)
		c.AssertExpectations(t)
	})

	t.Run("NodeDeleteError", func(t *testing.T) {
		nodeId := "uk-sa2341-anode-v0-a1a0"

		c.On("DeleteNode", nodeId).Return(
			errors.New("some unexpected error"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/%s", nodeEndpoint, nodeId), nil)

		r := NewRouter(c, routerConfig, arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		c.AssertExpectations(t)
	})
}
