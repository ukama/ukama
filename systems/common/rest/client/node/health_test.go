/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package node_test

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/tj/assert"

	"github.com/ukama/ukama/systems/common/rest/client"
	"github.com/ukama/ukama/systems/common/rest/client/node"
)

const (
	testNodeId        = "03cb753f-5e03-4c97-8e47-625115476c72"
	testReportId      = "7f8b9c0d-1e2f-3a4b-5c6d-7e8f9a0b1c2d"
	testInterfaceName = "switch"
)

func TestHealthClient_GetInterfaces(t *testing.T) {
	baseURL := "http://test-health-service.com"

	t.Run("InterfacesFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Contains(tt, req.URL.String(), baseURL+node.HealthEndpoint+"/interfaces?")
			q := req.URL.Query()
			assert.Equal(tt, testReportId, q.Get("reportId"))
			assert.Equal(tt, testNodeId, q.Get("nodeId"))
			assert.Equal(tt, testInterfaceName, q.Get("interfaceName"))

			body := `{"interfaces":{"switch":{"state":"active","policy":{"hash":"hash-123","source":"controller"}}}}`

			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",
				Header:     make(http.Header),
				Body:       io.NopCloser(bytes.NewBufferString(body)),
			}
		}

		testHealthClient := node.NewHealthClient(baseURL)
		testHealthClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		ifaces, err := testHealthClient.GetInterfaces(testInterfaceName, testNodeId, testReportId)

		assert.NoError(tt, err)
		assert.NotNil(tt, ifaces.Switch)
		assert.Equal(tt, "active", ifaces.Switch.State)
		assert.NotNil(tt, ifaces.Switch.Policy)
		assert.Equal(tt, "hash-123", ifaces.Switch.Policy.Hash)
		assert.Equal(tt, "controller", ifaces.Switch.Policy.Source)
	})

	t.Run("InvalidResponse", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Contains(tt, req.URL.String(), baseURL+node.HealthEndpoint+"/interfaces?")

			resp := `{"error":"internal server error"}`

			return &http.Response{
				StatusCode: 500,
				Body:       io.NopCloser(bytes.NewBufferString(resp)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}
		}

		testHealthClient := node.NewHealthClient(baseURL)
		testHealthClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		_, err := testHealthClient.GetInterfaces(testInterfaceName, testNodeId, testReportId)

		assert.Error(tt, err)
	})

	t.Run("DeserializationFailure", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			return &http.Response{
				StatusCode: 200,
				Header:     make(http.Header),
				Body:       io.NopCloser(bytes.NewBufferString("not-json")),
			}
		}

		testHealthClient := node.NewHealthClient(baseURL)
		testHealthClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		_, err := testHealthClient.GetInterfaces(testInterfaceName, testNodeId, testReportId)

		assert.Error(tt, err)
	})

	t.Run("RequestFailure", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			return nil
		}

		testHealthClient := node.NewHealthClient(baseURL)
		testHealthClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		_, err := testHealthClient.GetInterfaces(testInterfaceName, testNodeId, testReportId)

		assert.Error(tt, err)
	})
}
