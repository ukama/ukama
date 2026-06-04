/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/** Compact key/value row for detail cards (node-site-detail.jsx KV). */
import ErrorRounded from '@mui/icons-material/ErrorRounded';

export default function KV({
  k,
  v,
  warn,
  vColor,
}: {
  k: string;
  v: React.ReactNode;
  warn?: boolean;
  vColor?: string | null;
}) {
  return (
    <div
      style={{
        display: 'flex',
        justifyContent: 'space-between',
        alignItems: 'center',
        gap: 14,
        padding: '10px 0',
        borderBottom: '1px solid var(--uk-line-soft)',
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
    </div>
  );
}
