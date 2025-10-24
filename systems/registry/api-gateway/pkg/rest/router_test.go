/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package rest

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ukama/ukama/systems/common/rest"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/registry/api-gateway/pkg"
	"github.com/ukama/ukama/systems/registry/api-gateway/pkg/client"

	cconfig "github.com/ukama/ukama/systems/common/config"
	cmocks "github.com/ukama/ukama/systems/common/mocks"
	invpb "github.com/ukama/ukama/systems/registry/invitation/pb/gen"
	imocks "github.com/ukama/ukama/systems/registry/invitation/pb/gen/mocks"
	mpb "github.com/ukama/ukama/systems/registry/member/pb/gen"
	mmocks "github.com/ukama/ukama/systems/registry/member/pb/gen/mocks"
	netpb "github.com/ukama/ukama/systems/registry/network/pb/gen"
	netmocks "github.com/ukama/ukama/systems/registry/network/pb/gen/mocks"
	nodepb "github.com/ukama/ukama/systems/registry/node/pb/gen"
	nmocks "github.com/ukama/ukama/systems/registry/node/pb/gen/mocks"
	sitepb "github.com/ukama/ukama/systems/registry/site/pb/gen"
	sitmocks "github.com/ukama/ukama/systems/registry/site/pb/gen/mocks"
)

// Test data constants
const (
	TestName          = "John Doe"
	TestEmail         = "john@example.com"
	TestRole          = "member"
	TestNetworkName   = "test-network"
	TestSiteName      = "test-site"
	TestLocation      = "test-location"
	TestNodeName      = "test-node"
	TestLatitude      = 40.7128
	TestLongitude     = -74.0060
	TestBudget        = 1000.0
	TestOverdraft     = 100.0
	TestTrafficPolicy = 1
	TestInstallDate   = "2023-01-01"
	TestLink          = "http://dev.ukama.com"
	TestCompanyName   = "ukama"
	TestCompanyEmail  = "test@ukama.com"
)

var (
	// Test UUIDs - generated once for consistency
	TestUserId     = uuid.NewV4()
	TestMemberId   = uuid.NewV4()
	TestInviteId   = uuid.NewV4()
	TestNetworkId  = uuid.NewV4()
	TestSiteId     = uuid.NewV4()
	TestNodeId     = uuid.NewV4()
	TestBackhaulId = uuid.NewV4()
	TestPowerId    = uuid.NewV4()
	TestAccessId   = uuid.NewV4()
	TestSwitchId   = uuid.NewV4()
	TestSpectrumId = uuid.NewV4()
	TestAnodeL     = uuid.NewV4()
	TestAnodeR     = uuid.NewV4()
)

// Test utility functions
type TestMocks struct {
	Auth       *cmocks.AuthClient
	Network    *netmocks.NetworkServiceClient
	Node       *nmocks.NodeServiceClient
	Member     *mmocks.MemberServiceClient
	Site       *sitmocks.SiteServiceClient
	Invitation *imocks.InvitationServiceClient
}

func NewTestMocks() *TestMocks {
	return &TestMocks{
		Auth:       &cmocks.AuthClient{},
		Network:    &netmocks.NetworkServiceClient{},
		Node:       &nmocks.NodeServiceClient{},
		Member:     &mmocks.MemberServiceClient{},
		Site:       &sitmocks.SiteServiceClient{},
		Invitation: &imocks.InvitationServiceClient{},
	}
}

func (m *TestMocks) SetupAuth() {
	m.Auth.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)
}

func (m *TestMocks) CreateTestRouter() *gin.Engine {
	clients := &Clients{
		Node:       client.NewNodeFromClient(m.Node),
		Member:     client.NewRegistryFromClient(m.Member),
		Network:    client.NewNetworkRegistryFromClient(m.Network),
		Site:       client.NewSiteRegistryFromClient(m.Site),
		Invitation: client.NewInvitationRegistryFromClient(m.Invitation),
	}
	return NewRouter(clients, routerConfig, m.Auth.AuthenticateUser).f.Engine()
}

func (m *TestMocks) AssertAllExpectations(t *testing.T) {
	m.Auth.AssertExpectations(t)
	m.Network.AssertExpectations(t)
	m.Node.AssertExpectations(t)
	m.Member.AssertExpectations(t)
	m.Site.AssertExpectations(t)
	m.Invitation.AssertExpectations(t)
}

// Test data builders
func BuildMemberRequest(userId string) string {
	return fmt.Sprintf(`{"user_uuid": "%s", "role": "%s"}`, userId, TestRole)
}

func BuildInvitationRequest() string {
	return fmt.Sprintf(`{"name": "%s", "email": "%s", "role": "%s"}`, TestName, TestEmail, TestRole)
}

