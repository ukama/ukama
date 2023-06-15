import { commonData } from '@/app-recoil';
import {
  useAddDraftMutation,
  useGetDraftsQuery,
  useUpdateDraftNameMutation,
  useUpdateEventMutation,
  useUpdateSiteMutation,
} from '@/generated/planning-tool';
import styles from '@/styles/Site_Planning.module.css';
import { PageContainer } from '@/styles/global';
import { colors } from '@/styles/theme';
import { TCommonData, TSite } from '@/types';
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
import { useRecoilValue } from 'recoil';

const DEFAULT_CENTER = [38.907132, -77.036546];
const DRAFTS = [{ id: 1, name: 'Draft 1' }];

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
  const _commonData = useRecoilValue<TCommonData>(commonData);
  const [currentDraft, setCurrentDraft] = useState({
    id: '',
    name: '',
  });
  const [search, setSearch] = useState('');
  const [addSite, setAddSite] = useState(false);
  const [addLink, setAddLink] = useState(false);
  const [togglePower, setTogglePower] = useState(false);
  const [marker, setMarker] = useState([0, 0]);
  const [anchorSiteInfo, setAnchorSiteInfo] =
    useState<HTMLButtonElement | null>(null);

  const {
    data: getDraftsData,
    loading: getDraftsLoading,
    refetch: refetchDrafts,
  } = useGetDraftsQuery({
    variables: {
      userId: _commonData.userId,
    },
    onCompleted: (data) => {
      /* Save drafts in state */
      console.log(data);
      setCurrentDraft({
        id: data.getDrafts[0].id,
        name: data.getDrafts[0].name,
      });
    },
    onError: (error) => {
      /* Show error message */
    },
  });

  const [addDraftCall, { loading: addDraftLoading }] = useAddDraftMutation({
    onCompleted: (data) => {
      /* Show success message */
      refetchDrafts();
    },
    onError: (error) => {
      /* Show error message */
    },
  });
  const [updateDraftCall, { loading: updateDraftLoading }] =
    useUpdateDraftNameMutation({
      onCompleted: (data) => {
        /* Show success message */
      },
      onError: (error) => {
        /* Show error message */
      },
    });
  const [updateSiteCall, { loading: updateSiteLoading }] =
    useUpdateSiteMutation({
      onCompleted: (data) => {
        /* Show success message */
      },
      onError: (error) => {
        /* Show error message */
      },
    });
  const [updateEventCall, { loading: updateEventLoading }] =
    useUpdateEventMutation({
      onCompleted: (data) => {
        /* Show success message */
      },
      onError: (error) => {
        /* Show error message */
      },
    });

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
  const handleAddDraft = () => {
    addDraftCall({
      variables: {
        data: {
          name: 'New Draft',
          userId: _commonData.userId,
          lastSaved: new Date().getTime() / 1000,
        },
      },
    });
  };
  const handleDraftSelected = (draftId: string) => {
    setCurrentDraft({
      id: draftId,
      name:
        getDraftsData?.getDrafts.find(({ id }) => id === draftId)?.name || '',
    });
  };
  const handleDraftUpdated = (id: string, draft: string) => {
    updateDraftCall({
      variables: {
        name: draft,
        updateDraftNameId: id,
      },
    });
  };

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
        <DraftDropdown
          drafts={getDraftsData?.getDrafts || []}
          currentDraft={currentDraft}
          isLoading={getDraftsLoading}
          handleAddDraft={handleAddDraft}
          handleDraftUpdated={handleDraftUpdated}
          handleDraftSelected={handleDraftSelected}
        />
        <PageContainer sx={{ padding: 0, mt: '12px' }}>
          <Map
            id={'site-planning-map'}
            zoom={12}
            width={800}
            height={418}
            isAddSite={addSite}
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
