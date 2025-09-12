/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package rest

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ukama/ukama/systems/common/rest"
	"github.com/ukama/ukama/systems/nucleus/api-gateway/pkg"
	"github.com/ukama/ukama/systems/nucleus/api-gateway/pkg/client"

	cconfig "github.com/ukama/ukama/systems/common/config"
	cmocks "github.com/ukama/ukama/systems/common/mocks"
	orgpb "github.com/ukama/ukama/systems/nucleus/org/pb/gen"
	orgmocks "github.com/ukama/ukama/systems/nucleus/org/pb/gen/mocks"
	userspb "github.com/ukama/ukama/systems/nucleus/user/pb/gen"
	usermocks "github.com/ukama/ukama/systems/nucleus/user/pb/gen/mocks"
)

var (
	testClientSet *Clients

	defaultCors = cors.Config{
		AllowAllOrigins: true,
	}

	routerConfig = &RouterConfig{
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

	// Test data constants
	testOrgName       = "org-name"
	testOrgName2      = "test-org"
	testOrgId         = "org-id-123"
	testOrgId2        = "org-uuid-1"
	testUserId        = "user-id-123"
	testUserUUID      = "user-uuid-123"
	testUserUUID2     = "owner-uuid-123"
	testAuthId        = "auth-id-123"
	testOwner         = "owner"
	testEmail         = "test@example.com"
	testEmailNotFound = "notfound@example.com"
	testUserName      = "Test User"
	testPhone         = "1234567890"
	testCertificate   = "test-cert"
	testCountry       = "us"
	testCurrency      = "usd"
	testOrg1Name      = "org1"
	testOrg2Name      = "org2"

	// Test JSON data
	testOrgData = `{
		"name": "test-org",
		"owner_uuid": "owner-uuid-123",
		"certificate": "test-cert",
		"country": "us",
		"currency": "usd"
	}`

	testUserData = `{
		"name": "Test User",
		"email": "test@example.com",
		"phone": "1234567890",
		"auth_id": "auth-id-123"
	}`
)

func init() {
	gin.SetMode(gin.TestMode)
	// testClientSet = NewClientsSet(&pkg.GrpcEndpoints{
	// 	Timeout: 1 * time.Second,
	// })
}

func TestPingRoute(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	arc := &cmocks.AuthClient{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	req, _ := http.NewRequest("GET", "/ping", nil)
	r := NewRouter(testClientSet, routerConfig, arc.AuthenticateUser).f.Engine()
	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "pong")
}

func TestGetOrg_NotFound(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/orgs/"+testOrgName, nil)
	arc := &cmocks.AuthClient{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	o := &orgmocks.OrgServiceClient{}
	u := &usermocks.UserServiceClient{}
	o.On("GetByName", mock.Anything, mock.Anything).Return(nil, status.Error(codes.NotFound, "org not found"))

	r := NewRouter(&Clients{
		Organization: client.NewOrgRegistryFromClient(o),
		User:         client.NewUserRegistryFromClient(u),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusNotFound, w.Code)

}

func TestGetOrg(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/orgs/"+testOrgName, nil)
	o := &orgmocks.OrgServiceClient{}
	u := &usermocks.UserServiceClient{}
	arc := &cmocks.AuthClient{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	o.On("GetByName", mock.Anything, mock.Anything).Return(&orgpb.GetByNameResponse{
		Org: &orgpb.Organization{
			Name:  testOrgName,
			Owner: testOwner,
		},
	}, nil)

	r := NewRouter(&Clients{
		Organization: client.NewOrgRegistryFromClient(o),
		User:         client.NewUserRegistryFromClient(u),
	}, routerConfig, arc.AuthenticateUser).f.Engine()
	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	o.AssertExpectations(t)
	assert.Contains(t, w.Body.String(), `"name":"`+testOrgName+`"`)
}

func TestGetOrgs_Success(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/orgs?user_uuid="+testUserUUID, nil)
	o := &orgmocks.OrgServiceClient{}
	u := &usermocks.UserServiceClient{}
	arc := &cmocks.AuthClient{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	o.On("GetByUser", mock.Anything, mock.Anything).Return(&orgpb.GetByUserResponse{
		OwnerOf: []*orgpb.Organization{
			{
				Name:  testOrg1Name,
				Owner: testUserUUID,
			},
			{
				Name:  testOrg2Name,
				Owner: testUserUUID,
			},
		},
	}, nil)

	r := NewRouter(&Clients{
		Organization: client.NewOrgRegistryFromClient(o),
		User:         client.NewUserRegistryFromClient(u),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	o.AssertExpectations(t)
	assert.Contains(t, w.Body.String(), `"name":"`+testOrg1Name+`"`)
	assert.Contains(t, w.Body.String(), `"name":"`+testOrg2Name+`"`)
}

func TestGetOrgs_MissingUserUUID(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/orgs", nil)
	o := &orgmocks.OrgServiceClient{}
	u := &usermocks.UserServiceClient{}
	arc := &cmocks.AuthClient{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	r := NewRouter(&Clients{
		Organization: client.NewOrgRegistryFromClient(o),
		User:         client.NewUserRegistryFromClient(u),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "UserUuid")
}

func TestGetOrgs_NotFound(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/orgs?user_uuid="+testUserUUID, nil)
	o := &orgmocks.OrgServiceClient{}
	u := &usermocks.UserServiceClient{}
	arc := &cmocks.AuthClient{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	o.On("GetByUser", mock.Anything, mock.Anything).Return(nil, status.Error(codes.NotFound, "no orgs found"))

	r := NewRouter(&Clients{
		Organization: client.NewOrgRegistryFromClient(o),
		User:         client.NewUserRegistryFromClient(u),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPostOrg_Success(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/orgs", strings.NewReader(testOrgData))
	req.Header.Set("Content-Type", "application/json")
	o := &orgmocks.OrgServiceClient{}
	u := &usermocks.UserServiceClient{}
	arc := &cmocks.AuthClient{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	o.On("Add", mock.Anything, mock.Anything).Return(&orgpb.AddResponse{
		Org: &orgpb.Organization{
			Name:        testOrgName2,
			Owner:       testUserUUID2,
			Certificate: testCertificate,
			Country:     testCountry,
			Currency:    testCurrency,
		},
	}, nil)

	r := NewRouter(&Clients{
		Organization: client.NewOrgRegistryFromClient(o),
		User:         client.NewUserRegistryFromClient(u),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusCreated, w.Code)
	o.AssertExpectations(t)
	assert.Contains(t, w.Body.String(), `"name":"`+testOrgName2+`"`)
}

func TestPostOrg_Error(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/orgs", strings.NewReader(testOrgData))
	req.Header.Set("Content-Type", "application/json")
	o := &orgmocks.OrgServiceClient{}
	u := &usermocks.UserServiceClient{}
	arc := &cmocks.AuthClient{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	o.On("Add", mock.Anything, mock.Anything).Return(nil, status.Error(codes.Internal, "internal error"))

	r := NewRouter(&Clients{
		Organization: client.NewOrgRegistryFromClient(o),
		User:         client.NewUserRegistryFromClient(u),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestUpdateOrgToUser_Success(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/v1/orgs/"+testOrgId+"/users/"+testUserId, nil)
	o := &orgmocks.OrgServiceClient{}
	u := &usermocks.UserServiceClient{}
	arc := &cmocks.AuthClient{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	o.On("UpdateOrgForUser", mock.Anything, mock.Anything).Return(&orgpb.UpdateOrgForUserResponse{}, nil)

	r := NewRouter(&Clients{
		Organization: client.NewOrgRegistryFromClient(o),
		User:         client.NewUserRegistryFromClient(u),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusCreated, w.Code)
	o.AssertExpectations(t)
	// Response should be successful
}

func TestUpdateOrgToUser_Error(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/v1/orgs/"+testOrgId+"/users/"+testUserId, nil)
	o := &orgmocks.OrgServiceClient{}
	u := &usermocks.UserServiceClient{}
	arc := &cmocks.AuthClient{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	o.On("UpdateOrgForUser", mock.Anything, mock.Anything).Return(nil, status.Error(codes.NotFound, "user or org not found"))

	r := NewRouter(&Clients{
		Organization: client.NewOrgRegistryFromClient(o),
		User:         client.NewUserRegistryFromClient(u),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRemoveUserFromOrg_Success(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/v1/orgs/"+testOrgId+"/users/"+testUserId, nil)
	o := &orgmocks.OrgServiceClient{}
	u := &usermocks.UserServiceClient{}
	arc := &cmocks.AuthClient{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	o.On("RemoveOrgForUser", mock.Anything, mock.Anything).Return(&orgpb.RemoveOrgForUserResponse{}, nil)

	r := NewRouter(&Clients{
		Organization: client.NewOrgRegistryFromClient(o),
		User:         client.NewUserRegistryFromClient(u),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusCreated, w.Code)
	o.AssertExpectations(t)
	// Response should be successful
}

func TestRemoveUserFromOrg_Error(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/v1/orgs/"+testOrgId+"/users/"+testUserId, nil)
	o := &orgmocks.OrgServiceClient{}
	u := &usermocks.UserServiceClient{}
	arc := &cmocks.AuthClient{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	o.On("RemoveOrgForUser", mock.Anything, mock.Anything).Return(nil, status.Error(codes.NotFound, "user or org not found"))

	r := NewRouter(&Clients{
		Organization: client.NewOrgRegistryFromClient(o),
		User:         client.NewUserRegistryFromClient(u),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetUserByEmail_Success(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/users/email/"+testEmail, nil)
	o := &orgmocks.OrgServiceClient{}
	u := &usermocks.UserServiceClient{}
	arc := &cmocks.AuthClient{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	u.On("GetByEmail", mock.Anything, mock.Anything).Return(&userspb.GetResponse{
		User: &userspb.User{
			Id:    testUserUUID,
			Name:  testUserName,
			Email: testEmail,
		},
	}, nil)

	r := NewRouter(&Clients{
		Organization: client.NewOrgRegistryFromClient(o),
		User:         client.NewUserRegistryFromClient(u),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	u.AssertExpectations(t)
	assert.Contains(t, w.Body.String(), `"email":"`+testEmail+`"`)
}

func TestGetUserByEmail_NotFound(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/users/email/"+testEmailNotFound, nil)
	o := &orgmocks.OrgServiceClient{}
	u := &usermocks.UserServiceClient{}
	arc := &cmocks.AuthClient{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	u.On("GetByEmail", mock.Anything, mock.Anything).Return(nil, status.Error(codes.NotFound, "user not found"))

	r := NewRouter(&Clients{
		Organization: client.NewOrgRegistryFromClient(o),
		User:         client.NewUserRegistryFromClient(u),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetUser_Success(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/users/"+testUserId, nil)
	o := &orgmocks.OrgServiceClient{}
	u := &usermocks.UserServiceClient{}
	arc := &cmocks.AuthClient{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	u.On("Get", mock.Anything, mock.Anything).Return(&userspb.GetResponse{
		User: &userspb.User{
			Id:    testUserId,
			Name:  testUserName,
			Email: testEmail,
		},
	}, nil)

	r := NewRouter(&Clients{
		Organization: client.NewOrgRegistryFromClient(o),
		User:         client.NewUserRegistryFromClient(u),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	u.AssertExpectations(t)
	assert.Contains(t, w.Body.String(), `"id":"`+testUserId+`"`)
}

func TestGetUser_NotFound(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/users/"+testUserId, nil)
	o := &orgmocks.OrgServiceClient{}
	u := &usermocks.UserServiceClient{}
	arc := &cmocks.AuthClient{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	u.On("Get", mock.Anything, mock.Anything).Return(nil, status.Error(codes.NotFound, "user not found"))

	r := NewRouter(&Clients{
		Organization: client.NewOrgRegistryFromClient(o),
		User:         client.NewUserRegistryFromClient(u),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetUserByAuthId_Success(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/users/auth/"+testAuthId, nil)
	o := &orgmocks.OrgServiceClient{}
	u := &usermocks.UserServiceClient{}
	arc := &cmocks.AuthClient{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	u.On("GetByAuthId", mock.Anything, mock.Anything).Return(&userspb.GetResponse{
		User: &userspb.User{
			Id:     testUserUUID,
			Name:   testUserName,
			Email:  testEmail,
			AuthId: testAuthId,
		},
	}, nil)

	r := NewRouter(&Clients{
		Organization: client.NewOrgRegistryFromClient(o),
		User:         client.NewUserRegistryFromClient(u),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	u.AssertExpectations(t)
	assert.Contains(t, w.Body.String(), `"auth_id":"`+testAuthId+`"`)
}

func TestGetUserByAuthId_NotFound(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/users/auth/"+testAuthId, nil)
	o := &orgmocks.OrgServiceClient{}
	u := &usermocks.UserServiceClient{}
	arc := &cmocks.AuthClient{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	u.On("GetByAuthId", mock.Anything, mock.Anything).Return(nil, status.Error(codes.NotFound, "user not found"))

	r := NewRouter(&Clients{
		Organization: client.NewOrgRegistryFromClient(o),
		User:         client.NewUserRegistryFromClient(u),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestWhoami_Success(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/users/whoami/"+testUserId, nil)
	o := &orgmocks.OrgServiceClient{}
	u := &usermocks.UserServiceClient{}
	arc := &cmocks.AuthClient{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	u.On("Whoami", mock.Anything, mock.Anything).Return(&userspb.WhoamiResponse{
		User: &userspb.User{
			Id:    testUserId,
			Name:  testUserName,
			Email: testEmail,
		},
		OwnerOf: []*userspb.Organization{
			{
				Id:   testOrgId2,
				Name: testOrg1Name,
			},
		},
	}, nil)

	r := NewRouter(&Clients{
		Organization: client.NewOrgRegistryFromClient(o),
		User:         client.NewUserRegistryFromClient(u),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	u.AssertExpectations(t)
	assert.Contains(t, w.Body.String(), `"id":"`+testUserId+`"`)
	assert.Contains(t, w.Body.String(), `"name":"`+testOrg1Name+`"`)
}

func TestWhoami_NotFound(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/users/whoami/"+testUserId, nil)
	o := &orgmocks.OrgServiceClient{}
	u := &usermocks.UserServiceClient{}
	arc := &cmocks.AuthClient{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	u.On("Whoami", mock.Anything, mock.Anything).Return(nil, status.Error(codes.NotFound, "user not found"))

	r := NewRouter(&Clients{
		Organization: client.NewOrgRegistryFromClient(o),
		User:         client.NewUserRegistryFromClient(u),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPostUser_Success(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/users", strings.NewReader(testUserData))
	req.Header.Set("Content-Type", "application/json")
	o := &orgmocks.OrgServiceClient{}
	u := &usermocks.UserServiceClient{}
	arc := &cmocks.AuthClient{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	u.On("Add", mock.Anything, mock.Anything).Return(&userspb.AddResponse{
		User: &userspb.User{
			Id:     testUserUUID,
			Name:   testUserName,
			Email:  testEmail,
			Phone:  testPhone,
			AuthId: testAuthId,
		},
	}, nil)

	r := NewRouter(&Clients{
		Organization: client.NewOrgRegistryFromClient(o),
		User:         client.NewUserRegistryFromClient(u),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusCreated, w.Code)
	u.AssertExpectations(t)
	assert.Contains(t, w.Body.String(), `"name":"`+testUserName+`"`)
	assert.Contains(t, w.Body.String(), `"email":"`+testEmail+`"`)
}

func TestPostUser_Error(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/users", strings.NewReader(testUserData))
	req.Header.Set("Content-Type", "application/json")
	o := &orgmocks.OrgServiceClient{}
	u := &usermocks.UserServiceClient{}
	arc := &cmocks.AuthClient{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	u.On("Add", mock.Anything, mock.Anything).Return(nil, status.Error(codes.Internal, "internal error"))

	r := NewRouter(&Clients{
		Organization: client.NewOrgRegistryFromClient(o),
		User:         client.NewUserRegistryFromClient(u),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
