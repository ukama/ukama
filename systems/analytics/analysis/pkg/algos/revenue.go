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

// Revenue KPIs from settled payments (payments.processor.list). A payment
// record carries no network field, only the paying SIM in its metadata, so we
// attribute revenue to a network by joining the SIM to subscriber.sim.list
// (sim_id -> network_id). Payments whose SIM can't be resolved to a network
// land in the org bucket (empty scope), so no revenue is lost.
//
// REVENUE is one KPI with exact components, so the aggregator ops cover
// three dashboard cards at every span:
//   op=SUM   -> revenue collected (cents)
//   op=COUNT -> number of purchases
//   op=AVG   -> average purchase value (weighted, exact)

// Revenue (REVENUE @ scope network_id): payments that became settled in this
// window, grouped by the network their SIM belongs to. Value = collected
// cents; components carry count and per-payment min/max. Known networks are
// zero-filled so a network with no sales reads $0 rather than "-".
//
// Settlement is a state transition (settled now, not settled one window ago)
// rather than paid_at falling in the window, so each payment counts once.
func Revenue(win schema.Window, in Datasets, spec schema.KpiSpec) ([]Result, error) {
	payments, ok := in["payments"]
	if !ok {
		return nil, fmt.Errorf("REVENUE: missing input 'payments'")
	}

	previous, ok := in["payments_prev"]
	if !ok {
		return nil, fmt.Errorf("REVENUE: missing input 'payments_prev' (state_prev baseline)")
	}

	sims, ok := in["sims"]
	if !ok {
		return nil, fmt.Errorf("REVENUE: missing input 'sims'")
	}

	networks, ok := in["networks"]
	if !ok {
		return nil, fmt.Errorf("REVENUE: missing input 'networks'")
	}

	simNet := simNetwork(sims)

	// network_id ("" = SIM not resolvable to a network) -> revenue components.
	type revAgg struct{ sum, count, min, max float64 }

	byNet := map[string]*revAgg{}

	ensure := func(net string) *revAgg {
		a, ok := byNet[net]
		if !ok {
			a = &revAgg{min: math.Inf(1), max: math.Inf(-1)}
			byNet[net] = a
		}

		return a
	}

	// Always keep the org bucket ("") so a recompute overwrites any prior
	// org-wide row and the org series stays continuous even at $0.
	ensure("")

	// Zero-fill every known network so one with no sales reads $0, not "-".
	for _, n := range networks {
		if nid := str(n["network_id"]); nid != "" {
			ensure(nid)
		}
	}

	for _, p := range newlySettled(payments, previous) {
		cents := math.Round(num(p["amount"]) * 100)
		net := simNet[paymentSim(p)] // "" when the SIM can't be mapped

		a := ensure(net)
		a.sum += cents
		a.count++
		a.min = math.Min(a.min, cents)
		a.max = math.Max(a.max, cents)
	}

	results := make([]Result, 0, len(byNet))

	for net, a := range byNet {
		var scope map[string]string
		if net != "" {
			scope = map[string]string{"network_id": net}
		}

		if a.count == 0 {
			results = append(results, Result{Scope: scope}) // zero row keeps the series continuous
			continue
		}

		results = append(results, Result{
			Scope: scope,
			Value: a.sum,
			Sum:   a.sum,
			Count: a.count,
			Min:   a.min,
			Max:   a.max,
		})
	}

	return results, nil
}

// newlySettled returns payments that read settled now and did not one window
// ago.
func newlySettled(current, previous []map[string]interface{}) []map[string]interface{} {
	was := make(map[string]bool, len(previous))

	for _, p := range previous {
		if !isSettled(str(p["status"])) {
			continue
		}

		if k := entityKey(p, "payment_id"); k != "" {
			was[k] = true
		}
	}

	out := make([]map[string]interface{}, 0)

	for _, p := range current {
		if !isSettled(str(p["status"])) {
			continue
		}

		k := entityKey(p, "payment_id")
		if k == "" || was[k] {
			continue
		}

		out = append(out, p)
	}

	return out
}

// PaidCustomers (PAID_CUSTOMERS @ scope network_id): distinct SIMs with at
// least one settled payment month-to-date (month of the window, UTC — matches
// the aggregator's monthly span with the default UTC rollup timezone), grouped
// by the network the SIM belongs to. A customer can buy several packages
// (several payments) against the same SIM, so we dedupe by SIM id (from the
// payment's metadata) rather than counting payments or per-transaction payer
// contact fields. SIMs that can't be mapped to a network land in the org
// bucket. State gauge: read with op=LAST; the monthly trend gives "+N this
// month".
func PaidCustomers(win schema.Window, in Datasets, spec schema.KpiSpec) ([]Result, error) {
	payments, ok := in["payments"]
	if !ok {
		return nil, fmt.Errorf("PAID_CUSTOMERS: missing input 'payments'")
	}

	sims, ok := in["sims"]
	if !ok {
		return nil, fmt.Errorf("PAID_CUSTOMERS: missing input 'sims'")
	}

	networks, ok := in["networks"]
	if !ok {
		return nil, fmt.Errorf("PAID_CUSTOMERS: missing input 'networks'")
	}

	simNet := simNetwork(sims)
	monthStart := time.Date(win.Start.Year(), win.Start.Month(), 1, 0, 0, 0, 0, time.UTC)

	// network_id ("" = unresolved SIM) -> set of distinct paying SIMs.
	byNet := map[string]map[string]bool{}

	ensure := func(net string) map[string]bool {
		s, ok := byNet[net]
		if !ok {
			s = map[string]bool{}
			byNet[net] = s
		}

		return s
	}

	// Org bucket + zero-fill known networks (0 paid customers reads "0").
	ensure("")

	for _, n := range networks {
		if nid := str(n["network_id"]); nid != "" {
			ensure(nid)
		}
	}

	for _, p := range settledIn(payments, monthStart, win.End) {
		sim := paymentSim(p)
		if sim == "" {
			continue
		}

		ensure(simNet[sim])[sim] = true
	}

	results := make([]Result, 0, len(byNet))

	for net, set := range byNet {
		var scope map[string]string
		if net != "" {
			scope = map[string]string{"network_id": net}
		}

		results = append(results, CountResult(scope, float64(len(set))))
	}

	return results, nil
}

// simNetwork maps sim_id -> network_id from subscriber.sim.list, so a payment
// (which carries only a SIM id in its metadata) can be attributed to the
// network that SIM belongs to.
func simNetwork(sims []map[string]interface{}) map[string]string {
	out := make(map[string]string, len(sims))

	for _, s := range sims {
		if id := str(s["sim_id"]); id != "" {
			out[id] = str(s["network_id"])
		}
	}

	return out
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
