/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Resolves a human-readable address for coordinates via OpenStreetMap's
 * Nominatim service (no API key). Falls back to the formatted coordinates if
 * the lookup fails, so the step never blocks on it.
 */
'use client';

import { useEffect, useState } from 'react';

const coordLabel = (lat: number, lng: number): string =>
  `${lat.toFixed(5)}, ${lng.toFixed(5)}`;

export function useReverseGeocode(
  lat: number | null,
  lng: number | null,
): { address: string } {
  const [address, setAddress] = useState('');

  useEffect(() => {
    if (lat === null || lng === null) return;
    const fallback = coordLabel(lat, lng);
    const controller = new AbortController();
    let active = true;
    fetch(
      `https://nominatim.openstreetmap.org/reverse?format=json&lat=${lat}&lon=${lng}`,
      { signal: controller.signal, headers: { Accept: 'application/json' } },
    )
      .then((res) => (res.ok ? res.json() : null))
      .then((data: { display_name?: string } | null) => {
        if (active) setAddress(data?.display_name?.trim() || fallback);
      })
      .catch(() => {
        if (active) setAddress(fallback);
      });

    return () => {
      active = false;
      controller.abort();
    };
  }, [lat, lng]);

  return { address };
}
