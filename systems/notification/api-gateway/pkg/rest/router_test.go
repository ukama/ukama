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
	epb "github.com/ukama/ukama/systems/notification/event-notify/pb/gen"
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

func TestRouter_eventNotifications(t *testing.T) {
	m := &mmocks.MailerServiceClient{}
	n := &nmocks.NotifyServiceClient{}
	e := &emocks.EventToNotifyServiceClient{}
	d := &dmocks.DistributorServiceClient{}
	var arc = &cmocks.AuthClient{}
	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	ClientInit(m, n, e, d)
	r := NewRouter(testClientSet, routerConfig, arc.AuthenticateUser).f.Engine()

	notificationId := "test-notification-id"
	orgId := "test-org-id"
	networkId := "test-network-id"
	subscriberId := "test-subscriber-id"
	userId := "test-user-id"

	t.Run("getEventNotificationById", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/event-notification/"+notificationId, nil)

		// Mock the event notification service response
		notification := &epb.Notification{
			Id:          notificationId,
			Title:       "Test Notification",
			Description: "Test Description",
			OrgId:       orgId,
			UserId:      userId,
		}

		e.On("Get", mock.Anything, &epb.GetRequest{Id: notificationId}).Return(&epb.GetResponse{
			Notification: notification,
		}, nil)

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		e.AssertExpectations(t)
	})

	t.Run("getEventNotifications", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/event-notification?org_id="+orgId+"&network_id="+networkId+"&subscriber_id="+subscriberId+"&user_id="+userId, nil)

		// Mock the event notification service response
		notifications := []*epb.Notifications{
			{
				Id:          notificationId,
				Title:       "Test Notification 1",
				Description: "Test Description 1",
				Type:        "info",
				Scope:       "org",
				IsRead:      false,
			},
			{
				Id:          "test-notification-id-2",
				Title:       "Test Notification 2",
				Description: "Test Description 2",
				Type:        "warning",
				Scope:       "network",
				IsRead:      true,
			},
		}

		e.On("GetAll", mock.Anything, &epb.GetAllRequest{
			OrgId:        orgId,
			NetworkId:    networkId,
			SubscriberId: subscriberId,
			UserId:       userId,
		}).Return(&epb.GetAllResponse{
			Notifications: notifications,
		}, nil)

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		e.AssertExpectations(t)
	})

	t.Run("updateEventNotificationStatus", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/v1/event-notification/"+notificationId+"?is_read=true", nil)

		// Mock the event notification service response
		e.On("UpdateStatus", mock.Anything, &epb.UpdateStatusRequest{
			Id:     notificationId,
			IsRead: true,
		}).Return(&epb.UpdateStatusResponse{
			Id: notificationId,
		}, nil)

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		e.AssertExpectations(t)
	})

	t.Run("updateEventNotificationStatusToUnread", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/v1/event-notification/"+notificationId+"?is_read=false", nil)

		// Mock the event notification service response
		e.On("UpdateStatus", mock.Anything, &epb.UpdateStatusRequest{
			Id:     notificationId,
			IsRead: false,
		}).Return(&epb.UpdateStatusResponse{
			Id: notificationId,
		}, nil)

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		e.AssertExpectations(t)
	})
}

