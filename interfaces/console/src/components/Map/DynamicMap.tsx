/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import styles from '@/styles/Map.module.css';
import Leaflet, { LatLngLiteral } from 'leaflet';
import markerIcon from 'leaflet/dist/images/marker-icon.png';
import markerShadow from 'leaflet/dist/images/marker-shadow.png';
import 'leaflet/dist/leaflet.css';
import { Dispatch, SetStateAction, useEffect } from 'react';
import * as ReactLeaflet from 'react-leaflet';
import CustomMarker from './CustomMarker';

const { MapContainer } = ReactLeaflet;

const ICON = {
  iconUrl: markerIcon.src,
  iconRetinaUrl: markerIcon.src,
  shadowUrl: markerShadow.src,
};
interface IMap {
  data: any[];
  id: string;
  cursor: any;
  layer: string;
  links?: any[];
  children: any;
  linkSites: any;
  className?: string;
  isAddLink: boolean;
  center: LatLngLiteral;
  coverageLoading: boolean;
  zoom?: number | undefined;
  handleAction: (a: any) => void;
  selectedLink: string | undefined;
  handleLinkClick: (a: string) => void;
  handleDeleteSite: (a: string) => void;
  handleAddLinkToSite: (id: string) => void;
  setZoom: Dispatch<SetStateAction<number>>;
  handleGenerateAction: (a: string, b: any) => void;
  handleDragMarker: (l: LatLngLiteral, id: string) => void;
  handleAddMarker: (l: LatLngLiteral, b: string) => void;
}

const UkamaMap = ({
  id,
  zoom,
  layer,
  center,
  cursor,
  setZoom,
  children,
  linkSites,
  className,
  isAddLink,
  links = [],
  data: sites,
  selectedLink,
  handleAction,
  handleLinkClick,
  handleAddMarker,
  coverageLoading,
  handleDeleteSite,
  handleDragMarker,
  handleAddLinkToSite,
  handleGenerateAction,
}: IMap) => {
  let mapClassName = styles.map;

  if (className) {
    mapClassName = `${mapClassName} ${className} ${cursor ? styles.cursor : ''}`;
  }

  useEffect(() => {
    (function init() {
      Leaflet.Icon.Default.mergeOptions(ICON);
      Leaflet.Control.Zoom.prototype.options.position = 'bottomright';
    })();
  }, []);

  return (
    <MapContainer
      id={id}
      zoom={zoom}
      touchZoom={true}
      zoomControl={false}
      doubleClickZoom={true}
      scrollWheelZoom={true}
      className={mapClassName}
      attributionControl={false}
    >
      {children(ReactLeaflet, Leaflet)}
      <ReactLeaflet.ZoomControl position="bottomright" />
      <CustomMarker
        zoom={zoom}
        data={sites}
        layer={layer}
        links={links}
        center={center}
        setZoom={setZoom}
        linkSites={linkSites}
        isAddLink={isAddLink}
        selectedLink={selectedLink}
        handleAction={handleAction}
        coverageLoading={coverageLoading}
        handleLinkClick={handleLinkClick}
        handleAddMarker={handleAddMarker}
        handleDeleteSite={handleDeleteSite}
        handleDragMarker={handleDragMarker}
        handleAddLinkToSite={handleAddLinkToSite}
        handleGenerateAction={handleGenerateAction}
      />
    </MapContainer>
  );
};

export default UkamaMap;
