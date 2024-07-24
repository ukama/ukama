import React, { useState, useEffect } from 'react';
import { Grid, Paper, Stack, Typography, Skeleton, Alert } from '@mui/material';
import { LatLngTuple } from 'leaflet';
import { SiteDto } from '@/client/graphql/generated';
import dynamic from 'next/dynamic';

const SiteMapComponent = dynamic(() => import('../SiteMapComponent'), {
  ssr: false,
});

interface SiteInfoProps {
  selectedSite: SiteDto;
}

const SiteInfo: React.FC<SiteInfoProps> = ({ selectedSite }) => {
  const [isMapReady, setIsMapReady] = useState(false);
  const [mapError, setMapError] = useState<string | null>(null);

  const position: LatLngTuple = [
    selectedSite.latitude ?? 0,
    selectedSite.longitude ?? 0,
  ];

  const isValidPosition = (lat: number, lng: number) =>
    !isNaN(lat) && !isNaN(lng) && lat !== 0 && lng !== 0;

  useEffect(() => {
    try {
      if (isValidPosition(position[0], position[1])) {
        setIsMapReady(true);
      } else {
        setMapError('Invalid coordinates.');
      }
    } catch (error) {
      setMapError('Failed to load map.');
    }
  }, [position]);

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
        {mapError ? (
          <Alert severity="error">{mapError}</Alert>
        ) : isMapReady ? (
          <div style={{ flex: 1 }}>
            <SiteMapComponent
              posix={position}
              address={selectedSite.location || 'No address provided'}
            />
          </div>
        ) : (
          <Skeleton variant="rectangular" height="100%" />
        )}
      </Grid>
    </Grid>
  );
};

export default SiteInfo;
