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
	"github.com/ukama/ukama/systems/common/ukama"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	cconfig "github.com/ukama/ukama/systems/common/config"
	crest "github.com/ukama/ukama/systems/common/rest"
	"github.com/ukama/ukama/systems/node/api-gateway/pkg"
	"github.com/ukama/ukama/systems/node/api-gateway/pkg/client"
	cfgPb "github.com/ukama/ukama/systems/node/configurator/pb/gen"
	cmocks "github.com/ukama/ukama/systems/node/configurator/pb/gen/mocks"
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
		Timeout:      1 * time.Second,
		Controller:   "0.0.0.0:9092",
		Configurator: "0.0.0.0:9080",
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
	node := ukama.NewVirtualHomeNodeId().String()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/controllers/restartNode/"+node, nil)
	arc := &providers.AuthRestClient{}
	c := &nmocks.ControllerServiceClient{}
	cfg := &cmocks.ConfiguratorServiceClient{}

	c.On("RestartNode", mock.Anything, mock.Anything).Return(&cpb.RestartNodeResponse{
		Status: cpb.RestartStatus_RESTART_STATUS_SUCCESS},
		nil)

	r := NewRouter(&Clients{
		Controller:   client.NewControllerFromClient(c),
		Configurator: client.NewConfiguratorFromClient(cfg),
	}, routerConfig, arc.MockAuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusAccepted, w.Code)
	c.AssertExpectations(t)
}

func Test_postConfigEventHandler(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("POST", "/v1/configurator/config", strings.NewReader("{\"name\": \"config\"}"))
	arc := &providers.AuthRestClient{}
	c := &nmocks.ControllerServiceClient{}
	cfg := &cmocks.ConfiguratorServiceClient{}

	cfg.On("ConfigEvent", mock.Anything, mock.Anything).Return(&cfgPb.ConfigStoreEventResponse{},
		nil)

	r := NewRouter(&Clients{
		Controller:   client.NewControllerFromClient(c),
		Configurator: client.NewConfiguratorFromClient(cfg),
	}, routerConfig, arc.MockAuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusAccepted, w.Code)
	c.AssertExpectations(t)
}

func Test_postConfigApplyVersionHandler(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	hash := "1c924398265578d35e2b16adca25dcc021923c89"
	req, _ := http.NewRequest("POST", "/v1/configurator/config/apply/"+hash, nil)
	arc := &providers.AuthRestClient{}
	c := &nmocks.ControllerServiceClient{}
	cfg := &cmocks.ConfiguratorServiceClient{}

	cfg.On("ApplyConfig", mock.Anything, mock.Anything).Return(&cfgPb.ApplyConfigResponse{},
		nil)

	r := NewRouter(&Clients{
		Controller:   client.NewControllerFromClient(c),
		Configurator: client.NewConfiguratorFromClient(cfg),
	}, routerConfig, arc.MockAuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusAccepted, w.Code)
	c.AssertExpectations(t)
}

func Test_getRunningConfigVersionHandler(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	node := ukama.NewVirtualHomeNodeId().String()
	req, _ := http.NewRequest("GET", "/v1/configurator/config/node/"+node, nil)
	arc := &providers.AuthRestClient{}
	c := &nmocks.ControllerServiceClient{}
	cfg := &cmocks.ConfiguratorServiceClient{}

	cfg.On("GetConfigVersion", mock.Anything, mock.Anything).Return(&cfgPb.ConfigVersionResponse{
		NodeId:     node,
		Status:     "Success",
		Commit:     "1c924398265578d35e2b16adca25dcc021923c89",
		LastCommit: "1c924398265578d35e2b16adca25dcc021923c90",
		LastStatus: "Published",
	},
		nil)

	r := NewRouter(&Clients{
		Controller:   client.NewControllerFromClient(c),
		Configurator: client.NewConfiguratorFromClient(cfg),
	}, routerConfig, arc.MockAuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	if assert.Equal(t, http.StatusOK, w.Code) {
		assert.Contains(t, w.Body.String(), node)
	}
	c.AssertExpectations(t)
}
