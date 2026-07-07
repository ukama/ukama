/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import type { MetricsRangeQuery } from '@/client/graphql/range-metrics.generated';

/** One metric series straight from the BFF (metricsRange). */
export type MetricSeries = MetricsRangeQuery['metricsRange']['metrics'][number];

/** Latest value + presentation metadata for one KPI, derived from a series. */
export type LatestEntry = {
  value: number;
  success: boolean;
  label?: string | null;
  unit?: string | null;
  format?: string | null;
};

/** Collapse a series to its CURRENT KPI value: the value at the most recent
 *  timestamp. If that latest sample is a gap-fill placeholder (-1) or missing,
 *  the metric isn't reporting now → success:false, so the rail shows "—"
 *  instead of a stale earlier reading. */
export const seriesLatest = (m: MetricSeries): LatestEntry => {
  const vals = m.values ?? [];
  const raw = vals.length ? vals[vals.length - 1]?.[1] : null;
  const latest = raw != null && raw !== -1 ? raw : null;
  return {
    value: latest ?? 0,
    success: m.success !== false && latest != null,
    label: m.label,
    unit: m.unit,
    format: m.format,
  };
};
