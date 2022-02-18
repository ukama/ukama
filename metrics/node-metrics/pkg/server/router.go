package server

import (
	"fmt"
	"net/http"
	"strings"

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
	nodes.GET("/metrics", nil, tonic.Handler(r.metricListHandler, http.StatusOK))

}

func (r *Router) metricListHandler(c *gin.Context) ([]string, error) {
	return r.nodeMetrics.List(), nil
}

func (r *Router) metricHandler(c *gin.Context, in *GetNodeMetricsInput) error {
	if !r.nodeMetrics.MetricsExist(in.Metric) {
		return rest.HttpError{
			HttpCode: http.StatusNotFound,
			Message:  "Metric not found"}
	}

	httpCode, err := r.nodeMetrics.GetMetric(strings.ToLower(in.Metric), in.NodeID, &pkg.Interval{
		Start: in.From,
		End:   in.To,
		Step:  in.Step,
	}, c.Writer)

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
