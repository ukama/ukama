/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/** Horizontal bar list — Revenue by site / by package (biz-common.jsx). */
export interface BarListRow {
  name: string;
  value: number;
  color: string;
}

export default function BarList({
  rows,
  prefix = '$',
}: {
  rows: BarListRow[];
  prefix?: string;
}) {
  const max = Math.max(...rows.map((r) => r.value)) || 1;
  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: 18 }}>
      {rows.map((r) => (
        <div
          key={r.name}
          style={{
            display: 'grid',
            gridTemplateColumns: '120px 1fr auto',
            alignItems: 'center',
            gap: 14,
          }}
        >
          <span style={{ fontSize: 13.5, color: 'var(--uk-ink-2)' }}>{r.name}</span>
          <div style={{ position: 'relative', height: 18 }}>
            <span
              style={{
                position: 'absolute',
                left: 0,
                top: 0,
                height: 18,
                borderRadius: 999,
                width: (r.value / max) * 100 + '%',
                background: r.color,
                minWidth: 14,
              }}
            />
          </div>
          <span
            className="tnum"
            style={{
              fontSize: 13.5,
              fontWeight: 600,
              color: 'var(--uk-ink)',
              minWidth: 48,
              textAlign: 'right',
            }}
          >
            {prefix}
            {r.value.toLocaleString()}
          </span>
        </div>
      ))}
    </div>
  );
}
