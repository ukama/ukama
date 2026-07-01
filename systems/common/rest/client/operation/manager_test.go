/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package operation_test

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/tj/assert"

	"github.com/ukama/ukama/systems/common/rest/client"
	"github.com/ukama/ukama/systems/common/rest/client/operation"
)

const (
	testOperationId = "03cb753f-5e03-4c97-8e47-625115476c72"
	testResourceKey = "node/03cb753f-5e03-4c97-8e47-625115476c72"
)

func TestManagerClient_Start(t *testing.T) {
	t.Run("OperationStarted", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), operation.OperationsEndpoint)
			assert.Equal(tt, "POST", req.Method)

			body := `{"operation":{"id":"03cb753f-5e03-4c97-8e47-625115476c72","type":"deploy","system":"node","status":"PENDING","fencingToken":1,"resourceKey":"node/03cb753f-5e03-4c97-8e47-625115476c72","leaseExpiresAt":"2026-06-29T12:00:00Z","createdAt":"2026-06-29T11:00:00Z"}}`

			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",
				Body:       io.NopCloser(bytes.NewBufferString(body)),
				Header:     make(http.Header),
			}
		}

		testManagerClient := operation.NewManagerClient("")
		testManagerClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		resp, err := testManagerClient.Start(operation.StartRequest{
			Type:        "deploy",
			System:      "node",
			ResourceKey: testResourceKey,
		})

		assert.NoError(tt, err)
		assert.NotNil(tt, resp.Operation)
		assert.Equal(tt, testOperationId, resp.Operation.Id)
		assert.Equal(tt, operation.StatusPending, resp.Operation.Status)
		assert.Nil(tt, resp.ConflictingOperation)
	})

	t.Run("ConflictingOperation", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), operation.OperationsEndpoint)
			assert.Equal(tt, "POST", req.Method)

			body := `{"conflictingOperation":{"id":"7f8b9c0d-1e2f-3a4b-5c6d-7e8f9a0b1c2d","type":"deploy","system":"node","status":"RUNNING","fencingToken":2,"resourceKey":"node/03cb753f-5e03-4c97-8e47-625115476c72","leaseExpiresAt":"2026-06-29T12:00:00Z","createdAt":"2026-06-29T11:00:00Z"}}`

			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",
				Body:       io.NopCloser(bytes.NewBufferString(body)),
				Header:     make(http.Header),
			}
		}

		testManagerClient := operation.NewManagerClient("")
		testManagerClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		resp, err := testManagerClient.Start(operation.StartRequest{
			Type:        "deploy",
			System:      "node",
			ResourceKey: testResourceKey,
		})

		assert.NoError(tt, err)
		assert.Nil(tt, resp.Operation)
		assert.NotNil(tt, resp.ConflictingOperation)
		assert.Equal(tt, operation.StatusRunning, resp.ConflictingOperation.Status)
	})

	t.Run("StartFailed", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), operation.OperationsEndpoint)
			assert.Equal(tt, "POST", req.Method)

			resp := `{"error":"conflict"}`

			return &http.Response{
				StatusCode: 409,
				Status:     "409 CONFLICT",
				Body:       io.NopCloser(bytes.NewBufferString(resp)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}
		}

		testManagerClient := operation.NewManagerClient("")
		testManagerClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		resp, err := testManagerClient.Start(operation.StartRequest{
			Type:        "deploy",
			System:      "node",
			ResourceKey: testResourceKey,
		})

		assert.Error(tt, err)
		assert.Nil(tt, resp)
	})

	t.Run("InvalidResponsePayload", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), operation.OperationsEndpoint)
			assert.Equal(tt, "POST", req.Method)

			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",
				Body:       io.NopCloser(bytes.NewBufferString(`OK`)),
				Header:     make(http.Header),
			}
		}

		testManagerClient := operation.NewManagerClient("")
		testManagerClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		resp, err := testManagerClient.Start(operation.StartRequest{
			Type:        "deploy",
			System:      "node",
			ResourceKey: testResourceKey,
		})

		assert.Error(tt, err)
		assert.Nil(tt, resp)
	})

	t.Run("RequestFailure", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), operation.OperationsEndpoint)
			assert.Equal(tt, "POST", req.Method)

			return nil
		}

		testManagerClient := operation.NewManagerClient("")
		testManagerClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		resp, err := testManagerClient.Start(operation.StartRequest{
			Type:        "deploy",
			System:      "node",
			ResourceKey: testResourceKey,
		})

		assert.Error(tt, err)
		assert.Nil(tt, resp)
	})
}

