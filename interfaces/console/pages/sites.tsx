/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { colors } from '@/styles/theme';
import LoadingWrapper from '@/ui/molecules/LoadingWrapper';
// import Map from '@/ui/molecules/MapComponent';
import { Grid, Typography, Button } from '@mui/material';
import SiteCard from '@/ui/molecules/SiteCard';

export default function Page() {
  interface SiteInt {
    name: string;
    details: string;
    batteryStatus: 'charging' | 'notCharging';
    nodeStatus: 'online' | 'offline';
    towerStatus: 'online' | 'offline';
    numberOfPersonsConnected: number;
  }

  const fakeData: SiteInt[] = [
    {
      name: 'Site 1',
      details: 'Details for Site 1',
      batteryStatus: 'charging',
      nodeStatus: 'online',
      towerStatus: 'online',
      numberOfPersonsConnected: 3,
    },
    {
      name: 'Site 2',
      details: 'Details for Site 2',
      batteryStatus: 'notCharging',
      nodeStatus: 'offline',
      towerStatus: 'offline',
      numberOfPersonsConnected: 5,
    },
    {
      name: 'Site 2',
      details: 'Details for Site 2',
      batteryStatus: 'charging',
      nodeStatus: 'offline',
      towerStatus: 'offline',
      numberOfPersonsConnected: 5,
    },
  ];
  const handleAddSite = () => {};

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
        <Grid container spacing={2}>
          <Grid item xs={6}>
            <Typography variant="h6"> My sites</Typography>
          </Grid>
          <Grid item xs={6} container justifyContent={'flex-end'}>
            <Button variant="contained" onClick={handleAddSite}>
              ADD SITE
            </Button>
          </Grid>
          {fakeData.map((site, index) => (
            <Grid item xs={4} key={index}>
              <SiteCard sites={[site]} />
            </Grid>
          ))}
        </Grid>
      </LoadingWrapper>
    </>
  );
}
