/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package inventory_test

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ukama/ukama/systems/common/rest/client"
	"github.com/ukama/ukama/systems/common/rest/client/inventory"
)

const testUuid = "03cb753f-5e03-4c97-8e47-625115476c72"

func TestComponentClient_Get(t *testing.T) {
	t.Run("ComponentFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			// Test request parameters
			assert.Equal(tt, req.URL.String(), inventory.ComponentEndpoint+"/"+testUuid)

			// fake component info
			comp := `{"component":{"id": "03cb753f-5e03-4c97-8e47-625115476c72", "type": "backhaul"}}`

			// Send mock response
			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",

				// Send response to be tested
				Body: io.NopCloser(bytes.NewBufferString(comp)),

				// Must be set to non-nil value or it panics
				Header: make(http.Header),
			}
		}

		testComponentClient := inventory.NewComponentClient("")

		// We replace the transport mechanism by mocking the http request
		// so that the test stays a unit test e.g, no server/component call.
		testComponentClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		c, err := testComponentClient.Get(testUuid)

		assert.NoError(tt, err)
		assert.Equal(tt, testUuid, c.Id.String())
	})

	t.Run("ComponentNotFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), inventory.ComponentEndpoint+"/"+testUuid)

			// error payload
			resp := `{"error":"not found"}`

			return &http.Response{
				StatusCode: 404,
				Status:     "404 NOT FOUND",
				Body:       io.NopCloser(bytes.NewBufferString(resp)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}
		}

		testComponentClient := inventory.NewComponentClient("")

		testComponentClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		c, err := testComponentClient.Get(testUuid)

		assert.Error(tt, err)
		assert.Nil(tt, c)
	})

	t.Run("InvalidResponsePayload", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), inventory.ComponentEndpoint+"/"+testUuid)

			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",
				Body:       io.NopCloser(bytes.NewBufferString(`OK`)),
				Header:     make(http.Header),
			}
		}

		testComponentClient := inventory.NewComponentClient("")

		testComponentClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		c, err := testComponentClient.Get(testUuid)

		assert.Error(tt, err)
		assert.Nil(tt, c)
	})

	t.Run("RequestFailure", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), inventory.ComponentEndpoint+"/"+testUuid)

			return nil
		}

		testComponentClient := inventory.NewComponentClient("")

		testComponentClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		c, err := testComponentClient.Get(testUuid)

		assert.Error(tt, err)
		assert.Nil(tt, c)
	})
}

func TestComponentClient_List(t *testing.T) {
	t.Run("ComponentsFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			// Test request parameters
			expectedURL := inventory.ComponentEndpoint + "?id=" + testUuid + "&user_id=user123&part_number=PN001&category=backhaul"
			assert.Equal(tt, expectedURL, req.URL.String())

			// fake components list
			comps := `{"components":[{"id": "03cb753f-5e03-4c97-8e47-625115476c72", "type": "backhaul", "part_number": "PN001"}, {"id": "04cb753f-5e03-4c97-8e47-625115476c73", "type": "backhaul", "part_number": "PN002"}]}`

			// Send mock response
			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",

				// Send response to be tested
				Body: io.NopCloser(bytes.NewBufferString(comps)),

				// Must be set to non-nil value or it panics
				Header: make(http.Header),
			}
		}

		testComponentClient := inventory.NewComponentClient("")

		// We replace the transport mechanism by mocking the http request
		// so that the test stays a unit test e.g, no server/component call.
		testComponentClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		components, err := testComponentClient.List(testUuid, "user123", "PN001", "backhaul")

		assert.NoError(tt, err)
		assert.NotNil(tt, components)
		assert.Len(tt, components.Components, 2)
		assert.Equal(tt, testUuid, components.Components[0].Id.String())
		assert.Equal(tt, "backhaul", components.Components[0].Type)
		assert.Equal(tt, "PN001", components.Components[0].PartNumber)
	})

	t.Run("EmptyComponentsList", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			// Test request parameters with empty values
			expectedURL := inventory.ComponentEndpoint + "?id=&user_id=&part_number=&category="
			assert.Equal(tt, expectedURL, req.URL.String())

			// empty components list
			comps := `{"components":[]}`

			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",
				Body:       io.NopCloser(bytes.NewBufferString(comps)),
				Header:     make(http.Header),
			}
		}

		testComponentClient := inventory.NewComponentClient("")

		testComponentClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		components, err := testComponentClient.List("", "", "", "")

		assert.NoError(tt, err)
		assert.NotNil(tt, components)
		assert.Len(tt, components.Components, 0)
	})

	t.Run("ComponentsNotFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			expectedURL := inventory.ComponentEndpoint + "?id=invalid&user_id=user123&part_number=PN001&category=backhaul"
			assert.Equal(tt, expectedURL, req.URL.String())

			// error payload
			resp := `{"error":"not found"}`

			return &http.Response{
				StatusCode: 404,
				Status:     "404 NOT FOUND",
				Body:       io.NopCloser(bytes.NewBufferString(resp)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}
		}

		testComponentClient := inventory.NewComponentClient("")

		testComponentClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		components, err := testComponentClient.List("invalid", "user123", "PN001", "backhaul")

		assert.Error(tt, err)
		assert.Nil(tt, components)
	})

	t.Run("InvalidResponsePayload", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			expectedURL := inventory.ComponentEndpoint + "?id=&user_id=&part_number=&category="
			assert.Equal(tt, expectedURL, req.URL.String())

			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",
				Body:       io.NopCloser(bytes.NewBufferString(`OK`)),
				Header:     make(http.Header),
			}
		}

		testComponentClient := inventory.NewComponentClient("")

		testComponentClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		components, err := testComponentClient.List("", "", "", "")

		assert.Error(tt, err)
		assert.Nil(tt, components)
	})

	t.Run("RequestFailure", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			expectedURL := inventory.ComponentEndpoint + "?id=&user_id=&part_number=&category="
			assert.Equal(tt, expectedURL, req.URL.String())

			return nil
		}

		testComponentClient := inventory.NewComponentClient("")

		testComponentClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		components, err := testComponentClient.List("", "", "", "")

		assert.Error(tt, err)
		assert.Nil(tt, components)
	})

	t.Run("PartialParameters", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			// Test with only some parameters provided
			expectedURL := inventory.ComponentEndpoint + "?id=&user_id=user123&part_number=&category=backhaul"
			assert.Equal(tt, expectedURL, req.URL.String())

			// single component response
			comps := `{"components":[{"id": "03cb753f-5e03-4c97-8e47-625115476c72", "type": "backhaul", "user_id": "user123"}]}`

			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",
				Body:       io.NopCloser(bytes.NewBufferString(comps)),
				Header:     make(http.Header),
			}
		}

		testComponentClient := inventory.NewComponentClient("")

		testComponentClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		components, err := testComponentClient.List("", "user123", "", "backhaul")

		assert.NoError(tt, err)
		assert.NotNil(tt, components)
		assert.Len(tt, components.Components, 1)
		assert.Equal(tt, testUuid, components.Components[0].Id.String())
		assert.Equal(tt, "user123", components.Components[0].UserId)
	})
}
