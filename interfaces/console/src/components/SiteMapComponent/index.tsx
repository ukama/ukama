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
import { LatLngTuple } from 'leaflet';
import 'leaflet-defaulticon-compatibility';
import 'leaflet-defaulticon-compatibility/dist/leaflet-defaulticon-compatibility.webpack.css';
import 'leaflet/dist/leaflet.css';
import { MapContainer, Marker, Popup, TileLayer } from 'react-leaflet';

interface SiteMapProps {
  posix: LatLngTuple;
  address: string;
  height?: string;
}

const SiteMapComponent = ({ posix, address, height }: SiteMapProps) => {
  const { env } = useAppContext();
  return (
    <MapContainer
      zoomControl={false}
      preferCanvas={true}
      scrollWheelZoom={false}
      zoom={isValidLatLng(posix) ? 13 : 10}
      center={isValidLatLng(posix) ? posix : undefined}
      bounds={[
        [84.67351256610522, -174.0234375],
        [-58.995311187950925, 223.2421875],
      ]}
      style={{
        height: height ? `${height} ` : '100%',
        width: '100%',
        borderRadius: '5px',
      }}
    >
      <TileLayer
        url={`https://api.mapbox.com/styles/v1/salman-ukama/clxu9ic7z00ua01qr7hb93d2o/tiles/256/{z}/{x}/{y}@2x?access_token=${env.MAP_BOX_TOKEN}`}
        accessToken={env.MAP_BOX_TOKEN}
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
  );
};

export default SiteMapComponent;
