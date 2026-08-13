/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/** Map view preference — the base layer used by every console map. */
import UkamaMap from './Map/UkamaMap';
import { MAP_VIEWS } from './Map/basemaps';
import { useUiPrefs } from '@/lib/store';

/** Neutral preview location — enough terrain/roads to tell the layers apart. */
const PREVIEW_CENTER: [number, number] = [-15.3875, 28.3228];

export default function MapViewSetting() {
  const mapView = useUiPrefs((s) => s.mapView);
  const setMapView = useUiPrefs((s) => s.setMapView);

  return (
    <div className="card card-pad" style={{ gridColumn: '1 / -1' }}>
      <label className="flabel">Map view</label>
      <div style={{ fontSize: 13, color: 'var(--uk-ink-2)', marginBottom: 12 }}>
        Base layer used by every map in the console.
      </div>

      <div
        role="radiogroup"
        aria-label="Map view"
        style={{
          display: 'grid',
          gridTemplateColumns: 'repeat(auto-fit, minmax(180px, 1fr))',
          gap: 10,
        }}
      >
        {MAP_VIEWS.map((v) => {
          const active = v.id === mapView;
          return (
            <div
              key={v.id}
              role="radio"
              aria-checked={active}
              tabIndex={active ? 0 : -1}
              className={`card card-pad card-selectable${active ? ' is-active' : ''}`}
              onClick={() => setMapView(v.id)}
              onKeyDown={(e) => {
                if (e.key === 'Enter' || e.key === ' ') {
                  e.preventDefault();
                  setMapView(v.id);
                }
              }}
            >
              <div style={{ fontWeight: 600, fontSize: 14 }}>{v.label}</div>
              <div
                style={{ fontSize: 12, color: 'var(--uk-ink-2)', marginTop: 4 }}
              >
                {v.hint}
              </div>
            </div>
          );
        })}
      </div>

      <div
        style={{
          marginTop: 14,
          height: 180,
          borderRadius: 'var(--uk-r-md)',
          overflow: 'hidden',
          border: '1px solid var(--uk-line)',
        }}
      >
        <UkamaMap
          center={PREVIEW_CENTER}
          zoom={12}
          height="100%"
          interactive={false}
          fitToMarkers={false}
        />
      </div>
    </div>
  );
}
