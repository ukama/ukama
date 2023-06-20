import { Site } from '@/generated/planning-tool';
import { LatLngLiteral } from 'leaflet';
import dynamic from 'next/dynamic';
import { Dispatch, SetStateAction } from 'react';

const DynamicMap = dynamic(() => import('./DynamicMap'), {
  ssr: false,
});

const DEFAULT_WIDTH = 600;
const DEFAULT_HEIGHT = 600;

interface IMap {
  id: string;
  data: Site[];
  children: any;
  width?: number;
  height?: number;
  isAddSite: boolean;
  className?: string;
  center: LatLngLiteral;
  zoom?: number | undefined;
  handleAction: (a: Site) => void;
  setZoom: Dispatch<SetStateAction<number>>;
  handleAddMarker: (l: LatLngLiteral, b: string) => void;
  handleDragMarker: (l: LatLngLiteral, id: string) => void;
}

const Map = ({
  id,
  zoom,
  data,
  center,
  setZoom,
  children,
  isAddSite,
  className,
  handleAction,
  handleAddMarker,
  handleDragMarker,
  width = DEFAULT_WIDTH,
  height = DEFAULT_HEIGHT,
}: IMap) => {
  return (
    <div
      style={{
        aspectRatio: width / height,
        cursor: isAddSite ? 'pointer !important' : 'grab !important',
      }}
    >
      <DynamicMap
        id={id}
        zoom={zoom}
        data={data}
        center={center}
        setZoom={setZoom}
        cursor={isAddSite}
        className={className}
        handleAction={handleAction}
        handleAddMarker={handleAddMarker}
        handleDragMarker={handleDragMarker}
      >
        {children}
      </DynamicMap>
    </div>
  );
};

export default Map;
