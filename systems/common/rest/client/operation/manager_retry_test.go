/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package operation

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/tj/assert"

	"github.com/ukama/ukama/systems/common/rest/client"
)

type scriptedTransport struct {
	replies []int
	calls   int
}

func (s *scriptedTransport) RoundTrip(_ *http.Request) (*http.Response, error) {
	code := s.replies[s.calls]
	s.calls++
	if code == 0 {
		return nil, errors.New("connection reset")
	}
	body := `{"operation":{"id":"op","status":"SUCCESS"}}`
	if code != http.StatusOK {
		body = `{"error":"failed"}`
	}
	return &http.Response{
		StatusCode: code,
		Body:       io.NopCloser(bytes.NewBufferString(body)),
		Header:     http.Header{"Content-Type": []string{"application/json"}},
	}, nil
}

func newRetryClient(replies ...int) (*managerClient, *scriptedTransport) {
	t := &scriptedTransport{replies: replies}
	m := NewManagerClient("")
	m.backoff = 0
	m.R.C.SetTransport(t)
	return m, t
}

func TestManagerClient_TerminateRetry(t *testing.T) {
	calls := map[string]func(m *managerClient) (*OperationInfo, error){
		"Complete": func(m *managerClient) (*OperationInfo, error) {
			return m.Complete("op", 1, "controller", "dispatched")
		},
		"Fail": func(m *managerClient) (*OperationInfo, error) {
			return m.Fail("op", 1, "controller", "publish failed")
		},
	}

	for name, call := range calls {
		t.Run(name+"RetriesUnavailable", func(t *testing.T) {
			m, tr := newRetryClient(http.StatusServiceUnavailable, http.StatusOK)

			op, err := call(m)

			assert.NoError(t, err)
			assert.Equal(t, StatusSuccess, op.Status)
			assert.Equal(t, 2, tr.calls)
		})

		t.Run(name+"RetriesTransportError", func(t *testing.T) {
			m, tr := newRetryClient(0, http.StatusOK)

			_, err := call(m)

			assert.NoError(t, err)
			assert.Equal(t, 2, tr.calls)
		})

		t.Run(name+"GivesUpAfterThreeAttempts", func(t *testing.T) {
			m, tr := newRetryClient(http.StatusGatewayTimeout,
				http.StatusGatewayTimeout, http.StatusGatewayTimeout)

			_, err := call(m)

			assert.Error(t, err)
			assert.Equal(t, terminateAttempts, tr.calls)
		})

		t.Run(name+"DoesNotRetryRejection", func(t *testing.T) {
			m, tr := newRetryClient(http.StatusBadRequest)

			_, err := call(m)

			var status *client.ErrorStatus
			assert.True(t, errors.As(err, &status), err)
			assert.Equal(t, 1, tr.calls)
		})

		t.Run(name+"DoesNotRetryInternal", func(t *testing.T) {
			m, tr := newRetryClient(http.StatusInternalServerError)

			_, err := call(m)

			assert.Error(t, err)
			assert.Equal(t, 1, tr.calls)
		})
	}
}
