package nodemetrics

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const (
	MetricTypeCpu         = "cpu"
	MetricTypeMemory      = "memory"
	MetricTypeActiveUsers = "users"
)

var MetricTypes = []string{
	MetricTypeCpu,
	MetricTypeMemory,
	MetricTypeActiveUsers,
}

type metricQuery struct {
	NeedRate bool
	Metric   string
}

var prometheusMetric = map[string]*metricQuery{
	MetricTypeCpu:         &metricQuery{true, "system_process_cpu_seconds_total"},
	MetricTypeMemory:      &metricQuery{true, "system_process_virtual_memory_bytes"},
	MetricTypeActiveUsers: &metricQuery{false, "epc_active_ue"},
}

type Metrics struct {
	PrometheusUrl string
	Timeout       uint
	// Range vector duration used in Rate func https://prometheus.io/docs/prometheus/latest/querying/basics/#time-durations
	// Example: 1d or 5h, or 30s
	// Should be not less then ScrapeInterval*4 (that's a recommended value)
	// Default is 1h
	RateInterval string
}

type Interval struct {
	// Unix time
	Start int64
	// Unix time
	End int64
	// Step in seconds
	Step uint
}

// GetMetrics returns metrics for specified interval and metric type.
// metricType should be a value from  MetricTypes array (case-sensitive)
func (m *Metrics) GetMetric(metricType string, nodeId string, in *Interval, w io.Writer) (httpStatus int, err error) {

	if _, ok := prometheusMetric[metricType]; !ok {
		return http.StatusBadRequest, errors.New("unknown metric type")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(m.Timeout))
	defer cancel()

	u := fmt.Sprintf("%s/api/v1/query_range", strings.TrimSuffix(m.PrometheusUrl, "/"))

	data := url.Values{}
	data.Set("start", strconv.FormatInt(in.Start, 10))
	data.Set("end", strconv.FormatInt(in.End, 10))
	data.Set("step", strconv.FormatUint(uint64(in.Step), 10))
	data.Set("query", prometheusMetric[metricType].getQuery(nodeId, m.RateInterval))

	logrus.Infof("GetMetric query: %s", data.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, strings.NewReader(data.Encode()))
	if err != nil {
		return http.StatusInternalServerError, errors.Wrap(err, "failed to create request")
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return http.StatusInternalServerError, errors.Wrap(err, "failed to execute request")
	}

	_, err = io.Copy(w, res.Body)
	if err != nil {
		return http.StatusInternalServerError, errors.Wrap(err, "failed to copy response")
	}

	res.Body.Close()
	return res.StatusCode, nil
}

func (m *metricQuery) getQuery(nodeId string, rateInteral string) string {
	if len(rateInteral) == 0 {
		rateInteral = "1h"
	}
	if m.NeedRate {
		return fmt.Sprintf("avg(rate(%s {nodeid='%s'}[%s])) without (job, instance)", m.Metric,
			nodeId, rateInteral)
	}

	return fmt.Sprintf("avg(%s {nodeid='%s'}) without (job, instance)", m.Metric, nodeId)
}
