/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/**
 * Coverage / sites map (Business lens) — soft-blue canvas with curved
 * arcs and labelled status dots (biz-common.jsx SiteMap).
 */

export interface SiteMapSite {
  id: string;
  name: string;
  status: 'online' | 'warning' | 'degraded' | 'offline';
  x: number;
  y: number;
}

export const BIZ_DOT: Record<SiteMapSite['status'], string> = {
  online: 'var(--uk-success)',
  warning: 'var(--uk-warning)',
  degraded: 'var(--uk-warning)',
  offline: 'var(--uk-error)',
};

export function StatusDot({ status }: { status: SiteMapSite['status'] }) {
  return (
    <span
      style={{
        width: 9,
        height: 9,
        borderRadius: '50%',
        flex: 'none',
        background: BIZ_DOT[status] ?? 'var(--uk-ink-3)',
      }}
    />
  );
}

export default function SiteMap({
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
          ...(fill ? { flex: 1, minHeight: 300 } : { height }),
        }}
      >
        <div
          style={{
            position: 'absolute',
            inset: 0,
            background: 'linear-gradient(160deg, #EEF5FE, #E7F0FC)',
          }}
        >
          <svg
            width="100%"
            height="100%"
            viewBox="0 0 100 60"
            preserveAspectRatio="xMidYMid slice"
            style={{ position: 'absolute', inset: 0 }}
            aria-hidden="true"
          >
            <g fill="none" stroke="#BFD6F2" strokeWidth="0.35" opacity="0.7">
              {Array.from({ length: 9 }).map((_, i) => (
                <path
                  key={i}
                  d={`M${-5 + i * 9},62 Q${8 + i * 9},${40 - (i % 3) * 6} ${24 + i * 9},58`}
                />
              ))}
              {Array.from({ length: 7 }).map((_, i) => (
                <path
                  key={'b' + i}
                  d={`M${-2 + i * 13},20 Q${10 + i * 13},${4 + (i % 2) * 8} ${26 + i * 13},22`}
                  strokeWidth="0.3"
                  opacity="0.6"
                />
              ))}
            </g>
          </svg>
        </div>
        {sites.map((s) => {
          const sel = selected === s.id;
          const col = BIZ_DOT[s.status];
          return (
            <div
              key={s.id}
              style={{
                position: 'absolute',
                left: `${s.x}%`,
                top: `${s.y}%`,
                transform: 'translate(-50%,-50%)',
                display: 'flex',
                alignItems: 'center',
                zIndex: sel ? 5 : 2,
              }}
            >
              <button
                type="button"
                aria-label={s.name}
                onClick={() => onSelect?.(s)}
                style={{
                  width: 16,
                  height: 16,
                  borderRadius: '50%',
                  background: col,
                  border: '3px solid #fff',
                  boxShadow: sel
                    ? `0 0 0 4px color-mix(in srgb, ${col} 27%, transparent), 0 2px 6px rgba(0,0,0,.25)`
                    : '0 1px 4px rgba(0,0,0,.2)',
                  cursor: onSelect ? 'pointer' : 'default',
                  padding: 0,
                  flex: 'none',
                }}
              />
              <span
                style={{
                  marginLeft: 8,
                  background: '#fff',
                  border: '1px solid #e7e9ee',
                  borderRadius: 7,
                  padding: '4px 10px',
                  fontSize: 12.5,
                  fontWeight: 500,
                  color: '#1c1e22',
                  boxShadow: '0 1px 3px rgba(0,0,0,.08)',
                  whiteSpace: 'nowrap',
                }}
              >
                {s.name}
              </span>
            </div>
          );
        })}
      </div>
    </div>
  );
}
