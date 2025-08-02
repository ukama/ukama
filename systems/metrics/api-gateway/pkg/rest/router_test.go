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
		req, _ := http.NewRequest("GET", "/v1/nodes/node-id/metrics/cpu?from=1643106506&to=1644936312&step=3600", nil)

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, 200, w.Code)
		assert.Contains(t, body, c.MetricsConfig.Metrics["cpu"].Metric)
	})

	t.Run("MissingMetric", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/nodes/node-id/metrics/test-metrics-miss?from=1643106506&to=1644936312&step=3600", nil)

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

		for k := range c.MetricsConfig.Metrics {
			assert.Contains(t, w.Body.String(), k)
		}

	})
}
