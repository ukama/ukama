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

// rowFor returns the result whose scope network_id matches net ("" = org
// bucket / no scope). nil map reads yield "", so an unscoped Result matches "".
func rowFor(results []Result, net string) *Result {
	for i := range results {
		if results[i].Scope["network_id"] == net {
			return &results[i]
		}
	}

	return nil
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

// A payment that became settled in this window counts against its SIM's
// network; one already settled, and one not settled yet, do not.
func TestRevenue_CountsNewlySettledPayment(t *testing.T) {
	win := schema.Window{
		Start: time.Date(2026, 7, 22, 15, 0, 0, 0, time.UTC),
		End:   time.Date(2026, 7, 22, 16, 0, 0, 0, time.UTC),
	}
	in := Datasets{
		"payments": {
			{
				"payment_id": "pay-new",
				"amount":     "1.00",
				"status":     "completed",
				"paid_at":    "2026-07-22 15:25:42.989109 +0000 UTC",
				"metadata":   metaFor("sim-A"),
			},
			// excluded: already read settled one window ago
			{"payment_id": "pay-old", "amount": "5.00", "status": "completed", "metadata": metaFor("sim-A")},
			// excluded: not settled yet
			{"payment_id": "pay-pending", "amount": "9.00", "status": "pending", "metadata": metaFor("sim-A")},
		},
		"payments_prev": {
			{"payment_id": "pay-old", "status": "completed"},
			{"payment_id": "pay-pending", "status": "pending"},
		},
		"sims":     {{"sim_id": "sim-A", "network_id": "net-a"}},
		"networks": {{"network_id": "net-a"}},
	}

	results, err := Revenue(win, in, schema.KpiSpec{})
	assert.NoError(t, err)

	row := rowFor(results, "net-a")
	assert.NotNil(t, row, "net-a row present")
	assert.Equal(t, float64(100), row.Sum, "one $1.00 payment on net-a = 100 cents")
	assert.Equal(t, float64(1), row.Count)
}

// A late settlement counts in the window the transition is observed, once.
func TestRevenue_CountsLateSettlementExactlyOnce(t *testing.T) {
	win := schema.Window{
		Start: time.Date(2026, 7, 22, 15, 0, 0, 0, time.UTC),
		End:   time.Date(2026, 7, 22, 16, 0, 0, 0, time.UTC),
	}
	settled := map[string]interface{}{
		"payment_id": "pay-1", "amount": "3.00", "status": "completed",
		// hours before the window that finally observes the transition
		"paid_at": "2026-07-22T04:00:00Z", "metadata": metaFor("sim-A"),
	}
	sims := []map[string]interface{}{{"sim_id": "sim-A", "network_id": "net-a"}}
	networks := []map[string]interface{}{{"network_id": "net-a"}}

	first, err := Revenue(win, Datasets{
		"payments":      {settled},
		"payments_prev": {{"payment_id": "pay-1", "status": "pending"}},
		"sims":          sims,
		"networks":      networks,
	}, schema.KpiSpec{})
	assert.NoError(t, err)
	assert.Equal(t, float64(300), rowFor(first, "net-a").Sum)

	second, err := Revenue(win, Datasets{
		"payments":      {settled},
		"payments_prev": {settled},
		"sims":          sims,
		"networks":      networks,
	}, schema.KpiSpec{})
	assert.NoError(t, err)
	assert.Equal(t, float64(0), rowFor(second, "net-a").Sum, "not counted twice")
}

// Missing the lag-1 baseline is a spec error, not a silent zero.
func TestRevenue_RequiresPrevBaseline(t *testing.T) {
	_, err := Revenue(schema.Window{}, Datasets{
		"payments": {},
		"sims":     {},
		"networks": {},
	}, schema.KpiSpec{})
	assert.Error(t, err)
}

// Each payment is attributed to the network of its paying SIM; a SIM that
// can't be resolved lands in the org bucket (empty scope), not on a network.
func TestRevenue_AttributesBySim(t *testing.T) {
	win := schema.Window{
		Start: time.Date(2026, 7, 22, 0, 0, 0, 0, time.UTC),
		End:   time.Date(2026, 7, 23, 0, 0, 0, 0, time.UTC),
	}
	in := Datasets{
		"payments": {
			{"payment_id": "p1", "amount": "1.00", "status": "completed", "metadata": metaFor("sim-A")},
			{"payment_id": "p2", "amount": "2.00", "status": "completed", "metadata": metaFor("sim-B")},
			// SIM not in subscriber.sim.list -> org bucket
			{"payment_id": "p3", "amount": "4.00", "status": "completed", "metadata": metaFor("sim-Z")},
		},
		"payments_prev": {},
		"sims": {
			{"sim_id": "sim-A", "network_id": "net-a"},
			{"sim_id": "sim-B", "network_id": "net-b"},
		},
		"networks": {{"network_id": "net-a"}, {"network_id": "net-b"}},
	}

	results, err := Revenue(win, in, schema.KpiSpec{})
	assert.NoError(t, err)

	assert.Equal(t, float64(100), rowFor(results, "net-a").Sum, "net-a = $1")
	assert.Equal(t, float64(200), rowFor(results, "net-b").Sum, "net-b = $2")
	assert.Equal(t, float64(400), rowFor(results, "").Sum, "unresolved SIM in org bucket = $4")
}

// The SIM id lives in the payment's base64-encoded metadata.
func TestPaymentSim(t *testing.T) {
	meta := "eyJzaW0iOiAiYWNkNzE5MzUtNWIzZi00MjE5LThhMmEtMTUxNDg0ZTI1OTc5IiwgInByb3Zpc2lvbmVkIjogInRydWUifQ=="
	assert.Equal(t, "acd71935-5b3f-4219-8a2a-151484e25979",
		paymentSim(map[string]interface{}{"metadata": meta}))
	assert.Equal(t, "", paymentSim(map[string]interface{}{"metadata": ""}))
	assert.Equal(t, "", paymentSim(map[string]interface{}{}))
}

// PAID_CUSTOMERS dedupes by SIM within a network: multiple package payments on
// the same SIM count once (a customer buying several packages is still one
// paid customer).
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
		"sims": {
			{"sim_id": "sim-A", "network_id": "net-a"},
			{"sim_id": "sim-B", "network_id": "net-a"},
			{"sim_id": "sim-C", "network_id": "net-a"},
			{"sim_id": "sim-D", "network_id": "net-a"},
		},
		"networks": {{"network_id": "net-a"}},
	}

	results, err := PaidCustomers(win, in, schema.KpiSpec{})
	assert.NoError(t, err)

	row := rowFor(results, "net-a")
	assert.NotNil(t, row, "net-a row present")
	assert.Equal(t, float64(2), row.Value, "distinct SIMs A and B this month")
}
