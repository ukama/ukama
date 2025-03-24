/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package notification_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tj/assert"
	p "github.com/ukama/ukama/systems/common/rest/client/payments"
)

func TestListReq(t *testing.T) {
	expectedPayments := p.PaymentsRes{
		Payments: []*p.Payment{
			{
				Id:            "1",
				ItemId:        "item1",
				ItemType:      "type1",
				Amount:        "100.00",
				Currency:      "USD",
				PaymentMethod: "card",
				Status:        "completed",
			},
		},
	}
	responseJSON, err := json.Marshal(expectedPayments)
	if err != nil {
		t.Fatalf("Failed to marshal expected payments: %v", err)
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("item_id") != "item1" {
			http.Error(w, "bad item_id", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responseJSON)
	}))
	defer ts.Close()

	pc := p.NewPaymentClient(ts.URL)

	req := p.ListReq{
		ItemId:        "item1",
		ItemType:      "type1",
		PaymentMethod: "card",
		Status:        "completed",
		Count:         1,
		Sort:          true,
	}

	res, err := pc.ListReq(req)
	assert.NoError(t, err, "ListReq should not return an error")
	assert.NotNil(t, res, "Response should not be nil")
	assert.Equal(t, 1, len(res.Payments), "Expected one payment in the response")
	if len(res.Payments) > 0 {
		payment := res.Payments[0]
		assert.Equal(t, "1", payment.Id)
		assert.Equal(t, "item1", payment.ItemId)
		assert.Equal(t, "type1", payment.ItemType)
		assert.Equal(t, "100.00", payment.Amount)
		assert.Equal(t, "USD", payment.Currency)
		assert.Equal(t, "card", payment.PaymentMethod)
		assert.Equal(t, "completed", payment.Status)
	}
}
