/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/** Compact key/value row for detail cards (node-site-detail.jsx KV). */
import ErrorRounded from '@mui/icons-material/ErrorRounded';
import ChevronRightRounded from '@mui/icons-material/ChevronRightRounded';

export default function KV({
  k,
  v,
  warn,
  vColor,
  onClick,
  active,
}: {
  k: string;
  v: React.ReactNode;
  warn?: boolean;
  vColor?: string | null;
  onClick?: () => void;
  active?: boolean;
}) {
  return (
    <div
      className="kv-row"
      role={onClick ? 'button' : undefined}
      tabIndex={onClick ? 0 : undefined}
      onClick={onClick}
      onKeyDown={
        onClick
          ? (e) => {
              if (e.key === 'Enter' || e.key === ' ') {
                e.preventDefault();
                onClick();
              }
            }
          : undefined
      }
      style={{
        display: 'flex',
        justifyContent: 'space-between',
        alignItems: 'center',
        gap: 14,
        padding: onClick ? '10px 10px' : '10px 0',
        margin: onClick ? '0 -10px' : undefined,
        borderRadius: onClick ? 'var(--uk-r-sm)' : undefined,
        cursor: onClick ? 'pointer' : undefined,
        background: active ? 'var(--uk-ac-soft)' : undefined,
      }}
    >
      <span
        style={{
          fontSize: 13,
          color: 'var(--uk-ink-2)',
          display: 'flex',
          alignItems: 'center',
          gap: 6,
        }}
      >
        {warn && <ErrorRounded sx={{ fontSize: 15, color: 'var(--uk-orange)' }} />}
        {k}
      </span>
      <span style={{ display: 'flex', alignItems: 'center', gap: 6 }}>
        <span
          className="tnum"
          style={{
            fontSize: 13.5,
            fontWeight: 600,
            color: vColor ?? 'var(--uk-ink)',
            textAlign: 'right',
          }}
        >
          {v}
        </span>
        {onClick && (
          <ChevronRightRounded
            className="kv-chev"
            sx={{ fontSize: 18, color: 'var(--uk-ink-3)', flex: 'none' }}
          />
        )}
      </span>
    </div>
  );
}
