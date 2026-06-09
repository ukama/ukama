/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Deterministic metric mocks. Output is byte-shape-identical to the metric
 * service responses (latest = [ts, val]; range = [[ts, val], ...]) so the
 * console can't tell mock from live, and going live needs no console change.
 *
 * Determinism: a metric's series is seeded by key+scope+time-bucket, so it's
 * stable across polls within a bucket but differs per metric and per node.
 */
import { metricMeta } from "./catalog";

const hash = (s: string): number => {
  let h = 2166136261;
  for (let i = 0; i < s.length; i++) {
    h ^= s.charCodeAt(i);
    h = Math.imul(h, 16777619);
  }
  return h >>> 0;
};

/** mulberry32 PRNG — small, fast, deterministic. */
const prng = (seed: number) => () => {
  seed |= 0;
  seed = (seed + 0x6d2b79f5) | 0;
  let t = Math.imul(seed ^ (seed >>> 15), 1 | seed);
  t = (t + Math.imul(t ^ (t >>> 7), 61 | t)) ^ t;
  return ((t ^ (t >>> 14)) >>> 0) / 4294967296;
};

const clamp = (v: number, lo: number, hi: number) =>
  Math.min(hi, Math.max(lo, v));

const round = (v: number, format: string) =>
  format === "decimal" ? Math.round(v * 100) / 100 : Math.round(v);

/** One mocked value for a metric at point index `i` of `n`. */
const valueAt = (
  key: string,
  rand: () => number,
  i: number,
  n: number
): number => {
  const m = metricMeta(key);
  const drift = m.base * m.trend * (i / Math.max(1, n));
  const wobble =
    m.base * m.jitter * (Math.sin(i * 1.7) + Math.cos(i * 0.9)) * 0.5;
  const noise = m.base * m.jitter * (rand() - 0.5);
  return round(clamp(m.base + drift + wobble + noise, m.min, m.max), m.format);
};

/** Latest single value: [timestampSec, value]. */
export const mockLatest = (
  key: string,
  scopeId: string
): { value: [number, number]; success: boolean } => {
  const bucket = Math.floor(Date.now() / 60000); // 1-min stability bucket
  const rand = prng(hash(`${key}:${scopeId}:${bucket}`));
  const ts = Math.floor(Date.now() / 1000);
  // Draw a short walk and take the last point so it lines up with the range.
  let v = 0;
  for (let i = 0; i < 8; i++) v = valueAt(key, rand, i, 8);
  return { value: [ts, v], success: true };
};

/** Range series: [[tsSec, value], ...] ascending, `step`-spaced. */
export const mockRangeValues = (
  key: string,
  scopeId: string,
  from: number,
  to: number,
  step: number
): [number, number][] => {
  const safeStep = step > 0 ? step : 300;
  const n = Math.min(500, Math.max(2, Math.floor((to - from) / safeStep) + 1));
  const bucket = Math.floor(from / safeStep);
  const rand = prng(hash(`${key}:${scopeId}:${bucket}`));
  return Array.from({ length: n }, (_, i) => {
    const ts = from + i * safeStep;
    return [ts, valueAt(key, rand, i, n)] as [number, number];
  });
};
