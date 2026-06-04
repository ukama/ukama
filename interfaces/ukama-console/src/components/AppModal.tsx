/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/** Modal wrapper matching the prototype dialog anatomy (ds.jsx Modal). */
import Dialog from '@mui/material/Dialog';
import IconButton from '@mui/material/IconButton';
import Typography from '@mui/material/Typography';
import CloseRounded from '@mui/icons-material/CloseRounded';

export default function AppModal({
  title,
  onClose,
  children,
  footer,
  width = 460,
  open = true,
}: {
  title: string;
  onClose: () => void;
  children: React.ReactNode;
  footer?: React.ReactNode;
  width?: number;
  open?: boolean;
}) {
  return (
    <Dialog
      open={open}
      onClose={onClose}
      slotProps={{ paper: { sx: { width, maxWidth: '94vw', borderRadius: 3.5 } } }}
    >
      <div
        style={{
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'space-between',
          padding: '20px 24px 10px',
        }}
      >
        <Typography
          sx={{ fontFamily: 'var(--font-display)', fontSize: 20, fontWeight: 500 }}
        >
          {title}
        </Typography>
        <IconButton size="small" onClick={onClose} aria-label="Close" sx={{ color: 'text.disabled' }}>
          <CloseRounded />
        </IconButton>
      </div>
      <div style={{ padding: '4px 24px 8px' }}>{children}</div>
      {footer && (
        <div
          style={{
            display: 'flex',
            justifyContent: 'flex-end',
            gap: 10,
            padding: '16px 24px 22px',
          }}
        >
          {footer}
        </div>
      )}
    </Dialog>
  );
}
