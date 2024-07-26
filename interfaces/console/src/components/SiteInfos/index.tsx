import React, { useEffect } from 'react';
import { Grid, Paper, Stack, Typography, Skeleton } from '@mui/material';
import { SiteDto } from '@/client/graphql/generated';
import dynamic from 'next/dynamic';
import { useAppContext } from '@/context';
import { useFetchAddress } from '@/utils/useFetchAddress';

const SiteMapComponent = dynamic(() => import('../SiteMapComponent'), {
  ssr: false,
});

interface SiteInfoProps {
  selectedSite: SiteDto;
}

const SiteInfo: React.FC<SiteInfoProps> = ({ selectedSite }) => {
  const { setSnackbarMessage, setSelectedDefaultSite } = useAppContext();
  const { address, isLoading, error, fetchAddress } = useFetchAddress();

  useEffect(() => {
    const handleFetchAddress = async () => {
      setSnackbarMessage({
        id: 'fetching-address',
        type: 'info',
        show: true,
        message: 'Fetching address with coordinates',
      });
      await fetchAddress(selectedSite.latitude, selectedSite.longitude);
    };

    setSelectedDefaultSite(selectedSite.name);

    if (selectedSite) {
      handleFetchAddress();
    }
  }, [selectedSite, setSnackbarMessage, fetchAddress, setSelectedDefaultSite]);

  useEffect(() => {
    if (error) {
      setSnackbarMessage({
        id: 'error-fetching-address',
        type: 'error',
        show: true,
        message: 'Error fetching address from coordinates',
      });
    }
  }, [error, setSnackbarMessage]);

  return (
    <Grid container spacing={2} sx={{ height: '100%' }}>
      <Grid item xs={12} md={4} sx={{ display: 'flex', alignItems: 'stretch' }}>
        <Paper elevation={3} sx={{ p: 2, flex: 1 }}>
          <Typography variant="h6">Site Information</Typography>
          <Stack direction="column" spacing={2}>
            <Typography variant="subtitle1">
              Date created:{' '}
              {new Date(selectedSite.createdAt).toLocaleDateString()}
            </Typography>
            <Typography variant="subtitle1">
              Location: {selectedSite.location || 'N/A'}
            </Typography>
            <Typography variant="subtitle1">
              Coordinates: {selectedSite.latitude || 'N/A'},{' '}
              {selectedSite.longitude || 'N/A'}
            </Typography>
            <Typography variant="subtitle1">
              Address:{' '}
              {isLoading
                ? 'Loading...'
                : error
                  ? 'Error fetching address'
                  : address || 'N/A'}
            </Typography>
          </Stack>
        </Paper>
      </Grid>
      <Grid item xs={12} md={8} sx={{ display: 'flex', alignItems: 'stretch' }}>
        {isLoading ? (
          <Skeleton variant="rectangular" height="100%" />
        ) : error ? (
          <Typography variant="body1" color="error">
            Error loading map
          </Typography>
        ) : address ? (
          <SiteMapComponent
            posix={[selectedSite.latitude, selectedSite.longitude]}
            address={address}
          />
        ) : (
          <Skeleton variant="rectangular" height="100%" />
        )}
      </Grid>
    </Grid>
  );
};

export default SiteInfo;
