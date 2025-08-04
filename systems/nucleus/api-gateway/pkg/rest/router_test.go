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
	req, _ := http.NewRequest("GET", "/v1/orgs/org-name", nil)
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
	const orgName = "org-name"
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/orgs/"+orgName, nil)
	o := &orgmocks.OrgServiceClient{}
	u := &usermocks.UserServiceClient{}
	arc := &cmocks.AuthClient{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	o.On("GetByName", mock.Anything, mock.Anything).Return(&orgpb.GetByNameResponse{
		Org: &orgpb.Organization{
			Name:  orgName,
			Owner: "owner",
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
	assert.Contains(t, w.Body.String(), `"name":"org-name"`)
}
