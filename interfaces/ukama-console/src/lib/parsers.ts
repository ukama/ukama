/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/** Relative-time + data-size parsers for sortable columns (table-kit.jsx). */

/** "2m ago" → minutes; "just now" → 0; unknown/never → Infinity. */
export function parseSeen(s: string | undefined): number {
  if (!s) return Infinity;
  const v = String(s).toLowerCase();
  if (v.includes('just now')) return 0;
  if (v.includes('never')) return Infinity;
  const m = v.match(/(\d+)\s*(m|min|h|hour|d|day)/);
  if (!m || !m[1] || !m[2]) return Infinity;
  const n = +m[1];
  const u = m[2];
  return u[0] === 'm' ? n : u[0] === 'h' ? n * 60 : u[0] === 'd' ? n * 1440 : Infinity;
}

/**
 * Parses a backend timestamp to epoch ms, or NaN if unparseable. Handles the
 * Go format the gateway emits, e.g. "2026-06-06 23:04:53.705589 +0000 UTC",
 * as well as ISO strings.
 */
export function parseTimestamp(s: string | undefined): number {
  if (!s) return NaN;
  const direct = Date.parse(s);
  if (!Number.isNaN(direct)) return direct;
  // Normalize Go's "YYYY-MM-DD HH:mm:ss.ffffff +0000 UTC" to ISO.
  const normalized = s
    .replace(' ', 'T')
    .replace(/(\.\d{3})\d+/, '$1') // trim sub-millisecond precision
    .replace(/\s*\+0000 UTC$/, 'Z')
    .replace(/\s*UTC$/, 'Z')
    .replace(/\s*([+-]\d{2})(\d{2})$/, '$1:$2'); // +0000 → +00:00
  return Date.parse(normalized);
}

/** Human-friendly date for table cells, e.g. "Jun 6, 2026". '—' when absent. */
export function formatDate(s: string | undefined): string {
  const ms = parseTimestamp(s);
  if (Number.isNaN(ms)) return '—';
  return new Date(ms).toLocaleDateString('en-US', {
    day: 'numeric',
    month: 'short',
    year: 'numeric',
  });
}

/** "612 GB" / "1.8 TB" / "—" → GB float for sorting. */
export function parseData(d: string | undefined): number {
  if (!d || d === '—') return -1;
  const m = String(d).match(/([\d.]+)\s*(GB|MB|TB)/i);
  if (!m || !m[1] || !m[2]) return 0;
  let v = parseFloat(m[1]);
  const u = m[2].toUpperCase();
  if (u === 'MB') v /= 1024;
  if (u === 'TB') v *= 1024;
  return v;
}
