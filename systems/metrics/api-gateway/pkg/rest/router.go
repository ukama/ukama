package rest

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/sirupsen/logrus"
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

type Router struct {
	f       *fizz.Fizz
	clients *Clients
	config  *RouterConfig
	m       *pkg.Metrics
}

type RouterConfig struct {
	metricsServerConfig config.Metrics
	httpEndpoints       *pkg.HttpEndpoints
	debugMode           bool
	serverConf          *rest.HttpConfig
	metricsConf         *pkg.MetricsConfig
	auth                *config.Auth
}

type Clients struct {
	e exporter
}

type exporter interface {
	/* Yet to add RPC for exporteer.*/
	Dummy(req *pb.DummyParameter) (*pb.DummyParameter, error)
}

func NewClientsSet(endpoints *pkg.GrpcEndpoints, metricHost string, debug bool) *Clients {
	c := &Clients{}
	c.e = client.NewExporter(endpoints.Exporter, endpoints.Timeout)
	return c
}

func NewRouter(clients *Clients, config *RouterConfig, m *pkg.Metrics, authfunc func(*gin.Context, string) error) *Router {

	r := &Router{
		clients: clients,
		config:  config,
		m:       m,
	}

	if !config.debugMode {
		gin.SetMode(gin.ReleaseMode)
	}

	r.init(authfunc)
	return r
}

