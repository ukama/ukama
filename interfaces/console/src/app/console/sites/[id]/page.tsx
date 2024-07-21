/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
'use client';

import React, { useState, useEffect } from 'react';
import colors from '@/theme/colors';
import LoadingWrapper from '@/components/LoadingWrapper';
import SiteOverView from '@/components/SiteOverView';
import RestartSiteDialog from '@/components/RestartSiteDialog';
import { AlertColor, Grid, Paper } from '@mui/material';
import SiteOverallHealth from '@/components/SiteHealth';
import { useGetSiteLazyQuery } from '@/client/graphql/generated';
import { useAppContext } from '@/context';
interface SiteDetailsProps {
  params: {
    id: string;
  };
}
const Page: React.FC<SiteDetailsProps> = ({ params }) => {
  const { id } = params;
  console.log('SITE ID', id);
  const [open, setOpen] = useState(false);
  const [restartDialogOpen, setRestartDialogOpen] = useState(false);
  const { setSnackbarMessage } = useAppContext();

  const handleClose = () => {
    setOpen(false);
  };
  const handleRestartSite = (site: { id: string; name: string }) => {
    console.log(`Restarting site: ${id}`);
    setRestartDialogOpen(true);
  };

  const handleAddSite = () => {
    setOpen(true);
  };

  const handleSiteInstallation = (formData: any) => {
    console.log('Function not implemented.', formData);
  };
  const handleRestartDialogClose = () => {
    setRestartDialogOpen(false);
  };
  const handleConfirmRestart = (siteName: string) => {
    console.log(`Restarting site: ${siteName}`);
    setRestartDialogOpen(false);
  };

  const [getSite] = useGetSiteLazyQuery({
    onError: (error) => {
      setSnackbarMessage({
        id: 'sites-msg',
        message: error.message,
        type: 'error' as AlertColor,
        show: true,
      });
    },
  });
  useEffect(() => {
    getSite({ variables: { siteId: id } });
  }, []);

  return (
    <LoadingWrapper
      radius="small"
      width={'100%'}
      isLoading={false}
      cstyle={{
        backgroundColor: false ? colors.white : 'transparent',
      }}
    >
      <Grid container spacing={2} sx={{ mt: 1 }}>
        <SiteOverView restartSite={handleRestartSite} addSite={handleAddSite} />
      </Grid>
      <Paper elevation={3} sx={{ p: 4, mt: 2 }}>
        <SiteOverallHealth
          batteryInfo={[]}
          solarHealth={'good'}
          nodeHealth={'good'}
          switchHealth={'good'}
          controllerHealth={'good'}
          batteryHealth={'good'}
          backhaulHealth={'good'}
        />
      </Paper>

      {/* <ConfigureSiteDialog
        open={openSiteConfig}
        onClose={handleCloseSiteConfig}
        components={dummysComponents}
        networks={dummyNetworks}
      /> */}
      <RestartSiteDialog
        open={restartDialogOpen}
        onClose={handleRestartDialogClose}
        onConfirm={handleConfirmRestart}
        siteName={'selectedSite.name'}
      />
    </LoadingWrapper>
  );
};

export default Page;
