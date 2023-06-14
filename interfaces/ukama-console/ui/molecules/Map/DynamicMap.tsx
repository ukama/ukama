import Leaflet from 'leaflet';
import markerIcon from 'leaflet/dist/images/marker-icon.png';
import markerShadow from 'leaflet/dist/images/marker-shadow.png';
import 'leaflet/dist/leaflet.css';
import { useEffect } from 'react';
import * as ReactLeaflet from 'react-leaflet';
import styles from './Map.module.css';

const { MapContainer } = ReactLeaflet;

interface IMap {
  center?: any;
  children: any;
  zoom?: number;
  className?: string;
}

const ICON = {
  iconSize: [25, 41],
  iconAnchor: [10, 41],
  popupAnchor: [2, -40],
  iconUrl: markerIcon.src,
  iconRetinaUrl: markerIcon.src,
  shadowUrl: markerShadow.src,
};

const Map = ({ zoom, center, children, className }: IMap) => {
  let mapClassName = styles.map;

  if (className) {
    mapClassName = `${mapClassName} ${className}`;
  }

  useEffect(() => {
    (async function init() {
      Leaflet.Icon.Default.mergeOptions(ICON);
      Leaflet.Control.Zoom.prototype.options.position = 'bottomright';
    })();
  }, []);

  return (
    <MapContainer
      zoom={zoom}
      center={center}
      zoomControl={false}
      className={mapClassName}
    >
      {children(ReactLeaflet, Leaflet)}
      <ReactLeaflet.ZoomControl position="bottomright" />
    </MapContainer>
  );
};

export default Map;
