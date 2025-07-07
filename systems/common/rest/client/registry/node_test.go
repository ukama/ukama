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
	"github.com/ukama/ukama/systems/common/uuid"
)

const testNodeId = "uk-sa2341-hnode-v0-a1a0"

func TestNodeClient_Get(t *testing.T) {
	t.Run("NodeFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			// Test request parameters
			assert.Equal(tt, req.URL.String(), registry.NodeEndpoint+"/"+testNodeId)

			// fake node info
			node := `{"node":{"id": "uk-sa2341-hnode-v0-a1a0", "name": "Node-A"}}`

			// Send mock response
			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",

				// Send response to be tested
				Body: io.NopCloser(bytes.NewBufferString(node)),

				// Must be set to non-nil value or it panics
				Header: make(http.Header),
			}
		}

		testNodeClient := registry.NewNodeClient("")

		// We replace the transport mechanism by mocking the http request
		// so that the test stays a unit test e.g, no server/network call.
		testNodeClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		n, err := testNodeClient.Get(testNodeId)

		assert.NoError(tt, err)
		assert.Equal(tt, testNodeId, n.Id)
	})

	t.Run("NodeNotFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), registry.NodeEndpoint+"/"+testNodeId)

			// error payload
			resp := `{"error":"not found"}`

			return &http.Response{
				StatusCode: 404,
				Status:     "404 NOT FOUND",
				Body:       io.NopCloser(bytes.NewBufferString(resp)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}
		}

		testNodeClient := registry.NewNodeClient("")

		testNodeClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		n, err := testNodeClient.Get(testNodeId)

		assert.Error(tt, err)
		assert.Nil(tt, n)
	})

	t.Run("InvalidResponsePayload", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), registry.NodeEndpoint+"/"+testNodeId)

			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",
				Body:       io.NopCloser(bytes.NewBufferString(`OK`)),
				Header:     make(http.Header),
			}
		}

		testNodeClient := registry.NewNodeClient("")

		testNodeClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		n, err := testNodeClient.Get(testNodeId)

		assert.Error(tt, err)
		assert.Nil(tt, n)
	})

	t.Run("RequestFailure", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), registry.NodeEndpoint+"/"+testNodeId)

			return nil
		}

		testNodeClient := registry.NewNodeClient("")

		testNodeClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		n, err := testNodeClient.Get(testNodeId)

		assert.Error(tt, err)
		assert.Nil(tt, n)
	})
}

func TestNodeClient_Add(t *testing.T) {
	t.Run("NodeAdded", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			// Test request parameters
			assert.Equal(tt, req.URL.String(), registry.NodeEndpoint)

			// fake node info
			node := `{"node":{"id": "uk-sa2341-hnode-v0-a1a0", "name": "Node-A"}}`

			// Send mock response
			return &http.Response{
				StatusCode: 201,
				Status:     "201 CREATED",

				// Send response to be tested
				Body: io.NopCloser(bytes.NewBufferString(node)),

				// Must be set to non-nil value or it panics
				Header: make(http.Header),
			}
		}

		testNodeClient := registry.NewNodeClient("")

		// We replace the transport mechanism by mocking the http request
		// so that the test stays a unit test e.g, no server/network call.
		testNodeClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		n, err := testNodeClient.Add(
			registry.AddNodeRequest{
				NodeId: "uk-sa2341-hnode-v0-a1a0",
				Name:   "Node-A",
			})

		assert.NoError(tt, err)
		assert.Equal(tt, testNodeId, n.Id)
	})

	t.Run("InvalidResponseHeader", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), registry.NodeEndpoint)

			// error payload
			resp := `{"error":"internal server error"}`

			return &http.Response{
				StatusCode: 500,
				Status:     "500 INTERNAL SERVER ERROR",
				Body:       io.NopCloser(bytes.NewBufferString(resp)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}
		}

		testNodeClient := registry.NewNodeClient("")

		testNodeClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		n, err := testNodeClient.Add(
			registry.AddNodeRequest{
				NodeId: "uk-sa2341-hnode-v0-a1a0",
				Name:   "Node-A",
			},
		)

		assert.Error(tt, err)
		assert.Nil(tt, n)
	})

	t.Run("InvalidResponsePayload", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), registry.NodeEndpoint)

			return &http.Response{
				StatusCode: 201,
				Status:     "201 CREATED",
				Body:       io.NopCloser(bytes.NewBufferString(`CREATED`)),
				Header:     make(http.Header),
			}
		}

		testNodeClient := registry.NewNodeClient("")

		testNodeClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		n, err := testNodeClient.Add(
			registry.AddNodeRequest{
				NodeId: "uk-sa2341-hnode-v0-a1a0",
				Name:   "Node-A",
			},
		)

		assert.Error(tt, err)
		assert.Nil(tt, n)
	})

	t.Run("RequestFailure", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), registry.NodeEndpoint)

			return nil
		}

		testNodeClient := registry.NewNodeClient("")

		testNodeClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		n, err := testNodeClient.Add(
			registry.AddNodeRequest{
				NodeId: "uk-sa2341-hnode-v0-a1a0",
				Name:   "Node-A",
			},
		)

		assert.Error(tt, err)
		assert.Nil(tt, n)
	})
}

