/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain it at http://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package rest

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/wI2L/fizz"
	"github.com/wI2L/fizz/openapi"

	"github.com/ukama/ukama/systems/analytics/api-gateway/cmd/version"
	"github.com/ukama/ukama/systems/analytics/api-gateway/pkg"
	"github.com/ukama/ukama/systems/analytics/api-gateway/pkg/client"
	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/rest"

	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/analytics/aggregator/pb/gen"
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
	grpcEndpoints *pkg.GrpcEndpoints
	descriptions  *pkg.ServiceDescriptions
}

// aggregator is the interface the router needs; satisfied by client.Aggregator.
type aggregator interface {
	ListKpis() (*pb.ListKpisResponse, error)
	GetKpis(keys []string, span, op string, scope map[string]string) (*pb.GetKpisResponse, error)
	GetKpiTimeSeries(keys []string, span, op, from, to string, scope map[string]string) (*pb.GetKpiTimeSeriesResponse, error)
	GetKpiBreakdown(key, span, op, by string, top int32) (*pb.GetKpiBreakdownResponse, error)
	ListReports() (*pb.ListReportsResponse, error)
	GetPerformanceReport(report, span string, scope map[string]string, top int32) (*pb.GetPerformanceReportResponse, error)
}

type Clients struct {
	Aggregator aggregator
}

func NewClientsSet(endpoints *pkg.GrpcEndpoints) *Clients {
	return &Clients{
		Aggregator: client.NewAggregator(endpoints.Aggregator, endpoints.Timeout),
	}
}

func NewRouterConfig(svcConf *pkg.Config) *RouterConfig {
	return &RouterConfig{
		metricsConfig: svcConf.Metrics,
		serverConf:    &svcConf.Server,
		grpcEndpoints: &svcConf.Services,
		descriptions:  &svcConf.Descriptions,
		debugMode:     svcConf.DebugMode,
		auth:          svcConf.Auth,
	}
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

func (rt *Router) Run() {
	log.Info("Listening on port ", rt.config.serverConf.Port)

	err := rt.f.Engine().Run(fmt.Sprint(":", rt.config.serverConf.Port))
	if err != nil {
		log.Error(err)
	}
}

func (r *Router) init(f func(*gin.Context, string) error) {
	r.f = rest.NewFizzRouter(r.config.serverConf, pkg.SystemName, version.Version,
		r.config.debugMode, r.config.auth.AuthAppUrl+"?redirect=true")

	desc := r.config.descriptions
	if desc == nil {
		desc = &pkg.ServiceDescriptions{}
	}

	if r.config.grpcEndpoints != nil {
		rest.RegisterStatusEndpoint(r.f, pkg.SystemName, map[string]rest.StatusTarget{
			"aggregator": {Host: r.config.grpcEndpoints.Aggregator, Description: desc.Aggregator},
			"ingest":     {Host: r.config.grpcEndpoints.Ingest, Description: desc.Ingest},
			"analysis":   {Host: r.config.grpcEndpoints.Analysis, Description: desc.Analysis},
		}, r.config.grpcEndpoints.Timeout)
	}

	auth := r.f.Group("/v1", "Analytics API gateway", "Analytics system version v1", func(ctx *gin.Context) {
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
		kpis := auth.Group("/analytics/kpis", "KPIs", "Generic KPI read API")
		kpis.GET("", formatDoc("List KPIs", "Self-describing KPI registry: keys, domains, units, scopes, ops"),
			tonic.Handler(r.listKpisHandler, http.StatusOK))
		kpis.GET("/values", formatDoc("Get KPI values", "Latest rollup value per KPI/scope for a span, with trend"),
			tonic.Handler(r.getKpisHandler, http.StatusOK))
		kpis.GET("/timeseries", formatDoc("Get KPI time series", "One value per span bucket over a range"),
			tonic.Handler(r.getKpiTimeSeriesHandler, http.StatusOK))
		kpis.GET("/breakdown", formatDoc("Get KPI breakdown", "Top-N scope values for one KPI"),
			tonic.Handler(r.getKpiBreakdownHandler, http.StatusOK))

		reports := auth.Group("/analytics/reports", "Reports", "Resource performance reports")
		reports.GET("", formatDoc("List reports", "Performance report registry"),
			tonic.Handler(r.listReportsHandler, http.StatusOK))
		reports.GET("/:report", formatDoc("Get performance report", "Resource performance table from latest available KPI values"),
			tonic.Handler(r.getPerformanceReportHandler, http.StatusOK))
	}
}

func formatDoc(summary string, description string) []fizz.OperationOption {
	return []fizz.OperationOption{
		func(info *openapi.OperationInfo) {
			info.Summary = summary
			info.Description = description
		},
	}
}

func (r *Router) listKpisHandler(c *gin.Context) (*pb.ListKpisResponse, error) {
	return r.clients.Aggregator.ListKpis()
}

func (r *Router) getKpisHandler(c *gin.Context, req *GetKpisRequest) (*pb.GetKpisResponse, error) {
	return r.clients.Aggregator.GetKpis(splitKeys(req.Keys), req.Span, req.Op, scopeFilter(req.NetworkId, req.SiteId))
}

func (r *Router) getKpiTimeSeriesHandler(c *gin.Context, req *GetKpiTimeSeriesRequest) (*pb.GetKpiTimeSeriesResponse, error) {
	return r.clients.Aggregator.GetKpiTimeSeries(splitKeys(req.Keys), req.Span, req.Op,
		req.From, req.To, scopeFilter(req.NetworkId, req.SiteId))
}

func (r *Router) getKpiBreakdownHandler(c *gin.Context, req *GetKpiBreakdownRequest) (*pb.GetKpiBreakdownResponse, error) {
	return r.clients.Aggregator.GetKpiBreakdown(req.Key, req.Span, req.Op, req.By, req.Top)
}

func (r *Router) listReportsHandler(c *gin.Context) (*pb.ListReportsResponse, error) {
	return r.clients.Aggregator.ListReports()
}

func (r *Router) getPerformanceReportHandler(c *gin.Context, req *GetPerformanceReportRequest) (*pb.GetPerformanceReportResponse, error) {
	return r.clients.Aggregator.GetPerformanceReport(req.Report, req.Span, scopeFilter(req.NetworkId, ""), req.Top)
}

func splitKeys(csv string) []string {
	parts := strings.Split(csv, ",")
	out := make([]string, 0, len(parts))

	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}

	return out
}

func scopeFilter(networkID, siteID string) map[string]string {
	scope := map[string]string{}
	if networkID != "" {
		scope["network_id"] = networkID
	}

	if siteID != "" {
		scope["site_id"] = siteID
	}

	if len(scope) == 0 {
		return nil
	}

	return scope
}
