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
 * DR Congo region with always-on labelled status pins (BUILD-PLAN §7.1).
 */
import {
  ComposableMap,
  Geographies,
  Geography,
  Marker,
} from 'react-simple-maps';
import { DRC_GEO } from '@/data/reference/regions';
import { BIZ_DOT } from './siteMapShared';
import type { SiteMapSite } from './siteMapShared';

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
        style={{
          position: 'relative',
          margin: '14px 0 0',
          background: 'var(--uk-map-water)',
          ...(fill ? { flex: 1, minHeight: 300 } : { height }),
        }}
      >
        <ComposableMap
          projection="geoMercator"
          projectionConfig={{ center: [22.5, -3.5], scale: 1450 }}
          width={800}
          height={500}
          style={{ width: '100%', height: '100%' }}
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
                  <circle r={12} fill="none" stroke={col} strokeOpacity={0.3} strokeWidth={4} />
                )}
                <circle r={7} fill={col} stroke="#fff" strokeWidth={2.5} />
                <g>
                  <rect
                    x={12}
                    y={-11}
                    width={s.name.length * 6.6 + 18}
                    height={22}
                    rx={7}
                    fill="var(--uk-map-chip, #fff)"
                    stroke="var(--uk-map-border)"
                  />
                  <text
                    x={21}
                    y={4}
                    style={{
                      fontFamily: 'var(--font-body)',
                      fontSize: 12,
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
        </ComposableMap>
      </div>
    </div>
  );
}
