package rest

import (
	"net/http"
	"net/http/httptest"
	"strings"
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
	req, _ := http.NewRequest("POST", "/v1/controllers/nodes/60285a2a-fe1d-4261-a868-5be480075b8f/restart", nil)
	arc := &providers.AuthRestClient{}
	c := &nmocks.ControllerServiceClient{}

	c.On("RestartNode", mock.Anything, mock.Anything).Return(&cpb.RestartNodeResponse{
		Status: cpb.RestartStatus_ACCEPTED},
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

func Test_RestarteNodes(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	// Create a JSON payload with the necessary data.
	jsonPayload := `{"node_ids":["60285a2a-fe1d-4261-a868-5be480075b8f"]}`

	req, _ := http.NewRequest("POST", "/v1/controllers/networks/456b2743-4831-4d8d-9fbe-830df7bd59d4/restart-nodes", strings.NewReader(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	arc := &providers.AuthRestClient{}
	c := &nmocks.ControllerServiceClient{}

	restartNodeReq := &cpb.RestartNodesRequest{
		NetworkId: "456b2743-4831-4d8d-9fbe-830df7bd59d4",
		NodeIds:   []string{"60285a2a-fe1d-4261-a868-5be480075b8f"},
	}

	c.On("RestartNodes", mock.Anything, restartNodeReq).Return(&cpb.RestartNodesResponse{
		Status: cpb.RestartStatus_ACCEPTED,
	}, nil)

	r := NewRouter(&Clients{
		Controller: client.NewControllerFromClient(c),
	}, routerConfig, arc.MockAuthenticateUser).f.Engine()
	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	c.AssertExpectations(t)
}

func Test_RestarteSite(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/controllers/networks/0f37639d-3fd6-4741-b63b-9dd4f7ce55f0/sites/pamoja/restart", nil)
	arc := &providers.AuthRestClient{}
	c := &nmocks.ControllerServiceClient{}

	RestartSiteRequest := &cpb.RestartSiteRequest{
		SiteName:  "pamoja",
		NetworkId: "0f37639d-3fd6-4741-b63b-9dd4f7ce55f0",
	}

	c.On("RestartSite", mock.Anything, RestartSiteRequest).Return(&cpb.RestartSiteResponse{
		Status: cpb.RestartStatus_ACCEPTED},
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
