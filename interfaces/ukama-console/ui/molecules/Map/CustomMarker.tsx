import { Site } from '@/generated/planning-tool';
import Leaflet, { LatLngLiteral } from 'leaflet';
import { Dispatch, SetStateAction, useEffect, useState } from 'react';
import { Marker, Popup, useMap, useMapEvents } from 'react-leaflet';
import { v4 as uuidv4 } from 'uuid';
import SitePopup from '../SitePopup';

interface ICustomMarker {
  data: Site[];
  zoom?: number | undefined;
  center: LatLngLiteral | null;
  handleAction: (a: Site) => void;
  handleDeleteSite: (a: string) => void;
  setZoom: Dispatch<SetStateAction<number>>;
  handleAddMarker: (l: LatLngLiteral, b: string) => void;
  handleDragMarker: (l: LatLngLiteral, id: string) => void;
}

interface IMarker {
  id: string;
  lat: number;
  lng: number;
}

const CustomMarker = ({
  data,
  zoom,
  center,
  setZoom,
  handleAction,
  handleAddMarker,
  handleDeleteSite,
  handleDragMarker,
}: ICustomMarker) => {
  const map = useMap();
  const [markers, setMarkers] = useState<IMarker[]>([]);

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
