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
        To install your full site, please install your node(s), power, and
        backhaul components at their intended location(s). These three
        elements form a site — a full connection point to the network.
        <br />
        <br />
        If you&apos;d like to install your site later, skip this step — you
        can pick up where you left off anytime.
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
        label="I have installed my site"
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
