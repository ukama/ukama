/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain it at http://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package algos

import (
	"encoding/base64"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/ukama/ukama/systems/analytics/schema"
)

func metaFor(sim string) string {
	return base64.StdEncoding.EncodeToString([]byte(`{"sim":"` + sim + `","provisioned":"true"}`))
}

// The payments service marks paid transactions "completed" (not "success").
func TestIsSettled(t *testing.T) {
	assert.True(t, isSettled("completed"))
	assert.True(t, isSettled("COMPLETED"))
	assert.True(t, isSettled("success"))
	assert.False(t, isSettled("pending"))
	assert.False(t, isSettled("failed"))
	assert.False(t, isSettled(""))
}

// paid_at arrives as Go's default time.Time string, not RFC3339.
func TestParseTime_PaymentsFormat(t *testing.T) {
	got, err := parseTime("2026-07-22 15:25:42.989109 +0000 UTC")
	assert.NoError(t, err)
	assert.Equal(t, "2026-07-22T15:25:42Z", got.UTC().Format(time.RFC3339))

	// RFC3339 (package start_date) still parses — no regression.
	rfc, err := parseTime("2026-07-22T15:26:42Z")
	assert.NoError(t, err)
	assert.Equal(t, 2026, rfc.UTC().Year())

	_, err = parseTime("not-a-time")
	assert.Error(t, err)
}

// End-to-end: a "completed" payment with a Go-format paid_at inside the window
// is now counted (SUM = cents, COUNT = 1).
func TestRevenue_CountsCompletedPayment(t *testing.T) {
	win := schema.Window{
		Start: time.Date(2026, 7, 22, 15, 0, 0, 0, time.UTC),
		End:   time.Date(2026, 7, 22, 16, 0, 0, 0, time.UTC),
	}
	in := Datasets{
		"payments": {
			{
				"amount":  "1.00",
				"status":  "completed",
				"paid_at": "2026-07-22 15:25:42.989109 +0000 UTC",
			},
			// excluded: right status, but outside the window
			{"amount": "5.00", "status": "completed", "paid_at": "2026-07-21T09:00:00Z"},
			// excluded: in window, but not settled
			{"amount": "9.00", "status": "pending", "paid_at": "2026-07-22T15:30:00Z"},
		},
	}

	results, err := Revenue(win, in, schema.KpiSpec{})
	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, float64(100), results[0].Sum, "one $1.00 payment = 100 cents")
	assert.Equal(t, float64(1), results[0].Count)
}

// The SIM id lives in the payment's base64-encoded metadata.
func TestPaymentSim(t *testing.T) {
	meta := "eyJzaW0iOiAiYWNkNzE5MzUtNWIzZi00MjE5LThhMmEtMTUxNDg0ZTI1OTc5IiwgInByb3Zpc2lvbmVkIjogInRydWUifQ=="
	assert.Equal(t, "acd71935-5b3f-4219-8a2a-151484e25979",
		paymentSim(map[string]interface{}{"metadata": meta}))
	assert.Equal(t, "", paymentSim(map[string]interface{}{"metadata": ""}))
	assert.Equal(t, "", paymentSim(map[string]interface{}{}))
}

// PAID_CUSTOMERS dedupes by SIM: multiple package payments on the same SIM
// count once (a customer buying several packages is still one paid customer).
func TestPaidCustomers_DistinctSims(t *testing.T) {
	win := schema.Window{
		Start: time.Date(2026, 7, 22, 15, 0, 0, 0, time.UTC),
		End:   time.Date(2026, 7, 22, 16, 0, 0, 0, time.UTC),
	}
	in := Datasets{
		"payments": {
			// SIM A pays twice (two packages) -> counts once
			{"status": "completed", "paid_at": "2026-07-10T09:00:00Z", "metadata": metaFor("sim-A")},
			{"status": "completed", "paid_at": "2026-07-15T09:00:00Z", "metadata": metaFor("sim-A")},
			// SIM B pays once
			{"status": "completed", "paid_at": "2026-07-20T09:00:00Z", "metadata": metaFor("sim-B")},
			// not settled -> ignored
			{"status": "pending", "paid_at": "2026-07-21T09:00:00Z", "metadata": metaFor("sim-C")},
			// outside the month -> ignored
			{"status": "completed", "paid_at": "2026-06-30T09:00:00Z", "metadata": metaFor("sim-D")},
		},
	}

	results, err := PaidCustomers(win, in, schema.KpiSpec{})
	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, float64(2), results[0].Value, "distinct SIMs A and B this month")
}
