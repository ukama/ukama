/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ukama/ukama/systems/common/rest"

	bizpb "github.com/ukama/ukama/systems/analytics/business/pb/gen"
	colpb "github.com/ukama/ukama/systems/analytics/collector/pb/gen"
	custpb "github.com/ukama/ukama/systems/analytics/customer/pb/gen"
	netpb "github.com/ukama/ukama/systems/analytics/network/pb/gen"
	cconfig "github.com/ukama/ukama/systems/common/config"
)

var defaultCors = cors.Config{
	AllowAllOrigins: true,
}

var routerConfig = &RouterConfig{
	serverConf: &rest.HttpConfig{
		Cors: defaultCors,
	},
	auth: &cconfig.Auth{
		AuthAppUrl:    "http://localhost:4455",
		AuthServerUrl: "http://localhost:4434",
		AuthAPIGW:     "http://localhost:8080",
	},
}

func noAuth(c *gin.Context, s string) error {
	return nil
}

/* Stub clients implementing the narrow router interfaces. */

type stubBusiness struct {
	err error
}

func (s *stubBusiness) GetHome(networkId, period, from, to, tz string) (*bizpb.GetHomeResponse, error) {
	if s.err != nil {
		return nil, s.err
	}
	return &bizpb.GetHomeResponse{Kpis: []*bizpb.Kpi{{Key: "revenue_today", Value: 42}}}, nil
}

func (s *stubBusiness) GetSalesOverview(networkId, period, from, to, tz string) (*bizpb.GetSalesOverviewResponse, error) {
	return &bizpb.GetSalesOverviewResponse{}, s.err
}

func (s *stubBusiness) GetPackagePerformance(networkId, period, from, to, tz string, page, pageSize uint32) (*bizpb.GetPackagePerformanceResponse, error) {
	return &bizpb.GetPackagePerformanceResponse{}, s.err
}

func (s *stubBusiness) GetBillingSummary(period, from, to, tz string, page, pageSize uint32) (*bizpb.GetBillingSummaryResponse, error) {
	return &bizpb.GetBillingSummaryResponse{}, s.err
}

func (s *stubBusiness) GetSites(networkId, period, from, to, tz string, page, pageSize uint32) (*bizpb.GetSitesResponse, error) {
	return &bizpb.GetSitesResponse{}, s.err
}

func (s *stubBusiness) GetSite(siteId, period, from, to, tz string) (*bizpb.GetSiteResponse, error) {
	return &bizpb.GetSiteResponse{}, s.err
}

func (s *stubBusiness) GetInventoryReadiness(networkId string) (*bizpb.GetInventoryReadinessResponse, error) {
	return &bizpb.GetInventoryReadinessResponse{}, s.err
}

type stubCustomer struct {
	err error
}

func (s *stubCustomer) GetOverview(networkId, period, from, to, tz string) (*custpb.GetOverviewResponse, error) {
	if s.err != nil {
		return nil, s.err
	}
	return &custpb.GetOverviewResponse{Kpis: []*custpb.Kpi{{Key: "total", Value: 10}}}, nil
}

func (s *stubCustomer) List(networkId, siteId, status string, page, pageSize uint32) (*custpb.ListResponse, error) {
	return &custpb.ListResponse{}, s.err
}

func (s *stubCustomer) Search(query, networkId string, page, pageSize uint32) (*custpb.SearchResponse, error) {
	return &custpb.SearchResponse{}, s.err
}

func (s *stubCustomer) Get(customerId, period, from, to, tz string) (*custpb.GetResponse, error) {
	return &custpb.GetResponse{}, s.err
}

func (s *stubCustomer) GetSupport(customerId string) (*custpb.GetSupportResponse, error) {
	return &custpb.GetSupportResponse{}, s.err
}

func (s *stubCustomer) GetSims(networkId, status string, page, pageSize uint32) (*custpb.GetSimsResponse, error) {
	return &custpb.GetSimsResponse{}, s.err
}

func (s *stubCustomer) GetSimPool(networkId string) (*custpb.GetSimPoolResponse, error) {
	return &custpb.GetSimPoolResponse{}, s.err
}

type stubNetwork struct {
	err error
}

func (s *stubNetwork) GetOverview(networkId, period, from, to, tz string) (*netpb.GetOverviewResponse, error) {
	if s.err != nil {
		return nil, s.err
	}
	return &netpb.GetOverviewResponse{NetworkStatus: "healthy"}, nil
}

func (s *stubNetwork) GetTopology(networkId string) (*netpb.GetTopologyResponse, error) {
	return &netpb.GetTopologyResponse{}, s.err
}

func (s *stubNetwork) GetSites(networkId, status, period, from, to, tz string, page, pageSize uint32) (*netpb.GetSitesResponse, error) {
	return &netpb.GetSitesResponse{}, s.err
}

func (s *stubNetwork) GetSite(siteId, period, from, to, tz string) (*netpb.GetSiteResponse, error) {
	return &netpb.GetSiteResponse{}, s.err
}

