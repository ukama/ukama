/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package auth_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/tj/assert"

	"github.com/ukama/ukama/systems/common/rest/client"
	"github.com/ukama/ukama/systems/common/rest/client/auth"
)

func TestAuthClient_AuthenticateUser(t *testing.T) {
	// Mock the gin context
	ginContext, _ := gin.CreateTestContext(httptest.NewRecorder())

	// Mock the request
	ginContext.Request = httptest.NewRequest("GET", auth.AuthEndpoint, nil)

	t.Run("AuthFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			// Test request parameters
			assert.Equal(tt, auth.AuthEndpoint, req.URL.String())

			// Send mock response
			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",

				// Send response to be tested
				Body: io.NopCloser(bytes.NewBufferString("")),

				// Must be set to non-nil value or it panics
				Header: make(http.Header),
			}
		}

		testAuthClient := auth.NewAuthClient(auth.AuthEndpoint)

		// We replace the transport mechanism by mocking the http request
		// so that the test stays a unit test e.g, no server/auth call.
		testAuthClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		err := testAuthClient.AuthenticateUser(ginContext, auth.AuthEndpoint)

		assert.NoError(tt, err)
	})

	t.Run("AuthNotFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, auth.AuthEndpoint, req.URL.String())

			// error payload
			resp := `{"error":"not found"}`

			return &http.Response{
				StatusCode: 404,
				Status:     "404 NOT FOUND",
				Body:       io.NopCloser(bytes.NewBufferString(resp)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}
		}

		testAuthClient := auth.NewAuthClient(auth.AuthEndpoint)

		testAuthClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		err := testAuthClient.AuthenticateUser(ginContext, auth.AuthEndpoint)

		assert.Error(tt, err)
	})

	t.Run("RequestFailure", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, auth.AuthEndpoint, req.URL.String())

			return nil
		}

		testAuthClient := auth.NewAuthClient(auth.AuthEndpoint)

		testAuthClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		err := testAuthClient.AuthenticateUser(ginContext, auth.AuthEndpoint)

		assert.Error(tt, err)
	})
}
