/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { colors } from '@/styles/theme';
import { Site } from '@/types';
import LoadingWrapper from '@/ui/molecules/LoadingWrapper';
// import Map from '@/ui/molecules/MapComponent';
import SiteHeader from '@/ui/molecules/SiteHeader';
import SiteOverallHealth from '@/ui/molecules/SiteHealth';
import { Grid, Paper, Stack, Typography } from '@mui/material';
import { RoundedCard } from '@/styles/global';
import { SitePowerStatus } from '@/ui/molecules/SvgIcons';
import GroupIcon from '@mui/icons-material/Group';

const sites: Site[] = [
  { name: 'site1', health: 'online', duration: '3 days' },
  { name: 'site2', health: 'offline', duration: '1 week' },
  { name: 'site3', health: 'online', duration: '2 days' },
];

export default function Page() {
  const handleSiteSelect = (site: any): void => {};
  const handleAddSite = () => {
    // Logic to add a new site
  };
  const handleSiteRestart = () => {
    // Logic to restart a site
  };

  const batteryInfo = [
    { label: 'Model number', value: 'V1234' },
    { label: 'Current', value: '10 A' },
    { label: 'Charge', value: '80 %' },
    { label: 'Power', value: '100 W' },
    { label: 'Voltage', value: '12 V' },
  ];
  return (
    <>
      <LoadingWrapper
        radius="small"
        width={'100%'}
        isLoading={false}
        cstyle={{
          backgroundColor: false ? colors.white : 'transparent',
        }}
      >
        <SiteHeader
          sites={sites}
          sitesAction={handleSiteSelect}
          addSiteAction={handleAddSite}
          restartSiteAction={handleSiteRestart}
        />

        <Grid container spacing={2} sx={{ mt: 1 }}>
          <Grid item xs={3}>
            <RoundedCard>
              <Typography variant="h6" gutterBottom sx={{ py: 1 }}>
                Site details
              </Typography>
              <Stack direction="column" spacing={2}>
                <Stack direction="row" spacing={2} alignItems={'center'}>
                  <Typography variant="subtitle1">Date created:</Typography>
                  <Typography variant="body2"> July 13 2023</Typography>
                </Stack>
                <Stack direction="row" spacing={2} alignItems={'center'}>
                  <Typography variant="subtitle1"> Location:</Typography>
                  <Typography variant="body2"> 1000 Nelson Way</Typography>
                </Stack>
                <Stack direction="row" spacing={2} alignItems={'center'}>
                  <Typography variant="subtitle1"> Nodes:</Typography>
                  <Typography variant="body2"> 1000</Typography>
                </Stack>
              </Stack>
            </RoundedCard>
          </Grid>

          <Grid item xs={6}>
            <RoundedCard>
              <Typography variant="h6" sx={{ py: 1 }}>
                Site overview
              </Typography>
              <Stack direction="row" spacing={2} alignItems={'center'}>
                <Stack direction="row" spacing={1} alignItems={'center'}>
                  <SitePowerStatus />
                  <Typography variant="body1">Input power</Typography>
                  <SitePowerStatus />
                  <Typography variant="body1">Storage</Typography>
                  <SitePowerStatus />
                  <Typography variant="body1">Consumption</Typography>
                </Stack>
              </Stack>
            </RoundedCard>
          </Grid>
          <Grid item xs={3}>
            <RoundedCard>
              <Stack direction="row" spacing={1}>
                <GroupIcon />
                <Typography variant="body1">22</Typography>
              </Stack>
            </RoundedCard>
          </Grid>

          <Grid item xs={12}>
            <Paper sx={{ p: 2 }}>
              <Typography variant="h6" gutterBottom>
                Site components
              </Typography>

              <SiteOverallHealth
                solarHealth={'warning'}
                nodeHealth={'good'}
                switchHealth={'good'}
                controllerHealth={'good'}
                batteryHealth={'good'}
                backhaulHealth={'good'}
                batteryInfo={batteryInfo}
              />
            </Paper>
          </Grid>
        </Grid>
      </LoadingWrapper>
    </>
  );
}
