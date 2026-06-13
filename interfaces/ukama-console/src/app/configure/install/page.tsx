/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/** Install-site intro (Figma "Install site"). Skippable; physical step. */
'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import Checkbox from '@mui/material/Checkbox';
import FormControlLabel from '@mui/material/FormControlLabel';
import Skeleton from '@mui/material/Skeleton';

import ConfigureActions from '../_components/ConfigureActions';
import { stepUrl, useConfigureParams } from '../_components/state';

export default function ConfigureInstallPage() {
  const router = useRouter();
  const { flow, networkid } = useConfigureParams();
  const [installed, setInstalled] = useState(false);

  // Self-guard: this step installs a site into a specific network, so a missing
  // networkid (e.g. a stale resume URL or an entry point that didn't pick one)
  // sends the user to the select-network step first.
  useEffect(() => {
    if (!networkid) {
      router.replace(stepUrl('/configure/select-network', { flow }));
    }
  }, [networkid, flow, router]);

  if (!networkid) {
    return (
      <>
        <Skeleton width="60%" height={42} />
        <Skeleton width="100%" height={28} />
        <Skeleton width="40%" height={28} sx={{ mt: 2 }} />
      </>
    );
  }

  return (
    <>
      <h1 className="cfg-title">Install your site</h1>
      <p className="cfg-copy">
        A site is a full connection point to your network, made up of three
        Ukama units installed together along with their power and backhaul.
        <br />
        <br />
        Install all the units at their location, then switch them on and connect
        them to the internet. Once every unit is powered on and online,
        we&apos;ll detect your site automatically on the next step. If
        you&apos;d like to do this later, skip — you can pick up where you left
        off anytime.
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
