/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/** Setup overview — entry point of the onboarding flow. */
'use client';

import { useRouter } from 'next/navigation';
import CheckCircleOutlineRounded from '@mui/icons-material/CheckCircleOutlineRounded';

import ConfigureActions from './_components/ConfigureActions';

const STEPS = [
  'Name your first network',
  'Install your first site (node, power, backhaul)',
  'Upload SIMs to your SIM pool',
];

export default function ConfigureOverviewPage() {
  const router = useRouter();
  return (
    <>
      <h1 className="cfg-title">Let&apos;s set up your network</h1>
      <p className="cfg-copy">
        A few steps and your Ukama network is live. You can skip any step and
        pick up right where you left off — we&apos;ll keep your place.
      </p>
      <div className="cfg-fields">
        {STEPS.map((s) => (
          <div
            key={s}
            style={{ display: 'flex', alignItems: 'center', gap: 10 }}
          >
            <CheckCircleOutlineRounded
              sx={{ fontSize: 20, color: 'var(--uk-ink-3)' }}
            />
            <span style={{ fontSize: 14.5, color: 'var(--uk-ink)' }}>{s}</span>
          </div>
        ))}
      </div>
      <ConfigureActions
        nextLabel="Get started"
        onNext={() => router.push('/configure/network')}
      />
    </>
  );
}
