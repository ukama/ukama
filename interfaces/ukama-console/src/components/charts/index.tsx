/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/**
 * Lightweight SVG charts — ported 1:1 from the design prototype
 * (ds.jsx, biz-common.jsx, node-site-detail.jsx). Zero dependencies;
 * colors are theme CSS variables so charts recolor with mode/accent.
 */
import { useId } from 'react';

/* ---- Sparkline ---- */
export function Spark({
  data,
  color = 'var(--uk-ac)',
  w = 96,
  h = 30,
  fill = true,
  strokeW = 1.8,
}: {
  data: number[];
  color?: string;
  w?: number;
  h?: number;
  fill?: boolean;
  strokeW?: number;
}) {
  const gid = useId();
  const max = Math.max(...data);
  const min = Math.min(...data);
  const span = max - min || 1;
  const pts = data.map(
    (v, i) =>
      [
        (i / (data.length - 1)) * w,
        h - 3 - ((v - min) / span) * (h - 6),
      ] as const,
  );
  const d = pts
    .map((p, i) => (i ? 'L' : 'M') + p[0].toFixed(1) + ' ' + p[1].toFixed(1))
    .join(' ');
  const area = d + ` L${w} ${h} L0 ${h} Z`;
  const last = pts[pts.length - 1];

  return (
    <svg
      width={w}
      height={h}
      viewBox={`0 0 ${w} ${h}`}
      style={{ display: 'block', overflow: 'visible' }}
      aria-hidden="true"
    >
      <defs>
        <linearGradient id={gid} x1="0" y1="0" x2="0" y2="1">
          <stop offset="0" stopColor={color} stopOpacity="0.18" />
          <stop offset="1" stopColor={color} stopOpacity="0" />
        </linearGradient>
      </defs>
      {fill && <path d={area} fill={`url(#${gid})`} />}
      <path
        d={d}
        fill="none"
        stroke={color}
        strokeWidth={strokeW}
        strokeLinecap="round"
        strokeLinejoin="round"
      />
      {last && <circle cx={last[0]} cy={last[1]} r="2.6" fill={color} />}
    </svg>
  );
}

/* ---- Bar mini chart ---- */
export function MiniBars({
  data,
  color = 'var(--uk-ac)',
  w = 150,
  h = 46,
  gap = 3,
}: {
  data: number[];
  color?: string;
  w?: number;
  h?: number;
  gap?: number;
}) {
  const max = Math.max(...data) || 1;
  const bw = (w - gap * (data.length - 1)) / data.length;
  return (
    <svg width={w} height={h} viewBox={`0 0 ${w} ${h}`} style={{ display: 'block' }} aria-hidden="true">
      {data.map((v, i) => {
        const bh = Math.max(2, (v / max) * (h - 2));
        return (
          <rect
            key={i}
            x={i * (bw + gap)}
            y={h - bh}
            width={bw}
            height={bh}
            rx={Math.min(2.5, bw / 2)}
            fill={color}
            opacity={i === data.length - 1 ? 1 : 0.32}
          />
        );
      })}
    </svg>
  );
}

/* ---- Ring / gauge score ---- */
export function Ring({
  value,
  size = 132,
  stroke = 11,
  color,
  track = 'var(--uk-line)',
  label,
  sub,
}: {
  value: number;
  size?: number;
  stroke?: number;
  color?: string;
  track?: string;
  label?: React.ReactNode;
  sub?: React.ReactNode;
}) {
  const r = (size - stroke) / 2;
  const c = 2 * Math.PI * r;
  const off = c * (1 - value / 100);
  const col =
    color ??
    (value >= 90
      ? 'var(--uk-success-bright)'
      : value >= 70
        ? 'var(--uk-orange)'
        : 'var(--uk-error)');
  return (
    <div style={{ position: 'relative', width: size, height: size }}>
      <svg width={size} height={size} style={{ transform: 'rotate(-90deg)' }} aria-hidden="true">
        <circle cx={size / 2} cy={size / 2} r={r} fill="none" stroke={track} strokeWidth={stroke} />
        <circle
          cx={size / 2}
          cy={size / 2}
          r={r}
          fill="none"
          stroke={col}
          strokeWidth={stroke}
          strokeLinecap="round"
          strokeDasharray={c}
          strokeDashoffset={off}
          style={{ transition: 'stroke-dashoffset .8s cubic-bezier(.3,1,.4,1)' }}
        />
      </svg>
      <div
        style={{
          position: 'absolute',
          inset: 0,
          display: 'flex',
          flexDirection: 'column',
          alignItems: 'center',
          justifyContent: 'center',
        }}
      >
        <div
          className="tnum"
          style={{
            fontFamily: 'var(--font-display)',
            fontSize: size * 0.27,
            fontWeight: 500,
            lineHeight: 1,
            color: 'var(--uk-ink)',
          }}
        >
          {label ?? value}
        </div>
        {sub && (
          <div style={{ fontSize: 11.5, color: 'var(--uk-ink-3)', marginTop: 3 }}>{sub}</div>
        )}
      </div>
    </div>
  );
}

