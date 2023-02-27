package rest

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/config"
	"github.com/wI2L/fizz"
	"github.com/wI2L/fizz/openapi"

	"github.com/ukama/ukama/systems/common/rest"
	"github.com/ukama/ukama/systems/metrics/api-gateway/cmd/version"
	"github.com/ukama/ukama/systems/metrics/api-gateway/pkg"
	"github.com/ukama/ukama/systems/metrics/api-gateway/pkg/client"
	pb "github.com/ukama/ukama/systems/metrics/exporter/pb/gen"
)

const METRICS_URL_PARAMETER = "metrics"
const EXPORTER_URL_PARAMETER = "exporter"

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
	e  exporter
	ms metrics
}

type metrics interface {
}

type exporter interface {
	Dummy(req *pb.DummyParameter) (*pb.DummyParameter, error)
}

func NewClientsSet(endpoints *pkg.GrpcEndpoints, metricHost string, debug bool) *Clients {
	var err error
	c := &Clients{}
	c.e = client.NewExporter(endpoints.Exporter, endpoints.Timeout)
	c.ms, err = client.NewMetricsStore(metricHost, debug)
	if err != nil {
		log.Fatalf("Intialization failure %s", err.Error())
	}
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

	r.f = rest.NewFizzRouter(r.config.serverConf, pkg.SystemName, version.Version, r.config.debugMode)
	v1 := r.f.Group("/v1", "metrics system ", "metrics system version v1")

	metrics := v1.Group(METRICS_URL_PARAMETER, "metrics", "metrics")
	metrics.GET("", formatDoc("Get Orgs Credential", ""), tonic.Handler(r.getDummyHandler, http.StatusOK))

	exp := v1.Group(METRICS_URL_PARAMETER, "exporter", "exporter")
	exp.GET("", formatDoc("Dummy functions", ""), tonic.Handler(r.getDummyHandler, http.StatusOK))

}

func formatDoc(summary string, description string) []fizz.OperationOption {
	return []fizz.OperationOption{func(info *openapi.OperationInfo) {
		info.Summary = summary
		info.Description = description
	}}
}

func (r *Router) getDummyHandler(c *gin.Context, req *DummyParameters) (*pb.DummyParameter, error) {
	return r.clients.e.Dummy(&pb.DummyParameter{})
}
