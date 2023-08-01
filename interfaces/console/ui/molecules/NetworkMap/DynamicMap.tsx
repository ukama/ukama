import styles from '@/styles/Map.module.css';
import Leaflet from 'leaflet';
import markerIcon from 'leaflet/dist/images/marker-icon.png';
import markerShadow from 'leaflet/dist/images/marker-shadow.png';
import 'leaflet/dist/leaflet.css';
import { useEffect } from 'react';
import * as ReactLeaflet from 'react-leaflet';
import { MapLayer } from './MapLayer';

const { MapContainer } = ReactLeaflet;

const ICON = {
  iconUrl: markerIcon.src,
  iconRetinaUrl: markerIcon.src,
  shadowUrl: markerShadow.src,
};
interface IMap {
  id: string;
  zoom: number;
  children: any;
  className?: string;
}

const Map = ({ id, zoom, children }: IMap) => {
  let mapClassName = styles.map;
  let mapContainer = styles['leaflet-container'];

  useEffect(() => {
    (async function init() {
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
      attributionControl={false}
      center={[37.7780627, -121.9822475]}
      className={`${mapClassName} ${mapContainer}`}
    >
      {children(ReactLeaflet, Leaflet)}
      <ReactLeaflet.ZoomControl position="bottomright" />
      <MapLayer />
    </MapContainer>
  );
};

export default Map;
