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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ukama/ukama/systems/node/api-gateway/pkg"
	"github.com/ukama/ukama/systems/node/api-gateway/pkg/client"

	cconfig "github.com/ukama/ukama/systems/common/config"
	cmmocks "github.com/ukama/ukama/systems/common/mocks"
	ukamaPb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
	crest "github.com/ukama/ukama/systems/common/rest"
	cmocks "github.com/ukama/ukama/systems/node/configurator/pb/gen/mocks"
	cpb "github.com/ukama/ukama/systems/node/controller/pb/gen"
	nmocks "github.com/ukama/ukama/systems/node/controller/pb/gen/mocks"
	hpb "github.com/ukama/ukama/systems/node/health/pb/gen"
	hmocks "github.com/ukama/ukama/systems/node/health/pb/gen/mocks"
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

const healthNodeId = "uk-sa2341-hnode-v0-a1a0"

func newHealthRouter(h *hmocks.HealthServiceClient) *gin.Engine {
	arc := &cmmocks.AuthClient{}
	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	return NewRouter(&Clients{
		Health: client.NewHealthFromClient(h),
	}, routerConfig, arc.AuthenticateUser).f.Engine()
}

func TestGetNodeHealthReports(t *testing.T) {
	t.Run("LatestReport", func(t *testing.T) {
		h := &hmocks.HealthServiceClient{}
		h.On("ListReports", mock.Anything, &hpb.ListReportsRequest{
			NodeId:     healthNodeId,
			ReportId:   "",
			ReportedAt: 1779534357,
			Timeframe:  ukamaPb.FilterTimeframesType_LATEST,
		}).Return(&hpb.ListReportsResponse{
			Reports: []*hpb.HealthReport{{Id: "r1", NodeId: healthNodeId}},
		}, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/health/nodes/"+strings.ToUpper(healthNodeId)+"/reports?timeframe=latest&reportedAt=1779534357", nil)
		newHealthRouter(h).ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "r1")
		h.AssertExpectations(t)
	})

	t.Run("DefaultsToAllTimeframe", func(t *testing.T) {
		h := &hmocks.HealthServiceClient{}
		h.On("ListReports", mock.Anything, mock.MatchedBy(func(r *hpb.ListReportsRequest) bool {
			return r.NodeId == healthNodeId && r.Timeframe == ukamaPb.FilterTimeframesType_ALL
		})).Return(&hpb.ListReportsResponse{}, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/health/nodes/"+healthNodeId+"/reports", nil)
		newHealthRouter(h).ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		h.AssertExpectations(t)
	})

	t.Run("InvalidTimeframe", func(t *testing.T) {
		h := &hmocks.HealthServiceClient{}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/health/nodes/"+healthNodeId+"/reports?timeframe=weekly", nil)
		newHealthRouter(h).ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		h.AssertNotCalled(t, "ListReports", mock.Anything, mock.Anything)
	})

	t.Run("InvalidNodeId", func(t *testing.T) {
		h := &hmocks.HealthServiceClient{}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/health/nodes/not-a-node/reports", nil)
		newHealthRouter(h).ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		h.AssertNotCalled(t, "ListReports", mock.Anything, mock.Anything)
	})
}

func TestGetNodeApps(t *testing.T) {
	t.Run("FilterByAppName", func(t *testing.T) {
		h := &hmocks.HealthServiceClient{}
		h.On("ListApps", mock.Anything, &hpb.ListAppsRequest{
			NodeId:  healthNodeId,
			AppName: "noded",
		}).Return(&hpb.ListAppsResponse{
			Apps: []*hpb.App{{Name: "noded", Status: "running"}},
		}, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/health/nodes/"+healthNodeId+"/apps?appName=noded", nil)
		newHealthRouter(h).ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "noded")
		h.AssertExpectations(t)
	})

	t.Run("AppNotFound", func(t *testing.T) {
		h := &hmocks.HealthServiceClient{}
		h.On("ListApps", mock.Anything, mock.Anything).
			Return(nil, status.Error(codes.NotFound, `app "missing" not found`))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/health/nodes/"+healthNodeId+"/apps?appName=missing", nil)
		newHealthRouter(h).ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		h.AssertExpectations(t)
	})
}

func TestGetNodeInterfaces(t *testing.T) {
	h := &hmocks.HealthServiceClient{}
	h.On("ListInterfaces", mock.Anything, &hpb.ListInterfacesRequest{
		NodeId: healthNodeId,
	}).Return(&hpb.ListInterfacesResponse{
		Interfaces: &hpb.Interface{
			Cellular: &hpb.CellularInterface{Available: true, Service: "on"},
		},
	}, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/health/nodes/"+healthNodeId+"/interfaces", nil)
	newHealthRouter(h).ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "cellular")
	h.AssertExpectations(t)
}
