package rest

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	cconfig "github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/providers"
	"github.com/ukama/ukama/systems/common/rest"
	netmocks "github.com/ukama/ukama/systems/registry/network/pb/gen/mocks"
	orgpb "github.com/ukama/ukama/systems/registry/org/pb/gen"
	orgmocks "github.com/ukama/ukama/systems/registry/org/pb/gen/mocks"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ukama/ukama/systems/registry/api-gateway/pkg"
	"github.com/ukama/ukama/systems/registry/api-gateway/pkg/client"
)

var defaultCors = cors.Config{
	AllowAllOrigins: true,
}

var routerConfig = &RouterConfig{
	serverConf: &rest.HttpConfig{
		Cors: defaultCors,
	},
	httpEndpoints: &pkg.HttpEndpoints{
		NodeMetrics: "localhost:8080",
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
		Timeout: 1 * time.Second,
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

func TestGetOrg_NotFound(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/orgs/org-name", nil)
	arc := &providers.AuthRestClient{}
	n := &netmocks.NetworkServiceClient{}
	o := &orgmocks.OrgServiceClient{}
	o.On("GetByName", mock.Anything, mock.Anything).Return(nil, status.Error(codes.NotFound, "org not found"))

	r := NewRouter(&Clients{
		Registry: client.NewRegistryFromClient(n, o),
	}, routerConfig, arc.MockAuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusNotFound, w.Code)
	n.AssertExpectations(t)
}

func TestGetOrg(t *testing.T) {
	// arrange
	const orgName = "org-name"
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/orgs/"+orgName, nil)

	n := &netmocks.NetworkServiceClient{}
	o := &orgmocks.OrgServiceClient{}
	arc := &providers.AuthRestClient{}
	o.On("GetByName", mock.Anything, mock.Anything).Return(&orgpb.GetByNameResponse{
		Org: &orgpb.Organization{
			Name:  orgName,
			Owner: "owner",
		},
	}, nil)

	r := NewRouter(&Clients{
		Registry: client.NewRegistryFromClient(n, o),
	}, routerConfig, arc.MockAuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	o.AssertExpectations(t)
	assert.Contains(t, w.Body.String(), `"name":"org-name"`)
}
