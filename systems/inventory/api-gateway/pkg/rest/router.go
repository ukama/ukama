/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package rest

import (
	"fmt"
	"net/http"

	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/rest"
	"github.com/ukama/ukama/systems/inventory/api-gateway/cmd/version"
	"github.com/ukama/ukama/systems/inventory/api-gateway/pkg"
	"github.com/ukama/ukama/systems/inventory/api-gateway/pkg/client"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/wI2L/fizz"
	"github.com/wI2L/fizz/openapi"

	log "github.com/sirupsen/logrus"
	accountingpb "github.com/ukama/ukama/systems/inventory/accounting/pb/gen"
	componentpb "github.com/ukama/ukama/systems/inventory/component/pb/gen"
)

type Router struct {
	f       *fizz.Fizz
	clients *Clients
	config  *RouterConfig
}

type RouterConfig struct {
	metricsConfig config.Metrics
	httpEndpoints *pkg.HttpEndpoints
	debugMode     bool
	serverConf    *rest.HttpConfig
	auth          *config.Auth
}

type Clients struct {
	Component  component
	Accounting accounting
}

type component interface {
	Get(id string) (*componentpb.GetResponse, error)
	GetByCompany(c string, t string) (*componentpb.GetByCompanyResponse, error)
	SyncComponent() (*componentpb.SyncComponentsResponse, error)
}

type accounting interface {
	Get(id string) (*accountingpb.GetResponse, error)
	GetByCompany(c string) (*accountingpb.GetByCompanmyResponse, error)
	SyncAccounts() (*accountingpb.SyncAcountingResponse, error)
}

func NewClientsSet(endpoints *pkg.GrpcEndpoints) *Clients {
	c := &Clients{}
	c.Accounting = client.NewAccountingInventory(endpoints.Account, endpoints.Timeout)
	c.Component = client.NewComponentInventory(endpoints.Component, endpoints.Timeout)

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
		httpEndpoints: &svcConf.HttpServices,
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
	auth := r.f.Group("/v1", "API gateway", "Inventory system version v1", func(ctx *gin.Context) {
		if r.config.auth.BypassAuthMode {
			log.Info("Bypassing auth")
			return
		}
		s := fmt.Sprintf("%s, %s, %s", pkg.SystemName, ctx.Request.Method, ctx.Request.URL.Path)
		ctx.Request.Header.Set("Meta", s)
		err := f(ctx, r.config.auth.AuthAPIGW)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
			return
		}
		if err == nil {
			return
		}
	})
	auth.Use()
	{
		const component = "/components"
		components := auth.Group(component, "Component", "Operations on Component")
		components.GET("/:uuid", formatDoc("Get component", "Get component by id"), tonic.Handler(r.getComponentByIdHandler, http.StatusOK))
		components.GET("/company/:company", formatDoc("Get components", "Get components by company name"), tonic.Handler(r.getComponentsByCompanyHandler, http.StatusOK))
		components.GET("/sync", formatDoc("Sync components", "Sync components with repo"), tonic.Handler(r.syncComponentHandler, http.StatusOK))

		// Account routes
		const account = "/accounting"
		accounts := auth.Group(account, "Account", "Operations on Account")
		accounts.GET("/:uuid", formatDoc("Get accounting", "Get accounting by id"), tonic.Handler(r.getAccountByIdHandler, http.StatusOK))
		accounts.GET("/company/:company", formatDoc("Get accountings", "Get accountings by company name"), tonic.Handler(r.getAccountsByCompanyHandler, http.StatusOK))
		accounts.GET("/sync", formatDoc("Sync accounts", "Sync accounts with repo"), tonic.Handler(r.syncAccountsHandler, http.StatusOK))
	}
}

func (r *Router) getComponentByIdHandler(c *gin.Context, req *GetRequest) (*componentpb.GetResponse, error) {
	return r.clients.Component.Get(req.Uuid)
}

func (r *Router) getComponentsByCompanyHandler(c *gin.Context, req *GetComponents) (*componentpb.GetByCompanyResponse, error) {
	return r.clients.Component.GetByCompany(req.Company, req.Category)
}

func (r *Router) syncComponentHandler(c *gin.Context) (*componentpb.SyncComponentsResponse, error) {
	return r.clients.Component.SyncComponent()
}

func (r *Router) getAccountByIdHandler(c *gin.Context, req *GetRequest) (*accountingpb.GetResponse, error) {
	return r.clients.Accounting.Get(req.Uuid)
}

func (r *Router) getAccountsByCompanyHandler(c *gin.Context, req *GetAccounts) (*accountingpb.GetByCompanmyResponse, error) {
	return r.clients.Accounting.GetByCompany(req.Company)
}

func (r *Router) syncAccountsHandler(c *gin.Context) (*accountingpb.SyncAcountingResponse, error) {
	return r.clients.Accounting.SyncAccounts()
}

func formatDoc(summary string, description string) []fizz.OperationOption {
	return []fizz.OperationOption{func(info *openapi.OperationInfo) {
		info.Summary = summary
		info.Description = description
	}}
}
