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

	"github.com/ukama/ukama/systems/common/rest/client/registry"
)

func TestSiteClient_Get(t *testing.T) {
	t.Run("SiteFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			// Test request parameters
			assert.Equal(tt, req.URL.String(), registry.SiteEndpoint+"/"+testUuid)

			// fake site info
			site := `{"site":{"id": "03cb753f-5e03-4c97-8e47-625115476c72"}}`

			// Send mock response
			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",

				// Send response to be tested
				Body: io.NopCloser(bytes.NewBufferString(site)),

				// Must be set to non-nil value or it panics
				Header: make(http.Header),
			}
		}

		testNetworkClient := registry.NewSiteClient("")

		// We replace the transport mechanism by mocking the http request
		// so that the test stays a unit test e.g no server/network call.
		testNetworkClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		n, err := testNetworkClient.Get(testUuid)

		assert.NoError(tt, err)
		assert.Equal(tt, testUuid, n.Id)
	})

	t.Run("SiteNotFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), registry.SiteEndpoint+"/"+testUuid)

			// error payload
			resp := `{"error":"not found"}`

			return &http.Response{
				StatusCode: 404,
				Status:     "404 NOT FOUND",
				Body:       io.NopCloser(bytes.NewBufferString(resp)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}
		}

		testNetworkClient := registry.NewSiteClient("")

		testNetworkClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		n, err := testNetworkClient.Get(testUuid)

		assert.Error(tt, err)
		assert.Nil(tt, n)
	})
}
