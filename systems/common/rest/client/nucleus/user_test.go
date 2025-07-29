/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.user/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package nucleus_test

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/tj/assert"

	"github.com/ukama/ukama/systems/common/rest/client/nucleus"
)

func TestUserClient_GetById(t *testing.T) {
	t.Run("UserFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			// Test request parameters
			assert.Equal(tt, req.URL.String(), nucleus.UserEndpoint+"/"+testUuid)

			// fake user info
			user := `{"user":{"id": "03cb753f-5e03-4c97-8e47-625115476c72", "email": "john@example.com"}}`

			// Send mock response
			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",

				// Send response to be tested
				Body: io.NopCloser(bytes.NewBufferString(user)),

				// Must be set to non-nil value or it panics
				Header: make(http.Header),
			}
		}

		testUserClient := nucleus.NewUserClient("")

		// We replace the transport mechanism by mocking the http request
		// so that the test stays a unit test e.g, no server/user call.
		testUserClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		usr, err := testUserClient.GetById(testUuid)

		assert.NoError(tt, err)
		assert.Equal(tt, testUuid, usr.Id)
	})

	t.Run("UserNotFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), nucleus.UserEndpoint+"/"+testUuid)

			// error payload
			resp := `{"error":"not found"}`

			return &http.Response{
				StatusCode: 404,
				Status:     "404 NOT FOUND",
				Body:       io.NopCloser(bytes.NewBufferString(resp)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}
		}

		testUserClient := nucleus.NewUserClient("")

		testUserClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		usr, err := testUserClient.GetById(testUuid)

		assert.Error(tt, err)
		assert.Nil(tt, usr)
	})

	t.Run("InvalidResponsePayload", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), nucleus.UserEndpoint+"/"+testUuid)

			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",
				Body:       io.NopCloser(bytes.NewBufferString(`OK`)),
				Header:     make(http.Header),
			}
		}

		testUserClient := nucleus.NewUserClient("")

		testUserClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		usr, err := testUserClient.GetById(testUuid)

		assert.Error(tt, err)
		assert.Nil(tt, usr)
	})

	t.Run("RequestFailure", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), nucleus.UserEndpoint+"/"+testUuid)

			return nil
		}

		testUserClient := nucleus.NewUserClient("")

		testUserClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		usr, err := testUserClient.GetById(testUuid)

		assert.Error(tt, err)
		assert.Nil(tt, usr)
	})
}

func TestUserClient_GetByEmail(t *testing.T) {
	const testEmail = "john@example.com"

	t.Run("UserFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			// Test request parameters
			assert.Equal(tt, req.URL.String(), nucleus.UserEndpoint+"/email/"+testEmail)

			// fake user info
			user := `{"user":{"id": "03cb753f-5e03-4c97-8e47-625115476c72", "email": "john@example.com"}}`

			// Send mock response
			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",

				// Send response to be tested
				Body: io.NopCloser(bytes.NewBufferString(user)),

				// Must be set to non-nil value or it panics
				Header: make(http.Header),
			}
		}

		testUserClient := nucleus.NewUserClient("")

		// We replace the transport mechanism by mocking the http request
		// so that the test stays a unit test e.g, no server/user call.
		testUserClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		usr, err := testUserClient.GetByEmail(testEmail)

		assert.NoError(tt, err)
		assert.Equal(tt, testEmail, usr.Email)
	})

	t.Run("UserNotFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), nucleus.UserEndpoint+"/email/"+testEmail)

			// error payload
			resp := `{"error":"not found"}`

			return &http.Response{
				StatusCode: 404,
				Status:     "404 NOT FOUND",
				Body:       io.NopCloser(bytes.NewBufferString(resp)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}
		}

		testUserClient := nucleus.NewUserClient("")

		testUserClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		usr, err := testUserClient.GetByEmail(testEmail)

		assert.Error(tt, err)
		assert.Nil(tt, usr)
	})

	t.Run("InvalidResponsePayload", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), nucleus.UserEndpoint+"/email/"+testEmail)

			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",
				Body:       io.NopCloser(bytes.NewBufferString(`OK`)),
				Header:     make(http.Header),
			}
		}

		testUserClient := nucleus.NewUserClient("")

		testUserClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		usr, err := testUserClient.GetByEmail(testEmail)

		assert.Error(tt, err)
		assert.Nil(tt, usr)
	})

	t.Run("RequestFailure", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), nucleus.UserEndpoint+"/email/"+testEmail)

			return nil
		}

		testUserClient := nucleus.NewUserClient("")

		testUserClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		usr, err := testUserClient.GetByEmail(testEmail)

		assert.Error(tt, err)
		assert.Nil(tt, usr)
	})
}