func TestRouter_mailerErrorCases(t *testing.T) {
	m := &mmocks.MailerServiceClient{}
	n := &nmocks.NotifyServiceClient{}
	e := &emocks.EventToNotifyServiceClient{}
	d := &dmocks.DistributorServiceClient{}
	var arc = &cmocks.AuthClient{}
	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	ClientInit(m, n, e, d)
	r := NewRouter(testClientSet, routerConfig, arc.AuthenticateUser).f.Engine()

	t.Run("sendEmailWithInvalidJSON", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/v1/mailer/sendEmail", bytes.NewReader([]byte("invalid json")))

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("sendEmailWithMissingRequiredFields", func(t *testing.T) {
		data := map[string]interface{}{
			"to": []string{}, // Empty recipients
			// Missing template_name
		}

		jdata, err := json.Marshal(&data)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/v1/mailer/sendEmail", bytes.NewReader(jdata))

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("sendEmailServiceError", func(t *testing.T) {
		data := SendEmailReq{
			To:           []string{"test@ukama.com"},
			TemplateName: "test-template",
			Values:       map[string]interface{}{"Name": "test"},
			Attachments:  []Attachment{},
		}

		jdata, err := json.Marshal(&data)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/v1/mailer/sendEmail", bytes.NewReader(jdata))

		preq := &mailerpb.SendEmailRequest{
			To:           data.To,
			TemplateName: data.TemplateName,
			Values:       map[string]string{"Name": "test"},
			Attachments:  []*mailerpb.Attachment{},
		}

		m.On("SendEmail", mock.Anything, preq).Return(nil, fmt.Errorf("service unavailable"))

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		m.AssertExpectations(t)
	})

	t.Run("getEmailByIdNotFound", func(t *testing.T) {
		mailId := "non-existent-id"
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/mailer/"+mailId, nil)

		preq := &mailerpb.GetEmailByIdRequest{
			MailId: mailId,
		}
		m.On("GetEmailById", mock.Anything, preq).Return(nil, fmt.Errorf("email not found"))

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		m.AssertExpectations(t)
	})
}

func TestRouter_eventNotificationsErrorCases(t *testing.T) {
	m := &mmocks.MailerServiceClient{}
	n := &nmocks.NotifyServiceClient{}
	e := &emocks.EventToNotifyServiceClient{}
	d := &dmocks.DistributorServiceClient{}
	var arc = &cmocks.AuthClient{}
	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	ClientInit(m, n, e, d)
	r := NewRouter(testClientSet, routerConfig, arc.AuthenticateUser).f.Engine()

	notificationId := "test-notification-id"
	orgId := "test-org-id"
	networkId := "test-network-id"
	subscriberId := "test-subscriber-id"
	userId := "test-user-id"

	t.Run("getEventNotificationByIdNotFound", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/event-notification/"+notificationId, nil)

		e.On("Get", mock.Anything, &epb.GetRequest{Id: notificationId}).Return(nil, fmt.Errorf("notification not found"))

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		e.AssertExpectations(t)
	})

	t.Run("getEventNotificationsServiceError", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/event-notification?org_id="+orgId+"&network_id="+networkId+"&subscriber_id="+subscriberId+"&user_id="+userId, nil)

		e.On("GetAll", mock.Anything, &epb.GetAllRequest{
			OrgId:        orgId,
			NetworkId:    networkId,
			SubscriberId: subscriberId,
			UserId:       userId,
		}).Return(nil, fmt.Errorf("service unavailable"))

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		e.AssertExpectations(t)
	})

	t.Run("updateEventNotificationStatusError", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/v1/event-notification/"+notificationId+"?is_read=true", nil)

		e.On("UpdateStatus", mock.Anything, &epb.UpdateStatusRequest{
			Id:     notificationId,
			IsRead: true,
		}).Return(nil, fmt.Errorf("update failed"))

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		e.AssertExpectations(t)
	})
}

