/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
'use client';

import 'leaflet/dist/leaflet.css';
import 'leaflet-defaulticon-compatibility/dist/leaflet-defaulticon-compatibility.webpack.css';
import { LatLngExpression, LatLngTuple } from 'leaflet';
import 'leaflet-defaulticon-compatibility';
import { useEffect, useState } from 'react';
import { MapContainer, Marker, Popup, TileLayer } from 'react-leaflet';

interface SiteMapProps {
  posix: LatLngExpression | LatLngTuple;
  onAddressChange: (address: string) => void;
}

const SiteMapComponent = ({ posix, onAddressChange }: SiteMapProps) => {
  const [address, setAddress] = useState<string>('');

  useEffect(() => {
    const fetchAddress = async () => {
      const [lat, lng] = posix as LatLngTuple;
      const response = await fetch(
        `https://nominatim.openstreetmap.org/reverse?format=json&lat=${lat || 37.7749}&lon=${lng || -122.4194}`,
      );
      const data = await response.json();
      setAddress(data.display_name);
      onAddressChange(data.display_name);
    };

    fetchAddress();
  }, [posix, onAddressChange]);

  return (
    <MapContainer
      preferCanvas={true}
      center={posix}
      zoom={11}
      scrollWheelZoom={true}
      style={{ height: '200px', width: '100%' }}
    >
      <TileLayer
        attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
        url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
      />
      <Marker position={posix}>
        <Popup>{address || 'Fetching site location...'}</Popup>
      </Marker>
    </MapContainer>
  );
};

export default SiteMapComponent;
