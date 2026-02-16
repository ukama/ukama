/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package metric

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetPrometheusRequestUrl(t *testing.T) {
	t.Run("SingleMetric_NoOperation", func(t *testing.T) {
		pp := PrometheusPayload{
			Metrics:   []string{"cpu_usage"},
			Start:     "1700000000",
			End:       "1700003600",
			Step:      "60",
			Filters:   []Filter{{Key: "node_id", Value: "node-1"}},
			Operation: "",
		}
		prd := GetPrometheusRequestUrl("http://prometheus:9090", pp)
		assert.Equal(t, "http://prometheus:9090/api/v1/query_range", prd.Url)
		assert.Contains(t, prd.Query.Get("query"), "cpu_usage")
		assert.Contains(t, prd.Query.Get("query"), `node_id="node-1"`)
		assert.Equal(t, "1700000000", prd.Query.Get("start"))
		assert.Equal(t, "1700003600", prd.Query.Get("end"))
		assert.Equal(t, "60", prd.Query.Get("step"))
	})

	t.Run("MultipleMetrics_WithOperation", func(t *testing.T) {
		pp := PrometheusPayload{
			Metrics:   []string{"metric_a", "metric_b"},
			Start:     "1000",
			End:       "2000",
			Step:      "1",
			Filters:   []Filter{{Key: "job", Value: "metrics"}},
			Operation: "avg",
		}
		prd := GetPrometheusRequestUrl("http://localhost:9090/", pp)
		assert.Contains(t, prd.Query.Get("query"), "avg")
		assert.Contains(t, prd.Query.Get("query"), "metric_a")
		assert.Contains(t, prd.Query.Get("query"), "metric_b")
	})

	t.Run("BaseUrlTrimmed", func(t *testing.T) {
		pp := PrometheusPayload{Metrics: []string{"m"}, Start: "0", End: "1", Step: "1"}
		prd := GetPrometheusRequestUrl("http://prom:9090/", pp)
		assert.Equal(t, "http://prom:9090/api/v1/query_range", prd.Url)
	})
}

func TestBuildPrometheusRequest(t *testing.T) {
	t.Run("SingleMetric", func(t *testing.T) {
		mq := []MetricWithFilters{
			{Metric: "cpu_percent", Filters: []Filter{{Key: "node_id", Value: "tnode-1"}}},
		}
		prd := BuildPrometheusRequest("http://p:9090", "1000", "2000", "15", "", mq)
		assert.Equal(t, "http://p:9090/api/v1/query_range", prd.Url)
		assert.Contains(t, prd.Query.Get("query"), "cpu_percent")
		assert.Equal(t, "1000", prd.Query.Get("start"))
		assert.Equal(t, "2000", prd.Query.Get("end"))
		assert.Equal(t, "15", prd.Query.Get("step"))
	})

	t.Run("MultipleMetrics_DifferentFilters", func(t *testing.T) {
		mq := []MetricWithFilters{
			{Metric: "cpu", Filters: []Filter{{Key: "node_id", Value: "tnode-1"}}},
			{Metric: "memory", Filters: []Filter{{Key: "node_id", Value: "anode-1"}}},
		}
		prd := BuildPrometheusRequest("http://p:9090", "0", "1", "1", "sum", mq)
		assert.Contains(t, prd.Query.Get("query"), "sum")
		assert.Contains(t, prd.Query.Get("query"), "cpu")
		assert.Contains(t, prd.Query.Get("query"), "memory")
	})

	t.Run("EmptyFilters", func(t *testing.T) {
		mq := []MetricWithFilters{{Metric: "metric_only", Filters: nil}}
		prd := BuildPrometheusRequest("http://p:9090", "0", "1", "1", "", mq)
		assert.Contains(t, prd.Query.Get("query"), "metric_only")
	})
}

func TestToFilteredResponse(t *testing.T) {
	t.Run("ExtractsNodeIDAndMetric", func(t *testing.T) {
		pr := &PrometheusResponse{
			Status: "success",
			Data: PrometheusData{
				ResultType: "matrix",
				Result: []PrometheusResult{
					{
						Metric: map[string]interface{}{
							"node_id": "uk-sa-d3e1-0001-tnode",
							"metric":  "cpu_usage",
							"job":     "ignored",
						},
						Values: [][]interface{}{{1700000000.0, 45.5}},
					},
				},
			},
		}
		filtered := pr.ToFilteredResponse()
		require.Len(t, filtered.Data.Result, 1)
		assert.Equal(t, "uk-sa-d3e1-0001-tnode", filtered.Data.Result[0].Metric.NodeID)
		assert.Equal(t, "cpu_usage", filtered.Data.Result[0].Metric.Metric)
		assert.Equal(t, [][]interface{}{{1700000000.0, 45.5}}, filtered.Data.Result[0].Values)
	})

	t.Run("UsesNameWhenMetricLabelMissing", func(t *testing.T) {
		pr := &PrometheusResponse{
			Status: "success",
			Data: PrometheusData{
				ResultType: "matrix",
				Result: []PrometheusResult{
					{
						Metric: map[string]interface{}{
							"__name__": "custom_metric_total",
							"node_id":  "node-1",
						},
						Values: [][]interface{}{},
					},
				},
			},
		}
		filtered := pr.ToFilteredResponse()
		require.Len(t, filtered.Data.Result, 1)
		assert.Equal(t, "custom_metric_total", filtered.Data.Result[0].Metric.Metric)
	})

	t.Run("EmptyResult", func(t *testing.T) {
		pr := &PrometheusResponse{Status: "success", Data: PrometheusData{ResultType: "matrix", Result: nil}}
		filtered := pr.ToFilteredResponse()
		assert.Empty(t, filtered.Data.Result)
	})
}

