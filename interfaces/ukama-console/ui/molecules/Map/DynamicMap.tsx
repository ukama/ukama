import Leaflet from 'leaflet';
import markerIcon2x from 'leaflet/dist/images/marker-icon-2x.png';
import markerIcon from 'leaflet/dist/images/marker-icon.png';
import markerShadow from 'leaflet/dist/images/marker-shadow.png';
import 'leaflet/dist/leaflet.css';
import { useEffect } from 'react';
import * as ReactLeaflet from 'react-leaflet';
import styles from './Map.module.css';

const { MapContainer } = ReactLeaflet;

interface IMap {
  className?: string;
  center?: any;
  children: any;
  zoom?: number;
}

const Map = ({ children, className, zoom, center }: IMap) => {
  let mapClassName = styles.map;

  if (className) {
    mapClassName = `${mapClassName} ${className}`;
  }

  useEffect(() => {
    (async function init() {
      // delete Leaflet.Icon.Default.prototype._getIconUrl;
      Leaflet.Icon.Default.mergeOptions({
        iconUrl: markerIcon.src,
        iconRetinaUrl: markerIcon2x.src,
        shadowUrl: markerShadow.src,
      });
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
