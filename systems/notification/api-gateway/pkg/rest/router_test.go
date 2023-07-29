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



func TestRouter_notification(t *testing.T) {
	cmailer := &mmocks.MailerServiceClient{}
	arc := &providers.AuthRestClient{}

	r := NewRouter(&Clients{
		m: client.NewMailerFromClient(cmailer),
	}, routerConfig, arc.MockAuthenticateUser).f.Engine()
	
	m := &mailerpb.GetEmailByIdResponse{
		MailId: "65d969f7-d63e-44eb-b526-fd200e62a2b0",
		To:	 "test@ukama.com",
		TemplateName: "test-template",
		Values:  map[string]string{"Name": "test","Message": "welcome to ukama"},
		
	}
	t.Run("sendEmail", func(t *testing.T) {
		data := SendEmailReq{
			To:      []string{"test@ukama.com"}, 
			TemplateName: "test-template",
			Values:  map[string]interface{}{"Name": "test","Message": "welcome to ukama"},

		}

		jdata, err := json.Marshal(&data)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/v1/mailer/sendEmail", bytes.NewReader(jdata))
		assert.NoError(t, err)
		newValues:= make(map[string]string)
		for key, value := range data.Values {
			newValues[key] = fmt.Sprintf("%v", value)
		}
		preq := &mailerpb.SendEmailRequest{
			To:     data.To,
			TemplateName: data.TemplateName,
			Values: newValues,
			
		}

		cmailer.On("SendEmail", mock.Anything, preq).Return(&mailerpb.SendEmailResponse{
			Message: "email sent successfully",
			MailId:  "65d969f7-d63e-44eb-b526-fd200e62a2b0",
		}, nil)
		

		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusOK, w.Code)
		cmailer.AssertExpectations(t)
	})


	t.Run("getEmailById", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/mailer/"+m.MailId,nil)
			
		preq := &mailerpb.GetEmailByIdRequest{
			MailId: m.MailId,
		}
		cmailer.On("GetEmailById", mock.Anything, preq).Return(&mailerpb.GetEmailByIdResponse{
			MailId: m.MailId,
			To:	 "test@ukama.com",
			TemplateName: "test-template",
			Values:  map[string]string{"Name": "test","Message": "welcome to ukama"},
			Status: "sent",
		}, nil)

		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusOK, w.Code)
		cmailer.AssertExpectations(t)
	})



}