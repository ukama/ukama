import dynamic from 'next/dynamic';

const DynamicMap = dynamic(() => import('./DynamicMap'), {
  ssr: false,
});

const DEFAULT_WIDTH = 600;
const DEFAULT_HEIGHT = 600;

interface IMap {
  children: any;
  zoom?: number;
  width?: number;
  height?: number;
  center?: number[];
  className?: string;
}

const Map = ({
  zoom,
  center,
  children,
  className,
  width = DEFAULT_WIDTH,
  height = DEFAULT_HEIGHT,
}: IMap) => {
  return (
    <div style={{ aspectRatio: width / height }}>
      <DynamicMap zoom={zoom} center={center} className={className}>
        {children}
      </DynamicMap>
    </div>
  );
};

export default Map;
