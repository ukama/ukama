/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/**
 * Three-way section contract for composite queries (BFF plan §4.5):
 * - data null + error set            → failed: "—" / chip, retry if offered
 * - error.code === NOT_IMPLEMENTED   → backend gap: "—", NO retry
 * - data null/empty + error null     → genuinely empty: empty state
 *
 * Branch on `error.code`, never on message text.
 */
import ReplayRoundedIcon from '@mui/icons-material/ReplayRounded';
import { Box, IconButton, Tooltip, Typography } from '@mui/material';

import { SectionErrorCode } from '@/client/graphql/types';

export interface SectionErrorLike {
  code: SectionErrorCode;
  message: string;
}

/** True when a section value should render as data (no error, value set). */
export function sectionReady<T>(
  value: T | null | undefined,
  error?: SectionErrorLike | null
): value is T {
  return value != null && !error;
}

/**
 * Inline placeholder for a failed or not-implemented section value.
 * Renders "—" with a tooltip; failed (non-gap) sections may offer retry.
 */
export function SectionFallback({
  error,
  onRetry,
}: {
  error: SectionErrorLike;
  onRetry?: () => void;
}) {
  const isGap = error.code === SectionErrorCode.NotImplemented;
  return (
    <Tooltip title={isGap ? 'Not available yet' : error.message}>
      <Box component="span" sx={{ display: 'inline-flex', alignItems: 'center', gap: 0.5 }}>
        <Typography component="span" color="text.disabled">
          —
        </Typography>
        {!isGap && onRetry && (
          <IconButton size="small" aria-label="Retry" onClick={onRetry}>
            <ReplayRoundedIcon sx={{ fontSize: 14 }} />
          </IconButton>
        )}
      </Box>
    </Tooltip>
  );
}

/** "—" or the formatted value, per the §4.5 contract. */
export function sectionValue(
  value: string | number | null | undefined,
  error?: SectionErrorLike | null
): string {
  if (error || value == null) return '—';
  return String(value);
}
