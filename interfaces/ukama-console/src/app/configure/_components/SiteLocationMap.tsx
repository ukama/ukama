/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Small location preview for the site name step: an OpenStreetMap embed with
 * a marker at the tower's coordinates. Dependency-free (no Leaflet/Mapbox
 * token) so it stays light for a one-off onboarding view.
 */
'use client';

export default function SiteLocationMap({
  lat,
  lng,
  height = 150,
}: {
  lat: number;
  lng: number;
  height?: number;
}) {
  // A tight bbox zooms in on the marker (smaller = closer).
  const d = 0.004;
  const bbox = `${lng - d}%2C${lat - d}%2C${lng + d}%2C${lat + d}`;
  const src =
    `https://www.openstreetmap.org/export/embed.html?bbox=${bbox}` +
    `&layer=mapnik&marker=${lat}%2C${lng}`;

  return (
    <div
      style={{
        height,
        borderRadius: 'var(--uk-r-md)',
        overflow: 'hidden',
        border: '1px solid var(--uk-line)',
      }}
    >
      <iframe
        title="Site location"
        src={src}
        loading="lazy"
        style={{ width: '100%', height: '100%', border: 0 }}
      />
    </div>
  );
}
