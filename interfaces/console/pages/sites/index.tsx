/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { colors } from '@/styles/theme';
import LoadingWrapper from '@/ui/molecules/LoadingWrapper';
import { useState } from 'react';
// import Map from '@/ui/molecules/MapComponent';
import { Grid, Typography, AlertColor, Button } from '@mui/material';
import SiteCard from '@/ui/molecules/SiteCard';
import Link from 'next/link';
import AddSiteDialog from '@/ui/molecules/AddSiteDialog';
import { NetworkDto, useAddSiteMutation } from '@/generated';
import { useSetRecoilState } from 'recoil';
import { TSnackMessage } from '@/types';
import { snackbarMessage } from '@/app-recoil';

export default function Page() {
  const [isAddSiteDialogOpen, setIsAddSiteDialogOpen] = useState(false);
  const setSnackbarMessage = useSetRecoilState<TSnackMessage>(snackbarMessage);

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
  const mockNetwork: NetworkDto[] = [
    {
      __typename: 'NetworkDto',
      budget: 1000000.0,
      countries: ['Country1', 'Country2', 'Country3'],
      createdAt: '2022-01-16T12:00:00Z',
      id: '1234567890',
      isDeactivated: 'false',
      name: 'Sample Network',
      networks: ['Network1', 'Network2', 'Network3'],
      orgId: 'organization123',
    },
  ];
  const [addSite, { loading: addSiteLoading }] = useAddSiteMutation({
    onCompleted: () => {
      setSnackbarMessage({
        id: 'site-added-success',
        message: 'Site added successfully!',
        type: 'success' as AlertColor,
        show: true,
      });
    },
    onError: (error) => {
      setSnackbarMessage({
        id: 'site-added-error',
        message: error.message,
        type: 'error' as AlertColor,
        show: true,
      });
    },
  });

  const handleAddSite = async (data: any) => {
    setIsAddSiteDialogOpen(true);
    // await addSite({
    //   variables: {
    //     data: {
    //       site: data.site,
    //     },
    //   },
    // });
  };
  const handleCloseAction = () => setIsAddSiteDialogOpen(false);
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
              <Link href={`sites/687687`}>
                <SiteCard sites={[site]} />
              </Link>
            </Grid>
          ))}
        </Grid>
        <AddSiteDialog
          isOpen={isAddSiteDialogOpen}
          title={'ADD SITE'}
          description=""
          handleCloseAction={handleCloseAction}
          networks={mockNetwork}
          handleAddSite={handleAddSite}
        />
      </LoadingWrapper>
    </>
  );
}
