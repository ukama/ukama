/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
'use client';
import InstallSiteLoading from '@/components/InstallSiteLoading';
import colors from '@/theme/colors';
import {
  Button,
  Checkbox,
  FormControlLabel,
  Paper,
  Stack,
  SvgIcon,
  Typography,
} from '@mui/material';
import Link from 'next/link';
import { useRouter, useSearchParams } from 'next/navigation';
import { useState } from 'react';
import SiteInfo from '../../../public/svg/SiteInfo';

const Page = () => {
  const router = useRouter();
  const searchParams = useSearchParams();
  const flow = searchParams.get('flow') ?? 'onb';
  const step = parseInt(searchParams.get('step') ?? '1');
  const [isChecked, setIsChecked] = useState(false);
  const [checkForInstallation, setCheckForInstallation] = useState(false);
  const handleNext = () => {
    if (isChecked) {
      setCheckForInstallation(true);
    }
  };
  const handleSkip = () => router.push('/console/home');
  const handleOnInstalled = (isChecked: boolean) => {
    setIsChecked(isChecked);
  };
  const onInstallProgressComplete = () => {
    const id = 'uk-sa9001-tnode-a1-1234';
    router.push(
      `/configure/node/${id}?step=${step + 2}&flow=${flow}&lat=-4.322447&lng=15.307045`,
    );
  };

  return (
    <Paper elevation={0} sx={{ px: 4, py: 2 }}>
      {checkForInstallation ? (
        <InstallSiteLoading
          step={step + 1}
          flow={flow}
          duration={10}
          onCompleted={onInstallProgressComplete}
        />
      ) : (
        <Stack direction="column">
          <Stack direction={'row'}>
            <Typography variant="h6"> {'Install site'}</Typography>
            <Typography
              variant="h6"
              fontWeight={400}
              sx={{
                color: colors.black70,
              }}
            >
              {flow === 'onb' && <i>&nbsp;- optional</i>}&nbsp;({step}
              /6)
            </Typography>
          </Stack>

          <Stack mt={3} mb={3} direction={'column'} alignItems={'center'}>
            <Typography variant="body1" color={colors.vulcan}>
              If you would like to set up your network later, or if someone else
              will set up your network for you, skip this step.
              <br />
              <br />
              Install your node at the intended location, and ensure it is
              connected to a power and backhaul source. These three elements
              form a site, an abstracted representation of the aforementioned
              components. Each site can also hold up to three nodes for a
              stronger connection.
              <br />
              <br />
              You can follow the installation instructions in the provided
              manual, or in the PDF <Link href={''}>here</Link>.
            </Typography>
            <SvgIcon sx={{ width: 240, height: 176, mt: 2, mb: 4 }}>
              {SiteInfo}
            </SvgIcon>
            <FormControlLabel
              sx={{ alignSelf: 'baseline' }}
              control={
                <Checkbox
                  sx={{ p: 0, px: 1.3 }}
                  onChange={(e) => handleOnInstalled(e.target.checked)}
                />
              }
              label="I have installed my node, power, and backhaul components"
            />
          </Stack>
          <Stack
            mb={1}
            spacing={2}
            direction={'row'}
            justifyContent={'space-between'}
          >
            <Button
              variant="text"
              onClick={handleSkip}
              sx={{ color: colors.black70, p: 0 }}
            >
              {step !== 1 ? 'Back' : 'Skip'}
            </Button>
            <Button
              variant="contained"
              onClick={handleNext}
              disabled={!isChecked}
            >
              Next
            </Button>
          </Stack>
        </Stack>
      )}
    </Paper>
  );
};
export default Page;
