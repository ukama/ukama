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
	"github.com/ukama/ukama/systems/common/uuid"
	invpb "github.com/ukama/ukama/systems/registry/invitation/pb/gen"
	imocks "github.com/ukama/ukama/systems/registry/invitation/pb/gen/mocks"
	mpb "github.com/ukama/ukama/systems/registry/member/pb/gen"
	mmocks "github.com/ukama/ukama/systems/registry/member/pb/gen/mocks"
	netmocks "github.com/ukama/ukama/systems/registry/network/pb/gen/mocks"
	nmocks "github.com/ukama/ukama/systems/registry/node/pb/gen/mocks"

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
		Network: "network:9090",
		Member:  "member:9090",
		Node:    "node:9090",
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

func TestGetMembers(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/members", nil)
	arc := &providers.AuthRestClient{}
	net := &netmocks.NetworkServiceClient{}
	node := &nmocks.NodeServiceClient{}
	mem := &mmocks.MemberServiceClient{}
	OrgId := uuid.NewV4()
	UserId := uuid.NewV4()

	mem.On("GetMembers", mock.Anything, mock.Anything).Return(&mpb.GetMembersResponse{
		Members: []*mpb.Member{{
			OrgId:         OrgId.String(),
			UserId:        UserId.String(),
			IsDeactivated: false,
		}},
	}, nil)

	r := NewRouter(&Clients{
		Node:    client.NewNodeFromClient(node),
		Member:  client.NewRegistryFromClient(mem),
		Network: client.NewNetworkRegistryFromClient(net),
	}, routerConfig, arc.MockAuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	mem.AssertExpectations(t)
}

func TestGetInvitationByOrg(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/invitation", nil)
	arc := &providers.AuthRestClient{}
	net := &netmocks.NetworkServiceClient{}
	node := &nmocks.NodeServiceClient{}
	inv := &imocks.InvitationServiceClient{}
	mem := &mmocks.MemberServiceClient{}
	invId := uuid.NewV4()

	mem.On("GetInvitationByOrg", mock.Anything, mock.Anything).Return(&invpb.GetInvitationByOrgResponse{
		Invitations: []*invpb.Invitation{{
			Id:     invId.String(),
			Org:    "ukama",
			Name:   "ukama",
			Email:  "test@ukama.com",
			Role:   invpb.RoleType_Users,
			Status: invpb.StatusType_Pending,
		}},
	}, nil)

	r := NewRouter(&Clients{
		Node:       client.NewNodeFromClient(node),
		Member:     client.NewRegistryFromClient(mem),
		Network:    client.NewNetworkRegistryFromClient(net),
		Invitation: client.NewInvitationRegistryFromClient(inv),
	}, routerConfig, arc.MockAuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	mem.AssertExpectations(t)
}
