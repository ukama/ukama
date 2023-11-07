/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { NodesLocation } from '@/generated';
import { colors } from '@/styles/theme';
import Leaflet from 'leaflet';
import { useEffect } from 'react';
import * as ReactLeaflet from 'react-leaflet';
interface IMapLayer {
  data: NodesLocation | undefined;
}

const MapLayer = ({ data }: IMapLayer) => {
  const map = ReactLeaflet.useMap();

  useEffect(() => {
    map.setMaxBounds([
      [84.67351256610522, -174.0234375],
      [-58.995311187950925, 223.2421875],
    ]);

    Leaflet.tileLayer(
      'https://tiles.stadiamaps.com/tiles/alidade_smooth/{z}/{x}/{y}{r}.png',
      {
        noWrap: true,
        minZoom: 3,
        maxZoom: 20,
        maxNativeZoom: 20,
      },
    ).addTo(map);
  }, []);

  useEffect(() => {
    if (data && data.nodes.length > 0) {
      map.fitBounds(
        data.nodes.map((item) => [
          parseFloat(item?.lat) || 0,
          parseFloat(item?.lng) || 0,
        ]),
        {
          maxZoom: 16,
        },
      );
    }
  }, [data]);

  return (
    <div>
      {data &&
        data.nodes.length > 0 &&
        data.nodes.map((item) => {
          const svgIcon = Leaflet.divIcon({
            html: `<svg class="MuiSvgIcon-root MuiSvgIcon-fontSizeMedium MuiBox-root css-uqopch" focusable="false" aria-hidden="true" viewBox="0 0 24 24" data-testid="LocationOnIcon" fill=${colors.secondaryMain}><path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"></path></svg>`,
            className: '',
            iconSize: [28, 28],
            iconAnchor: [14, 26],
          });
          return (
            <ReactLeaflet.Marker
              autoPan
              draggable
              key={item.id}
              icon={svgIcon}
              title={`Node: ${item.id}`}
              position={{
                lat: parseFloat(item?.lat) || 0,
                lng: parseFloat(item?.lng) || 0,
              }}
              attribution={item.id}
            />
          );
        })}
    </div>
  );
};

export { MapLayer };
