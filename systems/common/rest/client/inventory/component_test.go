package inventory_test

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ukama/ukama/systems/common/rest/client"
	"github.com/ukama/ukama/systems/common/rest/client/inventory"
)

const testUuid = "03cb753f-5e03-4c97-8e47-625115476c72"

func TestComponentClient_Get(t *testing.T) {
	t.Run("ComponentFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			// Test request parameters
			assert.Equal(tt, req.URL.String(), inventory.ComponentEndpoint+"/"+testUuid)

			// fake component info
			comonent := `{"component":{"id": "03cb753f-5e03-4c97-8e47-625115476c72", "type": "backhaul"}}`

			// Send mock response
			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",

				// Send response to be tested
				Body: io.NopCloser(bytes.NewBufferString(comonent)),

				// Must be set to non-nil value or it panics
				Header: make(http.Header),
			}
		}

		testComponentClient := inventory.NewComponentClient("")

		// We replace the transport mechanism by mocking the http request
		// so that the test stays a unit test e.g no server/component call.
		testComponentClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		c, err := testComponentClient.Get(testUuid)

		assert.NoError(tt, err)
		assert.Equal(tt, testUuid, c.Id.String())
	})

	t.Run("ComponentNotFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), inventory.ComponentEndpoint+"/"+testUuid)

			// error payload
			resp := `{"error":"not found"}`

			return &http.Response{
				StatusCode: 404,
				Status:     "404 NOT FOUND",
				Body:       io.NopCloser(bytes.NewBufferString(resp)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}
		}

		testComponentClient := inventory.NewComponentClient("")

		testComponentClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		c, err := testComponentClient.Get(testUuid)

		assert.Error(tt, err)
		assert.Nil(tt, c)
	})

	t.Run("InvalidResponsePayload", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), inventory.ComponentEndpoint+"/"+testUuid)

			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",
				Body:       io.NopCloser(bytes.NewBufferString(`OK`)),
				Header:     make(http.Header),
			}
		}

		testComponentClient := inventory.NewComponentClient("")

		testComponentClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		c, err := testComponentClient.Get(testUuid)

		assert.Error(tt, err)
		assert.Nil(tt, c)
	})

	t.Run("RequestFailure", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), inventory.ComponentEndpoint+"/"+testUuid)

			return nil
		}

		testComponentClient := inventory.NewComponentClient("")

		testComponentClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		c, err := testComponentClient.Get(testUuid)

		assert.Error(tt, err)
		assert.Nil(tt, c)
	})
}
