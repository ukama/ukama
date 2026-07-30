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
 * MUI Card + sx (§7.2 C).
 */
import Box from '@mui/material/Box';
import Card from '@mui/material/Card';
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
  /** Truncate a long text value to one line (smaller font + ellipsis). */
  truncate?: boolean;
}

/** Trend delta — green up / red down (stat-delta). */
export function Delta({
  dir = 'up',
  children,
  fontSize = 12,
}: {
  dir?: 'up' | 'down';
  children: React.ReactNode;
  fontSize?: number;
}) {
  return (
    <Box
      component="span"
      sx={(t) => ({
        fontSize,
        fontWeight: 600,
        display: 'inline-flex',
        alignItems: 'center',
        gap: '2px',
        color: dir === 'down' ? 'var(--uk-error-deep, #cf121b)' : 'var(--uk-success)',
        ...t.applyStyles('dark', {
          color: dir === 'down' ? '#ff8a8a' : 'var(--uk-success-bright)',
        }),
      })}
    >
      {dir === 'down' ? (
        <ArrowDownwardRounded sx={{ fontSize: fontSize + 1 }} />
      ) : (
        <ArrowUpwardRounded sx={{ fontSize: fontSize + 1 }} />
      )}
      {children}
    </Box>
  );
}

// For now KPI cards show only the label + value (with unit); the trend delta
// and sub-note are intentionally not rendered. The delta/dir/sub/danger props
// remain on KpiProps so callers are unchanged and it's easy to restore later.
export function Kpi({ icon, color, label, value, unit, truncate }: KpiProps) {
  return (
    <Card
      sx={{
        p: '16px 18px',
        display: 'flex',
        flexDirection: 'column',
        gap: '9px',
        minWidth: 0,
      }}
    >
      <Box
        sx={{
          fontSize: 13,
          color: 'var(--uk-ink-2)',
          fontWeight: 500,
          display: 'flex',
          alignItems: 'center',
          gap: 1,
        }}
      >
        {icon && (
          <Ic name={icon} sx={{ fontSize: 18, color: color ?? 'var(--uk-ac)' }} />
        )}
        <span>{label}</span>
      </Box>
      <Box
        sx={{
          display: 'flex',
          alignItems: 'baseline',
          gap: '6px',
          whiteSpace: 'nowrap',
          minWidth: 0,
          overflow: 'hidden',
        }}
      >
        <Box
          component="span"
          className="tnum"
          title={truncate && typeof value === 'string' ? value : undefined}
          sx={{
            fontFamily: 'var(--font-display)',
            fontSize: truncate ? 19 : 26,
            fontWeight: 500,
            lineHeight: truncate ? 1.25 : 1,
            ...(truncate
              ? {
                  display: 'block',
                  overflow: 'hidden',
                  textOverflow: 'ellipsis',
                  maxWidth: '100%',
                }
              : {}),
          }}
        >
          {value}
        </Box>
        {unit && (
          <Box component="span" sx={{ fontSize: 14, color: 'var(--uk-ink-3)' }}>
            {unit}
          </Box>
        )}
      </Box>
    </Card>
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
