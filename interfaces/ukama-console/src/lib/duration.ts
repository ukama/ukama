/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Package duration helpers.
 *
 * Since #1496 the backend stores and interprets package `duration` in
 * MINUTES (previously days). The console presents validity dynamically
 * (minutes / hours / days / weeks / months), so these helpers convert and
 * format at the GraphQL/API boundary. Unit sizes mirror the backend
 * (systems/common/validation/date.go: a month = 30 days).
 */
export const MINUTES_PER_HOUR = 60;
export const MINUTES_PER_DAY = 1440;
export const MINUTES_PER_WEEK = 7 * MINUTES_PER_DAY; // 10080
export const MINUTES_PER_MONTH = 30 * MINUTES_PER_DAY; // 43200 — matches backend MinutesInMonth

/** Whole-day validity -> minutes the backend expects. */
export const daysToMinutes = (days: number): number =>
  Math.round(days) * MINUTES_PER_DAY;

/** Backend minute duration -> whole days (used by the day-based create form). */
export const minutesToDays = (minutes: number): number =>
  Math.round(minutes / MINUTES_PER_DAY);

/** Largest-to-smallest units used to render a duration. */
const DURATION_UNITS: ReadonlyArray<{ unit: string; minutes: number }> = [
  { unit: 'month', minutes: MINUTES_PER_MONTH },
  { unit: 'week', minutes: MINUTES_PER_WEEK },
  { unit: 'day', minutes: MINUTES_PER_DAY },
  { unit: 'hour', minutes: MINUTES_PER_HOUR },
  { unit: 'minute', minutes: 1 },
];

/**
 * Pick the largest unit that divides the duration evenly, so canonical plan
 * lengths render cleanly (43200 -> 1 month, 10080 -> 1 week, 1440 -> 1 day,
 * 120 -> 2 hours) and anything else falls back to minutes (90 -> 90 minutes).
 */
const pickDurationUnit = (
  minutes: number,
): { unit: string; count: number } => {
  const m = Math.max(0, Math.round(minutes));
  for (const u of DURATION_UNITS) {
    if (m >= u.minutes && m % u.minutes === 0) {
      return { unit: u.unit, count: m / u.minutes };
    }
  }
  return { unit: 'minute', count: m };
};

/** Human-readable, pluralized duration — e.g. "1 month", "2 weeks", "45 minutes". */
export const formatDuration = (minutes: number): string => {
  const { unit, count } = pickDurationUnit(minutes);
  return `${count} ${unit}${count === 1 ? '' : 's'}`;
};

/** Singular unit noun for a per-period price suffix — e.g. "month" in "$5 / month". */
export const durationUnitLabel = (minutes: number): string =>
  pickDurationUnit(minutes).unit;
