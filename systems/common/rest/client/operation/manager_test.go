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
	"errors"
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
			sent, _ := io.ReadAll(req.Body)
			assert.Contains(tt, string(sent), `"resource_keys":["`+testResourceKey+`"]`)

			body := `{"operations":[{"id":"03cb753f-5e03-4c97-8e47-625115476c72","type":"deploy","system":"node","status":"PENDING","fencing_token":1,"resource_key":"node/03cb753f-5e03-4c97-8e47-625115476c72","lease_expires_at":"2026-06-29T12:00:00Z","created_at":"2026-06-29T11:00:00Z"}]}`

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
			Type:         "deploy",
			System:       "node",
			ResourceKeys: []string{testResourceKey},
		})

		assert.NoError(tt, err)
		assert.Len(tt, resp.Operations, 1)
		assert.Equal(tt, testOperationId, resp.Operations[0].Id)
		assert.Equal(tt, operation.StatusPending, resp.Operations[0].Status)
		assert.Equal(tt, uint64(1), resp.Operations[0].FencingToken)
		assert.Equal(tt, testResourceKey, resp.Operations[0].ResourceKey)
	})

	t.Run("BusyIsBusyError", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), operation.OperationsEndpoint)
			assert.Equal(tt, "POST", req.Method)

			resp := `{"error":"node:c is busy with RestartNode operation op-1"}`

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
			Type:         "deploy",
			System:       "node",
			ResourceKeys: []string{testResourceKey},
		})

		var busy *operation.BusyError
		assert.True(tt, errors.As(err, &busy), err)
		assert.True(tt, errors.Is(err, operation.ErrBusy), err)
		assert.Equal(tt, "node:c is busy with RestartNode operation op-1", err.Error())
		assert.Nil(tt, resp)
	})

	t.Run("ServerErrorIsNotBusy", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			return &http.Response{
				StatusCode: 500,
				Body:       io.NopCloser(bytes.NewBufferString(`{"error":"db down"}`)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}
		}

		testManagerClient := operation.NewManagerClient("")
		testManagerClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		resp, err := testManagerClient.Start(operation.StartRequest{
			Type:         "deploy",
			System:       "node",
			ResourceKeys: []string{testResourceKey},
		})

		assert.Error(tt, err)
		assert.False(tt, errors.Is(err, operation.ErrBusy), err)
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
			Type:         "deploy",
			System:       "node",
			ResourceKeys: []string{testResourceKey},
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
			Type:         "deploy",
			System:       "node",
			ResourceKeys: []string{testResourceKey},
		})

		assert.Error(tt, err)
		assert.Nil(tt, resp)
	})
}

