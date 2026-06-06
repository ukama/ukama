/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/** Install-site intro (Figma "Install site"). Skippable; physical step. */
'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import Checkbox from '@mui/material/Checkbox';
import FormControlLabel from '@mui/material/FormControlLabel';

import ConfigureActions from '../_components/ConfigureActions';
import { stepUrl, useConfigureParams } from '../_components/state';

export default function ConfigureInstallPage() {
  const router = useRouter();
  const { flow, networkid } = useConfigureParams();
  const [installed, setInstalled] = useState(false);

  return (
    <>
      <h1 className="cfg-title">Install your site</h1>
      <p className="cfg-copy">
        A site is a full connection point to your network, made up of three
        Ukama units installed together along with their power and backhaul.
        <br />
        <br />
        Install all the units at their location, then switch them on and connect
        them to the internet. Once every unit is powered on and online, we&apos;ll
        detect your site automatically on the next step. If you&apos;d like to do
        this later, skip — you can pick up where you left off anytime.
      </p>
      <FormControlLabel
        sx={{ alignSelf: 'baseline' }}
        control={
          <Checkbox
            checked={installed}
            onChange={(e) => setInstalled(e.target.checked)}
            sx={{ p: 0, pr: 1.5 }}
          />
        }
        label="I've installed and powered on all my units"
      />
      <ConfigureActions
        nextLabel="Next"
        nextDisabled={!installed}
        onNext={() =>
          router.push(stepUrl('/configure/site', { flow, networkid }))
        }
      />
    </>
  );
}
