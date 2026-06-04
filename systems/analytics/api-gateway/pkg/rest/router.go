/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package rest

import (
	"fmt"
	"net/http"

	"github.com/ukama/ukama/systems/analytics/api-gateway/cmd/version"
	"github.com/ukama/ukama/systems/analytics/api-gateway/pkg"
	"github.com/ukama/ukama/systems/analytics/api-gateway/pkg/client"
	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/rest"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/wI2L/fizz"
	"github.com/wI2L/fizz/openapi"

	log "github.com/sirupsen/logrus"
	bizpb "github.com/ukama/ukama/systems/analytics/business/pb/gen"
	colpb "github.com/ukama/ukama/systems/analytics/collector/pb/gen"
	custpb "github.com/ukama/ukama/systems/analytics/customer/pb/gen"
	netpb "github.com/ukama/ukama/systems/analytics/network/pb/gen"
)

type Router struct {
	f       *fizz.Fizz
	clients *Clients
	config  *RouterConfig
}

type RouterConfig struct {
	metricsConfig config.Metrics
	debugMode     bool
	serverConf    *rest.HttpConfig
	auth          *config.Auth
}

type Clients struct {
	Business  business
	Customer  customer
	Network   network
	Collector collector
}

type business interface {
	GetHome(networkId, period, from, to, tz string) (*bizpb.GetHomeResponse, error)
	GetSalesOverview(networkId, period, from, to, tz string) (*bizpb.GetSalesOverviewResponse, error)
	GetPackagePerformance(networkId, period, from, to, tz string, page, pageSize uint32) (*bizpb.GetPackagePerformanceResponse, error)
	GetBillingSummary(period, from, to, tz string, page, pageSize uint32) (*bizpb.GetBillingSummaryResponse, error)
	GetSites(networkId, period, from, to, tz string, page, pageSize uint32) (*bizpb.GetSitesResponse, error)
	GetSite(siteId, period, from, to, tz string) (*bizpb.GetSiteResponse, error)
	GetInventoryReadiness(networkId string) (*bizpb.GetInventoryReadinessResponse, error)
}

type customer interface {
	GetOverview(networkId, period, from, to, tz string) (*custpb.GetOverviewResponse, error)
	List(networkId, siteId, status string, page, pageSize uint32) (*custpb.ListResponse, error)
	Search(query, networkId string, page, pageSize uint32) (*custpb.SearchResponse, error)
	Get(customerId, period, from, to, tz string) (*custpb.GetResponse, error)
	GetSupport(customerId string) (*custpb.GetSupportResponse, error)
	GetSims(networkId, status string, page, pageSize uint32) (*custpb.GetSimsResponse, error)
	GetSimPool(networkId string) (*custpb.GetSimPoolResponse, error)
}

type network interface {
	GetOverview(networkId, period, from, to, tz string) (*netpb.GetOverviewResponse, error)
	GetTopology(networkId string) (*netpb.GetTopologyResponse, error)
	GetSites(networkId, status, period, from, to, tz string, page, pageSize uint32) (*netpb.GetSitesResponse, error)
	GetSite(siteId, period, from, to, tz string) (*netpb.GetSiteResponse, error)
	GetNodes(networkId, siteId, status string, page, pageSize uint32) (*netpb.GetNodesResponse, error)
	GetNode(nodeId, period, from, to, tz string) (*netpb.GetNodeResponse, error)
	GetNodePool(networkId string) (*netpb.GetNodePoolResponse, error)
	GetRadio(networkId, siteId, nodeId, period, from, to, tz string) (*netpb.GetRadioResponse, error)
	GetBackhaul(networkId, siteId, period, from, to, tz string) (*netpb.GetBackhaulResponse, error)
	GetPower(networkId, siteId, period, from, to, tz string) (*netpb.GetPowerResponse, error)
	GetAlarms(networkId, siteId, severity, state string, page, pageSize uint32) (*netpb.GetAlarmsResponse, error)
	GetMetrics(networkId, siteId, nodeId, metric, period, from, to, tz string) (*netpb.GetMetricsResponse, error)
	GetEvents(networkId, siteId, nodeId, period, from, to, tz string, page, pageSize uint32) (*netpb.GetEventsResponse, error)
	SupportSearch(query, networkId string) (*netpb.SupportSearchResponse, error)
}

