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
	"github.com/ukama/ukama/systems/common/ukama"
)

// Package KPIs. Shared conventions:
//   - A package applies to a network when its network_id equals the
//     network's id OR is empty (org-level package).
//   - Money is integer org-currency cents (source amounts are float64).
//   - Revenue here is a SALES PROXY: assignments × package price. Settled
//     payment revenue comes later from billing/payments sources.

// pkgInfo is the decoded package catalog entry.
type pkgInfo struct {
	id          string
	networkID   string
	amountCents float64
	duration    float64 // days
	active      bool
	dataBytes   float64 // package data allowance, normalized to bytes
}

func decodePackages(in Datasets) (map[string]pkgInfo, error) {
	rows, ok := in["packages"]
	if !ok {
		return nil, fmt.Errorf("missing input 'packages'")
	}

	out := make(map[string]pkgInfo, len(rows))

	for _, p := range rows {
		id := str(p["package_id"])
		if id == "" {
			continue
		}

		out[id] = pkgInfo{
			id:          id,
			networkID:   str(p["network_id"]),
			amountCents: math.Round(num(p["amount"]) * 100),
			duration:    num(p["duration"]),
			active:      asBool(p["active"]),
			dataBytes:   dataVolumeBytes(p),
		}
	}

	return out, nil
}

// appliesTo: org-level packages (empty network_id) are valid for every
// network.
func (p pkgInfo) appliesTo(networkID string) bool {
	return p.networkID == "" || p.networkID == networkID
}

// mrrCents normalizes the package price to a 30-day recurring amount.
func (p pkgInfo) mrrCents() float64 {
	if p.duration <= 0 {
		return p.amountCents
	}

	return p.amountCents * 30 / p.duration
}

// PackageSales (PACKAGE_SALES @ scope network_id+package_id): sim-package
// assignments whose start_date falls inside the window. Daily SUM = sold
// that day. No zero-fill: absent scope rows read as 0 in SUM rollups.
func PackageSales(win schema.Window, in Datasets, spec schema.KpiSpec) ([]Result, error) {
	assignments, ok := in["sim_packages"]
	if !ok {
		return nil, fmt.Errorf("PACKAGE_SALES: missing input 'sim_packages'")
	}

	counts := map[[2]string]float64{} // [network, package] -> sold

	for _, a := range assignments {
		networkID, packageID := str(a["network_id"]), str(a["package_id"])
		if networkID == "" || packageID == "" {
			continue
		}

		// Count a sale when the package was PURCHASED (created_at), not when it
		// becomes active (start_date). A queued package with a future
		// start_date is still a sale now — this keeps PACKAGE_SALES /
		// PACKAGE_REVENUE consistent with REVENUE, which counts the payment at
		// paid time.
		created, err := parseTime(str(a["created_at"]))
		if err != nil {
			continue
		}

		if !created.Before(win.Start) && created.Before(win.End) {
			counts[[2]string{networkID, packageID}]++
		}
	}

	results := make([]Result, 0, len(counts))
	for key, n := range counts {
		results = append(results, CountResult(
			map[string]string{"network_id": key[0], "package_id": key[1]}, n))
	}

	return results, nil
}

// PackageRevenue (PACKAGE_REVENUE @ scope network_id+package_id): sales in
// the window × package price (org-currency cents). Sales proxy.
func PackageRevenue(win schema.Window, in Datasets, spec schema.KpiSpec) ([]Result, error) {
	packages, err := decodePackages(in)
	if err != nil {
		return nil, fmt.Errorf("PACKAGE_REVENUE: %w", err)
	}

	sales, err := PackageSales(win, in, spec)
	if err != nil {
		return nil, fmt.Errorf("PACKAGE_REVENUE: %w", err)
	}

	results := make([]Result, 0, len(sales))

	for _, s := range sales {
		p, ok := packages[s.Scope["package_id"]]
		if !ok || !p.appliesTo(s.Scope["network_id"]) {
			continue // assignment references an unknown/foreign package
		}

		results = append(results, CountResult(s.Scope, s.Value*p.amountCents))
	}

	return results, nil
}

