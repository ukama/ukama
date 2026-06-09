/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Setup alert bar (ACTIVATION-PLAN §5). Shown above the top bar while the
 * org's activation state is known-incomplete (no network or no site).
 * "Continue setup" resumes the /configure flow (saved progress wins).
 * Dismiss hides it for the session only; it disappears permanently once
 * derived state says activated — no dismissal bookkeeping.
 *
 * In 'hard-gate' mode this component redirects to the resume URL instead of
 * rendering a banner (dashboard unusable until activated).
 */
'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import Button from '@mui/material/Button';
import IconButton from '@mui/material/IconButton';
import CloseRounded from '@mui/icons-material/CloseRounded';
import InfoOutlined from '@mui/icons-material/InfoOutlined';

import { ACTIVATION_MODE, resolveResumeUrl, useActivation } from '@/lib/activation';
import { useUiPrefs } from '@/lib/store';

const DISMISS_KEY = 'uk-setup-bar-dismissed';

export default function SetupAlertBar() {
  const router = useRouter();
  const { status, needsSetup } = useActivation();
  const lastConfigureUrl = useUiPrefs((s) => s.lastConfigureUrl);
  // Lazy init is hydration-safe here: the banner can't render before the
  // first onboardingStatus result arrives (client-side only).
  const [dismissed, setDismissed] = useState(
    () =>
      typeof window !== 'undefined' &&
      sessionStorage.getItem(DISMISS_KEY) === '1',
  );

  // Hard gate: known-incomplete state never renders the dashboard.
  useEffect(() => {
    if (ACTIVATION_MODE === 'hard-gate' && needsSetup) {
      router.replace(resolveResumeUrl(status, lastConfigureUrl));
    }
  }, [needsSetup, status, lastConfigureUrl, router]);

  if (ACTIVATION_MODE === 'hard-gate') return null;
  if (!needsSetup || dismissed) return null;

  const message = status?.hasNetwork
    ? 'Almost there — install your first site to start seeing live data.'
    : 'Set up your first network to start seeing live data.';

  return (
    <div className="setup-bar" role="status">
      <InfoOutlined sx={{ fontSize: 18 }} />
      <span className="setup-bar-text">{message}</span>
      <Button
        size="small"
        variant="contained"
        disableElevation
        onClick={() =>
          router.push(resolveResumeUrl(status, lastConfigureUrl))
        }
      >
        Continue setup
      </Button>
      <IconButton
        size="small"
        aria-label="Dismiss setup reminder"
        onClick={() => {
          sessionStorage.setItem(DISMISS_KEY, '1');
          setDismissed(true);
        }}
        sx={{ color: 'inherit' }}
      >
        <CloseRounded sx={{ fontSize: 18 }} />
      </IconButton>
    </div>
  );
}
