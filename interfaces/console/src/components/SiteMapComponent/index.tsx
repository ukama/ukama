/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
'use client';

import { useAppContext } from '@/context';
import { isValidLatLng } from '@/utils';
import { LatLngTuple, Map } from 'leaflet';
import 'leaflet-defaulticon-compatibility';
import 'leaflet-defaulticon-compatibility/dist/leaflet-defaulticon-compatibility.webpack.css';
import 'leaflet/dist/leaflet.css';
import { MapContainer, Marker, Popup, TileLayer } from 'react-leaflet';
import { useEffect, useRef } from 'react';
import { Box, Typography } from '@mui/material';
import PeopleIcon from '@mui/icons-material/People';

type MapStyle = 'terrain' | 'satellite' | 'streets' | 'light' | 'dark';

interface SiteMapProps {
  posix: LatLngTuple;
  address: string;
  height?: string;
  mapStyle?: MapStyle;
  showUserCount?: boolean;
  userCount?: number;
}

const SiteMapComponent = ({
  posix,
  address,
  height,
  mapStyle = 'terrain',
  showUserCount = false,
  userCount = 0,
}: SiteMapProps) => {
  const { env } = useAppContext();
  const mapRef = useRef<Map | null>(null);

  const mapStyleUrls = {
    terrain: `https://api.mapbox.com/styles/v1/mapbox/outdoors-v11/tiles/256/{z}/{x}/{y}@2x?access_token=${env.MAP_BOX_TOKEN}`,
    satellite: `https://api.mapbox.com/styles/v1/mapbox/satellite-v9/tiles/256/{z}/{x}/{y}@2x?access_token=${env.MAP_BOX_TOKEN}`,
    streets: `https://api.mapbox.com/styles/v1/mapbox/streets-v11/tiles/256/{z}/{x}/{y}@2x?access_token=${env.MAP_BOX_TOKEN}`,
    light: `https://api.mapbox.com/styles/v1/mapbox/light-v10/tiles/256/{z}/{x}/{y}@2x?access_token=${env.MAP_BOX_TOKEN}`,
    dark: `https://api.mapbox.com/styles/v1/mapbox/dark-v10/tiles/256/{z}/{x}/{y}@2x?access_token=${env.MAP_BOX_TOKEN}`,
  };

  useEffect(() => {
    if (mapRef.current && isValidLatLng(posix)) {
      mapRef.current.setView(posix, 15);
    }
  }, [posix]);

  return (
    <Box
      sx={{
        position: 'relative',
        height: height ? height : '100%',
        width: '100%',
      }}
    >
      <MapContainer
        zoomControl={false}
        preferCanvas={true}
        scrollWheelZoom={false}
        zoom={isValidLatLng(posix) ? 15 : 10}
        center={isValidLatLng(posix) ? posix : undefined}
        bounds={[
          [84.67351256610522, -174.0234375],
          [-58.995311187950925, 223.2421875],
        ]}
        style={{
          height: '100%',
          width: '100%',
          borderRadius: '5px',
        }}
        ref={(map) => {
          if (map) {
            mapRef.current = map;
          }
        }}
      >
        <TileLayer
          url={mapStyleUrls[mapStyle]}
          maxNativeZoom={18}
          minZoom={10}
          maxZoom={18}
        />
        {isValidLatLng(posix) && (
          <Marker position={posix}>
            <Popup>{address || 'Fetching site location...'}</Popup>
          </Marker>
        )}
      </MapContainer>

      {showUserCount && (
        <Box
          sx={{
            position: 'absolute',
            bottom: '16px',
            left: '16px',
            backgroundColor: 'white',
            borderRadius: '28px',
            padding: '8px 16px',
            display: 'flex',
            alignItems: 'center',
            boxShadow: '0px 2px 4px rgba(0, 0, 0, 0.1)',
            zIndex: 1000,
          }}
        >
          <PeopleIcon sx={{ color: '#5F6368', mr: 1 }} />
          <Typography
            variant="body1"
            component="span"
            sx={{ fontWeight: 'bold', color: '#5F6368' }}
          >
            {userCount}
          </Typography>
        </Box>
      )}
    </Box>
  );
};

export default SiteMapComponent;
