/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/pkg/errors"
)

type Metrics struct {
	conf *MetricsConfig
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
	nodeId     string
	network    string
	subscriber string
	sim        string
	site       string
	operation  string
}

func NewFilter() *Filter {
	return &Filter{}
}

func (f *Filter) WithNodeId(nodeId string) *Filter {
	f.nodeId = nodeId
	return f
}

func (f *Filter) WithSite(site string) *Filter {
	f.site = site
	return f
}

func (f *Filter) WithSubscriber(network string, subscriber string) *Filter {
	f.network = network
	f.subscriber = subscriber
	return f
}

func (f *Filter) WithSim(network string, subscriber string, sim string) *Filter {
	f.network = network
	f.subscriber = subscriber
	f.sim = sim
	return f
}

func (f *Filter) WithAny(network string, subscriber string, sim string, site string, nodeId string, operation string) *Filter {
	f.network = network
	f.subscriber = subscriber
	f.sim = sim
	f.site = site
	f.nodeId = nodeId
	f.operation = operation
	return f
}

func (f *Filter) HasNetwork() bool {
	return f.network != ""
}

func (f *Filter) WithNetwork(network string) *Filter {
	f.network = network
	return f
}

// GetFilter returns a prometheus filter
func (f *Filter) GetFilter() string {
	var filter []string
	if f.nodeId != "" {
		filter = append(filter, fmt.Sprintf("node_id='%s'", f.nodeId))
	}
	if f.network != "" {
		filter = append(filter, fmt.Sprintf("network='%s'", f.network))
	}
	if f.subscriber != "" {
		filter = append(filter, fmt.Sprintf("subscriber='%s'", f.subscriber))
	}
	if f.sim != "" {
		filter = append(filter, fmt.Sprintf("sim='%s'", f.sim))
	}
	if f.site != "" {
		filter = append(filter, fmt.Sprintf("site='%s'", f.site))
	}
	return strings.Join(filter, ",")
}

// NewMetrics create new instance of metrics
// when metricsConfig in null then defaul config is used
func NewMetrics(config *MetricsConfig) (m *Metrics, err error) {
	return &Metrics{
		conf: config,
	}, nil
}

// GetMetrics returns metrics for specified interval and metric type.
// metricType should be a value from  MetricTypes array (case-sensitive)
func (m *Metrics) GetMetricRange(metricType string, metricFilter *Filter, in *Interval, w io.Writer) (httpStatus int, err error) {

	_, ok := m.conf.Metrics[metricType]
	if !ok {
		return http.StatusNotFound, errors.New("metric type not found")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(m.conf.Timeout))
	defer cancel()

	u := fmt.Sprintf("%s/api/v1/query_range", strings.TrimSuffix(m.conf.MetricsServer, "/"))
	log.Infof("GetMetricRange url: %s", u)

	data := url.Values{}
	data.Set("start", strconv.FormatInt(in.Start, 10))
	data.Set("end", strconv.FormatInt(in.End, 10))
	data.Set("step", strconv.FormatUint(uint64(in.Step), 10))
	data.Set("query", m.conf.Metrics[metricType].getQuery(metricFilter, m.conf.DefaultRateInterval, metricFilter.operation))

	log.Infof("GetMetricRange query: %s", data.Encode())

	return m.processPromRequest(ctx, metricType, u, data, w, false)
}

func (m *Metrics) GetMetric(metricType string, metricFilter *Filter, w io.Writer, formatting bool) (httpStatus int, err error) {

	_, ok := m.conf.Metrics[metricType]
	if !ok {
		return http.StatusNotFound, errors.New("metric type not found")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(m.conf.Timeout))
	defer cancel()

	u := fmt.Sprintf("%s/api/v1/query", strings.TrimSuffix(m.conf.MetricsServer, "/"))

	data := url.Values{}

	data.Set("query", m.conf.Metrics[metricType].getQuery(metricFilter, m.conf.DefaultRateInterval, metricFilter.operation))

	log.Infof("GetMetric query: %s", data.Encode())

	return m.processPromRequest(ctx, metricType, u, data, w, formatting)
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

	log.Infof("GetAggregateMetric query: %s", data.Encode())

	return m.processPromRequest(ctx, metricType, u, data, w, false)
}

func formatMetricsResponse(metricName string, w io.Writer, b io.ReadCloser) error {

	bytes, err := io.ReadAll(b)
	if err != nil {
		log.Errorf("Failed to read prometheus response for %s Error: %v", metricName, err)
		return err
	}

	rmap := map[string]interface{}{}
	err = json.Unmarshal([]byte(bytes), &rmap)
	if err != nil {
		log.Errorf("Failed to unmarshal prometheus response for %s Error: %v", metricName, err)
		return err
	}
	rmap["Name"] = metricName

	rb, err := json.Marshal(rmap)
	if err != nil {
		log.Errorf("Failed to marshal prometheus response for %s Error: %v", metricName, err)
		return err
	}

	n, err := w.Write(rb)
	if err != nil {
		log.Errorf("Failed to add prometheus response to ws response for %s Error: %v", metricName, err)
		return err
	}

	log.Infof("Updated %d bytes of response: %s", n, string(rb))
	return nil
}

func (m *Metrics) processPromRequest(ctx context.Context, metricName string, url string, data url.Values, w io.Writer, formatting bool) (httpStatusCode int, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(data.Encode()))
	if err != nil {
		return http.StatusInternalServerError, errors.Wrap(err, "failed to create request")
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	log.Infof("Request is: %v Body %+v", req, data.Encode())
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return http.StatusInternalServerError, errors.Wrap(err, "failed to execute request")
	}
	log.Infof("Response Body %+v", res.Body)
	if formatting {
		err = formatMetricsResponse(metricName, w, res.Body)
		if err != nil {
			return http.StatusInternalServerError, errors.Wrap(err, "failed to format response")
		}
	} else {
		_, err = io.Copy(w, res.Body)
		if err != nil {
			return http.StatusInternalServerError, errors.Wrap(err, "failed to copy response")
		}
	}

	err = res.Body.Close()
	if err != nil {
		log.Warnf("failed to properly close response body. Error: %v", err)
	}

	return res.StatusCode, nil
}

func UpdatedName(oldName string, slice string, newSlice string) string {
	return strings.Replace(oldName, slice, newSlice, 1)
}

func (m *Metrics) MetricsExist(metricType string) bool {
	_, ok := m.conf.Metrics[metricType]
	return ok
}

func (m *Metrics) MetricsCfg(metricType string) (Metric, bool) {
	me, ok := m.conf.Metrics[metricType]
	return me, ok
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

	if aggregateFunc == "sum" {
		return fmt.Sprintf("%s(%s {%s})", aggregateFunc, m.Metric, metricFilter.GetFilter())
	}

	if m.NeedRate {
		return fmt.Sprintf("%s(rate(%s {%s}[%s])) %s", aggregateFunc, m.Metric,
			metricFilter.GetFilter(), rateInterval, getExcludeStatements())
	}

	return fmt.Sprintf("%s(%s {%s}) %s", aggregateFunc, m.Metric, metricFilter.GetFilter(), getExcludeStatements())
}

func (m Metric) getAggregateQuery(filter *Filter, aggregateFunc string) string {
	exludSt := getExcludeStatements("node_id")

	if !filter.HasNetwork() {
		exludSt = getExcludeStatements("node_id", "network")
	}
	return fmt.Sprintf("%s(%s {%s}) %s", aggregateFunc, m.Metric, filter.GetFilter(), exludSt)
}
