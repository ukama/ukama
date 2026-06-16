/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/** Lazy, client-only Leaflet map (maps are code-split; never SSR'd). */
import dynamic from 'next/dynamic';
import Skeleton from '@mui/material/Skeleton';

export type { UkamaMapMarker, UkamaMapProps } from './UkamaMapImpl';

/** Default zoom for the home dashboards' site map. */
export const HOME_MAP_ZOOM = 13;

const UkamaMap = dynamic(() => import('./UkamaMapImpl'), {
  ssr: false,
  loading: () => (
    <Skeleton variant="rounded" sx={{ width: '100%', height: '100%', minHeight: 160 }} />
  ),
});

export default UkamaMap;
