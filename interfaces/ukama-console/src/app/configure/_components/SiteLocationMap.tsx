/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Small location preview for the site name step: a pin at the tower's
 * coordinates on the shared console base map. Static — no pan/zoom.
 */
'use client';

import UkamaMap from '@/components/Map/UkamaMap';

export default function SiteLocationMap({
  lat,
  lng,
  height = 150,
}: {
  lat: number;
  lng: number;
  height?: number;
}) {
  return (
    <div
      style={{
        height,
        borderRadius: 'var(--uk-r-md)',
        overflow: 'hidden',
        border: '1px solid var(--uk-line)',
      }}
    >
      <UkamaMap
        markers={[{ id: 'site', lat, lng }]}
        zoom={16}
        height="100%"
        interactive={false}
        fitToMarkers={false}
      />
    </div>
  );
}
