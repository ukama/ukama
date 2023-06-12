import dynamic from 'next/dynamic';

const DynamicMap = dynamic(() => import('./DynamicMap'), {
  ssr: false,
});

const DEFAULT_WIDTH = 600;
const DEFAULT_HEIGHT = 600;

interface IMap {
  width?: number;
  height?: number;
  className?: string;
  center?: number[];
  children: any;
  zoom?: number;
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
      <DynamicMap className={className} zoom={zoom} center={center}>
        {children}
      </DynamicMap>
    </div>
  );
};

export default Map;
