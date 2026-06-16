/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Helpers for the analytics `KpiDto[]` arrays. The analytics service returns
 * KPIs as a flat, string-keyed list (key/value/formatted/delta/...), so every
 * screen looks values up by key rather than by a fixed field. Keeping this in
 * one place means an unrecognised or not-yet-emitted key degrades to "—"
 * everywhere consistently — see docs/analytics-backend-gaps.md for the keys
 * each screen expects.
 */

/**
 * Central registry of analytics KPI keys the console reads, so the key strings
 * live in one place (no per-screen duplication). See docs/analytics-backend-gaps.md.
 */
export const KPI_KEYS = {
  networkUptime: 'network_uptime',
  activeCustomers: 'active_customers',
  customersTotal: 'customers_total',
  dataUsage: 'data_usage',
  revenueMonth: 'revenue_month',
  revenueCollected: 'revenue_collected',
  revenuePrevMonth: 'revenue_prev_month',
  revenuePending: 'revenue_pending',
} as const;

export interface Kpi {
  key: string;
  value: number;
  formatted?: string | null;
  delta?: number | null;
  deltaPeriod?: string | null;
  stale?: boolean | null;
  asOf?: string | null;
}

/** The KPI for `key`, or undefined if the backend didn't emit it. */
export const kpiByKey = (
  kpis: readonly Kpi[] | undefined,
  key: string,
): Kpi | undefined => kpis?.find((k) => k.key === key);

/**
 * Display string for a KPI: prefers the server-formatted value, else the raw
 * number, else a dash. Pass a `fallbackFormat` to render the raw number when
 * the server didn't pre-format (e.g. money/percent).
 */
export const kpiText = (
  kpis: readonly Kpi[] | undefined,
  key: string,
  fallbackFormat?: (value: number) => string,
): string => {
  const k = kpiByKey(kpis, key);
  if (!k) return '—';
  if (k.formatted) return k.formatted;
  return fallbackFormat ? fallbackFormat(k.value) : String(k.value);
};

/** Raw numeric value for a KPI, or undefined when absent. */
export const kpiValue = (
  kpis: readonly Kpi[] | undefined,
  key: string,
): number | undefined => kpiByKey(kpis, key)?.value;

/**
 * Money KPI display: formats the KPI's raw numeric value with the caller's
 * formatter (org currency symbol), deliberately IGNORING the backend's
 * pre-formatted string — that string hardcodes a "$", so it must not be used
 * for currency. Returns "—" when the KPI is absent.
 */
export const kpiAmount = (
  kpis: readonly Kpi[] | undefined,
  key: string,
  formatMoney: (value: number) => string,
): string => {
  const k = kpiByKey(kpis, key);
  return k ? formatMoney(k.value) : '—';
};
