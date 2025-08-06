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
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/rest"
	"github.com/ukama/ukama/systems/init/bootstrap/pb/gen"
	"github.com/ukama/ukama/systems/init/node-gateway/mocks"
	"github.com/ukama/ukama/systems/init/node-gateway/pkg"
)

var defaultCors = rest.HttpConfig{
	Cors: cors.Config{
		AllowAllOrigins: true,
	},
}

var routerConfig = &RouterConfig{
	serverConf: &defaultCors,
	httpEndpoints: &pkg.HttpEndpoints{
		NodeMetrics: "localhost:8080",
	},
	auth: &config.Auth{
		AuthAppUrl:    "http://localhost:4455",
		AuthServerUrl: "http://localhost:4434",
		AuthAPIGW:     "http://localhost:8080",
	},
}

var testClientSet *Clients

func init() {
	gin.SetMode(gin.TestMode)
	testClientSet = NewClientsSet(&pkg.GrpcEndpoints{
		Timeout:   1 * time.Second,
		Bootstrap: "localhost:8080",
	})
}

func TestRouter_GetNodeCredentials_Success(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/nodes/test-node-id", nil)

	mockBootstrap := &mocks.BootstrapEP{}
	expectedResponse := &gen.GetNodeCredentialsResponse{
		Id:          "test-node-id",
		OrgName:     "test-org",
		Ip:          "0.0.0.0",
		Certificate: "test-certificate-data",
	}

	mockBootstrap.On("GetNodeCredentials", mock.Anything).Return(expectedResponse, nil)

	clients := &Clients{
		Bootstrap: mockBootstrap,
	}

	r := NewRouter(clients, routerConfig).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "test-node-id")
	assert.Contains(t, w.Body.String(), "test-org")
	assert.Contains(t, w.Body.String(), "0.0.0.0")
	assert.Contains(t, w.Body.String(), "test-certificate-data")
	mockBootstrap.AssertExpectations(t)
}