type collector interface {
	Refresh(source string) (*colpb.RefreshResponse, error)
	GetRefreshState() (*colpb.GetRefreshStateResponse, error)
	RebuildRollups(family, from, to string) (*colpb.RebuildRollupsResponse, error)
	SeedDemo(sites, nodes, customers, days uint32) (*colpb.SeedDemoResponse, error)
}

func NewClientsSet(endpoints *pkg.GrpcEndpoints) *Clients {
	c := &Clients{}
	c.Business = client.NewBusinessAnalytics(endpoints.Business, endpoints.Timeout)
	c.Customer = client.NewCustomerAnalytics(endpoints.Customer, endpoints.Timeout)
	c.Network = client.NewNetworkAnalytics(endpoints.Network, endpoints.Timeout)
	c.Collector = client.NewCollectorAnalytics(endpoints.Collector, endpoints.Timeout)

	return c
}

func NewRouter(clients *Clients, config *RouterConfig, authfunc func(*gin.Context, string) error) *Router {
	r := &Router{
		clients: clients,
		config:  config,
	}

	if !config.debugMode {
		gin.SetMode(gin.ReleaseMode)
	}

	r.init(authfunc)
	return r
}

func NewRouterConfig(svcConf *pkg.Config) *RouterConfig {
	return &RouterConfig{
		metricsConfig: svcConf.Metrics,
		serverConf:    &svcConf.Server,
		debugMode:     svcConf.DebugMode,
		auth:          svcConf.Auth,
	}
}

func (rt *Router) Run() {
	log.Info("Listening on port ", rt.config.serverConf.Port)
	err := rt.f.Engine().Run(fmt.Sprint(":", rt.config.serverConf.Port))
	if err != nil {
		panic(err)
	}
}

