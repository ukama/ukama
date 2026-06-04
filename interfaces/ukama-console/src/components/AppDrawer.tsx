/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/** Right drawer + header/row primitives (ds.jsx Drawer, detail.jsx). */
import Drawer from '@mui/material/Drawer';
import IconButton from '@mui/material/IconButton';
import Typography from '@mui/material/Typography';
import CloseRounded from '@mui/icons-material/CloseRounded';

export default function AppDrawer({
  onClose,
  children,
  width = 460,
  open = true,
}: {
  onClose: () => void;
  children: React.ReactNode;
  width?: number;
  open?: boolean;
}) {
  return (
    <Drawer
      anchor="right"
      open={open}
      onClose={onClose}
      slotProps={{
        paper: {
          sx: {
            width,
            maxWidth: '96vw',
            display: 'flex',
            flexDirection: 'column',
          },
        },
      }}
    >
      {children}
    </Drawer>
  );
}

export function DrawerHead({
  title,
  sub,
  badge,
  onClose,
}: {
  title: React.ReactNode;
  sub?: React.ReactNode;
  badge?: React.ReactNode;
  onClose: () => void;
}) {
  return (
    <div style={{ padding: '20px 24px 16px', borderBottom: '1px solid var(--uk-line)' }}>
      <div
        style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}
      >
        <div>
          <div style={{ display: 'flex', alignItems: 'center', gap: 10 }}>
            <Typography
              sx={{ fontFamily: 'var(--font-display)', fontSize: 20, fontWeight: 500 }}
            >
              {title}
            </Typography>
            {badge}
          </div>
          {sub && (
            <div style={{ fontSize: 13, color: 'var(--uk-ink-2)', marginTop: 3 }}>{sub}</div>
          )}
        </div>
        <IconButton size="small" onClick={onClose} aria-label="Close" sx={{ color: 'text.disabled' }}>
          <CloseRounded />
        </IconButton>
      </div>
    </div>
  );
}

/** Key/value row used in drawers and detail rails. */
export function DetailRow({
  k,
  v,
  vColor,
}: {
  k: string;
  v: React.ReactNode;
  vColor?: string;
}) {
  return (
    <div
      style={{
        display: 'flex',
        justifyContent: 'space-between',
        alignItems: 'center',
        padding: '11px 0',
        borderBottom: '1px solid var(--uk-line-soft)',
      }}
    >
      <span style={{ fontSize: 13, color: 'var(--uk-ink-2)' }}>{k}</span>
      <span
        className="tnum"
        style={{ fontSize: 13.5, fontWeight: 600, color: vColor ?? 'var(--uk-ink)' }}
      >
        {v}
      </span>
    </div>
  );
}
