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
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/wI2L/fizz"
	"github.com/wI2L/fizz/openapi"

	"github.com/ukama/ukama/systems/billing/api-gateway/cmd/version"
	"github.com/ukama/ukama/systems/billing/api-gateway/pkg/client"
	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/rest"

	log "github.com/sirupsen/logrus"
	pkg "github.com/ukama/ukama/systems/billing/api-gateway/pkg"
	pb "github.com/ukama/ukama/systems/billing/report/pb/gen"
)

const (
	invoiceCreatedType = "invoice.created"
	invoiceObject      = "invoice"
)

type Router struct {
	f       *fizz.Fizz
	clients *Clients
	config  *RouterConfig
}

type RouterConfig struct {
	debugMode  bool
	serverConf *rest.HttpConfig
	auth       *config.Auth
}

type Clients struct {
	r client.Report
}

func NewClientsSet(grpcEndpoints *pkg.GrpcEndpoints, httpEndpoints *pkg.HttpEndpoints, debugMode bool) *Clients {
	c := &Clients{}

	c.r = client.NewReportClient(grpcEndpoints.Report, grpcEndpoints.Timeout)

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
		serverConf: &svcConf.Server,
		debugMode:  svcConf.DebugMode,
		auth:       svcConf.Auth,
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
	r.f = rest.NewFizzRouter(r.config.serverConf, pkg.SystemName,
		version.Version, r.config.debugMode, r.config.auth.AuthAppUrl+"?redirect=true")

	auth := r.f.Group("/v1", "API GW ", "Payments system version v1", func(ctx *gin.Context) {
		if r.config.auth.BypassAuthMode {
			log.Info("Bypassing auth")

			return
		}

		s := fmt.Sprintf("%s, %s, %s", pkg.SystemName, ctx.Request.Method, ctx.Request.URL.Path)
		ctx.Request.Header.Set("Meta", s)

		err := f(ctx, r.config.auth.AuthServerUrl)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
		}
	})

	auth.Use()
	{
		// Report routes
		reports := auth.Group("reports", "JSON Reports", "Operations on Reports")
		reports.GET("", formatDoc("Get Reports", "Get all Reports of a owner"), tonic.Handler(r.GetReports, http.StatusOK))
		reports.GET("/:report_id", formatDoc("Get Report", "Get a specific report"), tonic.Handler(r.GetReport, http.StatusOK))
		reports.POST("", formatDoc("Add Report", "Add a new report for a owner"), tonic.Handler(r.PostReport, http.StatusCreated))
		reports.PATCH("/:report_id", formatDoc("Update Report", "Update a specific report"), tonic.Handler(r.UpdateReport, http.StatusOK))
		reports.DELETE("/:report_id", formatDoc("Remove Report", "Remove a specific report"), tonic.Handler(r.RemoveReport, http.StatusOK))
	}
}

func (r *Router) PostReport(c *gin.Context, req *WebHookRequest) error {
	log.Infof("Webhook event of type %q for object %q received form billing provider",
		req.WebhookType, req.ObjectType)

	if req.WebhookType != invoiceCreatedType || req.ObjectType != invoiceObject {
		log.Infof("Discarding webhook event %q for object %q on reason: No handler set for webhook or object type",
			req.WebhookType, req.ObjectType)

		c.JSON(http.StatusOK, gin.H{
			"info": "webhook event discarded",
		})

		return nil
	}

	log.Infof("Handling webhook event %q for object %q", req.WebhookType, req.ObjectType)

	rwReportBytes, err := json.Marshal(req.Invoice)
	if err != nil {
		log.Errorf("Failed to marshal RawReport payload into rawReport JSON %v", err)

		return fmt.Errorf("failed to marshal RawReport payload into rawReport JSON %w", err)
	}

	resp, err := r.clients.r.Add(string(rwReportBytes))
	if err == nil {
		c.JSON(http.StatusCreated, resp)
	}

	return err
}

func (r *Router) GetReport(c *gin.Context, req *GetReportRequest) (*pb.ReportResponse, error) {
	return r.clients.r.Get(req.ReportId)
}

func (r *Router) GetReports(c *gin.Context, req *GetReportsRequest) (*pb.ListResponse, error) {
	return r.clients.r.List(req.OwnerId, req.OwnerType, req.NetworkId, req.ReortType, req.IsPaid, req.Count, req.Sort)
}

func (r *Router) UpdateReport(c *gin.Context, req *UpdateReportRequest) (*pb.ReportResponse, error) {
	return r.clients.r.Update(req.ReportId, req.IsPaid)
}

func (r *Router) RemoveReport(c *gin.Context, req *GetReportRequest) error {
	return r.clients.r.Remove(req.ReportId)
}

func formatDoc(summary string, description string) []fizz.OperationOption {
	return []fizz.OperationOption{func(info *openapi.OperationInfo) {
		info.Summary = summary
		info.Description = description
	}}
}