func TestNodeClient_GetAll(t *testing.T) {
	t.Run("NodesFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			// Test request parameters
			assert.Equal(tt, req.URL.String(), registry.NodeEndpoint)

			// fake nodes info
			nodes := `{"nodes":[{"id": "uk-sa2341-hnode-v0-a1a0", "name": "Node-A"}]}`

			// Send mock response
			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",

				// Send response to be tested
				Body: io.NopCloser(bytes.NewBufferString(nodes)),

				// Must be set to non-nil value or it panics
				Header: make(http.Header),
			}
		}

		testNodeClient := registry.NewNodeClient("")

		// We replace the transport mechanism by mocking the http request
		// so that the test stays a unit test e.g, no server/network call.
		testNodeClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		nodes, err := testNodeClient.GetAll()

		assert.NoError(tt, err)
		assert.Equal(tt, 1, len(nodes.Nodes))
		assert.Equal(tt, testNodeId, nodes.Nodes[0].Id)
	})

	t.Run("NodesNotFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), registry.NodeEndpoint)

			// error payload
			resp := `{"error":"not found"}`

			return &http.Response{
				StatusCode: 404,
				Status:     "404 NOT FOUND",
				Body:       io.NopCloser(bytes.NewBufferString(resp)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}
		}

		testNodeClient := registry.NewNodeClient("")

		testNodeClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		n, err := testNodeClient.GetAll()

		assert.Error(tt, err)
		assert.Nil(tt, n)
	})

	t.Run("InvalidResponsePayload", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), registry.NodeEndpoint)

			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",
				Body:       io.NopCloser(bytes.NewBufferString(`OK`)),
				Header:     make(http.Header),
			}
		}

		testNodeClient := registry.NewNodeClient("")

		testNodeClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		n, err := testNodeClient.GetAll()

		assert.Error(tt, err)
		assert.Nil(tt, n)
	})

	t.Run("RequestFailure", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), registry.NodeEndpoint)

			return nil
		}

		testNodeClient := registry.NewNodeClient("")

		testNodeClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		n, err := testNodeClient.GetAll()

		assert.Error(tt, err)
		assert.Nil(tt, n)
	})
}

