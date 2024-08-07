/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
'use client';

import {
  SiteDto,
  useAddSiteMutation,
  useGetComponentsByUserIdLazyQuery,
  useGetNetworksQuery,
  useGetSitesQuery,
  useUpdateSiteMutation,
} from '@/client/graphql/generated';
import ConfigureSiteDialog from '@/components/ConfigureSiteDialog';
import EditSiteDialog from '@/components/EditSiteDialog';
import SitesWrapper from '@/components/SitesWrapper';
import { useAppContext } from '@/context';
import { TSiteForm } from '@/types';
import {
  AlertColor,
  Box,
  Button,
  Grid,
  Paper,
  Typography,
} from '@mui/material';
import { formatISO } from 'date-fns';
import { useEffect, useState } from 'react';

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

const Sites = () => {
  const [sitesList, setSitesList] = useState<SiteDto[]>([]);
  const [componentsList, setComponentsList] = useState<any[]>([]);
  const { setSnackbarMessage, network } = useAppContext();
  const [openSiteConfig, setOpenSiteConfig] = useState(false);
  const [site, setSite] = useState<TSiteForm>(SITE_INIT);
  const [editSitedialogOpen, setEditSitedialogOpen] = useState(false);
  const [currentSite, setCurrentSite] = useState({
    siteName: '',
    siteId: '',
  });
  const { refetch: refetchSites, loading: sitesLoading } = useGetSitesQuery({
    skip: !network.id,
    variables: {
      networkId: network.id,
    },
    onCompleted: (res) => {
      setSitesList(res.getSites.sites);
    },
    onError: (error) => {
      setSnackbarMessage({
        id: 'fetching-sites-msg',
        message: error.message,
        type: 'error' as AlertColor,
        show: true,
      });
    },
  });

  const [addSite, { loading: addSiteLoading }] = useAddSiteMutation({
    onCompleted: () => {
      refetchSites().then((res) => {
        setSitesList(res.data.getSites.sites);
      });
      setSnackbarMessage({
        id: 'add-site-success',
        message: 'Site added successfully!',
        type: 'success' as AlertColor,
        show: true,
      });
    },
    onError: (error) => {
      setSnackbarMessage({
        id: 'add-site-error',
        message: error.message,
        type: 'error' as AlertColor,
        show: true,
      });
    },
  });
  const [updateSite, { loading: updateSiteLoading }] = useUpdateSiteMutation({
    onCompleted: () => {
      refetchSites().then((res) => {
        setSitesList(res.data.getSites.sites);
      });
      setSnackbarMessage({
        id: 'update-site-success',
        message: 'Site updated successfully!',
        type: 'success' as AlertColor,
        show: true,
      });
    },
    onError: (error) => {
      setSnackbarMessage({
        id: 'update-site-error',
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

  const { data: networks, loading: networksLoading } = useGetNetworksQuery({
    fetchPolicy: 'cache-and-network',
    onCompleted: (res) => {
      if (res.getNetworks.networks.length === 0) {
        setSnackbarMessage({
          id: 'no-network-msg',
          message: 'Please create a network first.',
          type: 'warning' as AlertColor,
          show: true,
        });
      }
    },
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
    if (!network.id)
      setSite({
        ...site,
        network: network.id,
      });
    getComponents({ variables: { category: 'switch' } });
  }, []);

  const handleSiteConfigOpen = () => {
    setOpenSiteConfig(true);
  };

  const handleCloseSiteConfig = () => {
    setSite(SITE_INIT);
    setOpenSiteConfig(false);
  };

  const handleSiteConfiguration = (values: TSiteForm) => {
    setSite(values);
    setOpenSiteConfig(false);
    addSite({
      variables: {
        data: {
          name: values.siteName,
          power_id: values.power,
          location: values.address,
          access_id: values.access,
          switch_id: values.switch,
          latitude: values.latitude,
          network_id: values.network,
          longitude: values.longitude,
          backhaul_id: values.backhaul,
          spectrum_id: values.spectrum,
          install_date: formatISO(new Date()),
        },
      },
    });
  };
  const handleSiteNameUpdate = (siteId: string, siteName: string) => {
    setCurrentSite((prevState) => ({
      ...prevState,
      siteId,
      siteName: siteName,
    }));
    setEditSitedialogOpen(true);
  };
  const handleSaveSiteName = (siteId: string, siteName: string) => {
    updateSite({
      variables: {
        siteId: siteId,
        data: {
          name: siteName,
        },
      },
    });
  };

  const closeEditSiteDialog = () => {
    setEditSitedialogOpen(false);
  };
  return (
    <Box mt={2}>
      <Paper
        sx={{
          p: 4,
          overflow: 'auto',
          padding: '24px 32px',
          borderRadius: '10px',
          height: 'calc(100vh - 228px)',
        }}
      >
        <Grid container rowSpacing={2}>
          <Grid item xs={6}>
            <Typography variant="h6" color="initial">
              My sites
            </Typography>
          </Grid>
          <Grid
            item
            xs={6}
            container
            justifyItems={'center'}
            justifyContent={'flex-end'}
          >
            <Button
              color="primary"
              variant="contained"
              onClick={handleSiteConfigOpen}
              disabled={networks?.getNetworks?.networks.length === 0}
            >
              ADD SITE
            </Button>
          </Grid>

          <Grid item xs={12} mt={1.5}>
            {
              <SitesWrapper
                loading={sitesLoading || networksLoading}
                sites={sitesList}
                handleSiteNameUpdate={handleSiteNameUpdate}
              />
            }
          </Grid>
        </Grid>
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
      <EditSiteDialog
        open={editSitedialogOpen}
        siteId={currentSite.siteId}
        currentSiteName={currentSite.siteName}
        onClose={closeEditSiteDialog}
        onSave={handleSaveSiteName}
        updateSiteLoading={updateSiteLoading}
      />
    </Box>
  );
};

export default Sites;
