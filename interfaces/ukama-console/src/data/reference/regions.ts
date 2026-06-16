/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Simplified operating-region outlines (BUILD-PLAN §7.1) — intentionally
 * low-fidelity "dummy" polygons good enough for the console's regional
 * maps. Replace with real GeoJSON when the API serves geographies.
 */

interface RegionFeature {
  type: 'Feature';
  properties: { name: string };
  geometry: { type: 'Polygon'; coordinates: [number, number][][] };
}

export interface RegionGeo {
  type: 'FeatureCollection';
  features: RegionFeature[];
}

/** Zambia — rough national outline (lng, lat). */
export const ZAMBIA_GEO: RegionGeo = {
  type: 'FeatureCollection',
  features: [
    {
      type: 'Feature',
      properties: { name: 'Zambia' },
      geometry: {
        type: 'Polygon',
        coordinates: [
          [
            [22.0, -16.5],
            [22.0, -13.1],
            [24.0, -12.9],
            [24.0, -11.3],
            [25.5, -11.2],
            [27.2, -12.0],
            [28.5, -12.8],
            [29.0, -13.5],
            [29.6, -13.2],
            [30.8, -13.5],
            [32.0, -13.6],
            [33.2, -13.9],
            [33.7, -14.5],
            [32.9, -15.4],
            [31.0, -16.0],
            [30.4, -16.0],
            [28.9, -16.0],
            [27.0, -17.0],
            [25.3, -17.8],
            [23.4, -17.6],
            [22.0, -16.5],
          ],
        ],
      },
    },
  ],
};
