/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package metric

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

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/errors"
)

type PrometheusPayload struct {
	Metrics []string `json:"metrics"`
	Start string `json:"start"`
	End string `json:"end"`
	Step string `json:"step"`
	Filters []Filter `json:"filters"`
	Operation string `json:"operation"`
}

type Filter struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// MetricWithFilters pairs a metric name with its label filters (e.g., node_id per metric).
// Used when multiple metrics in one query need different filters (e.g., tnode vs anode).
type MetricWithFilters struct {
	Metric  string
	Filters []Filter
}

type PrometheusRequestData struct {
	Url     string `json:"url"`
	Query   url.Values `json:"query"`
	Payload PrometheusPayload `json:"payload"`
}

func getFiltersQuery(filters []Filter) string {
	var fq []string
	for _, f := range filters {
		fq = append(fq, fmt.Sprintf("%s='%s'", f.Key, f.Value))
	}
	return strings.ReplaceAll(strings.Join(fq, ","), "'", "\"")
}

func GetPrometheusRequestUrl(
	baseUrl string,
	pp PrometheusPayload,
) PrometheusRequestData {

	queries := make([]string, 0)

	u := fmt.Sprintf("%s/api/v1/query_range", strings.TrimSuffix(baseUrl, "/"))

	data := url.Values{}
	data.Set("start", pp.Start)
	data.Set("end", pp.End)
	data.Set("step", pp.Step)

	filtersQuery := getFiltersQuery(pp.Filters)

	for _, metric := range pp.Metrics {
		queries = append(
			queries,
			fmt.Sprintf(`%s{%s}`, metric, filtersQuery),
		)
	}
	query := fmt.Sprintf("(%s)", strings.Join(queries, " or "))

	if pp.Operation != "" {
		data.Set("query", fmt.Sprintf("%s(%s)", pp.Operation, query))
	} else {
		data.Set("query", query)
	}

	return PrometheusRequestData{
		Url:     u,
		Query:   data,
		Payload: pp,
	}
}

// BuildPrometheusRequest builds a single Prometheus query_range request from metric+filter pairs.
// Each metric gets its own filters in the query: (m1{f1} or m2{f2} or ...).
func BuildPrometheusRequest(baseUrl, start, end, step, operation string, metricQueries []MetricWithFilters) PrometheusRequestData {
	queries := make([]string, 0, len(metricQueries))
	allMetrics := make([]string, 0, len(metricQueries))
	for _, mq := range metricQueries {
		filtersQuery := getFiltersQuery(mq.Filters)
		queries = append(queries, fmt.Sprintf(`%s{%s}`, mq.Metric, filtersQuery))
		allMetrics = append(allMetrics, mq.Metric)
	}
	fullQuery := fmt.Sprintf("(%s)", strings.Join(queries, " or "))
	if operation != "" {
		fullQuery = fmt.Sprintf("%s(%s)", operation, fullQuery)
	}

	u := fmt.Sprintf("%s/api/v1/query_range", strings.TrimSuffix(baseUrl, "/"))
	data := url.Values{}
	data.Set("start", start)
	data.Set("end", end)
	data.Set("step", step)
	data.Set("query", fullQuery)

	return PrometheusRequestData{
		Url:   u,
		Query: data,
		Payload: PrometheusPayload{
			Metrics:   allMetrics,
			Start:     start,
			End:       end,
			Step:      step,
			Operation: operation,
		},
	}
}

// PrometheusResponse represents the raw response from Prometheus API
type PrometheusResponse struct {
	Status string          `json:"status"`
	Data   PrometheusData  `json:"data"`
}

type PrometheusData struct {
	ResultType string            `json:"resultType"`
	Result     []PrometheusResult `json:"result"`
}

type PrometheusResult struct {
	Metric map[string]interface{} `json:"metric"` // __name__, instance, job, metric, node_id, etc.
	Values [][]interface{}       `json:"values"` // [[timestamp, value], ...]
}

