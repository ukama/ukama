/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
'use client';
import { INSTALLATION_FLOW, ONBOARDING_FLOW } from '@/constants';
import colors from '@/theme/colors';
import {
  Button,
  Checkbox,
  FormControlLabel,
  Stack,
  Typography,
} from '@mui/material';
import { useRouter, useSearchParams } from 'next/navigation';
import { useState } from 'react';

const Page = () => {
  const router = useRouter();
  const searchParams = useSearchParams();
  const flow = searchParams.get('flow') ?? INSTALLATION_FLOW;
  const totalSteps = flow === ONBOARDING_FLOW ? 5 : 4;
  const step = parseInt(searchParams.get('step') ?? '1');
  const [isChecked, setIsChecked] = useState(false);

  const handleNext = () => {
    if (isChecked) {
      const id = 'uk-sa9001-tnode-a1-1234';
      router.push(
        `/configure/node/${id}?step=${step + 1}&flow=${flow}&lat=-4.322447&lng=15.307045`,
      );
    }
  };
  const handleSkip = () => {
    //TODO: HANDLE SKIP LOGIC
    // router.push('/console/home');
  };

  const handleOnInstalled = (isChecked: boolean) => {
    setIsChecked(isChecked);
  };

  return (
    <Stack direction="column" spacing={2}>
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
