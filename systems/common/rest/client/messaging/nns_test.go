/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package messaging_test

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/tj/assert"

	"github.com/ukama/ukama/systems/common/rest/client"
	"github.com/ukama/ukama/systems/common/rest/client/messaging"
)

const testNodeID = "test-node-id"

const (
	contentTypeHeader = "Content-Type"
	contentTypeJSON   = "application/json"

	pathMesh = "/mesh/"
	pathNode = "/node/"
	pathList = "/list"

	statusOK       = "200 OK"
	statusNotFound = "404 NOT FOUND"
)

func TestNnsClientGetMesh(t *testing.T) {
	t.Run("MeshFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), messaging.NnsEndpoint+pathMesh+testNodeID)

			mesh := `{"meshIp":"10.0.0.1","meshPort":1234}`

			return &http.Response{
				StatusCode: 200,
				Status:     statusOK,
				Body:       io.NopCloser(bytes.NewBufferString(mesh)),
				Header:     make(http.Header),
			}
		}

		testNnsClient := messaging.NewNnsClient("")
		testNnsClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		m, err := testNnsClient.GetMesh(testNodeID)

		assert.NoError(tt, err)
		assert.Equal(tt, "10.0.0.1", m.MeshIp)
		assert.Equal(tt, 1234, m.MeshPort)
	})

	t.Run("MeshNotFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), messaging.NnsEndpoint+pathMesh+testNodeID)

			resp := `{"error":"not found"}`

			return &http.Response{
				StatusCode: 404,
				Status:     statusNotFound,
				Body:       io.NopCloser(bytes.NewBufferString(resp)),
				Header:     http.Header{contentTypeHeader: []string{contentTypeJSON}},
			}
		}

		testNnsClient := messaging.NewNnsClient("")
		testNnsClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		m, err := testNnsClient.GetMesh(testNodeID)

		assert.Error(tt, err)
		assert.Nil(tt, m)
	})

	t.Run("InvalidResponsePayload", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), messaging.NnsEndpoint+pathMesh+testNodeID)

			return &http.Response{
				StatusCode: 200,
				Status:     statusOK,
				Body:       io.NopCloser(bytes.NewBufferString(`OK`)),
				Header:     make(http.Header),
			}
		}

		testNnsClient := messaging.NewNnsClient("")
		testNnsClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		m, err := testNnsClient.GetMesh(testNodeID)

		assert.Error(tt, err)
		assert.Nil(tt, m)
	})

	t.Run("RequestFailure", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), messaging.NnsEndpoint+pathMesh+testNodeID)
			return nil
		}

		testNnsClient := messaging.NewNnsClient("")
		testNnsClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		m, err := testNnsClient.GetMesh(testNodeID)

		assert.Error(tt, err)
		assert.Nil(tt, m)
	})
}

func TestNnsClientGetNode(t *testing.T) {
	t.Run("NodeFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), messaging.NnsEndpoint+pathNode+testNodeID)

			node := `{"nodeId":"test-node-id","nodeIp":"192.168.1.10","nodePort":4321}`

			return &http.Response{
				StatusCode: 200,
				Status:     statusOK,
				Body:       io.NopCloser(bytes.NewBufferString(node)),
				Header:     make(http.Header),
			}
		}

		testNnsClient := messaging.NewNnsClient("")
		testNnsClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		n, err := testNnsClient.GetNode(testNodeID)

		assert.NoError(tt, err)
		assert.Equal(tt, testNodeID, n.NodeId)
		assert.Equal(tt, "192.168.1.10", n.NodeIp)
		assert.Equal(tt, 4321, n.NodePort)
	})

	t.Run("NodeNotFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), messaging.NnsEndpoint+pathNode+testNodeID)

			resp := `{"error":"not found"}`

			return &http.Response{
				StatusCode: 404,
				Status:     statusNotFound,
				Body:       io.NopCloser(bytes.NewBufferString(resp)),
				Header:     http.Header{contentTypeHeader: []string{contentTypeJSON}},
			}
		}

		testNnsClient := messaging.NewNnsClient("")
		testNnsClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		n, err := testNnsClient.GetNode(testNodeID)

		assert.Error(tt, err)
		assert.Nil(tt, n)
	})

	t.Run("InvalidResponsePayload", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), messaging.NnsEndpoint+pathNode+testNodeID)

			return &http.Response{
				StatusCode: 200,
				Status:     statusOK,
				Body:       io.NopCloser(bytes.NewBufferString(`OK`)),
				Header:     make(http.Header),
			}
		}

		testNnsClient := messaging.NewNnsClient("")
		testNnsClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		n, err := testNnsClient.GetNode(testNodeID)

		assert.Error(tt, err)
		assert.Nil(tt, n)
	})

	t.Run("RequestFailure", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), messaging.NnsEndpoint+pathNode+testNodeID)
			return nil
		}

		testNnsClient := messaging.NewNnsClient("")
		testNnsClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		n, err := testNnsClient.GetNode(testNodeID)

		assert.Error(tt, err)
		assert.Nil(tt, n)
	})
}

func TestNnsClientList(t *testing.T) {
	t.Run("ListFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), messaging.NnsEndpoint+pathList)

			list := `{"list":[{"nodeId":"node-1","nodeIp":"192.168.1.10","nodePort":4321,"meshIp":"10.0.0.1","meshPort":1234,"org":"ukama","network":"net-1","site":"site-1","meshHostName":"mesh-1"}]}`

			return &http.Response{
				StatusCode: 200,
				Status:     statusOK,
				Body:       io.NopCloser(bytes.NewBufferString(list)),
				Header:     make(http.Header),
			}
		}

		testNnsClient := messaging.NewNnsClient("")
		testNnsClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		l, err := testNnsClient.List()

		assert.NoError(tt, err)
		assert.NotNil(tt, l)
		assert.Equal(tt, 1, len(l.List))
		assert.Equal(tt, "node-1", l.List[0].NodeId)
		assert.Equal(tt, "10.0.0.1", l.List[0].MeshIp)
		assert.Equal(tt, "mesh-1", l.List[0].MeshHostName)
	})

	t.Run("ListNotFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), messaging.NnsEndpoint+pathList)

			resp := `{"error":"not found"}`

			return &http.Response{
				StatusCode: 404,
				Status:     statusNotFound,
				Body:       io.NopCloser(bytes.NewBufferString(resp)),
				Header:     http.Header{contentTypeHeader: []string{contentTypeJSON}},
			}
		}

		testNnsClient := messaging.NewNnsClient("")
		testNnsClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		l, err := testNnsClient.List()

		assert.Error(tt, err)
		assert.Nil(tt, l)
	})

	t.Run("InvalidResponsePayload", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), messaging.NnsEndpoint+pathList)

			return &http.Response{
				StatusCode: 200,
				Status:     statusOK,
				Body:       io.NopCloser(bytes.NewBufferString(`OK`)),
				Header:     make(http.Header),
			}
		}

		testNnsClient := messaging.NewNnsClient("")
		testNnsClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		l, err := testNnsClient.List()

		assert.Error(tt, err)
		assert.Nil(tt, l)
	})

	t.Run("RequestFailure", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), messaging.NnsEndpoint+pathList)
			return nil
		}

		testNnsClient := messaging.NewNnsClient("")
		testNnsClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		l, err := testNnsClient.List()

		assert.Error(tt, err)
		assert.Nil(tt, l)
	})
}

