package providers_test

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/tj/assert"
	"github.com/ukama/ukama/systems/subscriber/sim-manager/pkg/clients/providers"
)

const testUuid = "03cb753f-5e03-4c97-8e47-625115476c72"

func TestPackageClient_GetPackageInfo(t *testing.T) {
	t.Run("PackageFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			// Test request parameters
			assert.Equal(tt, req.URL.String(), providers.PackageEndpoint+testUuid)

			// fake package info
			pkg := `{"package":{"uuid": "03cb753f-5e03-4c97-8e47-625115476c72", "is_active": true}}`

			// Send mock response
			return &http.Response{
				StatusCode: 200,

				// Send response to be tested
				Body: io.NopCloser(bytes.NewBufferString(pkg)),

				// Must be set to non-nil value or it panics
				Header: make(http.Header),
			}
		}

		testPackageClient, err := providers.NewPackageClient("", false)

		assert.NoError(tt, err)

		// We replace the transport mechanism by mocking the http request
		// so that the test stays a unit test e.g no server/network call.
		testPackageClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		p, err := testPackageClient.GetPackageInfo(testUuid)

		assert.NoError(tt, err)
		assert.Equal(tt, testUuid, p.Id)
	})

	t.Run("PackageNotFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), providers.PackageEndpoint+testUuid)

			return &http.Response{
				StatusCode: 404,
				Header:     make(http.Header),
			}
		}

		testPackageClient, err := providers.NewPackageClient("", false)

		assert.NoError(tt, err)

		testPackageClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		p, err := testPackageClient.GetPackageInfo(testUuid)

		assert.Error(tt, err)
		assert.Nil(tt, p)
	})

	t.Run("InvalidResponsePayload", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), providers.PackageEndpoint+testUuid)

			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewBufferString(`OK`)),
				Header:     make(http.Header),
			}
		}

		testPackageClient, err := providers.NewPackageClient("", false)

		assert.NoError(tt, err)

		testPackageClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		p, err := testPackageClient.GetPackageInfo(testUuid)

		assert.Error(tt, err)
		assert.Nil(tt, p)
	})

	t.Run("RequestFailure", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), providers.PackageEndpoint+testUuid)
			return nil
		}

		testPackageClient, err := providers.NewPackageClient("", false)

		assert.NoError(tt, err)

		testPackageClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		p, err := testPackageClient.GetPackageInfo(testUuid)

		assert.Error(tt, err)
		assert.Nil(tt, p)
	})
}
