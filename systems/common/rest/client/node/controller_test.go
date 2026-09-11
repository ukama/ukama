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
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/tj/assert"

	"github.com/ukama/ukama/systems/common/rest/client"
	"github.com/ukama/ukama/systems/common/rest/client/node"
)

const (
	testOperationId = "03cb753f-5e03-4c97-8e47-625115476c72"
	testResourceKey = "node:" + testNodeId
)

const testOperationBody = `{"operation_id":"` + testOperationId + `","resource_key":"` + testResourceKey + `","status":"SUCCESS"}`

func TestControllerClient_RestartNode(t *testing.T) {
	baseURL := "http://test-controller-service.com"

	t.Run("Success", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, http.MethodPost, req.Method)
			assert.Equal(tt, node.ControllerEndpoint+"/nodes/"+testNodeId+"/restart", req.URL.Path)

			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",
				Header:     make(http.Header),
				Body:       io.NopCloser(bytes.NewBufferString(testOperationBody)),
			}
		}

		testControllerClient := node.NewNodeControllerClient(baseURL)
		testControllerClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		resp, err := testControllerClient.RestartNode(testNodeId)

		assert.NoError(tt, err)
		assert.Equal(tt, testOperationId, resp.OperationId)
		assert.Equal(tt, testResourceKey, resp.ResourceKey)
		assert.Equal(tt, "SUCCESS", resp.Status)
	})

	t.Run("InvalidResponse", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			resp := `{"error":"internal server error"}`

			return &http.Response{
				StatusCode: 500,
				Body:       io.NopCloser(bytes.NewBufferString(resp)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}
		}

		testControllerClient := node.NewNodeControllerClient(baseURL)
		testControllerClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		_, err := testControllerClient.RestartNode(testNodeId)

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

		testControllerClient := node.NewNodeControllerClient(baseURL)
		testControllerClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		_, err := testControllerClient.RestartNode(testNodeId)

		assert.Error(tt, err)
	})

	t.Run("RequestFailure", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			return nil
		}

		testControllerClient := node.NewNodeControllerClient(baseURL)
		testControllerClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		_, err := testControllerClient.RestartNode(testNodeId)

		assert.Error(tt, err)
	})
}

func TestControllerClient_ToggleSwitchPort(t *testing.T) {
	baseURL := "http://test-controller-service.com"

	t.Run("Success", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, http.MethodPost, req.Method)
			assert.Equal(tt, node.ControllerEndpoint+"/nodes/"+testNodeId+"/switch-port", req.URL.Path)
			var body node.ToggleSwitchPortRequest
			assert.NoError(tt, json.NewDecoder(req.Body).Decode(&body))
			assert.Equal(tt, true, body.Status)
			assert.Equal(tt, int32(3), body.Port)

			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",
				Header:     make(http.Header),
				Body:       io.NopCloser(bytes.NewBufferString(testOperationBody)),
			}
		}

		testControllerClient := node.NewNodeControllerClient(baseURL)
		testControllerClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		resp, err := testControllerClient.ToggleSwitchPort(testNodeId, node.ToggleSwitchPortRequest{Status: true, Port: 3})

		assert.NoError(tt, err)
		assert.Equal(tt, testOperationId, resp.OperationId)
		assert.Equal(tt, testResourceKey, resp.ResourceKey)
		assert.Equal(tt, "SUCCESS", resp.Status)
	})

	t.Run("InvalidResponse", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			resp := `{"error":"internal server error"}`

			return &http.Response{
				StatusCode: 500,
				Body:       io.NopCloser(bytes.NewBufferString(resp)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}
		}

		testControllerClient := node.NewNodeControllerClient(baseURL)
		testControllerClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		_, err := testControllerClient.ToggleSwitchPort(testNodeId, node.ToggleSwitchPortRequest{Status: true, Port: 3})

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

		testControllerClient := node.NewNodeControllerClient(baseURL)
		testControllerClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		_, err := testControllerClient.ToggleSwitchPort(testNodeId, node.ToggleSwitchPortRequest{Status: true, Port: 3})

		assert.Error(tt, err)
	})

	t.Run("RequestFailure", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			return nil
		}

		testControllerClient := node.NewNodeControllerClient(baseURL)
		testControllerClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		_, err := testControllerClient.ToggleSwitchPort(testNodeId, node.ToggleSwitchPortRequest{Status: true, Port: 3})

		assert.Error(tt, err)
	})
}

func TestControllerClient_ToggleRadio(t *testing.T) {
	baseURL := "http://test-controller-service.com"

	t.Run("Success", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, http.MethodPost, req.Method)
			assert.Equal(tt, node.ControllerEndpoint+"/nodes/"+testNodeId+"/radio/on", req.URL.Path)

			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",
				Header:     make(http.Header),
				Body:       io.NopCloser(bytes.NewBufferString(testOperationBody)),
			}
		}

		testControllerClient := node.NewNodeControllerClient(baseURL)
		testControllerClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		resp, err := testControllerClient.ToggleRadio(testNodeId, "on")

		assert.NoError(tt, err)
		assert.Equal(tt, testOperationId, resp.OperationId)
		assert.Equal(tt, testResourceKey, resp.ResourceKey)
		assert.Equal(tt, "SUCCESS", resp.Status)
	})

	t.Run("InvalidResponse", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			resp := `{"error":"internal server error"}`

			return &http.Response{
				StatusCode: 500,
				Body:       io.NopCloser(bytes.NewBufferString(resp)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}
		}

		testControllerClient := node.NewNodeControllerClient(baseURL)
		testControllerClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		_, err := testControllerClient.ToggleRadio(testNodeId, "on")

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

		testControllerClient := node.NewNodeControllerClient(baseURL)
		testControllerClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		_, err := testControllerClient.ToggleRadio(testNodeId, "on")

		assert.Error(tt, err)
	})

	t.Run("RequestFailure", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			return nil
		}

		testControllerClient := node.NewNodeControllerClient(baseURL)
		testControllerClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		_, err := testControllerClient.ToggleRadio(testNodeId, "on")

		assert.Error(tt, err)
	})
}

