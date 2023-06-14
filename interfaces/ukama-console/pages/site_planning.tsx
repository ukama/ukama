import styles from '@/styles/Site_Planning.module.css';
import { PageContainer } from '@/styles/global';
import { colors } from '@/styles/theme';
import { TSite } from '@/types';
import SitePopup from '@/ui/SitePopup';
import DraftDropdown from '@/ui/molecules/DraftDropdown';
import LoadingWrapper from '@/ui/molecules/LoadingWrapper';
import Map from '@/ui/molecules/Map';
import {
  LeftOverlayUI,
  RightOverlayUI,
  SiteSummary,
} from '@/ui/molecules/MapOverlayUI';
import { Popover } from '@mui/material';
import { useState } from 'react';

const DEFAULT_CENTER = [38.907132, -77.036546];

const Page = () => {
  const [site, setSite] = useState<TSite>({
    name: '',
    height: 0,
    solarUptime: 95,
    ap: 'ONE_TO_ONE',
    isBackhaul: true,
    location: {
      lat: 0,
      lng: 0,
      address: '',
    },
  });
  const [search, setSearch] = useState('');
  const [addSite, setAddSite] = useState(false);
  const [addLink, setAddLink] = useState(false);
  const [togglePower, setTogglePower] = useState(false);
  const [marker, setMarker] = useState([0, 0]);
  const [anchorSiteInfo, setAnchorSiteInfo] =
    useState<HTMLButtonElement | null>(null);

  const handleClick = (event: React.MouseEvent<HTMLButtonElement>) => {
    setAnchorSiteInfo(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorSiteInfo(null);
  };

  const open = Boolean(anchorSiteInfo);
  const id = open ? 'site-info-popover' : undefined;

  const handleMarkerDrag = (e: any) => {
    setMarker(e.latlng);
  };

  const handleMarkerAdd = (e: number[]) => {
    if (
      addSite &&
      (marker.length === 0 || (marker[0] === 0 && marker[1] === 0))
    ) {
      setAddSite(false);
      setMarker(e);
    }
  };

  const handleSiteAction = (action: string) => {
    if (action === 'add') {
    } else if (action === 'update') {
    }
  };

  const handleAddSite = () => setAddSite(true);
  const handleAddLink = () => setAddLink(true);
  const handleOnOff = () => setTogglePower(!togglePower);

  return (
    <>
      <Popover
        id={id}
        open={open}
        anchorOrigin={{
          vertical: 'bottom',
          horizontal: 'left',
        }}
        onClose={handleClose}
        anchorEl={anchorSiteInfo}
        sx={{ top: 4, left: -40 }}
        PaperProps={{
          sx: {
            width: '204px',
            padding: '16px 24px',
          },
        }}
      >
        <SiteSummary />
      </Popover>
      <LoadingWrapper
        radius="small"
        width={'100%'}
        isLoading={false}
        cstyle={{
          backgroundColor: false ? colors.white : 'transparent',
        }}
      >
        <DraftDropdown />
        <PageContainer sx={{ padding: 0, mt: '12px' }}>
          <Map
            id={'site-planning-map'}
            zoom={12}
            width={800}
            height={418}
            center={DEFAULT_CENTER}
            className={styles.homeMap}
            onMapClick={handleMarkerAdd}
          >
            {({ TileLayer, Marker, Popup }: any) => (
              <>
                <LeftOverlayUI
                  search={search}
                  setSearch={setSearch}
                  handleAddSite={handleAddSite}
                  handleAddLink={handleAddLink}
                />
                <RightOverlayUI
                  id={id}
                  handleClick={handleClick}
                  handleTogglePower={handleOnOff}
                />
                <TileLayer url="https://tiles.stadiamaps.com/tiles/alidade_smooth/{z}/{x}/{y}{r}.png" />
                <Marker draggable position={marker} ondrag={handleMarkerDrag}>
                  <Popup>
                    <SitePopup
                      data={site}
                      setData={setSite}
                      handleAction={handleSiteAction}
                    />
                  </Popup>
                </Marker>
              </>
            )}
          </Map>
        </PageContainer>
      </LoadingWrapper>
    </>
  );
};

export default Page;
