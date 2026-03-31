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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/node/api-gateway/pkg"
	"github.com/ukama/ukama/systems/node/api-gateway/pkg/client"

	cconfig "github.com/ukama/ukama/systems/common/config"
	cmmocks "github.com/ukama/ukama/systems/common/mocks"
	crest "github.com/ukama/ukama/systems/common/rest"
	cfgPb "github.com/ukama/ukama/systems/node/configurator/pb/gen"
	cmocks "github.com/ukama/ukama/systems/node/configurator/pb/gen/mocks"
	cpb "github.com/ukama/ukama/systems/node/controller/pb/gen"
	nmocks "github.com/ukama/ukama/systems/node/controller/pb/gen/mocks"
	spb "github.com/ukama/ukama/systems/node/software/pb/gen"
	smocks "github.com/ukama/ukama/systems/node/software/pb/gen/mocks"
	nspb "github.com/ukama/ukama/systems/node/state/pb/gen"
	stmocks "github.com/ukama/ukama/systems/node/state/pb/gen/mocks"
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

func TestRestartNodes(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	// Create a JSON payload with the necessary data.
	jsonPayload := `{"node_ids":["60285a2a-fe1d-4261-a868-5be480075b8f"]}`

	req, _ := http.NewRequest("POST", "/v1/controller/networks/456b2743-4831-4d8d-9fbe-830df7bd59d4/restart-nodes", strings.NewReader(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	arc := &cmmocks.AuthClient{}
	c := &nmocks.ControllerServiceClient{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	restartNodeReq := &cpb.RestartNodesRequest{
		NetworkId: "456b2743-4831-4d8d-9fbe-830df7bd59d4",
		NodeIds:   []string{"60285a2a-fe1d-4261-a868-5be480075b8f"},
	}

	c.On("RestartNodes", mock.Anything, restartNodeReq).Return(&cpb.RestartNodesResponse{}, nil)

	r := NewRouter(&Clients{
		Controller: client.NewControllerFromClient(c),
	}, routerConfig, arc.AuthenticateUser).f.Engine()
	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	c.AssertExpectations(t)
}

func TestRestartSite(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/controller/networks/0f37639d-3fd6-4741-b63b-9dd4f7ce55f0/sites/site-1/restart", nil)
	arc := &cmmocks.AuthClient{}
	c := &nmocks.ControllerServiceClient{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	RestartSiteRequest := &cpb.RestartSiteRequest{
		SiteId:    "site-1",
		NetworkId: "0f37639d-3fd6-4741-b63b-9dd4f7ce55f0",
	}

	c.On("RestartSite", mock.Anything, RestartSiteRequest).Return(&cpb.RestartSiteResponse{},
		nil)

	r := NewRouter(&Clients{
		Controller: client.NewControllerFromClient(c),
	}, routerConfig, arc.AuthenticateUser).f.Engine()
	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	c.AssertExpectations(t)
}

func TestPostConfigApplyVersionHandler(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	hash := "1c924398265578d35e2b16adca25dcc021923c89"
	req, _ := http.NewRequest("POST", "/v1/configurator/config/apply/"+hash, nil)
	arc := &cmmocks.AuthClient{}
	c := &nmocks.ControllerServiceClient{}
	cfg := &cmocks.ConfiguratorServiceClient{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	cfg.On("ApplyConfig", mock.Anything, mock.Anything).Return(&cfgPb.ApplyConfigResponse{},
		nil)

	r := NewRouter(&Clients{
		Controller:   client.NewControllerFromClient(c),
		Configurator: client.NewConfiguratorFromClient(cfg),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusAccepted, w.Code)
	cfg.AssertExpectations(t)
}

func TestGetPingNodeHandler(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	nodeId := "60285a2a-fe1d-4261-a868-5be480075b8f"
	req, _ := http.NewRequest("GET", "/v1/controller/nodes/"+nodeId+"/ping", nil)
	arc := &cmmocks.AuthClient{}
	c := &nmocks.ControllerServiceClient{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)
	c.On("PingNode", mock.Anything, &cpb.PingNodeRequest{NodeId: nodeId}).Return(&cpb.PingNodeResponse{}, nil)

	r := NewRouter(&Clients{
		Controller: client.NewControllerFromClient(c),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusAccepted, w.Code)
	c.AssertExpectations(t)
}

func TestPostToggleInternetSwitchHandler(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	siteId := "site-123"
	jsonPayload := `{"status": true, "port": 8080}`
	req, _ := http.NewRequest("POST", "/v1/controller/sites/"+siteId+"/toggle-internet-port", strings.NewReader(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	arc := &cmmocks.AuthClient{}
	c := &nmocks.ControllerServiceClient{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)
	c.On("ToggleInternetSwitch", mock.Anything, &cpb.ToggleInternetSwitchRequest{
		SiteId: siteId,
		Status: true,
		Port:   8080,
	}).Return(&cpb.ToggleInternetSwitchResponse{}, nil)

	r := NewRouter(&Clients{
		Controller: client.NewControllerFromClient(c),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	c.AssertExpectations(t)
}

func TestPostConfigEventHandler(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	body := `{"event": "push", "repo": "test-repo"}`
	req, _ := http.NewRequest("POST", "/v1/configurator/config", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	arc := &cmmocks.AuthClient{}
	cfg := &cmocks.ConfiguratorServiceClient{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)
	cfg.On("ConfigEvent", mock.Anything, &cfgPb.ConfigStoreEvent{Data: []byte(body)}).Return(&cfgPb.ConfigStoreEventResponse{}, nil)

	r := NewRouter(&Clients{
		Configurator: client.NewConfiguratorFromClient(cfg),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusAccepted, w.Code)
	cfg.AssertExpectations(t)
}

func TestGetListAppsHandler(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/software/apps", nil)
	arc := &cmmocks.AuthClient{}
	sw := &smocks.SoftwareServiceClient{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)
	sw.On("GetAppList", mock.Anything, &spb.GetAppListRequest{}).Return(&spb.GetAppListResponse{Apps: []*spb.App{}}, nil)

	r := NewRouter(&Clients{
		SoftwareManager: client.NewSoftwareManagerFromClient(sw),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	sw.AssertExpectations(t)
}

func TestGetListSoftwareHandler(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	nodeId := ukama.NewVirtualHomeNodeId().String()
	req, _ := http.NewRequest("GET", "/v1/software?node_id="+nodeId+"&app_name=ukama&status=up_to_date", nil)
	arc := &cmmocks.AuthClient{}
	sw := &smocks.SoftwareServiceClient{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)
	sw.On("GetSoftwareList", mock.Anything, mock.Anything).Return(&spb.GetSoftwareListResponse{Software: []*spb.Software{}}, nil)

	r := NewRouter(&Clients{
		SoftwareManager: client.NewSoftwareManagerFromClient(sw),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	sw.AssertExpectations(t)
}

func TestPostUpdateSoftwareHandler(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	nodeId := "60285a2a-fe1d-4261-a868-5be480075b8f"
	name := "ukama-node"
	tag := "v1.0.0"
	req, _ := http.NewRequest("POST", "/v1/software/update/"+name+"/"+tag+"/"+nodeId, nil)
	arc := &cmmocks.AuthClient{}
	sw := &smocks.SoftwareServiceClient{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)
	sw.On("UpdateSoftware", mock.Anything, &spb.UpdateSoftwareRequest{
		NodeId: nodeId,
		Name:   name,
		Tag:    tag,
	}).Return(&spb.UpdateSoftwareResponse{Message: "updated"}, nil)

	r := NewRouter(&Clients{
		SoftwareManager: client.NewSoftwareManagerFromClient(sw),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	sw.AssertExpectations(t)
}

func TestGetStatesHandler(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	nodeId := ukama.NewVirtualHomeNodeId().String()
	req, _ := http.NewRequest("POST", "/v1/state/"+nodeId, nil)
	arc := &cmmocks.AuthClient{}
	st := &stmocks.StateServiceClient{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)
	st.On("GetStates", mock.Anything, &nspb.GetStatesRequest{NodeId: nodeId}).Return(&nspb.GetStatesResponse{}, nil)

	r := NewRouter(&Clients{
		State: client.NewStateFromClient(st),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	st.AssertExpectations(t)
}

func TestGetStatesHistoryHandler(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	nodeId := ukama.NewVirtualHomeNodeId().String()
	req, _ := http.NewRequest("GET", "/v1/state/"+nodeId+"/history?page_size=10&page_number=1", nil)
	arc := &cmmocks.AuthClient{}
	st := &stmocks.StateServiceClient{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)
	st.On("GetStatesHistory", mock.Anything, &nspb.GetStatesHistoryRequest{
		NodeId:     nodeId,
		PageSize:   10,
		PageNumber: 1,
		StartTime:  "",
		EndTime:    "",
	}).Return(&nspb.GetStatesHistoryResponse{}, nil)

	r := NewRouter(&Clients{
		State: client.NewStateFromClient(st),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	st.AssertExpectations(t)
}

func TestEnforceStateTransitionHandler(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	nodeId := ukama.NewVirtualHomeNodeId().String()
	event := "activate"
	req, _ := http.NewRequest("POST", "/v1/state/"+nodeId+"/enforce/"+event, nil)
	arc := &cmmocks.AuthClient{}
	st := &stmocks.StateServiceClient{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)
	st.On("EnforceStateTransition", mock.Anything, &nspb.EnforceStateTransitionRequest{
		NodeId: nodeId,
		Event:  event,
	}).Return(&nspb.EnforceStateTransitionResponse{}, nil)

	r := NewRouter(&Clients{
		State: client.NewStateFromClient(st),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	st.AssertExpectations(t)
}

func TestGetRunningConfigVersionHandler(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	node := ukama.NewVirtualHomeNodeId().String()
	req, _ := http.NewRequest("GET", "/v1/configurator/config/node/"+node, nil)
	arc := &cmmocks.AuthClient{}
	c := &nmocks.ControllerServiceClient{}
	cfg := &cmocks.ConfiguratorServiceClient{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

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
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	if assert.Equal(t, http.StatusOK, w.Code) {
		assert.Contains(t, w.Body.String(), node)
	}
	c.AssertExpectations(t)
}

func TestPostToggleRfHandler(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	nodeId := "60285a2a-fe1d-4261-a868-5be480075b8f"
	jsonPayload := `{"state": "off"}`
	req, _ := http.NewRequest("POST", "/v1/controller/nodes/"+nodeId+"/toggle-radio", strings.NewReader(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	arc := &cmmocks.AuthClient{}
	c := &nmocks.ControllerServiceClient{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	c.On("ToggleRfSwitch", mock.Anything, &cpb.ToggleRfSwitchRequest{
		NodeId: nodeId,
		State:  "off",
	}).Return(&cpb.ToggleRfSwitchResponse{}, nil)

	r := NewRouter(&Clients{
		Controller: client.NewControllerFromClient(c),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	c.AssertExpectations(t)
}

func TestPostToggleNodeServiceHandler(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	nodeId := "uk-983794-hnode-78-7830"
	jsonPayload := `{"state": "on"}`
	req, _ := http.NewRequest("POST", "/v1/controller/nodes/"+nodeId+"/toggle-service", strings.NewReader(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	arc := &cmmocks.AuthClient{}
	c := &nmocks.ControllerServiceClient{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	c.On("ToggleNodeService", mock.Anything, &cpb.ToggleNodeServiceRequest{
		NodeId: nodeId,
		State:  "on",
	}).Return(&cpb.ToggleNodeServiceResponse{}, nil)

	r := NewRouter(&Clients{
		Controller: client.NewControllerFromClient(c),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	c.AssertExpectations(t)
}

func TestPostToggleNodeServiceHandlerInvalidState(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	nodeId := "uk-983794-hnode-78-7830"
	jsonPayload := `{"state": "invalid"}`
	req, _ := http.NewRequest("POST", "/v1/controller/nodes/"+nodeId+"/toggle-service", strings.NewReader(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	arc := &cmmocks.AuthClient{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	r := NewRouter(&Clients{
		Controller: client.NewControllerFromClient(&nmocks.ControllerServiceClient{}),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert - validation should reject invalid state
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
