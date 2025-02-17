/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package ukamaagent_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/tj/assert"
	"github.com/ukama/ukama/systems/common/rest/client"
	"github.com/ukama/ukama/systems/common/rest/client/ukamaagent"
)

const (
	testIccid = "890000000000000001234"
	testImsi  = "20000233489900"
	cdrType   = "data"
	from      = "2022-12-01T00:00:00Z"
	to        = "2023-12-01T00:00:00Z"
	region    = "meedan"
	bytesUsed = 28901234567
	cost      = 100.99
)

var req = client.AgentRequestData{
	Iccid:     testIccid,
	Imsi:      testImsi,
	NetworkId: "5248eefa-23a0-4222-b80b-e1af5047eaf8",
	SimId:     "0ba2f8d9-e888-4071-aa09-7300daa986aa",
}

func TestUkamaClient_GetSimInfo(t *testing.T) {
	t.Run("SimFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			// Test request parameters
			assert.Equal(tt, req.URL.String(), ukamaagent.UkamaSimsEndpoint+"/"+testIccid)

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

		testUkamaClient := ukamaagent.NewUkamaAgentClient("")

		// We replace the transport mechanism by mocking the http request
		// so that the test stays a unit test e.g no server/network call.
		testUkamaClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		s, err := testUkamaClient.GetSimInfo(testIccid)

		assert.NoError(tt, err)
		assert.Equal(tt, testIccid, s.Iccid)
	})

	t.Run("SimNotFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), ukamaagent.UkamaSimsEndpoint+"/"+testIccid)

			// error payload
			resp := `{"error":"not found"}`

			return &http.Response{
				StatusCode: 404,
				Body:       io.NopCloser(bytes.NewBufferString(resp)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}
		}

		testUkamaClient := ukamaagent.NewUkamaAgentClient("")

		testUkamaClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		s, err := testUkamaClient.GetSimInfo(testIccid)

		assert.Error(tt, err)
		assert.Nil(tt, s)
	})

	t.Run("InvalidResponsePayload", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), ukamaagent.UkamaSimsEndpoint+"/"+testIccid)

			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewBufferString(`OK`)),
				Header:     make(http.Header),
			}
		}

		testUkamaClient := ukamaagent.NewUkamaAgentClient("")

		testUkamaClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		s, err := testUkamaClient.GetSimInfo(testIccid)

		assert.Error(tt, err)
		assert.Nil(tt, s)
	})

	t.Run("RequestFailure", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), ukamaagent.UkamaSimsEndpoint+"/"+testIccid)

			return nil
		}

		testUkamaClient := ukamaagent.NewUkamaAgentClient("")

		testUkamaClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		s, err := testUkamaClient.GetSimInfo(testIccid)

		assert.Error(tt, err)
		assert.Nil(tt, s)
	})
}

func TestUkamaClient_GetUsages(t *testing.T) {
	t.Run("UsageFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			// Test request parameters
			assert.Equal(tt, req.URL.String(), fmt.Sprintf("%s?iccid=%s&cdr_type=%s&from=%s&to=%s&region=%s",
				ukamaagent.UkamaSimsEndpoint+"/usage/", testIccid, cdrType, from, to, region))

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

		testUkamaClient := ukamaagent.NewUkamaAgentClient("")

		// We replace the transport mechanism by mocking the http request
		// so that the test stays a unit test e.g no server/network call.
		testUkamaClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		u, c, err := testUkamaClient.GetUsages(testIccid, cdrType, from, to, region)

		assert.NoError(tt, err)
		assert.NotNil(tt, u[testIccid])
		assert.NotNil(tt, c[testIccid])
		assert.Equal(tt, 100.99, c[testIccid])
	})

	t.Run("InvalidResponsePayload", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), fmt.Sprintf("%s?iccid=%s&cdr_type=%s&from=%s&to=%s&region=%s",
				ukamaagent.UkamaSimsEndpoint+"/usage/", testIccid, cdrType, from, to, region))

			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewBufferString(`OK`)),
				Header:     make(http.Header),
			}
		}

		testUkamaClient := ukamaagent.NewUkamaAgentClient("")

		testUkamaClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		u, c, err := testUkamaClient.GetUsages(testIccid, cdrType, from, to, region)

		assert.Error(tt, err)
		assert.Nil(tt, u)
		assert.Nil(tt, c)
	})

	t.Run("RequestFailure", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), fmt.Sprintf("%s?iccid=%s&cdr_type=%s&from=%s&to=%s&region=%s",
				ukamaagent.UkamaSimsEndpoint+"/usage/", testIccid, cdrType, from, to, region))

			return nil
		}

		testUkamaClient := ukamaagent.NewUkamaAgentClient("")

		testUkamaClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		u, c, err := testUkamaClient.GetUsages(testIccid, cdrType, from, to, region)

		assert.Error(tt, err)
		assert.Nil(tt, u)
		assert.Nil(tt, c)
	})
}

func TestUkamaClient_ActivateSim(t *testing.T) {
	t.Run("SimFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), ukamaagent.UkamaSimsEndpoint+"/"+testIccid)

			return &http.Response{
				StatusCode: 201,
				Header:     make(http.Header),
			}
		}

		testUkamaClient := ukamaagent.NewUkamaAgentClient("")

		testUkamaClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		err := testUkamaClient.ActivateSim(req)

		assert.NoError(tt, err)
	})

	t.Run("RequestFailure", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), ukamaagent.UkamaSimsEndpoint+"/"+testIccid)

			return nil
		}

		testUkamaClient := ukamaagent.NewUkamaAgentClient("")

		testUkamaClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		err := testUkamaClient.ActivateSim(req)

		assert.Error(tt, err)
	})
}

func TestUkamaClient_DeactivateSim(t *testing.T) {
	t.Run("SimFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), ukamaagent.UkamaSimsEndpoint+"/"+testIccid)

			return &http.Response{
				StatusCode: 200,
				Header:     make(http.Header),
			}
		}

		testUkamaClient := ukamaagent.NewUkamaAgentClient("")

		testUkamaClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		err := testUkamaClient.DeactivateSim(req)

		assert.NoError(tt, err)
	})

	t.Run("RequestFailure", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), ukamaagent.UkamaSimsEndpoint+"/"+testIccid)

			return nil
		}

		testUkamaClient := ukamaagent.NewUkamaAgentClient("")

		testUkamaClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		err := testUkamaClient.DeactivateSim(req)

		assert.Error(tt, err)
	})
}

func TestUkamaClient_UpdateSimPackage(t *testing.T) {
	t.Run("SimFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), ukamaagent.UkamaSimsEndpoint+"/"+testIccid)

			return &http.Response{
				StatusCode: 200,
				Header:     make(http.Header),
			}
		}

		testUkamaClient := ukamaagent.NewUkamaAgentClient("")

		testUkamaClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		err := testUkamaClient.UpdatePackage(req)

		assert.NoError(tt, err)
	})

	t.Run("RequestFailure", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), ukamaagent.UkamaSimsEndpoint+"/"+testIccid)

			return nil
		}

		testUkamaClient := ukamaagent.NewUkamaAgentClient("")

		testUkamaClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		err := testUkamaClient.UpdatePackage(req)

		assert.Error(tt, err)
	})
}

type RoundTripFunc func(req *http.Request) *http.Response

func (r RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return r(req), nil
}
