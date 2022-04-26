package server

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/wI2L/fizz/openapi"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/sirupsen/logrus"
	"github.com/ukama/ukamaX/common/rest"
	"github.com/ukama/ukamaX/metrics/node-metrics/cmd/version"
	"github.com/ukama/ukamaX/metrics/node-metrics/pkg"
	"github.com/wI2L/fizz"
)

type Router struct {
	fizz        *fizz.Fizz
	port        int
	nodeMetrics *pkg.Metrics
}

func (r *Router) Run() {
	logrus.Info("Listening on port ", r.port)
	err := r.fizz.Engine().Run(fmt.Sprint(":", r.port))
	if err != nil {
		panic(err)
	}
}

func NewRouter(config *rest.HttpConfig, nodeMetrics *pkg.Metrics) *Router {
	f := rest.NewFizzRouter(config, pkg.ServiceName, version.Version, pkg.IsDebugMode)

	r := &Router{
		fizz:        f,
		port:        config.Port,
		nodeMetrics: nodeMetrics,
	}
	r.init()
	return r
}

func (r *Router) init() {

	nodes := r.fizz.Group("/nodes", "Nodes Metrics", "Query node metrics")

	nodes.GET("metrics/openapi.json", nil, r.fizz.OpenAPI(nil, "json"))
	// metrics
	nodes.GET(":node/metrics/:metric", []fizz.OperationOption{
		func(info *openapi.OperationInfo) {
			info.Description = "Get metrics for a node. Response has Prometheus data format https://prometheus.io/docs/prometheus/latest/querying/api/#range-vectors"
		}}, tonic.Handler(r.metricHandler, http.StatusOK))

	r.fizz.GET("orgs/:org/metrics/:metric", []fizz.OperationOption{
		func(info *openapi.OperationInfo) {
			info.Description = "Get metrics for an org. Response has Prometheus data format https://prometheus.io/docs/prometheus/latest/querying/api/#range-vectors"
		}}, tonic.Handler(r.orgMetricHandler, http.StatusOK))

	r.fizz.GET("orgs/:org/networks/:net/metrics/:metric", []fizz.OperationOption{
		func(info *openapi.OperationInfo) {
			info.Description = "Get metrics for a network. Response has Prometheus data format https://prometheus.io/docs/prometheus/latest/querying/api/#range-vectors"
		}}, tonic.Handler(r.netMetricHandler, http.StatusOK))

	nodes.GET("/metrics", nil, tonic.Handler(r.metricListHandler, http.StatusOK))

}

func (r *Router) metricListHandler(c *gin.Context) ([]string, error) {
	return r.nodeMetrics.List(), nil
}

func (r *Router) orgMetricHandler(c *gin.Context, in *GetOrgMetricsInput) error {
	httpCode, err := r.nodeMetrics.GetAggregateMetric(strings.ToLower(in.Metric), pkg.NewFilter().WithOrg(in.Org), c.Writer)
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

func (r *Router) netMetricHandler(c *gin.Context, in *GetNetMetricsInput) error {
	httpCode, err := r.nodeMetrics.GetAggregateMetric(strings.ToLower(in.Metric), pkg.NewFilter().WithNetwork(in.Org, in.Net), c.Writer)
	return httpErrorOrNil(httpCode, err)
}

func (r *Router) metricHandler(c *gin.Context, in *GetNodeMetricsInput) error {
	return r.requestMetricInternal(c.Writer, in.FilterBase, pkg.NewFilter().WithNodeId(in.NodeID))
}

func (r *Router) requestMetricInternal(writer io.Writer, filterBase FilterBase, filter *pkg.Filter) error {
	if !r.nodeMetrics.MetricsExist(filterBase.Metric) {
		return rest.HttpError{
			HttpCode: http.StatusNotFound,
			Message:  "Metric not found"}
	}

	to := filterBase.To
	if to == 0 {
		to = time.Now().Unix()
	}
	httpCode, err := r.nodeMetrics.GetMetric(strings.ToLower(filterBase.Metric), filter, &pkg.Interval{
		Start: filterBase.From,
		End:   to,
		Step:  filterBase.Step,
	}, writer)

	return httpErrorOrNil(httpCode, err)
}