func TestManagerClient_Get(t *testing.T) {
	t.Run("OperationFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), operation.OperationsEndpoint+"/"+testOperationId)

			body := `{"operation":{"id":"03cb753f-5e03-4c97-8e47-625115476c72","type":"deploy","system":"node","status":"SUCCESS","fencing_token":1,"resource_key":"node/03cb753f-5e03-4c97-8e47-625115476c72","lease_expires_at":"2026-06-29T12:00:00Z","created_at":"2026-06-29T11:00:00Z"}}`

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

			body := `{"operation":{"id":"03cb753f-5e03-4c97-8e47-625115476c72","type":"deploy","system":"node","status":"RUNNING","fencing_token":1,"resource_key":"node/03cb753f-5e03-4c97-8e47-625115476c72","lease_expires_at":"2026-06-29T12:00:00Z","created_at":"2026-06-29T11:00:00Z"}}`

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
			assert.Equal(tt, req.URL.String(),
				operation.OperationsEndpoint+"/"+testOperationId)
			assert.Equal(tt, "PATCH", req.Method)
			sent, _ := io.ReadAll(req.Body)
			assert.Contains(tt, string(sent), `"actor":"controller"`)

			body := `{"operation":{"id":"03cb753f-5e03-4c97-8e47-625115476c72","type":"deploy","system":"node","status":"RUNNING","fencing_token":42,"resource_key":"node/03cb753f-5e03-4c97-8e47-625115476c72","lease_expires_at":"2026-06-29T12:00:00Z","created_at":"2026-06-29T11:00:00Z"}}`

			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",
				Body:       io.NopCloser(bytes.NewBufferString(body)),
				Header:     make(http.Header),
			}
		}

		testManagerClient := operation.NewManagerClient("")
		testManagerClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		op, err := testManagerClient.MarkRunning(testOperationId, fencingToken, "controller")

		assert.NoError(tt, err)
		assert.Equal(tt, testOperationId, op.Id)
		assert.Equal(tt, operation.StatusRunning, op.Status)
		assert.Equal(tt, fencingToken, op.FencingToken)
	})

	for code, want := range map[int]error{
		http.StatusBadRequest:          nil,
		http.StatusNotFound:            operation.ErrNotRunnable,
		http.StatusConflict:            operation.ErrNotRunnable,
		http.StatusPreconditionFailed:  operation.ErrNotRunnable,
		http.StatusServiceUnavailable:  operation.ErrUnavailable,
		http.StatusGatewayTimeout:      operation.ErrUnavailable,
		http.StatusInternalServerError: nil,
	} {
		t.Run("MarkRunningFailed"+http.StatusText(code), func(tt *testing.T) {
			mockTransport := func(req *http.Request) *http.Response {
				body := bytes.NewBufferString(`{"error":"rejected"}`)
				return &http.Response{
					StatusCode: code,
					Body:       io.NopCloser(body),
					Header:     http.Header{"Content-Type": []string{"application/json"}},
				}
			}

			testManagerClient := operation.NewManagerClient("")
			testManagerClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

			op, err := testManagerClient.MarkRunning(testOperationId, fencingToken,
				"controller")

			var status *client.ErrorStatus
			assert.True(tt, errors.As(err, &status), err)
			assert.Equal(tt, code, status.StatusCode)
			assert.Equal(tt, want == operation.ErrNotRunnable,
				errors.Is(err, operation.ErrNotRunnable), err)
			assert.Equal(tt, want == operation.ErrUnavailable,
				errors.Is(err, operation.ErrUnavailable), err)
			assert.Nil(tt, op)
		})
	}

	t.Run("InvalidResponsePayload", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(),
				operation.OperationsEndpoint+"/"+testOperationId)
			assert.Equal(tt, "PATCH", req.Method)

			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",
				Body:       io.NopCloser(bytes.NewBufferString(`OK`)),
				Header:     make(http.Header),
			}
		}

		testManagerClient := operation.NewManagerClient("")
		testManagerClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		op, err := testManagerClient.MarkRunning(testOperationId, fencingToken, "controller")

		assert.Error(tt, err)
		assert.Nil(tt, op)
	})

	t.Run("RequestFailure", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(),
				operation.OperationsEndpoint+"/"+testOperationId)
			assert.Equal(tt, "PATCH", req.Method)

			return nil
		}

		testManagerClient := operation.NewManagerClient("")
		testManagerClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		op, err := testManagerClient.MarkRunning(testOperationId, fencingToken, "controller")

		assert.True(tt, errors.Is(err, operation.ErrUnavailable), err)
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
			assert.Equal(tt, req.URL.String(),
				operation.OperationsEndpoint+"/"+testOperationId+
					"/lock?actor=admin&reason=stale+lock")
			assert.Equal(tt, "DELETE", req.Method)

			body := `{"operation":{"id":"03cb753f-5e03-4c97-8e47-625115476c72","type":"deploy","system":"node","status":"CANCELLED","fencing_token":1,"resource_key":"node/03cb753f-5e03-4c97-8e47-625115476c72","lease_expires_at":"2026-06-29T12:00:00Z","created_at":"2026-06-29T11:00:00Z"}}`

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
			assert.Equal(tt, req.URL.String(),
				operation.OperationsEndpoint+"/"+testOperationId+
					"/lock?actor=admin&reason=stale+lock")
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
			assert.Equal(tt, req.URL.String(),
				operation.OperationsEndpoint+"/"+testOperationId+
					"/lock?actor=admin&reason=stale+lock")
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
			assert.Equal(tt, req.URL.String(),
				operation.OperationsEndpoint+"/"+testOperationId+
					"/lock?actor=admin&reason=stale+lock")
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

