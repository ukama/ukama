package rest

import (
	"fmt"
	"net/http"

	"github.com/ukama/ukama/systems/common/rest"

	"github.com/loopfz/gadgeto/tonic"
	"github.com/ukama/ukama/systems/billing/api-gateway/cmd/version"
	"github.com/ukama/ukama/systems/common/config"
	"github.com/wI2L/fizz"
	"github.com/wI2L/fizz/openapi"

	internal "github.com/ukama/ukama/systems/billing/api-gateway/internal"
	"github.com/ukama/ukama/systems/billing/api-gateway/internal/client"
	pb "github.com/ukama/ukama/systems/billing/invoice/pb/gen"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Router struct {
	f       *fizz.Fizz
	clients *Clients
	config  *RouterConfig
}

type RouterConfig struct {
	metricsConfig config.Metrics
	httpEndpoints *internal.HttpEndpoints
	debugMode     bool
	serverConf    *rest.HttpConfig
}

type Clients struct {
	Billing billing
}

type billing interface {
	AddInvoice(subscriberId string, rawInvoice string) (*pb.AddResponse, error)
	GetInvoice(invoiceId string, asPDF bool) (*pb.GetResponse, error)
	GetInvoices(subscriber string, asPDF bool) (*pb.GetBySubscriberResponse, error)
}

func NewClientsSet(endpoints *internal.GrpcEndpoints) *Clients {
	c := &Clients{}
	c.Billing = client.NewBilling(endpoints.Invoice, endpoints.Timeout)
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

func NewRouterConfig(svcConf *internal.Config) *RouterConfig {
	return &RouterConfig{
		metricsConfig: svcConf.Metrics,
		httpEndpoints: &svcConf.HttpServices,
		serverConf:    &svcConf.Server,
		debugMode:     svcConf.DebugMode,
	}
}

func (rt *Router) Run() {
	logrus.Info("Listening on port ", rt.config.serverConf.Port)

	err := rt.f.Engine().Run(fmt.Sprint(":", rt.config.serverConf.Port))
	if err != nil {
		panic(err)
	}
}

func (r *Router) init() {
	r.f = rest.NewFizzRouter(r.config.serverConf, internal.SystemName, version.Version, r.config.debugMode)
	v1 := r.f.Group("/v1", "API gateway", "Billing system version v1")

	// Invoice routes
	const invoice = "/invoices"

	invoices := v1.Group(invoice, "Invoices", "Operations on Invoices")
	invoices.GET("", formatDoc("Get Invoices", "Get all Invoices of a subscriber"), tonic.Handler(r.getInvoicesHandler, http.StatusOK))
	invoices.GET("/:invoice_id", formatDoc("Get Invoice", "Get a specific invoice"), tonic.Handler(r.GetInvocieHandler, http.StatusOK))
	invoices.POST("", formatDoc("Add Network", "Add a new network to an organization"), tonic.Handler(r.postInvoiceHandler, http.StatusCreated))
	// update invoice
	// delete invoice

}

func (r *Router) GetInvocieHandler(c *gin.Context, req *GetInvoiceRequest) (*pb.GetResponse, error) {
	asPDF := false

	pdf, ok := c.GetQuery("type")
	if ok && pdf == "pdf" {
		asPDF = true
	}

	return r.clients.Billing.GetInvoice(req.InvoiceId, asPDF)
}

func (r *Router) getInvoicesHandler(c *gin.Context, req *GetInvoicesRequest) (*pb.GetBySubscriberResponse, error) {
	subscriberId, ok := c.GetQuery("subscriber")
	if !ok {
		return nil, &rest.HttpError{HttpCode: http.StatusBadRequest,
			Message: "subscriber is a mandatory query parameter"}
	}

	asPDF := false

	pdf, ok := c.GetQuery("type")
	if ok && pdf == "pdf" {
		asPDF = true
	}

	return r.clients.Billing.GetInvoices(subscriberId, asPDF)
}

func (r *Router) postInvoiceHandler(c *gin.Context, req *AddInvoiceRequest) (*pb.AddResponse, error) {
	return r.clients.Billing.AddInvoice(req.SubscriberId, req.RawInvoice)
}

func formatDoc(summary string, description string) []fizz.OperationOption {
	return []fizz.OperationOption{func(info *openapi.OperationInfo) {
		info.Summary = summary
		info.Description = description
	}}
}
