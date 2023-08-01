import Leaflet from 'leaflet';
import { useEffect } from 'react';
import * as ReactLeaflet from 'react-leaflet';

const MapLayer = () => {
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
  return <div></div>;
};

export { MapLayer };
