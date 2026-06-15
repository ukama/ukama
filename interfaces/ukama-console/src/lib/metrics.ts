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

/** Collapse a series to its latest KPI value: the most recent sample that
 *  isn't a gap-fill placeholder (-1). The chart's last point IS the latest. */
export const seriesLatest = (m: MetricSeries): LatestEntry => {
  const vals = m.values ?? [];
  let last: number | null = null;
  for (let i = vals.length - 1; i >= 0; i--) {
    const v = vals[i]?.[1];
    if (v != null && v !== -1) {
      last = v;
      break;
    }
  }
  return {
    value: last ?? 0,
    success: m.success !== false && last != null,
    label: m.label,
    unit: m.unit,
    format: m.format,
  };
};
