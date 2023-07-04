import { Link, Site } from '@/generated/planning-tool';
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
  links?: Link[];
  children: any;
  linkSites: any;
  className?: string;
  isAddLink: boolean;
  center: LatLngLiteral;
  zoom?: number | undefined;
  handleAction: (a: Site) => void;
  handleLinkClick: (a: string) => void;
  handleDeleteSite: (a: string) => void;
  handleAddLinkToSite: (id: string) => void;
  setZoom: Dispatch<SetStateAction<number>>;
  handleDragMarker: (l: LatLngLiteral, id: string) => void;
  handleAddMarker: (l: LatLngLiteral, b: string) => void;
}

const Map = ({
  id,
  zoom,
  center,
  cursor,
  setZoom,
  children,
  linkSites,
  className,
  isAddLink,
  links = [],
  data: sites,
  handleAction,
  handleLinkClick,
  handleAddMarker,
  handleDeleteSite,
  handleDragMarker,
  handleAddLinkToSite,
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
      touchZoom={true}
      zoomControl={false}
      doubleClickZoom={true}
      scrollWheelZoom={true}
      className={mapClassName}
      attributionControl={false}
    >
      {children(ReactLeaflet, Leaflet)}
      <ReactLeaflet.ZoomControl position="bottomright" />
      <ReactLeaflet.LayersControl position="bottomleft"  />
      <CustomMarker
        data={sites}
        zoom={zoom}
        links={links}
        center={center}
        setZoom={setZoom}
        linkSites={linkSites}
        isAddLink={isAddLink}
        handleAction={handleAction}
        handleLinkClick={handleLinkClick}
        handleAddMarker={handleAddMarker}
        handleDeleteSite={handleDeleteSite}
        handleDragMarker={handleDragMarker}
        handleAddLinkToSite={handleAddLinkToSite}
      />
    </MapContainer>
  );
};

export default Map;
