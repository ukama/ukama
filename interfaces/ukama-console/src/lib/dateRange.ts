/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Shared date-range filter for the analytics screens' DateChip. The three
 * ranges are true ROLLING windows ending "now":
 *   Last 24h    -> last_24h
 *   Last 7 days -> last_7d
 *   Last 30 days-> last_30d
 * The analytics gateway serves these on /kpis/values by aggregating the exact
 * 5-min KPI windows over the trailing duration (see aggregator rolling.go), so
 * the number matches the label — unlike the old calendar spans, where "Last 7
 * days" actually returned the current ISO week.
 *
 * The performance-report endpoint is NOT driven by this filter: its composer
 * aggregates over its own config-driven rolling window (default 8 weeks) so
 * per-entity stats stay stable. Report queries omit span entirely.
 */
export const DATE_RANGES = ['Last 24h', 'Last 7 days', 'Last 30 days'] as const;

export type DateRange = (typeof DATE_RANGES)[number];

export const DEFAULT_RANGE: DateRange = 'Last 24h';

/** Rolling-window span tokens understood by the gateway's /kpis/values. */
export type KpiSpan = 'last_24h' | 'last_7d' | 'last_30d';

/** Map a DateChip range label to the rolling KPI span (for getKpiValues). */
export const rangeToSpan = (range: string): KpiSpan => {
  switch (range) {
    case 'Last 7 days':
      return 'last_7d';
    case 'Last 30 days':
      return 'last_30d';
    default:
      return 'last_24h';
  }
};

// The performance-report endpoint no longer takes a UI-filter span: the
// composer aggregates over its own config-driven rolling window (default 8
// weeks) so per-entity report stats are stable. Screens omit span on report
// queries and show the window returned by the API (see reportWindowLabel).
