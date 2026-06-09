/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Shared step footer: "Skip for now" (left, to the dashboard — the resume
 * point is already persisted by ConfigureShell) and the primary action
 * (right). Hide skip on steps that must not be skipped mid-mutation.
 */
'use client';

import Button from '@mui/material/Button';
import CircularProgress from '@mui/material/CircularProgress';
import { useRouter } from 'next/navigation';

export default function ConfigureActions({
  nextLabel,
  onNext,
  nextDisabled,
  busy,
  showSkip = true,
}: {
  nextLabel: string;
  onNext: () => void;
  nextDisabled?: boolean;
  busy?: boolean;
  showSkip?: boolean;
}) {
  const router = useRouter();
  return (
    <div className="cfg-actions">
      {showSkip ? (
        <Button
          variant="text"
          sx={{ color: 'var(--uk-ink-2)' }}
          onClick={() => router.push('/')}
          disabled={busy}
        >
          Skip for now
        </Button>
      ) : (
        <span />
      )}
      <Button
        variant="contained"
        onClick={onNext}
        disabled={nextDisabled || busy}
        startIcon={
          busy ? <CircularProgress size={16} color="inherit" /> : undefined
        }
      >
        {nextLabel}
      </Button>
    </div>
  );
}