func (r *Router) init(f func(*gin.Context, string) error) {
	r.f = rest.NewFizzRouter(r.config.serverConf, pkg.SystemName, version.Version, r.config.debugMode, r.config.auth.AuthAppUrl+"?redirect=true")
	auth := r.f.Group("/v1/analytics", "Analytics API gateway", "Analytics system version v1", func(ctx *gin.Context) {
		if r.config.auth.BypassAuthMode {
			log.Info("Bypassing auth")
			return
		}
		s := fmt.Sprintf("%s, %s, %s", pkg.SystemName, ctx.Request.Method, ctx.Request.URL.Path)
		ctx.Request.Header.Set("Meta", s)
		err := f(ctx, r.config.auth.AuthAPIGW)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
		}
	})
	auth.Use()
	{
		// Business routes
		const biz = "/business"
		business := auth.Group(biz, "Business", "Business analytics")
		business.GET("/home", formatDoc("Business Home", "KPIs, sites, top packages and recent activity for the home dashboard"), tonic.Handler(r.getBusinessHomeHandler, http.StatusOK))
		business.GET("/sales/overview", formatDoc("Sales Overview", "Sales KPIs, revenue trend and breakdowns"), tonic.Handler(r.getSalesOverviewHandler, http.StatusOK))
		business.GET("/sales/packages", formatDoc("Package Performance", "Per-package sales performance"), tonic.Handler(r.getPackagePerformanceHandler, http.StatusOK))
		business.GET("/packages", append(formatDoc("Package Performance", "Per-package sales performance (alias)"), fizz.ID("getBusinessPackagesHandler")), tonic.Handler(r.getPackagePerformanceHandler, http.StatusOK))
		business.GET("/billing", formatDoc("Billing Summary", "Billing KPIs and invoices"), tonic.Handler(r.getBillingSummaryHandler, http.StatusOK))
		business.GET("/sites", formatDoc("Business Sites", "Per-site business metrics"), tonic.Handler(r.getBusinessSitesHandler, http.StatusOK))
		business.GET("/sites/:site_id", formatDoc("Business Site", "Business metrics for a single site"), tonic.Handler(r.getBusinessSiteHandler, http.StatusOK))
		business.GET("/inventory", formatDoc("Inventory Readiness", "SIM and node inventory readiness KPIs"), tonic.Handler(r.getInventoryReadinessHandler, http.StatusOK))

		// Customer routes
		// NOTE: gin resolves static routes (/overview, /list, ...) with priority
		// over the :customer_id param sibling, so static routes are registered
		// first and customer ids matching a static segment are not supported.
		const cust = "/customers"
		customers := auth.Group(cust, "Customers", "Customer analytics")
		customers.GET("/overview", formatDoc("Customer Overview", "Customer KPI overview"), tonic.Handler(r.getCustomerOverviewHandler, http.StatusOK))
		customers.GET("/list", formatDoc("List Customers", "List customers with filters"), tonic.Handler(r.listCustomersHandler, http.StatusOK))
		customers.GET("/search", formatDoc("Search Customers", "Search customers by name, email, iccid"), tonic.Handler(r.searchCustomersHandler, http.StatusOK))
		customers.GET("/sims", formatDoc("Get Sims", "List sims with filters"), tonic.Handler(r.getSimsHandler, http.StatusOK))
		customers.GET("/sim-pool", formatDoc("Sim Pool", "Sim pool KPIs and batches"), tonic.Handler(r.getSimPoolHandler, http.StatusOK))
		customers.GET("/:customer_id", formatDoc("Get Customer", "Customer detail with KPIs and package history"), tonic.Handler(r.getCustomerHandler, http.StatusOK))
		customers.GET("/:customer_id/support", formatDoc("Customer Support View", "Support diagnosis for a customer"), tonic.Handler(r.getCustomerSupportHandler, http.StatusOK))

		// Network routes
		const net = "/network"
		network := auth.Group(net, "Network", "Network analytics")
		network.GET("/overview", formatDoc("Network Overview", "Network health status and KPIs"), tonic.Handler(r.getNetworkOverviewHandler, http.StatusOK))
		network.GET("/topology", formatDoc("Network Topology", "Sites and nodes topology"), tonic.Handler(r.getTopologyHandler, http.StatusOK))
		network.GET("/sites", formatDoc("Network Sites", "Per-site network health"), tonic.Handler(r.getNetworkSitesHandler, http.StatusOK))
		network.GET("/sites/:site_id", formatDoc("Network Site", "Network health for a single site"), tonic.Handler(r.getNetworkSiteHandler, http.StatusOK))
		network.GET("/nodes", formatDoc("Network Nodes", "Per-node network health"), tonic.Handler(r.getNetworkNodesHandler, http.StatusOK))
		network.GET("/nodes/:node_id", formatDoc("Network Node", "Network health for a single node"), tonic.Handler(r.getNetworkNodeHandler, http.StatusOK))
		network.GET("/node-pool", formatDoc("Node Pool", "Node inventory pool"), tonic.Handler(r.getNodePoolHandler, http.StatusOK))
		network.GET("/radio", formatDoc("Radio", "Radio KPIs and series"), tonic.Handler(r.getRadioHandler, http.StatusOK))
		network.GET("/backhaul", formatDoc("Backhaul", "Backhaul KPIs and series"), tonic.Handler(r.getBackhaulHandler, http.StatusOK))
		network.GET("/power", formatDoc("Power", "Power KPIs and series"), tonic.Handler(r.getPowerHandler, http.StatusOK))
		network.GET("/alarms", formatDoc("Alarms", "Alarms with filters"), tonic.Handler(r.getAlarmsHandler, http.StatusOK))
		network.GET("/metrics", formatDoc("Metrics", "Raw metric series"), tonic.Handler(r.getMetricsHandler, http.StatusOK))
		network.GET("/events", formatDoc("Events", "Network events"), tonic.Handler(r.getEventsHandler, http.StatusOK))
		network.GET("/support/search", formatDoc("Support Search", "Search sites and nodes for support diagnosis"), tonic.Handler(r.supportSearchHandler, http.StatusOK))

		// Collector routes
		const col = "/collector"
		collector := auth.Group(col, "Collector", "Collector operations")
		collector.POST("/refresh", formatDoc("Refresh", "Trigger a source refresh"), tonic.Handler(r.postRefreshHandler, http.StatusOK))
		collector.GET("/state", formatDoc("Refresh State", "Source refresh and rollup state"), tonic.Handler(r.getRefreshStateHandler, http.StatusOK))
		collector.POST("/rollups/rebuild", formatDoc("Rebuild Rollups", "Rebuild rollups for a family"), tonic.Handler(r.postRebuildRollupsHandler, http.StatusOK))

		// SeedDemo is intentionally only exposed when the gateway runs in
		// debug mode; it populates demo data and must never be reachable
		// in production deployments.
		if r.config.debugMode {
			collector.POST("/seed-demo", formatDoc("Seed Demo", "Seed demo data (debug only)"), tonic.Handler(r.postSeedDemoHandler, http.StatusOK))
		}
	}
}

