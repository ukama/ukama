/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Setup complete. ConfigureShell clears the saved resume point on this
 * route; onboardingStatus was refetched by the network/site mutations, so
 * the dashboard renders without the setup bar.
 */
'use client';

import { useRouter } from 'next/navigation';
import Button from '@mui/material/Button';
import CheckCircleRounded from '@mui/icons-material/CheckCircleRounded';

export default function ConfigureCompletePage() {
  const router = useRouter();
  return (
    <>
      <CheckCircleRounded
        sx={{ fontSize: 44, color: 'var(--uk-success-bright)' }}
      />
      <h1 className="cfg-title" style={{ marginTop: 14 }}>
        You&apos;re all set!
      </h1>
      <p className="cfg-copy">
        Your network is configured. Head to the Console to monitor your
        sites, manage data plans, and connect your first subscribers.
      </p>
      <div className="cfg-actions">
        <span />
        <Button variant="contained" onClick={() => router.push('/')}>
          Go to Console
        </Button>
      </div>
    </>
  );
}