// DataSold (DATA_SOLD @ scope network_id): bytes of package data allowance
// SOLD in the window — sum over sim-package assignments whose start_date
// falls in the window of the assigned package's data volume converted to
// bytes, grouped by network. This is a SALES proxy (allowance purchased),
// distinct from USAGE_BY_NETWORK which measures data actually consumed.
// Zero-filled across networks so the daily/monthly SUM series stays
// continuous.
func DataSold(win schema.Window, in Datasets, spec schema.KpiSpec) ([]Result, error) {
	packages, err := decodePackages(in)
	if err != nil {
		return nil, fmt.Errorf("DATA_SOLD: %w", err)
	}

	sales, err := PackageSales(win, in, spec)
	if err != nil {
		return nil, fmt.Errorf("DATA_SOLD: %w", err)
	}

	_, networks, err := packageInputs(in, "DATA_SOLD")
	if err != nil {
		return nil, err
	}

	bytesByNetwork := map[string]float64{}

	for _, s := range sales {
		p, ok := packages[s.Scope["package_id"]]
		if !ok || !p.appliesTo(s.Scope["network_id"]) {
			continue // assignment references an unknown/foreign package
		}

		bytesByNetwork[s.Scope["network_id"]] += s.Value * p.dataBytes
	}

	return zeroFilled(networks, func(networkID string) float64 {
		return bytesByNetwork[networkID]
	}), nil
}

// Mrr (MRR @ scope network_id): sum over currently-active assignments of the
// package price normalized to 30 days. State gauge (read with LAST).
func Mrr(win schema.Window, in Datasets, spec schema.KpiSpec) ([]Result, error) {
	packages, err := decodePackages(in)
	if err != nil {
		return nil, fmt.Errorf("MRR: %w", err)
	}

	assignments, networks, err := packageInputs(in, "MRR")
	if err != nil {
		return nil, err
	}

	mrr := map[string]float64{}

	for _, a := range activeOf(assignments) {
		networkID := str(a["network_id"])

		p, ok := packages[str(a["package_id"])]
		if !ok || !p.appliesTo(networkID) {
			continue
		}

		mrr[networkID] += p.mrrCents()
	}

	return zeroFilled(networks, func(networkID string) float64 {
		return math.Round(mrr[networkID])
	}), nil
}

// Arpu (ARPU @ scope network_id): committed spend ÷ distinct subscribers
// holding an active assignment, per network. Committed spend is the sum of the
// package's ACTUAL price (not the 30-day-normalized MRR figure) over every
// currently-active assignment, so a $1 plan bought by one subscriber reads $1
// regardless of the plan's duration — a per-user spend, not a monthly
// run-rate. The denominator matches CUSTOMERS_ON_PLAN (distinct subscribers on
// any active plan) so ARPU × customers reconciles to committed spend. State
// gauge; networks with no active subscriber read $0.
func Arpu(win schema.Window, in Datasets, spec schema.KpiSpec) ([]Result, error) {
	packages, err := decodePackages(in)
	if err != nil {
		return nil, fmt.Errorf("ARPU: %w", err)
	}

	assignments, networks, err := packageInputs(in, "ARPU")
	if err != nil {
		return nil, err
	}

	spend := map[string]float64{}               // network -> committed cents
	subscribers := map[string]map[string]bool{} // network -> distinct subscriber ids

	for _, a := range activeOf(assignments) {
		networkID, subscriberID := str(a["network_id"]), str(a["subscriber_id"])
		if networkID == "" || subscriberID == "" {
			continue
		}

		if subscribers[networkID] == nil {
			subscribers[networkID] = map[string]bool{}
		}

		subscribers[networkID][subscriberID] = true

		// Numerator counts only assignments we can price; an unresolvable
		// package contributes $0 but the subscriber still counts (data issue,
		// not a reason to shrink the denominator).
		if p, ok := packages[str(a["package_id"])]; ok && p.appliesTo(networkID) {
			spend[networkID] += p.amountCents
		}
	}

	return zeroFilled(networks, func(networkID string) float64 {
		n := float64(len(subscribers[networkID]))
		if n == 0 {
			return 0
		}

		return math.Round(spend[networkID] / n)
	}), nil
}

// CustomersOnPlan (CUSTOMERS_ON_PLAN @ scope network_id): distinct
// subscribers with at least one active assignment. State gauge.
func CustomersOnPlan(win schema.Window, in Datasets, spec schema.KpiSpec) ([]Result, error) {
	assignments, networks, err := packageInputs(in, "CUSTOMERS_ON_PLAN")
	if err != nil {
		return nil, err
	}

	subscribers := map[string]map[string]bool{}

	for _, a := range activeOf(assignments) {
		networkID, subscriberID := str(a["network_id"]), str(a["subscriber_id"])
		if networkID == "" || subscriberID == "" {
			continue
		}

		if subscribers[networkID] == nil {
			subscribers[networkID] = map[string]bool{}
		}

		subscribers[networkID][subscriberID] = true
	}

	return zeroFilled(networks, func(networkID string) float64 {
		return float64(len(subscribers[networkID]))
	}), nil
}

