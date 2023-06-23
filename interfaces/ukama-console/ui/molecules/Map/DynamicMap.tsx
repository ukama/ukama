import { GEO_DATA } from '@/constants';
import { Link, Site } from '@/generated/planning-tool';
import { colors } from '@/styles/theme';
import Leaflet, { LatLngLiteral, LatLngTuple } from 'leaflet';
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
  data: sites,
  links = [],
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

  const getLatLng = (
    sites: Site[],
    siteA: string,
    siteB: string,
  ): LatLngTuple[] => {
    if (sites && sites.length > 0) {
      const locs: LatLngTuple[] = [];
      for (let i = 0; i < sites.length; i++) {
        if (sites[i].id === siteA || sites[i].id === siteB) {
          locs.push([
            parseFloat(sites[i].location.lat),
            parseFloat(sites[i].location.lng),
          ]);
        }
      }
      return locs;
    }
    return [];
  };

  const getGeoData = (links: Link[]): any => {
    for (let i = 0; i < links.length; i++) {
      const geoData = GEO_DATA;
      const locTulp: LatLngTuple[] = getLatLng(
        sites,
        links[i].siteA,
        links[i].siteB,
      );
      geoData.features[0].geometry.coordinates = locTulp;
      console.log(JSON.stringify(geoData));
      return geoData;
    }
    return { type: 'FeatureCollection', features: [] };
  };

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
      {links && links.length && (
        <ReactLeaflet.GeoJSON
          data={getGeoData(links)}
          style={{
            color: colors.primaryMain,
          }}
        />
      )}
      <CustomMarker
        data={sites}
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