/* Business handlers */

func (r *Router) getBusinessHomeHandler(c *gin.Context, req *GetBusinessHomeRequest) (*bizpb.GetHomeResponse, error) {
	return r.clients.Business.GetHome(req.NetworkId, req.Period, req.From, req.To, req.Timezone)
}

func (r *Router) getSalesOverviewHandler(c *gin.Context, req *GetSalesOverviewRequest) (*bizpb.GetSalesOverviewResponse, error) {
	return r.clients.Business.GetSalesOverview(req.NetworkId, req.Period, req.From, req.To, req.Timezone)
}

func (r *Router) getPackagePerformanceHandler(c *gin.Context, req *GetPackagePerformanceRequest) (*bizpb.GetPackagePerformanceResponse, error) {
	return r.clients.Business.GetPackagePerformance(req.NetworkId, req.Period, req.From, req.To, req.Timezone, req.Page, req.PageSize)
}

func (r *Router) getBillingSummaryHandler(c *gin.Context, req *GetBillingSummaryRequest) (*bizpb.GetBillingSummaryResponse, error) {
	return r.clients.Business.GetBillingSummary(req.Period, req.From, req.To, req.Timezone, req.Page, req.PageSize)
}

func (r *Router) getBusinessSitesHandler(c *gin.Context, req *GetBusinessSitesRequest) (*bizpb.GetSitesResponse, error) {
	return r.clients.Business.GetSites(req.NetworkId, req.Period, req.From, req.To, req.Timezone, req.Page, req.PageSize)
}

func (r *Router) getBusinessSiteHandler(c *gin.Context, req *GetBusinessSiteRequest) (*bizpb.GetSiteResponse, error) {
	return r.clients.Business.GetSite(req.SiteId, req.Period, req.From, req.To, req.Timezone)
}

func (r *Router) getInventoryReadinessHandler(c *gin.Context, req *GetInventoryReadinessRequest) (*bizpb.GetInventoryReadinessResponse, error) {
	return r.clients.Business.GetInventoryReadiness(req.NetworkId)
}

/* Customer handlers */

func (r *Router) getCustomerOverviewHandler(c *gin.Context, req *GetCustomerOverviewRequest) (*custpb.GetOverviewResponse, error) {
	return r.clients.Customer.GetOverview(req.NetworkId, req.Period, req.From, req.To, req.Timezone)
}

func (r *Router) listCustomersHandler(c *gin.Context, req *ListCustomersRequest) (*custpb.ListResponse, error) {
	return r.clients.Customer.List(req.NetworkId, req.SiteId, req.Status, req.Page, req.PageSize)
}

func (r *Router) searchCustomersHandler(c *gin.Context, req *SearchCustomersRequest) (*custpb.SearchResponse, error) {
	return r.clients.Customer.Search(req.Query, req.NetworkId, req.Page, req.PageSize)
}

func (r *Router) getCustomerHandler(c *gin.Context, req *GetCustomerRequest) (*custpb.GetResponse, error) {
	return r.clients.Customer.Get(req.CustomerId, req.Period, req.From, req.To, req.Timezone)
}

func (r *Router) getCustomerSupportHandler(c *gin.Context, req *GetCustomerSupportRequest) (*custpb.GetSupportResponse, error) {
	return r.clients.Customer.GetSupport(req.CustomerId)
}

func (r *Router) getSimsHandler(c *gin.Context, req *GetSimsRequest) (*custpb.GetSimsResponse, error) {
	return r.clients.Customer.GetSims(req.NetworkId, req.Status, req.Page, req.PageSize)
}

func (r *Router) getSimPoolHandler(c *gin.Context, req *GetSimPoolRequest) (*custpb.GetSimPoolResponse, error) {
	return r.clients.Customer.GetSimPool(req.NetworkId)
}

/* Network handlers */

func (r *Router) getNetworkOverviewHandler(c *gin.Context, req *GetNetworkOverviewRequest) (*netpb.GetOverviewResponse, error) {
	return r.clients.Network.GetOverview(req.NetworkId, req.Period, req.From, req.To, req.Timezone)
}

func (r *Router) getTopologyHandler(c *gin.Context, req *GetTopologyRequest) (*netpb.GetTopologyResponse, error) {
	return r.clients.Network.GetTopology(req.NetworkId)
}

