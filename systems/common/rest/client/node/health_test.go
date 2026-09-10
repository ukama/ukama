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
	testNodeId   = "uk-sa2341-tnode-v0-a1a0"
	testReportId = "7f8b9c0d-1e2f-3a4b-5c6d-7e8f9a0b1c2d"
)

func TestHealthClient_GetInterfaces(t *testing.T) {
	baseURL := "http://test-health-service.com"

	t.Run("InterfacesFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, node.HealthEndpoint+"/nodes/"+testNodeId+"/interfaces", req.URL.Path)
			q := req.URL.Query()
			assert.Equal(tt, testReportId, q.Get("reportId"))
			assert.Empty(tt, q.Get("nodeId"))

			body := `{"interfaces":{"switch":{"state":"active","policy":{"hash":"hash-123","source":"controller"}}}}`

			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",
				Header:     make(http.Header),
				Body:       io.NopCloser(bytes.NewBufferString(body)),
			}
		}

		testHealthClient := node.NewNodeHealthClient(baseURL)
		testHealthClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		ifaces, err := testHealthClient.GetInterfaces(testNodeId, testReportId)

		assert.NoError(tt, err)
		assert.NotNil(tt, ifaces.Switch)
		assert.Equal(tt, "active", ifaces.Switch.State)
		assert.NotNil(tt, ifaces.Switch.Policy)
		assert.Equal(tt, "hash-123", ifaces.Switch.Policy.Hash)
		assert.Equal(tt, "controller", ifaces.Switch.Policy.Source)
	})

	t.Run("GPSWithNotAvailableTime", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			body := `{"interfaces":{"gps":{"available":true,"lock":true,"coordinates":"0.000000,-90.000000","time":"not-available"}}}`

			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",
				Header:     make(http.Header),
				Body:       io.NopCloser(bytes.NewBufferString(body)),
			}
		}

		testHealthClient := node.NewNodeHealthClient(baseURL)
		testHealthClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		ifaces, err := testHealthClient.GetInterfaces(testNodeId, testReportId)

		assert.NoError(tt, err)
		assert.NotNil(tt, ifaces.Gps)
		assert.Equal(tt, "0.000000,-90.000000", ifaces.Gps.Coordinates)
		assert.Equal(tt, "not-available", ifaces.Gps.Time)
	})

	t.Run("InvalidResponse", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, node.HealthEndpoint+"/nodes/"+testNodeId+"/interfaces", req.URL.Path)

			resp := `{"error":"internal server error"}`

			return &http.Response{
				StatusCode: 500,
				Body:       io.NopCloser(bytes.NewBufferString(resp)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}
		}

		testHealthClient := node.NewNodeHealthClient(baseURL)
		testHealthClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		_, err := testHealthClient.GetInterfaces(testNodeId, testReportId)

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

		testHealthClient := node.NewNodeHealthClient(baseURL)
		testHealthClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		_, err := testHealthClient.GetInterfaces(testNodeId, testReportId)

		assert.Error(tt, err)
	})

	t.Run("RequestFailure", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			return nil
		}

		testHealthClient := node.NewNodeHealthClient(baseURL)
		testHealthClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		_, err := testHealthClient.GetInterfaces(testNodeId, testReportId)

		assert.Error(tt, err)
	})

	t.Run("MissingNodeId", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			tt.Fatalf("unexpected request to %s", req.URL.String())

			return nil
		}

		testHealthClient := node.NewNodeHealthClient(baseURL)
		testHealthClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		_, err := testHealthClient.GetInterfaces("", testReportId)

		assert.Error(tt, err)
	})
}
