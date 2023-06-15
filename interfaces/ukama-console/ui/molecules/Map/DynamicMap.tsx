import Leaflet from 'leaflet';
import markerIcon from 'leaflet/dist/images/marker-icon.png';
import markerShadow from 'leaflet/dist/images/marker-shadow.png';
import 'leaflet/dist/leaflet.css';
import { useEffect } from 'react';
import * as ReactLeaflet from 'react-leaflet';
import styles from './Map.module.css';

const { MapContainer } = ReactLeaflet;

const ICON = {
  iconUrl: markerIcon.src,
  iconRetinaUrl: markerIcon.src,
  shadowUrl: markerShadow.src,
};
interface IMap {
  id: string;
  center?: any;
  zoom?: number;
  children: any;
  cursor: any;
  className?: string;
  onMapClick: Function;
}

interface ICustomMarker {
  saveMarkers: Function;
}

function CustomMarker({ saveMarkers }: ICustomMarker) {
  ReactLeaflet.useMapEvents({
    click: (e) => {
      const { lat, lng } = e.latlng;
      saveMarkers([lat, lng]);
      Leaflet.tooltip().openTooltip();
    },
  });
  return null;
}

const Map = ({
  id,
  zoom,
  center,
  cursor,
  children,
  className,
  onMapClick,
}: IMap) => {
  let mapClassName = styles.map;

  if (className) {
    mapClassName = `${mapClassName} ${className} ${
      cursor ? styles.cursor : ''
    }`;
  }

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
      center={center}
      touchZoom={false}
      zoomControl={false}
      scrollWheelZoom={false}
      className={mapClassName}
    >
      {children(ReactLeaflet, Leaflet)}
      <ReactLeaflet.ZoomControl position="bottomright" />
      <CustomMarker saveMarkers={onMapClick} />
    </MapContainer>
  );
};

export default Map;
