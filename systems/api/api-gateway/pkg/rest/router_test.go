package rest

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ukama/ukama/systems/common/providers"

	"github.com/stretchr/testify/assert"
	"github.com/ukama/ukama/systems/api/api-gateway/pkg"

	// mmocks "github.com/ukama/ukama/systems/api/mailer/pb/gen/mocks"
	// nmocks "github.com/ukama/ukama/systems/api/notify/pb/gen/mocks"
	cconfig "github.com/ukama/ukama/systems/common/config"
	crest "github.com/ukama/ukama/systems/common/rest"
)

const apiEndpoint = "/v1/apis"

var defaultCors = cors.Config{
	AllowAllOrigins: true,
}

var routerConfig = &RouterConfig{
	serverConf: &crest.HttpConfig{
		Cors: defaultCors,
	},
	auth: &cconfig.Auth{
		AuthAppUrl:    "http://localhost:4455",
		AuthServerUrl: "http://localhost:4434",
		AuthAPIGW:     "http://localhost:8080",
	},
}

var testClientSet *Clients

func init() {
	gin.SetMode(gin.TestMode)
	testClientSet = NewClientsSet(&pkg.GrpcEndpoints{
		Timeout: 1 * time.Second,
		Mailer:  "0.0.0.0:9092",
		Notify:  "0.0.0.0:9093",
	})
}

func TestRouter_PingRoute(t *testing.T) {
	// var m = &mmocks.MailerServiceClient{}
	// var n = &nocks.NotifyServiceClient{}
	var arc = &providers.AuthRestClient{}

	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)

	r := NewRouter(&Clients{
		// m: client.NewMailerFromClient(m),
		// n: client.NewNotifyFromClient(n),
	}, routerConfig, arc.MockAuthenticateUser).f.Engine()

	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "pong")
}
