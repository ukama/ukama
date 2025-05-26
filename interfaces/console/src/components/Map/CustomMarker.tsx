/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { colors } from '@/theme';
import Leaflet, { LatLngLiteral, LatLngTuple, Layer, Polyline } from 'leaflet';
import {
  Dispatch,
  SetStateAction,
  useCallback,
  useEffect,
  useState,
} from 'react';
import { Marker, Popup, useMap, useMapEvents } from 'react-leaflet';
import SitePopup from '../SitePopup';

interface ICustomMarker {
  data: any[];
  layer: string;
  links: any[];
  linkSites: any;
  zoom?: number;
  isAddLink: boolean;
  coverageLoading: boolean;
  center: LatLngLiteral | null;
  handleAction: (a: any) => void;
  selectedLink: string | undefined;
  handleLinkClick: (a: string) => void;
  handleDeleteSite: (a: string) => void;
  setZoom: Dispatch<SetStateAction<number>>;
  handleAddLinkToSite: (id: string) => void;
  handleGenerateAction: (a: string, b: any) => void;
  handleAddMarker: (l: LatLngLiteral, b: string) => void;
  handleDragMarker: (l: LatLngLiteral, id: string) => void;
}

interface IMarker {
  id: string;
  lat: number;
  lng: number;
}

interface ILink {
  id: string;
  latlng: LatLngTuple[];
}

const getLatLng = (sites: any[], links: any[]): ILink[] => {
  const data: ILink[] = [];
  if (sites && sites.length > 0) {
    for (let i = 0; i < links.length; i++) {
      const siteA = links[i].siteA;
      const siteB = links[i].siteB;
      const locs: LatLngTuple[] = [];
      sites.forEach(() => {
        // if (site.id === siteA ?? site.id === siteB) {
        //   locs.push([
        //     parseFloat(site.location.lat),
        //     parseFloat(site.location.lng),
        //   ]);
        // }
      });
      if (locs.length > 1)
        data.push({ id: `${siteA}*${siteB}*${links[i].id}`, latlng: locs });
    }
  }
  return data.length > 0 ? data : [];
};

const addRasterData = async (url: string, _: any, __: string) => {
  await fetch(url).then((response) => response.arrayBuffer());
};

const getKey = (lat: string, lng: string) =>
  `lat${lat.replace('.', '_')}lon${lng.replace('.', '_')}`;

