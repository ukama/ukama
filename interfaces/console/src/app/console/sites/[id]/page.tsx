/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
'use client';

import {
  useAddSiteMutation,
  useGetComponentsByUserIdLazyQuery,
  useGetNetworksQuery,
  useGetSiteLazyQuery,
} from '@/client/graphql/generated';
import LoadingWrapper from '@/components/LoadingWrapper';
import RestartSiteDialog from '@/components/RestartSiteDialog';
import SiteOverallHealth from '@/components/SiteHealth';
import SiteOverView from '@/components/SiteOverView';
import { useAppContext } from '@/context';
import colors from '@/theme/colors';
import { AlertColor, Grid, Paper } from '@mui/material';
import React, { useEffect, useState } from 'react';

import ConfigureSiteDialog from '@/components/ConfigureSiteDialog';
import { TSiteForm } from '@/types';

const SITE_INIT = {
  switch: '',
  power: '',
  access: '',
  backhaul: '',
  address: '',
  spectrum: '',
  siteName: '',
  latitude: NaN,
  longitude: NaN,
  network: '',
};

interface SiteDetailsProps {
  params: {
    id: string;
  };
}
const Page: React.FC<SiteDetailsProps> = ({ params }) => {
  const { id } = params;
  const [site, setSite] = useState<TSiteForm>(SITE_INIT);
  const [restartDialogOpen, setRestartDialogOpen] = useState(false);
  const { setSnackbarMessage } = useAppContext();
  const [openSiteConfig, setOpenSiteConfig] = useState(false);
  const [componentsList, setComponentsList] = useState<any[]>([]);

  const handleSiteConfigOpen = () => {
    setOpenSiteConfig(true);
  };

  const handleCloseSiteConfig = () => {
    setOpenSiteConfig(false);
  };

  const handleRestartSite = (site: { id: string; name: string }) => {
    //console.log(`Restarting site: ${id}`);
    setRestartDialogOpen(true);
  };

  const handleRestartDialogClose = () => {
    setRestartDialogOpen(false);
  };
  const handleConfirmRestart = (siteName: string) => {
    setRestartDialogOpen(false);
  };
  const [addSite, { loading: addSiteLoading }] = useAddSiteMutation({
    onCompleted: (res) => {
      setSnackbarMessage({
        id: 'add-site-success',
        message: 'Site added successfully!',
        type: 'success' as AlertColor,
        show: true,
      });
    },
    onError: (error) => {
      setSnackbarMessage({
        id: 'add-subscriber-error',
        message: error.message,
        type: 'error' as AlertColor,
        show: true,
      });
    },
  });
  const getCurrentDateInISOFormat = () => {
    const date = new Date();
    return date.toISOString().split('T')[0] + 'T00:00:00Z';
  };

  const handleSiteConfiguration = async (data: any) => {
    const variables = {
      access_id: data.access,
      backhaul_id: data.backhaul,
      install_date: getCurrentDateInISOFormat(),
      latitude: data.coordinates.lat,
      location: data.location,
      longitude: data.coordinates.lng,
      name: data.siteName,
      network_id: data.selectedNetwork,
      power_id: data.power,
      spectrum_id: data.spectrumId || '',
      switch_id: data.switch,
      is_deactivated: data.is_deactivated || false,
    };

    try {
      await addSite({ variables: { data: variables } });
    } catch (error) {
      console.error('Error submitting site configuration:', error);
    }
  };

  const { data: networkList, loading: networkLoading } = useGetNetworksQuery({
    fetchPolicy: 'cache-and-network',
    onError: (error) => {
      setSnackbarMessage({
        id: 'networks-msg',
        message: error.message,
        type: 'error' as AlertColor,
        show: true,
      });
    },
  });

  const [getComponents] = useGetComponentsByUserIdLazyQuery({
    onCompleted: (res) => {
      if (res.getComponentsByUserId) {
        setComponentsList(res.getComponentsByUserId.components);
      }
    },
    onError: (error) => {
      setSnackbarMessage({
        id: 'components-msg',
        message: error.message,
        type: 'error' as AlertColor,
        show: true,
      });
    },
  });

  const { data: networks } = useGetNetworksQuery({
    fetchPolicy: 'cache-and-network',

    onError: (error) => {
      setSnackbarMessage({
        id: 'networks-msg',
        message: error.message,
        type: 'error' as AlertColor,
        show: true,
      });
    },
  });
  useEffect(() => {
    getComponents({ variables: { category: 'switch' } });
  }, []);

  const [getSite] = useGetSiteLazyQuery({
    onCompleted: (res) => {
      setSite({
        power: res.getSite.powerId,
        siteName: res.getSite.name,
        switch: res.getSite.switchId,
        access: res.getSite.accessId,
        address: res.getSite.location,
        latitude: res.getSite.latitude,
        network: res.getSite.networkId,
        spectrum: res.getSite.spectrumId,
        backhaul: res.getSite.backhaulId,
        longitude: res.getSite.longitude,
      });
    },
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
        <SiteOverView
          restartSite={handleRestartSite}
          addSite={handleSiteConfigOpen}
        />
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

      <ConfigureSiteDialog
        site={site}
        open={openSiteConfig}
        addSiteLoading={addSiteLoading}
        onClose={handleCloseSiteConfig}
        components={componentsList || []}
        networks={networks?.getNetworks?.networks || []}
        handleSiteConfiguration={handleSiteConfiguration}
      />
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
