/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ukama/ukama/systems/notification/api-gateway/pkg/client"

	cconfig "github.com/ukama/ukama/systems/common/config"
	cmocks "github.com/ukama/ukama/systems/common/mocks"
	crest "github.com/ukama/ukama/systems/common/rest"
	dmocks "github.com/ukama/ukama/systems/notification/distributor/pb/gen/mocks"
	emocks "github.com/ukama/ukama/systems/notification/event-notify/pb/gen/mocks"
	mailerpb "github.com/ukama/ukama/systems/notification/mailer/pb/gen"
	mmocks "github.com/ukama/ukama/systems/notification/mailer/pb/gen/mocks"
	nmocks "github.com/ukama/ukama/systems/notification/notify/pb/gen/mocks"
)

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

func Inti() {
	gin.SetMode(gin.TestMode)
}

func ClientInit(m *mmocks.MailerServiceClient, n *nmocks.NotifyServiceClient, e *emocks.EventToNotifyServiceClient, d *dmocks.DistributorServiceClient) {
	testClientSet = &Clients{
		m: client.NewMailerFromClient(m),
		e: client.NewEventToNotifyFromClient(e),
		d: client.NewDistributorFromClient(d),
	}

}

func TestRouter_PingRoute(t *testing.T) {
	m := &mmocks.MailerServiceClient{}
	n := &nmocks.NotifyServiceClient{}
	e := &emocks.EventToNotifyServiceClient{}
	d := &dmocks.DistributorServiceClient{}
	var arc = &cmocks.AuthClient{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	ClientInit(m, n, e, d)

	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)

	r := NewRouter(testClientSet, routerConfig, arc.AuthenticateUser).f.Engine()

	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "pong")
}

func TestRouter_mailer(t *testing.T) {
	m := &mmocks.MailerServiceClient{}
	n := &nmocks.NotifyServiceClient{}
	e := &emocks.EventToNotifyServiceClient{}
	d := &dmocks.DistributorServiceClient{}
	var arc = &cmocks.AuthClient{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	ClientInit(m, n, e, d)
	r := NewRouter(testClientSet, routerConfig, arc.AuthenticateUser).f.Engine()

	mailer := &mailerpb.GetEmailByIdResponse{
		MailId:       "65d969f7-d63e-44eb-b526-fd200e62a2b0",
		To:           "test@ukama.com",
		TemplateName: "test-template",
		Values:       map[string]string{"Name": "test", "Message": "welcome to ukama"},
	}
	t.Run("sendEmail", func(t *testing.T) {
		data := SendEmailReq{
			To:           []string{"test@ukama.com"},
			TemplateName: "test-template",
			Values:       map[string]interface{}{"Name": "test", "Message": "welcome to ukama"},
			Attachments:  []Attachment{},
		}

		jdata, err := json.Marshal(&data)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/v1/mailer/sendEmail", bytes.NewReader(jdata))
		assert.NoError(t, err)
		newValues := make(map[string]string)
		for key, value := range data.Values {
			newValues[key] = fmt.Sprintf("%v", value)
		}
		pbAttachments := make([]*mailerpb.Attachment, len(data.Attachments))
		for i, att := range data.Attachments {
			pbAttachments[i] = &mailerpb.Attachment{
				Filename:    att.Filename,
				Content:     att.Content,
				ContentType: att.ContentType,
			}
		}
		preq := &mailerpb.SendEmailRequest{
			To:           data.To,
			TemplateName: data.TemplateName,
			Values:       newValues,
			Attachments:  pbAttachments,
		}

		m.On("SendEmail", mock.Anything, preq).Return(&mailerpb.SendEmailResponse{
			Message: "email sent successfully",
			MailId:  "65d969f7-d63e-44eb-b526-fd200e62a2b0",
		}, nil)

		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusOK, w.Code)
		m.AssertExpectations(t)
	})

	t.Run("getEmailById", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/mailer/"+mailer.MailId, nil)

		preq := &mailerpb.GetEmailByIdRequest{
			MailId: mailer.MailId,
		}
		m.On("GetEmailById", mock.Anything, preq).Return(&mailerpb.GetEmailByIdResponse{
			MailId:       mailer.MailId,
			To:           "test@ukama.com",
			TemplateName: "test-template",
			Values:       map[string]string{"Name": "test", "Message": "welcome to ukama"},
		}, nil)

		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusOK, w.Code)
		m.AssertExpectations(t)
	})
	t.Run("sendEmailWithAttachments", func(t *testing.T) {
		data := SendEmailReq{
			To:           []string{"test@ukama.com"},
			TemplateName: "test-template",
			Values:       map[string]interface{}{"Name": "test", "Message": "welcome to ukama"},
			Attachments: []Attachment{
				{
					Filename:    "test.txt",
					Content:     []byte("dGVzdCBjb250ZW50"),
					ContentType: "text/plain",
				},
			},
		}

		jdata, err := json.Marshal(&data)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/v1/mailer/sendEmail", bytes.NewReader(jdata))
		assert.NoError(t, err)

		newValues := make(map[string]string)
		for key, value := range data.Values {
			newValues[key] = fmt.Sprintf("%v", value)
		}

		pbAttachments := make([]*mailerpb.Attachment, len(data.Attachments))
		for i, att := range data.Attachments {
			pbAttachments[i] = &mailerpb.Attachment{
				Filename:    att.Filename,
				Content:     att.Content,
				ContentType: att.ContentType,
			}
		}

		preq := &mailerpb.SendEmailRequest{
			To:           data.To,
			TemplateName: data.TemplateName,
			Values:       newValues,
			Attachments:  pbAttachments,
		}

		m.On("SendEmail", mock.Anything, preq).Return(&mailerpb.SendEmailResponse{
			Message: "email sent successfully",
			MailId:  "65d969f7-d63e-44eb-b526-fd200e62a2b0",
		}, nil)

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		m.AssertExpectations(t)
	})
}