func BuildNetworkRequest() string {
	return fmt.Sprintf(`{
		"network_name": "%s",
		"allowed_countries": ["US", "CA"],
		"allowed_networks": ["mesh"],
		"budget": %.1f,
		"overdraft": %.1f,
		"traffic_policy": %d,
		"payment_links": true
	}`, TestNetworkName, TestBudget, TestOverdraft, TestTrafficPolicy)
}

func BuildSiteRequest(networkId string) string {
	return fmt.Sprintf(`{
		"network_id": "%s",
		"site": "%s",
		"location": "%s",
		"backhaul_id": "%s",
		"power_id": "%s",
		"access_id": "%s",
		"switch_id": "%s",
		"spectrum_id": "%s",
		"is_deactivated": false,
		"latitude": %.4f,
		"longitude": %.4f,
		"install_date": "%s"
	}`, networkId, TestSiteName, TestLocation, TestBackhaulId.String(),
		TestPowerId.String(), TestAccessId.String(), TestSwitchId.String(),
		TestSpectrumId.String(), TestLatitude, TestLongitude, TestInstallDate)
}

func BuildNodeRequest(nodeId string) string {
	return fmt.Sprintf(`{
		"node_id": "%s",
		"name": "%s",
		"state": "operational",
		"latitude": %.4f,
		"longitude": %.4f
	}`, nodeId, TestNodeName, TestLatitude, TestLongitude)
}

func BuildAttachNodesRequest() string {
	return fmt.Sprintf(`{
		"anodel": "%s",
		"anoder": "%s"
	}`, TestAnodeL.String(), TestAnodeR.String())
}

func BuildNodeToSiteRequest(networkId, siteId string) string {
	return fmt.Sprintf(`{
		"network_id": "%s",
		"site_id": "%s"
	}`, networkId, siteId)
}

// Common test assertions
func AssertHTTPStatus(t *testing.T, w *httptest.ResponseRecorder, expectedStatus int) {
	assert.Equal(t, expectedStatus, w.Code)
}

func AssertResponseContains(t *testing.T, w *httptest.ResponseRecorder, expectedContent string) {
	assert.Contains(t, w.Body.String(), expectedContent)
}

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
		Timeout:    1 * time.Second,
		Network:    "network:9090",
		Member:     "member:9090",
		Node:       "node:9090",
		Site:       "site:9090",
		Invitation: "invitation:9090",
	})
}

