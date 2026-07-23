/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Data-usage parsing/formatting. The backend always reports data as a raw byte
 * count (as a string); the UI stores bytes and formats on render.
 *
 * BINARY units (1 KB = 1024 B) to match Ukama's data-unit convention
 * (systems/common/ukama), where package allowances and usage are defined in
 * binary — a "1 GB" package is 1024 MB = 1,073,741,824 bytes. Formatting in
 * decimal here would render that same 1 GB package as "1.07 GB".
 */

const KB = 1024;
const MB = KB * 1024;
const GB = MB * 1024;
const TB = GB * 1024;

/** Sentinel stored in Subscriber.usage when usage is unknown/not loaded. */
export const USAGE_UNKNOWN = -1;

/**
 * Parse a raw byte value (string or number) to a non-negative number of
 * bytes. Returns null for missing/invalid input so callers can fall back to
 * the "unknown" sentinel.
 */
export function parseUsageBytes(
  value: string | number | null | undefined,
): number | null {
  if (value == null || value === '') return null;
  const n = typeof value === 'number' ? value : Number(value);
  return Number.isFinite(n) && n >= 0 ? n : null;
}

/** Bytes → GB (binary), for ratios against a plan cap expressed in GB. */
export function bytesToGB(bytes: number): number {
  return bytes / GB;
}

/**
 * Human-readable bytes in the most fitting unit, e.g. "0 B", "512 MB",
 * "2.34 GB". Trailing zeros are trimmed (2.00 → "2").
 */
export function formatBytes(bytes: number, decimals = 2): string {
  if (!Number.isFinite(bytes) || bytes <= 0) return '0 B';
  const units: [number, string][] = [
    [TB, 'TB'],
    [GB, 'GB'],
    [MB, 'MB'],
    [KB, 'KB'],
  ];
  for (const [factor, unit] of units) {
    if (bytes >= factor) {
      const value = Math.round((bytes / factor) * 10 ** decimals) / 10 ** decimals;
      return `${value} ${unit}`;
    }
  }
  return `${Math.round(bytes)} B`;
}
