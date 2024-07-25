import React, { useState, useEffect } from 'react';
import { Grid, Paper, Stack, Typography, Skeleton } from '@mui/material';
import { SiteDto } from '@/client/graphql/generated';
import dynamic from 'next/dynamic';
import { useAppContext } from '@/context';

const SiteMapComponent = dynamic(() => import('../SiteMapComponent'), {
  ssr: false,
});

interface SiteInfoProps {
  selectedSite: SiteDto;
}

const SiteInfo: React.FC<SiteInfoProps> = ({ selectedSite }) => {
  const { setSnackbarMessage, setSelectedDefaultSite } = useAppContext();
  const [address, setAddress] = useState('');

  const fetchAddress = async (lat: number, lng: number) => {
    return await fetch(
      `https://nominatim.openstreetmap.org/reverse?format=json&lat=${lat || 37.7749}&lon=${lng || -122.4194}`,
      {
        cache: 'force-cache',
      },
    )
      .then((res) => res.json())
      .then((data) => data.display_name)
      .catch(() => 'Location not found');
  };

  useEffect(() => {
    const handleFetchAddress = async (lat: number, lng: number) => {
      setSnackbarMessage({
        id: 'fetching-address',
        type: 'success',
        show: true,
        message: 'Fetching address with coordinates',
      });
      const addressData = await fetchAddress(lat, lng);
      setAddress(addressData);
    };
    setSelectedDefaultSite(selectedSite.name);

    if (selectedSite) {
      handleFetchAddress(selectedSite.latitude, selectedSite.longitude);
    }
  }, [selectedSite, setSnackbarMessage]);

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
          </Stack>
        </Paper>
      </Grid>
      <Grid item xs={12} md={8} sx={{ display: 'flex', alignItems: 'stretch' }}>
        {address ? (
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