/* ---- Donut (multi-segment) ---- */
export interface DonutSegment {
  value: number;
  color: string;
}
export function Donut({
  segments,
  size = 96,
  stroke = 14,
}: {
  segments: DonutSegment[];
  size?: number;
  stroke?: number;
}) {
  const r = (size - stroke) / 2;
  const c = 2 * Math.PI * r;
  const total = segments.reduce((s, x) => s + x.value, 0) || 1;
  const fracs = segments.map((s) => s.value / total);
  const offsets = fracs.map((_, i) =>
    fracs.slice(0, i).reduce((a, b) => a + b, 0),
  );
  return (
    <svg width={size} height={size} style={{ transform: 'rotate(-90deg)' }} aria-hidden="true">
      <circle cx={size / 2} cy={size / 2} r={r} fill="none" stroke="var(--uk-line)" strokeWidth={stroke} />
      {segments.map((s, i) => {
        const frac = fracs[i] ?? 0;
        const dash = `${(frac * c).toFixed(2)} ${c}`;
        return (
          <circle
            key={i}
            cx={size / 2}
            cy={size / 2}
            r={r}
            fill="none"
            stroke={s.color}
            strokeWidth={stroke}
            strokeDasharray={dash}
            strokeDashoffset={-(offsets[i] ?? 0) * c}
          />
        );
      })}
    </svg>
  );
}

/* ---- Clean gridded line chart (revenue trend, node health) ---- */
export function LineChart({
  data,
  height = 240,
  color = 'var(--uk-ac)',
}: {
  data: number[];
  height?: number;
  color?: string;
}) {
  const w = 900;
  const h = height;
  const padL = 8;
  const padR = 12;
  const padT = 16;
  const padB = 14;
  const max = Math.max(...data) * 1.08;
  const min = Math.min(...data) * 0.82;
  const span = max - min || 1;
  const X = (i: number) => padL + (i * (w - padL - padR)) / (data.length - 1);
  const Y = (v: number) => padT + (1 - (v - min) / span) * (h - padT - padB);
  const line = data
    .map((v, i) => `${i ? 'L' : 'M'}${X(i).toFixed(1)} ${Y(v).toFixed(1)}`)
    .join(' ');
  const grid = [0, 0.25, 0.5, 0.75, 1].map((f) => min + f * span);
  return (
    <svg
      viewBox={`0 0 ${w} ${h}`}
      width="100%"
      height={h}
      preserveAspectRatio="none"
      style={{ display: 'block', overflow: 'visible' }}
      aria-hidden="true"
    >
      {grid.map((g, i) => (
        <line key={i} x1={padL} x2={w - padR} y1={Y(g)} y2={Y(g)} stroke="var(--uk-line-soft)" strokeWidth="1" />
      ))}
      <path d={line} fill="none" stroke={color} strokeWidth="2.4" strokeLinecap="round" strokeLinejoin="round" />
      {data.map((v, i) => (
        <circle key={i} cx={X(i)} cy={Y(v)} r="3.4" fill={color} />
      ))}
    </svg>
  );
}

/* ---- Combo chart: translucent bars + line (site/switch overview) ---- */
export interface ComboPoint {
  bar: number;
  line: number;
}
export function ComboChart({
  data,
  height = 168,
  accent = 'var(--uk-ac)',
}: {
  data: ComboPoint[];
  height?: number;
  accent?: string;
}) {
  const w = 900;
  const h = height;
  const padT = 10;
  const padB = 16;
  const padX = 8;
  const max = Math.max(...data.map((d) => d.bar)) * 1.18 || 1;
  const n = data.length;
  const cx = (i: number) => padX + ((i + 0.5) * (w - padX * 2)) / n;
  const bw = ((w - padX * 2) / n) * 0.46;
  const yl = (v: number) => padT + (1 - v / max) * (h - padT - padB);
  const line = data
    .map((d, i) => `${i ? 'L' : 'M'}${cx(i).toFixed(1)} ${yl(d.line).toFixed(1)}`)
    .join(' ');
  return (
    <svg
      viewBox={`0 0 ${w} ${h}`}
      width="100%"
      height={h}
      preserveAspectRatio="none"
      style={{ display: 'block' }}
      aria-hidden="true"
    >
      {data.map((d, i) => {
        const bh = ((h - padT - padB) * d.bar) / max;
        return (
          <rect key={i} x={cx(i) - bw / 2} y={h - padB - bh} width={bw} height={bh} fill={accent} opacity="0.18" />
        );
      })}
      <path d={line} fill="none" stroke={accent} strokeWidth="2.2" strokeLinejoin="round" strokeLinecap="round" />
    </svg>
  );
}

/* ---- Mini area spark (switch port voltage/current/power) ---- */
export function MiniSpark({
  data,
  height = 44,
  accent = 'var(--uk-ac)',
}: {
  data: number[];
  height?: number;
  accent?: string;
}) {
  const w = 320;
  const h = height;
  const pad = 3;
  const max = Math.max(...data) * 1.1;
  const min = Math.min(...data) * 0.9;
  const span = max - min || 1;
  const X = (i: number) => pad + (i * (w - 2 * pad)) / (data.length - 1);
  const Y = (v: number) => pad + (1 - (v - min) / span) * (h - 2 * pad);
  const ln = data
    .map((v, i) => `${i ? 'L' : 'M'}${X(i).toFixed(1)} ${Y(v).toFixed(1)}`)
    .join(' ');
  const area = `${ln} L${X(data.length - 1).toFixed(1)} ${h} L${X(0).toFixed(1)} ${h} Z`;
  return (
    <svg
      viewBox={`0 0 ${w} ${h}`}
      width="100%"
      height={h}
      preserveAspectRatio="none"
      style={{ display: 'block' }}
      aria-hidden="true"
    >
      <path d={area} fill={accent} opacity="0.12" />
      <path d={ln} fill="none" stroke={accent} strokeWidth="2" />
    </svg>
  );
}

/* ---- series helper (fabricates a plausible series, from data.jsx) ---- */
export const series = (base: number, n = 14, jitter = 0.12, trend = 0): number[] =>
  Array.from({ length: n }, (_, i) =>
    Math.max(
      0,
      Math.round(
        (base +
          base * trend * (i / n) +
          base * jitter * (Math.sin(i * 1.7) + Math.cos(i * 0.9))) *
          10,
      ) / 10,
    ),
  );
