package providers

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/tj/assert"
)

const apiEndpoint = "/v1/sims/"
const testIccid = "890000000000000001234"

type RoundTripFunc func(req *http.Request) *http.Response

func (r RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return r(req), nil
}

func TestOperatorClient_GetSimInfo(t *testing.T) {
	t.Run("SimFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			// Test request parameters
			assert.Equal(tt, req.URL.String(), apiEndpoint+testIccid)

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

		testOperatorClient, err := NewOperatorClient("", false)

		assert.NoError(tt, err)

		// We replace the transport mechanism by mocking the http request
		// so that the test stays a unit test e.g no server/network call.
		testOperatorClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		s, err := testOperatorClient.GetSimInfo(testIccid)

		assert.NoError(tt, err)
		assert.Equal(tt, testIccid, s.Iccid)
	})

	t.Run("SimNotFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), apiEndpoint+testIccid)

			return &http.Response{
				StatusCode: 404,
				Header:     make(http.Header),
			}
		}

		testOperatorClient, err := NewOperatorClient("", false)

		assert.NoError(tt, err)

		testOperatorClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		s, err := testOperatorClient.GetSimInfo(testIccid)

		assert.Error(tt, err)
		assert.Nil(tt, s)
	})

	t.Run("InvalidResponsePayload", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), apiEndpoint+testIccid)

			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewBufferString(`OK`)),
				Header:     make(http.Header),
			}
		}

		testOperatorClient, err := NewOperatorClient("", false)

		assert.NoError(tt, err)

		testOperatorClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		s, err := testOperatorClient.GetSimInfo(testIccid)

		assert.Error(tt, err)
		assert.Nil(tt, s)
	})

	t.Run("RequestFailure", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), apiEndpoint+testIccid)
			return nil
		}

		testOperatorClient, err := NewOperatorClient("", false)

		assert.NoError(tt, err)

		testOperatorClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		s, err := testOperatorClient.GetSimInfo(testIccid)

		assert.Error(tt, err)
		assert.Nil(tt, s)
	})
}

func TestOperatorClient_ActivateSim(t *testing.T) {
	t.Run("SimFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), apiEndpoint+testIccid)

			return &http.Response{
				StatusCode: 200,
				Header:     make(http.Header),
			}
		}

		testOperatorClient, err := NewOperatorClient("", false)

		assert.NoError(tt, err)

		testOperatorClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		err = testOperatorClient.ActivateSim(testIccid)

		assert.NoError(tt, err)
	})

	t.Run("SimNotFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), apiEndpoint+testIccid)

			return &http.Response{
				StatusCode: 404,
				Header:     make(http.Header),
			}
		}

		testOperatorClient, err := NewOperatorClient("", false)

		assert.NoError(tt, err)

		testOperatorClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		err = testOperatorClient.ActivateSim(testIccid)

		assert.Error(tt, err)
	})

	t.Run("RequestFailure", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), apiEndpoint+testIccid)
			return nil
		}

		testOperatorClient, err := NewOperatorClient("", false)

		assert.NoError(tt, err)

		testOperatorClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		err = testOperatorClient.ActivateSim(testIccid)

		assert.Error(tt, err)
	})
}

func TestOperatorClient_DeactivateSim(t *testing.T) {
	t.Run("SimFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), apiEndpoint+testIccid)

			return &http.Response{
				StatusCode: 200,
				Header:     make(http.Header),
			}
		}

		testOperatorClient, err := NewOperatorClient("", false)

		assert.NoError(tt, err)

		testOperatorClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		err = testOperatorClient.DeactivateSim(testIccid)

		assert.NoError(tt, err)
	})

	t.Run("SimNotFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), apiEndpoint+testIccid)

			return &http.Response{
				StatusCode: 404,
				Header:     make(http.Header),
			}
		}

		testOperatorClient, err := NewOperatorClient("", false)

		assert.NoError(tt, err)

		testOperatorClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		err = testOperatorClient.DeactivateSim(testIccid)

		assert.Error(tt, err)
	})

	t.Run("RequestFailure", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), apiEndpoint+testIccid)
			return nil
		}

		testOperatorClient, err := NewOperatorClient("", false)

		assert.NoError(tt, err)

		testOperatorClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		err = testOperatorClient.DeactivateSim(testIccid)

		assert.Error(tt, err)
	})
}

func TestOperatorClient_TerminateSim(t *testing.T) {
	t.Run("SimFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), apiEndpoint+testIccid)

			return &http.Response{
				StatusCode: 200,
				Header:     make(http.Header),
			}
		}

		testOperatorClient, err := NewOperatorClient("", false)

		assert.NoError(tt, err)

		testOperatorClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		err = testOperatorClient.TerminateSim(testIccid)

		assert.NoError(tt, err)
	})

	t.Run("SimNotFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), apiEndpoint+testIccid)

			return &http.Response{
				StatusCode: 404,
				Header:     make(http.Header),
			}
		}

		testOperatorClient, err := NewOperatorClient("", false)

		assert.NoError(tt, err)

		testOperatorClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		err = testOperatorClient.TerminateSim(testIccid)

		assert.Error(tt, err)
	})

	t.Run("RequestFailure", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), apiEndpoint+testIccid)
			return nil
		}

		testOperatorClient, err := NewOperatorClient("", false)

		assert.NoError(tt, err)

		testOperatorClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		err = testOperatorClient.TerminateSim(testIccid)

		assert.Error(tt, err)
	})
}
