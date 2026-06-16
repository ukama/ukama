/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Parse a site's latitude/longitude into valid map coordinates. Some records
 * store the pair swapped (latitude holding the longitude), which puts the
 * marker off-world and blanks the map — if the "latitude" is outside the
 * valid range we swap. Returns null when either value is missing/non-numeric.
 */
export function normalizeCoords(
  latRaw?: string | number | null,
  lngRaw?: string | number | null,
): { lat: number; lng: number } | null {
  const a = typeof latRaw === 'number' ? latRaw : Number.parseFloat(latRaw ?? '');
  const b = typeof lngRaw === 'number' ? lngRaw : Number.parseFloat(lngRaw ?? '');
  if (!Number.isFinite(a) || !Number.isFinite(b)) return null;
  // Latitude must be within ±90; if it isn't but the other value is, swap.
  if (Math.abs(a) > 90 && Math.abs(b) <= 90) return { lat: b, lng: a };
  return { lat: a, lng: b };
}
