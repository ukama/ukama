/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain it at http://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package rest

// GetKpisRequest: latest values for KPIs at a span.
type GetKpisRequest struct {
	Keys      string `form:"keys" query:"keys" binding:"required" validate:"required"` // csv KPI keys
	Span      string `form:"span" query:"span"`                                        // daily|weekly|monthly (default daily)
	Op        string `form:"op" query:"op"`                                            // optional; defaults per KPI spec
	NetworkId string `form:"network_id" query:"network_id"`                            // optional scope filter
	SiteId    string `form:"site_id" query:"site_id"`                                   // optional scope filter (with network_id, e.g. SITE_UPTIME)
}

// GetKpiTimeSeriesRequest: one value per span bucket over a range.
type GetKpiTimeSeriesRequest struct {
	Keys      string `form:"keys" query:"keys" binding:"required" validate:"required"`
	Span      string `form:"span" query:"span"`
	Op        string `form:"op" query:"op"`
	From      string `form:"from" query:"from"` // RFC3339 inclusive
	To        string `form:"to" query:"to"`     // RFC3339 exclusive
	NetworkId string `form:"network_id" query:"network_id"`
	SiteId    string `form:"site_id" query:"site_id"`
}

// GetPerformanceReportRequest: resource performance table from the latest
// available KPI values.
type GetPerformanceReportRequest struct {
	Report    string `path:"report" json:"report" validate:"required"`
	Span      string `form:"span" query:"span"` // daily|weekly|monthly (default daily)
	NetworkId string `form:"network_id" query:"network_id"`
	Top       int32  `form:"top" query:"top"`
}

// GetKpiBreakdownRequest: top-N scope values for one KPI.
type GetKpiBreakdownRequest struct {
	Key  string `form:"key" query:"key" binding:"required" validate:"required"`
	Span string `form:"span" query:"span"`
	Op   string `form:"op" query:"op"`
	By   string `form:"by" query:"by" binding:"required" validate:"required"` // e.g. network_id
	Top  int32  `form:"top" query:"top"`
}
