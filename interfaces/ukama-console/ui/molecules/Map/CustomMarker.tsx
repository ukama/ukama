import { Site } from '@/generated/planning-tool';
import { randomUUID } from 'crypto';
import Leaflet, { LatLngLiteral } from 'leaflet';
import { Dispatch, SetStateAction, useEffect, useState } from 'react';
import { Marker, Popup, useMapEvents } from 'react-leaflet';
import SitePopup from '../SitePopup';

interface ICustomMarker {
  data: Site[];
  handleAction: (a: Site) => void;
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
  setZoom,
  handleAction,
  handleAddMarker,
  handleDragMarker,
}: ICustomMarker) => {
  const [markers, setMarkers] = useState<IMarker[]>([]);

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
  }, []);

  useMapEvents({
    click: (e) => {
      const { lat, lng } = e.latlng;
      const id = randomUUID();
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
              opacity={parseFloat(item.location.lat) === 0 ? 0 : 1}
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
                <SitePopup site={item} handleAction={handleAction} />
              </Popup>
            </Marker>
          );
        })}
    </div>
  );
};

export default CustomMarker;