func TestPingRoute(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	mocks := NewTestMocks()
	mocks.SetupAuth()

	req, _ := http.NewRequest("GET", "/ping", nil)
	r := NewRouter(testClientSet, routerConfig, mocks.Auth.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	AssertHTTPStatus(t, w, http.StatusOK)
	AssertResponseContains(t, w, "pong")
}

func TestGetMembers(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/members", nil)
	mocks := NewTestMocks()
	mocks.SetupAuth()

	mocks.Member.On("GetMembers", mock.Anything, mock.Anything).Return(&mpb.GetMembersResponse{
		Members: []*mpb.Member{{
			UserId:        TestUserId.String(),
			IsDeactivated: false,
		}},
	}, nil)

	r := mocks.CreateTestRouter()

	// act
	r.ServeHTTP(w, req)

	// assert
	AssertHTTPStatus(t, w, http.StatusOK)
	mocks.Member.AssertExpectations(t)
}

func TestGetInvitation_NotFound(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/invitations/"+TestInviteId.String(), nil)
	mocks := NewTestMocks()
	mocks.SetupAuth()

	mocks.Invitation.On("Get", mock.Anything, mock.Anything).Return(nil, status.Error(codes.NotFound, "invitation not found"))

	r := mocks.CreateTestRouter()

	// act
	r.ServeHTTP(w, req)

	// assert
	AssertHTTPStatus(t, w, http.StatusNotFound)
}
func TestGetInvitation_Found(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/invitations/"+TestInviteId.String(), nil)
	mocks := NewTestMocks()
	mocks.SetupAuth()

	mocks.Invitation.On("Get", mock.Anything, mock.Anything).Return(&invpb.GetResponse{
		Invitation: &invpb.Invitation{
			Id:    TestInviteId.String(),
			Link:  TestLink,
			Name:  TestCompanyName,
			Email: TestCompanyEmail,
		},
	}, nil)

	r := mocks.CreateTestRouter()

	// act
	r.ServeHTTP(w, req)

	// assert
	AssertHTTPStatus(t, w, http.StatusOK)
	mocks.Invitation.AssertExpectations(t)
	AssertResponseContains(t, w, TestInviteId.String())
}

// ===== MEMBER ENDPOINT TESTS =====

func TestGetMemberByUserId(t *testing.T) {
	// arrange
	userId := uuid.NewV4()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/members/user/"+userId.String(), nil)
	arc := &cmocks.AuthClient{}
	net := &netmocks.NetworkServiceClient{}
	node := &nmocks.NodeServiceClient{}
	mem := &mmocks.MemberServiceClient{}
	site := &sitmocks.SiteServiceClient{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	mem.On("GetMemberByUserId", mock.Anything, mock.Anything).Return(&mpb.GetMemberByUserIdResponse{
		Member: &mpb.Member{
			UserId:        userId.String(),
			IsDeactivated: false,
		},
	}, nil)

	r := NewRouter(&Clients{
		Node:    client.NewNodeFromClient(node),
		Member:  client.NewRegistryFromClient(mem),
		Network: client.NewNetworkRegistryFromClient(net),
		Site:    client.NewSiteRegistryFromClient(site),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	mem.AssertExpectations(t)
}

func TestGetMember(t *testing.T) {
	// arrange
	memberId := uuid.NewV4()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/members/"+memberId.String(), nil)
	arc := &cmocks.AuthClient{}
	net := &netmocks.NetworkServiceClient{}
	node := &nmocks.NodeServiceClient{}
	mem := &mmocks.MemberServiceClient{}
	site := &sitmocks.SiteServiceClient{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	mem.On("GetMember", mock.Anything, mock.Anything).Return(&mpb.MemberResponse{
		Member: &mpb.Member{
			UserId:        uuid.NewV4().String(),
			IsDeactivated: false,
		},
	}, nil)

	r := NewRouter(&Clients{
		Node:    client.NewNodeFromClient(node),
		Member:  client.NewRegistryFromClient(mem),
		Network: client.NewNetworkRegistryFromClient(net),
		Site:    client.NewSiteRegistryFromClient(site),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	mem.AssertExpectations(t)
}

func TestPostMember(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	reqBody := BuildMemberRequest(TestUserId.String())
	req, _ := http.NewRequest("POST", "/v1/members", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	mocks := NewTestMocks()
	mocks.SetupAuth()

	mocks.Member.On("AddMember", mock.Anything, mock.Anything).Return(&mpb.MemberResponse{
		Member: &mpb.Member{
			UserId:        TestUserId.String(),
			IsDeactivated: false,
		},
	}, nil)

	r := mocks.CreateTestRouter()

	// act
	r.ServeHTTP(w, req)

	// assert
	AssertHTTPStatus(t, w, http.StatusCreated)
	mocks.Member.AssertExpectations(t)
}

func TestPatchMember(t *testing.T) {
	// arrange
	memberId := uuid.NewV4()
	w := httptest.NewRecorder()
	reqBody := `{"isDeactivated": true, "role": "admin"}`
	req, _ := http.NewRequest("PATCH", "/v1/members/"+memberId.String(), strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	arc := &cmocks.AuthClient{}
	net := &netmocks.NetworkServiceClient{}
	node := &nmocks.NodeServiceClient{}
	mem := &mmocks.MemberServiceClient{}
	site := &sitmocks.SiteServiceClient{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	mem.On("UpdateMember", mock.Anything, mock.Anything).Return(&mpb.MemberResponse{}, nil)

	r := NewRouter(&Clients{
		Node:    client.NewNodeFromClient(node),
		Member:  client.NewRegistryFromClient(mem),
		Network: client.NewNetworkRegistryFromClient(net),
		Site:    client.NewSiteRegistryFromClient(site),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	mem.AssertExpectations(t)
}

func TestRemoveMember(t *testing.T) {
	// arrange
	memberId := uuid.NewV4()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/v1/members/"+memberId.String(), nil)
	arc := &cmocks.AuthClient{}
	net := &netmocks.NetworkServiceClient{}
	node := &nmocks.NodeServiceClient{}
	mem := &mmocks.MemberServiceClient{}
	site := &sitmocks.SiteServiceClient{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	mem.On("RemoveMember", mock.Anything, mock.Anything).Return(&mpb.MemberResponse{}, nil)

	r := NewRouter(&Clients{
		Node:    client.NewNodeFromClient(node),
		Member:  client.NewRegistryFromClient(mem),
		Network: client.NewNetworkRegistryFromClient(net),
		Site:    client.NewSiteRegistryFromClient(site),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	mem.AssertExpectations(t)
}

// ===== INVITATION ENDPOINT TESTS =====

func TestPostInvitation(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	reqBody := BuildInvitationRequest()
	req, _ := http.NewRequest("POST", "/v1/invitations", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	mocks := NewTestMocks()
	mocks.SetupAuth()

	mocks.Invitation.On("Add", mock.Anything, mock.Anything).Return(&invpb.AddResponse{
		Invitation: &invpb.Invitation{
			Id:    TestInviteId.String(),
			Name:  TestName,
			Email: TestEmail,
		},
	}, nil)

	r := mocks.CreateTestRouter()

	// act
	r.ServeHTTP(w, req)

	// assert
	AssertHTTPStatus(t, w, http.StatusCreated)
	mocks.Invitation.AssertExpectations(t)
}

func TestPatchInvitation(t *testing.T) {
	// arrange
	invId := uuid.NewV4()
	w := httptest.NewRecorder()
	reqBody := `{"email": "john@example.com", "status": "accepted"}`
	req, _ := http.NewRequest("PATCH", "/v1/invitations/"+invId.String(), strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	arc := &cmocks.AuthClient{}
	inv := &imocks.InvitationServiceClient{}
	net := &netmocks.NetworkServiceClient{}
	node := &nmocks.NodeServiceClient{}
	mem := &mmocks.MemberServiceClient{}
	site := &sitmocks.SiteServiceClient{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	inv.On("UpdateStatus", mock.Anything, mock.Anything).Return(&invpb.UpdateStatusResponse{}, nil)

	r := NewRouter(&Clients{
		Node:       client.NewNodeFromClient(node),
		Member:     client.NewRegistryFromClient(mem),
		Network:    client.NewNetworkRegistryFromClient(net),
		Invitation: client.NewInvitationRegistryFromClient(inv),
		Site:       client.NewSiteRegistryFromClient(site),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	inv.AssertExpectations(t)
}

func TestRemoveInvitation(t *testing.T) {
	// arrange
	invId := uuid.NewV4()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/v1/invitations/"+invId.String(), nil)
	arc := &cmocks.AuthClient{}
	inv := &imocks.InvitationServiceClient{}
	net := &netmocks.NetworkServiceClient{}
	node := &nmocks.NodeServiceClient{}
	mem := &mmocks.MemberServiceClient{}
	site := &sitmocks.SiteServiceClient{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	inv.On("Delete", mock.Anything, mock.Anything).Return(&invpb.DeleteResponse{}, nil)

	r := NewRouter(&Clients{
		Node:       client.NewNodeFromClient(node),
		Member:     client.NewRegistryFromClient(mem),
		Network:    client.NewNetworkRegistryFromClient(net),
		Invitation: client.NewInvitationRegistryFromClient(inv),
		Site:       client.NewSiteRegistryFromClient(site),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	inv.AssertExpectations(t)
}

func TestGetAllInvitations(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/invitations/", nil)
	arc := &cmocks.AuthClient{}
	inv := &imocks.InvitationServiceClient{}
	net := &netmocks.NetworkServiceClient{}
	node := &nmocks.NodeServiceClient{}
	mem := &mmocks.MemberServiceClient{}
	site := &sitmocks.SiteServiceClient{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	inv.On("GetAll", mock.Anything, mock.Anything).Return(&invpb.GetAllResponse{
		Invitations: []*invpb.Invitation{{
			Id:    uuid.NewV4().String(),
			Name:  "John Doe",
			Email: "john@example.com",
		}},
	}, nil)

	r := NewRouter(&Clients{
		Node:       client.NewNodeFromClient(node),
		Member:     client.NewRegistryFromClient(mem),
		Network:    client.NewNetworkRegistryFromClient(net),
		Invitation: client.NewInvitationRegistryFromClient(inv),
		Site:       client.NewSiteRegistryFromClient(site),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	inv.AssertExpectations(t)
}

func TestGetInvitationsByEmail(t *testing.T) {
	// arrange
	email := "john@example.com"
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/invitations/user/"+email, nil)
	arc := &cmocks.AuthClient{}
	inv := &imocks.InvitationServiceClient{}
	net := &netmocks.NetworkServiceClient{}
	node := &nmocks.NodeServiceClient{}
	mem := &mmocks.MemberServiceClient{}
	site := &sitmocks.SiteServiceClient{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	inv.On("GetByEmail", mock.Anything, mock.Anything).Return(&invpb.GetByEmailResponse{
		Invitation: &invpb.Invitation{
			Id:    uuid.NewV4().String(),
			Name:  "John Doe",
			Email: email,
		},
	}, nil)

	r := NewRouter(&Clients{
		Node:       client.NewNodeFromClient(node),
		Member:     client.NewRegistryFromClient(mem),
		Network:    client.NewNetworkRegistryFromClient(net),
		Invitation: client.NewInvitationRegistryFromClient(inv),
		Site:       client.NewSiteRegistryFromClient(site),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	inv.AssertExpectations(t)
}

// ===== NETWORK ENDPOINT TESTS =====

func TestGetNetworks(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/networks", nil)
	arc := &cmocks.AuthClient{}
	net := &netmocks.NetworkServiceClient{}
	node := &nmocks.NodeServiceClient{}
	mem := &mmocks.MemberServiceClient{}
	site := &sitmocks.SiteServiceClient{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	net.On("GetAll", mock.Anything, mock.Anything).Return(&netpb.GetNetworksResponse{
		Networks: []*netpb.Network{{
			Id:   uuid.NewV4().String(),
			Name: "test-network",
		}},
	}, nil)

	r := NewRouter(&Clients{
		Node:    client.NewNodeFromClient(node),
		Member:  client.NewRegistryFromClient(mem),
		Network: client.NewNetworkRegistryFromClient(net),
		Site:    client.NewSiteRegistryFromClient(site),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	net.AssertExpectations(t)
}

func TestGetDefaultNetwork(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/networks/default", nil)
	arc := &cmocks.AuthClient{}
	net := &netmocks.NetworkServiceClient{}
	node := &nmocks.NodeServiceClient{}
	mem := &mmocks.MemberServiceClient{}
	site := &sitmocks.SiteServiceClient{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	net.On("GetDefault", mock.Anything, mock.Anything).Return(&netpb.GetDefaultResponse{
		Network: &netpb.Network{
			Id:   uuid.NewV4().String(),
			Name: "default-network",
		},
	}, nil)

	r := NewRouter(&Clients{
		Node:    client.NewNodeFromClient(node),
		Member:  client.NewRegistryFromClient(mem),
		Network: client.NewNetworkRegistryFromClient(net),
		Site:    client.NewSiteRegistryFromClient(site),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	net.AssertExpectations(t)
}

func TestPostNetwork(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	reqBody := BuildNetworkRequest()
	req, _ := http.NewRequest("POST", "/v1/networks", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	mocks := NewTestMocks()
	mocks.SetupAuth()

	mocks.Network.On("Add", mock.Anything, mock.Anything).Return(&netpb.AddResponse{
		Network: &netpb.Network{
			Id:   TestNetworkId.String(),
			Name: TestNetworkName,
		},
	}, nil)

	r := mocks.CreateTestRouter()

	// act
	r.ServeHTTP(w, req)

	// assert
	AssertHTTPStatus(t, w, http.StatusCreated)
	mocks.Network.AssertExpectations(t)
}

func TestGetNetwork(t *testing.T) {
	// arrange
	networkId := uuid.NewV4()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/networks/"+networkId.String(), nil)
	arc := &cmocks.AuthClient{}
	net := &netmocks.NetworkServiceClient{}
	node := &nmocks.NodeServiceClient{}
	mem := &mmocks.MemberServiceClient{}
	site := &sitmocks.SiteServiceClient{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	net.On("Get", mock.Anything, mock.Anything).Return(&netpb.GetResponse{
		Network: &netpb.Network{
			Id:   networkId.String(),
			Name: "test-network",
		},
	}, nil)

	r := NewRouter(&Clients{
		Node:    client.NewNodeFromClient(node),
		Member:  client.NewRegistryFromClient(mem),
		Network: client.NewNetworkRegistryFromClient(net),
		Site:    client.NewSiteRegistryFromClient(site),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	net.AssertExpectations(t)
}

func TestSetNetworkDefault(t *testing.T) {
	// arrange
	networkId := uuid.NewV4()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/v1/networks/"+networkId.String(), nil)
	arc := &cmocks.AuthClient{}
	net := &netmocks.NetworkServiceClient{}
	node := &nmocks.NodeServiceClient{}
	mem := &mmocks.MemberServiceClient{}
	site := &sitmocks.SiteServiceClient{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	net.On("SetDefault", mock.Anything, mock.Anything).Return(&netpb.SetDefaultResponse{}, nil)

	r := NewRouter(&Clients{
		Node:    client.NewNodeFromClient(node),
		Member:  client.NewRegistryFromClient(mem),
		Network: client.NewNetworkRegistryFromClient(net),
		Site:    client.NewSiteRegistryFromClient(site),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	net.AssertExpectations(t)
}

// ===== SITE ENDPOINT TESTS =====

func TestGetSites(t *testing.T) {
	// arrange
	networkId := uuid.NewV4()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/sites?network_id="+networkId.String()+"&is_deactivated=false", nil)
	arc := &cmocks.AuthClient{}
	net := &netmocks.NetworkServiceClient{}
	node := &nmocks.NodeServiceClient{}
	mem := &mmocks.MemberServiceClient{}
	site := &sitmocks.SiteServiceClient{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	site.On("List", mock.Anything, mock.Anything).Return(&sitepb.ListResponse{
		Sites: []*sitepb.Site{{
			Id:   uuid.NewV4().String(),
			Name: "test-site",
		}},
	}, nil)

	r := NewRouter(&Clients{
		Node:    client.NewNodeFromClient(node),
		Member:  client.NewRegistryFromClient(mem),
		Network: client.NewNetworkRegistryFromClient(net),
		Site:    client.NewSiteRegistryFromClient(site),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	site.AssertExpectations(t)
}

func TestPostSite(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	reqBody := BuildSiteRequest(TestNetworkId.String())
	req, _ := http.NewRequest("POST", "/v1/sites", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	mocks := NewTestMocks()
	mocks.SetupAuth()

	mocks.Site.On("Add", mock.Anything, mock.Anything).Return(&sitepb.AddResponse{
		Site: &sitepb.Site{
			Id:   TestSiteId.String(),
			Name: TestSiteName,
		},
	}, nil)

	r := mocks.CreateTestRouter()

	// act
	r.ServeHTTP(w, req)

	// assert
	AssertHTTPStatus(t, w, http.StatusCreated)
	mocks.Site.AssertExpectations(t)
}

func TestGetSite(t *testing.T) {
	// arrange
	siteId := uuid.NewV4()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/sites/"+siteId.String(), nil)
	arc := &cmocks.AuthClient{}
	net := &netmocks.NetworkServiceClient{}
	node := &nmocks.NodeServiceClient{}
	mem := &mmocks.MemberServiceClient{}
	site := &sitmocks.SiteServiceClient{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	site.On("Get", mock.Anything, mock.Anything).Return(&sitepb.GetResponse{
		Site: &sitepb.Site{
			Id:   siteId.String(),
			Name: "test-site",
		},
	}, nil)

	r := NewRouter(&Clients{
		Node:    client.NewNodeFromClient(node),
		Member:  client.NewRegistryFromClient(mem),
		Network: client.NewNetworkRegistryFromClient(net),
		Site:    client.NewSiteRegistryFromClient(site),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	site.AssertExpectations(t)
}

func TestUpdateSite(t *testing.T) {
	// arrange
	siteId := uuid.NewV4()
	w := httptest.NewRecorder()
	reqBody := `{"name": "updated-site"}`
	req, _ := http.NewRequest("PATCH", "/v1/sites/"+siteId.String(), strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	arc := &cmocks.AuthClient{}
	net := &netmocks.NetworkServiceClient{}
	node := &nmocks.NodeServiceClient{}
	mem := &mmocks.MemberServiceClient{}
	site := &sitmocks.SiteServiceClient{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	site.On("Update", mock.Anything, mock.Anything).Return(&sitepb.UpdateResponse{
		Site: &sitepb.Site{
			Id:   siteId.String(),
			Name: "updated-site",
		},
	}, nil)

	r := NewRouter(&Clients{
		Node:    client.NewNodeFromClient(node),
		Member:  client.NewRegistryFromClient(mem),
		Network: client.NewNetworkRegistryFromClient(net),
		Site:    client.NewSiteRegistryFromClient(site),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	site.AssertExpectations(t)
}

// ===== NODE ENDPOINT TESTS =====

func TestGetNodes(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/nodes", nil)
	arc := &cmocks.AuthClient{}
	net := &netmocks.NetworkServiceClient{}
	node := &nmocks.NodeServiceClient{}
	mem := &mmocks.MemberServiceClient{}
	site := &sitmocks.SiteServiceClient{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	node.On("GetNodes", mock.Anything, mock.Anything).Return(&nodepb.GetNodesResponse{
		Nodes: []*nodepb.Node{{
			Id:   uuid.NewV4().String(),
			Name: "test-node",
		}},
	}, nil)

	r := NewRouter(&Clients{
		Node:    client.NewNodeFromClient(node),
		Member:  client.NewRegistryFromClient(mem),
		Network: client.NewNetworkRegistryFromClient(net),
		Site:    client.NewSiteRegistryFromClient(site),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	node.AssertExpectations(t)
}

func TestGetNodesByState(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/nodes/state?connectivity=online&state=operational", nil)
	arc := &cmocks.AuthClient{}
	net := &netmocks.NetworkServiceClient{}
	node := &nmocks.NodeServiceClient{}
	mem := &mmocks.MemberServiceClient{}
	site := &sitmocks.SiteServiceClient{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	node.On("GetNodesByState", mock.Anything, mock.Anything, mock.Anything).Return(&nodepb.GetNodesResponse{
		Nodes: []*nodepb.Node{{
			Id:   uuid.NewV4().String(),
			Name: "test-node",
		}},
	}, nil)

	r := NewRouter(&Clients{
		Node:    client.NewNodeFromClient(node),
		Member:  client.NewRegistryFromClient(mem),
		Network: client.NewNetworkRegistryFromClient(net),
		Site:    client.NewSiteRegistryFromClient(site),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	node.AssertExpectations(t)
}

func TestGetNode(t *testing.T) {
	// arrange
	nodeId := uuid.NewV4()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/nodes/"+nodeId.String(), nil)
	arc := &cmocks.AuthClient{}
	net := &netmocks.NetworkServiceClient{}
	node := &nmocks.NodeServiceClient{}
	mem := &mmocks.MemberServiceClient{}
	site := &sitmocks.SiteServiceClient{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	node.On("GetNode", mock.Anything, mock.Anything).Return(&nodepb.GetNodeResponse{
		Node: &nodepb.Node{
			Id:   nodeId.String(),
			Name: "test-node",
		},
	}, nil)

	r := NewRouter(&Clients{
		Node:    client.NewNodeFromClient(node),
		Member:  client.NewRegistryFromClient(mem),
		Network: client.NewNetworkRegistryFromClient(net),
		Site:    client.NewSiteRegistryFromClient(site),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	node.AssertExpectations(t)
}

func TestGetSiteNodes(t *testing.T) {
	// arrange
	siteId := uuid.NewV4()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/nodes/sites/"+siteId.String(), nil)
	arc := &cmocks.AuthClient{}
	net := &netmocks.NetworkServiceClient{}
	node := &nmocks.NodeServiceClient{}
	mem := &mmocks.MemberServiceClient{}
	site := &sitmocks.SiteServiceClient{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	node.On("GetNodesForSite", mock.Anything, mock.Anything).Return(&nodepb.GetBySiteResponse{
		Nodes: []*nodepb.Node{{
			Id:   uuid.NewV4().String(),
			Name: "test-node",
		}},
	}, nil)

	r := NewRouter(&Clients{
		Node:    client.NewNodeFromClient(node),
		Member:  client.NewRegistryFromClient(mem),
		Network: client.NewNetworkRegistryFromClient(net),
		Site:    client.NewSiteRegistryFromClient(site),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	node.AssertExpectations(t)
}

func TestGetNetworkNodes(t *testing.T) {
	// arrange
	networkId := uuid.NewV4()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/nodes/networks/"+networkId.String(), nil)
	arc := &cmocks.AuthClient{}
	net := &netmocks.NetworkServiceClient{}
	node := &nmocks.NodeServiceClient{}
	mem := &mmocks.MemberServiceClient{}
	site := &sitmocks.SiteServiceClient{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	node.On("GetNodesForNetwork", mock.Anything, mock.Anything).Return(&nodepb.GetByNetworkResponse{
		Nodes: []*nodepb.Node{{
			Id:   uuid.NewV4().String(),
			Name: "test-node",
		}},
	}, nil)

	r := NewRouter(&Clients{
		Node:    client.NewNodeFromClient(node),
		Member:  client.NewRegistryFromClient(mem),
		Network: client.NewNetworkRegistryFromClient(net),
		Site:    client.NewSiteRegistryFromClient(site),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	node.AssertExpectations(t)
}

func TestListNodes(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/nodes/list?type=mesh&state=operational&connectivity=online", nil)
	arc := &cmocks.AuthClient{}
	net := &netmocks.NetworkServiceClient{}
	node := &nmocks.NodeServiceClient{}
	mem := &mmocks.MemberServiceClient{}
	site := &sitmocks.SiteServiceClient{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	node.On("List", mock.Anything, mock.Anything).Return(&nodepb.ListResponse{
		Nodes: []*nodepb.Node{{
			Id:   uuid.NewV4().String(),
			Name: "test-node",
		}},
	}, nil)

	r := NewRouter(&Clients{
		Node:    client.NewNodeFromClient(node),
		Member:  client.NewRegistryFromClient(mem),
		Network: client.NewNetworkRegistryFromClient(net),
		Site:    client.NewSiteRegistryFromClient(site),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	node.AssertExpectations(t)
}

func TestPostAddNode(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	reqBody := BuildNodeRequest(TestNodeId.String())
	req, _ := http.NewRequest("POST", "/v1/nodes", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	mocks := NewTestMocks()
	mocks.SetupAuth()

	mocks.Node.On("AddNode", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&nodepb.AddNodeResponse{
		Node: &nodepb.Node{
			Id:   TestNodeId.String(),
			Name: TestNodeName,
		},
	}, nil)

	r := mocks.CreateTestRouter()

	// act
	r.ServeHTTP(w, req)

	// assert
	AssertHTTPStatus(t, w, http.StatusCreated)
	mocks.Node.AssertExpectations(t)
}

func TestPutUpdateNode(t *testing.T) {
	// arrange
	nodeId := uuid.NewV4()
	w := httptest.NewRecorder()
	reqBody := `{
		"name": "updated-node",
		"latitude": 40.7128,
		"longitude": -74.0060
	}`
	req, _ := http.NewRequest("PUT", "/v1/nodes/"+nodeId.String(), strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	arc := &cmocks.AuthClient{}
	net := &netmocks.NetworkServiceClient{}
	node := &nmocks.NodeServiceClient{}
	mem := &mmocks.MemberServiceClient{}
	site := &sitmocks.SiteServiceClient{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	node.On("UpdateNode", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&nodepb.UpdateNodeResponse{
		Node: &nodepb.Node{
			Id:   nodeId.String(),
			Name: "updated-node",
		},
	}, nil)

	r := NewRouter(&Clients{
		Node:    client.NewNodeFromClient(node),
		Member:  client.NewRegistryFromClient(mem),
		Network: client.NewNetworkRegistryFromClient(net),
		Site:    client.NewSiteRegistryFromClient(site),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	node.AssertExpectations(t)
}

func TestPatchUpdateNodeState(t *testing.T) {
	// arrange
	nodeId := uuid.NewV4()
	w := httptest.NewRecorder()
	reqBody := `{"state": "faulty"}`
	req, _ := http.NewRequest("PATCH", "/v1/nodes/"+nodeId.String(), strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	arc := &cmocks.AuthClient{}
	net := &netmocks.NetworkServiceClient{}
	node := &nmocks.NodeServiceClient{}
	mem := &mmocks.MemberServiceClient{}
	site := &sitmocks.SiteServiceClient{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	node.On("UpdateNodeState", mock.Anything, mock.Anything, mock.Anything).Return(&nodepb.UpdateNodeResponse{
		Node: &nodepb.Node{
			Id:   nodeId.String(),
			Name: "test-node",
		},
	}, nil)

	r := NewRouter(&Clients{
		Node:    client.NewNodeFromClient(node),
		Member:  client.NewRegistryFromClient(mem),
		Network: client.NewNetworkRegistryFromClient(net),
		Site:    client.NewSiteRegistryFromClient(site),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	node.AssertExpectations(t)
}

func TestDeleteNode(t *testing.T) {
	// arrange
	nodeId := uuid.NewV4()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/v1/nodes/"+nodeId.String(), nil)
	arc := &cmocks.AuthClient{}
	net := &netmocks.NetworkServiceClient{}
	node := &nmocks.NodeServiceClient{}
	mem := &mmocks.MemberServiceClient{}
	site := &sitmocks.SiteServiceClient{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	node.On("DeleteNode", mock.Anything, mock.Anything).Return(&nodepb.DeleteNodeResponse{}, nil)

	r := NewRouter(&Clients{
		Node:    client.NewNodeFromClient(node),
		Member:  client.NewRegistryFromClient(mem),
		Network: client.NewNetworkRegistryFromClient(net),
		Site:    client.NewSiteRegistryFromClient(site),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	node.AssertExpectations(t)
}

func TestPostAttachedNodes(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	reqBody := BuildAttachNodesRequest()
	req, _ := http.NewRequest("POST", "/v1/nodes/"+TestNodeId.String()+"/attach", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	mocks := NewTestMocks()
	mocks.SetupAuth()

	mocks.Node.On("AttachNodes", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&nodepb.AttachNodesResponse{}, nil)

	r := mocks.CreateTestRouter()

	// act
	r.ServeHTTP(w, req)

	// assert
	AssertHTTPStatus(t, w, http.StatusCreated)
	mocks.Node.AssertExpectations(t)
}

func TestDeleteAttachedNode(t *testing.T) {
	// arrange
	nodeId := uuid.NewV4()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/v1/nodes/"+nodeId.String()+"/attach", nil)
	arc := &cmocks.AuthClient{}
	net := &netmocks.NetworkServiceClient{}
	node := &nmocks.NodeServiceClient{}
	mem := &mmocks.MemberServiceClient{}
	site := &sitmocks.SiteServiceClient{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	node.On("DetachNode", mock.Anything, mock.Anything).Return(&nodepb.DetachNodeResponse{}, nil)

	r := NewRouter(&Clients{
		Node:    client.NewNodeFromClient(node),
		Member:  client.NewRegistryFromClient(mem),
		Network: client.NewNetworkRegistryFromClient(net),
		Site:    client.NewSiteRegistryFromClient(site),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	node.AssertExpectations(t)
}

func TestPostNodeToSite(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	reqBody := BuildNodeToSiteRequest(TestNetworkId.String(), TestSiteId.String())
	req, _ := http.NewRequest("POST", "/v1/nodes/"+TestNodeId.String()+"/sites", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	mocks := NewTestMocks()
	mocks.SetupAuth()

	mocks.Node.On("AddNodeToSite", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&nodepb.AddNodeToSiteResponse{}, nil)

	r := mocks.CreateTestRouter()

	// act
	r.ServeHTTP(w, req)

	// assert
	AssertHTTPStatus(t, w, http.StatusCreated)
	mocks.Node.AssertExpectations(t)
}

func TestDeleteNodeFromSite(t *testing.T) {
	// arrange
	nodeId := uuid.NewV4()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/v1/nodes/"+nodeId.String()+"/sites", nil)
	arc := &cmocks.AuthClient{}
	net := &netmocks.NetworkServiceClient{}
	node := &nmocks.NodeServiceClient{}
	mem := &mmocks.MemberServiceClient{}
	site := &sitmocks.SiteServiceClient{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	node.On("ReleaseNodeFromSite", mock.Anything, mock.Anything).Return(&nodepb.ReleaseNodeFromSiteResponse{}, nil)

	r := NewRouter(&Clients{
		Node:    client.NewNodeFromClient(node),
		Member:  client.NewRegistryFromClient(mem),
		Network: client.NewNetworkRegistryFromClient(net),
		Site:    client.NewSiteRegistryFromClient(site),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	node.AssertExpectations(t)
}
