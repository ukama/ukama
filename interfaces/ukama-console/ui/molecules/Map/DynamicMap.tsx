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
  center: LatLngLiteral;
  marker: LatLngLiteral;
  zoom?: number | undefined;
  handleAction: (a: Site) => void;
  setData: Dispatch<SetStateAction<any>>;
  setZoom: Dispatch<SetStateAction<number>>;
  handleDragMarker: (l: LatLngLiteral, id: string) => void;
  handleAddMarker: (l: LatLngLiteral) => void;
}

const Map = ({
  id,
  zoom,
  data,
  cursor,
  marker,
  center,
  setData,
  setZoom,
  children,
  className,
  handleAction,
  handleAddMarker,
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
      center={center}
      touchZoom={false}
      zoomControl={false}
      doubleClickZoom={true}
      scrollWheelZoom={false}
      className={mapClassName}
    >
      {children(ReactLeaflet, Leaflet)}
      <ReactLeaflet.ZoomControl position="bottomright" />
      <CustomMarker
        zoom={zoom}
        data={data}
        marker={marker}
        setData={setData}
        setZoom={setZoom}
        handleAction={handleAction}
        handleAddMarker={handleAddMarker}
        handleDragMarker={handleDragMarker}
      />
    </MapContainer>
  );
};

export default Map;
