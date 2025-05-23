/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
'use client';
import { INSTALLATION_FLOW } from '@/constants';
import colors from '@/theme/colors';
import { setQueryParam } from '@/utils';
import {
  Button,
  Checkbox,
  FormControlLabel,
  Stack,
  Typography,
} from '@mui/material';
import { usePathname, useRouter, useSearchParams } from 'next/navigation';
import { useState } from 'react';

const Page = () => {
  const router = useRouter();
  const pathname = usePathname();
  const searchParams = useSearchParams();
  const flow = searchParams.get('flow') ?? INSTALLATION_FLOW;
  const [isChecked, setIsChecked] = useState(false);

  const handleNext = () => {
    if (isChecked) {
      const p = setQueryParam('flow', flow, searchParams.toString(), pathname);
      router.push(`/configure/check?${p}`);
    }
  };

  const handleSkip = () => {
    const p = setQueryParam('flow', flow, searchParams.toString(), pathname);
    router.push(`/configure/sims?flow=${p}`);
  };

  const handleOnInstalled = (isChecked: boolean) => {
    setIsChecked(isChecked);
  };

  return (
    <Stack direction="column" spacing={2} overflow={'scroll'}>
      <Typography variant="h4" fontWeight={500}>
        Install site
      </Typography>

      <Stack spacing={{ xs: 2, md: 4 }}>
        <Typography variant="body1" color={colors.vulcan}>
          If you would like to install your site later, please skip this step.
          <br />
          <br />
          To install your full site, please install your node(s), power, and
          backhaul components at their intended location(s). These three
          elements form a site, which represents a full connection point to the
          network. Each site can also hold up to three nodes for a stronger
          connection.
        </Typography>

        <FormControlLabel
          sx={{ alignSelf: 'baseline' }}
          control={
            <Checkbox
              sx={{ p: 0, pr: 1.5 }}
              onChange={(e) => handleOnInstalled(e.target.checked)}
            />
          }
          label="I have installed my site"
        />
      </Stack>

      <Stack
        direction={'row'}
        pt={{ xs: 4, md: 6 }}
        justifyContent={'space-between'}
      >
        <Button
          variant="text"
          onClick={handleSkip}
          sx={{ color: colors.black70, p: 0 }}
        >
          Skip
        </Button>
        <Button variant="contained" onClick={handleNext} disabled={!isChecked}>
          Next
        </Button>
      </Stack>
    </Stack>
  );
};
export default Page;
