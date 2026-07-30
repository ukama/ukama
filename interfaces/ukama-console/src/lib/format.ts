/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/** ISO-3166 alpha-2 code → English country name; passes other strings through. */
export const countryLabel = (country: string): string => {
  if (/^[A-Z]{2}$/.test(country)) {
    try {
      return (
        new Intl.DisplayNames(['en'], { type: 'region' }).of(country) ?? country
      );
    } catch {
      return country;
    }
  }
  return country;
};

/** Up to two uppercase initials from a name (falling back to an email), else '?'. */
export const initials = (name?: string, email?: string): string => {
  const base = (name?.trim() || email || '').trim() || '?';
  return (
    base
      .split(/\s+/)
      .filter(Boolean)
      .slice(0, 2)
      .map((part) => part[0]?.toUpperCase() ?? '')
      .join('') || '?'
  );
};

/**
 * Human label for the analytics report-window token returned by
 * getPerformanceReport (e.g. "8w" -> "Last 8 weeks", "3d" -> "Last 3 days").
 * The report uses a config-driven rolling window independent of the UI filter,
 * so screens show this to make the window explicit. "" for empty/unknown.
 */
export const reportWindowLabel = (span?: string | null): string => {
  if (!span) return '';
  const m = /^(\d+)([wd])$/.exec(span);
  if (!m) return '';
  const n = Number(m[1]);
  const unit = m[2] === 'w' ? 'week' : 'day';
  return `Last ${n} ${unit}${n === 1 ? '' : 's'}`;
};