// ActivePlans (ACTIVE_PLANS @ scope network_id): distinct packages with at
// least one active assignment. State gauge.
func ActivePlans(win schema.Window, in Datasets, spec schema.KpiSpec) ([]Result, error) {
	assignments, networks, err := packageInputs(in, "ACTIVE_PLANS")
	if err != nil {
		return nil, err
	}

	plans := map[string]map[string]bool{}

	for _, a := range activeOf(assignments) {
		networkID, packageID := str(a["network_id"]), str(a["package_id"])
		if networkID == "" || packageID == "" {
			continue
		}

		if plans[networkID] == nil {
			plans[networkID] = map[string]bool{}
		}

		plans[networkID][packageID] = true
	}

	return zeroFilled(networks, func(networkID string) float64 {
		return float64(len(plans[networkID]))
	}), nil
}

// PackageDataUsed (PACKAGE_DATA_USED @ scope network_id+package_id): window
// usage per sim attributed to the sim's currently-active package.
func PackageDataUsed(win schema.Window, in Datasets, spec schema.KpiSpec) ([]Result, error) {
	usage, ok := in["usage"]
	if !ok {
		return nil, fmt.Errorf("PACKAGE_DATA_USED: missing input 'usage'")
	}

	assignments, ok := in["sim_packages"]
	if !ok {
		return nil, fmt.Errorf("PACKAGE_DATA_USED: missing input 'sim_packages'")
	}

	activePackageOfSim := map[string]string{}
	for _, a := range activeOf(assignments) {
		activePackageOfSim[str(a["sim_id"])] = str(a["package_id"])
	}

	bytes := map[[2]string]float64{}

	for _, rec := range usage {
		networkID, simID := str(rec["network_id"]), str(rec["sim_id"])

		packageID := activePackageOfSim[simID]
		if networkID == "" || packageID == "" {
			continue
		}

		bytes[[2]string{networkID, packageID}] += sumNumeric(rec["usage"])
	}

	results := make([]Result, 0, len(bytes))
	for key, b := range bytes {
		results = append(results, CountResult(
			map[string]string{"network_id": key[0], "package_id": key[1]}, b))
	}

	return results, nil
}

// --- shared helpers ---

func packageInputs(in Datasets, kpi string) (assignments, networks []map[string]interface{}, err error) {
	assignments, ok := in["sim_packages"]
	if !ok {
		return nil, nil, fmt.Errorf("%s: missing input 'sim_packages'", kpi)
	}

	networks, ok = in["networks"]
	if !ok {
		return nil, nil, fmt.Errorf("%s: missing input 'networks'", kpi)
	}

	return assignments, networks, nil
}

func activeOf(assignments []map[string]interface{}) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(assignments))

	for _, a := range assignments {
		if asBool(a["is_active"]) && !asBool(a["as_expired"]) {
			out = append(out, a)
		}
	}

	return out
}

func zeroFilled(networks []map[string]interface{}, value func(networkID string) float64) []Result {
	results := make([]Result, 0, len(networks))

	for _, n := range networks {
		networkID := str(n["network_id"])
		if networkID == "" {
			continue
		}

		results = append(results, CountResult(
			map[string]string{"network_id": networkID}, value(networkID)))
	}

	return results
}

func asBool(v interface{}) bool {
	switch t := v.(type) {
	case bool:
		return t
	case string:
		return t == "true" || t == "1"
	case float64:
		return t != 0
	default:
		return false
	}
}

func num(v interface{}) float64 {
	switch t := v.(type) {
	case float64:
		return t
	case string:
		var f float64
		if _, err := fmt.Sscanf(t, "%g", &f); err == nil {
			return f
		}

		return 0
	default:
		return 0
	}
}

// dataVolumeBytes converts a package's declared data allowance to bytes using
// its data unit (bytes|kilobytes|megabytes|gigabytes, or b/kb/mb/gb — parsed
// by the shared ukama converter). An absent or unrecognised unit is treated
// as bytes so a mis-tagged package still contributes its raw volume rather
// than silently dropping to zero.
func dataVolumeBytes(p map[string]interface{}) float64 {
	perUnit := float64(ukama.ReturnDataUnitsInBytes(ukama.ParseDataUnitType(str(p["data_unit"]))))
	if perUnit == 0 {
		perUnit = 1
	}

	return num(p["data_volume"]) * perUnit
}

// parseTime accepts the timestamp formats seen across the source systems:
// RFC3339 (registry/subscriber, e.g. package start_date) and Go's default
// time.Time string layout (payments paid_at, e.g.
// "2026-07-22 15:25:42.989109 +0000 UTC").
func parseTime(s string) (time.Time, error) {
	s = strings.TrimSpace(s)

	layouts := []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02T15:04:05Z0700",
		"2006-01-02 15:04:05.999999999 -0700 MST",
		"2006-01-02 15:04:05 -0700 MST",
	}

	for _, l := range layouts {
		if t, err := time.Parse(l, s); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unrecognized time %q", s)
}
