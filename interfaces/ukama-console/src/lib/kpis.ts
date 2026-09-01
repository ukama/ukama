/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Helpers for the analytics `getKpiValues` result. The analytics gateway
 * returns KPIs as a flat, string-keyed list of `KpiValueDto` (kpi/value/span/
 * unit/symbol/trend/...), so every screen looks values up by key rather than by
 * a fixed field. Keeping this in one place means an unrecognised or
 * not-yet-emitted key degrades to "—" everywhere consistently.
 */

/**
 * Central registry of the analytics KPI keys the console reads, so the key
 * strings live in one place (no per-screen duplication). These map 1:1 to the
 * gateway's KPI registry — see systems/analytics/aggregator/configs/kpis.
 */
export const KPI_KEYS = {
  networkUptime: 'NETWORK_UPTIME',
  siteUptime: 'SITE_UPTIME',
  sitesOnline: 'SITES_ONLINE',
  // ACTIVE_CUSTOMERS counts SIMs in active service; CUSTOMERS counts every
  // SIM whatever its service state (active + inactive). Both are per network.
  activeCustomers: 'ACTIVE_CUSTOMERS',
  customers: 'CUSTOMERS',
  paidCustomers: 'PAID_CUSTOMERS',
  revenue: 'REVENUE',
  mrr: 'MRR',
  arpu: 'ARPU',
  // DATA_USAGE is the single generic usage KPI (it replaced USAGE_BY_NETWORK).
  // It is scoped per network×site×package×assignment×iccid; passing a
  // networkId filter folds it to one per-network total, so "by network" is a
  // read-time filter rather than a separate key.
  dataUsage: 'DATA_USAGE',
  dataSold: 'DATA_SOLD',
  packageSales: 'PACKAGE_SALES',
} as const;

/** Period-over-period comparison carried by a KPI value. */
export interface KpiTrend {
  direction?: string | null;
  changePct?: number | null;
  changeAbs?: number | null;
  prevValue?: number | null;
  hasPrevious?: boolean | null;
}

/** A single KPI reading, matching the `KpiValueDto` GraphQL shape. */
export interface Kpi {
  kpi: string;
  value: number;
  span?: string | null;
  unit?: string | null;
  symbol?: string | null;
  isPartial?: boolean | null;
  trend?: KpiTrend | null;
}

/** The KPI for `key`, or undefined if the backend didn't emit it. */
export const kpiByKey = (
  kpis: readonly Kpi[] | undefined,
  key: string,
): Kpi | undefined => kpis?.find((k) => k.kpi === key);

/** Raw numeric value for a KPI, or undefined when absent. */
export const kpiValue = (
  kpis: readonly Kpi[] | undefined,
  key: string,
): number | undefined => kpiByKey(kpis, key)?.value;

/** Period-over-period % change for a KPI (trend.changePct), or undefined. */
export const kpiDelta = (
  kpis: readonly Kpi[] | undefined,
  key: string,
): number | undefined => kpiByKey(kpis, key)?.trend?.changePct ?? undefined;

/** Period-over-period absolute change for a KPI (trend.changeAbs), or undefined. */
export const kpiChangeAbs = (
  kpis: readonly Kpi[] | undefined,
  key: string,
): number | undefined => kpiByKey(kpis, key)?.trend?.changeAbs ?? undefined;

/** Previous-period value for a KPI (trend.prevValue), or undefined. */
export const kpiPrev = (
  kpis: readonly Kpi[] | undefined,
  key: string,
): number | undefined => kpiByKey(kpis, key)?.trend?.prevValue ?? undefined;

/**
 * Display string for a KPI: the raw numeric value rendered with the caller's
 * `fallbackFormat` (e.g. percent / units), else the plain number, else a dash.
 * The gateway does not pre-format values, so a formatter should be passed for
 * anything other than a bare count.
 */
export const kpiText = (
  kpis: readonly Kpi[] | undefined,
  key: string,
  fallbackFormat?: (value: number) => string,
): string => {
  const k = kpiByKey(kpis, key);
  if (!k) return '—';
  return fallbackFormat ? fallbackFormat(k.value) : String(k.value);
};

/**
 * Money KPI display: formats the KPI's value with the caller's formatter (org
 * currency symbol). Money KPIs are reported in MINOR units (unit: "cents"), so
 * divide by 100 before formatting. Returns "—" when the KPI is absent.
 */
export const kpiAmount = (
  kpis: readonly Kpi[] | undefined,
  key: string,
  formatMoney: (value: number) => string,
): string => {
  const k = kpiByKey(kpis, key);
  if (!k) return '—';
  const value = (k.unit ?? '').toLowerCase() === 'cents' ? k.value / 100 : k.value;
  return formatMoney(value);
};

/**
 * Money display of a KPI's PREVIOUS-period value (trend.prevValue), cents-aware
 * like kpiAmount. Returns "—" when the KPI or its previous value is absent.
 */
export const kpiAmountPrev = (
  kpis: readonly Kpi[] | undefined,
  key: string,
  formatMoney: (value: number) => string,
): string => {
  const k = kpiByKey(kpis, key);
  const prev = k?.trend?.prevValue;
  if (k == null || prev == null) return '—';
  const value = (k.unit ?? '').toLowerCase() === 'cents' ? prev / 100 : prev;
  return formatMoney(value);
};

/** Value for a report cell column, or undefined when the column is absent. */
export const cellValue = (
  cells: readonly { column: string; value: number }[] | undefined,
  column: string,
): number | undefined => cells?.find((c) => c.column === column)?.value;

/**
 * Money value for a report cell in MAJOR currency units. The gateway reports
 * money columns in minor units (unit: "cents"), so divide by 100 when the cell
 * says so. Returns 0 when the column is absent.
 */
export const cellMoney = (
  cells:
    | readonly { column: string; value: number; unit?: string | null }[]
    | undefined,
  column: string,
): number => {
  const c = cells?.find((x) => x.column === column);
  if (!c) return 0;
  return (c.unit ?? '').toLowerCase() === 'cents' ? c.value / 100 : c.value;
};

/** Value for a report-row attribute, or undefined when the attribute is absent. */
export const attrValue = (
  attributes: readonly { key: string; value: string }[] | undefined,
  key: string,
): string | undefined => attributes?.find((a) => a.key === key)?.value;
