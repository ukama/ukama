/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package notification_test

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/tj/assert"

	"github.com/ukama/ukama/systems/common/rest/client/notification"
)
 
 func TestMailerClient_Send(t *testing.T) {
	 baseURL := "http://test-mailer-service.com"
 
	 t.Run("MailSent", func(tt *testing.T) {
		 mockTransport := func(req *http.Request) *http.Response {
			 // Test request parameters
			 assert.Equal(tt, req.URL.String(), baseURL+notification.MailerEndpoint+"/sendEmail")
 
			 // Send mock response
			 return &http.Response{
				 StatusCode: 200,
				 Status:     "200 OK",
				 // Must be set to non-nil value or it panics
				 Header: make(http.Header),
				 Body:   io.NopCloser(bytes.NewBufferString("")),
			 }
		 }
 
		 testMailerClient := notification.NewMailerClient(baseURL)
 
		 // We replace the transport mechanism by mocking the http request
		 // so that the test stays a unit test e.g no server/network call.
		 testMailerClient.R.C.SetTransport(RoundTripFunc(mockTransport))
 
		 err := testMailerClient.SendEmail(
			 notification.SendEmailReq{
				 To:           []string{"johndoe@example.com"},
				 TemplateName: "mail.html.tmpl",
				 Values:       map[string]interface{}{},
			 })
 
		 assert.NoError(tt, err)
	 })
 
	 t.Run("InvalidResponseHeader", func(tt *testing.T) {
		 mockTransport := func(req *http.Request) *http.Response {
			 assert.Equal(tt, req.URL.String(), baseURL+notification.MailerEndpoint+"/sendEmail")
 
			 // error payload
			 resp := `{"error":"internal server error"}`
 
			 return &http.Response{
				 StatusCode: 500,
				 Status:     "500 INTERNAL SERVER ERROR",
				 Body:       io.NopCloser(bytes.NewBufferString(resp)),
				 Header:     http.Header{"Content-Type": []string{"application/json"}},
			 }
		 }
 
		 testMailerClient := notification.NewMailerClient(baseURL)
 
		 testMailerClient.R.C.SetTransport(RoundTripFunc(mockTransport))
 
		 err := testMailerClient.SendEmail(
			 notification.SendEmailReq{
				 To:           []string{"johndoe@example.com"},
				 TemplateName: "mail.html.tmpl",
				 Values:       map[string]interface{}{},
			 })
 
		 assert.Error(tt, err)
	 })
 
	 t.Run("RequestFailure", func(tt *testing.T) {
		 mockTransport := func(req *http.Request) *http.Response {
			 assert.Equal(tt, req.URL.String(), baseURL+notification.MailerEndpoint+"/sendEmail")
 
			 return nil
		 }
 
		 testMailerClient := notification.NewMailerClient(baseURL)
 
		 testMailerClient.R.C.SetTransport(RoundTripFunc(mockTransport))
 
		 err := testMailerClient.SendEmail(
			 notification.SendEmailReq{
				 To:           []string{"johndoe@example.com"},
				 TemplateName: "mail.html.tmpl",
				 Values:       map[string]interface{}{},
			 })
 
		 assert.Error(tt, err)
	 })
 }
 
 type RoundTripFunc func(req *http.Request) *http.Response
 
 func (r RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	 return r(req), nil
 }