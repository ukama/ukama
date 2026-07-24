/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/**
 * Revenue trend line chart (Recharts). Plots one point per weekly rollup
 * bucket from the analytics `getKpiTimeSeries` result, with a date X-axis and
 * a money Y-axis. Presentation only — the caller maps KPI buckets to points.
 */
import {
  CartesianGrid,
  Line,
  LineChart,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from 'recharts';

export interface TrendPoint {
  /** Bucket start, epoch milliseconds. */
  t: number;
  /** Value in major units (e.g. dollars). */
  value: number;
}

const fmtDate = (ms: number) =>
  new Date(ms).toLocaleDateString('en-US', { month: 'short', day: 'numeric' });

export default function RevenueTrendChart({
  points,
  formatValue,
  height = 300,
}: {
  points: TrendPoint[];
  formatValue: (v: number) => string;
  height?: number | string;
}) {
  const ys = points.map((p) => p.value);
  const top = Math.ceil((ys.length ? Math.max(...ys) : 1) * 1.12) || 1;

  return (
    <ResponsiveContainer width="100%" height={height}>
      <LineChart data={points} margin={{ top: 10, right: 16, bottom: 4, left: 0 }}>
        <CartesianGrid strokeDasharray="4 4" stroke="var(--uk-line-soft)" />
        <XAxis
          dataKey="t"
          type="number"
          scale="time"
          domain={['dataMin', 'dataMax']}
          tickFormatter={fmtDate}
          tick={{ fontSize: 12, fill: 'var(--uk-ink-3)' }}
          stroke="var(--uk-line)"
          minTickGap={40}
          tickMargin={8}
        />
        <YAxis
          domain={[0, top]}
          tickFormatter={formatValue}
          tick={{ fontSize: 12, fill: 'var(--uk-ink-3)' }}
          stroke="var(--uk-line)"
          width={72}
          tickMargin={6}
        />
        <Tooltip
          cursor={{ stroke: 'var(--uk-success)', strokeWidth: 1 }}
          content={(props: {
            active?: boolean;
            label?: string | number;
            payload?: { value?: number | string }[];
          }) => {
            const { active, payload, label } = props;
            if (!active || !payload || payload.length === 0) return null;
            const v = Number(payload[0]?.value ?? 0);
            return (
              <div
                style={{
                  background: '#1b2430',
                  color: '#fff',
                  borderRadius: 10,
                  padding: '10px 14px',
                  boxShadow: '0 8px 28px rgba(0,0,0,.35)',
                }}
              >
                <div style={{ fontWeight: 600, fontSize: 13 }}>
                  {fmtDate(Number(label))}
                </div>
                <div style={{ fontSize: 13, color: '#c7ced8', marginTop: 3 }}>
                  Revenue: {formatValue(v)}
                </div>
              </div>
            );
          }}
        />
        <Line
          type="monotone"
          dataKey="value"
          stroke="var(--uk-ac)"
          strokeWidth={2.5}
          dot={{ r: 3, fill: 'var(--uk-ac)' }}
          activeDot={{ r: 5, stroke: 'var(--uk-success)', strokeWidth: 2, fill: '#fff' }}
          isAnimationActive={false}
        />
      </LineChart>
    </ResponsiveContainer>
  );
}
