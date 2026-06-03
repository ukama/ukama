/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/**
 * Stylized network map (Network lens) — ported from the prototype's
 * MapPanel (shell.jsx): abstract terrain SVG, status pins, zoom chrome
 * and legend. Replaces Leaflet per BUILD-PLAN §7.1 (no tiles, no tokens).
 */
import AddRounded from '@mui/icons-material/AddRounded';
import RemoveRounded from '@mui/icons-material/RemoveRounded';
import type { Site } from '@/data';

const PIN_COLOR: Record<Site['status'], string> = {
  online: 'var(--uk-success-bright)',
  degraded: 'var(--uk-warning)',
  offline: 'var(--uk-error)',
};

export default function MapPanel({
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
  return (
    <div
      style={{
        position: 'relative',
        width: '100%',
        height: '100%',
        borderRadius: 'var(--uk-r-md)',
        overflow: 'hidden',
        background: 'linear-gradient(160deg,#eef2f6,#e6ebf0)',
      }}
    >
      <svg
        width="100%"
        height="100%"
        style={{ position: 'absolute', inset: 0 }}
        preserveAspectRatio="none"
        viewBox="0 0 100 100"
        aria-hidden="true"
      >
        <rect width="100" height="100" fill="#e9eef3" />
        <path d="M0,44 Q30,38 52,52 T100,48" stroke="#d3dde6" strokeWidth="1" fill="none" />
        <path d="M18,0 Q24,40 44,62 T56,100" stroke="#d9e2ea" strokeWidth="0.7" fill="none" />
        <path d="M0,74 Q40,66 70,80 T100,76" stroke="#d9e2ea" strokeWidth="0.7" fill="none" />
        <path d="M62,8 L100,30 L100,0 Z" fill="#dfe9e2" opacity="0.55" />
        <circle cx="31" cy="34" r="15" fill="#dce8ef" opacity="0.5" />
        <circle cx="52" cy="20" r="11" fill="#dce8ef" opacity="0.45" />
        <path d="M52,18 L70,26 L62,46 L44,40 Z" fill="#e2ebe6" opacity="0.5" />
      </svg>
      {sites.map((s) => {
        const col = PIN_COLOR[s.status];
        const sel = selected === s.id;
        return (
          <button
            key={s.id}
            type="button"
            onClick={() => onSelect?.(s)}
            title={s.name}
            style={{
              position: 'absolute',
              left: `${s.x}%`,
              top: `${s.y}%`,
              transform: 'translate(-50%,-50%)',
              background: 'none',
              border: 'none',
              cursor: 'pointer',
              padding: 0,
              zIndex: sel ? 5 : 2,
            }}
          >
            <span
              style={{
                display: 'block',
                width: sel ? 18 : 14,
                height: sel ? 18 : 14,
                borderRadius: '50%',
                background: col,
                border: '3px solid #fff',
                boxShadow: sel
                  ? `0 0 0 4px color-mix(in srgb, ${col} 33%, transparent), 0 2px 6px rgba(0,0,0,.25)`
                  : '0 2px 5px rgba(0,0,0,.25)',
                transition: 'all .15s',
              }}
            />
            {!compact && (sel || s.status !== 'online') && (
              <span
                style={{
                  position: 'absolute',
                  top: '-30px',
                  left: '50%',
                  transform: 'translateX(-50%)',
                  whiteSpace: 'nowrap',
                  background: '#15181F',
                  color: '#fff',
                  fontSize: 11.5,
                  fontWeight: 600,
                  padding: '3px 8px',
                  borderRadius: 6,
                  boxShadow: 'var(--uk-shadow)',
                }}
              >
                {s.name}
              </span>
            )}
          </button>
        );
      })}
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
        <button type="button" className="map-zoom" aria-label="Zoom in">
          <AddRounded sx={{ fontSize: 18 }} />
        </button>
        <hr className="divider" />
        <button type="button" className="map-zoom" aria-label="Zoom out">
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
    </div>
  );
}
