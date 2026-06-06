/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Parses a node's latitude/longitude strings into valid map coordinates.
 * Some upstreams report the pair swapped, so we orient by range: latitude
 * must be within ±90, longitude within ±180. Returns null when neither
 * ordering is valid (e.g. unset/zeroed coordinates).
 */
export interface LatLng {
  lat: number;
  lng: number;
}

const inLat = (v: number): boolean => Math.abs(v) <= 90;
const inLng = (v: number): boolean => Math.abs(v) <= 180;

export function parseCoords(
  latStr?: string | null,
  lngStr?: string | null,
): LatLng | null {
  const a = Number(latStr);
  const b = Number(lngStr);
  if (!Number.isFinite(a) || !Number.isFinite(b)) return null;
  if (inLat(a) && inLng(b)) return { lat: a, lng: b };
  // Stored the other way round — swap so the latitude is the in-range value.
  if (inLat(b) && inLng(a)) return { lat: b, lng: a };
  return null;
}
