/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

import { useEffect, useState } from 'react';
import Skeleton from '@mui/material/Skeleton';
import { useMetricsRangeQuery } from '@/client/graphql/range-metrics.generated';
import MetricLineChart, {
  ChartMessage,
  thresholdLegendRows,
} from '@/components/MetricLineChart';
import RangeToggle from '@/components/RangeToggle';
import SectionCard from '@/components/SectionCard';
import { metricLabel } from '@/lib/labels';
import { type LatestEntry, seriesLatest } from '@/lib/metrics';
import { RANGE_SECONDS, type Range } from '@/lib/ranges';

const DEFAULT_HEIGHT = 300;

function LegendDot({ color, label }: { color: string; label: string }) {
  return (
    <span
      style={{
        display: 'inline-flex',
        alignItems: 'center',
        gap: 6,
        fontSize: 12,
        color: 'var(--uk-ink-2)',
      }}
    >
      <span style={{ width: 9, height: 9, borderRadius: 3, background: color }} />
      {label}
    </span>
  );
}

/**
 * One self-fetching metric chart: its own Day/Week/Month filter, loading /
 * empty / error states, line chart and threshold legend. Used by both the node
 * and site detail screens.
 *
 * - `nodeId` scopes the query to a node (omit for org/site scope).
 * - `off` zeroes the series (offline node visual).
 * - `onLatest` reports the series' latest sample so a left-rail value can reuse
 *   it without a second query.
 * - `titleOverride` forces the card title; otherwise the server label / a
 *   humanized key is used.
 */
export default function MetricChartCard({
  metricKey,
  nodeId,
  fallbackLabel,
  titleOverride,
  off = false,
  height = DEFAULT_HEIGHT,
  onLatest,
}: {
  metricKey: string;
  nodeId?: string | null;
  fallbackLabel?: string;
  titleOverride?: string;
  off?: boolean;
  height?: number;
  onLatest?: (key: string, entry: LatestEntry) => void;
}) {
  const [range, setRange] = useState<Range>('Day');
  const [nowSec] = useState(() => Math.floor(Date.now() / 1000));
  const to = nowSec;
  const from = nowSec - RANGE_SECONDS[range];
  const { data, loading, error } = useMetricsRangeQuery({
    variables: {
      data: { keys: [metricKey], from, to, ...(nodeId ? { nodeId } : {}) },
    },
  });

  const m = data?.metricsRange.metrics?.[0];
  // The chart already holds the series — feed its latest sample upward so the
  // same metric isn't fetched twice (chart + a separate latest query).
  useEffect(() => {
    if (m && onLatest) onLatest(metricKey, seriesLatest(m));
  }, [m, metricKey, onLatest]);

  const hasData = !!m && m.values.length > 0 && m.success !== false;
  const values: [number, number][] = hasData
    ? off
      ? m!.values.map((v) => [v[0] ?? 0, 0])
      : m!.values.map((v) => [v[0] ?? 0, v[1] ?? 0])
    : [];
  const legend = thresholdLegendRows(m?.threshold ?? null, m?.unit);
  const title = titleOverride ?? metricLabel(m?.label, metricKey, fallbackLabel);

  return (
    <SectionCard
      title={title}
      right={<RangeToggle value={range} onChange={setRange} />}
    >
      {error ? (
        <ChartMessage kind="error" message={error.message} height={height} />
      ) : loading && !m ? (
        <Skeleton variant="rounded" sx={{ height }} />
      ) : !hasData ? (
        <ChartMessage kind="empty" height={height} />
      ) : (
        <>
          <MetricLineChart
            values={values}
            title={title}
            unit={m?.unit}
            format={m?.format}
            threshold={m?.threshold ?? null}
            height={height}
          />
          <div
            style={{
              display: 'flex',
              gap: 18,
              justifyContent: 'center',
              marginTop: 10,
              flexWrap: 'wrap',
            }}
          >
            {legend.map((l) => (
              <LegendDot key={l.label} {...l} />
            ))}
          </div>
        </>
      )}
    </SectionCard>
  );
}
