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
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/ukama/ukama/systems/analytics/schema"
)

// Revenue KPIs from settled payments (payments.processor.list). Payments
// carry no network attribution, so these are org-wide (empty scope).
//
// REVENUE is one KPI with exact components, so the aggregator ops cover
// three dashboard cards at every span:
//   op=SUM   -> revenue collected (cents)
//   op=COUNT -> number of purchases
//   op=AVG   -> average purchase value (weighted, exact)

// Revenue (REVENUE @ org scope): settled payments with paid_at inside the
// window. Value = collected cents; components carry count and per-payment
// min/max.
func Revenue(win schema.Window, in Datasets, spec schema.KpiSpec) ([]Result, error) {
	payments, ok := in["payments"]
	if !ok {
		return nil, fmt.Errorf("REVENUE: missing input 'payments'")
	}

	sum, count := 0.0, 0.0
	min, max := math.Inf(1), math.Inf(-1)

	for _, p := range settledIn(payments, win.Start, win.End) {
		cents := math.Round(num(p["amount"]) * 100)

		sum += cents
		count++
		min = math.Min(min, cents)
		max = math.Max(max, cents)
	}

	if count == 0 {
		return []Result{{Scope: nil}}, nil // zero row keeps series/trends continuous
	}

	return []Result{{
		Scope: nil,
		Value: sum,
		Sum:   sum,
		Count: count,
		Min:   min,
		Max:   max,
	}}, nil
}

// PaidCustomers (PAID_CUSTOMERS @ org scope): distinct SIMs with at least one
// settled payment month-to-date (month of the window, UTC — matches the
// aggregator's monthly span with the default UTC rollup timezone). A customer
// can buy several packages (several payments) against the same SIM, so we
// dedupe by SIM id (from the payment's metadata) rather than counting payments
// or per-transaction payer contact fields. State gauge: read with op=LAST; the
// monthly trend gives "+N this month".
func PaidCustomers(win schema.Window, in Datasets, spec schema.KpiSpec) ([]Result, error) {
	payments, ok := in["payments"]
	if !ok {
		return nil, fmt.Errorf("PAID_CUSTOMERS: missing input 'payments'")
	}

	monthStart := time.Date(win.Start.Year(), win.Start.Month(), 1, 0, 0, 0, 0, time.UTC)

	sims := map[string]bool{}

	for _, p := range settledIn(payments, monthStart, win.End) {
		if sim := paymentSim(p); sim != "" {
			sims[sim] = true
		}
	}

	return []Result{CountResult(nil, float64(len(sims)))}, nil
}

// settledIn returns settled payments whose paid_at falls in [from, to).
func settledIn(payments []map[string]interface{}, from, to time.Time) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(payments))

	for _, p := range payments {
		if !isSettled(str(p["status"])) {
			continue
		}

		paidAt, err := parseTime(str(p["paid_at"]))
		if err != nil {
			continue
		}

		if !paidAt.Before(from) && paidAt.Before(to) {
			out = append(out, p)
		}
	}

	return out
}

// isSettled reports whether a payment status counts as collected revenue. The
// payments service marks a paid transaction "completed"; we accept a few
// synonyms (case-insensitive) so a status-vocabulary change upstream doesn't
// silently zero revenue. Not settled: pending, failed, etc.
func isSettled(status string) bool {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "completed", "success", "succeeded", "paid":
		return true
	default:
		return false
	}
}

// paymentSim extracts the SIM id a payment was made for. The payments service
// puts it in the record's `metadata`, which arrives base64-encoded JSON
// (e.g. {"sim": "acd719...", "provisioned": "true"}). Returns "" when absent
// or unparseable (uncounted). Tolerates padded/unpadded base64 and a metadata
// that's already a decoded object.
func paymentSim(p map[string]interface{}) string {
	switch m := p["metadata"].(type) {
	case string:
		if m == "" {
			return ""
		}

		decoded, err := base64.StdEncoding.DecodeString(m)
		if err != nil {
			decoded, err = base64.RawStdEncoding.DecodeString(m)
			if err != nil {
				return ""
			}
		}

		var obj struct {
			Sim string `json:"sim"`
		}
		if json.Unmarshal(decoded, &obj) == nil {
			return obj.Sim
		}

		return ""
	case map[string]interface{}:
		return str(m["sim"])
	default:
		return ""
	}
}