func TestNodeClient_GetNodesBySite(t *testing.T) {
	t.Run("NodesFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			// Test request parameters
			assert.Equal(tt, req.URL.String(), registry.NodeEndpoint+"/sites/"+testUuid)

			// fake nodes info
			nodes := `{"nodes":[{"id": "uk-sa2341-hnode-v0-a1a0", "name": "Node-A"}]}`

			// Send mock response
			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",

				// Send response to be tested
				Body: io.NopCloser(bytes.NewBufferString(nodes)),

				// Must be set to non-nil value or it panics
				Header: make(http.Header),
			}
		}

		testNodeClient := registry.NewNodeClient("")

		// We replace the transport mechanism by mocking the http request
		// so that the test stays a unit test e.g, no server/network call.
		testNodeClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		nodes, err := testNodeClient.GetNodesBySite(testUuid)

		assert.NoError(tt, err)
		assert.Equal(tt, 1, len(nodes.Nodes))
		assert.Equal(tt, testNodeId, nodes.Nodes[0].Id)
	})

	t.Run("NodesNotFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), registry.NodeEndpoint+"/sites/"+testUuid)

			// error payload
			resp := `{"error":"not found"}`

			return &http.Response{
				StatusCode: 404,
				Status:     "404 NOT FOUND",
				Body:       io.NopCloser(bytes.NewBufferString(resp)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}
		}

		testNodeClient := registry.NewNodeClient("")

		testNodeClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		n, err := testNodeClient.GetNodesBySite(testUuid)

		assert.Error(tt, err)
		assert.Nil(tt, n)
	})

	t.Run("InvalidResponsePayload", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), registry.NodeEndpoint+"/sites/"+testUuid)

			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",
				Body:       io.NopCloser(bytes.NewBufferString(`OK`)),
				Header:     make(http.Header),
			}
		}

		testNodeClient := registry.NewNodeClient("")

		testNodeClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		n, err := testNodeClient.GetNodesBySite(testUuid)

		assert.Error(tt, err)
		assert.Nil(tt, n)
	})

	t.Run("RequestFailure", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), registry.NodeEndpoint+"/sites/"+testUuid)

			return nil
		}

		testNodeClient := registry.NewNodeClient("")

		testNodeClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		n, err := testNodeClient.GetNodesBySite(testUuid)

		assert.Error(tt, err)
		assert.Nil(tt, n)
	})
}
func TestNodeClient_Delete(t *testing.T) {
	t.Run("NodeFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			// Test request parameters
			assert.Equal(tt, req.URL.String(), registry.NodeEndpoint+"/"+testNodeId)

			// Send mock response
			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",
				Header:     make(http.Header),
			}
		}

		testNodeClient := registry.NewNodeClient("")

		// We replace the transport mechanism by mocking the http request
		// so that the test stays a unit test e.g, no server/network call.
		testNodeClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		err := testNodeClient.Delete(testNodeId)

		assert.NoError(tt, err)
	})

	t.Run("NodeNotFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), registry.NodeEndpoint+"/"+testNodeId)

			// error payload
			resp := `{"error":"not found"}`

			return &http.Response{
				StatusCode: 404,
				Status:     "404 NOT FOUND",
				Body:       io.NopCloser(bytes.NewBufferString(resp)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}
		}

		testNodeClient := registry.NewNodeClient("")

		testNodeClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		err := testNodeClient.Delete(testNodeId)

		assert.Error(tt, err)
	})

	t.Run("RequestFailure", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), registry.NodeEndpoint+"/"+testNodeId)

			return nil
		}

		testNodeClient := registry.NewNodeClient("")

		testNodeClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		err := testNodeClient.Delete(testNodeId)

		assert.Error(tt, err)
	})
}

func TestNodeClient_Attach(t *testing.T) {
	t.Run("NodeAttached", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			// Test request parameters
			assert.Equal(tt, req.URL.String(), registry.NodeEndpoint+"/"+testNodeId+"/attach")

			// Send mock response
			return &http.Response{
				StatusCode: 201,

				// Must be set to non-nil value or it panics
				Header: make(http.Header),
			}
		}

		testNodeClient := registry.NewNodeClient("")

		// We replace the transport mechanism by mocking the http request
		// so that the test stays a unit test e.g, no server/network call.
		testNodeClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		err := testNodeClient.Attach(testNodeId,
			registry.AttachNodesRequest{
				AmpNodeL: "uk-sa2341-anode-v0-a1a0",
				AmpNodeR: "uk-sa2341-anode-v0-a1a1",
			})

		assert.NoError(tt, err)
	})

	t.Run("InvalidResponseHeader", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), registry.NodeEndpoint+"/"+testNodeId+"/attach")

			// error payload
			resp := `{"error":"internal server error"}`

			return &http.Response{
				StatusCode: 500,
				Status:     "500 INTERNAL SERVER ERROR",
				Body:       io.NopCloser(bytes.NewBufferString(resp)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}
		}

		testNodeClient := registry.NewNodeClient("")

		testNodeClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		err := testNodeClient.Attach(testNodeId,
			registry.AttachNodesRequest{
				AmpNodeL: "uk-sa2341-anode-v0-a1a0",
				AmpNodeR: "uk-sa2341-anode-v0-a1a1",
			})

		assert.Error(tt, err)
	})

	t.Run("RequestFailure", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), registry.NodeEndpoint+"/"+testNodeId+"/attach")

			return nil
		}

		testNodeClient := registry.NewNodeClient("")

		testNodeClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		err := testNodeClient.Attach(testNodeId,
			registry.AttachNodesRequest{
				AmpNodeL: "uk-sa2341-anode-v0-a1a0",
				AmpNodeR: "uk-sa2341-anode-v0-a1a1",
			})

		assert.Error(tt, err)
	})
}

