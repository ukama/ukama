package client_test

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/ukama/ukama/systems/billing/invoice/pkg/client"

	"github.com/tj/assert"
)

type RoundTripFunc func(req *http.Request) *http.Response

func (r RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return r(req), nil
}

func TestSubscriberClient_Get(t *testing.T) {
	testUuid := "03cb753f-5e03-4c97-8e47-625115476c72"

	t.Run("SubscriberFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			// Test request parameters
			assert.Equal(tt, req.URL.String(), client.SubscriberEndpoint+"/"+testUuid)

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

		testSubscriberClient := client.NewSubscriberClient("", false)

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

		testSubscriberClient := client.NewSubscriberClient("", false)

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

		testSubscriberClient := client.NewSubscriberClient("", false)

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

		testSubscriberClient := client.NewSubscriberClient("", false)

		testSubscriberClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		s, err := testSubscriberClient.Get(testUuid)

		assert.Error(tt, err)
		assert.Nil(tt, s)
	})
}
