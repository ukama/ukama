/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
'use client';
import { INSTALLATION_FLOW, ONBOARDING_FLOW } from '@/constants';
import { Button, Paper, Stack, Typography } from '@mui/material';
import { useRouter, useSearchParams } from 'next/navigation';

const SiteSuccess = () => {
  const router = useRouter();
  const searchParams = useSearchParams();
  const flow = searchParams.get('flow') ?? INSTALLATION_FLOW;

  return (
    <Paper elevation={0} sx={{ px: { xs: 2, md: 4 }, py: { xs: 1, md: 2 } }}>
      <Stack direction={'column'} spacing={2}>
        <Typography variant="h6">Network setup complete</Typography>
        <Typography variant="body1">
          {flow === ONBOARDING_FLOW
            ? 'Congratulations, you have successfully created your first network, and almost ready to experience reliable, fast, connectivity!'
            : 'Congratulations, you have successfully created network, and almost ready to experience reliable, fast, connectivity!'}
        </Typography>

        <Typography variant="body1">
          {flow === ONBOARDING_FLOW
            ? 'To get connected to the network, you still need to create a custom data plan, and add subscribers to your network.'
            : 'SOME TEXT HERE FOR NON ONBOARDING FLOW'}
        </Typography>
        <br />
        <Button
          variant="contained"
          sx={{ width: 'fit-content', alignSelf: 'flex-end' }}
          onClick={() => router.push('/')}
        >
          Continue to console
        </Button>
      </Stack>
    </Paper>
  );
};

export default SiteSuccess;
