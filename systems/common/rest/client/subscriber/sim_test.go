/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package subscriber_test

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/tj/assert"
	"github.com/ukama/ukama/systems/common/rest/client"
	"github.com/ukama/ukama/systems/common/rest/client/subscriber"
)

const (
	testUuid = "03cb753f-5e03-4c97-8e47-625115476c72"
)

func TestSimClient_Get(t *testing.T) {
	t.Run("SimFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			// Test request parameters
			assert.Equal(tt, req.URL.String(), subscriber.SimEndpoint+"/"+testUuid)

			// fake sim info
			sim := `{"sim":{"id": "03cb753f-5e03-4c97-8e47-625115476c72", "is_physical": false}}`

			// Send mock response
			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",

				// Send response to be tested
				Body: io.NopCloser(bytes.NewBufferString(sim)),

				// Must be set to non-nil value or it panics
				Header: make(http.Header),
			}
		}

		testSimClient := subscriber.NewSimClient("")

		// We replace the transport mechanism by mocking the http request
		// so that the test stays a unit test e.g no server/network call.
		testSimClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		s, err := testSimClient.Get(testUuid)

		assert.NoError(tt, err)
		assert.Equal(tt, testUuid, s.Id)
	})

	t.Run("SimNotFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), subscriber.SimEndpoint+"/"+testUuid)

			// error payload
			resp := `{"error":"not found"}`

			return &http.Response{
				StatusCode: 404,
				Status:     "404 NOT FOUND",
				Body:       io.NopCloser(bytes.NewBufferString(resp)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}
		}

		testSimClient := subscriber.NewSimClient("")

		testSimClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		s, err := testSimClient.Get(testUuid)

		assert.Error(tt, err)
		assert.Nil(tt, s)
	})

	t.Run("InvalidResponsePayload", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), subscriber.SimEndpoint+"/"+testUuid)

			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",
				Body:       io.NopCloser(bytes.NewBufferString(`OK`)),
				Header:     make(http.Header),
			}
		}

		testSimClient := subscriber.NewSimClient("")

		testSimClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		s, err := testSimClient.Get(testUuid)

		assert.Error(tt, err)
		assert.Nil(tt, s)
	})

	t.Run("RequestFailure", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), subscriber.SimEndpoint+"/"+testUuid)

			return nil
		}

		testSimClient := subscriber.NewSimClient("")

		testSimClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		s, err := testSimClient.Get(testUuid)

		assert.Error(tt, err)
		assert.Nil(tt, s)
	})
}

func TestSimClient_Add(t *testing.T) {
	t.Run("SimAdded", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			// Test request parameters
			assert.Equal(tt, req.URL.String(), subscriber.SimEndpoint)

			// fake sim info
			sim := `{"sim":{"id": "03cb753f-5e03-4c97-8e47-625115476c72", "is_physical": false}}`

			// Send mock response
			return &http.Response{
				StatusCode: 201,
				Status:     "201 CREATED",

				// Send response to be tested
				Body: io.NopCloser(bytes.NewBufferString(sim)),

				// Must be set to non-nil value or it panics
				Header: make(http.Header),
			}
		}

		testSimClient := subscriber.NewSimClient("")

		// We replace the transport mechanism by mocking the http request
		// so that the test stays a unit test e.g no server/network call.
		testSimClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		s, err := testSimClient.Add(
			subscriber.AddSimRequest{
				SubscriberId: "some-subscriber_Id",
				PackageId:    "some-package_id"},
		)

		assert.NoError(tt, err)
		assert.Equal(tt, testUuid, s.Id)
	})

	t.Run("InvalidResponseHeader", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), subscriber.SimEndpoint)

			// error payload
			resp := `{"error":"internal server error"}`

			return &http.Response{
				StatusCode: 500,
				Status:     "500 INTERNAL SERVER ERROR",
				Body:       io.NopCloser(bytes.NewBufferString(resp)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}
		}

		testSimClient := subscriber.NewSimClient("")

		testSimClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		s, err := testSimClient.Add(
			subscriber.AddSimRequest{
				SubscriberId: "some-subscriber_Id",
				PackageId:    "some-package_id"},
		)

		assert.Error(tt, err)
		assert.Nil(tt, s)
	})

	t.Run("InvalidResponsePayload", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), subscriber.SimEndpoint)

			return &http.Response{
				StatusCode: 201,
				Status:     "201 CREATED",
				Body:       io.NopCloser(bytes.NewBufferString(`CREATED`)),
				Header:     make(http.Header),
			}
		}

		testSimClient := subscriber.NewSimClient("")

		testSimClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		s, err := testSimClient.Add(
			subscriber.AddSimRequest{
				SubscriberId: "some-subscriber_Id",
				PackageId:    "some-package_id"},
		)

		assert.Error(tt, err)
		assert.Nil(tt, s)
	})

	t.Run("RequestFailure", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), subscriber.SimEndpoint)

			return nil
		}

		testSimClient := subscriber.NewSimClient("")

		testSimClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		s, err := testSimClient.Add(
			subscriber.AddSimRequest{
				SubscriberId: "some-subscriber_Id",
				PackageId:    "some-package_id"},
		)

		assert.Error(tt, err)
		assert.Nil(tt, s)
	})
}
