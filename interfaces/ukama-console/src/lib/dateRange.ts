/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Shared date-range filter for the analytics screens' DateChip. The analytics
 * gateway rolls KPIs up per span (daily/weekly/monthly), so the three ranges
 * map onto the closest span rather than an arbitrary from/to window:
 *   Today       -> daily
 *   Last 7 days -> weekly
 *   Last 30 days-> monthly
 * Keeping the option list + mapping here means every page shows the same chip.
 */
export const DATE_RANGES = ['Today', 'Last 7 days', 'Last 30 days'] as const;

export type DateRange = (typeof DATE_RANGES)[number];

export const DEFAULT_RANGE: DateRange = 'Today';

export type KpiSpan = 'daily' | 'weekly' | 'monthly';

/** Map a DateChip range label to the KPI rollup span the gateway supports. */
export const rangeToSpan = (range: string): KpiSpan => {
  switch (range) {
    case 'Last 7 days':
      return 'weekly';
    case 'Last 30 days':
      return 'monthly';
    default:
      return 'daily';
  }
};
