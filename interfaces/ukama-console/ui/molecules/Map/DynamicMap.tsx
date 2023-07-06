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
  layer: string;
  links?: Link[];
  children: any;
  linkSites: any;
  className?: string;
  isAddLink: boolean;
  center: LatLngLiteral;
  coverageLoading: boolean;
  zoom?: number | undefined;
  handleAction: (a: Site) => void;
  selectedLink: string | undefined;
  handleLinkClick: (a: string) => void;
  handleDeleteSite: (a: string) => void;
  handleAddLinkToSite: (id: string) => void;
  setZoom: Dispatch<SetStateAction<number>>;
  handleGenerateAction: (a: string, b: Site) => void;
  handleDragMarker: (l: LatLngLiteral, id: string) => void;
  handleAddMarker: (l: LatLngLiteral, b: string) => void;
}

const Map = ({
  id,
  zoom,
  layer,
  center,
  cursor,
  setZoom,
  children,
  linkSites,
  className,
  isAddLink,
  links = [],
  data: sites,
  selectedLink,
  handleAction,
  handleLinkClick,
  handleAddMarker,
  coverageLoading,
  handleDeleteSite,
  handleDragMarker,
  handleAddLinkToSite,
  handleGenerateAction,
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
      <CustomMarker
        zoom={zoom}
        data={sites}
        layer={layer}
        links={links}
        center={center}
        setZoom={setZoom}
        linkSites={linkSites}
        isAddLink={isAddLink}
        selectedLink={selectedLink}
        handleAction={handleAction}
        coverageLoading={coverageLoading}
        handleLinkClick={handleLinkClick}
        handleAddMarker={handleAddMarker}
        handleDeleteSite={handleDeleteSite}
        handleDragMarker={handleDragMarker}
        handleAddLinkToSite={handleAddLinkToSite}
        handleGenerateAction={handleGenerateAction}
      />
    </MapContainer>
  );
};

export default Map;
