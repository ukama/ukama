/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

import Button, { type ButtonProps } from '@mui/material/Button';
import CircularProgress from '@mui/material/CircularProgress';

/**
 * Primary action button that reflects an async operation-lock state. While
 * `busy`, it shows a spinner + `busyLabel` and is inert. `reason` becomes the
 * tooltip/aria hint (e.g. "Restart in progress by alan") so the disabled state
 * always explains itself.
 */
export interface ActionButtonProps extends Omit<ButtonProps, 'disabled'> {
  busy?: boolean;
  busyLabel?: string;
  /** Disabled without a spinner: the action is unavailable, not in progress. */
  blocked?: boolean;
  /** Human reason shown as tooltip + aria hint when busy/blocked. */
  reason?: string;
}

export function ActionButton({
  busy = false,
  busyLabel,
  blocked = false,
  reason,
  children,
  startIcon,
  onClick,
  sx,
  ...rest
}: ActionButtonProps) {
  return (
    <Button
      {...rest}
      startIcon={busy ? <CircularProgress size={15} color="inherit" /> : startIcon}
      aria-disabled={busy || blocked || undefined}
      aria-label={reason && (busy || blocked) ? reason : undefined}
      title={reason && (busy || blocked) ? reason : undefined}
      onClick={busy || blocked ? undefined : onClick}
      sx={[
        ...(Array.isArray(sx) ? sx : [sx]),
        busy || blocked ? { opacity: 0.7, pointerEvents: 'none' } : {},
      ]}
    >
      {busy && busyLabel ? busyLabel : children}
    </Button>
  );
}
