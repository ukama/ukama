/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { useAppContext } from '@/context';
import colors from '@/theme/colors';
import Leaflet from 'leaflet';
import { useEffect } from 'react';
import * as ReactLeaflet from 'react-leaflet';
interface IMapLayer {
  posix: [string, string];
}

const MapLayer = ({ posix }: IMapLayer) => {
  const map = ReactLeaflet.useMap();
  const { env } = useAppContext();

  useEffect(() => {
    map.setMaxBounds([
      [84.67351256610522, -174.0234375],
      [-58.995311187950925, 223.2421875],
    ]);

    Leaflet.tileLayer(
      `https://api.mapbox.com/styles/v1/mapbox/satellite-v9/tiles/256/{z}/{x}/{y}@2x?access_token=${env.MAP_BOX_TOKEN}`,
      {
        minZoom: 3,
        maxZoom: 15,
        noWrap: true,
        maxNativeZoom: 15,
        attribution: '&copy; <a href="https://www.mapbox.com">Mapbox</a> ',
      },
    ).addTo(map);
  }, [map]);

  useEffect(() => {
    if (posix) {
      map.fitBounds(
        [[Number.parseFloat(posix[0]), Number.parseFloat(posix[1])]],
        {
          maxZoom: 8,
        },
      );
    }
  }, [posix, map]);

  return (
    <div>
      {posix && (
        <ReactLeaflet.Marker
          autoPan
          key={`${posix[0]}-${posix[1]}`}
          icon={Leaflet.divIcon({
            html: `<svg class="MuiSvgIcon-root MuiSvgIcon-fontSizeMedium MuiSvgIcon-root MuiSvgIcon-fontSizeMedium svg-icon css-5zsjn4" focusable="false" aria-hidden="true" viewBox="0 0 24 24" tabindex="-1" title="CellTower" fill="${colors.secondaryMain}"><path d="m7.3 14.7 1.2-1.2c-1-1-1.5-2.3-1.5-3.5 0-1.3.5-2.6 1.5-3.5L7.3 5.3c-1.3 1.3-2 3-2 4.7s.7 3.4 2 4.7M19.1 2.9l-1.2 1.2c1.6 1.6 2.4 3.8 2.4 5.9s-.8 4.3-2.4 5.9l1.2 1.2c2-2 2.9-4.5 2.9-7.1s-1-5.1-2.9-7.1"></path><path d="M6.1 4.1 4.9 2.9C3 4.9 2 7.4 2 10s1 5.1 2.9 7.1l1.2-1.2c-1.6-1.6-2.4-3.8-2.4-5.9s.8-4.3 2.4-5.9m10.6 10.6c1.3-1.3 2-3 2-4.7-.1-1.7-.7-3.4-2-4.7l-1.2 1.2c1 1 1.5 2.3 1.5 3.5 0 1.3-.5 2.6-1.5 3.5zM14.5 10c0-1.38-1.12-2.5-2.5-2.5S9.5 8.62 9.5 10c0 .76.34 1.42.87 1.88L7 22h2l.67-2h4.67l.66 2h2l-3.37-10.12c.53-.46.87-1.12.87-1.88m-4.17 8L12 13l1.67 5z"></path></svg>`,
            className: '',
            iconSize: [28, 28],
            iconAnchor: [14, 26],
          })}
          title={`Site: ${posix[0]}`}
          position={{
            lat: Number.parseFloat(posix[0]),
            lng: Number.parseFloat(posix[1]),
          }}
          attribution={`${posix[0]}-${posix[1]}`}
        />
      )}
    </div>
  );
};

export { MapLayer };
