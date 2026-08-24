/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package rest

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ukama/ukama/systems/metrics/api-gateway/pkg"
	"github.com/ukama/ukama/systems/metrics/api-gateway/pkg/client"
	"github.com/ukama/ukama/systems/metrics/exporter/pb/gen/mocks"

	log "github.com/sirupsen/logrus"
	cmocks "github.com/ukama/ukama/systems/common/mocks"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func Test_RouterPing(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	c := pkg.NewConfig()

	m, err := pkg.NewMetrics(c.MetricsConfig)
	if err != nil {
		t.Error(err)
	}
	rc := NewRouterConfig(c)
	cl := &Clients{}
	cl.e = client.NewExporterFromClient(&mocks.ExporterServiceClient{})

	arc := &cmocks.AuthClient{}
	r := NewRouter(cl, rc, m, arc.AuthenticateUser).f.Engine()

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "pong")
}

func Test_GetMetrics(t *testing.T) {
	// arrange
	body := ""
	testSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info(r.URL.String())
		b, err := io.ReadAll(r.Body)
		if err != nil {
			t.Error(err)
		}
		body = string(b)
	}))
	c := pkg.NewConfig()
	pkg.ApplyMetricsFromEnvOverride(c) // load default metrics (e.g. "cpu") so handler can resolve metric keys
	c.MetricsConfig.MetricsServer = testSrv.URL
	m, err := pkg.NewMetrics(c.MetricsConfig)
	if err != nil {
		t.Error(err)
	}
	rc := NewRouterConfig(c)
	cl := &Clients{}
	cl.e = client.NewExporterFromClient(&mocks.ExporterServiceClient{})

	arc := &cmocks.AuthClient{}
	r := NewRouter(cl, rc, m, arc.AuthenticateUser).f.Engine()

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	t.Run("NodeMetrics", func(t *testing.T) {
		w := httptest.NewRecorder()
		// Node ID must contain a recognised node type token so ExtractNodeType resolves "tnode".
		req, _ := http.NewRequest("GET", "/v1/nodes/uk-sa0001-tnode-v0-test/metrics/cpu?from=1643106506&to=1644936312&step=3600", nil)

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, 200, w.Code)
		assert.Contains(t, body, c.MetricsConfig.Metrics["tnode"]["cpu"].Metric)
	})

	t.Run("MissingMetric", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/nodes/uk-sa0001-tnode-v0-test/metrics/test-metrics-miss?from=1643106506&to=1644936312&step=3600", nil)

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, 404, w.Code)
	})

	t.Run("List", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/metrics", nil)

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, 200, w.Code)

		// List returns deduplicated generic keys across all node-type buckets.
		for _, keyMap := range c.MetricsConfig.Metrics {
			for k := range keyMap {
				assert.Contains(t, w.Body.String(), k)
			}
		}

	})

	t.Run("LastMetricWithFilters", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET",
			"/v1/last/metrics/data_usage?package=pkg-1&iccid=8910309414559836625&network=net-1&lookback=24h",
			nil)

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, 200, w.Code)

		q, err := url.ParseQuery(body)
		assert.NoError(t, err)
		assert.Equal(t,
			"last_over_time(data_usage {network='net-1',package='pkg-1',iccid='8910309414559836625'}[24h])",
			q.Get("query"))
	})

	// Unfiltered: matches the hand-run `last_over_time(data_usage[7d])` curl —
	// every series returned, no aggregation wrapper.
	t.Run("LastMetricDefaultsLookback", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/last/metrics/data_usage", nil)

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, 200, w.Code)

		q, err := url.ParseQuery(body)
		assert.NoError(t, err)
		assert.Equal(t, "last_over_time(data_usage {}[7d])", q.Get("query"))
	})

	t.Run("LastMetricWithDeltaFn", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET",
			"/v1/last/metrics/com_uptime?node=uk-sa2634-tnode-v0-4945&lookback=5m&fn=delta", nil)

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, 200, w.Code)

		q, err := url.ParseQuery(body)
		assert.NoError(t, err)
		assert.Equal(t, "delta(com_node_uptime {node_id='uk-sa2634-tnode-v0-4945'}[5m])", q.Get("query"))
	})

	// com_uptime is defined only in the system bucket: a node id must not 404.
	t.Run("LastMetricWithNodeFilterResolvesSystemMetric", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET",
			"/v1/last/metrics/com_uptime?node=uk-sa2634-tnode-v0-4945&lookback=5m", nil)

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, 200, w.Code)

		q, err := url.ParseQuery(body)
		assert.NoError(t, err)
		assert.Equal(t, "last_over_time(com_node_uptime {node_id='uk-sa2634-tnode-v0-4945'}[5m])", q.Get("query"))
	})

	t.Run("LastMetricRejectsUnknownFn", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/last/metrics/data_usage?fn=sum", nil)

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, 400, w.Code)
	})

	t.Run("LastMetricRejectsBadLookback", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/last/metrics/data_usage?lookback=7", nil)

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, 400, w.Code)
	})

	t.Run("LastMetricUnknownMetric", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/last/metrics/test-metrics-miss", nil)

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, 404, w.Code)
	})
}
