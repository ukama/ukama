/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Base map layers shared by every console map. The selected view lives in the
 * UI prefs store (`useUiPrefs().mapView`).
 *
 * Street (CARTO) follows the app theme and needs an API key
 * (NEXT_PUBLIC_CARTO_API_KEY); satellite/terrain are keyless and naturally
 * dark in both modes.
 */

import { publicEnv } from '@/lib/runtime-env';

export type MapView = 'street' | 'satellite' | 'terrain';

export interface Basemap {
  id: MapView;
  /** Layer name shown in the Settings picker. */
  label: string;
  /** One-line description for the Settings picker. */
  hint: string;
  url: string;
  /** Optional dark-mode variant; falls back to `url`. */
  darkUrl?: string;
  /** Appended as `?key=` when NEXT_PUBLIC_CARTO_API_KEY is set. */
  needsKey?: boolean;
}

export const MAP_VIEWS: readonly Basemap[] = [
  {
    id: 'street',
    label: 'Street',
    hint: 'Roads and place names. Follows light/dark theme.',
    url: 'https://{s}.basemaps.cartocdn.com/light_all/{z}/{x}/{y}.png',
    darkUrl: 'https://{s}.basemaps.cartocdn.com/dark_all/{z}/{x}/{y}.png',
    needsKey: true,
  },
  {
    id: 'satellite',
    label: 'Satellite',
    hint: 'Aerial imagery — useful for checking site surroundings.',
    url: 'https://server.arcgisonline.com/ArcGIS/rest/services/World_Imagery/MapServer/tile/{z}/{y}/{x}',
  },
  {
    id: 'terrain',
    label: 'Terrain',
    hint: 'Elevation and contours — useful for coverage planning.',
    url: 'https://{s}.tile.opentopomap.org/{z}/{x}/{y}.png',
  },
] as const;

export const DEFAULT_MAP_VIEW: MapView = 'terrain';

/** Tile URL for a view, honouring dark mode where the source has a variant. */
export function basemapUrl(view: MapView, dark: boolean): string {
  const base = MAP_VIEWS.find((v) => v.id === view) ?? MAP_VIEWS[0]!;
  const url = (dark && base.darkUrl) || base.url;
  const key = publicEnv().cartoApiKey;
  return base.needsKey && key ? `${url}?key=${encodeURIComponent(key)}` : url;
}
