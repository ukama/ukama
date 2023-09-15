package rest

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ukama/ukama/systems/common/providers"
	"github.com/ukama/ukama/systems/common/rest"
	"github.com/ukama/ukama/systems/common/uuid"

	"github.com/stretchr/testify/assert"
	"github.com/ukama/ukama/systems/api/api-gateway/mocks"
	"github.com/ukama/ukama/systems/api/api-gateway/pkg/client"

	cconfig "github.com/ukama/ukama/systems/common/config"
	crest "github.com/ukama/ukama/systems/common/rest"
)

const netEndpoint = "/v1/networks"

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
	resRepo := &mocks.ResourceRepo{}
	netClient := &mocks.NetworkClient{}

	gin.SetMode(gin.TestMode)
	testClientSet = client.NewClientsSet(resRepo, netClient)
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

	t.Run("NetworkFoundAndStatusFailed", func(t *testing.T) {
		netId := uuid.NewV4()

		c.On("GetNetwork", netId.String()).Return(nil,
			rest.HttpError{
				HttpCode: http.StatusUnprocessableEntity,
				Message:  "inconsistent state. request has failed",
			})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/%s", netEndpoint, netId), nil)

		r := NewRouter(c, routerConfig, arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
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
