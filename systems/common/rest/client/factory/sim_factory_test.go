/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package factory_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"testing"

	"github.com/tj/assert"

	"github.com/ukama/ukama/systems/common/rest/client"
	"github.com/ukama/ukama/systems/common/rest/client/factory"
)

const (
	testSimIccid = "0123456789012345678912"
	testImsi     = "012345678912345"
	testBatchId  = "batchId"
	testOrgName  = "test-org"
	testcount    = 0
	testSort     = false
	simType      = "test"
)

func TestNewSimFactoryClient(t *testing.T) {
	t.Run("ValidHost", func(tt *testing.T) {
		client := factory.NewSimFactoryClient(testFactoryHost)

		assert.NotNil(tt, client)
		assert.NotNil(tt, client.R)
	})

	t.Run("WithOptions", func(tt *testing.T) {
		client := factory.NewSimFactoryClient(testFactoryHost, client.WithDebug(true))

		assert.NotNil(tt, client)
		assert.NotNil(tt, client.R)
	})
}

func TestSimFactory_ReadSim(t *testing.T) {

	simInfo := factory.SimCardInfo{
		Iccid:          testSimIccid,
		Imsi:           testImsi,
		Op:             []byte("0123456789012345"),
		Key:            []byte("0123456789012345"),
		Amf:            []byte("800"),
		AlgoType:       1,
		UeDlAmbrBps:    2000000,
		UeUlAmbrBps:    2000000,
		Sqn:            1,
		CsgIdPrsent:    false,
		CsgId:          0,
		DefaultApnName: "ukama",
	}

	sim := factory.Sim{
		SimCardInfo: &simInfo,
	}

	t.Run("ReadSimCardInfo_Success", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			// Test request parameters
			expectedURL := testFactoryHost + factory.SimFactoryEndpoint + "/" + testSimIccid
			assert.Equal(tt, expectedURL, req.URL.String())
			assert.Equal(tt, "GET", req.Method)

			// Serialize expected response
			responseBody, _ := json.Marshal(sim)

			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",
				Body:       io.NopCloser(bytes.NewBuffer(responseBody)),
				Header:     make(http.Header),
			}
		}

		// Act
		testClient := factory.NewSimFactoryClient(testFactoryHost)
		testClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		s, err := testClient.ReadSimCardInfo(testSimIccid)

		assert.NoError(t, err)
		assert.NotNil(t, s)
		assert.Equal(t, s.Iccid, testSimIccid)
	})

	t.Run("SimNotFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			expectedURL := testFactoryHost + factory.SimFactoryEndpoint + "/" + testSimIccid
			assert.Equal(tt, expectedURL, req.URL.String())

			errorResponse := `{"error": "node not found"}`

			return &http.Response{
				StatusCode: 404,
				Status:     "404 NOT FOUND",
				Body:       io.NopCloser(bytes.NewBufferString(errorResponse)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}
		}

		testClient := factory.NewSimFactoryClient(testFactoryHost)
		testClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		result, err := testClient.ReadSimCardInfo(testSimIccid)

		assert.Error(tt, err)
		assert.Nil(tt, result)
	})

	t.Run("InvalidResponsePayload", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			expectedURL := testFactoryHost + factory.SimFactoryEndpoint + "/" + testSimIccid
			assert.Equal(tt, expectedURL, req.URL.String())

			// Return invalid JSON
			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",
				Body:       io.NopCloser(bytes.NewBufferString(`{"invalid": json`)),
				Header:     make(http.Header),
			}
		}

		testClient := factory.NewSimFactoryClient(testFactoryHost)
		testClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		result, err := testClient.ReadSimCardInfo(testSimIccid)

		assert.Error(tt, err)
		assert.Nil(tt, result)
	})

	t.Run("RequestFailure", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			expectedURL := testFactoryHost + factory.SimFactoryEndpoint + "/" + testSimIccid
			assert.Equal(tt, expectedURL, req.URL.String())

			// Return nil to simulate network failure
			return nil
		}

		testClient := factory.NewSimFactoryClient(testFactoryHost)
		testClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		result, err := testClient.ReadSimCardInfo(testSimIccid)

		assert.Error(tt, err)
		assert.Nil(tt, result)
	})

	t.Run("EmptyResponseBody", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			expectedURL := testFactoryHost + factory.SimFactoryEndpoint + "/" + testSimIccid
			assert.Equal(tt, expectedURL, req.URL.String())

			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",
				Body:       io.NopCloser(bytes.NewBufferString("")),
				Header:     make(http.Header),
			}
		}

		testClient := factory.NewSimFactoryClient(testFactoryHost)
		testClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		result, err := testClient.ReadSimCardInfo(testSimIccid)

		assert.Error(tt, err)
		assert.Nil(tt, result)
	})
}