func TestRouter_mailerEdgeCases(t *testing.T) {
	var arc = &cmocks.AuthClient{}
	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	t.Run("sendEmailWithNonStringValues", func(t *testing.T) {
		m := &mmocks.MailerServiceClient{}
		n := &nmocks.NotifyServiceClient{}
		e := &emocks.EventToNotifyServiceClient{}
		d := &dmocks.DistributorServiceClient{}

		ClientInit(m, n, e, d)
		r := NewRouter(testClientSet, routerConfig, arc.AuthenticateUser).f.Engine()

		data := SendEmailReq{
			To:           []string{"test@ukama.com"},
			TemplateName: "test-template",
			Values: map[string]interface{}{
				"string_value": "test",
				"int_value":    123,
				"bool_value":   true,
				"float_value":  3.14,
			},
			Attachments: []Attachment{},
		}

		jdata, err := json.Marshal(&data)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/v1/mailer/sendEmail", bytes.NewReader(jdata))

		// Only string values are processed by the handler
		expectedValues := map[string]string{
			"string_value": "test",
			"int_value":    "123",
			"bool_value":   "true",
			"float_value":  "3.14",
		}

		preq := &mailerpb.SendEmailRequest{
			To:           data.To,
			TemplateName: data.TemplateName,
			Values:       expectedValues,
			Attachments:  []*mailerpb.Attachment{},
		}

		m.On("SendEmail", mock.Anything, preq).Return(&mailerpb.SendEmailResponse{
			Message: "email sent successfully",
			MailId:  "test-mail-id",
		}, nil)

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		m.AssertExpectations(t)
	})

	t.Run("sendEmailWithMultipleAttachments", func(t *testing.T) {
		m := &mmocks.MailerServiceClient{}
		n := &nmocks.NotifyServiceClient{}
		e := &emocks.EventToNotifyServiceClient{}
		d := &dmocks.DistributorServiceClient{}

		ClientInit(m, n, e, d)
		r := NewRouter(testClientSet, routerConfig, arc.AuthenticateUser).f.Engine()

		data := SendEmailReq{
			To:           []string{"test@ukama.com"},
			TemplateName: "test-template",
			Values:       map[string]interface{}{"Name": "test"},
			Attachments: []Attachment{
				{
					Filename:    "document1.pdf",
					Content:     []byte("pdf content 1"),
					ContentType: "application/pdf",
				},
				{
					Filename:    "document2.txt",
					Content:     []byte("text content 2"),
					ContentType: "text/plain",
				},
				{
					Filename:    "image.jpg",
					Content:     []byte("image content"),
					ContentType: "image/jpeg",
				},
			},
		}

		jdata, err := json.Marshal(&data)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/v1/mailer/sendEmail", bytes.NewReader(jdata))

		pbAttachments := []*mailerpb.Attachment{
			{
				Filename:    "document1.pdf",
				Content:     []byte("pdf content 1"),
				ContentType: "application/pdf",
			},
			{
				Filename:    "document2.txt",
				Content:     []byte("text content 2"),
				ContentType: "text/plain",
			},
			{
				Filename:    "image.jpg",
				Content:     []byte("image content"),
				ContentType: "image/jpeg",
			},
		}

		preq := &mailerpb.SendEmailRequest{
			To:           data.To,
			TemplateName: data.TemplateName,
			Values:       map[string]string{"Name": "test"},
			Attachments:  pbAttachments,
		}

		m.On("SendEmail", mock.Anything, preq).Return(&mailerpb.SendEmailResponse{
			Message: "email sent successfully",
			MailId:  "test-mail-id",
		}, nil)

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		m.AssertExpectations(t)
	})

	t.Run("sendEmailWithEmptyValues", func(t *testing.T) {
		m := &mmocks.MailerServiceClient{}
		n := &nmocks.NotifyServiceClient{}
		e := &emocks.EventToNotifyServiceClient{}
		d := &dmocks.DistributorServiceClient{}

		ClientInit(m, n, e, d)
		r := NewRouter(testClientSet, routerConfig, arc.AuthenticateUser).f.Engine()

		data := SendEmailReq{
			To:           []string{"test@ukama.com"},
			TemplateName: "test-template",
			Values:       map[string]interface{}{}, // Empty values
			Attachments:  []Attachment{},
		}

		jdata, err := json.Marshal(&data)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/v1/mailer/sendEmail", bytes.NewReader(jdata))

		preq := &mailerpb.SendEmailRequest{
			To:           data.To,
			TemplateName: data.TemplateName,
			Values:       map[string]string{}, // Should be empty
			Attachments:  []*mailerpb.Attachment{},
		}

		m.On("SendEmail", mock.Anything, preq).Return(&mailerpb.SendEmailResponse{
			Message: "email sent successfully",
			MailId:  "test-mail-id",
		}, nil)

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		m.AssertExpectations(t)
	})
}
