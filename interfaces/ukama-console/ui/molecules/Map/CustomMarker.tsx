import { Site } from '@/generated/planning-tool';
import Leaflet, { LatLngLiteral } from 'leaflet';
import { Dispatch, SetStateAction } from 'react';
import { Marker, Popup, useMapEvents } from 'react-leaflet';
import SitePopup from '../SitePopup';

interface ICustomMarker {
  data: Site[];
  handleAction: (a: Site) => void;
  setZoom: Dispatch<SetStateAction<number>>;
  handleAddMarker: (l: LatLngLiteral) => void;
  handleDragMarker: (l: LatLngLiteral, id: string) => void;
}

const CustomMarker = ({
  data,
  setZoom,
  handleAction,
  handleAddMarker,
  handleDragMarker,
}: ICustomMarker) => {
  useMapEvents({
    click: (e) => {
      const { lat, lng } = e.latlng;
      handleAddMarker({ lat, lng });
      Leaflet.tooltip().openTooltip();
    },
    zoom: (e) => {
      console.log(e.target.getZoom());
      setZoom(e.target.getZoom());
    },
  });

  return (
    <div>
      {data.length > 0 &&
        data.map((item) => (
          <Marker
            key={item.id}
            title={item.id}
            autoPan
            draggable
            position={{
              lat: parseFloat(item.location.lat),
              lng: parseFloat(item.location.lng),
            }}
            opacity={parseFloat(item.location.lat) === 0 ? 0 : 1}
            eventHandlers={{
              moveend: (event: any) =>
                handleDragMarker(event.target.getLatLng(), item.location.id),
            }}
          >
            <Popup>
              <SitePopup site={item} handleAction={handleAction} />
            </Popup>
          </Marker>
        ))}
    </div>
  );
};

export default CustomMarker;
