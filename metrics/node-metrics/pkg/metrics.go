package pkg

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/pkg/errors"
)

type Metric struct {
	NeedRate bool   `json:"needRate"`
	Metric   string `json:"metric"`
	// Range vector duration used in Rate func https://prometheus.io/docs/prometheus/latest/querying/basics/#time-durations
	// if NeedRate is false then this field is ignored
	// Example: 1d or 5h, or 30s
	RateInterval string `json:"rateInterval"`
}

type Metrics struct {
	conf *NodeMetricsConfig
}

type Interval struct {
	// Unix time
	Start int64
	// Unix time
	End int64
	// Step in seconds
	Step uint
}

// NewMetrics create new instance of metrics
// when metricsConfig in null then defaul config is used
func NewMetrics(config *NodeMetricsConfig) (m *Metrics, err error) {
	return &Metrics{
		conf: config,
	}, nil
}

// GetMetrics returns metrics for specified interval and metric type.
// metricType should be a value from  MetricTypes array (case-sensitive)
func (m *Metrics) GetMetric(metricType string, nodeId string, in *Interval, w io.Writer) (httpStatus int, err error) {

	mi, ok := m.conf.Metrics[metricType]
	if !ok {
		return http.StatusNotFound, errors.New("metric type not found")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(m.conf.Timeout))
	defer cancel()

	u := fmt.Sprintf("%s/api/v1/query_range", strings.TrimSuffix(m.conf.MetricsServer, "/"))

	data := url.Values{}
	data.Set("start", strconv.FormatInt(in.Start, 10))
	data.Set("end", strconv.FormatInt(in.End, 10))
	data.Set("step", strconv.FormatUint(uint64(in.Step), 10))
	data.Set("query", m.conf.Metrics[metricType].getQuery(nodeId, mi.RateInterval, m.conf.DefaultRateInterval))

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

func (m *Metrics) MetricsExist(metricType string) bool {
	_, ok := m.conf.Metrics[metricType]
	return ok
}

func (m *Metrics) List() (r []string) {
	for k := range m.conf.Metrics {
		r = append(r, k)
	}
	return r
}

func (m Metric) getQuery(nodeId string, rateInteral string, defaultRateInterval string) string {
	if m.NeedRate && len(rateInteral) == 0 {
		rateInteral = defaultRateInterval
	}
	if m.NeedRate {
		return fmt.Sprintf("avg(rate(%s {nodeid='%s'}[%s])) without (job, instance)", m.Metric,
			nodeId, rateInteral)
	}

	return fmt.Sprintf("avg(%s {nodeid='%s'}) without (job, instance)", m.Metric, nodeId)
}
