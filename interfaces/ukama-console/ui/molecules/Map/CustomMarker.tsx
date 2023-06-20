import { Site } from '@/generated/planning-tool';
import Leaflet, { LatLngLiteral } from 'leaflet';
import { Dispatch, SetStateAction, useEffect } from 'react';
import { Marker, Popup, useMap, useMapEvents } from 'react-leaflet';
import SitePopup from '../SitePopup';

const DEFAULT_CENTER = { lat: 37.7780627, lng: -121.9822475 };

interface ICustomMarker {
  data: Site[];
  marker: LatLngLiteral;
  zoom: number | undefined;
  handleAction: (a: Site) => void;
  setData: Dispatch<SetStateAction<any>>;
  setZoom: Dispatch<SetStateAction<number>>;
  handleAddMarker: (l: LatLngLiteral) => void;
  handleDragMarker: (l: LatLngLiteral, id: string) => void;
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
  // useEffect(() => {
  //   map.eachLayer((layer) => {
  //     if (layer instanceof Leaflet.Marker) {
  //       layer.remove();
  //     }
  //   });
  // }, [data]);

  // useEffect(() => {
  //   if (m.length > 0){
  //     m.forEach((marker) => {
  //       Leaflet.marker(marker)
  //       .addTo(map)
  //       .bindPopup(
  //         ReactDOMServer.renderToStaticMarkup(
  //           <SitePopup
  //             latlng={}
  //             data={data}
  //             setData={setData}
  //             handleAction={(a, b) => console.log(a, b)}
  //           />,
  //         ),
  //       );
  //     });
  //   }

  //   // .openPopup();
  // }, [m]);

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
    // popupopen: (e) => {
    //   const ev: any = e.popup;
    //   const { lat, lng } = ev._latlng;
    //   console.log('Site: ', data);
    //   console.log('Popup open', lat, lng);
    // },
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
