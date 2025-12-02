/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2025-present, Ukama Inc.
 */

package factory_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tj/assert"

	"github.com/ukama/ukama/systems/common/rest/client"
	"github.com/ukama/ukama/systems/common/rest/client/factory"
)

const (
	testNodeId = "test-node-123"
	testHost   = "http://test-host:8080"
)

func TestMain(m *testing.M) {
	// Suppress log output during tests
	logrus.SetLevel(logrus.PanicLevel)
	os.Exit(m.Run())
}

func TestNewNodeFactoryClient(t *testing.T) {
	t.Run("ValidHost", func(tt *testing.T) {
		client := factory.NewNodeFactoryClient(testHost)

		assert.NotNil(tt, client)
		assert.NotNil(tt, client.R)
	})

	t.Run("WithOptions", func(tt *testing.T) {
		client := factory.NewNodeFactoryClient(testHost, client.WithDebug(true))

		assert.NotNil(tt, client)
		assert.NotNil(tt, client.R)
	})
}

func TestNodeFactoryClient_Get(t *testing.T) {
	t.Run("Success", func(tt *testing.T) {
		expectedNode := factory.NodeFactoryInfo{
			Id:            testNodeId,
			Type:          "tnode",
			OrgName:       "test-org",
			IsProvisioned: true,
			ProvisionedAt: time.Now().UTC(),
		}

		mockTransport := func(req *http.Request) *http.Response {
			// Test request parameters
			expectedURL := testHost + factory.FactoryEndpoint + "/node/" + testNodeId
			assert.Equal(tt, expectedURL, req.URL.String())
			assert.Equal(tt, "GET", req.Method)

			// Serialize expected response
			responseBody, _ := json.Marshal(expectedNode)

			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",
				Body:       io.NopCloser(bytes.NewBuffer(responseBody)),
				Header:     make(http.Header),
			}
		}

		testClient := factory.NewNodeFactoryClient(testHost)
		testClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		result, err := testClient.Get(testNodeId)

		assert.NoError(tt, err)
		assert.NotNil(tt, result)
		assert.Equal(tt, expectedNode.Id, result.Id)
		assert.Equal(tt, expectedNode.Type, result.Type)
		assert.Equal(tt, expectedNode.OrgName, result.OrgName)
		assert.Equal(tt, expectedNode.IsProvisioned, result.IsProvisioned)
		assert.True(tt, result.ProvisionedAt.Equal(expectedNode.ProvisionedAt))
	})

	t.Run("NodeNotFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			expectedURL := testHost + factory.FactoryEndpoint + "/node/" + testNodeId
			assert.Equal(tt, expectedURL, req.URL.String())

			errorResponse := `{"error": "node not found"}`

			return &http.Response{
				StatusCode: 404,
				Status:     "404 NOT FOUND",
				Body:       io.NopCloser(bytes.NewBufferString(errorResponse)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}
		}

		testClient := factory.NewNodeFactoryClient(testHost)
		testClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		result, err := testClient.Get(testNodeId)

		assert.Error(tt, err)
		assert.Nil(tt, result)
	})

	t.Run("InvalidResponsePayload", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			expectedURL := testHost + factory.FactoryEndpoint + "/node/" + testNodeId
			assert.Equal(tt, expectedURL, req.URL.String())

			// Return invalid JSON
			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",
				Body:       io.NopCloser(bytes.NewBufferString(`{"invalid": json`)),
				Header:     make(http.Header),
			}
		}

		testClient := factory.NewNodeFactoryClient(testHost)
		testClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		result, err := testClient.Get(testNodeId)

		assert.Error(tt, err)
		assert.Nil(tt, result)
		assert.Contains(tt, err.Error(), "node info deserialization failure")
	})

	t.Run("RequestFailure", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			expectedURL := testHost + factory.FactoryEndpoint + "/node/" + testNodeId
			assert.Equal(tt, expectedURL, req.URL.String())

			// Return nil to simulate network failure
			return nil
		}

		testClient := factory.NewNodeFactoryClient(testHost)
		testClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		result, err := testClient.Get(testNodeId)

		assert.Error(tt, err)
		assert.Nil(tt, result)
	})

	t.Run("EmptyResponseBody", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			expectedURL := testHost + factory.FactoryEndpoint + "/node/" + testNodeId
			assert.Equal(tt, expectedURL, req.URL.String())

			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",
				Body:       io.NopCloser(bytes.NewBufferString("")),
				Header:     make(http.Header),
			}
		}

		testClient := factory.NewNodeFactoryClient(testHost)
		testClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		result, err := testClient.Get(testNodeId)

		assert.Error(tt, err)
		assert.Nil(tt, result)
		assert.Contains(tt, err.Error(), "node info deserialization failure")
	})

}

func TestNodeFactoryClient_List(t *testing.T) {
	t.Run("ListNodesSuccess", func(tt *testing.T) {
		expectedNodes := factory.Nodes{
			Nodes: []*factory.NodeFactoryInfo{
				{
					Id:            testNodeId,
					Type:          "tnode",
					OrgName:       "test-org",
					IsProvisioned: true,
					ProvisionedAt: time.Now().UTC(),
				},
			},
		}

		mockTransport := func(req *http.Request) *http.Response {
			expectedURL := testHost + factory.FactoryEndpoint + "/nodes?type=tnode&orgName=test-org&isProvisioned=true"
			assert.Equal(tt, expectedURL, req.URL.String())

			// Serialize expected response
			responseBody, _ := json.Marshal(expectedNodes)

			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",
				Body:       io.NopCloser(bytes.NewBuffer(responseBody)),
				Header:     make(http.Header),
			}
		}

		testClient := factory.NewNodeFactoryClient(testHost)
		testClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		result, err := testClient.List("tnode", "test-org", true)

		assert.NoError(tt, err)
		assert.NotNil(tt, result)
		assert.Equal(tt, expectedNodes.Nodes[0].Id, result.Nodes[0].Id)
		assert.Equal(tt, expectedNodes.Nodes[0].Type, result.Nodes[0].Type)
		assert.Equal(tt, expectedNodes.Nodes[0].OrgName, result.Nodes[0].OrgName)
		assert.Equal(tt, expectedNodes.Nodes[0].IsProvisioned, result.Nodes[0].IsProvisioned)
		assert.True(tt, result.Nodes[0].ProvisionedAt.Equal(expectedNodes.Nodes[0].ProvisionedAt))
	})

	t.Run("ListNodesNotFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			expectedURL := testHost + factory.FactoryEndpoint + "/nodes?type=tnode&orgName=test-org&isProvisioned=true"
			assert.Equal(tt, expectedURL, req.URL.String())

			return &http.Response{
				StatusCode: 404,
				Status:     "404 NOT FOUND",
				Body:       io.NopCloser(bytes.NewBufferString(`{"error": "nodes not found"}`)),
				Header:     make(http.Header),
			}
		}

		testClient := factory.NewNodeFactoryClient(testHost)
		testClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		result, err := testClient.List("tnode", "test-org", true)

		assert.Error(tt, err)
		assert.Nil(tt, result)
	})
}
