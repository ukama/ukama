/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Unified KPI / stat card — single anatomy across every screen (ds.jsx):
 * label (optional tinted icon) · big value (+unit) · trend delta OR sub note.
 */
import ArrowDownwardRounded from '@mui/icons-material/ArrowDownwardRounded';
import ArrowUpwardRounded from '@mui/icons-material/ArrowUpwardRounded';
import { Ic } from '@/app/(dashboard)/_components/icons';

export interface KpiProps {
  icon?: string;
  color?: string;
  label: string;
  value: React.ReactNode;
  unit?: string;
  delta?: string;
  dir?: 'up' | 'down';
  sub?: React.ReactNode;
  danger?: boolean;
}

export function Kpi({
  icon,
  color,
  label,
  value,
  unit,
  delta,
  dir,
  sub,
  danger,
}: KpiProps) {
  return (
    <div className="card stat">
      <div className="stat-label">
        {icon && (
          <Ic name={icon} sx={{ fontSize: 18, color: color ?? 'var(--uk-ac)' }} />
        )}
        <span>{label}</span>
      </div>
      <div className="stat-figure">
        <span className="stat-val tnum">{value}</span>
        {unit && <span className="stat-unit">{unit}</span>}
      </div>
      {delta != null ? (
        <span className={`stat-delta ${dir === 'down' ? 'down' : 'up'}`}>
          {dir === 'down' ? (
            <ArrowDownwardRounded sx={{ fontSize: 13 }} />
          ) : (
            <ArrowUpwardRounded sx={{ fontSize: 13 }} />
          )}
          {delta}
        </span>
      ) : sub != null ? (
        <span className={`stat-sub${danger ? ' danger' : ''}`}>{sub}</span>
      ) : null}
    </div>
  );
}

export function KpiRow({ items, cols }: { items: KpiProps[]; cols?: number }) {
  return (
    <div
      className="tile-grid kpi-row"
      style={{
        gridTemplateColumns: `repeat(${cols ?? items.length}, minmax(0,1fr))`,
      }}
    >
      {items.map((k, i) => (
        <Kpi key={i} {...k} />
      ))}
    </div>
  );
}
