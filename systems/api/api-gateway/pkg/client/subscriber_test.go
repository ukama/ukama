package client_test

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/tj/assert"
	"github.com/ukama/ukama/systems/api/api-gateway/pkg/client"
)

func TestSubscriberClient_Get(t *testing.T) {
	t.Run("SubscriberFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			// Test request parameters
			assert.Equal(tt, req.URL.String(), client.SubscriberEndpoint+"/"+testUuid)

			// fake network info
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

		testSubscriberClient := client.NewSubscriberClient("")

		// We replace the transport mechanism by mocking the http request
		// so that the test stays a unit test e.g no server/network call.
		testSubscriberClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		s, err := testSubscriberClient.Get(testUuid)

		assert.NoError(tt, err)
		assert.Equal(tt, testUuid, s.SubscriberId.String())
	})

	t.Run("SubscriberNotFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), client.SubscriberEndpoint+"/"+testUuid)

			return &http.Response{
				StatusCode: 404,
				Status:     "404 NOT FOUND",
				Header:     make(http.Header),
			}
		}

		testSubscriberClient := client.NewSubscriberClient("")

		testSubscriberClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		s, err := testSubscriberClient.Get(testUuid)

		assert.Error(tt, err)
		assert.Nil(tt, s)
	})

	t.Run("InvalidResponsePayload", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), client.SubscriberEndpoint+"/"+testUuid)

			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",
				Body:       io.NopCloser(bytes.NewBufferString(`OK`)),
				Header:     make(http.Header),
			}
		}

		testSubscriberClient := client.NewSubscriberClient("")

		testSubscriberClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		s, err := testSubscriberClient.Get(testUuid)

		assert.Error(tt, err)
		assert.Nil(tt, s)
	})

	t.Run("RequestFailure", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), client.SubscriberEndpoint+"/"+testUuid)

			return nil
		}

		testSubscriberClient := client.NewSubscriberClient("")

		testSubscriberClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		s, err := testSubscriberClient.Get(testUuid)

		assert.Error(tt, err)
		assert.Nil(tt, s)
	})
}

func TestSubscriberClient_Add(t *testing.T) {
	t.Run("NetworkAdded", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			// Test request parameters
			assert.Equal(tt, req.URL.String(), client.SubscriberEndpoint)

			// fake network info
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

		testSubscriberClient := client.NewSubscriberClient("")

		// We replace the transport mechanism by mocking the http request
		// so that the test stays a unit test e.g no server/network call.
		testSubscriberClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		s, err := testSubscriberClient.Add(
			client.AddSubscriberRequest{
				OrgId:    testUuid,
				LastName: "Foo"},
		)

		assert.NoError(tt, err)
		assert.Equal(tt, testUuid, s.SubscriberId.String())
	})

	t.Run("InvalidResponseHeader", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), client.SubscriberEndpoint)

			return &http.Response{
				StatusCode: 500,
				Status:     "500 INTERNAL SERVER ERROR",
				Body:       io.NopCloser(bytes.NewBufferString(`INTERNAL SERVER ERROR`)),
				Header:     make(http.Header),
			}
		}

		testSubscriberClient := client.NewSubscriberClient("")

		testSubscriberClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		s, err := testSubscriberClient.Add(
			client.AddSubscriberRequest{
				OrgId:    testUuid,
				LastName: "Foo"},
		)

		assert.Error(tt, err)
		assert.Nil(tt, s)
	})

	t.Run("InvalidResponsePayload", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), client.SubscriberEndpoint)

			return &http.Response{
				StatusCode: 201,
				Status:     "201 CREATED",
				Body:       io.NopCloser(bytes.NewBufferString(`CREATED`)),
				Header:     make(http.Header),
			}
		}

		testSubscriberClient := client.NewSubscriberClient("")

		testSubscriberClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		n, err := testSubscriberClient.Add(
			client.AddSubscriberRequest{
				OrgId:    testUuid,
				LastName: "Foo"},
		)

		assert.Error(tt, err)
		assert.Nil(tt, n)
	})

	t.Run("RequestFailure", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), client.SubscriberEndpoint)

			return nil
		}

		testSubscriberClient := client.NewSubscriberClient("")

		testSubscriberClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		s, err := testSubscriberClient.Add(
			client.AddSubscriberRequest{
				OrgId:    testUuid,
				LastName: "Foo"},
		)

		assert.Error(tt, err)
		assert.Nil(tt, s)
	})
}
