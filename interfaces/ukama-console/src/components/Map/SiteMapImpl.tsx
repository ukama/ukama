/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/**
 * Coverage / sites map (Business lens) — react-simple-maps over the
 * DR Congo region with always-on labelled status pins, zoom/pan and
 * zoom chrome (BUILD-PLAN §7.1).
 */
import { useState } from 'react';
import AddRounded from '@mui/icons-material/AddRounded';
import RemoveRounded from '@mui/icons-material/RemoveRounded';
import {
  ComposableMap,
  Geographies,
  Geography,
  Marker,
  ZoomableGroup,
} from 'react-simple-maps';
import { DRC_GEO } from '@/data/reference/regions';
import { BIZ_DOT } from './siteMapShared';
import type { SiteMapSite } from './siteMapShared';
import { useElementSize } from './useElementSize';

const CENTER: [number, number] = [22.5, -3.5];

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

export default function SiteMapImpl({
  sites,
  title,
  height = 380,
  fill,
  action,
  selected,
  onSelect,
}: {
  sites: SiteMapSite[];
  title?: string;
  height?: number;
  fill?: boolean;
  action?: React.ReactNode;
  selected?: string | null;
  onSelect?: (site: SiteMapSite) => void;
}) {
  const [zoom, setZoom] = useState(1);
  const [sizeRef, { width: w, height: h }] = useElementSize();
  const k = 1 / Math.sqrt(zoom);

  return (
    <div
      className="card"
      style={{
        padding: 0,
        overflow: 'hidden',
        ...(fill ? { display: 'flex', flexDirection: 'column', height: '100%' } : {}),
      }}
    >
      {(title || action) && (
        <div
          style={{
            padding: '18px 22px 0',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'space-between',
            gap: 12,
          }}
        >
          {title ? <div className="sec-title">{title}</div> : <span />}
          {action}
        </div>
      )}
      <div
        ref={sizeRef}
        style={{
          position: 'relative',
          margin: '14px 0 0',
          background: 'var(--uk-map-water)',
          ...(fill ? { flex: 1, minHeight: 300 } : { height }),
        }}
      >
        {w > 0 && h > 0 && (
        <ComposableMap
          projection="geoMercator"
          projectionConfig={{ center: CENTER, scale: Math.min(w * 1.6, h * 2.8) }}
          width={w}
          height={h}
          style={{ width: '100%', height: '100%', display: 'block' }}
        >
          <ZoomableGroup
            center={CENTER}
            zoom={zoom}
            minZoom={0.8}
            maxZoom={12}
            onMoveEnd={({ zoom: z }) => setZoom(z)}
            // zoom only via the overlay buttons — no wheel / pinch / dblclick
            filterZoomEvent={(e) => e?.type !== 'wheel' && e?.type !== 'dblclick'}
          >
            <Geographies geography={DRC_GEO}>
              {({ geographies }) =>
                geographies.map((geo) => (
                  <Geography key={geo.rsmKey} geography={geo} style={GEO_STYLE} tabIndex={-1} />
                ))
              }
            </Geographies>
            {sites.map((s) => {
              const sel = selected === s.id;
              const col = BIZ_DOT[s.status];
              return (
                <Marker
                  key={s.id}
                  coordinates={[s.lng, s.lat]}
                  onClick={() => onSelect?.(s)}
                  style={{ default: { cursor: onSelect ? 'pointer' : 'default' } }}
                >
                  {sel && (
                    <circle
                      r={9 * k}
                      fill="none"
                      stroke={col}
                      strokeOpacity={0.3}
                      strokeWidth={3 * k}
                    />
                  )}
                  <circle r={5.5 * k} fill={col} stroke="#fff" strokeWidth={1.8 * k} />
                  <g transform={`scale(${k})`}>
                    <rect
                      x={10}
                      y={-9}
                      width={s.name.length * 6 + 16}
                      height={18}
                      rx={6}
                      fill="var(--uk-map-chip, #fff)"
                      stroke="var(--uk-map-border)"
                    />
                    <text
                      x={18}
                      y={4}
                      style={{
                        fontFamily: 'var(--font-body)',
                        fontSize: 11,
                        fontWeight: 500,
                        fill: 'var(--uk-ink)',
                      }}
                    >
                      {s.name}
                    </text>
                  </g>
                </Marker>
              );
            })}
          </ZoomableGroup>
        </ComposableMap>
        )}

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
            background: 'var(--uk-panel)',
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
      </div>
    </div>
  );
}
