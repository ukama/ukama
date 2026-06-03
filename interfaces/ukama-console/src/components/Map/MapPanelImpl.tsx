/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/**
 * Network map (ops lens) — react-simple-maps over the Zambia region
 * (BUILD-PLAN §7.1), with the design's status pins, labels for selected /
 * problem sites, zoom chrome and legend. Coordinates are dummy lat/lng
 * until the API serves real geo.
 */
import { useMemo, useState } from 'react';
import AddRounded from '@mui/icons-material/AddRounded';
import RemoveRounded from '@mui/icons-material/RemoveRounded';
import {
  ComposableMap,
  Geographies,
  Geography,
  Marker,
  ZoomableGroup,
} from 'react-simple-maps';
import type { Site } from '@/data';
import { ZAMBIA_GEO } from '@/data/reference/regions';

const PIN_COLOR: Record<Site['status'], string> = {
  online: 'var(--uk-success-bright)',
  degraded: 'var(--uk-warning)',
  offline: 'var(--uk-error)',
};

const GEO_STYLE = {
  default: {
    fill: 'var(--uk-map-land)',
    stroke: 'var(--uk-map-border)',
    strokeWidth: 0.6,
    outline: 'none',
  },
  hover: { fill: 'var(--uk-map-land)', outline: 'none' },
  pressed: { fill: 'var(--uk-map-land)', outline: 'none' },
};

export default function MapPanelImpl({
  sites,
  selected,
  onSelect,
  compact,
}: {
  sites: Site[];
  selected?: string | null;
  onSelect?: (site: Site) => void;
  compact?: boolean;
}) {
  const [zoom, setZoom] = useState(compact ? 5 : 1);

  const center = useMemo<[number, number]>(() => {
    if (sites.length === 0) return [27.8, -14.8];
    if (sites.length === 1 && sites[0]) return [sites[0].lng, sites[0].lat];
    return [27.8, -14.8]; // Zambia
  }, [sites]);

  return (
    <div
      style={{
        position: 'relative',
        width: '100%',
        height: '100%',
        borderRadius: 'var(--uk-r-md)',
        overflow: 'hidden',
        background: 'var(--uk-map-water)',
      }}
    >
      <ComposableMap
        projection="geoMercator"
        projectionConfig={{ center, scale: 2300 }}
        width={800}
        height={520}
        style={{ width: '100%', height: '100%' }}
      >
        <ZoomableGroup
          center={center}
          zoom={zoom}
          minZoom={0.8}
          maxZoom={12}
          onMoveEnd={({ zoom: z }) => setZoom(z)}
        >
          <Geographies geography={ZAMBIA_GEO}>
            {({ geographies }) =>
              geographies.map((geo) => (
                <Geography key={geo.rsmKey} geography={geo} style={GEO_STYLE} tabIndex={-1} />
              ))
            }
          </Geographies>
          {sites.map((s) => {
            const col = PIN_COLOR[s.status];
            const sel = selected === s.id;
            const showLabel = !compact && (sel || s.status !== 'online');
            const r = (sel ? 8 : 6.5) / Math.sqrt(zoom);
            return (
              <Marker
                key={s.id}
                coordinates={[s.lng, s.lat]}
                onClick={() => onSelect?.(s)}
                style={{ default: { cursor: 'pointer' } }}
              >
                {sel && (
                  <circle
                    r={r + 5 / Math.sqrt(zoom)}
                    fill="none"
                    stroke={col}
                    strokeOpacity={0.35}
                    strokeWidth={4 / Math.sqrt(zoom)}
                  />
                )}
                <circle
                  r={r}
                  fill={col}
                  stroke="#fff"
                  strokeWidth={2.5 / Math.sqrt(zoom)}
                />
                {showLabel && (
                  <g transform={`scale(${1 / Math.sqrt(zoom)})`}>
                    <rect
                      x={-(s.name.length * 3.4 + 8)}
                      y={-32}
                      width={s.name.length * 6.8 + 16}
                      height={20}
                      rx={6}
                      fill="#15181F"
                    />
                    <text
                      textAnchor="middle"
                      y={-18}
                      style={{
                        fontFamily: 'var(--font-body)',
                        fontSize: 11.5,
                        fontWeight: 600,
                        fill: '#fff',
                      }}
                    >
                      {s.name}
                    </text>
                  </g>
                )}
              </Marker>
            );
          })}
        </ZoomableGroup>
      </ComposableMap>

      {!compact && (
        <>
          <div
            style={{
              position: 'absolute',
              top: 14,
              right: 14,
              display: 'flex',
              flexDirection: 'column',
              borderRadius: 8,
              overflow: 'hidden',
              boxShadow: 'var(--uk-shadow)',
              background: '#fff',
            }}
          >
            <button
              type="button"
              className="map-zoom"
              aria-label="Zoom in"
              onClick={() => setZoom((z) => Math.min(12, z * 1.6))}
            >
              <AddRounded sx={{ fontSize: 18 }} />
            </button>
            <hr className="divider" />
            <button
              type="button"
              className="map-zoom"
              aria-label="Zoom out"
              onClick={() => setZoom((z) => Math.max(0.8, z / 1.6))}
            >
              <RemoveRounded sx={{ fontSize: 18 }} />
            </button>
          </div>
          <div
            style={{
              position: 'absolute',
              left: 14,
              bottom: 14,
              display: 'flex',
              gap: 14,
              background: 'rgba(255,255,255,.92)',
              padding: '6px 12px',
              borderRadius: 8,
              boxShadow: 'var(--uk-shadow)',
              fontSize: 11.5,
              color: '#5a5e66',
            }}
          >
            {(
              [
                ['Online', 'var(--uk-success-bright)'],
                ['Degraded', 'var(--uk-warning)'],
                ['Offline', 'var(--uk-error)'],
              ] as const
            ).map(([l, c]) => (
              <span key={l} style={{ display: 'flex', alignItems: 'center', gap: 6 }}>
                <span style={{ width: 8, height: 8, borderRadius: '50%', background: c }} />
                {l}
              </span>
            ))}
          </div>
        </>
      )}
    </div>
  );
}
