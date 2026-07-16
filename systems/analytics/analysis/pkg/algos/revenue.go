/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain it at http://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package algos

import (
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

// Revenue (REVENUE @ org scope): successful payments with paid_at inside
// the window. Value = collected cents; components carry count and per-
// payment min/max.
func Revenue(win schema.Window, in Datasets, spec schema.KpiSpec) ([]Result, error) {
	payments, ok := in["payments"]
	if !ok {
		return nil, fmt.Errorf("REVENUE: missing input 'payments'")
	}

	sum, count := 0.0, 0.0
	min, max := math.Inf(1), math.Inf(-1)

	for _, p := range successfulIn(payments, win.Start, win.End) {
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

// PaidCustomers (PAID_CUSTOMERS @ org scope): distinct payers with at least
// one successful payment month-to-date (month of the window, UTC — matches
// the aggregator's monthly span with the default UTC rollup timezone).
// State gauge: read with op=LAST; the monthly trend gives "+N this month".
func PaidCustomers(win schema.Window, in Datasets, spec schema.KpiSpec) ([]Result, error) {
	payments, ok := in["payments"]
	if !ok {
		return nil, fmt.Errorf("PAID_CUSTOMERS: missing input 'payments'")
	}

	monthStart := time.Date(win.Start.Year(), win.Start.Month(), 1, 0, 0, 0, 0, time.UTC)

	payers := map[string]bool{}

	for _, p := range successfulIn(payments, monthStart, win.End) {
		if key := payerKey(p); key != "" {
			payers[key] = true
		}
	}

	return []Result{CountResult(nil, float64(len(payers)))}, nil
}

// successfulIn returns successful payments whose paid_at falls in [from, to).
func successfulIn(payments []map[string]interface{}, from, to time.Time) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(payments))

	for _, p := range payments {
		if !strings.EqualFold(str(p["status"]), "success") {
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

// payerKey identifies a payer: email, else phone, else empty (uncounted).
func payerKey(p map[string]interface{}) string {
	if email := str(p["payer_email"]); email != "" {
		return strings.ToLower(email)
	}

	return str(p["payer_phone"])
}
