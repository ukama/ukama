/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package operatoragent_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/tj/assert"
	"github.com/ukama/ukama/systems/common/rest/client/operatoragent"
)

const (
	testIccid = "890000000000000001234"
	cdrType   = "data"
	from      = "2022-12-01T00:00:00Z"
	to        = "2023-12-01T00:00:00Z"
	region    = "meedan"
	bytesUsed = 28901234567
	cost      = 100.99
)

func TestOperatorClient_GetSimInfo(t *testing.T) {
	t.Run("SimFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			// Test request parameters
			assert.Equal(tt, req.URL.String(), operatoragent.OperatorSimsEndpoint+"/"+testIccid)

			// fake sim info
			sim := `{"sim":{"iccid": "890000000000000001234", "imsi": "20000233489900"}}`

			// Send mock response
			return &http.Response{
				StatusCode: 200,

				// Send response to be tested
				Body: io.NopCloser(bytes.NewBufferString(sim)),
				// Must be set to non-nil value or it panics
				Header: make(http.Header),
			}
		}

		testOperatorClient := operatoragent.NewOperatorAgentClient("")

		// We replace the transport mechanism by mocking the http request
		// so that the test stays a unit test e.g no server/network call.
		testOperatorClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		s, err := testOperatorClient.GetSimInfo(testIccid)

		assert.NoError(tt, err)
		assert.Equal(tt, testIccid, s.Iccid)
	})

	t.Run("SimNotFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), operatoragent.OperatorSimsEndpoint+"/"+testIccid)

			// error payload
			resp := `{"error":"not found"}`

			return &http.Response{
				StatusCode: 404,
				Body:       io.NopCloser(bytes.NewBufferString(resp)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}
		}

		testOperatorClient := operatoragent.NewOperatorAgentClient("")

		testOperatorClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		s, err := testOperatorClient.GetSimInfo(testIccid)

		assert.Error(tt, err)
		assert.Nil(tt, s)
	})

	t.Run("InvalidResponsePayload", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), operatoragent.OperatorSimsEndpoint+"/"+testIccid)

			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewBufferString(`OK`)),
				Header:     make(http.Header),
			}
		}

		testOperatorClient := operatoragent.NewOperatorAgentClient("")

		testOperatorClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		s, err := testOperatorClient.GetSimInfo(testIccid)

		assert.Error(tt, err)
		assert.Nil(tt, s)
	})

	t.Run("RequestFailure", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), operatoragent.OperatorSimsEndpoint+"/"+testIccid)

			return nil
		}

		testOperatorClient := operatoragent.NewOperatorAgentClient("")

		testOperatorClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		s, err := testOperatorClient.GetSimInfo(testIccid)

		assert.Error(tt, err)
		assert.Nil(tt, s)
	})
}

func TestOperatorClient_GetUsages(t *testing.T) {
	t.Run("UsageFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			// Test request parameters
			assert.Equal(tt, req.URL.String(), fmt.Sprintf("%s?iccid=%s&cdr_type=%s&from=%s&to=%s&region=%s",
				operatoragent.OperatorUsagesEndpoint, testIccid, cdrType, from, to, region))

			// fake usage usage
			usage := `{"usage":{"890000000000000001234": 28901234567}, "cost":{"890000000000000001234":100.99}}`

			// Send mock response
			return &http.Response{
				StatusCode: 200,

				// Send response to be tested
				Body: io.NopCloser(bytes.NewBufferString(usage)),
				// Must be set to non-nil value or it panics
				Header: make(http.Header),
			}
		}

		testOperatorClient := operatoragent.NewOperatorAgentClient("")

		// We replace the transport mechanism by mocking the http request
		// so that the test stays a unit test e.g no server/network call.
		testOperatorClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		u, c, err := testOperatorClient.GetUsages(testIccid, cdrType, from, to, region)

		assert.NoError(tt, err)
		assert.NotNil(tt, u[testIccid])
		assert.NotNil(tt, c[testIccid])
		assert.Equal(tt, 100.99, c[testIccid])
	})

	t.Run("InvalidResponsePayload", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), fmt.Sprintf("%s?iccid=%s&cdr_type=%s&from=%s&to=%s&region=%s",
				operatoragent.OperatorUsagesEndpoint, testIccid, cdrType, from, to, region))

			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewBufferString(`OK`)),
				Header:     make(http.Header),
			}
		}

		testOperatorClient := operatoragent.NewOperatorAgentClient("")

		testOperatorClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		u, c, err := testOperatorClient.GetUsages(testIccid, cdrType, from, to, region)

		assert.Error(tt, err)
		assert.Nil(tt, u)
		assert.Nil(tt, c)
	})

	t.Run("InvalidParameterFrom", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), fmt.Sprintf("%s?iccid=%s&cdr_type=%s&from=%s&to=%s&region=%s",
				operatoragent.OperatorUsagesEndpoint, testIccid, cdrType, "lol", to, region))

			return nil
		}

		testOperatorClient := operatoragent.NewOperatorAgentClient("")

		testOperatorClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		u, c, err := testOperatorClient.GetUsages(testIccid, cdrType, "lol", to, region)

		assert.Error(tt, err)
		assert.Nil(tt, u)
		assert.Nil(tt, c)
	})

	t.Run("InvalidParameterTo", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), fmt.Sprintf("%s?iccid=%s&cdr_type=%s&from=%s&to=%s&region=%s",
				operatoragent.OperatorUsagesEndpoint, testIccid, cdrType, from, "lol", region))

			return nil
		}

		testOperatorClient := operatoragent.NewOperatorAgentClient("")

		testOperatorClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		u, c, err := testOperatorClient.GetUsages(testIccid, cdrType, from, "lol", region)

		assert.Error(tt, err)
		assert.Nil(tt, u)
		assert.Nil(tt, c)
	})

	t.Run("RequestFailure", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), fmt.Sprintf("%s?iccid=%s&cdr_type=%s&from=%s&to=%s&region=%s",
				operatoragent.OperatorUsagesEndpoint, testIccid, cdrType, from, to, region))

			return nil
		}

		testOperatorClient := operatoragent.NewOperatorAgentClient("")

		testOperatorClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		u, c, err := testOperatorClient.GetUsages(testIccid, cdrType, from, to, region)

		assert.Error(tt, err)
		assert.Nil(tt, u)
		assert.Nil(tt, c)
	})
}

