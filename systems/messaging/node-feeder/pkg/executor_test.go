/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package pkg

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ukama/ukama/systems/messaging/node-feeder/mocks"

	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func answer(code int) roundTripFunc {
	return func(*http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: code, Body: io.NopCloser(strings.NewReader("{}"))},
			nil
	}
}

func newTestExecutor(t *testing.T, resolveErr error, rt http.RoundTripper) *requestExecutor {
	resolver := mocks.NewNodeIpResolver(t)
	resolver.On("Resolve", mock.Anything).Return("10.0.0.1", resolveErr).Maybe()
	e := NewRequestExecutor(resolver, 8082, 1)
	e.transport = rt
	return e
}

func TestExecutor_ReturnsNodeAnswer(t *testing.T) {
	for name, tc := range map[string]struct {
		code    int
		errType any
	}{
		"Accepted":    {http.StatusAccepted, nil},
		"Conflict":    {http.StatusConflict, Device4xxServerError{}},
		"Unavailable": {http.StatusServiceUnavailable, Device5xxServerError{}},
	} {
		t.Run(name, func(t *testing.T) {
			var sent *http.Request
			e := newTestExecutor(t, nil, roundTripFunc(func(r *http.Request) (*http.Response,
				error) {
				sent = r
				return answer(tc.code)(r)
			}))

			code, err := e.Execute(&epb.NodeFeederMessage{Target: testTarget,
				HttpMethod: "POST", Path: "/device/v1/service", Msg: []byte(`{}`),
				NodeId: testNodeId})

			assert.Equal(t, tc.code, code)
			switch tc.errType.(type) {
			case nil:
				assert.NoError(t, err)
			case Device4xxServerError:
				assert.IsType(t, Device4xxServerError{}, err)
			case Device5xxServerError:
				assert.IsType(t, Device5xxServerError{}, err)
			}
			assert.Equal(t, "http://10.0.0.1:8082/device/v1/service", sent.URL.String())
		})
	}
}

func TestExecutor_NodeNotFound(t *testing.T) {
	e := newTestExecutor(t, status.Error(codes.NotFound, "no node"), answer(http.StatusOK))

	code, err := e.Execute(&epb.NodeFeederMessage{Target: testTarget, HttpMethod: "GET",
		Path: "/device/v1/ping"})

	assert.Zero(t, code)
	assert.True(t, errors.Is(err, ErrNodeNotFound))
}

func TestExecutor_Failures(t *testing.T) {
	transportErr := roundTripFunc(func(*http.Request) (*http.Response, error) {
		return nil, errors.New("connection refused")
	})

	for name, tc := range map[string]struct {
		target     string
		resolveErr error
		rt         http.RoundTripper
	}{
		"BadTarget":      {"ukama.org.node", nil, answer(http.StatusOK)},
		"BadNodeId":      {"ukama.org.node.not-a-node", nil, answer(http.StatusOK)},
		"ResolveFails":   {testTarget, errors.New("nns down"), answer(http.StatusOK)},
		"TransportFails": {testTarget, nil, transportErr},
	} {
		t.Run(name, func(t *testing.T) {
			e := newTestExecutor(t, tc.resolveErr, tc.rt)

			code, err := e.Execute(&epb.NodeFeederMessage{Target: tc.target,
				HttpMethod: "GET", Path: "/device/v1/ping"})

			assert.Zero(t, code)
			assert.Error(t, err)
			assert.False(t, errors.Is(err, ErrNodeNotFound))
		})
	}
}