func TestManagerClient_Get(t *testing.T) {
	t.Run("OperationFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), operation.OperationsEndpoint+"/"+testOperationId)

			body := `{"operation":{"id":"03cb753f-5e03-4c97-8e47-625115476c72","type":"deploy","system":"node","status":"SUCCESS","fencingToken":1,"resourceKey":"node/03cb753f-5e03-4c97-8e47-625115476c72","leaseExpiresAt":"2026-06-29T12:00:00Z","createdAt":"2026-06-29T11:00:00Z"}}`

			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",
				Body:       io.NopCloser(bytes.NewBufferString(body)),
				Header:     make(http.Header),
			}
		}

		testManagerClient := operation.NewManagerClient("")
		testManagerClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		op, err := testManagerClient.Get(testOperationId)

		assert.NoError(tt, err)
		assert.Equal(tt, testOperationId, op.Id)
		assert.Equal(tt, operation.StatusSuccess, op.Status)
	})

	t.Run("OperationNotFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), operation.OperationsEndpoint+"/"+testOperationId)

			resp := `{"error":"not found"}`

			return &http.Response{
				StatusCode: 404,
				Status:     "404 NOT FOUND",
				Body:       io.NopCloser(bytes.NewBufferString(resp)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}
		}

		testManagerClient := operation.NewManagerClient("")
		testManagerClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		op, err := testManagerClient.Get(testOperationId)

		assert.Error(tt, err)
		assert.Nil(tt, op)
	})

	t.Run("InvalidResponsePayload", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), operation.OperationsEndpoint+"/"+testOperationId)

			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",
				Body:       io.NopCloser(bytes.NewBufferString(`OK`)),
				Header:     make(http.Header),
			}
		}

		testManagerClient := operation.NewManagerClient("")
		testManagerClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		op, err := testManagerClient.Get(testOperationId)

		assert.Error(tt, err)
		assert.Nil(tt, op)
	})

	t.Run("RequestFailure", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), operation.OperationsEndpoint+"/"+testOperationId)

			return nil
		}

		testManagerClient := operation.NewManagerClient("")
		testManagerClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		op, err := testManagerClient.Get(testOperationId)

		assert.Error(tt, err)
		assert.Nil(tt, op)
	})
}

func TestManagerClient_GetByResource(t *testing.T) {
	t.Run("OperationFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, operation.OperationsEndpoint, req.URL.Path)
			assert.Equal(tt, testResourceKey, req.URL.Query().Get("resource_key"))

			body := `{"operation":{"id":"03cb753f-5e03-4c97-8e47-625115476c72","type":"deploy","system":"node","status":"RUNNING","fencingToken":1,"resourceKey":"node/03cb753f-5e03-4c97-8e47-625115476c72","leaseExpiresAt":"2026-06-29T12:00:00Z","createdAt":"2026-06-29T11:00:00Z"}}`

			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",
				Body:       io.NopCloser(bytes.NewBufferString(body)),
				Header:     make(http.Header),
			}
		}

		testManagerClient := operation.NewManagerClient("")
		testManagerClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		op, err := testManagerClient.GetByResource(testResourceKey)

		assert.NoError(tt, err)
		assert.Equal(tt, testOperationId, op.Id)
		assert.Equal(tt, testResourceKey, op.ResourceKey)
	})

	t.Run("OperationNotFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, operation.OperationsEndpoint, req.URL.Path)
			assert.Equal(tt, testResourceKey, req.URL.Query().Get("resource_key"))

			resp := `{"error":"not found"}`

			return &http.Response{
				StatusCode: 404,
				Status:     "404 NOT FOUND",
				Body:       io.NopCloser(bytes.NewBufferString(resp)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}
		}

		testManagerClient := operation.NewManagerClient("")
		testManagerClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		op, err := testManagerClient.GetByResource(testResourceKey)

		assert.Error(tt, err)
		assert.Nil(tt, op)
	})

	t.Run("InvalidResponsePayload", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, operation.OperationsEndpoint, req.URL.Path)
			assert.Equal(tt, testResourceKey, req.URL.Query().Get("resource_key"))

			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",
				Body:       io.NopCloser(bytes.NewBufferString(`OK`)),
				Header:     make(http.Header),
			}
		}

		testManagerClient := operation.NewManagerClient("")
		testManagerClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		op, err := testManagerClient.GetByResource(testResourceKey)

		assert.Error(tt, err)
		assert.Nil(tt, op)
	})

	t.Run("RequestFailure", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, operation.OperationsEndpoint, req.URL.Path)
			assert.Equal(tt, testResourceKey, req.URL.Query().Get("resource_key"))

			return nil
		}

		testManagerClient := operation.NewManagerClient("")
		testManagerClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		op, err := testManagerClient.GetByResource(testResourceKey)

		assert.Error(tt, err)
		assert.Nil(tt, op)
	})
}

