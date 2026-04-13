/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import styles from '@/styles/Map.module.css';
import PeopleIcon from '@mui/icons-material/People';
import { Box, Typography } from '@mui/material';
import Leaflet from 'leaflet';
import markerIcon from 'leaflet/dist/images/marker-icon.png';
import markerShadow from 'leaflet/dist/images/marker-shadow.png';
import 'leaflet/dist/leaflet.css';
import { useEffect } from 'react';
import * as ReactLeaflet from 'react-leaflet';
import { MapLayer } from './MapLayer';

const { MapContainer } = ReactLeaflet;

const ICON = {
  iconUrl: markerIcon.src,
  iconRetinaUrl: markerIcon.src,
  shadowUrl: markerShadow.src,
};
interface IMap {
  id: string;
  zoom?: number;
  address?: string;
  height?: string;
  mapStyle?: 'terrain' | 'satellite' | 'streets' | 'light' | 'dark';
  showUserCount?: boolean;
  userCount?: number;
  posix: [string, string];
}

const SiteMap = ({ id, showUserCount = false, userCount, posix }: IMap) => {
  const mapClassName = styles.map;
  const mapContainer = styles['leaflet-container'];

  useEffect(() => {
    (function init() {
      Leaflet.Icon.Default.mergeOptions(ICON);
      Leaflet.Control.Zoom.prototype.options.position = 'bottomright';
    })();
  }, []);

  return (
    <MapContainer
      id={id}
      zoom={8}
      touchZoom={false}
      zoomControl={false}
      doubleClickZoom={false}
      scrollWheelZoom={false}
      attributionControl={false}
      center={[37.7780627, -121.9822475]}
      className={`${mapClassName} ${mapContainer}`}
    >
      <ReactLeaflet.ZoomControl position="bottomright" />
      <MapLayer posix={posix} />
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
    </MapContainer>
  );
};

export default SiteMap;