func (r *Router) getNetworkSitesHandler(c *gin.Context, req *GetNetworkSitesRequest) (*netpb.GetSitesResponse, error) {
	return r.clients.Network.GetSites(req.NetworkId, req.Status, req.Period, req.From, req.To, req.Timezone, req.Page, req.PageSize)
}

func (r *Router) getNetworkSiteHandler(c *gin.Context, req *GetNetworkSiteRequest) (*netpb.GetSiteResponse, error) {
	return r.clients.Network.GetSite(req.SiteId, req.Period, req.From, req.To, req.Timezone)
}

func (r *Router) getNetworkNodesHandler(c *gin.Context, req *GetNetworkNodesRequest) (*netpb.GetNodesResponse, error) {
	return r.clients.Network.GetNodes(req.NetworkId, req.SiteId, req.Status, req.Page, req.PageSize)
}

func (r *Router) getNetworkNodeHandler(c *gin.Context, req *GetNetworkNodeRequest) (*netpb.GetNodeResponse, error) {
	return r.clients.Network.GetNode(req.NodeId, req.Period, req.From, req.To, req.Timezone)
}

func (r *Router) getNodePoolHandler(c *gin.Context, req *GetNodePoolRequest) (*netpb.GetNodePoolResponse, error) {
	return r.clients.Network.GetNodePool(req.NetworkId)
}

func (r *Router) getRadioHandler(c *gin.Context, req *GetRadioRequest) (*netpb.GetRadioResponse, error) {
	return r.clients.Network.GetRadio(req.NetworkId, req.SiteId, req.NodeId, req.Period, req.From, req.To, req.Timezone)
}

func (r *Router) getBackhaulHandler(c *gin.Context, req *GetBackhaulRequest) (*netpb.GetBackhaulResponse, error) {
	return r.clients.Network.GetBackhaul(req.NetworkId, req.SiteId, req.Period, req.From, req.To, req.Timezone)
}

func (r *Router) getPowerHandler(c *gin.Context, req *GetPowerRequest) (*netpb.GetPowerResponse, error) {
	return r.clients.Network.GetPower(req.NetworkId, req.SiteId, req.Period, req.From, req.To, req.Timezone)
}

func (r *Router) getAlarmsHandler(c *gin.Context, req *GetAlarmsRequest) (*netpb.GetAlarmsResponse, error) {
	return r.clients.Network.GetAlarms(req.NetworkId, req.SiteId, req.Severity, req.State, req.Page, req.PageSize)
}

func (r *Router) getMetricsHandler(c *gin.Context, req *GetMetricsRequest) (*netpb.GetMetricsResponse, error) {
	return r.clients.Network.GetMetrics(req.NetworkId, req.SiteId, req.NodeId, req.Metric, req.Period, req.From, req.To, req.Timezone)
}

func (r *Router) getEventsHandler(c *gin.Context, req *GetEventsRequest) (*netpb.GetEventsResponse, error) {
	return r.clients.Network.GetEvents(req.NetworkId, req.SiteId, req.NodeId, req.Period, req.From, req.To, req.Timezone, req.Page, req.PageSize)
}

func (r *Router) supportSearchHandler(c *gin.Context, req *SupportSearchRequest) (*netpb.SupportSearchResponse, error) {
	return r.clients.Network.SupportSearch(req.Query, req.NetworkId)
}

/* Collector handlers */

func (r *Router) postRefreshHandler(c *gin.Context, req *PostRefreshRequest) (*colpb.RefreshResponse, error) {
	return r.clients.Collector.Refresh(req.Source)
}

func (r *Router) getRefreshStateHandler(c *gin.Context) (*colpb.GetRefreshStateResponse, error) {
	return r.clients.Collector.GetRefreshState()
}

func (r *Router) postRebuildRollupsHandler(c *gin.Context, req *PostRebuildRollupsRequest) (*colpb.RebuildRollupsResponse, error) {
	return r.clients.Collector.RebuildRollups(req.Family, req.From, req.To)
}

func (r *Router) postSeedDemoHandler(c *gin.Context, req *PostSeedDemoRequest) (*colpb.SeedDemoResponse, error) {
	return r.clients.Collector.SeedDemo(req.Sites, req.Nodes, req.Customers, req.Days)
}

func formatDoc(summary string, description string) []fizz.OperationOption {
	return []fizz.OperationOption{func(info *openapi.OperationInfo) {
		info.Summary = summary
		info.Description = description
	}}
}