func TestParsePrometheusResponse(t *testing.T) {
	t.Run("ValidResponse", func(t *testing.T) {
		body := []byte(`{"status":"success","data":{"resultType":"matrix","result":[{"metric":{"node_id":"n1","__name__":"m1"},"values":[[1000,"10"]]}]}}`)
		pr, err := parsePrometheusResponse(body)
		require.NoError(t, err)
		assert.Equal(t, "success", pr.Status)
		require.Len(t, pr.Data.Result, 1)
		assert.Equal(t, "n1", pr.Data.Result[0].Metric["node_id"])
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		_, err := parsePrometheusResponse([]byte(`not json`))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid prometheus response format")
	})
}

func TestErrPrometheusError(t *testing.T) {
	t.Run("WithBody", func(t *testing.T) {
		e := &ErrPrometheusError{StatusCode: 500, Status: "Internal Server Error", Body: "error details"}
		assert.Contains(t, e.Error(), "500")
		assert.Contains(t, e.Error(), "error details")
	})

	t.Run("WithoutBody", func(t *testing.T) {
		e := &ErrPrometheusError{StatusCode: 404, Status: "Not Found", Body: ""}
		assert.Contains(t, e.Error(), "404")
		assert.NotContains(t, e.Error(), "body")
	})
}

func TestFilterResultsByMetric(t *testing.T) {
	t.Run("FiltersByMetricKey", func(t *testing.T) {
		r := &FilteredPrometheusResponse{
			Status: "success",
			Data: FilteredPrometheusData{
				ResultType: "matrix",
				Result: []FilteredPrometheusResult{
					{Metric: FilteredMetric{NodeID: "n1", Metric: "cpu"}, Values: [][]interface{}{}},
					{Metric: FilteredMetric{NodeID: "n1", Metric: "memory"}, Values: [][]interface{}{}},
					{Metric: FilteredMetric{NodeID: "n2", Metric: "cpu"}, Values: [][]interface{}{}},
				},
			},
		}
		filtered := r.FilterResultsByMetric("cpu")
		require.Len(t, filtered.Data.Result, 2)
		assert.Equal(t, "cpu", filtered.Data.Result[0].Metric.Metric)
		assert.Equal(t, "cpu", filtered.Data.Result[1].Metric.Metric)
	})

	t.Run("NoMatches_ReturnsEmpty", func(t *testing.T) {
		r := &FilteredPrometheusResponse{
			Data: FilteredPrometheusData{
				Result: []FilteredPrometheusResult{
					{Metric: FilteredMetric{Metric: "cpu"}, Values: [][]interface{}{}},
				},
			},
		}
		filtered := r.FilterResultsByMetric("nonexistent")
		assert.Empty(t, filtered.Data.Result)
	})
}

func TestProcessPromRequest(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"status":"success","data":{"resultType":"matrix","result":[{"metric":{"node_id":"n1","__name__":"cpu"},"values":[[1000,50]]}]}}`))
		}))
		defer server.Close()

		prd := PrometheusRequestData{
			Url:   server.URL,
			Query: map[string][]string{"query": {"cpu"}},
		}
		res, err := ProcessPromRequest(context.Background(), prd)
		require.NoError(t, err)
		require.Len(t, res.Data.Result, 1)
		assert.Equal(t, "cpu", res.Data.Result[0].Metric.Metric)
	})

	t.Run("HTTPError", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte("service down"))
		}))
		defer server.Close()

		prd := PrometheusRequestData{Url: server.URL, Query: map[string][]string{}}
		_, err := ProcessPromRequest(context.Background(), prd)
		require.Error(t, err)
		var promErr *ErrPrometheusError
		assert.ErrorAs(t, err, &promErr)
		assert.Equal(t, 503, promErr.StatusCode)
	})

	t.Run("PrometheusErrorStatus", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"status":"error","data":null}`))
		}))
		defer server.Close()

		prd := PrometheusRequestData{Url: server.URL, Query: map[string][]string{}}
		_, err := ProcessPromRequest(context.Background(), prd)
		require.Error(t, err)
		var promErr *ErrPrometheusError
		assert.ErrorAs(t, err, &promErr)
	})
}
