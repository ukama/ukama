/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain it at http://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package rest

// ScopeParams are the generic scope filters accepted by the KPI read
// endpoints. Every non-empty param becomes a scope filter key; the
// aggregator validates keys against the requested KPIs' scope dimensions
// (400 on a key no requested KPI carries — never a silent empty result).
//
// Filtering also sets the answer's grain: for component ops (SUM/COUNT/AVG/
// MIN/MAX) matching rows are folded into one row per distinct combination
// of the filter keys — DATA_USAGE?network_id=X returns ONE total for the
// network, not one row per sim-series. No group_by needed for that.
//
// package_id is the CATALOG package (data plan product); sim_package_id is
// a sim's package assignment instance.
type ScopeParams struct {
	NetworkId    string `form:"network_id" query:"network_id"`
	SiteId       string `form:"site_id" query:"site_id"`
	PackageId    string `form:"package_id" query:"package_id"`
	SimPackageId string `form:"sim_package_id" query:"sim_package_id"`
	Iccid        string `form:"iccid" query:"iccid"`
}

// QueryRequest is the single question shape: KPIs + filters + dimensions +
// range + granularity. Aggregation comes from each KPI's kind — no op to
// pick. Filters always fold: DATA_USAGE?network_id=X&range=this_month is
// ONE number; add group_by=package_id for a per-package breakdown; add
// granularity=day for a chart; sort=-value&top=5 for a top-N.
type QueryRequest struct {
	Kpis string `form:"kpis" query:"kpis" binding:"required" validate:"required"` // csv KPI keys
	// Range: today|this_week|this_month|last_24h|last_7d|last_30d
	// (default this_month), or from/to RFC3339 for a custom period.
	Range string `form:"range" query:"range"`
	From  string `form:"from" query:"from"`
	To    string `form:"to" query:"to"`
	// Granularity: total (default, one point) | day | week | month (series).
	Granularity string `form:"granularity" query:"granularity"`
	GroupBy     string `form:"group_by" query:"group_by"` // csv dimensions
	Sort        string `form:"sort" query:"sort"`         // value|-value
	Top         int32  `form:"top" query:"top"`
	Agg         string `form:"agg" query:"agg"` // expert override: sum|avg|min|max|count
	ScopeParams
}

// GetKpisRequest: latest values for KPIs at a span.
type GetKpisRequest struct {
	Keys string `form:"keys" query:"keys" binding:"required" validate:"required"` // csv KPI keys
	Span string `form:"span" query:"span"`                                        // daily|weekly|monthly|last_24h|last_7d|last_30d (default daily)
	Op   string `form:"op" query:"op"`                                            // optional; defaults per KPI spec
	// GroupBy is the explicit fold grain for breakdown-style reads WITHOUT a
	// filter (csv of scope keys), e.g. DATA_USAGE group_by=network_id → one
	// row per network. With a filter, the fold happens automatically at the
	// filter's grain; group_by overrides it when set.
	GroupBy string `form:"group_by" query:"group_by"`
	ScopeParams
}

// GetKpiTimeSeriesRequest: one value per span bucket over a range.
type GetKpiTimeSeriesRequest struct {
	Keys    string `form:"keys" query:"keys" binding:"required" validate:"required"`
	Span    string `form:"span" query:"span"`
	Op      string `form:"op" query:"op"`
	From    string `form:"from" query:"from"` // RFC3339 inclusive
	To      string `form:"to" query:"to"`     // RFC3339 exclusive
	GroupBy string `form:"group_by" query:"group_by"`
	ScopeParams
}

// GetPerformanceReportRequest: resource performance table from the latest
// available KPI values.
type GetPerformanceReportRequest struct {
	Report    string `path:"report" json:"report" validate:"required"`
	Span      string `form:"span" query:"span"` // daily|weekly|monthly (default daily)
	NetworkId string `form:"network_id" query:"network_id"`
	Top       int32  `form:"top" query:"top"`
}

// GetKpiBreakdownRequest: top-N scope values for one KPI, optionally
// filtered (e.g. top packages WITHIN a network: by=package_id&network_id=X).
type GetKpiBreakdownRequest struct {
	Key  string `form:"key" query:"key" binding:"required" validate:"required"`
	Span string `form:"span" query:"span"`
	Op   string `form:"op" query:"op"`
	By   string `form:"by" query:"by" binding:"required" validate:"required"` // e.g. network_id
	Top  int32  `form:"top" query:"top"`
	ScopeParams
}
