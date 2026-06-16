/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/** Turn a snake_case metric key into Title Case
 *  (e.g. subscribers_active → "Subscribers Active"). */
export const humanizeKey = (key: string): string =>
  key
    .split('_')
    .filter(Boolean)
    .map((w) => w.charAt(0).toUpperCase() + w.slice(1))
    .join(' ');

/** True when a string is a raw snake_case key rather than a display label. */
const isRawKey = (s?: string | null): boolean =>
  !!s && /^[a-z0-9]+(_[a-z0-9]+)+$/.test(s);

/**
 * Best display title for a metric: a clean server label if present, else the
 * provided fallback, else a humanized key. Never returns a raw snake_case key
 * (e.g. an enrich() fallback where label === the metric key).
 */
export const metricLabel = (
  label: string | null | undefined,
  key: string,
  fallback?: string,
): string => {
  if (label && !isRawKey(label)) return label;
  return fallback || humanizeKey(key);
};
