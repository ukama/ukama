'use client';

import colors from '@/theme/colors';
import LoadingWrapper from '@/components/LoadingWrapper';
import SiteCard from '@/components/SiteCard';
import { Grid, Paper, Typography, Button, AlertColor } from '@mui/material';
import ConfigureSiteDialog from '@/components/ConfigureSiteDialog';
import StepperDialog from '@/components/ConfigureSiteDialog';

import { useEffect, useState } from 'react';
import {
  useGetSitesLazyQuery,
  useGetNetworksQuery,
  useGetComponentsByUserIdLazyQuery,
  useAddSiteMutation,
} from '@/client/graphql/generated';
import { useAppContext } from '@/context';

const Sites = () => {
  const [sitesList, setSitesList] = useState<any[]>([]);
  const [componentsList, setComponentsList] = useState<any[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const { setSnackbarMessage } = useAppContext();
  const [openSiteConfig, setOpenSiteConfig] = useState(false);

  const handleSiteConfigOpen = () => {
    setOpenSiteConfig(true);
  };

  const handleCloseSiteConfig = () => {
    setOpenSiteConfig(false);
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
  // const handleSiteConfiguration = (data: any) => {
  //   console.log('FINAL DATA RES :', data);
  //   addSite({});
  // };
  const handleSiteConfiguration = async (data: any) => {
    console.log('SUBMITED DATA :', data);
    await addSite({
      variables: {
        data: {
          accessId: '',
          backhaulId: '',
          installDate: '',
          latitude: 0,
          location: '',
          longitude: 0,
          name: '',
          networkId: '',
          powerId: '',
          spectrumId: '',
          switchId: '',
        },
      },
    });
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

  const [getSites] = useGetSitesLazyQuery({
    onError: (error) => {
      setSnackbarMessage({
        id: 'sites-msg',
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

  useEffect(() => {
    const fetchAllSites = async () => {
      if (networkList && networkList.getNetworks.networks) {
        setIsLoading(true);
        const allSitesPromises = networkList.getNetworks.networks.map(
          (network) => getSites({ variables: { networkId: network.id } }),
        );

        try {
          const results = await Promise.all(allSitesPromises);
          const allSites = results.flatMap(
            (result) => result.data?.getSites.sites || [],
          );
          setSitesList(allSites);
        } catch (error) {
          console.error('Error fetching sites:', error);
          setSnackbarMessage({
            id: 'sites-fetch-error',
            message: 'Error fetching sites',
            type: 'error' as AlertColor,
            show: true,
          });
        } finally {
          setIsLoading(false);
        }
      }
    };

    fetchAllSites();
  }, [networkList, getSites]);

  const handleMenuClick = (siteId: string) => {
    console.log(`Menu clicked for siteId: ${siteId}`);
  };

  return (
    <Grid container spacing={0} sx={{ mt: 1 }}>
      <Grid item xs={12}>
        <Paper sx={{ p: 4 }}>
          <Grid container spacing={0} sx={{ mb: 2 }}>
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
                variant="contained"
                color="primary"
                onClick={handleSiteConfigOpen}
              >
                ADD SITE
              </Button>
            </Grid>
          </Grid>

          <Grid container spacing={2}>
            {sitesList.length === 0 ? (
              <Grid item xs={12} style={{ textAlign: 'center' }}>
                <Typography variant="body1">No sites available.</Typography>
              </Grid>
            ) : (
              sitesList.map((site, index) => (
                <Grid item xs={12} md={4} lg={4} key={index}>
                  <SiteCard
                    siteId={site.id}
                    name={site.name}
                    address={site.location}
                    users={site.users || ''}
                    siteStatus={site.isDeactivated}
                    onClickMenu={handleMenuClick}
                    loading={isLoading || networkLoading}
                    status={{
                      online: false,
                      charging: false,
                      signal: '',
                    }}
                  />
                </Grid>
              ))
            )}
          </Grid>
        </Paper>
      </Grid>
      <ConfigureSiteDialog
        open={openSiteConfig}
        onClose={handleCloseSiteConfig}
        components={componentsList || []}
        networks={networks?.getNetworks?.networks || []}
        handleSiteConfiguration={handleSiteConfiguration}
      />
    </Grid>
  );
};

export default Sites;
