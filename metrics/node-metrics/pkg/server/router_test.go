package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"

	"github.com/gin-contrib/cors"
	"github.com/stretchr/testify/assert"
	"github.com/ukama/ukamaX/common/rest"
	"github.com/ukama/ukamaX/metrics/node-metrics/pkg"
)

func init() {
	pkg.IsDebugMode = true
}

var defaultCongif = &rest.HttpConfig{
	Cors: cors.Config{
		AllowAllOrigins: true,
	},
}

func Test_RouterPing(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	c := pkg.NewConfig()

	m, err := pkg.NewMetrics(c.NodeMetrics)
	if err != nil {
		t.Error(err)
	}

	r := NewRouter(defaultCongif, m).fizz.Engine()

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
		logrus.Info(r.URL.String())
		b, err := io.ReadAll(r.Body)
		if err != nil {
			t.Error(err)
		}
		body = string(b)
	}))
	c := pkg.NewConfig()
	c.NodeMetrics.MetricsServer = testSrv.URL

	m, err := pkg.NewMetrics(c.NodeMetrics)
	if err != nil {
		t.Error(err)
	}

	r := NewRouter(defaultCongif, m).fizz.Engine()

	t.Run("NodeMetrics", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/nodes/node-id/metrics/cpu?from=1643106506&to=1644936312&step=3600", nil)

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, 200, w.Code)
		assert.Contains(t, body, c.NodeMetrics.Metrics["cpu"].Metric)
	})

	t.Run("OrgMetrics", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/orgs/org-id/metrics/cpu?from=1643106506&to=1644936312&step=3600", nil)

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, 200, w.Code)
		assert.Contains(t, body, c.NodeMetrics.Metrics["cpu"].Metric)
	})

	t.Run("MissingMetric", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/nodes/node-id/metrics/test-metrics-miss?from=1643106506&to=1644936312&step=3600", nil)

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, 404, w.Code)
	})

	t.Run("List", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/nodes/metrics", nil)

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, 200, w.Code)

		for k := range c.NodeMetrics.Metrics {
			assert.Contains(t, w.Body.String(), k)
		}

	})

}
