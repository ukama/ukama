/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package dataplan_test

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/tj/assert"

	"github.com/ukama/ukama/systems/common/rest/client/dataplan"
	"github.com/ukama/ukama/systems/common/uuid"
)

const testUuid = "03cb753f-5e03-4c97-8e47-625115476c72"

func TestPackageClient_Get(t *testing.T) {
	t.Run("PackageFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			// Test request parameters
			assert.Equal(tt, req.URL.String(), dataplan.PackageEndpoint+"/"+testUuid)

			// fake package info
			pkg := `{"package":{"uuid": "03cb753f-5e03-4c97-8e47-625115476c72", "active": true}}`

			// Send mock response
			return &http.Response{
				StatusCode: 200,

				// Send response to be tested
				Body: io.NopCloser(bytes.NewBufferString(pkg)),

				// Must be set to non-nil value or it panics
				Header: make(http.Header),
			}
		}

		testPackageClient := dataplan.NewPackageClient("")

		// We replace the transport mechanism by mocking the http request
		// so that the test stays a unit test e.g, no server/network call.
		testPackageClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		p, err := testPackageClient.Get(testUuid)

		assert.NoError(tt, err)
		assert.Equal(tt, testUuid, p.Id)
	})

	t.Run("PackageNotFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), dataplan.PackageEndpoint+"/"+testUuid)

			// error payload
			resp := `{"error":"not found"}`

			return &http.Response{
				StatusCode: 404,
				Body:       io.NopCloser(bytes.NewBufferString(resp)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}
		}

		testPackageClient := dataplan.NewPackageClient("")

		testPackageClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		p, err := testPackageClient.Get(testUuid)

		assert.Error(tt, err)
		assert.Nil(tt, p)
	})

	t.Run("InvalidResponsePayload", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), dataplan.PackageEndpoint+"/"+testUuid)

			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewBufferString(`OK`)),
				Header:     make(http.Header),
			}
		}

		testPackageClient := dataplan.NewPackageClient("")

		testPackageClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		p, err := testPackageClient.Get(testUuid)

		assert.Error(tt, err)
		assert.Nil(tt, p)
	})

	t.Run("RequestFailure", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), dataplan.PackageEndpoint+"/"+testUuid)

			return nil
		}

		testPackageClient := dataplan.NewPackageClient("")

		testPackageClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		p, err := testPackageClient.Get(testUuid)

		assert.Error(tt, err)
		assert.Nil(tt, p)
	})
}

func TestPackageClient_Add(t *testing.T) {
	t.Run("PackageAdded", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			// Test request parameters
			assert.Equal(tt, req.URL.String(), dataplan.PackageEndpoint)

			// fake package info
			pkg := `{"package":{"uuid": "03cb753f-5e03-4c97-8e47-625115476c72", "active": true}}`

			// Send mock response
			return &http.Response{
				StatusCode: 201,
				Status:     "201 CREATED",

				// Send response to be tested
				Body: io.NopCloser(bytes.NewBufferString(pkg)),

				// Must be set to non-nil value or it panics
				Header: make(http.Header),
			}
		}

		testPackageClient := dataplan.NewPackageClient("")

		// We replace the transport mechanism by mocking the http request
		// so that the test stays a unit test e.g, no server/network call.
		testPackageClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		p, err := testPackageClient.Add(
			dataplan.AddPackageRequest{
				Name:        "Monthly Data",
				OrgId:       uuid.NewV4().String(),
				OwnerId:     uuid.NewV4().String(),
				From:        "2023-04-01T00:00:00Z",
				To:          "2023-04-01T00:00:00Z",
				BaserateId:  uuid.NewV4().String(),
				VoiceVolume: 0,
				Active:      true,
				DataVolume:  1024,
				SmsVolume:   0,
				DataUnit:    "MegaBytes",
				VoiceUnit:   "seconds",
				SimType:     "test",
				Apn:         "ukama.tel",
				Markup:      0,
				Type:        "postpaid",
				Flatrate:    false,
				Amount:      0,
			},
		)

		assert.NoError(tt, err)
		assert.Equal(tt, testUuid, p.Id)
	})

	t.Run("InvalidResponseHeader", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), dataplan.PackageEndpoint)

			// error payload
			resp := `{"error":"internal server error"}`

			return &http.Response{
				StatusCode: 500,
				Status:     "500 INTERNAL SERVER ERROR",
				Body:       io.NopCloser(bytes.NewBufferString(resp)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}
		}

		testPackageClient := dataplan.NewPackageClient("")

		testPackageClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		p, err := testPackageClient.Add(
			dataplan.AddPackageRequest{
				Name:        "Monthly Data",
				OrgId:       uuid.NewV4().String(),
				OwnerId:     uuid.NewV4().String(),
				From:        "2023-04-01T00:00:00Z",
				To:          "2023-04-01T00:00:00Z",
				BaserateId:  uuid.NewV4().String(),
				VoiceVolume: 0,
				Active:      true,
				DataVolume:  1024,
				SmsVolume:   0,
				DataUnit:    "MegaBytes",
				VoiceUnit:   "seconds",
				SimType:     "test",
				Apn:         "ukama.tel",
				Markup:      0,
				Type:        "postpaid",
				Flatrate:    false,
				Amount:      0,
			},
		)

		assert.Error(tt, err)
		assert.Nil(tt, p)
	})

	t.Run("InvalidResponsePayload", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), dataplan.PackageEndpoint)

			return &http.Response{
				StatusCode: 201,
				Status:     "201 CREATED",
				Body:       io.NopCloser(bytes.NewBufferString(`CREATED`)),
				Header:     make(http.Header),
			}
		}

		testPackageClient := dataplan.NewPackageClient("")

		testPackageClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		p, err := testPackageClient.Add(
			dataplan.AddPackageRequest{
				Name:        "Monthly Data",
				OrgId:       uuid.NewV4().String(),
				OwnerId:     uuid.NewV4().String(),
				From:        "2023-04-01T00:00:00Z",
				To:          "2023-04-01T00:00:00Z",
				BaserateId:  uuid.NewV4().String(),
				VoiceVolume: 0,
				Active:      true,
				DataVolume:  1024,
				SmsVolume:   0,
				DataUnit:    "MegaBytes",
				VoiceUnit:   "seconds",
				SimType:     "test",
				Apn:         "ukama.tel",
				Markup:      0,
				Type:        "postpaid",
				Flatrate:    false,
				Amount:      0,
			},
		)

		assert.Error(tt, err)
		assert.Nil(tt, p)
	})

	t.Run("RequestFailure", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), dataplan.PackageEndpoint)

			return nil
		}

		testPackageClient := dataplan.NewPackageClient("")

		testPackageClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		p, err := testPackageClient.Add(
			dataplan.AddPackageRequest{
				Name:        "Monthly Data",
				OrgId:       uuid.NewV4().String(),
				OwnerId:     uuid.NewV4().String(),
				From:        "2023-04-01T00:00:00Z",
				To:          "2023-04-01T00:00:00Z",
				BaserateId:  uuid.NewV4().String(),
				VoiceVolume: 0,
				Active:      true,
				DataVolume:  1024,
				SmsVolume:   0,
				DataUnit:    "MegaBytes",
				VoiceUnit:   "seconds",
				SimType:     "test",
				Apn:         "ukama.tel",
				Markup:      0,
				Type:        "postpaid",
				Flatrate:    false,
				Amount:      0,
			},
		)

		assert.Error(tt, err)
		assert.Nil(tt, p)
	})
}

type RoundTripFunc func(req *http.Request) *http.Response

func (r RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return r(req), nil
}