func TestOperatorClient_ActivateSim(t *testing.T) {
	t.Run("SimFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), operatoragent.OperatorSimsEndpoint+"/"+testIccid)

			return &http.Response{
				StatusCode: 201,
				Header:     make(http.Header),
			}
		}

		testOperatorClient := operatoragent.NewOperatorAgentClient("")

		testOperatorClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		err := testOperatorClient.ActivateSim(testIccid)

		assert.NoError(tt, err)
	})

	t.Run("SimNotFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), operatoragent.OperatorSimsEndpoint+"/"+testIccid)

			// error payload
			resp := `{"error":"not found"}`

			return &http.Response{
				StatusCode: 404,
				Body:       io.NopCloser(bytes.NewBufferString(resp)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}
		}

		testOperatorClient := operatoragent.NewOperatorAgentClient("")

		testOperatorClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		err := testOperatorClient.ActivateSim(testIccid)

		assert.Error(tt, err)
	})

	t.Run("RequestFailure", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), operatoragent.OperatorSimsEndpoint+"/"+testIccid)

			return nil
		}

		testOperatorClient := operatoragent.NewOperatorAgentClient("")

		testOperatorClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		err := testOperatorClient.ActivateSim(testIccid)

		assert.Error(tt, err)
	})
}

func TestOperatorClient_DeactivateSim(t *testing.T) {
	t.Run("SimFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), operatoragent.OperatorSimsEndpoint+"/"+testIccid)

			return &http.Response{
				StatusCode: 200,
				Header:     make(http.Header),
			}
		}

		testOperatorClient := operatoragent.NewOperatorAgentClient("")

		testOperatorClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		err := testOperatorClient.DeactivateSim(testIccid)

		assert.NoError(tt, err)
	})

	t.Run("SimNotFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), operatoragent.OperatorSimsEndpoint+"/"+testIccid)

			// error payload
			resp := `{"error":"not found"}`

			return &http.Response{
				StatusCode: 404,
				Body:       io.NopCloser(bytes.NewBufferString(resp)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}
		}

		testOperatorClient := operatoragent.NewOperatorAgentClient("")

		testOperatorClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		err := testOperatorClient.DeactivateSim(testIccid)

		assert.Error(tt, err)
	})

	t.Run("RequestFailure", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), operatoragent.OperatorSimsEndpoint+"/"+testIccid)

			return nil
		}

		testOperatorClient := operatoragent.NewOperatorAgentClient("")

		testOperatorClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		err := testOperatorClient.DeactivateSim(testIccid)

		assert.Error(tt, err)
	})
}

func TestOperatorClient_TerminateSim(t *testing.T) {
	t.Run("SimFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), operatoragent.OperatorSimsEndpoint+"/"+testIccid)

			return &http.Response{
				StatusCode: 200,
				Header:     make(http.Header),
			}
		}

		testOperatorClient := operatoragent.NewOperatorAgentClient("")

		testOperatorClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		err := testOperatorClient.TerminateSim(testIccid)

		assert.NoError(tt, err)
	})

	t.Run("SimNotFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), operatoragent.OperatorSimsEndpoint+"/"+testIccid)

			// error payload
			resp := `{"error":"not found"}`

			return &http.Response{
				StatusCode: 404,
				Body:       io.NopCloser(bytes.NewBufferString(resp)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}
		}

		testOperatorClient := operatoragent.NewOperatorAgentClient("")

		testOperatorClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		err := testOperatorClient.TerminateSim(testIccid)

		assert.Error(tt, err)
	})

	t.Run("RequestFailure", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), operatoragent.OperatorSimsEndpoint+"/"+testIccid)

			return nil
		}

		testOperatorClient := operatoragent.NewOperatorAgentClient("")

		testOperatorClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		err := testOperatorClient.TerminateSim(testIccid)

		assert.Error(tt, err)
	})
}

type RoundTripFunc func(req *http.Request) *http.Response

func (r RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return r(req), nil
}
