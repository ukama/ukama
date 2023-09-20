package rest

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ukama/ukama/systems/common/providers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	cconfig "github.com/ukama/ukama/systems/common/config"
	crest "github.com/ukama/ukama/systems/common/rest"
	"github.com/ukama/ukama/systems/node/api-gateway/pkg"
	"github.com/ukama/ukama/systems/node/api-gateway/pkg/client"
	cpb "github.com/ukama/ukama/systems/node/controller/pb/gen"
	nmocks "github.com/ukama/ukama/systems/node/controller/pb/gen/mocks"
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

var testClientSet *Clients

func init() {
	gin.SetMode(gin.TestMode)
	testClientSet = NewClientsSet(&pkg.GrpcEndpoints{
		Timeout:    1 * time.Second,
		Controller: "0.0.0.0:9092",
	})
}
func TestPingRoute(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	arc := &providers.AuthRestClient{}
	req, _ := http.NewRequest("GET", "/ping", nil)
	r := NewRouter(testClientSet, routerConfig, arc.MockAuthenticateUser).f.Engine()
	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "pong")
}

func Test_RestarteNode(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/controllers", nil)
	arc := &providers.AuthRestClient{}
	c := &nmocks.ControllerServiceClient{}

	c.On("RestartNode", mock.Anything, mock.Anything).Return(&cpb.RestartNodeResponse{
		Status: cpb.RestartStatus_RESTART_STATUS_SUCCESS},
		nil)

	r := NewRouter(&Clients{
		Controller: client.NewControllerFromClient(c),
	}, routerConfig, arc.MockAuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	c.AssertExpectations(t)
}

// func TestRouter_PingRoute(t *testing.T) {
// 	var c = &nmocks.ControllerServiceClient{}
// 	var arc = &providers.AuthRestClient{}

// 	// arrange
// 	w := httptest.NewRecorder()
// 	req, _ := http.NewRequest("GET", "/ping", nil)

// 	r := NewRouter(&Clients{
// 		Controller: client.NewControllerFromClient(c),
// 	}, routerConfig, arc.MockAuthenticateUser).f.Engine()

// 	r.ServeHTTP(w, req)

// 	assert.Equal(t, 200, w.Code)
// 	assert.Contains(t, w.Body.String(), "pong")
// }

// func TestRouter_RestartNode(t *testing.T) {
// 	var c = &nmocks.ControllerServiceClient{}
// 	var arc = &providers.AuthRestClient{}
// 	var cr = RestartNodeRequest{
// 		NodeId: uuid.NewV4().String(),
// 	}

// 	t.Run("NodeHasRestared", func(t *testing.T) {
// 		body, err := json.Marshal(cr)
// 		if err != nil {
// 			t.Errorf("fail to marshal request data: %v. Error: %v", cr, err)
// 		}

// 		w := httptest.NewRecorder()
// 		req, _ := http.NewRequest("POST", nodeApiEndpoint, bytes.NewReader(body))

// 		controllerReq := &cpb.RestartNodeRequest{
// 			NodeId: cr.NodeId,
// 		}

// 		c.On("RestartNode", controllerReq).Return(&cpb.RestartNodeResponse{Status: cpb.RestartStatus_RESTART_STATUS_SUCCESS}, nil)

// 		r := NewRouter(&Clients{
// 			Controller: client.NewControllerFromClient(c),
// 		}, routerConfig, arc.MockAuthenticateUser).f.Engine()

// 		// act
// 		r.ServeHTTP(w, req)

// 		// assert
// 		assert.Equal(t, http.StatusCreated, w.Code)
// 		c.AssertExpectations(t)
// 	})

// }
