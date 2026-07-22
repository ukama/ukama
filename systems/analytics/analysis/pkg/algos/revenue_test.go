/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain it at http://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package algos

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/ukama/ukama/systems/analytics/schema"
)

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
