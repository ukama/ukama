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
	crest "github.com/ukama/ukama/systems/common/rest"
)

const netEndpoint = "/v1/networks"
const simEndpoint = "/v1/sims"

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

var testClientSet client.Client

func init() {
	netClient := &mocks.NetworkClient{}
	subscriberClient := &mocks.SubscriberClient{}
	simClient := &mocks.SimClient{}

	gin.SetMode(gin.TestMode)
	testClientSet = client.NewClientsSet(netClient, subscriberClient, simClient)
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
			Id:   netId,
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
			Id:   netId,
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
		paymentLinks := false

		var ntwk = AddNetworkReq{
			OrgName:          orgName,
			NetName:          netName,
			AllowedCountries: countries,
			AllowedNetworks:  networks,
			PaymentLinks:     paymentLinks,
		}

		netInfo := &client.NetworkInfo{
			Id:   netId,
			Name: netName,
		}

		body, err := json.Marshal(ntwk)
		if err != nil {
			t.Errorf("fail to marshal request data: %v. Error: %v", ntwk, err)
		}

		c.On("CreateNetwork", orgName, netName, countries, networks, paymentLinks).Return(netInfo, nil)

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
		paymentLinks := false

		var ntwk = AddNetworkReq{
			OrgName:          orgName,
			NetName:          netName,
			AllowedCountries: countries,
			AllowedNetworks:  networks,
			PaymentLinks:     paymentLinks,
		}

		body, err := json.Marshal(ntwk)
		if err != nil {
			t.Errorf("fail to marshal request data: %v. Error: %v", ntwk, err)
		}

		c.On("CreateNetwork", orgName, netName, countries, networks, paymentLinks).Return(nil,
			errors.New("some unexpected error occured"))

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

func TestRouter_GetSim(t *testing.T) {
	c := &mocks.Client{}
	arc := &providers.AuthRestClient{}

	subscriberId := uuid.NewV4()

	t.Run("SimFoundAndStatusCompleted", func(t *testing.T) {
		simId := uuid.NewV4()

		simInfo := &client.SimInfo{
			Id:           simId,
			SubscriberId: subscriberId,
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
			Id:           simId,
			SubscriberId: subscriberId,
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

	t.Run("SimConfiguredAndStatusUpdated", func(t *testing.T) {
		simId := uuid.NewV4()
		subscriberId := uuid.NewV4()
		networkId := uuid.NewV4()
		packageId := uuid.NewV4()
		simType := "some-sim-type"
		simToken := "some-sim-token"

		var sim = AddSimReq{
			SubscriberId: subscriberId.String(),
			NetworkId:    networkId.String(),
			PackageId:    packageId.String(),
			SimType:      simType,
			SimToken:     simToken,
		}

		simInfo := &client.SimInfo{
			Id:           simId,
			SubscriberId: subscriberId,
		}

		body, err := json.Marshal(sim)
		if err != nil {
			t.Errorf("fail to marshal request data: %v. Error: %v", sim, err)
		}

		c.On("ConfigureSim", subscriberId.String(), networkId.String(),
			packageId.String(), simType, simToken).
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
		networkId := uuid.NewV4()
		packageId := uuid.NewV4()
		simType := "some-sim-type"
		simToken := "some-sim-token"

		var sim = AddSimReq{
			SubscriberId: subscriberId.String(),
			NetworkId:    networkId.String(),
			PackageId:    packageId.String(),
			SimType:      simType,
			SimToken:     simToken,
		}

		body, err := json.Marshal(sim)
		if err != nil {
			t.Errorf("fail to marshal request data: %v. Error: %v", sim, err)
		}

		c.On("ConfigureSim", subscriberId.String(), networkId.String(),
			packageId.String(), simType, simToken).
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