func TestNodeFactoryClient_ListSims(t *testing.T) {
	t.Run("ListSimsSuccess", func(tt *testing.T) {
		expectedSim := factory.Sims{
			Sims: []*factory.SimCardInfo{
				{
					Iccid:   testSimIccid,
					Imsi:    testImsi,
					BatchId: testBatchId,
					SimType: simType,
					OrgName: testOrgName,
				},
			},
		}

		mockTransport := func(req *http.Request) *http.Response {
			expectedURL := testFactoryHost + fmt.Sprintf("%s?imsi=%s&batch_id=%s&org_name=%s&sim_type=%s&count=%d&sort=%s",
				factory.SimFactoryEndpoint, testImsi, testBatchId, testOrgName, simType, testcount, strconv.FormatBool(testSort))
			assert.Equal(tt, expectedURL, req.URL.String())

			// Serialize expected response
			responseBody, _ := json.Marshal(expectedSim)

			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",
				Body:       io.NopCloser(bytes.NewBuffer(responseBody)),
				Header:     make(http.Header),
			}
		}

		testClient := factory.NewSimFactoryClient(testFactoryHost)
		testClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		sims, err := testClient.ListSims(testImsi, testBatchId, testOrgName, simType, 0, false)

		assert.NoError(tt, err)
		assert.NotNil(tt, sims)
		assert.Equal(tt, expectedSim.Sims[0].Iccid, sims[0].Iccid)
		assert.Equal(tt, expectedSim.Sims[0].Imsi, sims[0].Imsi)
		assert.Equal(tt, expectedSim.Sims[0].OrgName, sims[0].OrgName)
		assert.Equal(tt, expectedSim.Sims[0].BatchId, sims[0].BatchId)
		assert.Equal(tt, expectedSim.Sims[0].SimType, sims[0].SimType)
	})

	t.Run("ListSimsNotFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			expectedURL := testFactoryHost + fmt.Sprintf("%s?imsi=%s&batch_id=%s&org_name=%s&sim_type=%s&count=%d&sort=%s",
				factory.SimFactoryEndpoint, testImsi, testBatchId, testOrgName, simType, testcount, strconv.FormatBool(testSort))
			assert.Equal(tt, expectedURL, req.URL.String())

			return &http.Response{
				StatusCode: 404,
				Status:     "404 NOT FOUND",
				Body:       io.NopCloser(bytes.NewBufferString(`{"error": "nodes not found"}`)),
				Header:     make(http.Header),
			}
		}

		testClient := factory.NewSimFactoryClient(testFactoryHost)
		testClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		sims, err := testClient.ListSims(testImsi, testBatchId, testOrgName, simType, 0, false)

		assert.Error(tt, err)
		assert.Nil(tt, sims)
	})

	t.Run("InvalidResponsePayload", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			expectedURL := testFactoryHost + fmt.Sprintf("%s?imsi=%s&batch_id=%s&org_name=%s&sim_type=%s&count=%d&sort=%s",
				factory.SimFactoryEndpoint, testImsi, testBatchId, testOrgName, simType, testcount, strconv.FormatBool(testSort))
			assert.Equal(tt, expectedURL, req.URL.String())

			// Return invalid JSON
			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",
				Body:       io.NopCloser(bytes.NewBufferString(`{"invalid": json`)),
				Header:     make(http.Header),
			}
		}

		testClient := factory.NewSimFactoryClient(testFactoryHost)
		testClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		sims, err := testClient.ListSims(testImsi, testBatchId, testOrgName, simType, 0, false)

		assert.Error(tt, err)
		assert.Nil(tt, sims)
	})

	t.Run("RequestFailure", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			expectedURL := testFactoryHost + fmt.Sprintf("%s?imsi=%s&batch_id=%s&org_name=%s&sim_type=%s&count=%d&sort=%s",
				factory.SimFactoryEndpoint, testImsi, testBatchId, testOrgName, simType, testcount, strconv.FormatBool(testSort))
			assert.Equal(tt, expectedURL, req.URL.String())

			// Return nil to simulate network failure
			return nil
		}

		testClient := factory.NewSimFactoryClient(testFactoryHost)
		testClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		sims, err := testClient.ListSims(testImsi, testBatchId, testOrgName, simType, 0, false)

		assert.Error(tt, err)
		assert.Nil(tt, sims)
	})

	t.Run("EmptyResponseBody", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			expectedURL := testFactoryHost + fmt.Sprintf("%s?imsi=%s&batch_id=%s&org_name=%s&sim_type=%s&count=%d&sort=%s",
				factory.SimFactoryEndpoint, testImsi, testBatchId, testOrgName, simType, testcount, strconv.FormatBool(testSort))
			assert.Equal(tt, expectedURL, req.URL.String())

			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",
				Body:       io.NopCloser(bytes.NewBufferString("")),
				Header:     make(http.Header),
			}
		}

		testClient := factory.NewSimFactoryClient(testFactoryHost)
		testClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		sims, err := testClient.ListSims(testImsi, testBatchId, testOrgName, simType, 0, false)

		assert.Error(tt, err)
		assert.Nil(tt, sims)
	})

	t.Run("UnknownError", func(tt *testing.T) {
		var expectedSim factory.Sims

		mockTransport := func(req *http.Request) *http.Response {
			expectedURL := testFactoryHost + fmt.Sprintf("%s?imsi=%s&batch_id=%s&org_name=%s&sim_type=%s&count=%d&sort=%s",
				factory.SimFactoryEndpoint, testImsi, testBatchId, testOrgName, simType, testcount, strconv.FormatBool(testSort))
			assert.Equal(tt, expectedURL, req.URL.String())

			// Serialize expected response
			responseBody, _ := json.Marshal(expectedSim)

			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",
				Body:       io.NopCloser(bytes.NewBuffer(responseBody)),
				Header:     make(http.Header),
			}
		}

		testClient := factory.NewSimFactoryClient(testFactoryHost)
		testClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		sims, err := testClient.ListSims(testImsi, testBatchId, testOrgName, simType, 0, false)

		assert.Error(tt, err)
		assert.Nil(tt, sims)
	})
}
