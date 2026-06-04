/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/** Thin meter bar — MUI LinearProgress themed to the design (.meter). */
import LinearProgress from '@mui/material/LinearProgress';
import type { SxProps } from '@mui/material/styles';

export default function Meter({
  value,
  color,
  height,
  sx,
}: {
  /** 0–100 */
  value: number;
  /** bar color (CSS value); defaults to the accent */
  color?: string;
  height?: number;
  sx?: SxProps;
}) {
  return (
    <LinearProgress
      variant="determinate"
      value={Math.max(0, Math.min(100, value))}
      sx={[
        {
          ...(height != null && { height }),
          '& .MuiLinearProgress-bar': {
            backgroundColor: color ?? 'var(--uk-ac)',
          },
        },
        ...(Array.isArray(sx) ? sx : sx ? [sx] : []),
      ]}
    />
  );
}
