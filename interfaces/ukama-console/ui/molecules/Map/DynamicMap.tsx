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
  onMapClick: Function;
}

const ICON = {
  iconUrl: markerIcon.src,
  iconRetinaUrl: markerIcon.src,
  shadowUrl: markerShadow.src,
};

function MyComponent({ saveMarkers }: any) {
  ReactLeaflet.useMapEvents({
    click: (e) => {
      const { lat, lng } = e.latlng;
      saveMarkers([lat, lng]);
    },
  });
  return null;
}

const Map = ({ zoom, center, children, className, onMapClick }: IMap) => {
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
      <MyComponent saveMarkers={onMapClick} />
    </MapContainer>
  );
};

export default Map;
