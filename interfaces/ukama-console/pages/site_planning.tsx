import styles from '@/styles/Site_Planning.module.css';
import { PageContainer } from '@/styles/global';
import { colors } from '@/styles/theme';
import LoadingWrapper from '@/ui/molecules/LoadingWrapper';
import Map from '@/ui/molecules/Map';

const DEFAULT_CENTER = [38.907132, -77.036546];

export default function Page() {
  return (
    <LoadingWrapper
      radius="small"
      width={'100%'}
      isLoading={false}
      cstyle={{
        backgroundColor: false ? colors.white : 'transparent',
      }}
    >
      <PageContainer sx={{ padding: 0 }}>
        <Map
          zoom={12}
          width={800}
          height={418}
          center={DEFAULT_CENTER}
          className={styles.homeMap}
        >
          {({ TileLayer, Marker, Popup }: any) => (
            <>
              <TileLayer
                attribution='&copy; <a href="https://stadiamaps.com/">Stadia Maps</a>, &copy; <a href="https://openmaptiles.org/">OpenMapTiles</a> &copy; <a href="http://openstreetmap.org">OpenStreetMap</a> contributors'
                url="https://tiles.stadiamaps.com/tiles/alidade_smooth/{z}/{x}/{y}{r}.png"
              />
              <Marker position={DEFAULT_CENTER}>
                <Popup>Site Info</Popup>
              </Marker>
            </>
          )}
        </Map>
      </PageContainer>
    </LoadingWrapper>
  );
}
