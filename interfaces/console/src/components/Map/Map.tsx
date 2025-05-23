/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { LatLngLiteral } from 'leaflet';
import dynamic from 'next/dynamic';
import { Dispatch, SetStateAction } from 'react';

const DynamicMap = dynamic(() => import('./DynamicMap'), {
  ssr: false,
});

const DEFAULT_WIDTH = 600;
const DEFAULT_HEIGHT = 600;

interface IMap {
  id: string;
  data: any[];
  children: any;
  layer: string;
  width?: number;
  height?: number;
  links?: any[];
  linkSites: any;
  isAddLink: boolean;
  isAddSite: boolean;
  className?: string;
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
  handleAddMarker: (l: LatLngLiteral, b: string) => void;
  handleDragMarker: (l: LatLngLiteral, id: string) => void;
}

const MapComponent = ({
  id,
  zoom,
  data,
  links,
  layer,
  center,
  setZoom,
  children,
  linkSites,
  isAddSite,
  isAddLink,
  className,
  selectedLink,
  handleAction,
  coverageLoading,
  handleLinkClick,
  handleAddMarker,
  handleDeleteSite,
  handleDragMarker,
  handleAddLinkToSite,
  handleGenerateAction,
  width = DEFAULT_WIDTH,
  height = DEFAULT_HEIGHT,
}: IMap) => {
  return (
    <div
      style={{
        aspectRatio: width / height,
        cursor: isAddSite ? 'pointer !important' : 'grab !important',
      }}
    >
      <DynamicMap
        id={id}
        zoom={zoom}
        data={data}
        links={links}
        layer={layer}
        center={center}
        setZoom={setZoom}
        cursor={isAddSite}
        isAddLink={isAddLink}
        linkSites={linkSites}
        className={className}
        selectedLink={selectedLink}
        handleAction={handleAction}
        coverageLoading={coverageLoading}
        handleLinkClick={handleLinkClick}
        handleAddMarker={handleAddMarker}
        handleDeleteSite={handleDeleteSite}
        handleDragMarker={handleDragMarker}
        handleAddLinkToSite={handleAddLinkToSite}
        handleGenerateAction={handleGenerateAction}
      >
        {children}
      </DynamicMap>
    </div>
  );
};

export default MapComponent;
