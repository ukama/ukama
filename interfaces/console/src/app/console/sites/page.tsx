'use client';

import colors from '@/theme/colors';
import LoadingWrapper from '@/components/LoadingWrapper';
import SiteCard from '@/components/SiteCard';
import { Grid, Paper, Typography, Button, AlertColor } from '@mui/material';
import SiteConfigurationStepperDialog from '@/components/SiteConfigurationStepperDialog';
import { useEffect, useState } from 'react';
import {
  useGetSitesLazyQuery,
  useGetNetworksQuery,
  useGetComponentsByUserIdLazyQuery,
} from '@/client/graphql/generated';
import { useAppContext } from '@/context';

const Sites = () => {
  const [open, setOpen] = useState(false);
  const [sitesList, setSitesList] = useState<any[]>([]);
  const [componentsList, setComponentsList] = useState<any[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const { setSnackbarMessage } = useAppContext();

  const handleOpen = () => {
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
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

  const handleFormDataSubmit = (formData: any) => {
    console.log('Form data submitted:', formData);
  };
  console.log('SITES LIST', sitesList);
  return (
    <LoadingWrapper
      radius="small"
      width={'100%'}
      isLoading={isLoading || networkLoading}
      cstyle={{
        backgroundColor: false ? colors.white : 'transparent',
      }}
    >
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
                  onClick={handleOpen}
                >
                  ADD SITE
                </Button>
              </Grid>
            </Grid>

            <Grid container spacing={2}>
              {sitesList.map((site, index) => (
                <Grid item xs={12} md={4} lg={4} key={index}>
                  <SiteCard
                    siteId={site.id}
                    name={site.name}
                    address={site.location}
                    users={site.users || ''}
                    siteStatus={site.isDeactivated}
                    onClickMenu={handleMenuClick}
                    status={{
                      online: false,
                      charging: false,
                      signal: '',
                    }}
                  />
                </Grid>
              ))}
            </Grid>
          </Paper>
        </Grid>
        <SiteConfigurationStepperDialog
          open={open}
          handleClose={handleClose}
          components={componentsList}
          handleFormDataSubmit={handleFormDataSubmit}
        />
      </Grid>
    </LoadingWrapper>
  );
};

export default Sites;
