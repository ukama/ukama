/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package rest

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/wI2L/fizz"
	"github.com/wI2L/fizz/openapi"

	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/rest"
	"github.com/ukama/ukama/systems/report/api-gateway/cmd/version"
	"github.com/ukama/ukama/systems/report/api-gateway/internal/client"

	log "github.com/sirupsen/logrus"
	internal "github.com/ukama/ukama/systems/report/api-gateway/internal"
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
	p client.Pdf
}

func NewClientsSet(grpcEndpoints *internal.GrpcEndpoints, httpEndpoints *internal.HttpEndpoints, debugMode bool) *Clients {
	c := &Clients{}

	c.p = client.NewPdfClient(httpEndpoints.Files, debugMode)

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

func NewRouterConfig(svcConf *internal.Config) *RouterConfig {
	return &RouterConfig{
		serverConf: &svcConf.Server,
		debugMode:  svcConf.DebugMode,
		auth:       svcConf.Auth,
	}
}

func (rt *Router) Run() {
	log.Info("Listening on port ", rt.config.serverConf.Port)

	if err := rt.f.Engine().Run(fmt.Sprint(":", rt.config.serverConf.Port)); err != nil {
		panic(err)
	}
}

func (r *Router) init(f func(*gin.Context, string) error) {
	r.f = rest.NewFizzRouter(r.config.serverConf, internal.SystemName,
		version.Version, r.config.debugMode, r.config.auth.AuthAppUrl+"?redirect=true")

	auth := r.f.Group("/v1", "API GW ", "Payments system version v1", func(ctx *gin.Context) {
		if r.config.auth.BypassAuthMode {
			log.Info("Bypassing auth")

			return
		}

		s := fmt.Sprintf("%s, %s, %s", internal.SystemName, ctx.Request.Method, ctx.Request.URL.Path)
		ctx.Request.Header.Set("Meta", s)

		err := f(ctx, r.config.auth.AuthServerUrl)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
		}
	})

	auth.Use()
	{
		// pdf file routes
		pdfs := auth.Group("pdf", "PDF Reports", "Operations on report PDF files")
		pdfs.GET("/:report_id", formatDoc("Get Report PDF file", "Get a specific report file as PDF"), tonic.Handler(r.Pdf, http.StatusOK))
	}
}

func (r *Router) Pdf(c *gin.Context, req *GetReportRequest) error {
	content, err := r.clients.p.GetPdf(req.ReportId)
	if err != nil {
		if errors.Is(err, client.ErrInvoicePDFNotFound) {
			c.Status(http.StatusNotFound)
		}

		return err
	}

	fileName := req.ReportId + ".pdf"
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Type", "application/pdf")
	c.Header("Accept-Length", fmt.Sprintf("%d", len(content)))

	_, err = c.Writer.Write(content)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, gin.H{
		"info": "report pdf file downloaded successfully",
	})

	return nil
}

func formatDoc(summary string, description string) []fizz.OperationOption {
	return []fizz.OperationOption{func(info *openapi.OperationInfo) {
		info.Summary = summary
		info.Description = description
	}}
}
