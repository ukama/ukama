/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * SIMs step — never blocks activation. Points at the SIM pool tooling in
 * the console; Continue finishes the flow.
 */
'use client';

import { useRouter } from 'next/navigation';
import SimCardRounded from '@mui/icons-material/SimCardRounded';

import ConfigureActions from '../_components/ConfigureActions';

export default function ConfigureSimsPage() {
  const router = useRouter();
  return (
    <>
      <h1 className="cfg-title">Upload SIMs</h1>
      <p className="cfg-copy">
        Your SIM pool holds the SIMs you&apos;ll assign to subscribers. You
        can upload SIMs anytime from the Console under{' '}
        <strong>Manage&nbsp;→&nbsp;SIM&nbsp;pool</strong> — it&apos;s not
        required to finish setup.
      </p>
      <div
        style={{
          display: 'flex',
          alignItems: 'center',
          gap: 10,
          color: 'var(--uk-ink-2)',
          fontSize: 14,
        }}
      >
        <SimCardRounded sx={{ fontSize: 22, color: 'var(--uk-ac)' }} />
        Have your SIM batch file ready? Head to the SIM pool after setup.
      </div>
      <ConfigureActions
        nextLabel="Finish setup"
        onNext={() => router.push('/configure/complete')}
        showSkip={false}
      />
    </>
  );
}
