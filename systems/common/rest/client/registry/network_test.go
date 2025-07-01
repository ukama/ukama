/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package registry_test

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/tj/assert"

	"github.com/ukama/ukama/systems/common/rest/client"
	"github.com/ukama/ukama/systems/common/rest/client/registry"
)

const testUuid = "03cb753f-5e03-4c97-8e47-625115476c72"

func TestNetworkClient_Get(t *testing.T) {
	t.Run("NetworkFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			// Test request parameters
			assert.Equal(tt, req.URL.String(), registry.NetworkEndpoint+"/"+testUuid)

			// fake network info
			ntwk := `{"network":{"id": "03cb753f-5e03-4c97-8e47-625115476c72", "is_deactivated": true}}`

			// Send mock response
			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",

				// Send response to be tested
				Body: io.NopCloser(bytes.NewBufferString(ntwk)),

				// Must be set to non-nil value or it panics
				Header: make(http.Header),
			}
		}

		testNetworkClient := registry.NewNetworkClient("")

		// We replace the transport mechanism by mocking the http request
		// so that the test stays a unit test e.g, no server/network call.
		testNetworkClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		n, err := testNetworkClient.Get(testUuid)

		assert.NoError(tt, err)
		assert.Equal(tt, testUuid, n.Id)
	})

	t.Run("NetworkNotFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), registry.NetworkEndpoint+"/"+testUuid)

			// error payload
			resp := `{"error":"not found"}`

			return &http.Response{
				StatusCode: 404,
				Status:     "404 NOT FOUND",
				Body:       io.NopCloser(bytes.NewBufferString(resp)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}
		}

		testNetworkClient := registry.NewNetworkClient("")

		testNetworkClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		n, err := testNetworkClient.Get(testUuid)

		assert.Error(tt, err)
		assert.Nil(tt, n)
	})

	t.Run("InvalidResponsePayload", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), registry.NetworkEndpoint+"/"+testUuid)

			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",
				Body:       io.NopCloser(bytes.NewBufferString(`OK`)),
				Header:     make(http.Header),
			}
		}

		testNetworkClient := registry.NewNetworkClient("")

		testNetworkClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		n, err := testNetworkClient.Get(testUuid)

		assert.Error(tt, err)
		assert.Nil(tt, n)
	})

	t.Run("RequestFailure", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), registry.NetworkEndpoint+"/"+testUuid)

			return nil
		}

		testNetworkClient := registry.NewNetworkClient("")

		testNetworkClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		n, err := testNetworkClient.Get(testUuid)

		assert.Error(tt, err)
		assert.Nil(tt, n)
	})
}

func TestNetworkClient_GetDefault(t *testing.T) {
	t.Run("DefaultNetworkFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			// Test request parameters
			assert.Equal(tt, req.URL.String(), registry.NetworkEndpoint+"/default")

			// fake network info
			ntwk := `{"network":{"id": "03cb753f-5e03-4c97-8e47-625115476c72", "is_deactivated": true}}`

			// Send mock response
			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",

				// Send response to be tested
				Body: io.NopCloser(bytes.NewBufferString(ntwk)),

				// Must be set to non-nil value or it panics
				Header: make(http.Header),
			}
		}

		testNetworkClient := registry.NewNetworkClient("")

		// We replace the transport mechanism by mocking the http request
		// so that the test stays a unit test e.g, no server/network call.
		testNetworkClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		n, err := testNetworkClient.GetDefault()

		assert.NoError(tt, err)
		assert.Equal(tt, testUuid, n.Id)
	})

	t.Run("DefaultNetworkNotFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), registry.NetworkEndpoint+"/default")

			// error payload
			resp := `{"error":"not found"}`

			return &http.Response{
				StatusCode: 404,
				Status:     "404 NOT FOUND",
				Body:       io.NopCloser(bytes.NewBufferString(resp)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}
		}

		testNetworkClient := registry.NewNetworkClient("")

		testNetworkClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		n, err := testNetworkClient.GetDefault()

		assert.Error(tt, err)
		assert.Nil(tt, n)
	})

	t.Run("InvalidResponsePayload", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), registry.NetworkEndpoint+"/default")

			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",
				Body:       io.NopCloser(bytes.NewBufferString(`OK`)),
				Header:     make(http.Header),
			}
		}

		testNetworkClient := registry.NewNetworkClient("")

		testNetworkClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		n, err := testNetworkClient.GetDefault()

		assert.Error(tt, err)
		assert.Nil(tt, n)
	})

	t.Run("RequestFailure", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), registry.NetworkEndpoint+"/default")

			return nil
		}

		testNetworkClient := registry.NewNetworkClient("")

		testNetworkClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		n, err := testNetworkClient.GetDefault()

		assert.Error(tt, err)
		assert.Nil(tt, n)
	})
}

func TestNetworkClient_Add(t *testing.T) {
	t.Run("NetworkAdded", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			// Test request parameters
			assert.Equal(tt, req.URL.String(), registry.NetworkEndpoint)

			// fake network info
			ntwk := `{"network":{"id": "03cb753f-5e03-4c97-8e47-625115476c72","Name": "net-1", "is_deactivated": true}}`

			// Send mock response
			return &http.Response{
				StatusCode: 201,
				Status:     "201 CREATED",

				// Send response to be tested
				Body: io.NopCloser(bytes.NewBufferString(ntwk)),

				// Must be set to non-nil value or it panics
				Header: make(http.Header),
			}
		}

		testNetworkClient := registry.NewNetworkClient("")

		// We replace the transport mechanism by mocking the http request
		// so that the test stays a unit test e.g, no server/network call.
		testNetworkClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		n, err := testNetworkClient.Add(
			registry.AddNetworkRequest{NetName: "net-1",
				OrgName: "Ukama"},
		)

		assert.NoError(tt, err)
		assert.Equal(tt, testUuid, n.Id)
	})

	t.Run("InvalidResponseHeader", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), registry.NetworkEndpoint)

			// error payload
			resp := `{"error":"internal server error"}`

			return &http.Response{
				StatusCode: 500,
				Status:     "500 INTERNAL SERVER ERROR",
				Body:       io.NopCloser(bytes.NewBufferString(resp)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}
		}

		testNetworkClient := registry.NewNetworkClient("")

		testNetworkClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		n, err := testNetworkClient.Add(
			registry.AddNetworkRequest{
				NetName: "net-1",
				OrgName: "Ukama",
			},
		)

		assert.Error(tt, err)
		assert.Nil(tt, n)
	})

	t.Run("InvalidResponsePayload", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), registry.NetworkEndpoint)

			return &http.Response{
				StatusCode: 201,
				Status:     "201 CREATED",
				Body:       io.NopCloser(bytes.NewBufferString(`CREATED`)),
				Header:     make(http.Header),
			}
		}

		testNetworkClient := registry.NewNetworkClient("")

		testNetworkClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		n, err := testNetworkClient.Add(
			registry.AddNetworkRequest{
				NetName: "net-1",
				OrgName: "Ukama",
			},
		)

		assert.Error(tt, err)
		assert.Nil(tt, n)
	})

	t.Run("RequestFailure", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), registry.NetworkEndpoint)

			return nil
		}

		testNetworkClient := registry.NewNetworkClient("")

		testNetworkClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		n, err := testNetworkClient.Add(
			registry.AddNetworkRequest{
				NetName: "net-1",
				OrgName: "Ukama",
			},
		)

		assert.Error(tt, err)
		assert.Nil(tt, n)
	})
}