// FilteredPrometheusResponse represents the filtered response with only node_id and metric labels
type FilteredPrometheusResponse struct {
	Status string                    `json:"status"`
	Data   FilteredPrometheusData    `json:"data"`
}

type FilteredPrometheusData struct {
	ResultType string                    `json:"resultType"`
	Result     []FilteredPrometheusResult `json:"result"`
}

type FilteredPrometheusResult struct {
	Metric FilteredMetric  `json:"metric"`
	Values [][]interface{} `json:"values"`
}

type FilteredMetric struct {
	NodeID string `json:"node_id"`
	Metric string `json:"metric"`
}

func getStringFromMap(m map[string]interface{}, keys ...string) string {
	for _, k := range keys {
		if v, ok := m[k]; ok {
			if s, ok := v.(string); ok {
				return s
			}
		}
	}
	return ""
}

// ToFilteredResponse converts a raw Prometheus response to the filtered format (only node_id and metric labels)
func (p *PrometheusResponse) ToFilteredResponse() *FilteredPrometheusResponse {
	filtered := &FilteredPrometheusResponse{
		Status: p.Status,
		Data: FilteredPrometheusData{
			ResultType: p.Data.ResultType,
			Result:     make([]FilteredPrometheusResult, 0, len(p.Data.Result)),
		},
	}
	for _, r := range p.Data.Result {
		fm := FilteredMetric{
			NodeID: getStringFromMap(r.Metric, "node_id"),
			Metric: getStringFromMap(r.Metric, "metric", "__name__"),
		}
		filtered.Data.Result = append(filtered.Data.Result, FilteredPrometheusResult{
			Metric: fm,
			Values: r.Values,
		})
	}
	return filtered
}

// parsePrometheusResponse parses raw Prometheus API response body into PrometheusResponse.
// Uses UseNumber to preserve epoch timestamps as integers.
func parsePrometheusResponse(body []byte) (*PrometheusResponse, error) {
	dec := json.NewDecoder(bytes.NewReader(body))
	dec.UseNumber()
	var pr PrometheusResponse
	if err := dec.Decode(&pr); err != nil {
		return nil, errors.Wrap(err, "invalid prometheus response format")
	}
	return &pr, nil
}

// ErrPrometheusError is returned when Prometheus API returns an error status
type ErrPrometheusError struct {
	StatusCode int
	Status     string
	Body       string
}

func (e *ErrPrometheusError) Error() string {
	if e.Body != "" {
		return fmt.Sprintf("prometheus returned error: status=%s, http_status=%d, body=%s", e.Status, e.StatusCode, e.Body)
	}
	return fmt.Sprintf("prometheus returned error: status=%s, http_status=%d", e.Status, e.StatusCode)
}

// ProcessPromRequest makes an HTTP request to Prometheus, parses the response,
// and returns it in FilteredPrometheusResponse format (only node_id and metric labels).
func ProcessPromRequest(ctx context.Context, prd PrometheusRequestData) (*FilteredPrometheusResponse, error) {
	body := strings.NewReader(prd.Query.Encode())
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, prd.Url, body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create prometheus request")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(prd.Query.Encode())))

	log.Debugf("Prometheus request: %s %s", req.Method, prd.Url)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute prometheus request")
	}
	defer func() { _ = res.Body.Close() }()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read prometheus response body")
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, &ErrPrometheusError{
			StatusCode: res.StatusCode,
			Status:     res.Status,
			Body:       truncateString(string(bodyBytes), 500),
		}
	}

	pr, err := parsePrometheusResponse(bodyBytes)
	if err != nil {
		return nil, err
	}

	if pr.Status != "success" {
		return nil, &ErrPrometheusError{
			StatusCode: res.StatusCode,
			Status:     pr.Status,
			Body:       truncateString(string(bodyBytes), 500),
		}
	}

	filtered := pr.ToFilteredResponse()
	log.Debugf("Prometheus response: %d results", len(filtered.Data.Result))
	return filtered, nil
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}