/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package subscriber_test

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/tj/assert"

	"github.com/ukama/ukama/systems/common/rest/client"
	"github.com/ukama/ukama/systems/common/rest/client/subscriber"
)

const (
	testEmail = "foo@example.com"
)

func TestSubscriberClient_Get(t *testing.T) {
	t.Run("SubscriberFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			// Test request parameters
			assert.Equal(tt, req.URL.String(), subscriber.SubscriberEndpoint+"/"+testUuid)

			// fake subscriber info
			subscriber := `{"subscriber":{"subscriber_id": "03cb753f-5e03-4c97-8e47-625115476c72", "last_name": "Foo"}}`

			// Send mock response
			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",

				// Send response to be tested
				Body: io.NopCloser(bytes.NewBufferString(subscriber)),

				// Must be set to non-nil value or it panics
				Header: make(http.Header),
			}
		}

		testSubscriberClient := subscriber.NewSubscriberClient("")

		// We replace the transport mechanism by mocking the http request
		// so that the test stays a unit test e.g no server/network call.
		testSubscriberClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		s, err := testSubscriberClient.Get(testUuid)

		assert.NoError(tt, err)
		assert.Equal(tt, testUuid, s.SubscriberId.String())
	})

	t.Run("SubscriberNotFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), subscriber.SubscriberEndpoint+"/"+testUuid)

			// error payload
			resp := `{"error":"not found"}`

			return &http.Response{
				StatusCode: 404,
				Status:     "404 NOT FOUND",
				Body:       io.NopCloser(bytes.NewBufferString(resp)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}
		}

		testSubscriberClient := subscriber.NewSubscriberClient("")

		testSubscriberClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		s, err := testSubscriberClient.Get(testUuid)

		assert.Error(tt, err)
		assert.Nil(tt, s)
	})

	t.Run("InvalidResponsePayload", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), subscriber.SubscriberEndpoint+"/"+testUuid)

			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",
				Body:       io.NopCloser(bytes.NewBufferString(`OK`)),
				Header:     make(http.Header),
			}
		}

		testSubscriberClient := subscriber.NewSubscriberClient("")

		testSubscriberClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		s, err := testSubscriberClient.Get(testUuid)

		assert.Error(tt, err)
		assert.Nil(tt, s)
	})

	t.Run("RequestFailure", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), subscriber.SubscriberEndpoint+"/"+testUuid)

			return nil
		}

		testSubscriberClient := subscriber.NewSubscriberClient("")

		testSubscriberClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		s, err := testSubscriberClient.Get(testUuid)

		assert.Error(tt, err)
		assert.Nil(tt, s)
	})
}

