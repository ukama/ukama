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
  isAddLink: boolean;
  zoom?: number | undefined;
  center: LatLngLiteral | null;
  handleAction: (a: Site) => void;
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

const getLatLng = (sites: Site[], links: Link[]): LatLngTuple[] => {
  if (sites && sites.length > 0) {
    const locs: LatLngTuple[] = [];
    for (let i = 0; i < links.length; i++) {
      const siteA = links[i].siteA;
      const siteB = links[i].siteB;
      sites.forEach((site) => {
        if (site.id === siteA || site.id === siteB) {
          const l: LatLngTuple = [
            parseFloat(site.location.lat),
            parseFloat(site.location.lng),
          ];
          if (!locs.toString().includes(l.toString()))
            locs.push([
              parseFloat(site.location.lat),
              parseFloat(site.location.lng),
            ]);
        }
      });
    }
    return locs.length === 1 ? [] : locs;
  }
  return [];
};

const CustomMarker = ({
  data,
  zoom,
  links,
  center,
  setZoom,
  isAddLink,
  handleAction,
  handleAddMarker,
  handleDeleteSite,
  handleDragMarker,
  handleAddLinkToSite,
}: ICustomMarker) => {
  const map = useMap();
  const [markers, setMarkers] = useState<IMarker[]>([]);
  const [polylines, setPolylines] = useState<Polyline>();

  useEffect(() => {
    map.setMaxBounds([
      [84.67351256610522, -174.0234375],
      [-58.995311187950925, 223.2421875],
    ]);
    Leaflet.tileLayer(
      'https://tiles.stadiamaps.com/tiles/alidade_smooth/{z}/{x}/{y}{r}.png',
      {
        noWrap: true,
        minZoom: 3,
        maxZoom: 20,
        maxNativeZoom: 20,
      },
    ).addTo(map);
  }, []);

  useEffect(() => {
    if (center) map.setView(center, zoom);
  }, [center]);

  useEffect(() => {
    polylines?.removeFrom(map);
    var latlngs = getLatLng(data, links);
    const linesLayer = Leaflet.polyline(latlngs, {
      color: colors.primaryMain,
    }).addTo(map);
    setPolylines(linesLayer);

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
          const m = markers.find((m) => m.id === item.location.id);
          return (
            <Marker
              autoPan
              draggable
              key={item.id}
              title={item.name}
              position={{
                lat: m?.lat || 0,
                lng: m?.lng || 0,
              }}
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