func TestManagerClient_MarkRunning(t *testing.T) {
	const fencingToken = uint64(42)

	t.Run("OperationMarkedRunning", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), operation.OperationsEndpoint+"/"+testOperationId+"/run")
			assert.Equal(tt, "POST", req.Method)

			body := `{"operation":{"id":"03cb753f-5e03-4c97-8e47-625115476c72","type":"deploy","system":"node","status":"RUNNING","fencingToken":42,"resourceKey":"node/03cb753f-5e03-4c97-8e47-625115476c72","leaseExpiresAt":"2026-06-29T12:00:00Z","createdAt":"2026-06-29T11:00:00Z"}}`

			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",
				Body:       io.NopCloser(bytes.NewBufferString(body)),
				Header:     make(http.Header),
			}
		}

		testManagerClient := operation.NewManagerClient("")
		testManagerClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		op, err := testManagerClient.MarkRunning(testOperationId, fencingToken)

		assert.NoError(tt, err)
		assert.Equal(tt, testOperationId, op.Id)
		assert.Equal(tt, operation.StatusRunning, op.Status)
		assert.Equal(tt, fencingToken, op.FencingToken)
	})

	t.Run("MarkRunningFailed", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), operation.OperationsEndpoint+"/"+testOperationId+"/run")
			assert.Equal(tt, "POST", req.Method)

			resp := `{"error":"invalid fencing token"}`

			return &http.Response{
				StatusCode: 409,
				Status:     "409 CONFLICT",
				Body:       io.NopCloser(bytes.NewBufferString(resp)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}
		}

		testManagerClient := operation.NewManagerClient("")
		testManagerClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		op, err := testManagerClient.MarkRunning(testOperationId, fencingToken)

		assert.Error(tt, err)
		assert.Nil(tt, op)
	})

	t.Run("InvalidResponsePayload", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), operation.OperationsEndpoint+"/"+testOperationId+"/run")
			assert.Equal(tt, "POST", req.Method)

			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",
				Body:       io.NopCloser(bytes.NewBufferString(`OK`)),
				Header:     make(http.Header),
			}
		}

		testManagerClient := operation.NewManagerClient("")
		testManagerClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		op, err := testManagerClient.MarkRunning(testOperationId, fencingToken)

		assert.Error(tt, err)
		assert.Nil(tt, op)
	})

	t.Run("RequestFailure", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), operation.OperationsEndpoint+"/"+testOperationId+"/run")
			assert.Equal(tt, "POST", req.Method)

			return nil
		}

		testManagerClient := operation.NewManagerClient("")
		testManagerClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		op, err := testManagerClient.MarkRunning(testOperationId, fencingToken)

		assert.Error(tt, err)
		assert.Nil(tt, op)
	})
}

func TestManagerClient_ForceUnlock(t *testing.T) {
	const (
		actor  = "admin"
		reason = "stale lock"
	)

	t.Run("OperationForceUnlocked", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), operation.OperationsEndpoint+"/"+testOperationId)
			assert.Equal(tt, "DELETE", req.Method)

			body := `{"operation":{"id":"03cb753f-5e03-4c97-8e47-625115476c72","type":"deploy","system":"node","status":"CANCELLED","fencingToken":1,"resourceKey":"node/03cb753f-5e03-4c97-8e47-625115476c72","leaseExpiresAt":"2026-06-29T12:00:00Z","createdAt":"2026-06-29T11:00:00Z"}}`

			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",
				Body:       io.NopCloser(bytes.NewBufferString(body)),
				Header:     make(http.Header),
			}
		}

		testManagerClient := operation.NewManagerClient("")
		testManagerClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		op, err := testManagerClient.ForceUnlock(testOperationId, actor, reason)

		assert.NoError(tt, err)
		assert.Equal(tt, testOperationId, op.Id)
		assert.Equal(tt, operation.StatusCancelled, op.Status)
	})

	t.Run("ForceUnlockFailed", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), operation.OperationsEndpoint+"/"+testOperationId)
			assert.Equal(tt, "DELETE", req.Method)

			resp := `{"error":"not found"}`

			return &http.Response{
				StatusCode: 404,
				Status:     "404 NOT FOUND",
				Body:       io.NopCloser(bytes.NewBufferString(resp)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}
		}

		testManagerClient := operation.NewManagerClient("")
		testManagerClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		op, err := testManagerClient.ForceUnlock(testOperationId, actor, reason)

		assert.Error(tt, err)
		assert.Nil(tt, op)
	})

	t.Run("InvalidResponsePayload", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), operation.OperationsEndpoint+"/"+testOperationId)
			assert.Equal(tt, "DELETE", req.Method)

			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",
				Body:       io.NopCloser(bytes.NewBufferString(`OK`)),
				Header:     make(http.Header),
			}
		}

		testManagerClient := operation.NewManagerClient("")
		testManagerClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		op, err := testManagerClient.ForceUnlock(testOperationId, actor, reason)

		assert.Error(tt, err)
		assert.Nil(tt, op)
	})

	t.Run("RequestFailure", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), operation.OperationsEndpoint+"/"+testOperationId)
			assert.Equal(tt, "DELETE", req.Method)

			return nil
		}

		testManagerClient := operation.NewManagerClient("")
		testManagerClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		op, err := testManagerClient.ForceUnlock(testOperationId, actor, reason)

		assert.Error(tt, err)
		assert.Nil(tt, op)
	})
}
