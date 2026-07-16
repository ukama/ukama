/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package payments_test

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/tj/assert"

	"github.com/ukama/ukama/systems/common/rest/client/payments"
	"github.com/ukama/ukama/systems/common/uuid"
)

const testPaymentUuid = "2f2262d1-3d95-472d-9418-844b92c75b05"

func newCashPackageRequest() payments.AddPaymentRequest {
	return payments.AddPaymentRequest{
		ItemId:        uuid.NewV4().String(),
		ItemType:      "package",
		Amount:        "1.00",
		Currency:      "usd",
		PaymentMethod: "cash",
		Metadata:      map[string]string{"sim": uuid.NewV4().String()},
	}
}

func TestPaymentClient_Add(t *testing.T) {
	t.Run("PaymentAdded", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), payments.PaymentEndpoint)
			assert.Equal(tt, http.MethodPost, req.Method)

			body := `{"payment":{"id":"2f2262d1-3d95-472d-9418-844b92c75b05",` +
				`"item_type":"package","payment_method":"cash","amount":"1.00",` +
				`"currency":"usd","status":"completed"}}`

			return &http.Response{
				StatusCode: 201,
				Status:     "201 CREATED",
				Body:       io.NopCloser(bytes.NewBufferString(body)),
				Header:     make(http.Header),
			}
		}

		testPaymentClient := payments.NewPaymentClient("")
		testPaymentClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		p, err := testPaymentClient.Add(newCashPackageRequest())

		assert.NoError(tt, err)
		assert.Equal(tt, testPaymentUuid, p.Id)
		assert.Equal(tt, "completed", p.Status)
		assert.Equal(tt, "cash", p.PaymentMethod)
	})

	t.Run("InvalidResponseHeader", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), payments.PaymentEndpoint)

			resp := `{"error":"internal server error"}`

			return &http.Response{
				StatusCode: 500,
				Status:     "500 INTERNAL SERVER ERROR",
				Body:       io.NopCloser(bytes.NewBufferString(resp)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}
		}

		testPaymentClient := payments.NewPaymentClient("")
		testPaymentClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		p, err := testPaymentClient.Add(newCashPackageRequest())

		assert.Error(tt, err)
		assert.Nil(tt, p)
	})

	t.Run("InvalidResponsePayload", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), payments.PaymentEndpoint)

			return &http.Response{
				StatusCode: 201,
				Status:     "201 CREATED",
				Body:       io.NopCloser(bytes.NewBufferString(`CREATED`)),
				Header:     make(http.Header),
			}
		}

		testPaymentClient := payments.NewPaymentClient("")
		testPaymentClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		p, err := testPaymentClient.Add(newCashPackageRequest())

		assert.Error(tt, err)
		assert.Nil(tt, p)
	})

	t.Run("RequestFailure", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), payments.PaymentEndpoint)

			return nil
		}

		testPaymentClient := payments.NewPaymentClient("")
		testPaymentClient.R.C.SetTransport(RoundTripFunc(mockTransport))

		p, err := testPaymentClient.Add(newCashPackageRequest())

		assert.Error(tt, err)
		assert.Nil(tt, p)
	})
}

type RoundTripFunc func(req *http.Request) *http.Response

func (r RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return r(req), nil
}
