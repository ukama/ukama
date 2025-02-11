/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package rest

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	cconfig "github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/providers"

	"github.com/tj/assert"
	crest "github.com/ukama/ukama/systems/common/rest"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/messaging/api-gateway/pkg"

	"github.com/ukama/ukama/systems/messaging/api-gateway/pkg/client"
	nnmocks "github.com/ukama/ukama/systems/messaging/nns/pb/gen/mocks"
	ngr "github.com/ukama/ukama/systems/node/node-gateway/pkg/rest"
)

const notifyApiEndpoint = "/v1/notify"

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
		Timeout: 1 * time.Second,
		Nns:     "0.0.0.0:9092",
	})
}
func TestPingRoute(t *testing.T) {
	arc := &providers.AuthRestClient{}
	nn := &nnmocks.NnsClient{}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)

	r := NewRouter(&Clients{
		n: client.NewNnsFromClient(nn),
	}, routerConfig, arc.MockAuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "pong")
}

var nodeId = ukama.NewVirtualHomeNodeId().String()
var nt = ngr.AddNotificationReq{
	NodeId:      nodeId,
	Severity:    "high",
	Type:        "event",
	ServiceName: "noded",
	Status:      8300,
	Time:        uint32(time.Now().Unix()),
	Details:     json.RawMessage(`{"reason":"testing","component":"router_test"}`),
}
