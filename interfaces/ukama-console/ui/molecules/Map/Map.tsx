import dynamic from 'next/dynamic';

const DynamicMap = dynamic(() => import('./DynamicMap'), {
  ssr: false,
});

const DEFAULT_WIDTH = 600;
const DEFAULT_HEIGHT = 600;

interface IMap {
  id: string;
  children: any;
  zoom?: number;
  width?: number;
  height?: number;
  center?: number[];
  isAddSite: boolean;
  className?: string;
  onMapClick: Function;
}

const Map = ({
  id,
  zoom,
  center,
  children,
  isAddSite,
  className,
  onMapClick,
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
        center={center}
        cursor={isAddSite}
        className={className}
        onMapClick={onMapClick}
      >
        {children}
      </DynamicMap>
    </div>
  );
};

export default Map;