func TestSubscriberClient_Add(t *testing.T) {
	t.Run("NetworkAdded", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			// Test request parameters
			assert.Equal(tt, req.URL.String(), subscriber.SubscriberEndpoint)

			// fake subscriber info
			subscriber := `{"subscriber":{"subscriber_id": "03cb753f-5e03-4c97-8e47-625115476c72", "last_name": "Foo"}}`

			// Send mock response
			return &http.Response{
				StatusCode: 201,
				Status:     "201 CREATED",

				// Send response to be tested
				Body: io.NopCloser(bytes.NewBufferString(subscriber)),

				// Must be set to non-nil value or it panics
				Header: make(http.Header),
			}
		}

		testSubscriberClient := subscriber.NewSubscriberClient("")

		// We replace the transport mechanism by mocking the http request
		// so that the test stays a unit test e.g no server/network call.
		testSubscriberClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		s, err := testSubscriberClient.Add(
			subscriber.AddSubscriberRequest{
				OrgId: testUuid,
				Name:  "Foo"},
		)

		assert.NoError(tt, err)
		assert.Equal(tt, testUuid, s.SubscriberId.String())
	})

	t.Run("InvalidResponseHeader", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), subscriber.SubscriberEndpoint)

			// error payload
			resp := `{"error":"internal server error"}`

			return &http.Response{
				StatusCode: 500,
				Status:     "500 INTERNAL SERVER ERROR",
				Body:       io.NopCloser(bytes.NewBufferString(resp)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}
		}

		testSubscriberClient := subscriber.NewSubscriberClient("")

		testSubscriberClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		s, err := testSubscriberClient.Add(
			subscriber.AddSubscriberRequest{
				OrgId: testUuid,
				Name:  "Foo"},
		)

		assert.Error(tt, err)
		assert.Nil(tt, s)
	})

	t.Run("InvalidResponsePayload", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), subscriber.SubscriberEndpoint)

			return &http.Response{
				StatusCode: 201,
				Status:     "201 CREATED",
				Body:       io.NopCloser(bytes.NewBufferString(`CREATED`)),
				Header:     make(http.Header),
			}
		}

		testSubscriberClient := subscriber.NewSubscriberClient("")

		testSubscriberClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		s, err := testSubscriberClient.Add(
			subscriber.AddSubscriberRequest{
				OrgId: testUuid,
				Name:  "Foo"},
		)

		assert.Error(tt, err)
		assert.Nil(tt, s)
	})

	t.Run("RequestFailure", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), subscriber.SubscriberEndpoint)

			return nil
		}

		testSubscriberClient := subscriber.NewSubscriberClient("")

		testSubscriberClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		s, err := testSubscriberClient.Add(
			subscriber.AddSubscriberRequest{
				OrgId: testUuid,
				Name:  "Foo"},
		)

		assert.Error(tt, err)
		assert.Nil(tt, s)
	})
}

func TestSubscriberClient_GetByEmail(t *testing.T) {
	t.Run("SubscriberFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			// Test request parameters
			assert.Equal(tt, req.URL.String(), subscriber.SubscriberEndpoint+"/email/"+testEmail)

			// fake subscriber info
			subscriber := `{"subscriber":{"subscriber_id": "03cb753f-5e03-4c97-8e47-625115476c72", "email": "foo@example.com"}}`

			// Send mock response
			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",

				// Send response to be tested
				Body: io.NopCloser(bytes.NewBufferString(subscriber)),

				// Must be set to non-nil value or it panics
				Header: make(http.Header),
			}
		}

		testSubscriberClient := subscriber.NewSubscriberClient("")

		// We replace the transport mechanism by mocking the http request
		// so that the test stays a unit test e.g no server/network call.
		testSubscriberClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		s, err := testSubscriberClient.GetByEmail(testEmail)

		assert.NoError(tt, err)
		assert.Equal(tt, testUuid, s.SubscriberId.String())
	})

	t.Run("SubscriberNotFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), subscriber.SubscriberEndpoint+"/email/"+testEmail)

			// error payload
			resp := `{"error":"not found"}`

			return &http.Response{
				StatusCode: 404,
				Status:     "404 NOT FOUND",
				Body:       io.NopCloser(bytes.NewBufferString(resp)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}
		}

		testSubscriberClient := subscriber.NewSubscriberClient("")

		testSubscriberClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		s, err := testSubscriberClient.GetByEmail(testEmail)

		assert.Error(tt, err)
		assert.Nil(tt, s)
	})

	t.Run("InvalidResponsePayload", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), subscriber.SubscriberEndpoint+"/email/"+testEmail)

			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",
				Body:       io.NopCloser(bytes.NewBufferString(`OK`)),
				Header:     make(http.Header),
			}
		}

		testSubscriberClient := subscriber.NewSubscriberClient("")

		testSubscriberClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		s, err := testSubscriberClient.GetByEmail(testEmail)

		assert.Error(tt, err)
		assert.Nil(tt, s)
	})

	t.Run("RequestFailure", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), subscriber.SubscriberEndpoint+"/"+testUuid)

			return nil
		}

		testSubscriberClient := subscriber.NewSubscriberClient("")

		testSubscriberClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		s, err := testSubscriberClient.Get(testUuid)

		assert.Error(tt, err)
		assert.Nil(tt, s)
	})
}
