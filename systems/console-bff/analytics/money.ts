/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Monetary unit normalisation for the analytics subgraph.
 *
 * The analytics pipeline stores money in MINOR units (integer cents) for
 * arithmetic exactness — see analysis/pkg/algos/{packages,revenue}.go, which
 * are the only two places money is scaled by 100. Every other console-bff
 * module (payment, package, dashboard/commerceView, subscriberView) exposes
 * MAJOR units. This module is the presentation boundary that reconciles the
 * two, so every monetary value on the GraphQL surface is in major units.
 *
 * The `unit` field moves with the value: cells and KPI values that are
 * converted come back tagged with the org's currency ISO code (e.g. "usd")
 * instead of "cents". Clients that key off `unit === "cents"` — the console's
 * `kpiAmount`/`cellMoney` helpers do — therefore stop dividing on their own,
 * with no client change required. Never convert the value without rewriting
 * the unit: that renders every money figure 100x too small.
 */

/** Unit strings the analytics gateway uses for minor-unit money. */
const MINOR_UNITS = new Set(["cents"]);

/** Fallback unit when the request carries no org currency claim. */
const DEFAULT_MAJOR_UNIT = "major";

/** Report-column `format` that marks a cell as monetary, per the report spec. */
const MONEY_FORMAT = "money";

/**
 * Ops whose result is a count, not an amount. The gateway tags `unit` per KPI
 * rather than per op, so e.g. REVENUE/COUNT arrives labelled "cents" while
 * actually being a purchase count. Converting it would divide a count by 100.
 */
const NON_MONETARY_OPS = new Set(["COUNT"]);

const isMinorUnit = (unit?: string | null): boolean =>
  MINOR_UNITS.has((unit ?? "").toLowerCase());

const isCountOp = (op?: string | null): boolean =>
  NON_MONETARY_OPS.has((op ?? "").toUpperCase());

/**
 * Minor units -> major units.
 *
 * A pure scale, deliberately without rounding. Consumers verify relationships
 * *between* the values we convert — ukama-lab's `kpi_contract` with
 * `require_trend_consistency` recomputes `value - prevValue` and compares it
 * to `changeAbs` — so rounding each member independently would break the
 * invariant whenever AVG rollups produce fractional cents
 * (aggregator/pkg/rollup/engine.go). Scaling alone preserves it.
 *
 * Integer cents, which is every SUM/LAST op, divide to clean two-decimal
 * values. Fractional cents inherit the usual binary-float representation, well
 * inside the tolerances every consumer already applies.
 */
const toMajor = (value: number): number => value / 100;

const toMajorOrNull = (value?: number | null): number | null | undefined =>
  value == null ? value : toMajor(value);

/** The unit string to emit in place of "cents". */
export const majorUnit = (currency?: string): string =>
  (currency ?? "").trim().toLowerCase() || DEFAULT_MAJOR_UNIT;

interface TrendLike {
  direction?: string;
  changePct?: number | null;
  changeAbs?: number | null;
  prevValue?: number | null;
  hasPrevious?: boolean;
}

/**
 * `changeAbs` and `prevValue` carry the same unit as the value they describe
 * and must convert with it. `changePct` is a ratio and must not.
 */
const normalizeTrend = <T extends TrendLike>(trend?: T): T | undefined =>
  trend == null
    ? trend
    : {
        ...trend,
        changeAbs: toMajorOrNull(trend.changeAbs),
        prevValue: toMajorOrNull(trend.prevValue),
      };

interface KpiValueLike {
  value: number;
  op?: string;
  unit?: string;
  trend?: TrendLike;
}

/**
 * Normalises one KPI value in place of the gateway's minor units. Non-monetary
 * KPIs (count/bytes/percent) and count ops pass through untouched; a count op
 * additionally gets its mislabelled "cents" unit corrected to "count".
 */
export const normalizeKpiValue = <T extends KpiValueLike>(
  value: T,
  currency?: string
): T => {
  if (!isMinorUnit(value.unit)) return value;
  if (isCountOp(value.op)) return { ...value, unit: "count" };
  return {
    ...value,
    value: toMajor(value.value),
    unit: majorUnit(currency),
    trend: normalizeTrend(value.trend),
  };
};

interface ReportCellLike {
  value: number;
  unit?: string;
  format?: string;
  trend?: TrendLike;
}

/**
 * Normalises one report cell.
 *
 * `format === "money"` is the primary discriminator here, not `unit`: an
 * entity with no rollup rows in the window gets an empty `unit` from the
 * composer's zero-cell early return (aggregator/pkg/performance/composer.go),
 * while `format` is always set from the report spec. Such cells carry value 0,
 * so the conversion is a no-op numerically, but tagging them consistently
 * keeps the unit meaningful for clients.
 */
export const normalizeReportCell = <T extends ReportCellLike>(
  cell: T,
  currency?: string
): T => {
  const monetary =
    isMinorUnit(cell.unit) ||
    (cell.format ?? "").toLowerCase() === MONEY_FORMAT;
  if (!monetary) return cell;
  return {
    ...cell,
    value: toMajor(cell.value),
    unit: majorUnit(currency),
    trend: normalizeTrend(cell.trend),
  };
};
