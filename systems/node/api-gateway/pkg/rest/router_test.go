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
	spb "github.com/ukama/ukama/systems/node/software/pb/gen"
	smocks "github.com/ukama/ukama/systems/node/software/pb/gen/mocks"
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
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/controller/nodes/60285a2a-fe1d-4261-a868-5be480075b8f/restart", nil)
	arc := &providers.AuthRestClient{}
	c := &nmocks.ControllerServiceClient{}
	cfg := &cmocks.ConfiguratorServiceClient{}

	c.On("RestartNode", mock.Anything, mock.Anything).Return(&cpb.RestartNodeResponse{},
		nil)

	r := NewRouter(&Clients{
		Controller:   client.NewControllerFromClient(c),
		Configurator: client.NewConfiguratorFromClient(cfg),
	}, routerConfig, arc.MockAuthenticateUser).f.Engine()
	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	c.AssertExpectations(t)
}

func Test_RestarteNodes(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	// Create a JSON payload with the necessary data.
	jsonPayload := `{"node_ids":["60285a2a-fe1d-4261-a868-5be480075b8f"]}`

	req, _ := http.NewRequest("POST", "/v1/controller/networks/456b2743-4831-4d8d-9fbe-830df7bd59d4/restart-nodes", strings.NewReader(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	arc := &providers.AuthRestClient{}
	c := &nmocks.ControllerServiceClient{}

	restartNodeReq := &cpb.RestartNodesRequest{
		NetworkId: "456b2743-4831-4d8d-9fbe-830df7bd59d4",
		NodeIds:   []string{"60285a2a-fe1d-4261-a868-5be480075b8f"},
	}

	c.On("RestartNodes", mock.Anything, restartNodeReq).Return(&cpb.RestartNodesResponse{}, nil)

	r := NewRouter(&Clients{
		Controller: client.NewControllerFromClient(c),
	}, routerConfig, arc.MockAuthenticateUser).f.Engine()
	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	c.AssertExpectations(t)
}

func Test_SoftwareUpdate(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/software/update/space1/name1/tag1/uk-983794-hnode-78-7830", nil)
	arc := &providers.AuthRestClient{}
	c := &smocks.SoftwareServiceClient{}

	c.On("UpdateSoftware", mock.Anything, mock.Anything).Return(&spb.UpdateSoftwareResponse{},
		nil)

	r := NewRouter(&Clients{
		SoftwareManager: client.NewSoftwareManagerFromClient(c),
	}, routerConfig, arc.MockAuthenticateUser).f.Engine()
	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	c.AssertExpectations(t)
}

func Test_RestarteSite(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/controller/networks/0f37639d-3fd6-4741-b63b-9dd4f7ce55f0/sites/site-1/restart", nil)
	arc := &providers.AuthRestClient{}
	c := &nmocks.ControllerServiceClient{}

	RestartSiteRequest := &cpb.RestartSiteRequest{
		SiteId:    "site-1",
		NetworkId: "0f37639d-3fd6-4741-b63b-9dd4f7ce55f0",
	}

	c.On("RestartSite", mock.Anything, RestartSiteRequest).Return(&cpb.RestartSiteResponse{},
		nil)

	r := NewRouter(&Clients{
		Controller: client.NewControllerFromClient(c),
	}, routerConfig, arc.MockAuthenticateUser).f.Engine()
	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
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

func Test_postToggleRfHandler(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	nodeId := "60285a2a-fe1d-4261-a868-5be480075b8f"
	jsonPayload := `{"status": false}`
	req, _ := http.NewRequest("POST", "/v1/controller/nodes/"+nodeId+"/toggle-rf", strings.NewReader(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	arc := &providers.AuthRestClient{}
	c := &nmocks.ControllerServiceClient{}

	c.On("ToggleRfSwitch", mock.Anything, &cpb.ToggleRfSwitchRequest{
		NodeId: nodeId,
		Status: false,
	}).Return(&cpb.ToggleRfSwitchResponse{}, nil)

	r := NewRouter(&Clients{
		Controller: client.NewControllerFromClient(c),
	}, routerConfig, arc.MockAuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	c.AssertExpectations(t)
}