func NewRouterConfig(svcConf *pkg.Config) *RouterConfig {

	return &RouterConfig{
		metricsServerConfig: svcConf.MetricsServer,
		httpEndpoints:       &svcConf.HttpServices,
		serverConf:          &svcConf.Server,
		metricsConf:         svcConf.MetricsConfig,
		debugMode:           svcConf.DebugMode,
		auth:                svcConf.Auth,
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
	auth := r.f.Group("/v1", "metrics system ", "metrics system version v1", func(ctx *gin.Context) {
		if r.config.auth.BypassAuthMode {
			logrus.Info("Bypassing auth")
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

		auth.GET("/metrics", formatDoc("Get Metrics", ""), tonic.Handler(r.metricListHandler, http.StatusOK))

		auth.GET("/metrics/:metric", []fizz.OperationOption{
			func(info *openapi.OperationInfo) {
				info.Description = "Get metrics. Response has Prometheus data format https://prometheus.io/docs/prometheus/latest/querying/api/#range-vectors"
			}}, tonic.Handler(r.metricHandler, http.StatusOK))

		auth.GET("/subscriber/:subscriber/orgs/:org/networks/:network/metrics/:metric", []fizz.OperationOption{
			func(info *openapi.OperationInfo) {
				info.Description = "Get metrics for a susbcriber. Response has Prometheus data format https://prometheus.io/docs/prometheus/latest/querying/api/#range-vectors"
			}}, tonic.Handler(r.subscriberMetricHandler, http.StatusOK))

		auth.GET("/sims/:sim/orgs/:org/networks/:network/subscribers/:subscriber/metrics/:metric", []fizz.OperationOption{
			func(info *openapi.OperationInfo) {
				info.Description = "Get metrics for a sim. Response has Prometheus data format https://prometheus.io/docs/prometheus/latest/querying/api/#range-vectors"
			}}, tonic.Handler(r.simMetricHandler, http.StatusOK))

		auth.GET("/orgs/:org/metrics/:metric", []fizz.OperationOption{
			func(info *openapi.OperationInfo) {
				info.Description = "Get metrics for an org. Response has Prometheus data format https://prometheus.io/docs/prometheus/latest/querying/api/#range-vectors"
			}}, tonic.Handler(r.orgMetricHandler, http.StatusOK))

		auth.GET("/networks/:network/metrics/:metric", []fizz.OperationOption{
			func(info *openapi.OperationInfo) {
				info.Description = "Get metrics for an network. Response has Prometheus data format https://prometheus.io/docs/prometheus/latest/querying/api/#range-vectors"
			}}, tonic.Handler(r.networkMetricHandler, http.StatusOK))

		auth.GET("/nodes/:node/metrics/:metric", []fizz.OperationOption{
			func(info *openapi.OperationInfo) {
				info.Description = "Get metrics for anode. Response has Prometheus data format https://prometheus.io/docs/prometheus/latest/querying/api/#range-vectors"
			}}, tonic.Handler(r.nodeMetricHandler, http.StatusOK))

		exp := auth.Group("/exporter", "exporter", "exporter")
		exp.GET("", formatDoc("Dummy functions", ""), tonic.Handler(r.getDummyHandler, http.StatusOK))
	}
}

func formatDoc(summary string, description string) []fizz.OperationOption {
	return []fizz.OperationOption{func(info *openapi.OperationInfo) {
		info.Summary = summary
		info.Description = description
	}}
}

func (r *Router) metricHandler(c *gin.Context, in *GetMetricsInput) error {
	httpCode, err := r.m.GetAggregateMetric(strings.ToLower(in.Metric), pkg.NewFilter(), c.Writer)
	return httpErrorOrNil(httpCode, err)
}

func (r *Router) subscriberMetricHandler(c *gin.Context, in *GetSubscriberMetricsInput) error {
	return r.requestMetricInternal(c.Writer, in.FilterBase, pkg.NewFilter().WithSubscriber(in.Org, in.Network, in.Subscriber))
}

func (r *Router) simMetricHandler(c *gin.Context, in *GetSimMetricsInput) error {
	logrus.Infof("Request Sim metrics: %+v", in)

	return r.requestMetricInternal(c.Writer, in.FilterBase, pkg.NewFilter().WithSim(in.Org, in.Network, in.Subscriber, in.Sim))
}

func (r *Router) networkMetricHandler(c *gin.Context, in *GetNetworkMetricsInput) error {
	httpCode, err := r.m.GetAggregateMetric(strings.ToLower(in.Metric), pkg.NewFilter().WithNetwork(in.Network), c.Writer)
	return httpErrorOrNil(httpCode, err)
}
func (r *Router) metricListHandler(c *gin.Context) ([]string, error) {
	return r.m.List(), nil
}

func (r *Router) orgMetricHandler(c *gin.Context, in *GetOrgMetricsInput) error {
	httpCode, err := r.m.GetAggregateMetric(strings.ToLower(in.Metric), pkg.NewFilter().WithOrg(in.Org), c.Writer)
	return httpErrorOrNil(httpCode, err)
}

func httpErrorOrNil(httpCode int, err error) error {
	if err != nil {
		return rest.HttpError{
			HttpCode: httpCode,
			Message:  err.Error()}
	}

	if httpCode != 200 {
		return rest.HttpError{
			HttpCode: httpCode,
			Message:  "Failed to get metric"}
	}
	return nil
}

func (r *Router) nodeMetricHandler(c *gin.Context, in *GetNodeMetricsInput) error {
	return r.requestMetricInternal(c.Writer, in.FilterBase, pkg.NewFilter().WithNodeId(in.NodeID))
}

func (r *Router) requestMetricInternal(writer io.Writer, filterBase FilterBase, filter *pkg.Filter) error {

	ok := r.m.MetricsExist(filterBase.Metric)
	if !ok {
		return rest.HttpError{
			HttpCode: http.StatusNotFound,
			Message:  "Metric not found"}
	}

	to := filterBase.To
	if to == 0 {
		to = time.Now().Unix()
	}

	logrus.Infof("Metrics request with filters: %+v", filter)
	httpCode, err := r.m.GetMetric(strings.ToLower(filterBase.Metric), filter, &pkg.Interval{
		Start: filterBase.From,
		End:   to,
		Step:  filterBase.Step,
	}, writer)

	return httpErrorOrNil(httpCode, err)
}

func (r *Router) getDummyHandler(c *gin.Context, req *DummyParameters) (*pb.DummyParameter, error) {
	return r.clients.e.Dummy(&pb.DummyParameter{})
}