const CustomMarker = ({
  data,
  zoom,
  layer,
  links,
  center,
  setZoom,
  handleAction,
  selectedLink,
  handleAddMarker,
  coverageLoading,
  handleDeleteSite,
  handleDragMarker,
  handleLinkClick,
  handleGenerateAction,
}: ICustomMarker) => {
  const map = useMap();
  const [markers, setMarkers] = useState<IMarker[]>([]);
  const [polylines, setPolylines] = useState<Polyline[]>();

  const memoizedHandleLinkClick = useCallback(handleLinkClick, [
    handleLinkClick,
  ]);

  useEffect(() => {
    if (!map) return;

    map.setMaxBounds([
      [84.67351256610522, -174.0234375],
      [-58.995311187950925, 223.2421875],
    ]);

    Leaflet.tileLayer(
      layer === 'satellite'
        ? 'https://tiles.stadiamaps.com/tiles/alidade_smooth/{z}/{x}/{y}{r}.png'
        : 'https://server.arcgisonline.com/arcgis/rest/services/World_Imagery/MapServer/tile/{z}/{y}/{x}',
      {
        subdomains: ['mt0', 'mt1', 'mt2', 'mt3'],
        noWrap: true,
        minZoom: 3,
        maxZoom: 20,
        maxNativeZoom: 20,
      },
    ).addTo(map);
  }, [layer, map]);

  useEffect(() => {
    if (center) map.setView(center, zoom);
  }, [center, map, zoom]);

  useEffect(() => {
    if (!map) return;

    polylines?.forEach((p) => p.removeFrom(map));
    const layers: any = [];
    const latlngs = getLatLng(data, links);
    latlngs.forEach(({ id, latlng }) => {
      const p = Leaflet.polyline(latlng, {
        color: colors.primaryLight,
        weight: 2,
        attribution: id,
      })
        .setStyle({
          interactive: true,
        })
        .addEventListener('click', (e) => {
          memoizedHandleLinkClick(e.target.options.attribution);
        })
        .addTo(map);
      layers.push(p);
    });
    setPolylines(layers);

    const m: any = [];
    data.map((item) => {
      m.push({
        id: item.location.id,
        lat: parseFloat(item.location.lat),
        lng: parseFloat(item.location.lng),
      });
      map.eachLayer((layer: Layer) => {
        if (
          layer.options.attribution ===
          getKey(item.location.lat, item.location.lng)
        )
          map.removeLayer(layer);
      });

      if (item.url)
        addRasterData(
          item.url,
          map,
          getKey(item.location.lat, item.location.lng),
        );
    });
    setMarkers(m);
  }, [data, links, map, memoizedHandleLinkClick, polylines]);

  useEffect(() => {
    polylines?.forEach((p) =>
      p.setStyle({ color: colors.primaryLight, weight: 2 }),
    );
    if (selectedLink) {
      const p = polylines?.find(
        (p) => p.options.attribution?.split('*')[2] === selectedLink,
      );
      if (p) {
        p.setStyle({ color: colors.secondaryMain, weight: 3 });
      }
    }
  }, [selectedLink, polylines]);

  useEffect(() => {
    if (coverageLoading) map.closePopup();
  }, [coverageLoading, map]);

  useMapEvents({
    click: (e) => {
      const { lat, lng } = e.latlng;
      const id = '';
      handleAddMarker({ lat, lng }, id);
      Leaflet.tooltip().openTooltip();
      setMarkers([
        ...markers,
        {
          id: id,
          lat,
          lng,
        },
      ]);
    },
    zoom: (e) => {
      setZoom(e.target.getZoom());
    },
  });

  return (
    <div>
      {data.length > 0 &&
        markers.length > 0 &&
        data.map((item) => {
          const color = colors.black38;
          // (linkSites.siteA === item.id ?? linkSites.siteB === item.id)
          //   ? colors.secondaryMain
          //   : colors.black38;
          const m = markers.find((m) => m.id === item.location.id);
          const svgIcon = Leaflet.divIcon({
            html: `<svg class="MuiSvgIcon-root MuiSvgIcon-fontSizeMedium MuiBox-root css-uqopch" focusable="false" aria-hidden="true" viewBox="0 0 24 24" data-testid="LocationOnIcon" fill=${color}><path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"></path></svg>`,
            className: '',
            iconSize: [28, 28],
            iconAnchor: [14, 26],
          });
          return (
            <Marker
              autoPan
              draggable
              key={item.id}
              icon={svgIcon}
              title={`Population Covered: ${item.populationCovered}`}
              position={{
                lat: m?.lat ?? 0,
                lng: m?.lng ?? 0,
              }}
              attribution={item.id}
              opacity={m?.lat === 0 ? 0 : 1}
              eventHandlers={{
                moveend: (event: any) => {
                  setMarkers([
                    ...markers.filter((m) => m.id !== item.location.id),
                    {
                      id: item.location.id,
                      lat: event.target.getLatLng().lat,
                      lng: event.target.getLatLng().lng,
                    },
                  ]);
                  handleDragMarker(event.target.getLatLng(), item.location.id);
                },
                popupopen: () => {
                  // if (isAddLink) {
                  //   event.target.closePopup();
                  //   const { lat, lng } = event.target.getLatLng();
                  //   const s = data.find(
                  //     (d) =>
                  //       d.location.lat.includes(`${lat}`) &&
                  //       d.location.lng.includes(`${lng}`),
                  //   );
                  //   s?.id && handleAddLinkToSite(s?.id);
                  // }
                },
              }}
            >
              <Popup>
                <SitePopup
                  site={item}
                  handleAction={handleAction}
                  coverageLoading={coverageLoading}
                  handleDeleteSite={handleDeleteSite}
                  handleGenerateAction={handleGenerateAction}
                />
              </Popup>
            </Marker>
          );
        })}
    </div>
  );
};

export default CustomMarker;
