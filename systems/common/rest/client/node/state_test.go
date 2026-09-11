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

const testStateId = "1a2b3c4d-5e6f-4a7b-8c9d-0e1f2a3b4c5d"

func TestStateClient_GetLatestState(t *testing.T) {
	baseURL := "http://test-state-service.com"

	t.Run("StateFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, http.MethodGet, req.Method)
			assert.Equal(tt, node.StateEndpoint+"/"+testNodeId+"/latest", req.URL.Path)

			body := `{"State":{"id":"` + testStateId + `","nodeId":"` + testNodeId + `","currentState":"Operational","subState":["sub-1"],"events":["event-1"],"nodeType":"tnode","createdAt":"2026-09-11T10:00:00Z","updatedAt":"2026-09-11T10:00:00Z"}}`

			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",
				Header:     make(http.Header),
				Body:       io.NopCloser(bytes.NewBufferString(body)),
			}
		}

		testStateClient := node.NewNodeStateClient(baseURL)
		testStateClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		state, err := testStateClient.GetLatestState(testNodeId)

		assert.NoError(tt, err)
		assert.NotNil(tt, state)
		assert.Equal(tt, testStateId, state.Id)
		assert.Equal(tt, testNodeId, state.NodeId)
		assert.Equal(tt, "Operational", state.CurrentState)
		assert.Equal(tt, []string{"sub-1"}, state.SubState)
		assert.Equal(tt, []string{"event-1"}, state.Events)
		assert.Equal(tt, "tnode", state.NodeType)
	})

	t.Run("NoStateRecord", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			body := `{"State":null}`

			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",
				Header:     make(http.Header),
				Body:       io.NopCloser(bytes.NewBufferString(body)),
			}
		}

		testStateClient := node.NewNodeStateClient(baseURL)
		testStateClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		state, err := testStateClient.GetLatestState(testNodeId)

		assert.NoError(tt, err)
		assert.Nil(tt, state)
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

		testStateClient := node.NewNodeStateClient(baseURL)
		testStateClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		_, err := testStateClient.GetLatestState(testNodeId)

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

		testStateClient := node.NewNodeStateClient(baseURL)
		testStateClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		_, err := testStateClient.GetLatestState(testNodeId)

		assert.Error(tt, err)
	})

	t.Run("RequestFailure", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			return nil
		}

		testStateClient := node.NewNodeStateClient(baseURL)
		testStateClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		_, err := testStateClient.GetLatestState(testNodeId)

		assert.Error(tt, err)
	})
}
