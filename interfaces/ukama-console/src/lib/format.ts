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
