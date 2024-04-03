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
	"errors"
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
	pb "github.com/ukama/ukama/systems/billing/invoice/pb/gen"
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
	i client.Invoice
	p client.Pdf
}

func NewClientsSet(grpcEndpoints *pkg.GrpcEndpoints, httpEndpoints *pkg.HttpEndpoints, debugMode bool) *Clients {
	c := &Clients{}
	c.i = client.NewInvoiceClient(grpcEndpoints.Invoice, grpcEndpoints.Timeout)
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

			return
		}

		if err == nil {
			return
		}
	})

	auth.Use()
	{
		// Invoice routes
		invoices := auth.Group("invoices", "JSON Invoices", "Operations on Invoices")
		invoices.GET("", formatDoc("Get Invoices", "Get all Invoices of a invoicee"), tonic.Handler(r.GetInvoices, http.StatusOK))
		invoices.GET("/:invoice_id", formatDoc("Get Invoice", "Get a specific invoice"), tonic.Handler(r.GetInvoice, http.StatusOK))
		invoices.POST("", formatDoc("Add Invoice", "Add a new invoice for a invoicee"), tonic.Handler(r.PostInvoice, http.StatusCreated))
		// update invoice
		invoices.DELETE("/:invoice_id", formatDoc("Remove Invoice", "Remove a specific invoice"), tonic.Handler(r.RemoveInvoice, http.StatusOK))

		// pdf file routes
		pdfs := auth.Group("pdf", "PDF Invoices", "Operations on invoice PDF files")
		pdfs.GET("/:invoice_id", formatDoc("Get Invoice PDF file", "Get a specific invoice file as PDF"), tonic.Handler(r.Pdf, http.StatusOK))
	}
}

func (r *Router) GetInvoice(c *gin.Context, req *GetInvoiceRequest) (*pb.GetResponse, error) {
	return r.clients.i.Get(req.InvoiceId, req.AsPdf)
}

func (r *Router) GetInvoices(c *gin.Context, req *GetInvoicesRequest) (*pb.ListResponse, error) {
	return r.clients.i.List(req.InvoiceeId, req.InvoiceeType, req.NetworkId, req.IsPaid, req.Count, req.Sort)
}

func (r *Router) RemoveInvoice(c *gin.Context, req *GetInvoiceRequest) error {
	return r.clients.i.Remove(req.InvoiceId)
}

func (r *Router) PostInvoice(c *gin.Context, req *WebHookRequest) error {
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

	rwInvoiceBytes, err := json.Marshal(req.Invoice)
	if err != nil {
		log.Errorf("Failed to marshal RawInvoice payload into rawInvoice JSON %v", err)

		return fmt.Errorf("failed to marshal RawInvoice payload into rawInvoice JSON %w", err)
	}

	resp, err := r.clients.i.Add(string(rwInvoiceBytes))
	if err == nil {
		c.JSON(http.StatusCreated, resp)
	}

	return err
}

func (r *Router) Pdf(c *gin.Context, req *GetInvoiceRequest) error {
	content, err := r.clients.p.GetPdf(req.InvoiceId)
	if err != nil {
		if errors.Is(err, client.ErrInvoicePDFNotFound) {
			c.Status(http.StatusNotFound)
		}

		return err
	}

	fileName := req.InvoiceId + ".pdf"
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Type", "application/pdf")
	c.Header("Accept-Length", fmt.Sprintf("%d", len(content)))

	_, err = c.Writer.Write(content)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, gin.H{
		"info": "download invoice pdf file successfully",
	})

	return nil
}

func formatDoc(summary string, description string) []fizz.OperationOption {
	return []fizz.OperationOption{func(info *openapi.OperationInfo) {
		info.Summary = summary
		info.Description = description
	}}
}
