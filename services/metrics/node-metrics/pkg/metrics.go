package pkg

import (
	"bytes"
	"context"
	"encoding/json"
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

	// Aggregate function used to aggregate metrics that involve multiple nodes
	AggregateFunc string `json:"aggregateFunc"`
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

type Filter struct {
	nodeId  string
	org     string
	network string
}

func NewFilter() *Filter {
	return &Filter{}
}

func (f *Filter) WithNodeId(nodeId string) *Filter {
	f.nodeId = nodeId
	return f
}

func (f *Filter) WithOrg(org string) *Filter {
	f.org = org
	return f
}

func (f *Filter) HasNetwork() bool {
	return f.network != ""
}

func (f *Filter) WithNetwork(org string, network string) *Filter {
	f.org = org
	f.network = network
	return f
}

// GetFilter returns a prometheus filter
func (f *Filter) GetFilter() string {
	var filter []string
	if f.nodeId != "" {
		filter = append(filter, fmt.Sprintf("nodeid='%s'", f.nodeId))
	}
	if f.org != "" {
		filter = append(filter, fmt.Sprintf("org='%s'", f.org))
	}
	if f.network != "" {
		filter = append(filter, fmt.Sprintf("network='%s'", f.network))
	}
	return strings.Join(filter, ",")
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
func (m *Metrics) GetMetric(metricType string, metricFilter *Filter, in *Interval, w io.Writer) (httpStatus int, err error) {

	_, ok := m.conf.Metrics[metricType]
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
	data.Set("query", m.conf.Metrics[metricType].getQuery(metricFilter, m.conf.DefaultRateInterval, "avg"))

	logrus.Infof("GetMetric query: %s", data.Encode())

	return m.processPromRequest(ctx, u, data, w)
}

// GetAggregateMetric returns aggregated value of a metric based on filter
// uses Sum function by default
func (m *Metrics) GetAggregateMetric(metricType string, metricFilter *Filter, w io.Writer) (httpStatus int, err error) {
	_, ok := m.conf.Metrics[metricType]
	if !ok {
		return http.StatusNotFound, errors.New("metric type not found")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(m.conf.Timeout))
	defer cancel()

	u := fmt.Sprintf("%s/api/v1/query", strings.TrimSuffix(m.conf.MetricsServer, "/"))

	data := url.Values{}
	data.Set("query", m.conf.Metrics[metricType].getAggregateQuery(metricFilter, "sum"))

	logrus.Infof("GetAggregateMetric query: %s", data.Encode())

	return m.processPromRequest(ctx, u, data, w)
}

func (m *Metrics) GetLatestMetric(metricType string, metricFilter *Filter, w io.Writer) (httpStatus int, err error) {
	_, ok := m.conf.Metrics[metricType]
	if !ok {
		return http.StatusNotFound, errors.New("metric type not found")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(m.conf.Timeout))
	defer cancel()

	u := fmt.Sprintf("%s/api/v1/query", strings.TrimSuffix(m.conf.MetricsServer, "/"))

	data := url.Values{}
	data.Set("query", m.conf.Metrics[metricType].getLatestQuery(metricFilter))

	logrus.Infof("GetLatestMetric query: %s", data.Encode())

	buf := bytes.Buffer{}
	st, err := m.processPromRequest(ctx, u, data, &buf)
	if err != nil {
		return st, err
	}

	if st != http.StatusOK {
		return st, nil
	}

	// response struct
	type result struct {
		Data struct {
			Result []struct {
				Value []any `json:"value"`
			} `json:"result"`
		} `json:"data"`
	}
	var res result
	err = json.Unmarshal(buf.Bytes(), &res)
	if err != nil {
		logrus.Errorf("unable to unmarshal response: %s", err)
		return http.StatusInternalServerError, fmt.Errorf("unable to unmarshal response")
	}

	if len(res.Data.Result) == 1 && len(res.Data.Result[0].Value) == 2 {
		_, err = fmt.Fprintf(w, `{ "time": "%.2f", "value": "%v"  }`,
			res.Data.Result[0].Value[0], res.Data.Result[0].Value[1])
		if err != nil {
			logrus.Errorf("unable to write response: %s", err)
			return http.StatusInternalServerError, fmt.Errorf("unable to write response")
		}

	} else {
		return http.StatusInternalServerError, errors.New("unexpected response from server")
	}

	return http.StatusOK, nil
}

func (m *Metrics) processPromRequest(ctx context.Context, url string, data url.Values, w io.Writer) (httpStatusCode int, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(data.Encode()))
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

func getExcludeStatements(labels ...string) string {
	el := []string{"job", "instance", "receive", "tenant_id"}
	el = append(el, labels...)
	return fmt.Sprintf("without (%s)", strings.Join(el, ","))
}

func (m Metric) getQuery(metricFilter *Filter, defaultRateInterval string, aggregateFunc string) string {
	rateInterval := m.RateInterval
	if m.NeedRate && len(rateInterval) == 0 {
		rateInterval = defaultRateInterval
	}

	if m.NeedRate {
		return fmt.Sprintf("%s(rate(%s {%s}[%s])) %s", aggregateFunc, m.Metric,
			metricFilter.GetFilter(), rateInterval, getExcludeStatements())
	}

	return fmt.Sprintf("%s(%s {%s}) %s", aggregateFunc, m.Metric, metricFilter.GetFilter(), getExcludeStatements())
}

func (m Metric) getAggregateQuery(filter *Filter, aggregateFunc string) string {
	exludSt := getExcludeStatements("nodeid")

	// org only filter
	if !filter.HasNetwork() {
		exludSt = getExcludeStatements("nodeid", "network")
	}
	af := m.AggregateFunc
	if len(af) == 0 {
		af = "sum"
	}

	return fmt.Sprintf("%s(%s {%s}) %s", af, m.Metric, filter.GetFilter(), exludSt)
}

func (m Metric) getLatestQuery(filter *Filter) string {
	return fmt.Sprintf("%s {%s}", m.Metric, filter.GetFilter())
}
