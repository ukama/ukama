package rest

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/stretchr/testify/assert"
	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/providers"
	"github.com/ukama/ukama/systems/common/rest"
	"github.com/ukama/ukama/systems/metrics/api-gateway/pkg"
	"github.com/ukama/ukama/systems/metrics/api-gateway/pkg/client"
	"github.com/ukama/ukama/systems/metrics/exporter/pb/gen/mocks"
)

var defaultCors = cors.Config{
	AllowAllOrigins: true,
}

var routerConfig = &RouterConfig{
	serverConf: &rest.HttpConfig{
		Cors: defaultCors,
	},
	httpEndpoints: &pkg.HttpEndpoints{
		NodeMetrics: "localhost:8080",
	},
	auth: &config.Auth{
		AuthAppUrl:    "http://localhost:4455",
		AuthServerUrl: "http://localhost:4434",
		AuthAPIGW:     "http://localhost:8080",
	},
}

func init() {
	gin.SetMode(gin.TestMode)
}

// var defaultCongif = &rest.HttpConfig{
// 	Cors: cors.Config{
// 		AllowAllOrigins: true,
// 	},
// }

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

	arc := &providers.AuthRestClient{}
	r := NewRouter(cl, rc, m, arc.MockAuthenticateUser).f.Engine()

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
	c.MetricsConfig.MetricsServer = testSrv.URL
	m, err := pkg.NewMetrics(c.MetricsConfig)
	if err != nil {
		t.Error(err)
	}
	rc := NewRouterConfig(c)
	cl := &Clients{}
	cl.e = client.NewExporterFromClient(&mocks.ExporterServiceClient{})

	arc := &providers.AuthRestClient{}
	r := NewRouter(cl, rc, m, arc.MockAuthenticateUser).f.Engine()

	t.Run("NodeMetrics", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/nodes/node-id/metrics/cpu?from=1643106506&to=1644936312&step=3600", nil)

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, 200, w.Code)
		assert.Contains(t, body, c.MetricsConfig.Metrics["cpu"].Metric)
	})

	t.Run("OrgMetrics", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/orgs/org-id/metrics/cpu?from=1643106506&to=1644936312&step=3600", nil)

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