func TestControllerClient_ToggleService(t *testing.T) {
	baseURL := "http://test-controller-service.com"

	t.Run("Success", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, http.MethodPost, req.Method)
			assert.Equal(tt, node.ControllerEndpoint+"/nodes/"+testNodeId+"/service/off", req.URL.Path)

			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",
				Header:     make(http.Header),
				Body:       io.NopCloser(bytes.NewBufferString(testOperationBody)),
			}
		}

		testControllerClient := node.NewNodeControllerClient(baseURL)
		testControllerClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		resp, err := testControllerClient.ToggleService(testNodeId, "off")

		assert.NoError(tt, err)
		assert.Equal(tt, testOperationId, resp.OperationId)
		assert.Equal(tt, testResourceKey, resp.ResourceKey)
		assert.Equal(tt, "SUCCESS", resp.Status)
	})

	t.Run("InvalidResponse", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			resp := `{"error":"internal server error"}`

			return &http.Response{
				StatusCode: 500,
				Body:       io.NopCloser(bytes.NewBufferString(resp)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}
		}

		testControllerClient := node.NewNodeControllerClient(baseURL)
		testControllerClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		_, err := testControllerClient.ToggleService(testNodeId, "off")

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

		testControllerClient := node.NewNodeControllerClient(baseURL)
		testControllerClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		_, err := testControllerClient.ToggleService(testNodeId, "off")

		assert.Error(tt, err)
	})

	t.Run("RequestFailure", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			return nil
		}

		testControllerClient := node.NewNodeControllerClient(baseURL)
		testControllerClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		_, err := testControllerClient.ToggleService(testNodeId, "off")

		assert.Error(tt, err)
	})
}

func TestControllerClient_ConfigNode(t *testing.T) {
	baseURL := "http://test-controller-service.com"

	t.Run("Success", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, http.MethodPost, req.Method)
			assert.Equal(tt, node.ControllerEndpoint+"/nodes/"+testNodeId+"/config", req.URL.Path)

			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",
				Header:     make(http.Header),
				Body:       io.NopCloser(bytes.NewBufferString(testOperationBody)),
			}
		}

		testControllerClient := node.NewNodeControllerClient(baseURL)
		testControllerClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		resp, err := testControllerClient.ConfigNode(testNodeId)

		assert.NoError(tt, err)
		assert.Equal(tt, testOperationId, resp.OperationId)
		assert.Equal(tt, testResourceKey, resp.ResourceKey)
		assert.Equal(tt, "SUCCESS", resp.Status)
	})

	t.Run("InvalidResponse", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			resp := `{"error":"internal server error"}`

			return &http.Response{
				StatusCode: 500,
				Body:       io.NopCloser(bytes.NewBufferString(resp)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}
		}

		testControllerClient := node.NewNodeControllerClient(baseURL)
		testControllerClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		_, err := testControllerClient.ConfigNode(testNodeId)

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

		testControllerClient := node.NewNodeControllerClient(baseURL)
		testControllerClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		_, err := testControllerClient.ConfigNode(testNodeId)

		assert.Error(tt, err)
	})

	t.Run("RequestFailure", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			return nil
		}

		testControllerClient := node.NewNodeControllerClient(baseURL)
		testControllerClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		_, err := testControllerClient.ConfigNode(testNodeId)

		assert.Error(tt, err)
	})
}

func TestControllerClient_DeleteNodeConfig(t *testing.T) {
	baseURL := "http://test-controller-service.com"

	t.Run("Success", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, http.MethodDelete, req.Method)
			assert.Equal(tt, node.ControllerEndpoint+"/nodes/"+testNodeId+"/config", req.URL.Path)

			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",
				Header:     make(http.Header),
				Body:       io.NopCloser(bytes.NewBufferString(testOperationBody)),
			}
		}

		testControllerClient := node.NewNodeControllerClient(baseURL)
		testControllerClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		resp, err := testControllerClient.DeleteNodeConfig(testNodeId)

		assert.NoError(tt, err)
		assert.Equal(tt, testOperationId, resp.OperationId)
		assert.Equal(tt, testResourceKey, resp.ResourceKey)
		assert.Equal(tt, "SUCCESS", resp.Status)
	})

	t.Run("InvalidResponse", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			resp := `{"error":"internal server error"}`

			return &http.Response{
				StatusCode: 500,
				Body:       io.NopCloser(bytes.NewBufferString(resp)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}
		}

		testControllerClient := node.NewNodeControllerClient(baseURL)
		testControllerClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		_, err := testControllerClient.DeleteNodeConfig(testNodeId)

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

		testControllerClient := node.NewNodeControllerClient(baseURL)
		testControllerClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		_, err := testControllerClient.DeleteNodeConfig(testNodeId)

		assert.Error(tt, err)
	})

	t.Run("RequestFailure", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			return nil
		}

		testControllerClient := node.NewNodeControllerClient(baseURL)
		testControllerClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		_, err := testControllerClient.DeleteNodeConfig(testNodeId)

		assert.Error(tt, err)
	})
}
