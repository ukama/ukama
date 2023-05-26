package client_test

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/tj/assert"
	"github.com/ukama/ukama/systems/billing/api-gateway/internal/client"
)

// fake pdf file
const pdfContent = `This \\ is \\ a \\ fake \\ PDF \\ file.`

const invoiceId = "03cb753f-5e03-4c97-8e47-625115476c72"

type RoundTripFunc func(req *http.Request) *http.Response

func (r RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return r(req), nil
}

func TestPdfClient_GetPdf(t *testing.T) {
	t.Run("PdfFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			// Test request parameters
			assert.Equal(tt, req.URL.String(), client.FileServerEndpoint+invoiceId+".pdf")

			// Send mock response
			return &http.Response{
				StatusCode: 200,

				// Send response to be tested
				Body: io.NopCloser(bytes.NewBufferString(pdfContent)),

				// Must be set to non-nil value or it panics
				Header: make(http.Header),
			}
		}

		testPdfClient, err := client.NewPdfClient("", false)

		assert.NoError(tt, err)

		// We replace the transport mechanism by mocking the http request
		// so that the test stays a unit test e.g no server/network call.
		testPdfClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		p, err := testPdfClient.GetPdf(invoiceId)

		assert.NoError(tt, err)
		assert.Equal(tt, pdfContent, string(p))
	})

	t.Run("PdfNotFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), client.FileServerEndpoint+invoiceId+".pdf")

			return &http.Response{
				StatusCode: 404,
				Header:     make(http.Header),
			}
		}

		testPdfClient, err := client.NewPdfClient("", false)

		assert.NoError(tt, err)

		testPdfClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		p, err := testPdfClient.GetPdf(invoiceId)

		assert.Error(tt, err)
		assert.Nil(tt, p)
	})

	t.Run("RequestFailure", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), client.FileServerEndpoint+invoiceId+".pdf")
			return nil
		}

		testPdfClient, err := client.NewPdfClient("", false)

		assert.NoError(tt, err)

		testPdfClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		p, err := testPdfClient.GetPdf(invoiceId)

		assert.Error(tt, err)
		assert.Nil(tt, p)
	})
}
