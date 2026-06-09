/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/**
 * Time-series metric chart (Recharts). Renders a smooth line with a time
 * X-axis, value Y-axis, threshold reference lines, a dark hover tooltip that
 * names the zone (normal/high/critical), and a threshold-derived legend.
 * Driven entirely by the BFF metric shape — no per-metric config here.
 */
import {
  CartesianGrid,
  Line,
  LineChart,
  ReferenceLine,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from 'recharts';

import { EmptyState } from './EmptyState';

export interface MetricThreshold {
  min: number;
  normal: number;
  max: number;
}

/** Centered fallback shown inside a chart card instead of an empty plot. */
export function ChartMessage({
  kind,
  message,
  height = 300,
}: {
  kind: 'error' | 'empty';
  message?: string;
  height?: number | string;
}) {
  return (
    <div style={{ height, display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
      <EmptyState
        art={kind === 'error' ? 'error' : 'search'}
        title={kind === 'error' ? "Couldn't load metric" : 'No data'}
        sub={
          message ??
          (kind === 'error'
            ? 'The metric service didn’t respond.'
            : 'No data for the selected period.')
        }
      />
    </div>
  );
}

export interface MetricLineChartProps {
  /** [timestampSeconds, value] pairs, ascending. */
  values: [number, number][];
  /** Series title used in the tooltip (e.g. "Temperature"). */
  title: string;
  unit?: string | null;
  format?: string | null;
  threshold?: MetricThreshold | null;
  height?: number | string;
}

const fmtTime = (ms: number) =>
  new Date(ms).toLocaleTimeString('en-US', { hour: 'numeric', minute: '2-digit' });

const fmtFull = (ms: number) =>
  new Date(ms).toLocaleString('en-US', {
    weekday: 'short',
    month: 'short',
    day: 'numeric',
    year: 'numeric',
    hour: 'numeric',
    minute: '2-digit',
  });

const zoneOf = (v: number, t?: MetricThreshold | null): string | null => {
  if (!t) return null;
  if (v < t.normal) return 'normal';
  if (v < t.max) return 'high';
  return 'critical';
};

/** Legend rows mirroring the legacy console: Below N normal / N–max high /
 *  Above max critical (colours grey / pink / red). */
export function thresholdLegendRows(
  t?: MetricThreshold | null,
  unit?: string | null,
): { color: string; label: string }[] {
  if (!t) return [{ color: 'var(--uk-ink-3)', label: 'Trend' }];
  const u = unit ? ` ${unit}` : '';
  return [
    { color: 'var(--uk-ink-3)', label: `Below ${t.normal}${u}: Normal` },
    { color: 'var(--uk-orange)', label: `${t.normal}–${t.max}${u}: High` },
    { color: 'var(--uk-error)', label: `Above ${t.max}${u}: Critical` },
  ];
}

export default function MetricLineChart({
  values,
  title,
  unit,
  format,
  threshold,
  height = 300,
}: MetricLineChartProps) {
  const data = values.map(([ts, value]) => ({ ts: ts * 1000, value }));
  const fmtVal = (v: number) =>
    format === 'decimal' ? v.toFixed(2) : String(Math.round(v));
  const u = unit ? ` ${unit}` : '';

  const ys = data.map((d) => d.value);
  const dataMax = ys.length ? Math.max(...ys) : 1;
  const top = Math.ceil(Math.max(threshold?.max ?? 0, dataMax) * 1.06);

  return (
    <ResponsiveContainer width="100%" height={height}>
      <LineChart data={data} margin={{ top: 10, right: 16, bottom: 4, left: 0 }}>
        <CartesianGrid strokeDasharray="4 4" stroke="var(--uk-line-soft)" />
        <XAxis
          dataKey="ts"
          type="number"
          scale="time"
          domain={['dataMin', 'dataMax']}
          tickFormatter={fmtTime}
          tick={{ fontSize: 12, fill: 'var(--uk-ink-3)' }}
          stroke="var(--uk-line)"
          minTickGap={48}
          tickMargin={8}
        />
        <YAxis
          domain={[0, top]}
          tickFormatter={(v: number) => `${Math.round(v)}${unit ?? ''}`}
          tick={{ fontSize: 12, fill: 'var(--uk-ink-3)' }}
          stroke="var(--uk-line)"
          width={72}
          tickMargin={6}
        />
        {threshold && (
          <>
            <ReferenceLine
              y={threshold.normal}
              stroke="var(--uk-orange)"
              strokeDasharray="5 5"
              strokeOpacity={0.55}
            />
            <ReferenceLine
              y={threshold.max}
              stroke="var(--uk-error)"
              strokeDasharray="5 5"
              strokeOpacity={0.55}
            />
          </>
        )}
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
            const zone = zoneOf(v, threshold);
            return (
              <div
                style={{
                  background: '#1b2430',
                  color: '#fff',
                  borderRadius: 10,
                  padding: '10px 14px',
                  boxShadow: '0 8px 28px rgba(0,0,0,.35)',
                  maxWidth: 280,
                }}
              >
                <div style={{ fontWeight: 600, fontSize: 13 }}>
                  {fmtFull(Number(label))}
                </div>
                <div style={{ fontSize: 13, color: '#c7ced8', marginTop: 3 }}>
                  {title}: {fmtVal(v)}
                  {u}
                  {zone ? ` (${zone})` : ''}
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
          dot={false}
          activeDot={{ r: 5, stroke: 'var(--uk-success)', strokeWidth: 2, fill: '#fff' }}
          isAnimationActive={false}
        />
      </LineChart>
    </ResponsiveContainer>
  );
}