func (s *stubNetwork) GetNodes(networkId, siteId, status string, page, pageSize uint32) (*netpb.GetNodesResponse, error) {
	return &netpb.GetNodesResponse{}, s.err
}

func (s *stubNetwork) GetNode(nodeId, period, from, to, tz string) (*netpb.GetNodeResponse, error) {
	return &netpb.GetNodeResponse{}, s.err
}

func (s *stubNetwork) GetNodePool(networkId string) (*netpb.GetNodePoolResponse, error) {
	return &netpb.GetNodePoolResponse{}, s.err
}

func (s *stubNetwork) GetRadio(networkId, siteId, nodeId, period, from, to, tz string) (*netpb.GetRadioResponse, error) {
	return &netpb.GetRadioResponse{}, s.err
}

func (s *stubNetwork) GetBackhaul(networkId, siteId, period, from, to, tz string) (*netpb.GetBackhaulResponse, error) {
	return &netpb.GetBackhaulResponse{}, s.err
}

func (s *stubNetwork) GetPower(networkId, siteId, period, from, to, tz string) (*netpb.GetPowerResponse, error) {
	return &netpb.GetPowerResponse{}, s.err
}

func (s *stubNetwork) GetAlarms(networkId, siteId, severity, state string, page, pageSize uint32) (*netpb.GetAlarmsResponse, error) {
	return &netpb.GetAlarmsResponse{}, s.err
}

func (s *stubNetwork) GetMetrics(networkId, siteId, nodeId, metric, period, from, to, tz string) (*netpb.GetMetricsResponse, error) {
	return &netpb.GetMetricsResponse{}, s.err
}

func (s *stubNetwork) GetEvents(networkId, siteId, nodeId, period, from, to, tz string, page, pageSize uint32) (*netpb.GetEventsResponse, error) {
	return &netpb.GetEventsResponse{}, s.err
}

func (s *stubNetwork) SupportSearch(query, networkId string) (*netpb.SupportSearchResponse, error) {
	return &netpb.SupportSearchResponse{}, s.err
}

type stubCollector struct {
	err error
}

func (s *stubCollector) Refresh(source string) (*colpb.RefreshResponse, error) {
	if s.err != nil {
		return nil, s.err
	}
	return &colpb.RefreshResponse{States: []*colpb.SourceState{{Source: source, Status: "running"}}}, nil
}

func (s *stubCollector) GetRefreshState() (*colpb.GetRefreshStateResponse, error) {
	return &colpb.GetRefreshStateResponse{}, s.err
}

func (s *stubCollector) RebuildRollups(family, from, to string) (*colpb.RebuildRollupsResponse, error) {
	return &colpb.RebuildRollupsResponse{}, s.err
}

func (s *stubCollector) SeedDemo(sites, nodes, customers, days uint32) (*colpb.SeedDemoResponse, error) {
	return &colpb.SeedDemoResponse{}, s.err
}

func newTestRouter(clients *Clients) *gin.Engine {
	gin.SetMode(gin.TestMode)
	return NewRouter(clients, routerConfig, noAuth).f.Engine()
}

func defaultClients() *Clients {
	return &Clients{
		Business:  &stubBusiness{},
		Customer:  &stubCustomer{},
		Network:   &stubNetwork{},
		Collector: &stubCollector{},
	}
}

func TestPingRoute(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)

	newTestRouter(defaultClients()).ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "pong")
}

func TestGetBusinessHome(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/analytics/business/home?period=week", nil)

	newTestRouter(defaultClients()).ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "revenue_today")
}

func TestGetCustomerOverview(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/analytics/customers/overview", nil)

	newTestRouter(defaultClients()).ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "total")
}

func TestGetNetworkOverview(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/analytics/network/overview", nil)

	newTestRouter(defaultClients()).ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "healthy")
}

func TestPostCollectorRefresh(t *testing.T) {
	w := httptest.NewRecorder()
	body := strings.NewReader(`{"source": "all"}`)
	req, _ := http.NewRequest("POST", "/v1/analytics/collector/refresh", body)
	req.Header.Set("Content-Type", "application/json")

	newTestRouter(defaultClients()).ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "running")
}

func TestPostCollectorRefresh_MissingSource(t *testing.T) {
	// source is required; missing it must yield a 400.
	w := httptest.NewRecorder()
	body := strings.NewReader(`{}`)
	req, _ := http.NewRequest("POST", "/v1/analytics/collector/refresh", body)
	req.Header.Set("Content-Type", "application/json")

	newTestRouter(defaultClients()).ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSearchCustomers_MissingQuery(t *testing.T) {
	// q is required; missing it must yield a 400.
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/analytics/customers/search", nil)

	newTestRouter(defaultClients()).ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetBusinessHome_ClientError(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/analytics/business/home", nil)

	clients := defaultClients()
	clients.Business = &stubBusiness{err: status.Error(codes.Internal, "boom")}

	newTestRouter(clients).ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetCustomer_ById(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/analytics/customers/3f1f3d6e-1111-4222-8333-944444444444", nil)

	newTestRouter(defaultClients()).ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
