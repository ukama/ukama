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
	"github.com/stretchr/testify/mock"
	"github.com/tj/assert"

	"github.com/ukama/ukama/systems/messaging/api-gateway/pkg"
	"github.com/ukama/ukama/systems/messaging/api-gateway/pkg/client"

	cconfig "github.com/ukama/ukama/systems/common/config"
	cmocks "github.com/ukama/ukama/systems/common/mocks"
	crest "github.com/ukama/ukama/systems/common/rest"
	nnmocks "github.com/ukama/ukama/systems/messaging/nns/pb/gen/mocks"
)

var (
	testClientSet *Clients

	defaultCors = cors.Config{
		AllowAllOrigins: true,
	}

	routerConfig = &RouterConfig{
		serverConf: &crest.HttpConfig{
			Cors: defaultCors,
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
	testClientSet = NewClientsSet(&pkg.GrpcEndpoints{
		Timeout: 1 * time.Second,
		Nns:     "0.0.0.0:9092",
	})
}

func TestPingRoute(t *testing.T) {
	arc := &cmocks.AuthClient{}
	nn := &nnmocks.NnsClient{}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	r := NewRouter(&Clients{
		n: client.NewNnsFromClient(nn),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "pong")
}
