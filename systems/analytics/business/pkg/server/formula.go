/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package server

import (
	"fmt"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/ukama/ukama/systems/analytics/business/pb/gen"
	"github.com/ukama/ukama/systems/analytics/business/pkg/db"
)

// pctDelta is the percent change of current vs previous. Returns 0 when the
// previous value is 0 (no meaningful baseline).
func pctDelta(current, previous float64) float64 {
	if previous == 0 {
		return 0
	}

	return (current - previous) / previous * 100
}

// avgPurchase = revenue / purchases, guarded against division by zero.
func avgPurchase(revenue float64, purchases uint32) float64 {
	if purchases == 0 {
		return 0
	}

	return revenue / float64(purchases)
}

// mrr = sum over packages of ActiveSubscribers * Price.
func mrr(packages []db.PackageSnapshot) float64 {
	var total float64

	for _, p := range packages {
		total += float64(p.ActiveSubscribers) * p.Price
	}

	return total
}

// activeSubscribers = sum of ActiveSubscribers across package snapshots.
func activeSubscribers(packages []db.PackageSnapshot) uint32 {
	var total uint32

	for _, p := range packages {
		total += p.ActiveSubscribers
	}

	return total
}

// arpu = revenue over the window divided by the number of active customers.
// Returns 0 when there are no active customers.
func arpu(revenue float64, activeCustomers uint32) float64 {
	if activeCustomers == 0 {
		return 0
	}

	return revenue / float64(activeCustomers)
}

// topPlanByRevenue aggregates package rollups by package and returns the top
// plan's id, revenue and its share (%) of the total revenue. ok is false when
// there are no rollups or no revenue at all.
func topPlanByRevenue(rollups []db.BusinessPackageRollupDaily) (topId string, topRevenue float64, share float64, ok bool) {
	totals := map[string]float64{}
	var grandTotal float64

	for _, r := range rollups {
		id := r.PackageId.String()
		totals[id] += r.Revenue
		grandTotal += r.Revenue
	}

	for id, rev := range totals {
		if topId == "" || rev > topRevenue || (rev == topRevenue && id < topId) {
			topId = id
			topRevenue = rev
		}
	}

	if topId == "" || grandTotal == 0 {
		return "", 0, 0, false
	}

	return topId, topRevenue, topRevenue / grandTotal * 100, true
}

// makeKpi builds a Kpi message with formatted display value.
func makeKpi(key string, value float64, formatted string, delta float64, deltaPeriod string, asOf time.Time) *pb.Kpi {
	return &pb.Kpi{
		Key:         key,
		Value:       value,
		Formatted:   formatted,
		Delta:       delta,
		DeltaPeriod: deltaPeriod,
		// TODO: wire stale flag from analytics_refresh_states once the
		// collector exposes refresh state to read services.
		Stale: false,
		AsOf:  timestamppb.New(asOf),
	}
}

func formatMoney(v float64) string {
	return fmt.Sprintf("$%.2f", v)
}

func formatCount(v uint32) string {
	return fmt.Sprintf("%d", v)
}

func formatPercent(v float64) string {
	return fmt.Sprintf("%.1f%%", v)
}
