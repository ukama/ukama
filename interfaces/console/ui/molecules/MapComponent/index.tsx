import React from 'react';
import { Box, Stack, Typography, Paper } from '@mui/material';
import { GoogleMap, Marker } from '@react-google-maps/api';
import GroupIcon from '@mui/icons-material/Group';

interface MapProps {
  site: string;
  users: number;
}

const Map: React.FC<MapProps> = ({ site, users }) => {
  return (
    <Box sx={{ position: 'relative', borderRadius: '30%' }}>
      <GoogleMap
        center={{ lat: 37.7749, lng: -122.4194 }}
        zoom={10}
        mapContainerStyle={{ height: '140px', width: '100%' }}
        mapTypeId="terrain"
      >
        <Marker position={{ lat: 37.7749, lng: -122.4194 }} />
      </GoogleMap>
      <Box sx={{ position: 'absolute', bottom: 0, left: 0, ml: 1, mb: 1 }}>
        <Paper sx={{ p: 1, borderRadius: '10px' }}>
          <Stack direction="row" spacing={1} alignItems={'center'}>
            <GroupIcon />
            <Typography variant="body1">{users}</Typography>
          </Stack>
        </Paper>
      </Box>
    </Box>
  );
};

export default Map;
