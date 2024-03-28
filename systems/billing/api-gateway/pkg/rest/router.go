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

	"github.com/ukama/ukama/systems/billing/api-gateway/cmd/version"
	"github.com/ukama/ukama/systems/billing/api-gateway/pkg/client"
	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/rest"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/wI2L/fizz"
	"github.com/wI2L/fizz/openapi"

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
	metricsConfig config.Metrics
	httpEndpoints *pkg.HttpEndpoints
	debugMode     bool
	serverConf    *rest.HttpConfig
}

type Clients struct {
	Billing billing
}

type billing interface {
	AddInvoice(rawInvoice string) (*pb.AddResponse, error)
	GetInvoice(invoiceId string, asPDF bool) (*pb.GetResponse, error)
	GetInvoicesByInvoicee(invoicee string) (*pb.GetByInvoiceeResponse, error)
	GetInvoicesByNetwork(network string) (*pb.GetByNetworkResponse, error)
	RemoveInvoice(invoiceId string) error
	GetInvoicePDF(invoiceId string) ([]byte, error)
}

func NewClientsSet(endpoints *pkg.GrpcEndpoints) *Clients {
	c := &Clients{}
	c.Billing = client.NewBilling(endpoints.Invoice, endpoints.Files, endpoints.Timeout)

	return c
}

func NewRouter(clients *Clients, config *RouterConfig) *Router {
	r := &Router{
		clients: clients,
		config:  config,
	}

	if !config.debugMode {
		gin.SetMode(gin.ReleaseMode)
	}

	r.init()

	return r
}

func NewRouterConfig(svcConf *pkg.Config) *RouterConfig {
	return &RouterConfig{
		metricsConfig: svcConf.Metrics,
		httpEndpoints: &svcConf.HttpServices,
		serverConf:    &svcConf.Server,
		debugMode:     svcConf.DebugMode,
	}
}

func (rt *Router) Run() {
	log.Info("Listening on port ", rt.config.serverConf.Port)

	err := rt.f.Engine().Run(fmt.Sprint(":", rt.config.serverConf.Port))
	if err != nil {
		panic(err)
	}
}

func (r *Router) init() {
	r.f = rest.NewFizzRouter(r.config.serverConf, pkg.SystemName, version.Version, r.config.debugMode, "")
	v1 := r.f.Group("/v1", "API gateway", "Billing system version v1")

	// Invoice routes
	const invoice = "/invoices"

	invoices := v1.Group(invoice, "JSON Invoice", "Operations on Invoices")
	invoices.GET("", formatDoc("Get Invoices", "Get all Invoices of a invoicee"), tonic.Handler(r.getInvoicesHandler, http.StatusOK))
	invoices.GET("/:invoice_id", formatDoc("Get Invoice", "Get a specific invoice"), tonic.Handler(r.GetInvoiceHandler, http.StatusOK))
	invoices.POST("", formatDoc("Add Invoice", "Add a new invoice for a invoicee"), tonic.Handler(r.postInvoiceHandler, http.StatusCreated))
	// update invoice
	invoices.DELETE("/:invoice_id", formatDoc("Remove Invoice", "Remove a specific invoice"), tonic.Handler(r.removeInvoiceHandler, http.StatusOK))

	const pdf = "/pdf"
	pdfs := v1.Group(pdf, "PDF Invoices", "Operations on invoice PDF files")
	pdfs.GET("/:invoice_id", formatDoc("Get Invoice PDF file", "Get a specific invoice file as PDF"), tonic.Handler(r.GetInvoicePdfHandler, http.StatusOK))
}

func (r *Router) GetInvoiceHandler(c *gin.Context, req *GetInvoiceRequest) (*pb.GetResponse, error) {
	asPDF := false

	pdf, ok := c.GetQuery("type")
	if ok && pdf == "pdf" {
		asPDF = true
	}

	return r.clients.Billing.GetInvoice(req.InvoiceId, asPDF)
}

func (r *Router) getInvoicesHandler(c *gin.Context, req *GetInvoicesRequest) ([]*pb.Invoice, error) {
	var invoices []*pb.Invoice
	var err error

	invoiceeId, sOK := c.GetQuery("invoicee")
	networkId, nOK := c.GetQuery("network")

	if sOK == nOK {
		return nil, rest.HttpError{
			HttpCode: http.StatusBadRequest,
			Message:  "either invoicee or network must be provided as a mandatory query parameter, but not both"}
	}

	if sOK {
		res, sErr := r.clients.Billing.GetInvoicesByInvoicee(invoiceeId)
		if sErr == nil {
			invoices = res.Invoices
		}
		err = sErr
	} else {
		res, nErr := r.clients.Billing.GetInvoicesByNetwork(networkId)
		if nErr == nil {
			invoices = res.Invoices
		}
		err = nErr
	}

	return invoices, err
}

func (r *Router) postInvoiceHandler(c *gin.Context, req *WebHookRequest) error {
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

	resp, err := r.clients.Billing.AddInvoice(string(rwInvoiceBytes))
	if err == nil {
		c.JSON(http.StatusCreated, resp)
	}

	return err
}

func (r *Router) removeInvoiceHandler(c *gin.Context, req *GetInvoiceRequest) error {
	return r.clients.Billing.RemoveInvoice(req.InvoiceId)
}

func (r *Router) GetInvoicePdfHandler(c *gin.Context, req *GetInvoiceRequest) error {
	content, err := r.clients.Billing.GetInvoicePDF(req.InvoiceId)
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