func TestManagerClient_CompleteAndFail(t *testing.T) {
	const (
		fencingToken = uint64(7)
		actor        = "controller"
		reason       = "done"
	)

	type terminate func(operation.ManagerClient) (*operation.OperationInfo, error)

	calls := []struct {
		name   string
		status string
		call   terminate
	}{
		{"Complete", "SUCCESS", func(c operation.ManagerClient) (*operation.OperationInfo,
			error) {
			return c.Complete(testOperationId, fencingToken, actor, reason)
		}},
		{"Fail", "FAILED", func(c operation.ManagerClient) (*operation.OperationInfo,
			error) {
			return c.Fail(testOperationId, fencingToken, actor, reason)
		}},
	}

	for _, tc := range calls {
		t.Run(tc.name+"SendsFencingToken", func(tt *testing.T) {
			mockTransport := func(req *http.Request) *http.Response {
				assert.Equal(tt, req.URL.String(),
					operation.OperationsEndpoint+"/"+testOperationId)
				assert.Equal(tt, "PATCH", req.Method)

				sent, err := io.ReadAll(req.Body)
				assert.NoError(tt, err)
				assert.JSONEq(tt, `{"status":"`+tc.status+
					`","fencing_token":7,"actor":"controller","reason":"done"}`,
					string(sent))

				body := `{"operation":{"id":"` + testOperationId + `","fencing_token":7}}`

				return &http.Response{
					StatusCode: 200,
					Status:     "200 OK",
					Body:       io.NopCloser(bytes.NewBufferString(body)),
					Header:     make(http.Header),
				}
			}

			testManagerClient := operation.NewManagerClient("")
			testManagerClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

			op, err := tc.call(testManagerClient)

			assert.NoError(tt, err)
			assert.Equal(tt, testOperationId, op.Id)
		})

		t.Run(tc.name+"RejectedByManager", func(tt *testing.T) {
			mockTransport := func(req *http.Request) *http.Response {
				return &http.Response{
					StatusCode: 412,
					Status:     "412 PRECONDITION FAILED",
					Body:       io.NopCloser(bytes.NewBufferString(`{"error":"mismatch"}`)),
					Header:     http.Header{"Content-Type": []string{"application/json"}},
				}
			}

			testManagerClient := operation.NewManagerClient("")
			testManagerClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

			op, err := tc.call(testManagerClient)

			assert.Error(tt, err)
			assert.Nil(tt, op)
		})

		t.Run(tc.name+"InvalidResponsePayload", func(tt *testing.T) {
			mockTransport := func(req *http.Request) *http.Response {
				return &http.Response{
					StatusCode: 200,
					Status:     "200 OK",
					Body:       io.NopCloser(bytes.NewBufferString(`OK`)),
					Header:     make(http.Header),
				}
			}

			testManagerClient := operation.NewManagerClient("")
			testManagerClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

			op, err := tc.call(testManagerClient)

			assert.Error(tt, err)
			assert.Nil(tt, op)
		})
	}
}

func TestManagerClient_Delete(t *testing.T) {
	t.Run("SendsActorAsQuery", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, operation.OperationsEndpoint+"/"+testOperationId+
				"?actor=admin", req.URL.String())
			assert.Equal(tt, "DELETE", req.Method)

			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString("")),
				Header:     make(http.Header),
			}
		}

		testManagerClient := operation.NewManagerClient("")
		testManagerClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		assert.NoError(tt, testManagerClient.Delete(testOperationId, "admin"))
	})

	t.Run("ActiveOperationIsRejected", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			body := bytes.NewBufferString(`{"error":"operation is RUNNING"}`)
			return &http.Response{
				StatusCode: http.StatusConflict,
				Body:       io.NopCloser(body),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}
		}

		testManagerClient := operation.NewManagerClient("")
		testManagerClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		err := testManagerClient.Delete(testOperationId, "admin")

		var status *client.ErrorStatus
		assert.True(tt, errors.As(err, &status), err)
		assert.Equal(tt, http.StatusConflict, status.StatusCode)
	})
}