func TestNodeClient_Detach(t *testing.T) {
	t.Run("NodeFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			// Test request parameters
			assert.Equal(tt, req.URL.String(), registry.NodeEndpoint+"/"+testNodeId+"/attach")

			// Send mock response
			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",
				Header:     make(http.Header),
			}
		}

		testNodeClient := registry.NewNodeClient("")

		// We replace the transport mechanism by mocking the http request
		// so that the test stays a unit test e.g, no server/network call.
		testNodeClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		err := testNodeClient.Detach(testNodeId)

		assert.NoError(tt, err)
	})

	t.Run("NodeNotFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), registry.NodeEndpoint+"/"+testNodeId+"/attach")

			// error payload
			resp := `{"error":"not found"}`

			return &http.Response{
				StatusCode: 404,
				Status:     "404 NOT FOUND",
				Body:       io.NopCloser(bytes.NewBufferString(resp)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}
		}

		testNodeClient := registry.NewNodeClient("")

		testNodeClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		err := testNodeClient.Detach(testNodeId)

		assert.Error(tt, err)
	})

	t.Run("RequestFailure", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), registry.NodeEndpoint+"/"+testNodeId+"/attach")

			return nil
		}

		testNodeClient := registry.NewNodeClient("")

		testNodeClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		err := testNodeClient.Detach(testNodeId)

		assert.Error(tt, err)
	})
}

func TestNodeClient_AddToSite(t *testing.T) {
	t.Run("NodeAdded", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			// Test request parameters
			assert.Equal(tt, req.URL.String(), registry.NodeEndpoint+"/"+testNodeId+"/sites")

			// Send mock response
			return &http.Response{
				StatusCode: 201,

				// Must be set to non-nil value or it panics
				Header: make(http.Header),
			}
		}

		testNodeClient := registry.NewNodeClient("")

		// We replace the transport mechanism by mocking the http request
		// so that the test stays a unit test e.g, no server/network call.
		testNodeClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		err := testNodeClient.AddToSite(testNodeId,
			registry.AddToSiteRequest{
				NetworkId: uuid.NewV4().String(),
				SiteId:    uuid.NewV4().String(),
			})

		assert.NoError(tt, err)
	})

	t.Run("InvalidResponseHeader", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), registry.NodeEndpoint+"/"+testNodeId+"/sites")

			// error payload
			resp := `{"error":"internal server error"}`

			return &http.Response{
				StatusCode: 500,
				Status:     "500 INTERNAL SERVER ERROR",
				Body:       io.NopCloser(bytes.NewBufferString(resp)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}
		}

		testNodeClient := registry.NewNodeClient("")

		testNodeClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		err := testNodeClient.AddToSite(testNodeId,
			registry.AddToSiteRequest{
				NetworkId: uuid.NewV4().String(),
				SiteId:    uuid.NewV4().String(),
			})

		assert.Error(tt, err)
	})

	t.Run("RequestFailure", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), registry.NodeEndpoint+"/"+testNodeId+"/sites")

			return nil
		}

		testNodeClient := registry.NewNodeClient("")

		testNodeClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		err := testNodeClient.AddToSite(testNodeId,
			registry.AddToSiteRequest{
				NetworkId: uuid.NewV4().String(),
				SiteId:    uuid.NewV4().String(),
			})

		assert.Error(tt, err)
	})
}

func TestNodeClient_RemoveFromSite(t *testing.T) {
	t.Run("NodeFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			// Test request parameters
			assert.Equal(tt, req.URL.String(), registry.NodeEndpoint+"/"+testNodeId+"/sites")

			// Send mock response
			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",
				Header:     make(http.Header),
			}
		}

		testNodeClient := registry.NewNodeClient("")

		// We replace the transport mechanism by mocking the http request
		// so that the test stays a unit test e.g, no server/network call.
		testNodeClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		err := testNodeClient.RemoveFromSite(testNodeId)

		assert.NoError(tt, err)
	})

	t.Run("NodeNotFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), registry.NodeEndpoint+"/"+testNodeId+"/sites")

			// error payload
			resp := `{"error":"not found"}`

			return &http.Response{
				StatusCode: 404,
				Status:     "404 NOT FOUND",
				Body:       io.NopCloser(bytes.NewBufferString(resp)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}
		}

		testNodeClient := registry.NewNodeClient("")

		testNodeClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		err := testNodeClient.RemoveFromSite(testNodeId)

		assert.Error(tt, err)
	})

	t.Run("RequestFailure", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), registry.NodeEndpoint+"/"+testNodeId+"/sites")

			return nil
		}

		testNodeClient := registry.NewNodeClient("")

		testNodeClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		err := testNodeClient.RemoveFromSite(testNodeId)

		assert.Error(tt, err)
	})
}
