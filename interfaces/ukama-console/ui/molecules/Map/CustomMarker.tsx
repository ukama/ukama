import { Site } from '@/generated/planning-tool';
import Leaflet, { LatLngLiteral } from 'leaflet';
import { Dispatch, SetStateAction, useEffect } from 'react';
import { Marker, Popup, useMap, useMapEvents } from 'react-leaflet';
import SitePopup from '../SitePopup';

const DEFAULT_CENTER = { lat: 37.7780627, lng: -121.9822475 };

interface ICustomMarker {
  data: Site;
  marker: LatLngLiteral;
  handleAction: () => void;
  zoom: number | undefined;
  setData: Dispatch<SetStateAction<any>>;
  setZoom: Dispatch<SetStateAction<number>>;
  handleDragMarker: (l: LatLngLiteral) => void;
  handleAddMarker: (l: LatLngLiteral) => void;
}

const CustomMarker = ({
  zoom,
  data,
  marker,
  setData,
  setZoom,
  handleAction,
  handleAddMarker,
  handleDragMarker,
}: ICustomMarker) => {
  const map = useMap();
  useEffect(() => {
    map.setView(marker.lat === 0 ? DEFAULT_CENTER : marker, zoom);
  }, [marker]);

  useMapEvents({
    click: (e) => {
      const { lat, lng } = e.latlng;
      handleAddMarker({ lat, lng });
      Leaflet.tooltip().openTooltip();
    },
    zoom: (e) => {
      setZoom(e.target.getZoom());
    },
  });
  return (
    <div>
      <Marker
        autoPan
        draggable
        position={marker}
        opacity={marker.lat === 0 ? 0 : 1}
        eventHandlers={{
          moveend: (event: any) => handleDragMarker(event.target.getLatLng()),
        }}
      >
        <Popup>
          <SitePopup
            data={data}
            setData={setData}
            handleAction={handleAction}
          />
        </Popup>
      </Marker>
    </div>
  );
};

export default CustomMarker;
