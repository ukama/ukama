import { Site } from '@/generated/planning-tool';
import Leaflet, { LatLngLiteral } from 'leaflet';
import markerIcon from 'leaflet/dist/images/marker-icon.png';
import markerShadow from 'leaflet/dist/images/marker-shadow.png';
import 'leaflet/dist/leaflet.css';
import { Dispatch, SetStateAction, useEffect } from 'react';
import * as ReactLeaflet from 'react-leaflet';
import CustomMarker from './CustomMarker';
import styles from './Map.module.css';

const { MapContainer } = ReactLeaflet;

const ICON = {
  iconUrl: markerIcon.src,
  iconRetinaUrl: markerIcon.src,
  shadowUrl: markerShadow.src,
};
interface IMap {
  data: Site[];
  id: string;
  cursor: any;
  children: any;
  className?: string;
  zoom?: number | undefined;
  center: LatLngLiteral;
  handleAction: (a: Site) => void;
  handleDeleteSite: (a: string) => void;
  setZoom: Dispatch<SetStateAction<number>>;
  handleDragMarker: (l: LatLngLiteral, id: string) => void;
  handleAddMarker: (l: LatLngLiteral, b: string) => void;
}

const Map = ({
  id,
  zoom,
  data,
  center,
  cursor,
  setZoom,
  children,
  className,
  handleAction,
  handleAddMarker,
  handleDeleteSite,
  handleDragMarker,
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
      touchZoom={false}
      zoomControl={false}
      doubleClickZoom={true}
      scrollWheelZoom={false}
      className={mapClassName}
    >
      {children(ReactLeaflet, Leaflet)}
      <ReactLeaflet.ZoomControl position="bottomright" />
      <CustomMarker
        data={data}
        zoom={zoom}
        center={center}
        setZoom={setZoom}
        handleAction={handleAction}
        handleAddMarker={handleAddMarker}
        handleDeleteSite={handleDeleteSite}
        handleDragMarker={handleDragMarker}
      />
    </MapContainer>
  );
};

export default Map;
