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
import AddSiteDialog from '@/ui/molecules/AddSiteDialog';
import { NetworkDto, useAddSiteMutation } from '@/generated';
import { useSetRecoilState } from 'recoil';
import { TSnackMessage } from '@/types';
import { snackbarMessage } from '@/app-recoil';
import DeleteConfirmation from '@/ui/molecules/DeleteSiteDialog';

export default function Page() {
  const [isAddSiteDialogOpen, setIsAddSiteDialogOpen] = useState(false);
  const setSnackbarMessage = useSetRecoilState<TSnackMessage>(snackbarMessage);
  const [isConfirmationOpen, setIsConfirmationOpen] = useState(false);
  const [siteId, setSiteId] = useState<any>();

  interface SiteInt {
    id: string;
    name: string;
    details: string;
    batteryStatus: 'charging' | 'notCharging';
    nodeStatus: 'online' | 'offline';
    towerStatus: 'online' | 'offline';
    numberOfPersonsConnected: number;
  }

  const fakeData: SiteInt[] = [
    {
      id: '9eb408c6-cdf0-4bc3-8802-6f546b7bede1',
      name: 'Site 1',
      details: 'Details for Site 1',
      batteryStatus: 'charging',
      nodeStatus: 'online',
      towerStatus: 'online',
      numberOfPersonsConnected: 3,
    },
    {
      id: '7eb408c6-cdf0-4bc3-8802-6f546b7bede1',
      name: 'Site 2',
      details: 'Details for Site 2',
      batteryStatus: 'notCharging',
      nodeStatus: 'offline',
      towerStatus: 'offline',
      numberOfPersonsConnected: 5,
    },
    {
      id: '9ec408c6-cdf0-4bc3-8802-6f546b7bede1',
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
  const handleDelete = async () => {};
  const handleCloseAction = () => setIsAddSiteDialogOpen(false);
  const handleDeleteSite = (siteId?: string) => {
    setIsConfirmationOpen(true);
    console.log('SITE ID', siteId);
    setSiteId(siteId);
  };
  const handleCancel = () => {
    // Handle cancel operation or close the confirmation dialog
    setIsConfirmationOpen(false);
  };

  const handleOpenConfirmation = () => {
    // Open the confirmation dialog
    setIsConfirmationOpen(true);
  };
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
            <Grid item xs={12} key={index} md={6} lg={4}>
              <SiteCard sites={[site]} handleDeleteSite={handleDeleteSite} />
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
        <DeleteConfirmation
          open={isConfirmationOpen}
          onDelete={handleDelete}
          onCancel={handleCancel}
          itemName={siteId}
          loading={false}
        />
      </LoadingWrapper>
    </>
  );
}
