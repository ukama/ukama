package rest

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	cconfig "github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/providers"
	"github.com/ukama/ukama/systems/common/rest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ukama/ukama/systems/notification/api-gateway/pkg"
	"github.com/ukama/ukama/systems/notification/api-gateway/pkg/client"
	mailerpb "github.com/ukama/ukama/systems/notification/mailer/pb/gen"
	mmocks "github.com/ukama/ukama/systems/notification/mailer/pb/gen/mocks"
)

var defaultCors = cors.Config{
	AllowAllOrigins: true,
}

var routerConfig = &RouterConfig{
	serverConf: &rest.HttpConfig{
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
	})
}

func TestPingRoute(t *testing.T) {
	// arrange
	m := &mmocks.MailerServiceClient{}

	arc := &providers.AuthRestClient{}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)

	r := NewRouter(&Clients{
		m: client.NewMailerFromClient(m),
	}, routerConfig, arc.MockAuthenticateUser).f.Engine()

	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "pong")
}

func TestRouter_sendEmail(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/mailer/sendEmail",
		strings.NewReader(`{"to": ["brackley@ukama.com"], "subject": "test", "body": "welcome to ukama"}`))

	m := &mmocks.MailerServiceClient{}
	arc := &providers.AuthRestClient{}
	preq := &mailerpb.SendEmailRequest{
		To:      []string{"brackley@ukama.com"},
		Subject: "test",
		Body:    "welcome to ukama",
		Values:  nil,
	}
	m.On("SendEmail", mock.Anything, preq).Return(&mailerpb.SendEmailResponse{
		Message: "email sent successfully",
	}, nil)

	r := NewRouter(&Clients{
		m: client.NewMailerFromClient(m),
	}, routerConfig, arc.MockAuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	m.AssertExpectations(t)

}
