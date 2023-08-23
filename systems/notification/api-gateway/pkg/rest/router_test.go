package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ukama/ukama/systems/common/providers"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ukama/ukama/systems/notification/api-gateway/pkg"
	"github.com/ukama/ukama/systems/notification/api-gateway/pkg/client"

	cconfig "github.com/ukama/ukama/systems/common/config"
	crest "github.com/ukama/ukama/systems/common/rest"
	mailerpb "github.com/ukama/ukama/systems/notification/mailer/pb/gen"
	mmocks "github.com/ukama/ukama/systems/notification/mailer/pb/gen/mocks"
	npb "github.com/ukama/ukama/systems/notification/notify/pb/gen"
	nmocks "github.com/ukama/ukama/systems/notification/notify/pb/gen/mocks"
)

const notifyApiEndpoint = "/v1/notifications"

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
	var m = &mmocks.MailerServiceClient{}
	var n = &nmocks.NotifyServiceClient{}
	var arc = &providers.AuthRestClient{}

	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)

	r := NewRouter(&Clients{
		m: client.NewMailerFromClient(m),
		n: client.NewNotifyFromClient(n),
	}, routerConfig, arc.MockAuthenticateUser).f.Engine()

	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "pong")
}


	m := &mmocks.MailerServiceClient{}
	arc := &providers.AuthRestClient{}
	preq := &mailerpb.SendEmailRequest{
		To:      []string{"brackley@ukama.com"},
		TemplateName:"test-template",
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