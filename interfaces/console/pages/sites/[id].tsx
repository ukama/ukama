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
      <Typography variant="body1" color="initial">
        hello{' '}
      </Typography>
    </>
  );
}
