import { Link, Site } from '@/generated/planning-tool';
import { colors } from '@/styles/theme';
import Leaflet, { LatLngLiteral, LatLngTuple, Polyline } from 'leaflet';
import { Dispatch, SetStateAction, useEffect, useState } from 'react';
import { Marker, Popup, useMap, useMapEvents } from 'react-leaflet';
import { v4 as uuidv4 } from 'uuid';
import SitePopup from '../SitePopup';
interface ICustomMarker {
  data: Site[];
  links: Link[];
  linkSites: any;
  isAddLink: boolean;
  zoom?: number | undefined;
  center: LatLngLiteral | null;
  handleAction: (a: Site) => void;
  handleLinkClick: (a: string) => void;
  handleDeleteSite: (a: string) => void;
  setZoom: Dispatch<SetStateAction<number>>;
  handleAddLinkToSite: (id: string) => void;
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

const getLatLng = (sites: Site[], links: Link[]): ILink[] => {
  const data: ILink[] = [];
  if (sites && sites.length > 0) {
    for (let i = 0; i < links.length; i++) {
      const siteA = links[i].siteA;
      const siteB = links[i].siteB;
      const locs: LatLngTuple[] = [];
      sites.forEach((site) => {
        if (site.id === siteA || site.id === siteB) {
          // const l: LatLngTuple = [
          //   parseFloat(site.location.lat),
          //   parseFloat(site.location.lng),
          // ];
          // if (!locs.toString().includes(l.toString()))
          locs.push([
            parseFloat(site.location.lat),
            parseFloat(site.location.lng),
          ]);
        }
      });
      if (locs.length > 1)
        data.push({ id: `${siteA}-${siteB}-${links[i].id}`, latlng: locs });
    }
  }
  return data.length > 0 ? data : [];
};

const CustomMarker = ({
  data,
  zoom,
  links,
  center,
  setZoom,
  isAddLink,
  linkSites,
  handleAction,
  handleAddMarker,
  handleDeleteSite,
  handleDragMarker,
  handleLinkClick,
  handleAddLinkToSite,
}: ICustomMarker) => {
  const map = useMap();
  const [markers, setMarkers] = useState<IMarker[]>([]);
  const [polylines, setPolylines] = useState<Polyline[]>();

  useEffect(() => {
    map.setMaxBounds([
      [84.67351256610522, -174.0234375],
      [-58.995311187950925, 223.2421875],
    ]);

    Leaflet.tileLayer(
      // 'https://tiles.stadiamaps.com/tiles/alidade_smooth/{z}/{x}/{y}{r}.png',//grey
      // 'http://{s}.google.com/vt/lyrs=p&x={x}&y={y}&z={z}', //Terain
      // 'http://{s}.google.com/vt/lyrs=s&x={x}&y={y}&z={z}', //Satellite
      'http://{s}.google.com/vt/lyrs=s,h&x={x}&y={y}&z={z}', //Hybrid
      // 'https://server.arcgisonline.com/arcgis/rest/services/World_Imagery/MapServer/tile/{z}/{y}/{x}',
      {
        subdomains: ['mt0', 'mt1', 'mt2', 'mt3'],
        noWrap: true,
        minZoom: 3,
        maxZoom: 20,
        maxNativeZoom: 20,
      },
    ).addTo(map);
    Leaflet.imageOverlay(
      'one.png',
      [
        [-5.056886378670788, 26.404207],
        [-4.517721, 26.94504066615385],
      ],
      {},
    ).addTo(map);
  }, []);

  useEffect(() => {
    if (center) map.setView(center, zoom);
  }, [center]);

  useEffect(() => {
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
          e.target.setStyle({ color: colors.secondaryMain, weight: 3 });
          handleLinkClick(e.target.options.attribution);
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
    });
    setMarkers(m);
  }, [data]);

  useMapEvents({
    click: (e) => {
      const { lat, lng } = e.latlng;
      const id = uuidv4();
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
          const color =
            linkSites.siteA === item.id || linkSites.siteB === item.id
              ? colors.secondaryMain
              : colors.primaryMain;
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
              title={item.name}
              position={{
                lat: m?.lat || 0,
                lng: m?.lng || 0,
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
                popupopen: (event: any) => {
                  if (isAddLink) {
                    event.target.closePopup();
                    const { lat, lng } = event.target.getLatLng();
                    const s = data.find(
                      (d) =>
                        d.location.lat.includes(`${lat}`) &&
                        d.location.lng.includes(`${lng}`),
                    );
                    s?.id && handleAddLinkToSite(s?.id);
                  }
                },
              }}
            >
              <Popup>
                <SitePopup
                  site={item}
                  handleAction={handleAction}
                  handleDeleteSite={handleDeleteSite}
                />
              </Popup>
            </Marker>
          );
        })}
    </div>
  );
};

export default CustomMarker;
