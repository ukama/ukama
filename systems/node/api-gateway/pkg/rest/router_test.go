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

	"github.com/ukama/ukama/systems/node/api-gateway/pkg"
	"github.com/ukama/ukama/systems/node/api-gateway/pkg/client"

	cconfig "github.com/ukama/ukama/systems/common/config"
	cmmocks "github.com/ukama/ukama/systems/common/mocks"
	crest "github.com/ukama/ukama/systems/common/rest"
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
		Software:     "0.0.0.0:9091",
	})
}
func TestPingRoute(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	arc := &cmmocks.AuthClient{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	req, _ := http.NewRequest("GET", "/ping", nil)
	r := NewRouter(testClientSet, routerConfig, arc.AuthenticateUser).f.Engine()
	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "pong")
}

func TestRestartNode(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/controller/nodes/60285a2a-fe1d-4261-a868-5be480075b8f/restart", nil)
	arc := &cmmocks.AuthClient{}
	c := &nmocks.ControllerServiceClient{}
	cfg := &cmocks.ConfiguratorServiceClient{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	c.On("RestartNode", mock.Anything, mock.Anything).Return(&cpb.RestartNodeResponse{},
		nil)

	r := NewRouter(&Clients{
		Controller:   client.NewControllerFromClient(c),
		Configurator: client.NewConfiguratorFromClient(cfg),
	}, routerConfig, arc.AuthenticateUser).f.Engine()
	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	c.AssertExpectations(t)
}
