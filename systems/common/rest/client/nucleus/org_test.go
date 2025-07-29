/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
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

const testUuid = "03cb753f-5e03-4c97-8e47-625115476c72"

func TestOrgClient_Get(t *testing.T) {
	t.Run("OrgFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			// Test request parameters
			assert.Equal(tt, req.URL.String(), nucleus.OrgEndpoint+"/"+testUuid)

			// fake org info
			org := `{"org":{"id": "03cb753f-5e03-4c97-8e47-625115476c72", "is_deactivated": true}}`

			// Send mock response
			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",

				// Send response to be tested
				Body: io.NopCloser(bytes.NewBufferString(org)),

				// Must be set to non-nil value or it panics
				Header: make(http.Header),
			}
		}

		testOrgClient := nucleus.NewOrgClient("")

		// We replace the transport mechanism by mocking the http request
		// so that the test stays a unit test e.g, no server/org call.
		testOrgClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		o, err := testOrgClient.Get(testUuid)

		assert.NoError(tt, err)
		assert.Equal(tt, testUuid, o.Id)
	})

	t.Run("OrgNotFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), nucleus.OrgEndpoint+"/"+testUuid)

			// error payload
			resp := `{"error":"not found"}`

			return &http.Response{
				StatusCode: 404,
				Status:     "404 NOT FOUND",
				Body:       io.NopCloser(bytes.NewBufferString(resp)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}

		}

		testOrgClient := nucleus.NewOrgClient("")

		testOrgClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		o, err := testOrgClient.Get(testUuid)

		assert.Error(tt, err)
		assert.Nil(tt, o)
	})

	t.Run("InvalidResponsePayload", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), nucleus.OrgEndpoint+"/"+testUuid)

			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",
				Body:       io.NopCloser(bytes.NewBufferString(`OK`)),
				Header:     make(http.Header),
			}
		}

		testOrgClient := nucleus.NewOrgClient("")

		testOrgClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		o, err := testOrgClient.Get(testUuid)

		assert.Error(tt, err)
		assert.Nil(tt, o)
	})

	t.Run("RequestFailure", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), nucleus.OrgEndpoint+"/"+testUuid)

			return nil
		}

		testOrgClient := nucleus.NewOrgClient("")

		testOrgClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		o, err := testOrgClient.Get(testUuid)

		assert.Error(tt, err)
		assert.Nil(tt, o)
	})
}

func TestOrgClient_AddUser(t *testing.T) {
	t.Run("OrgAndUserFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), nucleus.OrgEndpoint+"/"+testUuid+"/users/"+testUuid)

			return &http.Response{
				StatusCode: 201,
				Header:     make(http.Header),
			}
		}

		testOrgClient := nucleus.NewOrgClient("")

		testOrgClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		err := testOrgClient.AddUser(testUuid, testUuid)

		assert.NoError(tt, err)
	})

	t.Run("OrgOrUserNotFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), nucleus.OrgEndpoint+"/"+testUuid+"/users/"+testUuid)

			// error payload
			resp := `{"error":"not found"}`

			return &http.Response{
				StatusCode: 404,
				Status:     "404 NOT FOUND",
				Body:       io.NopCloser(bytes.NewBufferString(resp)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}
		}

		testOrgClient := nucleus.NewOrgClient("")

		testOrgClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		err := testOrgClient.AddUser(testUuid, testUuid)

		assert.Error(tt, err)
	})

	t.Run("RequestFailure", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), nucleus.OrgEndpoint+"/"+testUuid+"/users/"+testUuid)

			return nil
		}

		testOrgClient := nucleus.NewOrgClient("")

		testOrgClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		err := testOrgClient.AddUser(testUuid, testUuid)

		assert.Error(tt, err)
	})
}

func TestOrgClient_RemoveUser(t *testing.T) {
	t.Run("OrgAndUserFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), nucleus.OrgEndpoint+"/"+testUuid+"/users/"+testUuid)

			return &http.Response{
				StatusCode: 200,
				Header:     make(http.Header),
			}
		}

		testOrgClient := nucleus.NewOrgClient("")

		testOrgClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		err := testOrgClient.RemoveUser(testUuid, testUuid)

		assert.NoError(tt, err)
	})

	t.Run("OrgOrUserNotFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), nucleus.OrgEndpoint+"/"+testUuid+"/users/"+testUuid)

			// error payload
			resp := `{"error":"not found"}`

			return &http.Response{
				StatusCode: 404,
				Status:     "404 NOT FOUND",
				Body:       io.NopCloser(bytes.NewBufferString(resp)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}
		}

		testOrgClient := nucleus.NewOrgClient("")

		testOrgClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		err := testOrgClient.RemoveUser(testUuid, testUuid)

		assert.Error(tt, err)
	})

	t.Run("RequestFailure", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), nucleus.OrgEndpoint+"/"+testUuid+"/users/"+testUuid)

			return nil
		}

		testOrgClient := nucleus.NewOrgClient("")

		testOrgClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		err := testOrgClient.RemoveUser(testUuid, testUuid)

		assert.Error(tt, err)
	})
}

type RoundTripFunc func(req *http.Request) *http.Response

func (r RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return r(req), nil
}
