package client_test

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/tj/assert"

	"github.com/ukama/ukama/systems/api/api-gateway/pkg/client"
)

func TestSimClient_Get(t *testing.T) {
	t.Run("SimFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			// Test request parameters
			assert.Equal(tt, req.URL.String(), client.SimEndpoint+"/"+testUuid)

			// fake network info
			sim := `{"sim":{"id": "03cb753f-5e03-4c97-8e47-625115476c72", "is_physical": false}}`

			// Send mock response
			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",

				// Send response to be tested
				Body: io.NopCloser(bytes.NewBufferString(sim)),

				// Must be set to non-nil value or it panics
				Header: make(http.Header),
			}
		}

		testSimClient := client.NewSimClient("")

		// We replace the transport mechanism by mocking the http request
		// so that the test stays a unit test e.g no server/network call.
		testSimClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		s, err := testSimClient.Get(testUuid)

		assert.NoError(tt, err)
		assert.Equal(tt, testUuid, s.Id)
	})

	t.Run("SimNotFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), client.SimEndpoint+"/"+testUuid)

			return &http.Response{
				StatusCode: 404,
				Status:     "404 NOT FOUND",
				Header:     make(http.Header),
			}
		}

		testSimClient := client.NewSimClient("")

		testSimClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		s, err := testSimClient.Get(testUuid)

		assert.Error(tt, err)
		assert.Nil(tt, s)
	})

	t.Run("InvalidResponsePayload", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), client.SimEndpoint+"/"+testUuid)

			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",
				Body:       io.NopCloser(bytes.NewBufferString(`OK`)),
				Header:     make(http.Header),
			}
		}

		testSimClient := client.NewSimClient("")

		testSimClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		s, err := testSimClient.Get(testUuid)

		assert.Error(tt, err)
		assert.Nil(tt, s)
	})

	t.Run("RequestFailure", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), client.SimEndpoint+"/"+testUuid)

			return nil
		}

		testSimClient := client.NewSimClient("")

		testSimClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		s, err := testSimClient.Get(testUuid)

		assert.Error(tt, err)
		assert.Nil(tt, s)
	})
}

func TestSimClient_Add(t *testing.T) {
	t.Run("SimAdded", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			// Test request parameters
			assert.Equal(tt, req.URL.String(), client.SimEndpoint)

			// fake network info
			sim := `{"sim":{"id": "03cb753f-5e03-4c97-8e47-625115476c72", "is_physical": false}}`

			// Send mock response
			return &http.Response{
				StatusCode: 201,
				Status:     "201 CREATED",

				// Send response to be tested
				Body: io.NopCloser(bytes.NewBufferString(sim)),

				// Must be set to non-nil value or it panics
				Header: make(http.Header),
			}
		}

		testSimClient := client.NewSimClient("")

		// We replace the transport mechanism by mocking the http request
		// so that the test stays a unit test e.g no server/network call.
		testSimClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		s, err := testSimClient.Add(
			client.AddSimRequest{
				SubscriberId: "some-subscriber_Id",
				PackageId:    "some-package_id"},
		)

		assert.NoError(tt, err)
		assert.Equal(tt, testUuid, s.Id)
	})

	t.Run("InvalidResponseHeader", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), client.SimEndpoint)

			return &http.Response{
				StatusCode: 500,
				Status:     "500 INTERNAL SERVER ERROR",
				Body:       io.NopCloser(bytes.NewBufferString(`INTERNAL SERVER ERROR`)),
				Header:     make(http.Header),
			}
		}

		testSimClient := client.NewSimClient("")

		testSimClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		s, err := testSimClient.Add(
			client.AddSimRequest{
				SubscriberId: "some-subscriber_Id",
				PackageId:    "some-package_id"},
		)

		assert.Error(tt, err)
		assert.Nil(tt, s)
	})

	t.Run("InvalidResponsePayload", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), client.SimEndpoint)

			return &http.Response{
				StatusCode: 201,
				Status:     "201 CREATED",
				Body:       io.NopCloser(bytes.NewBufferString(`CREATED`)),
				Header:     make(http.Header),
			}
		}

		testSimClient := client.NewSimClient("")

		testSimClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		s, err := testSimClient.Add(
			client.AddSimRequest{
				SubscriberId: "some-subscriber_Id",
				PackageId:    "some-package_id"},
		)

		assert.Error(tt, err)
		assert.Nil(tt, s)
	})

	t.Run("RequestFailure", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), client.SimEndpoint)

			return nil
		}

		testSimClient := client.NewSimClient("")

		testSimClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		s, err := testSimClient.Add(
			client.AddSimRequest{
				SubscriberId: "some-subscriber_Id",
				PackageId:    "some-package_id"},
		)

		assert.Error(tt, err)
		assert.Nil(tt, s)
	})
}
